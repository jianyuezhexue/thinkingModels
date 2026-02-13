package thinking

import (
	"strconv"

	"thinkingModels/api"
	"thinkingModels/domain/thinking/model"
	"thinkingModels/logic/thinking"

	"github.com/gin-gonic/gin"
)

type Model struct {
	api.Base
}

func NewModel() *Model {
	return &Model{}
}

// Create 创建思维模型
func (a *Model) Create(ctx *gin.Context) {
	a.Ctx = ctx
	req := &model.CreateModel{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewModelLogic(ctx)
	res, err := logic.Create(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "创建成功")
}

// Update 更新思维模型
func (a *Model) Update(ctx *gin.Context) {
	a.Ctx = ctx
	req := &model.UpdateModel{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewModelLogic(ctx)
	res, err := logic.Update(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "更新成功")
}

// Get 获取思维模型详情
func (a *Model) Get(ctx *gin.Context) {
	a.Ctx = ctx
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewModelLogic(ctx)
	res, err := logic.Get(id)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}

// GetByCode 根据编码获取思维模型
func (a *Model) GetByCode(ctx *gin.Context) {
	a.Ctx = ctx
	code := ctx.Param("code")

	logic := thinking.NewModelLogic(ctx)
	res, err := logic.GetByCode(code)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}

// List 获取思维模型列表
func (a *Model) List(ctx *gin.Context) {
	a.Ctx = ctx
	req := &model.SearchModel{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewModelLogic(ctx)
	res, err := logic.List(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}

// ListMy 获取我的思维模型
func (a *Model) ListMy(ctx *gin.Context) {
	a.Ctx = ctx
	req := &model.SearchModel{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewModelLogic(ctx)
	res, err := logic.ListMy(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}

// Del 删除思维模型
func (a *Model) Del(ctx *gin.Context) {
	a.Ctx = ctx
	req := &model.DelModel{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewModelLogic(ctx)
	err := logic.Del(req.Ids)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(nil, "删除成功")
}

// Publish 发布思维模型
func (a *Model) Publish(ctx *gin.Context) {
	a.Ctx = ctx
	req := &model.PublishModel{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewModelLogic(ctx)
	res, err := logic.Publish(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "发布成功")
}

// Unpublish 下架思维模型
func (a *Model) Unpublish(ctx *gin.Context) {
	a.Ctx = ctx
	req := &model.UnpublishModel{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewModelLogic(ctx)
	res, err := logic.Unpublish(req.Id)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "下架成功")
}

// Fork 引用创建思维模型
func (a *Model) Fork(ctx *gin.Context) {
	a.Ctx = ctx
	req := &model.ForkModel{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewModelLogic(ctx)
	res, err := logic.Fork(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "创建成功")
}
