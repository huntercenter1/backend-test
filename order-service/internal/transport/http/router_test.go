package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/huntercenter1/backend-test/order-service/internal/models"
	"github.com/huntercenter1/backend-test/order-service/internal/service"
)

type memSvc struct{
	o *models.Order
	it []models.OrderItem
	byUser []models.Order
}
func (m *memSvc) Create(_ context.Context, userID string, items []service.CreateItem) (*models.Order, []models.OrderItem, error) {
	m.o = &models.Order{ID:"o1", UserID:userID, Status:"pending", Total:100}
	m.it = []models.OrderItem{{ID:"i1", OrderID:"o1", ProductID:items[0].ProductID, Quantity:items[0].Quantity, Price:100}}
	return m.o, m.it, nil
}
func (m *memSvc) Get(_ context.Context, id string) (*models.Order, error) { return m.o, nil }
func (m *memSvc) Items(_ context.Context, id string) ([]models.OrderItem, error) { return m.it, nil }
func (m *memSvc) ByUser(_ context.Context, userID string) ([]models.Order, error) { if m.o!=nil { m.byUser = []models.Order{*m.o} }; return m.byUser, nil }
func (m *memSvc) UpdateStatus(_ context.Context, id, status string) (*models.Order, error) { m.o.Status = status; return m.o, nil }

func setupOrderRouter() (*gin.Engine, *Router, *memSvc) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	s := &memSvc{}
	rt := New(s)
	rt.Register(r)
	return r, rt, s
}

func TestOrderHandlers(t *testing.T){
	r, _, s := setupOrderRouter()

	// create
	body := []byte(`{"user_id":"u1","items":[{"product_id":"p1","quantity":1}]}`)
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(body))
	req.Header.Set("Content-Type","application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated { t.Fatalf("create code=%d", w.Code) }

	// get
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/orders/o1", nil))
	if w.Code != http.StatusOK { t.Fatalf("get code=%d", w.Code) }

	// items
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/orders/o1/items", nil))
	if w.Code != http.StatusOK { t.Fatalf("items code=%d", w.Code) }

	// by user
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/orders/user/u1", nil))
	if w.Code != http.StatusOK { t.Fatalf("user code=%d", w.Code) }
	var resp struct{ Orders []models.Order `json:"orders"` }
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp.Orders) == 0 { t.Fatalf("expected orders") }

	// update status
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPut, "/orders/o1/status", bytes.NewReader([]byte(`{"status":"paid"}`)))
	req.Header.Set("Content-Type","application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK { t.Fatalf("status code=%d", w.Code) }
	if s.o.Status != "paid" { t.Fatalf("status not updated") }
}

func TestOrderBadRequests(t *testing.T){
	r, _, _ := setupOrderRouter()

	// body inválido
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader([]byte(`{}`))))
	if w.Code != http.StatusBadRequest { t.Fatalf("bad body code=%d", w.Code) }

	// status inválido
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPut, "/orders/o1/status", bytes.NewReader([]byte(`{}`))))
	if w.Code != http.StatusBadRequest { t.Fatalf("bad status code=%d", w.Code) }
}
