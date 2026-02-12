package subject

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"thinkingModels/component/db"
	"thinkingModels/domain/market/model"
	"thinkingModels/domain/subject/analysis"
	"thinkingModels/domain/subject/topic"
	"thinkingModels/logic"
)

// TopicLogic 课题业务逻辑
type TopicLogic struct {
	logic.BaseLogic
}

// NewTopicLogic 初始化TopicLogic
func NewTopicLogic(ctx *gin.Context) *TopicLogic {
	return &TopicLogic{BaseLogic: logic.BaseLogic{Ctx: ctx}}
}

// Create 创建课题
func (l *TopicLogic) Create(req *topic.CreateTopic) (*topic.TopicInfo, error) {
	// 实例化模型
	entity := topic.NewTopicEntity(l.Ctx)

	// 设置数据
	if concreteEntity, ok := entity.(*topic.TopicEntity); ok {
		concreteEntity.Title = req.Title
		concreteEntity.Description = req.Description
		concreteEntity.Priority = req.Priority
		if concreteEntity.Priority == 0 {
			concreteEntity.Priority = 1
		}
		concreteEntity.Tags = req.Tags

		// 解析截止日期
		if req.Deadline != "" {
			deadline, err := time.Parse(time.RFC3339, req.Deadline)
			if err == nil {
				concreteEntity.Deadline = db.LocalTime(deadline)
			}
		}

		// 设置默认状态
		concreteEntity.Status = 0 // 进行中

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

	// 保存数据
	res, err := entity.Create()
	if err != nil {
		return nil, err
	}

	// 转换为DTO返回
	return convertToTopicInfo(res), nil
}

// Update 更新课题
func (l *TopicLogic) Update(req *topic.UpdateTopic) (*topic.TopicInfo, error) {
	// 实例化模型
	entity := topic.NewTopicEntity(l.Ctx)

	// 先加载旧数据
	oldData, err := entity.LoadById(req.Id)
	if err != nil {
		return nil, err
	}

	// 设置数据
	if concreteEntity, ok := entity.(*topic.TopicEntity); ok {
		concreteEntity.Title = req.Title
		concreteEntity.Description = req.Description
		concreteEntity.Priority = req.Priority
		concreteEntity.Tags = req.Tags

		// 解析截止日期
		if req.Deadline != "" {
			deadline, err := time.Parse(time.RFC3339, req.Deadline)
			if err == nil {
				concreteEntity.Deadline = db.LocalTime(deadline)
			}
		} else {
			concreteEntity.Deadline = db.LocalTime{} // 清空截止日期
		}

		// 保持原有状态和用户
		concreteEntity.Status = oldData.Status
		concreteEntity.UserId = oldData.UserId
		concreteEntity.ModelId = oldData.ModelId
		concreteEntity.CompleteTime = oldData.CompleteTime
	}

	// 数据校验
	if err := entity.Validate(); err != nil {
		return nil, err
	}

	// 更新数据
	res, err := entity.Update()
	if err != nil {
		return nil, err
	}

	// 转换为DTO返回
	return convertToTopicInfo(res), nil
}

// Get 查询课题详情（包含分析记录）
func (l *TopicLogic) Get(id uint64) (*topic.TopicDetail, error) {
	// 获取课题信息
	entity := topic.NewTopicEntity(l.Ctx)
	topicData, err := entity.LoadById(id)
	if err != nil {
		return nil, err
	}

	// 获取课题的分析记录列表
	analysisLogic := NewAnalysisLogic(l.Ctx)
	analyses, err := analysisLogic.GetByTopic(id)
	if err != nil {
		analyses = make([]*analysis.AnalysisInfo, 0) // 如果出错，返回空列表
	}

	// 转换为简要信息
	analysisBriefs := make([]*topic.TopicAnalysisBrief, 0, len(analyses))
	for _, a := range analyses {
		analysisBriefs = append(analysisBriefs, &topic.TopicAnalysisBrief{
			Id:        a.Id,
			ModelId:   a.ModelId,
			ModelName: a.ModelName,
			Version:   a.Version,
			IsCurrent: a.IsCurrent == 1,
			CreatedAt: a.CreatedAt,
		})
	}

	topicInfo := convertToTopicInfo(topicData)

	return &topic.TopicDetail{
		TopicInfo: *topicInfo,
		Analyses:  analysisBriefs,
	}, nil
}

// GetSimple 查询课题简单信息
func (l *TopicLogic) GetSimple(id uint64) (*topic.TopicInfo, error) {
	entity := topic.NewTopicEntity(l.Ctx)
	topicData, err := entity.LoadById(id)
	if err != nil {
		return nil, err
	}
	return convertToTopicInfo(topicData), nil
}

// List 查询课题列表
func (l *TopicLogic) List(req *topic.SearchTopic) (*topic.ListTopicResponse, error) {
	entity := topic.NewTopicEntity(l.Ctx)

	cond := entity.MakeConditon(*req)
	total, err := entity.Count(cond)
	if err != nil {
		return nil, err
	}

	list, err := entity.List(cond)
	if err != nil {
		return nil, err
	}

	// 转换为DTO列表
	topicList := make([]*topic.TopicInfo, 0, len(list))
	for _, item := range list {
		topicList = append(topicList, convertToTopicInfo(item))
	}

	return &topic.ListTopicResponse{
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
		List:     topicList,
	}, nil
}

// Del 删除课题
func (l *TopicLogic) Del(ids []uint64) (any, error) {
	entity := topic.NewTopicEntity(l.Ctx)
	err := entity.Del(ids...)
	return nil, err
}

// UpdateStatus 更新课题状态
func (l *TopicLogic) UpdateStatus(req *topic.UpdateTopicStatus) (*topic.TopicInfo, error) {
	entity := topic.NewTopicEntity(l.Ctx)

	_, err := entity.LoadById(req.Id)
	if err != nil {
		return nil, err
	}

	if concreteEntity, ok := entity.(*topic.TopicEntity); ok {
		concreteEntity.Status = req.Status

		// 如果标记为已完成，设置完成时间
		if req.Status == 1 {
			concreteEntity.CompleteTime = db.LocalTime(time.Now())
		} else if req.Status == 0 {
			// 如果重新进行中，清空完成时间
			concreteEntity.CompleteTime = db.LocalTime{}
		}
	}

	res, err := entity.Update()
	if err != nil {
		return nil, err
	}

	return convertToTopicInfo(res), nil
}

// SelectModel 为课题选择思维模型
func (l *TopicLogic) SelectModel(req *topic.UpdateTopicModel) (*topic.TopicInfo, error) {
	// 验证模型是否存在
	modelEntity := model.NewModelEntity(l.Ctx)
	modelData, err := modelEntity.LoadById(req.ModelId)
	if err != nil {
		return nil, errors.New("思维模型不存在")
	}

	// 检查模型状态
	if modelData.Status != 1 {
		return nil, errors.New("该思维模型未发布，无法选用")
	}

	// 更新课题的模型
	entity := topic.NewTopicEntity(l.Ctx)
	_, err = entity.LoadById(req.Id)
	if err != nil {
		return nil, err
	}

	if concreteEntity, ok := entity.(*topic.TopicEntity); ok {
		concreteEntity.ModelId = req.ModelId
	}

	res, err := entity.Update()
	if err != nil {
		return nil, err
	}

	return convertToTopicInfo(res), nil
}

// RemoveModel 移除课题的思维模型
func (l *TopicLogic) RemoveModel(id uint64) (*topic.TopicInfo, error) {
	entity := topic.NewTopicEntity(l.Ctx)
	_, err := entity.LoadById(id)
	if err != nil {
		return nil, err
	}

	if concreteEntity, ok := entity.(*topic.TopicEntity); ok {
		concreteEntity.ModelId = 0
	}

	res, err := entity.Update()
	if err != nil {
		return nil, err
	}

	return convertToTopicInfo(res), nil
}

// Complete 完成课题
func (l *TopicLogic) Complete(req *topic.CompleteTopic) (*topic.TopicInfo, error) {
	entity := topic.NewTopicEntity(l.Ctx)
	_, err := entity.LoadById(req.Id)
	if err != nil {
		return nil, err
	}

	if concreteEntity, ok := entity.(*topic.TopicEntity); ok {
		concreteEntity.Status = 1 // 已完成
		concreteEntity.CompleteTime = db.LocalTime(time.Now())
	}

	res, err := entity.Update()
	if err != nil {
		return nil, err
	}

	return convertToTopicInfo(res), nil
}

// Archive 归档课题
func (l *TopicLogic) Archive(req *topic.ArchiveTopic) (*topic.TopicInfo, error) {
	entity := topic.NewTopicEntity(l.Ctx)
	_, err := entity.LoadById(req.Id)
	if err != nil {
		return nil, err
	}

	if concreteEntity, ok := entity.(*topic.TopicEntity); ok {
		concreteEntity.Status = 2 // 已归档
	}

	res, err := entity.Update()
	if err != nil {
		return nil, err
	}

	return convertToTopicInfo(res), nil
}

// Reopen 重新打开课题
func (l *TopicLogic) Reopen(id uint64) (*topic.TopicInfo, error) {
	entity := topic.NewTopicEntity(l.Ctx)
	_, err := entity.LoadById(id)
	if err != nil {
		return nil, err
	}

	if concreteEntity, ok := entity.(*topic.TopicEntity); ok {
		concreteEntity.Status = 0 // 进行中
		concreteEntity.CompleteTime = db.LocalTime{} // 清空完成时间
	}

	res, err := entity.Update()
	if err != nil {
		return nil, err
	}

	return convertToTopicInfo(res), nil
}

// GetStatistics 获取课题统计
func (l *TopicLogic) GetStatistics(userId uint64) (*topic.TopicStatistics, error) {
	entity := topic.NewTopicEntity(l.Ctx)

	// 统计总数
	search := &topic.SearchTopic{UserId: userId}
	cond := entity.MakeConditon(*search)
	total, err := entity.Count(cond)
	if err != nil {
		total = 0
	}

	// 统计各状态数量
	search.Status = 0
	cond = entity.MakeConditon(*search)
	active, _ := entity.Count(cond)

	search.Status = 1
	cond = entity.MakeConditon(*search)
	completed, _ := entity.Count(cond)

	search.Status = 2
	cond = entity.MakeConditon(*search)
	archived, _ := entity.Count(cond)

	// 统计高优先级
	search.Status = -1
	search.Priority = 3
	cond = entity.MakeConditon(*search)
	highPriority, _ := entity.Count(cond)

	return &topic.TopicStatistics{
		Total:        total,
		Active:       active,
		Completed:    completed,
		Archived:     archived,
		HighPriority: highPriority,
	}, nil
}

// 辅助函数：转换为TopicInfo
func convertToTopicInfo(entity *topic.TopicEntity) *topic.TopicInfo {
	if entity == nil {
		return nil
	}

	// 获取模型名称（如果需要可以在这里关联查询）
	modelName := ""
	if entity.ModelId > 0 {
		// 这里简化处理，实际可以通过关联查询获取
		modelName = ""
	}

	// 分割标签
	tagsList := make([]string, 0)
	if entity.Tags != "" {
		tagsList = strings.Split(entity.Tags, ",")
		// 清理空白
		cleaned := make([]string, 0, len(tagsList))
		for _, tag := range tagsList {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				cleaned = append(cleaned, tag)
			}
		}
		tagsList = cleaned
	}

	return &topic.TopicInfo{
		Id:           entity.Id,
		Title:        entity.Title,
		Description:  entity.Description,
		Status:       entity.Status,
		StatusText:   entity.GetStatusText(),
		Priority:     entity.Priority,
		PriorityText: entity.GetPriorityText(),
		Tags:         entity.Tags,
		TagsList:     tagsList,
		Deadline:     entity.Deadline.String(),
		CompleteTime: entity.CompleteTime.String(),
		UserId:       entity.UserId,
		ModelId:      entity.ModelId,
		ModelName:    modelName,
		CreatedAt:    entity.CreatedAt.String(),
		UpdatedAt:    entity.UpdatedAt.String(),
	}
}
