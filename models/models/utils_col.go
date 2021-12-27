package models

import (
	"github.com/luke513009828/crawlab-core/interfaces"
	"github.com/luke513009828/crawlab-core/utils/binders"
)

func GetModelColName(id interfaces.ModelId) (colName string) {
	return binders.NewColNameBinder(id).MustBindString()
}
