package api

import (
	"context"
	"errors"
	"fmt"
	"golang_todo/pkg/config"
	logging "golang_todo/pkg/logger"
	"golang_todo/pkg/migrations"
	"sync"

	"github.com/uptrace/bun"
	"golang.org/x/time/rate"

	//"golang_todo/pkg/migrations"

	"golang_todo/pkg/routes"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var once sync.Once

func doOnce(f func()) {
	once.Do(f)
}

// StartServer starts the server in a go routine
func StartServer() {

	logging.InitLogger(config.Envs.UptraceDsn)

	// Initialize the database
	db, err := config.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	// Run migrations once
	doOnce(func() {
		log.Println("running migrations")
		migrations.Migrate(db)
	})
	// drop and recreate the DB
	//migrations.Drop(db)
	// Set up Gin server
	server := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	server.Use(rateLimiter)

	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowOriginFunc: func(origin string) bool {
			// Allow specific origins
			allowedOrigins := []string{
				"https://vue-todo-nine-henna.vercel.app",
				"http://localhost:3000",
			}
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					return true
				}
			}
			return false
		},
		AllowMethods:     []string{"GET", "DELETE", "POST", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Inject DB into router
	routes.SetupRoutes(server, db)
	server.RemoveExtraSlash = true
	var PORT string = config.Envs.ServerPort
	srv := &http.Server{
		Addr:         "0.0.0.0:" + PORT,
		Handler:      server,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Channel to listen for OS signals (Ctrl+C, SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Start server in a goroutine
	go func() {
		fmt.Printf("🚀 Server running on http://localhost:%s/api/users/test\n", PORT)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logging.Logger.Warn("Server failed", "error", err)
		}
		// defer db.Close()
	}()
	// go func() {
	// 	http.ListenAndServe(":6060", nil)
	// }()
	// Wait for termination signal
	<-quit
	fmt.Println("\n⚠️ Shutting down server...")

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down the server: %v", err)
	}

	// Close DB connection
	defer func(db *bun.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	fmt.Println("✅ Server gracefully stopped")
}

var limiter = rate.NewLimiter(1, 10)

func rateLimiter(c *gin.Context) {
	if !limiter.Allow() {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
		c.Abort()
		return
	}
	c.Next()
}
