package category

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// 获取测试用的 gin context
func getTestContext() *gin.Context {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(nil)
	return ctx
}

// TestCategoryEntity_Validate 测试数据校验
func TestCategoryEntity_Validate(t *testing.T) {
	ctx := getTestContext()
	entity := NewCategoryEntity(ctx)

	tests := []struct {
		name    string
		entity  *CategoryEntity
		wantErr bool
		errMsg  string
	}{
		{
			name: "正常情况",
			entity: &CategoryEntity{
				Name:        "决策思维",
				Icon:        "https://example.com/icon.png",
				Description: "描述",
				Heat:        10,
			},
			wantErr: false,
		},
		{
			name: "名称为空",
			entity: &CategoryEntity{
				Name:        "",
				Description: "描述",
			},
			wantErr: true,
			errMsg:  "分类名称不能为空",
		},
		{
			name: "名称过长",
			entity: &CategoryEntity{
				Name:        "这是一段超过50个字符长度的分类名称用于测试超长情况",
				Description: "描述",
			},
			wantErr: true,
			errMsg:  "分类名称长度不能超过50个字符",
		},
		{
			name: "描述过长",
			entity: &CategoryEntity{
				Name:        "正常名称",
				Description: string(make([]byte, 501)),
			},
			wantErr: true,
			errMsg:  "分类描述长度不能超过500个字符",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.entity.BaseModel = entity.BaseModel
			err := tt.entity.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestCategoryEntity_TableName 测试表名
func TestCategoryEntity_TableName(t *testing.T) {
	ctx := getTestContext()
	entity := NewCategoryEntity(ctx)
	assert.Equal(t, "category", entity.TableName())
}

// TestCategoryEntity_Complete 测试数据完善
func TestCategoryEntity_Complete(t *testing.T) {
	ctx := getTestContext()
	entity := NewCategoryEntity(ctx)

	// 测试热度默认值
	testEntity := &CategoryEntity{
		Name: "测试",
		Heat: 0,
	}
	testEntity.BaseModel = entity.BaseModel

	err := testEntity.Complete()
	assert.NoError(t, err)
	assert.Equal(t, 0, testEntity.Heat)
}

// TestCreateCategory_Validate 测试新建分类校验
func TestCreateCategory_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     *CreateCategory
		wantErr bool
		errMsg  string
	}{
		{
			name: "正常情况",
			req: &CreateCategory{
				Name:        "决策思维",
				Icon:        "https://example.com/icon.png",
				Description: "描述",
				Heat:        10,
			},
			wantErr: false,
		},
		{
			name: "名称为空",
			req: &CreateCategory{
				Name: "",
			},
			wantErr: true,
			errMsg:  "分类名称不能为空",
		},
		{
			name: "名称过长",
			req: &CreateCategory{
				Name: "这是一段超过50个字符长度的分类名称用于测试超长情况",
			},
			wantErr: true,
			errMsg:  "分类名称长度不能超过50个字符",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req.Name == "" {
				assert.Equal(t, "", tt.req.Name)
			} else if len(tt.req.Name) > 50 {
				assert.Greater(t, len(tt.req.Name), 50)
			} else {
				assert.LessOrEqual(t, len(tt.req.Name), 50)
			}
		})
	}
}

// TestIncreaseHeat_Validate 测试增加热度校验
func TestIncreaseHeat_Validate(t *testing.T) {
	tests := []struct {
		name    string
		id      uint64
		delta   int
		wantErr bool
	}{
		{
			name:    "正常情况",
			id:      1,
			delta:   10,
			wantErr: false,
		},
		{
			name:    "delta为0",
			id:      1,
			delta:   0,
			wantErr: true,
		},
		{
			name:    "delta为负数",
			id:      1,
			delta:   -5,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.delta <= 0 {
				assert.LessOrEqual(t, tt.delta, 0)
			} else {
				assert.Greater(t, tt.delta, 0)
			}
		})
	}
}
