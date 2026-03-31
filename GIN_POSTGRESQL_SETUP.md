# WMS Backend Setup Guide - Gin + PostgreSQL

Dokumentasi lengkap untuk setup dan menjalankan Warehouse Management System (WMS) dengan framework Gin dan database PostgreSQL.

## 📋 Prasyarat

- **Go**: Version 1.20 atau lebih tinggi
- **PostgreSQL**: Version 12 atau lebih tinggi
- **Docker & Docker Compose**: (Opsional, untuk PostgreSQL dalam kontainer)
- **Git**: Untuk clone repository

## 🚀 Quick Start

### 1. Setup Environment Variables

```bash
cp .env.example .env
```

Edit file `.env` dan sesuaikan dengan konfigurasi lokal:

```env
# Database
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=wms_db
DB_PORT=5432
DB_SSLMODE=disable

# Server
PORT=8080
GIN_MODE=debug
```

### 2. Start PostgreSQL Database

**Opsi A: Menggunakan Docker Compose (Recommended)**

```bash
make docker-up
```

Ini akan menjalankan:

- PostgreSQL database (port 5432)
- pgAdmin interface (port 5050) - akses di http://localhost:5050

**Opsi B: PostgreSQL Lokal**

Pastikan PostgreSQL sudah berjalan, kemudian buat database:

```bash
psql -U postgres -c "CREATE DATABASE wms_db;"
```

### 3. Install Dependencies

```bash
make install
```

atau

```bash
go mod tidy
go mod download
```

### 4. Run Migrations

Migrations akan berjalan otomatis saat aplikasi startup. Models yang di-migrate:

- **User Management**: User, Tax
- **Products**: Category, Sticker, ProductDocument, ProductMaster, ProductLog
- **Warehouse**: Rack, Store, StoreCrew
- **Transfers**: StoreTransfer, StoreTransferBag
- **Inventory**: Cargo, Bag, SlowMoving, SlowMovingItem
- **Orders**: Order, OrderItem, OrderCargo
- **Classification**: Class, Buyer, UserClassLog

### 5. Start Development Server

```bash
make dev
```

atau

```bash
GIN_MODE=debug go run main.go
```

Server akan berjalan di: **http://localhost:8080**

## 📡 Testing API

### Health Check

```bash
curl http://localhost:8080/health
```

Respons:

```json
{
  "status": "success",
  "message": "WMS API is running"
}
```

### Menggunakan REST Client

File `api-requests.rest` berisi contoh semua API requests. Buka dengan VS Code REST Client extension:

```bash
# Buka file api-requests.rest dan klik "Send Request" pada endpoint yang ingin ditest
```

### Contoh: Create User

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "081234567890",
    "password": "securepassword"
  }'
```

## 📁 Project Structure

```
wms-v3/
├── config/              # Configuration files
│   ├── database.go      # Database initialization
│   ├── env.go           # Environment variables
│   └── middleware.go    # Gin middlewares (CORS, Logging, etc)
├── controller/          # HTTP controllers (Gin handlers)
│   ├── user_controller.go
│   ├── product_controller.go
│   └── order_controller.go
├── models/              # Database models (GORM)
│   ├── user.go
│   ├── product.go
│   ├── order.go
│   └── ... (other models)
├── repositories/        # Data access layer
│   ├── user_repository.go
│   ├── product_repository.go
│   └── order_repository.go
├── services/            # Business logic layer
│   ├── user_service.go
│   ├── product_service.go
│   └── order_service.go
├── routes/              # Route definitions (Gin routes)
│   └── routes.go
├── utils/               # Utility functions
│   ├── response.go      # API response helpers
│   ├── validators.go    # Input validation
│   └── ... (other utilities)
├── db/                  # Database scripts
│   └── schema.sql       # SQL schema reference
├── main.go              # Application entry point
├── Makefile             # Make commands
├── docker-compose.yml   # Docker configuration
└── .env.example         # Environment template
```

## 🏗️ Architecture

WMS menggunakan **Clean Architecture** pattern:

```
HTTP Request
    ↓
Routes (Gin Router)
    ↓
Controllers (HTTP handlers)
    ↓
Services (Business logic)
    ↓
Repositories (Database access)
    ↓
Models (GORM)
    ↓
PostgreSQL Database
```

### Layer Descriptions

1. **Controllers**: Handle HTTP requests/responses menggunakan Gin context
2. **Services**: Contain business logic dan rules
3. **Repositories**: Handle database queries dengan GORM
4. **Models**: Define database schema dan relationships

## 🔌 Database Connection

WMS menggunakan **GORM** sebagai ORM dengan konfigurasi:

```go
// Connection Pooling
MaxIdleConns: 10
MaxOpenConns: 100
ConnMaxLifetime: 1 hour
```

Parameter ini dapat dikonfigurasi di `config/database.go`.

## 🛠️ Make Commands

```bash
make help          # Show all available commands
make install       # Install dependencies
make build         # Build for Linux
make run           # Run application
make dev           # Run in development mode
make test          # Run tests
make clean         # Clean build artifacts
make docker-up     # Start PostgreSQL + pgAdmin
make docker-down   # Stop Docker containers
make docker-logs   # View PostgreSQL logs
make lint          # Run linter (requires golangci-lint)
make fmt           # Format code
make vet           # Run go vet
```

## 📊 Accessing pgAdmin

Jika menggunakan Docker Compose, pgAdmin tersedia di: **http://localhost:5050**

**Login credentials:**

- Email: `admin@example.com`
- Password: `admin`

**Menambah PostgreSQL Server:**

1. Klik "Add New Server"
2. General tab:
   - Name: `WMS Database`
3. Connection tab:
   - Host: `postgres` (nama service di docker-compose)
   - Port: `5432`
   - Username: `postgres`
   - Password: `password`

## 🔍 API Overview

### Available Endpoints

#### Users

- `POST /api/v1/users` - Create user
- `GET /api/v1/users` - List users (paginated)
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

#### Products

- `POST /api/v1/products` - Create product
- `GET /api/v1/products` - List products
- `GET /api/v1/products/:id` - Get product by ID
- `GET /api/v1/products/barcode/:barcode` - Get by barcode
- `GET /api/v1/products/category/:categoryID` - Get by category
- `PUT /api/v1/products/:id` - Update product
- `DELETE /api/v1/products/:id` - Delete product

#### Orders

- `POST /api/v1/orders` - Create order
- `GET /api/v1/orders` - List orders
- `GET /api/v1/orders/:id` - Get order by ID
- `GET /api/v1/orders/code/:code` - Get by code
- `GET /api/v1/orders/status/:status` - Get by status
- `PUT /api/v1/orders/:id` - Update order
- `DELETE /api/v1/orders/:id` - Delete order

## 🌍 Environment Variables Reference

| Variable      | Default   | Description                   |
| ------------- | --------- | ----------------------------- |
| `DB_HOST`     | localhost | PostgreSQL host               |
| `DB_USER`     | postgres  | Database user                 |
| `DB_PASSWORD` | password  | Database password             |
| `DB_NAME`     | wms_db    | Database name                 |
| `DB_PORT`     | 5432      | Database port                 |
| `DB_SSLMODE`  | disable   | SSL mode (disable/require)    |
| `PORT`        | 8080      | Server port                   |
| `GIN_MODE`    | debug     | Gin mode (debug/release/test) |
| `API_HOST`    | localhost | API host                      |

## 📝 Middleware

WMS menggunakan beberapa middleware Gin:

1. **Logging**: Logs semua HTTP requests dan responses
2. **Recovery**: Catches panic dan returns 500 error
3. **CORS**: Handles cross-origin requests
4. **Request ID**: Adds unique ID untuk setiap request

## 🔐 Security Notes

⚠️ **Production Setup:**

1. Ubah password PostgreSQL di `.env`
2. Set `GIN_MODE=release` untuk production
3. Gunakan `DB_SSLMODE=require` untuk PostgreSQL production
4. Implementasi JWT authentication di middleware
5. Validate dan sanitize semua inputs
6. Use environment variables untuk sensitive data

## 🐛 Troubleshooting

### Connection Refused

```
Failed to connect to database: connection refused
```

**Solution**: Pastikan PostgreSQL running dan credentials di `.env` benar.

### Database Already Exists

```
database "wms_db" already exists
```

**Solution**: Gunakan existing database atau drop dulu:

```bash
psql -U postgres -c "DROP DATABASE wms_db;"
```

### Port Already in Use

```
listen tcp :8080: bind: address already in use
```

**Solution**: Ubah PORT di `.env` atau kill process:

```bash
# Windows
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# Linux/Mac
lsof -i :8080
kill -9 <PID>
```

### CORS Issues

Solution sudah di-handle oleh middleware di `config/middleware.go`.

## 📚 Next Steps

Setelah setup berhasil:

1. **Implementasi Authentication**: JWT, OAuth
2. **Validation**: Input validation, business rules
3. **Testing**: Unit tests, integration tests
4. **Documentation**: Swagger/OpenAPI
5. **Deployment**: Docker, Kubernetes, Cloud platforms
6. **Monitoring**: Logging, Metrics, Tracing

## 📞 Support & Resources

- **Go Documentation**: https://golang.org/doc
- **Gin Documentation**: https://gin-gonic.com
- **GORM Documentation**: https://gorm.io
- **PostgreSQL Documentation**: https://www.postgresql.org/docs
