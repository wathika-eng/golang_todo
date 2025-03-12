package api

import (
	"context"
	"fmt"
	"golang_todo/pkg/config"
	logging "golang_todo/pkg/logger"
	"golang_todo/pkg/migrations"
	"golang_todo/pkg/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// starts the server in a go routine
func StartServer() {
	logging.InitLogger(config.Envs.UPTRACE_DSN)

	// Initialize the database
	db := config.InitDB()

	// Run migrations once
	migrations.Migrate(db)
	// drop and recreate the DB
	//migrations.Drop(db)
	// Set up Gin server
	server := gin.Default()
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://162.245.188.225:3000/*"},
		AllowMethods:     []string{"GET", "DELETE", "POST", "PATCH"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Inject DB into routesr)
	routes.SetupRoutes(server, db)
	server.RemoveExtraSlash = true
	var PORT string = config.Envs.SERVER_PORT
	srv := &http.Server{
		Addr:         ":" + PORT,
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
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Logger.Warn("Server failed", "error", err)
		}
		defer db.Close()
	}()

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
	db.Close()

	fmt.Println("✅ Server gracefully stopped")
}
