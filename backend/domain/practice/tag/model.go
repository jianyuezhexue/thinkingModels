package tag

import (
	"errors"

	"thinkingModels/component/db"

	"github.com/gin-gonic/gin"
	"github.com/jianyuezhexue/base"
)

// TagEntityInterface 标签实体接口
type TagEntityInterface interface {
	base.BaseModelInterface[TagEntity]
	GetByModel(modelId uint64) ([]*TagEntity, error)
	AddToModel(modelId uint64, tags []string) error
	RemoveFromModel(modelId uint64, tags []string) error
	GetHotTags(limit int) ([]*HotTag, error)
}

// TagEntity 标签实体
type TagEntity struct {
	base.BaseModel[TagEntity]
	ModelId uint64 `json:"modelId" type:"db" comment:"模型ID"`
	TagName string `json:"tagName" type:"db" comment:"标签名称"`
}

// NewTagEntity 实例化标签实体
func NewTagEntity(ctx *gin.Context, opt ...base.Option[TagEntity]) TagEntityInterface {
	entity := &TagEntity{}
	entity.BaseModel = base.NewBaseModel(ctx, db.InitDb(), entity.TableName(), entity)
	if len(opt) > 0 {
		for _, fc := range opt {
			fc(&entity.BaseModel)
		}
	}
	return entity
}

// TableName 数据表名
func (t *TagEntity) TableName() string {
	return "model_tags"
}

// Validate 数据校验
func (t *TagEntity) Validate() error {
	if t.ModelId == 0 {
		return errors.New("模型ID不能为空")
	}
	if t.TagName == "" {
		return errors.New("标签名称不能为空")
	}
	if len(t.TagName) > 50 {
		return errors.New("标签名称不能超过50字符")
	}
	return nil
}

// Repair 数据修复
func (t *TagEntity) Repair() error {
	return nil
}

// Complete 数据完善
func (t *TagEntity) Complete() error {
	return nil
}

// GetByModel 获取模型的标签
func (t *TagEntity) GetByModel(modelId uint64) ([]*TagEntity, error) {
	cond := t.MakeConditon(SearchTag{ModelId: modelId})
	return t.List(cond)
}

// AddToModel 为模型添加标签
func (t *TagEntity) AddToModel(modelId uint64, tags []string) error {
	for _, tagName := range tags {
		newTag := &TagEntity{
			ModelId: modelId,
			TagName: tagName,
		}
		newTag.BaseModel = t.BaseModel
		if _, err := newTag.Create(); err != nil {
			return err
		}
	}
	return nil
}

// RemoveFromModel 从模型移除标签
func (t *TagEntity) RemoveFromModel(modelId uint64, tags []string) error {
	for _, tagName := range tags {
		cond := t.MakeConditon(SearchTag{ModelId: modelId, TagName: tagName})
		list, err := t.List(cond)
		if err != nil {
			return err
		}
		for _, tag := range list {
			if err := t.Del(tag.Id); err != nil {
				return err
			}
		}
	}
	return nil
}

// GetHotTags 获取热门标签
func (t *TagEntity) GetHotTags(limit int) ([]*HotTag, error) {
	// 使用原生SQL查询热门标签
	var hotTags []*HotTag
	err := db.InitDb().Table(t.TableName()).
		Select("tag_name, COUNT(*) as count").
		Group("tag_name").
		Order("count DESC").
		Limit(limit).
		Find(&hotTags).Error
	return hotTags, err
}
