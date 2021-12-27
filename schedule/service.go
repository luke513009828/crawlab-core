package schedule

import (
	"github.com/luke513009828/crawlab-core/config"
	"github.com/luke513009828/crawlab-core/interfaces"
	"github.com/luke513009828/crawlab-core/models/delegate"
	"github.com/luke513009828/crawlab-core/models/models"
	"github.com/luke513009828/crawlab-core/models/service"
	"github.com/luke513009828/crawlab-core/spider/admin"
	"github.com/luke513009828/crawlab-core/utils"
	"github.com/crawlab-team/go-trace"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/dig"
	"sync"
	"time"
)

type Service struct {
	// dependencies
	interfaces.WithConfigPath
	modelSvc service.ModelService
	adminSvc interfaces.SpiderAdminService

	// settings variables
	loc            *time.Location
	delay          bool
	skip           bool
	updateInterval time.Duration

	// internals
	cron      *cron.Cron
	logger    cron.Logger
	schedules []models.Schedule
	stopped   bool
	mu        sync.Mutex
}

func (svc *Service) GetLocation() (loc *time.Location) {
	return svc.loc
}

func (svc *Service) SetLocation(loc *time.Location) {
	svc.loc = loc
}

func (svc *Service) GetDelay() (delay bool) {
	return svc.delay
}

func (svc *Service) SetDelay(delay bool) {
	svc.delay = delay
}

func (svc *Service) GetSkip() (skip bool) {
	return svc.skip
}

func (svc *Service) SetSkip(skip bool) {
	svc.skip = skip
}

func (svc *Service) GetUpdateInterval() (interval time.Duration) {
	return svc.updateInterval
}

func (svc *Service) SetUpdateInterval(interval time.Duration) {
	svc.updateInterval = interval
}

func (svc *Service) Init() (err error) {
	return svc.fetch()
}

func (svc *Service) Start() {
	svc.cron.Start()
	go svc.Update()
}

func (svc *Service) Wait() {
	utils.DefaultWait()
	svc.Stop()
}

func (svc *Service) Stop() {
	svc.stopped = true
	svc.cron.Stop()
}

func (svc *Service) Enable(s interfaces.Schedule, args ...interface{}) (err error) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	id, err := svc.cron.AddFunc(s.GetCron(), svc.schedule(s.GetId()))
	if err != nil {
		return trace.TraceError(err)
	}
	s.SetEnabled(true)
	s.SetEntryId(id)
	u := utils.GetUserFromArgs(args...)
	return delegate.NewModelDelegate(s, u).Save()
}

func (svc *Service) Disable(s interfaces.Schedule, args ...interface{}) (err error) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	svc.cron.Remove(s.GetEntryId())
	s.SetEnabled(false)
	s.SetEntryId(-1)
	u := utils.GetUserFromArgs(args...)
	return delegate.NewModelDelegate(s, u).Save()
}

func (svc *Service) Update() {
	for {
		if svc.stopped {
			return
		}

		svc.update()

		time.Sleep(svc.updateInterval)
	}
}

func (svc *Service) GetCron() (c *cron.Cron) {
	return svc.cron
}

func (svc *Service) update() {
	// fetch enabled schedules
	if err := svc.fetch(); err != nil {
		trace.PrintError(err)
		return
	}

	// entry id map
	entryIdsMap := svc.getEntryIdsMap()

	// iterate enabled schedules
	for _, s := range svc.schedules {
		_, ok := entryIdsMap[s.EntryId]
		if ok {
			entryIdsMap[s.EntryId] = true
		} else {
			if err := svc.Enable(&s); err != nil {
				trace.PrintError(err)
				continue
			}
		}
	}

	// remove non-existent entries
	for id, ok := range entryIdsMap {
		if !ok {
			svc.cron.Remove(id)
		}
	}
}

func (svc *Service) getEntryIdsMap() (res map[cron.EntryID]bool) {
	res = map[cron.EntryID]bool{}
	for _, e := range svc.cron.Entries() {
		res[e.ID] = false
	}
	return res
}

func (svc *Service) fetch() (err error) {
	query := bson.M{
		"enabled": true,
	}
	svc.schedules, err = svc.modelSvc.GetScheduleList(query, nil)
	if err != nil {
		return err
	}
	return nil
}

func (svc *Service) schedule(id primitive.ObjectID) (fn func()) {
	return func() {
		// schedule
		s, err := svc.modelSvc.GetScheduleById(id)
		if err != nil {
			trace.PrintError(err)
			return
		}

		// spider
		spider, err := svc.modelSvc.GetSpiderById(s.GetSpiderId())
		if err != nil {
			trace.PrintError(err)
			return
		}

		// options
		opts := &interfaces.SpiderRunOptions{
			Mode:       s.GetMode(),
			NodeIds:    s.GetNodeIds(),
			NodeTags:   s.GetNodeTags(),
			Cmd:        s.GetCmd(),
			Param:      s.GetParam(),
			Priority:   s.GetPriority(),
			ScheduleId: s.GetId(),
			UserId:     s.UserId,
		}

		// normalize options
		if opts.Mode == "" {
			opts.Mode = spider.Mode
		}
		if len(opts.NodeIds) == 0 {
			opts.NodeIds = spider.NodeIds
		}
		if len(opts.NodeTags) == 0 {
			opts.NodeTags = spider.NodeTags
		}
		if opts.Cmd == "" {
			opts.Cmd = spider.Cmd
		}
		if opts.Param == "" {
			opts.Param = spider.Param
		}
		if opts.Priority == 0 {
			if spider.Priority > 0 {
				opts.Priority = spider.Priority
			} else {
				opts.Priority = 5
			}
		}

		// schedule
		if err := svc.adminSvc.Schedule(s.GetSpiderId(), opts); err != nil {
			trace.PrintError(err)
		}
	}
}

func NewScheduleService(opts ...Option) (svc2 interfaces.ScheduleService, err error) {
	// service
	svc := &Service{
		WithConfigPath: config.NewConfigPathService(),
		loc:            time.Local,
		// TODO: implement delay and skip
		delay:          false,
		skip:           false,
		updateInterval: 1 * time.Minute,
	}

	// apply options
	for _, opt := range opts {
		opt(svc)
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.GetService); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Provide(admin.ProvideSpiderAdminService(svc.GetConfigPath())); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Invoke(func(
		modelSvc service.ModelService,
		adminSvc interfaces.SpiderAdminService,
	) {
		svc.modelSvc = modelSvc
		svc.adminSvc = adminSvc
	}); err != nil {
		return nil, trace.TraceError(err)
	}

	// logger
	svc.logger = NewLogger()

	// cron
	svc.cron = cron.New(
		cron.WithLogger(svc.logger),
		cron.WithLocation(svc.loc),
		cron.WithChain(cron.Recover(svc.logger)),
	)

	// initialize
	if err := svc.Init(); err != nil {
		return nil, err
	}

	return svc, nil
}

func ProvideScheduleService(path string, opts ...Option) func() (svc interfaces.ScheduleService, err error) {
	opts = append(opts, WithConfigPath(path))
	return func() (svc interfaces.ScheduleService, err error) {
		return NewScheduleService(opts...)
	}
}

var store = sync.Map{}

func GetScheduleService(path string, opts ...Option) (svc interfaces.ScheduleService, err error) {
	if path == "" {
		path = config.DefaultConfigPath
	}
	opts = append(opts, WithConfigPath(path))
	res, ok := store.Load(path)
	if ok {
		svc, ok = res.(interfaces.ScheduleService)
		if ok {
			return svc, nil
		}
	}
	svc, err = NewScheduleService(opts...)
	if err != nil {
		return nil, err
	}
	store.Store(path, svc)
	return svc, nil
}

func ProvideGetScheduleService(path string, opts ...Option) func() (svr interfaces.ScheduleService, err error) {
	return func() (svr interfaces.ScheduleService, err error) {
		return GetScheduleService(path, opts...)
	}
}
