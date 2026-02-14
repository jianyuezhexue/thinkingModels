# 业务模型文档示例

以下是一个完整的单个业务模型 README.md 示例（以 thinkingModel 领域的 Model 为例）：

```markdown
# 思维模型实体（Model）

## 模型概述

### 基本信息

| 项目 | 说明 |
|------|------|
| **模型名称** | 思维模型 |
| **英文标识** | Model |
| **所属领域** | thinkingModel |
| **对应表名** | `thinking_models` |
| **模型文件** | `domain/thinkingModel/model/model.go` |

### 模型职责

管理平台级思维模型的元数据，包括模型名称、编码、分类、内容、价格、状态等信息。为用户的课题分析提供可选用的标准化思维模型。

### 业务场景

1. **模型创建**：作者创建新的思维模型，填写基本信息和内容
2. **模型发布**：将草稿状态的模型发布到平台供用户选用
3. **模型选用**：用户在创建课题时从模型库中选择思维模型
4. **模型更新**：修改模型信息、内容或状态

---

## 实体定义

```go
// ModelEntity 思维模型实体
type ModelEntity struct {
    base.BaseModel[ModelEntity]
    // 基础信息
    Name          string  `json:"name" type:"db" comment:"模型名称"`
    Code          string  `json:"code" type:"db" comment:"模型编码"`
    Description   string  `json:"description" type:"db" comment:"模型描述"`
    CoverImage    string  `json:"coverImage" type:"db" comment:"封面图片URL"`
    Icon          string  `json:"icon" type:"db" comment:"模型图标"`
    // 分类与价格
    CategoryId    uint64  `json:"categoryId" type:"db" comment:"所属分类ID"`
    Price         float64 `json:"price" type:"db" comment:"价格，0表示免费"`
    // 模型内容
    Content       string  `json:"content" type:"db" comment:"模型内容JSON"`
    // 难度与时间
    Difficulty    int     `json:"difficulty" type:"db" comment:"难度: 1=入门, 2=进阶, 3=高级"`
    EstimatedTime int     `json:"estimatedTime" type:"db" comment:"预计完成时间（分钟）"`
    // 统计信息
    UsageCount    int     `json:"usageCount" type:"db" comment:"使用次数统计"`
    // 状态管理
    Status        int     `json:"status" type:"db" comment:"状态: 0=草稿, 1=已发布, 2=已下架"`
    PublishTime   db.LocalTime `json:"publishTime" type:"db" comment:"发布时间"`
    // 版本信息
    Version       string  `json:"version" type:"db" comment:"版本号"`
    // 作者信息
    AuthorId      uint64  `json:"authorId" type:"db" comment:"作者ID"`
    AuthorName    string  `json:"authorName" type:"db" comment:"作者名称"`
}
```

---

## 字段说明

| 字段名 | Go类型 | 数据库类型 | JSON字段 | 约束 | 说明 |
|--------|--------|------------|----------|------|------|
| id | uint64 | bigint unsigned | id | 主键、自增 | 主键ID |
| name | string | varchar(255) | name | 非空 | 模型名称 |
| code | string | varchar(100) | code | 非空、唯一 | 模型编码 |
| description | string | text | description | 可空 | 模型描述 |
| cover_image | string | varchar(500) | coverImage | 可空 | 封面图片URL |
| icon | string | varchar(255) | icon | 可空 | 模型图标 |
| category_id | uint64 | bigint unsigned | categoryId | 外键、默认0 | 所属分类ID |
| price | float64 | decimal(10,2) | price | 默认0.00 | 价格，0表示免费 |
| content | string | json | content | 非空 | 模型内容JSON |
| difficulty | int | tinyint | difficulty | 默认1 | 难度：1=入门，2=进阶，3=高级 |
| estimated_time | int | int | estimatedTime | 默认30 | 预计完成时间（分钟） |
| usage_count | int | int | usageCount | 默认0 | 使用次数统计 |
| status | int | tinyint | status | 默认0 | 状态：0=草稿，1=已发布，2=已下架 |
| publish_time | db.LocalTime | datetime | publishTime | 可空 | 发布时间 |
| version | string | varchar(20) | version | 默认"1.0.0" | 版本号 |
| author_id | uint64 | bigint unsigned | authorId | 外键、默认0 | 作者ID |
| author_name | string | varchar(255) | authorName | 可空 | 作者名称 |
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
| 删除 | Del | 软删除记录 |
| 详情查询 | LoadById | 根据ID查询详情 |
| 列表查询 | List | 分页列表查询 |
| 统计 | Count | 统计记录数量 |
| 数据校验 | Validate | 校验数据合法性 |
| 数据修复 | Repair | 修复数据默认值 |
| 数据完善 | Complete | 完善数据 |

### 定制化能力

| 能力 | 方法 | 说明 | 来源 |
|------|------|------|------|
| 发布模型 | Publish | 将模型状态设为已发布，设置发布时间 | abilitys.go |
| 下架模型 | Unpublish | 将模型状态设为已下架 | abilitys.go |
| 更新内容 | UpdateContent | 更新模型内容JSON | abilitys.go |
| 增加使用次数 | IncrementUsageCount | 使用次数+1 | abilitys.go |

---

## 关联关系

### 与 Category 的关系

- **关系类型**：多对一
- **关联方式**：通过外键字段 `category_id` 关联
- **说明**：一个模型属于一个分类，一个分类可以包含多个模型

### 与 Topic（Subject 领域）的关系

- **关系类型**：一对多
- **关联方式**：Topic 通过 `model_id` 外键关联本模型
- **说明**：一个模型可以被多个课题选用

---

## 数据库表结构

```sql
CREATE TABLE `thinking_models` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(255) NOT NULL COMMENT '模型名称',
  `code` varchar(100) NOT NULL COMMENT '模型编码',
  `description` text COMMENT '模型描述',
  `cover_image` varchar(500) DEFAULT NULL COMMENT '封面图片URL',
  `icon` varchar(255) DEFAULT NULL COMMENT '模型图标',
  `category_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '所属分类ID',
  `price` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '价格，0表示免费',
  `content` json NOT NULL COMMENT '模型内容JSON',
  `difficulty` tinyint NOT NULL DEFAULT '1' COMMENT '难度：1=入门,2=进阶,3=高级',
  `estimated_time` int NOT NULL DEFAULT '30' COMMENT '预计完成时间（分钟）',
  `usage_count` int NOT NULL DEFAULT '0' COMMENT '使用次数统计',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态：0=草稿,1=已发布,2=已下架',
  `publish_time` datetime DEFAULT NULL COMMENT '发布时间',
  `version` varchar(20) NOT NULL DEFAULT '1.0.0' COMMENT '版本号',
  `author_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '作者ID',
  `author_name` varchar(255) DEFAULT NULL COMMENT '作者名称',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_status` (`status`),
  KEY `idx_author_id` (`author_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='思维模型表';
```

### 索引说明

| 索引名 | 字段 | 类型 | 说明 |
|--------|------|------|------|
| PRIMARY | id | 主键 | 自增主键 |
| uk_code | code | 唯一 | 编码唯一索引 |
| idx_category_id | category_id | 普通 | 分类ID索引 |
| idx_status | status | 普通 | 状态索引 |
| idx_author_id | author_id | 普通 | 作者ID索引 |
| idx_created_at | created_at | 普通 | 创建时间索引 |

---

## 类型定义

相关 DTO 和请求/响应类型定义在 `types.go`：

```go
// CreateModel 创建思维模型请求
type CreateModel struct {
    Name          string  `json:"name" binding:"required"`
    Code          string  `json:"code" binding:"required"`
    Description   string  `json:"description"`
    CoverImage    string  `json:"coverImage"`
    Icon          string  `json:"icon"`
    CategoryId    uint64  `json:"categoryId" binding:"required"`
    Price         float64 `json:"price"`
    Content       string  `json:"content" binding:"required"`
    Difficulty    int     `json:"difficulty"`
    EstimatedTime int     `json:"estimatedTime"`
}

// UpdateModel 更新思维模型请求
type UpdateModel struct {
    Id            uint64  `json:"id" binding:"required"`
    Name          string  `json:"name"`
    Code          string  `json:"code"`
    Description   string  `json:"description"`
    CoverImage    string  `json:"coverImage"`
    Icon          string  `json:"icon"`
    CategoryId    uint64  `json:"categoryId"`
    Price         float64 `json:"price"`
    Content       string  `json:"content"`
    Difficulty    int     `json:"difficulty"`
    EstimatedTime int     `json:"estimatedTime"`
}

// ModelInfo 思维模型信息响应
type ModelInfo struct {
    Id            uint64  `json:"id"`
    Name          string  `json:"name"`
    Code          string  `json:"code"`
    Description   string  `json:"description"`
    CoverImage    string  `json:"coverImage"`
    Icon          string  `json:"icon"`
    CategoryId    uint64  `json:"categoryId"`
    Price         float64 `json:"price"`
    Difficulty    int     `json:"difficulty"`
    EstimatedTime int     `json:"estimatedTime"`
    UsageCount    int     `json:"usageCount"`
    Status        int     `json:"status"`
    Version       string  `json:"version"`
    AuthorId      uint64  `json:"authorId"`
    AuthorName    string  `json:"authorName"`
    CreatedAt     string  `json:"createdAt"`
    UpdatedAt     string  `json:"updatedAt"`
}

// SearchModel 查询思维模型请求
type SearchModel struct {
    Name       string  `json:"name"`
    Code       string  `json:"code"`
    CategoryId uint64  `json:"categoryId"`
    Status     int     `json:"status"`
    Page       int     `json:"page"`
    PageSize   int     `json:"pageSize"`
}
```

---

## 变更日志

| 日期 | 版本 | 变更内容 | 作者 |
|------|------|----------|------|
| 2024-02-14 | v1.0 | 初始版本 | Claude |
```
