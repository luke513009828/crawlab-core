package service

import (
	"github.com/apex/log"
	config2 "github.com/luke513009828/crawlab-core/config"
	"github.com/luke513009828/crawlab-core/constants"
	"github.com/luke513009828/crawlab-core/errors"
	"github.com/luke513009828/crawlab-core/grpc/server"
	"github.com/luke513009828/crawlab-core/interfaces"
	"github.com/luke513009828/crawlab-core/models/common"
	"github.com/luke513009828/crawlab-core/models/delegate"
	"github.com/luke513009828/crawlab-core/models/models"
	"github.com/luke513009828/crawlab-core/models/service"
	"github.com/luke513009828/crawlab-core/node/config"
	"github.com/luke513009828/crawlab-core/plugin"
	"github.com/luke513009828/crawlab-core/schedule"
	"github.com/luke513009828/crawlab-core/task/handler"
	"github.com/luke513009828/crawlab-core/task/scheduler"
	"github.com/luke513009828/crawlab-core/utils"
	grpc "github.com/crawlab-team/crawlab-grpc"
	"github.com/crawlab-team/go-trace"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/dig"
	"time"
)

type MasterService struct {
	// dependencies
	modelSvc     service.ModelService
	cfgSvc       interfaces.NodeConfigService
	server       interfaces.GrpcServer
	schedulerSvc interfaces.TaskSchedulerService
	handlerSvc   interfaces.TaskHandlerService
	scheduleSvc  interfaces.ScheduleService
	pluginSvc    interfaces.PluginService

	// settings
	cfgPath         string
	address         interfaces.Address
	monitorInterval time.Duration
	stopOnError     bool
}

func (svc *MasterService) Init() (err error) {
	// do nothing
	return nil
}

func (svc *MasterService) Start() {
	// create indexes
	common.CreateIndexes()

	// start grpc server
	if err := svc.server.Start(); err != nil {
		panic(err)
	}

	// register to db
	if err := svc.Register(); err != nil {
		panic(err)
	}

	// start monitoring worker nodes
	go svc.Monitor()

	// start task handler
	go svc.handlerSvc.Start()

	// start task scheduler
	go svc.schedulerSvc.Start()

	// start schedule service
	go svc.scheduleSvc.Start()

	// start plugin service
	go svc.pluginSvc.Start()

	// wait for quit signal
	svc.Wait()

	// stop
	svc.Stop()
}

func (svc *MasterService) Wait() {
	utils.DefaultWait()
}

func (svc *MasterService) Stop() {
	_ = svc.server.Stop()
	log.Infof("master[%s] service has stopped", svc.GetConfigService().GetNodeKey())
}

func (svc *MasterService) Monitor() {
	log.Infof("master[%s] monitoring started", svc.GetConfigService().GetNodeKey())
	for {
		if err := svc.monitor(); err != nil {
			trace.PrintError(err)
			if svc.stopOnError {
				log.Errorf("master[%s] monitor error, now stopping...", svc.GetConfigService().GetNodeKey())
				svc.Stop()
				return
			}
		}

		time.Sleep(svc.monitorInterval)
	}
}

func (svc *MasterService) GetConfigService() (cfgSvc interfaces.NodeConfigService) {
	return svc.cfgSvc
}

func (svc *MasterService) GetConfigPath() (path string) {
	return svc.cfgPath
}

func (svc *MasterService) SetConfigPath(path string) {
	svc.cfgPath = path
}

func (svc *MasterService) GetAddress() (address interfaces.Address) {
	return svc.address
}

func (svc *MasterService) SetAddress(address interfaces.Address) {
	svc.address = address
}

func (svc *MasterService) SetMonitorInterval(duration time.Duration) {
	svc.monitorInterval = duration
}

func (svc *MasterService) Register() (err error) {
	nodeKey := svc.GetConfigService().GetNodeKey()
	nodeName := svc.GetConfigService().GetNodeName()
	node, err := svc.modelSvc.GetNodeByKey(nodeKey, nil)
	if err != nil && err.Error() == mongo2.ErrNoDocuments.Error() {
		// not exists
		log.Infof("master[%s] does not exist in db", nodeKey)
		node := &models.Node{
			Key:        nodeKey,
			Name:       nodeName,
			MaxRunners: config.DefaultConfigOptions.MaxRunners,
			IsMaster:   true,
			Status:     constants.NodeStatusOnline,
			Enabled:    true,
			Active:     true,
			ActiveTs:   time.Now(),
		}
		if viper.GetInt("task.handler.maxRunners") > 0 {
			node.MaxRunners = viper.GetInt("task.handler.maxRunners")
		}
		nodeD := delegate.NewModelNodeDelegate(node)
		if err := nodeD.Add(); err != nil {
			return err
		}
		log.Infof("added master[%s] in db. id: %s", nodeKey, nodeD.GetModel().GetId().Hex())
		return nil
	} else if err == nil {
		// exists
		log.Infof("master[%s] exists in db", nodeKey)
		nodeD := delegate.NewModelNodeDelegate(node)
		if err := nodeD.UpdateStatusOnline(); err != nil {
			return err
		}
		log.Infof("updated master[%s] in db. id: %s", nodeKey, nodeD.GetModel().GetId().Hex())
		return nil
	} else {
		// error
		return err
	}
}

func (svc *MasterService) StopOnError() {
	svc.stopOnError = true
}

func (svc *MasterService) GetServer() (svr interfaces.GrpcServer) {
	return svc.server
}

func (svc *MasterService) monitor() (err error) {
	// update master node status in db
	if err := svc.updateMasterNodeStatus(); err != nil {
		if err.Error() == mongo2.ErrNoDocuments.Error() {
			return nil
		}
		return err
	}

	// all worker nodes
	query := bson.M{
		"key":    bson.M{"$ne": svc.cfgSvc.GetNodeKey()}, // not self
		"active": true,                                   // active
	}
	nodes, err := svc.modelSvc.GetNodeList(query, nil)
	if err != nil {
		if err == mongo2.ErrNoDocuments {
			return nil
		}
		return trace.TraceError(err)
	}

	// error flag
	isErr := false

	// iterate all nodes
	for _, n := range nodes {
		// subscribe
		_, err := svc.server.GetSubscribe("node:" + n.GetKey())
		if err != nil {
			trace.PrintError(err)
			isErr = true
			if err := svc.setWorkerNodeOffline(&n); err != nil {
				trace.PrintError(err)
			}
			continue
		}

		// PING client
		if err := svc.server.SendStreamMessage("node:"+n.GetKey(), grpc.StreamMessageCode_PING); err != nil {
			log.Errorf("cannot ping worker[%s]: %v", n.GetKey(), err)
			trace.PrintError(err)
			isErr = true
			if err := svc.setWorkerNodeOffline(&n); err != nil {
				trace.PrintError(err)
			}
			continue
		}
	}

	if isErr {
		return trace.TraceError(errors.ErrorNodeMonitorError)
	}

	return nil
}

func (svc *MasterService) updateMasterNodeStatus() (err error) {
	nodeKey := svc.GetConfigService().GetNodeKey()
	node, err := svc.modelSvc.GetNodeByKey(nodeKey, nil)
	if err != nil {
		return err
	}
	nodeD := delegate.NewModelNodeDelegate(node)
	return nodeD.UpdateStatusOnline()
}

func (svc *MasterService) setWorkerNodeOffline(n interfaces.Node) (err error) {
	return delegate.NewModelNodeDelegate(n).UpdateStatusOffline()
}

func NewMasterService(opts ...Option) (res interfaces.NodeMasterService, err error) {
	// master service
	svc := &MasterService{
		cfgPath:         config2.DefaultConfigPath,
		monitorInterval: 15 * time.Second,
		stopOnError:     false,
	}

	// apply options
	for _, opt := range opts {
		opt(svc)
	}

	// server options
	var serverOpts []server.Option
	if svc.address != nil {
		serverOpts = append(serverOpts, server.WithAddress(svc.address))
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.NewService); err != nil {
		return nil, err
	}
	if err := c.Provide(config.ProvideConfigService(svc.cfgPath)); err != nil {
		return nil, err
	}
	if err := c.Provide(server.ProvideGetServer(svc.cfgPath, serverOpts...)); err != nil {
		return nil, err
	}
	if err := c.Provide(scheduler.ProvideGetTaskSchedulerService(svc.cfgPath)); err != nil {
		return nil, err
	}
	if err := c.Provide(handler.ProvideGetTaskHandlerService(svc.cfgPath)); err != nil {
		return nil, err
	}
	if err := c.Provide(schedule.ProvideGetScheduleService(svc.cfgPath)); err != nil {
		return nil, err
	}
	if err := c.Provide(plugin.ProvideGetPluginService(svc.cfgPath)); err != nil {
		return nil, err
	}
	if err := c.Invoke(func(
		cfgSvc interfaces.NodeConfigService,
		modelSvc service.ModelService,
		server interfaces.GrpcServer,
		schedulerSvc interfaces.TaskSchedulerService,
		handlerSvc interfaces.TaskHandlerService,
		scheduleSvc interfaces.ScheduleService,
		pluginSvc interfaces.PluginService,
	) {
		svc.cfgSvc = cfgSvc
		svc.modelSvc = modelSvc
		svc.server = server
		svc.schedulerSvc = schedulerSvc
		svc.handlerSvc = handlerSvc
		svc.scheduleSvc = scheduleSvc
		svc.pluginSvc = pluginSvc
	}); err != nil {
		return nil, err
	}

	// init
	if err := svc.Init(); err != nil {
		return nil, err
	}

	return svc, nil
}

func ProvideMasterService(path string, opts ...Option) func() (interfaces.NodeMasterService, error) {
	if path != "" {
		opts = append(opts, WithConfigPath(path))
	}
	return func() (interfaces.NodeMasterService, error) {
		return NewMasterService(opts...)
	}
}
