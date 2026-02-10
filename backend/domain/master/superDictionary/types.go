package superDictionary

// 新增实体
type CreateSuperDictionary struct {
	ParentId    int64  `json:"parentId" uri:"parentId"`       // 父级ID
	DictValue   string `json:"dictValue" uri:"dictValue"`     // 字典值
	DictName    string `json:"dictName" uri:"dictName"`       // 字典名称
	Level       int64  `json:"level" uri:"level"`             // 层级
	LevelName   string `json:"levelName" uri:"levelName"`     // 层级名称
	Description string `json:"description" uri:"description"` // 字典描述
	Eval        string `json:"eval" uri:"eval"`               // eval
	ExtSchema   string `json:"extSchema" uri:"extSchema"`     // 拓展Schema
	ExtJson     string `json:"extJson" uri:"extJson"`         // 拓展json

}

// 更新实体
type UpdateSuperDictionary struct {
	Id          uint64 `json:"id" binding:"required"`
	ParentId    int64  `json:"parentId" uri:"parentId"`       // 父级ID
	DictValue   string `json:"dictValue" uri:"dictValue"`     // 字典值
	DictName    string `json:"dictName" uri:"dictName"`       // 字典名称
	Level       int64  `json:"level" uri:"level"`             // 层级
	LevelName   string `json:"levelName" uri:"levelName"`     // 层级名称
	Description string `json:"description" uri:"description"` // 字典描述
	Eval        string `json:"eval" uri:"eval"`               // eval
	ExtSchema   string `json:"extSchema" uri:"extSchema"`     // 拓展Schema
	ExtJson     string `json:"extJson" uri:"extJson"`         // 拓展json

}

// 搜索实体
type SearchSuperDictionary struct {
	Page        int64    `json:"page" uri:"page" search:"page"`                                                            // 分页
	PageSize    int64    `json:"pageSize" uri:"pageSize" search:"pageSize"`                                                // 分页大小
	Id          uint64   `json:"id" uri:"id" search:"type:eq;column:id;table:super_dictionary"`                            // ID
	Ids         []uint64 `json:"ids" uri:"ids" search:"type:in;column:id;table:super_dictionary"`                          // IDs
	CreatedAt   []string `json:"createdAt" uri:"createdAt" search:"type:between;column:created_at;table:super_dictionary"` // 创建时间
	UpdatedAt   []string `json:"updatedAt" uri:"updatedAt" search:"type:between;column:updated_at;table:super_dictionary"` // 创建时间
	ParentId    int64    `json:"parentId" uri:"parentId" search:"type:eq;column:parent_id;table:super_dictionary"`         // 父级ID
	DictValue   string   `json:"dictValue" uri:"dictValue" search:"type:eq;column:dict_value;table:super_dictionary"`      // 字典值
	DictName    string   `json:"dictName" uri:"dictName" search:"type:eq;column:dict_name;table:super_dictionary"`         // 字典名称
	Level       int64    `json:"level" uri:"level" search:"type:eq;column:level;table:super_dictionary"`                   // 层级
	LevelName   string   `json:"levelName" uri:"levelName" search:"type:eq;column:level_name;table:super_dictionary"`      // 层级名称
	Description string   `json:"description" uri:"description" search:"type:eq;column:description;table:super_dictionary"` // 字典描述
	Eval        string   `json:"eval" uri:"eval" search:"type:eq;column:eval;table:super_dictionary"`                      // eval
	ExtSchema   string   `json:"extSchema" uri:"extSchema" search:"type:eq;column:ext_schema;table:super_dictionary"`      // 拓展Schema
	ExtJson     string   `json:"extJson" uri:"extJson" search:"type:eq;column:ext_json;table:super_dictionary"`            // 拓展json

}

// 删除实体
type DelSuperDictionary struct {
	Ids []uint64 `json:"ids" binding:"required"`
}

// 分页返回
type ListReap struct {
	Page     int64                    `json:"page" comment:"页数"`
	PageSize int64                    `json:"pageSize" comment:"每页数量"`
	Total    int64                    `json:"total" comment:"总条数"`
	List     []*SuperDictionaryEntity `json:"list" comment:"数据"`
}

// 树形结构请求
type TreeSuperDictionary struct {
	ParentId int64 `json:"parentId" uri:"parentId"` // 父级ID,默认为0表示查询根节点
}

// 树形节点
type TreeNode struct {
	*SuperDictionaryEntity
	Children []*TreeNode `json:"children"` // 子节点
}

// 获取子节点请求
type ChildrenSuperDictionary struct {
	ParentId int64 `json:"parentId" uri:"parentId" binding:"required"` // 父级ID
}
