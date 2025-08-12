package service

import (
	"context"
	"testing"

	"github.com/huntercenter1/backend-test/order-service/internal/clients"
	"github.com/huntercenter1/backend-test/order-service/internal/models"
	"github.com/huntercenter1/backend-test/order-service/internal/repo"
)

type fakeUC struct{ ok bool; err error }
func (f fakeUC) Validate(ctx context.Context, id string)(bool,error){ return f.ok, f.err }

type fakePC struct{
	price float64; stock int; err error
}
func (f fakePC) Get(ctx context.Context, id string)(*clients.Product, error){
	if f.err != nil { return nil, f.err }
	return &clients.Product{ID:"p1", Price:f.price, Stock:f.stock}, nil
}
func (f fakePC) ApplyStockDelta(ctx context.Context, id string, delta int)(*clients.Product, error){
	return &clients.Product{ID:"p1", Price:f.price, Stock:f.stock - delta}, nil
}

type fakeRepo struct{}
func (f fakeRepo) CreateOrder(ctx context.Context, o *models.Order, items []models.OrderItem)(*models.Order, []models.OrderItem, error){ return o, items, nil }
func (f fakeRepo) GetOrder(ctx context.Context, id string)(*models.Order, error){ return nil, repo.ErrNotFound }
func (f fakeRepo) GetItems(ctx context.Context, id string)([]models.OrderItem, error){ return nil, nil }
func (f fakeRepo) ListByUser(ctx context.Context, userID string)([]models.Order, error){ return nil, nil }
func (f fakeRepo) UpdateStatus(ctx context.Context, id, status string)(*models.Order, error){ return nil, nil }

func TestCreateComputesTotal(t *testing.T){
	s := New(fakeRepo{}, fakeUC{ok:true}, fakePC{price:100, stock:10})
	o, items, err := s.Create(context.Background(), "u1", []CreateItem{{ProductID:"p1", Quantity:3}})
	if err != nil { t.Fatal(err) }
	if o.Total != 300 { t.Fatalf("want total=300 got %v", o.Total) }
	if len(items) != 1 || items[0].Price != 100 { t.Fatalf("items wrong") }
}

func TestCreateInvalidUser(t *testing.T){
	s := New(fakeRepo{}, fakeUC{ok:false}, fakePC{price:100, stock:10})
	if _, _, err := s.Create(context.Background(), "u1", []CreateItem{{ProductID:"p1", Quantity:1}}); err == nil {
		t.Fatalf("expected error")
	}
}

func TestCreateInsufficientStock(t *testing.T){
	s := New(fakeRepo{}, fakeUC{ok:true}, fakePC{price:100, stock:0})
	if _, _, err := s.Create(context.Background(), "u1", []CreateItem{{ProductID:"p1", Quantity:1}}); err == nil {
		t.Fatalf("expected error")
	}
}
