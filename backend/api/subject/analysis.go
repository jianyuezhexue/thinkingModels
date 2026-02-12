package subject

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"thinkingModels/api"
	"thinkingModels/domain/subject/analysis"
	"thinkingModels/logic/subject"
)

type Analysis struct {
	api.Base
}

func NewAnalysis() *Analysis {
	return &Analysis{}
}

// Create 创建分析记录
func (a *Analysis) Create(ctx *gin.Context) {
	req := &analysis.CreateAnalysis{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := subject.NewAnalysisLogic(ctx)
	res, err := logic.Create(req)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "创建成功")
}

// Update 更新分析记录
func (a *Analysis) Update(ctx *gin.Context) {
	req := &analysis.UpdateAnalysis{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := subject.NewAnalysisLogic(ctx)
	res, err := logic.Update(req)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "更新成功")
}

// SaveWithAi 保存分析记录并更新AI结果
func (a *Analysis) SaveWithAi(ctx *gin.Context) {
	req := &analysis.SaveAnalysisWithAi{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := subject.NewAnalysisLogic(ctx)
	res, err := logic.SaveWithAi(req)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "保存成功")
}

// Get 查询分析记录详情
func (a *Analysis) Get(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		a.Error(err)
		return
	}

	logic := subject.NewAnalysisLogic(ctx)
	res, err := logic.Get(id)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "查询成功")
}

// GetCurrent 获取当前分析记录
func (a *Analysis) GetCurrent(ctx *gin.Context) {
	req := &analysis.GetCurrentAnalysisRequest{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := subject.NewAnalysisLogic(ctx)
	res, err := logic.GetCurrentByTopic(req.TopicId)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "查询成功")
}

// GetLatest 获取课题最新分析记录
func (a *Analysis) GetLatest(ctx *gin.Context) {
	req := &analysis.GetLatestAnalysisRequest{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := subject.NewAnalysisLogic(ctx)
	res, err := logic.GetLatestByTopic(req.TopicId)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "查询成功")
}

// List 查询分析记录列表
func (a *Analysis) List(ctx *gin.Context) {
	req := &analysis.SearchAnalysis{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := subject.NewAnalysisLogic(ctx)
	res, err := logic.List(req)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "查询列表成功")
}

// ListByTopic 查询课题的所有分析记录
func (a *Analysis) ListByTopic(ctx *gin.Context) {
	topicIdStr := ctx.Param("topicId")
	topicId, err := strconv.ParseUint(topicIdStr, 10, 64)
	if err != nil {
		a.Error(err)
		return
	}

	logic := subject.NewAnalysisLogic(ctx)
	res, err := logic.GetByTopic(topicId)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "查询成功")
}

// GetHistory 获取分析记录历史（同一课题同一模型）
func (a *Analysis) GetHistory(ctx *gin.Context) {
	topicIdStr := ctx.Param("topicId")
	modelIdStr := ctx.Param("modelId")

	topicId, err := strconv.ParseUint(topicIdStr, 10, 64)
	if err != nil {
		a.Error(err)
		return
	}

	modelId, err := strconv.ParseUint(modelIdStr, 10, 64)
	if err != nil {
		a.Error(err)
		return
	}

	logic := subject.NewAnalysisLogic(ctx)
	res, err := logic.GetHistory(topicId, modelId)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "查询成功")
}

// Del 删除分析记录
func (a *Analysis) Del(ctx *gin.Context) {
	req := &analysis.DelAnalysis{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := subject.NewAnalysisLogic(ctx)
	_, err := logic.Del(req.Ids)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(nil, "删除成功")
}

// SetCurrent 设置当前版本
func (a *Analysis) SetCurrent(ctx *gin.Context) {
	req := &analysis.SetCurrentVersion{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := subject.NewAnalysisLogic(ctx)
	res, err := logic.SetCurrent(req)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "设置成功")
}

// ListByUser 查询当前用户的分析记录
func (a *Analysis) ListByUser(ctx *gin.Context) {
	req := &analysis.SearchAnalysis{}
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

	logic := subject.NewAnalysisLogic(ctx)
	res, err := logic.List(req)
	if err != nil {
		a.Error(err)
		return
	}

	a.Success(res, "查询列表成功")
}
