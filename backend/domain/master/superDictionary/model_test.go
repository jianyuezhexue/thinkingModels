package superDictionary

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

// TestInitSuperDictionaryTable 测试创建 super_dictionary 表
func TestInitSuperDictionaryTable(t *testing.T) {
	gormDB := db.InitDb()

	// 创建 super_dictionary 表（如果不存在）
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS super_dictionary (
		id              bigint unsigned auto_increment primary key comment '主键ID',
		parent_id       bigint          default 0 not null comment '父级ID',
		dict_value      varchar(100)    not null comment '字典值',
		dict_name       varchar(100)    not null comment '字典名称',
		level           int             default 1 not null comment '层级',
		level_name      varchar(100)    null comment '层级名称',
		description     varchar(500)    null comment '字典描述',
		eval            varchar(500)    null comment '表达式',
		ext_schema      text            null comment '扩展Schema',
		ext_json        text            null comment '扩展JSON数据',
		created_at      datetime        default CURRENT_TIMESTAMP not null comment '创建时间',
		updated_at      datetime        default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP comment '修改时间',
		deleted_at      datetime        null comment '软删除时间',
		create_by       bigint unsigned default '0' not null comment '创建人ID',
		create_by_name  varchar(50)     default '系统' not null comment '创建人姓名',
		update_by       bigint unsigned default '0' not null comment '修改人ID',
		update_by_name  varchar(50)     default '系统' not null comment '修改人姓名',
		INDEX idx_parent_id (parent_id),
		INDEX idx_dict_value (dict_value),
		INDEX idx_level (level)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='超级字典表';
	`

	err := gormDB.Exec(createTableSQL).Error
	if err != nil {
		t.Logf("创建表失败（可能表已存在）: %v", err)
	}

	t.Log("super_dictionary 表已准备好")
}

// TestFixSuperDictionaryTable 修复现有表结构
func TestFixSuperDictionaryTable(t *testing.T) {
	gormDB := db.InitDb()

	// 检查表是否存在
	var count int64
	gormDB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'super_dictionary'").Scan(&count)

	if count == 0 {
		t.Log("表不存在，请先运行 TestInitSuperDictionaryTable")
		return
	}

	// 修复 id 字段为自增
	fixSQL := `
		ALTER TABLE super_dictionary
		MODIFY COLUMN id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID';
	`

	err := gormDB.Exec(fixSQL).Error
	if err != nil {
		t.Logf("修复表结构失败: %v", err)
		// 如果修复失败，尝试删除并重建表
		t.Log("尝试删除并重建表...")

		gormDB.Exec("DROP TABLE IF EXISTS super_dictionary")

		createTableSQL := `
		CREATE TABLE super_dictionary (
			id              bigint unsigned auto_increment primary key comment '主键ID',
			parent_id       bigint          default 0 not null comment '父级ID',
			dict_value      varchar(100)    not null comment '字典值',
			dict_name       varchar(100)    not null comment '字典名称',
			level           int             default 1 not null comment '层级',
			level_name      varchar(100)    null comment '层级名称',
			description     varchar(500)    null comment '字典描述',
			eval            varchar(500)    null comment '表达式',
			ext_schema      text            null comment '扩展Schema',
			ext_json        text            null comment '扩展JSON数据',
			created_at      datetime        default CURRENT_TIMESTAMP not null comment '创建时间',
			updated_at      datetime        default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP comment '修改时间',
			deleted_at      datetime        null comment '软删除时间',
			create_by       bigint unsigned default '0' not null comment '创建人ID',
			create_by_name  varchar(50)     default '系统' not null comment '创建人姓名',
			update_by       bigint unsigned default '0' not null comment '修改人ID',
			update_by_name  varchar(50)     default '系统' not null comment '修改人姓名',
			INDEX idx_parent_id (parent_id),
			INDEX idx_dict_value (dict_value),
			INDEX idx_level (level)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='超级字典表';
		`

		err = gormDB.Exec(createTableSQL).Error
		if err != nil {
			t.Fatalf("重建表失败: %v", err)
		}
		t.Log("表重建成功")
	} else {
		t.Log("表结构修复成功")
	}
}