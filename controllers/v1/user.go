package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guths/zpe/constants"
	"github.com/guths/zpe/datatransfers"
	"github.com/guths/zpe/handlers"
	"github.com/guths/zpe/models"
)

func GETUser(c *gin.Context) {
	var err error
	var userInfo datatransfers.UserInfo
	if err = c.ShouldBindUri(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.Response{Error: err.Error()})
		return
	}

	user, err := handlers.Handler.RetrieveUser(userInfo.Email)

	if err != nil {
		c.JSON(http.StatusNotFound, datatransfers.Response{
			Code:  400,
			Error: "user not found",
		})

		return
	}

	c.JSON(http.StatusOK, datatransfers.Response{
		Code: 200,
		Data: user,
	})
}

func DELETEUser(c *gin.Context) {
	var err error
	var userInfo datatransfers.UserInfo
	if err = c.ShouldBindUri(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.Response{Error: err.Error()})
	}

	user, err := handlers.Handler.RetrieveUser(userInfo.Email)

	if err != nil {
		c.JSON(http.StatusNotFound, datatransfers.Response{Error: "user not found"})
		return
	}

	currentLvl := c.GetInt(constants.UserRoleLvl)

	lvl := models.GetMaxRoleLvl(user.Roles)

	if lvl < currentLvl || lvl == 0 {
		c.JSON(http.StatusUnauthorized, datatransfers.Response{Error: "permission denied"})
		return
	}

	res := handlers.Handler.DeleteUser(userInfo.Email)

	if res.Error != "" {
		c.JSON(http.StatusNotFound, res)
		return
	}

	c.JSON(http.StatusOK, res)
}

func PUTUser(c *gin.Context) {
	var err error
	var userInfo datatransfers.UserUpdate
	if err = c.ShouldBind(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.Response{Error: err.Error()})
		return
	}

	user, err := handlers.Handler.RetrieveUser(userInfo.Email)

	if err != nil {
		c.JSON(http.StatusNotFound, datatransfers.Response{Error: "user not found"})
		return
	}

	currentLvl := c.GetInt(constants.UserRoleLvl)

	lvl := models.GetMaxRoleLvl(user.Roles)

	if lvl < currentLvl || lvl == 0 {
		c.JSON(http.StatusUnauthorized, datatransfers.Response{Error: "permission denied"})
		return
	}

	if !handlers.Handler.ValidateUserRoles(lvl, userInfo.Roles) {
		c.JSON(http.StatusUnauthorized, datatransfers.Response{Error: "permission denied, user trying to put a higher role than is permitted"})
		return
	}

	if err = handlers.Handler.UpdateUser(uint(c.GetInt(constants.IsAuthenticatedKey)), user); err != nil {
		c.JSON(http.StatusNotModified, datatransfers.Response{Error: "failed updating user"})
		return
	}
	c.JSON(http.StatusOK, datatransfers.Response{Data: user})
}
