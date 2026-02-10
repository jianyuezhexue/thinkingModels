package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// 代码生成路由,自定义代码不要写在这里，否则会被覆盖
func genCodeRouters() {
	genCodeRouters := func(router *gin.Engine) {
		v1 := router.Group("/v1")
		fmt.Println(v1)
	}
	Routers = append(Routers, genCodeRouters)
}
