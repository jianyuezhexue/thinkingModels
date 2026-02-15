package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestUserInfo_WithValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建测试 token
	claims := &UserClaims{
		UserID:   9,
		Username: "vben",
		RoleIds:  "",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "thinkingModels",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	assert.NoError(t, err)

	// 创建测试上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.Header.Set("Authorization", "Bearer "+tokenString)

	// 执行中间件
	UserInfo()(c)

	// 验证结果
	currUserId, exists := c.Get("currUserId")
	assert.True(t, exists)
	assert.Equal(t, "9", currUserId)

	currUserName, exists := c.Get("currUserName")
	assert.True(t, exists)
	assert.Equal(t, "vben", currUserName)
}

func TestUserInfo_WithoutToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建测试上下文（无 token）
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}

	// 执行中间件
	UserInfo()(c)

	// 验证结果 - 应该设置默认系统用户
	currUserId, exists := c.Get("currUserId")
	assert.True(t, exists)
	assert.Equal(t, "0", currUserId)

	currUserName, exists := c.Get("currUserName")
	assert.True(t, exists)
	assert.Equal(t, "系统", currUserName)
}

func TestUserInfo_WithInvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建测试上下文（无效 token）
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.Header.Set("Authorization", "Bearer invalid-token")

	// 执行中间件
	UserInfo()(c)

	// 验证结果 - 应该设置默认系统用户
	currUserId, exists := c.Get("currUserId")
	assert.True(t, exists)
	assert.Equal(t, "0", currUserId)

	currUserName, exists := c.Get("currUserName")
	assert.True(t, exists)
	assert.Equal(t, "系统", currUserName)
}

func TestUserInfo_Float64UserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 模拟 JWT 解析后的 float64 类型用户ID
	claims := &UserClaims{
		UserID:   float64(123),
		Username: "testuser",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "thinkingModels",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.Header.Set("Authorization", "Bearer "+tokenString)

	UserInfo()(c)

	currUserId, exists := c.Get("currUserId")
	assert.True(t, exists)
	// 验证 float64 类型的 ID 正确转换为字符串
	assert.Equal(t, "123", currUserId)
}