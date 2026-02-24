package model

import (
	"errors"
)

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

// ==================== 发布相关能力 ====================

// Publish 发布模型
func (m *ModelEntity) Publish() error {
	if m.Status == StatusPublished {
		return errors.New("模型已发布")
	}
	m.Status = StatusPublished
	return nil
}

// Unpublish 下架模型
func (m *ModelEntity) Unpublish() error {
	if m.Status != StatusPublished {
		return errors.New("模型未发布，无法下架")
	}
	m.Status = StatusUnpublish
	return nil
}

// ==================== 审核相关能力 ====================

// SubmitForReview 提交审核
func (m *ModelEntity) SubmitForReview() error {
	if m.Status == StatusReviewing {
		return errors.New("模型正在审核中")
	}
	// if m.Status == StatusPublished {
	// 	return errors.New("模型已发布，无需审核")
	// }
	m.Status = StatusReviewing
	return nil
}

// Approve 审核通过
func (m *ModelEntity) Approve(reviewerId uint64, reviewerName, note string) {
	m.Status = StatusPublished
	m.ReviewerId = reviewerId
	m.ReviewerName = reviewerName
	m.ReviewNote = note
}

// Reject 审核驳回
func (m *ModelEntity) Reject(reviewerId uint64, reviewerName, note string) error {
	if note == "" {
		return errors.New("驳回时必须填写审核意见")
	}
	m.Status = StatusRejected
	m.ReviewerId = reviewerId
	m.ReviewerName = reviewerName
	m.ReviewNote = note
	return nil
}

// ==================== 统计相关能力 ====================

// IncrementUsageCount 增加使用次数
func (m *ModelEntity) IncrementUsageCount() {
	m.UsageCount++
}

// IncrementAdoptCount 增加采纳次数
func (m *ModelEntity) IncrementAdoptCount() {
	m.AdoptCount++
}

// IncrementLikeCount 增加点赞数
func (m *ModelEntity) IncrementLikeCount() {
	m.LikeCount++
}

// IncrementCommentCount 增加评论数
func (m *ModelEntity) IncrementCommentCount() {
	m.CommentCount++
}

// CountByStatus 按状态统计数量
func (m *ModelEntity) CountByStatus(status int) (int64, error) {
	var count int64
	err := m.Db.Model(&ModelEntity{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

// ==================== 派生相关能力 ====================

// Fork 派生模型
func (m *ModelEntity) Fork() ModelEntityInterface {
	newEntity := &ModelEntity{
		Name:          m.Name + " (派生)",
		Description:   m.Description,
		Icon:          m.Icon,
		CategoryId:    m.CategoryId,
		Content:       m.Content,
		UsageGuide:    m.UsageGuide,
		Examples:      m.Examples,
		AiPrompt:      m.AiPrompt,
		Difficulty:    m.Difficulty,
		EstimatedTime: m.EstimatedTime,
		Status:        0, // 草稿
		Version:       "1.0.0",
		IsOfficial:    false,
		SourceModelId: m.Id,
	}
	newEntity.BaseModel = m.BaseModel
	return newEntity
}
