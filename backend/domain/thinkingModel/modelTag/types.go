package modelTag

// ==================== 请求REQ ====================

// AddTagsToModelRequest 为模型添加标签请求
type AddTagsToModelRequest struct {
	ModelId uint64    `json:"modelId" binding:"required"`    // 模型ID
	Tags    []TagItem `json:"tags" binding:"required,min=1"` // 标签列表
}

// TagItem 标签项
type TagItem struct {
	TagId   uint64 `json:"tagId"`                      // 标签ID（来自 master 领域公共 tag）
	TagName string `json:"tagName" binding:"required"` // 标签名称
}

// RemoveTagsFromModelRequest 从模型移除标签请求
type RemoveTagsFromModelRequest struct {
	ModelId uint64   `json:"modelId" binding:"required"`      // 模型ID
	TagIds  []uint64 `json:"tagIds" binding:"required,min=1"` // 标签ID列表
}

// SearchModelTag 模型标签搜索条件
type SearchModelTag struct {
	Page     int64  `json:"page" form:"page" search:"page"`
	PageSize int64  `json:"pageSize" form:"pageSize" search:"pageSize"`
	ModelId  uint64 `json:"modelId" form:"modelId" search:"type:eq;column:model_id;table:thinking_model_tags"`
	TagId    uint64 `json:"tagId" form:"tagId" search:"type:eq;column:tag_id;table:thinking_model_tags"`
	TagName  string `json:"tagName" form:"tagName" search:"type:like;column:tag_name;table:thinking_model_tags"`
}

// DelModelTag 删除模型标签请求
type DelModelTag struct {
	Ids []uint64 `json:"ids" binding:"required,min=1"`
}

// IncreaseHeatRequest 增加热度请求
type IncreaseHeatRequest struct {
	Id    uint64 `json:"id" binding:"required"`          // 模型标签ID
	Delta int    `json:"delta" binding:"required,min=1"` // 增加的热度值
}

// ==================== 响应RESP ====================

// TagInfo 标签信息DTO
type TagInfo struct {
	Id        uint64 `json:"id"`
	ModelId   uint64 `json:"modelId"`
	TagId     uint64 `json:"tagId"`
	TagName   string `json:"tagName"`
	Heat      int    `json:"heat"`
	CreatedAt string `json:"createdAt"`
}

// ListModelTagResponse 模型标签列表响应
type ListModelTagResponse struct {
	Page     int64      `json:"page"`
	PageSize int64      `json:"pageSize"`
	Total    int64      `json:"total"`
	List     []*TagInfo `json:"list"`
}

// HotTagResponse 热门标签响应
type HotTagResponse struct {
	List []*HotTagItem `json:"list"`
}

// HotTagItem 热门标签项
type HotTagItem struct {
	TagName string `json:"tagName"`
	Heat    int    `json:"heat"`
	Count   int64  `json:"count"` // 使用该标签的模型数量
}
