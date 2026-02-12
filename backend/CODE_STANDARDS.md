# Backend Code Standards

## Project Overview

This is a Go backend project using the Gin framework with a simplified layered architecture.

## Architecture

The project follows a **simplified layered architecture** (not DDD):

```
backend/
├── api/              # API Layer (Controllers)
├── domain/           # Domain Layer (Models)
│   └── {domain}/
│       └── {entity}/
│           ├── model.go      # Business model
│           ├── types.go      # DTOs (Request/Response)
│           ├── abilitys.go   # Ability interfaces
│           └── model_test.go # Unit tests
├── logic/            # Business Logic Layer
├── router/           # Router Layer
├── component/        # Common Components (DB, Redis)
├── config/           # Configuration
├── middleware/       # Middleware
└── tool/             # Utilities
```

## Layer Responsibilities

| Layer | Path | Responsibility |
|-------|------|----------------|
| **API** | `api/{domain}/{entity}.go` | HTTP request handling, parameter binding, response formatting |
| **Router** | `router/v1.go` | Route registration and grouping |
| **Logic** | `logic/{domain}/{entity}.go` | Business logic orchestration |
| **Domain** | `domain/{domain}/{entity}/` | Data models, validation, business rules |
| **Component** | `component/` | Shared components (database, cache) |
| **Tool** | `tool/` | Utility functions |

## Naming Conventions

| Type | Pattern | Example |
|------|---------|---------|
| Entity | `{Entity}Entity` | `UserEntity`, `TopicEntity` |
| Logic | `{Entity}Logic` | `UserLogic`, `TopicLogic` |
| API Controller | `{Entity}` | `User`, `Topic` |
| DTO Create | `Create{Entity}` | `CreateUser`, `CreateTopic` |
| DTO Update | `Update{Entity}` | `UpdateUser`, `UpdateTopic` |
| DTO Search | `Search{Entity}` | `SearchUser`, `SearchTopic` |
| DTO Delete | `Del{Entity}` | `DelUser`, `DelTopic` |

## File Templates

### 1. domain/{entity}/model.go

```go
package {entity}

import (
    "github.com/gin-gonic/gin"
    "github.com/jianyuezhexue/base"
    "thinkingModels/component/db"
)

// {Entity}Entity business model
type {Entity}Entity struct {
    base.BaseModel[{Entity}Entity]
    // Add fields here
    Name string `json:"name" gorm:"column:name"`
}

// New{Entity}Entity creates a new entity instance
func New{Entity}Entity(ctx *gin.Context) *{Entity}Entity {
    entity := &{Entity}Entity{}
    entity.BaseModel = base.NewBaseModel(ctx, db.InitDb(), entity.TableName(), entity)
    return entity
}

// TableName returns the database table name
func (m *{Entity}Entity) TableName() string {
    return "{table_name}"
}

// Validate validates the entity data
func (m *{Entity}Entity) Validate() error {
    // Add validation logic
    return nil
}
```

### 2. domain/{entity}/types.go

```go
package {entity}

// Create{Entity} DTO for creating entity
type Create{Entity} struct {
    Name string `json:"name" binding:"required"`
}

// Update{Entity} DTO for updating entity
type Update{Entity} struct {
    Id   uint64 `json:"id" binding:"required"`
    Name string `json:"name"`
}

// Search{Entity} DTO for searching entities
type Search{Entity} struct {
    Page     int64  `json:"page" form:"page" search:"page"`
    PageSize int64  `json:"pageSize" form:"pageSize" search:"pageSize"`
    Keyword  string `json:"keyword" form:"keyword" search:"like"`
}

// Del{Entity} DTO for deleting entities
type Del{Entity} struct {
    Ids []uint64 `json:"ids" binding:"required"`
}

// ListResp pagination response
type ListResp struct {
    Page     int64              `json:"page"`
    PageSize int64              `json:"pageSize"`
    Total    int64              `json:"total"`
    List     []*{Entity}Entity `json:"list"`
}
```

### 3. domain/{entity}/abilitys.go

```go
package {entity}

// {Entity}Ability defines entity capabilities
type {Entity}Ability interface {
    // Define ability methods
    CanEdit(userId uint64) bool
    CanDelete(userId uint64) bool
}
```

### 4. domain/{entity}/model_test.go

```go
package {entity}

import (
    "testing"
)

func Test{Entity}_Basic(t *testing.T) {
    // Unit tests
    entity := &{Entity}Entity{}
    if entity.TableName() != "{table_name}" {
        t.Errorf("Expected table name {table_name}, got %s", entity.TableName())
    }
}
```

### 5. logic/{domain}/{entity}.go

```go
package {domain}

import (
    "github.com/gin-gonic/gin"
    "thinkingModels/domain/{domain}/{entity}"
    "thinkingModels/logic"
)

// {Entity}Logic business logic structure
type {Entity}Logic struct {
    logic.BaseLogic
}

// New{Entity}Logic creates a new logic instance
func New{Entity}Logic(ctx *gin.Context) *{Entity}Logic {
    return &{Entity}Logic{BaseLogic: logic.BaseLogic{Ctx: ctx}}
}

// Create creates a new entity
func (l *{Entity}Logic) Create(req *{entity}.Create{Entity}) (*{entity}.{Entity}Entity, error) {
    entity := {entity}.New{Entity}Entity(l.Ctx)

    _, err := entity.SetData(req)
    if err != nil {
        return nil, err
    }

    err = entity.Validate()
    if err != nil {
        return nil, err
    }

    res, err := entity.Create()
    if err != nil {
        return nil, err
    }

    return res, nil
}

// Update updates an entity
func (l *{Entity}Logic) Update(req *{entity}.Update{Entity}) (*{entity}.{Entity}Entity, error) {
    entity := {entity}.New{Entity}Entity(l.Ctx)

    _, err := entity.LoadById(req.Id)
    if err != nil {
        return nil, err
    }

    _, err = entity.SetData(req)
    if err != nil {
        return nil, err
    }

    err = entity.Validate()
    if err != nil {
        return nil, err
    }

    res, err := entity.Update()
    return res, err
}

// Get retrieves a single entity
func (l *{Entity}Logic) Get(req *{entity}.Search{Entity}) (*{entity}.{Entity}Entity, error) {
    entity := {entity}.New{Entity}Entity(l.Ctx)
    cond := entity.MakeConditon(*req)
    res, err := entity.LoadData(cond)
    return res, err
}

// List retrieves a paginated list
func (l *{Entity}Logic) List(req *{entity}.Search{Entity}) (*{entity}.ListResp, error) {
    entity := {entity}.New{Entity}Entity(l.Ctx)
    cond := entity.MakeConditon(*req)

    total, err := entity.Count(cond)
    if err != nil {
        return nil, err
    }

    list, err := entity.List(cond)
    if err != nil {
        return nil, err
    }

    res := &{entity}.ListResp{
        Page:     req.Page,
        PageSize: req.PageSize,
        Total:    total,
        List:     list,
    }
    return res, nil
}

// Del deletes entities
func (l *{Entity}Logic) Del(req *{entity}.Del{Entity}) (any, error) {
    entity := {entity}.New{Entity}Entity(l.Ctx)
    err := entity.Del(req.Ids...)
    return nil, err
}
```

### 6. api/{domain}/{entity}.go

```go
package {domain}

import (
    "github.com/gin-gonic/gin"
    "thinkingModels/api"
    "thinkingModels/domain/{domain}/{entity}"
    "thinkingModels/logic/{domain}"
)

// {Entity} API controller
type {Entity} struct {
    api.Base
}

// New{Entity} creates a new API controller
func New{Entity}() *{Entity} {
    return &{Entity}{}
}

// Create handles POST /{domain}/{entity}
func (a {Entity}) Create(ctx *gin.Context) {
    req := &{entity}.Create{Entity}{}
    err := a.Bind(ctx, req)
    if err != nil {
        a.Error(err)
        return
    }

    logic := {domain}.New{Entity}Logic(ctx)
    res, err := logic.Create(req)
    if err != nil {
        a.Error(err)
        return
    }

    a.Success(res, "Created successfully")
}

// Update handles PUT /{domain}/{entity}
func (a {Entity}) Update(ctx *gin.Context) {
    req := &{entity}.Update{Entity}{}
    err := a.Bind(ctx, req)
    if err != nil {
        a.Error(err)
        return
    }

    logic := {domain}.New{Entity}Logic(ctx)
    res, err := logic.Update(req)
    if err != nil {
        a.Error(err)
        return
    }

    a.Success(res, "Updated successfully")
}

// Get handles GET /{domain}/{entity}/:id
func (a {Entity}) Get(ctx *gin.Context) {
    req := &{entity}.Search{Entity}{}
    err := a.Bind(ctx, req)
    if err != nil {
        a.Error(err)
        return
    }

    logic := {domain}.New{Entity}Logic(ctx)
    res, err := logic.Get(req)
    if err != nil {
        a.Error(err)
        return
    }

    a.Success(res, "Retrieved successfully")
}

// List handles POST /{domain}/{entity}/list
func (a {Entity}) List(ctx *gin.Context) {
    req := &{entity}.Search{Entity}{}
    err := a.Bind(ctx, req)
    if err != nil {
        a.Error(err)
        return
    }

    logic := {domain}.New{Entity}Logic(ctx)
    res, err := logic.List(req)
    if err != nil {
        a.Error(err)
        return
    }

    a.Success(res, "List retrieved successfully")
}

// Del handles DELETE /{domain}/{entity}
func (a {Entity}) Del(ctx *gin.Context) {
    req := &{entity}.Del{Entity}{}
    err := a.Bind(ctx, req)
    if err != nil {
        a.Error(err)
        return
    }

    logic := {domain}.New{Entity}Logic(ctx)
    res, err := logic.Del(req)
    if err != nil {
        a.Error(err)
        return
    }

    a.Success(res, "Deleted successfully")
}
```

### 7. router/v1.go

```go
package router

import (
    "github.com/gin-gonic/gin"
    "thinkingModels/api/{domain}"
)

// AuthorizedRouters registers authorized routes
func AuthorizedRouters() {
    authorizedRouters := func(router *gin.Engine) {
        v1 := router.Group("/v1")

        // {domain} routes
        {entity}Api := {domain}.New{Entity}()
        {entity}Group := v1.Group("/{domain}/{entity}")
        {
            {entity}Group.POST("", {entity}Api.Create)
            {entity}Group.PUT("", {entity}Api.Update)
            {entity}Group.GET("/:id", {entity}Api.Get)
            {entity}Group.POST("/list", {entity}Api.List)
            {entity}Group.DELETE("", {entity}Api.Del)
        }
    }

    Routers = append(Routers, authorizedRouters)
}
```

## Design Principles

1. **One domain, one router**: Each domain registers its routes independently
2. **One domain, one API**: Each domain has one API file
3. **One domain, one Logic**: Each domain has one Logic file
4. **One table, one directory**: Each table has its own directory under domain
5. **Ability definition**: Define capabilities in `abilitys.go`, implement in `model.go`
6. **Logic orchestration**: Logic layer calls different model abilities
7. **Clear layering**: API handles HTTP, Logic handles business, Domain handles rules

## Adding a New Module

1. Create directory: `domain/{domain}/{entity}/`
2. Create files: `model.go`, `types.go`, `abilitys.go`, `model_test.go`
3. Create logic: `logic/{domain}/{entity}.go`
4. Create API: `api/{domain}/{entity}.go`
5. Register routes in `router/v1.go`

## Best Practices

- Use `base.BaseModel` for common CRUD operations
- Define validation in `Validate()` method
- Use DTOs for all API inputs/outputs
- Handle errors at each layer appropriately
- Write unit tests for model logic
- Keep functions small and focused
- Use meaningful variable names
- Add comments for complex business logic
