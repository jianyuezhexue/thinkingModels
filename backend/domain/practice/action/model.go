package action

import (
	"errors"
	"time"

	"thinkingModels/component/db"

	"github.com/gin-gonic/gin"
	"github.com/jianyuezhexue/base"
)

// ActionEntityInterface 行动实体接口
type ActionEntityInterface interface {
	base.BaseModelInterface[ActionEntity]
	UpdateProgress(progress int, note string) error
	MarkComplete() error
	Cancel() error
	CheckOverdue() bool
	IncrementFollowUpCount()
	DecrementFollowUpCount()
}

// ActionEntity 行动实体
type ActionEntity struct {
	base.BaseModel[ActionEntity]
	UserId         uint64       `json:"userId" type:"db" comment:"所属用户ID"`
	TopicId        uint64       `json:"topicId" type:"db" comment:"关联课题ID"`
	TopicTitle     string       `json:"topicTitle" type:"db" comment:"课题标题快照"`
	AnalysisId     uint64       `json:"analysisId" type:"db" comment:"来源分析ID"`
	Title          string       `json:"title" type:"db" comment:"行动标题"`
	Description    string       `json:"description" type:"db" comment:"行动描述"`
	Priority       int          `json:"priority" type:"db" comment:"优先级: 1=低, 2=中, 3=高"`
	Status         int          `json:"status" type:"db" comment:"状态: 0=待执行, 1=进行中, 2=已完成, 3=已取消"`
	Progress       int          `json:"progress" type:"db" comment:"完成进度(0-100)"`
	ExpectedResult string       `json:"expectedResult" type:"db" comment:"预期结果"`
	ActualResult   string       `json:"actualResult" type:"db" comment:"实际结果"`
	Deadline       db.LocalTime `json:"deadline" type:"db" comment:"截止日期"`
	CompletedAt    db.LocalTime `json:"completedAt" type:"db" comment:"完成时间"`
	FollowUpCount  int          `json:"followUpCount" type:"db" comment:"跟进记录数"`
}

// NewActionEntity 实例化行动实体
func NewActionEntity(ctx *gin.Context, opt ...base.Option[ActionEntity]) ActionEntityInterface {
	entity := &ActionEntity{}
	entity.BaseModel = base.NewBaseModel(ctx, db.InitDb(), entity.TableName(), entity)
	if len(opt) > 0 {
		for _, fc := range opt {
			fc(&entity.BaseModel)
		}
	}
	return entity
}

// TableName 数据表名
func (a *ActionEntity) TableName() string {
	return "actions"
}

// Validate 数据校验
func (a *ActionEntity) Validate() error {
	if a.Title == "" {
		return errors.New("行动标题不能为空")
	}
	if len(a.Title) > 200 {
		return errors.New("行动标题不能超过200字符")
	}
	if a.UserId == 0 {
		return errors.New("用户ID不能为空")
	}
	if a.Priority < 1 || a.Priority > 3 {
		a.Priority = 2
	}
	if a.Progress < 0 || a.Progress > 100 {
		return errors.New("进度必须在0-100之间")
	}
	return nil
}

// Repair 数据修复
func (a *ActionEntity) Repair() error {
	if a.Priority == 0 {
		a.Priority = 2
	}
	if a.Status < 0 || a.Status > 3 {
		a.Status = 0
	}
	return nil
}

// Complete 数据完善
func (a *ActionEntity) Complete() error {
	return nil
}

// UpdateProgress 更新进度
func (a *ActionEntity) UpdateProgress(progress int, note string) error {
	if progress < 0 || progress > 100 {
		return errors.New("进度必须在0-100之间")
	}
	a.Progress = progress
	// 自动更新状态
	if progress == 0 {
		a.Status = 0 // 待执行
	} else if progress == 100 {
		a.Status = 2 // 已完成
		a.CompletedAt = db.LocalTime(time.Now())
	} else {
		a.Status = 1 // 进行中
	}
	return nil
}

// MarkComplete 标记完成
func (a *ActionEntity) MarkComplete() error {
	a.Status = 2
	a.Progress = 100
	a.CompletedAt = db.LocalTime(time.Now())
	return nil
}

// Cancel 取消行动
func (a *ActionEntity) Cancel() error {
	a.Status = 3
	return nil
}

// CheckOverdue 检查是否逾期
func (a *ActionEntity) CheckOverdue() bool {
	if time.Time(a.Deadline).IsZero() {
		return false
	}
	return time.Now().After(time.Time(a.Deadline)) && a.Status != 2
}

// IncrementFollowUpCount 增加跟进数
func (a *ActionEntity) IncrementFollowUpCount() {
	a.FollowUpCount++
}

// DecrementFollowUpCount 减少跟进数
func (a *ActionEntity) DecrementFollowUpCount() {
	if a.FollowUpCount > 0 {
		a.FollowUpCount--
	}
}
