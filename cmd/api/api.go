package api

import (
	"context"
	"fmt"
	"golang_todo/pkg/config"
	logging "golang_todo/pkg/logger"
	"golang_todo/pkg/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	PORT       = config.Envs.SERVER_PORT
	SECRET_KEY = config.Envs.SECRET_KEY
)

func StartServer() {
	logging.InitLogger()
	// initializer the database
	db := config.InitDB()
	defer db.Close()

	//
	server := gin.Default()
	// Routes
	routes.SetupRoutes(server, db)

	// Create HTTP Server with timeout settings
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
			server.Run(srv.Addr)
		}
	}()

	// Wait for termination signal
	<-quit
	fmt.Println("\n⚠️ Shutting down server...")

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("error shutdowning the server: %v", err)
	}

	fmt.Println("✅ Server gracefully stopped")
}
