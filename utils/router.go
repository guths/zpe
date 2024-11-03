package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guths/zpe/constants"
	"github.com/guths/zpe/datatransfers"
)

const (
	ADMIN_ROLE_LVL    = 1
	MODIFIER_ROLE_LVL = 2
	WATCHER_ROLE_LVL  = 3
)

func AuthOnly(c *gin.Context) {
	if !c.GetBool(constants.IsAuthenticatedKey) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, datatransfers.Response{Error: "user not authenticated"})
	}
}

func AdminOnly(c *gin.Context) {
	if !checkLvl(c, ADMIN_ROLE_LVL) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, datatransfers.Response{Error: "permission denied"})
		return
	}
}

func ModifierOnly(c *gin.Context) {
	if !checkLvl(c, MODIFIER_ROLE_LVL) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, datatransfers.Response{Error: "permission denied"})
		return
	}
}

func WatcherOnly(c *gin.Context) {
	if !checkLvl(c, WATCHER_ROLE_LVL) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, datatransfers.Response{Error: "permission denied"})
		return
	}
}

func checkLvl(c *gin.Context, routeLvl int) bool {
	roleLvl := c.GetInt(constants.UserRoleLvl)

	if roleLvl == 0 {
		return false
	}

	if roleLvl > routeLvl {
		return false
	}

	return true
}
