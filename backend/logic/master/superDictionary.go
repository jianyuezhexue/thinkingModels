package master

import (
	"github.com/gin-gonic/gin"

	"thinkingModels/domain/master/superDictionary"
	"thinkingModels/logic"
)

// 公共
type SuperDictionaryLogic struct {
	logic.BaseLogic
}

// 初始化SuperDictionaryLogic
func NewSuperDictionaryLogic(ctx *gin.Context) *SuperDictionaryLogic {
	return &SuperDictionaryLogic{BaseLogic: logic.BaseLogic{Ctx: ctx}}
}

// 新增
func (l *SuperDictionaryLogic) Create(req *superDictionary.CreateSuperDictionary) (*superDictionary.SuperDictionaryEntity, error) {

	// 实例化模型
	superDictionaryEntity := superDictionary.NewSuperDictionaryEntity(l.Ctx)

	// 数据赋值
	_, err := superDictionaryEntity.SetData(req)
	if err != nil {
		return nil, err
	}

	// 数据校验
	err = superDictionaryEntity.Validate()
	if err != nil {
		return nil, err
	}

	// 保存数据
	res, err := superDictionaryEntity.Create()
	if err != nil {
		return nil, err
	}

	return res, nil
}

// 更新
func (l *SuperDictionaryLogic) Update(req *superDictionary.UpdateSuperDictionary) (*superDictionary.SuperDictionaryEntity, error) {

	// 实例化模型
	superDictionaryEntity := superDictionary.NewSuperDictionaryEntity(l.Ctx)

	// 先加载旧数据
	_, err := superDictionaryEntity.LoadById(req.Id)
	if err != nil {
		return nil, err
	}

	// 数据赋值
	_, err = superDictionaryEntity.SetData(req)
	if err != nil {
		return nil, err
	}

	// 数据校验
	err = superDictionaryEntity.Validate()
	if err != nil {
		return nil, err
	}

	// 更新数据
	res, err := superDictionaryEntity.Update()
	return res, err
}

// 查询
func (l *SuperDictionaryLogic) Get(req *superDictionary.SearchSuperDictionary) (*superDictionary.SuperDictionaryEntity, error) {

	// 实例化模型
	superDictionaryEntity := superDictionary.NewSuperDictionaryEntity(l.Ctx)

	// 获取搜索条件
	cond := superDictionaryEntity.MakeConditon(*req)

	// 查询数据
	res, err := superDictionaryEntity.LoadData(cond)
	if err != nil {
		return nil, err
	}

	return res, err
}

// 查询列表
func (l *SuperDictionaryLogic) List(req *superDictionary.SearchSuperDictionary) (*superDictionary.ListReap, error) {

	// 实例化模型
	superDictionaryEntity := superDictionary.NewSuperDictionaryEntity(l.Ctx)

	// 获取搜索条件
	cond := superDictionaryEntity.MakeConditon(*req)

	// 查询总数
	total, err := superDictionaryEntity.Count(cond)
	if err != nil {
		return nil, err
	}

	// 获取列表
	list, err := superDictionaryEntity.List(cond)
	if err != nil {
		return nil, err
	}

	// 返回数据
	res := &superDictionary.ListReap{Page: req.Page, PageSize: req.PageSize, Total: total, List: list}
	return res, nil
}

// 删除
func (l *SuperDictionaryLogic) Del(req *superDictionary.DelSuperDictionary) (any, error) {

	// 实例化模型
	superDictionaryEntity := superDictionary.NewSuperDictionaryEntity(l.Ctx)

	// 删除数据
	err := superDictionaryEntity.Del(req.Ids...)
	return nil, err
}

// 树形结构
func (l *SuperDictionaryLogic) Tree(req *superDictionary.TreeSuperDictionary) ([]*superDictionary.TreeNode, error) {

	// 实例化模型
	superDictionaryEntity := superDictionary.NewSuperDictionaryEntity(l.Ctx)

	// 查询所有数据
	searchReq := &superDictionary.SearchSuperDictionary{}
	cond := superDictionaryEntity.MakeConditon(*searchReq)
	list, err := superDictionaryEntity.List(cond)
	if err != nil {
		return nil, err
	}

	// 构建树形结构
	return buildTree(list, req.ParentId), nil
}

// 构建树形结构
func buildTree(list []*superDictionary.SuperDictionaryEntity, parentId int64) []*superDictionary.TreeNode {
	// 使用 map 记录已访问的节点，防止循环引用
	return buildTreeWithVisited(list, parentId, make(map[uint64]bool))
}

// 带访问记录的递归构建树
func buildTreeWithVisited(list []*superDictionary.SuperDictionaryEntity, parentId int64, visited map[uint64]bool) []*superDictionary.TreeNode {
	var tree []*superDictionary.TreeNode

	for _, item := range list {
		// 检查是否匹配父节点
		if item.ParentId == parentId {
			// 防止循环引用：如果节点已被访问过，跳过
			if visited[item.Id] {
				continue
			}

			// 防止自引用：如果 ParentId == Id，跳过
			if item.ParentId == int64(item.Id) {
				continue
			}

			// 标记为已访问
			visited[item.Id] = true

			node := &superDictionary.TreeNode{
				SuperDictionaryEntity: item,
				Children:              buildTreeWithVisited(list, int64(item.Id), visited),
			}
			tree = append(tree, node)
		}
	}

	return tree
}

// 获取子节点
func (l *SuperDictionaryLogic) Children(req *superDictionary.ChildrenSuperDictionary) ([]*superDictionary.SuperDictionaryEntity, error) {

	// 实例化模型
	superDictionaryEntity := superDictionary.NewSuperDictionaryEntity(l.Ctx)

	// 构建查询条件
	searchReq := &superDictionary.SearchSuperDictionary{
		ParentId: req.ParentId,
	}
	cond := superDictionaryEntity.MakeConditon(*searchReq)

	// 查询直接子节点
	list, err := superDictionaryEntity.List(cond)
	if err != nil {
		return nil, err
	}

	return list, nil
}
