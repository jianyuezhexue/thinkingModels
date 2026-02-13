package thinking

import (
	"thinkingModels/domain/thinking/analysis"
	"thinkingModels/logic"

	"github.com/gin-gonic/gin"
)

// AnalysisLogic 分析业务逻辑
type AnalysisLogic struct {
	logic.BaseLogic
}

// NewAnalysisLogic 初始化AnalysisLogic
func NewAnalysisLogic(ctx *gin.Context) *AnalysisLogic {
	return &AnalysisLogic{BaseLogic: logic.BaseLogic{Ctx: ctx}}
}

// Create 创建分析
func (l *AnalysisLogic) Create(req *analysis.CreateAnalysis) (*analysis.AnalysisInfo, error) {
	entity := analysis.NewAnalysisEntity(l.Ctx)
	_, err := entity.SetData(req)
	if err != nil {
		return nil, err
	}

	if concreteEntity, ok := entity.(*analysis.AnalysisEntity); ok {
		concreteEntity.Version = 1
		concreteEntity.IsCurrent = true
		concreteEntity.Status = 0
	}

	if err := entity.Validate(); err != nil {
		return nil, err
	}

	res, err := entity.Create()
	if err != nil {
		return nil, err
	}

	return convertToAnalysisInfo(res), nil
}

// Update 更新分析
func (l *AnalysisLogic) Update(req *analysis.UpdateAnalysis) (*analysis.AnalysisInfo, error) {
	entity := analysis.NewAnalysisEntity(l.Ctx)
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

	return convertToAnalysisInfo(res), nil
}

// Get 查询分析详情
func (l *AnalysisLogic) Get(id uint64) (*analysis.AnalysisDetail, error) {
	entity := analysis.NewAnalysisEntity(l.Ctx)
	res, err := entity.LoadById(id)
	if err != nil {
		return nil, err
	}
	return convertToAnalysisDetail(res), nil
}

// List 查询分析列表
func (l *AnalysisLogic) List(req *analysis.SearchAnalysis) (*logic.ListReap, error) {
	entity := analysis.NewAnalysisEntity(l.Ctx)
	cond := entity.MakeConditon(*req)

	total, err := entity.Count(cond)
	if err != nil {
		return nil, err
	}

	list, err := entity.List(cond)
	if err != nil {
		return nil, err
	}

	analysisInfoList := make([]*analysis.AnalysisInfo, 0, len(list))
	for _, item := range list {
		analysisInfoList = append(analysisInfoList, convertToAnalysisInfo(item))
	}

	return &logic.ListReap{
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
		List:     analysisInfoList,
	}, nil
}

// ListMy 查询我的分析
func (l *AnalysisLogic) ListMy(req *analysis.SearchAnalysis) (*logic.ListReap, error) {
	return l.List(req)
}

// Del 删除分析
func (l *AnalysisLogic) Del(req *analysis.DelAnalysis) error {
	entity := analysis.NewAnalysisEntity(l.Ctx)
	return entity.Del(req.Ids...)
}

// SaveWithAi 保存并AI分析
func (l *AnalysisLogic) SaveWithAi(req *analysis.SaveWithAi) (*analysis.AnalysisInfo, error) {
	entity := analysis.NewAnalysisEntity(l.Ctx)

	if req.Id > 0 {
		_, err := entity.LoadById(req.Id)
		if err != nil {
			return nil, err
		}
	}

	_, err := entity.SetData(req)
	if err != nil {
		return nil, err
	}

	var res any
	if req.Id > 0 {
		res, err = entity.Update()
	} else {
		res, err = entity.Create()
	}

	if err != nil {
		return nil, err
	}

	return convertToAnalysisInfo(res), nil
}

// GetCurrent 获取当前版本分析
func (l *AnalysisLogic) GetCurrent(topicId, modelId uint64) (*analysis.AnalysisInfo, error) {
	entity := analysis.NewAnalysisEntity(l.Ctx)
	isCurrent := true
	cond := entity.MakeConditon(analysis.SearchAnalysis{
		TopicId:   topicId,
		ModelId:   modelId,
		IsCurrent: &isCurrent,
	})

	list, err := entity.List(cond)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	return convertToAnalysisInfo(list[0]), nil
}

// GetLatest 获取最新分析
func (l *AnalysisLogic) GetLatest(topicId uint64) (*analysis.AnalysisInfo, error) {
	entity := analysis.NewAnalysisEntity(l.Ctx)
	cond := entity.MakeConditon(analysis.SearchAnalysis{
		TopicId:  topicId,
		Page:     1,
		PageSize: 1,
	})

	list, err := entity.List(cond)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	return convertToAnalysisInfo(list[0]), nil
}

// ListByTopic 按课题查询分析
func (l *AnalysisLogic) ListByTopic(topicId uint64) ([]*analysis.AnalysisInfo, error) {
	entity := analysis.NewAnalysisEntity(l.Ctx)
	cond := entity.MakeConditon(analysis.SearchAnalysis{TopicId: topicId})

	list, err := entity.List(cond)
	if err != nil {
		return nil, err
	}

	result := make([]*analysis.AnalysisInfo, 0, len(list))
	for _, item := range list {
		result = append(result, convertToAnalysisInfo(item))
	}

	return result, nil
}

// GetHistory 获取历史版本
func (l *AnalysisLogic) GetHistory(topicId, modelId uint64) ([]*analysis.AnalysisInfo, error) {
	entity := analysis.NewAnalysisEntity(l.Ctx)
	cond := entity.MakeConditon(analysis.SearchAnalysis{
		TopicId: topicId,
		ModelId: modelId,
	})

	list, err := entity.List(cond)
	if err != nil {
		return nil, err
	}

	result := make([]*analysis.AnalysisInfo, 0, len(list))
	for _, item := range list {
		result = append(result, convertToAnalysisInfo(item))
	}

	return result, nil
}

// SetCurrent 设为当前版本
func (l *AnalysisLogic) SetCurrent(req *analysis.SetCurrentAnalysis) error {
	entity := analysis.NewAnalysisEntity(l.Ctx)
	_, err := entity.LoadById(req.Id)
	if err != nil {
		return err
	}

	if err := entity.SetAsCurrent(); err != nil {
		return err
	}

	_, err = entity.Update()
	return err
}

func convertToAnalysisInfo(entity any) *analysis.AnalysisInfo {
	e, ok := entity.(*analysis.AnalysisEntity)
	if !ok {
		return nil
	}
	return &analysis.AnalysisInfo{
		Id:        e.Id,
		TopicId:   e.TopicId,
		ModelId:   e.ModelId,
		Version:   e.Version,
		IsCurrent: e.IsCurrent,
		Status:    e.Status,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func convertToAnalysisDetail(entity any) *analysis.AnalysisDetail {
	e, ok := entity.(*analysis.AnalysisEntity)
	if !ok {
		return nil
	}
	return &analysis.AnalysisDetail{
		Id:           e.Id,
		TopicId:      e.TopicId,
		ModelId:      e.ModelId,
		InputContent: e.InputContent,
		AiResult:     e.AiResult,
		UserResult:   e.UserResult,
		Conclusion:   e.Conclusion,
		Version:      e.Version,
		IsCurrent:    e.IsCurrent,
		Status:       e.Status,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}
