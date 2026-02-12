package analysis

// ==================== 请求DTO ====================

// CreateAnalysis 创建分析记录请求
type CreateAnalysis struct {
	TopicId      uint64 `json:"topicId" binding:"required"`       // 所属课题ID
	ModelId      uint64 `json:"modelId" binding:"required"`       // 使用的思维模型ID
	Content      string `json:"content" binding:"required"`       // 用户填写的分析内容JSON
	AiAnalysis   string `json:"aiAnalysis"`                       // AI分析结果（可选，创建时可能为空）
	AiSuggestions string `json:"aiSuggestions"`                   // AI建议（可选）
}

// UpdateAnalysis 更新分析记录请求
type UpdateAnalysis struct {
	Id            uint64 `json:"id" binding:"required"`
	Content       string `json:"content" binding:"required"`
	AiAnalysis    string `json:"aiAnalysis"`
	AiSuggestions string `json:"aiSuggestions"`
}

// SaveAnalysisWithAi 保存分析记录并更新AI结果请求
type SaveAnalysisWithAi struct {
	Id            uint64 `json:"id" binding:"required"`
	Content       string `json:"content" binding:"required"`
	AiAnalysis    string `json:"aiAnalysis" binding:"required"`
	AiSuggestions string `json:"aiSuggestions" binding:"required"`
}

// SetCurrentVersion 设置当前版本请求
type SetCurrentVersion struct {
	Id      uint64 `json:"id" binding:"required"`
	TopicId uint64 `json:"topicId" binding:"required"`
}

// SearchAnalysis 分析记录搜索条件
type SearchAnalysis struct {
	Page     int64  `json:"page" form:"page" search:"page"`
	PageSize int64  `json:"pageSize" form:"pageSize" search:"pageSize"`
	Id       uint64 `json:"id" form:"id" search:"type:eq;column:id;table:topic_analyses"`
	TopicId  uint64 `json:"topicId" form:"topicId" search:"type:eq;column:topic_id;table:topic_analyses"`
	ModelId  uint64 `json:"modelId" form:"modelId" search:"type:eq;column:model_id;table:topic_analyses"`
	UserId   uint64 `json:"userId" form:"userId" search:"type:eq;column:user_id;table:topic_analyses"`
	IsCurrent int   `json:"isCurrent" form:"isCurrent" search:"type:eq;column:is_current;table:topic_analyses"`
}

// DelAnalysis 删除分析记录请求
type DelAnalysis struct {
	Ids []uint64 `json:"ids" binding:"required,min=1"`
}

// GetCurrentAnalysisRequest 获取当前分析记录请求
type GetCurrentAnalysisRequest struct {
	TopicId uint64 `json:"topicId" form:"topicId" binding:"required"`
}

// GetLatestAnalysisRequest 获取课题最新分析记录请求
type GetLatestAnalysisRequest struct {
	TopicId uint64 `json:"topicId" form:"topicId" binding:"required"`
}

// ==================== 响应DTO ====================

// AnalysisInfo 分析记录信息DTO
type AnalysisInfo struct {
	Id            uint64 `json:"id"`
	TopicId       uint64 `json:"topicId"`
	ModelId       uint64 `json:"modelId"`
	ModelName     string `json:"modelName"`     // 模型名称（快照）
	Content       string `json:"content"`       // 用户填写的分析内容JSON
	AiAnalysis    string `json:"aiAnalysis"`    // AI分析结果
	AiSuggestions string `json:"aiSuggestions"` // AI建议
	Version       int    `json:"version"`
	IsCurrent     int    `json:"isCurrent"`
	IsCurrentBool bool   `json:"isCurrentBool"` // 布尔值格式
	UserId        uint64 `json:"userId"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

// AnalysisDetail 分析记录详情DTO
type AnalysisDetail struct {
	AnalysisInfo
	TopicTitle string `json:"topicTitle"` // 关联课题标题
	// 解析后的内容
	ParsedContent map[string]interface{} `json:"parsedContent"` // 解析后的Content JSON
}

// ListAnalysisResponse 分析记录列表响应
type ListAnalysisResponse struct {
	Page     int64           `json:"page"`
	PageSize int64           `json:"pageSize"`
	Total    int64           `json:"total"`
	List     []*AnalysisInfo `json:"list"`
}

// AnalysisHistory 分析记录历史（同一课题同一模型的所有版本）
type AnalysisHistory struct {
	TopicId   uint64          `json:"topicId"`
	ModelId   uint64          `json:"modelId"`
	ModelName string          `json:"modelName"`
	Versions  []*AnalysisInfo `json:"versions"`
}

// CreateAnalysisResponse 创建分析记录响应
type CreateAnalysisResponse struct {
	AnalysisInfo
	Version int `json:"version"` // 版本号
}

// ==================== 内容结构定义 ====================

// AnalysisContent 分析内容结构
type AnalysisContent struct {
	Fields map[string]interface{} `json:"fields"` // 字段值映射 fieldKey -> value
	Steps  []StepContent          `json:"steps"`  // 步骤内容
}

// StepContent 步骤内容
type StepContent struct {
	StepId    string                 `json:"stepId"`
	StepTitle string                 `json:"stepTitle"`
	Fields    map[string]interface{} `json:"fields"`
}

// AiResult AI分析结果结构
type AiResult struct {
	Summary     string   `json:"summary"`     // 总结
	KeyPoints   []string `json:"keyPoints"`   // 关键点
	Suggestions []string `json:"suggestions"` // 建议列表
	NextSteps   []string `json:"nextSteps"`   // 下一步行动
}
