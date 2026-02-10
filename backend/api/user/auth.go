package user

import (
	"github.com/gin-gonic/gin"
	"thinkingModels/api"
	userDomain "thinkingModels/domain/user/user"
	"thinkingModels/logic/user"
)

// AuthController 认证控制器
type AuthController struct {
	api.Base
}

// NewAuthController 创建认证控制器
func NewAuthController() *AuthController {
	return &AuthController{}
}

// Login 用户登录
func (a AuthController) Login(ctx *gin.Context) {
	req := &userDomain.LoginRequest{}
	err := a.Bind(ctx, req)
	if err != nil {
		a.Error(err)
		return
	}

	logic := user.NewUserLogic(ctx)
	res, err := logic.Login(req)
	if err != nil {
		a.Error(err)
		return
	}

	// 返回前端期望的格式
	a.Success(gin.H{"accessToken": res.AccessToken}, "登录成功")
}

// Logout 用户登出
func (a AuthController) Logout(ctx *gin.Context) {
	// TODO: 从Header获取Token并登出
	a.Success(nil, "登出成功")
}

// Refresh 刷新Token
func (a AuthController) Refresh(ctx *gin.Context) {
	// TODO: 实现Token刷新逻辑
	a.Success(nil, "Token刷新成功")
}
