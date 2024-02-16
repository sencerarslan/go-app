// main.go

package main

import (
	"go-app/routes"
	"go-app/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	utils.ConnectDB()

	router := gin.Default()

	// Global middleware
	router.Use(gin.Recovery())

	// Routes
	routes.SetupRoutes(router)

	// Start server
	router.Run(":8080")
}
