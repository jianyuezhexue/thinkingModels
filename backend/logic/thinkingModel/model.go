package thinkingModel

import (
	"errors"

	"github.com/gin-gonic/gin"
	"thinkingModels/domain/thinkingModel/model"
	"thinkingModels/logic"
)

// ModelLogic 思维模型业务逻辑
type ModelLogic struct {
	logic.BaseLogic
}

// NewModelLogic 初始化ModelLogic
func NewModelLogic(ctx *gin.Context) *ModelLogic {
	return &ModelLogic{BaseLogic: logic.BaseLogic{Ctx: ctx}}
}

// Create 创建思维模型
func (l *ModelLogic) Create(req *model.CreateModel) (*model.ModelInfo, error) {
	entity := model.NewModelEntity(l.Ctx)
	_, err := entity.SetData(req)
	if err != nil {
		return nil, err
	}

	if concreteEntity, ok := entity.(*model.ModelEntity); ok {
		concreteEntity.Status = 0
		concreteEntity.Version = "1.0.0"
		if concreteEntity.Difficulty == 0 {
			concreteEntity.Difficulty = 1
		}
		if concreteEntity.EstimatedTime == 0 {
			concreteEntity.EstimatedTime = 30
		}
	}

	if err := entity.Validate(); err != nil {
		return nil, err
	}

	res, err := entity.Create()
	if err != nil {
		return nil, err
	}

	return convertToModelInfo(res), nil
}

// Update 更新思维模型
func (l *ModelLogic) Update(req *model.UpdateModel) (*model.ModelInfo, error) {
	entity := model.NewModelEntity(l.Ctx)
	_, err := entity.LoadById(req.Id)
	if err != nil {
		return nil, err
	}

	_, err = entity.SetData(req)
	if err != nil {
		return nil, err
	}

	if err := entity.Validate(); err != nil {
		return nil, err
	}

	res, err := entity.Update()
	if err != nil {
		return nil, err
	}

	return convertToModelInfo(res), nil
}

// Get 查询思维模型详情
func (l *ModelLogic) Get(id uint64) (*model.ModelDetail, error) {
	entity := model.NewModelEntity(l.Ctx)
	res, err := entity.LoadById(id)
	if err != nil {
		return nil, err
	}
	return convertToModelDetail(res), nil
}

// GetByCode 按编码查询（已废弃，请使用 Get）
func (l *ModelLogic) GetByCode(code string) (*model.ModelDetail, error) {
	return nil, errors.New("该方法已废弃")
}

// List 查询思维模型列表
func (l *ModelLogic) List(req *model.SearchModel) (*model.ListModelResponse, error) {
	entity := model.NewModelEntity(l.Ctx)
	cond := entity.MakeConditon(*req)

	total, err := entity.Count(cond)
	if err != nil {
		return nil, err
	}

	list, err := entity.List(cond)
	if err != nil {
		return nil, err
	}

	modelInfoList := make([]*model.ModelInfo, 0, len(list))
	for _, item := range list {
		modelInfoList = append(modelInfoList, convertToModelInfo(item))
	}

	return &model.ListModelResponse{
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
		List:     modelInfoList,
	}, nil
}

// ListMy 查询我的思维模型
func (l *ModelLogic) ListMy(req *model.SearchModel) (*model.ListModelResponse, error) {
	return l.List(req)
}

// Del 删除思维模型
func (l *ModelLogic) Del(ids []uint64) error {
	entity := model.NewModelEntity(l.Ctx)
	return entity.Del(ids...)
}

// Publish 发布思维模型（提交审核）
func (l *ModelLogic) Publish(req *model.PublishModel) (*model.ModelInfo, error) {
	entity := model.NewModelEntity(l.Ctx)
	_, err := entity.LoadById(req.Id)
	if err != nil {
		return nil, err
	}

	// 提交审核而不是直接发布
	if err := entity.SubmitForReview(); err != nil {
		return nil, err
	}

	res, err := entity.Update()
	if err != nil {
		return nil, err
	}

	return convertToModelInfo(res), nil
}

// Unpublish 下架思维模型
func (l *ModelLogic) Unpublish(id uint64) (*model.ModelInfo, error) {
	entity := model.NewModelEntity(l.Ctx)
	_, err := entity.LoadById(id)
	if err != nil {
		return nil, err
	}

	if err := entity.Unpublish(); err != nil {
		return nil, err
	}

	res, err := entity.Update()
	if err != nil {
		return nil, err
	}

	return convertToModelInfo(res), nil
}

// SubmitForReview 提交审核
func (l *ModelLogic) SubmitForReview(id uint64) (*model.ModelInfo, error) {
	entity := model.NewModelEntity(l.Ctx)
	_, err := entity.LoadById(id)
	if err != nil {
		return nil, err
	}

	if err := entity.SubmitForReview(); err != nil {
		return nil, err
	}

	res, err := entity.Update()
	if err != nil {
		return nil, err
	}

	return convertToModelInfo(res), nil
}

// Review 审核模型
func (l *ModelLogic) Review(req *model.ReviewModel) (*model.ModelInfo, error) {
	entity := model.NewModelEntity(l.Ctx)
	_, err := entity.LoadById(req.Id)
	if err != nil {
		return nil, err
	}

	concreteEntity, ok := entity.(*model.ModelEntity)
	if !ok {
		return nil, errors.New("模型实体转换失败")
	}

	// 获取当前用户信息
	currUserId, _ := l.Ctx.Get("currUserId")
	currUserName, _ := l.Ctx.Get("currUserName")

	var reviewerId uint64
	if id, ok := currUserId.(uint64); ok {
		reviewerId = id
	}
	reviewerName := ""
	if name, ok := currUserName.(string); ok {
		reviewerName = name
	}

	if req.Approved {
		concreteEntity.Approve(reviewerId, reviewerName, req.Note)
	} else {
		if err := concreteEntity.Reject(reviewerId, reviewerName, req.Note); err != nil {
			return nil, err
		}
	}

	res, err := entity.Update()
	if err != nil {
		return nil, err
	}

	return convertToModelInfo(res), nil
}

// Fork 派生思维模型
func (l *ModelLogic) Fork(req *model.ForkModel) (*model.ModelInfo, error) {
	entity := model.NewModelEntity(l.Ctx)
	_, err := entity.LoadById(req.SourceModelId)
	if err != nil {
		return nil, err
	}

	forkedEntity := entity.Fork()
	if concreteEntity, ok := forkedEntity.(*model.ModelEntity); ok {
		concreteEntity.Name = req.Name
		concreteEntity.Description = req.Description
	}

	res, err := forkedEntity.Create()
	if err != nil {
		return nil, err
	}

	return convertToModelInfo(res), nil
}

func convertToModelInfo(entity any) *model.ModelInfo {
	e, ok := entity.(*model.ModelEntity)
	if !ok {
		return nil
	}
	return &model.ModelInfo{
		Id:            e.Id,
		Name:          e.Name,
		Description:   e.Description,
		Overview:      e.Overview,
		Icon:          e.Icon,
		CategoryId:    e.CategoryId,
		Difficulty:    e.Difficulty,
		EstimatedTime: e.EstimatedTime,
		Status:        e.Status,
		Version:       e.Version,
		IsOfficial:    e.IsOfficial,
		Stats: model.ModelStats{
			UsageCount:   int(e.UsageCount),
			AdoptCount:   int(e.AdoptCount),
			LikeCount:    int(e.LikeCount),
			CommentCount: int(e.CommentCount),
		},
		ReviewNote:   e.ReviewNote,
		ReviewerName: e.ReviewerName,
		CreatedAt:    e.CreatedAt.String(),
		UpdatedAt:    e.UpdatedAt.String(),
	}
}

func convertToModelDetail(entity any) *model.ModelDetail {
	e, ok := entity.(*model.ModelEntity)
	if !ok {
		return nil
	}
	info := convertToModelInfo(entity)
	if info == nil {
		return nil
	}
	return &model.ModelDetail{
		ModelInfo: *info,
		Content:   e.Content,
	}
}