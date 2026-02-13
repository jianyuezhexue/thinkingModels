package thinking

import (
	"thinkingModels/domain/thinking/followup"
	"thinkingModels/logic"

	"github.com/gin-gonic/gin"
)

// FollowUpLogic 跟进记录业务逻辑
type FollowUpLogic struct {
	logic.BaseLogic
}

// NewFollowUpLogic 初始化FollowUpLogic
func NewFollowUpLogic(ctx *gin.Context) *FollowUpLogic {
	return &FollowUpLogic{BaseLogic: logic.BaseLogic{Ctx: ctx}}
}

// Create 创建跟进记录
func (l *FollowUpLogic) Create(req *followup.CreateFollowUp) (*followup.FollowUpInfo, error) {
	entity := followup.NewFollowUpEntity(l.Ctx)
	_, err := entity.SetData(req)
	if err != nil {
		return nil, err
	}

	if concreteEntity, ok := entity.(*followup.FollowUpEntity); ok {
		concreteEntity.SetProgressChange(req.ProgressBefore, req.ProgressAfter)
	}

	if err := entity.Validate(); err != nil {
		return nil, err
	}

	res, err := entity.Create()
	if err != nil {
		return nil, err
	}

	return convertToFollowUpInfo(res), nil
}

// Update 更新跟进记录
func (l *FollowUpLogic) Update(req *followup.UpdateFollowUp) (*followup.FollowUpInfo, error) {
	entity := followup.NewFollowUpEntity(l.Ctx)
	_, err := entity.LoadById(req.Id)
	if err != nil {
		return nil, err
	}

	_, err = entity.SetData(req)
	if err != nil {
		return nil, err
	}

	if err := entity.Validate(); err != nil {
		return nil, err
	}

	res, err := entity.Update()
	if err != nil {
		return nil, err
	}

	return convertToFollowUpInfo(res), nil
}

// Get 查询跟进记录详情
func (l *FollowUpLogic) Get(id uint64) (*followup.FollowUpInfo, error) {
	entity := followup.NewFollowUpEntity(l.Ctx)
	res, err := entity.LoadById(id)
	if err != nil {
		return nil, err
	}
	return convertToFollowUpInfo(res), nil
}

// Del 删除跟进记录
func (l *FollowUpLogic) Del(req *followup.DelFollowUp) error {
	entity := followup.NewFollowUpEntity(l.Ctx)
	return entity.Del(req.Ids...)
}

// ListByAction 按行动项查询跟进记录
func (l *FollowUpLogic) ListByAction(actionId uint64) ([]*followup.FollowUpInfo, error) {
	entity := followup.NewFollowUpEntity(l.Ctx)
	cond := entity.MakeConditon(followup.SearchFollowUp{ActionId: actionId})

	list, err := entity.List(cond)
	if err != nil {
		return nil, err
	}

	result := make([]*followup.FollowUpInfo, 0, len(list))
	for _, item := range list {
		result = append(result, convertToFollowUpInfo(item))
	}

	return result, nil
}

func convertToFollowUpInfo(entity any) *followup.FollowUpInfo {
	e, ok := entity.(*followup.FollowUpEntity)
	if !ok {
		return nil
	}
	return &followup.FollowUpInfo{
		Id:             e.Id,
		ActionId:       e.ActionId,
		Content:        e.Content,
		ProgressBefore: e.ProgressBefore,
		ProgressAfter:  e.ProgressAfter,
		ProgressDelta:  e.ProgressAfter - e.ProgressBefore,
		CreatedAt:      e.CreatedAt.String(),
		UpdatedAt:      e.UpdatedAt.String(),
	}
}
