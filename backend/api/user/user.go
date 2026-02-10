package user

import (
	"github.com/gin-gonic/gin"
	"thinkingModels/api"
	userDomain "thinkingModels/domain/user/user"
	"thinkingModels/logic/user"
)

type User struct {
	api.Base
}

func NewUser() *User {
	return &User{}
}

// Create 创建用户（注册）
func (a User) Create(ctx *gin.Context) {
	// 参数校验
	req := &userDomain.RegisterRequest{}
	err := a.Bind(ctx, req)
	if err != nil {
		a.Error(err)
		return
	}

	// 实例化逻辑层
	logic := user.NewUserLogic(ctx)
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
	req := &userDomain.UpdateUserRequest{}
	err := a.Bind(ctx, req)
	if err != nil {
		a.Error(err)
		return
	}

	// 实例化逻辑层
	logic := user.NewUserLogic(ctx)
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
	logic := user.NewUserLogic(ctx)
	res, err := logic.Get(0) // TODO: 使用实际的id
	if err != nil {
		a.Error(err)
		return
	}

	// 接口返回
	a.Success(res, "查询成功")
}

// List 查询用户列表
func (a User) List(ctx *gin.Context) {
	// 参数校验
	req := &userDomain.SearchUser{}
	err := a.Bind(ctx, req)
	if err != nil {
		a.Error(err)
		return
	}

	// 实例化逻辑层
	logic := user.NewUserLogic(ctx)
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
	req := &userDomain.DelUser{}
	err := a.Bind(ctx, req)
	if err != nil {
		a.Error(err)
		return
	}

	// 实例化逻辑层
	logic := user.NewUserLogic(ctx)
	res, err := logic.Del(req.Ids)
	if err != nil {
		a.Error(err)
		return
	}

	// 接口返回
	a.Success(res, "删除成功")
}
