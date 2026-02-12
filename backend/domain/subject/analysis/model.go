package analysis

import (
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/jianyuezhexue/base"
	"thinkingModels/component/db"
)

// AnalysisEntityInterface 分析记录实体接口
type AnalysisEntityInterface interface {
	base.BaseModelInterface[AnalysisEntity]
	// 版本管理
	SetAsCurrent() error
	UnsetAsCurrent() error
	// 内容管理
	UpdateContent(content string) error
	UpdateAiResult(analysis, suggestions string) error
	// 数据解析
	ParseContent() (map[string]interface{}, error)
	GetFieldValue(fieldKey string) (interface{}, bool)
}

// AnalysisEntity 分析记录实体
type AnalysisEntity struct {
	base.BaseModel[AnalysisEntity]
	TopicId       uint64 `json:"topicId" type:"db" comment:"所属课题ID"`
	ModelId       uint64 `json:"modelId" type:"db" comment:"使用的思维模型ID"`
	ModelName     string `json:"modelName" type:"db" comment:"模型名称（快照）"`
	Content       string `json:"content" type:"db" comment:"用户填写的分析内容JSON"`
	AiAnalysis    string `json:"aiAnalysis" type:"db" comment:"AI分析结果"`
	AiSuggestions string `json:"aiSuggestions" type:"db" comment:"AI建议"`
	Version       int    `json:"version" type:"db" comment:"版本号"`
	IsCurrent     int    `json:"isCurrent" type:"db" comment:"是否为当前版本: 0=否, 1=是"`
	UserId        uint64 `json:"userId" type:"db" comment:"用户ID"`
}

// NewAnalysisEntity 实例化分析记录实体
func NewAnalysisEntity(ctx *gin.Context, opt ...base.Option[AnalysisEntity]) AnalysisEntityInterface {
	entity := &AnalysisEntity{}
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
func (a *AnalysisEntity) TableName() string {
	return "topic_analyses"
}

// Validate 数据校验
func (a *AnalysisEntity) Validate() error {
	if a.TopicId == 0 {
		return errors.New("课题ID不能为空")
	}
	if a.ModelId == 0 {
		return errors.New("思维模型ID不能为空")
	}
	if a.Content == "" {
		return errors.New("分析内容不能为空")
	}
	// 校验Content是否为有效的JSON
	var contentMap map[string]interface{}
	if err := json.Unmarshal([]byte(a.Content), &contentMap); err != nil {
		return errors.New("分析内容必须是有效的JSON格式")
	}
	return nil
}

// Repair 数据修复
func (a *AnalysisEntity) Repair() error {
	if a.Version == 0 {
		a.Version = 1
	}
	if a.IsCurrent != 0 && a.IsCurrent != 1 {
		a.IsCurrent = 1
	}
	return nil
}

// Complete 数据完善
func (a *AnalysisEntity) Complete() error {
	return nil
}

// SetAsCurrent 设置为当前版本
func (a *AnalysisEntity) SetAsCurrent() error {
	a.IsCurrent = 1
	return nil
}

// UnsetAsCurrent 取消当前版本设置
func (a *AnalysisEntity) UnsetAsCurrent() error {
	a.IsCurrent = 0
	return nil
}

// UpdateContent 更新分析内容
func (a *AnalysisEntity) UpdateContent(content string) error {
	// 验证内容是否为有效的JSON
	var contentMap map[string]interface{}
	if err := json.Unmarshal([]byte(content), &contentMap); err != nil {
		return err
	}
	a.Content = content
	return nil
}

// UpdateAiResult 更新AI分析结果
func (a *AnalysisEntity) UpdateAiResult(analysis, suggestions string) error {
	a.AiAnalysis = analysis
	a.AiSuggestions = suggestions
	return nil
}

// ParseContent 解析内容JSON
func (a *AnalysisEntity) ParseContent() (map[string]interface{}, error) {
	var result map[string]interface{}
	if a.Content == "" {
		return result, nil
	}
	err := json.Unmarshal([]byte(a.Content), &result)
	return result, err
}

// GetFieldValue 获取指定字段的值
func (a *AnalysisEntity) GetFieldValue(fieldKey string) (interface{}, bool) {
	content, err := a.ParseContent()
	if err != nil {
		return nil, false
	}
	val, ok := content[fieldKey]
	return val, ok
}

// GetContentField 获取内容中的字段值（支持嵌套路径，如 "step1.field1"）
func (a *AnalysisEntity) GetContentField(path string) (interface{}, bool) {
	content, err := a.ParseContent()
	if err != nil {
		return nil, false
	}

	// 简单路径直接查找
	val, ok := content[path]
	if ok {
		return val, true
	}

	// TODO: 支持嵌套路径解析

	return nil, false
}
