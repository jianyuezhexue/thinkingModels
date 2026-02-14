package master

import (
	"github.com/gin-gonic/gin"
	"thinkingModels/domain/master/category"
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

// All 获取全量分类列表（按热度降序）
// 用于思维模型列表页的筛选条件
func (l *CategoryLogic) All() (*category.AllReap, error) {
	entity := category.NewCategoryEntity(l.Ctx)
	list, err := entity.All()
	if err != nil {
		return nil, err
	}
	return &category.AllReap{List: list}, nil
}

// Create 创建分类
func (l *CategoryLogic) Create(req *category.CreateCategory) (*category.CategoryEntity, error) {
	entity := category.NewCategoryEntity(l.Ctx)
	return entity.CreateCategory(req)
}

// Update 更新分类
func (l *CategoryLogic) Update(req *category.UpdateCategory) (*category.CategoryEntity, error) {
	// 实例化模型
	entity := category.NewCategoryEntity(l.Ctx)

	// 加载旧数据
	_, err := entity.LoadById(req.Id)
	if err != nil {
		return nil, err
	}

	// 设置新数据
	_, err = entity.SetData(req)
	if err != nil {
		return nil, err
	}

	// 数据校验
	err = entity.Validate()
	if err != nil {
		return nil, err
	}

	// 更新数据
	res, err := entity.Update()
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Get 查询分类详情
func (l *CategoryLogic) Get(id uint64) (*category.CategoryEntity, error) {
	entity := category.NewCategoryEntity(l.Ctx)
	res, err := entity.LoadById(id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// List 查询分类列表（分页）
func (l *CategoryLogic) List(req *category.SearchCategory) (*category.ListReap, error) {
	entity := category.NewCategoryEntity(l.Ctx)

	// 构造查询条件
	cond := entity.MakeConditon(*req)

	// 统计总数
	total, err := entity.Count(cond)
	if err != nil {
		return nil, err
	}

	// 查询列表
	list, err := entity.List(cond)
	if err != nil {
		return nil, err
	}

	return &category.ListReap{
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
		List:     list,
	}, nil
}

// Del 删除分类
func (l *CategoryLogic) Del(req *category.DelCategory) (any, error) {
	entity := category.NewCategoryEntity(l.Ctx)
	err := entity.Del(req.Ids...)
	return nil, err
}

// IncreaseHeat 增加分类热度
func (l *CategoryLogic) IncreaseHeat(req *category.IncreaseHeatRequest) (*category.CategoryEntity, error) {
	entity := category.NewCategoryEntity(l.Ctx)

	// 增加热度
	err := entity.IncreaseHeat(req.Id, req.Delta)
	if err != nil {
		return nil, err
	}

	// 返回更新后的数据
	res, err := entity.LoadById(req.Id)
	if err != nil {
		return nil, err
	}

	return res, nil
}
