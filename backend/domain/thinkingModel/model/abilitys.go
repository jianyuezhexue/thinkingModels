package model

// ModelAbility 思维模型能力接口定义
type ModelAbility interface {
	// 基础能力由 BaseModel 提供
	// Create() 创建模型
	// Update() 更新模型
	// Delete() 删除模型
	// LoadById() 根据ID加载
	// List() 列表查询
	// Count() 统计数量

	// 扩展能力
	// Publish() 发布模型
	// Unpublish() 下架模型
	// SubmitForReview() 提交审核
	// Approve() 审核通过
	// Reject() 审核驳回
	// Fork() 派生模型
	// IncrementUsageCount() 增加使用次数
	// IncrementAdoptCount() 增加采纳次数
	// IncrementLikeCount() 增加点赞数
	// IncrementCommentCount() 增加评论数
}