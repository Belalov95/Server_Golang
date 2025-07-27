package main

import (
	"context"
	"example/web-service-gin/internal/app"
	"log/slog"
	"os"
)

func main() {
	if err := app.Run(context.Background()); err != nil {
		slog.Error("app run", slog.Any("error", err))
		os.Exit(1)
	}
}
