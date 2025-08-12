# Backend Test — Go (Microservicios)

## Servicios
- **user-service** (gRPC): CRUD de usuarios + Authenticate/Validate
- **product-service** (REST): CRUD, búsqueda y stock
- **order-service** (REST): crea pedidos, valida usuario (gRPC), verifica stock (HTTP)

## Requisitos
- Docker + Docker Compose
- Make, curl, jq (opcional pero útil)
- Go 1.23+ (solo si vas a compilar local)

## Levantar todo
```bash
make proto                 # generar stubs gRPC
docker compose build
docker compose up -d

Health

curl -s http://localhost:8081/health
curl -s http://localhost:8082/health
# user-service es gRPC (usa grpcurl)
grpcurl -plaintext localhost:50051 list


Pruebas rápidas


# Crear usuario
grpcurl -plaintext -d '{"username":"demo","email":"demo@example.com","password":"123456"}' \
  localhost:50051 user.UserService/CreateUser

# Crear producto
curl -s -X POST http://localhost:8081/products -H "Content-Type: application/json" \
  -d '{"name":"Laptop","description":"15 inch","price":1200,"stock":5}'

# Crear orden
# (pon USER_ID y PRODUCT_ID reales)
curl -s -X POST http://localhost:8082/orders -H "Content-Type: application/json" \
  -d '{"user_id":"<USER_ID>","items":[{"product_id":"<PRODUCT_ID>","quantity":2}]}'

------------------
 Swagger

## 📜 OpenAPI / Swagger

Este proyecto incluye la documentación OpenAPI para los servicios **Product** y **Order**.

- `docs/openapi-product.yaml`
- `docs/openapi-order.yaml`

### Ver en Swagger UI (Docker)

#### Order Service
```bash
docker run -p 8089:8080 \
  -e SWAGGER_JSON=/docs/openapi-order.yaml \
  -v $(pwd)/docs:/docs swaggerapi/swagger-ui
Abrir en navegador: http://localhost:8089

### Product Service

docker run -p 8090:8080 \
  -e SWAGGER_JSON=/docs/openapi-product.yaml \
  -v $(pwd)/docs:/docs swaggerapi/swagger-ui
Abrir en navegador: http://localhost:8090


### Levantar ambos Swagger UI al mismo tiempo

```bash
docker run -d --name swagger-order -p 8089:8080 \
  -e SWAGGER_JSON=/docs/openapi-order.yaml \
  -v $(pwd)/docs:/docs swaggerapi/swagger-ui

docker run -d --name swagger-product -p 8090:8080 \
  -e SWAGGER_JSON=/docs/openapi-product.yaml \
  -v $(pwd)/docs:/docs swaggerapi/swagger-ui

Para detenerlos:
docker stop swagger-order swagger-product && docker rm swagger-order swagger-product


- gRPC (User Service): ver `docs/user-service-grpc.md` o abrir `docs/user-service-grpc.html`


-------------

Migraciones
Goose corre automáticamente al iniciar cada servicio.

Migraciones en */migrations.
--------------

Configuración (env)
Ver docker-compose.yml. Puertos: 50051 (gRPC), 8081, 8082.

Postgres: 5433/5434/5435.
--------------

Estructura
proto/: .proto + stubs generados

*-service/internal/...: capas (db, repo, service, transport, clients)

deploy/*.Dockerfile: build multi-stage
--------------
Ejecutar test

make test

--------------
POSTMAN

Cómo hacer los requests gRPC en Postman (para crear usuarios)
Esto se hace una vez en la UI de Postman (no viene en el JSON de collection estándar).

Abre Postman → New → gRPC Request.

En “Server URL” pon: {{USER_GRPC_HOST}} (asegúrate de seleccionar el entorno backend-test).

Haz clic en Import .proto → selecciona proto/user.proto del repo.

En Method, elige user.UserService/CreateUser.

En Message pega un JSON de ejemplo:

{
  "username": "demo2",
  "email": "demo2@example.com",
  "password": "123456"
}

Click Invoke. Verás la respuesta con el id del usuario.

Repite para:

user.UserService/AuthenticateUser:

{ "username":"demo2", "password":"123456" }
user.UserService/GetUser:

{ "id": "{{USER_ID}}" }
user.UserService/ValidateUser:

{ "user_id": "{{USER_ID}}" }
Si Postman te pide reflection y no aparece el servicio, asegúrate:

que el contenedor de user-service tiene reflection habilitado (ya lo tenías activo), o

que importaste proto/user.proto en la request gRPC (es lo más seguro).
