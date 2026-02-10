package master

import (
	"github.com/gin-gonic/gin"
	"thinkingModels/api"
	"thinkingModels/domain/master/superDictionary"
	"thinkingModels/logic/master"
)

type SuperDictionary struct {
	api.Base
}

func NewSuperDictionary() *SuperDictionary {
	return &SuperDictionary{}
}

// 新建
func (a SuperDictionary) Create(ctx *gin.Context) {

	// 参数校验
	req := &superDictionary.CreateSuperDictionary{}
	err := a.Bind(ctx, req)
	if err != nil {
		a.Error(err)
		return
	}

	// 实例化逻辑层
	logic := master.NewSuperDictionaryLogic(ctx)
	res, err := logic.Create(req)
	if err != nil {
		a.Error(err)
		return
	}

	// 接口返回
	a.Success(res, "新建成功")
}

// 更新
func (a SuperDictionary) Update(ctx *gin.Context) {
	// 参数校验
	req := &superDictionary.UpdateSuperDictionary{}
	err := a.Bind(ctx, req)
	if err != nil {
		a.Error(err)
		return
	}

	// 实例化逻辑层
	logic := master.NewSuperDictionaryLogic(ctx)
	res, err := logic.Update(req)
	if err != nil {
		a.Error(err)
		return
	}

	// 接口返回
	a.Success(res, "更新成功")
}

// 查询
func (a SuperDictionary) Get(ctx *gin.Context) {
	req := &superDictionary.SearchSuperDictionary{}
	err := a.Bind(ctx, req)
	if err != nil {
		a.Error(err)
		return
	}

	// 实例化逻辑层
	logic := master.NewSuperDictionaryLogic(ctx)
	res, err := logic.Get(req)
	if err != nil {
		a.Error(err)
		return
	}

	// 接口返回
	a.Success(res, "查询成功")
}

// 列表
func (a SuperDictionary) List(ctx *gin.Context) {

	// 参数校验
	req := &superDictionary.SearchSuperDictionary{}
	err := a.Bind(ctx, req)
	if err != nil {
		a.Error(err)
		return
	}

	// 实例化逻辑层
	logic := master.NewSuperDictionaryLogic(ctx)
	res, err := logic.List(req)
	if err != nil {
		a.Error(err)
		return
	}

	// 接口返回
	a.Success(res, "查询列表成功")
}

// 删除
func (a SuperDictionary) Del(ctx *gin.Context) {
	req := &superDictionary.DelSuperDictionary{}
	err := a.Bind(ctx, req)
	if err != nil {
		a.Error(err)
		return
	}

	// 实例化逻辑层
	logic := master.NewSuperDictionaryLogic(ctx)
	res, err := logic.Del(req)
	if err != nil {
		a.Error(err)
		return
	}

	// 接口返回
	a.Success(res, "删除成功")
}

// 树形结构
func (a SuperDictionary) Tree(ctx *gin.Context) {
	// 参数校验
	req := &superDictionary.TreeSuperDictionary{}
	err := a.Bind(ctx, req)
	if err != nil {
		a.Error(err)
		return
	}

	// 实例化逻辑层
	logic := master.NewSuperDictionaryLogic(ctx)
	res, err := logic.Tree(req)
	if err != nil {
		a.Error(err)
		return
	}

	// 接口返回
	a.Success(res, "查询树形结构成功")
}

// 获取子节点
func (a SuperDictionary) Children(ctx *gin.Context) {
	// 参数校验
	req := &superDictionary.ChildrenSuperDictionary{}
	err := a.Bind(ctx, req)
	if err != nil {
		a.Error(err)
		return
	}

	// 实例化逻辑层
	logic := master.NewSuperDictionaryLogic(ctx)
	res, err := logic.Children(req)
	if err != nil {
		a.Error(err)
		return
	}

	// 接口返回
	a.Success(res, "查询子节点成功")
}
