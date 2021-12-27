package models

import "github.com/luke513009828/crawlab-core/interfaces"

type modelRelation struct {
	d       interfaces.Model
	id      interfaces.ModelId
	colName string
}

var ModelRelations = []modelRelation{
	{d: &Job{}, id: interfaces.ModelIdJob, colName: interfaces.ModelColNameJob},
	{d: &Node{}, id: interfaces.ModelIdNode, colName: interfaces.ModelColNameNode},
}
