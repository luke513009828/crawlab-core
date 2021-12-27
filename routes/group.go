package routes

import (
	"github.com/luke513009828/crawlab-core/middlewares"
	"github.com/gin-gonic/gin"
)

type RouterGroups struct {
	AuthGroup      *gin.RouterGroup
	AnonymousGroup *gin.RouterGroup
	FilerGroup     *gin.RouterGroup
}

func NewRouterGroups(app *gin.Engine) (groups *RouterGroups) {
	return &RouterGroups{
		AuthGroup:      app.Group("/", middlewares.AuthorizationMiddleware()),
		AnonymousGroup: app.Group("/"),
		FilerGroup:     app.Group("/filer", middlewares.FilerAuthorizationMiddleware()),
	}
}
