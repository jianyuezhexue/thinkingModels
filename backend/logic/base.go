package logic

import "github.com/gin-gonic/gin"

// todo
// 这里可以放一些公共信息，如租户信息，用户信息等

type BaseLogic struct {
	Ctx      *gin.Context // 上下文
	CurrUser any          // 当前登陆用户
}

// todo

// 列表返回
type ListReap struct {
	List     any   `json:"list"`
	PageSize int64 `json:"pageSize"`
	Page     int64 `json:"page"`
	Total    int64 `json:"total"`
}
