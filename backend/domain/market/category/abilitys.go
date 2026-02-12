package category

// CategoryAbility 分类能力接口
type CategoryAbility interface {
	// 树形结构操作
	GetChildren() ([]*CategoryEntity, error)
	GetParent() (*CategoryEntity, error)
	GetFullPath() ([]*CategoryEntity, error)

	// 层级管理
	UpdateLevel() error
	BuildPath() error

	// 禁用/启用
	Enable() error
	Disable() error
}
