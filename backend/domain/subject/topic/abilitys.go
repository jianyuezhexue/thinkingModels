package topic

// TopicAbility 课题能力接口
type TopicAbility interface {
	// 状态管理
	Activate() error   // 激活/开始
	Complete() error   // 完成
	Archive() error    // 归档
	Reopen() error     // 重新打开

	// 模型关联
	SelectModel(modelId uint64) error
	RemoveModel() error

	// 标签管理
	AddTag(tag string) error
	RemoveTag(tag string) error
	SetTags(tags []string)
}
