package user

import "github.com/luke513009828/crawlab-core/interfaces"

type Option func(svc interfaces.UserService)

func WithJwtSecret(secret string) Option {
	return func(svc interfaces.UserService) {
		svc.SetJwtSecret(secret)
	}
}
