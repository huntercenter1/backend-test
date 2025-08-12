FROM golang:1.23 AS build
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
WORKDIR /src
COPY product-service/ product-service/
WORKDIR /src/product-service
RUN go mod tidy
RUN go mod download
RUN go build -trimpath -ldflags="-s -w" -o /out/product-service ./cmd/server

FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=build /out/product-service /app/product-service
COPY product-service/migrations /app/migrations
ENV APP_PORT=:8081
ENTRYPOINT ["/app/product-service"]
