package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/sencerarslan/go-app/controllers"
	"github.com/sencerarslan/go-app/middleware"
)

func AuthMenuRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.POST("/menu/show", controller.ShowMenu())

	menu := incomingRoutes.Group("/menu")
	menu.POST("", middleware.Authenticate(), controller.GetMenu())
	menu.POST("/add", middleware.Authenticate(), controller.AddUpdateMenu())
	menu.POST("/delete", middleware.Authenticate(), controller.DeleteMenu())

	menuGroup := incomingRoutes.Group("/menu/group")
	menuGroup.POST("", middleware.Authenticate(), controller.GetGroup())
	menuGroup.POST("/add", middleware.Authenticate(), controller.AddUpdateGroup())
	menuGroup.POST("/delete", middleware.Authenticate(), controller.DeleteGroup())

	menuGroupItem := incomingRoutes.Group("/menu/group/item")
	menuGroupItem.POST("", middleware.Authenticate(), controller.GetItem())
	menuGroupItem.POST("/add", middleware.Authenticate(), controller.AddUpdateItem())
	menuGroupItem.POST("/delete", middleware.Authenticate(), controller.DeleteItem())
}
