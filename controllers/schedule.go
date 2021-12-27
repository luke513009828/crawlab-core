package controllers

import (
	"github.com/luke513009828/crawlab-core/config"
	"github.com/luke513009828/crawlab-core/interfaces"
	"github.com/luke513009828/crawlab-core/models/delegate"
	"github.com/luke513009828/crawlab-core/models/models"
	"github.com/luke513009828/crawlab-core/models/service"
	"github.com/luke513009828/crawlab-core/schedule"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/dig"
	"net/http"
)

var ScheduleController *scheduleController

func getScheduleActions() []Action {
	scheduleCtx := newScheduleContext()
	return []Action{
		{
			Method:      http.MethodPost,
			Path:        "/:id/enable",
			HandlerFunc: scheduleCtx.enable,
		},
		{
			Method:      http.MethodPost,
			Path:        "/:id/disable",
			HandlerFunc: scheduleCtx.disable,
		},
	}
}

type scheduleController struct {
	ListActionControllerDelegate
	d   ListActionControllerDelegate
	ctx *scheduleContext
}

func (ctr *scheduleController) Put(c *gin.Context) {
	var s models.Schedule
	if err := c.ShouldBindJSON(&s); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	if err := delegate.NewModelDelegate(&s, GetUserFromContext(c)).Add(); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	if s.Enabled {
		if err := ctr.ctx.scheduleSvc.Enable(&s, GetUserFromContext(c)); err != nil {
			HandleErrorInternalServerError(c, err)
			return
		}
	}
	HandleSuccessWithData(c, s)
}

func (ctr *scheduleController) Delete(c *gin.Context) {
	id := c.Param("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	s, err := ctr.ctx.modelSvc.GetScheduleById(oid)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	if err := ctr.ctx.scheduleSvc.Disable(s); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	if err := delegate.NewModelDelegate(s, GetUserFromContext(c)).Delete(); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
}

func (ctr *scheduleController) DeleteList(c *gin.Context) {
	payload, err := NewJsonBinder(interfaces.ModelIdSchedule).BindBatchRequestPayload(c)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	for _, id := range payload.Ids {
		s, err := ctr.ctx.modelSvc.GetScheduleById(id)
		if err != nil {
			HandleErrorInternalServerError(c, err)
			return
		}
		if err := ctr.ctx.scheduleSvc.Disable(s); err != nil {
			HandleErrorInternalServerError(c, err)
			return
		}
	}
	if err := ctr.ctx.modelSvc.GetBaseService(interfaces.ModelIdSchedule).DeleteList(bson.M{
		"_id": bson.M{
			"$in": payload.Ids,
		},
	}); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
}

func (ctx *scheduleContext) enable(c *gin.Context) {
	s, err := ctx._getSchedule(c)
	if err != nil {
		return
	}
	if err := ctx.scheduleSvc.Enable(s, GetUserFromContext(c)); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
}

func (ctx *scheduleContext) disable(c *gin.Context) {
	s, err := ctx._getSchedule(c)
	if err != nil {
		return
	}
	if err := ctx.scheduleSvc.Disable(s, GetUserFromContext(c)); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
}

func (ctx *scheduleContext) _getSchedule(c *gin.Context) (s *models.Schedule, err error) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	s, err = ctx.modelSvc.GetScheduleById(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	return s, nil
}

type scheduleContext struct {
	modelSvc    service.ModelService
	scheduleSvc interfaces.ScheduleService
}

func newScheduleContext() *scheduleContext {
	// context
	ctx := &scheduleContext{}

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.NewService); err != nil {
		panic(err)
	}
	if err := c.Provide(schedule.ProvideGetScheduleService(config.DefaultConfigPath)); err != nil {
		panic(err)
	}
	if err := c.Invoke(func(
		modelSvc service.ModelService,
		scheduleSvc interfaces.ScheduleService,
	) {
		ctx.modelSvc = modelSvc
		ctx.scheduleSvc = scheduleSvc
	}); err != nil {
		panic(err)
	}

	return ctx
}

func newScheduleController() *scheduleController {
	actions := getScheduleActions()
	modelSvc, err := service.GetService()
	if err != nil {
		panic(err)
	}

	ctr := NewListPostActionControllerDelegate(ControllerIdSchedule, modelSvc.GetBaseService(interfaces.ModelIdSchedule), actions)
	d := NewListPostActionControllerDelegate(ControllerIdSchedule, modelSvc.GetBaseService(interfaces.ModelIdSchedule), actions)
	ctx := newScheduleContext()

	return &scheduleController{
		ListActionControllerDelegate: *ctr,
		d:                            *d,
		ctx:                          ctx,
	}
}
