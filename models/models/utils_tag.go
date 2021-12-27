package models

import (
	"github.com/apex/log"
	"github.com/luke513009828/crawlab-core/errors"
	"github.com/luke513009828/crawlab-core/interfaces"
	"github.com/crawlab-team/go-trace"
)

func convertInterfacesToTags(tags []interfaces.Tag) (res []Tag) {
	if tags == nil {
		return nil
	}
	for _, t := range tags {
		tag, ok := t.(*Tag)
		if !ok {
			log.Warnf("%v: cannot convert tag", trace.TraceError(errors.ErrorModelInvalidType))
			return nil
		}
		if tag == nil {
			log.Warnf("%v: cannot convert tag", trace.TraceError(errors.ErrorModelInvalidType))
			return nil
		}
		res = append(res, *tag)
	}
	return res
}

func convertTagsToInterfaces(tags []Tag) (res []interfaces.Tag) {
	for _, t := range tags {
		res = append(res, &t)
	}
	return res
}
