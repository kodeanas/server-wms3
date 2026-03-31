# WMS (Warehouse Management System) API

Sistem manajemen gudang berbasis Go dengan Gin framework dan PostgreSQL database. Aplikasi ini menyediakan API lengkap untuk mengelola inventori, pesanan, pengguna, dan operasi gudang lainnya.

## 📋 Fitur Utama

### User Management

- Manajemen pengguna dengan role dan permission
- Tracking log perubahan kelas pengguna
- Sistem tax management

### Product Management

- Katalog produk dengan barcode support
- Kategorisasi dan label produk dengan stiker
- Tracking harga dan kuantitas per warehouse
- Product log untuk audit trail

### Order Management

- Pembuatan dan pengelolaan pesanan
- Order items dengan pricing detail
- Shipping/Cargo tracking
- Order status tracking

### Warehouse Operations

- Manajemen rak (rack) gudang
- Store management
- Transfer antar gudang dengan tracking bag
- Slow moving product tracking

### Cargo & Bag Management

- Cargo tracking dan management
- Bag grouping untuk pengiriman
- Status tracking real-time

## 🗄️ Database Schema

### Core Tables

- **users**: Data pengguna sistem
- **categories**: Kategori produk
- **stickers**: Label/stiker produk
- **classes**: Kelas pembeli
- **buyers**: Data pembeli

### Product Tables

- **product_documents**: Dokumen produk masuk
- **product_masters**: Master data produk
- **product_logs**: Audit log perubahan produk

### Warehouse Tables

- **racks**: Rak gudang
- **stores**: Lokasi toko/gudang
- **store_crews**: Staff toko
- **store_transfers**: Pengiriman antar toko
- **store_transfer_bags**: Item dalam pengiriman

### Business Tables

- **orders**: Pesanan
- **order_items**: Item dalam pesanan
- **order_cargos**: Kargo pesanan
- **cargos**: Data kargo
- **bags**: Paket/tas barang
- **slow_movings**: Tracking slow moving products

### Transactional Tables

- **taxes**: Data pajak
- **user_class_logs**: History perubahan kelas user

## 🚀 Quick Start

### Prerequisites

- Go 1.23+
- PostgreSQL 12+
- Git

### Installation

1. **Clone repository**

```bash
cd wms-v3
```

2. **Set up environment variables**

```bash
cp .env.example .env
# Edit .env dengan konfigurasi database Anda
```

3. **Install dependencies**

```bash
go mod download
go mod tidy
```

4. **Create database**

```bash
createdb wms_db
```

5. **Run application**

```bash
go run main.go
```

Server akan berjalan di `http://localhost:8080`

## 📡 API Endpoints

### Users (`/api/v1/users`)

```
POST   /api/v1/users              - Create user
GET    /api/v1/users              - List users (paginated)
GET    /api/v1/users/:id          - Get user by ID
PUT    /api/v1/users/:id          - Update user
DELETE /api/v1/users/:id          - Delete user
```

### Products (`/api/v1/products`)

```
POST   /api/v1/products                    - Create product
GET    /api/v1/products                    - List products (paginated)
GET    /api/v1/products/:id                - Get product by ID
GET    /api/v1/products/barcode/:barcode   - Get product by barcode
GET    /api/v1/products/category/:categoryID - Get products by category
PUT    /api/v1/products/:id                - Update product
DELETE /api/v1/products/:id                - Delete product
```

### Orders (`/api/v1/orders`)

```
POST   /api/v1/orders                  - Create order
GET    /api/v1/orders                  - List orders (paginated)
GET    /api/v1/orders/:id              - Get order by ID
GET    /api/v1/orders/code/:code       - Get order by code
GET    /api/v1/orders/status/:status   - Get orders by status
PUT    /api/v1/orders/:id              - Update order
DELETE /api/v1/orders/:id              - Delete order
```

## 📁 Project Structure

```
wms-v3/
├── config/              - Database & environment configuration
├── controller/          - HTTP request handlers
├── db/                  - Database migrations & scripts
├── models/              - Data models & structs
├── repositories/        - Data access layer
├── routes/              - API route definitions
├── services/            - Business logic layer
├── utils/               - Utility functions
├── main.go              - Application entry point
├── go.mod               - Go module definition
├── go.sum               - Go dependencies lock file
└── .env.example         - Environment variables template
```

## 🏗️ Architecture

Aplikasi ini menggunakan **Clean Architecture** dengan layers:

1. **Models** (`models/`) - Entity definitions
2. **Repositories** (`repositories/`) - Data access abstraction
3. **Services** (`services/`) - Business logic
4. **Controllers** (`controller/`) - HTTP handlers
5. **Routes** (`routes/`) - Route configuration

### Data Flow

```
Request → Controller → Service → Repository → Database
Response ← Controller ← Service ← Repository ← Database
```

## 🔑 Key Components

### Database Connection (config/database.go)

- GORM ORM untuk database abstraction
- Auto migration support
- Connection pooling

### Repository Pattern (repositories/)

- Interface-based design untuk testability
- Generic CRUD operations
- Pre-loading relationships

### Service Layer (services/)

- Business logic centralization
- Transaction handling
- Data validation

### Controllers (controller/)

- Request/response handling
- Error handling
- Pagination support

## 🛠️ Development

### Adding New Feature

1. **Create Model** di `models/`

   ```go
   type MyEntity struct {
       ID        string    `gorm:"primaryKey;type:uuid"`
       Name      string    `gorm:"type:varchar(255)"`
       CreatedAt time.Time
   }
   ```

2. **Create Repository Interface** di `repositories/`

   ```go
   type MyEntityRepository interface {
       Create(ctx context.Context, entity *MyEntity) error
       GetByID(ctx context.Context, id string) (*MyEntity, error)
       // ... other methods
   }
   ```

3. **Create Repository Implementation**

   ```go
   type myEntityRepository struct{}

   func NewMyEntityRepository() MyEntityRepository {
       return &myEntityRepository{}
   }
   ```

4. **Create Service**

   ```go
   type MyEntityService interface {
       CreateEntity(ctx context.Context, entity *MyEntity) error
       // ... other methods
   }
   ```

5. **Create Controller**

   ```go
   type MyEntityController struct {
       service MyEntityService
   }

   func (c *MyEntityController) CreateEntity(ctx *gin.Context) {
       // handler logic
   }
   ```

6. **Register Routes** di `routes/routes.go`

## 📊 Database Relationships

### One-to-Many

- User → Orders, Products, Cargos
- Category → Products
- Cargo → Bags
- Store → StoreCrew, StoreTransfers

### Many-to-Many

- Category ↔ Stickers (via category_stickers)
- Product ↔ Stickers (via product_stickers)

### Constraints

- Cascading deletes untuk maintaining referential integrity
- Foreign key indexes untuk query performance

## 🔒 Security Considerations

- Perlu implementasi authentication (JWT recommended)
- Input validation untuk semua endpoints
- Role-based access control (RBAC)
- SQL injection protection via GORM parameterized queries
- Password hashing untuk user passwords

## 📝 API Response Format

### Success Response

```json
{
  "status": "success",
  "data": { ... },
  "message": "Operation successful"
}
```

### Error Response

```json
{
  "status": "error",
  "message": "Error description"
}
```

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestUserCreate ./...
```

## 📦 Deployment

### Using Docker

```dockerfile
FROM golang:1.23-alpine
WORKDIR /app
COPY . .
RUN go build -o wms main.go
EXPOSE 8080
CMD ["./wms"]
```

### Environment Variables untuk Production

- Set `GIN_MODE=release`
- Configure proper database URL
- Set production PORT
- Add authentication secrets

## 📚 Dependencies

- **gin-gonic/gin** v1.9.0 - Web framework
- **gorm.io/gorm** v1.25.7 - ORM
- **gorm.io/driver/postgres** v1.5.7 - PostgreSQL driver
- **google/uuid** v1.6.0 - UUID generation
- **joho/godotenv** v1.5.1 - Environment variables

## 🐛 Troubleshooting

### Database Connection Error

- Pastikan PostgreSQL running
- Cek credentials di `.env`
- Verify database exists

### Migrations Failed

- Check table normalization
- Delete existing tables jika perlu reset
- Review migration constraints

## 📞 Support

Untuk pertanyaan atau issue, dokumentasi schema tersedia di attachments.

## 📄 License

Project ini menggunakan struktur untuk enterprise WMS system.

---

**Catatan**: Fondasi ini siap untuk expansion dengan fitur-fitur tambahan seperti authentication, advanced reporting, dan real-time tracking.
