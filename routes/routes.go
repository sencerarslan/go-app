// routes.go

package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	SetupUserRoutes(router)
}
