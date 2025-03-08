package middleware

import (
	"context"
	"golang_todo/pkg/services"
	redisservices "golang_todo/pkg/services/redis"
	"golang_todo/pkg/types"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

func AuthMiddleware(authService services.Auth, db *bun.DB, redisService redisservices.Redis) gin.HandlerFunc {
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
		if redisService.IsTokenBlacklisted(token) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "token revoked",
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
		userIDStr, exists := claims["user_id"].(string)
		userEmail := claims["sub"].(string)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "user_id not found in token",
			})
			return
		}
		userID, _ := uuid.Parse(userIDStr)
		// if err != nil {
		// 	return nil, fmt.Errorf("invalid user_id format: %v", err)
		// }
		// bug when db is dropped, token is still valid
		var user types.User
		err = db.NewSelect().Model(&user).Where("id = ?", userID).Scan(context.Background())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": true, "message": "user not found"})
			return
		}
		expUnix := int64(claims["exp"].(float64))
		expirationTime := time.Until(time.Unix(expUnix, 0))
		// Store userID in context
		c.Set("user_id", userID)
		c.Set("user_email", userEmail)
		c.Set("exp_time", expirationTime)
		c.Set("user_token", token)
		c.Next()
	}
}
