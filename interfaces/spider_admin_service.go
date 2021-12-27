package interfaces

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SpiderAdminService interface {
	WithConfigPath
	// Schedule a new task of the spider
	Schedule(id primitive.ObjectID, opts *SpiderRunOptions) (err error)

	// Schedule a new task of the spider and return task id
	ScheduleWithTaskId(id primitive.ObjectID, opts *SpiderRunOptions) (taskIds []primitive.ObjectID, err error)

	// Clone the spider
	Clone(id primitive.ObjectID, opts *SpiderCloneOptions) (err error)
	// Delete the spider
	Delete(id primitive.ObjectID) (err error)
}
