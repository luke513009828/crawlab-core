package scheduler

import (
	"github.com/luke513009828/crawlab-core/interfaces"
	"time"
)

type Option func(svc interfaces.TaskSchedulerService)

func WithConfigPath(path string) Option {
	return func(svc interfaces.TaskSchedulerService) {
		svc.SetConfigPath(path)
	}
}

func WithInterval(interval time.Duration) Option {
	return func(svc interfaces.TaskSchedulerService) {
		svc.SetInterval(interval)
	}
}
