package category

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/jianyuezhexue/base"
	"thinkingModels/component/db"
)

// CategoryEntity 分类实体业务模型
type CategoryEntity struct {
	base.BaseModel[CategoryEntity]
	Name        string `json:"name" type:"db" comment:"分类名称"`         // 分类名称
	Icon        string `json:"icon" type:"db" comment:"分类图标"`           // 分类图标URL
	Description string `json:"description" type:"db" comment:"分类描述"`    // 分类描述
	Heat        int    `json:"heat" type:"db" comment:"热度值"`             // 热度值，用于排序
}

// 实例化领域业务模型
func NewCategoryEntity(ctx *gin.Context) *CategoryEntity {
	entity := &CategoryEntity{}
	entity.BaseModel = base.NewBaseModel(ctx, db.InitDb(), entity.TableName(), entity)
	return entity
}

// 数据表名
func (m *CategoryEntity) TableName() string {
	return "category"
}

// ValidateFunc 数据校验
func (m *CategoryEntity) Validate() error {
	if m.Name == "" {
		return errors.New("分类名称不能为空")
	}
	if len(m.Name) > 50 {
		return errors.New("分类名称长度不能超过50个字符")
	}
	if len(m.Description) > 500 {
		return errors.New("分类描述长度不能超过500个字符")
	}
	return nil
}

// Repair 数据修复
func (m *CategoryEntity) Repair() error {
	// 数据修复逻辑
	return nil
}

// Complete 数据完善
func (m *CategoryEntity) Complete() error {
	// 热度默认值为0
	if m.Heat == 0 {
		m.Heat = 0
	}
	return nil
}
