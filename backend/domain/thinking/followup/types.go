package followup

// ==================== 请求DTO ====================

// CreateFollowUp 创建跟进请求
type CreateFollowUp struct {
	ActionId       uint64 `json:"actionId" binding:"required"`
	Content        string `json:"content" binding:"required"`
	ProgressBefore int    `json:"progressBefore"`
	ProgressAfter  int    `json:"progressAfter" binding:"min=0,max=100"`
}

// UpdateFollowUp 更新跟进请求
type UpdateFollowUp struct {
	Id      uint64 `json:"id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// SearchFollowUp 跟进搜索条件
type SearchFollowUp struct {
	Page     int64  `json:"page" form:"page" search:"page"`
	PageSize int64  `json:"pageSize" form:"pageSize" search:"pageSize"`
	Id       uint64 `json:"id" form:"id" search:"type:eq;column:id;table:action_followups"`
	ActionId uint64 `json:"actionId" form:"actionId" search:"type:eq;column:action_id;table:action_followups"`
	UserId   uint64 `json:"userId" form:"userId" search:"type:eq;column:user_id;table:action_followups"`
}

// DelFollowUp 删除跟进请求
type DelFollowUp struct {
	Ids []uint64 `json:"ids" binding:"required,min=1"`
}

// ==================== 响应DTO ====================

// FollowUpInfo 跟进信息DTO
type FollowUpInfo struct {
	Id             uint64 `json:"id"`
	ActionId       uint64 `json:"actionId"`
	UserId         uint64 `json:"userId"`
	UserName       string `json:"userName"`
	Content        string `json:"content"`
	ProgressBefore int    `json:"progressBefore"`
	ProgressAfter  int    `json:"progressAfter"`
	ProgressDelta  int    `json:"progressDelta"` // 进度变化
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}

// ListFollowUpResponse 跟进列表响应
type ListFollowUpResponse struct {
	Page     int64           `json:"page"`
	PageSize int64           `json:"pageSize"`
	Total    int64           `json:"total"`
	List     []*FollowUpInfo `json:"list"`
}
