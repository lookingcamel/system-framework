# HTTP 处理器模块

提供 HTTP 请求处理的示例和最佳实践。

## 功能特性

- 📝 **RESTful API**: 标准 REST 接口设计
- 🔄 **统一响应**: 统一的响应格式
- ✅ **参数验证**: 请求参数自动验证
- 🏥 **健康检查**: 存活和就绪探针
- 📊 **分页支持**: 列表查询分页
- 🛡️ **错误处理**: 统一错误处理

## 快速开始

### 创建处理器

```go
import "github.com/lookingcamel/system-framework/internal/handler"

func main() {
    exampleHandler := handler.NewExampleHandler()

    v1 := engine.Group("/api/v1")
    {
        v1.GET("/example/:id", exampleHandler.GetExample)
        v1.GET("/examples", exampleHandler.ListExample)
        v1.POST("/example", exampleHandler.CreateExample)
        v1.PUT("/example/:id", exampleHandler.UpdateExample)
        v1.DELETE("/example/:id", exampleHandler.DeleteExample)
    }
}
```

## 响应格式

### 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "example"
  }
}
```

### 错误响应

```json
{
  "code": -1,
  "message": "error message",
  "data": null
}
```

### 分页响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [...],
    "total": 100,
    "page": 1,
    "page_size": 10,
    "total_pages": 10
  }
}
```

## 处理器示例

### 1. 健康检查处理器

```go
type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
    return &HealthHandler{}
}

func (h *HealthHandler) Health(c *gin.Context) {
    c.JSON(200, gin.H{
        "code":    0,
        "message": "healthy",
        "data": gin.H{
            "status": "up",
            "time":   time.Now().Unix(),
        },
    })
}

func (h *HealthHandler) Ready(c *gin.Context) {
    // 检查数据库连接
    if err := database.Ping(); err != nil {
        c.JSON(503, gin.H{
            "code":    -1,
            "message": "not ready: database unavailable",
        })
        return
    }

    c.JSON(200, gin.H{
        "code":    0,
        "message": "ready",
    })
}
```

### 2. CRUD 处理器

```go
type ExampleHandler struct{}

func NewExampleHandler() *ExampleHandler {
    return &ExampleHandler{}
}

// 获取单个
func (h *ExampleHandler) GetExample(c *gin.Context) {
    id := c.Param("id")

    example, err := getExampleByID(id)
    if err != nil {
        c.JSON(404, gin.H{
            "code":    -1,
            "message": "example not found",
        })
        return
    }

    c.JSON(200, gin.H{
        "code":    0,
        "message": "success",
        "data":    example,
    })
}

// 列表查询
func (h *ExampleHandler) ListExample(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

    examples, total, err := listExamples(page, pageSize)
    if err != nil {
        c.JSON(500, gin.H{
            "code":    -1,
            "message": err.Error(),
        })
        return
    }

    c.JSON(200, gin.H{
        "code":    0,
        "message": "success",
        "data": gin.H{
            "items":      examples,
            "total":      total,
            "page":       page,
            "page_size":  pageSize,
            "total_pages": (total + pageSize - 1) / pageSize,
        },
    })
}

// 创建
func (h *ExampleHandler) CreateExample(c *gin.Context) {
    var req CreateExampleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{
            "code":    -1,
            "message": "invalid request: " + err.Error(),
        })
        return
    }

    example, err := createExample(&req)
    if err != nil {
        c.JSON(500, gin.H{
            "code":    -1,
            "message": err.Error(),
        })
        return
    }

    c.JSON(201, gin.H{
        "code":    0,
        "message": "success",
        "data":    example,
    })
}

// 更新
func (h *ExampleHandler) UpdateExample(c *gin.Context) {
    id := c.Param("id")

    var req UpdateExampleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{
            "code":    -1,
            "message": "invalid request: " + err.Error(),
        })
        return
    }

    example, err := updateExample(id, &req)
    if err != nil {
        c.JSON(500, gin.H{
            "code":    -1,
            "message": err.Error(),
        })
        return
    }

    c.JSON(200, gin.H{
        "code":    0,
        "message": "success",
        "data":    example,
    })
}

// 删除
func (h *ExampleHandler) DeleteExample(c *gin.Context) {
    id := c.Param("id")

    if err := deleteExample(id); err != nil {
        c.JSON(500, gin.H{
            "code":    -1,
            "message": err.Error(),
        })
        return
    }

    c.JSON(200, gin.H{
        "code":    0,
        "message": "success",
    })
}
```

## 请求模型

```go
type CreateExampleRequest struct {
    Name        string `json:"name" binding:"required,min=1,max=100"`
    Description string `json:"description" binding:"max=500"`
    Status      string `json:"status" binding:"required,oneof=active inactive"`
}

type UpdateExampleRequest struct {
    Name        string `json:"name" binding:"omitempty,min=1,max=100"`
    Description string `json:"description" binding:"max=500"`
    Status      string `json:"status" binding:"omitempty,oneof=active inactive"`
}

type ListQuery struct {
    Page     int    `form:"page" binding:"min=1"`
    PageSize int    `form:"page_size" binding:"min=1,max=100"`
    Keyword  string `form:"keyword"`
    Sort     string `form:"sort" binding:"omitempty,oneof=created_at updated_at"`
    Order    string `form:"order" binding:"omitempty,oneof=asc desc"`
}
```

## 错误处理

```go
type AppError struct {
    Code    int
    Message string
    Cause   error
}

func (e *AppError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Cause)
    }
    return e.Message
}

func HandleError(c *gin.Context, err error) {
    if appErr, ok := err.(*AppError); ok {
        c.JSON(appErr.Code, gin.H{
            "code":    -1,
            "message": appErr.Message,
        })
        return
    }

    c.JSON(500, gin.H{
        "code":    -1,
        "message": "internal server error",
    })
}
```

## 中间件组合

```go
v1 := engine.Group("/api/v1")

// 需要认证的接口
authorized := v1.Group("")
authorized.Use(auth.JWT())
{
    authorized.GET("/example/:id", exampleHandler.GetExample)
    authorized.POST("/example", exampleHandler.CreateExample)
    authorized.PUT("/example/:id", exampleHandler.UpdateExample)
    authorized.DELETE("/example/:id", exampleHandler.DeleteExample)
}

// 公开接口
v1.GET("/examples", exampleHandler.ListExample)
v1.GET("/health", healthHandler.Health)
```

## 路由组织

### 按功能分组

```
/api/v1
├── users/           # 用户相关
│   ├── GET    /          # 列表
│   ├── POST   /          # 创建
│   ├── GET    /:id       # 获取
│   ├── PUT    /:id       # 更新
│   └── DELETE /:id       # 删除
├── orders/          # 订单相关
│   ├── GET    /          # 列表
│   ├── POST   /          # 创建
│   └── GET    /:id       # 获取
└── products/        # 产品相关
    ├── GET    /          # 列表
    ├── POST   /          # 创建
    └── GET    /:id       # 获取
```

### 按版本分组

```
/api/v1/users     # v1 API
/api/v2/users     # v2 API
```

## 测试处理器

```go
func TestGetExample(t *testing.T) {
    engine := gin.New()
    handler := NewExampleHandler()
    engine.GET("/example/:id", handler.GetExample)

    req, _ := http.NewRequest("GET", "/example/1", nil)
    w := httptest.NewRecorder()
    engine.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Code)
}

func TestCreateExample(t *testing.T) {
    engine := gin.New()
    handler := NewExampleHandler()
    engine.POST("/example", handler.CreateExample)

    body := `{"name": "test", "status": "active"}`
    req, _ := http.NewRequest("POST", "/example", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    engine.ServeHTTP(w, req)

    assert.Equal(t, 201, w.Code)
}
```

## 最佳实践

1. **RESTful 设计**
   - 使用标准 HTTP 方法
   - 使用名词表示资源
   - 使用复数形式

2. **参数验证**
   - 使用 binding tag
   - 返回详细错误信息
   - 拒绝非法输入

3. **错误处理**
   - 统一错误格式
   - 区分客户端和服务器错误
   - 记录详细日志

4. **性能优化**
   - 使用分页查询
   - 避免 N+1 查询
   - 合理使用缓存

## 相关文档

- [中间件模块](../middleware/README.md)
- [认证模块](../auth/README.md)
- [数据库模块](../database/README.md)
- [项目总览](../../README.md)
