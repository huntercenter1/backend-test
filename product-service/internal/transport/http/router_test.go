package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/huntercenter1/backend-test/product-service/internal/models"
	"github.com/huntercenter1/backend-test/product-service/internal/repo"
)

type memRepo struct {
	data map[string]*models.Product
}

func newMemRepo() *memRepo { return &memRepo{data: map[string]*models.Product{}} }

func (m *memRepo) Create(ctx context.Context, p *models.Product) (*models.Product, error) {
	if p.ID == "" { p.ID = uuid.NewString() }
	now := time.Now().UTC()
	if p.CreatedAt.IsZero() { p.CreatedAt = now }
	p.UpdatedAt = now
	cp := *p
	m.data[p.ID] = &cp
	return &cp, nil
}

func (m *memRepo) GetByID(ctx context.Context, id string) (*models.Product, error) {
	if p, ok := m.data[id]; ok { cp := *p; return &cp, nil }
	return nil, repo.ErrNotFound
}

func (m *memRepo) Update(ctx context.Context, p *models.Product) (*models.Product, error) {
	if _, ok := m.data[p.ID]; !ok { return nil, repo.ErrNotFound }
	p.UpdatedAt = time.Now().UTC()
	cp := *p; m.data[p.ID] = &cp
	return &cp, nil
}

func (m *memRepo) Delete(ctx context.Context, id string) error {
	if _, ok := m.data[id]; !ok { return repo.ErrNotFound }
	delete(m.data, id); return nil
}

func (m *memRepo) List(ctx context.Context, limit, offset int) ([]models.Product, int, error) {
	ids := make([]string, 0, len(m.data))
	for id := range m.data { ids = append(ids, id) }
	sort.Strings(ids)
	total := len(ids)
	end := offset + limit
	if end > total { end = total }
	ids = ids[offset:end]
	out := make([]models.Product, 0, len(ids))
	for _, id := range ids { out = append(out, *m.data[id]) }
	return out, total, nil
}

func (m *memRepo) Search(ctx context.Context, q string, limit, offset int) ([]models.Product, int, error) {
	return m.List(ctx, limit, offset)
}

func (m *memRepo) UpdateStock(ctx context.Context, id string, delta int) (*models.Product, error) {
	p, ok := m.data[id]; if !ok { return nil, repo.ErrNotFound }
	p.Stock += delta; if p.Stock < 0 { p.Stock = 0 }
	p.UpdatedAt = time.Now().UTC()
	cp := *p; m.data[id] = &cp
	return &cp, nil
}

func setupRouter(t *testing.T) (*gin.Engine, *Router, *memRepo) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	rt := New(nil) // db nil en test
	mem := newMemRepo()
	rt.repo = mem // inyectamos fake repo
	rt.Register(r)
	return r, rt, mem
}

func TestProductCRUDAndSearch(t *testing.T) {
	r, _, mem := setupRouter(t)

	// create
	body := []byte(`{"name":"Headset","description":"Wireless","price":99.9,"stock":15}`)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
	req.Header.Set("Content-Type","application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated { t.Fatalf("create code=%d", w.Code) }
	var created models.Product
	_ = json.Unmarshal(w.Body.Bytes(), &created)

	// get
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/products/"+created.ID, nil))
	if w.Code != http.StatusOK { t.Fatalf("get code=%d", w.Code) }

	// list
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/products?limit=10&offset=0", nil))
	if w.Code != http.StatusOK { t.Fatalf("list code=%d", w.Code) }

	// search
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/products/search?q=head", nil))
	if w.Code != http.StatusOK { t.Fatalf("search code=%d", w.Code) }

	// update
	created.Price = 79.5
	buf, _ := json.Marshal(created)
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPut, "/products/"+created.ID, bytes.NewReader(buf))
	req.Header.Set("Content-Type","application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK { t.Fatalf("update code=%d", w.Code) }

	// stock delta
	w = httptest.NewRecorder()
	delta := []byte(`{"delta": -3}`)
	req = httptest.NewRequest(http.MethodPut, "/products/"+created.ID+"/stock", bytes.NewReader(delta))
	req.Header.Set("Content-Type","application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK { t.Fatalf("stock code=%d", w.Code) }
	if mem.data[created.ID].Stock != 12 { t.Fatalf("want stock 12 got %d", mem.data[created.ID].Stock) }

	// delete
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/products/"+created.ID, nil))
	if w.Code != http.StatusNoContent { t.Fatalf("delete code=%d", w.Code) }
}

func TestPaginationParams(t *testing.T) {
	r, _, _ := setupRouter(t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/products?limit=200&offset=-1", nil) // normaliza
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK { t.Fatalf("pag code=%d", w.Code) }

	// parámetros edge
	q := "/products?limit=" + strconv.Itoa(0) + "&offset=" + strconv.Itoa(0)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, q, nil))
	if w.Code != http.StatusOK { t.Fatalf("pag2 code=%d", w.Code) }
}
