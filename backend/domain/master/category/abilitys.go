package category

import (
	"errors"
)

// CategoryAbility 分类能力接口定义
type CategoryAbility interface {
	// All 获取全量分类列表（按热度降序）
	All() ([]*CategoryEntity, error)
	// CreateCategory 新建分类（自动校验名称唯一性）
	CreateCategory(req *CreateCategory) (*CategoryEntity, error)
	// IncreaseHeat 增加分类热度值
	IncreaseHeat(id uint64, delta int) error
}

// All 获取全量分类列表（按热度降序）
// 用于思维模型列表页的筛选条件展示
func (m *CategoryEntity) All() ([]*CategoryEntity, error) {
	var list []*CategoryEntity
	err := m.Tx().Model(m).
		Where("deleted_at IS NULL").
		Order("heat DESC, id ASC").
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

// CreateCategory 新建分类
// 自动校验名称唯一性，初始化热度为0
func (m *CategoryEntity) CreateCategory(req *CreateCategory) (*CategoryEntity, error) {
	// 校验名称必填
	if req.Name == "" {
		return nil, errors.New("分类名称不能为空")
	}

	// 校验名称长度
	if len(req.Name) > 50 {
		return nil, errors.New("分类名称长度不能超过50个字符")
	}

	// 校验名称唯一性
	exists, err := m.CheckBusinessCodeExist("name", req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("分类名称已存在")
	}

	// 初始化热度为0
	heat := req.Heat
	if heat == 0 {
		heat = 0
	}

	// 创建实体
	entity := &CategoryEntity{
		Name:        req.Name,
		Icon:        req.Icon,
		Description: req.Description,
		Heat:        heat,
	}
	entity.BaseModel = m.BaseModel

	// 数据校验
	err = entity.Validate()
	if err != nil {
		return nil, err
	}

	// 完善数据
	err = entity.Complete()
	if err != nil {
		return nil, err
	}

	// 保存数据
	res, err := entity.Create()
	if err != nil {
		return nil, err
	}

	return res, nil
}

// IncreaseHeat 增加分类热度值
// 用于热门分类排序，delta 为增加的热度值
func (m *CategoryEntity) IncreaseHeat(id uint64, delta int) error {
	if delta <= 0 {
		return errors.New("热度增加值必须大于0")
	}

	// 加载实体
	entity, err := m.LoadById(id)
	if err != nil {
		return err
	}

	// 更新热度
	entity.Heat += delta

	// 保存更新
	_, err = entity.Update()
	return err
}

// InitCategoryAbility 初始化分类能力
// 在需要时使用能力模式初始化
func InitCategoryAbility(m *CategoryEntity) CategoryAbility {
	return m
}
