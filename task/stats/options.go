package stats

import "github.com/luke513009828/crawlab-core/interfaces"

type Option func(service interfaces.TaskStatsService)

func WithConfigPath(path string) Option {
	return func(svc interfaces.TaskStatsService) {
		svc.SetConfigPath(path)
	}
}
