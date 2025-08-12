package service

import (
	"context"
	"errors"

	"github.com/huntercenter1/backend-test/order-service/internal/clients"
	"github.com/huntercenter1/backend-test/order-service/internal/models"
	"github.com/huntercenter1/backend-test/order-service/internal/repo"
)

type CreateItem struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
}

type Service interface {
	Create(ctx context.Context, userID string, items []CreateItem) (*models.Order, []models.OrderItem, error)
	Get(ctx context.Context, id string) (*models.Order, error)
	Items(ctx context.Context, id string) ([]models.OrderItem, error)
	ByUser(ctx context.Context, userID string) ([]models.Order, error)
	UpdateStatus(ctx context.Context, id, status string) (*models.Order, error)
}

type service struct {
	repo repo.Repo
	uc   clients.UserClient
	pc   clients.ProductClient
}

func New(r repo.Repo, uc clients.UserClient, pc clients.ProductClient) Service {
	return &service{repo: r, uc: uc, pc: pc}
}

func (s *service) Create(ctx context.Context, userID string, items []CreateItem) (*models.Order, []models.OrderItem, error) {
	if userID == "" || len(items) == 0 { return nil, nil, errors.New("invalid payload") }

	// 1) validar usuario
	ok, err := s.uc.Validate(ctx, userID)
	if err != nil || !ok { return nil, nil, errors.New("invalid user") }

	// 2) verificar stock y total
	var orderItems []models.OrderItem
	var total float64
	for _, it := range items {
		if it.Quantity <= 0 { return nil, nil, errors.New("quantity must be > 0") }
		p, err := s.pc.Get(ctx, it.ProductID)
		if err != nil { return nil, nil, err }
		if p.Stock < it.Quantity { return nil, nil, errors.New("insufficient stock") }
		line := models.OrderItem{
			ProductID: it.ProductID,
			Quantity:  it.Quantity,
			Price:     p.Price,
		}
		total += p.Price * float64(it.Quantity)
		orderItems = append(orderItems, line)
	}

	// 3) crear orden
	o := &models.Order{ UserID: userID, Status: "pending", Total: total }
	o, orderItems, err = s.repo.CreateOrder(ctx, o, orderItems)
	if err != nil { return nil, nil, err }

	// 4) descontar stock (delta negativo) por cada item
	for _, it := range items {
		if _, err := s.pc.ApplyStockDelta(ctx, it.ProductID, -it.Quantity); err != nil {
			// Nota: en un caso real, aquí haríamos compensación/cola
			return o, orderItems, nil // devolvemos la orden creada aunque haya fallos de stock update
		}
	}

	return o, orderItems, nil
}

func (s *service) Get(ctx context.Context, id string) (*models.Order, error) {
	return s.repo.GetOrder(ctx, id)
}

func (s *service) Items(ctx context.Context, id string) ([]models.OrderItem, error) {
	return s.repo.GetItems(ctx, id)
}

func (s *service) ByUser(ctx context.Context, userID string) ([]models.Order, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *service) UpdateStatus(ctx context.Context, id, status string) (*models.Order, error) {
	if status == "" { return nil, errors.New("status required") }
	return s.repo.UpdateStatus(ctx, id, status)
}
