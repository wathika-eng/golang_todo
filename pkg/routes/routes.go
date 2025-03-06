package routes

import (
	"golang_todo/pkg/config"
	"golang_todo/pkg/handlers"
	"golang_todo/pkg/repository"
	"golang_todo/pkg/services"
	notesservices "golang_todo/pkg/services/notes_services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

var secretKey = config.Envs.SECRET_KEY
var refreshKey = config.Envs.REFRESH_KEY
var resendApiKey = config.Envs.RESEND_API_KEY

func SetupRoutes(s *gin.Engine, db *bun.DB) {
	userRepo := repository.NewUserRepo(db)
	services := services.NewUserServices([]byte(secretKey), []byte(refreshKey), resendApiKey)
	userHandler := handlers.NewUserHandler(userRepo, services)
	api := s.Group("/api")
	users := api.Group("/users")
	{
		users.GET("/test", utest)
		users.POST("/signup", userHandler.SignUp)
		users.POST("/login", userHandler.Login)
		users.POST("/refresh", userHandler.RefreshAccess)
	}

	notesRepo := repository.NewNotesRepo(db)
	notesServices := notesservices.NewNotesServices()
	notesHandler := handlers.NewNotesHandler(notesRepo, notesServices)
	notes := api.Group("/notes")
	// notes.Use(middleware.AuthMiddleware(services))
	{
		notes.GET("/profile")
		notes.GET("/test", notesHandler.NotesTest)
		notes.POST("/create", notesHandler.NotesTest)
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
