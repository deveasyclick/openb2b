# OpenB2B Copilot Instructions

## Project Overview
OpenB2B is an open-source, multi-tenant ordering and invoicing platform for small wholesalers.  
It allows businesses to manage products, process orders, and track invoices.  
Supports multi-user roles per organization and ensures relational integrity between orders, products, and customers.

---

## Technology Stack
- **Backend:** Golang (`net/http`, `pgx`, PostgreSQL)
- **Frontend:** React + TypeScript
- **Validation:** Zod
- **State Management:** React Query or Context API
- **Styling:** TailwindCSS
- **PDF Generation:** Go libraries for server-side invoice generation

---

## Coding Conventions
- Types/interfaces: `PascalCase`
- Go: idiomatic code; services contain business logic, handlers are minimal
- Handle all errors explicitly; avoid `panic` except for fatal conditions
- Functional React components using hooks
- Multi-step forms: React Hook Form + Zod
- Keep business logic in services, not in components or handlers
- Use concise comments explaining **why**, not **what**

---

## Architecture Guidance
- **Models:** only DB fields; avoid methods that touch other models
- **Services:** implement business logic, validation, cross-model operations
- **Handlers/Controllers:** minimal logic; call services
- **Frontend components:** modular, reusable; separate UI and logic
- **Multi-tenancy:** all queries/actions respect `organization_id`
- Use transactions for multi-step DB operations (e.g., create order + items)

---

## Code Generation Preferences
- Generate type-safe code (TS interfaces, Go structs)
- Include pagination and filtering for listing endpoints
- Include role-based access checks
- Include unit tests (Go: services, React: components)
- Accessible HTML (semantic elements, ARIA roles)
- Use descriptive variable names
- Prefer simple, readable, maintainable solutions

---

## File Structure

```bash
├── README.md
├── api
│   ├── cmd/api/main.go
│   ├── db/migrations
│   ├── docs
│   ├── internal
│   │   ├── config
│   │   ├── db
│   │   ├── middleware
│   │   ├── model
│   │   ├── modules
│   │   │   ├── org
│   │   │   ├── product
│   │   │   ├── user
│   │   │   └── webhook
│   │   ├── routes
│   │   ├── shared
│   │   └── utils
│   ├── main
│   ├── pkg
│   │   ├── interfaces
│   │   ├── logger
│   │   └── svix
│   └── test/integration
├── docs
└── makefile
```

---

## Example Prompts
- "Generate a Go service to create an order with multiple items and validate stock"
- "Create a React component to list products with filters and pagination"
- "Generate PDF invoice for an order in Go using existing structs"
- "Add role-based access check middleware in Go"
- "Implement a multi-step product creation form using React Hook Form and Zod"

---

## Models

```go
type BaseModel struct {
    ID        uint
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt
}

type Address struct {
    Zip     string
    State   string
    City    string
    Country string
    Address string
}

type Customer struct {
    BaseModel
    FirstName   string
    LastName    string
    PhoneNumber string
    Email       *string
    State       string
    City        string
    Country     string
    Address     string
    Company     string
    OrgID       uint
    Org         Org
    Orders      []Order
}

type Invoice struct {
    BaseModel
    OrderID       uint
    Order         Order
    InvoiceNumber string
    PDFPath       string
    OrgID         uint
}

type Org struct {
    BaseModel
    Name             string
    Logo             string
    OrganizationName string
    OrganizationUrl  string
    Email            string
    Phone            string
    Address          *Address
    OnboardedAt      bool
    Users            []Customer
    Products         []Product
    Customers        []Customer
    Orders           []Order
}

type Product struct {
    BaseModel
    Name        string
    Category    string
    OrgID       uint
    Org         *Org
    ImageURL    string
    Description string
    Variants    []Variant
}

type Variant struct {
    ID        uint
    ProductID uint
    Color     string
    Size      string
    Price     float64
    Stock     int
    SKU       string
    OrgID     uint
}

type User struct {
    BaseModel
    ClerkID   string
    FirstName string
    LastName  string
    Email     string
    Phone     *string
    Role      Role
    OrgID     *uint
    Org       *Org
    Address   *Address
}
```

---

## DTOs

```go
package dto

type CreateProductDTO struct {
    Field string `json:"field" validate:"required,min=2,max=50"`
}

func (p *CreateProductDTO) ToModel(orgID uint) model.Product {
    return model.Product{
        Field: p.Field,
    }
}
```

---

## Repository

```go
package product

type repository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) interfaces.ProductRepository {
    return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, product *model.Product) error {
    return r.db.WithContext(ctx).Create(product).Error
}
```

---

## Service

```go
package product

type service struct {
    repo interfaces.ProductRepository
}

func NewService(repo interfaces.ProductRepository) interfaces.ProductService {
    return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, product *model.Product) error {
    return s.repo.Create(ctx, product)
}
```

---

## Handler

```go
package product

type ProductHandler struct {
    service interfaces.ProductService
    appCtx  *deps.AppContext
}

func NewHandler(service interfaces.ProductService, appCtx *deps.AppContext) interfaces.ProductHandler {
    return &ProductHandler{service: service, appCtx: appCtx}
}

// Filter godoc
// @Summary      List products with filtering and pagination
// @Description  Returns a paginated list of products. Supports filtering, sorting, searching, and preloading.
func (h *ProductHandler) Filter(w http.ResponseWriter, r *http.Request) {
    response.WriteJSONSuccess(w, http.StatusOK, resp, h.appCtx.Logger)
}
```

---

## Response Helpers

```go
package response

func WriteJSONSuccess(w http.ResponseWriter, statusCode int, data any, logger interfaces.Logger) {
    resp := APIResponse{
        Code:    statusCode,
        Message: "success",
        Data:    data,
    }

    w.WriteHeader(statusCode)
    if err := json.NewEncoder(w).Encode(resp); err != nil {
        logger.Error(apperrors.ErrEncodeResponse, "error", err)
    }
}
```

---

## Main Entry

```go
package main

func main() {
    r := chi.NewRouter()
    logger := logger.New(os.Getenv("ENV"))

    cfg, err := config.LoadConfig(logger)
    if err != nil {
        logger.Fatal("failed to load config", "err", err)
    }

    clerk.SetKey(cfg.ClerkSecret)

    dbConn := db.New(db.DBConfig{
        Host:     cfg.DBHost,
        Port:     cfg.DBPort,
        User:     cfg.DBUser,
        Password: cfg.DBPassword,
        Name:     cfg.DBName,
        Env:      cfg.Env,
    }, logger)

    appCtx := &deps.AppContext{
        DB:     dbConn,
        Config: cfg,
        Logger: logger,
        Cache:  nil,
    }

    middlewares := middleware.New()
    routes.Register(r, appCtx, middlewares)

    port := cfg.Port
    if port == 0 {
        port = 8080
    }

    srv := http.Server{
        Addr:    fmt.Sprintf(":%d", port),
        Handler: r,
    }
}
```

---

## Interfaces

```go
package interfaces

type ProductService interface {
    Create(ctx context.Context, org *model.Product) error
}

type ProductRepository interface {
    Create(ctx context.Context, product *model.Product) error
}

type ProductHandler interface {
    Create(w http.ResponseWriter, r *http.Request)
}

type Middleware interface {
    Recover(logger Logger) func(http.Handler) http.Handler
    ValidateJWT(opts ...clerkHttp.AuthorizationOption) func(http.Handler) http.Handler
}
```

---

## API Response Types

```go
type APIResponse struct {
    Code    int    `json:"code" example:"200"`
    Message string `json:"message" example:"success"`
    Data    any    `json:"data,omitempty"`
}

type FilterResponse[T any] struct {
    Items      []T                   `json:"items"`
    Pagination pagination.Pagination `json:"pagination"`
}

type APIErrorResponse struct {
    Code    int    `json:"code" example:"400"`
    Message string `json:"message" example:"invalid request body"`
}

type APIError struct {
    Code        int
    Message     string
    InternalMsg string // detailed message for logs if code is 500
}
```

---

## Routes

```go
package routes

func registerProductRoutes(router chi.Router, productHandler interfaces.ProductHandler) {
    router.Route("/products", func(r chi.Router) {
        r.Get("/", productHandler.Filter)
    })
}

func Register(r chi.Router, appCtx *deps.AppContext, middleware interfaces.Middleware) {
    productRepository := product.NewRepository(appCtx.DB)
    productService := product.NewService(productRepository)
    productHandler := product.NewHandler(productService, appCtx)

    r.Group(func(r chi.Router) {
        r.Use(middleware.ValidateJWT())
        registerProductRoutes(r, productHandler)
    })
}
```


