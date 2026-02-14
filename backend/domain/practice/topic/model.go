package topic

import (
	"errors"

	"thinkingModels/component/db"

	"github.com/gin-gonic/gin"
	"github.com/jianyuezhexue/base"
)

// TopicEntityInterface 课题实体接口
type TopicEntityInterface interface {
	base.BaseModelInterface[TopicEntity]
	SelectModel(modelId uint64, modelName string) error
	RemoveModel() error
	MarkComplete() error
	Archive() error
	Reopen() error
	IncrementAnalysisCount()
	IncrementActionCount()
}

// TopicEntity 课题实体
type TopicEntity struct {
	base.BaseModel[TopicEntity]
	UserId        uint64       `json:"userId" type:"db" comment:"所属用户ID"`
	Title         string       `json:"title" type:"db" comment:"课题标题"`
	Description   string       `json:"description" type:"db" comment:"课题描述"`
	Background    string       `json:"background" type:"db" comment:"背景说明"`
	Goal          string       `json:"goal" type:"db" comment:"期望目标"`
	Constraints   string       `json:"constraints" type:"db" comment:"约束条件"`
	ModelId       uint64       `json:"modelId" type:"db" comment:"选用模型ID"`
	ModelName     string       `json:"modelName" type:"db" comment:"模型名称快照"`
	Status        int          `json:"status" type:"db" comment:"状态: 0=草稿, 1=进行中, 2=已完成, 3=已归档"`
	Priority      int          `json:"priority" type:"db" comment:"优先级: 1=低, 2=中, 3=高"`
	Deadline      db.LocalTime `json:"deadline" type:"db" comment:"截止日期"`
	AnalysisCount int          `json:"analysisCount" type:"db" comment:"分析次数"`
	ActionCount   int          `json:"actionCount" type:"db" comment:"行动数量"`
}

// NewTopicEntity 实例化课题实体
func NewTopicEntity(ctx *gin.Context, opt ...base.Option[TopicEntity]) TopicEntityInterface {
	entity := &TopicEntity{}
	entity.BaseModel = base.NewBaseModel(ctx, db.InitDb(), entity.TableName(), entity)
	if len(opt) > 0 {
		for _, fc := range opt {
			fc(&entity.BaseModel)
		}
	}
	return entity
}

// TableName 数据表名
func (t *TopicEntity) TableName() string {
	return "topics"
}

// Validate 数据校验
func (t *TopicEntity) Validate() error {
	if t.Title == "" {
		return errors.New("课题标题不能为空")
	}
	if len(t.Title) > 200 {
		return errors.New("课题标题不能超过200字符")
	}
	if t.UserId == 0 {
		return errors.New("用户ID不能为空")
	}
	return nil
}

// Repair 数据修复
func (t *TopicEntity) Repair() error {
	if t.Priority == 0 {
		t.Priority = 2
	}
	if t.Status < 0 || t.Status > 3 {
		t.Status = 0
	}
	return nil
}

// Complete 数据完善
func (t *TopicEntity) Complete() error {
	return nil
}

// SelectModel 选择模型
func (t *TopicEntity) SelectModel(modelId uint64, modelName string) error {
	t.ModelId = modelId
	t.ModelName = modelName
	if t.Status == 0 {
		t.Status = 1 // 自动进入进行中状态
	}
	return nil
}

// RemoveModel 移除模型
func (t *TopicEntity) RemoveModel() error {
	t.ModelId = 0
	t.ModelName = ""
	return nil
}

// MarkComplete 标记完成
func (t *TopicEntity) MarkComplete() error {
	if t.Status == 2 {
		return errors.New("课题已完成")
	}
	t.Status = 2
	return nil
}

// Archive 归档
func (t *TopicEntity) Archive() error {
	if t.Status == 3 {
		return errors.New("课题已归档")
	}
	t.Status = 3
	return nil
}

// Reopen 重新打开
func (t *TopicEntity) Reopen() error {
	if t.Status == 1 {
		return errors.New("课题已在进行中")
	}
	t.Status = 1
	return nil
}

// IncrementAnalysisCount 增加分析次数
func (t *TopicEntity) IncrementAnalysisCount() {
	t.AnalysisCount++
}

// IncrementActionCount 增加行动数量
func (t *TopicEntity) IncrementActionCount() {
	t.ActionCount++
}
