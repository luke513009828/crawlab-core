package service

import (
	"github.com/luke513009828/crawlab-core/errors"
	"github.com/luke513009828/crawlab-core/interfaces"
	models2 "github.com/luke513009828/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeDataCollection(d interface{}, err error) (res *models2.DataCollection, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*models2.DataCollection)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetDataCollectionById(id primitive.ObjectID) (res *models2.DataCollection, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdDataCollection).GetById(id)
	return convertTypeDataCollection(d, err)
}

func (svc *Service) GetDataCollection(query bson.M, opts *mongo.FindOptions) (res *models2.DataCollection, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdDataCollection).Get(query, opts)
	return convertTypeDataCollection(d, err)
}

func (svc *Service) GetDataCollectionList(query bson.M, opts *mongo.FindOptions) (res []models2.DataCollection, err error) {
	err = svc.getListSerializeTarget(interfaces.ModelIdDataCollection, query, opts, &res)
	return res, err
}

func (svc *Service) GetDataCollectionByName(name string, opts *mongo.FindOptions) (res *models2.DataCollection, err error) {
	query := bson.M{"name": name}
	return svc.GetDataCollection(query, opts)
}
