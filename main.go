package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/Mallbrusss/BackEntryMiddle/internal/server"
	"github.com/joho/godotenv"
)

func init() {

}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	log.Println("Starting server...")
	serv := server.NewServer()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

	go func() {
		if err := serv.Run(); err != nil {
			log.Fatalf("Server encountered an error: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := serv.Shutdown(ctx); err != nil {
		log.Fatalf("Graceful shutdown failed: %v", err)
	}

	log.Println("Server stopped")
}
