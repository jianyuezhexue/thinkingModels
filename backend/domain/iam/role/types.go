package role

// ==================== 请求REQ ====================

// CreateRole 创建角色请求
type CreateRole struct {
	RoleName    string `json:"roleName" binding:"required,min=2,max=20"`          // 角色名称
	RoleCode    string `json:"roleCode" binding:"required,alphanum,min=2,max=20"` // 角色编码
	Description string `json:"description" binding:"max=200"`                     // 角色描述
	Status      int    `json:"status"`                                            // 状态：0=禁用，1=正常
	Sort        int    `json:"sort"`                                              // 排序
	MenuIds     string `json:"menuIds"`                                           // 菜单ID列表，逗号分隔
}

// UpdateRole 更新角色请求
type UpdateRole struct {
	ID          uint64 `json:"id" binding:"required"`                    // 角色ID
	RoleName    string `json:"roleName" binding:"required,min=2,max=20"` // 角色名称
	Description string `json:"description" binding:"max=200"`            // 角色描述
	Status      int    `json:"status"`                                   // 状态：0=禁用，1=正常
	Sort        int    `json:"sort"`                                     // 排序
	MenuIds     string `json:"menuIds"`                                  // 菜单ID列表，逗号分隔
}

// UpdateRolePermission 更新角色权限请求
type UpdateRolePermission struct {
	ID      uint64   `json:"id" binding:"required"` // 角色ID
	MenuIds []string `json:"menuIds"`               // 菜单ID列表
}

// ==================== 响应RESP ====================

// RoleInfo 角色信息DTO
type RoleInfo struct {
	ID          uint64 `json:"id"`
	RoleName    string `json:"roleName"`
	RoleCode    string `json:"roleCode"`
	Description string `json:"description"`
	Status      int    `json:"status"`
	Sort        int    `json:"sort"`
	MenuIds     string `json:"menuIds"`
	UserCount   int64  `json:"userCount"` // 关联用户数
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// ==================== 搜索条件 ====================

// SearchRole 角色搜索条件
type SearchRole struct {
	Page     int64  `json:"page" form:"page" search:"page"`                                           // 分页
	PageSize int64  `json:"pageSize" form:"pageSize" search:"pageSize"`                               // 分页大小
	RoleName string `json:"roleName" form:"roleName" search:"type:like;column:role_name;table:roles"` // 角色名称模糊查询
	RoleCode string `json:"roleCode" form:"roleCode" search:"type:eq;column:role_code;table:roles"`   // 角色编码精确查询
	Status   int    `json:"status" form:"status" search:"type:eq;column:status;table:roles"`          // 状态
}

// ListRoleResponse 角色列表响应
type ListRoleResponse struct {
	Page     int64       `json:"page" comment:"页数"`
	PageSize int64       `json:"pageSize" comment:"每页数量"`
	Total    int64       `json:"total" comment:"总条数"`
	List     []*RoleInfo `json:"list" comment:"数据"`
}

// DelRole 删除角色请求
type DelRole struct {
	Ids []uint64 `json:"ids" binding:"required"`
}

// ==================== 菜单树 ====================

// MenuNode 菜单树节点
type MenuNode struct {
	ID       string      `json:"id"`
	Label    string      `json:"label"`
	Children []*MenuNode `json:"children,omitempty"`
}
