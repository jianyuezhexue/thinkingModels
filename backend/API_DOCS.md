# API 接口文档

本项目使用 **Swagger/OpenAPI 2.0** 自动生成 API 接口文档。

## 快速开始

### 1. 启动服务

```bash
# 进入后端目录
cd backend

# 启动服务
go run main.go
```

服务启动后，Swagger UI 可通过以下地址访问：
- **Swagger UI**: http://localhost:2500/swagger/index.html
- **Swagger JSON**: http://localhost:2500/swagger/doc.json
- **Swagger YAML**: http://localhost:2500/swagger/doc.yaml

### 2. 生成/更新 API 文档

当修改 API 代码后，需要重新生成 Swagger 文档：

```bash
# 使用 swag 命令生成
swag init --parseDependency --parseInternal

# 或使用 Makefile 命令
make swagger
```

## 文档结构

### 已实现的 API 模块

| 模块 | 路径前缀 | 描述 |
|------|---------|------|
| 用户认证 | `/auth/*` | 用户注册、登录、登出 |
| 用户管理 | `/user/*` | 用户信息管理 |
| 课题管理 | `/subject/topic/*` | 课题的 CRUD、状态管理 |
| 思维模型 | `/thinking/model/*` | 思维模型的创建、发布、管理 |

### 认证方式

API 使用 **JWT Bearer Token** 进行认证。

1. 首先调用 `/auth/login` 获取 Token
2. 在 Swagger UI 点击 **"Authorize"** 按钮
3. 输入格式：`Bearer eyJhbGciOiJIUzI1NiIs...`
4. 点击 **"Authorize"** 确认

### 响应格式

所有 API 返回统一的响应结构：

```json
{
  "code": 0,        // 状态码：0 表示成功，-1 表示失败
  "msg": "success", // 消息描述
  "data": {},       // 响应数据
  "trace_id": ""    // 追踪 ID
}
```

## API 列表

### 用户认证模块

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | `/auth/register` | 用户注册 | 否 |
| POST | `/auth/login` | 用户登录 | 否 |
| POST | `/auth/logout` | 用户登出 | 是 |
| GET | `/user/info` | 获取当前用户信息 | 是 |
| POST | `/user/list` | 查询用户列表 | 是 |

### 课题管理模块

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | `/subject/topic` | 创建课题 | 是 |
| PUT | `/subject/topic` | 更新课题 | 是 |
| GET | `/subject/topic/:id` | 获取课题详情 | 是 |
| GET | `/subject/topic/list` | 查询课题列表 | 是 |
| DELETE | `/subject/topic` | 删除课题 | 是 |
| POST | `/subject/topic/status` | 更新课题状态 | 是 |
| POST | `/subject/topic/complete` | 完成课题 | 是 |
| POST | `/subject/topic/archive` | 归档课题 | 是 |

### 思维模型模块

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | `/thinking/model` | 创建思维模型 | 是 |
| PUT | `/thinking/model` | 更新思维模型 | 是 |
| GET | `/thinking/model/:id` | 获取模型详情 | 是 |
| GET | `/thinking/model/list` | 查询模型列表 | 否 |
| GET | `/thinking/model/my` | 获取我的模型 | 是 |
| POST | `/thinking/model/publish` | 发布模型 | 是 |
| POST | `/thinking/model/unpublish/:id` | 下架模型 | 是 |

## 开发规范

### 添加新 API 文档

在 handler 函数上添加 Swagger 注释：

```go
// Create 创建思维模型
// @Summary 创建思维模型
// @Description 创建一个新的思维模型
// @Tags 思维模型
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.CreateModel true "请求参数"
// @Success 200 {object} api.Response{data=model.ModelDetail} "创建成功"
// @Failure 400 {object} api.Response "参数错误"
// @Router /thinking/model [post]
func (a *Model) Create(ctx *gin.Context) {
    // 实现代码...
}
```

### 常用注释说明

| 注释 | 说明 |
|------|------|
| `@Summary` | API 简短摘要 |
| `@Description` | API 详细描述 |
| `@Tags` | API 分类标签 |
| `@Accept` | 请求 Content-Type |
| `@Produce` | 响应 Content-Type |
| `@Security` | 认证方式 |
| `@Param` | 请求参数 |
| `@Success` | 成功响应 |
| `@Failure` | 失败响应 |
| `@Router` | 路由路径和方法 |

## 常用命令

```bash
# 生成 Swagger 文档
make swagger

# 生成并自动打开浏览器（如已配置）
make swagger-serve

# 清理生成的文档
make swagger-clean

# 完整开发环境初始化
make dev-setup
```

## 参考资料

- [Swaggo GitHub](https://github.com/swaggo/swag)
- [Gin Swagger](https://github.com/swaggo/gin-swagger)
- [OpenAPI 2.0 Specification](https://swagger.io/specification/v2/)

## 注意事项

1. 修改 API 后务必重新生成文档
2. 确保所有请求/响应结构体都有完整的字段标签
3. 使用 `@Security Bearer` 标记需要认证的接口
4. 枚举类型在描述中说明可选值
