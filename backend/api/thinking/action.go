package thinking

import (
	"strconv"

	"thinkingModels/api"
	"thinkingModels/domain/thinking/action"
	"thinkingModels/logic/thinking"

	"github.com/gin-gonic/gin"
)

type Action struct {
	api.Base
}

func NewAction() *Action {
	return &Action{}
}

// Create 创建行动项
func (a *Action) Create(ctx *gin.Context) {
	a.Ctx = ctx
	req := &action.CreateAction{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewActionLogic(ctx)
	res, err := logic.Create(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "创建成功")
}

// Update 更新行动项
func (a *Action) Update(ctx *gin.Context) {
	a.Ctx = ctx
	req := &action.UpdateAction{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewActionLogic(ctx)
	res, err := logic.Update(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "更新成功")
}

// Get 查询行动项详情
func (a *Action) Get(ctx *gin.Context) {
	a.Ctx = ctx
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewActionLogic(ctx)
	res, err := logic.Get(id)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}

// List 查询行动项列表
func (a *Action) List(ctx *gin.Context) {
	a.Ctx = ctx
	req := &action.SearchAction{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewActionLogic(ctx)
	res, err := logic.List(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}

// ListMy 查询我的行动项
func (a *Action) ListMy(ctx *gin.Context) {
	a.Ctx = ctx
	req := &action.SearchAction{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewActionLogic(ctx)
	res, err := logic.ListMy(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}

// Del 删除行动项
func (a *Action) Del(ctx *gin.Context) {
	a.Ctx = ctx
	req := &action.DelAction{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewActionLogic(ctx)
	err := logic.Del(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(nil, "删除成功")
}

// CreateFromAnalysis 从分析创建行动项
func (a *Action) CreateFromAnalysis(ctx *gin.Context) {
	a.Ctx = ctx
	req := &action.CreateActionsFromAnalysis{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewActionLogic(ctx)
	res, err := logic.CreateFromAnalysis(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "创建成功")
}

// ListByTopic 按课题查询行动项
func (a *Action) ListByTopic(ctx *gin.Context) {
	a.Ctx = ctx
	topicIdStr := ctx.Param("topicId")
	topicId, err := strconv.ParseUint(topicIdStr, 10, 64)
	if err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewActionLogic(ctx)
	res, err := logic.ListByTopic(topicId)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}

// ListByAnalysis 按分析查询行动项
func (a *Action) ListByAnalysis(ctx *gin.Context) {
	a.Ctx = ctx
	analysisIdStr := ctx.Param("analysisId")
	analysisId, err := strconv.ParseUint(analysisIdStr, 10, 64)
	if err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewActionLogic(ctx)
	res, err := logic.ListByAnalysis(analysisId)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}

// UpdateProgress 更新进度
func (a *Action) UpdateProgress(ctx *gin.Context) {
	a.Ctx = ctx
	req := &action.UpdateProgress{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewActionLogic(ctx)
	err := logic.UpdateProgress(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(nil, "更新成功")
}

// Complete 完成行动项
func (a *Action) Complete(ctx *gin.Context) {
	a.Ctx = ctx
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewActionLogic(ctx)
	err = logic.Complete(id)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(nil, "完成成功")
}

// Cancel 取消行动项
func (a *Action) Cancel(ctx *gin.Context) {
	a.Ctx = ctx
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewActionLogic(ctx)
	err = logic.Cancel(id)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(nil, "取消成功")
}

// Statistics 获取统计信息
func (a *Action) Statistics(ctx *gin.Context) {
	a.Ctx = ctx

	logic := thinking.NewActionLogic(ctx)
	res, err := logic.Statistics()
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}
