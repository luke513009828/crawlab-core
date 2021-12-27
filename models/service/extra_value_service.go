package service

import (
	"github.com/luke513009828/crawlab-core/errors"
	"github.com/luke513009828/crawlab-core/interfaces"
	"github.com/luke513009828/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeExtraValue(d interface{}, err error) (res *models.ExtraValue, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*models.ExtraValue)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetExtraValueById(id primitive.ObjectID) (res *models.ExtraValue, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdExtraValue).GetById(id)
	return convertTypeExtraValue(d, err)
}

func (svc *Service) GetExtraValue(query bson.M, opts *mongo.FindOptions) (res *models.ExtraValue, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdExtraValue).Get(query, opts)
	return convertTypeExtraValue(d, err)
}

func (svc *Service) GetExtraValueList(query bson.M, opts *mongo.FindOptions) (res []models.ExtraValue, err error) {
	err = svc.getListSerializeTarget(interfaces.ModelIdExtraValue, query, opts, &res)
	return res, err
}

func (svc *Service) GetExtraValueByObjectIdModel(oid primitive.ObjectID, m string, opts *mongo.FindOptions) (res *models.ExtraValue, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdExtraValue).Get(bson.M{"oid": oid, "m": m}, opts)
	return convertTypeExtraValue(d, err)
}
