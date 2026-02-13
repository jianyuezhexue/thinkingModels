package model

import (
	"errors"
	"time"

	"thinkingModels/component/db"

	"github.com/gin-gonic/gin"
	"github.com/jianyuezhexue/base"
)

// ModelEntityInterface 思维模型实体接口
type ModelEntityInterface interface {
	base.BaseModelInterface[ModelEntity]
	Publish() error
	Unpublish() error
	IncrementUsageCount()
	IncrementAdoptCount()
	IncrementLikeCount()
	IncrementCommentCount()
	Fork() ModelEntityInterface
}

// ModelEntity 思维模型实体
type ModelEntity struct {
	base.BaseModel[ModelEntity]
	Code          string `json:"code" type:"db" comment:"模型编码"`
	Name          string `json:"name" type:"db" comment:"模型名称"`
	Description   string `json:"description" type:"db" comment:"简短描述"`
	Overview      string `json:"overview" type:"db" comment:"概述"`
	Icon          string `json:"icon" type:"db" comment:"图标"`
	CategoryId    uint64 `json:"categoryId" type:"db" comment:"分类ID"`
	Content       string `json:"content" type:"db" comment:"模型内容(JSON)"`
	UsageGuide    string `json:"usageGuide" type:"db" comment:"使用指南"`
	Examples      string `json:"examples" type:"db" comment:"案例(JSON)"`
	AiPrompt      string `json:"aiPrompt" type:"db" comment:"AI提示词模板"`
	Difficulty    int    `json:"difficulty" type:"db" comment:"难度: 1=简单, 2=中等, 3=困难"`
	EstimatedTime int    `json:"estimatedTime" type:"db" comment:"预计用时(分钟)"`
	Status        int    `json:"status" type:"db" comment:"状态: 0=草稿, 1=已发布, 2=已下架"`
	Version       string `json:"version" type:"db" comment:"版本号"`
	IsOfficial    bool   `json:"isOfficial" type:"db" comment:"是否官方"`
	SourceModelId uint64 `json:"sourceModelId" type:"db" comment:"派生来源ID"`
	UsageCount    int64  `json:"usageCount" type:"db" comment:"使用次数"`
	AdoptCount    int64  `json:"adoptCount" type:"db" comment:"采纳次数"`
	LikeCount     int64  `json:"likeCount" type:"db" comment:"点赞数"`
	CommentCount  int64  `json:"commentCount" type:"db" comment:"评论数"`
}

// NewModelEntity 实例化思维模型实体
func NewModelEntity(ctx *gin.Context, opt ...base.Option[ModelEntity]) ModelEntityInterface {
	entity := &ModelEntity{}
	entity.BaseModel = base.NewBaseModel(ctx, db.InitDb(), entity.TableName(), entity)
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
	if m.Name == "" {
		return errors.New("模型名称不能为空")
	}
	if len(m.Name) > 100 {
		return errors.New("模型名称不能超过100字符")
	}
	if m.Difficulty < 1 || m.Difficulty > 3 {
		m.Difficulty = 1
	}
	return nil
}

// Repair 数据修复
func (m *ModelEntity) Repair() error {
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
	return nil
}

// Publish 发布模型
func (m *ModelEntity) Publish() error {
	if m.Status == 1 {
		return errors.New("模型已发布")
	}
	m.Status = 1
	return nil
}

// Unpublish 下架模型
func (m *ModelEntity) Unpublish() error {
	if m.Status != 1 {
		return errors.New("模型未发布，无法下架")
	}
	m.Status = 2
	return nil
}

// IncrementUsageCount 增加使用次数
func (m *ModelEntity) IncrementUsageCount() {
	m.UsageCount++
}

// IncrementAdoptCount 增加采纳次数
func (m *ModelEntity) IncrementAdoptCount() {
	m.AdoptCount++
}

// IncrementLikeCount 增加点赞数
func (m *ModelEntity) IncrementLikeCount() {
	m.LikeCount++
}

// IncrementCommentCount 增加评论数
func (m *ModelEntity) IncrementCommentCount() {
	m.CommentCount++
}

// Fork 派生模型
func (m *ModelEntity) Fork() ModelEntityInterface {
	newEntity := &ModelEntity{
		Code:          m.Code + "_fork_" + time.Now().Format("20060102150405"),
		Name:          m.Name + " (派生)",
		Description:   m.Description,
		Overview:      m.Overview,
		Icon:          m.Icon,
		CategoryId:    m.CategoryId,
		Content:       m.Content,
		UsageGuide:    m.UsageGuide,
		Examples:      m.Examples,
		AiPrompt:      m.AiPrompt,
		Difficulty:    m.Difficulty,
		EstimatedTime: m.EstimatedTime,
		Status:        0, // 草稿
		Version:       "1.0.0",
		IsOfficial:    false,
		SourceModelId: m.Id,
	}
	newEntity.BaseModel = m.BaseModel
	return newEntity
}
