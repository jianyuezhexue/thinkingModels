package engine

import (
	"fmt"

	"thinkingModels/config"
	"thinkingModels/router"

	"github.com/gin-gonic/gin"
)

// 启动服务
func Run() {

	// todo 设置模式
	// gin.SetMode(gin.ReleaseMode)

	// 实例化引擎
	r := gin.Default()

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 使用全局限流中间件

	// 注册所有路由
	router.InitRouter(r)

	// 启动定时脚本

	// 2500端口
	r.Run(fmt.Sprintf("%s:%s", config.Config.Host, config.Config.Port))
}
