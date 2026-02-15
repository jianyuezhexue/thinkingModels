package modelTag

// ModelTagAbility 模型标签能力接口定义
type ModelTagAbility interface {
	// 基础能力由 BaseModel 提供

	// 扩展能力
	// GetByModelId(modelId uint64) 根据模型ID获取标签列表
	// AddTagsToModel(modelId uint64, tags []TagInput) 为模型添加标签
	// RemoveTagsFromModel(modelId uint64, tagIds []uint64) 从模型移除标签
	// GetHotTags(limit int) 获取热门标签
	// IncrementHeat(id uint64, delta int) 增加标签热度
}