package web

import (
	"encoding/json"
	"net/http"
	"strconv"

	"expensif/internal/service"
)

type APIHandler struct {
	svc *service.Service
}

func NewAPIHandler(svc *service.Service) *APIHandler {
	return &APIHandler{svc: svc}
}

type apiResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *APIHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if n := r.URL.Query().Get("n"); n != "" {
		if v, err := strconv.Atoi(n); err == nil && v > 0 {
			limit = v
		}
	}
	expenses, err := h.svc.ListExpenses(r.Context(), limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, apiResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Data: expenses})
}

func (h *APIHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Amount      float64 `json:"amount"`
		Category    string  `json:"category"`
		Description string  `json:"description,omitempty"`
		Date        string  `json:"date,omitempty"`
		Currency    string  `json:"currency,omitempty"`
		PaidBy      string  `json:"paid_by,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Error: "invalid JSON"})
		return
	}
	id, err := h.svc.CreateExpense(r.Context(), req.Amount, req.Category, req.Description, req.Date, req.Currency, req.PaidBy)
	if err != nil {
		if isValidationErr(err) {
			writeJSON(w, http.StatusBadRequest, apiResponse{Error: err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, apiResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, apiResponse{Data: map[string]int64{"id": id}})
}

func (h *APIHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	exp, err := h.svc.GetExpense(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, apiResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Data: exp})
}

func (h *APIHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	var req struct {
		Amount      float64 `json:"amount"`
		Category    string  `json:"category"`
		Description string  `json:"description,omitempty"`
		Date        string  `json:"date,omitempty"`
		Currency    string  `json:"currency,omitempty"`
		PaidBy      string  `json:"paid_by,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Error: "invalid JSON"})
		return
	}
	err := h.svc.UpdateExpense(r.Context(), id, req.Amount, req.Category, req.Description, req.Date, req.Currency, req.PaidBy)
	if err != nil {
		if isValidationErr(err) {
			writeJSON(w, http.StatusBadRequest, apiResponse{Error: err.Error()})
			return
		}
		writeJSON(w, http.StatusNotFound, apiResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Data: map[string]int64{"updated": id}})
}

func (h *APIHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err := h.svc.DeleteExpense(r.Context(), id); err != nil {
		writeJSON(w, http.StatusNotFound, apiResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Data: map[string]int64{"deleted": id}})
}

func (h *APIHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {
	cats, err := h.svc.ListCategories(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, apiResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Data: cats})
}

func (h *APIHandler) HandleSummary(w http.ResponseWriter, r *http.Request) {
	total, err := h.svc.TotalExpenses(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, apiResponse{Error: err.Error()})
		return
	}
	byCat, err := h.svc.SummaryByCategory(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, apiResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Data: map[string]interface{}{
		"total":      total,
		"byCategory": byCat,
	}})
}

func isValidationErr(err error) bool {
	return err == service.ErrInvalidAmount || err == service.ErrMissingCategory || err == service.ErrMissingDescription
}
