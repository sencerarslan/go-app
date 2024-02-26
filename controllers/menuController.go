package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sencerarslan/go-app/database"
	"github.com/sencerarslan/go-app/middleware"
	"github.com/sencerarslan/go-app/models"

	helper "github.com/sencerarslan/go-app/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")
var menuItemCollection *mongo.Collection = database.OpenCollection(database.Client, "menu-item")

func ShowMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var responseData models.Menu
		if err := c.BindJSON(&responseData); err != nil {
			errorResponse := helper.ErrorResponse(nil, err.Error())
			errorResponse.SendJSON(c.Writer, http.StatusBadRequest)
			return
		}

		var menuID = responseData.Menu_id
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var menu models.Menu
		err := menuCollection.FindOne(ctx, bson.M{"menu_id": menuID}).Decode(&menu)
		if err != nil {
			errorResponse := helper.ErrorResponse(nil, err.Error())
			errorResponse.SendJSON(c.Writer, http.StatusInternalServerError)
			cancel()
			return
		}

		var menuItems []models.MenuItem
		cursor, err := menuItemCollection.Find(ctx, bson.M{"menu_id": menuID})
		if err != nil {
			errorResponse := helper.ErrorResponse(nil, err.Error())
			errorResponse.SendJSON(c.Writer, http.StatusInternalServerError)
			cancel()
			return
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var menuItem models.MenuItem
			if err := cursor.Decode(&menuItem); err != nil {
				errorResponse := helper.ErrorResponse(nil, err.Error())
				errorResponse.SendJSON(c.Writer, http.StatusInternalServerError)
				cancel()
				return
			}
			menuItems = append(menuItems, menuItem)
		}

		if err := cursor.Err(); err != nil {
			errorResponse := helper.ErrorResponse(nil, err.Error())
			errorResponse.SendJSON(c.Writer, http.StatusInternalServerError)
			cancel()
			return
		}

		response := gin.H{
			"menu_id":    menu.Menu_id,
			"name":       menu.Name,
			"menu_items": menuItems,
		}

		successResponse := helper.SuccessResponse(response, "")
		successResponse.SendJSON(c.Writer, http.StatusOK)
		cancel()
	}
}
func AllGetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		middleware.Authenticate()(c)

		userID := c.GetString("uid")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var menus []models.Menu
		cursor, err := menuCollection.Find(ctx, bson.M{"user_id": userID})
		if err != nil {
			helper.ErrorResponse(c, err.Error())
			cancel()
			return
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var menu models.Menu
			if err := cursor.Decode(&menu); err != nil {
				helper.ErrorResponse(c, err.Error())
				cancel()
				return
			}

			var menuItems []models.MenuItem
			menuItemCursor, err := menuItemCollection.Find(ctx, bson.M{"menu_id": menu.Menu_id})
			if err != nil {
				helper.ErrorResponse(c, err.Error())
				cancel()
				return
			}
			defer menuItemCursor.Close(ctx)

			for menuItemCursor.Next(ctx) {
				var menuItem models.MenuItem
				if err := menuItemCursor.Decode(&menuItem); err != nil {
					helper.ErrorResponse(c, err.Error())
					cancel()
					return
				}
				menuItems = append(menuItems, menuItem)
			}

			if err := menuItemCursor.Err(); err != nil {
				helper.ErrorResponse(c, err.Error())
				cancel()
				return
			}

			menu.MenuItem = menuItems
			menus = append(menus, menu)
		}

		if err := cursor.Err(); err != nil {
			helper.ErrorResponse(c, err.Error())
			cancel()
			return
		}

		response := menus
		successResponse := helper.SuccessResponse(response, "")
		successResponse.SendJSON(c.Writer, http.StatusOK)
		cancel()
	}
}

func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var responseData models.Menu
		if err := c.BindJSON(&responseData); err != nil {
			response := helper.ErrorResponse(nil, err.Error())
			response.SendJSON(c.Writer, http.StatusBadRequest)
			return
		}

		var menuID = responseData.Menu_id
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var menu models.Menu
		err := menuCollection.FindOne(ctx, bson.M{"menu_id": menuID}).Decode(&menu)
		if err != nil {
			response := helper.ErrorResponse(nil, err.Error())
			response.SendJSON(c.Writer, http.StatusInternalServerError)
			return
		}

		var menuItems []models.MenuItem
		cursor, err := menuItemCollection.Find(ctx, bson.M{"menu_id": menuID})
		if err != nil {
			response := helper.ErrorResponse(nil, err.Error())
			response.SendJSON(c.Writer, http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var menuItem models.MenuItem
			if err := cursor.Decode(&menuItem); err != nil {
				response := helper.ErrorResponse(nil, err.Error())
				response.SendJSON(c.Writer, http.StatusInternalServerError)
				return
			}
			menuItems = append(menuItems, menuItem)
		}

		if err := cursor.Err(); err != nil {
			response := helper.ErrorResponse(nil, err.Error())
			response.SendJSON(c.Writer, http.StatusInternalServerError)
			return
		}

		response := gin.H{
			"ID":         menu.ID,
			"menu_id":    menu.Menu_id,
			"user_id":    menu.User_id,
			"name":       menu.Name,
			"created_at": menu.Created_at,
			"updated_at": menu.Updated_at,
			"menu_items": menuItems,
		}

		successResponse := helper.SuccessResponse(response, "")
		successResponse.SendJSON(c.Writer, http.StatusOK)
		cancel()
	}
}

func AddMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		middleware.Authenticate()(c)

		userID := c.GetString("uid")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
		if err != nil {
			response := helper.ErrorResponse(nil, "User not found")
			response.SendJSON(c.Writer, http.StatusInternalServerError)
			return
		}

		var menu models.Menu
		if err := c.BindJSON(&menu); err != nil {
			response := helper.ErrorResponse(nil, err.Error())
			response.SendJSON(c.Writer, http.StatusBadRequest)
			return
		}

		validationErr := validate.Struct(menu)
		if validationErr != nil {
			response := helper.ErrorResponse(nil, validationErr.Error())
			response.SendJSON(c.Writer, http.StatusBadRequest)
			return
		}

		menuID := primitive.NewObjectID().Hex()
		menu.ID = primitive.NewObjectID()
		menu.User_id = &userID
		menu.Menu_id = &menuID
		menu.Created_at = time.Now()
		menu.Updated_at = time.Now()

		_, err = menuCollection.InsertOne(ctx, menu)
		if err != nil {
			response := helper.ErrorResponse(nil, "Error while adding menu")
			response.SendJSON(c.Writer, http.StatusInternalServerError)
			return
		}

		responseData := gin.H{
			"message": "Menu added successfully",
			"menu":    menu,
		}
		response := helper.SuccessResponse(responseData, "")
		response.SendJSON(c.Writer, http.StatusOK)
	}
}

func AddMenuItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		middleware.Authenticate()(c)

		userID := c.GetString("uid")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
		if err != nil {
			response := helper.ErrorResponse(nil, "User not found")
			response.SendJSON(c.Writer, http.StatusInternalServerError)
			return
		}

		var menuItem models.MenuItem
		if err := c.BindJSON(&menuItem); err != nil {
			response := helper.ErrorResponse(nil, err.Error())
			response.SendJSON(c.Writer, http.StatusBadRequest)
			return
		}

		validationErr := validate.Struct(menuItem)
		if validationErr != nil {
			response := helper.ErrorResponse(nil, validationErr.Error())
			response.SendJSON(c.Writer, http.StatusBadRequest)
			return
		}

		if menuItem.ID != primitive.NilObjectID {
			filter := bson.M{"_id": menuItem.ID}
			update := bson.M{
				"$set": bson.M{
					"name":        menuItem.Name,
					"description": menuItem.Description,
					"price":       menuItem.Price,
					"imageurl":    menuItem.ImageURL,
					"updated_at":  time.Now(),
				},
			}

			updateResult, err := menuItemCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				response := helper.ErrorResponse(nil, "Error while updating menu item")
				response.SendJSON(c.Writer, http.StatusInternalServerError)
				return
			}

			if updateResult.ModifiedCount == 0 {
				response := helper.ErrorResponse(nil, "Menu item not found")
				response.SendJSON(c.Writer, http.StatusNotFound)
				return
			}

			responseData := gin.H{
				"message":   "Menu item updated successfully",
				"menu_item": menuItem,
			}
			response := helper.SuccessResponse(responseData, "")
			response.SendJSON(c.Writer, http.StatusOK)
			return
		}

		menuItem.ID = primitive.NewObjectID()
		menuItem.Created_at = time.Now()
		menuItem.Updated_at = time.Now()

		_, err = menuItemCollection.InsertOne(ctx, menuItem)
		if err != nil {
			response := helper.ErrorResponse(nil, "Error while adding menu item")
			response.SendJSON(c.Writer, http.StatusInternalServerError)
			return
		}

		responseData := gin.H{
			"message":   "Menu item added successfully",
			"menu_item": menuItem,
		}
		response := helper.SuccessResponse(responseData, "")
		response.SendJSON(c.Writer, http.StatusOK)
	}
}
func DeleteMenuItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		middleware.Authenticate()(c)

		var menuItem models.MenuItem
		if err := c.BindJSON(&menuItem); err != nil {
			response := helper.ErrorResponse(nil, err.Error())
			response.SendJSON(c.Writer, http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		filter := bson.M{"_id": menuItem.ID}

		deleteResult, err := menuItemCollection.DeleteOne(ctx, filter)
		if err != nil {
			response := helper.ErrorResponse(nil, "Error while deleting menu item")
			response.SendJSON(c.Writer, http.StatusInternalServerError)
			return
		}

		if deleteResult.DeletedCount == 0 {
			response := helper.ErrorResponse(nil, "Menu item not found")
			response.SendJSON(c.Writer, http.StatusNotFound)
			return
		}

		response := helper.SuccessResponse(nil, "Menu item deleted successfully")
		response.SendJSON(c.Writer, http.StatusOK)
	}
}
