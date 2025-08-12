package repo

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/huntercenter1/backend-test/product-service/internal/models"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

func testDB(t *testing.T) *bun.DB {
	sqlDB, err := sql.Open(sqliteshim.DriverName(), "file:memdb1?mode=memory&cache=shared")
	if err != nil { t.Fatal(err) }
	db := bun.NewDB(sqlDB, sqlitedialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(false)))

	// Tabla compatible con SQLite
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS products(
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			price REAL NOT NULL,
			stock INTEGER NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
	`)
	if err != nil { t.Fatal(err) }
	return db
}

func TestUpdateStock(t *testing.T){
	db := testDB(t)
	r := New(db)

	now := time.Now().UTC()
	p := &models.Product{
		ID:        uuid.NewString(),
		Name:      "A",
		Price:     10,
		Stock:     5,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if _, err := r.Create(context.Background(), p); err != nil { t.Fatalf("create: %v", err) }

	if _, err := r.UpdateStock(context.Background(), p.ID, -3); err != nil { t.Fatalf("update: %v", err) }

	out, err := r.GetByID(context.Background(), p.ID)
	if err != nil { t.Fatalf("get: %v", err) }
	if out.Stock != 2 { t.Fatalf("want 2 got %d", out.Stock) }
}
