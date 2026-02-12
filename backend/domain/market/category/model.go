package category

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jianyuezhexue/base"
	"thinkingModels/component/db"
)

// CategoryEntityInterface 分类实体接口
type CategoryEntityInterface interface {
	base.BaseModelInterface[CategoryEntity]
	// 树形结构操作
	GetChildren() ([]*CategoryEntity, error)
	GetParent() (*CategoryEntity, error)
	GetFullPath() ([]*CategoryEntity, error)
	// 层级管理
	UpdateLevel() error
	BuildPath() error
	// 禁用/启用
	Enable() error
	Disable() error
}

// CategoryEntity 分类实体
type CategoryEntity struct {
	base.BaseModel[CategoryEntity]
	ParentId    uint64 `json:"parentId" type:"db" comment:"父分类ID，0为顶级分类"`
	Name        string `json:"name" type:"db" comment:"分类名称"`
	Code        string `json:"code" type:"db" comment:"分类编码"`
	Description string `json:"description" type:"db" comment:"分类描述"`
	Icon        string `json:"icon" type:"db" comment:"分类图标URL"`
	SortOrder   int    `json:"sortOrder" type:"db" comment:"排序顺序"`
	Level       int    `json:"level" type:"db" comment:"层级深度"`
	Path        string `json:"path" type:"db" comment:"分类路径，如：/1/2/3/"`
	Status      int    `json:"status" type:"db" comment:"状态: 0=禁用, 1=启用"`
}

// NewCategoryEntity 实例化分类实体
func NewCategoryEntity(ctx *gin.Context, opt ...base.Option[CategoryEntity]) CategoryEntityInterface {
	entity := &CategoryEntity{}
	entity.BaseModel = base.NewBaseModel(ctx, db.InitDb(), entity.TableName(), entity)

	// 自定义配置选项
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
	if len(c.Name) > 100 {
		return errors.New("分类名称不能超过100字符")
	}
	if c.Code == "" {
		return errors.New("分类编码不能为空")
	}
	if len(c.Code) > 50 {
		return errors.New("分类编码不能超过50字符")
	}
	// 校验层级是否合理
	if c.Level < 1 || c.Level > 10 {
		return errors.New("分类层级必须在1-10之间")
	}
	return nil
}

// Repair 数据修复
func (c *CategoryEntity) Repair() error {
	if c.ParentId == 0 {
		c.Level = 1
		c.Path = fmt.Sprintf("/%d/", c.Id)
	}
	if c.Status == 0 && c.Status != 1 {
		c.Status = 1 // 默认启用
	}
	if c.SortOrder == 0 {
		c.SortOrder = 0
	}
	return nil
}

// Complete 数据完善
func (c *CategoryEntity) Complete() error {
	return nil
}

// GetChildren 获取子分类列表
func (c *CategoryEntity) GetChildren() ([]*CategoryEntity, error) {
	search := &SearchCategory{ParentId: c.Id, Status: -1} // -1表示查询所有状态
	cond := c.MakeConditon(*search)
	return c.List(cond)
}

// GetParent 获取父分类
func (c *CategoryEntity) GetParent() (*CategoryEntity, error) {
	if c.ParentId == 0 {
		return nil, nil // 没有父分类
	}
	return c.LoadById(c.ParentId)
}

// GetFullPath 获取完整路径（从根到当前）
func (c *CategoryEntity) GetFullPath() ([]*CategoryEntity, error) {
	path := make([]*CategoryEntity, 0)

	// 解析路径
	if c.Path == "" {
		return path, nil
	}

	// 移除首尾的/并分割
	trimmedPath := strings.Trim(c.Path, "/")
	if trimmedPath == "" {
		return path, nil
	}

	ids := strings.Split(trimmedPath, "/")
	for _, idStr := range ids {
		var id uint64
		_, err := fmt.Sscanf(idStr, "%d", &id)
		if err != nil {
			continue
		}
		parent, err := c.LoadById(id)
		if err != nil {
			continue
		}
		path = append(path, parent)
	}

	return path, nil
}

// UpdateLevel 更新层级
func (c *CategoryEntity) UpdateLevel() error {
	if c.ParentId == 0 {
		c.Level = 1
	} else {
		parent, err := c.GetParent()
		if err != nil {
			return err
		}
		if parent != nil {
			c.Level = parent.Level + 1
		}
	}
	return nil
}

// BuildPath 构建路径
func (c *CategoryEntity) BuildPath() error {
	if c.ParentId == 0 {
		c.Path = fmt.Sprintf("/%d/", c.Id)
	} else {
		parent, err := c.GetParent()
		if err != nil {
			return err
		}
		if parent != nil {
			c.Path = fmt.Sprintf("%s%d/", parent.Path, c.Id)
		}
	}
	return nil
}

// Enable 启用分类
func (c *CategoryEntity) Enable() error {
	c.Status = 1
	return nil
}

// Disable 禁用分类
func (c *CategoryEntity) Disable() error {
	c.Status = 0
	return nil
}

// IsRoot 是否为根分类
func (c *CategoryEntity) IsRoot() bool {
	return c.ParentId == 0
}
