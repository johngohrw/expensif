package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"expensif/internal/repository"
	"expensif/internal/service"
)

func newTestAPIHandler() (*APIHandler, *mockRepo) {
	repo := newMockRepo()
	svc := service.New(repo)
	return NewAPIHandler(svc), repo
}

func parseAPIResponse(t *testing.T, rr *httptest.ResponseRecorder) apiResponse {
	t.Helper()
	var resp apiResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}
	return resp
}

// ---------- LIST ----------

func TestAPIList_Empty(t *testing.T) {
	h, _ := newTestAPIHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/expenses", nil)
	rr := httptest.NewRecorder()

	h.HandleList(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	resp := parseAPIResponse(t, rr)
	data, ok := resp.Data.([]interface{})
	if !ok {
		t.Fatalf("expected array data, got %T", resp.Data)
	}
	if len(data) != 0 {
		t.Fatalf("expected 0 expenses, got %d", len(data))
	}
}

func TestAPIList_WithData(t *testing.T) {
	h, repo := newTestAPIHandler()
	repo.seed()

	req := httptest.NewRequest(http.MethodGet, "/api/expenses", nil)
	rr := httptest.NewRecorder()
	h.HandleList(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	resp := parseAPIResponse(t, rr)
	data, ok := resp.Data.([]interface{})
	if !ok {
		t.Fatalf("expected array data, got %T", resp.Data)
	}
	if len(data) != 3 {
		t.Fatalf("expected 3 expenses, got %d", len(data))
	}
}

func TestAPIList_LimitQuery(t *testing.T) {
	h, repo := newTestAPIHandler()
	repo.seed()

	req := httptest.NewRequest(http.MethodGet, "/api/expenses?n=2", nil)
	rr := httptest.NewRecorder()
	h.HandleList(rr, req)

	resp := parseAPIResponse(t, rr)
	data := resp.Data.([]interface{})
	if len(data) != 2 {
		t.Fatalf("expected 2 expenses with limit=2, got %d", len(data))
	}
}

// ---------- CREATE ----------

func TestAPICreate_Success(t *testing.T) {
	h, _ := newTestAPIHandler()
	body := map[string]interface{}{
		"amount":      25.0,
		"category":    "food",
		"description": "dinner",
		"date":        "2024-06-15",
		"currency":    "MYR",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/expenses", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.HandleCreate(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rr.Code, rr.Body.String())
	}
	resp := parseAPIResponse(t, rr)
	m, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map data, got %T", resp.Data)
	}
	if m["id"] == nil {
		t.Fatalf("expected id in response, got nil")
	}
}

func TestAPICreate_MissingCategory(t *testing.T) {
	h, _ := newTestAPIHandler()
	body := map[string]interface{}{"amount": 10.0, "description": "test"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/expenses", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.HandleCreate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
	resp := parseAPIResponse(t, rr)
	if resp.Error == "" {
		t.Fatal("expected error message")
	}
	if !strings.Contains(resp.Error, "category") {
		t.Fatalf("expected category validation error, got: %s", resp.Error)
	}
}

func TestAPICreate_InvalidAmount(t *testing.T) {
	h, _ := newTestAPIHandler()
	body := map[string]interface{}{"amount": -5, "category": "food", "description": "test"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/expenses", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.HandleCreate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
	resp := parseAPIResponse(t, rr)
	if !strings.Contains(resp.Error, "amount") {
		t.Fatalf("expected amount validation error, got: %s", resp.Error)
	}
}

func TestAPICreate_MissingDescription(t *testing.T) {
	h, _ := newTestAPIHandler()
	body := map[string]interface{}{"amount": 10.0, "category": "food"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/expenses", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.HandleCreate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestAPICreate_InvalidJSON(t *testing.T) {
	h, _ := newTestAPIHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/expenses", strings.NewReader("not json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.HandleCreate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
	resp := parseAPIResponse(t, rr)
	if resp.Error == "" {
		t.Fatal("expected error message")
	}
}

// ---------- GET ----------

func TestAPIGet_Success(t *testing.T) {
	h, repo := newTestAPIHandler()
	repo.seed()

	req := httptest.NewRequest(http.MethodGet, "/api/expenses/1", nil)
	req.SetPathValue("id", "1")
	rr := httptest.NewRecorder()

	h.HandleGet(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	resp := parseAPIResponse(t, rr)
	m, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", resp.Data)
	}
	if m["ID"] != float64(1) {
		t.Fatalf("expected id 1, got %v", m["ID"])
	}
}

func TestAPIGet_NotFound(t *testing.T) {
	h, _ := newTestAPIHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/expenses/999", nil)
	req.SetPathValue("id", "999")
	rr := httptest.NewRecorder()

	h.HandleGet(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

// ---------- UPDATE ----------

func TestAPIUpdate_Success(t *testing.T) {
	h, repo := newTestAPIHandler()
	repo.seed()

	body := map[string]interface{}{
		"amount":      99.0,
		"category":    "updated",
		"description": "updated desc",
		"date":        "2024-12-01",
		"currency":    "JPY",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/api/expenses/1", bytes.NewReader(b))
	req.SetPathValue("id", "1")
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.HandleUpdate(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	resp := parseAPIResponse(t, rr)
	m := resp.Data.(map[string]interface{})
	if m["updated"] != float64(1) {
		t.Fatalf("expected updated id 1, got %v", m["updated"])
	}

	// verify
	exp, _ := repo.GetExpense(nil, 1)
	if exp.Category != "updated" || exp.Amount != 99.0 {
		t.Fatalf("update not persisted: %+v", exp)
	}
}

func TestAPIUpdate_NotFound(t *testing.T) {
	h, _ := newTestAPIHandler()
	body := map[string]interface{}{"amount": 10, "category": "x", "description": "y"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/api/expenses/999", bytes.NewReader(b))
	req.SetPathValue("id", "999")
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.HandleUpdate(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestAPIUpdate_InvalidJSON(t *testing.T) {
	h, repo := newTestAPIHandler()
	repo.seed()
	req := httptest.NewRequest(http.MethodPut, "/api/expenses/1", strings.NewReader("bad"))
	req.SetPathValue("id", "1")
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.HandleUpdate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestAPIUpdate_ValidationError(t *testing.T) {
	h, repo := newTestAPIHandler()
	repo.seed()
	body := map[string]interface{}{"amount": -1, "category": "x", "description": "y"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/api/expenses/1", bytes.NewReader(b))
	req.SetPathValue("id", "1")
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.HandleUpdate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

// ---------- DELETE ----------

func TestAPIDelete_Success(t *testing.T) {
	h, repo := newTestAPIHandler()
	repo.seed()

	req := httptest.NewRequest(http.MethodDelete, "/api/expenses/1", nil)
	req.SetPathValue("id", "1")
	rr := httptest.NewRecorder()

	h.HandleDelete(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	resp := parseAPIResponse(t, rr)
	m := resp.Data.(map[string]interface{})
	if m["deleted"] != float64(1) {
		t.Fatalf("expected deleted 1, got %v", m["deleted"])
	}

	_, err := repo.GetExpense(nil, 1)
	if err == nil {
		t.Fatal("expected expense to be deleted")
	}
}

func TestAPIDelete_NotFound(t *testing.T) {
	h, _ := newTestAPIHandler()
	req := httptest.NewRequest(http.MethodDelete, "/api/expenses/999", nil)
	req.SetPathValue("id", "999")
	rr := httptest.NewRecorder()

	h.HandleDelete(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

// ---------- CATEGORIES ----------

func TestAPICategories(t *testing.T) {
	h, repo := newTestAPIHandler()
	repo.seed()

	req := httptest.NewRequest(http.MethodGet, "/api/categories", nil)
	rr := httptest.NewRecorder()

	h.HandleCategories(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	resp := parseAPIResponse(t, rr)
	cats, ok := resp.Data.([]interface{})
	if !ok {
		t.Fatalf("expected array, got %T", resp.Data)
	}
	if len(cats) == 0 {
		t.Fatal("expected some categories")
	}
}

// ---------- SUMMARY ----------

func TestAPISummary(t *testing.T) {
	h, repo := newTestAPIHandler()
	repo.seed()

	req := httptest.NewRequest(http.MethodGet, "/api/summary", nil)
	rr := httptest.NewRecorder()

	h.HandleSummary(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	resp := parseAPIResponse(t, rr)
	m, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", resp.Data)
	}
	if m["total"] == nil {
		t.Fatal("expected total")
	}
	if m["byCategory"] == nil {
		t.Fatal("expected byCategory")
	}
}

// Ensure the interface is satisfied so mockRepo compiles against Repository.
var _ repository.Repository = (*mockRepo)(nil)
