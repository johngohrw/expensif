package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	initDB()
	defer closeDB()
	loadPreferences()
	parseTemplates()

	mux := http.NewServeMux()

	// HTML routes
	mux.HandleFunc("GET /{$}", handleListHTML)
	mux.HandleFunc("GET /daily", handleDailyHTML)
	mux.HandleFunc("GET /preferences", handlePreferencesHTML)
	mux.HandleFunc("POST /preferences", handleSavePreferences)
	mux.HandleFunc("GET /expenses/new", handleAddHTML)
	mux.HandleFunc("POST /expenses/new", handleCreateHTML)
	mux.HandleFunc("GET /expenses/edit/{id}", handleEditHTML)
	mux.HandleFunc("POST /expenses/edit/{id}", handleUpdateHTML)
	mux.HandleFunc("POST /expenses/delete/{id}", handleDeleteHTML)

	// JSON API routes
	mux.HandleFunc("GET /api/expenses", handleList)
	mux.HandleFunc("POST /api/expenses", handleCreate)
	mux.HandleFunc("GET /api/expenses/{id}", handleGet)
	mux.HandleFunc("PUT /api/expenses/{id}", handleUpdate)
	mux.HandleFunc("DELETE /api/expenses/{id}", handleDelete)
	mux.HandleFunc("GET /api/categories", handleCategories)
	mux.HandleFunc("GET /api/summary", handleSummary)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
