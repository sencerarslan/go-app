// controllers/user_controller.go

package controllers

import (
	"context"
	"fmt"
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

func FetchUserByID(userID string) (models.User, error) {
	collection := utils.GetDB().Collection("users")
	user := models.User{}
	err := collection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
	return user, err
}
func FetchUserByEmail(email string) (models.User, error) {
	collection := utils.GetDB().Collection("users")
	user := models.User{}
	err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	return user, err
}

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func ComparePasswords(hashedPassword []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}

func SendError(c *gin.Context, err error, message string, statusCode int) {
	if err == mongo.ErrNoDocuments {
		message = "Resource not found"
	}
	apiErr := utils.NewAPIError(message, statusCode)
	apiErr.Send(c)
}

func UserDataWithoutPassword(user models.User) gin.H {
	return gin.H{
		"id":       user.ID,
		"name":     user.Name,
		"lastname": user.Lastname,
		"email":    user.Email,
	}
}

func GetUserByID(c *gin.Context) {
	userID := c.Param("id") // URL'den parametre olarak gelen kullanıcı kimliği
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		// Geçersiz ObjectID, hata döndür
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	collection := utils.GetDB().Collection("users")

	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Kullanıcı bulunamadı
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		// Veritabanı hatası
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		log.Fatal(err)
		return
	}

	resUserData := UserDataWithoutPassword(user)

	resp := utils.Response{Data: resUserData, Message: "Successfully", Status: http.StatusCreated}
	resp.Send(c)
}

func AllUsers(c *gin.Context) {
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

	user, err := FetchUserByEmail(loginData.Email)
	if err != nil {
		SendError(c, err, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	err = ComparePasswords([]byte(user.Password), loginData.Password)
	if err != nil {
		SendError(c, err, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	// Başarılı giriş, JWT oluştur
	token, err := utils.CreateToken(user.ID.Hex())
	if err != nil {
		SendError(c, err, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	// Şifre eşleşti, giriş başarılı

	resp := utils.Response{Data: gin.H{
		"token": token,
		"a":     "ea",
	}, Message: "Login successful", Status: http.StatusOK}
	resp.Send(c)
}

func Me(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is missing"})
		return
	}

	userID, err := utils.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized2"})
		return
	}
	fmt.Println("userID:", userID)

	// Kullanıcı verilerini veritabanından al
	user, err := FetchUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user data"})
		return
	}

	// Kullanıcı verilerini döndür
	resUserData := UserDataWithoutPassword(user)

	resp := utils.Response{Data: resUserData, Message: "Successfully", Status: http.StatusCreated}
	resp.Send(c)
}

func Register(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		apiErr := utils.NewAPIError(err.Error(), http.StatusBadRequest)
		apiErr.Send(c)
		return
	}

	_, err := FetchUserByEmail(user.Email)
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

	collection := utils.GetDB().Collection("users")

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		apiErr := utils.NewAPIError("Failed to hash password", http.StatusInternalServerError)
		apiErr.Send(c)
		log.Fatal(err)
		return
	}

	// Yeni bir ObjectID oluştur
	user.ID = primitive.NewObjectID()
	user.Password = string(hashedPassword)

	// Kullanıcıyı veritabanına ekle
	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		apiErr := utils.NewAPIError("User creation failed", http.StatusInternalServerError)
		apiErr.Send(c)
		log.Fatal(err)
		return
	}

	resUserData := UserDataWithoutPassword(user)

	resp := utils.Response{Data: resUserData, Message: "User created successfully", Status: http.StatusCreated}
	resp.Send(c)
}
