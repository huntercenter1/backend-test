package db

import (
    "context"
    "database/sql"
    "os"
    "time"

    "github.com/uptrace/bun"
    "github.com/uptrace/bun/dialect/pgdialect"
    "github.com/uptrace/bun/extra/bundebug"
    _ "github.com/jackc/pgx/v5/stdlib"
)

func New(ctx context.Context) (*bun.DB, error) {
    dsn := os.Getenv("DB_DSN")
    sqldb, err := sql.Open("pgx", dsn)
    if err != nil { return nil, err }
    sqldb.SetMaxOpenConns(10)
    sqldb.SetMaxIdleConns(5)
    sqldb.SetConnMaxLifetime(30 * time.Minute)

    db := bun.NewDB(sqldb, pgdialect.New())
    if os.Getenv("APP_ENV") == "local" {
        db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
    }
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := db.PingContext(ctx); err != nil { return nil, err }
    return db, nil
}
