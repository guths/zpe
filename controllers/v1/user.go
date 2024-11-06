package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guths/zpe/constants"
	"github.com/guths/zpe/datatransfers"
	"github.com/guths/zpe/handlers"
	"github.com/guths/zpe/models"
)

// simple query by user email
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

// delete a user with a higher permission than the auth user is not allowed
// is not possible to delete a user that not exists
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

	err = handlers.Handler.DeleteUser(userInfo.Email)

	if err != nil {
		c.JSON(http.StatusNotFound, datatransfers.Response{Error: "error deleting the user"})
		return
	}

	c.JSON(http.StatusOK, datatransfers.Response{Message: "user sucessfully deleted"})
}

// the updated user can not have a higher permission than the authenticated user
// is not possible to delete a user that not exists
func PUTUser(c *gin.Context) {
	var err error
	var userURI datatransfers.UserUpdateURI
	var userInfo datatransfers.UserUpdate
	if err = c.ShouldBindUri(&userURI); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.Response{Error: err.Error()})
		return
	}

	if err = c.ShouldBindJSON(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.Response{Error: err.Error()})
		return
	}

	user, err := handlers.Handler.RetrieveUser(userURI.Email)

	if err != nil {
		c.JSON(http.StatusNotFound, datatransfers.Response{Error: "user not found"})
		return
	}

	currentLvl := c.GetInt(constants.UserRoleLvl)

	lvl := models.GetMaxRoleLvl(user.Roles)

	fmt.Println(currentLvl, lvl)

	if lvl < currentLvl || lvl == 0 {
		c.JSON(http.StatusUnauthorized, datatransfers.Response{Error: "permission denied"})
		return
	}

	if !handlers.Handler.ValidateUserRoles(lvl, userInfo.Roles) {
		c.JSON(http.StatusUnauthorized, datatransfers.Response{Error: "permission denied, user trying to put a higher role than is permitted"})
		return
	}

	updatedUser, err := handlers.Handler.UpdateUser(user.ID, userInfo)

	if err != nil {
		c.JSON(http.StatusNotModified, datatransfers.Response{Error: "failed updating user"})
		return
	}

	c.JSON(http.StatusOK, datatransfers.Response{Data: updatedUser})
}

// the user that is being created can not have a higher permission than the authenticated user
func POSTUser(c *gin.Context) {
	var err error
	var userInfo datatransfers.UserSignup

	if err = c.ShouldBindJSON(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.Response{Error: err.Error()})
		return
	}

	currentLvl := c.GetInt(constants.UserRoleLvl)

	if !handlers.Handler.ValidateUserRoles(currentLvl, userInfo.Roles) {
		c.JSON(http.StatusUnauthorized, datatransfers.Response{Error: "permission denied, user trying to put a higher role than is permitted"})
		return
	}

	res := handlers.Handler.RegisterUser(userInfo)

	if res.Error != "" {
		c.JSON(http.StatusBadRequest, datatransfers.Response{Error: res.Error})
		return
	}

	c.JSON(http.StatusOK, datatransfers.Response{Data: res.Data})
}
