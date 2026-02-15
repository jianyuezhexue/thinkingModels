package model

// ==================== 状态常量 ====================

const (
	StatusDraft     = 0 // 草稿
	StatusPublished = 1 // 已发布
	StatusUnpublish = 2 // 已下架
	StatusReviewing = 3 // 审核中
	StatusRejected  = 4 // 已驳回
)

// ==================== 请求DTO ====================

// CreateModel 创建思维模型请求
type CreateModel struct {
	Name          string  `json:"name" binding:"required,max=100"`
	Code          string  `json:"code" binding:"required,max=50"`
	Description   string  `json:"description" binding:"max=2000"`
	CoverImage    string  `json:"coverImage" binding:"max=500"`
	Icon          string  `json:"icon" binding:"max=200"`
	CategoryId    uint64  `json:"categoryId" binding:"required"`
	Price         float64 `json:"price" binding:"min=0"`
	Content       string  `json:"content"`
	Overview      string  `json:"overview"`
	Difficulty    int     `json:"difficulty" binding:"oneof=1 2 3"`
	EstimatedTime int     `json:"estimatedTime" binding:"min=1"`
	Version       string  `json:"version" binding:"max=20"`
	AuthorName    string  `json:"authorName" binding:"max=100"`
	IsOfficial    int     `json:"isOfficial"`
}

// UpdateModel 更新思维模型请求
type UpdateModel struct {
	Id            uint64  `json:"id" binding:"required"`
	Name          string  `json:"name" binding:"required,max=100"`
	Description   string  `json:"description" binding:"max=2000"`
	CoverImage    string  `json:"coverImage" binding:"max=500"`
	Icon          string  `json:"icon" binding:"max=200"`
	CategoryId    uint64  `json:"categoryId" binding:"required"`
	Price         float64 `json:"price" binding:"min=0"`
	Content       string  `json:"content"`
	Overview      string  `json:"overview"`
	Difficulty    int     `json:"difficulty" binding:"oneof=1 2 3"`
	EstimatedTime int     `json:"estimatedTime" binding:"min=1"`
	Version       string  `json:"version" binding:"max=20"`
	AuthorName    string  `json:"authorName" binding:"max=100"`
}

// PublishModel 发布模型请求
type PublishModel struct {
	Id uint64 `json:"id" binding:"required"`
}

// UnpublishModel 下架模型请求
type UnpublishModel struct {
	Id uint64 `json:"id" binding:"required"`
}

// ReviewModel 审核模型请求
type ReviewModel struct {
	Id       uint64 `json:"id" binding:"required"`
	Approved bool   `json:"approved"`       // true=通过, false=驳回
	Note     string `json:"note" binding:"max=500"` // 审核意见
}

// ForkModel 引用创建模型请求
type ForkModel struct {
	SourceModelId uint64 `json:"sourceModelId" binding:"required"`
	Name          string `json:"name" binding:"required,max=100"`
	Description   string `json:"description" binding:"max=2000"`
}

// SearchModel 思维模型搜索条件
type SearchModel struct {
	Page       int64   `json:"page" form:"page" search:"page"`
	PageSize   int64   `json:"pageSize" form:"pageSize" search:"pageSize"`
	Id         uint64  `json:"id" form:"id" search:"type:eq;column:id;table:thinking_models"`
	Name       string  `json:"name" form:"name" search:"type:like;column:name;table:thinking_models"`
	Code       string  `json:"code" form:"code" search:"type:eq;column:code;table:thinking_models"`
	CategoryId uint64  `json:"categoryId" form:"categoryId" search:"type:eq;column:category_id;table:thinking_models"`
	Status     int     `json:"status" form:"status" search:"type:eq;column:status;table:thinking_models"`
	Difficulty int     `json:"difficulty" form:"difficulty" search:"type:eq;column:difficulty;table:thinking_models"`
	MinPrice   float64 `json:"minPrice" form:"minPrice" search:"type:gte;column:price;table:thinking_models"`
	MaxPrice   float64 `json:"maxPrice" form:"maxPrice" search:"type:lte;column:price;table:thinking_models"`
	AuthorId   uint64  `json:"authorId" form:"authorId" search:"type:eq;column:author_id;table:thinking_models"`
	IsOfficial int     `json:"isOfficial" form:"isOfficial" search:"type:eq;column:is_official;table:thinking_models"`
	IsFree     *bool   `json:"isFree" form:"isFree"`
	SortBy     string  `json:"sortBy" form:"sortBy"`
}

// DelModel 删除思维模型请求
type DelModel struct {
	Ids []uint64 `json:"ids" binding:"required,min=1"`
}

// ==================== 响应DTO ====================

// ModelAuthor 模型作者信息
type ModelAuthor struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

// ModelStats 模型统计数据
type ModelStats struct {
	UsageCount   int `json:"usageCount"`
	AdoptCount   int `json:"adoptCount"`
	LikeCount    int `json:"likeCount"`
	CommentCount int `json:"commentCount"`
}

// ModelInfo 思维模型信息DTO
type ModelInfo struct {
	Id            uint64      `json:"id"`
	Name          string      `json:"name"`
	Code          string      `json:"code"`
	Description   string      `json:"description"`
	CoverImage    string      `json:"coverImage"`
	Icon          string      `json:"icon"`
	CategoryId    uint64      `json:"categoryId"`
	Price         float64     `json:"price"`
	IsFree        bool        `json:"isFree"`
	Overview      string      `json:"overview"`
	Difficulty    int         `json:"difficulty"`
	EstimatedTime int         `json:"estimatedTime"`
	Status        int         `json:"status"`
	Version       string      `json:"version"`
	IsOfficial    bool        `json:"isOfficial"`
	SourceModelId uint64      `json:"sourceModelId,omitempty"`
	Author        ModelAuthor `json:"author"`
	Stats         ModelStats  `json:"stats"`
	Tags          []string    `json:"tags"`
	PublishTime   string      `json:"publishTime,omitempty"`
	ReviewNote    string      `json:"reviewNote,omitempty"`    // 审核意见
	ReviewerName  string      `json:"reviewerName,omitempty"` // 审核人
	ReviewTime    string      `json:"reviewTime,omitempty"`   // 审核时间
	CreatedAt     string      `json:"createdAt"`
	UpdatedAt     string      `json:"updatedAt"`
}

// ModelDetail 思维模型详情DTO（包含完整内容）
type ModelDetail struct {
	ModelInfo
	Content string `json:"content"`
}

// ListModelResponse 思维模型列表响应
type ListModelResponse struct {
	Page     int64        `json:"page"`
	PageSize int64        `json:"pageSize"`
	Total    int64        `json:"total"`
	List     []*ModelInfo `json:"list"`
}
