package routes

import (
	"golang_todo/pkg/config"
	"golang_todo/pkg/handlers"
	"golang_todo/pkg/middleware"
	"golang_todo/pkg/repository"
	"golang_todo/pkg/response"
	"golang_todo/pkg/services"
	notesservices "golang_todo/pkg/services/notes_services"
	redisservices "golang_todo/pkg/services/redis"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

var secretKey = config.Envs.SecretKey
var refreshKey = config.Envs.RefreshKey
var resendApiKey = config.Envs.ResendApiKey
var redisURL = config.Envs.RedisUrl

func SetupRoutes(s *gin.Engine, db *bun.DB) {
	userRepo := repository.NewUserRepo(db)
	userServices := services.NewUserServices([]byte(secretKey), []byte(refreshKey), resendApiKey)
	newResponse := response.NewResponse()
	userHandler := handlers.NewUserHandler(userRepo, userServices, newResponse)
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
	notesHandler := handlers.NewNotesHandler(notesRepo, notesServices, redisServices, newResponse)
	notes := api.Group("/notes")
	notes.Use(middleware.AuthMiddleware(userServices, db, redisServices))
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
	userProfile.Use(middleware.AuthMiddleware(userServices, db, redisServices))
	{
		userProfile.GET("/", userHandler.UserProfile)
	}
}
