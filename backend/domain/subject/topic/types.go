package topic

// ==================== 请求DTO ====================

// CreateTopic 创建课题请求
type CreateTopic struct {
	Title       string `json:"title" binding:"required,max=200"`  // 课题标题
	Description string `json:"description" binding:"max=2000"`    // 课题描述
	Priority    int    `json:"priority" binding:"oneof=1 2 3"`    // 优先级: 1=低, 2=中, 3=高
	Tags        string `json:"tags" binding:"max=500"`            // 标签，逗号分隔
	Deadline    string `json:"deadline"`                          // 截止日期（ISO格式字符串）
}

// UpdateTopic 更新课题请求
type UpdateTopic struct {
	Id          uint64 `json:"id" binding:"required"`
	Title       string `json:"title" binding:"required,max=200"`
	Description string `json:"description" binding:"max=2000"`
	Priority    int    `json:"priority" binding:"oneof=1 2 3"`
	Tags        string `json:"tags" binding:"max=500"`
	Deadline    string `json:"deadline"`
}

// UpdateTopicStatus 更新课题状态请求
type UpdateTopicStatus struct {
	Id     uint64 `json:"id" binding:"required"`
	Status int    `json:"status" binding:"oneof=0 1 2"` // 0=进行中, 1=已完成, 2=已归档
}

// UpdateTopicModel 更新课题使用的模型请求
type UpdateTopicModel struct {
	Id      uint64 `json:"id" binding:"required"`
	ModelId uint64 `json:"modelId" binding:"required"`
}

// SearchTopic 课题搜索条件
type SearchTopic struct {
	Page       int64  `json:"page" form:"page" search:"page"`
	PageSize   int64  `json:"pageSize" form:"pageSize" search:"pageSize"`
	Id         uint64 `json:"id" form:"id" search:"type:eq;column:id;table:topics"`
	Title      string `json:"title" form:"title" search:"type:like;column:title;table:topics"`
	Status     int    `json:"status" form:"status" search:"type:eq;column:status;table:topics"`
	Priority   int    `json:"priority" form:"priority" search:"type:eq;column:priority;table:topics"`
	UserId     uint64 `json:"userId" form:"userId" search:"type:eq;column:user_id;table:topics"`
	ModelId    uint64 `json:"modelId" form:"modelId" search:"type:eq;column:model_id;table:topics"`
	Tags       string `json:"tags" form:"tags" search:"type:like;column:tags;table:topics"`
}

// DelTopic 删除课题请求
type DelTopic struct {
	Ids []uint64 `json:"ids" binding:"required,min=1"`
}

// CompleteTopic 完成课题请求
type CompleteTopic struct {
	Id uint64 `json:"id" binding:"required"`
}

// ArchiveTopic 归档课题请求
type ArchiveTopic struct {
	Id uint64 `json:"id" binding:"required"`
}

// ==================== 响应DTO ====================

// TopicInfo 课题信息DTO
type TopicInfo struct {
	Id            uint64 `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Status        int    `json:"status"`
	StatusText    string `json:"statusText"`
	Priority      int    `json:"priority"`
	PriorityText  string `json:"priorityText"`
	Tags          string `json:"tags"`
	TagsList      []string `json:"tagsList"`
	Deadline      string `json:"deadline"`
	CompleteTime  string `json:"completeTime"`
	UserId        uint64 `json:"userId"`
	ModelId       uint64 `json:"modelId"`
	ModelName     string `json:"modelName"` // 关联的模型名称
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

// TopicDetail 课题详情DTO（包含分析记录）
type TopicDetail struct {
	TopicInfo
	Analyses []*TopicAnalysisBrief `json:"analyses"` // 课题的分析记录列表
}

// TopicAnalysisBrief 课题分析记录简要信息
type TopicAnalysisBrief struct {
	Id           uint64 `json:"id"`
	ModelId      uint64 `json:"modelId"`
	ModelName    string `json:"modelName"`
	Version      int    `json:"version"`
	IsCurrent    bool   `json:"isCurrent"`
	CreatedAt    string `json:"createdAt"`
}

// ListTopicResponse 课题列表响应
type ListTopicResponse struct {
	Page     int64        `json:"page"`
	PageSize int64        `json:"pageSize"`
	Total    int64        `json:"total"`
	List     []*TopicInfo `json:"list"`
}

// TopicStatistics 课题统计信息
type TopicStatistics struct {
	Total      int64 `json:"total"`
	Active     int64 `json:"active"`     // 进行中
	Completed  int64 `json:"completed"`  // 已完成
	Archived   int64 `json:"archived"`   // 已归档
	HighPriority int64 `json:"highPriority"` // 高优先级
}
