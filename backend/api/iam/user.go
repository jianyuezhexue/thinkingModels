package iam

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"thinkingModels/api"
	"thinkingModels/domain/iam/user"
	"thinkingModels/logic/iam"
)

type User struct {
	api.Base
}

func NewUser() *User {
	return &User{}
}

// Create 创建用户（注册）
// @Summary 用户注册
// @Description 用户注册接口，创建新用户账号
// @Tags 用户认证
// @Accept json
// @Produce json
// @Param request body user.RegisterRequest true "注册请求参数"
// @Success 200 {object} api.Response{data=user.UserInfo} "注册成功"
// @Failure 400 {object} api.Response "参数错误"
// @Failure 409 {object} api.Response "用户名已存在"
// @Router /auth/register [post]
func (a User) Create(ctx *gin.Context) {
	// 参数校验
	req := &user.RegisterRequest{}
	err := a.Bind(ctx, req)
	if err != nil {
		a.Error(err)
		return
	}

	// 实例化逻辑层
	logic := iam.NewUserLogic(ctx)
	res, err := logic.Create(req)
	if err != nil {
		a.Error(err)
		return
	}

	// 接口返回
	a.Success(res, "注册成功")
}

// Update 更新用户信息
func (a User) Update(ctx *gin.Context) {
	// 参数校验
	req := &user.UpdateUserRequest{}
	err := a.Bind(ctx, req)
	if err != nil {
		a.Error(err)
		return
	}

	// 实例化逻辑层
	logic := iam.NewUserLogic(ctx)
	res, err := logic.Update(req)
	if err != nil {
		a.Error(err)
		return
	}

	// 接口返回
	a.Success(res, "更新成功")
}

// Get 查询用户详情
func (a User) Get(ctx *gin.Context) {
	// TODO: 从路径参数获取id
	// idStr := ctx.Param("id")
	// 实例化逻辑层
	logic := iam.NewUserLogic(ctx)
	res, err := logic.Get(0) // TODO: 使用实际的id
	if err != nil {
		a.Error(err)
		return
	}

	// 接口返回
	a.Success(res, "查询成功")
}

// List 查询用户列表
// @Summary 查询用户列表
// @Description 分页查询用户列表，支持多种筛选条件
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码，默认1"
// @Param pageSize query int false "每页数量，默认10"
// @Param username query string false "用户名模糊查询"
// @Param status query int false "状态：0-禁用，1-启用"
// @Success 200 {object} api.Response{data=user.ListUserResponse} "查询成功"
// @Router /user/list [post]
func (a User) List(ctx *gin.Context) {
	// 参数校验
	req := &user.SearchUser{}
	err := a.Bind(ctx, req)
	if err != nil {
		a.Error(err)
		return
	}

	// 实例化逻辑层
	logic := iam.NewUserLogic(ctx)
	res, err := logic.List(req)
	if err != nil {
		a.Error(err)
		return
	}

	// 接口返回
	a.Success(res, "查询列表成功")
}

// Del 删除用户
func (a User) Del(ctx *gin.Context) {
	req := &user.DelUser{}
	err := a.Bind(ctx, req)
	if err != nil {
		a.Error(err)
		return
	}

	// 实例化逻辑层
	logic := iam.NewUserLogic(ctx)
	res, err := logic.Del(req.Ids)
	if err != nil {
		a.Error(err)
		return
	}

	// 接口返回
	a.Success(res, "删除成功")
}

// Info 获取当前登录用户信息
// @Summary 获取当前用户信息
// @Description 获取当前登录用户的详细信息
// @Tags 用户认证
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} api.Response{data=user.UserInfo} "获取成功"
// @Failure 401 {object} api.Response "未登录或token无效"
// @Router /user/info [get]
func (a User) Info(ctx *gin.Context) {
	// 设置上下文（重要：必须在调用其他方法前设置）
	a.Ctx = ctx

	// 从上下文获取当前用户ID（由JWT中间件设置）
	userIDStr, exists := ctx.Get("currUserId")
	if !exists {
		a.Error(errors.New("未登录或token无效"))
		return
	}

	// 将字符串ID转换为uint64
	userID, err := strconv.ParseUint(userIDStr.(string), 10, 64)
	if err != nil {
		a.Error(errors.New("用户ID格式无效"))
		return
	}

	// 实例化逻辑层
	logic := iam.NewUserLogic(ctx)
	res, err := logic.Get(userID)
	if err != nil {
		a.Error(err)
		return
	}

	// 接口返回
	a.Success(res, "获取成功")
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口，返回 JWT Token
// @Tags 用户认证
// @Accept json
// @Produce json
// @Param request body user.LoginRequest true "登录请求参数"
// @Success 200 {object} api.Response{data=user.LoginResponse} "登录成功"
// @Failure 400 {object} api.Response "参数错误"
// @Failure 401 {object} api.Response "用户名或密码错误"
// @Router /auth/login [post]
func (a User) Login(ctx *gin.Context) {
	req := &user.LoginRequest{}
	err := a.Bind(ctx, req)
	if err != nil {
		a.Error(err)
		return
	}

	logic := iam.NewUserLogic(ctx)
	res, err := logic.Login(req)
	if err != nil {
		a.Error(err)
		return
	}

	// 返回前端期望的格式
	a.Success(res, "登录成功")
}

// Logout 用户登出
func (a User) Logout(ctx *gin.Context) {
	// 设置上下文
	a.Ctx = ctx

	// 从请求头获取Token
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		a.Success(nil, "登出成功")
		return
	}

	// 提取 token（Bearer token）
	// 格式: Bearer <token>
	token := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		token = authHeader
	}

	// TODO: 将Token加入Redis黑名单，设置过期时间为token剩余有效期
	// 这里暂时直接返回成功，实际业务中需要：
	// 1. 解析token获取过期时间
	// 2. 将token加入redis黑名单，key为 blacklisted_token:<token>
	// 3. 设置过期时间为 token 剩余有效期
	_ = token

	a.Success(nil, "登出成功")
}

// Refresh 刷新Token
func (a User) Refresh(ctx *gin.Context) {
	// TODO: 实现Token刷新逻辑
	a.Success(nil, "Token刷新成功")
}

// Codes 获取用户权限码列表
func (a User) Codes(ctx *gin.Context) {
	// 设置上下文
	a.Ctx = ctx

	// 从上下文获取当前用户角色信息
	roleIds, exists := ctx.Get("currRoleIds")
	if !exists {
		// 如果没有角色信息，返回空数组
		a.Success([]string{}, "获取成功")
		return
	}

	// 根据角色ID解析权限码
	// 这里可以根据实际业务从数据库查询权限码
	// 简单实现：默认返回一些常用权限码
	codes := []string{
		"ACCOUNT",
		"ACCOUNT_ADD",
		"ACCOUNT_EDIT",
		"ACCOUNT_DELETE",
		"ACCOUNT_VIEW",
		"ROLE",
		"ROLE_ADD",
		"ROLE_EDIT",
		"ROLE_DELETE",
		"ROLE_VIEW",
		"DICT",
		"DICT_ADD",
		"DICT_EDIT",
		"DICT_DELETE",
		"DICT_VIEW",
	}

	// 如果有特定角色，可以添加更多权限码
	if roleIdsStr, ok := roleIds.(string); ok && roleIdsStr != "" {
		// 管理员角色拥有所有权限
		if roleIdsStr == "1" || roleIdsStr == "admin" {
			codes = append(codes, "ADMIN", "SYSTEM", "SYSTEM_SETTING")
		}
	}

	a.Success(codes, "获取成功")
}
