package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guths/zpe/datatransfers"
	"github.com/guths/zpe/handlers"
)

func POSTRegister(c *gin.Context) {
	var err error
	var user datatransfers.UserSignup
	if err = c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.Response{Error: err.Error()})
	}

	res := handlers.Handler.RegisterUser(user)

	if res.Error != "" {
		c.JSON(res.Code, res)
		return
	}

	c.JSON(res.Code, res)
}

func POSTLogin(c *gin.Context) {
	var err error
	var user datatransfers.UserLogin
	if err = c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.Response{Error: err.Error()})
		return
	}

	res := handlers.Handler.AuthenticateUser(user)

	if res.Error != "" {
		c.JSON(res.Code, res)
		return
	}

	c.JSON(http.StatusOK, res)
}
