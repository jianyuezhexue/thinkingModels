package domain

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"thinkingModels/component/db"
	"gorm.io/gorm"
)

// 搜索条件
type SearchConditon = func(db *gorm.DB) *gorm.DB

// 充血模型基础接口
type BaseModelInterface[T any] interface {
	DB() *gorm.DB                                                                // 获取DB
	SetData(data any) error                                                      // 设置数据
	Create() error                                                               // 新增数据
	Update() (*T, error)                                                         // 更新数据
	MakeConditon(data any) SearchConditon                                        // 构造查询条件
	LoadData(SearchConditon) error                                               // 加载数据
	GetById(Id uint64) (*T, error)                                               // 根据Id查询数据
	Repair() error                                                               // 修复数据
	Count(SearchConditon) (int64, error)                                         // 统计数据条数
	List(SearchConditon) ([]*T, error)                                           // 查询列表数据
	Import(data any) error                                                       // 导入数据
	Export(data any) error                                                       // 导出数据
	Complete() error                                                             // 完善数据
	Del(ids ...uint64) error                                                     // 删除数据
	GetMaxId(tableName string) (uint64, error)                                   // 获取最大ID
	CheckBusinessNoRepeat(tableName, filedName, businessNo string) (bool, error) // 业务单号重复校验
	Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error
}

// 公共模型方法
type BaseModel[T any] struct {
	Id        uint64                   `json:"id" uri:"id" search:"exact" gorm:"primarykey"` // 主键
	CreatedAt db.LocalTime             `json:"createdAt" search:"gte"`                       // 创建时间
	UpdatedAt db.LocalTime             `json:"updatedAt" search:"lte"`                       // 更新时间
	DeletedAt gorm.DeletedAt           `json:"-" gorm:"index" search:"-"`                    // 删除标记
	Page      int64                    `json:"-"  uri:"page" gorm:"-" search:"page"`         // 分页
	PageSize  int64                    `json:"-"  uri:"pageSize" gorm:"-" search:"pageSize"` // 分页大小
	Db        *gorm.DB                 `json:"-" gorm:"-" search:"-"`                        // 数据库连接
	Ctx       *gin.Context             `json:"-" gorm:"-" search:"-"`                        // 上下文
	Preloads  map[string][]interface{} `json:"-" gorm:"-" search:"-"`                        // 预加载
	CurrTime  time.Time                `json:"-" gorm:"-" search:"-"`                        // 当前时间
}

// 获取最大ID
func (m *BaseModel[T]) GetMaxId(tableName string) (uint64, error) {
	var maxId uint64
	err := m.Db.Table(tableName).Select("COALESCE(max(id), 0)").Scan(&maxId).Error
	return maxId, err
}

// 业务单号重复校验
func (m *BaseModel[T]) CheckBusinessNoRepeat(tableName, filedName, businessNo string) (bool, error) {
	var count int64
	model := new(T)
	err := m.Db.Model(model).Where(fmt.Sprintf("%s = ?", filedName), businessNo).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ------------- DB相关函数 -----------------
// 清除分页和偏移量
func (b *BaseModel[T]) ClearOffset() SearchConditon {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Limit(-1).Offset(-1)
		return db
	}
}

// 获取Db
func (m *BaseModel[T]) DB() *gorm.DB {
	db, exist := m.Ctx.Get("txDb")
	if exist && db != nil {
		return db.(*gorm.DB)
	}
	return m.Db
}

// 开启事务
func (m *BaseModel[T]) Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	err := m.Db.Transaction(func(tx *gorm.DB) error {
		// 预埋事务Db
		m.Ctx.Set("txDb", tx)

		// 执行事务逻辑代码
		if err := fc(tx); err != nil {
			return err
		}

		// 回收事务Db
		m.Ctx.Set("txDb", nil)
		return nil
	}, opts...)
	return err
}
