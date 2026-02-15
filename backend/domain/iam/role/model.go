package role

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/jianyuezhexue/base"
	"thinkingModels/component/db"
)

// RoleEntityInterface 角色实体接口
type RoleEntityInterface interface {
	base.BaseModelInterface[RoleEntity]
}

// RoleEntity 角色实体
type RoleEntity struct {
	base.BaseModel[RoleEntity]
	RoleName     string `json:"roleName" type:"db" comment:"角色名称"`
	RoleCode     string `json:"roleCode" type:"db" comment:"角色编码"`
	Description  string `json:"description" type:"db" comment:"角色描述"`
	Status       int    `json:"status" type:"db" comment:"状态:0=禁用,1=正常"`
	Sort         int    `json:"sort" type:"db" comment:"排序"`
	MenuIds      string `json:"menuIds" type:"db" comment:"菜单ID列表，逗号分隔"`
	EnterpriseID uint64 `json:"enterpriseId" type:"db" comment:"企业ID"`
}

// 实例化角色实体
func NewRoleEntity(ctx *gin.Context, opt ...base.Option[RoleEntity]) RoleEntityInterface {
	entity := &RoleEntity{}
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
func (m *RoleEntity) TableName() string {
	return "roles"
}

// Validate 数据校验
func (m *RoleEntity) Validate() error {
	// 角色名称不能为空
	if m.RoleName == "" {
		return errors.New("角色名称不能为空")
	}
	// 角色编码不能为空
	if m.RoleCode == "" {
		return errors.New("角色编码不能为空")
	}
	return nil
}

// Repair 数据修复
func (m *RoleEntity) Repair() error {
	// 设置默认状态
	if m.Status == 0 {
		m.Status = 1
	}
	return nil
}

// Complete 数据完善
func (m *RoleEntity) Complete() error {
	// 数据完善逻辑
	return nil
}