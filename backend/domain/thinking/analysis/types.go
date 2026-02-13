package analysis

// ==================== 请求DTO ====================

// CreateAnalysis 创建分析请求
type CreateAnalysis struct {
	TopicId   uint64 `json:"topicId" binding:"required"`
	ModelId   uint64 `json:"modelId" binding:"required"`
	ModelName string `json:"modelName"`
	Content   string `json:"content"`
}

// UpdateAnalysis 更新分析请求
type UpdateAnalysis struct {
	Id      uint64 `json:"id" binding:"required"`
	Content string `json:"content"`
}

// SaveWithAi 保存并AI分析请求
type SaveWithAi struct {
	Id        uint64 `json:"id"`
	TopicId   uint64 `json:"topicId" binding:"required"`
	ModelId   uint64 `json:"modelId" binding:"required"`
	ModelName string `json:"modelName"`
	Content   string `json:"content" binding:"required"`
}

// SetCurrentAnalysis 设为当前版本请求
type SetCurrentAnalysis struct {
	Id uint64 `json:"id" binding:"required"`
}

// SearchAnalysis 分析搜索条件
type SearchAnalysis struct {
	Page      int64  `json:"page" form:"page" search:"page"`
	PageSize  int64  `json:"pageSize" form:"pageSize" search:"pageSize"`
	Id        uint64 `json:"id" form:"id" search:"type:eq;column:id;table:topic_analyses"`
	TopicId   uint64 `json:"topicId" form:"topicId" search:"type:eq;column:topic_id;table:topic_analyses"`
	ModelId   uint64 `json:"modelId" form:"modelId" search:"type:eq;column:model_id;table:topic_analyses"`
	UserId    uint64 `json:"userId" form:"userId" search:"type:eq;column:user_id;table:topic_analyses"`
	IsCurrent *bool  `json:"isCurrent" form:"isCurrent" search:"type:eq;column:is_current;table:topic_analyses"`
	Status    int    `json:"status" form:"status" search:"type:eq;column:status;table:topic_analyses"`
}

// DelAnalysis 删除分析请求
type DelAnalysis struct {
	Ids []uint64 `json:"ids" binding:"required,min=1"`
}

// ==================== 响应DTO ====================

// AnalysisInfo 分析信息DTO
type AnalysisInfo struct {
	Id            uint64      `json:"id"`
	TopicId       uint64      `json:"topicId"`
	ModelId       uint64      `json:"modelId"`
	ModelName     string      `json:"modelName"`
	Content       string      `json:"content"`
	AiAnalysis    string      `json:"aiAnalysis"`
	AiSuggestions string      `json:"aiSuggestions"`
	Version       int         `json:"version"`
	IsCurrent     bool        `json:"isCurrent"`
	UserId        uint64      `json:"userId"`
	Status        int         `json:"status"`
	StatusText    string      `json:"statusText"`
	CreatedAt     interface{} `json:"createdAt"`
	UpdatedAt     interface{} `json:"updatedAt"`
}

// AnalysisDetail 分析详情DTO
type AnalysisDetail struct {
	Id           uint64      `json:"id"`
	TopicId      uint64      `json:"topicId"`
	ModelId      uint64      `json:"modelId"`
	ModelName    string      `json:"modelName"`
	InputContent string      `json:"inputContent"`
	AiResult     string      `json:"aiResult"`
	UserResult   string      `json:"userResult"`
	Conclusion   string      `json:"conclusion"`
	Version      int         `json:"version"`
	IsCurrent    bool        `json:"isCurrent"`
	UserId       uint64      `json:"userId"`
	Status       int         `json:"status"`
	StatusText   string      `json:"statusText"`
	CreatedAt    interface{} `json:"createdAt"`
	UpdatedAt    interface{} `json:"updatedAt"`
}

// ListAnalysisResponse 分析列表响应
type ListAnalysisResponse struct {
	Page     int64           `json:"page"`
	PageSize int64           `json:"pageSize"`
	Total    int64           `json:"total"`
	List     []*AnalysisInfo `json:"list"`
}

// AnalysisHistory 分析历史版本
type AnalysisHistory struct {
	TopicId   string          `json:"topicId"`
	ModelId   string          `json:"modelId"`
	ModelName string          `json:"modelName"`
	Versions  []*AnalysisInfo `json:"versions"`
}
