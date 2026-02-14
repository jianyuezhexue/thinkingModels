package thinking

import (
	"strconv"

	"thinkingModels/api"
	"thinkingModels/domain/practice/topic"
	"thinkingModels/logic/practice"

	"github.com/gin-gonic/gin"
)

type Topic struct {
	api.Base
}

func NewTopic() *Topic {
	return &Topic{}
}

// Create 创建课题
func (a *Topic) Create(ctx *gin.Context) {
	a.Ctx = ctx
	req := &topic.CreateTopic{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewTopicLogic(ctx)
	res, err := logic.Create(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "创建成功")
}

// Update 更新课题
func (a *Topic) Update(ctx *gin.Context) {
	a.Ctx = ctx
	req := &topic.UpdateTopic{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewTopicLogic(ctx)
	res, err := logic.Update(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "更新成功")
}

// Get 查询课题详情
func (a *Topic) Get(ctx *gin.Context) {
	a.Ctx = ctx
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewTopicLogic(ctx)
	res, err := logic.Get(id)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}

// List 查询课题列表
func (a *Topic) List(ctx *gin.Context) {
	a.Ctx = ctx
	req := &topic.SearchTopic{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewTopicLogic(ctx)
	res, err := logic.List(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}

// ListMy 查询我的课题
func (a *Topic) ListMy(ctx *gin.Context) {
	a.Ctx = ctx
	req := &topic.SearchTopic{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewTopicLogic(ctx)
	res, err := logic.ListMy(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}

// Del 删除课题
func (a *Topic) Del(ctx *gin.Context) {
	a.Ctx = ctx
	req := &topic.DelTopic{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewTopicLogic(ctx)
	err := logic.Del(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(nil, "删除成功")
}

// SelectModel 为课题选用模型
func (a *Topic) SelectModel(ctx *gin.Context) {
	a.Ctx = ctx
	req := &topic.SelectModelForTopic{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewTopicLogic(ctx)
	err := logic.SelectModel(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(nil, "选用成功")
}

// RemoveModel 移除课题模型
func (a *Topic) RemoveModel(ctx *gin.Context) {
	a.Ctx = ctx
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewTopicLogic(ctx)
	err = logic.RemoveModel(id)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(nil, "移除成功")
}

// UpdateStatus 更新课题状态
func (a *Topic) UpdateStatus(ctx *gin.Context) {
	a.Ctx = ctx
	req := &topic.UpdateTopicStatus{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewTopicLogic(ctx)
	err := logic.UpdateStatus(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(nil, "更新成功")
}

// Complete 完成课题
func (a *Topic) Complete(ctx *gin.Context) {
	a.Ctx = ctx
	req := &topic.CompleteTopic{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewTopicLogic(ctx)
	err := logic.Complete(req.Id)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(nil, "完成成功")
}

// Archive 归档课题
func (a *Topic) Archive(ctx *gin.Context) {
	a.Ctx = ctx
	req := &topic.ArchiveTopic{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewTopicLogic(ctx)
	err := logic.Archive(req.Id)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(nil, "归档成功")
}

// Reopen 重新打开课题
func (a *Topic) Reopen(ctx *gin.Context) {
	a.Ctx = ctx
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewTopicLogic(ctx)
	err = logic.Reopen(id)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(nil, "重新打开成功")
}

// Statistics 获取课题统计
func (a *Topic) Statistics(ctx *gin.Context) {
	a.Ctx = ctx
	logic := thinking.NewTopicLogic(ctx)
	res, err := logic.Statistics()
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}
