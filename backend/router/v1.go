package router

import (
	"thinkingModels/api/iam"
	"thinkingModels/api/market"
	"thinkingModels/api/master"
	"thinkingModels/api/subject"

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

		// ==================== 思维模型市场模块 ====================
		// 思维模型管理
		modelApi := market.NewModel()
		modelGroup := api.Group("/market/model")
		// 注意：具体路由必须放在参数路由（/:id）之前
		modelGroup.GET("/list", modelApi.List)
		modelGroup.GET("/code/:code", modelApi.GetByCode)
		modelGroup.POST("", modelApi.Create)
		modelGroup.PUT("", modelApi.Update)
		modelGroup.GET("/:id", modelApi.Get)
		modelGroup.DELETE("", modelApi.Del)
		modelGroup.POST("/publish", modelApi.Publish)
		modelGroup.POST("/unpublish", modelApi.Unpublish)
		modelGroup.POST("/status", modelApi.UpdateStatus)

		// 模型分类管理
		categoryApi := market.NewCategory()
		categoryGroup := api.Group("/market/category")
		// 注意：具体路由必须放在参数路由（/:id）之前
		categoryGroup.GET("/list", categoryApi.List)
		categoryGroup.GET("/tree", categoryApi.Tree)
		categoryGroup.GET("/tree/all", categoryApi.TreeWithRoot)
		categoryGroup.GET("/children/:id", categoryApi.Children)
		categoryGroup.GET("/path/:id", categoryApi.Path)
		categoryGroup.POST("/move", categoryApi.Move)
		categoryGroup.POST("/status", categoryApi.UpdateStatus)
		categoryGroup.POST("", categoryApi.Create)
		categoryGroup.PUT("", categoryApi.Update)
		categoryGroup.GET("/:id", categoryApi.Get)
		categoryGroup.DELETE("", categoryApi.Del)

		// ==================== 课题管理模块 ====================
		// 课题管理
		topicApi := subject.NewTopic()
		topicGroup := api.Group("/subject/topic")
		// 注意：具体路由必须放在参数路由（/:id）之前
		topicGroup.GET("/list", topicApi.List)
		topicGroup.GET("/my-list", topicApi.ListByUser)
		topicGroup.POST("/status", topicApi.UpdateStatus)
		topicGroup.POST("/select-model", topicApi.SelectModel)
		topicGroup.POST("/remove-model/:id", topicApi.RemoveModel)
		topicGroup.POST("/complete", topicApi.Complete)
		topicGroup.POST("/archive", topicApi.Archive)
		topicGroup.POST("/reopen/:id", topicApi.Reopen)
		topicGroup.GET("/statistics", topicApi.GetStatistics)
		topicGroup.POST("", topicApi.Create)
		topicGroup.PUT("", topicApi.Update)
		topicGroup.GET("/:id", topicApi.Get)
		topicGroup.DELETE("", topicApi.Del)

		// 课题分析记录
		analysisApi := subject.NewAnalysis()
		analysisGroup := api.Group("/subject/analysis")
		// 注意：具体路由必须放在参数路由（/:id）之前
		analysisGroup.POST("/save-with-ai", analysisApi.SaveWithAi)
		analysisGroup.GET("/list", analysisApi.List)
		analysisGroup.GET("/my-list", analysisApi.ListByUser)
		analysisGroup.GET("/current", analysisApi.GetCurrent)
		analysisGroup.GET("/latest", analysisApi.GetLatest)
		analysisGroup.GET("/by-topic/:topicId", analysisApi.ListByTopic)
		analysisGroup.GET("/history/:topicId/:modelId", analysisApi.GetHistory)
		analysisGroup.POST("/set-current", analysisApi.SetCurrent)
		analysisGroup.POST("", analysisApi.Create)
		analysisGroup.PUT("", analysisApi.Update)
		analysisGroup.GET("/:id", analysisApi.Get)
		analysisGroup.DELETE("", analysisApi.Del)
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
