package market

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"thinkingModels/api"
	"thinkingModels/domain/market/model"
	"thinkingModels/logic/market"
)

type Model struct {
	api.Base
}

func NewModel() *Model {
	return &Model{}
}

// Create 创建思维模型
func (a *Model) Create(ctx *gin.Context) {
	// 设置上下文
	a.Ctx = ctx

	// 参数校验
	req := &model.CreateModel{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	// 实例化逻辑层
	logic := market.NewModelLogic(ctx)
	res, err := logic.Create(req)
	if err != nil {
		a.Error(err)
		return
	}

	// 接口返回
	a.Success(res, "创建成功")
}

// Update 更新思维模型
func (a *Model) Update(ctx *gin.Context) {
	// 设置上下文
	a.Ctx = ctx

	// 参数校验
	req := &model.UpdateModel{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	// 实例化逻辑层
	logic := market.NewModelLogic(ctx)
	res, err := logic.Update(req)
	if err != nil {
		a.Error(err)
		return
	}

	// 接口返回
	a.Success(res, "更新成功")
}

// Get 查询思维模型详情
func (a *Model) Get(ctx *gin.Context) {
	// 设置上下文（必须在调用其他方法前设置）
	a.Ctx = ctx

	// 从路径参数获取id
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		a.Error(err)
		return
	}

	// 实例化逻辑层
	logic := market.NewModelLogic(ctx)
	res, err := logic.Get(id)
	if err != nil {
		a.Error(err)
		return
	}

	// 接口返回
	a.Success(res, "查询成功")
}

// List 查询思维模型列表
func (a *Model) List(ctx *gin.Context) {
	// 设置上下文
	a.Ctx = ctx

	// 参数校验
	req := &model.SearchModel{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	// 实例化逻辑层
	logic := market.NewModelLogic(ctx)
	res, err := logic.List(req)
	if err != nil {
		a.Error(err)
		return
	}

	// 接口返回
	a.Success(res, "查询列表成功")
}

// Del 删除思维模型
func (a *Model) Del(ctx *gin.Context) {
	// 设置上下文
	a.Ctx = ctx

	req := &model.DelModel{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	// 实例化逻辑层
	logic := market.NewModelLogic(ctx)
	_, err := logic.Del(req.Ids)
	if err != nil {
		a.Error(err)
		return
	}

	// 接口返回
	a.Success(nil, "删除成功")
}

// Publish 发布思维模型
func (a *Model) Publish(ctx *gin.Context) {
	req := &model.PublishModel{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	// 实例化逻辑层
	logic := market.NewModelLogic(ctx)
	res, err := logic.Publish(req)
	if err != nil {
		a.Error(err)
		return
	}

	// 接口返回
	a.Success(res, "发布成功")
}

// Unpublish 下架思维模型
func (a *Model) Unpublish(ctx *gin.Context) {
	req := &model.UnpublishModel{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	// 实例化逻辑层
	logic := market.NewModelLogic(ctx)
	res, err := logic.Unpublish(req)
	if err != nil {
		a.Error(err)
		return
	}

	// 接口返回
	a.Success(res, "下架成功")
}

// UpdateStatus 更新思维模型状态
func (a *Model) UpdateStatus(ctx *gin.Context) {
	req := &model.UpdateModelStatus{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	// 实例化逻辑层
	logic := market.NewModelLogic(ctx)
	res, err := logic.UpdateStatus(req)
	if err != nil {
		a.Error(err)
		return
	}

	// 接口返回
	a.Success(res, "状态更新成功")
}

// GetByCode 根据编码查询思维模型
func (a *Model) GetByCode(ctx *gin.Context) {
	code := ctx.Param("code")
	if code == "" {
		a.Error(nil)
		return
	}

	// 实例化逻辑层
	logic := market.NewModelLogic(ctx)
	res, err := logic.GetByCode(code)
	if err != nil {
		a.Error(err)
		return
	}

	// 接口返回
	a.Success(res, "查询成功")
}
