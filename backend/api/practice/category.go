package thinking

import (
	"strconv"

	"thinkingModels/api"
	"thinkingModels/domain/practice/category"
	"thinkingModels/logic/practice"

	"github.com/gin-gonic/gin"
)

type Category struct {
	api.Base
}

func NewCategory() *Category {
	return &Category{}
}

// Create 创建分类
func (a *Category) Create(ctx *gin.Context) {
	a.Ctx = ctx
	req := &category.CreateCategory{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewCategoryLogic(ctx)
	res, err := logic.Create(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "创建成功")
}

// Update 更新分类
func (a *Category) Update(ctx *gin.Context) {
	a.Ctx = ctx
	req := &category.UpdateCategory{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewCategoryLogic(ctx)
	res, err := logic.Update(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "更新成功")
}

// Get 查询分类详情
func (a *Category) Get(ctx *gin.Context) {
	a.Ctx = ctx
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewCategoryLogic(ctx)
	res, err := logic.Get(id)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}

// List 查询分类列表
func (a *Category) List(ctx *gin.Context) {
	a.Ctx = ctx
	req := &category.SearchCategory{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewCategoryLogic(ctx)
	res, err := logic.List(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}

// Tree 获取分类树
func (a *Category) Tree(ctx *gin.Context) {
	a.Ctx = ctx
	logic := thinking.NewCategoryLogic(ctx)
	res, err := logic.Tree()
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}

// Children 获取子分类
func (a *Category) Children(ctx *gin.Context) {
	a.Ctx = ctx
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewCategoryLogic(ctx)
	res, err := logic.GetChildren(id)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}

// Move 移动分类
func (a *Category) Move(ctx *gin.Context) {
	a.Ctx = ctx
	req := &category.MoveCategory{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewCategoryLogic(ctx)
	err := logic.Move(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(nil, "移动成功")
}

// UpdateStatus 更新分类状态
func (a *Category) UpdateStatus(ctx *gin.Context) {
	a.Ctx = ctx
	req := &category.UpdateCategoryStatus{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewCategoryLogic(ctx)
	err := logic.UpdateStatus(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(nil, "更新成功")
}

// Del 删除分类
func (a *Category) Del(ctx *gin.Context) {
	a.Ctx = ctx
	req := &category.DelCategory{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewCategoryLogic(ctx)
	err := logic.Del(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(nil, "删除成功")
}
