package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jehufrayle/grimoire/internal/database"
	"github.com/jehufrayle/grimoire/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	// Initialize the godotenv package to load environment variables from a .env file
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("No .env file found, using environment variables directly")
	}

	// Connect to the database
	database.Connect()

	defer database.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start the server
	server.StartServer(ctx, ":8080")
}
