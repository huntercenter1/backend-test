package repo

import (
	"context"
	"errors"
	"time"

	"github.com/uptrace/bun"

	"github.com/huntercenter1/backend-test/order-service/internal/models"
)

var (
	ErrNotFound = errors.New("not found")
	timeout     = 5 * time.Second
)

type Repo interface {
	CreateOrder(ctx context.Context, o *models.Order, items []models.OrderItem) (*models.Order, []models.OrderItem, error)
	GetOrder(ctx context.Context, id string) (*models.Order, error)
	GetItems(ctx context.Context, orderID string) ([]models.OrderItem, error)
	ListByUser(ctx context.Context, userID string) ([]models.Order, error)
	UpdateStatus(ctx context.Context, id, status string) (*models.Order, error)
}

type repo struct{ db *bun.DB }

func New(db *bun.DB) Repo { return &repo{db: db} }

func (r *repo) CreateOrder(ctx context.Context, o *models.Order, items []models.OrderItem) (*models.Order, []models.OrderItem, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout); defer cancel()

	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(o).Exec(ctx); err != nil {
			return err
		}
		for i := range items {
			items[i].OrderID = o.ID
		}
		if _, err := tx.NewInsert().Model(&items).Exec(ctx); err != nil {
			return err
		}
		return nil
	})
	return o, items, err
}

func (r *repo) GetOrder(ctx context.Context, id string) (*models.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout); defer cancel()
	var o models.Order
	if err := r.db.NewSelect().Model(&o).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, ErrNotFound
	}
	return &o, nil
}

func (r *repo) GetItems(ctx context.Context, orderID string) ([]models.OrderItem, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout); defer cancel()
	var items []models.OrderItem
	if err := r.db.NewSelect().Model(&items).Where("order_id = ?", orderID).Scan(ctx); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *repo) ListByUser(ctx context.Context, userID string) ([]models.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout); defer cancel()
	var orders []models.Order
	if err := r.db.NewSelect().Model(&orders).Where("user_id = ?", userID).Order("created_at DESC").Scan(ctx); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *repo) UpdateStatus(ctx context.Context, id, status string) (*models.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout); defer cancel()
	o, err := r.GetOrder(ctx, id)
	if err != nil { return nil, err }
	o.Status = status
	o.UpdatedAt = time.Now()
	if _, err := r.db.NewUpdate().Model(o).Column("status", "updated_at").WherePK().Exec(ctx); err != nil {
		return nil, err
	}
	return o, nil
}
