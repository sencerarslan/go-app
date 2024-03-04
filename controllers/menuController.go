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
var menuGroupCollection *mongo.Collection = database.OpenCollection(database.Client, "menu-group")
var menuItemCollection *mongo.Collection = database.OpenCollection(database.Client, "menu-item")

func useContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	return ctx, cancel
}

func getMenuByID(menuID primitive.ObjectID) (models.Menu, error) {
	ctx, cancel := useContext()
	defer cancel()

	var data models.Menu
	err := menuCollection.FindOne(ctx, bson.M{"_id": menuID}).Decode(&data)
	if err != nil {
		return models.Menu{}, err
	}

	return data, nil
}
func getMenuGroupByAll(menuID string) ([]models.MenuGroup, error) {
	ctx, cancel := useContext()
	defer cancel()

	cursor, err := menuGroupCollection.Find(ctx, bson.M{"menuid": menuID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items = make([]models.MenuGroup, 0)
	for cursor.Next(ctx) {
		var item models.MenuGroup
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func getMenuItemByID(groupID string) ([]models.MenuItem, error) {
	ctx, cancel := useContext()
	defer cancel()

	cursor, err := menuItemCollection.Find(ctx, bson.M{"groupid": groupID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	items := make([]models.MenuItem, 0)
	for cursor.Next(ctx) {
		var item models.MenuItem
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func ShowMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var responseData models.Menu
		if err := c.BindJSON(&responseData); err != nil {
			response := helper.ErrorResponse(nil, err.Error())
			response.SendJSON(c.Writer, http.StatusBadRequest)
			return
		}

		menuID := responseData.ID

		menu, err := getMenuByID(menuID)
		if err != nil {
			response := helper.ErrorResponse(nil, err.Error())
			response.SendJSON(c.Writer, http.StatusInternalServerError)
			return
		}

		menuGroups, err := getMenuGroupByAll(menuID.Hex())
		if err != nil {
			response := helper.ErrorResponse(nil, err.Error())
			response.SendJSON(c.Writer, http.StatusInternalServerError)
			return
		}

		menuItems := make([][]models.MenuItem, len(menuGroups))
		for i, group := range menuGroups {
			groupID := group.ID.Hex()
			menuItems[i], err = getMenuItemByID(groupID)
			if err != nil {
				response := helper.ErrorResponse(nil, err.Error())
				response.SendJSON(c.Writer, http.StatusInternalServerError)
				return
			}
		}

		menuGroupsArray := make([]gin.H, len(menuGroups))
		for i, group := range menuGroups {
			menuItemsArray := make([]gin.H, len(menuItems[i]))
			for j, item := range menuItems[i] {
				menuItem := gin.H{
					"id":          item.ID,
					"group_id":    item.GroupID,
					"name":        item.Name,
					"price":       item.Price,
					"description": item.Description,
					"image_url":   item.ImageURL,
				}
				menuItemsArray[j] = menuItem
			}

			menuGroup := gin.H{
				"id":    group.ID,
				"name":  group.Name,
				"items": menuItemsArray,
			}
			menuGroupsArray[i] = menuGroup
		}

		response := gin.H{
			"id":          menu.ID,
			"name":        menu.Name,
			"logo":        menu.Logo,
			"banner":      menu.Banner,
			"menu_groups": menuGroupsArray,
		}

		successResponse := helper.SuccessResponse(response, "")
		successResponse.SendJSON(c.Writer, http.StatusOK)
	}
}

func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := useContext()
		defer cancel()

		userID := c.GetString("uid")

		cursor, err := menuCollection.Find(ctx, bson.M{"userid": userID})
		if err != nil {
			response := helper.ErrorResponse(nil, err.Error())
			response.SendJSON(c.Writer, http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)

		var items []gin.H
		for cursor.Next(ctx) {
			var item models.Menu
			if err := cursor.Decode(&item); err != nil {
				response := helper.ErrorResponse(nil, err.Error())
				response.SendJSON(c.Writer, http.StatusInternalServerError)
				return
			}
			responseItem := gin.H{
				"id":     item.ID,
				"name":   item.Name,
				"logo":   item.Logo,
				"banner": item.Banner,
			}
			items = append(items, responseItem)
		}

		successResponse := helper.SuccessResponse(items, "")
		successResponse.SendJSON(c.Writer, http.StatusOK)
	}
}
func AddUpdateMenu() gin.HandlerFunc {
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

		if menu.ID != primitive.NilObjectID {
			filter := bson.M{"_id": menu.ID}
			update := bson.M{
				"$set": bson.M{
					"name":       menu.Name,
					"logo":       menu.Logo,
					"banner":     menu.Banner,
					"updated_at": time.Now(),
				},
			}

			updateResult, err := menuCollection.UpdateOne(ctx, filter, update)
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
				"menu_item": menu,
			}
			response := helper.SuccessResponse(responseData, "")
			response.SendJSON(c.Writer, http.StatusOK)
			return
		}

		menu.ID = primitive.NewObjectID()
		menu.UserID = &userID
		menu.MenuGroup = make([]models.MenuGroup, 0)
		menu.CreatedAt = time.Now()
		menu.UpdatedAt = time.Now()

		_, err = menuCollection.InsertOne(ctx, menu)
		if err != nil {
			response := helper.ErrorResponse(nil, "Error while adding menu")
			response.SendJSON(c.Writer, http.StatusInternalServerError)
			return
		}

		response := helper.SuccessResponse(menu, "Menu added successfully")
		response.SendJSON(c.Writer, http.StatusOK)
	}
}
func DeleteMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		middleware.Authenticate()(c)

		var menu models.Menu
		if err := c.BindJSON(&menu); err != nil {
			response := helper.ErrorResponse(nil, err.Error())
			response.SendJSON(c.Writer, http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		filter := bson.M{"_id": menu.ID}

		deleteResult, err := menuCollection.DeleteOne(ctx, filter)
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

func GetGroup() gin.HandlerFunc {
	return func(c *gin.Context) {

		var responseData models.MenuGroup
		if err := c.BindJSON(&responseData); err != nil {
			response := helper.ErrorResponse(nil, err.Error())
			response.SendJSON(c.Writer, http.StatusBadRequest)
			return
		}

		menuID := responseData.ID.Hex()

		data, err := getMenuGroupByAll(menuID)
		if err != nil {
			response := helper.ErrorResponse(nil, err.Error())
			response.SendJSON(c.Writer, http.StatusInternalServerError)
			return
		}

		successResponse := helper.SuccessResponse(data, "")
		successResponse.SendJSON(c.Writer, http.StatusOK)
	}
}
func AddUpdateGroup() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := useContext()
		defer cancel()

		userID := c.GetString("uid")

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
		if err != nil {
			response := helper.ErrorResponse(nil, "User not found")
			response.SendJSON(c.Writer, http.StatusInternalServerError)
			return
		}

		var menuGroup models.MenuGroup
		if err := c.BindJSON(&menuGroup); err != nil {
			response := helper.ErrorResponse(nil, err.Error())
			response.SendJSON(c.Writer, http.StatusBadRequest)
			return
		}

		validationErr := validate.Struct(menuGroup)
		if validationErr != nil {
			response := helper.ErrorResponse(nil, validationErr.Error())
			response.SendJSON(c.Writer, http.StatusBadRequest)
			return
		}

		if menuGroup.ID != primitive.NilObjectID {
			filter := bson.M{"_id": menuGroup.ID}
			update := bson.M{
				"$set": bson.M{
					"name":       menuGroup.Name,
					"updated_at": time.Now(),
				},
			}

			updateResult, err := menuGroupCollection.UpdateOne(ctx, filter, update)
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
				"menu_item": menuGroup,
			}
			response := helper.SuccessResponse(responseData, "")
			response.SendJSON(c.Writer, http.StatusOK)
			return
		}

		menuGroup.ID = primitive.NewObjectID()
		menuGroup.MenuItem = make([]models.MenuItem, 0)
		menuGroup.CreatedAt = time.Now()
		menuGroup.UpdatedAt = time.Now()

		_, err = menuGroupCollection.InsertOne(ctx, menuGroup)
		if err != nil {
			response := helper.ErrorResponse(nil, "Error while adding menu item")
			response.SendJSON(c.Writer, http.StatusInternalServerError)
			return
		}

		response := helper.SuccessResponse(menuGroup, "")
		response.SendJSON(c.Writer, http.StatusOK)
	}
}
func DeleteGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		middleware.Authenticate()(c)

		var menuGroup models.MenuGroup
		if err := c.BindJSON(&menuGroup); err != nil {
			response := helper.ErrorResponse(nil, err.Error())
			response.SendJSON(c.Writer, http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		filter := bson.M{"_id": menuGroup.ID}

		deleteResult, err := menuGroupCollection.DeleteOne(ctx, filter)
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

func GetItem() gin.HandlerFunc {
	return func(c *gin.Context) {

		var responseData models.MenuItem
		if err := c.BindJSON(&responseData); err != nil {
			response := helper.ErrorResponse(nil, err.Error())
			response.SendJSON(c.Writer, http.StatusBadRequest)
			return
		}

		menuGroupID := responseData.ID.Hex()

		data, err := getMenuItemByID(menuGroupID)
		if err != nil {
			response := helper.ErrorResponse(nil, err.Error())
			response.SendJSON(c.Writer, http.StatusInternalServerError)
			return
		}

		successResponse := helper.SuccessResponse(data, "")
		successResponse.SendJSON(c.Writer, http.StatusOK)
	}
}
func AddUpdateItem() gin.HandlerFunc {
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
		menuItem.CreatedAt = time.Now()
		menuItem.UpdatedAt = time.Now()

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
func DeleteItem() gin.HandlerFunc {
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
