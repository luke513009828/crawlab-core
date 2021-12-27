package result

import (
	"github.com/luke513009828/crawlab-core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Option func(svc interfaces.ResultService)

func WithId(id primitive.ObjectID) Option {
	return func(svc interfaces.ResultService) {
		svc.SetId(id)
	}
}
