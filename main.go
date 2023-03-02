package main

import (
	config "github.com/jsandis/shopping-list-go-rest-api/configs"
	routes "github.com/jsandis/shopping-list-go-rest-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	port := config.EnvPort()

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.ShoppingListRoutes(router)

	router.Run(":" + port)
}
