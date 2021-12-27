package handler

import (
	"github.com/luke513009828/crawlab-core/interfaces"
	"time"
)

type Option func(svc interfaces.TaskHandlerService)

func WithConfigPath(path string) Option {
	return func(svc interfaces.TaskHandlerService) {
		svc.SetConfigPath(path)
	}
}

func WithExitWatchDuration(duration time.Duration) Option {
	return func(svc interfaces.TaskHandlerService) {
		svc.SetExitWatchDuration(duration)
	}
}

func WithReportInterval(interval time.Duration) Option {
	return func(svc interfaces.TaskHandlerService) {
		svc.SetReportInterval(interval)
	}
}

type RunnerOption func(r interfaces.TaskRunner)

func WithLogDriverType(driverType string) RunnerOption {
	return func(r interfaces.TaskRunner) {
		r.SetLogDriverType(driverType)
	}
}

func WithSubscribeTimeout(timeout time.Duration) RunnerOption {
	return func(r interfaces.TaskRunner) {
		r.SetSubscribeTimeout(timeout)
	}
}
