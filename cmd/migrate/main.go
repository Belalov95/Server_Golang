package main

import (
	"example/web-service-gin/database"
	"log"
	"os"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is required")
	}

	log.Printf("Applying migrations to %s", dbURL)
	if err := database.Migrate(dbURL); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migrations applied")
}
