package modelTag

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/jianyuezhexue/base"
	"thinkingModels/component/db"
)

// ModelTagEntityInterface 模型标签关联实体接口
type ModelTagEntityInterface interface {
	base.BaseModelInterface[ModelTagEntity]
	GetByModelId(modelId uint64) ([]*ModelTagEntity, error)
	AddTagsToModel(modelId uint64, tags []TagInput) error
	RemoveTagsFromModel(modelId uint64, tagIds []uint64) error
	GetHotTags(limit int) ([]*HotTag, error)
	IncrementHeat(id uint64, delta int) error
}

// TagInput 标签输入参数
type TagInput struct {
	TagId   uint64 `json:"tagId"`   // 标签ID（来自 master 领域的公共 tag）
	TagName string `json:"tagName"` // 标签名称
}

// HotTag 热门标签
type HotTag struct {
	TagName string `json:"tagName"`
	Heat    int    `json:"heat"`
	Count   int64  `json:"count"`
}

// ModelTagEntity 模型标签关联实体
type ModelTagEntity struct {
	base.BaseModel[ModelTagEntity]
	ModelId uint64 `json:"modelId" type:"db" comment:"模型ID"`   // 模型ID
	TagId   uint64 `json:"tagId" type:"db" comment:"标签ID"`     // 标签ID（关联 master 领域公共 tag）
	TagName string `json:"tagName" type:"db" comment:"标签名称"` // 标签名称（冗余存储，便于查询）
	Heat    int    `json:"heat" type:"db" comment:"热度值"`       // 热度值
}

// NewModelTagEntity 实例化模型标签关联实体
func NewModelTagEntity(ctx *gin.Context, opt ...base.Option[ModelTagEntity]) ModelTagEntityInterface {
	entity := &ModelTagEntity{}
	entity.BaseModel = base.NewBaseModel(ctx, db.InitDb(), entity.TableName(), entity)
	if len(opt) > 0 {
		for _, fc := range opt {
			fc(&entity.BaseModel)
		}
	}
	return entity
}

// TableName 数据表名
func (m *ModelTagEntity) TableName() string {
	return "thinking_model_tags"
}

// Validate 数据校验
func (m *ModelTagEntity) Validate() error {
	if m.ModelId == 0 {
		return errors.New("模型ID不能为空")
	}
	if m.TagName == "" {
		return errors.New("标签名称不能为空")
	}
	return nil
}

// Repair 数据修复
func (m *ModelTagEntity) Repair() error {
	// 热度默认值为0
	if m.Heat == 0 {
		m.Heat = 0
	}
	return nil
}

// Complete 数据完善
func (m *ModelTagEntity) Complete() error {
	return nil
}

// GetByModelId 根据模型ID获取标签列表
func (m *ModelTagEntity) GetByModelId(modelId uint64) ([]*ModelTagEntity, error) {
	var list []*ModelTagEntity
	err := m.Tx().Model(m).
		Where("model_id = ? AND deleted_at IS NULL", modelId).
		Order("heat DESC, id ASC").
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

// AddTagsToModel 为模型添加标签
func (m *ModelTagEntity) AddTagsToModel(modelId uint64, tags []TagInput) error {
	if modelId == 0 {
		return errors.New("模型ID不能为空")
	}
	if len(tags) == 0 {
		return errors.New("标签列表不能为空")
	}

	for _, tag := range tags {
		// 检查是否已存在
		var existing ModelTagEntity
		err := m.Tx().Model(&ModelTagEntity{}).
			Where("model_id = ? AND tag_id = ? AND deleted_at IS NULL", modelId, tag.TagId).
			First(&existing).Error
		if err == nil {
			// 已存在，跳过
			continue
		}

		// 创建新的关联
		entity := NewModelTagEntity(m.Ctx)
		if e, ok := entity.(*ModelTagEntity); ok {
			e.ModelId = modelId
			e.TagId = tag.TagId
			e.TagName = tag.TagName
			e.Heat = 0

			if err := e.Validate(); err != nil {
				return err
			}

			if _, err := e.Create(); err != nil {
				return err
			}
		}
	}

	return nil
}

// RemoveTagsFromModel 从模型移除标签
func (m *ModelTagEntity) RemoveTagsFromModel(modelId uint64, tagIds []uint64) error {
	if modelId == 0 {
		return errors.New("模型ID不能为空")
	}
	if len(tagIds) == 0 {
		return errors.New("标签ID列表不能为空")
	}

	// 软删除
	err := m.Tx().Model(&ModelTagEntity{}).
		Where("model_id = ? AND tag_id IN ?", modelId, tagIds).
		Delete(&ModelTagEntity{}).Error
	return err
}

// GetHotTags 获取热门标签
func (m *ModelTagEntity) GetHotTags(limit int) ([]*HotTag, error) {
	if limit <= 0 {
		limit = 10
	}

	var hotTags []*HotTag
	err := m.Tx().Model(&ModelTagEntity{}).
		Select("tag_name, SUM(heat) as heat, COUNT(*) as count").
		Where("deleted_at IS NULL").
		Group("tag_name").
		Order("heat DESC, count DESC").
		Limit(limit).
		Find(&hotTags).Error
	if err != nil {
		return nil, err
	}

	return hotTags, nil
}

// IncrementHeat 增加标签热度
func (m *ModelTagEntity) IncrementHeat(id uint64, delta int) error {
	if delta <= 0 {
		return errors.New("热度增加值必须大于0")
	}

	entity, err := m.LoadById(id)
	if err != nil {
		return err
	}

	entity.Heat += delta
	_, err = entity.Update()
	return err
}