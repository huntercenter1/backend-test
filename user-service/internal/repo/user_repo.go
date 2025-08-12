package repo

import (
	"context"
	"errors"
	"time"

	"github.com/uptrace/bun"

	"github.com/huntercenter1/backend-test/user-service/internal/models"
)

var (
	ErrNotFound      = errors.New("user not found")
	ErrDuplicate     = errors.New("user duplicate username/email")
	defaultTimeout   = 5 * time.Second
)

type UserRepo interface {
	Create(ctx context.Context, u *models.User) (*models.User, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Update(ctx context.Context, u *models.User) (*models.User, error)
	Delete(ctx context.Context, id string) error
}

type userRepo struct {
	db *bun.DB
}

func NewUserRepo(db *bun.DB) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, u *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	_, err := r.db.NewInsert().Model(u).Exec(ctx)
	if err != nil {
		// índice unique en username/email → tratamos como duplicado
		return nil, ErrDuplicate
	}
	return u, nil
}

func (r *userRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var u models.User
	err := r.db.NewSelect().Model(&u).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, ErrNotFound
	}
	return &u, nil
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var u models.User
	err := r.db.NewSelect().Model(&u).Where("username = ?", username).Scan(ctx)
	if err != nil {
		return nil, ErrNotFound
	}
	return &u, nil
}

func (r *userRepo) Update(ctx context.Context, u *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	u.UpdatedAt = time.Now()
	_, err := r.db.NewUpdate().Model(u).
		Column("username", "email", "password_hash", "updated_at").
		WherePK().
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepo) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	res, err := r.db.NewDelete().Model((*models.User)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return ErrNotFound
	}
	return nil
}
