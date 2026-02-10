# 用户领域开发需求文档

## 1. 概述

### 1.1 目标
构建完整的用户领域模块，包含用户管理、角色权限、企业管理等核心功能，为平台提供统一的认证授权体系。

### 1.2 范围
第一阶段（本期实现）：
- 用户基础功能（登录/登出/注册）
- JWT Token 认证机制
- 用户上下文中间件
- 密码加密存储

第二阶段（后续迭代）：
- 角色权限管理（RBAC）
- 企业/租户管理
- 菜单权限控制
- 数据权限隔离

---

## 2. 领域结构设计

### 2.1 领域层级
```
domain/user/
├── user/           # 用户实体
├── role/           # 角色实体
├── permission/     # 权限实体
├── enterprise/     # 企业/租户实体
└── auth/           # 认证相关（JWT/Token）
```

### 2.2 实体关系
```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│ Enterprise  │────<│    User     │>────│    Role     │
│  (企业/租户)  │ 1:n │   (用户)     │ n:m │   (角色)     │
└─────────────┘     └─────────────┘     └─────────────┘
                                               │
                                               │ n:m
                                               ▼
                                        ┌─────────────┐
                                        │ Permission  │
                                        │   (权限)     │
                                        └─────────────┘
```

---

## 3. 功能需求

### 3.1 用户实体 (User)

#### 3.1.1 字段定义
| 字段名 | 类型 | 说明 | 约束 |
|--------|------|------|------|
| id | uint64 | 主键ID | 自增 |
| username | string | 登录用户名 | 唯一，必填，3-20字符 |
| password | string | 加密密码 | 必填，bcrypt加密 |
| nickname | string | 用户昵称 | 可选，最多50字符 |
| email | string | 邮箱 | 可选，格式校验 |
| phone | string | 手机号 | 可选，格式校验 |
| avatar | string | 头像URL | 可选 |
| status | int | 状态 | 0=禁用, 1=正常 |
| last_login_time | datetime | 最后登录时间 | 自动更新 |
| last_login_ip | string | 最后登录IP | 自动记录 |
| enterprise_id | uint64 | 所属企业ID | 外键关联 |
| created_at | datetime | 创建时间 | 默认CURRENT_TIMESTAMP，非空 |
| updated_at | datetime | 更新时间 | 默认CURRENT_TIMESTAMP |
| deleted_at | datetime | 删除时间 | 软删除标记，可空 |
| create_by | uint64 | 创建人ID | 默认0，非空 |
| create_by_name | string | 创建人姓名 | 默认'系统'，非空，20字符 |
| update_by | uint64 | 更新人ID | 默认0，非空 |
| update_by_name | string | 更新人姓名 | 默认'系统'，非空，20字符 |

#### 3.1.2 功能列表
- [ ] 用户注册
- [ ] 用户登录（返回JWT Token）
- [ ] 用户登出（Token黑名单/失效）
- [ ] 获取当前用户信息
- [ ] 修改密码
- [ ] 更新用户信息
- [ ] 用户列表查询（分页）

### 3.2 认证机制 (Auth)

#### 3.2.1 JWT Token 设计
```
Header: {
  "alg": "HS256",
  "typ": "JWT"
}

Payload: {
  "sub": "user_id",           // 用户ID
  "username": "zhangsan",     // 用户名
  "enterprise_id": 1,         // 企业ID（多租户支持）
  "roles": ["admin", "user"], // 角色列表
  "iat": 1700000000,          // 签发时间
  "exp": 1700003600,          // 过期时间（默认1小时）
  "jti": "unique_token_id"    // Token唯一标识（用于登出）
}
```

#### 3.2.2 Token 策略
| 类型 | 有效期 | 用途 | 存储位置 |
|------|--------|------|----------|
| Access Token | 1小时 | API访问认证 | 前端内存/LocalStorage |
| Refresh Token | 7天 | 刷新Access Token | HttpOnly Cookie |

#### 3.2.3 功能列表
- [ ] JWT Token 生成
- [ ] JWT Token 验证
- [ ] Token 刷新机制
- [ ] Token 黑名单（登出失效）

### 3.3 中间件 (Middleware)

#### 3.3.1 用户上下文中间件
功能：解析Token，将用户信息注入Gin上下文
```go
// 用户信息上下文结构
type UserContext struct {
    UserID       uint64   // 用户ID
    Username     string   // 用户名
    EnterpriseID uint64   // 企业ID
    Roles        []string // 角色列表
    IsAdmin      bool     // 是否管理员
}
```

#### 3.3.2 认证中间件
- 验证请求头中的 Authorization: Bearer <token>
- Token无效/过期时返回401
- 将解析后的用户信息存入 *gin.Context

#### 3.3.3 功能列表
- [ ] JWT认证中间件
- [ ] 用户上下文注入
- [ ] 登录状态检查

---

## 4. 技术实现方案

### 4.1 目录结构
```
backend/
├── api/user/
│   ├── user.go          # 用户API接口
│   └── auth.go          # 认证API接口（登录/登出）
├── domain/user/
│   └── user/
│       ├── model.go     # 用户实体模型（包含JWT生成/密码验证等充血方法）
│       ├── types.go     # 类型定义（入参DTO、返回DTO、搜索条件等）
│       └── ability.go   # 领域接口（保持简洁）
├── logic/user/
│   ├── user.go          # 用户业务逻辑
│   └── auth.go          # 认证业务逻辑
├── middleware/
│   └── auth.go          # 认证中间件
└── router/
    └── user.go          # 用户路由注册
```

### 4.2 接口设计

#### 4.2.1 认证接口
| 方法 | 路径 | 说明 | 是否需要认证 |
|------|------|------|-------------|
| POST | /v1/auth/login | 用户登录 | 否 |
| POST | /v1/auth/logout | 用户登出 | 是 |
| POST | /v1/auth/refresh | 刷新Token | 否（需Refresh Token）|
| GET | /v1/auth/captcha | 获取验证码 | 否 |

#### 4.2.2 用户接口
| 方法 | 路径 | 说明 | 是否需要认证 |
|------|------|------|-------------|
| POST | /v1/user | 创建用户（注册） | 否（或管理员） |
| GET | /v1/user/profile | 获取当前用户信息 | 是 |
| PUT | /v1/user/profile | 更新当前用户信息 | 是 |
| PUT | /v1/user/password | 修改密码 | 是 |
| GET | /v1/users | 用户列表（分页） | 是（管理员） |
| GET | /v1/user/:id | 用户详情 | 是 |
| PUT | /v1/user/:id | 更新用户 | 是（管理员或本人） |
| DELETE | /v1/user/:id | 删除用户 | 是（管理员） |

### 4.3 数据库表结构

#### 用户表 (users)
```sql
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '登录用户名',
    password VARCHAR(255) NOT NULL COMMENT '加密密码',
    nickname VARCHAR(100) DEFAULT '' COMMENT '用户昵称',
    email VARCHAR(100) DEFAULT '' COMMENT '邮箱',
    phone VARCHAR(20) DEFAULT '' COMMENT '手机号',
    avatar VARCHAR(255) DEFAULT '' COMMENT '头像URL',
    status TINYINT DEFAULT 1 COMMENT '状态: 0=禁用, 1=正常',
    last_login_time DATETIME NULL COMMENT '最后登录时间',
    last_login_ip VARCHAR(50) DEFAULT '' COMMENT '最后登录IP',
    enterprise_id BIGINT UNSIGNED DEFAULT 0 COMMENT '企业ID',
    -- 标准审计字段（所有表必须包含）
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NULL COMMENT '修改时间',
    deleted_at DATETIME NULL COMMENT '软删除时间',
    create_by BIGINT UNSIGNED DEFAULT 0 NOT NULL COMMENT '创建人ID',
    create_by_name VARCHAR(20) DEFAULT '系统' NOT NULL COMMENT '创建人姓名',
    update_by BIGINT UNSIGNED DEFAULT 0 NOT NULL COMMENT '修改人ID',
    update_by_name VARCHAR(20) DEFAULT '系统' NOT NULL COMMENT '修改人姓名',
    INDEX idx_username (username),
    INDEX idx_enterprise (enterprise_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
```

#### Token黑名单表 (token_blacklist) - 用于登出
```sql
CREATE TABLE token_blacklist (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    jti VARCHAR(100) NOT NULL UNIQUE COMMENT 'Token唯一标识',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    expires_at DATETIME NOT NULL COMMENT 'Token过期时间',
    -- 标准审计字段（所有表必须包含）
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NULL COMMENT '修改时间',
    deleted_at DATETIME NULL COMMENT '软删除时间',
    create_by BIGINT UNSIGNED DEFAULT 0 NOT NULL COMMENT '创建人ID',
    create_by_name VARCHAR(20) DEFAULT '系统' NOT NULL COMMENT '创建人姓名',
    update_by BIGINT UNSIGNED DEFAULT 0 NOT NULL COMMENT '修改人ID',
    update_by_name VARCHAR(20) DEFAULT '系统' NOT NULL COMMENT '修改人姓名',
    INDEX idx_jti (jti),
    INDEX idx_expires (expires_at),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Token黑名单';
```

### 4.4 核心代码实现

#### 4.4.1 用户实体模型 (domain/user/user/model.go)
> **充血模型设计**：JWT相关能力作为User实体的方法，而非独立工具类

```go
package user

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
    "github.com/jianyuezhexue/base"
    "thinkingModels/component/db"
    "thinkingModels/config"
)

// ========== 实体定义 ==========

type User struct {
    base.BaseModel[User]
    Username      string `json:"username" type:"db" comment:"登录用户名"`
    Password      string `json:"-" type:"db" comment:"加密密码"` // json:"-" 不序列化
    Nickname      string `json:"nickname" type:"db" comment:"用户昵称"`
    Email         string `json:"email" type:"db" comment:"邮箱"`
    Phone         string `json:"phone" type:"db" comment:"手机号"`
    Avatar        string `json:"avatar" type:"db" comment:"头像URL"`
    Status        int    `json:"status" type:"db" comment:"状态:0=禁用,1=正常"`
    LastLoginTime db.LocalTime `json:"lastLoginTime" type:"db" comment:"最后登录时间"`
    LastLoginIP   string `json:"lastLoginIp" type:"db" comment:"最后登录IP"`
    EnterpriseID  uint64 `json:"enterpriseId" type:"db" comment:"企业ID"`
}

func (u *User) TableName() string {
    return "users"
}

// ========== 认证相关方法（充血模型）==========

// UserClaims JWT Claims定义
type UserClaims struct {
    UserID       uint64   `json:"sub"`
    Username     string   `json:"username"`
    EnterpriseID uint64   `json:"enterprise_id"`
    Roles        []string `json:"roles"`
    jwt.RegisteredClaims
}

// TokenPair Token对
type TokenPair struct {
    AccessToken  string
    RefreshToken string
    ExpiresIn    int64 // 过期时间（秒）
}

// GenerateToken 为用户生成JWT Token（实体方法）
func (u *User) GenerateToken() (*TokenPair, error) {
    cfg := config.Config.JWT
    now := time.Now()

    // Access Token
    accessClaims := UserClaims{
        UserID:       u.ID,
        Username:     u.Username,
        EnterpriseID: u.EnterpriseID,
        Roles:        []string{}, // TODO: 从角色关联加载
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(cfg.AccessExpiry) * time.Second)),
            IssuedAt:  jwt.NewNumericDate(now),
            NotBefore: jwt.NewNumericDate(now),
            Issuer:    cfg.Issuer,
            ID:        uuid.New().String(), // JTI
        },
    }
    accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(cfg.Secret))
    if err != nil {
        return nil, err
    }

    // Refresh Token
    refreshClaims := jwt.RegisteredClaims{
        ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(cfg.RefreshExpiry) * time.Second)),
        IssuedAt:  jwt.NewNumericDate(now),
        ID:        uuid.New().String(),
    }
    refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(cfg.Secret))
    if err != nil {
        return nil, err
    }

    return &TokenPair{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresIn:    cfg.AccessExpiry,
    }, nil
}

// VerifyPassword 验证密码（bcrypt比对）
func (u *User) VerifyPassword(plainPassword string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword))
    return err == nil
}

// HashPassword 密码加密（bcrypt）
func (u *User) HashPassword() error {
    hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    u.Password = string(hashed)
    return nil
}

// UpdateLoginInfo 更新登录信息
func (u *User) UpdateLoginInfo(ip string) {
    u.LastLoginTime = db.LocalTime(time.Now())
    u.LastLoginIP = ip
}

// ToUserInfo 转换为脱敏的用户信息DTO
func (u *User) ToUserInfo() *UserInfo {
    return &UserInfo{
        ID:            u.ID,
        Username:      u.Username,
        Nickname:      u.Nickname,
        Email:         u.Email,
        Phone:         u.Phone,
        Avatar:        u.Avatar,
        Status:        u.Status,
        LastLoginTime: u.LastLoginTime.String(),
        CreatedAt:     u.CreatedAt.String(),
    }
}

// ========== 静态方法/包级函数 ==========

// ParseToken 解析JWT Token（包级函数，不依赖实体）
func ParseToken(tokenString string) (*UserClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(config.Config.JWT.Secret), nil
    })
    if err != nil {
        return nil, err
    }
    if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
        return claims, nil
    }
    return nil, errors.New("invalid token")
}

// IsTokenBlacklisted 检查Token是否在黑名单（包级函数）
func IsTokenBlacklisted(jti string, db *gorm.DB) bool {
    var count int64
    db.Table("token_blacklist").Where("jti = ?", jti).Count(&count)
    return count > 0
}

// BlacklistToken 将Token加入黑名单（登出）
func BlacklistToken(jti string, userID uint64, expiresAt time.Time, db *gorm.DB) error {
    return db.Table("token_blacklist").Create(map[string]interface{}{
        "jti":        jti,
        "user_id":    userID,
        "expires_at": expiresAt,
    }).Error
}
func (j *JWTUtil) ParseToken(tokenString string) (*UserClaims, error)

// ValidateToken 验证Token是否有效（检查黑名单）
func (j *JWTUtil) ValidateToken(tokenString string) (*UserClaims, error)

// BlacklistToken 将Token加入黑名单（登出用）
func (j *JWTUtil) BlacklistToken(jti string, userID uint64, expiresAt time.Time) error
```

#### 4.4.2 类型定义 (domain/user/user/types.go)
> 除JWT相关的Claims和TokenPair已在model.go定义外，其他请求/响应类型在此定义
> 所有入参DTO、返回DTO、搜索条件均在 domain/types.go 中定义

```go
package user

// ==================== 请求DTO ====================

// LoginRequest 登录请求
type LoginRequest struct {
    Username string `json:"username" binding:"required"` // 用户名
    Password string `json:"password" binding:"required"` // 密码
    Captcha  string `json:"captcha"`                     // 验证码（可选）
}

// RegisterRequest 注册请求
type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=20"`  // 用户名
    Password string `json:"password" binding:"required,min=8"`         // 密码
    Nickname string `json:"nickname" binding:"max=50"`                 // 昵称
    Email    string `json:"email" binding:"omitempty,email"`           // 邮箱
    Phone    string `json:"phone" binding:"omitempty,len=11,numeric"`  // 手机号
}

// UpdateUserRequest 更新用户信息请求
type UpdateUserRequest struct {
    ID       uint64 `json:"id" binding:"required"`
    Nickname string `json:"nickname" binding:"max=50"`
    Email    string `json:"email" binding:"omitempty,email"`
    Phone    string `json:"phone" binding:"omitempty,len=11,numeric"`
    Avatar   string `json:"avatar" binding:"max=255"`
}

// UpdatePasswordRequest 修改密码请求
type UpdatePasswordRequest struct {
    OldPassword string `json:"oldPassword" binding:"required"`
    NewPassword string `json:"newPassword" binding:"required,min=8"`
}

// ==================== 响应DTO ====================

// LoginResponse 登录响应
type LoginResponse struct {
    AccessToken  string   `json:"accessToken"`  // Access Token
    RefreshToken string   `json:"refreshToken"` // Refresh Token
    ExpiresIn    int64    `json:"expiresIn"`    // 过期时间（秒）
    UserInfo     UserInfo `json:"userInfo"`     // 用户信息
}

// UserInfo 用户信息DTO（脱敏，不含密码）
type UserInfo struct {
    ID           uint64 `json:"id"`
    Username     string `json:"username"`
    Nickname     string `json:"nickname"`
    Email        string `json:"email"`
    Phone        string `json:"phone"`
    Avatar       string `json:"avatar"`
    Status       int    `json:"status"`
    LastLoginTime string `json:"lastLoginTime"`
    CreatedAt    string `json:"createdAt"`
}

// ==================== 搜索条件 ====================

// SearchUser 用户搜索条件
type SearchUser struct {
    Page       int64    `json:"page" form:"page" search:"page"`                               // 分页
    PageSize   int64    `json:"pageSize" form:"pageSize" search:"pageSize"`                   // 分页大小
    ID         uint64   `json:"id" form:"id" search:"type:eq;column:id;table:users"`          // ID精确匹配
    Username   string   `json:"username" form:"username" search:"type:like;column:username;table:users"` // 用户名模糊查询
    Status     int      `json:"status" form:"status" search:"type:eq;column:status;table:users"`         // 状态
    Email      string   `json:"email" form:"email" search:"type:eq;column:email;table:users"`            // 邮箱
    Phone      string   `json:"phone" form:"phone" search:"type:eq;column:phone;table:users"`            // 手机号
    CreatedAt  []string `json:"createdAt" form:"createdAt" search:"type:between;column:created_at;table:users"`` // 创建时间范围
}

// ListUserResponse 用户列表响应
type ListUserResponse struct {
    Page     int64      `json:"page"`
    PageSize int64      `json:"pageSize"`
    Total    int64      `json:"total"`
    List     []*UserInfo `json:"list"`
}
```

#### 4.4.3 认证中间件 (middleware/auth.go)
```go
package middleware

import (
    "strings"
    "github.com/gin-gonic/gin"
    userDomain "thinkingModels/domain/user/user" // 注意包名引用
)

const UserContextKey = "user_context"

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 从Header获取Token
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"code": 401, "msg": "未提供认证信息"})
            c.Abort()
            return
        }

        // 2. 解析 Bearer Token
        parts := strings.SplitN(authHeader, " ", 2)
        if !(len(parts) == 2 && parts[0] == "Bearer") {
            c.JSON(401, gin.H{"code": 401, "msg": "认证格式错误"})
            c.Abort()
            return
        }

        // 3. 解析Token（使用user领域的包级函数）
        claims, err := userDomain.ParseToken(parts[1])
        if err != nil {
            c.JSON(401, gin.H{"code": 401, "msg": "登录已过期或无效"})
            c.Abort()
            return
        }

        // 4. 检查Token是否在黑名单
        if userDomain.IsTokenBlacklisted(claims.ID, db.InitDb()) {
            c.JSON(401, gin.H{"code": 401, "msg": "Token已被注销"})
            c.Abort()
            return
        }

        // 4. 将用户信息存入上下文
        c.Set(UserContextKey, &UserContext{
            UserID:       claims.UserID,
            Username:     claims.Username,
            EnterpriseID: claims.EnterpriseID,
            Roles:        claims.Roles,
        })

        c.Next()
    }
}

// GetUserContext 从上下文中获取用户信息
func GetUserContext(c *gin.Context) *UserContext {
    if val, exists := c.Get(UserContextKey); exists {
        return val.(*UserContext)
    }
    return nil
}
```

#### 4.4.4 登录逻辑 (logic/user/auth.go)
```go
package user

import (
    "errors"
    userDomain "thinkingModels/domain/user/user"
)

// Login 用户登录
func (l *AuthLogic) Login(req *userDomain.LoginRequest) (*userDomain.LoginResponse, error) {
    // 1. 根据用户名查询用户
    userEntity := userDomain.NewUserEntity(l.Ctx)
    user, err := userEntity.FindByUsername(req.Username)
    if err != nil {
        return nil, errors.New("用户名或密码错误")
    }

    // 2. 检查用户状态
    if user.Status == 0 {
        return nil, errors.New("账号已被禁用")
    }

    // 3. 验证密码（实体方法）
    if !user.VerifyPassword(req.Password) {
        return nil, errors.New("用户名或密码错误")
    }

    // 4. 更新登录信息（实体方法）
    user.UpdateLoginInfo(l.Ctx.ClientIP())
    user.Update() // 保存到数据库

    // 5. 生成JWT Token（实体方法 - 充血模型）
    tokenPair, err := user.GenerateToken()
    if err != nil {
        return nil, errors.New("Token生成失败")
    }

    return &LoginResponse{
        AccessToken:  tokenPair.AccessToken,
        RefreshToken: tokenPair.RefreshToken,
        ExpiresIn:    tokenPair.ExpiresIn,
        UserInfo:     user.ToUserInfo(), // 实体方法：转换为脱敏DTO
    }, nil
}

// Logout 用户登出
func (l *AuthLogic) Logout(token string) error {
    // 1. 解析Token获取JTI（包级函数）
    claims, err := userDomain.ParseToken(token)
    if err != nil {
        return err
    }

    // 2. 将Token加入黑名单（包级函数）
    return userDomain.BlacklistToken(
        claims.ID,
        claims.UserID,
        claims.ExpiresAt.Time,
        db.InitDb(),
    )
}
}
```

### 4.5 配置文件
```yaml
# config.yaml 新增JWT配置
jwt:
  secret: "your-secret-key-here-change-in-production"  # JWT密钥
  access_expiry: 3600    # Access Token有效期（秒）
  refresh_expiry: 604800 # Refresh Token有效期（秒）
  issuer: "thinkingModels" # 签发者
```

---

## 5. 开发计划

### Phase 1: 基础用户功能（本期实现）
- [x] 1. 创建用户领域目录结构
- [ ] 2. 创建用户表（users）
- [ ] 3. 实现JWT工具类
- [ ] 4. 实现用户实体模型
- [ ] 5. 实现登录/登出API
- [ ] 6. 实现认证中间件
- [ ] 7. 实现用户上下文注入
- [ ] 8. 实现用户注册API
- [ ] 9. 实现用户信息获取/更新API
- [ ] 10. 添加密码加密（bcrypt）
- [ ] 11. 测试联调

### Phase 2: 角色权限管理（后续迭代）
- [ ] 1. 创建角色表、权限表、用户角色关联表
- [ ] 2. 实现RBAC权限模型
- [ ] 3. 实现角色管理API
- [ ] 4. 实现权限管理API
- [ ] 5. 实现权限校验中间件

### Phase 3: 企业/租户管理（后续迭代）
- [ ] 1. 创建企业表
- [ ] 2. 实现企业隔离（数据权限）
- [ ] 3. 实现企业管理API

---

## 6. 安全考虑

1. **密码安全**
   - 使用 bcrypt 进行密码加密（自适应哈希，自动加盐）
   - 密码强度校验（最少8位，包含大小写字母+数字）

2. **Token安全**
   - Access Token 有效期短（1小时）
   - Refresh Token 使用 HttpOnly Cookie
   - Token 支持黑名单机制（登出失效）
   - 每个Token有唯一标识（JTI）

3. **传输安全**
   - 生产环境强制 HTTPS
   - 密码等敏感信息不打印日志

4. **接口安全**
   - 登录接口添加验证码/限流（防暴力破解）
   - 敏感操作需二次验证

---

## 7. 依赖项

需要添加的Go依赖：
```bash
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt
```

---

## 8. 表结构设计公约（强制规范）

> **所有数据库表必须包含以下标准审计字段**，这是项目强制规范。

### 8.1 标准字段定义

| 字段名 | 类型 | 默认值 | 约束 | 说明 |
|--------|------|--------|------|------|
| `created_at` | DATETIME | CURRENT_TIMESTAMP | NOT NULL | 创建时间，自动填充 |
| `updated_at` | DATETIME | CURRENT_TIMESTAMP | NULL | 修改时间，自动更新（ON UPDATE） |
| `deleted_at` | DATETIME | NULL | NULL | 软删除标记，NULL表示未删除 |
| `create_by` | BIGINT UNSIGNED | 0 | NOT NULL | 创建人ID，默认0表示系统 |
| `create_by_name` | VARCHAR(20) | '系统' | NOT NULL | 创建人姓名 |
| `update_by` | BIGINT UNSIGNED | 0 | NOT NULL | 更新人ID，默认0表示系统 |
| `update_by_name` | VARCHAR(20) | '系统' | NOT NULL | 更新人姓名 |

### 8.2 完整 SQL 模板

```sql
CREATE TABLE example_table (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    -- 业务字段放在这里...

    -- 标准审计字段（所有表必须包含，复制粘贴使用）
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NULL COMMENT '修改时间',
    deleted_at DATETIME NULL COMMENT '软删除时间',
    create_by BIGINT UNSIGNED DEFAULT 0 NOT NULL COMMENT '创建人ID',
    create_by_name VARCHAR(20) DEFAULT '系统' NOT NULL COMMENT '创建人姓名',
    update_by BIGINT UNSIGNED DEFAULT 0 NOT NULL COMMENT '修改人ID',
    update_by_name VARCHAR(20) DEFAULT '系统' NOT NULL COMMENT '修改人姓名',

    -- 索引建议
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='表示例表';
```

### 8.3 GORM 模型定义

```go
// BaseModel 基础模型（所有实体必须嵌入）
type BaseModel struct {
    ID           uint64         `json:"id" gorm:"primaryKey;comment:主键ID"`
    CreatedAt    LocalTime      `json:"createdAt" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
    UpdatedAt    LocalTime      `json:"updatedAt" gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:修改时间"`
    DeletedAt    gorm.DeletedAt `json:"-" gorm:"index;comment:软删除时间"`
    CreateBy     uint64         `json:"createBy" gorm:"not null;default:0;comment:创建人ID"`
    CreateByName string         `json:"createByName" gorm:"size:20;not null;default:系统;comment:创建人姓名"`
    UpdateBy     uint64         `json:"updateBy" gorm:"not null;default:0;comment:修改人ID"`
    UpdateByName string         `json:"updateByName" gorm:"size:20;not null;default:系统;comment:修改人姓名"`
}
```

### 8.4 公约检查清单

创建新表时必须检查：
- [ ] 包含 `created_at` 字段，DATETIME 类型，NOT NULL，默认 CURRENT_TIMESTAMP
- [ ] 包含 `updated_at` 字段，DATETIME 类型，可空，默认 CURRENT_TIMESTAMP ON UPDATE
- [ ] 包含 `deleted_at` 字段，DATETIME 类型，可空，用于软删除
- [ ] 包含 `create_by` 字段，BIGINT UNSIGNED，默认 0，NOT NULL
- [ ] 包含 `create_by_name` 字段，VARCHAR(20)，默认 '系统'，NOT NULL
- [ ] 包含 `update_by` 字段，BIGINT UNSIGNED，默认 0，NOT NULL
- [ ] 包含 `update_by_name` 字段，VARCHAR(20)，默认 '系统'，NOT NULL
- [ ] 为 `created_at` 和 `deleted_at` 字段添加索引

---

**文档版本**: v1.1
**更新日期**: 2026-02-10
**作者**: Claude Code
