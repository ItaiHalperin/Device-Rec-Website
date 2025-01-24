package main

import (
	"context"
	"deviceRecommendationWebsite/internal/handlers"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("/Users/itaihalperin/GolandProjects/deviceRecommendationWebsite/internal/web")))

	mux.HandleFunc("/submit-filters", handlers.SendMessageHandler)

	server := &http.Server{
		Addr:         ":80",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		fmt.Println("Server starting on http://localhost:80")
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Printf("Server error: %v", err)
	case sig := <-shutdown:
		log.Printf("Start shutdown, signal: %v", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)
			if err := server.Close(); err != nil {
				log.Printf("Forcing server close: %v", err)
			}
		}
	}
}
