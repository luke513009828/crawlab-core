package client

import (
	"github.com/luke513009828/crawlab-core/errors"
	"github.com/luke513009828/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskServiceDelegate struct {
	interfaces.GrpcClientModelBaseService
}

func (svc *TaskServiceDelegate) GetTaskById(id primitive.ObjectID) (t interfaces.Task, err error) {
	res, err := svc.GetById(id)
	if err != nil {
		return nil, err
	}
	s, ok := res.(interfaces.Task)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return s, nil
}

func (svc *TaskServiceDelegate) GetTask(query bson.M, opts *mongo.FindOptions) (t interfaces.Task, err error) {
	res, err := svc.Get(query, opts)
	if err != nil {
		return nil, err
	}
	s, ok := res.(interfaces.Task)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return s, nil
}

func (svc *TaskServiceDelegate) GetTaskList(query bson.M, opts *mongo.FindOptions) (res []interfaces.Task, err error) {
	list, err := svc.GetList(query, opts)
	if err != nil {
		return nil, err
	}
	for _, item := range list.Values() {
		s, ok := item.(interfaces.Task)
		if !ok {
			return nil, errors.ErrorModelInvalidType
		}
		res = append(res, s)
	}
	return res, nil
}

func NewTaskServiceDelegate(opts ...ModelBaseServiceDelegateOption) (svc2 interfaces.GrpcClientModelTaskService, err error) {
	// apply options
	opts = append(opts, WithBaseServiceModelId(interfaces.ModelIdTask))

	// base service
	baseSvc, err := NewBaseServiceDelegate(opts...)
	if err != nil {
		return nil, err
	}

	// service
	svc := &TaskServiceDelegate{baseSvc}

	return svc, nil
}

func ProvideTaskServiceDelegate(path string, opts ...ModelBaseServiceDelegateOption) func() (svc interfaces.GrpcClientModelTaskService, err error) {
	if path != "" {
		opts = append(opts, WithBaseServiceConfigPath(path))
	}
	return func() (svc interfaces.GrpcClientModelTaskService, err error) {
		return NewTaskServiceDelegate(opts...)
	}
}
