package router

import (
	"github.com/gin-gonic/gin"
	"thinkingModels/middleware"
)

var Routers = make([]func(router *gin.Engine), 0)

// 实例化所有路由
func InitRouter(router *gin.Engine) {

	// 全局中间件 | 允许跨域
	router.Use(middleware.Cors()).Use(middleware.UserInfo())

	// 代码自动生成路由
	genCodeRouters()

	// 鉴权路由
	AuthorizedRouters()

	// 非鉴权路由
	UnAuthorizedRouters()

	// Oauth2路由
	Oauth2Routers()

	// 模块单元测试路由
	UnitTestRouters()

	// 批量注册
	for _, fn := range Routers {
		fn(router)
	}
}
