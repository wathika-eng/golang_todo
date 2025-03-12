package routes

import (
	"golang_todo/pkg/config"
	"golang_todo/pkg/handlers"
	"golang_todo/pkg/middleware"
	"golang_todo/pkg/repository"
	"golang_todo/pkg/services"
	notesservices "golang_todo/pkg/services/notes_services"
	redisservices "golang_todo/pkg/services/redis"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

var secretKey = config.Envs.SECRET_KEY
var refreshKey = config.Envs.REFRESH_KEY
var resendApiKey = config.Envs.RESEND_API_KEY
var redisURL = config.Envs.REDIS_URL

func SetupRoutes(s *gin.Engine, db *bun.DB) {
	userRepo := repository.NewUserRepo(db)
	services := services.NewUserServices([]byte(secretKey), []byte(refreshKey), resendApiKey)
	userHandler := handlers.NewUserHandler(userRepo, services)
	api := s.Group("/api")
	users := api.Group("/users")
	{
		users.GET("/test", userHandler.UserTest)
		users.POST("/signup", userHandler.SignUp)
		users.POST("/login", userHandler.Login)
		users.POST("/refresh", userHandler.RefreshAccess)
	}

	notesRepo := repository.NewNotesRepo(db)
	notesServices := notesservices.NewNotesServices()
	redisServices := redisservices.NewRedisClient(redisURL)
	notesHandler := handlers.NewNotesHandler(notesRepo, notesServices, redisServices)
	notes := api.Group("/notes")
	notes.Use(middleware.AuthMiddleware(services, db, redisServices))
	{
		notes.GET("/test", notesHandler.NotesTest)
		notes.POST("/create", notesHandler.CreateNotes)
		notes.GET("/", notesHandler.GetNotes)
		notes.GET("/:id", notesHandler.GetNoteByID)
		notes.PATCH("/:id", notesHandler.UpdateNotes)
		notes.DELETE("/:id", notesHandler.DeleteNotes)
		notes.POST("/logout", notesHandler.Logout)
		notes.GET("/recent/del", notesHandler.RecentDeletions)
	}
	userProfile := api.Group("/profile")
	userProfile.Use(middleware.AuthMiddleware(services, db, redisServices))
	{
		userProfile.GET("/", userHandler.UserProfile)
	}
}
