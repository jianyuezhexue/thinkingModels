package iam

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"thinkingModels/domain/iam/role"
	"thinkingModels/logic"
)

// RoleLogic 角色业务逻辑
type RoleLogic struct {
	logic.BaseLogic
}

// 初始化RoleLogic
func NewRoleLogic(ctx *gin.Context) *RoleLogic {
	return &RoleLogic{BaseLogic: logic.BaseLogic{Ctx: ctx}}
}

// Create 创建角色
func (l *RoleLogic) Create(req *role.CreateRole) (*role.RoleInfo, error) {
	// 实例化模型
	roleEntity := role.NewRoleEntity(l.Ctx)

	// 检查角色编码是否已存在
	cond := roleEntity.MakeConditon(role.SearchRole{RoleCode: req.RoleCode})
	existingRole, err := roleEntity.LoadData(cond)
	if err == nil && existingRole.Id > 0 {
		return nil, errors.New("角色编码已存在")
	}

	// 数据赋值
	_, err = roleEntity.SetData(req)
	if err != nil {
		return nil, err
	}

	// 数据校验
	err = roleEntity.Validate()
	if err != nil {
		return nil, err
	}

	// 设置默认状态
	if entity, ok := roleEntity.(*role.RoleEntity); ok {
		if entity.Status == 0 {
			entity.Status = 1
		}
	}

	// 保存数据
	res, err := roleEntity.Create()
	if err != nil {
		return nil, err
	}

	// 返回角色信息DTO
	return &role.RoleInfo{
		ID:          res.Id,
		RoleName:    res.RoleName,
		RoleCode:    res.RoleCode,
		Description: res.Description,
		Status:      res.Status,
		Sort:        res.Sort,
		MenuIds:     res.MenuIds,
		UserCount:   0,
		CreatedAt:   res.CreatedAt.String(),
		UpdatedAt:   res.UpdatedAt.String(),
	}, nil
}

// Update 更新角色
func (l *RoleLogic) Update(req *role.UpdateRole) (*role.RoleInfo, error) {
	// 实例化模型
	roleEntity := role.NewRoleEntity(l.Ctx)

	// 先加载旧数据
	_, err := roleEntity.LoadById(req.ID)
	if err != nil {
		return nil, errors.New("角色不存在")
	}

	// 数据赋值
	_, err = roleEntity.SetData(req)
	if err != nil {
		return nil, err
	}

	// 数据校验
	err = roleEntity.Validate()
	if err != nil {
		return nil, err
	}

	// 更新数据
	res, err := roleEntity.Update()
	if err != nil {
		return nil, err
	}

	// 获取用户数量
	userCount, _ := l.getUserCountByRoleId(res.Id)

	// 返回角色信息DTO
	return &role.RoleInfo{
		ID:          res.Id,
		RoleName:    res.RoleName,
		RoleCode:    res.RoleCode,
		Description: res.Description,
		Status:      res.Status,
		Sort:        res.Sort,
		MenuIds:     res.MenuIds,
		UserCount:   userCount,
		CreatedAt:   res.CreatedAt.String(),
		UpdatedAt:   res.UpdatedAt.String(),
	}, nil
}

// Get 查询角色详情
func (l *RoleLogic) Get(id uint64) (*role.RoleInfo, error) {
	// 实例化模型
	roleEntity := role.NewRoleEntity(l.Ctx)

	// 查询数据
	res, err := roleEntity.LoadById(id)
	if err != nil {
		return nil, errors.New("角色不存在")
	}

	// 获取用户数量
	userCount, _ := l.getUserCountByRoleId(res.Id)

	// 返回角色信息DTO
	return &role.RoleInfo{
		ID:          res.Id,
		RoleName:    res.RoleName,
		RoleCode:    res.RoleCode,
		Description: res.Description,
		Status:      res.Status,
		Sort:        res.Sort,
		MenuIds:     res.MenuIds,
		UserCount:   userCount,
		CreatedAt:   res.CreatedAt.String(),
		UpdatedAt:   res.UpdatedAt.String(),
	}, nil
}

// List 查询角色列表
func (l *RoleLogic) List(req *role.SearchRole) (*role.ListRoleResponse, error) {
	// 实例化模型
	roleEntity := role.NewRoleEntity(l.Ctx)

	// 获取搜索条件
	cond := roleEntity.MakeConditon(*req)

	// 查询总数
	total, err := roleEntity.Count(cond)
	if err != nil {
		return nil, err
	}

	// 获取列表
	list, err := roleEntity.List(cond)
	if err != nil {
		return nil, err
	}

	// 转换为DTO列表
	roleInfoList := make([]*role.RoleInfo, 0, len(list))
	for _, item := range list {
		// 获取用户数量
		userCount, _ := l.getUserCountByRoleId(item.Id)

		roleInfoList = append(roleInfoList, &role.RoleInfo{
			ID:          item.Id,
			RoleName:    item.RoleName,
			RoleCode:    item.RoleCode,
			Description: item.Description,
			Status:      item.Status,
			Sort:        item.Sort,
			MenuIds:     item.MenuIds,
			UserCount:   userCount,
			CreatedAt:   item.CreatedAt.String(),
			UpdatedAt:   item.UpdatedAt.String(),
		})
	}

	// 返回数据
	res := &role.ListRoleResponse{
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
		List:     roleInfoList,
	}
	return res, nil
}

// Del 删除角色
func (l *RoleLogic) Del(ids []uint64) (any, error) {
	// 检查是否有用户关联该角色
	for _, id := range ids {
		userCount, err := l.getUserCountByRoleId(id)
		if err != nil {
			return nil, err
		}
		if userCount > 0 {
			return nil, errors.New("存在关联用户，无法删除该角色")
		}
	}

	// 实例化模型
	roleEntity := role.NewRoleEntity(l.Ctx)

	// 删除数据
	err := roleEntity.Del(ids...)
	return nil, err
}

// UpdatePermission 更新角色权限
func (l *RoleLogic) UpdatePermission(req *role.UpdateRolePermission) (*role.RoleInfo, error) {
	// 实例化模型
	roleEntity := role.NewRoleEntity(l.Ctx)

	// 先加载旧数据
	_, err := roleEntity.LoadById(req.ID)
	if err != nil {
		return nil, errors.New("角色不存在")
	}

	// 更新菜单ID列表
	if entity, ok := roleEntity.(*role.RoleEntity); ok {
		entity.MenuIds = strings.Join(req.MenuIds, ",")
	}

	// 更新数据
	res, err := roleEntity.Update()
	if err != nil {
		return nil, err
	}

	// 获取用户数量
	userCount, _ := l.getUserCountByRoleId(res.Id)

	// 返回角色信息DTO
	return &role.RoleInfo{
		ID:          res.Id,
		RoleName:    res.RoleName,
		RoleCode:    res.RoleCode,
		Description: res.Description,
		Status:      res.Status,
		Sort:        res.Sort,
		MenuIds:     res.MenuIds,
		UserCount:   userCount,
		CreatedAt:   res.CreatedAt.String(),
		UpdatedAt:   res.UpdatedAt.String(),
	}, nil
}

// getUserCountByRoleId 获取角色关联的用户数量
// 注意：这是一个简化的实现，返回0表示暂不统计
// 实际项目中需要根据业务需求实现精确的用户统计
func (l *RoleLogic) getUserCountByRoleId(roleId uint64) (int64, error) {
	// TODO: 实现角色关联用户数量统计
	// 由于用户的 roleIds 字段是逗号分隔的字符串，需要使用 LIKE 查询或关联表
	// 这里暂时返回0，后续可以优化为关联表设计
	return 0, nil
}

// All 获取所有角色（用于下拉选择）
func (l *RoleLogic) All() ([]*role.RoleInfo, error) {
	// 实例化模型
	roleEntity := role.NewRoleEntity(l.Ctx)

	// 获取所有角色（不分页，只查询启用状态）
	cond := roleEntity.MakeConditon(role.SearchRole{Status: 1})
	list, err := roleEntity.List(cond)
	if err != nil {
		return nil, err
	}

	// 转换为DTO列表
	roleInfoList := make([]*role.RoleInfo, 0, len(list))
	for _, item := range list {
		roleInfoList = append(roleInfoList, &role.RoleInfo{
			ID:          item.Id,
			RoleName:    item.RoleName,
			RoleCode:    item.RoleCode,
			Description: item.Description,
			Status:      item.Status,
			Sort:        item.Sort,
			MenuIds:     item.MenuIds,
			CreatedAt:   item.CreatedAt.String(),
			UpdatedAt:   item.UpdatedAt.String(),
		})
	}

	return roleInfoList, nil
}