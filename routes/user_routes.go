// routes/user_routes.go

package routes

import (
	"go-app/controllers"
	"go-app/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.Engine) {

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/validateToken", controllers.Me)
	}

	userRoutes := router.Group("/user")
	{
		userRoutes.POST("/delete", middleware.AuthMiddleware(), controllers.DeleteUser)
		userRoutes.GET("/", middleware.AuthMiddleware(), controllers.AllUsers)
		userRoutes.GET("/:id", middleware.AuthMiddleware(), controllers.GetUserByID)
	}

	accountRoutes := router.Group("/account")
	{
		accountRoutes.POST("/login", controllers.Login)
		accountRoutes.POST("/register", controllers.Register)
	}

}
