package model

// 能力接口定义
// 在此定义实体可执行的操作

// ModelAbility 思维模型能力接口
type ModelAbility interface {
	// 发布管理
	Publish() error
	Unpublish() error

	// 内容管理
	UpdateContent(content string) error

	// 统计相关
	IncrementUsageCount()
}
