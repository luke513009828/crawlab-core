package service

import (
	"github.com/luke513009828/crawlab-core/errors"
	"github.com/luke513009828/crawlab-core/interfaces"
	models2 "github.com/luke513009828/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeGit(d interface{}, err error) (res *models2.Git, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*models2.Git)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetGitById(id primitive.ObjectID) (res *models2.Git, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdGit).GetById(id)
	return convertTypeGit(d, err)
}

func (svc *Service) GetGit(query bson.M, opts *mongo.FindOptions) (res *models2.Git, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdGit).Get(query, opts)
	return convertTypeGit(d, err)
}

func (svc *Service) GetGitList(query bson.M, opts *mongo.FindOptions) (res []models2.Git, err error) {
	err = svc.getListSerializeTarget(interfaces.ModelIdGit, query, opts, &res)
	return res, err
}
