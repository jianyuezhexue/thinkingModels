package thinking

import (
	"thinkingModels/domain/practice/action"
	"thinkingModels/logic"

	"github.com/gin-gonic/gin"
)

// ActionLogic 行动项业务逻辑
type ActionLogic struct {
	logic.BaseLogic
}

// NewActionLogic 初始化ActionLogic
func NewActionLogic(ctx *gin.Context) *ActionLogic {
	return &ActionLogic{BaseLogic: logic.BaseLogic{Ctx: ctx}}
}

// Create 创建行动项
func (l *ActionLogic) Create(req *action.CreateAction) (*action.ActionInfo, error) {
	entity := action.NewActionEntity(l.Ctx)
	_, err := entity.SetData(req)
	if err != nil {
		return nil, err
	}

	if concreteEntity, ok := entity.(*action.ActionEntity); ok {
		concreteEntity.Status = 0
		concreteEntity.Progress = 0
	}

	if err := entity.Validate(); err != nil {
		return nil, err
	}

	res, err := entity.Create()
	if err != nil {
		return nil, err
	}

	return convertToActionInfo(res), nil
}

// Update 更新行动项
func (l *ActionLogic) Update(req *action.UpdateAction) (*action.ActionInfo, error) {
	entity := action.NewActionEntity(l.Ctx)
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

	return convertToActionInfo(res), nil
}

// Get 查询行动项详情
func (l *ActionLogic) Get(id uint64) (*action.ActionDetail, error) {
	entity := action.NewActionEntity(l.Ctx)
	res, err := entity.LoadById(id)
	if err != nil {
		return nil, err
	}
	return convertToActionDetail(res), nil
}

// List 查询行动项列表
func (l *ActionLogic) List(req *action.SearchAction) (*logic.ListReap, error) {
	entity := action.NewActionEntity(l.Ctx)
	cond := entity.MakeConditon(*req)

	total, err := entity.Count(cond)
	if err != nil {
		return nil, err
	}

	list, err := entity.List(cond)
	if err != nil {
		return nil, err
	}

	actionInfoList := make([]*action.ActionInfo, 0, len(list))
	for _, item := range list {
		actionInfoList = append(actionInfoList, convertToActionInfo(item))
	}

	return &logic.ListReap{
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
		List:     actionInfoList,
	}, nil
}

// ListMy 查询我的行动项
func (l *ActionLogic) ListMy(req *action.SearchAction) (*logic.ListReap, error) {
	return l.List(req)
}

// Del 删除行动项
func (l *ActionLogic) Del(req *action.DelAction) error {
	entity := action.NewActionEntity(l.Ctx)
	return entity.Del(req.Ids...)
}

// CreateFromAnalysis 从分析创建行动项
func (l *ActionLogic) CreateFromAnalysis(req *action.CreateActionsFromAnalysis) ([]*action.ActionInfo, error) {
	result := make([]*action.ActionInfo, 0, len(req.Actions))

	for _, actionReq := range req.Actions {
		entity := action.NewActionEntity(l.Ctx)
		_, err := entity.SetData(actionReq)
		if err != nil {
			return nil, err
		}

		if concreteEntity, ok := entity.(*action.ActionEntity); ok {
			concreteEntity.TopicId = req.TopicId
			concreteEntity.AnalysisId = req.AnalysisId
			concreteEntity.Status = 0
			concreteEntity.Progress = 0
		}

		res, err := entity.Create()
		if err != nil {
			return nil, err
		}

		result = append(result, convertToActionInfo(res))
	}

	return result, nil
}

// ListByTopic 按课题查询行动项
func (l *ActionLogic) ListByTopic(topicId uint64) ([]*action.ActionInfo, error) {
	entity := action.NewActionEntity(l.Ctx)
	cond := entity.MakeConditon(action.SearchAction{TopicId: topicId})

	list, err := entity.List(cond)
	if err != nil {
		return nil, err
	}

	result := make([]*action.ActionInfo, 0, len(list))
	for _, item := range list {
		result = append(result, convertToActionInfo(item))
	}

	return result, nil
}

// ListByAnalysis 按分析查询行动项
func (l *ActionLogic) ListByAnalysis(analysisId uint64) ([]*action.ActionInfo, error) {
	entity := action.NewActionEntity(l.Ctx)
	cond := entity.MakeConditon(action.SearchAction{AnalysisId: analysisId})

	list, err := entity.List(cond)
	if err != nil {
		return nil, err
	}

	result := make([]*action.ActionInfo, 0, len(list))
	for _, item := range list {
		result = append(result, convertToActionInfo(item))
	}

	return result, nil
}

// UpdateProgress 更新进度
func (l *ActionLogic) UpdateProgress(req *action.UpdateProgress) error {
	entity := action.NewActionEntity(l.Ctx)
	_, err := entity.LoadById(req.Id)
	if err != nil {
		return err
	}

	if err := entity.UpdateProgress(req.Progress, req.Note); err != nil {
		return err
	}

	_, err = entity.Update()
	return err
}

// Complete 完成行动项
func (l *ActionLogic) Complete(id uint64) error {
	entity := action.NewActionEntity(l.Ctx)
	_, err := entity.LoadById(id)
	if err != nil {
		return err
	}

	if err := entity.MarkComplete(); err != nil {
		return err
	}

	_, err = entity.Update()
	return err
}

// Cancel 取消行动项
func (l *ActionLogic) Cancel(id uint64) error {
	entity := action.NewActionEntity(l.Ctx)
	_, err := entity.LoadById(id)
	if err != nil {
		return err
	}

	if err := entity.Cancel(); err != nil {
		return err
	}

	_, err = entity.Update()
	return err
}

// Statistics 获取统计信息
func (l *ActionLogic) Statistics() (*action.ActionStatistics, error) {
	entity := action.NewActionEntity(l.Ctx)

	pendingCond := entity.MakeConditon(action.SearchAction{Status: intPtrAction(0)})
	pendingCount, _ := entity.Count(pendingCond)

	progressCond := entity.MakeConditon(action.SearchAction{Status: intPtrAction(1)})
	progressCount, _ := entity.Count(progressCond)

	completeCond := entity.MakeConditon(action.SearchAction{Status: intPtrAction(2)})
	completeCount, _ := entity.Count(completeCond)

	cancelCond := entity.MakeConditon(action.SearchAction{Status: intPtrAction(3)})
	cancelCount, _ := entity.Count(cancelCond)

	overdueCond := entity.MakeConditon(action.SearchAction{IsOverdue: boolPtr(true)})
	overdueCount, _ := entity.Count(overdueCond)

	return &action.ActionStatistics{
		PendingCount:  pendingCount,
		ProgressCount: progressCount,
		CompleteCount: completeCount,
		CancelCount:   cancelCount,
		OverdueCount:  overdueCount,
	}, nil
}

func intPtrAction(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}

func convertToActionInfo(entity any) *action.ActionInfo {
	e, ok := entity.(*action.ActionEntity)
	if !ok {
		return nil
	}
	return &action.ActionInfo{
		Id:            e.Id,
		TopicId:       e.TopicId,
		AnalysisId:    e.AnalysisId,
		Title:         e.Title,
		Description:   e.Description,
		Priority:      e.Priority,
		Status:        e.Status,
		Progress:      e.Progress,
		Deadline:      e.Deadline,
		FollowUpCount: e.FollowUpCount,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
}

func convertToActionDetail(entity any) *action.ActionDetail {
	e, ok := entity.(*action.ActionEntity)
	if !ok {
		return nil
	}
	return &action.ActionDetail{
		Id:             e.Id,
		TopicId:        e.TopicId,
		AnalysisId:     e.AnalysisId,
		Title:          e.Title,
		Description:    e.Description,
		Priority:       e.Priority,
		Status:         e.Status,
		Progress:       e.Progress,
		ExpectedResult: e.ExpectedResult,
		ActualResult:   e.ActualResult,
		Deadline:       e.Deadline,
		CompletedAt:    e.CompletedAt,
		FollowUpCount:  e.FollowUpCount,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}
