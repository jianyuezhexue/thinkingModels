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
// @Summary 创建思维模型
// @Description 创建一个新的思维模型
// @Tags 思维模型
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.CreateModel true "创建思维模型请求参数"
// @Success 200 {object} api.Response{data=model.ModelDetail} "创建成功"
// @Failure 400 {object} api.Response "参数错误"
// @Failure 401 {object} api.Response "未登录"
// @Router /thinking/model [post]
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
// @Summary 获取思维模型详情
// @Description 根据ID获取思维模型详细信息
// @Tags 思维模型
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "模型ID"
// @Success 200 {object} api.Response{data=model.ModelDetail} "查询成功"
// @Failure 404 {object} api.Response "模型不存在"
// @Router /thinking/model/{id} [get]
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
// @Summary 获取思维模型列表
// @Description 分页查询思维模型列表，支持多种筛选条件
// @Tags 思维模型
// @Accept json
// @Produce json
// @Param page query int false "页码，默认1"
// @Param pageSize query int false "每页数量，默认10"
// @Param name query string false "名称模糊查询"
// @Param categoryId query int false "分类ID"
// @Param difficulty query int false "难度：1-简单，2-中等，3-困难"
// @Param isFree query bool false "是否免费"
// @Success 200 {object} api.Response{data=model.ListModelResponse} "查询成功"
// @Router /thinking/model/list [get]
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
// @Summary 发布思维模型
// @Description 将思维模型发布到市场
// @Tags 思维模型
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.PublishModel true "发布请求参数"
// @Success 200 {object} api.Response{data=model.ModelInfo} "发布成功"
// @Failure 400 {object} api.Response "参数错误"
// @Failure 401 {object} api.Response "未登录"
// @Router /thinking/model/publish [post]
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
