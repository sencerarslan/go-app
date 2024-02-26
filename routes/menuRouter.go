package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/sencerarslan/go-app/controllers"
	"github.com/sencerarslan/go-app/middleware"
)

func AuthMenuRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())

	incomingRoutes.GET("/all-menu", controller.AllGetMenu())
	incomingRoutes.POST("/menu", controller.GetMenu())
	incomingRoutes.POST("/menu/add", controller.AddMenu())
	incomingRoutes.POST("/menu/item/add", controller.AddMenuItem())
	incomingRoutes.POST("/menu/item/delete", controller.DeleteMenuItem())

}
