package thinkingModel

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"thinkingModels/api"
	"thinkingModels/domain/thinkingModel/model"
	"thinkingModels/logic/thinkingModel"
)

type Model struct {
	api.Base
}

func NewModel() *Model {
	return &Model{}
}

// Create 创建思维模型
func (a *Model) Create(ctx *gin.Context) {
	req := &model.CreateModel{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinkingModel.NewModelLogic(ctx)
	res, err := logic.Create(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "创建成功")
}

// Update 更新思维模型
func (a *Model) Update(ctx *gin.Context) {
	req := &model.UpdateModel{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinkingModel.NewModelLogic(ctx)
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

	logic := thinkingModel.NewModelLogic(ctx)
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

	logic := thinkingModel.NewModelLogic(ctx)
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

	logic := thinkingModel.NewModelLogic(ctx)
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

	logic := thinkingModel.NewModelLogic(ctx)
	res, err := logic.ListMy(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}

// Del 删除思维模型
func (a *Model) Del(ctx *gin.Context) {
	req := &model.DelModel{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinkingModel.NewModelLogic(ctx)
	err := logic.Del(req.Ids)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(nil, "删除成功")
}

// Publish 发布思维模型
func (a *Model) Publish(ctx *gin.Context) {
	req := &model.PublishModel{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinkingModel.NewModelLogic(ctx)
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
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		a.Error(err)
		return
	}

	logic := thinkingModel.NewModelLogic(ctx)
	res, err := logic.Unpublish(id)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "下架成功")
}

// Fork 引用创建思维模型
func (a *Model) Fork(ctx *gin.Context) {
	req := &model.ForkModel{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinkingModel.NewModelLogic(ctx)
	res, err := logic.Fork(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "创建成功")
}

// Review 审核思维模型
func (a *Model) Review(ctx *gin.Context) {
	req := &model.ReviewModel{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinkingModel.NewModelLogic(ctx)
	res, err := logic.Review(req)
	if err != nil {
		a.Error(err)
		return
	}

	if req.Approved {
		a.Success(res, "审核通过")
	} else {
		a.Success(res, "已驳回")
	}
}

// StatusCounts 获取各状态的模型数量统计
func (a *Model) StatusCounts(ctx *gin.Context) {
	a.Ctx = ctx

	logic := thinkingModel.NewModelLogic(ctx)
	res, err := logic.StatusCounts()
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "查询成功")
}