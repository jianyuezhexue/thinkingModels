package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Base struct {
	Ctx *gin.Context
}

type Response struct {
	Code    int64  `json:"code"`
	Msg     string `json:"msg"`
	Data    any    `json:"data"`
	TraceId any    `json:"trace_id"`
}

// 返回失败
func (a *Base) Error(err error, ext ...any) {
	tracdId := a.Ctx.Request.Header.Get("x-trace-id")
	res := Response{TraceId: tracdId, Code: -1, Msg: err.Error()}
	a.Ctx.JSON(200, res)
}

// 返回成功
func (a *Base) Success(data interface{}, ext ...any) {
	msg := "success"
	if len(ext) > 0 {
		msg = ext[0].(string)
	}

	var code int64
	if len(ext) > 1 {
		code = ext[1].(int64)
	}

	tracdId := a.Ctx.Request.Header.Get("x-trace-id")
	res := Response{TraceId: tracdId, Code: code, Msg: msg, Data: data}
	a.Ctx.JSON(200, res)
}

// 绑定参数
func (a *Base) Bind(ctx *gin.Context, d any, bindings ...binding.Binding) error {
	a.Ctx = ctx

	if d == nil {
		return nil
	}

	err := a.Ctx.ShouldBind(d)
	if err != nil {
		return err
	}
	return nil
}
