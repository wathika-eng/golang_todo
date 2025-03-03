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

	"github.com/gin-gonic/gin"
)

func StartServer() {
	logging.InitLogger()

	// Initialize the database
	db := config.InitDB()

	// Run migrations
	migrations.Migrate(db)
	//migrations.Drop(db)
	// Set up Gin server
	server := gin.Default()

	// Inject DB into routes
	routes.SetupRoutes(server, db)
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
		fmt.Printf("🚀 Server running on http://localhost:%s\n", PORT)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Logger.Warn("Server failed: ", err)
		}
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
