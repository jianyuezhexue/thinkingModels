package thinking

import (
	"strconv"

	"thinkingModels/api"
	"thinkingModels/domain/practice/analysis"
	"thinkingModels/logic/practice"

	"github.com/gin-gonic/gin"
)

type Analysis struct {
	api.Base
}

func NewAnalysis() *Analysis {
	return &Analysis{}
}

// Create 创建分析
func (a *Analysis) Create(ctx *gin.Context) {
	a.Ctx = ctx
	req := &analysis.CreateAnalysis{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewAnalysisLogic(ctx)
	res, err := logic.Create(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "创建成功")
}

// Update 更新分析
func (a *Analysis) Update(ctx *gin.Context) {
	a.Ctx = ctx
	req := &analysis.UpdateAnalysis{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewAnalysisLogic(ctx)
	res, err := logic.Update(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "更新成功")
}

// Get 查询分析详情
func (a *Analysis) Get(ctx *gin.Context) {
	a.Ctx = ctx
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewAnalysisLogic(ctx)
	res, err := logic.Get(id)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}

// List 查询分析列表
func (a *Analysis) List(ctx *gin.Context) {
	a.Ctx = ctx
	req := &analysis.SearchAnalysis{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewAnalysisLogic(ctx)
	res, err := logic.List(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}

// ListMy 查询我的分析
func (a *Analysis) ListMy(ctx *gin.Context) {
	a.Ctx = ctx
	req := &analysis.SearchAnalysis{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewAnalysisLogic(ctx)
	res, err := logic.ListMy(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}

// Del 删除分析
func (a *Analysis) Del(ctx *gin.Context) {
	a.Ctx = ctx
	req := &analysis.DelAnalysis{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewAnalysisLogic(ctx)
	err := logic.Del(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(nil, "删除成功")
}

// SaveWithAi 保存并AI分析
func (a *Analysis) SaveWithAi(ctx *gin.Context) {
	a.Ctx = ctx
	req := &analysis.SaveWithAi{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewAnalysisLogic(ctx)
	res, err := logic.SaveWithAi(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "分析成功")
}

// GetCurrent 获取当前版本分析
func (a *Analysis) GetCurrent(ctx *gin.Context) {
	a.Ctx = ctx
	topicIdStr := ctx.Query("topicId")
	modelIdStr := ctx.Query("modelId")
	topicId, _ := strconv.ParseUint(topicIdStr, 10, 64)
	modelId, _ := strconv.ParseUint(modelIdStr, 10, 64)

	logic := thinking.NewAnalysisLogic(ctx)
	res, err := logic.GetCurrent(topicId, modelId)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}

// GetLatest 获取最新分析
func (a *Analysis) GetLatest(ctx *gin.Context) {
	a.Ctx = ctx
	topicIdStr := ctx.Query("topicId")
	topicId, _ := strconv.ParseUint(topicIdStr, 10, 64)

	logic := thinking.NewAnalysisLogic(ctx)
	res, err := logic.GetLatest(topicId)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}

// ListByTopic 按课题查询分析
func (a *Analysis) ListByTopic(ctx *gin.Context) {
	a.Ctx = ctx
	topicIdStr := ctx.Param("topicId")
	topicId, err := strconv.ParseUint(topicIdStr, 10, 64)
	if err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewAnalysisLogic(ctx)
	res, err := logic.ListByTopic(topicId)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}

// GetHistory 获取历史版本
func (a *Analysis) GetHistory(ctx *gin.Context) {
	a.Ctx = ctx
	topicIdStr := ctx.Param("topicId")
	modelIdStr := ctx.Param("modelId")
	topicId, _ := strconv.ParseUint(topicIdStr, 10, 64)
	modelId, _ := strconv.ParseUint(modelIdStr, 10, 64)

	logic := thinking.NewAnalysisLogic(ctx)
	res, err := logic.GetHistory(topicId, modelId)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(res, "查询成功")
}

// SetCurrent 设为当前版本
func (a *Analysis) SetCurrent(ctx *gin.Context) {
	a.Ctx = ctx
	req := &analysis.SetCurrentAnalysis{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	logic := thinking.NewAnalysisLogic(ctx)
	err := logic.SetCurrent(req)
	if err != nil {
		a.Error(err)
		return
	}
	a.Success(nil, "设置成功")
}
