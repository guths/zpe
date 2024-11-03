package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guths/zpe/config"
	"github.com/guths/zpe/constants"
	"github.com/guths/zpe/datatransfers"
	"github.com/guths/zpe/models"
)

func getRoles(c *gin.Context) ([]models.Role, error) {
	roleNames := c.GetStringSlice(constants.UserRoles)

	roleO := models.NewRoleOrmer(config.DBManager.GetDB())

	roles, err := roleO.GetManyByName(roleNames)

	if err != nil {
		return []models.Role{}, err
	}

	return roles, nil
}

func getRouteRoleLvl(c *gin.Context) (int, error) {
	var err error
	var userInfo datatransfers.UserInfo

	if err = c.ShouldBindUri(&userInfo); err != nil {
		return 0, fmt.Errorf("user not found")
	}

	userO := models.NewUserOrmer(config.DBManager.GetDB())

	user, err := userO.GetOneByEmail(userInfo.Email)

	if err != nil {
		c.JSON(http.StatusNotFound, datatransfers.Response{Error: "user not found"})
	}

	return models.GetMaxRoleLvl(user.Roles), nil
}

func RoleMiddleware(c *gin.Context) {
	roles, err := getRoles(c)

	if err != nil || len(roles) == 0 {
		c.Next()
		return
	}

	userMaxLvl := models.GetMaxRoleLvl(roles)

	c.Set(constants.UserRoleLvl, userMaxLvl)

	c.Next()
}
