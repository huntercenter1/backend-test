.PHONY: proto up down run test

proto:
	@echo "Generating gRPC code..."
	protoc -I proto \
		--go_out=paths=source_relative:proto \
		--go-grpc_out=paths=source_relative:proto \
		proto/user.proto

up:
	docker compose up -d --build

down:
	docker compose down -v

run:
	go run ./user-service/cmd/server & \
	go run ./product-service/cmd/server & \
	go run ./order-service/cmd/server

.PHONY: test

test:
	@echo ">> Ejecutando tests por servicio con cobertura..."
	@rm -f cover*.out cover.out
	@( cd user-service    && go test ./... -coverprofile=../cover_user.out    -covermode=atomic )
	@( cd product-service && go test ./... -coverprofile=../cover_product.out -covermode=atomic )
	@( cd order-service   && go test ./... -coverprofile=../cover_order.out   -covermode=atomic )

	@echo ">> Uniendo reportes..."
	@echo "mode: atomic" > cover.out
	@tail -n +2 cover_user.out    >> cover.out
	@tail -n +2 cover_product.out >> cover.out
	@tail -n +2 cover_order.out   >> cover.out

	@echo ">> Mostrando resumen..."
	@go tool cover -func=cover.out

	@echo ">> Generando reporte HTML..."
	@go tool cover -html=cover.out -o cover.html
	@echo "✅ Reporte generado en: $$(pwd)/cover.html"

