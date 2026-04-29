package web

import (
	"context"
	"net/http"
	"testing"
	"time"

	"expensif/internal/repository"
	"expensif/internal/service"
)

func TestServerShutdownBeforeRun(t *testing.T) {
	repo := newMockRepo()
	svc := service.New(repository.Repos{
		Expenses:    repo,
		Users:       repo,
		Preferences: repo,
		Rates:       repo,
	})
	api := NewAPIHandler(svc)
	html := NewHTMLHandler(svc, nil)

	server := NewServer(api, html, "9999")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestServerRunAndShutdown(t *testing.T) {
	repo := newMockRepo()
	svc := service.New(repository.Repos{
		Expenses:    repo,
		Users:       repo,
		Preferences: repo,
		Rates:       repo,
	})
	api := NewAPIHandler(svc)

	renderer, err := NewRenderer("../../templates", false, nil)
	if err != nil {
		t.Fatalf("renderer init failed: %v", err)
	}
	html := NewHTMLHandler(svc, renderer)
	server := NewServer(api, html, "0") // port 0 = random available port

	go func() {
		if err := server.Run(); err != nil && err != http.ErrServerClosed {
			t.Errorf("run: %v", err)
		}
	}()

	// Give the server time to start listening
	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		t.Fatalf("shutdown: %v", err)
	}
}
