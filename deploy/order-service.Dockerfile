FROM golang:1.23 AS build
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOWORK=off GOFLAGS="-mod=vendor"
WORKDIR /src

# Copiamos el servicio (incluye vendor) y los protos
COPY order-service/ order-service/
COPY proto/ proto/

WORKDIR /src/order-service

# Compila usando vendor (sin tidy ni download)
RUN go build -trimpath -ldflags="-s -w" -o /out/order-service ./cmd/server

FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=build /out/order-service /app/order-service
COPY order-service/migrations /app/migrations
ENV APP_PORT=:8082
ENTRYPOINT ["/app/order-service"]
