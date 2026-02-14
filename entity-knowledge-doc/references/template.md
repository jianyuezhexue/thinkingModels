# 业务模型文档模板

以下是单个业务模型 README.md 的标准模板，根据实际模型情况填写各部分内容。

```markdown
# [模型中文名称]（[模型英文标识]）

## 模型概述

### 基本信息

| 项目 | 说明 |
|------|------|
| **模型名称** | [中文名] |
| **英文标识** | [英文名，如 Model/Category/Tag] |
| **所属领域** | [domain名称，如 thinkingModel/practice/master] |
| **对应表名** | [数据库表名] |
| **模型文件** | `domain/[领域]/[模型]/model.go` |

### 模型职责

[描述该业务模型负责的核心职责，1-2 句话概括]

### 业务场景

[列举该模型的典型使用场景]

---

## 实体定义

```go
type [ModelEntity] struct {
    base.BaseModel[[ModelEntity]]
    // 基础字段
    Name   string `json:"name" type:"db" comment:"名称"`
    Code   string `json:"code" type:"db" comment:"编码"`
    // ... 其他字段
}
```

---

## 字段说明

| 字段名 | Go类型 | 数据库类型 | JSON字段 | 约束 | 说明 |
|--------|--------|------------|----------|------|------|
| id | uint64 | bigint unsigned | id | 主键、自增 | 主键ID |
| name | string | varchar(255) | name | 非空 | 名称 |
| code | string | varchar(100) | code | 非空、唯一 | 编码 |
| ... | ... | ... | ... | ... | ... |
| created_at | db.LocalTime | datetime | createdAt | 可空 | 创建时间 |
| updated_at | db.LocalTime | datetime | updatedAt | 可空 | 更新时间 |
| deleted_at | db.LocalTime | datetime | deletedAt | 可空 | 删除时间（软删除） |

### 审计字段

所有模型自动包含以下审计字段（来自 base.BaseModel）：

| 字段名 | 类型 | 说明 |
|--------|------|------|
| created_at | LocalTime | 创建时间 |
| updated_at | LocalTime | 更新时间 |
| deleted_at | LocalTime | 删除时间（软删除标记） |

---

## 业务能力

### 常规能力（继承自 base.BaseModel）

| 能力 | 方法 | 说明 |
|------|------|------|
| 创建 | Create | 创建新记录 |
| 更新 | Update | 更新已有记录 |
| 删除 | Del/Delete | 软删除记录 |
| 详情查询 | LoadById | 根据ID查询详情 |
| 列表查询 | List | 分页列表查询 |
| 统计 | Count | 统计记录数量 |
| 数据校验 | Validate | 校验数据合法性 |
| 数据修复 | Repair | 修复数据默认值 |
| 数据完善 | Complete | 完善数据 |

### 定制化能力

| 能力 | 方法 | 说明 | 来源 |
|------|------|------|------|
| [能力1] | [MethodName] | [说明] | abilitys.go |
| [能力2] | [MethodName] | [说明] | abilitys.go |

---

## 关联关系

### 与 [其他模型A] 的关系

- **关系类型**：[一对一/一对多/多对多]
- **关联方式**：[通过外字段/中间表关联]
- **说明**：[描述关联逻辑]

### 与 [其他模型B] 的关系

...

---

## 数据库表结构

```sql
CREATE TABLE `[table_name]` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(255) NOT NULL COMMENT '名称',
  `code` varchar(100) NOT NULL COMMENT '编码',
  ...
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_name` (`name`),
  ...
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='[表中文说明]';
```

### 索引说明

| 索引名 | 字段 | 类型 | 说明 |
|--------|------|------|------|
| PRIMARY | id | 主键 | 自增主键 |
| uk_code | code | 唯一 | 编码唯一索引 |
| idx_name | name | 普通 | 名称索引 |

---

## 类型定义

相关 DTO 和请求/响应类型定义在 `types.go`：

```go
// CreateModel 创建请求
type CreateModel struct {
    Name string `json:"name" binding:"required"`
    Code string `json:"code" binding:"required"`
    ...
}

// ModelInfo 响应信息
type ModelInfo struct {
    Id   uint64 `json:"id"`
    Name string `json:"name"`
    ...
}
```

---

## 变更日志

| 日期 | 版本 | 变更内容 | 作者 |
|------|------|----------|------|
| YYYY-MM-DD | v1.0 | 初始版本 | [作者] |
```
