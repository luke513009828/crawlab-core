package service

import (
	"github.com/luke513009828/crawlab-core/errors"
	"github.com/luke513009828/crawlab-core/interfaces"
	models2 "github.com/luke513009828/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeTask(d interface{}, err error) (res *models2.Task, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*models2.Task)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetTaskById(id primitive.ObjectID) (res *models2.Task, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdTask).GetById(id)
	return convertTypeTask(d, err)
}

func (svc *Service) GetTask(query bson.M, opts *mongo.FindOptions) (res *models2.Task, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdTask).Get(query, opts)
	return convertTypeTask(d, err)
}

func (svc *Service) GetTaskList(query bson.M, opts *mongo.FindOptions) (res []models2.Task, err error) {
	err = svc.getListSerializeTarget(interfaces.ModelIdTask, query, opts, &res)
	return res, err
}
