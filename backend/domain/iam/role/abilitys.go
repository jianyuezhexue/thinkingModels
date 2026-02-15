package role

// RoleAbility 角色能力接口定义
// 在此定义角色实体可执行的操作
type RoleAbility interface {
	// 基础能力（由 BaseModel 提供）
	// Create() 创建角色
	// Update() 更新角色
	// Delete() 删除角色
	// LoadById() 根据ID加载
	// List() 列表查询
	// Count() 统计数量

	// 扩展能力
	// GetUserCount() 获取角色关联的用户数量
	// CheckRoleCodeExists() 检查角色编码是否存在
}