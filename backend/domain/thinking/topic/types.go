package topic

// ==================== 请求DTO ====================

// CreateTopic 创建课题请求
type CreateTopic struct {
	Title       string `json:"title" binding:"required,max=200"`
	Description string `json:"description"`
	Background  string `json:"background"`
	Goal        string `json:"goal"`
	Constraints string `json:"constraints"`
	Priority    int    `json:"priority" binding:"oneof=1 2 3"`
	Tags        string `json:"tags" binding:"max=500"`
	Deadline    string `json:"deadline"`
}

// UpdateTopic 更新课题请求
type UpdateTopic struct {
	Id          uint64 `json:"id" binding:"required"`
	Title       string `json:"title" binding:"required,max=200"`
	Description string `json:"description"`
	Background  string `json:"background"`
	Goal        string `json:"goal"`
	Constraints string `json:"constraints"`
	Priority    int    `json:"priority" binding:"oneof=1 2 3"`
	Tags        string `json:"tags" binding:"max=500"`
	Deadline    string `json:"deadline"`
}

// UpdateTopicStatus 更新课题状态请求
type UpdateTopicStatus struct {
	Id     uint64 `json:"id" binding:"required"`
	Status int    `json:"status" binding:"oneof=0 1 2 3"`
}

// SelectModelForTopic 为课题选用模型请求
type SelectModelForTopic struct {
	TopicId   uint64 `json:"topicId" binding:"required"`
	ModelId   uint64 `json:"modelId" binding:"required"`
	ModelName string `json:"modelName"`
}

// CompleteTopic 完成课题请求
type CompleteTopic struct {
	Id uint64 `json:"id" binding:"required"`
}

// ArchiveTopic 归档课题请求
type ArchiveTopic struct {
	Id uint64 `json:"id" binding:"required"`
}

// SearchTopic 课题搜索条件
type SearchTopic struct {
	Page     int64  `json:"page" form:"page" search:"page"`
	PageSize int64  `json:"pageSize" form:"pageSize" search:"pageSize"`
	Id       uint64 `json:"id" form:"id" search:"type:eq;column:id;table:topics"`
	Title    string `json:"title" form:"title" search:"type:like;column:title;table:topics"`
	UserId   uint64 `json:"userId" form:"userId" search:"type:eq;column:user_id;table:topics"`
	ModelId  uint64 `json:"modelId" form:"modelId" search:"type:eq;column:model_id;table:topics"`
	Status   *int   `json:"status" form:"status" search:"type:eq;column:status;table:topics"`
	Priority int    `json:"priority" form:"priority" search:"type:eq;column:priority;table:topics"`
}

// DelTopic 删除课题请求
type DelTopic struct {
	Ids []uint64 `json:"ids" binding:"required,min=1"`
}

// ==================== 响应DTO ====================

// TopicInfo 课题信息DTO
type TopicInfo struct {
	Id            uint64      `json:"id"`
	Title         string      `json:"title"`
	Description   string      `json:"description"`
	Background    string      `json:"background"`
	Goal          string      `json:"goal"`
	Constraints   string      `json:"constraints"`
	Status        int         `json:"status"`
	StatusText    string      `json:"statusText"`
	UserId        uint64      `json:"userId"`
	ModelId       uint64      `json:"modelId"`
	ModelName     string      `json:"modelName"`
	Priority      int         `json:"priority"`
	PriorityText  string      `json:"priorityText"`
	Tags          []string    `json:"tags"`
	AnalysisCount int         `json:"analysisCount"`
	ActionCount   int         `json:"actionCount"`
	Deadline      interface{} `json:"deadline,omitempty"`
	CompleteTime  interface{} `json:"completeTime,omitempty"`
	CreatedAt     interface{} `json:"createdAt"`
	UpdatedAt     interface{} `json:"updatedAt"`
}

// TopicDetail 课题详情DTO
type TopicDetail struct {
	Id            uint64      `json:"id"`
	Title         string      `json:"title"`
	Description   string      `json:"description"`
	Background    string      `json:"background"`
	Goal          string      `json:"goal"`
	Constraints   string      `json:"constraints"`
	Status        int         `json:"status"`
	StatusText    string      `json:"statusText"`
	UserId        uint64      `json:"userId"`
	ModelId       uint64      `json:"modelId"`
	ModelName     string      `json:"modelName"`
	Priority      int         `json:"priority"`
	PriorityText  string      `json:"priorityText"`
	Tags          []string    `json:"tags"`
	AnalysisCount int         `json:"analysisCount"`
	ActionCount   int         `json:"actionCount"`
	Deadline      interface{} `json:"deadline,omitempty"`
	CompleteTime  interface{} `json:"completeTime,omitempty"`
	CreatedAt     interface{} `json:"createdAt"`
	UpdatedAt     interface{} `json:"updatedAt"`
}

// ListTopicResponse 课题列表响应
type ListTopicResponse struct {
	Page     int64        `json:"page"`
	PageSize int64        `json:"pageSize"`
	Total    int64        `json:"total"`
	List     []*TopicInfo `json:"list"`
}

// TopicStatistics 课题统计
type TopicStatistics struct {
	Total         int64 `json:"total"`
	DraftCount    int64 `json:"draftCount"`    // 草稿
	ProgressCount int64 `json:"progressCount"` // 进行中
	CompleteCount int64 `json:"completeCount"` // 已完成
	ArchiveCount  int64 `json:"archiveCount"`  // 已归档
}
