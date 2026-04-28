package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"expensif/internal/db"
	"expensif/internal/repository"
	"expensif/internal/service"
	"expensif/internal/web"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))

	database, err := db.New()
	if err != nil {
		log.Fatalf("db init failed: %v", err)
	}
	defer database.Close()

	repo := repository.NewSQLite(database)
	svc := service.New(repo)

	// Background rate refresh: fetch immediately, then every 6 hours
	go func() {
		ctx := context.Background()
		for {
			if err := svc.RefreshRates(ctx); err != nil {
				slog.Error("rate refresh failed", "error", err)
			} else {
				slog.Info("exchange rates refreshed")
			}
			time.Sleep(6 * time.Hour)
		}
	}()

	renderer, err := web.NewRenderer("templates")
	if err != nil {
		log.Fatalf("template init failed: %v", err)
	}

	apiHandler := web.NewAPIHandler(svc)
	htmlHandler := web.NewHTMLHandler(svc, renderer)
	server := web.NewServer(apiHandler, htmlHandler)

	if err := server.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
