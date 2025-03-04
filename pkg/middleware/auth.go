package middleware

import (
	"golang_todo/pkg/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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
		verifiedToken, err := authService.ValidateToken(token, false)
		if err != nil || verifiedToken == nil || !verifiedToken.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "invalid token",
			})
			return
		}
		c.Next()
	}
}
