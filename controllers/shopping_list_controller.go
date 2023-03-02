package controllers

import (
	"context"
	"net/http"
	"time"

	config "github.com/jsandis/shopping-list-go-rest-api/configs"
	helper "github.com/jsandis/shopping-list-go-rest-api/helpers"
	model "github.com/jsandis/shopping-list-go-rest-api/models"
	response "github.com/jsandis/shopping-list-go-rest-api/responses"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

var listCollection *mongo.Collection = config.GetCollection(config.ConnectDB(), config.EnvMongoDBCollectionNameList())
var listItemCollection *mongo.Collection = config.GetCollection(config.ConnectDB(), config.EnvMongoDBCollectionNameListItem())

func GetShoppingList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var list model.List
		var user_id string = c.GetString("user_id")
		var list_id string

		err := listCollection.FindOne(ctx, bson.M{"user_id": user_id}).Decode(&list)
		defer cancel()
		if err != nil {
			id := primitive.NewObjectID()
			newShoppingList := model.List{
				ID:      id,
				User_id: user_id}

			_, err := listCollection.InsertOne(ctx, newShoppingList)
			if err != nil {
				c.JSON(http.StatusInternalServerError, response.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			list_id = id.Hex()
		} else {
			list_id = list.ID.Hex()
		}

		var shoppingList []model.ListItem
		defer cancel()

		results, err := listItemCollection.Find(ctx, bson.M{"list_id": list_id})

		if err != nil {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		defer results.Close(ctx)
		for results.Next(ctx) {
			var shoppingListItem model.ListItem
			if err = results.Decode(&shoppingListItem); err != nil {
				c.JSON(http.StatusInternalServerError, response.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			shoppingList = append(shoppingList, shoppingListItem)
		}

		token, refreshToken, _ := helper.GenerateAllTokens(user_id)

		helper.UpdateAllTokens(token, refreshToken, user_id)

		c.JSON(http.StatusOK, shoppingList)
	}
}

func GetShoppingListItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		shoppingListItemId := c.Param("item_id")
		var shoppingListItem model.ListItem
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(shoppingListItemId)

		err := listItemCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&shoppingListItem)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, shoppingListItem)
	}
}

func AddListItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var list model.List
		var user_id string = c.GetString("user_id")
		var list_id string

		err := listCollection.FindOne(ctx, bson.M{"user_id": user_id}).Decode(&list)
		defer cancel()
		if err != nil {
			id := primitive.NewObjectID()
			newShoppingList := model.List{
				ID:      id,
				User_id: user_id}

			_, err := listCollection.InsertOne(ctx, newShoppingList)
			if err != nil {
				c.JSON(http.StatusInternalServerError, response.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			list_id = id.Hex()
		} else {
			list_id = list.ID.Hex()
		}

		var shoppingListItem model.ListItem
		if err := c.BindJSON(&shoppingListItem); err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if validationErr := validate.Struct(&shoppingListItem); validationErr != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newShoppingListItem := model.ListItem{
			ID:      primitive.NewObjectID(),
			Name:    shoppingListItem.Name,
			Status:  false,
			List_id: list_id}

		result, err := listItemCollection.InsertOne(ctx, newShoppingListItem)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		token, refreshToken, _ := helper.GenerateAllTokens(user_id)
		helper.UpdateAllTokens(token, refreshToken, user_id)

		c.JSON(http.StatusCreated, result)
	}
}

func EditShoppingListItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		shoppingListItemId := c.Param("item_id")
		var shoppingListItem model.ListItem
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(shoppingListItemId)

		if err := c.BindJSON(&shoppingListItem); err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if validationErr := validate.Struct(&shoppingListItem); validationErr != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{"name": shoppingListItem.Name, "status": shoppingListItem.Status}

		result, err := listItemCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		var updatedShoppingListItem model.ListItem
		if result.MatchedCount == 1 {
			err := listItemCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedShoppingListItem)
			if err != nil {
				c.JSON(http.StatusInternalServerError, response.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		var user_id string = c.GetString("user_id")
		token, refreshToken, _ := helper.GenerateAllTokens(user_id)
		helper.UpdateAllTokens(token, refreshToken, user_id)

		c.JSON(http.StatusOK, updatedShoppingListItem)
	}
}

func DeleteShoppingListItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		shoppingListItemId := c.Param("item_id")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(shoppingListItemId)

		result, err := listItemCollection.DeleteOne(ctx, bson.M{"_id": objId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				response.ErrorResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Shopping list item with specified ID not found!"}},
			)
			return
		}

		var user_id string = c.GetString("user_id")
		token, refreshToken, _ := helper.GenerateAllTokens(user_id)
		helper.UpdateAllTokens(token, refreshToken, user_id)

		c.JSON(http.StatusOK, "Shopping list item successfully deleted!")
	}
}
