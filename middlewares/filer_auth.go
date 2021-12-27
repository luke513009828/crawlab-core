package middlewares

import (
	"github.com/luke513009828/crawlab-core/controllers"
	"github.com/luke513009828/crawlab-core/errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func FilerAuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// auth key
		authKey := c.GetHeader("Authorization")

		// server auth key
		svrAuthKey := viper.GetString("fs.filer.authKey")

		// skip to next if no server auth key is provided
		if svrAuthKey == "" {
			c.Next()
			return
		}

		// validate
		if authKey != svrAuthKey {
			// validation failed, return error response
			controllers.HandleErrorUnauthorized(c, errors.ErrorHttpUnauthorized)
			return
		}

		// validation success
		c.Next()
	}
}
