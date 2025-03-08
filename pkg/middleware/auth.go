package middleware

import (
	"golang_todo/pkg/services"
	redisservices "golang_todo/pkg/services/redis"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(authService services.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "missing auth header",
			})
			c.Redirect(http.StatusPermanentRedirect, "/login")
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if redisservices.NewRedisClient().IsTokenBlacklisted(token) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "token revoked",
			})
			return
		}

		verifiedToken, err := authService.ValidateToken(token, false)
		if err != nil || verifiedToken == nil || !verifiedToken.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "invalid token",
			})
			return
		}
		// extract token
		claims, ok := verifiedToken.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "invalid token claims",
			})
			return
		}
		userID, exists := claims["user_id"].(float64)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "user_id not found in token",
			})
			return
		}
		// userID := uint(userIDFloat)

		// Store userID in context
		c.Set("user_id", uint(userID))
		c.Next()
	}
}
