package thinkingModel

import (
	"github.com/gin-gonic/gin"
	"thinkingModels/domain/thinkingModel/modelTag"
	"thinkingModels/logic"
)

// ModelTagLogic 模型标签业务逻辑
type ModelTagLogic struct {
	logic.BaseLogic
}

// NewModelTagLogic 初始化ModelTagLogic
func NewModelTagLogic(ctx *gin.Context) *ModelTagLogic {
	return &ModelTagLogic{BaseLogic: logic.BaseLogic{Ctx: ctx}}
}

// GetByModel 根据模型ID获取标签列表
func (l *ModelTagLogic) GetByModel(modelId uint64) ([]*modelTag.TagInfo, error) {
	entity := modelTag.NewModelTagEntity(l.Ctx)
	list, err := entity.GetByModelId(modelId)
	if err != nil {
		return nil, err
	}

	tagInfoList := make([]*modelTag.TagInfo, 0, len(list))
	for _, item := range list {
		tagInfoList = append(tagInfoList, &modelTag.TagInfo{
			Id:        item.Id,
			ModelId:   item.ModelId,
			TagId:     item.TagId,
			TagName:   item.TagName,
			Heat:      item.Heat,
			CreatedAt: item.CreatedAt.String(),
		})
	}

	return tagInfoList, nil
}

// AddToModel 为模型添加标签
func (l *ModelTagLogic) AddToModel(req *modelTag.AddTagsToModelRequest) error {
	entity := modelTag.NewModelTagEntity(l.Ctx)

	// 转换标签输入
	tags := make([]modelTag.TagInput, 0, len(req.Tags))
	for _, t := range req.Tags {
		tags = append(tags, modelTag.TagInput{
			TagId:   t.TagId,
			TagName: t.TagName,
		})
	}

	return entity.AddTagsToModel(req.ModelId, tags)
}

// RemoveFromModel 从模型移除标签
func (l *ModelTagLogic) RemoveFromModel(req *modelTag.RemoveTagsFromModelRequest) error {
	entity := modelTag.NewModelTagEntity(l.Ctx)
	return entity.RemoveTagsFromModel(req.ModelId, req.TagIds)
}

// Hot 获取热门标签
func (l *ModelTagLogic) Hot(limit int) (*modelTag.HotTagResponse, error) {
	if limit <= 0 {
		limit = 10
	}

	entity := modelTag.NewModelTagEntity(l.Ctx)
	hotTags, err := entity.GetHotTags(limit)
	if err != nil {
		return nil, err
	}

	list := make([]*modelTag.HotTagItem, 0, len(hotTags))
	for _, t := range hotTags {
		list = append(list, &modelTag.HotTagItem{
			TagName: t.TagName,
			Heat:    t.Heat,
			Count:   t.Count,
		})
	}

	return &modelTag.HotTagResponse{List: list}, nil
}

// List 查询模型标签列表
func (l *ModelTagLogic) List(req *modelTag.SearchModelTag) (*modelTag.ListModelTagResponse, error) {
	entity := modelTag.NewModelTagEntity(l.Ctx)
	cond := entity.MakeConditon(*req)

	total, err := entity.Count(cond)
	if err != nil {
		return nil, err
	}

	list, err := entity.List(cond)
	if err != nil {
		return nil, err
	}

	tagInfoList := make([]*modelTag.TagInfo, 0, len(list))
	for _, item := range list {
		tagInfoList = append(tagInfoList, &modelTag.TagInfo{
			Id:        item.Id,
			ModelId:   item.ModelId,
			TagId:     item.TagId,
			TagName:   item.TagName,
			Heat:      item.Heat,
			CreatedAt: item.CreatedAt.String(),
		})
	}

	return &modelTag.ListModelTagResponse{
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
		List:     tagInfoList,
	}, nil
}

// Del 删除模型标签
func (l *ModelTagLogic) Del(ids []uint64) error {
	entity := modelTag.NewModelTagEntity(l.Ctx)
	return entity.Del(ids...)
}

// IncreaseHeat 增加标签热度
func (l *ModelTagLogic) IncreaseHeat(req *modelTag.IncreaseHeatRequest) error {
	entity := modelTag.NewModelTagEntity(l.Ctx)
	return entity.IncrementHeat(req.Id, req.Delta)
}