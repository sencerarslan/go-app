package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/sencerarslan/go-app/controllers"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("signup", controller.Signup())
	incomingRoutes.POST("login", controller.Login())
}
