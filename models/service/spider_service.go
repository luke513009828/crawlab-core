package service

import (
	"github.com/luke513009828/crawlab-core/errors"
	"github.com/luke513009828/crawlab-core/interfaces"
	models2 "github.com/luke513009828/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeSpider(d interface{}, err error) (res *models2.Spider, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*models2.Spider)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetSpiderById(id primitive.ObjectID) (res *models2.Spider, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdSpider).GetById(id)
	return convertTypeSpider(d, err)
}

func (svc *Service) GetSpider(query bson.M, opts *mongo.FindOptions) (res *models2.Spider, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdSpider).Get(query, opts)
	return convertTypeSpider(d, err)
}

func (svc *Service) GetSpiderList(query bson.M, opts *mongo.FindOptions) (res []models2.Spider, err error) {
	err = svc.getListSerializeTarget(interfaces.ModelIdSpider, query, opts, &res)
	return res, err
}
