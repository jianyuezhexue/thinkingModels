package superDictionary

import (
	"github.com/gin-gonic/gin"
	"github.com/jianyuezhexue/base"
	"thinkingModels/component/db"
)

// 业务模型接口定义
type SuperDictionaryEntityInterface interface {
	base.BaseModelInterface[SuperDictionaryEntity]
}

// SuperDictionary 实体业务模型
type SuperDictionaryEntity struct {
	base.BaseModel[SuperDictionaryEntity]
	ParentId    int64  `json:"parentId" type:"db" comment:"父级ID"`      // 父级ID
	DictValue   string `json:"dictValue" type:"db" comment:"字典值"`      // 字典值
	DictName    string `json:"dictName" type:"db" comment:"字典名称"`      // 字典名称
	Level       int64  `json:"level" type:"db" comment:"层级"`           // 层级
	LevelName   string `json:"levelName" type:"db" comment:"层级名称"`     // 层级名称
	Description string `json:"description" type:"db" comment:"字典描述"`   // 字典描述
	Eval        string `json:"eval" type:"db" comment:"eval"`          // eval
	ExtSchema   string `json:"extSchema" type:"db" comment:"拓展Schema"` // 拓展Schema
	ExtJson     string `json:"extJson" type:"db" comment:"拓展json"`     // 拓展json

}

// 实例化领域业务模型
func NewSuperDictionaryEntity(ctx *gin.Context, opt ...base.Option[SuperDictionaryEntity]) SuperDictionaryEntityInterface {
	entity := &SuperDictionaryEntity{}
	entity.BaseModel = base.NewBaseModel(ctx, db.InitDb(), entity.TableName(), entity)

	// 自定义配置选项
	if len(opt) > 0 {
		for _, fc := range opt {
			fc(&entity.BaseModel)
		}
	}
	return entity
}

// 数据表名
func (m *SuperDictionaryEntity) TableName() string {
	return "super_dictionary"
}

// ValidateFunc 数据校验
func (m *SuperDictionaryEntity) Validate() error {
	// more...
	return nil
}

// Repair 数据修复
func (m *SuperDictionaryEntity) Repair() error {
	// more...
	return nil
}

// Complete 数据完善
func (m *SuperDictionaryEntity) Complete() error {
	// more...
	return nil
}

// more abilits...
