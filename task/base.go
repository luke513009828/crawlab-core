package task

import (
	"fmt"
	"github.com/luke513009828/crawlab-core/config"
	"github.com/luke513009828/crawlab-core/constants"
	"github.com/luke513009828/crawlab-core/interfaces"
	"github.com/luke513009828/crawlab-core/models/delegate"
	"github.com/luke513009828/crawlab-core/models/service"
	"github.com/luke513009828/crawlab-core/utils"
	"github.com/crawlab-team/go-trace"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/dig"
)

type BaseService struct {
	// dependencies
	interfaces.WithConfigPath
	modelSvc service.ModelService

	// internals
	stopped bool
}

func (svc *BaseService) Init() error {
	// implement me
	return nil
}

func (svc *BaseService) Start() {
	// implement me
}

func (svc *BaseService) Wait() {
	utils.DefaultWait()
}

func (svc *BaseService) Stop() {
	svc.stopped = true
}

// SaveTask deprecated
func (svc *BaseService) SaveTask(t interfaces.Task, status string) (err error) {
	// normalize status
	if status == "" {
		status = constants.TaskStatusPending
	}

	// set task status
	t.SetStatus(status)

	// attempt to get task from database
	_, err = svc.modelSvc.GetTaskById(t.GetId())
	if err != nil {
		// if task does not exist, add to database
		if err == mongo.ErrNoDocuments {
			if err := delegate.NewModelDelegate(t).Add(); err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	} else {
		// otherwise, update
		if err := delegate.NewModelDelegate(t).Save(); err != nil {
			return err
		}
		return nil
	}
}

func (svc *BaseService) IsStopped() (res bool) {
	return svc.stopped
}

func (svc *BaseService) GetQueue(nodeId primitive.ObjectID) (queue string) {
	if nodeId.IsZero() {
		return fmt.Sprintf("%s", constants.TaskListQueuePrefixPublic)
	} else {
		return fmt.Sprintf("%s:%s", constants.TaskListQueuePrefixNodes, nodeId.Hex())
	}
}

func NewBaseService() (svc2 interfaces.TaskBaseService, err error) {
	svc := &BaseService{}

	// dependency injection
	c := dig.New()
	if err := c.Provide(config.NewConfigPathService); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Provide(service.NewService); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Invoke(func(cfgPath interfaces.WithConfigPath, modelSvc service.ModelService) {
		svc.WithConfigPath = cfgPath
		svc.modelSvc = modelSvc
	}); err != nil {
		return nil, trace.TraceError(err)
	}

	return svc, nil
}
