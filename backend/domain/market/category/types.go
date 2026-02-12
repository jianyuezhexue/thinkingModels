package category

// ==================== 请求DTO ====================

// CreateCategory 创建分类请求
type CreateCategory struct {
	ParentId    uint64 `json:"parentId"`              // 父分类ID，0为顶级分类
	Name        string `json:"name" binding:"required,max=100"`  // 分类名称
	Code        string `json:"code" binding:"required,max=50"`   // 分类编码
	Description string `json:"description" binding:"max=500"`    // 分类描述
	Icon        string `json:"icon" binding:"max=255"`           // 分类图标
	SortOrder   int    `json:"sortOrder"`             // 排序顺序
}

// UpdateCategory 更新分类请求
type UpdateCategory struct {
	Id          uint64 `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"max=500"`
	Icon        string `json:"icon" binding:"max=255"`
	SortOrder   int    `json:"sortOrder"`
}

// MoveCategory 移动分类请求（修改父级）
type MoveCategory struct {
	Id       uint64 `json:"id" binding:"required"`
	ParentId uint64 `json:"parentId" binding:"required"` // 新的父分类ID
}

// UpdateCategoryStatus 更新分类状态请求
type UpdateCategoryStatus struct {
	Id     uint64 `json:"id" binding:"required"`
	Status int    `json:"status" binding:"oneof=0 1"` // 0=禁用, 1=启用
}

// SearchCategory 分类搜索条件
type SearchCategory struct {
	Page       int64  `json:"page" form:"page" search:"page"`
	PageSize   int64  `json:"pageSize" form:"pageSize" search:"pageSize"`
	Id         uint64 `json:"id" form:"id" search:"type:eq;column:id;table:model_categories"`
	Name       string `json:"name" form:"name" search:"type:like;column:name;table:model_categories"`
	Code       string `json:"code" form:"code" search:"type:eq;column:code;table:model_categories"`
	ParentId   uint64 `json:"parentId" form:"parentId" search:"type:eq;column:parent_id;table:model_categories"`
	Status     int    `json:"status" form:"status" search:"type:eq;column:status;table:model_categories"`
	Level      int    `json:"level" form:"level" search:"type:eq;column:level;table:model_categories"`
}

// DelCategory 删除分类请求
type DelCategory struct {
	Ids []uint64 `json:"ids" binding:"required,min=1"`
}

// ==================== 响应DTO ====================

// CategoryInfo 分类信息DTO
type CategoryInfo struct {
	Id          uint64 `json:"id"`
	ParentId    uint64 `json:"parentId"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	SortOrder   int    `json:"sortOrder"`
	Level       int    `json:"level"`
	Path        string `json:"path"`
	Status      int    `json:"status"`
	StatusText  string `json:"statusText"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// CategoryTreeNode 分类树节点DTO
type CategoryTreeNode struct {
	CategoryInfo
	Children []*CategoryTreeNode `json:"children"` // 子分类列表
}

// ListCategoryResponse 分类列表响应
type ListCategoryResponse struct {
	Page     int64           `json:"page"`
	PageSize int64           `json:"pageSize"`
	Total    int64           `json:"total"`
	List     []*CategoryInfo `json:"list"`
}

// CategoryTreeResponse 分类树响应
type CategoryTreeResponse struct {
	List []*CategoryTreeNode `json:"list"`
}

// CategoryPathResponse 分类路径响应
type CategoryPathResponse struct {
	Path []*CategoryInfo `json:"path"`
}

// CategoryWithChildren 带子分类的分类信息
type CategoryWithChildren struct {
	CategoryInfo
	ChildrenCount int `json:"childrenCount"`
	ModelCount    int `json:"modelCount"`
}
