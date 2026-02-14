# Category（模型分类）

## 摘要

Category 是思维模型平台的分类体系核心模型，用于对思维模型进行归类管理，支持用户在列表页按分类筛选和在创建思维模型时标记分类。

## 1. 业务模型名称和所属领域

- **模型名称**: CategoryEntity
- **所属领域**: master（主数据领域）
- **文件路径**: backend/domain/master/category/

## 2. 业务知识描述

### 2.1 业务场景

Category（模型分类）是思维模型平台的基础主数据，主要用于以下场景：

1. **思维模型列表筛选** - 在思维模型列表页面，用户可以按照分类进行筛选，快速定位特定类型的思维模型
2. **新建思维模型标记** - 用户在创建新的思维模型时，需要选择或新建所属分类，对模型进行归类标记
3. **分类导航浏览** - 支持通过分类体系浏览不同类别的思维模型，提升内容发现效率

### 2.2 核心价值

- **内容组织**: 为海量思维模型提供结构化的分类体系，便于内容组织和管理
- **用户体验**: 通过分类筛选和导航，帮助用户快速找到感兴趣的思维模型类型
- **扩展性**: 支持动态新建分类，适应不断新增的思维模型类型
- **热度导向**: 通过热度字段识别热门分类，优化内容推荐和排序策略

## 3. 数据表结构

### 3.1 表名

`category`

### 3.2 字段说明

| 字段名      | 类型           | 标签       | 描述           |
|------------|---------------|-----------|---------------|
| id         | int unsigned  | 主键ID     | 自增主键        |
| name       | varchar(50)   | 分类名称    | 分类的中文名称   |
| icon       | varchar(255)  | 图标       | 分类的图标URL    |
| description| varchar(500)  | 分类描述    | 分类的详细说明   |
| heat       | int           | 热度       | 分类的热门程度，用于排序 |
| created_at | datetime      | 创建时间    | 标准审计字段     |
| updated_at | datetime      | 修改时间    | 标准审计字段     |
| deleted_at | datetime      | 软删除时间  | 标准审计字段     |
| create_by  | bigint unsigned | 创建人ID  | 标准审计字段     |
| create_by_name | varchar(20) | 创建人姓名 | 标准审计字段    |
| update_by  | bigint unsigned | 修改人ID  | 标准审计字段     |
| update_by_name | varchar(20) | 修改人姓名 | 标准审计字段    |

### 3.3 建表 SQL

```sql
CREATE TABLE category (
    id              int unsigned auto_increment primary key comment '主键ID',

    -- 业务字段
    name            varchar(50)     not null comment '分类名称',
    icon            varchar(255)    null comment '分类图标URL',
    description     varchar(500)    null comment '分类描述',
    heat            int             default '0' not null comment '热度值，用于排序',

    -- 标准审计字段
    created_at      datetime        default CURRENT_TIMESTAMP not null comment '创建时间',
    updated_at      datetime        default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP comment '修改时间',
    deleted_at      datetime        null comment '软删除时间',
    create_by       bigint unsigned default '0' not null comment '创建人ID',
    create_by_name  varchar(20)     default '系统' not null comment '创建人姓名',
    update_by       bigint unsigned default '0' not null comment '修改人ID',
    update_by_name  varchar(20)     default '系统' not null comment '修改人姓名'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='思维模型分类表';

-- 索引建议
CREATE INDEX idx_name ON category(name);
CREATE INDEX idx_heat ON category(heat DESC);
```

## 4. 业务能力

### 4.1 常规业务能力（继承自 base.BaseModel）

| 能力 | 方法名 | 描述 |
|------|--------|------|
| 获取表名 | TableName() string | 返回数据表名称 "category" |
| 获取事务DB | Tx() *gorm.DB | 获取当前事务的数据库连接 |
| 事务处理 | Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error | 在事务中执行操作 |
| 设置数据 | SetData(data any) (*T, error) | 设置模型数据 |
| 数据校验 | Validate() error | 验证数据合法性（校验名称必填、长度等） |
| 完善数据 | Complete() error | 补充完善数据（如填充默认热度为0） |
| 新增数据 | Create() (*T, error) | 插入新分类记录 |
| 保存数据 | CreateWithData(data *T) (*T, error) | 使用传入对象保存 |
| 更新数据 | Update() (*T, error) | 更新当前分类数据 |
| 使用对象更新 | UpdateWithData(data *T) (*T, error) | 使用传入对象更新 |
| 加载数据 | LoadData(cond SearchCondition, preloads ...PreloadsType) (*T, error) | 根据条件加载单条分类 |
| 根据ID加载 | LoadById(id uint64, preloads ...PreloadsType) (*T, error) | 根据ID加载分类 |
| 根据业务编码加载 | LoadByBusinessCode(fieldName, fieldValue string, preloads ...PreloadsType) (*T, error) | 根据业务编码加载 |
| 查询单条 | GetById(Id uint64, preloads ...PreloadsType) (*T, error) | 根据ID查询单条分类 |
| 批量查询 | GetByIds(Ids []uint64, preloads ...PreloadsType) ([]*T, error) | 根据ID列表查询分类 |
| 修复数据 | Repair() error | 修复异常分类数据 |
| 统计条数 | Count(conds ...SearchCondition) (int64, error) | 统计满足条件的分类条数 |
| 查询列表 | List(conds ...SearchCondition) ([]*T, error) | 查询分类列表 |
| 根据ID列表查询 | ListByIds(Ids []uint64, preloads ...PreloadsType) ([]*T, error) | 根据ID列表查询分类 |
| 根据业务编码查询列表 | ListByBusinessCode(fieldName, fieldValue string, preloads ...PreloadsType) ([]*T, error) | 根据业务编码查询分类列表 |
| 根据业务编码统计 | CountByBusinessCode(fieldName, fieldValue string) (int64, error) | 根据业务编码统计分类数量 |
| 根据业务编码列表查询 | ListByBusinessCodes(fieldName string, fieldValues []string, preloads ...PreloadsType) ([]*T, error) | 根据业务编码列表查询 |
| 根据业务编码列表统计 | CountByBusinessCodes(fieldName string, fieldValues []string) (int64, error) | 根据业务编码列表统计数量 |
| 获取最大ID | MaxId() (int64, error) | 获取表中最大ID |
| 删除数据 | Del(ids ...uint64) error | 软删除分类 |
| 检查业务编码是否存在 | CheckBusinessCodeExist(fieldName, businessCode string) (bool, error) | 检查分类名称是否已存在 |
| 业务编码不能重复 | BusinessCodeCannotRepeat(fieldName, businessCode string) error | 检查分类名称不重复 |
| 批量检查业务编码 | CheckBusinessCodesExist(fieldName string, values []string) (map[int]bool, error) | 批量检查分类名称是否存在 |
| 检查唯一键重复 | CheckUniqueKeysExist(fieldNames []string, values []string) (bool, error) | 检查唯一键组合是否重复 |
| 批量检查唯一键 | CheckUniqueKeysExistBatch(fieldNames []string, values [][]string, withOutIds ...uint64) ([]bool, error) | 批量检查唯一键组合是否重复 |
| 构造查询条件 | MakeCondition(data any) func(db *gorm.DB) *gorm.DB | 根据数据构造查询条件 |
| 重置模型 | ReInit(entity *T, baseModel *BaseModel[T]) error | 重置模型中的Context和Db |
| 初始化状态机 | InitStateMachine(initStatus string, events []fsm.EventDesc, afterEvent fsm.Callback, callbacks ...map[string]fsm.Callback) error | 初始化状态机 |
| 执行状态事件 | EventExecution(initStatus, event, eventZhName string, args ...any) error | 执行状态机事件 |

### 4.2 定制化业务能力

| 能力 | 方法名 | 描述 |
|------|--------|------|
| 全量列表 | All() ([]*CategoryEntity, error) | 获取所有有效的分类列表（按热度降序），用于思维模型列表页筛选条件 |
| 新建分类 | CreateCategory(req CreateCategoryRequest) (*CategoryEntity, error) | 新建分类，自动校验名称唯一性，初始化热度为0 |
| 增加热度值 | IncreaseHeat(id uint64, delta int) error | 增加指定分类的热度值，用于热门分类排序 |

## 5. 关联关系

- **一对多关联**: 一个 Category 对应多个 ThinkingModel（思维模型）
- **关联字段**: `thinking_model.category_id` -> `category.id`

## 6. 业务规则

1. **名称唯一性**: 分类名称不能重复，新增和编辑时需校验
2. **热度默认值**: 新建分类时热度默认为 0
3. **排序规则**: 分类列表按热度降序排列，热度高的靠前显示
4. **软删除**: 删除分类时执行软删除，保留历史数据关联
5. **必填字段**: 分类名称为必填项，图标和描述为可选项
