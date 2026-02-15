package role

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"thinkingModels/component/db"
)

// 获取测试用的 gin context
func getTestContext() *gin.Context {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(nil)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}
	ctx.Set("currUserId", "1")
	ctx.Set("currUserName", "系统")
	return ctx
}

// TestInitRolesTable 测试创建 roles 表
func TestInitRolesTable(t *testing.T) {
	gormDB := db.InitDb()

	// 创建 roles 表
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS roles (
		id              bigint unsigned auto_increment primary key comment '主键ID',
		role_name       varchar(50)     not null comment '角色名称',
		role_code       varchar(50)     not null comment '角色编码',
		description     varchar(200)    null comment '角色描述',
		status          int             default 1 not null comment '状态:0=禁用,1=正常',
		sort            int             default 0 not null comment '排序',
		menu_ids        text            null comment '菜单ID列表，逗号分隔',
		enterprise_id   bigint unsigned default 0 not null comment '企业ID',
		created_at      datetime        default CURRENT_TIMESTAMP not null comment '创建时间',
		updated_at      datetime        default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP comment '修改时间',
		deleted_at      datetime        null comment '软删除时间',
		create_by       bigint unsigned default '0' not null comment '创建人ID',
		create_by_name  varchar(20)     default '系统' not null comment '创建人姓名',
		update_by       bigint unsigned default '0' not null comment '修改人ID',
		update_by_name  varchar(20)     default '系统' not null comment '修改人姓名',
		UNIQUE KEY uk_role_code (role_code),
		INDEX idx_role_name (role_name),
		INDEX idx_enterprise_id (enterprise_id),
		INDEX idx_deleted_at (deleted_at)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';
	`

	err := gormDB.Exec(createTableSQL).Error
	if err != nil {
		t.Logf("创建表失败（可能表已存在）: %v", err)
	} else {
		t.Log("roles 表创建成功")
	}

	// 初始化默认角色数据
	initRolesSQL := `
	INSERT IGNORE INTO roles (id, role_name, role_code, description, status, sort) VALUES
	(1, '超级管理员', 'admin', '拥有系统所有权限', 1, 1),
	(2, '运营人员', 'operator', '负责日常运营工作', 1, 2),
	(3, '普通用户', 'user', '普通注册用户', 1, 3),
	(4, 'VIP用户', 'vip', '付费VIP用户', 1, 4),
	(5, '测试人员', 'tester', '负责测试工作', 1, 5);
	`

	err = gormDB.Exec(initRolesSQL).Error
	if err != nil {
		t.Logf("初始化角色数据失败: %v", err)
	} else {
		t.Log("默认角色数据初始化成功")
	}
}

// TestDropRolesTable 删除 roles 表（谨慎使用）
func TestDropRolesTable(t *testing.T) {
	gormDB := db.InitDb()

	err := gormDB.Exec("DROP TABLE IF EXISTS roles").Error
	if err != nil {
		t.Fatalf("删除表失败: %v", err)
	}
	t.Log("roles 表已删除")
}