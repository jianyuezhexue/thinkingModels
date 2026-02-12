package market

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"thinkingModels/domain/market/category"
	"thinkingModels/logic"
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
	// 实例化模型
	entity := category.NewCategoryEntity(l.Ctx)

	// 设置数据
	if concreteEntity, ok := entity.(*category.CategoryEntity); ok {
		concreteEntity.ParentId = req.ParentId
		concreteEntity.Name = req.Name
		concreteEntity.Code = req.Code
		concreteEntity.Description = req.Description
		concreteEntity.Icon = req.Icon
		concreteEntity.SortOrder = req.SortOrder
		concreteEntity.Status = 1 // 默认启用

		// 计算层级
		if req.ParentId == 0 {
			concreteEntity.Level = 1
			concreteEntity.Path = "/"
		} else {
			// 查询父分类
			parent, err := entity.LoadById(req.ParentId)
			if err != nil {
				return nil, errors.New("父分类不存在")
			}
			concreteEntity.Level = parent.Level + 1
			concreteEntity.Path = fmt.Sprintf("%s%d/", parent.Path, parent.Id)
		}
	}

	// 数据校验
	if err := entity.Validate(); err != nil {
		return nil, err
	}

	// 保存数据（此时还没有ID）
	res, err := entity.Create()
	if err != nil {
		return nil, err
	}

	// 更新路径（包含自己的ID）
	if concreteEntity, ok := entity.(*category.CategoryEntity); ok {
		if concreteEntity.ParentId == 0 {
			concreteEntity.Path = fmt.Sprintf("/%d/", res.Id)
		} else {
			concreteEntity.Path = fmt.Sprintf("%s%d/", concreteEntity.Path, res.Id)
		}
		_, err = entity.Update()
		if err != nil {
			return nil, err
		}
	}

	// 转换为DTO返回
	return convertToCategoryInfo(res), nil
}

// Update 更新分类
func (l *CategoryLogic) Update(req *category.UpdateCategory) (*category.CategoryInfo, error) {
	// 实例化模型
	entity := category.NewCategoryEntity(l.Ctx)

	// 先加载旧数据
	oldData, err := entity.LoadById(req.Id)
	if err != nil {
		return nil, err
	}

	// 设置数据
	if concreteEntity, ok := entity.(*category.CategoryEntity); ok {
		concreteEntity.ParentId = oldData.ParentId // 保持父级不变
		concreteEntity.Name = req.Name
		concreteEntity.Description = req.Description
		concreteEntity.Icon = req.Icon
		concreteEntity.SortOrder = req.SortOrder
		concreteEntity.Level = oldData.Level
		concreteEntity.Path = oldData.Path
		concreteEntity.Status = oldData.Status
		concreteEntity.Code = oldData.Code // 编码不允许修改
	}

	// 数据校验
	if err := entity.Validate(); err != nil {
		return nil, err
	}

	// 更新数据
	res, err := entity.Update()
	if err != nil {
		return nil, err
	}

	// 转换为DTO返回
	return convertToCategoryInfo(res), nil
}

// Move 移动分类（修改父级）
func (l *CategoryLogic) Move(req *category.MoveCategory) (*category.CategoryInfo, error) {
	// 实例化模型
	entity := category.NewCategoryEntity(l.Ctx)

	// 加载当前分类
	_, err := entity.LoadById(req.Id)
	if err != nil {
		return nil, err
	}

	// 检查目标父分类是否存在（如果parentId不为0）
	if req.ParentId != 0 {
		_, err = entity.LoadById(req.ParentId)
		if err != nil {
			return nil, errors.New("目标父分类不存在")
		}
	}

	// 检查是否移动到自己下面（循环引用检查）
	if req.ParentId == req.Id {
		return nil, errors.New("不能将分类移动到自己下面")
	}

	// 检查是否会形成循环引用
	if req.ParentId != 0 {
		children, err := l.GetAllChildren(req.Id)
		if err != nil {
			return nil, err
		}
		for _, child := range children {
			if child.Id == req.ParentId {
				return nil, errors.New("不能将分类移动到其子分类下，会形成循环引用")
			}
		}
	}

	// 更新父级和层级
	if concreteEntity, ok := entity.(*category.CategoryEntity); ok {
		concreteEntity.ParentId = req.ParentId

		// 重新计算层级和路径
		if req.ParentId == 0 {
			concreteEntity.Level = 1
			concreteEntity.Path = fmt.Sprintf("/%d/", req.Id)
		} else {
			parent, _ := entity.LoadById(req.ParentId)
			if parent != nil {
				concreteEntity.Level = parent.Level + 1
				concreteEntity.Path = fmt.Sprintf("%s%d/", parent.Path, req.Id)
			}
		}
	}

	// 更新数据
	res, err := entity.Update()
	if err != nil {
		return nil, err
	}

	// 递归更新所有子分类的路径
	if err := l.updateChildrenPath(req.Id, res.Path, res.Level); err != nil {
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
func (l *CategoryLogic) List(req *category.SearchCategory) (*category.ListCategoryResponse, error) {
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

	infoList := make([]*category.CategoryInfo, 0, len(list))
	for _, item := range list {
		infoList = append(infoList, convertToCategoryInfo(item))
	}

	return &category.ListCategoryResponse{
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
		List:     infoList,
	}, nil
}

// Tree 获取分类树
func (l *CategoryLogic) Tree() (*category.CategoryTreeResponse, error) {
	// 查询所有分类
	entity := category.NewCategoryEntity(l.Ctx)
	search := &category.SearchCategory{Status: 1, PageSize: 1000}
	cond := entity.MakeConditon(*search)

	list, err := entity.List(cond)
	if err != nil {
		return nil, err
	}

	// 构建树
	tree := buildCategoryTree(list)

	return &category.CategoryTreeResponse{List: tree}, nil
}

// TreeWithRoot 获取分类树（包含根节点）
func (l *CategoryLogic) TreeWithRoot() (*category.CategoryTreeResponse, error) {
	// 查询所有分类
	entity := category.NewCategoryEntity(l.Ctx)
	search := &category.SearchCategory{Status: -1, PageSize: 1000} // -1表示所有状态
	cond := entity.MakeConditon(*search)

	list, err := entity.List(cond)
	if err != nil {
		return nil, err
	}

	// 构建树
	tree := buildCategoryTree(list)

	return &category.CategoryTreeResponse{List: tree}, nil
}

// GetChildren 获取子分类列表
func (l *CategoryLogic) GetChildren(parentId uint64) ([]*category.CategoryInfo, error) {
	entity := category.NewCategoryEntity(l.Ctx)
	search := &category.SearchCategory{ParentId: parentId, Status: 1}
	cond := entity.MakeConditon(*search)

	list, err := entity.List(cond)
	if err != nil {
		return nil, err
	}

	infoList := make([]*category.CategoryInfo, 0, len(list))
	for _, item := range list {
		infoList = append(infoList, convertToCategoryInfo(item))
	}

	return infoList, nil
}

// GetAllChildren 递归获取所有子分类
func (l *CategoryLogic) GetAllChildren(parentId uint64) ([]*category.CategoryEntity, error) {
	entity := category.NewCategoryEntity(l.Ctx)
	search := &category.SearchCategory{Status: -1, PageSize: 1000}
	cond := entity.MakeConditon(*search)

	allList, err := entity.List(cond)
	if err != nil {
		return nil, err
	}

	result := make([]*category.CategoryEntity, 0)
	findChildren(allList, parentId, &result)

	return result, nil
}

// Path 获取分类完整路径
func (l *CategoryLogic) Path(id uint64) (*category.CategoryPathResponse, error) {
	entity := category.NewCategoryEntity(l.Ctx)
	current, err := entity.LoadById(id)
	if err != nil {
		return nil, err
	}

	// 使用实体的GetFullPath方法
	path, err := current.GetFullPath()
	if err != nil {
		return nil, err
	}

	pathInfo := make([]*category.CategoryInfo, 0, len(path))
	for _, item := range path {
		pathInfo = append(pathInfo, convertToCategoryInfo(item))
	}

	return &category.CategoryPathResponse{Path: pathInfo}, nil
}

// Del 删除分类
func (l *CategoryLogic) Del(ids []uint64) (any, error) {
	// 检查是否有子分类
	for _, id := range ids {
		children, err := l.GetChildren(id)
		if err != nil {
			return nil, err
		}
		if len(children) > 0 {
			return nil, errors.New("存在子分类，无法删除")
		}
	}

	entity := category.NewCategoryEntity(l.Ctx)
	err := entity.Del(ids...)
	return nil, err
}

// UpdateStatus 更新分类状态
func (l *CategoryLogic) UpdateStatus(req *category.UpdateCategoryStatus) (*category.CategoryInfo, error) {
	entity := category.NewCategoryEntity(l.Ctx)

	_, err := entity.LoadById(req.Id)
	if err != nil {
		return nil, err
	}

	if concreteEntity, ok := entity.(*category.CategoryEntity); ok {
		concreteEntity.Status = req.Status
	}

	res, err := entity.Update()
	if err != nil {
		return nil, err
	}

	return convertToCategoryInfo(res), nil
}

// 递归更新子分类的路径
func (l *CategoryLogic) updateChildrenPath(parentId uint64, parentPath string, parentLevel int) error {
	entity := category.NewCategoryEntity(l.Ctx)
	search := &category.SearchCategory{ParentId: parentId, Status: -1, PageSize: 1000}
	cond := entity.MakeConditon(*search)

	children, err := entity.List(cond)
	if err != nil {
		return err
	}

	for _, child := range children {
		// 更新当前子分类
		newPath := fmt.Sprintf("%s%d/", parentPath, child.Id)
		newLevel := parentLevel + 1

		if concreteChild, ok := interface{}(child).(*category.CategoryEntity); ok {
			concreteChild.Path = newPath
			concreteChild.Level = newLevel
			_, err := concreteChild.Update()
			if err != nil {
				return err
			}
		}

		// 递归更新孙子分类
		if err := l.updateChildrenPath(child.Id, newPath, newLevel); err != nil {
			return err
		}
	}

	return nil
}

// 构建分类树
func buildCategoryTree(list []*category.CategoryEntity) []*category.CategoryTreeNode {
	// 创建节点映射
	nodeMap := make(map[uint64]*category.CategoryTreeNode)
	for _, item := range list {
		nodeMap[item.Id] = &category.CategoryTreeNode{
			CategoryInfo: *convertToCategoryInfo(item),
			Children:     make([]*category.CategoryTreeNode, 0),
		}
	}

	// 构建树结构
	var tree []*category.CategoryTreeNode
	for _, item := range list {
		node := nodeMap[item.Id]
		if item.ParentId == 0 {
			// 根节点
			tree = append(tree, node)
		} else {
			// 子节点
			if parent, ok := nodeMap[item.ParentId]; ok {
				parent.Children = append(parent.Children, node)
			}
		}
	}

	return tree
}

// 递归查找所有子分类
func findChildren(list []*category.CategoryEntity, parentId uint64, result *[]*category.CategoryEntity) {
	for _, item := range list {
		if item.ParentId == parentId {
			*result = append(*result, item)
			findChildren(list, item.Id, result)
		}
	}
}

// 辅助函数：转换为CategoryInfo
func convertToCategoryInfo(entity *category.CategoryEntity) *category.CategoryInfo {
	if entity == nil {
		return nil
	}

	statusText := "禁用"
	if entity.Status == 1 {
		statusText = "启用"
	}

	return &category.CategoryInfo{
		Id:          entity.Id,
		ParentId:    entity.ParentId,
		Name:        entity.Name,
		Code:        entity.Code,
		Description: entity.Description,
		Icon:        entity.Icon,
		SortOrder:   entity.SortOrder,
		Level:       entity.Level,
		Path:        entity.Path,
		Status:      entity.Status,
		StatusText:  statusText,
		CreatedAt:   entity.CreatedAt.String(),
		UpdatedAt:   entity.UpdatedAt.String(),
	}
}
