// middleware/session_middleware.go

package middleware

import (
	"github.com/gin-gonic/gin"
)

func SessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Oturum kontrolü yapılabilir
		// Örneğin, burada oturum bilgilerini kontrol edebilir ve giriş yapmamış kullanıcıları yönlendirebilirsiniz

		c.Next()
	}
}
