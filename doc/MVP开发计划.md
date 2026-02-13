# 思维模型平台 MVP 开发计划

> 本文档定义 AI 编程助手执行 thinking 领域 MVP 开发的任务清单、执行动作、规范和检查标准

---

## 一、MVP 范围

### 1.1 实体清单

| 序号 | 模块 | 表名 | 当前状态 |
|------|------|------|----------|
| 1 | 思维模型 | thinking_models | 待检查 |
| 2 | 模型分类 | model_categories | 待检查 |
| 3 | 模型标签 | model_tags | 待新建 |
| 4 | 课题管理 | topics | 待检查 |
| 5 | 分析记录 | topic_analyses | 待检查 |
| 6 | 行动管理 | actions | 待新建 |
| 7 | 跟进记录 | action_followups | 待新建 |

### 1.2 开发流程（每个模块）

```
数据表检查 → 数据表调整/新建 → 后端代码实现 → 接口自测 → 前端联调
```

---

## 二、通用规范

### 2.1 数据表规范

**审计字段（必须包含）：**
```sql
created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',
create_by BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建人ID',
create_by_name VARCHAR(64) NOT NULL DEFAULT '' COMMENT '创建人姓名',
update_by BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新人ID',
update_by_name VARCHAR(64) NOT NULL DEFAULT '' COMMENT '更新人姓名'
```

**主键规范：**
- 使用 `BIGINT UNSIGNED AUTO_INCREMENT`
- 命名为 `id`

**索引规范：**
- 外键字段必须建索引
- 常用查询字段建索引
- 复合索引遵循最左前缀原则

### 2.2 后端代码规范

**目录结构：**
```
backend/
├── domain/thinking/{module}/
│   └── model.go          # Entity定义 + 充血方法
├── api/thinking/
│   └── {module}.go       # API层，参数校验
├── logic/thinking/
│   └── {module}.go       # Logic层，业务编排
└── router/v1.go          # 路由注册
```

**Entity 必须包含：**
```go
type Entity struct {
    base.BaseModel[Entity]
    // 业务字段...
}

func (e *Entity) TableName() string { return "表名" }
func (e *Entity) Validate() error   // 数据校验
func (e *Entity) Repair() error     // 数据修复
```

**API 层规范：**
- 使用 `tool/resp` 统一响应格式
- 参数使用 struct 定义，添加 binding tag
- 调用 logic 层处理业务

**Logic 层规范：**
- 继承 `logic.BaseLogic`
- 调用 domain 层方法
- 处理事务和业务编排

### 2.3 接口响应规范

```json
{
  "code": 0,
  "msg": "success",
  "data": {}
}
```

---

## 三、任务清单

---

### 任务1：思维模型（thinking_models）

#### 1.1 数据表检查

**执行动作：**
1. 读取 `backend/init.sql` 或数据库，获取 `thinking_models` 表结构
2. 对比技术方案中的字段定义
3. 输出差异清单

**检查标准：**
- [ ] 表名是否为 `thinking_models`
- [ ] 是否包含7个审计字段
- [ ] 字段类型是否与技术方案一致
- [ ] 索引是否完整（category_id, author_id, status）

**技术方案字段：**
```
id, name, code, description, cover_image, icon, category_id, 
price, content, overview, difficulty, estimated_time,
usage_count, adopt_count, like_count, comment_count,
status, publish_time, version, author_id, author_name,
is_official, source_model_id, + 7个审计字段
```

#### 1.2 数据表调整

**执行动作：**
1. 根据检查结果生成 ALTER TABLE 语句
2. 执行 SQL 或更新 init.sql
3. 确认表结构与技术方案一致

**检查标准：**
- [ ] 执行 `DESCRIBE thinking_models` 确认结构正确

#### 1.3 后端代码实现

**执行动作：**
1. 检查/创建 `domain/thinking/model/model.go`
   - Entity 结构体（字段与表一致）
   - TableName() 方法
   - Validate() 方法
   - Repair() 方法
   - Publish() / Unpublish() 方法
   - IncrementUsageCount() 等统计方法

2. 检查/创建 `api/thinking/model.go`
   - Create / Update / Delete / Get / List
   - GetMy（我的模型）
   - Publish / Unpublish
   - Fork（引用创建）

3. 检查/创建 `logic/thinking/model.go`
   - 对应 API 的业务逻辑

4. 更新 `router/v1.go` 注册路由

**检查标准：**
- [ ] Entity 字段与数据表一致
- [ ] 包含 Validate/Repair 方法
- [ ] API 层参数校验完整
- [ ] 路由已注册
- [ ] 编译通过 `go build`

#### 1.4 接口自测

**执行动作：**
1. 启动服务
2. 使用 curl 测试各接口
3. 验证响应格式和数据正确性

**测试用例：**
```bash
# 创建模型
curl -X POST localhost:8080/api/v1/thinking/model \
  -H "Content-Type: application/json" \
  -d '{"name":"测试模型","code":"test_model","category_id":1}'

# 模型列表
curl localhost:8080/api/v1/thinking/model/list

# 模型详情
curl localhost:8080/api/v1/thinking/model/1

# 发布模型
curl -X POST localhost:8080/api/v1/thinking/model/publish \
  -d '{"id":1}'
```

**检查标准：**
- [ ] 创建返回新记录ID
- [ ] 列表返回分页数据
- [ ] 详情返回完整字段
- [ ] 发布后 status 变为已发布

---

### 任务2：模型分类（model_categories）

#### 2.1 数据表检查

**执行动作：**
1. 获取 `model_categories` 表结构
2. 对比技术方案字段

**技术方案字段：**
```
id, parent_id, name, code, icon, description, 
sort, level, path, status, model_count, + 7个审计字段
```

**检查标准：**
- [ ] 包含父子关系字段（parent_id, level, path）
- [ ] 包含7个审计字段
- [ ] 索引：parent_id, status

#### 2.2 数据表调整

**执行动作：**
1. 生成并执行 ALTER TABLE 语句

#### 2.3 后端代码实现

**执行动作：**
1. `domain/thinking/category/model.go`
   - Entity + BuildPath() + GetChildren() + UpdateModelCount()

2. `api/thinking/category.go`
   - Tree / Children / Create / Update / Delete / Move

3. `logic/thinking/category.go`

4. 路由注册

**检查标准：**
- [ ] Tree 接口返回树形结构
- [ ] 移动分类后 path 自动更新

#### 2.4 接口自测

```bash
# 分类树
curl localhost:8080/api/v1/thinking/category/tree

# 子分类
curl localhost:8080/api/v1/thinking/category/children/1
```

---

### 任务3：模型标签（model_tags）

#### 3.1 数据表新建

**执行动作：**
1. 生成建表 SQL
2. 执行创建表

**建表语句：**
```sql
CREATE TABLE model_tags (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(32) NOT NULL DEFAULT '' COMMENT '标签名称',
    code VARCHAR(32) NOT NULL DEFAULT '' COMMENT '标签编码',
    description VARCHAR(255) NOT NULL DEFAULT '' COMMENT '标签描述',
    color VARCHAR(16) NOT NULL DEFAULT '' COMMENT '标签颜色',
    sort INT NOT NULL DEFAULT 0 COMMENT '排序',
    use_count INT NOT NULL DEFAULT 0 COMMENT '使用次数',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态 0禁用 1启用',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    create_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
    create_by_name VARCHAR(64) NOT NULL DEFAULT '',
    update_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
    update_by_name VARCHAR(64) NOT NULL DEFAULT '',
    UNIQUE KEY uk_code (code),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='模型标签表';
```

**检查标准：**
- [ ] 表创建成功
- [ ] 唯一索引 code 生效

#### 3.2 模型-标签关联表新建

**建表语句：**
```sql
CREATE TABLE model_tag_relations (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    model_id BIGINT UNSIGNED NOT NULL COMMENT '模型ID',
    tag_id BIGINT UNSIGNED NOT NULL COMMENT '标签ID',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    create_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
    UNIQUE KEY uk_model_tag (model_id, tag_id),
    INDEX idx_tag_id (tag_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='模型标签关联表';
```

#### 3.3 后端代码实现

**执行动作：**
1. `domain/thinking/tag/model.go`
   - TagEntity + TagRelationEntity
   - Validate() + IncrementUseCount() + DecrementUseCount()

2. `api/thinking/tag.go`
   - CRUD + Hot（热门标签）
   - AddToModel / RemoveFromModel / GetByModel

3. `logic/thinking/tag.go`

**检查标准：**
- [ ] 打标签时 use_count 自增
- [ ] 移除标签时 use_count 自减
- [ ] 热门标签按 use_count 排序

#### 3.4 接口自测

```bash
# 创建标签
curl -X POST localhost:8080/api/v1/thinking/tag \
  -d '{"name":"决策","code":"decision","color":"#FF5722"}'

# 给模型打标签
curl -X POST localhost:8080/api/v1/thinking/model/tags \
  -d '{"model_id":1,"tag_ids":[1,2,3]}'

# 热门标签
curl localhost:8080/api/v1/thinking/tag/hot?limit=10
```

---

### 任务4：课题管理（topics）

#### 4.1 数据表检查

**技术方案字段：**
```
id, title, description, background, goal, constraints,
status, user_id, model_id, model_name, priority, tags,
analysis_count, action_count, deadline, complete_time,
+ 7个审计字段
```

**检查标准：**
- [ ] 包含模型关联字段（model_id, model_name）
- [ ] 包含统计字段（analysis_count, action_count）
- [ ] 索引：user_id, model_id, status

#### 4.2 数据表调整

**执行动作：**
1. 对比现有表结构，生成 ALTER 语句

#### 4.3 后端代码实现

**执行动作：**
1. `domain/thinking/topic/model.go`
   - Entity + SelectModel() + RemoveModel() + Complete() + Archive()

2. `api/thinking/topic.go`
   - CRUD + GetMy + SelectModel + RemoveModel + UpdateStatus + Complete + Archive + Statistics

3. `logic/thinking/topic.go`

**检查标准：**
- [ ] 选用模型后 model_id 和 model_name 正确填充
- [ ] 完成课题时 complete_time 自动设置
- [ ] 统计接口返回各状态数量

#### 4.4 接口自测

```bash
# 创建课题
curl -X POST localhost:8080/api/v1/thinking/topic \
  -d '{"title":"职业发展规划","description":"分析未来3年职业方向"}'

# 选用模型
curl -X POST localhost:8080/api/v1/thinking/topic/select-model \
  -d '{"topic_id":1,"model_id":1}'

# 我的课题
curl localhost:8080/api/v1/thinking/topic/my

# 课题统计
curl localhost:8080/api/v1/thinking/topic/statistics
```

---

### 任务5：分析记录（topic_analyses）

#### 5.1 数据表检查

**技术方案字段：**
```
id, topic_id, model_id, model_name, content, ai_analysis,
ai_suggestions, version, is_current, user_id, status,
+ 7个审计字段
```

**检查标准：**
- [ ] content 字段为 TEXT 类型（存储 JSON）
- [ ] 包含版本管理字段（version, is_current）
- [ ] 索引：topic_id, model_id, user_id

#### 5.2 数据表调整

**执行动作：**
1. 确认字段类型正确
2. 补充缺失字段

#### 5.3 后端代码实现

**执行动作：**
1. `domain/thinking/analysis/model.go`
   - Entity + SetAsCurrent() + IncrementVersion() + ParseContent() + SetAiResult()

2. `api/thinking/analysis.go`
   - Create + Update + Get + ByTopic + Current + History + SetCurrent
   - **SaveWithAI**（核心：保存并调用 AI 分析）

3. `logic/thinking/analysis.go`
   - SaveWithAI 逻辑：保存内容 → 调用AI → 更新结果

**检查标准：**
- [ ] 创建分析时 version=1, is_current=1
- [ ] 新版本创建时旧版本 is_current=0
- [ ] SaveWithAI 返回 AI 分析结果

#### 5.4 接口自测

```bash
# 创建分析
curl -X POST localhost:8080/api/v1/thinking/analysis \
  -d '{"topic_id":1,"model_id":1,"content":"{\"step1\":\"xxx\"}"}'

# AI分析（核心）
curl -X POST localhost:8080/api/v1/thinking/analysis/save-with-ai \
  -d '{"topic_id":1,"model_id":1,"content":"{\"situation\":\"当前状态\"}"}'

# 当前版本
curl localhost:8080/api/v1/thinking/analysis/current?topic_id=1&model_id=1

# 历史版本
curl localhost:8080/api/v1/thinking/analysis/history/1/1
```

---

### 任务6：行动管理（actions）

#### 6.1 数据表新建

**建表语句：**
```sql
CREATE TABLE actions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(128) NOT NULL DEFAULT '' COMMENT '行动标题',
    description TEXT COMMENT '行动描述',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    topic_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '关联课题ID',
    topic_title VARCHAR(128) NOT NULL DEFAULT '' COMMENT '课题标题',
    analysis_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '来源分析ID',
    priority TINYINT NOT NULL DEFAULT 2 COMMENT '优先级 1低 2中 3高',
    status TINYINT NOT NULL DEFAULT 0 COMMENT '状态 0待执行 1进行中 2已完成 3已取消',
    progress TINYINT NOT NULL DEFAULT 0 COMMENT '进度 0-100',
    deadline TIMESTAMP NULL COMMENT '截止时间',
    complete_time TIMESTAMP NULL COMMENT '完成时间',
    guide_principle TEXT COMMENT '指导原则',
    followup_count INT NOT NULL DEFAULT 0 COMMENT '跟进记录数',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    create_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
    create_by_name VARCHAR(64) NOT NULL DEFAULT '',
    update_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
    update_by_name VARCHAR(64) NOT NULL DEFAULT '',
    INDEX idx_user_id (user_id),
    INDEX idx_topic_id (topic_id),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='行动表';
```

**检查标准：**
- [ ] 表创建成功
- [ ] 包含所有索引

#### 6.2 后端代码实现

**执行动作：**
1. `domain/thinking/action/model.go`
   - Entity + Validate() + UpdateProgress() + Complete() + Cancel() + CheckOverdue()

2. `api/thinking/action.go`
   - CRUD + GetMy + FromAnalysis + UpdateProgress + Complete + Statistics

3. `logic/thinking/action.go`
   - FromAnalysis：从分析结果批量导出行动

**检查标准：**
- [ ] 完成时自动设置 complete_time 和 status=2
- [ ] 进度更新时自动更新 status（100%→完成）
- [ ] FromAnalysis 能批量创建行动

#### 6.3 接口自测

```bash
# 从分析导出行动
curl -X POST localhost:8080/api/v1/thinking/action/from-analysis \
  -d '{"analysis_id":1,"actions":[{"title":"行动1"},{"title":"行动2"}]}'

# 更新进度
curl -X POST localhost:8080/api/v1/thinking/action/progress \
  -d '{"id":1,"progress":50}'

# 完成行动
curl -X POST localhost:8080/api/v1/thinking/action/complete \
  -d '{"id":1}'

# 我的行动
curl localhost:8080/api/v1/thinking/action/my?status=1

# 行动统计
curl localhost:8080/api/v1/thinking/action/statistics
```

---

### 任务7：跟进记录（action_followups）

#### 7.1 数据表新建

**建表语句：**
```sql
CREATE TABLE action_followups (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    action_id BIGINT UNSIGNED NOT NULL COMMENT '行动ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    content TEXT NOT NULL COMMENT '跟进内容',
    progress_before TINYINT NOT NULL DEFAULT 0 COMMENT '跟进前进度',
    progress_after TINYINT NOT NULL DEFAULT 0 COMMENT '跟进后进度',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    create_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
    create_by_name VARCHAR(64) NOT NULL DEFAULT '',
    update_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
    update_by_name VARCHAR(64) NOT NULL DEFAULT '',
    INDEX idx_action_id (action_id),
    INDEX idx_user_id (user_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='行动跟进记录表';
```

#### 7.2 后端代码实现

**执行动作：**
1. `domain/thinking/followup/model.go`
   - Entity + Validate() + SetProgressChange()

2. `api/thinking/followup.go`
   - Create + Update + Delete + Get + ByAction

3. `logic/thinking/followup.go`
   - 创建跟进时同步更新 action 的 progress 和 followup_count

**检查标准：**
- [ ] 创建跟进后 action.followup_count 自增
- [ ] 创建跟进后 action.progress 更新为 progress_after
- [ ] 删除跟进后 followup_count 自减

#### 7.3 接口自测

```bash
# 添加跟进
curl -X POST localhost:8080/api/v1/thinking/followup \
  -d '{"action_id":1,"content":"完成第一阶段","progress_after":30}'

# 行动的跟进列表
curl localhost:8080/api/v1/thinking/followup/by-action/1
```

---

## 四、执行顺序

```
任务1（模型）──┬── 任务2（分类）
              └── 任务3（标签）
                      │
任务4（课题）─────────┘
       │
任务5（分析）
       │
任务6（行动）
       │
任务7（跟进）
```

**并行可能性：**
- 任务1、2、3、4 可并行（无依赖）
- 任务5 依赖任务4（课题）
- 任务6 依赖任务5（分析）
- 任务7 依赖任务6（行动）

---

## 五、总检查清单

### 数据表检查

- [ ] thinking_models 结构正确
- [ ] model_categories 结构正确
- [ ] model_tags 已创建
- [ ] model_tag_relations 已创建
- [ ] topics 结构正确
- [ ] topic_analyses 结构正确
- [ ] actions 已创建
- [ ] action_followups 已创建

### 后端代码检查

- [ ] domain/thinking/ 包含 7 个模块
- [ ] api/thinking/ 包含 7 个接口文件
- [ ] logic/thinking/ 包含 7 个逻辑文件
- [ ] router/v1.go 注册所有路由
- [ ] `go build` 编译通过
- [ ] `go vet` 无警告

### 接口测试检查

- [ ] 模型 CRUD + 发布/下架 正常
- [ ] 分类树 + 移动 正常
- [ ] 标签 CRUD + 打标签 正常
- [ ] 课题 CRUD + 选模型 正常
- [ ] 分析 CRUD + AI分析 正常
- [ ] 行动 CRUD + 从分析导出 正常
- [ ] 跟进 CRUD + 同步更新行动 正常

---

## 六、开始执行

准备好后，请告诉我 **"开始任务1"** 或指定任务编号，我将按照上述流程执行。
