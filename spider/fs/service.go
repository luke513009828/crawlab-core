package fs

import (
	"fmt"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/fs"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/go-trace"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/dig"
	"os"
	"sync"
)

// Service implementation of interfaces.SpiderFsService
// It is a wrapper of interfaces.FsService that manages a spider's fs related functions
type Service struct {
	// settings
	cfgPath           string
	fsPathBase        string
	workspacePathBase string
	repoPathBase      string

	// dependencies
	modelSvc service.ModelService
	fsSvc    interfaces.FsService

	// internals
	id primitive.ObjectID
	s  *models.Spider
}

func (svc *Service) Init() (err error) {
	// workspace
	if _, err := os.Stat(svc.GetWorkspacePath()); err != nil {
		if err := os.MkdirAll(svc.GetWorkspacePath(), os.FileMode(0766)); err != nil {
			return trace.TraceError(err)
		}
	}

	return nil
}

func (svc *Service) GetConfigPath() string {
	return svc.cfgPath
}

func (svc *Service) SetConfigPath(path string) {
	svc.cfgPath = path
}

func (svc *Service) SetId(id primitive.ObjectID) {
	svc.id = id
}

func (svc *Service) GetFsPath() (res string) {
	return fmt.Sprintf("%s/%s", svc.fsPathBase, svc.id.Hex())
}

func (svc *Service) GetWorkspacePath() (res string) {
	return fmt.Sprintf("%s/%s", svc.workspacePathBase, svc.id.Hex())
}

func (svc *Service) GetRepoPath() (res string) {
	return fmt.Sprintf("%s/%s", svc.repoPathBase, svc.id.Hex())
}

func (svc *Service) SetFsPathBase(path string) {
	svc.fsPathBase = path
}

func (svc *Service) SetWorkspacePathBase(path string) {
	svc.workspacePathBase = path
}

func (svc *Service) SetRepoPathBase(path string) {
	svc.repoPathBase = path
}

func (svc *Service) GetFsService() (fsSvc interfaces.FsService) {
	return svc.fsSvc
}

func NewSpiderFsService(id primitive.ObjectID, opts ...Option) (svc2 interfaces.SpiderFsService, err error) {
	// service
	svc := &Service{
		fsPathBase:        fs.DefaultFsPath,
		workspacePathBase: fs.DefaultWorkspacePath,
		repoPathBase:      fs.DefaultRepoPath,
		id:                id,
	}

	// apply options
	for _, opt := range opts {
		opt(svc)
	}

	// validate
	if svc.id.IsZero() {
		return nil, trace.TraceError(errors.ErrorSpiderMissingRequiredOption)
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.NewService); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Invoke(func(modelSvc service.ModelService) {
		svc.modelSvc = modelSvc
	}); err != nil {
		return nil, trace.TraceError(err)
	}

	// spider
	svc.s, err = svc.modelSvc.GetSpiderById(svc.id)
	if err != nil {
		return nil, err
	}

	// fs service
	var fsOpts []fs.Option
	fsOpts = append(fsOpts, fs.WithConfigPath(svc.cfgPath))
	fsOpts = append(fsOpts, fs.WithFsPath(svc.GetFsPath()))
	fsOpts = append(fsOpts, fs.WithWorkspacePath(svc.GetWorkspacePath()))
	if svc.repoPathBase != "" {
		fsOpts = append(fsOpts, fs.WithRepoPath(svc.GetRepoPath()))
	}
	svc.fsSvc, err = fs.NewFsService(fsOpts...)
	if err != nil {
		return nil, err
	}

	// initialize
	if err := svc.Init(); err != nil {
		return nil, err
	}

	return svc, nil
}

func ProvideSpiderFsService(id primitive.ObjectID, opts ...Option) func() (svc interfaces.SpiderFsService, err error) {
	return func() (svc interfaces.SpiderFsService, err error) {
		return NewSpiderFsService(id, opts...)
	}
}

var spiderFsSvcCache = sync.Map{}

func GetSpiderFsService(id primitive.ObjectID, opts ...Option) (svc interfaces.SpiderFsService, err error) {
	// cache key consisted of id and config path
	// FIXME: this is only for testing purpose
	cfgPath := getConfigPathFromOptions(opts...)
	key := getHashStringFromIdAndConfigPath(id, cfgPath)

	// attempt to load from cache
	res, ok := spiderFsSvcCache.Load(key)
	if !ok {
		// not exists in cache, create a new service
		svc, err = NewSpiderFsService(id, opts...)
		if err != nil {
			return nil, err
		}

		// store in cache
		spiderFsSvcCache.Store(key, svc)
		return svc, nil
	}

	// load from cache successful
	svc, ok = res.(interfaces.SpiderFsService)
	if !ok {
		return nil, trace.TraceError(errors.ErrorFsInvalidType)
	}

	return svc, nil
}

func ProvideGetSpiderFsService(id primitive.ObjectID, opts ...Option) func() (svc interfaces.SpiderFsService, err error) {
	return func() (svc interfaces.SpiderFsService, err error) {
		return GetSpiderFsService(id, opts...)
	}
}
