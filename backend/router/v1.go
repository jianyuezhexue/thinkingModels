package router

import (
	"thinkingModels/api/iam"
	"thinkingModels/api/master"
	"thinkingModels/api/subject"
	"thinkingModels/api/practice"

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

		// ==================== 主数据模块 (Master) ====================
		// 模型分类管理
		masterCategoryApi := master.NewCategory()
		masterCategoryGroup := api.Group("/master/category")
		masterCategoryGroup.GET("/all", masterCategoryApi.All)     // 全量列表（按热度降序）
		masterCategoryGroup.GET("/list", masterCategoryApi.List)   // 分页列表
		masterCategoryGroup.POST("", masterCategoryApi.Create)     // 新建分类
		masterCategoryGroup.PUT("", masterCategoryApi.Update)    // 更新分类
		masterCategoryGroup.GET("/:id", masterCategoryApi.Get)     // 查询详情
		masterCategoryGroup.DELETE("", masterCategoryApi.Del)     // 删除分类
		masterCategoryGroup.POST("/increaseHeat", masterCategoryApi.IncreaseHeat) // 增加热度

		// ==================== 课题管理模块 ====================
		// 课题管理
		topicApi := subject.NewTopic()
		topicGroup := api.Group("/subject/topic")
		// 注意：具体路由必须放在参数路由（/:id）之前
		topicGroup.GET("/list", topicApi.List)
		topicGroup.GET("/my", topicApi.ListByUser)
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
		analysisGroup.GET("/my", analysisApi.ListByUser)
		analysisGroup.GET("/current", analysisApi.GetCurrent)
		analysisGroup.GET("/latest", analysisApi.GetLatest)
		analysisGroup.GET("/by-topic/:topicId", analysisApi.ListByTopic)
		analysisGroup.GET("/history/:topicId/:modelId", analysisApi.GetHistory)
		analysisGroup.POST("/set-current", analysisApi.SetCurrent)
		analysisGroup.POST("", analysisApi.Create)
		analysisGroup.PUT("", analysisApi.Update)
		analysisGroup.GET("/:id", analysisApi.Get)
		analysisGroup.DELETE("", analysisApi.Del)

		// ==================== Thinking 模块 (MVP) ====================
		// 思维模型管理
		thinkingModelApi := thinking.NewModel()
		thinkingModelGroup := api.Group("/thinking/model")
		thinkingModelGroup.GET("/list", thinkingModelApi.List)
		thinkingModelGroup.GET("/my", thinkingModelApi.ListMy)
		thinkingModelGroup.GET("/code/:code", thinkingModelApi.GetByCode)
		thinkingModelGroup.POST("", thinkingModelApi.Create)
		thinkingModelGroup.PUT("", thinkingModelApi.Update)
		thinkingModelGroup.GET("/:id", thinkingModelApi.Get)
		thinkingModelGroup.DELETE("", thinkingModelApi.Del)
		thinkingModelGroup.POST("/publish", thinkingModelApi.Publish)
		thinkingModelGroup.POST("/unpublish/:id", thinkingModelApi.Unpublish)
		thinkingModelGroup.POST("/fork", thinkingModelApi.Fork)

		// 模型分类管理
		thinkingCategoryApi := thinking.NewCategory()
		thinkingCategoryGroup := api.Group("/thinking/category")
		thinkingCategoryGroup.GET("/list", thinkingCategoryApi.List)
		thinkingCategoryGroup.GET("/tree", thinkingCategoryApi.Tree)
		thinkingCategoryGroup.GET("/children/:id", thinkingCategoryApi.Children)
		thinkingCategoryGroup.POST("/move", thinkingCategoryApi.Move)
		thinkingCategoryGroup.POST("/status", thinkingCategoryApi.UpdateStatus)
		thinkingCategoryGroup.POST("", thinkingCategoryApi.Create)
		thinkingCategoryGroup.PUT("", thinkingCategoryApi.Update)
		thinkingCategoryGroup.GET("/:id", thinkingCategoryApi.Get)
		thinkingCategoryGroup.DELETE("", thinkingCategoryApi.Del)

		// 模型标签管理
		thinkingTagApi := thinking.NewTag()
		thinkingTagGroup := api.Group("/thinking/tag")
		thinkingTagGroup.GET("/model/:modelId", thinkingTagApi.GetByModel)
		thinkingTagGroup.POST("/model", thinkingTagApi.AddToModel)
		thinkingTagGroup.DELETE("/model", thinkingTagApi.RemoveFromModel)
		thinkingTagGroup.GET("/hot", thinkingTagApi.Hot)

		// 课题管理
		thinkingTopicApi := thinking.NewTopic()
		thinkingTopicGroup := api.Group("/thinking/topic")
		thinkingTopicGroup.GET("/list", thinkingTopicApi.List)
		thinkingTopicGroup.GET("/my", thinkingTopicApi.ListMy)
		thinkingTopicGroup.POST("/select-model", thinkingTopicApi.SelectModel)
		thinkingTopicGroup.POST("/remove-model/:id", thinkingTopicApi.RemoveModel)
		thinkingTopicGroup.POST("/status", thinkingTopicApi.UpdateStatus)
		thinkingTopicGroup.POST("/complete/:id", thinkingTopicApi.Complete)
		thinkingTopicGroup.POST("/archive/:id", thinkingTopicApi.Archive)
		thinkingTopicGroup.POST("/reopen/:id", thinkingTopicApi.Reopen)
		thinkingTopicGroup.GET("/statistics", thinkingTopicApi.Statistics)
		thinkingTopicGroup.POST("", thinkingTopicApi.Create)
		thinkingTopicGroup.PUT("", thinkingTopicApi.Update)
		thinkingTopicGroup.GET("/:id", thinkingTopicApi.Get)
		thinkingTopicGroup.DELETE("", thinkingTopicApi.Del)

		// 分析记录管理
		thinkingAnalysisApi := thinking.NewAnalysis()
		thinkingAnalysisGroup := api.Group("/thinking/analysis")
		thinkingAnalysisGroup.POST("/save-with-ai", thinkingAnalysisApi.SaveWithAi)
		thinkingAnalysisGroup.GET("/list", thinkingAnalysisApi.List)
		thinkingAnalysisGroup.GET("/my", thinkingAnalysisApi.ListMy)
		thinkingAnalysisGroup.GET("/current", thinkingAnalysisApi.GetCurrent)
		thinkingAnalysisGroup.GET("/latest", thinkingAnalysisApi.GetLatest)
		thinkingAnalysisGroup.GET("/by-topic/:topicId", thinkingAnalysisApi.ListByTopic)
		thinkingAnalysisGroup.GET("/history/:topicId/:modelId", thinkingAnalysisApi.GetHistory)
		thinkingAnalysisGroup.POST("/set-current", thinkingAnalysisApi.SetCurrent)
		thinkingAnalysisGroup.POST("", thinkingAnalysisApi.Create)
		thinkingAnalysisGroup.PUT("", thinkingAnalysisApi.Update)
		thinkingAnalysisGroup.GET("/:id", thinkingAnalysisApi.Get)
		thinkingAnalysisGroup.DELETE("", thinkingAnalysisApi.Del)

		// 行动项管理
		thinkingActionApi := thinking.NewAction()
		thinkingActionGroup := api.Group("/thinking/action")
		thinkingActionGroup.GET("/list", thinkingActionApi.List)
		thinkingActionGroup.GET("/my", thinkingActionApi.ListMy)
		thinkingActionGroup.POST("/from-analysis", thinkingActionApi.CreateFromAnalysis)
		thinkingActionGroup.GET("/by-topic/:topicId", thinkingActionApi.ListByTopic)
		thinkingActionGroup.GET("/by-analysis/:analysisId", thinkingActionApi.ListByAnalysis)
		thinkingActionGroup.POST("/progress", thinkingActionApi.UpdateProgress)
		thinkingActionGroup.POST("/complete/:id", thinkingActionApi.Complete)
		thinkingActionGroup.POST("/cancel/:id", thinkingActionApi.Cancel)
		thinkingActionGroup.GET("/statistics", thinkingActionApi.Statistics)
		thinkingActionGroup.POST("", thinkingActionApi.Create)
		thinkingActionGroup.PUT("", thinkingActionApi.Update)
		thinkingActionGroup.GET("/:id", thinkingActionApi.Get)
		thinkingActionGroup.DELETE("", thinkingActionApi.Del)

		// 跟进记录管理
		thinkingFollowUpApi := thinking.NewFollowUp()
		thinkingFollowUpGroup := api.Group("/thinking/followup")
		thinkingFollowUpGroup.GET("/by-action/:actionId", thinkingFollowUpApi.ListByAction)
		thinkingFollowUpGroup.POST("", thinkingFollowUpApi.Create)
		thinkingFollowUpGroup.PUT("", thinkingFollowUpApi.Update)
		thinkingFollowUpGroup.GET("/:id", thinkingFollowUpApi.Get)
		thinkingFollowUpGroup.DELETE("", thinkingFollowUpApi.Del)
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
