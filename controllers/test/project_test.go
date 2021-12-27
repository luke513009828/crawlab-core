package test

import (
	"encoding/json"
	"fmt"
	"github.com/luke513009828/crawlab-core/constants"
	"github.com/luke513009828/crawlab-core/entity"
	"github.com/luke513009828/crawlab-core/models/delegate"
	"github.com/luke513009828/crawlab-core/models/models"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"testing"
	"time"
)

func TestProjectController_Get(t *testing.T) {
	T.Setup(t)
	e := T.NewExpect(t)

	p := models.Project{
		Name: "test project",
	}
	res := e.PUT("/projects").WithJSON(p).Expect().Status(http.StatusOK).JSON().Object()
	res.Path("$.data._id").NotNull()
	id := res.Path("$.data._id").String().Raw()
	oid, err := primitive.ObjectIDFromHex(id)
	require.Nil(t, err)
	require.False(t, oid.IsZero())

	res = e.GET("/projects/" + id).WithJSON(p).Expect().Status(http.StatusOK).JSON().Object()
	res.Path("$.data._id").NotNull()
	res.Path("$.data.name").Equal("test project")
}

func TestProjectController_Post(t *testing.T) {
	T.Setup(t)
	e := T.NewExpect(t)

	p := models.Project{
		Name:        "old name",
		Description: "old description",
	}

	// add
	res := e.PUT("/projects").
		WithJSON(p).
		Expect().Status(http.StatusOK).
		JSON().Object()
	res.Path("$.data._id").NotNull()
	id := res.Path("$.data._id").String().Raw()
	oid, err := primitive.ObjectIDFromHex(id)
	require.Nil(t, err)
	require.False(t, oid.IsZero())

	// change object
	p.Id = oid
	p.Name = "new name"
	p.Description = "new description"

	// update
	e.POST("/projects/" + id).
		WithJSON(p).
		Expect().Status(http.StatusOK)

	// check
	res = e.GET("/projects/" + id).Expect().Status(http.StatusOK).JSON().Object()
	res.Path("$.data._id").Equal(id)
	res.Path("$.data.name").Equal("new name")
	res.Path("$.data.description").Equal("new description")
}

func TestProjectController_Put(t *testing.T) {
	T.Setup(t)
	e := T.NewExpect(t)

	p := models.Project{
		Name:        "test project",
		Description: "this is a test project",
	}

	res := e.PUT("/projects").WithJSON(p).Expect().Status(http.StatusOK).JSON().Object()
	res.Path("$.data._id").NotNull()
	res.Path("$.data.name").Equal("test project")
	res.Path("$.data.description").Equal("this is a test project")
}

func TestProjectController_Delete(t *testing.T) {
	T.Setup(t)
	e := T.NewExpect(t)

	p := models.Project{
		Name:        "test project",
		Description: "this is a test project",
	}

	// add
	res := e.PUT("/projects").
		WithJSON(p).
		Expect().Status(http.StatusOK).
		JSON().Object()
	res.Path("$.data._id").NotNull()
	id := res.Path("$.data._id").String().Raw()
	oid, err := primitive.ObjectIDFromHex(id)
	require.Nil(t, err)
	require.False(t, oid.IsZero())

	// get
	res = e.GET("/projects/" + id).
		Expect().Status(http.StatusOK).
		JSON().Object()
	res.Path("$.data._id").NotNull()
	id = res.Path("$.data._id").String().Raw()
	oid, err = primitive.ObjectIDFromHex(id)
	require.Nil(t, err)
	require.False(t, oid.IsZero())

	// delete
	e.DELETE("/projects/" + id).
		Expect().Status(http.StatusOK).
		JSON().Object()

	// get
	e.GET("/projects/" + id).
		Expect().Status(http.StatusNotFound)
}

func TestProjectController_GetList(t *testing.T) {
	T.Setup(t)
	e := T.NewExpect(t)

	n := 100 // total
	bn := 10 // batch

	for i := 0; i < n; i++ {
		p := models.Project{
			Name: fmt.Sprintf("test name %d", i+1),
		}
		obj := e.PUT("/projects").WithJSON(p).Expect().Status(http.StatusOK).JSON().Object()
		obj.Path("$.data._id").NotNull()
	}

	f := entity.Filter{
		//IsOr: false,
		Conditions: []entity.Condition{
			{Key: "name", Op: constants.FilterOpContains, Value: "test name"},
		},
	}
	condBytes, err := json.Marshal(&f.Conditions)
	require.Nil(t, err)

	pagination := entity.Pagination{
		Page: 1,
		Size: bn,
	}

	// get list with pagination
	res := e.GET("/projects").
		WithQuery("conditions", string(condBytes)).
		WithQueryObject(pagination).
		Expect().Status(http.StatusOK).JSON().Object()
	res.Path("$.data").Array().Length().Equal(bn)
	res.Path("$.total").Number().Equal(n)

	data := res.Path("$.data").Array()
	for i := 0; i < bn; i++ {
		obj := data.Element(i)
		obj.Path("$.name").Equal(fmt.Sprintf("test name %d", i+1))
	}

	// get all
	res = e.GET("/projects").
		WithQuery("all", "1").
		Expect().Status(http.StatusOK).JSON().Object()
	res.Path("$.data").Array().Length().Equal(n)
	res.Path("$.total").Number().Equal(n)

}

func TestProjectController_PutList(t *testing.T) {
	T.Setup(t)
	e := T.NewExpect(t)

	n := 10
	var docs []models.Project
	for i := 0; i < n; i++ {
		docs = append(docs, models.Project{
			Name:        fmt.Sprintf("project %d", i+1),
			Description: "this is a project",
		})
	}

	e.PUT("/projects/batch").WithJSON(docs).Expect().Status(http.StatusOK)

	res := e.GET("/projects").
		WithQueryObject(entity.Pagination{Page: 1, Size: 100}).
		Expect().Status(http.StatusOK).
		JSON().Object()
	res.Path("$.data").Array().Length().Equal(n)
}

func TestProjectController_DeleteList(t *testing.T) {
	T.Setup(t)
	e := T.NewExpect(t)

	n := 10
	var docs []models.Project
	for i := 0; i < n; i++ {
		docs = append(docs, models.Project{
			Name:        fmt.Sprintf("project %d", i+1),
			Description: "this is a project",
		})
	}

	// add
	res := e.PUT("/projects/batch").WithJSON(docs).Expect().Status(http.StatusOK).
		JSON().Object()
	var ids []primitive.ObjectID
	data := res.Path("$.data").Array()
	for i := 0; i < n; i++ {
		obj := data.Element(i)
		id := obj.Path("$._id").String().Raw()
		oid, err := primitive.ObjectIDFromHex(id)
		require.Nil(t, err)
		require.False(t, oid.IsZero())
		ids = append(ids, oid)
	}

	// delete
	payload := entity.BatchRequestPayload{
		Ids: ids,
	}
	e.DELETE("/projects").
		WithJSON(payload).
		Expect().Status(http.StatusOK)

	// check
	res = e.GET("/projects").
		Expect().Status(http.StatusOK).JSON().Object()
	res.Path("$.data").Array().Empty()
	res.Path("$.total").Number().Equal(0)

}

func TestProjectController_PostList(t *testing.T) {
	T.Setup(t)
	e := T.NewExpect(t)

	// now
	now := time.Now()

	n := 10
	var docs []models.Project
	for i := 0; i < n; i++ {
		docs = append(docs, models.Project{
			Name:        "old name",
			Description: "old description",
		})
	}

	// add
	res := e.PUT("/projects/batch").WithJSON(docs).Expect().Status(http.StatusOK).
		JSON().Object()
	var ids []primitive.ObjectID
	data := res.Path("$.data").Array()
	for i := 0; i < n; i++ {
		obj := data.Element(i)
		id := obj.Path("$._id").String().Raw()
		oid, err := primitive.ObjectIDFromHex(id)
		require.Nil(t, err)
		require.False(t, oid.IsZero())
		ids = append(ids, oid)
	}

	// wait for 100 millisecond
	time.Sleep(100 * time.Millisecond)

	// update
	p := models.Project{
		Name:        "new name",
		Description: "new description",
	}
	dataBytes, err := json.Marshal(&p)
	require.Nil(t, err)
	payload := entity.BatchRequestPayloadWithStringData{
		Ids:  ids,
		Data: string(dataBytes),
		Fields: []string{
			"name",
			"description",
		},
	}
	e.POST("/projects").
		WithJSON(payload).
		Expect().Status(http.StatusOK)

	// check response data
	res = e.GET("/projects").
		WithQueryObject(entity.Pagination{
			Page: 1,
			Size: 100,
		}).
		Expect().Status(http.StatusOK).JSON().Object()
	data = res.Path("$.data").Array()
	for i := 0; i < n; i++ {
		obj := data.Element(i)
		obj.Path("$.name").Equal("new name")
		obj.Path("$.description").Equal("new description")
	}

	// check artifacts
	pl, err := T.modelSvc.GetProjectList(bson.M{"_id": bson.M{"$in": ids}}, nil)
	require.Nil(t, err)
	for _, p := range pl {
		a, err := delegate.NewModelDelegate(&p).GetArtifact()
		require.Nil(t, err)
		require.True(t, a.GetSys().GetUpdateTs().After(now))
		require.True(t, a.GetSys().GetUpdateTs().After(a.GetSys().GetCreateTs()))
	}
}
