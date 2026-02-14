package thinking

import (
	"thinkingModels/domain/practice/tag"
	"thinkingModels/logic"

	"github.com/gin-gonic/gin"
)

// TagLogic 标签业务逻辑
type TagLogic struct {
	logic.BaseLogic
}

// NewTagLogic 初始化TagLogic
func NewTagLogic(ctx *gin.Context) *TagLogic {
	return &TagLogic{BaseLogic: logic.BaseLogic{Ctx: ctx}}
}

// GetByModel 获取模型的标签
func (l *TagLogic) GetByModel(modelId uint64) ([]*tag.TagInfo, error) {
	entity := tag.NewTagEntity(l.Ctx)
	tags, err := entity.GetByModel(modelId)
	if err != nil {
		return nil, err
	}

	result := make([]*tag.TagInfo, 0, len(tags))
	for _, t := range tags {
		result = append(result, &tag.TagInfo{
			Id:      t.Id,
			ModelId: t.ModelId,
			TagName: t.TagName,
		})
	}
	return result, nil
}

// AddToModel 为模型添加标签
func (l *TagLogic) AddToModel(req *tag.AddTagsToModel) error {
	entity := tag.NewTagEntity(l.Ctx)
	return entity.AddToModel(req.ModelId, req.TagNames)
}

// RemoveFromModel 从模型移除标签
func (l *TagLogic) RemoveFromModel(req *tag.RemoveTagsFromModel) error {
	entity := tag.NewTagEntity(l.Ctx)
	return entity.RemoveFromModel(req.ModelId, req.TagNames)
}

// Hot 获取热门标签
func (l *TagLogic) Hot(limit int) ([]*tag.HotTag, error) {
	entity := tag.NewTagEntity(l.Ctx)
	return entity.GetHotTags(limit)
}
