package iam

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"thinkingModels/api"
	"thinkingModels/domain/iam/role"
	"thinkingModels/logic/iam"
)

type Role struct {
	api.Base
}

func NewRole() *Role {
	return &Role{}
}

// Create 创建角色
// @Summary 创建角色
// @Description 创建新的角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body role.CreateRole true "创建角色请求"
// @Success 200 {object} api.Response{data=role.RoleInfo} "创建成功"
// @Failure 400 {object} api.Response "参数错误"
// @Failure 409 {object} api.Response "角色编码已存在"
// @Router /role [post]
func (a Role) Create(ctx *gin.Context) {
	// 参数校验
	req := &role.CreateRole{}
	err := a.Bind(ctx, req)
	if err != nil {
		a.Error(err)
		return
	}

	// 实例化逻辑层
	logic := iam.NewRoleLogic(ctx)
	res, err := logic.Create(req)
	if err != nil {
		a.Error(err)
		return
	}

	// 接口返回
	a.Success(res, "创建成功")
}

// Update 更新角色
// @Summary 更新角色
// @Description 更新角色信息
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body role.UpdateRole true "更新角色请求"
// @Success 200 {object} api.Response{data=role.RoleInfo} "更新成功"
// @Failure 400 {object} api.Response "参数错误"
// @Failure 404 {object} api.Response "角色不存在"
// @Router /role [put]
func (a Role) Update(ctx *gin.Context) {
	// 参数校验
	req := &role.UpdateRole{}
	err := a.Bind(ctx, req)
	if err != nil {
		a.Error(err)
		return
	}

	// 实例化逻辑层
	logic := iam.NewRoleLogic(ctx)
	res, err := logic.Update(req)
	if err != nil {
		a.Error(err)
		return
	}

	// 接口返回
	a.Success(res, "更新成功")
}

// Get 查询角色详情
// @Summary 查询角色详情
// @Description 根据ID查询角色详情
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "角色ID"
// @Success 200 {object} api.Response{data=role.RoleInfo} "查询成功"
// @Failure 404 {object} api.Response "角色不存在"
// @Router /role/{id} [get]
func (a Role) Get(ctx *gin.Context) {
	// 从路径参数获取id
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		a.Error(err)
		return
	}

	// 实例化逻辑层
	logic := iam.NewRoleLogic(ctx)
	res, err := logic.Get(id)
	if err != nil {
		a.Error(err)
		return
	}

	// 接口返回
	a.Success(res, "查询成功")
}

// List 查询角色列表
// @Summary 查询角色列表
// @Description 分页查询角色列表，支持按角色名称模糊查询
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码，默认1"
// @Param pageSize query int false "每页数量，默认10"
// @Param roleName query string false "角色名称模糊查询"
// @Param status query int false "状态：0-禁用，1-启用"
// @Success 200 {object} api.Response{data=role.ListRoleResponse} "查询成功"
// @Router /role/list [post]
func (a Role) List(ctx *gin.Context) {
	// 参数校验
	req := &role.SearchRole{}
	err := a.Bind(ctx, req)
	if err != nil {
		a.Error(err)
		return
	}

	// 设置默认分页
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 实例化逻辑层
	logic := iam.NewRoleLogic(ctx)
	res, err := logic.List(req)
	if err != nil {
		a.Error(err)
		return
	}

	// 接口返回
	a.Success(res, "查询列表成功")
}

// Del 删除角色
// @Summary 删除角色
// @Description 批量删除角色，如果角色下有用户则无法删除
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body role.DelRole true "删除角色请求"
// @Success 200 {object} api.Response "删除成功"
// @Failure 400 {object} api.Response "参数错误"
// @Failure 409 {object} api.Response "角色下存在用户，无法删除"
// @Router /role [delete]
func (a Role) Del(ctx *gin.Context) {
	req := &role.DelRole{}
	err := a.Bind(ctx, req)
	if err != nil {
		a.Error(err)
		return
	}

	// 实例化逻辑层
	logic := iam.NewRoleLogic(ctx)
	res, err := logic.Del(req.Ids)
	if err != nil {
		a.Error(err)
		return
	}

	// 接口返回
	a.Success(res, "删除成功")
}

// All 获取所有角色
// @Summary 获取所有角色
// @Description 获取所有启用的角色列表，用于下拉选择
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} api.Response{data=[]role.RoleInfo} "查询成功"
// @Router /role/all [get]
func (a Role) All(ctx *gin.Context) {
	// 必须调用 Bind 以设置上下文
	req := &struct{}{}
	if err := a.Bind(ctx, req); err != nil {
		a.Error(err)
		return
	}

	// 实例化逻辑层
	logic := iam.NewRoleLogic(ctx)
	res, err := logic.All()
	if err != nil {
		a.Error(err)
		return
	}

	// 接口返回
	a.Success(res, "查询成功")
}

// UpdatePermission 更新角色权限
// @Summary 更新角色权限
// @Description 更新角色的菜单权限
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body role.UpdateRolePermission true "更新权限请求"
// @Success 200 {object} api.Response{data=role.RoleInfo} "更新成功"
// @Failure 400 {object} api.Response "参数错误"
// @Failure 404 {object} api.Response "角色不存在"
// @Router /role/permission [put]
func (a Role) UpdatePermission(ctx *gin.Context) {
	// 参数校验
	req := &role.UpdateRolePermission{}
	err := a.Bind(ctx, req)
	if err != nil {
		a.Error(err)
		return
	}

	// 实例化逻辑层
	logic := iam.NewRoleLogic(ctx)
	res, err := logic.UpdatePermission(req)
	if err != nil {
		a.Error(err)
		return
	}

	// 接口返回
	a.Success(res, "权限更新成功")
}