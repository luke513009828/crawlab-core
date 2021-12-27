package controllers

import (
	"github.com/luke513009828/crawlab-core/interfaces"
	"github.com/luke513009828/crawlab-core/models/delegate"
	"github.com/luke513009828/crawlab-core/models/models"
	"github.com/luke513009828/crawlab-core/models/service"
	"github.com/gin-gonic/gin"
)

var SettingController *settingController

type settingController struct {
	ListControllerDelegate
}

func (ctr *settingController) Get(c *gin.Context) {
	// key
	key := c.Param("id")

	// model service
	modelSvc, err := service.NewService()
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// setting
	s, err := modelSvc.GetSettingByKey(key, nil)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithData(c, s)
}

func (ctr *settingController) Post(c *gin.Context) {
	// key
	key := c.Param("id")

	// settings
	var s models.Setting
	if err := c.ShouldBindJSON(&s); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// model service
	modelSvc, err := service.NewService()
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// setting
	_s, err := modelSvc.GetSettingByKey(key, nil)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// save
	_s.Value = s.Value
	if err := delegate.NewModelDelegate(_s).Save(); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func newSettingController() *settingController {
	modelSvc, err := service.GetService()
	if err != nil {
		panic(err)
	}

	ctr := NewListControllerDelegate(ControllerIdSetting, modelSvc.GetBaseService(interfaces.ModelIdSetting))

	return &settingController{
		ListControllerDelegate: *ctr,
	}
}
