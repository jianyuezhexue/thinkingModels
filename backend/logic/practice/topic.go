package thinking

import (
	"thinkingModels/domain/practice/topic"
	"thinkingModels/logic"

	"github.com/gin-gonic/gin"
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
	entity := topic.NewTopicEntity(l.Ctx)
	_, err := entity.SetData(req)
	if err != nil {
		return nil, err
	}

	if concreteEntity, ok := entity.(*topic.TopicEntity); ok {
		concreteEntity.Status = 0
	}

	if err := entity.Validate(); err != nil {
		return nil, err
	}

	res, err := entity.Create()
	if err != nil {
		return nil, err
	}

	return convertToTopicInfo(res), nil
}

// Update 更新课题
func (l *TopicLogic) Update(req *topic.UpdateTopic) (*topic.TopicInfo, error) {
	entity := topic.NewTopicEntity(l.Ctx)
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

	return convertToTopicInfo(res), nil
}

// Get 查询课题详情
func (l *TopicLogic) Get(id uint64) (*topic.TopicDetail, error) {
	entity := topic.NewTopicEntity(l.Ctx)
	res, err := entity.LoadById(id)
	if err != nil {
		return nil, err
	}
	return convertToTopicDetail(res), nil
}

// List 查询课题列表
func (l *TopicLogic) List(req *topic.SearchTopic) (*logic.ListReap, error) {
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

	topicInfoList := make([]*topic.TopicInfo, 0, len(list))
	for _, item := range list {
		topicInfoList = append(topicInfoList, convertToTopicInfo(item))
	}

	return &logic.ListReap{
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
		List:     topicInfoList,
	}, nil
}

// ListMy 查询我的课题
func (l *TopicLogic) ListMy(req *topic.SearchTopic) (*logic.ListReap, error) {
	return l.List(req)
}

// Del 删除课题
func (l *TopicLogic) Del(req *topic.DelTopic) error {
	entity := topic.NewTopicEntity(l.Ctx)
	return entity.Del(req.Ids...)
}

// SelectModel 选择模型
func (l *TopicLogic) SelectModel(req *topic.SelectModelForTopic) error {
	entity := topic.NewTopicEntity(l.Ctx)
	_, err := entity.LoadById(req.TopicId)
	if err != nil {
		return err
	}
	if err := entity.SelectModel(req.ModelId, req.ModelName); err != nil {
		return err
	}
	_, err = entity.Update()
	return err
}

// RemoveModel 移除模型
func (l *TopicLogic) RemoveModel(topicId uint64) error {
	entity := topic.NewTopicEntity(l.Ctx)
	_, err := entity.LoadById(topicId)
	if err != nil {
		return err
	}
	if err := entity.RemoveModel(); err != nil {
		return err
	}
	_, err = entity.Update()
	return err
}

// UpdateStatus 更新状态
func (l *TopicLogic) UpdateStatus(req *topic.UpdateTopicStatus) error {
	entity := topic.NewTopicEntity(l.Ctx)
	_, err := entity.LoadById(req.Id)
	if err != nil {
		return err
	}

	if concreteEntity, ok := entity.(*topic.TopicEntity); ok {
		concreteEntity.Status = req.Status
	}

	_, err = entity.Update()
	return err
}

// Complete 完成课题
func (l *TopicLogic) Complete(id uint64) error {
	entity := topic.NewTopicEntity(l.Ctx)
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

// Archive 归档课题
func (l *TopicLogic) Archive(id uint64) error {
	entity := topic.NewTopicEntity(l.Ctx)
	_, err := entity.LoadById(id)
	if err != nil {
		return err
	}

	if err := entity.Archive(); err != nil {
		return err
	}

	_, err = entity.Update()
	return err
}

// Reopen 重新打开课题
func (l *TopicLogic) Reopen(id uint64) error {
	entity := topic.NewTopicEntity(l.Ctx)
	_, err := entity.LoadById(id)
	if err != nil {
		return err
	}

	if err := entity.Reopen(); err != nil {
		return err
	}

	_, err = entity.Update()
	return err
}

// Statistics 获取统计信息
func (l *TopicLogic) Statistics() (*topic.TopicStatistics, error) {
	entity := topic.NewTopicEntity(l.Ctx)

	draftCond := entity.MakeConditon(topic.SearchTopic{Status: intPtr(0)})
	draftCount, _ := entity.Count(draftCond)

	progressCond := entity.MakeConditon(topic.SearchTopic{Status: intPtr(1)})
	progressCount, _ := entity.Count(progressCond)

	completeCond := entity.MakeConditon(topic.SearchTopic{Status: intPtr(2)})
	completeCount, _ := entity.Count(completeCond)

	archiveCond := entity.MakeConditon(topic.SearchTopic{Status: intPtr(3)})
	archiveCount, _ := entity.Count(archiveCond)

	return &topic.TopicStatistics{
		DraftCount:    draftCount,
		ProgressCount: progressCount,
		CompleteCount: completeCount,
		ArchiveCount:  archiveCount,
	}, nil
}

func intPtr(i int) *int {
	return &i
}

func convertToTopicInfo(entity any) *topic.TopicInfo {
	e, ok := entity.(*topic.TopicEntity)
	if !ok {
		return nil
	}
	return &topic.TopicInfo{
		Id:            e.Id,
		Title:         e.Title,
		Description:   e.Description,
		Background:    e.Background,
		Goal:          e.Goal,
		Constraints:   e.Constraints,
		ModelId:       e.ModelId,
		ModelName:     e.ModelName,
		Status:        e.Status,
		Priority:      e.Priority,
		Deadline:      e.Deadline,
		AnalysisCount: e.AnalysisCount,
		ActionCount:   e.ActionCount,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
}

func convertToTopicDetail(entity any) *topic.TopicDetail {
	e, ok := entity.(*topic.TopicEntity)
	if !ok {
		return nil
	}
	return &topic.TopicDetail{
		Id:            e.Id,
		Title:         e.Title,
		Description:   e.Description,
		Background:    e.Background,
		Goal:          e.Goal,
		Constraints:   e.Constraints,
		ModelId:       e.ModelId,
		ModelName:     e.ModelName,
		Status:        e.Status,
		Priority:      e.Priority,
		Deadline:      e.Deadline,
		AnalysisCount: e.AnalysisCount,
		ActionCount:   e.ActionCount,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
}
