package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/sencerarslan/go-app/controllers"
	"github.com/sencerarslan/go-app/middleware"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/users", middleware.Authenticate(), controller.GetUsers())
	incomingRoutes.GET("/users/:user_id", middleware.Authenticate(), controller.GetUser())
}
