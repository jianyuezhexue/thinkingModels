# 思维模型平台 MVP 开发计划

> 基于技术方案文档，聚焦 **thinking 领域（思维模型核心）** 的 MVP 功能开发

---

## 一、MVP 范围定义

### 1.1 核心功能闭环

```
用户 → 浏览/采纳模型 → 创建课题 → 选用模型分析 → AI辅助分析 → 导出行动 → 跟踪执行
```

### 1.2 MVP 实体清单

| 模块 | 实体 | 表名 | 开发状态 |
|------|------|------|----------|
| 思维模型 | ModelEntity | thinking_models | 🟡 部分完成 |
| 模型分类 | CategoryEntity | model_categories | 🟡 部分完成 |
| 课题管理 | TopicEntity | topics | 🟡 部分完成 |
| 分析记录 | AnalysisEntity | topic_analyses | 🟡 部分完成 |
| 行动管理 | ActionEntity | actions | 🔴 待开发 |
| 跟进记录 | FollowUpEntity | action_followups | 🔴 待开发 |
| 模型标签 | TagEntity | model_tags | 🔴 待开发 |

### 1.3 MVP 不包含

- 模型评论系统（model_comments）
- 付费购买流程（订单/支付）
- AI对话基础设施（ai领域）
- 社群功能（community领域）

---

## 二、现有代码盘点

### 2.1 已完成基础设施

| 模块 | 路径 | 状态 |
|------|------|------|
| 用户认证 | `domain/iam/user/` | ✅ 完成 |
| 超级字典 | `domain/master/superDictionary/` | ✅ 完成 |
| 数据库连接 | `component/db/` | ✅ 完成 |
| Redis | `component/redis/` | ✅ 完成 |
| 路由框架 | `router/` | ✅ 完成 |
| 中间件 | `middleware/` | ✅ 完成 |

### 2.2 thinking 领域现状

需要检查的现有代码：

```
backend/
├── domain/
│   ├── market/
│   │   ├── model/        # 思维模型 → 迁移到 thinking/model
│   │   └── category/     # 模型分类 → 迁移到 thinking/category
│   └── subject/
│       ├── topic/        # 课题管理 → 迁移到 thinking/topic
│       └── analysis/     # 分析记录 → 迁移到 thinking/analysis
├── api/
│   ├── market/           # → 迁移到 thinking/
│   └── subject/          # → 迁移到 thinking/
└── logic/
    ├── market/           # → 迁移到 thinking/
    └── subject/          # → 迁移到 thinking/
```

---

## 三、开发步骤

### 阶段一：代码结构重组（1天）

#### Step 1.1 创建 thinking 领域目录

```bash
mkdir -p backend/domain/thinking/{model,category,topic,analysis,action,followup}
mkdir -p backend/api/thinking
mkdir -p backend/logic/thinking
```

#### Step 1.2 迁移现有代码

| 源路径 | 目标路径 | 操作 |
|--------|----------|------|
| `domain/market/model/` | `domain/thinking/model/` | 移动并更新包名 |
| `domain/market/category/` | `domain/thinking/category/` | 移动并更新包名 |
| `domain/subject/topic/` | `domain/thinking/topic/` | 移动并更新包名 |
| `domain/subject/analysis/` | `domain/thinking/analysis/` | 移动并更新包名 |
| `api/market/` | `api/thinking/` | 移动并更新import |
| `api/subject/` | `api/thinking/` | 移动并更新import |
| `logic/market/` | `logic/thinking/` | 移动并更新import |
| `logic/subject/` | `logic/thinking/` | 移动并更新import |

#### Step 1.3 更新路由配置

修改 `router/v1.go`：
- `/market/model` → `/thinking/model`
- `/market/category` → `/thinking/category`
- `/subject/topic` → `/thinking/topic`
- `/subject/analysis` → `/thinking/analysis`

#### Step 1.4 清理旧目录

删除空的 `domain/market/`、`domain/subject/`、`api/market/`、`api/subject/` 等目录

---

### 阶段二：思维模型完善（2天）

#### Step 2.1 检查 ModelEntity 完整性

对照技术方案检查字段：

```go
// 需要确认的字段
type Model struct {
    // 基础信息
    Name, Code, Description, CoverImage, Icon string
    CategoryID uint64
    // 定价
    Price decimal.Decimal
    // 内容
    Content, Overview string
    // 属性
    Difficulty, EstimatedTime int
    // 统计
    UsageCount, AdoptCount, LikeCount, CommentCount int
    // 状态
    Status int  // 0草稿,1审核中,2已发布,3已下架,4审核拒绝
    PublishTime *time.Time
    Version string
    // 作者
    AuthorID uint64
    AuthorName string
    IsOfficial int
    SourceModelID uint64
}
```

#### Step 2.2 完善 ModelEntity 能力方法

```go
// 需要实现的充血模型方法
func (e *Entity) Validate() error           // 数据校验
func (e *Entity) Repair() error             // 数据修复
func (e *Entity) Publish() error            // 发布模型
func (e *Entity) Unpublish() error          // 下架模型
func (e *Entity) IncrementUsageCount()      // 增加使用次数
func (e *Entity) IncrementAdoptCount()      // 增加采纳次数
func (e *Entity) CalculateStats() error     // 计算统计数据
```

#### Step 2.3 完善模型接口

| 接口 | 方法 | 路径 | 状态 |
|------|------|------|------|
| 创建模型 | POST | /thinking/model | 检查 |
| 更新模型 | PUT | /thinking/model | 检查 |
| 模型详情 | GET | /thinking/model/:id | 检查 |
| 模型列表 | GET | /thinking/model/list | 检查 |
| 我的模型 | GET | /thinking/model/my | 检查 |
| 删除模型 | DELETE | /thinking/model | 检查 |
| 发布模型 | POST | /thinking/model/publish | 新增 |
| 下架模型 | POST | /thinking/model/unpublish | 新增 |
| 引用创建 | POST | /thinking/model/fork | 新增 |

---

### 阶段三：模型标签系统（1天）

#### Step 3.1 创建 TagEntity

新建文件：`domain/thinking/tag/model.go`

```go
type Tag struct {
    base.BaseModel[Tag]
    Name string           // 标签名称
    Code string           // 标签编码
    Description string    // 标签描述
    Color string          // 标签颜色（十六进制）
    Sort int              // 排序
    UseCount int          // 使用次数
    Status int            // 0禁用,1启用
}
```

#### Step 3.2 实现标签能力方法

```go
func (e *Entity) Validate() error
func (e *Entity) IncrementUseCount()
func (e *Entity) DecrementUseCount()
```

#### Step 3.3 创建标签接口

新建文件：`api/thinking/tag.go`、`logic/thinking/tag.go`

| 接口 | 方法 | 路径 |
|------|------|------|
| 创建标签 | POST | /thinking/tag |
| 更新标签 | PUT | /thinking/tag |
| 标签详情 | GET | /thinking/tag/:id |
| 标签列表 | GET | /thinking/tag/list |
| 删除标签 | DELETE | /thinking/tag/:id |
| 热门标签 | GET | /thinking/tag/hot |

#### Step 3.4 模型-标签关联

在模型接口中支持标签操作：

| 接口 | 方法 | 路径 |
|------|------|------|
| 给模型打标签 | POST | /thinking/model/tags |
| 移除模型标签 | DELETE | /thinking/model/tags |
| 按标签查模型 | GET | /thinking/model/by-tag/:tagId |

---

### 阶段四：模型分类完善（1天）

#### Step 3.1 检查 CategoryEntity 完整性

```go
type Category struct {
    ParentID uint64
    Name, Code, Icon, Description string
    Sort, Level int
    Path string
    Status int
    ModelCount int
}
```

#### Step 3.2 完善分类能力方法

```go
func (e *Entity) BuildPath() string         // 构建路径
func (e *Entity) GetChildren() []*Entity    // 获取子分类
func (e *Entity) UpdateModelCount() error   // 更新模型数量
func (e *Entity) Move(newParentID uint64)   // 移动分类
```

#### Step 3.3 完善分类接口

| 接口 | 方法 | 路径 | 状态 |
|------|------|------|------|
| 分类树 | GET | /thinking/category/tree | 检查 |
| 子分类 | GET | /thinking/category/children/:id | 新增 |
| 移动分类 | POST | /thinking/category/move | 新增 |

---

### 阶段五：课题管理完善（2天）

#### Step 4.1 检查 TopicEntity 完整性

```go
type Topic struct {
    Title, Description, Background, Goal, Constraints string
    Status int  // 0草稿,1进行中,2已完成,3已归档
    UserID uint64
    ModelID uint64
    ModelName string
    Priority int
    Tags string
    AnalysisCount, ActionCount int
    Deadline, CompleteTime *time.Time
}
```

#### Step 4.2 完善课题能力方法

```go
func (e *Entity) Validate() error
func (e *Entity) SelectModel(modelID uint64, modelName string) error
func (e *Entity) RemoveModel() error
func (e *Entity) UpdateStatus(status int) error
func (e *Entity) Complete() error
func (e *Entity) Archive() error
func (e *Entity) IncrementAnalysisCount()
func (e *Entity) IncrementActionCount()
```

#### Step 4.3 完善课题接口

| 接口 | 方法 | 路径 | 状态 |
|------|------|------|------|
| 创建课题 | POST | /thinking/topic | 检查 |
| 更新课题 | PUT | /thinking/topic | 检查 |
| 课题详情 | GET | /thinking/topic/:id | 检查 |
| 课题列表 | GET | /thinking/topic/list | 检查 |
| 我的课题 | GET | /thinking/topic/my | 检查 |
| 选用模型 | POST | /thinking/topic/select-model | 新增 |
| 移除模型 | POST | /thinking/topic/remove-model/:id | 新增 |
| 更新状态 | POST | /thinking/topic/status | 新增 |
| 完成课题 | POST | /thinking/topic/complete | 新增 |
| 归档课题 | POST | /thinking/topic/archive | 新增 |
| 课题统计 | GET | /thinking/topic/statistics | 新增 |

---

### 阶段六：分析记录完善（2天）

#### Step 5.1 检查 AnalysisEntity 完整性

```go
type Analysis struct {
    TopicID, ModelID uint64
    ModelName string
    Content string      // JSON格式用户填写内容
    AiAnalysis string   // AI分析结果
    AiSuggestions string
    Version int
    IsCurrent int       // 是否当前版本
    UserID uint64
    Status int          // 0分析中,1已完成,2失败
}
```

#### Step 5.2 完善分析能力方法

```go
func (e *Entity) Validate() error
func (e *Entity) SetAsCurrent() error
func (e *Entity) IncrementVersion() int
func (e *Entity) ParseContent() (map[string]interface{}, error)
func (e *Entity) GenerateAiPrompt(model *ModelEntity) string
func (e *Entity) SetAiResult(analysis, suggestions string) error
```

#### Step 5.3 完善分析接口

| 接口 | 方法 | 路径 | 状态 |
|------|------|------|------|
| 创建分析 | POST | /thinking/analysis | 检查 |
| AI分析 | POST | /thinking/analysis/save-with-ai | 新增（核心） |
| 更新分析 | PUT | /thinking/analysis | 检查 |
| 分析详情 | GET | /thinking/analysis/:id | 检查 |
| 课题分析列表 | GET | /thinking/analysis/by-topic/:topicId | 新增 |
| 当前版本 | GET | /thinking/analysis/current | 新增 |
| 历史版本 | GET | /thinking/analysis/history/:topicId/:modelId | 新增 |
| 设为当前 | POST | /thinking/analysis/set-current | 新增 |

---

### 阶段七：行动管理开发（3天）

#### Step 6.1 创建 ActionEntity

新建文件：`domain/thinking/action/model.go`

```go
type Action struct {
    base.BaseModel[Action]
    Title, Description string
    UserID, TopicID uint64
    TopicTitle string
    AnalysisID uint64
    Priority int        // 1低,2中,3高
    Status int          // 0待执行,1进行中,2已完成,3已取消
    Progress int        // 0-100
    Deadline, CompleteTime *time.Time
    GuidePrinciple string
    FollowupCount int
}
```

#### Step 6.2 实现行动能力方法

```go
func (e *Entity) Validate() error
func (e *Entity) UpdateProgress(progress int) error
func (e *Entity) Complete() error
func (e *Entity) Cancel() error
func (e *Entity) SetGuidePrinciple(principle string)
func (e *Entity) IncrementFollowupCount()
func (e *Entity) CheckOverdue() bool
```

#### Step 6.3 创建行动接口

新建文件：`api/thinking/action.go`、`logic/thinking/action.go`

| 接口 | 方法 | 路径 |
|------|------|------|
| 创建行动 | POST | /thinking/action |
| 从分析导出 | POST | /thinking/action/from-analysis |
| 更新行动 | PUT | /thinking/action |
| 行动详情 | GET | /thinking/action/:id |
| 行动列表 | GET | /thinking/action/list |
| 我的行动 | GET | /thinking/action/my |
| 删除行动 | DELETE | /thinking/action |
| 更新进度 | POST | /thinking/action/progress |
| 完成行动 | POST | /thinking/action/complete |
| 行动统计 | GET | /thinking/action/statistics |

---

### 阶段八：跟进记录开发（1天）

#### Step 7.1 创建 FollowUpEntity

新建文件：`domain/thinking/followup/model.go`

```go
type FollowUp struct {
    base.BaseModel[FollowUp]
    ActionID, UserID uint64
    Content string
    ProgressBefore, ProgressAfter int
}
```

#### Step 7.2 实现跟进能力方法

```go
func (e *Entity) Validate() error
func (e *Entity) SetProgressChange(before, after int)
```

#### Step 7.3 创建跟进接口

新建文件：`api/thinking/followup.go`、`logic/thinking/followup.go`

| 接口 | 方法 | 路径 |
|------|------|------|
| 添加跟进 | POST | /thinking/followup |
| 跟进详情 | GET | /thinking/followup/:id |
| 行动跟进列表 | GET | /thinking/followup/by-action/:actionId |
| 更新跟进 | PUT | /thinking/followup |
| 删除跟进 | DELETE | /thinking/followup/:id |

---

### 阶段九：数据库迁移（0.5天）

#### Step 8.1 执行建表语句

从技术方案中提取以下表的DDL并执行：

- [x] `thinking_models` - 检查现有表结构是否一致
- [x] `model_categories` - 检查现有表结构是否一致
- [x] `topics` - 检查现有表结构是否一致
- [x] `topic_analyses` - 检查现有表结构是否一致
- [ ] `model_tags` - 新建
- [ ] `actions` - 新建
- [ ] `action_followups` - 新建

#### Step 8.2 数据迁移（如需要）

检查是否有测试数据需要迁移

---

### 阶段十：集成测试（1天）

#### Step 9.1 接口测试

使用 Postman/curl 测试所有接口：

1. **模型流程**：创建模型 → 发布 → 查询列表
2. **课题流程**：创建课题 → 选用模型 → 创建分析 → AI分析
3. **行动流程**：从分析导出行动 → 更新进度 → 添加跟进 → 完成

#### Step 9.2 业务流程验收

完整走通核心流程：
```
浏览模型市场 → 采纳模型 → 创建课题 → 选用模型分析 → 
填写分析内容 → AI辅助分析 → 导出行动清单 → 
更新行动进度 → 添加跟进记录 → 完成行动
```

---

## 四、开发时间估算

| 阶段 | 内容 | 预估工时 |
|------|------|----------|
| 阶段一 | 代码结构重组 | 1天 |
| 阶段二 | 思维模型完善 | 2天 |
| 阶段三 | 模型标签系统 | 1天 |
| 阶段四 | 模型分类完善 | 1天 |
| 阶段五 | 课题管理完善 | 2天 |
| 阶段六 | 分析记录完善 | 2天 |
| 阶段七 | 行动管理开发 | 3天 |
| 阶段八 | 跟进记录开发 | 1天 |
| 阶段九 | 数据库迁移 | 0.5天 |
| 阶段十 | 集成测试 | 1天 |
| **合计** | | **14.5天** |

---

## 五、开发顺序依赖

```
阶段一（重组）
    │
    ├─→ 阶段二（模型）─→ 阶段三（标签）─→ 阶段四（分类）
    │
    └─→ 阶段五（课题）─→ 阶段六（分析）─→ 阶段七（行动）─→ 阶段八（跟进）
                                              │
                                              └─→ 阶段九（数据库）
                                                      │
                                                      └─→ 阶段十（测试）
```

---

## 六、风险与注意事项

### 6.1 代码迁移风险

- 包名修改后需要全局搜索替换 import 路径
- 路由变更后前端需要同步修改
- 确保 git 提交记录清晰，便于回滚

### 6.2 AI分析集成

- MVP 阶段可先使用固定的 AI 服务调用（如直接调用 OpenAI）
- AI 领域的完整基础设施可在后续阶段开发
- 需要预留 AI 调用的接口抽象

### 6.3 数据兼容性

- 检查现有测试数据是否符合新的字段定义
- 确保审计字段（create_by, update_by 等）正确填充

---

## 七、下一步行动

1. **确认现有代码状态** - 检查 market/subject 目录下的代码完成度
2. **确定迁移策略** - 是重写还是移动修改
3. **开始阶段一** - 代码结构重组
