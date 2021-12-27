package service

import (
	"github.com/luke513009828/crawlab-core/errors"
	"github.com/luke513009828/crawlab-core/interfaces"
	models2 "github.com/luke513009828/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeTaskQueueItem(d interface{}, err error) (res *models2.TaskQueueItem, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*models2.TaskQueueItem)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetTaskQueueItemById(id primitive.ObjectID) (res *models2.TaskQueueItem, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdTaskQueue).GetById(id)
	return convertTypeTaskQueueItem(d, err)
}

func (svc *Service) GetTaskQueueItem(query bson.M, opts *mongo.FindOptions) (res *models2.TaskQueueItem, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdTaskQueue).Get(query, opts)
	return convertTypeTaskQueueItem(d, err)
}

func (svc *Service) GetTaskQueueItemList(query bson.M, opts *mongo.FindOptions) (res []models2.TaskQueueItem, err error) {
	err = svc.getListSerializeTarget(interfaces.ModelIdTaskQueue, query, opts, &res)
	return res, err
}
