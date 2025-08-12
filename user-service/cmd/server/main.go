package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	userpb "github.com/huntercenter1/backend-test/proto"
	dbpkg "github.com/huntercenter1/backend-test/user-service/internal/db"
	"github.com/huntercenter1/backend-test/user-service/internal/repo"
	"github.com/huntercenter1/backend-test/user-service/internal/service"
	grpcsvr "github.com/huntercenter1/backend-test/user-service/internal/transport/grpc"
)

func main() {
	if err := dbpkg.Migrate(os.Getenv("DB_DSN"), os.Getenv("MIGRATIONS_DIR")); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	db, err := dbpkg.New(context.Background())
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer func() { _ = db.Close() }()

	r := repo.NewUserRepo(db)
	svc := service.New(r)
	h := grpcsvr.NewServer(svc)

	addr := getenv("APP_PORT", ":50051")
	lis, err := net.Listen("tcp", addr)
	if err != nil { log.Fatalf("listen: %v", err) }

	s := grpc.NewServer()
	userpb.RegisterUserServiceServer(s, h)

	// SIEMPRE habilitar reflection para debug
	reflection.Register(s)
	log.Println("gRPC reflection enabled")

	go func() {
		log.Printf("user-service gRPC listening on %s", addr)
		if err := s.Serve(lis); err != nil { log.Fatalf("grpc serve: %v", err) }
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("shutting down gRPC server...")
	stopped := make(chan struct{})
	go func() { s.GracefulStop(); close(stopped) }()
	select {
	case <-stopped:
	case <-time.After(10 * time.Second):
		s.Stop()
	}
}

func getenv(k, d string) string { v := os.Getenv(k); if v == "" { return d }; return v }
