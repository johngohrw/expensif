package web

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

type Server struct {
	mux *http.ServeMux
}

func NewServer(api *APIHandler, html *HTMLHandler) *Server {
	mux := http.NewServeMux()

	// HTML routes
	mux.HandleFunc("GET /{$}", html.HandleList)
	mux.HandleFunc("GET /daily", html.HandleDaily)
	mux.HandleFunc("GET /preferences", html.HandlePreferences)
	mux.HandleFunc("POST /preferences", html.HandleSavePreferences)
	mux.HandleFunc("GET /expenses/new", html.HandleAdd)
	mux.HandleFunc("POST /expenses/new", html.HandleCreate)
	mux.HandleFunc("GET /expenses/edit/{id}", html.HandleEdit)
	mux.HandleFunc("POST /expenses/edit/{id}", html.HandleUpdate)
	mux.HandleFunc("POST /expenses/delete/{id}", html.HandleDelete)

	// JSON API routes
	mux.HandleFunc("GET /api/expenses", api.HandleList)
	mux.HandleFunc("POST /api/expenses", api.HandleCreate)
	mux.HandleFunc("GET /api/expenses/{id}", api.HandleGet)
	mux.HandleFunc("PUT /api/expenses/{id}", api.HandleUpdate)
	mux.HandleFunc("DELETE /api/expenses/{id}", api.HandleDelete)
	mux.HandleFunc("GET /api/categories", api.HandleCategories)
	mux.HandleFunc("GET /api/summary", api.HandleSummary)

	return &Server{mux: mux}
}

func (s *Server) Run() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	slog.Info("server running", "url", fmt.Sprintf("http://localhost:%s", port))
	return http.ListenAndServe(":"+port, s.mux)
}
