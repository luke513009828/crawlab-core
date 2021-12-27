package controllers

import (
	"github.com/luke513009828/crawlab-core/interfaces"
	"github.com/luke513009828/crawlab-core/models/service"
)

func InitControllers() (err error) {
	modelSvc, err := service.GetService()
	if err != nil {
		return err
	}

	NodeController = newNodeController()
	ProjectController = newProjectController()
	SpiderController = newSpiderController()
	TaskController = newTaskController()
	UserController = newUserController()
	TagController = NewListControllerDelegate(ControllerIdTag, modelSvc.GetBaseService(interfaces.ModelIdTag))
	SettingController = newSettingController()
	LoginController = NewActionControllerDelegate(ControllerIdLogin, getLoginActions())
	ColorController = NewActionControllerDelegate(ControllerIdColor, getColorActions())
	PluginController = newPluginController()
	DataCollectionController = NewListControllerDelegate(ControllerIdDataCollection, modelSvc.GetBaseService(interfaces.ModelIdDataCollection))
	ResultController = NewActionControllerDelegate(ControllerIdResult, getResultActions())
	ScheduleController = newScheduleController()
	StatsController = NewActionControllerDelegate(ControllerIdStats, getStatsActions())
	TokenController = newTokenController()
	FilerController = NewActionControllerDelegate(ControllerIdFiler, getFilerActions())
	PluginProxyController = NewActionControllerDelegate(ControllerIdPluginDo, getPluginProxyActions())
	GitController = NewListControllerDelegate(ControllerIdGit, modelSvc.GetBaseService(interfaces.ModelIdGit))
	VersionController = NewActionControllerDelegate(ControllerIdVersion, getVersionActions())

	return nil
}
