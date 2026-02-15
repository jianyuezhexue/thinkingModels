package thinkingModel

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"thinkingModels/api"
	"thinkingModels/domain/thinkingModel/modelTag"
	"thinkingModels/logic/thinkingModel"
)

type ModelTag struct {
	api.Base
}

func NewModelTag() *ModelTag {
	return &ModelTag{}
}

// GetByModel 根据模型ID获取标签列表
func (a *ModelTag) GetByModel(ctx *gin.Context) {
	modelIdStr := ctx.Param("modelId")
	modelId, err := strconv.ParseUint(modelIdStr, 10, 64)
	if err != nil {
		a.Error(err)
		return
	}

	logic := thinkingModel.NewModelTagLogic(ctx)
	res, err := logic.GetByModel(modelId)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}

// AddToModel 为模型添加标签
func (a *ModelTag) AddToModel(ctx *gin.Context) {
	req := &modelTag.AddTagsToModelRequest{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinkingModel.NewModelTagLogic(ctx)
	err := logic.AddToModel(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(nil, "添加成功")
}

// RemoveFromModel 从模型移除标签
func (a *ModelTag) RemoveFromModel(ctx *gin.Context) {
	req := &modelTag.RemoveTagsFromModelRequest{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinkingModel.NewModelTagLogic(ctx)
	err := logic.RemoveFromModel(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(nil, "移除成功")
}

// Hot 获取热门标签
func (a *ModelTag) Hot(ctx *gin.Context) {
	req := &struct{}{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	limit := 10
	if limitStr := ctx.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	logic := thinkingModel.NewModelTagLogic(ctx)
	res, err := logic.Hot(limit)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}

// List 查询模型标签列表
func (a *ModelTag) List(ctx *gin.Context) {
	req := &modelTag.SearchModelTag{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		a.Error(err)
		return
	}

	logic := thinkingModel.NewModelTagLogic(ctx)
	res, err := logic.List(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}

// Del 删除模型标签
func (a *ModelTag) Del(ctx *gin.Context) {
	req := &modelTag.DelModelTag{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinkingModel.NewModelTagLogic(ctx)
	err := logic.Del(req.Ids)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(nil, "删除成功")
}

// IncreaseHeat 增加标签热度
func (a *ModelTag) IncreaseHeat(ctx *gin.Context) {
	req := &modelTag.IncreaseHeatRequest{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinkingModel.NewModelTagLogic(ctx)
	err := logic.IncreaseHeat(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(nil, "热度更新成功")
}