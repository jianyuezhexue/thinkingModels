package thinking

import (
	"thinkingModels/domain/practice/category"
	"thinkingModels/logic"

	"github.com/gin-gonic/gin"
)

// CategoryLogic 分类业务逻辑
type CategoryLogic struct {
	logic.BaseLogic
}

// NewCategoryLogic 初始化CategoryLogic
func NewCategoryLogic(ctx *gin.Context) *CategoryLogic {
	return &CategoryLogic{BaseLogic: logic.BaseLogic{Ctx: ctx}}
}

// Create 创建分类
func (l *CategoryLogic) Create(req *category.CreateCategory) (*category.CategoryInfo, error) {
	entity := category.NewCategoryEntity(l.Ctx)
	_, err := entity.SetData(req)
	if err != nil {
		return nil, err
	}

	if err := entity.Validate(); err != nil {
		return nil, err
	}

	res, err := entity.Create()
	if err != nil {
		return nil, err
	}

	return convertToCategoryInfo(res), nil
}

// Update 更新分类
func (l *CategoryLogic) Update(req *category.UpdateCategory) (*category.CategoryInfo, error) {
	entity := category.NewCategoryEntity(l.Ctx)
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

	return convertToCategoryInfo(res), nil
}

// Get 查询分类详情
func (l *CategoryLogic) Get(id uint64) (*category.CategoryInfo, error) {
	entity := category.NewCategoryEntity(l.Ctx)
	res, err := entity.LoadById(id)
	if err != nil {
		return nil, err
	}
	return convertToCategoryInfo(res), nil
}

// List 查询分类列表
func (l *CategoryLogic) List(req *category.SearchCategory) (*logic.ListReap, error) {
	entity := category.NewCategoryEntity(l.Ctx)
	cond := entity.MakeConditon(*req)

	total, err := entity.Count(cond)
	if err != nil {
		return nil, err
	}

	list, err := entity.List(cond)
	if err != nil {
		return nil, err
	}

	categoryInfoList := make([]*category.CategoryInfo, 0, len(list))
	for _, item := range list {
		categoryInfoList = append(categoryInfoList, convertToCategoryInfo(item))
	}

	return &logic.ListReap{
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
		List:     categoryInfoList,
	}, nil
}

// Del 删除分类
func (l *CategoryLogic) Del(req *category.DelCategory) error {
	entity := category.NewCategoryEntity(l.Ctx)
	return entity.Del(req.Ids...)
}

// Tree 获取分类树
func (l *CategoryLogic) Tree() ([]*category.CategoryTree, error) {
	entity := category.NewCategoryEntity(l.Ctx)
	cond := entity.MakeConditon(category.SearchCategory{})

	list, err := entity.List(cond)
	if err != nil {
		return nil, err
	}

	return buildCategoryTree(list, 0), nil
}

func buildCategoryTree(list []*category.CategoryEntity, parentId uint64) []*category.CategoryTree {
	tree := make([]*category.CategoryTree, 0)
	for _, item := range list {
		if item.ParentId == parentId {
			node := &category.CategoryTree{
				Id:          item.Id,
				Code:        item.Code,
				Name:        item.Name,
				Icon:        item.Icon,
				Description: item.Description,
				ParentId:    item.ParentId,
				Sort:        item.Sort,
				Level:       item.Level,
				Path:        item.Path,
				ModelCount:  item.ModelCount,
				Status:      item.Status,
				Children:    buildCategoryTree(list, item.Id),
			}
			tree = append(tree, node)
		}
	}
	return tree
}

// GetChildren 获取子分类
func (l *CategoryLogic) GetChildren(parentId uint64) ([]*category.CategoryInfo, error) {
	entity := category.NewCategoryEntity(l.Ctx)
	cond := entity.MakeConditon(category.SearchCategory{ParentId: parentId})
	list, err := entity.List(cond)
	if err != nil {
		return nil, err
	}

	result := make([]*category.CategoryInfo, 0, len(list))
	for _, child := range list {
		result = append(result, convertToCategoryInfo(child))
	}
	return result, nil
}

// Move 移动分类
func (l *CategoryLogic) Move(req *category.MoveCategory) error {
	entity := category.NewCategoryEntity(l.Ctx)
	_, err := entity.LoadById(req.Id)
	if err != nil {
		return err
	}

	if concreteEntity, ok := entity.(*category.CategoryEntity); ok {
		concreteEntity.ParentId = req.TargetParentId
		concreteEntity.Sort = req.Sort
	}

	_, err = entity.Update()
	return err
}

// UpdateStatus 更新状态
func (l *CategoryLogic) UpdateStatus(req *category.UpdateCategoryStatus) error {
	entity := category.NewCategoryEntity(l.Ctx)
	_, err := entity.LoadById(req.Id)
	if err != nil {
		return err
	}

	if concreteEntity, ok := entity.(*category.CategoryEntity); ok {
		concreteEntity.Status = req.Status
	}

	_, err = entity.Update()
	return err
}

func convertToCategoryInfo(entity any) *category.CategoryInfo {
	e, ok := entity.(*category.CategoryEntity)
	if !ok {
		return nil
	}
	return &category.CategoryInfo{
		Id:          e.Id,
		Code:        e.Code,
		Name:        e.Name,
		Icon:        e.Icon,
		Description: e.Description,
		ParentId:    e.ParentId,
		Sort:        e.Sort,
		Level:       e.Level,
		Path:        e.Path,
		ModelCount:  e.ModelCount,
		Status:      e.Status,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}
