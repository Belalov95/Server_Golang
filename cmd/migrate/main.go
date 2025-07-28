package main

import (
	"example/web-service-gin/database"
	"example/web-service-gin/internal/config"
	"log"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := database.Migrate(cfg.GetConnStr()); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration successful")
}
