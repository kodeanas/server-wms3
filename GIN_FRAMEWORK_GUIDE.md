# WMS Backend - Gin Framework Implementation

Dokumentasi implementasi **Gin Web Framework** di WMS Backend.

## 🎯 Gin Framework Overview

**Gin** adalah web framework yang cepat dan minimalist untuk Go. WMS menggunakan Gin untuk:

- HTTP routing dan handling
- Middleware management
- Request/response binding
- JSON marshaling/unmarshaling
- Built-in logging dan recovery

## 📦 Dependencies

```go
// go.mod
github.com/gin-gonic/gin v1.9.0          // Gin framework
gorm.io/driver/postgres v1.5.7           // PostgreSQL driver
gorm.io/gorm v1.25.7                     // ORM
```

## 🏗️ Implementation Details

### 1. Router Setup (main.go)

```go
// Create Gin router dengan default middleware (logger + recovery)
r := gin.Default()

// Add custom middlewares
r.Use(config.CORSMiddleware())
r.Use(config.RequestIDMiddleware())

// Setup all routes
routes.SetupRoutes(r)

// Start server
r.Run(":8080")
```

### 2. Middleware (config/middleware.go)

#### CORS Middleware

```go
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        // ... handle CORS headers
        c.Next()
    }
}
```

#### Request ID Middleware

```go
func RequestIDMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := c.GetHeader("X-Request-ID")
        c.Set("RequestID", requestID)
        c.Header("X-Request-ID", requestID)
        c.Next()
    }
}
```

### 3. Route Definition (routes/routes.go)

```go
func SetupRoutes(r *gin.Engine) {
    // Initialize repositories, services, controllers

    // User routes
    userGroup := r.Group("/api/v1/users")
    {
        userGroup.POST("", userCtrl.CreateUser)
        userGroup.GET("", userCtrl.ListUsers)
        userGroup.GET("/:id", userCtrl.GetUser)
        userGroup.PUT("/:id", userCtrl.UpdateUser)
        userGroup.DELETE("/:id", userCtrl.DeleteUser)
    }

    // Product routes
    productGroup := r.Group("/api/v1/products")
    // ...

    // Order routes
    orderGroup := r.Group("/api/v1/orders")
    // ...
}
```

### 4. Controllers dengan Gin Context

```go
// User Controller Example
type UserController struct {
    service services.UserService
}

func (c *UserController) CreateUser(ctx *gin.Context) {
    var user models.User

    // Bind JSON request body ke struct
    if err := ctx.ShouldBindJSON(&user); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request body",
        })
        return
    }

    // Call service business logic
    if err := c.service.CreateUser(ctx.Request.Context(), &user); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to create user",
        })
        return
    }

    // Send success response
    ctx.JSON(http.StatusCreated, gin.H{
        "status": "success",
        "data":   user,
    })
}

func (c *UserController) GetUser(ctx *gin.Context) {
    // Get path parameter
    id := ctx.Param("id")

    // Get query parameter
    limit := ctx.Query("limit")

    // Get header
    token := ctx.GetHeader("Authorization")

    // Set response header
    ctx.Header("X-Custom-Header", "value")

    // Send JSON response
    ctx.JSON(http.StatusOK, gin.H{
        "data": user,
    })
}
```

### 5. Gin Context Methods

| Method                        | Deskripsi                  |
| ----------------------------- | -------------------------- |
| `ctx.Param(key)`              | Get URL path parameter     |
| `ctx.Query(key)`              | Get query string parameter |
| `ctx.PostForm(key)`           | Get POST form data         |
| `ctx.GetHeader(key)`          | Get HTTP header            |
| `ctx.Header(key, val)`        | Set HTTP header            |
| `ctx.ShouldBindJSON(obj)`     | Bind & validate JSON body  |
| `ctx.JSON(code, obj)`         | Send JSON response         |
| `ctx.String(code, fmt, args)` | Send text response         |
| `ctx.File(path)`              | Send file                  |
| `ctx.Redirect(code, url)`     | Redirect                   |
| `ctx.AbortWithStatus(code)`   | Abort with status          |
| `ctx.Request.Context()`       | Get context.Context        |
| `ctx.Set(key, val)`           | Set value in context       |
| `ctx.Get(key)`                | Get value from context     |

### 6. Response Helpers (utils/response.go)

```go
// Success Response
func SendSuccess(ctx *gin.Context, data interface{}, message string, statusCode ...int) {
    code := http.StatusOK
    if len(statusCode) > 0 {
        code = statusCode[0]
    }
    ctx.JSON(code, gin.H{
        "status":  "success",
        "message": message,
        "data":    data,
    })
}

// Error Response
func SendError(ctx *gin.Context, statusCode int, message string) {
    ctx.JSON(statusCode, gin.H{
        "status":  "error",
        "message": message,
    })
}

// Paginated Response
func SendPaginated(ctx *gin.Context, data interface{}, total int64, page, limit int) {
    ctx.JSON(http.StatusOK, gin.H{
        "status": "success",
        "data":   data,
        "pagination": gin.H{
            "total":       total,
            "page":        page,
            "limit":       limit,
            "total_pages": (total + int64(limit) - 1) / int64(limit),
        },
    })
}
```

### 7. Database Integration

WMS menggunakan **GORM** dengan Gin context:

```go
// Repositories menerima context dari Gin
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
    var user models.User
    if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return &user, nil
}
```

## 🔄 Request-Response Flow

```
1. HTTP Request
   └─→ Gin Router (routes.go)

2. Routing Match
   └─→ Controller Handler (gin.Context)

3. Request Binding
   └─→ ctx.ShouldBindJSON()
   └─→ utils.BindJSONOrFail()

4. Business Logic
   └─→ Service Layer

5. Data Access
   └─→ Repository Layer
   └─→ GORM Queries

6. Response Building
   └─→ utils.SendSuccess()
   └─→ utils.SendError()

7. HTTP Response
   └─→ ctx.JSON()
```

## 📝 Request Binding Examples

### JSON Binding

```go
// Auto bind JSON body to struct
var user models.User
if err := ctx.ShouldBindJSON(&user); err != nil {
    ctx.JSON(400, err)
    return
}
```

### Query Parameters

```go
// Get query params: GET /products?category=electronics&status=active
category := ctx.Query("category")
status := ctx.DefaultQuery("status", "active")
```

### URL Parameters

```go
// Get path params: GET /users/:id
id := ctx.Param("id")
```

### Form Data

```go
// POST /upload (form-data)
file, _ := ctx.FormFile("file")
name := ctx.PostForm("name")
```

## 🌳 Router Groups

Gin supports grouped routes dengan common prefix:

```go
// Group dengan common prefix dan shared middleware
v1 := r.Group("/api/v1")
v1.Use(authMiddleware)
{
    v1.POST("/users", createUser)
    v1.GET("/users", listUsers)

    products := v1.Group("/products")
    {
        products.POST("", createProduct)
        products.GET("", listProducts)
    }
}
```

## ⚡ Performance Features

### Built-in Logging

```
[GIN-debug] Listening and serving HTTP on :8080
[GIN] 2024/01/01 - 12:00:00 | 200 |     123.45ms | 192.168.1.1 | POST  /api/v1/users
```

### Connection Pooling

```go
// Database connection pooling di config/database.go
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

### Middleware Order

1. CORS Middleware (allow cross-origin)
2. Request ID Middleware (tracking)
3. Logging (built-in)
4. Recovery (panic handling)
5. Business logic

## 🔐 Security Best Practices

### 1. Input Validation

```go
var user models.User
if err := ctx.ShouldBindJSON(&user); err != nil {
    ctx.JSON(400, gin.H{"error": err.Error()})
    return
}
// Validate email
if !utils.ValidateEmail(user.Email) {
    ctx.JSON(400, gin.H{"error": "Invalid email"})
    return
}
```

### 2. Error Handling

```go
if err != nil {
    // Log error internally
    log.Printf("Error: %v", err)
    // Return safe error message to client
    ctx.JSON(500, gin.H{"error": "Internal server error"})
}
```

### 3. CORS Configuration

```go
// Current: Allow all origins
c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

// Production: Restrict origins
c.Writer.Header().Set("Access-Control-Allow-Origin", "https://trusted-domain.com")
```

## 📚 Common Patterns

### 1. Create Resource

```go
func (c *Controller) Create(ctx *gin.Context) {
    var req CreateRequest
    if !utils.BindJSONOrFail(ctx, &req) {
        return
    }

    resource := &models.Resource{...}
    if err := c.service.Create(ctx.Request.Context(), resource); err != nil {
        utils.SendError(ctx, 500, "Failed to create")
        return
    }

    utils.SendSuccess(ctx, resource, "Created", 201)
}
```

### 2. Get by ID

```go
func (c *Controller) GetByID(ctx *gin.Context) {
    id := ctx.Param("id")

    resource, err := c.service.GetByID(ctx.Request.Context(), id)
    if err != nil {
        utils.SendError(ctx, 404, "Not found")
        return
    }

    utils.SendSuccess(ctx, resource, "OK")
}
```

### 3. List with Pagination

```go
func (c *Controller) List(ctx *gin.Context) {
    limit, offset := utils.GetPaginationParams(ctx, 10)

    resources, total, err := c.service.List(ctx.Request.Context(), limit, offset)
    if err != nil {
        utils.SendError(ctx, 500, "Failed to fetch")
        return
    }

    utils.SendPaginated(ctx, resources, total, (offset/limit)+1, limit)
}
```

## 🚀 Deployment Tips

### Development Mode

```bash
GIN_MODE=debug go run main.go
```

### Production Mode

```bash
GIN_MODE=release go run main.go
```

### Docker

```dockerfile
FROM golang:1.20-alpine
WORKDIR /app
COPY . .
RUN go build -o wms main.go
EXPOSE 8080
CMD ["./wms"]
```

---

**Dokumentasi Gin**: https://gin-gonic.com/docs/
**Dokumentasi GORM**: https://gorm.io/docs/
