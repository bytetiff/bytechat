package middleware

import (
	"net/http"
	"strings"

	"github.com/bytetiff/bytechat/internal/auth"
	"github.com/bytetiff/bytechat/internal/repository"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(tokenRepo repository.RefreshTokenRepository, jwtKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
			c.Abort()
			return
		}

		tokenString := authHeader[7:]
		claims, err := auth.ParseJWT(tokenString, jwtKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Добавляем user_id в контекст
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
