FROM golang:1.23 AS build
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
WORKDIR /src
COPY user-service/ user-service/
COPY proto/ proto/
WORKDIR /src/user-service
RUN go mod tidy
RUN go mod download
RUN go build -trimpath -ldflags="-s -w" -o /out/user-service ./cmd/server

FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=build /out/user-service /app/user-service
COPY user-service/migrations /app/migrations
ENV APP_PORT=:50051
ENTRYPOINT ["/app/user-service"]
