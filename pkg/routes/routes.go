package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

func SetupRoutes(s *gin.Engine, db *bun.DB) {
	api := s.Group("/api")
	users := api.Group("/users")
	{
		users.GET("/test", utest)
	}
	notes := api.Group("/notes")
	{
		notes.GET("/test", ntest)
	}
}

func utest(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "API is up and running",
		"time":    time.Now().Local(),
	})
}

func ntest(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "notes API is up and running",
		"time":    time.Now().Local(),
	})
}
