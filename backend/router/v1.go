package router

import (
	"thinkingModels/api/iam"
	"thinkingModels/api/master"
	"thinkingModels/api/thinkingModel"

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

		// 角色管理
		roleApi := iam.NewRole()
		roleGroup := api.Group("/role")
		roleGroup.GET("/all", roleApi.All)                     // 全量列表（用于下拉选择）
		roleGroup.POST("/list", roleApi.List)                  // 分页列表
		roleGroup.POST("", roleApi.Create)                     // 新建角色
		roleGroup.PUT("", roleApi.Update)                      // 更新角色
		roleGroup.PUT("/permission", roleApi.UpdatePermission) // 更新权限
		roleGroup.GET("/:id", roleApi.Get)                     // 查询详情
		roleGroup.DELETE("", roleApi.Del)                      // 删除角色

		// 认证相关（非鉴权）
		authGroup := api.Group("/auth")
		authGroup.POST("/login", userApi.Login)
		authGroup.POST("/logout", userApi.Logout)
		authGroup.POST("/refresh", userApi.Refresh)
		authGroup.GET("/codes", userApi.Codes) // 获取用户权限码

		// ==================== 主数据模块 (Master) ====================

		// ==================== ThinkingModel 模块 ====================
		// 思维模型管理
		tmModelApi := thinkingModel.NewModel()
		tmModelGroup := api.Group("/thinkingModel/model")
		tmModelGroup.GET("/list", tmModelApi.List)
		tmModelGroup.GET("/my", tmModelApi.ListMy)
		tmModelGroup.GET("/code/:code", tmModelApi.GetByCode)
		tmModelGroup.POST("", tmModelApi.Create)
		tmModelGroup.PUT("", tmModelApi.Update)
		tmModelGroup.GET("/:id", tmModelApi.Get)
		tmModelGroup.DELETE("", tmModelApi.Del)
		tmModelGroup.POST("/publish", tmModelApi.Publish)
		tmModelGroup.POST("/unpublish/:id", tmModelApi.Unpublish)
		tmModelGroup.POST("/fork", tmModelApi.Fork)

		// 思维模型分类管理
		tmCategoryApi := thinkingModel.NewCategory()
		tmCategoryGroup := api.Group("/thinkingModel/category")
		tmCategoryGroup.GET("/all", tmCategoryApi.All)
		tmCategoryGroup.GET("/list", tmCategoryApi.List)
		tmCategoryGroup.POST("", tmCategoryApi.Create)
		tmCategoryGroup.PUT("", tmCategoryApi.Update)
		tmCategoryGroup.GET("/:id", tmCategoryApi.Get)
		tmCategoryGroup.DELETE("", tmCategoryApi.Del)
		tmCategoryGroup.POST("/increaseHeat", tmCategoryApi.IncreaseHeat)
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
