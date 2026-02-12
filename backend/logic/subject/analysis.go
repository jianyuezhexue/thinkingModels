package subject

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"thinkingModels/domain/market/model"
	"thinkingModels/domain/subject/analysis"
	"thinkingModels/logic"
)

// AnalysisLogic 分析记录业务逻辑
type AnalysisLogic struct {
	logic.BaseLogic
}

// NewAnalysisLogic 初始化AnalysisLogic
func NewAnalysisLogic(ctx *gin.Context) *AnalysisLogic {
	return &AnalysisLogic{BaseLogic: logic.BaseLogic{Ctx: ctx}}
}

// Create 创建分析记录
func (l *AnalysisLogic) Create(req *analysis.CreateAnalysis) (*analysis.CreateAnalysisResponse, error) {
	// 实例化模型
	entity := analysis.NewAnalysisEntity(l.Ctx)

	// 验证课题是否存在
	topicLogic := NewTopicLogic(l.Ctx)
	_, err := topicLogic.GetSimple(req.TopicId)
	if err != nil {
		return nil, errors.New("课题不存在")
	}

	// 验证思维模型是否存在
	modelEntity := model.NewModelEntity(l.Ctx)
	modelData, err := modelEntity.LoadById(req.ModelId)
	if err != nil {
		return nil, errors.New("思维模型不存在")
	}

	// 计算版本号
	latestVersion := l.getLatestVersion(req.TopicId, req.ModelId)
	newVersion := latestVersion + 1

	// 设置数据
	if concreteEntity, ok := entity.(*analysis.AnalysisEntity); ok {
		concreteEntity.TopicId = req.TopicId
		concreteEntity.ModelId = req.ModelId
		concreteEntity.ModelName = modelData.Name // 保存模型名称快照
		concreteEntity.Content = req.Content
		concreteEntity.AiAnalysis = req.AiAnalysis
		concreteEntity.AiSuggestions = req.AiSuggestions
		concreteEntity.Version = newVersion

		// 新记录设置为当前版本，同时取消该课题该模型的其他记录的当前版本标记
		concreteEntity.IsCurrent = 1

		// 从上下文中获取当前用户ID
		if userIDStr, exists := l.Ctx.Get("currUserId"); exists {
			if userID, err := strconv.ParseUint(userIDStr.(string), 10, 64); err == nil {
				concreteEntity.UserId = userID
			}
		}
	}

	// 数据校验
	if err := entity.Validate(); err != nil {
		return nil, err
	}

	// 先将同一课题同一模型的其他记录设置为非当前
	if err := l.unsetCurrentVersion(req.TopicId, req.ModelId); err != nil {
		return nil, err
	}

	// 保存数据
	res, err := entity.Create()
	if err != nil {
		return nil, err
	}

	// 返回响应
	info := convertToAnalysisInfo(res)
	return &analysis.CreateAnalysisResponse{
		AnalysisInfo: *info,
		Version:      newVersion,
	}, nil
}

// Update 更新分析记录
func (l *AnalysisLogic) Update(req *analysis.UpdateAnalysis) (*analysis.AnalysisInfo, error) {
	entity := analysis.NewAnalysisEntity(l.Ctx)

	// 加载旧数据
	_, err := entity.LoadById(req.Id)
	if err != nil {
		return nil, err
	}

	// 更新数据
	if concreteEntity, ok := entity.(*analysis.AnalysisEntity); ok {
		concreteEntity.Content = req.Content
		if req.AiAnalysis != "" {
			concreteEntity.AiAnalysis = req.AiAnalysis
		}
		if req.AiSuggestions != "" {
			concreteEntity.AiSuggestions = req.AiSuggestions
		}
	}

	// 数据校验
	if err := entity.Validate(); err != nil {
		return nil, err
	}

	res, err := entity.Update()
	if err != nil {
		return nil, err
	}

	return convertToAnalysisInfo(res), nil
}

// SaveWithAi 保存分析记录并更新AI结果
func (l *AnalysisLogic) SaveWithAi(req *analysis.SaveAnalysisWithAi) (*analysis.AnalysisInfo, error) {
	entity := analysis.NewAnalysisEntity(l.Ctx)

	// 加载旧数据
	_, err := entity.LoadById(req.Id)
	if err != nil {
		return nil, err
	}

	// 更新数据
	if concreteEntity, ok := entity.(*analysis.AnalysisEntity); ok {
		concreteEntity.Content = req.Content
		concreteEntity.AiAnalysis = req.AiAnalysis
		concreteEntity.AiSuggestions = req.AiSuggestions
	}

	// 数据校验
	if err := entity.Validate(); err != nil {
		return nil, err
	}

	res, err := entity.Update()
	if err != nil {
		return nil, err
	}

	return convertToAnalysisInfo(res), nil
}

// Get 查询分析记录详情
func (l *AnalysisLogic) Get(id uint64) (*analysis.AnalysisDetail, error) {
	entity := analysis.NewAnalysisEntity(l.Ctx)
	res, err := entity.LoadById(id)
	if err != nil {
		return nil, err
	}

	return convertToAnalysisDetail(res), nil
}

// GetByTopic 查询课题的所有分析记录
func (l *AnalysisLogic) GetByTopic(topicId uint64) ([]*analysis.AnalysisInfo, error) {
	entity := analysis.NewAnalysisEntity(l.Ctx)
	search := &analysis.SearchAnalysis{TopicId: topicId, PageSize: 1000}
	cond := entity.MakeConditon(*search)

	list, err := entity.List(cond)
	if err != nil {
		return nil, err
	}

	infoList := make([]*analysis.AnalysisInfo, 0, len(list))
	for _, item := range list {
		infoList = append(infoList, convertToAnalysisInfo(item))
	}

	return infoList, nil
}

// GetCurrentByTopic 获取课题当前使用的分析记录
func (l *AnalysisLogic) GetCurrentByTopic(topicId uint64) (*analysis.AnalysisInfo, error) {
	entity := analysis.NewAnalysisEntity(l.Ctx)
	search := &analysis.SearchAnalysis{
		TopicId:   topicId,
		IsCurrent: 1,
		PageSize:  1,
	}
	cond := entity.MakeConditon(*search)

	list, err := entity.List(cond)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, errors.New("未找到当前分析记录")
	}

	return convertToAnalysisInfo(list[0]), nil
}

// GetLatestByTopic 获取课题最新的分析记录（按版本号）
func (l *AnalysisLogic) GetLatestByTopic(topicId uint64) (*analysis.AnalysisInfo, error) {
	entity := analysis.NewAnalysisEntity(l.Ctx)
	search := &analysis.SearchAnalysis{
		TopicId:  topicId,
		PageSize: 1000,
	}
	cond := entity.MakeConditon(*search)

	list, err := entity.List(cond)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, errors.New("未找到分析记录")
	}

	// 找出最新版本
	var latest *analysis.AnalysisEntity
	maxVersion := 0
	for _, item := range list {
		if item.Version > maxVersion {
			maxVersion = item.Version
			latest = item
		}
	}

	return convertToAnalysisInfo(latest), nil
}

// List 查询分析记录列表
func (l *AnalysisLogic) List(req *analysis.SearchAnalysis) (*analysis.ListAnalysisResponse, error) {
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

	infoList := make([]*analysis.AnalysisInfo, 0, len(list))
	for _, item := range list {
		infoList = append(infoList, convertToAnalysisInfo(item))
	}

	return &analysis.ListAnalysisResponse{
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
		List:     infoList,
	}, nil
}

// GetHistory 获取分析记录历史（同一课题同一模型的所有版本）
func (l *AnalysisLogic) GetHistory(topicId, modelId uint64) (*analysis.AnalysisHistory, error) {
	entity := analysis.NewAnalysisEntity(l.Ctx)
	search := &analysis.SearchAnalysis{
		TopicId:  topicId,
		ModelId:  modelId,
		PageSize: 1000,
	}
	cond := entity.MakeConditon(*search)

	list, err := entity.List(cond)
	if err != nil {
		return nil, err
	}

	infoList := make([]*analysis.AnalysisInfo, 0, len(list))
	var modelName string
	for _, item := range list {
		infoList = append(infoList, convertToAnalysisInfo(item))
		if modelName == "" {
			modelName = item.ModelName
		}
	}

	return &analysis.AnalysisHistory{
		TopicId:   topicId,
		ModelId:   modelId,
		ModelName: modelName,
		Versions:  infoList,
	}, nil
}

// SetCurrent 设置当前版本
func (l *AnalysisLogic) SetCurrent(req *analysis.SetCurrentVersion) (*analysis.AnalysisInfo, error) {
	entity := analysis.NewAnalysisEntity(l.Ctx)

	// 加载目标分析记录
	target, err := entity.LoadById(req.Id)
	if err != nil {
		return nil, err
	}

	// 验证topicId匹配
	if target.TopicId != req.TopicId {
		return nil, errors.New("课题ID不匹配")
	}

	// 先将同一课题同一模型的其他记录设置为非当前
	if err := l.unsetCurrentVersion(target.TopicId, target.ModelId); err != nil {
		return nil, err
	}

	// 设置当前记录为当前版本
	if concreteEntity, ok := entity.(*analysis.AnalysisEntity); ok {
		concreteEntity.IsCurrent = 1
	}

	res, err := entity.Update()
	if err != nil {
		return nil, err
	}

	return convertToAnalysisInfo(res), nil
}

// Del 删除分析记录
func (l *AnalysisLogic) Del(ids []uint64) (any, error) {
	entity := analysis.NewAnalysisEntity(l.Ctx)
	err := entity.Del(ids...)
	return nil, err
}

// 获取最新版本号
func (l *AnalysisLogic) getLatestVersion(topicId, modelId uint64) int {
	entity := analysis.NewAnalysisEntity(l.Ctx)
	search := &analysis.SearchAnalysis{
		TopicId:  topicId,
		ModelId:  modelId,
		PageSize: 1000,
	}
	cond := entity.MakeConditon(*search)

	list, err := entity.List(cond)
	if err != nil {
		return 0
	}

	maxVersion := 0
	for _, item := range list {
		if item.Version > maxVersion {
			maxVersion = item.Version
		}
	}

	return maxVersion
}

// 取消指定课题和模型的所有记录的当前版本标记
func (l *AnalysisLogic) unsetCurrentVersion(topicId, modelId uint64) error {
	entity := analysis.NewAnalysisEntity(l.Ctx)
	search := &analysis.SearchAnalysis{
		TopicId:   topicId,
		ModelId:   modelId,
		IsCurrent: 1,
		PageSize:  1000,
	}
	cond := entity.MakeConditon(*search)

	list, err := entity.List(cond)
	if err != nil {
		return err
	}

	for _, item := range list {
		if concreteItem, ok := interface{}(item).(*analysis.AnalysisEntity); ok {
			concreteItem.IsCurrent = 0
			concreteItem.Update()
		}
	}

	return nil
}

// 辅助函数：转换为AnalysisInfo
func convertToAnalysisInfo(entity *analysis.AnalysisEntity) *analysis.AnalysisInfo {
	if entity == nil {
		return nil
	}

	return &analysis.AnalysisInfo{
		Id:            entity.Id,
		TopicId:       entity.TopicId,
		ModelId:       entity.ModelId,
		ModelName:     entity.ModelName,
		Content:       entity.Content,
		AiAnalysis:    entity.AiAnalysis,
		AiSuggestions: entity.AiSuggestions,
		Version:       entity.Version,
		IsCurrent:     entity.IsCurrent,
		IsCurrentBool: entity.IsCurrent == 1,
		UserId:        entity.UserId,
		CreatedAt:     entity.CreatedAt.String(),
		UpdatedAt:     entity.UpdatedAt.String(),
	}
}

// 辅助函数：转换为AnalysisDetail
func convertToAnalysisDetail(entity *analysis.AnalysisEntity) *analysis.AnalysisDetail {
	if entity == nil {
		return nil
	}

	info := convertToAnalysisInfo(entity)

	// 解析内容JSON
	var parsedContent map[string]interface{}
	if entity.Content != "" {
		json.Unmarshal([]byte(entity.Content), &parsedContent)
	}

	return &analysis.AnalysisDetail{
		AnalysisInfo:  *info,
		TopicTitle:    "", // 可以通过关联查询获取
		ParsedContent: parsedContent,
	}
}
