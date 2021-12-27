package binders

import (
	"github.com/luke513009828/crawlab-core/errors"
	"github.com/luke513009828/crawlab-core/interfaces"
)

func NewColNameBinder(id interfaces.ModelId) (b *ColNameBinder) {
	return &ColNameBinder{id: id}
}

type ColNameBinder struct {
	id interfaces.ModelId
}

func (b *ColNameBinder) Bind() (res interface{}, err error) {
	switch b.id {
	// system models
	case interfaces.ModelIdArtifact:
		return interfaces.ModelColNameArtifact, nil
	case interfaces.ModelIdTag:
		return interfaces.ModelColNameTag, nil

	// operation models
	case interfaces.ModelIdNode:
		return interfaces.ModelColNameNode, nil
	case interfaces.ModelIdProject:
		return interfaces.ModelColNameProject, nil
	case interfaces.ModelIdSpider:
		return interfaces.ModelColNameSpider, nil
	case interfaces.ModelIdTask:
		return interfaces.ModelColNameTask, nil
	case interfaces.ModelIdJob:
		return interfaces.ModelColNameJob, nil
	case interfaces.ModelIdSchedule:
		return interfaces.ModelColNameSchedule, nil
	case interfaces.ModelIdUser:
		return interfaces.ModelColNameUser, nil
	case interfaces.ModelIdSetting:
		return interfaces.ModelColNameSetting, nil
	case interfaces.ModelIdToken:
		return interfaces.ModelColNameToken, nil
	case interfaces.ModelIdVariable:
		return interfaces.ModelColNameVariable, nil
	case interfaces.ModelIdTaskQueue:
		return interfaces.ModelColNameTaskQueue, nil
	case interfaces.ModelIdTaskStat:
		return interfaces.ModelColNameTaskStat, nil
	case interfaces.ModelIdPlugin:
		return interfaces.ModelColNamePlugin, nil
	case interfaces.ModelIdSpiderStat:
		return interfaces.ModelColNameSpiderStat, nil
	case interfaces.ModelIdDataSource:
		return interfaces.ModelColNameDataSource, nil
	case interfaces.ModelIdDataCollection:
		return interfaces.ModelColNameDataCollection, nil
	case interfaces.ModelIdPassword:
		return interfaces.ModelColNamePasswords, nil
	case interfaces.ModelIdExtraValue:
		return interfaces.ModelColNameExtraValues, nil
	case interfaces.ModelIdPluginStatus:
		return interfaces.ModelColNamePluginStatus, nil
	case interfaces.ModelIdGit:
		return interfaces.ModelColNameGit, nil

	// invalid
	default:
		return res, errors.ErrorModelNotImplemented
	}
}

func (b *ColNameBinder) MustBind() (res interface{}) {
	res, err := b.Bind()
	if err != nil {
		panic(err)
	}
	return res
}

func (b *ColNameBinder) BindString() (res string, err error) {
	res_, err := b.Bind()
	if err != nil {
		return "", err
	}
	res = res_.(string)
	return res, nil
}

func (b *ColNameBinder) MustBindString() (res string) {
	return b.MustBind().(string)
}
