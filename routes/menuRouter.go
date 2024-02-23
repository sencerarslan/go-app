package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/sencerarslan/go-app/controllers"
	"github.com/sencerarslan/go-app/middleware"
)

func MenuRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/menu/:menu_id", controller.GetMenu())
}

func AuthMenuRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.POST("/menu/new", controller.AddMenu())
}
