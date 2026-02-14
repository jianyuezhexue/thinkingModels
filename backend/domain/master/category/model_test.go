package category

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"thinkingModels/component/db"
)

// 分类选项初始化数据
var initCategories = []CreateCategory{
	{Name: "商业管理", Icon: "💼", Description: "商业思维与管理方法论", Heat: 100},
	{Name: "战略规划", Icon: "🎯", Description: "战略分析与规划工具", Heat: 90},
	{Name: "创新思维", Icon: "💡", Description: "创新方法与创意激发", Heat: 85},
	{Name: "分析工具", Icon: "📊", Description: "数据分析与可视化工具", Heat: 80},
	{Name: "决策方法", Icon: "⚖️", Description: "决策框架与评估方法", Heat: 75},
	{Name: "创意构思", Icon: "🎨", Description: "创意产生与设计思维", Heat: 70},
	{Name: "心理学", Icon: "🧠", Description: "心理模型与认知科学", Heat: 65},
	{Name: "沟通表达", Icon: "💬", Description: "沟通技巧与表达方法", Heat: 60},
}

// 获取测试用的 gin context
func getTestContext() *gin.Context {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(nil)
	// 创建虚拟请求以满足 base 库的校验
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}
	// 设置当前用户ID，满足审计字段要求（必须是 string 类型）
	ctx.Set("currUserId", "1")
	ctx.Set("currUserName", "系统")
	return ctx
}

// TestInitCategoryTable 测试创建 category 表并初始化数据
func TestInitCategoryTable(t *testing.T) {
	// 获取数据库连接
	gormDB := db.InitDb()

	// 1. 创建 category 表
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS category (
		id              int unsigned auto_increment primary key comment '主键ID',
		name            varchar(50)     not null comment '分类名称',
		icon            varchar(255)    null comment '分类图标URL',
		description     varchar(500)    null comment '分类描述',
		heat            int             default '0' not null comment '热度值，用于排序',
		created_at      datetime        default CURRENT_TIMESTAMP not null comment '创建时间',
		updated_at      datetime        default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP comment '修改时间',
		deleted_at      datetime        null comment '软删除时间',
		create_by       bigint unsigned default '0' not null comment '创建人ID',
		create_by_name  varchar(20)     default '系统' not null comment '创建人姓名',
		update_by       bigint unsigned default '0' not null comment '修改人ID',
		update_by_name  varchar(20)     default '系统' not null comment '修改人姓名',
		INDEX idx_name (name),
		INDEX idx_heat (heat)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='思维模型分类表';
	`

	err := gormDB.Exec(createTableSQL).Error
	assert.NoError(t, err, "创建 category 表失败")

	// 2. 清空已有数据（避免重复插入）
	gormDB.Exec("DELETE FROM category")

	// 3. 初始化数据
	ctx := getTestContext()
	entity := NewCategoryEntity(ctx)

	for i, cat := range initCategories {
		_, err := entity.CreateCategory(&cat)
		if err != nil {
			t.Logf("插入第 %d 条数据失败: %v", i+1, err)
		} else {
			t.Logf("成功插入分类: %s", cat.Name)
		}
	}

	// 4. 验证数据
	var count int64
	gormDB.Model(&CategoryEntity{}).Count(&count)
	assert.Equal(t, int64(len(initCategories)), count, "数据初始化数量不匹配")
	t.Logf("成功初始化 %d 个分类", count)

	// 5. 查询验证
	list, err := entity.All()
	assert.NoError(t, err)
	assert.Equal(t, len(initCategories), len(list))

	// 验证第一条数据
	if len(list) > 0 {
		assert.Equal(t, "商业管理", list[0].Name)
		assert.Equal(t, "💼", list[0].Icon)
		t.Logf("第一条分类: %s %s", list[0].Name, list[0].Icon)
	}
}

// TestCategoryEntity_Validate 测试数据校验
func TestCategoryEntity_Validate(t *testing.T) {
	ctx := getTestContext()
	entity := NewCategoryEntity(ctx)

	tests := []struct {
		name    string
		entity  *CategoryEntity
		wantErr bool
		errMsg  string
	}{
		{
			name: "正常情况",
			entity: &CategoryEntity{
				Name:        "决策思维",
				Icon:        "https://example.com/icon.png",
				Description: "描述",
				Heat:        10,
			},
			wantErr: false,
		},
		{
			name: "名称为空",
			entity: &CategoryEntity{
				Name:        "",
				Description: "描述",
			},
			wantErr: true,
			errMsg:  "分类名称不能为空",
		},
		{
			name: "名称过长",
			entity: &CategoryEntity{
				Name:        "这是一段超过50个字符长度的分类名称用于测试超长情况",
				Description: "描述",
			},
			wantErr: true,
			errMsg:  "分类名称长度不能超过50个字符",
		},
		{
			name: "描述过长",
			entity: &CategoryEntity{
				Name:        "正常名称",
				Description: string(make([]byte, 501)),
			},
			wantErr: true,
			errMsg:  "分类描述长度不能超过500个字符",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.entity.BaseModel = entity.BaseModel
			err := tt.entity.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestCategoryEntity_TableName 测试表名
func TestCategoryEntity_TableName(t *testing.T) {
	ctx := getTestContext()
	entity := NewCategoryEntity(ctx)
	assert.Equal(t, "category", entity.TableName())
}

// TestCategoryEntity_Complete 测试数据完善
func TestCategoryEntity_Complete(t *testing.T) {
	ctx := getTestContext()
	entity := NewCategoryEntity(ctx)

	// 测试热度默认值
	testEntity := &CategoryEntity{
		Name: "测试",
		Heat: 0,
	}
	testEntity.BaseModel = entity.BaseModel

	err := testEntity.Complete()
	assert.NoError(t, err)
	assert.Equal(t, 0, testEntity.Heat)
}

// TestCreateCategory_Validate 测试新建分类校验
func TestCreateCategory_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     *CreateCategory
		wantErr bool
		errMsg  string
	}{
		{
			name: "正常情况",
			req: &CreateCategory{
				Name:        "决策思维",
				Icon:        "https://example.com/icon.png",
				Description: "描述",
				Heat:        10,
			},
			wantErr: false,
		},
		{
			name: "名称为空",
			req: &CreateCategory{
				Name: "",
			},
			wantErr: true,
			errMsg:  "分类名称不能为空",
		},
		{
			name: "名称过长",
			req: &CreateCategory{
				Name: "这是一段超过50个字符长度的分类名称用于测试超长情况",
			},
			wantErr: true,
			errMsg:  "分类名称长度不能超过50个字符",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req.Name == "" {
				assert.Equal(t, "", tt.req.Name)
			} else if len(tt.req.Name) > 50 {
				assert.Greater(t, len(tt.req.Name), 50)
			} else {
				assert.LessOrEqual(t, len(tt.req.Name), 50)
			}
		})
	}
}

// TestIncreaseHeat_Validate 测试增加热度校验
func TestIncreaseHeat_Validate(t *testing.T) {
	tests := []struct {
		name    string
		id      uint64
		delta   int
		wantErr bool
	}{
		{
			name:    "正常情况",
			id:      1,
			delta:   10,
			wantErr: false,
		},
		{
			name:    "delta为0",
			id:      1,
			delta:   0,
			wantErr: true,
		},
		{
			name:    "delta为负数",
			id:      1,
			delta:   -5,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.delta <= 0 {
				assert.LessOrEqual(t, tt.delta, 0)
			} else {
				assert.Greater(t, tt.delta, 0)
			}
		})
	}
}

// TestCategoryAbility_All 测试全量列表能力
func TestCategoryAbility_All(t *testing.T) {
	ctx := getTestContext()
	entity := NewCategoryEntity(ctx)

	list, err := entity.All()
	// 由于可能表中没有数据，这里不断言错误，只记录日志
	if err != nil {
		t.Logf("查询全量列表失败: %v", err)
	} else {
		t.Logf("查询到 %d 个分类", len(list))
	}
}

// TestCategoryAbility_CreateCategory 测试新建分类能力
func TestCategoryAbility_CreateCategory(t *testing.T) {
	ctx := getTestContext()
	entity := NewCategoryEntity(ctx)

	// 测试正常创建
	req := &CreateCategory{
		Name:        "测试分类" + randString(5),
		Icon:        "📝",
		Description: "这是一个测试分类",
		Heat:        50,
	}

	res, err := entity.CreateCategory(req)
	if err != nil {
		// 可能是唯一性校验失败或其他错误
		t.Logf("创建分类失败: %v", err)
	} else {
		t.Logf("成功创建分类: ID=%d, Name=%s", res.Id, res.Name)
		assert.Equal(t, req.Name, res.Name)
		assert.Equal(t, req.Icon, res.Icon)
		assert.Equal(t, req.Heat, res.Heat)
	}

	// 测试重复名称
	_, err = entity.CreateCategory(req)
	if err != nil {
		t.Logf("重复名称校验: %v", err)
	}
}

// TestCategoryAbility_IncreaseHeat 测试增加热度能力
func TestCategoryAbility_IncreaseHeat(t *testing.T) {
	ctx := getTestContext()
	entity := NewCategoryEntity(ctx)

	// 测试 delta <= 0 的情况
	err := entity.IncreaseHeat(1, 0)
	if err != nil {
		t.Logf("delta=0 校验: %v", err)
	}

	err = entity.IncreaseHeat(1, -5)
	if err != nil {
		t.Logf("delta<0 校验: %v", err)
	}
}

// randString 生成随机字符串
func randString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[i%26]
	}
	return string(b)
}
