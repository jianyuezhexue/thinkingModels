package middleware

import (
	"github.com/gin-gonic/gin"
)

// 预埋用户信息
func UserInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 模拟从jwt中读出了用户信息
		currUserId := "1"
		currUserName := "buildBlock"

		// 预埋到上下文
		c.Set("currUserId", currUserId)
		c.Set("currUserName", currUserName)

		// 继续执行后续的中间件或者处理器函数
		c.Next()
	}
}
