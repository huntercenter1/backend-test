package db

import (
	"fmt"
	"os"
	"time"

	goose "github.com/pressly/goose/v3"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Migrate(dsn, dir string) error {
	if dir == "" { dir = os.Getenv("MIGRATIONS_DIR") }
	if err := goose.SetDialect("postgres"); err != nil { return err }

	deadline := time.Now().Add(60 * time.Second)
	var lastErr error
	for time.Now().Before(deadline) {
		db, err := goose.OpenDBWithDriver("pgx", dsn)
		if err == nil {
			err = goose.Up(db, dir)
			db.Close()
			if err == nil {
				return nil
			}
		}
		lastErr = err
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("migrate failed after retries: %w", lastErr)
}
