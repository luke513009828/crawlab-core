package service_test

import (
	"github.com/luke513009828/crawlab-core/interfaces"
	"github.com/luke513009828/crawlab-core/models/delegate"
	models2 "github.com/luke513009828/crawlab-core/models/models"
	"github.com/luke513009828/crawlab-core/models/service"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTagService_GetModel(t *testing.T) {
	SetupTest(t)

	node := &models2.Node{
		Name:     "test node",
		IsMaster: true,
		Tags: []models2.Tag{
			{Name: "tag 1", Color: "red"},
		},
	}
	err := delegate.NewModelNodeDelegate(node).Add()
	require.Nil(t, err)

	svc, err := service.NewService()
	require.Nil(t, err)

	tag, err := svc.GetTag(nil, nil)
	require.Nil(t, err)
	require.False(t, tag.Id.IsZero())
	require.Equal(t, "tag 1", tag.Name)
	require.Equal(t, "red", tag.Color)
	require.Equal(t, interfaces.ModelColNameNode, tag.Col)
}

func TestTagService_GetModelById(t *testing.T) {
	SetupTest(t)

	node := &models2.Node{
		Name:     "test node",
		IsMaster: true,
		Tags: []models2.Tag{
			{Name: "tag 1", Color: "red"},
		},
	}
	err := delegate.NewModelNodeDelegate(node).Add()
	require.Nil(t, err)

	svc, err := service.NewService()
	require.Nil(t, err)

	tag, err := svc.GetTag(nil, nil)
	require.Nil(t, err)
	require.False(t, tag.Id.IsZero())
	require.Equal(t, "tag 1", tag.Name)
	require.Equal(t, "red", tag.Color)
	require.Equal(t, interfaces.ModelColNameNode, tag.Col)
}
