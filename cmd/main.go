package main

import (
	"example/web-service-gin/internal/app"
	"log/slog"
	"os"
)

func main() {
	if err := app.Run(); err != nil {
		slog.Error("Problem starting the program: %v\n", slog.Any("error", err))
		os.Exit(1)
	}
}
