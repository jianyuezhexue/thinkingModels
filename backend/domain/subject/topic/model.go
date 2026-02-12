package topic

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jianyuezhexue/base"
	"thinkingModels/component/db"
)

// TopicEntityInterface 课题实体接口
type TopicEntityInterface interface {
	base.BaseModelInterface[TopicEntity]
	// 状态管理
	Activate() error
	Complete() error
	Archive() error
	Reopen() error
	// 模型关联
	SelectModel(modelId uint64) error
	RemoveModel() error
	// 标签管理
	SetTags(tags []string)
	GetTagsList() []string
}

// TopicEntity 课题实体
type TopicEntity struct {
	base.BaseModel[TopicEntity]
	Title        string       `json:"title" type:"db" comment:"课题标题"`
	Description  string       `json:"description" type:"db" comment:"课题描述"`
	Status       int          `json:"status" type:"db" comment:"状态: 0=进行中, 1=已完成, 2=已归档"`
	UserId       uint64       `json:"userId" type:"db" comment:"所属用户ID"`
	ModelId      uint64       `json:"modelId" type:"db" comment:"当前使用的思维模型ID"`
	Priority     int          `json:"priority" type:"db" comment:"优先级: 1=低, 2=中, 3=高"`
	Tags         string       `json:"tags" type:"db" comment:"标签，逗号分隔"`
	Deadline     db.LocalTime `json:"deadline" type:"db" comment:"截止日期"`
	CompleteTime db.LocalTime `json:"completeTime" type:"db" comment:"完成时间"`
}

// NewTopicEntity 实例化课题实体
func NewTopicEntity(ctx *gin.Context, opt ...base.Option[TopicEntity]) TopicEntityInterface {
	entity := &TopicEntity{}
	entity.BaseModel = base.NewBaseModel(ctx, db.InitDb(), entity.TableName(), entity)

	// 自定义配置选项
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
	if len(t.Description) > 2000 {
		return errors.New("课题描述不能超过2000字符")
	}
	// 优先级校验
	if t.Priority < 1 || t.Priority > 3 {
		t.Priority = 1 // 默认为低优先级
	}
	// 标签长度校验
	if len(t.Tags) > 500 {
		return errors.New("标签总长度不能超过500字符")
	}
	return nil
}

// Repair 数据修复
func (t *TopicEntity) Repair() error {
	if t.Status < 0 || t.Status > 2 {
		t.Status = 0 // 默认为进行中
	}
	if t.Priority == 0 {
		t.Priority = 1
	}
	return nil
}

// Complete 数据完善
func (t *TopicEntity) Complete() error {
	return nil
}

// Activate 激活课题（进行中）
func (t *TopicEntity) Activate() error {
	t.Status = 0
	return nil
}

// Finish 完成课题
func (t *TopicEntity) Finish() error {
	t.Status = 1
	t.CompleteTime = db.LocalTime(time.Now())
	return nil
}

// Archive 归档课题
func (t *TopicEntity) Archive() error {
	t.Status = 2
	return nil
}

// Reopen 重新打开课题
func (t *TopicEntity) Reopen() error {
	t.Status = 0
	t.CompleteTime = db.LocalTime{} // 清空完成时间
	return nil
}

// SelectModel 选择思维模型
func (t *TopicEntity) SelectModel(modelId uint64) error {
	t.ModelId = modelId
	return nil
}

// RemoveModel 移除思维模型
func (t *TopicEntity) RemoveModel() error {
	t.ModelId = 0
	return nil
}

// SetTags 设置标签（传入切片，内部处理为逗号分隔）
func (t *TopicEntity) SetTags(tags []string) {
	if len(tags) == 0 {
		t.Tags = ""
		return
	}
	// 过滤空标签并去重
	seen := make(map[string]bool)
	filtered := make([]string, 0, len(tags))
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag != "" && !seen[tag] {
			seen[tag] = true
			filtered = append(filtered, tag)
		}
	}
	t.Tags = strings.Join(filtered, ",")
}

// GetTagsList 获取标签列表（将逗号分隔的标签转换为切片）
func (t *TopicEntity) GetTagsList() []string {
	if t.Tags == "" {
		return []string{}
	}
	parts := strings.Split(t.Tags, ",")
	result := make([]string, 0, len(parts))
	for _, tag := range parts {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			result = append(result, tag)
		}
	}
	return result
}

// AddTag 添加单个标签
func (t *TopicEntity) AddTag(tag string) error {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return nil
	}

	tags := t.GetTagsList()
	// 检查是否已存在
	for _, existing := range tags {
		if existing == tag {
			return nil // 已存在，不重复添加
		}
	}

	tags = append(tags, tag)
	t.SetTags(tags)
	return nil
}

// RemoveTag 移除单个标签
func (t *TopicEntity) RemoveTag(tag string) error {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return nil
	}

	tags := t.GetTagsList()
	newTags := make([]string, 0, len(tags))
	for _, existing := range tags {
		if existing != tag {
			newTags = append(newTags, existing)
		}
	}
	t.SetTags(newTags)
	return nil
}

// IsOverdue 检查是否逾期
func (t *TopicEntity) IsOverdue() bool {
	if t.Deadline.IsZero() || t.Status == 1 || t.Status == 2 {
		return false // 无截止日期或已完成/已归档，不算逾期
	}
	return time.Time(t.Deadline).Before(time.Now())
}

// GetStatusText 获取状态文本
func (t *TopicEntity) GetStatusText() string {
	switch t.Status {
	case 0:
		return "进行中"
	case 1:
		return "已完成"
	case 2:
		return "已归档"
	default:
		return "未知"
	}
}

// GetPriorityText 获取优先级文本
func (t *TopicEntity) GetPriorityText() string {
	switch t.Priority {
	case 1:
		return "低"
	case 2:
		return "中"
	case 3:
		return "高"
	default:
		return "未知"
	}
}
