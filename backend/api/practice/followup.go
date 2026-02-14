package thinking

import (
	"strconv"

	"thinkingModels/api"
	"thinkingModels/domain/practice/followup"
	"thinkingModels/logic/practice"

	"github.com/gin-gonic/gin"
)

type FollowUp struct {
	api.Base
}

func NewFollowUp() *FollowUp {
	return &FollowUp{}
}

// Create 创建跟进记录
func (f *FollowUp) Create(ctx *gin.Context) {
	f.Ctx = ctx
	req := &followup.CreateFollowUp{}
	if err := f.Bind(ctx, req); err != nil {
		f.Error(err)
		return
	}

	logic := thinking.NewFollowUpLogic(ctx)
	res, err := logic.Create(req)
	if err != nil {
		f.Error(err)
		return
	}
	f.Success(res, "创建成功")
}

// Update 更新跟进记录
func (f *FollowUp) Update(ctx *gin.Context) {
	f.Ctx = ctx
	req := &followup.UpdateFollowUp{}
	if err := f.Bind(ctx, req); err != nil {
		f.Error(err)
		return
	}

	logic := thinking.NewFollowUpLogic(ctx)
	res, err := logic.Update(req)
	if err != nil {
		f.Error(err)
		return
	}
	f.Success(res, "更新成功")
}

// Get 查询跟进记录详情
func (f *FollowUp) Get(ctx *gin.Context) {
	f.Ctx = ctx
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		f.Error(err)
		return
	}

	logic := thinking.NewFollowUpLogic(ctx)
	res, err := logic.Get(id)
	if err != nil {
		f.Error(err)
		return
	}
	f.Success(res, "查询成功")
}

// Del 删除跟进记录
func (f *FollowUp) Del(ctx *gin.Context) {
	f.Ctx = ctx
	req := &followup.DelFollowUp{}
	if err := f.Bind(ctx, req); err != nil {
		f.Error(err)
		return
	}

	logic := thinking.NewFollowUpLogic(ctx)
	err := logic.Del(req)
	if err != nil {
		f.Error(err)
		return
	}
	f.Success(nil, "删除成功")
}

// ListByAction 按行动项查询跟进记录
func (f *FollowUp) ListByAction(ctx *gin.Context) {
	f.Ctx = ctx
	actionIdStr := ctx.Param("actionId")
	actionId, err := strconv.ParseUint(actionIdStr, 10, 64)
	if err != nil {
		f.Error(err)
		return
	}

	logic := thinking.NewFollowUpLogic(ctx)
	res, err := logic.ListByAction(actionId)
	if err != nil {
		f.Error(err)
		return
	}
	f.Success(res, "查询成功")
}
