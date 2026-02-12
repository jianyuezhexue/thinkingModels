package model

// ==================== 请求DTO ====================

// CreateModel 创建思维模型请求
type CreateModel struct {
	Name          string  `json:"name" binding:"required,max=100"`           // 模型名称
	Code          string  `json:"code" binding:"required,max=50"`            // 模型编码
	Description   string  `json:"description" binding:"max=2000"`            // 模型描述
	CoverImage    string  `json:"coverImage" binding:"max=255"`              // 封面图片
	Icon          string  `json:"icon" binding:"max=255"`                    // 模型图标
	CategoryId    uint64  `json:"categoryId" binding:"required"`             // 所属分类ID
	Price         float64 `json:"price" binding:"min=0"`                     // 价格
	Content       string  `json:"content" binding:"required"`                // 模型内容JSON
	Difficulty    int     `json:"difficulty" binding:"oneof=1 2 3"`          // 难度
	EstimatedTime int     `json:"estimatedTime" binding:"min=1"`             // 预计时间
	Version       string  `json:"version" binding:"max=20"`                  // 版本号
	AuthorName    string  `json:"authorName" binding:"max=100"`              // 作者名称
}

// UpdateModel 更新思维模型请求
type UpdateModel struct {
	Id            uint64  `json:"id" binding:"required"`
	Name          string  `json:"name" binding:"required,max=100"`
	Description   string  `json:"description" binding:"max=2000"`
	CoverImage    string  `json:"coverImage" binding:"max=255"`
	Icon          string  `json:"icon" binding:"max=255"`
	CategoryId    uint64  `json:"categoryId" binding:"required"`
	Price         float64 `json:"price" binding:"min=0"`
	Content       string  `json:"content" binding:"required"`
	Difficulty    int     `json:"difficulty" binding:"oneof=1 2 3"`
	EstimatedTime int     `json:"estimatedTime" binding:"min=1"`
	Version       string  `json:"version" binding:"max=20"`
	AuthorName    string  `json:"authorName" binding:"max=100"`
}

// UpdateModelStatus 更新模型状态请求
type UpdateModelStatus struct {
	Id     uint64 `json:"id" binding:"required"`
	Status int    `json:"status" binding:"oneof=0 1 2"` // 0=草稿, 1=已发布, 2=已下架
}

// PublishModel 发布模型请求
type PublishModel struct {
	Id uint64 `json:"id" binding:"required"`
}

// UnpublishModel 下架模型请求
type UnpublishModel struct {
	Id uint64 `json:"id" binding:"required"`
}

// SearchModel 思维模型搜索条件
type SearchModel struct {
	Page        int64   `json:"page" form:"page" search:"page"`                                    // 分页
	PageSize    int64   `json:"pageSize" form:"pageSize" search:"pageSize"`                        // 分页大小
	Id          uint64  `json:"id" form:"id" search:"type:eq;column:id;table:thinking_models"`      // ID精确匹配
	Name        string  `json:"name" form:"name" search:"type:like;column:name;table:thinking_models"` // 名称模糊查询
	Code        string  `json:"code" form:"code" search:"type:eq;column:code;table:thinking_models"`   // 编码精确匹配
	CategoryId  uint64  `json:"categoryId" form:"categoryId" search:"type:eq;column:category_id;table:thinking_models"` // 分类ID
	Category    string  `json:"category" form:"category"`                                            // 分类字符串（前端传入，如 'business'）
	Status      int     `json:"status" form:"status" search:"type:eq;column:status;table:thinking_models"` // 状态
	Difficulty  int     `json:"difficulty" form:"difficulty" search:"type:eq;column:difficulty;table:thinking_models"` // 难度
	MinPrice    float64 `json:"minPrice" form:"minPrice" search:"type:gte;column:price;table:thinking_models"` // 最低价格
	MaxPrice    float64 `json:"maxPrice" form:"maxPrice" search:"type:lte;column:price;table:thinking_models"` // 最高价格
	AuthorId    uint64  `json:"authorId" form:"authorId" search:"type:eq;column:author_id;table:thinking_models"` // 作者ID
	IsFree      *bool   `json:"isFree" form:"isFree"`                                                // 是否免费
	SortBy      string  `json:"sortBy" form:"sortBy"`                                                // 排序方式
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
	Bio    string `json:"bio,omitempty"`
}

// ModelStats 模型统计数据
type ModelStats struct {
	Adoptions   int `json:"adoptions"`
	Practices   int `json:"practices"`
	Discussions int `json:"discussions"`
	Forks       int `json:"forks"`
	Likes       int `json:"likes"`
}

// ModelInfo 思维模型信息DTO - 匹配前端期望的数据结构
type ModelInfo struct {
	Id          string      `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Cover       string      `json:"cover"`
	Author      ModelAuthor `json:"author"`
	IsFree      bool        `json:"isFree"`
	Price       float64     `json:"price,omitempty"`
	Category    string      `json:"category"`
	CategoryId  string      `json:"categoryId,omitempty"`
	Tags        []string    `json:"tags"`
	Stats       ModelStats  `json:"stats"`
	UpdatedAt   string      `json:"updatedAt"`
	Status      int         `json:"status,omitempty"`
	Content     string      `json:"content,omitempty"`
}

// ModelDetail 思维模型详情DTO（包含完整内容）
type ModelDetail struct {
	ModelInfo
}

// ListModelResponse 思维模型列表响应
type ListModelResponse struct {
	Page     int64        `json:"page"`
	PageSize int64        `json:"pageSize"`
	Total    int64        `json:"total"`
	List     []*ModelInfo `json:"list"`
}

// ModelContent 模型内容结构（用于解析JSON）
type ModelContent struct {
	Title       string            `json:"title"`       // 模型标题
	Description string            `json:"description"` // 模型描述
	Fields      []ModelField      `json:"fields"`      // 字段定义
	Steps       []ModelStep       `json:"steps"`       // 步骤定义
	Metadata    map[string]interface{} `json:"metadata"` // 元数据
}

// ModelField 模型字段定义
type ModelField struct {
	Key         string `json:"key"`         // 字段标识
	Label       string `json:"label"`       // 字段标签
	Type        string `json:"type"`        // 字段类型: text, textarea, number, select, radio, checkbox
	Required    bool   `json:"required"`    // 是否必填
	Placeholder string `json:"placeholder"` // 提示文本
	Options     []Option `json:"options"`   // 选项（针对select/radio/checkbox）
	DefaultValue interface{} `json:"defaultValue"` // 默认值
	HelpText    string `json:"helpText"`    // 帮助文本
}

// Option 选项定义
type Option struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// ModelStep 模型步骤定义
type ModelStep struct {
	Id          string `json:"id"`          // 步骤ID
	Title       string `json:"title"`       // 步骤标题
	Description string `json:"description"` // 步骤描述
	Fields      []string `json:"fields"`    // 该步骤包含的字段key列表
	Order       int    `json:"order"`       // 顺序
}
