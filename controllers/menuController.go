package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sencerarslan/go-app/database"
	"github.com/sencerarslan/go-app/middleware"
	"github.com/sencerarslan/go-app/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		menuId := c.Param("menu_id")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var menu models.Menu
		err := menuCollection.FindOne(ctx, bson.M{"menu_id": menuId}).Decode(&menu)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, menu)
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
