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

		// Menu verisini al
		var menu models.Menu
		err := menuCollection.FindOne(ctx, bson.M{"menu_id": menuID}).Decode(&menu)
		if err != nil {
			errorResponse := helper.ErrorResponse(nil, err.Error())
			errorResponse.SendJSON(c.Writer, http.StatusInternalServerError)
			cancel()
			return
		}

		// MenuItem verilerini al
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

		// Menu ve MenuItem verilerini birleştirerek yanıt oluştur
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
		// Authenticate middleware'ini kullanarak kimlik doğrulaması yap
		middleware.Authenticate()(c)

		// Kullanıcı kimliğini al
		userID := c.GetString("uid")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		// Menü verilerini al
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

			// MenuItem verilerini al
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

		// Menü ve MenuItem verilerini birleştirerek yanıt oluştur
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
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var menuID = responseData.Menu_id
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		// Menu verisini al
		var menu models.Menu
		err := menuCollection.FindOne(ctx, bson.M{"menu_id": menuID}).Decode(&menu)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			cancel()
			return
		}

		// MenuItem verilerini al
		var menuItems []models.MenuItem
		cursor, err := menuItemCollection.Find(ctx, bson.M{"menu_id": menuID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			cancel()
			return
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var menuItem models.MenuItem
			if err := cursor.Decode(&menuItem); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				cancel()
				return
			}
			menuItems = append(menuItems, menuItem)
		}

		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			cancel()
			return
		}

		// Menu ve MenuItem verilerini birleştirerek yanıt oluştur
		response := gin.H{
			"menu": gin.H{
				"ID":         menu.ID,
				"menu_id":    menu.Menu_id,
				"user_id":    menu.User_id,
				"name":       menu.Name,
				"created_at": menu.Created_at,
				"updated_at": menu.Updated_at,
				"menu_items": menuItems,
			},
		}

		c.JSON(http.StatusOK, response)
		cancel()
	}
}

func AddMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authenticate middleware'ini kullanarak kimlik doğrulaması yap
		middleware.Authenticate()(c)

		// Kullanıcı kimliğini al
		userID := c.GetString("uid")

		// Kullanıcı kimliğiyle ilgili kullanıcıyı veritabanından al
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}

		// Gelen isteği Menu yapısına bind et
		var menu models.Menu
		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Gelen veriyi doğrula
		validationErr := validate.Struct(menu)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Menu verisini hazırla
		menuID := primitive.NewObjectID().Hex() // Rastgele bir ID oluştur
		menu.ID = primitive.NewObjectID()
		menu.User_id = &userID
		menu.Menu_id = &menuID
		menu.Created_at = time.Now()
		menu.Updated_at = time.Now()

		// Menu verisini veritabanına ekle
		_, err = menuCollection.InsertOne(ctx, menu)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while adding menu"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Menu added successfully", "menu": menu})
	}
}
func AddMenuItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authenticate middleware'ini kullanarak kimlik doğrulaması yap
		middleware.Authenticate()(c)

		// Kullanıcı kimliğini al
		userID := c.GetString("uid")

		// Kullanıcı kimliğiyle ilgili kullanıcıyı veritabanından al
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}

		// Gelen isteği MenuItem yapısına bind et
		var menuItem models.MenuItem
		if err := c.BindJSON(&menuItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Gelen veriyi doğrula
		validationErr := validate.Struct(menuItem)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		if menuItem.ID != primitive.NilObjectID {

			// Güncelleme işlemi yapılacak

			// Güncellenecek belgeyi belirlemek için bir filtre oluşturun
			filter := bson.M{"_id": menuItem.ID}

			// Yeni değerlerin atanacağı bir döküman oluşturun

			update := bson.M{
				"$set": bson.M{
					"name":        menuItem.Name,
					"description": menuItem.Description,
					"price":       menuItem.Price,
					"imageurl":    menuItem.ImageURL,
					"updated_at":  time.Now(),
					// Diğer alanları da güncelleyebilirsiniz
				},
			}

			// UpdateOne fonksiyonunu kullanarak güncelleme işlemini gerçekleştirin
			updateResult, err := menuItemCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating menu item"})
				return
			}

			// Güncellenen belgenin sayısını kontrol edin
			if updateResult.ModifiedCount == 0 {
				c.JSON(http.StatusNotFound, gin.H{"error": "Menu item not found"})
				return
			}

			// Başarılı güncelleme durumunda mesajı ve güncellenen menü öğesini döndürün
			c.JSON(http.StatusOK, gin.H{"message": "Menu item updated successfully", "menu_item": menuItem})
			return
		}

		// MenuItem verisini hazırla
		menuItem.ID = primitive.NewObjectID()
		menuItem.Created_at = time.Now()
		menuItem.Updated_at = time.Now()

		// MenuItem verisini veritabanına ekle
		_, err = menuItemCollection.InsertOne(ctx, menuItem)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while adding menu item"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Menu item added successfully", "menu_item": menuItem})
	}
}
func DeleteMenuItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authenticate middleware'ini kullanarak kimlik doğrulaması yap
		middleware.Authenticate()(c)

		// Gelen isteği MenuItem yapısına bind et
		var menuItem models.MenuItem
		if err := c.BindJSON(&menuItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Bir context oluşturun
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Silinecek belgeyi belirlemek için bir filtre oluşturun
		filter := bson.M{"_id": menuItem.ID}

		// DeleteOne fonksiyonunu kullanarak silme işlemini gerçekleştirin
		deleteResult, err := menuItemCollection.DeleteOne(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while deleting menu item"})
			return
		}

		// Silinen belgenin sayısını kontrol edin
		if deleteResult.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Menu item not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Menu item deleted successfully"})
	}
}
