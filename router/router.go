package router

import (
	"github.com/gin-gonic/gin"
	"github.com/guths/zpe/controllers/middleware"
	v1 "github.com/guths/zpe/controllers/v1"
	"github.com/guths/zpe/utils"
)

func InitializeRouter() (router *gin.Engine) {
	router = gin.Default()
	v1route := router.Group("/api/v1")
	v1route.Use(
		middleware.CORSMiddleware,
		middleware.AuthMiddleware,
		middleware.RoleMiddleware,
	)
	{
		auth := v1route.Group("/auth")
		{
			auth.POST("/login", v1.POSTLogin)
		}

		user := v1route.Group("/user")
		{
			//place middleware for auth and role to allow only lvl 2 -> modifiier
			//only get by email allow all roles to access
			user.POST("/", utils.ModifierOnly, utils.AuthOnly, v1.POSTUser)
			user.GET("/:email", utils.AuthOnly, v1.GETUser)
			user.DELETE("/:email", utils.ModifierOnly, utils.AuthOnly, v1.DELETEUser)
			user.PUT("/:email", utils.ModifierOnly, utils.AuthOnly, v1.PUTUser)
		}
	}
	return
}
