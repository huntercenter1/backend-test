package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	dbpkg "github.com/huntercenter1/backend-test/order-service/internal/db"
	"github.com/huntercenter1/backend-test/order-service/internal/clients"
	"github.com/huntercenter1/backend-test/order-service/internal/repo"
	"github.com/huntercenter1/backend-test/order-service/internal/service"
	httpr "github.com/huntercenter1/backend-test/order-service/internal/transport/http"
)

func main() {
	if err := dbpkg.Migrate(os.Getenv("DB_DSN"), os.Getenv("MIGRATIONS_DIR")); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	db, err := dbpkg.New(context.Background())
	if err != nil { log.Fatalf("db: %v", err) }
	defer db.Close()

	// clients
	userAddr := getenv("USER_GRPC_ADDR", "user-service:50051")
	uc, closeUC, err := clients.NewUserClient(userAddr)
	if err != nil { log.Fatalf("user client: %v", err) }
	defer closeUC()
	pc := clients.NewProductClient(getenv("PRODUCT_BASE_URL", "http://product-service:8081"))

	// wiring
	rp := repo.New(db)
	svc := service.New(rp, uc, pc)
	rt := httpr.New(svc)

	// http
	r := gin.New()
	rt.Register(r)

	srv := &http.Server{ Addr: getenv("APP_PORT", ":8082"), Handler: r }
	go func() {
		log.Printf("order-service HTTP listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second); defer cancel()
	_ = srv.Shutdown(ctx)
}

func getenv(k, d string) string { v := os.Getenv(k); if v == "" { return d }; return v }
