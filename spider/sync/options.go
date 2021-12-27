package sync

import (
	"github.com/luke513009828/crawlab-core/interfaces"
)

type Option func(svc interfaces.SpiderSyncService)

func WithConfigPath(path string) Option {
	return func(svc interfaces.SpiderSyncService) {
		svc.SetConfigPath(path)
	}
}

func WithFsPathBase(path string) Option {
	return func(svc interfaces.SpiderSyncService) {
		svc.SetFsPathBase(path)
	}
}

func WithWorkspacePathBase(path string) Option {
	return func(svc interfaces.SpiderSyncService) {
		svc.SetWorkspacePathBase(path)
	}
}

func WithRepoPathBase(path string) Option {
	return func(svc interfaces.SpiderSyncService) {
		svc.SetRepoPathBase(path)
	}
}
