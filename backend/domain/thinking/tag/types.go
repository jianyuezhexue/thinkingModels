package tag

// ==================== 请求DTO ====================

// AddTagsToModel 给模型添加标签请求
type AddTagsToModel struct {
	ModelId  uint64   `json:"modelId" binding:"required"`
	TagNames []string `json:"tagNames" binding:"required,min=1"`
}

// RemoveTagsFromModel 从模型移除标签请求
type RemoveTagsFromModel struct {
	ModelId  uint64   `json:"modelId" binding:"required"`
	TagNames []string `json:"tagNames" binding:"required,min=1"`
}

// SearchTag 标签搜索条件
type SearchTag struct {
	Page     int64  `json:"page" form:"page" search:"page"`
	PageSize int64  `json:"pageSize" form:"pageSize" search:"pageSize"`
	ModelId  uint64 `json:"modelId" form:"modelId" search:"type:eq;column:model_id;table:model_tags"`
	TagName  string `json:"tagName" form:"tagName" search:"type:like;column:tag_name;table:model_tags"`
}

// ==================== 响应DTO ====================

// TagInfo 标签信息DTO
type TagInfo struct {
	Id        uint64 `json:"id"`
	ModelId   uint64 `json:"modelId"`
	TagName   string `json:"tagName"`
	Count     int    `json:"count,omitempty"` // 使用次数（热门标签用）
	CreatedAt string `json:"createdAt"`
}

// ListTagResponse 标签列表响应
type ListTagResponse struct {
	Page     int64      `json:"page"`
	PageSize int64      `json:"pageSize"`
	Total    int64      `json:"total"`
	List     []*TagInfo `json:"list"`
}

// HotTagResponse 热门标签响应
type HotTagResponse struct {
	TagName string `json:"tagName"`
	Count   int    `json:"count"`
}

// HotTag 热门标签（数据库查询用）
type HotTag struct {
	TagName string `json:"tagName" gorm:"column:tag_name"`
	Count   int    `json:"count" gorm:"column:count"`
}
