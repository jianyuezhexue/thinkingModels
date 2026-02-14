# 字段分析指南

本指南说明如何从 Go 业务模型代码中提取字段信息。

## 字段定义分析

### Go 结构体字段

```go
type ModelEntity struct {
    base.BaseModel[ModelEntity]
    Name        string  `json:"name" type:"db" comment:"模型名称"`
    Code        string  `json:"code" type:"db" comment:"模型编码"`
    CategoryId  uint64  `json:"categoryId" type:"db" comment:"所属分类ID"`
    Price       float64 `json:"price" type:"db" comment:"价格"`
    Status      int     `json:"status" type:"db" comment:"状态"`
}
```

### Tag 解析

| Tag | 含义 | 示例 |
|-----|------|------|
| `json:"xxx"` | JSON 序列化名称 | `json:"categoryId"` |
| `type:"db"` | 数据库字段标记 | 表示该字段映射到数据库表 |
| `comment:"xxx"` | 字段中文说明 | `comment:"模型名称"` |
| `gorm:"primaryKey"` | 主键标记 | 主键字段 |
| `gorm:"index"` | 索引标记 | 普通索引 |
| `gorm:"uniqueIndex"` | 唯一索引标记 | 唯一索引 |

## Go 类型映射

| Go 类型 | MySQL 类型 | 说明 |
|---------|------------|------|
| uint64 | bigint unsigned | 主键ID、外键 |
| string | varchar(255) | 短文本 |
| string | text | 长文本（内容字段） |
| int | int | 整数 |
| int8 | tinyint | 小整数（状态、枚举） |
| float64 | decimal(10,2) | 金额、价格 |
| bool | tinyint(1) | 布尔值 |
| time.Time | datetime | 时间戳 |
| db.LocalTime | datetime | 本地时间（项目自定义类型） |

## 标准审计字段

所有继承 base.BaseModel 的实体自动包含：

```go
type BaseModel struct {
    Id        uint64     // 主键ID
    CreatedAt LocalTime  // 创建时间
    UpdatedAt LocalTime  // 更新时间
    DeletedAt LocalTime  // 删除时间（软删除）
}
```

对应 SQL：

```sql
`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
`created_at` datetime DEFAULT NULL COMMENT '创建时间',
`updated_at` datetime DEFAULT NULL COMMENT '更新时间',
`deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
```

## 字段约束提取

### 从 Validate 方法提取

```go
func (m *ModelEntity) Validate() error {
    if m.Name == "" {
        return errors.New("模型名称不能为空")
    }
    if m.Code == "" {
        return errors.New("模型编码不能为空")
    }
    if m.CategoryId == 0 {
        return errors.New("所属分类不能为空")
    }
    return nil
}
```

约束提取：
- `Name != ""` → 非空
- `Code != ""` → 非空
- `CategoryId != 0` → 非空

### 从 Repair 方法提取默认值

```go
func (m *ModelEntity) Repair() error {
    if m.Version == "" {
        m.Version = "1.0.0"
    }
    if m.Difficulty == 0 {
        m.Difficulty = 1
    }
    return nil
}
```

默认值提取：
- `Version` 默认 `"1.0.0"`
- `Difficulty` 默认 `1`

## 注意事项

1. **软删除**：通过 `deleted_at` 字段实现，查询时自动过滤
2. **时间字段**：使用 `db.LocalTime` 替代 `time.Time`
3. **外键关系**：通过 `xxx_id` 字段命名体现
4. **枚举字段**：使用 int 类型，通过注释说明取值含义
