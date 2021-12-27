package client

import (
	"encoding/json"
	"github.com/luke513009828/crawlab-core/errors"
	"github.com/luke513009828/crawlab-core/interfaces"
	"github.com/luke513009828/crawlab-core/models/models"
	grpc "github.com/crawlab-team/crawlab-grpc"
	"github.com/crawlab-team/go-trace"
)

func NewBasicBinder(id interfaces.ModelId, res *grpc.Response) (b interfaces.GrpcModelBinder) {
	return &BasicBinder{
		id:  id,
		res: res,
	}
}

type BasicBinder struct {
	id  interfaces.ModelId
	res *grpc.Response
}

func (b *BasicBinder) Bind() (res interfaces.Model, err error) {
	m := models.NewModelMap()

	switch b.id {
	case interfaces.ModelIdArtifact:
		return b.Process(&m.Artifact)
	case interfaces.ModelIdTag:
		return b.Process(&m.Tag)
	case interfaces.ModelIdNode:
		return b.Process(&m.Node)
	case interfaces.ModelIdProject:
		return b.Process(&m.Project)
	case interfaces.ModelIdSpider:
		return b.Process(&m.Spider)
	case interfaces.ModelIdTask:
		return b.Process(&m.Task)
	case interfaces.ModelIdJob:
		return b.Process(&m.Job)
	case interfaces.ModelIdSchedule:
		return b.Process(&m.Schedule)
	case interfaces.ModelIdUser:
		return b.Process(&m.User)
	case interfaces.ModelIdSetting:
		return b.Process(&m.Setting)
	case interfaces.ModelIdToken:
		return b.Process(&m.Token)
	case interfaces.ModelIdVariable:
		return b.Process(&m.Variable)
	case interfaces.ModelIdTaskQueue:
		return b.Process(&m.TaskQueueItem)
	case interfaces.ModelIdTaskStat:
		return b.Process(&m.TaskStat)
	case interfaces.ModelIdPlugin:
		return b.Process(&m.Plugin)
	case interfaces.ModelIdSpiderStat:
		return b.Process(&m.SpiderStat)
	case interfaces.ModelIdDataSource:
		return b.Process(&m.DataSource)
	case interfaces.ModelIdDataCollection:
		return b.Process(&m.DataCollection)
	case interfaces.ModelIdResult:
		return b.Process(&m.Result)
	case interfaces.ModelIdPassword:
		return b.Process(&m.Password)
	case interfaces.ModelIdExtraValue:
		return b.Process(&m.ExtraValue)
	case interfaces.ModelIdPluginStatus:
		return b.Process(&m.PluginStatus)
	case interfaces.ModelIdGit:
		return b.Process(&m.Git)
	default:
		return nil, errors.ErrorModelInvalidModelId
	}
}

func (b *BasicBinder) MustBind() (res interfaces.Model) {
	res, err := b.Bind()
	if err != nil {
		panic(err)
	}
	return res
}

func (b *BasicBinder) Process(d interfaces.Model) (res interfaces.Model, err error) {
	if err := json.Unmarshal(b.res.Data, d); err != nil {
		return nil, trace.TraceError(err)
	}
	return d, nil
}
