package controllers

import (
	"github.com/luke513009828/crawlab-core/interfaces"
)

func NewListPostActionControllerDelegate(id ControllerId, svc interfaces.ModelBaseService, actions []Action) (d *ListActionControllerDelegate) {
	return &ListActionControllerDelegate{
		NewListControllerDelegate(id, svc),
		NewActionControllerDelegate(id, actions),
	}
}

type ListActionControllerDelegate struct {
	ListController
	ActionController
}
