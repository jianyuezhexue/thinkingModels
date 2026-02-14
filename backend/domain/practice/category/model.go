package category

import (
	"errors"

	"thinkingModels/component/db"

	"github.com/gin-gonic/gin"
	"github.com/jianyuezhexue/base"
)

// CategoryEntityInterface 分类实体接口
type CategoryEntityInterface interface {
	base.BaseModelInterface[CategoryEntity]
	GetChildren() ([]*CategoryEntity, error)
	GetParent() (*CategoryEntity, error)
	BuildPath() string
	UpdateModelCount(delta int)
}

// CategoryEntity 分类实体
type CategoryEntity struct {
	base.BaseModel[CategoryEntity]
	Code        string `json:"code" type:"db" comment:"分类编码"`
	Name        string `json:"name" type:"db" comment:"分类名称"`
	Icon        string `json:"icon" type:"db" comment:"分类图标"`
	Description string `json:"description" type:"db" comment:"分类描述"`
	ParentId    uint64 `json:"parentId" type:"db" comment:"父分类ID"`
	Sort        int    `json:"sort" type:"db" comment:"排序"`
	Level       int    `json:"level" type:"db" comment:"层级"`
	Path        string `json:"path" type:"db" comment:"路径"`
	ModelCount  int    `json:"modelCount" type:"db" comment:"模型数量"`
	Status      int    `json:"status" type:"db" comment:"状态: 0=禁用, 1=启用"`
}

// NewCategoryEntity 实例化分类实体
func NewCategoryEntity(ctx *gin.Context, opt ...base.Option[CategoryEntity]) CategoryEntityInterface {
	entity := &CategoryEntity{}
	entity.BaseModel = base.NewBaseModel(ctx, db.InitDb(), entity.TableName(), entity)
	if len(opt) > 0 {
		for _, fc := range opt {
			fc(&entity.BaseModel)
		}
	}
	return entity
}

// TableName 数据表名
func (c *CategoryEntity) TableName() string {
	return "model_categories"
}

// Validate 数据校验
func (c *CategoryEntity) Validate() error {
	if c.Name == "" {
		return errors.New("分类名称不能为空")
	}
	if len(c.Name) > 50 {
		return errors.New("分类名称不能超过50字符")
	}
	return nil
}

// Repair 数据修复
func (c *CategoryEntity) Repair() error {
	if c.Level == 0 {
		if c.ParentId == 0 {
			c.Level = 1
		}
	}
	if c.Status < 0 || c.Status > 1 {
		c.Status = 1
	}
	return nil
}

// Complete 数据完善
func (c *CategoryEntity) Complete() error {
	return nil
}

// GetChildren 获取子分类
func (c *CategoryEntity) GetChildren() ([]*CategoryEntity, error) {
	cond := c.MakeConditon(SearchCategory{ParentId: c.Id})
	return c.List(cond)
}

// GetParent 获取父分类
func (c *CategoryEntity) GetParent() (*CategoryEntity, error) {
	if c.ParentId == 0 {
		return nil, nil
	}
	return c.LoadById(c.ParentId)
}

// BuildPath 构建路径
func (c *CategoryEntity) BuildPath() string {
	if c.ParentId == 0 {
		return "/" + c.Code
	}
	parent, err := c.GetParent()
	if err != nil || parent == nil {
		return "/" + c.Code
	}
	return parent.BuildPath() + "/" + c.Code
}

// UpdateModelCount 更新模型数量
func (c *CategoryEntity) UpdateModelCount(delta int) {
	c.ModelCount += delta
	if c.ModelCount < 0 {
		c.ModelCount = 0
	}
}
