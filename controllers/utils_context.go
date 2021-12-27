package controllers

import (
	"github.com/luke513009828/crawlab-core/constants"
	"github.com/luke513009828/crawlab-core/interfaces"
	"github.com/gin-gonic/gin"
)

func GetUserFromContext(c *gin.Context) (u interfaces.User) {
	value, ok := c.Get(constants.ContextUser)
	if !ok {
		return nil
	}
	u, ok = value.(interfaces.User)
	if !ok {
		return nil
	}
	return u
}
