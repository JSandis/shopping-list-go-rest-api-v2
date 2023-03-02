package routes

import (
	controller "github.com/jsandis/shopping-list-go-rest-api/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	route := incomingRoutes.Group("/api/v2/user")

	route.POST("/signup", controller.CreateUser())

	route.POST("/login", controller.Login())
}
