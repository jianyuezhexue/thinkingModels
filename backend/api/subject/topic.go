package subject

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"thinkingModels/api"
	"thinkingModels/domain/subject/topic"
	"thinkingModels/logic/subject"
)

type Topic struct {
	api.Base
}

func NewTopic() *Topic {
	return &Topic{}
}

// Create 创建课题
// @Summary 创建课题
// @Description 创建一个新的思考课题
// @Tags 课题管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body topic.CreateTopic true "创建课题请求参数"
// @Success 200 {object} api.Response{data=topic.TopicInfo} "创建成功"
// @Failure 400 {object} api.Response "参数错误"
// @Failure 401 {object} api.Response "未登录"
// @Router /subject/topic [post]
func (a *Topic) Create(ctx *gin.Context) {
	req := &topic.CreateTopic{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := subject.NewTopicLogic(ctx)
	res, err := logic.Create(req)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "创建成功")
}

// Update 更新课题
func (a *Topic) Update(ctx *gin.Context) {
	req := &topic.UpdateTopic{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := subject.NewTopicLogic(ctx)
	res, err := logic.Update(req)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "更新成功")
}

// Get 查询课题详情
// @Summary 获取课题详情
// @Description 根据ID获取课题详细信息
// @Tags 课题管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "课题ID"
// @Success 200 {object} api.Response{data=topic.TopicInfo} "查询成功"
// @Failure 404 {object} api.Response "课题不存在"
// @Router /subject/topic/{id} [get]
func (a *Topic) Get(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		a.Error(err)
		return
	}

	logic := subject.NewTopicLogic(ctx)
	res, err := logic.Get(id)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "查询成功")
}

// List 查询课题列表
// @Summary 查询课题列表
// @Description 分页查询课题列表，支持多种筛选条件
// @Tags 课题管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码，默认1"
// @Param pageSize query int false "每页数量，默认10"
// @Param title query string false "标题模糊查询"
// @Param status query int false "状态：0-进行中，1-已完成，2-已归档"
// @Param priority query int false "优先级：1-低，2-中，3-高"
// @Success 200 {object} api.Response{data=topic.ListTopicResponse} "查询成功"
// @Router /subject/topic/list [get]
func (a *Topic) List(ctx *gin.Context) {
	req := &topic.SearchTopic{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := subject.NewTopicLogic(ctx)
	res, err := logic.List(req)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "查询列表成功")
}

// ListByUser 查询当前用户的课题列表
func (a *Topic) ListByUser(ctx *gin.Context) {
	req := &topic.SearchTopic{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	// 从上下文中获取当前用户ID
	userIDStr, exists := ctx.Get("currUserId")
	if exists {
		if userID, err := strconv.ParseUint(userIDStr.(string), 10, 64); err == nil {
			req.UserId = userID
		}
	}

	logic := subject.NewTopicLogic(ctx)
	res, err := logic.List(req)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "查询列表成功")
}

// Del 删除课题
func (a *Topic) Del(ctx *gin.Context) {
	req := &topic.DelTopic{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := subject.NewTopicLogic(ctx)
	_, err := logic.Del(req.Ids)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(nil, "删除成功")
}

// UpdateStatus 更新课题状态
func (a *Topic) UpdateStatus(ctx *gin.Context) {
	req := &topic.UpdateTopicStatus{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := subject.NewTopicLogic(ctx)
	res, err := logic.UpdateStatus(req)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "状态更新成功")
}

// SelectModel 为课题选择思维模型
func (a *Topic) SelectModel(ctx *gin.Context) {
	req := &topic.UpdateTopicModel{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := subject.NewTopicLogic(ctx)
	res, err := logic.SelectModel(req)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "模型选择成功")
}

// RemoveModel 移除课题的思维模型
func (a *Topic) RemoveModel(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		a.Error(err)
		return
	}

	logic := subject.NewTopicLogic(ctx)
	res, err := logic.RemoveModel(id)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "模型移除成功")
}

// Complete 完成课题
func (a *Topic) Complete(ctx *gin.Context) {
	req := &topic.CompleteTopic{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := subject.NewTopicLogic(ctx)
	res, err := logic.Complete(req)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "课题已完成")
}

// Archive 归档课题
func (a *Topic) Archive(ctx *gin.Context) {
	req := &topic.ArchiveTopic{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := subject.NewTopicLogic(ctx)
	res, err := logic.Archive(req)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "课题已归档")
}

// Reopen 重新打开课题
func (a *Topic) Reopen(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		a.Error(err)
		return
	}

	logic := subject.NewTopicLogic(ctx)
	res, err := logic.Reopen(id)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "课题已重新打开")
}

// GetStatistics 获取课题统计
func (a *Topic) GetStatistics(ctx *gin.Context) {
	userIDStr, exists := ctx.Get("currUserId")
	var userId uint64 = 0
	if exists {
		if userID, err := strconv.ParseUint(userIDStr.(string), 10, 64); err == nil {
			userId = userID
		}
	}

	logic := subject.NewTopicLogic(ctx)
	res, err := logic.GetStatistics(userId)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "查询成功")
}

// parseTime 解析时间字符串
func parseTime(timeStr string) (time.Time, error) {
	if timeStr == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, timeStr)
}
