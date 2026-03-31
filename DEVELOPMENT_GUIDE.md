# WMS Development Guide

Panduan lengkap untuk development, testing, dan deployment WMS API.

## 📋 Daftar Isi

1. [Setup Development Environment](#setup-development-environment)
2. [Development Workflow](#development-workflow)
3. [Adding New Features](#adding-new-features)
4. [Testing](#testing)
5. [Git Workflow](#git-workflow)
6. [Deployment](#deployment)
7. [Best Practices](#best-practices)

## Setup Development Environment

### 1. Prerequisites

- Go 1.23+
- PostgreSQL 12+
- Git
- Docker & Docker Compose (opsional)

### 2. Initial Setup

```bash
# Clone repository
git clone <repository-url>
cd wms-v3

# Copy environment file
cp .env.example .env

# Edit .env dengan konfigurasi lokal
nano .env

# Install dependencies
make install

# Start PostgreSQL
make docker-up

# Wait for database to be ready
sleep 5

# Run migrations
go run main.go
```

### 3. Verify Setup

```bash
# Check API health
curl http://localhost:8080/health

# Expected response:
# {
#   "status": "success",
#   "message": "WMS API is running"
# }
```

## Development Workflow

### Daily Development

```bash
# Start server in development mode
make dev
# atau
GIN_MODE=debug go run main.go

# In another terminal, test API
curl http://localhost:8080/api/v1/users

# Stop server
Ctrl+C

# Run tests
make test

# Check code quality
make lint
make vet
```

### Making Changes

1. **Create feature branch**

   ```bash
   git checkout -b feature/new-feature-name
   ```

2. **Make changes** (see Adding New Features section)

3. **Test changes**

   ```bash
   make test
   ```

4. **Format code**

   ```bash
   make fmt
   go mod tidy
   ```

5. **Commit changes**
   ```bash
   git add .
   git commit -m "feature: add new feature"
   git push origin feature/new-feature-name
   ```

## Adding New Features

### Example: Add Category Management

#### Step 1: Create/Verify Model

Model sudah ada di `models/category.go`:

```go
type Category struct {
    ID        string    `gorm:"primaryKey;type:uuid" json:"id"`
    Name      string    `gorm:"type:varchar(255);not null" json:"name"`
    Slug      string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
    // ... fields
}
```

#### Step 2: Create Repository Interface

Tambahkan ke `repositories/interface.go`:

```go
type CategoryRepository interface {
    Create(ctx context.Context, category *models.Category) error
    GetByID(ctx context.Context, id string) (*models.Category, error)
    GetBySlug(ctx context.Context, slug string) (*models.Category, error)
    GetAll(ctx context.Context, limit, offset int) ([]models.Category, int64, error)
    Update(ctx context.Context, id string, category *models.Category) error
    Delete(ctx context.Context, id string) error
}
```

#### Step 3: Implement Repository

Buat file `repositories/category_repository.go`:

```go
package repositories

import (
    "context"
    "wms/config"
    "wms/models"
)

type categoryRepository struct{}

func NewCategoryRepository() CategoryRepository {
    return &categoryRepository{}
}

func (r *categoryRepository) Create(ctx context.Context, category *models.Category) error {
    return config.DB.WithContext(ctx).Create(category).Error
}

// ... implement other methods
```

#### Step 4: Create Service Interface & Implementation

Buat file `services/category_service.go`:

```go
package services

import (
    "context"
    "wms/models"
    "wms/repositories"
)

type CategoryService interface {
    CreateCategory(ctx context.Context, category *models.Category) error
    GetCategory(ctx context.Context, id string) (*models.Category, error)
    GetCategoryBySlug(ctx context.Context, slug string) (*models.Category, error)
    ListCategories(ctx context.Context, limit, offset int) ([]models.Category, int64, error)
    UpdateCategory(ctx context.Context, id string, category *models.Category) error
    DeleteCategory(ctx context.Context, id string) error
}

type categoryService struct {
    repo repositories.CategoryRepository
}

func NewCategoryService(repo repositories.CategoryRepository) CategoryService {
    return &categoryService{repo: repo}
}

// ... implement methods
```

#### Step 5: Create Controller

Buat file `controller/category_controller.go`:

```go
package controller

import (
    "github.com/gin-gonic/gin"
    "wms/models"
    "wms/services"
    "wms/utils"
)

type CategoryController struct {
    service services.CategoryService
}

func NewCategoryController(service services.CategoryService) *CategoryController {
    return &CategoryController{service: service}
}

func (c *CategoryController) CreateCategory(ctx *gin.Context) {
    var category models.Category
    if err := ctx.ShouldBindJSON(&category); err != nil {
        utils.SendError(ctx, 400, "Invalid request body")
        return
    }

    if err := c.service.CreateCategory(ctx.Request.Context(), &category); err != nil {
        utils.SendError(ctx, 500, "Failed to create category")
        return
    }

    utils.SendSuccess(ctx, category, "Category created successfully", 201)
}

// ... implement other handlers
```

#### Step 6: Register Routes

Update `routes/routes.go`:

```go
// Include this in setupRoutes function
categoryRepo := repositories.NewCategoryRepository()
categoryService := services.NewCategoryService(categoryRepo)
categoryCtrl := controller.NewCategoryController(categoryService)

categoryGroup := r.Group("/api/v1/categories")
{
    categoryGroup.POST("", categoryCtrl.CreateCategory)
    categoryGroup.GET("", categoryCtrl.ListCategories)
    categoryGroup.GET("/:id", categoryCtrl.GetCategory)
    categoryGroup.PUT("/:id", categoryCtrl.UpdateCategory)
    categoryGroup.DELETE("/:id", categoryCtrl.DeleteCategory)
}
```

## Testing

### Unit Tests

Struktur test file:

```
somefeature_test.go
```

```go
package services

import (
    "context"
    "testing"
    "wms/models"
)

func TestCreateCategory(t *testing.T) {
    // Arrange
    mockRepo := NewMockCategoryRepository()
    service := NewCategoryService(mockRepo)

    category := &models.Category{
        Name: "Electronics",
        Slug: "electronics",
    }

    // Act
    err := service.CreateCategory(context.Background(), category)

    // Assert
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
}
```

### Running Tests

```bash
# Run all tests
make test

# Run specific test
go test -run TestCreateCategory ./services

# With coverage
make test-coverage

# Verbose output
go test -v ./...
```

### Integration Tests

Gunakan database testing dengan fixtures:

```go
// Test dengan real database
func TestIntegration_CreateCategory(t *testing.T) {
    // Setup database
    db := setupTestDB(t)
    defer teardownTestDB(db)

    repo := repositories.NewCategoryRepository()
    service := services.NewCategoryService(repo)

    // Test logic
    // ...
}
```

## Git Workflow

### Branching Strategy

```
main (production-ready)
├── develop (integration branch)
│   ├── feature/user-management
│   ├── feature/product-catalog
│   ├── bugfix/issue-123
│   └── release/v1.0.0
```

### Commit Messages

Format: `<type>(<scope>): <subject>`

```bash
# Examples
git commit -m "feat(auth): add JWT authentication"
git commit -m "fix(product): correct price calculation"
git commit -m "docs(api): update endpoint documentation"
git commit -m "refactor(repo): simplify database queries"
git commit -m "test(order): add integration tests"
```

Types: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`, `style`

### Pull Request Process

1. Create branch dari `develop`
2. Make changes dan test
3. Push ke remote
4. Create Pull Request
5. Code review oleh team
6. Merge ke `develop` dengan squash commits

## Deployment

### Development Deployment

```bash
# Build binary
make build

# Run binary
./bin/wms
```

### Docker Deployment

```bash
# Build Docker image
docker build -t wms:latest .

# Run container
docker run -p 8080:8080 \
  -e DB_HOST=postgres \
  -e DB_USER=postgres \
  -e DB_PASSWORD=password \
  -e DB_NAME=wms_db \
  wms:latest
```

### Production Deployment

1. **Set environment variables**

   ```
   GIN_MODE=release
   PORT=8080
   DB_HOST=prod-db.example.com
   etc...
   ```

2. **Run migrations**

   ```bash
   ./wms # akan automatic migrate
   ```

3. **Monitor logs**
   ```bash
   tail -f /var/log/wms/app.log
   ```

## Best Practices

### Code Quality

1. **Follow Go conventions**
   - Use `golangci-lint` untuk linting
   - Format dengan `gofmt`
   - Vet dengan `go vet`

2. **Error Handling**

   ```go
   // Good
   if err != nil {
       return fmt.Errorf("failed to create user: %w", err)
   }

   // Avoid
   if err != nil {
       panic(err) // Jangan gunakan di production
   }
   ```

3. **Context Usage**

   ```go
   // Selalu gunakan context parameter
   func (s *userService) CreateUser(ctx context.Context, user *models.User) error {
       return s.repo.Create(ctx, user)
   }
   ```

4. **Logging**
   ```go
   // Use standard logging atau structured logging
   log.Printf("User created: %s", user.ID)
   ```

### Database

1. **Migrations**
   - Use GORM AutoMigrate
   - Keep migrations idempotent
   - Never hardcoded business logic dalam migration

2. **Queries**
   - Use parameterized queries (GORM handles this)
   - Index frequently queried fields
   - Use preloading untuk relationships

3. **Transactions**
   ```go
   // Multi-step operations require transactions
   tx := config.DB.BeginTx(ctx, nil)
   // ... operations
   tx.Commit()
   ```

### API Design

1. **RESTful Principles**
   - POST untuk create
   - GET untuk read
   - PUT untuk update
   - DELETE untuk delete

2. **Consistent Response Format**
   - Selalu gunakan Response struct
   - Consistent error messages
   - Proper HTTP status codes

3. **Validation**
   ```go
   // Validate input
   if !utils.ValidateEmail(user.Email) {
       return utils.NewBadRequestError("Invalid email")
   }
   ```

### Performance

1. **Database**
   - Use indexes untuk frequently queried fields
   - Limit result sets dengan pagination
   - Preload relationships secara bijak

2. **Caching** (Future)
   - Cache frequently accessed data
   - Set appropriate TTL
   - Invalidate cache when data changes

3. **Concurrency**
   ```go
   // Use goroutines untuk I/O operations
   go func() {
       // Background processing
   }()
   ```

## Troubleshooting

### Database Connection Issues

```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# View logs
docker logs wms_postgres

# Restart services
make docker-down
make docker-up
```

### Port Already in Use

```bash
# Find process using port 8080
lsof -i :8080

# Kill process
kill -9 <PID>
```

### Module Not Found

```bash
# Clear Go cache
go clean -modcache

# Get dependencies
go get -u ./...

# Verify
go mod verify
```

## IDE Setup

### VS Code

Recommended extensions:

- Go (golang.go)
- PostgreSQL (cweijan.vscode-postgresql-client2)
- REST Client (humao.rest-client)

Extensions configuration:

```json
{
  "[go]": {
    "editor.formatOnSave": true,
    "editor.defaultFormatter": "golang.go",
    "editor.codeActionsOnSave": {
      "source.organizeImports": true
    }
  }
}
```

### GoLand/IntelliJ

- Built-in Go support
- Integrated debugger
- Database tools

## Resources

- [Go Documentation](https://golang.org/doc/)
- [GORM Documentation](https://gorm.io/)
- [Gin Framework](https://gin-gonic.com/)
- [PostgreSQL Docs](https://www.postgresql.org/docs/)

---

**Last Updated**: 2024
