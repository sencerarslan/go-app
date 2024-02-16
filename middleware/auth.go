// middleware/auth.go

package middleware

import (
	"go-app/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Middleware işlevi: Token doğrulama
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is missing"})
			c.Abort()
			return
		}
		userID, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		// Kullanıcı kimliğini context'e ekle
		c.Set("userID", userID)
		c.Next()
	}
}
