package action

// ==================== 请求DTO ====================

// CreateAction 创建行动请求
type CreateAction struct {
	Title          string `json:"title" binding:"required,max=200"`
	Description    string `json:"description"`
	TopicId        uint64 `json:"topicId"`
	TopicTitle     string `json:"topicTitle"`
	AnalysisId     uint64 `json:"analysisId"`
	Priority       int    `json:"priority" binding:"oneof=1 2 3"`
	Deadline       string `json:"deadline"`
	GuidePrinciple string `json:"guidePrinciple"`
}

// UpdateAction 更新行动请求
type UpdateAction struct {
	Id             uint64 `json:"id" binding:"required"`
	Title          string `json:"title" binding:"required,max=200"`
	Description    string `json:"description"`
	Priority       int    `json:"priority" binding:"oneof=1 2 3"`
	Deadline       string `json:"deadline"`
	GuidePrinciple string `json:"guidePrinciple"`
}

// UpdateActionProgress 更新行动进度请求
type UpdateActionProgress struct {
	Id       uint64 `json:"id" binding:"required"`
	Progress int    `json:"progress" binding:"min=0,max=100"`
}

// CompleteAction 完成行动请求
type CompleteAction struct {
	Id uint64 `json:"id" binding:"required"`
}

// CancelAction 取消行动请求
type CancelAction struct {
	Id uint64 `json:"id" binding:"required"`
}

// CreateActionsFromAnalysis 从分析导出行动请求
type CreateActionsFromAnalysis struct {
	AnalysisId uint64               `json:"analysisId" binding:"required"`
	TopicId    uint64               `json:"topicId"`
	TopicTitle string               `json:"topicTitle"`
	Actions    []*CreateActionBatch `json:"actions" binding:"required,min=1"`
}

// CreateActionBatch 批量创建行动项
type CreateActionBatch struct {
	Title          string `json:"title" binding:"required,max=200"`
	Description    string `json:"description"`
	Priority       int    `json:"priority"`
	Deadline       string `json:"deadline"`
	GuidePrinciple string `json:"guidePrinciple"`
}

// SearchAction 行动搜索条件
type SearchAction struct {
	Page       int64  `json:"page" form:"page" search:"page"`
	PageSize   int64  `json:"pageSize" form:"pageSize" search:"pageSize"`
	Id         uint64 `json:"id" form:"id" search:"type:eq;column:id;table:actions"`
	Title      string `json:"title" form:"title" search:"type:like;column:title;table:actions"`
	UserId     uint64 `json:"userId" form:"userId" search:"type:eq;column:user_id;table:actions"`
	TopicId    uint64 `json:"topicId" form:"topicId" search:"type:eq;column:topic_id;table:actions"`
	AnalysisId uint64 `json:"analysisId" form:"analysisId" search:"type:eq;column:analysis_id;table:actions"`
	Status     *int   `json:"status" form:"status" search:"type:eq;column:status;table:actions"`
	Priority   int    `json:"priority" form:"priority" search:"type:eq;column:priority;table:actions"`
	IsOverdue  *bool  `json:"isOverdue" form:"isOverdue"`
}

// DelAction 删除行动请求
type DelAction struct {
	Ids []uint64 `json:"ids" binding:"required,min=1"`
}

// ==================== 响应DTO ====================

// ActionInfo 行动信息DTO
type ActionInfo struct {
	Id             uint64      `json:"id"`
	Title          string      `json:"title"`
	Description    string      `json:"description"`
	UserId         uint64      `json:"userId"`
	TopicId        uint64      `json:"topicId"`
	TopicTitle     string      `json:"topicTitle"`
	AnalysisId     uint64      `json:"analysisId"`
	Priority       int         `json:"priority"`
	PriorityText   string      `json:"priorityText"`
	Status         int         `json:"status"`
	StatusText     string      `json:"statusText"`
	Progress       int         `json:"progress"`
	Deadline       interface{} `json:"deadline,omitempty"`
	CompletedAt    interface{} `json:"completedAt,omitempty"`
	GuidePrinciple string      `json:"guidePrinciple"`
	FollowUpCount  int         `json:"followUpCount"`
	IsOverdue      bool        `json:"isOverdue"`
	CreatedAt      interface{} `json:"createdAt"`
	UpdatedAt      interface{} `json:"updatedAt"`
}

// ListActionResponse 行动列表响应
type ListActionResponse struct {
	Page     int64         `json:"page"`
	PageSize int64         `json:"pageSize"`
	Total    int64         `json:"total"`
	List     []*ActionInfo `json:"list"`
}

// ActionStatistics 行动统计
type ActionStatistics struct {
	Total         int64 `json:"total"`
	PendingCount  int64 `json:"pendingCount"`  // 待执行
	ProgressCount int64 `json:"progressCount"` // 进行中
	CompleteCount int64 `json:"completeCount"` // 已完成
	CancelCount   int64 `json:"cancelCount"`   // 已取消
	OverdueCount  int64 `json:"overdueCount"`  // 逾期
}

// UpdateProgress 更新进度请求
type UpdateProgress struct {
	Id       uint64 `json:"id" binding:"required"`
	Progress int    `json:"progress" binding:"min=0,max=100"`
	Note     string `json:"note"`
}

// ActionDetail 行动详情DTO
type ActionDetail struct {
	Id             uint64      `json:"id"`
	TopicId        uint64      `json:"topicId"`
	AnalysisId     uint64      `json:"analysisId"`
	Title          string      `json:"title"`
	Description    string      `json:"description"`
	Priority       int         `json:"priority"`
	Status         int         `json:"status"`
	Progress       int         `json:"progress"`
	ExpectedResult string      `json:"expectedResult"`
	ActualResult   string      `json:"actualResult"`
	Deadline       interface{} `json:"deadline"`
	CompletedAt    interface{} `json:"completedAt"`
	FollowUpCount  int         `json:"followUpCount"`
	CreatedAt      interface{} `json:"createdAt"`
	UpdatedAt      interface{} `json:"updatedAt"`
}
