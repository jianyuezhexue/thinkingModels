package router

import (
	"thinkingModels/api/master"
	"thinkingModels/api/user"

	"github.com/gin-gonic/gin"
)

// 鉴权路由
func AuthorizedRouters() {
	authorizedRouters := func(router *gin.Engine) {
		v1 := router.Group("/v1")

		// 超级字典
		superDictionaryApi := master.NewSuperDictionary()
		superDictionaryGroup := v1.Group("/master/superDictionary")
		superDictionaryGroup.POST("", superDictionaryApi.Create)
		superDictionaryGroup.PUT("", superDictionaryApi.Update)
		superDictionaryGroup.POST("/:id", superDictionaryApi.Get)
		superDictionaryGroup.POST("/list", superDictionaryApi.List)
		superDictionaryGroup.DELETE("", superDictionaryApi.Del)
		superDictionaryGroup.GET("/tree", superDictionaryApi.Tree)
		superDictionaryGroup.POST("/children", superDictionaryApi.Children)

		// 用户管理
		userApi := user.NewUser()
		userGroup := v1.Group("/user")
		userGroup.POST("", userApi.Create)
		userGroup.PUT("", userApi.Update)
		userGroup.POST("/:id", userApi.Get)
		userGroup.POST("/list", userApi.List)
		userGroup.DELETE("", userApi.Del)

		// 认证相关（非鉴权）
		authApi := user.NewAuthController()
		authGroup := v1.Group("/auth")
		authGroup.POST("/login", authApi.Login)
		authGroup.POST("/logout", authApi.Logout)
		authGroup.POST("/refresh", authApi.Refresh)
	}
	Routers = append(Routers, authorizedRouters)
}

// 非鉴权路由
func UnAuthorizedRouters() {
	unAuthorizedRouters := func(router *gin.Engine) {
		// 无公开路由
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
