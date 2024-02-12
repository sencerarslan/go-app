// controllers/user_controller.go

package controllers

import (
	"context"
	"log"
	"net/http"

	"go-app/models"
	"go-app/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers(c *gin.Context) {
	var users []models.User

	collection := utils.GetDB().Collection("users")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}

func DeleteUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Kullanıcı kimliğini al
	objectID := user.ID

	// MongoDB bağlantısını al
	collection := utils.GetDB().Collection("users")

	// Kullanıcıyı sil
	filter := bson.M{"_id": objectID}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete user"})
		log.Fatal(err)
		return
	}

	// Kullanıcı silindi mi kontrol et
	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
func Login(c *gin.Context) {
	var loginData models.User
	if err := c.BindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Kullanıcıyı veritabanından bul
	collection := utils.GetDB().Collection("users")

	user := models.User{}
	err := collection.FindOne(context.Background(), bson.M{"email": loginData.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Kullanıcı bulunamadı
			apiErr := utils.NewAPIError("Invalid email or password", http.StatusUnauthorized)
			apiErr.Send(c)
			return
		}
		// Veritabanı hatası
		apiErr := utils.NewAPIError("Database error", http.StatusInternalServerError)
		apiErr.Send(c)
		log.Fatal(err)
		return
	}

	// Kullanıcının girdiği şifreyi bcrypt ile hashleme
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		// Şifre eşleşmedi
		apiErr := utils.NewAPIError("Invalid email or password", http.StatusUnauthorized)
		apiErr.Send(c)
		return
	}

	// Şifre eşleşti, giriş başarılı
	// Kullanıcı verisinden şifre bilgisini kaldır
	user.Password = ""

	resp := utils.Response{Data: gin.H{
		"id":       user.ID,
		"name":     user.Name,
		"lastname": user.Lastname,
		"email":    user.Email,
	}, Message: "Login successful", Status: http.StatusOK}
	resp.Send(c)
}

func Register(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		apiErr := utils.NewAPIError(err.Error(), http.StatusBadRequest)
		apiErr.Send(c)
		return
	}

	// E-posta adresi benzersiz olmalıdır, bu yüzden veritabanında var olup olmadığını kontrol ediyoruz
	collection := utils.GetDB().Collection("users")
	existingUser := models.User{}
	err := collection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		// E-posta adresi zaten kullanımda, hata döndür
		apiErr := utils.NewAPIError("Email address is already in use", http.StatusBadRequest)
		apiErr.Send(c)
		return
	} else if err != nil && err != mongo.ErrNoDocuments {
		// Veritabanı hatası
		apiErr := utils.NewAPIError("Database error", http.StatusInternalServerError)
		apiErr.Send(c)
		log.Fatal(err)
		return
	}

	// Şifreyi bcrypt ile hashle
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		apiErr := utils.NewAPIError("Failed to hash password", http.StatusInternalServerError)
		apiErr.Send(c)
		log.Fatal(err)
		return
	}

	// Yeni bir ObjectID oluştur
	id := primitive.NewObjectID()
	user.ID = id
	user.Password = string(hashedPassword)

	// Kullanıcıyı veritabanına ekle
	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		apiErr := utils.NewAPIError("User creation failed", http.StatusInternalServerError)
		apiErr.Send(c)
		log.Fatal(err)
		return
	}

	user.Password = ""
	resUserData := gin.H{
		"id":       user.ID,
		"name":     user.Name,
		"lastname": user.Lastname,
		"email":    user.Email,
	}

	resp := utils.Response{Data: resUserData, Message: "User created successfully", Status: http.StatusCreated}
	resp.Send(c)
}
