package thinkingModel

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"thinkingModels/api"
	"thinkingModels/domain/thinkingModel/category"
	"thinkingModels/logic/thinkingModel"
)

type Category struct {
	api.Base
}

func NewCategory() *Category {
	return &Category{}
}

// All 获取全量分类列表
func (a *Category) All(ctx *gin.Context) {
	req := &struct{}{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinkingModel.NewCategoryLogic(ctx)
	res, err := logic.All()
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}

// List 获取分类列表（分页）
func (a *Category) List(ctx *gin.Context) {
	req := &category.SearchCategory{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		a.Error(err)
		return
	}

	logic := thinkingModel.NewCategoryLogic(ctx)
	res, err := logic.List(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}

// Create 创建分类
func (a *Category) Create(ctx *gin.Context) {
	req := &category.CreateCategory{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinkingModel.NewCategoryLogic(ctx)
	res, err := logic.Create(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "创建成功")
}

// Update 更新分类
func (a *Category) Update(ctx *gin.Context) {
	req := &category.UpdateCategory{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinkingModel.NewCategoryLogic(ctx)
	res, err := logic.Update(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "更新成功")
}

// Get 获取分类详情
func (a *Category) Get(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		a.Error(err)
		return
	}

	logic := thinkingModel.NewCategoryLogic(ctx)
	res, err := logic.Get(id)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}

// Del 删除分类
func (a *Category) Del(ctx *gin.Context) {
	req := &category.DelCategory{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinkingModel.NewCategoryLogic(ctx)
	_, err := logic.Del(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(nil, "删除成功")
}

// IncreaseHeat 增加分类热度
func (a *Category) IncreaseHeat(ctx *gin.Context) {
	req := &category.IncreaseHeatRequest{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinkingModel.NewCategoryLogic(ctx)
	res, err := logic.IncreaseHeat(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "热度更新成功")
}