package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/guths/zpe/config"
	"github.com/guths/zpe/constants"
	"github.com/guths/zpe/models"
)

//here the user role is validated, the level of the role is measured

func getRoles(c *gin.Context) ([]models.Role, error) {
	roleNames := c.GetStringSlice(constants.UserRoles)

	roleO := models.NewRoleOrmer(config.DBManager.GetDB())

	roles, err := roleO.GetManyByName(roleNames)

	if err != nil {
		return []models.Role{}, err
	}

	return roles, nil
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
