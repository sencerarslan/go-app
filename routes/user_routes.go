// routes/user_routes.go

package routes

import (
	"go-app/controllers"
	"go-app/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.Engine) {
	userRoutes := router.Group("/user")
	{
		userRoutes.POST("/delete", middleware.SessionMiddleware(), controllers.DeleteUser)
		userRoutes.GET("/", middleware.SessionMiddleware(), controllers.GetUsers)
	}

	accountRoutes := router.Group("/account")
	{
		accountRoutes.POST("/login", middleware.SessionMiddleware(), controllers.Login)
		userRoutes.POST("/register", middleware.SessionMiddleware(), controllers.Register)
	}
}
