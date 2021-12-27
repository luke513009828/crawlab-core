package service

import (
	"github.com/luke513009828/crawlab-core/errors"
	"github.com/luke513009828/crawlab-core/interfaces"
	models2 "github.com/luke513009828/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeDataSource(d interface{}, err error) (res *models2.DataSource, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*models2.DataSource)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetDataSourceById(id primitive.ObjectID) (res *models2.DataSource, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdDataSource).GetById(id)
	return convertTypeDataSource(d, err)
}

func (svc *Service) GetDataSource(query bson.M, opts *mongo.FindOptions) (res *models2.DataSource, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdDataSource).Get(query, opts)
	return convertTypeDataSource(d, err)
}

func (svc *Service) GetDataSourceList(query bson.M, opts *mongo.FindOptions) (res []models2.DataSource, err error) {
	err = svc.getListSerializeTarget(interfaces.ModelIdDataSource, query, opts, &res)
	return res, err
}
