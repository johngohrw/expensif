package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type APIResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

type CreateExpenseReq struct {
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description,omitempty"`
	Date        string  `json:"date,omitempty"`
	Currency    string  `json:"currency,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func handleCreate(w http.ResponseWriter, r *http.Request) {
	var req CreateExpenseReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Error: "invalid JSON"})
		return
	}
	if req.Amount <= 0 || req.Category == "" || strings.TrimSpace(req.Description) == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Error: "amount, description and category required"})
		return
	}

	id, err := addExpense(req.Amount, strings.TrimSpace(req.Category), strings.TrimSpace(req.Description), req.Date, req.Currency)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, APIResponse{Data: map[string]int64{"id": id}})
}

func handleList(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if n := r.URL.Query().Get("n"); n != "" {
		if v, err := strconv.Atoi(n); err == nil && v > 0 {
			limit = v
		}
	}

	expenses, err := listExpenses(limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Data: expenses})
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	exp, err := getExpense(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Data: exp})
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	var req CreateExpenseReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Error: "invalid JSON"})
		return
	}
	if req.Amount <= 0 || req.Category == "" || strings.TrimSpace(req.Description) == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Error: "amount, description and category required"})
		return
	}
	if err := updateExpense(id, req.Amount, strings.TrimSpace(req.Category), strings.TrimSpace(req.Description), req.Date, req.Currency); err != nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Data: map[string]int64{"updated": id}})
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err := deleteExpense(id); err != nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Data: map[string]int64{"deleted": id}})
}

func handleCategories(w http.ResponseWriter, r *http.Request) {
	cats, err := listCategories()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Data: cats})
}

func handleSummary(w http.ResponseWriter, r *http.Request) {
	total, err := totalExpenses()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Error: err.Error()})
		return
	}
	byCat, err := summaryByCategory()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Data: map[string]interface{}{
		"total":      total,
		"byCategory": byCat,
	}})
}


