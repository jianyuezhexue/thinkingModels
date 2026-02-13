package followup

import (
	"errors"

	"thinkingModels/component/db"

	"github.com/gin-gonic/gin"
	"github.com/jianyuezhexue/base"
)

// FollowUpEntityInterface 跟进记录实体接口
type FollowUpEntityInterface interface {
	base.BaseModelInterface[FollowUpEntity]
	SetProgressChange(before, after int)
}

// FollowUpEntity 跟进记录实体
type FollowUpEntity struct {
	base.BaseModel[FollowUpEntity]
	UserId         uint64 `json:"userId" type:"db" comment:"所属用户ID"`
	ActionId       uint64 `json:"actionId" type:"db" comment:"关联行动ID"`
	Content        string `json:"content" type:"db" comment:"跟进内容"`
	ProgressBefore int    `json:"progressBefore" type:"db" comment:"跟进前进度"`
	ProgressAfter  int    `json:"progressAfter" type:"db" comment:"跟进后进度"`
}

// NewFollowUpEntity 实例化跟进记录实体
func NewFollowUpEntity(ctx *gin.Context, opt ...base.Option[FollowUpEntity]) FollowUpEntityInterface {
	entity := &FollowUpEntity{}
	entity.BaseModel = base.NewBaseModel(ctx, db.InitDb(), entity.TableName(), entity)
	if len(opt) > 0 {
		for _, fc := range opt {
			fc(&entity.BaseModel)
		}
	}
	return entity
}

// TableName 数据表名
func (f *FollowUpEntity) TableName() string {
	return "action_followups"
}

// Validate 数据校验
func (f *FollowUpEntity) Validate() error {
	if f.ActionId == 0 {
		return errors.New("行动ID不能为空")
	}
	if f.UserId == 0 {
		return errors.New("用户ID不能为空")
	}
	if f.Content == "" {
		return errors.New("跟进内容不能为空")
	}
	if f.ProgressBefore < 0 || f.ProgressBefore > 100 {
		return errors.New("跟进前进度必须在0-100之间")
	}
	if f.ProgressAfter < 0 || f.ProgressAfter > 100 {
		return errors.New("跟进后进度必须在0-100之间")
	}
	return nil
}

// Repair 数据修复
func (f *FollowUpEntity) Repair() error {
	return nil
}

// Complete 数据完善
func (f *FollowUpEntity) Complete() error {
	return nil
}

// SetProgressChange 设置进度变化
func (f *FollowUpEntity) SetProgressChange(before, after int) {
	f.ProgressBefore = before
	f.ProgressAfter = after
}
