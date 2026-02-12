package market

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"thinkingModels/api"
	"thinkingModels/domain/market/category"
	"thinkingModels/logic/market"
)

type Category struct {
	api.Base
}

func NewCategory() *Category {
	return &Category{}
}

// Create 创建分类
func (a *Category) Create(ctx *gin.Context) {
	req := &category.CreateCategory{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := market.NewCategoryLogic(ctx)
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

	logic := market.NewCategoryLogic(ctx)
	res, err := logic.Update(req)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "更新成功")
}

// Move 移动分类
func (a *Category) Move(ctx *gin.Context) {
	req := &category.MoveCategory{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := market.NewCategoryLogic(ctx)
	res, err := logic.Move(req)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "移动成功")
}

// Get 查询分类详情
func (a *Category) Get(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		a.Error(err)
		return
	}

	logic := market.NewCategoryLogic(ctx)
	res, err := logic.Get(id)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "查询成功")
}

// List 查询分类列表
func (a *Category) List(ctx *gin.Context) {
	req := &category.SearchCategory{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := market.NewCategoryLogic(ctx)
	res, err := logic.List(req)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "查询列表成功")
}

// Tree 获取分类树
func (a *Category) Tree(ctx *gin.Context) {
	logic := market.NewCategoryLogic(ctx)
	res, err := logic.Tree()
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "查询成功")
}

// TreeWithRoot 获取分类树（包含所有状态）
func (a *Category) TreeWithRoot(ctx *gin.Context) {
	logic := market.NewCategoryLogic(ctx)
	res, err := logic.TreeWithRoot()
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "查询成功")
}

// Children 获取子分类列表
func (a *Category) Children(ctx *gin.Context) {
	idStr := ctx.Param("id")
	parentId, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		// 如果没有指定ID，默认查询顶级分类
		parentId = 0
	}

	logic := market.NewCategoryLogic(ctx)
	res, err := logic.GetChildren(parentId)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "查询成功")
}

// Path 获取分类路径
func (a *Category) Path(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		a.Error(err)
		return
	}

	logic := market.NewCategoryLogic(ctx)
	res, err := logic.Path(id)
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

	logic := market.NewCategoryLogic(ctx)
	_, err := logic.Del(req.Ids)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(nil, "删除成功")
}

// UpdateStatus 更新分类状态
func (a *Category) UpdateStatus(ctx *gin.Context) {
	req := &category.UpdateCategoryStatus{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := market.NewCategoryLogic(ctx)
	res, err := logic.UpdateStatus(req)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "状态更新成功")
}
