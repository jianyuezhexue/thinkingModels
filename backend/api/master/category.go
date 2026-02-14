package master

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"thinkingModels/api"
	"thinkingModels/domain/master/category"
	"thinkingModels/logic/master"
)

// Category 分类API控制器
type Category struct {
	api.Base
}

// NewCategory 初始化Category控制器
func NewCategory() *Category {
	return &Category{}
}

// All 获取全量分类列表（按热度降序）
func (a *Category) All(ctx *gin.Context) {
	req := &struct{}{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := master.NewCategoryLogic(ctx)
	res, err := logic.All()
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

	logic := master.NewCategoryLogic(ctx)
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

	logic := master.NewCategoryLogic(ctx)
	res, err := logic.Update(req)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "更新成功")
}

// Get 查询分类详情
func (a *Category) Get(ctx *gin.Context) {
	req := &struct{}{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		a.Error(err)
		return
	}

	logic := master.NewCategoryLogic(ctx)
	res, err := logic.Get(id)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "查询成功")
}

// List 查询分类列表（分页）
func (a *Category) List(ctx *gin.Context) {
	req := &category.SearchCategory{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := master.NewCategoryLogic(ctx)
	res, err := logic.List(req)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "查询列表成功")
}

// Del 删除分类
func (a *Category) Del(ctx *gin.Context) {
	req := &category.DelCategory{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := master.NewCategoryLogic(ctx)
	res, err := logic.Del(req)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "删除成功")
}

// IncreaseHeat 增加分类热度
func (a *Category) IncreaseHeat(ctx *gin.Context) {
	req := &category.IncreaseHeatRequest{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := master.NewCategoryLogic(ctx)
	res, err := logic.IncreaseHeat(req)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "热度增加成功")
}
