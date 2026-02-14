package thinking

import (
	"strconv"

	"thinkingModels/api"
	"thinkingModels/domain/practice/tag"
	"thinkingModels/logic/practice"

	"github.com/gin-gonic/gin"
)

type Tag struct {
	api.Base
}

func NewTag() *Tag {
	return &Tag{}
}

// GetByModel 获取模型的标签
func (a *Tag) GetByModel(ctx *gin.Context) {
	a.Ctx = ctx
	idStr := ctx.Param("modelId")
	modelId, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewTagLogic(ctx)
	res, err := logic.GetByModel(modelId)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}

// AddToModel 给模型添加标签
func (a *Tag) AddToModel(ctx *gin.Context) {
	a.Ctx = ctx
	req := &tag.AddTagsToModel{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewTagLogic(ctx)
	err := logic.AddToModel(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(nil, "添加成功")
}

// RemoveFromModel 从模型移除标签
func (a *Tag) RemoveFromModel(ctx *gin.Context) {
	a.Ctx = ctx
	req := &tag.RemoveTagsFromModel{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewTagLogic(ctx)
	err := logic.RemoveFromModel(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(nil, "移除成功")
}

// Hot 获取热门标签
func (a *Tag) Hot(ctx *gin.Context) {
	a.Ctx = ctx
	limitStr := ctx.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	logic := thinking.NewTagLogic(ctx)
	res, err := logic.Hot(limit)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}
