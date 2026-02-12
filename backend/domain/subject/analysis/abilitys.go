package analysis

// AnalysisAbility 分析记录能力接口
type AnalysisAbility interface {
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
