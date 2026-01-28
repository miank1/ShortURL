package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"url-shortnener/internal/handler"
	"url-shortnener/internal/service"
	"url-shortnener/internal/store"

	"github.com/gin-gonic/gin"
)

func main() {
	connStr := "postgresql://neondb_owner:npg_I32xocTdRqwO@ep-floral-cherry-ah84bezz-pooler.c-3.us-east-1.aws.neon.tech/tinyurl?sslmode=require&channel_binding=require"

	pgStore, err := store.NewPostgresStore(connStr)
	if err != nil {
		log.Fatal("failed to connect to postgres:", err)
	}

	defer pgStore.Close()

	urlService := service.NewURLService(pgStore)
	urlHandler := handler.NewURLHandler(urlService)

	r := gin.Default()

	r.GET("/health", urlHandler.Health)
	r.POST("/shorten", urlHandler.Shorten)
	r.GET("/:code", urlHandler.Redirect)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// starting the server
	go func() {
		log.Println("Server running on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutdown signal received")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Println("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")

}
