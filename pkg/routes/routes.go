package routes

import (
	"golang_todo/pkg/handlers"
	"golang_todo/pkg/repository"
	"golang_todo/pkg/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

func SetupRoutes(s *gin.Engine, db *bun.DB) {
	userRepo := repository.NewUserRepo(db)
	services := services.NewUserServices()
	userHandler := handlers.NewUserHandler(userRepo, services)
	api := s.Group("/api")
	users := api.Group("/users")
	{
		users.GET("/test", utest)
		users.POST("/signup", userHandler.SignUp)
		users.POST("/login", userHandler.Login)
	}
	notes := api.Group("/notes")
	{
		notes.GET("/test", ntest)
		notes.POST("/", handlers.CreateNotes)
		notes.GET("/", handlers.GetNotes)
		notes.GET("/:id", handlers.GetNoteByID)
		notes.PATCH("/:id", handlers.UpdateNotes)
		notes.DELETE("/:id", handlers.DeleteNotes)
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
