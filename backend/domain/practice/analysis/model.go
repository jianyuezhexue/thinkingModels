package analysis

import (
	"errors"

	"thinkingModels/component/db"

	"github.com/gin-gonic/gin"
	"github.com/jianyuezhexue/base"
)

// AnalysisEntityInterface 分析实体接口
type AnalysisEntityInterface interface {
	base.BaseModelInterface[AnalysisEntity]
	SetAsCurrent() error
	IncrementVersion()
	SetAiResult(aiResult, userResult string)
	MarkComplete() error
	MarkFailed() error
}

// AnalysisEntity 分析实体
type AnalysisEntity struct {
	base.BaseModel[AnalysisEntity]
	UserId       uint64 `json:"userId" type:"db" comment:"所属用户ID"`
	TopicId      uint64 `json:"topicId" type:"db" comment:"关联课题ID"`
	TopicTitle   string `json:"topicTitle" type:"db" comment:"课题标题快照"`
	ModelId      uint64 `json:"modelId" type:"db" comment:"使用模型ID"`
	ModelName    string `json:"modelName" type:"db" comment:"模型名称快照"`
	InputContent string `json:"inputContent" type:"db" comment:"用户输入内容(JSON)"`
	AiResult     string `json:"aiResult" type:"db" comment:"AI分析结果(JSON)"`
	UserResult   string `json:"userResult" type:"db" comment:"用户编辑结果(JSON)"`
	Conclusion   string `json:"conclusion" type:"db" comment:"分析结论"`
	Version      int    `json:"version" type:"db" comment:"版本号"`
	IsCurrent    bool   `json:"isCurrent" type:"db" comment:"是否当前版本"`
	Status       int    `json:"status" type:"db" comment:"状态: 0=草稿, 1=分析中, 2=已完成, 3=失败"`
}

// NewAnalysisEntity 实例化分析实体
func NewAnalysisEntity(ctx *gin.Context, opt ...base.Option[AnalysisEntity]) AnalysisEntityInterface {
	entity := &AnalysisEntity{}
	entity.BaseModel = base.NewBaseModel(ctx, db.InitDb(), entity.TableName(), entity)
	if len(opt) > 0 {
		for _, fc := range opt {
			fc(&entity.BaseModel)
		}
	}
	return entity
}

// TableName 数据表名
func (a *AnalysisEntity) TableName() string {
	return "topic_analyses"
}

// Validate 数据校验
func (a *AnalysisEntity) Validate() error {
	if a.TopicId == 0 {
		return errors.New("课题ID不能为空")
	}
	if a.UserId == 0 {
		return errors.New("用户ID不能为空")
	}
	return nil
}

// Repair 数据修复
func (a *AnalysisEntity) Repair() error {
	if a.Version == 0 {
		a.Version = 1
	}
	if a.Status < 0 || a.Status > 3 {
		a.Status = 0
	}
	return nil
}

// Complete 数据完善
func (a *AnalysisEntity) Complete() error {
	return nil
}

// SetAsCurrent 设为当前版本
func (a *AnalysisEntity) SetAsCurrent() error {
	// 先将同课题同模型的其他分析设为非当前
	// TODO: 需要批量更新其他记录
	a.IsCurrent = true
	return nil
}

// IncrementVersion 增加版本号
func (a *AnalysisEntity) IncrementVersion() {
	a.Version++
}

// SetAiResult 设置AI分析结果
func (a *AnalysisEntity) SetAiResult(aiResult, userResult string) {
	a.AiResult = aiResult
	a.UserResult = userResult
}

// MarkComplete 标记完成
func (a *AnalysisEntity) MarkComplete() error {
	a.Status = 2
	return nil
}

// MarkFailed 标记失败
func (a *AnalysisEntity) MarkFailed() error {
	a.Status = 3
	return nil
}
