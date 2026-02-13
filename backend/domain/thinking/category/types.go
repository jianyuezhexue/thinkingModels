package category

// ==================== 请求DTO ====================

// CreateCategory 创建分类请求
type CreateCategory struct {
	ParentId    uint64 `json:"parentId"`
	Name        string `json:"name" binding:"required,max=50"`
	Code        string `json:"code" binding:"required,max=50"`
	Description string `json:"description" binding:"max=200"`
	Icon        string `json:"icon" binding:"max=100"`
	Sort        int    `json:"sort"`
}

// UpdateCategory 更新分类请求
type UpdateCategory struct {
	Id          uint64 `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required,max=50"`
	Description string `json:"description" binding:"max=200"`
	Icon        string `json:"icon" binding:"max=100"`
	Sort        int    `json:"sort"`
}

// MoveCategory 移动分类请求
type MoveCategory struct {
	Id             uint64 `json:"id" binding:"required"`
	TargetParentId uint64 `json:"targetParentId"`
	Sort           int    `json:"sort"`
}

// UpdateCategoryStatus 更新分类状态请求
type UpdateCategoryStatus struct {
	Id     uint64 `json:"id" binding:"required"`
	Status int    `json:"status" binding:"oneof=0 1"`
}

// SearchCategory 分类搜索条件
type SearchCategory struct {
	Page     int64  `json:"page" form:"page" search:"page"`
	PageSize int64  `json:"pageSize" form:"pageSize" search:"pageSize"`
	Id       uint64 `json:"id" form:"id" search:"type:eq;column:id;table:model_categories"`
	ParentId uint64 `json:"parentId" form:"parentId" search:"type:eq;column:parent_id;table:model_categories"`
	Name     string `json:"name" form:"name" search:"type:like;column:name;table:model_categories"`
	Code     string `json:"code" form:"code" search:"type:eq;column:code;table:model_categories"`
	Status   int    `json:"status" form:"status" search:"type:eq;column:status;table:model_categories"`
	Level    int    `json:"level" form:"level" search:"type:eq;column:level;table:model_categories"`
}

// DelCategory 删除分类请求
type DelCategory struct {
	Ids []uint64 `json:"ids" binding:"required,min=1"`
}

// ==================== 响应DTO ====================

// CategoryInfo 分类信息DTO
type CategoryInfo struct {
	Id          uint64      `json:"id"`
	ParentId    uint64      `json:"parentId"`
	Name        string      `json:"name"`
	Code        string      `json:"code"`
	Description string      `json:"description"`
	Icon        string      `json:"icon"`
	Sort        int         `json:"sort"`
	Level       int         `json:"level"`
	Path        string      `json:"path"`
	Status      int         `json:"status"`
	ModelCount  int         `json:"modelCount"`
	CreatedAt   interface{} `json:"createdAt"`
	UpdatedAt   interface{} `json:"updatedAt"`
}

// CategoryTree 分类树节点
type CategoryTree struct {
	Id          uint64          `json:"id"`
	ParentId    uint64          `json:"parentId"`
	Name        string          `json:"name"`
	Code        string          `json:"code"`
	Description string          `json:"description"`
	Icon        string          `json:"icon"`
	Sort        int             `json:"sort"`
	Level       int             `json:"level"`
	Path        string          `json:"path"`
	Status      int             `json:"status"`
	ModelCount  int             `json:"modelCount"`
	Children    []*CategoryTree `json:"children"`
}

// ListCategoryResponse 分类列表响应
type ListCategoryResponse struct {
	Page     int64           `json:"page"`
	PageSize int64           `json:"pageSize"`
	Total    int64           `json:"total"`
	List     []*CategoryInfo `json:"list"`
}
