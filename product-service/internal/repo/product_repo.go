package repo

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/uptrace/bun"

	"github.com/huntercenter1/backend-test/product-service/internal/models"
)

var (
	ErrNotFound = errors.New("product not found")
	timeout     = 5 * time.Second
)

type ProductRepo interface {
	Create(ctx context.Context, p *models.Product) (*models.Product, error)
	GetByID(ctx context.Context, id string) (*models.Product, error)
	Update(ctx context.Context, p *models.Product) (*models.Product, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]models.Product, int, error)
	Search(ctx context.Context, q string, limit, offset int) ([]models.Product, int, error)
	UpdateStock(ctx context.Context, id string, delta int) (*models.Product, error)
}

type productRepo struct{ db *bun.DB }

func New(db *bun.DB) ProductRepo { return &productRepo{db: db} }

func (r *productRepo) Create(ctx context.Context, p *models.Product) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout); defer cancel()
	_, err := r.db.NewInsert().Model(p).Exec(ctx)
	return p, err
}

func (r *productRepo) GetByID(ctx context.Context, id string) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout); defer cancel()
	var p models.Product
	if err := r.db.NewSelect().Model(&p).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, ErrNotFound
	}
	return &p, nil
}

func (r *productRepo) Update(ctx context.Context, p *models.Product) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout); defer cancel()
	p.UpdatedAt = time.Now()
	_, err := r.db.NewUpdate().Model(p).Column("name", "description", "price", "stock", "updated_at").WherePK().Exec(ctx)
	return p, err
}

func (r *productRepo) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, timeout); defer cancel()
	res, err := r.db.NewDelete().Model((*models.Product)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil { return err }
	if n, _ := res.RowsAffected(); n == 0 { return ErrNotFound }
	return nil
}

func (r *productRepo) List(ctx context.Context, limit, offset int) ([]models.Product, int, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout); defer cancel()
	var items []models.Product
	q := r.db.NewSelect().Model(&items).Order("created_at DESC").Limit(limit).Offset(offset)
	if err := q.Scan(ctx); err != nil { return nil, 0, err }
	total, err := r.db.NewSelect().Model((*models.Product)(nil)).Count(ctx)
	return items, total, err
}

func (r *productRepo) Search(ctx context.Context, qstr string, limit, offset int) ([]models.Product, int, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout); defer cancel()
	qstr = strings.TrimSpace(qstr)
	var items []models.Product
	q := r.db.NewSelect().Model(&items).
		Where("name ILIKE ? OR description ILIKE ?", "%"+qstr+"%", "%"+qstr+"%").
		Order("created_at DESC").Limit(limit).Offset(offset)
	if err := q.Scan(ctx); err != nil { return nil, 0, err }
	total, err := r.db.NewSelect().Model((*models.Product)(nil)).
		Where("name ILIKE ? OR description ILIKE ?", "%"+qstr+"%", "%"+qstr+"%").Count(ctx)
	return items, total, err
}

func (r *productRepo) UpdateStock(ctx context.Context, id string, delta int) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout); defer cancel()
	p, err := r.GetByID(ctx, id)
	if err != nil { return nil, err }
	p.Stock += delta
	if p.Stock < 0 { p.Stock = 0 }
	_, err = r.db.NewUpdate().Model(p).Column("stock", "updated_at").WherePK().Exec(ctx)
	return p, err
}
