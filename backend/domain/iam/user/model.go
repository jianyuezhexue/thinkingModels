package user

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jianyuezhexue/base"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"thinkingModels/component/db"
)

// ========== JWT相关类型定义 ==========

// UserClaims JWT Claims定义
type UserClaims struct {
	UserID       uint64   `json:"sub"`
	Username     string   `json:"username"`
	EnterpriseID uint64   `json:"enterprise_id"`
	RoleIds      string   `json:"role_ids"`
	jwt.RegisteredClaims
}

// TokenPair Token对
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64 // 过期时间（秒）
}

// ========== 用户实体 ==========

// UserEntityInterface 用户实体接口
type UserEntityInterface interface {
	base.BaseModelInterface[UserEntity]
	// 认证相关方法
	VerifyPassword(plainPassword string) bool
	HashPassword() error
	GenerateToken() (*TokenPair, error)
	UpdateLoginInfo(ip string)
}

// UserEntity 用户实体
type UserEntity struct {
	base.BaseModel[UserEntity]
	Username       string       `json:"username" type:"db" comment:"登录用户名"`
	Password       string       `json:"-" type:"db" comment:"加密密码"`
	Nickname       string       `json:"nickname" type:"db" comment:"用户昵称"`
	Email          string       `json:"email" type:"db" comment:"邮箱"`
	Phone          string       `json:"phone" type:"db" comment:"手机号"`
	Avatar         string       `json:"avatar" type:"db" comment:"头像URL"`
	Status         int          `json:"status" type:"db" comment:"状态:0=禁用,1=正常"`
	LastLoginTime  db.LocalTime `json:"lastLoginTime" type:"db" comment:"最后登录时间"`
	LastLoginIP    string       `json:"lastLoginIp" type:"db" comment:"最后登录IP"`
	EnterpriseID   uint64       `json:"enterpriseId" type:"db" comment:"企业ID"`
	RoleIds        string       `json:"roleIds" type:"db" comment:"角色ID列表，逗号分隔"`
}

// 实例化用户实体
func NewUserEntity(ctx *gin.Context, opt ...base.Option[UserEntity]) UserEntityInterface {
	entity := &UserEntity{}
	entity.BaseModel = base.NewBaseModel(ctx, db.InitDb(), entity.TableName(), entity)

	// 自定义配置选项
	if len(opt) > 0 {
		for _, fc := range opt {
			fc(&entity.BaseModel)
		}
	}
	return entity
}

// 数据表名
func (m *UserEntity) TableName() string {
	return "users"
}

// Validate 数据校验
func (m *UserEntity) Validate() error {
	// more...
	return nil
}

// Repair 数据修复
func (m *UserEntity) Repair() error {
	// more...
	return nil
}

// Complete 数据完善（在实体内部完成数据补充、脱敏等面向对象操作）
func (m *UserEntity) Complete() error {
	// 数据完善逻辑：如设置默认值、计算字段、数据脱敏准备等
	// 所有需要在展示前完成的数据补充工作都在这里处理
	return nil
}

// ========== 认证相关方法（充血模型）==========

// VerifyPassword 验证密码（bcrypt比对）
func (u *UserEntity) VerifyPassword(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword))
	return err == nil
}

// HashPassword 密码加密（bcrypt）
func (u *UserEntity) HashPassword() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}

// GenerateToken 为用户生成JWT Token（实体方法）
func (u *UserEntity) GenerateToken() (*TokenPair, error) {
	// TODO: 从配置文件读取JWT密钥
	secret := "your-secret-key-change-in-production"
	now := time.Now()

	// Access Token (1小时有效期)
	accessClaims := UserClaims{
		UserID:       u.Id,
		Username:     u.Username,
		EnterpriseID: u.EnterpriseID,
		RoleIds:      u.RoleIds,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "thinkingModels",
			ID:        generateUniqueID(),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	// Refresh Token (7天有效期)
	refreshClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(7 * 24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        generateUniqueID(),
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600, // 1小时
	}, nil
}

// UpdateLoginInfo 更新登录信息
func (u *UserEntity) UpdateLoginInfo(ip string) {
	u.LastLoginTime = db.LocalTime(time.Now())
	u.LastLoginIP = ip
}

// generateUniqueID 生成唯一ID
func generateUniqueID() string {
	// 生成6字节随机数
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		// 如果随机数生成失败，回退到纳秒时间戳
		return time.Now().Format("20060102150405") + "-" + fmt.Sprintf("%d", time.Now().Nanosecond())
	}
	return time.Now().Format("20060102150405") + "-" + hex.EncodeToString(b)
}
