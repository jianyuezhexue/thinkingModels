package model

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestCreateModel_WithCurrUserId 测试创建模型时 currUserId 正确传递
func TestCreateModel_WithCurrUserId(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建测试上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}

	// 设置用户信息（模拟中间件）
	c.Set("currUserId", "9")
	c.Set("currUserName", "vben")

	// 创建实体
	entity := NewModelEntity(c)

	// 验证 base model 中的 operator 信息
	assert.NotNil(t, entity)

	// 获取底层实体
	concreteEntity, ok := entity.(*ModelEntity)
	assert.True(t, ok)

	// 验证 context 正确传递
	assert.NotNil(t, concreteEntity.Ctx)
}

// TestCreateModel_WithoutCurrUserId 测试无用户信息时使用默认值
func TestCreateModel_WithoutCurrUserId(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建测试上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}

	// 设置默认系统用户（模拟中间件对无 token 的处理）
	c.Set("currUserId", "0")
	c.Set("currUserName", "系统")

	// 创建实体
	entity := NewModelEntity(c)
	assert.NotNil(t, entity)

	concreteEntity, ok := entity.(*ModelEntity)
	assert.True(t, ok)
	assert.NotNil(t, concreteEntity.Ctx)
}

// TestModelEntity_Fields 测试模型字段定义
func TestModelEntity_Fields(t *testing.T) {
	entity := &ModelEntity{
		Name:          "测试模型",
		Description:   "这是一个测试模型",
		CoverImage:    "https://example.com/cover.jpg",
		Price:         99.0,
		AuthorId:      1,
		AuthorName:    "测试作者",
		CategoryId:    1,
		Difficulty:    1,
		EstimatedTime: 30,
	}

	assert.Equal(t, "测试模型", entity.Name)
	assert.Equal(t, 99.0, entity.Price)
	assert.Equal(t, uint64(1), entity.AuthorId)
	assert.Equal(t, "测试作者", entity.AuthorName)
}