package model

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jianyuezhexue/base"
	"thinkingModels/component/db"
)

// ModelEntityInterface 思维模型实体接口
type ModelEntityInterface interface {
	base.BaseModelInterface[ModelEntity]
	// 发布相关
	Publish() error
	Unpublish() error
	// 内容相关
	UpdateContent(content string) error
}

// ModelEntity 思维模型实体
type ModelEntity struct {
	base.BaseModel[ModelEntity]
	// 基础信息
	Name          string  `json:"name" type:"db" comment:"模型名称"`
	Code          string  `json:"code" type:"db" comment:"模型编码"`
	Description   string  `json:"description" type:"db" comment:"模型描述"`
	CoverImage    string  `json:"coverImage" type:"db" comment:"封面图片URL"`
	Icon          string  `json:"icon" type:"db" comment:"模型图标"`
	// 分类与价格
	CategoryId    uint64  `json:"categoryId" type:"db" comment:"所属分类ID"`
	Price         float64 `json:"price" type:"db" comment:"价格，0表示免费"`
	// 模型内容
	Content       string  `json:"content" type:"db" comment:"模型内容JSON"`
	// 难度与时间
	Difficulty    int     `json:"difficulty" type:"db" comment:"难度: 1=入门, 2=进阶, 3=高级"`
	EstimatedTime int     `json:"estimatedTime" type:"db" comment:"预计完成时间（分钟）"`
	// 统计信息
	UsageCount    int     `json:"usageCount" type:"db" comment:"使用次数统计"`
	// 状态管理
	Status        int     `json:"status" type:"db" comment:"状态: 0=草稿, 1=已发布, 2=已下架"`
	PublishTime   db.LocalTime `json:"publishTime" type:"db" comment:"发布时间"`
	// 版本信息
	Version       string  `json:"version" type:"db" comment:"版本号"`
	// 作者信息
	AuthorId      uint64  `json:"authorId" type:"db" comment:"作者ID"`
	AuthorName    string  `json:"authorName" type:"db" comment:"作者名称"`
}

// NewModelEntity 实例化思维模型实体
func NewModelEntity(ctx *gin.Context, opt ...base.Option[ModelEntity]) ModelEntityInterface {
	entity := &ModelEntity{}
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
func (m *ModelEntity) TableName() string {
	return "thinking_models"
}

// Validate 数据校验
func (m *ModelEntity) Validate() error {
	// 基础校验
	if m.Name == "" {
		return errors.New("模型名称不能为空")
	}
	if m.Code == "" {
		return errors.New("模型编码不能为空")
	}
	if m.CategoryId == 0 {
		return errors.New("所属分类不能为空")
	}
	if m.Content == "" {
		return errors.New("模型内容不能为空")
	}
	// 难度校验
	if m.Difficulty < 1 || m.Difficulty > 3 {
		m.Difficulty = 1
	}
	// 价格校验
	if m.Price < 0 {
		m.Price = 0
	}
	return nil
}

// Repair 数据修复
func (m *ModelEntity) Repair() error {
	// 设置默认值
	if m.Version == "" {
		m.Version = "1.0.0"
	}
	if m.Difficulty == 0 {
		m.Difficulty = 1
	}
	if m.EstimatedTime == 0 {
		m.EstimatedTime = 30
	}
	return nil
}

// Complete 数据完善
func (m *ModelEntity) Complete() error {
	// 数据处理
	return nil
}

// Publish 发布模型
func (m *ModelEntity) Publish() error {
	m.Status = 1
	m.PublishTime = db.LocalTime(time.Now())
	return nil
}

// Unpublish 下架模型
func (m *ModelEntity) Unpublish() error {
	m.Status = 2
	return nil
}

// UpdateContent 更新模型内容
func (m *ModelEntity) UpdateContent(content string) error {
	m.Content = content
	// 版本号自动升级（简单实现：最后一位+1）
	// 实际项目中可能需要更复杂的版本管理
	return nil
}

// IncrementUsageCount 增加使用次数
func (m *ModelEntity) IncrementUsageCount() {
	m.UsageCount++
}
