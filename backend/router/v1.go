package router

import (
	"thinkingModels/api/iam"
	"thinkingModels/api/master"

	"github.com/gin-gonic/gin"
)

// 鉴权路由
func AuthorizedRouters() {
	authorizedRouters := func(router *gin.Engine) {
		api := router.Group("")

		// 超级字典
		superDictionaryApi := master.NewSuperDictionary()
		superDictionaryGroup := api.Group("/master/superDictionary")
		superDictionaryGroup.POST("", superDictionaryApi.Create)
		superDictionaryGroup.PUT("", superDictionaryApi.Update)
		superDictionaryGroup.POST("/:id", superDictionaryApi.Get)
		superDictionaryGroup.POST("/list", superDictionaryApi.List)
		superDictionaryGroup.DELETE("", superDictionaryApi.Del)
		superDictionaryGroup.GET("/tree", superDictionaryApi.Tree)
		superDictionaryGroup.POST("/children", superDictionaryApi.Children)

		// 用户管理
		userApi := iam.NewUser()
		userGroup := api.Group("/user")
		userGroup.GET("/info", userApi.Info) // 获取当前登录用户信息
		userGroup.POST("", userApi.Create)
		userGroup.PUT("", userApi.Update)
		userGroup.POST("/:id", userApi.Get)
		userGroup.POST("/list", userApi.List)
		userGroup.DELETE("", userApi.Del)

		// 认证相关（非鉴权）
		authGroup := api.Group("/auth")
		authGroup.POST("/login", userApi.Login)
		authGroup.POST("/logout", userApi.Logout)
		authGroup.POST("/refresh", userApi.Refresh)
		authGroup.GET("/codes", userApi.Codes) // 获取用户权限码
	}
	Routers = append(Routers, authorizedRouters)
}

// 非鉴权路由
func UnAuthorizedRouters() {
	unAuthorizedRouters := func(router *gin.Engine) {
		api := router.Group("")

		// 用户注册（无需鉴权）
		userApi := iam.NewUser()
		api.POST("/auth/register", userApi.Create)
	}
	Routers = append(Routers, unAuthorizedRouters)
}

// oauth2路由
func Oauth2Routers() {
	oauth2Routers := func(router *gin.Engine) {
		// todo
	}
	Routers = append(Routers, oauth2Routers)
}

// 单元测试路由
func UnitTestRouters() {
	unitTestRouters := func(router *gin.Engine) {

	}
	Routers = append(Routers, unitTestRouters)
}
