package routes

import (
	controller "github.com/jsandis/shopping-list-go-rest-api/controllers"
	middleware "github.com/jsandis/shopping-list-go-rest-api/middleware"

	"github.com/gin-gonic/gin"
)

func ShoppingListRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authorize())

	route := incomingRoutes.Group("/api/v2/shoppinglist")

	route.GET("/", controller.GetShoppingList())

	route.GET("/item/:item_id", controller.GetShoppingListItem())

	route.POST("/item", controller.AddListItem())

	route.PATCH("/item/:item_id", controller.EditShoppingListItem())

	route.DELETE("/item/:item_id", controller.DeleteShoppingListItem())
}
