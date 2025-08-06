package main

import (
	"context"
	"log"
	"os"

	"page-state-saver/pagestate"
)

func main() {
	ctx := context.Background()

	dbConnStr := os.Getenv("DATABASE_URL")

	if dbConnStr == "" {
		dbConnStr = "postgres://postgres:postgres@localhost:5432/page-state-saver?sslmode=disable"
	}

	repo, err := pagestate.NewRepository(ctx, dbConnStr)

	if err != nil {
		log.Fatal("Failed to create repository:", err)
	}

	server := pagestate.NewServer("8080", repo)

	log.Println("Starting page state server on port 8080")

	if err := server.Start(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
