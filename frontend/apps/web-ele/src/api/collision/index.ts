export namespace CollisionApi {
  /** 用户信息 */
  export interface UserInfo {
    id: string;
    name: string;
    avatar: string;
    bio?: string;
    interests?: string[];
  }

  /** 话题分类 */
  export type Category = 'inspiration' | 'methodology' | 'case' | 'question' | 'share';

  // ==================== 找人聊聊 Meetup ====================
  /** 约见状态 */
  export type MeetupStatus = 'open' | 'pending' | 'confirmed' | 'completed' | 'cancelled';

  /** 费用承担方式 */
  export type CostSplit = 'host' | 'aa' | 'guest';

  /** 约见主题类型 */
  export type MeetupTopic = 'career' | 'startup' | 'technology' | 'investment' | 'life' | 'other';

  /** 约见信息 */
  export interface Meetup {
    id: string;
    title: string;
    topic: MeetupTopic;
    description: string;
    thoughts: string; // 发起人的基本思考
    host: UserInfo;
    modelId?: string;
    modelName?: string;
    tags: string[];
    city: string;
    location?: string; // 具体地点（确认后显示）
    preferredTime: string; // 期望时间描述
    scheduledTime?: string; // 确定的时间
    costSplit: CostSplit;
    maxGuests: number; // 最多约几人
    currentGuests: number; // 已报名人数
    status: MeetupStatus;
    viewCount: number;
    interestedCount: number; // 感兴趣人数
    isInterested?: boolean;
    createdAt: string;
    updatedAt: string;
  }

  /** 约见申请 */
  export interface MeetupApplication {
    id: string;
    meetupId: string;
    applicant: UserInfo;
    message: string; // 申请留言
    status: 'pending' | 'accepted' | 'rejected';
    createdAt: string;
  }

  /** 约见列表查询参数 */
  export interface MeetupListParams {
    page?: number;
    pageSize?: number;
    topic?: MeetupTopic | 'all';
    city?: string;
    costSplit?: CostSplit | 'all';
    status?: MeetupStatus | 'all';
    keyword?: string;
    sortBy?: 'latest' | 'popular' | 'soonest';
  }

  /** 约见列表响应 */
  export interface MeetupListResult {
    list: Meetup[];
    total: number;
    page: number;
    pageSize: number;
  }

  /** 创建约见参数 */
  export interface CreateMeetupParams {
    title: string;
    topic: MeetupTopic;
    description: string;
    thoughts: string;
    city: string;
    preferredTime: string;
    costSplit: CostSplit;
    maxGuests: number;
    tags?: string[];
    modelId?: string;
  }

  /** 申请约见参数 */
  export interface ApplyMeetupParams {
    meetupId: string;
    message: string;
  }

  // ==================== 付费咨询 Consultation ====================
  /** 咨询状态 */
  export type ConsultationStatus = 'open' | 'matched' | 'inProgress' | 'completed' | 'cancelled' | 'expired';

  /** 咨询领域 */
  export type ConsultationField = 'career' | 'startup' | 'technology' | 'product' | 'investment' | 'management' | 'psychology' | 'other';

  /** 咨询方式 */
  export type ConsultationMode = 'online' | 'offline' | 'both';

  /** 咨询需求 */
  export interface Consultation {
    id: string;
    title: string;
    field: ConsultationField;
    description: string;
    background: string; // 问题背景
    expectation: string; // 期望获得的帮助
    requester: UserInfo;
    modelId?: string;
    modelName?: string;
    tags: string[];
    reward: number; // 悬赏金额（元）
    mode: ConsultationMode;
    city?: string; // 线下时需要
    deadline: string; // 截止日期
    status: ConsultationStatus;
    viewCount: number;
    applicationCount: number; // 申请人数
    selectedExpert?: UserInfo; // 选中的专家
    createdAt: string;
    updatedAt: string;
  }

  /** 专家信息（扩展用户信息） */
  export interface Expert extends UserInfo {
    title?: string; // 职称/头衔
    company?: string;
    experience: number; // 从业年限
    expertise: string[]; // 擅长领域
    consultCount: number; // 已完成咨询数
    rating: number; // 评分
    responseRate: number; // 响应率
    hourlyRate?: number; // 时薪参考
  }

  /** 咨询申请 */
  export interface ConsultationApplication {
    id: string;
    consultationId: string;
    expert: Expert;
    proposal: string; // 申请方案
    estimatedTime: string; // 预计所需时间
    quotation?: number; // 报价（可选，可能接受悬赏价）
    status: 'pending' | 'accepted' | 'rejected';
    createdAt: string;
  }

  /** 咨询列表查询参数 */
  export interface ConsultationListParams {
    page?: number;
    pageSize?: number;
    field?: ConsultationField | 'all';
    mode?: ConsultationMode | 'all';
    status?: ConsultationStatus | 'all';
    minReward?: number;
    maxReward?: number;
    keyword?: string;
    sortBy?: 'latest' | 'reward' | 'deadline' | 'popular';
  }

  /** 咨询列表响应 */
  export interface ConsultationListResult {
    list: Consultation[];
    total: number;
    page: number;
    pageSize: number;
  }

  /** 创建咨询参数 */
  export interface CreateConsultationParams {
    title: string;
    field: ConsultationField;
    description: string;
    background: string;
    expectation: string;
    reward: number;
    mode: ConsultationMode;
    city?: string;
    deadline: string;
    tags?: string[];
    modelId?: string;
  }

  /** 申请咨询参数 */
  export interface ApplyConsultationParams {
    consultationId: string;
    proposal: string;
    estimatedTime: string;
    quotation?: number;
  }

  /** 话题 */
  export interface Discussion {
    id: string;
    title: string;
    content: string;
    summary: string;
    cover?: string;
    user: UserInfo;
    modelId?: string;
    modelName?: string;
    category: Category;
    tags: string[];
    viewCount: number;
    likeCount: number;
    commentCount: number;
    favoriteCount: number;
    isTop: boolean;
    isFeatured: boolean;
    isLiked?: boolean;
    isFavorited?: boolean;
    status: number; // 0:草稿 1:已发布 2:已关闭
    publishTime: string;
    createdAt: string;
    updatedAt: string;
  }

  /** 评论 */
  export interface Comment {
    id: string;
    discussionId: string;
    parentId?: string;
    replyToId?: string;
    replyToName?: string;
    content: string;
    user: UserInfo;
    likeCount: number;
    replyCount: number;
    isLiked?: boolean;
    createdAt: string;
    replies?: Comment[];
  }

  /** 话题列表查询参数 */
  export interface DiscussionListParams {
    page?: number;
    pageSize?: number;
    category?: Category | 'all';
    keyword?: string;
    modelId?: string;
    userId?: string;
    sortBy?: 'latest' | 'popular' | 'mostCommented' | 'mostLiked';
  }

  /** 话题列表响应 */
  export interface DiscussionListResult {
    list: Discussion[];
    total: number;
    page: number;
    pageSize: number;
  }

  /** 创建话题参数 */
  export interface CreateDiscussionParams {
    title: string;
    content: string;
    summary?: string;
    cover?: string;
    category: Category;
    tags?: string[];
    modelId?: string;
    status?: number; // 0:草稿 1:发布
  }

  /** 更新话题参数 */
  export interface UpdateDiscussionParams extends Partial<CreateDiscussionParams> {
    id: string;
  }

  /** 评论列表查询参数 */
  export interface CommentListParams {
    discussionId: string;
    parentId?: string;
    page?: number;
    pageSize?: number;
  }

  /** 评论列表响应 */
  export interface CommentListResult {
    list: Comment[];
    total: number;
    page: number;
    pageSize: number;
  }

  /** 创建评论参数 */
  export interface CreateCommentParams {
    discussionId: string;
    parentId?: string;
    replyToId?: string;
    content: string;
  }
}

// ==================== 模拟数据 ====================
const mockUsers: CollisionApi.UserInfo[] = [
  { id: '1', name: '思维探索者', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=explorer' },
  { id: '2', name: '创新先锋', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=pioneer' },
  { id: '3', name: '逻辑大师', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=logic' },
  { id: '4', name: '灵感捕手', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=catch' },
  { id: '5', name: '问题终结者', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=solver' },
];

const mockDiscussions: CollisionApi.Discussion[] = [
  {
    id: '1',
    title: '如何用第一性原理思考职业发展？',
    content: `最近在思考职业发展的问题，尝试用第一性原理来分析。

传统观念告诉我们要按部就班：读书 → 找工作 → 升职加薪。但如果回归本质，职业发展的目的是什么？

我的思考是：
1. **经济独立** - 满足基本生活需求
2. **价值实现** - 发挥个人能力，创造社会价值
3. **持续成长** - 不断学习，保持竞争力

从这三个基本原则出发，传统的职业路径并不是唯一选择。自由职业、创业、甚至躺平一段时间专注学习，都可能是更优解。

想问问大家：
- 你们是如何定义职业成功的？
- 有没有尝试过非传统的职业路径？`,
    summary: '尝试用第一性原理分析职业发展，探讨传统职业路径之外的可能性',
    cover: 'https://images.unsplash.com/photo-1507679799987-c73779587ccf?w=800',
    user: mockUsers[0]!,
    modelId: '3',
    modelName: '第一性原理思维',
    category: 'methodology',
    tags: ['职业发展', '第一性原理', '人生规划'],
    viewCount: 1234,
    likeCount: 89,
    commentCount: 23,
    favoriteCount: 45,
    isTop: true,
    isFeatured: true,
    status: 1,
    publishTime: '2024-02-10T10:30:00Z',
    createdAt: '2024-02-10T10:30:00Z',
    updatedAt: '2024-02-12T15:20:00Z',
  },
  {
    id: '2',
    title: '💡 灵感分享：用 SWOT 分析法评估创业想法',
    content: `最近有个创业想法，用 SWOT 分析了一下，分享给大家。

**项目：社区团购小程序**

**优势 (Strengths):**
- 技术团队经验丰富
- 有社区资源和人脉
- 启动成本相对较低

**劣势 (Weaknesses):**
- 没有供应链经验
- 品牌知名度为零
- 资金有限

**机会 (Opportunities):**
- 社区经济崛起
- 疫情改变消费习惯
- 本地化服务需求增加

**威胁 (Threats):**
- 巨头已经入场
- 价格战激烈
- 用户忠诚度低

分析完感觉挑战不小，但也看到了一些差异化机会。大家觉得这个方向可行吗？`,
    summary: '用 SWOT 分析社区团购创业项目的可行性',
    user: mockUsers[1]!,
    modelId: '1',
    modelName: 'SWOT 分析模型',
    category: 'case',
    tags: ['创业', 'SWOT', '社区团购', '案例分析'],
    viewCount: 856,
    likeCount: 56,
    commentCount: 34,
    favoriteCount: 28,
    isTop: false,
    isFeatured: true,
    status: 1,
    publishTime: '2024-02-09T14:20:00Z',
    createdAt: '2024-02-09T14:20:00Z',
    updatedAt: '2024-02-11T09:15:00Z',
  },
  {
    id: '3',
    title: '为什么我们总是无法做出满意的决策？',
    content: `发现一个有趣的现象：即使有了决策矩阵这样的工具，很多人（包括我自己）还是经常后悔自己的选择。

想讨论几个问题：

1. **信息过载** - 选项越多越难选？
2. **完美主义** - 总想找到"最优解"？
3. **损失厌恶** - 害怕做出错误选择？
4. **后见之明** - 用结果反推决策质量？

有没有什么方法可以帮助我们更平和地接受自己的决策？`,
    summary: '探讨为什么有了决策工具还是难以做出满意决策',
    cover: 'https://images.unsplash.com/photo-1516321318423-f06f85e504b3?w=800',
    user: mockUsers[2]!,
    category: 'question',
    tags: ['决策', '心理学', '选择困难'],
    viewCount: 2341,
    likeCount: 167,
    commentCount: 89,
    favoriteCount: 112,
    isTop: false,
    isFeatured: false,
    status: 1,
    publishTime: '2024-02-08T18:45:00Z',
    createdAt: '2024-02-08T18:45:00Z',
    updatedAt: '2024-02-10T20:30:00Z',
  },
  {
    id: '4',
    title: '六顶思考帽在团队会议中的实践心得',
    content: `在团队中实践六顶思考帽一个月了，分享一些心得。

**效果明显的场景：**
- 头脑风暴会议
- 项目复盘
- 棘手问题讨论

**踩过的坑：**
1. 不要在一个问题上停留太久
2. 主持人要控制好节奏
3. 黑帽（批判）容易过度使用

**意外收获：**
团队氛围变好了！因为每个人都有机会表达不同角度的想法，减少了"意见不合"的对立感。

有实践经验的朋友欢迎交流！`,
    summary: '六顶思考帽在团队会议中的一个月实践心得',
    user: mockUsers[3]!,
    modelId: '6',
    modelName: '六顶思考帽',
    category: 'share',
    tags: ['六顶思考帽', '团队管理', '会议技巧', '经验分享'],
    viewCount: 678,
    likeCount: 45,
    commentCount: 12,
    favoriteCount: 34,
    isTop: false,
    isFeatured: false,
    status: 1,
    publishTime: '2024-02-07T11:00:00Z',
    createdAt: '2024-02-07T11:00:00Z',
    updatedAt: '2024-02-09T14:45:00Z',
  },
  {
    id: '5',
    title: '突然想到：AI 时代，什么样的思维能力最重要？',
    content: `刚才在用 AI 写代码，突然想到一个问题：

AI 能做越来越多的事情，那人类还需要学习什么思维技能？

我的猜想：
- **提问能力** - 问对问题比找答案更重要
- **系统思维** - 理解事物间的关联
- **批判性思维** - 判断信息的可靠性
- **创造性思维** - AI 难以复制的灵感

你们觉得呢？欢迎碰撞！💭`,
    summary: 'AI 时代人类最需要什么样的思维能力？',
    cover: 'https://images.unsplash.com/photo-1677442136019-21780ecad995?w=800',
    user: mockUsers[4]!,
    category: 'inspiration',
    tags: ['AI', '思维能力', '未来趋势', '灵感'],
    viewCount: 3456,
    likeCount: 234,
    commentCount: 156,
    favoriteCount: 189,
    isTop: true,
    isFeatured: true,
    status: 1,
    publishTime: '2024-02-06T20:15:00Z',
    createdAt: '2024-02-06T20:15:00Z',
    updatedAt: '2024-02-12T08:00:00Z',
  },
  {
    id: '6',
    title: '读书笔记：《思考，快与慢》中的认知偏差',
    content: `最近重读了卡尼曼的《思考，快与慢》，整理了一些常见的认知偏差，分享给大家。

**系统1（快思考）容易产生的偏差：**

1. **锚定效应** - 第一印象影响后续判断
2. **可得性偏差** - 容易想到的例子被高估
3. **代表性偏差** - 以貌取人
4. **损失厌恶** - 损失的痛苦大于收益的快乐

**如何应对？**
- 放慢决策速度
- 主动寻找反例
- 用数据代替直觉
- 建立检查清单

有没有朋友愿意分享其他认知偏差的例子？`,
    summary: '《思考，快与慢》读书笔记：常见认知偏差及应对方法',
    user: mockUsers[0]!,
    category: 'share',
    tags: ['读书笔记', '认知偏差', '思考快与慢', '心理学'],
    viewCount: 1567,
    likeCount: 123,
    commentCount: 45,
    favoriteCount: 178,
    isTop: false,
    isFeatured: true,
    status: 1,
    publishTime: '2024-02-05T16:30:00Z',
    createdAt: '2024-02-05T16:30:00Z',
    updatedAt: '2024-02-08T12:00:00Z',
  },
];

const mockComments: CollisionApi.Comment[] = [
  {
    id: '1',
    discussionId: '1',
    content: '非常认同！我也经历过类似的思考过程。后来选择了 gap year，虽然当时家人不理解，但确实让我更清楚自己想要什么。',
    user: mockUsers[1]!,
    likeCount: 12,
    replyCount: 2,
    createdAt: '2024-02-10T12:00:00Z',
    replies: [
      {
        id: '1-1',
        discussionId: '1',
        parentId: '1',
        content: 'Gap year 那段时间都在做什么？有什么收获吗？',
        user: mockUsers[2]!,
        likeCount: 3,
        replyCount: 0,
        createdAt: '2024-02-10T14:30:00Z',
      },
      {
        id: '1-2',
        discussionId: '1',
        parentId: '1',
        replyToId: '3',
        replyToName: '逻辑大师',
        content: '@逻辑大师 主要是旅行和学习新技能，最大的收获是明确了自己真正的兴趣方向',
        user: mockUsers[1]!,
        likeCount: 5,
        replyCount: 0,
        createdAt: '2024-02-10T15:00:00Z',
      },
    ],
  },
  {
    id: '2',
    discussionId: '1',
    content: '第一性原理确实是个好框架，但我觉得也不能完全忽视经验和传统路径的价值。毕竟有些坑前人已经帮我们踩过了。',
    user: mockUsers[3]!,
    likeCount: 8,
    replyCount: 1,
    createdAt: '2024-02-10T16:00:00Z',
    replies: [
      {
        id: '2-1',
        discussionId: '1',
        parentId: '2',
        content: '说得对，关键是要在创新和借鉴之间找到平衡',
        user: mockUsers[0]!,
        likeCount: 6,
        replyCount: 0,
        createdAt: '2024-02-10T17:00:00Z',
      },
    ],
  },
  {
    id: '3',
    discussionId: '1',
    content: '我觉得职业成功的定义因人而异。有人追求财务自由，有人追求社会影响力，有人追求工作生活平衡。没有标准答案。',
    user: mockUsers[4]!,
    likeCount: 15,
    replyCount: 0,
    createdAt: '2024-02-11T09:00:00Z',
  },
];

// ==================== API 函数 ====================

/**
 * 获取话题列表
 */
export async function getDiscussionListApi(params: CollisionApi.DiscussionListParams = {}): Promise<CollisionApi.DiscussionListResult> {
  // 模拟 API 延迟
  await new Promise(resolve => setTimeout(resolve, 300));

  let filtered = [...mockDiscussions];

  // 关键词搜索
  if (params.keyword) {
    const kw = params.keyword.toLowerCase();
    filtered = filtered.filter(d =>
      d.title.toLowerCase().includes(kw) ||
      d.content.toLowerCase().includes(kw) ||
      d.tags.some(t => t.toLowerCase().includes(kw))
    );
  }

  // 分类筛选
  if (params.category && params.category !== 'all') {
    filtered = filtered.filter(d => d.category === params.category);
  }

  // 模型筛选
  if (params.modelId) {
    filtered = filtered.filter(d => d.modelId === params.modelId);
  }

  // 用户筛选
  if (params.userId) {
    filtered = filtered.filter(d => d.user.id === params.userId);
  }

  // 排序
  switch (params.sortBy) {
    case 'popular':
      filtered.sort((a, b) => b.viewCount - a.viewCount);
      break;
    case 'mostCommented':
      filtered.sort((a, b) => b.commentCount - a.commentCount);
      break;
    case 'mostLiked':
      filtered.sort((a, b) => b.likeCount - a.likeCount);
      break;
    default: // latest
      filtered.sort((a, b) => {
        // 置顶优先
        if (a.isTop !== b.isTop) return a.isTop ? -1 : 1;
        return new Date(b.publishTime).getTime() - new Date(a.publishTime).getTime();
      });
  }

  // 分页
  const page = params.page || 1;
  const pageSize = params.pageSize || 10;
  const start = (page - 1) * pageSize;
  const list = filtered.slice(start, start + pageSize);

  return {
    list,
    total: filtered.length,
    page,
    pageSize,
  };
}

/**
 * 获取话题详情
 */
export async function getDiscussionDetailApi(id: string): Promise<CollisionApi.Discussion | null> {
  await new Promise(resolve => setTimeout(resolve, 200));
  const discussion = mockDiscussions.find(d => d.id === id);
  if (discussion) {
    // 模拟增加浏览量
    discussion.viewCount += 1;
  }
  return discussion || null;
}

/**
 * 创建话题
 */
export async function createDiscussionApi(data: CollisionApi.CreateDiscussionParams): Promise<CollisionApi.Discussion> {
  await new Promise(resolve => setTimeout(resolve, 500));

  const newDiscussion: CollisionApi.Discussion = {
    id: String(Date.now()),
    title: data.title,
    content: data.content,
    summary: data.summary || data.content.slice(0, 100),
    cover: data.cover,
    user: mockUsers[0]!,
    modelId: data.modelId,
    modelName: data.modelId ? `模型 ${data.modelId}` : undefined,
    category: data.category,
    tags: data.tags || [],
    viewCount: 0,
    likeCount: 0,
    commentCount: 0,
    favoriteCount: 0,
    isTop: false,
    isFeatured: false,
    status: data.status || 1,
    publishTime: new Date().toISOString(),
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  };

  mockDiscussions.unshift(newDiscussion);
  return newDiscussion;
}

/**
 * 更新话题
 */
export async function updateDiscussionApi(data: CollisionApi.UpdateDiscussionParams): Promise<CollisionApi.Discussion | null> {
  await new Promise(resolve => setTimeout(resolve, 300));

  const index = mockDiscussions.findIndex(d => d.id === data.id);
  if (index === -1) return null;

  const discussion = mockDiscussions[index]!;
  Object.assign(discussion, {
    ...data,
    updatedAt: new Date().toISOString(),
  });

  return discussion;
}

/**
 * 删除话题
 */
export async function deleteDiscussionApi(id: string): Promise<boolean> {
  await new Promise(resolve => setTimeout(resolve, 300));

  const index = mockDiscussions.findIndex(d => d.id === id);
  if (index === -1) return false;

  mockDiscussions.splice(index, 1);
  return true;
}

/**
 * 点赞话题
 */
export async function likeDiscussionApi(id: string): Promise<{ liked: boolean; likeCount: number }> {
  await new Promise(resolve => setTimeout(resolve, 200));

  const discussion = mockDiscussions.find(d => d.id === id);
  if (!discussion) throw new Error('话题不存在');

  discussion.isLiked = !discussion.isLiked;
  discussion.likeCount += discussion.isLiked ? 1 : -1;

  return {
    liked: discussion.isLiked,
    likeCount: discussion.likeCount,
  };
}

/**
 * 收藏话题
 */
export async function favoriteDiscussionApi(id: string): Promise<{ favorited: boolean; favoriteCount: number }> {
  await new Promise(resolve => setTimeout(resolve, 200));

  const discussion = mockDiscussions.find(d => d.id === id);
  if (!discussion) throw new Error('话题不存在');

  discussion.isFavorited = !discussion.isFavorited;
  discussion.favoriteCount += discussion.isFavorited ? 1 : -1;

  return {
    favorited: discussion.isFavorited,
    favoriteCount: discussion.favoriteCount,
  };
}

/**
 * 获取评论列表
 */
export async function getCommentListApi(params: CollisionApi.CommentListParams): Promise<CollisionApi.CommentListResult> {
  await new Promise(resolve => setTimeout(resolve, 300));

  let filtered = mockComments.filter(c => c.discussionId === params.discussionId);

  if (params.parentId) {
    filtered = filtered.filter(c => c.parentId === params.parentId);
  } else {
    filtered = filtered.filter(c => !c.parentId);
  }

  const page = params.page || 1;
  const pageSize = params.pageSize || 10;
  const start = (page - 1) * pageSize;
  const list = filtered.slice(start, start + pageSize);

  return {
    list,
    total: filtered.length,
    page,
    pageSize,
  };
}

/**
 * 创建评论
 */
export async function createCommentApi(data: CollisionApi.CreateCommentParams): Promise<CollisionApi.Comment> {
  await new Promise(resolve => setTimeout(resolve, 300));

  const newComment: CollisionApi.Comment = {
    id: String(Date.now()),
    discussionId: data.discussionId,
    parentId: data.parentId,
    replyToId: data.replyToId,
    content: data.content,
    user: mockUsers[0]!,
    likeCount: 0,
    replyCount: 0,
    createdAt: new Date().toISOString(),
  };

  // 更新话题评论数
  const discussion = mockDiscussions.find(d => d.id === data.discussionId);
  if (discussion) {
    discussion.commentCount += 1;
  }

  // 更新父评论回复数
  if (data.parentId) {
    const parent = mockComments.find(c => c.id === data.parentId);
    if (parent) {
      parent.replyCount += 1;
      if (!parent.replies) parent.replies = [];
      parent.replies.push(newComment);
    }
  } else {
    mockComments.push(newComment);
  }

  return newComment;
}

/**
 * 点赞评论
 */
export async function likeCommentApi(id: string): Promise<{ liked: boolean; likeCount: number }> {
  await new Promise(resolve => setTimeout(resolve, 200));

  const findComment = (comments: CollisionApi.Comment[]): CollisionApi.Comment | null => {
    for (const c of comments) {
      if (c.id === id) return c;
      if (c.replies) {
        const found = findComment(c.replies);
        if (found) return found;
      }
    }
    return null;
  };

  const comment = findComment(mockComments);
  if (!comment) throw new Error('评论不存在');

  comment.isLiked = !comment.isLiked;
  comment.likeCount += comment.isLiked ? 1 : -1;

  return {
    liked: comment.isLiked,
    likeCount: comment.likeCount,
  };
}

/**
 * 获取热门标签
 */
export async function getHotTagsApi(): Promise<string[]> {
  await new Promise(resolve => setTimeout(resolve, 100));
  const tagCounts = new Map<string, number>();
  
  mockDiscussions.forEach(d => {
    d.tags.forEach(tag => {
      tagCounts.set(tag, (tagCounts.get(tag) || 0) + 1);
    });
  });

  return Array.from(tagCounts.entries())
    .sort((a, b) => b[1] - a[1])
    .slice(0, 10)
    .map(([tag]) => tag);
}

/**
 * 获取活跃用户
 */
export async function getActiveUsersApi(): Promise<CollisionApi.UserInfo[]> {
  await new Promise(resolve => setTimeout(resolve, 100));
  return mockUsers.slice(0, 5);
}

// ==================== 找人聊聊 模拟数据 ====================
const mockMeetups: CollisionApi.Meetup[] = [
  {
    id: 'm1',
    title: '聊聊用第一性原理做职业规划',
    topic: 'career',
    description: '想找几位朋友线下交流一下，如何用第一性原理来重新思考职业发展路径。',
    thoughts: `最近一直在思考一个问题：传统的职业发展路径真的适合所有人吗？

我尝试用第一性原理来分析：
1. 工作的本质是价值交换
2. 技能是可以迁移和组合的
3. 收入来源不一定是单一的

基于这些思考，我觉得"跳出公司-职级"的框架来看职业发展，会有不一样的风景。

希望找到志同道合的朋友，一起探讨这个话题。`,
    host: { ...mockUsers[0]!, bio: '10年产品经理，正在探索职业转型', interests: ['思维模型', '职业发展', '自由职业'] },
    modelId: '3',
    modelName: '第一性原理思维',
    tags: ['职业规划', '第一性原理', '人生选择'],
    city: '北京',
    preferredTime: '周末下午 2-5 点',
    costSplit: 'host',
    maxGuests: 3,
    currentGuests: 1,
    status: 'open',
    viewCount: 234,
    interestedCount: 12,
    createdAt: '2024-02-08T10:00:00Z',
    updatedAt: '2024-02-10T15:30:00Z',
  },
  {
    id: 'm2',
    title: '创业想法头脑风暴',
    topic: 'startup',
    description: '有一个AI相关的创业想法，想约几位技术和产品背景的朋友聊聊可行性。',
    thoughts: `想法背景：
我观察到一个有趣的现象——很多人在使用 AI 工具时，不知道如何写好 prompt。

痛点：
- 普通用户不懂提示词工程
- 好的 prompt 很难被发现和复用
- 不同场景需要不同的 prompt 策略

我的初步想法是做一个"AI 助手的助手"，帮用户优化和管理 prompt。

希望听听大家的看法，验证一下这个方向是否有价值。`,
    host: { ...mockUsers[1]!, bio: '全栈工程师，连续创业者', interests: ['AI', '创业', '产品设计'] },
    modelId: '1',
    modelName: 'SWOT 分析',
    tags: ['创业', 'AI', '产品验证', 'MVP'],
    city: '上海',
    preferredTime: '工作日晚上 7-9 点',
    costSplit: 'aa',
    maxGuests: 4,
    currentGuests: 2,
    status: 'open',
    viewCount: 456,
    interestedCount: 28,
    createdAt: '2024-02-05T14:00:00Z',
    updatedAt: '2024-02-11T09:00:00Z',
  },
  {
    id: 'm3',
    title: '深度聊聊投资认知',
    topic: 'investment',
    description: '想和有投资经验的朋友交流一下投资思维框架和认知升级。',
    thoughts: `我的投资经历比较曲折，走过很多弯路。最近在反思：

以前的问题：
- 追涨杀跌，情绪化决策
- 没有自己的投资框架
- 信息茧房，只看自己想看的

现在的思考：
- 投资本质是认知变现
- 需要建立自己的决策系统
- 多元信息源 + 独立思考

想找几位朋友深度交流，互相学习投资思维。

要求：有实际投资经验，愿意分享真实案例（不要抱着推销产品的目的来）`,
    host: { ...mockUsers[2]!, bio: '金融从业者，价值投资践行者', interests: ['价值投资', '财务分析', '行业研究'] },
    tags: ['投资', '认知升级', '思维框架'],
    city: '深圳',
    preferredTime: '周六全天',
    costSplit: 'host',
    maxGuests: 2,
    currentGuests: 0,
    status: 'open',
    viewCount: 189,
    interestedCount: 8,
    createdAt: '2024-02-09T16:00:00Z',
    updatedAt: '2024-02-09T16:00:00Z',
  },
  {
    id: 'm4',
    title: '技术人如何提升产品思维',
    topic: 'technology',
    description: '作为技术出身，想和有产品经验的朋友聊聊如何培养产品思维。',
    thoughts: `困惑：
作为一个写了8年代码的程序员，我发现自己陷入了一个怪圈：
- 技术能力越来越强
- 但对产品的理解越来越机械
- 不知道怎么从"做得好"进化到"做对的事"

想要获得：
- 产品思维是怎么训练出来的？
- 技术人常见的产品认知误区？
- 有没有推荐的学习路径？

希望和有产品背景的朋友交流，打开视野。`,
    host: { ...mockUsers[3]!, bio: '8年后端开发，想向产品方向拓展', interests: ['技术架构', '产品设计', '用户体验'] },
    modelId: '5',
    modelName: '用户思维',
    tags: ['产品思维', '技术转型', '能力拓展'],
    city: '杭州',
    preferredTime: '周末均可',
    costSplit: 'aa',
    maxGuests: 3,
    currentGuests: 1,
    status: 'open',
    viewCount: 312,
    interestedCount: 15,
    createdAt: '2024-02-07T11:00:00Z',
    updatedAt: '2024-02-10T20:00:00Z',
  },
  {
    id: 'm5',
    title: '35岁职业焦虑怎么破？',
    topic: 'life',
    description: '想找同龄人聊聊35岁职业焦虑这个话题，分享彼此的应对策略。',
    thoughts: `最近"35岁危机"这个话题很火，我自己也快到这个年龄了。

我的焦虑来源：
- 行业内卷严重
- 新技术更新太快
- 体力和精力不如从前
- 家庭责任越来越重

我尝试的应对方法：
- 建立被动收入
- 发展第二技能
- 保持学习习惯
- 调整心态预期

想听听其他朋友是怎么看待和应对这个问题的。大家互相打打气，也许能找到新的思路。`,
    host: { ...mockUsers[4]!, bio: '34岁互联网人，正在寻找职业第二曲线', interests: ['职业发展', '个人成长', '生活平衡'] },
    tags: ['35岁危机', '职业焦虑', '人生规划', '心态调整'],
    city: '广州',
    preferredTime: '周日下午',
    costSplit: 'host',
    maxGuests: 4,
    currentGuests: 3,
    status: 'pending',
    viewCount: 567,
    interestedCount: 42,
    createdAt: '2024-02-03T09:00:00Z',
    updatedAt: '2024-02-11T14:00:00Z',
  },
];

const mockApplications: CollisionApi.MeetupApplication[] = [
  {
    id: 'a1',
    meetupId: 'm1',
    applicant: mockUsers[2]!,
    message: '我也在思考类似的问题，目前在尝试做自由职业，很想交流一下经验。',
    status: 'accepted',
    createdAt: '2024-02-09T10:00:00Z',
  },
  {
    id: 'a2',
    meetupId: 'm2',
    applicant: mockUsers[3]!,
    message: '做过几年产品，对 AI 领域也很感兴趣，希望能参与讨论。',
    status: 'accepted',
    createdAt: '2024-02-06T18:00:00Z',
  },
  {
    id: 'a3',
    meetupId: 'm2',
    applicant: mockUsers[4]!,
    message: '正在做 AI 工具类产品，有一些实战经验可以分享。',
    status: 'accepted',
    createdAt: '2024-02-07T09:00:00Z',
  },
];

// ==================== 找人聊聊 API ====================

/**
 * 获取约见列表
 */
export async function getMeetupListApi(params: CollisionApi.MeetupListParams = {}): Promise<CollisionApi.MeetupListResult> {
  await new Promise(resolve => setTimeout(resolve, 400));

  let filtered = [...mockMeetups];

  // 主题筛选
  if (params.topic && params.topic !== 'all') {
    filtered = filtered.filter(m => m.topic === params.topic);
  }

  // 城市筛选
  if (params.city) {
    filtered = filtered.filter(m => m.city === params.city);
  }

  // 费用方式筛选
  if (params.costSplit && params.costSplit !== 'all') {
    filtered = filtered.filter(m => m.costSplit === params.costSplit);
  }

  // 状态筛选
  if (params.status && params.status !== 'all') {
    filtered = filtered.filter(m => m.status === params.status);
  }

  // 关键词搜索
  if (params.keyword) {
    const kw = params.keyword.toLowerCase();
    filtered = filtered.filter(m =>
      m.title.toLowerCase().includes(kw) ||
      m.description.toLowerCase().includes(kw) ||
      m.host.name.toLowerCase().includes(kw) ||
      m.tags.some(t => t.toLowerCase().includes(kw))
    );
  }

  // 排序
  if (params.sortBy === 'popular') {
    filtered.sort((a, b) => b.interestedCount - a.interestedCount);
  } else if (params.sortBy === 'soonest') {
    // 按时间排序（这里简化处理）
    filtered.sort((a, b) => new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime());
  } else {
    // 默认最新
    filtered.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime());
  }

  const page = params.page || 1;
  const pageSize = params.pageSize || 10;
  const start = (page - 1) * pageSize;
  const list = filtered.slice(start, start + pageSize);

  return {
    list,
    total: filtered.length,
    page,
    pageSize,
  };
}

/**
 * 获取约见详情
 */
export async function getMeetupDetailApi(id: string): Promise<CollisionApi.Meetup | null> {
  await new Promise(resolve => setTimeout(resolve, 300));
  const meetup = mockMeetups.find(m => m.id === id);
  if (meetup) {
    meetup.viewCount += 1;
  }
  return meetup || null;
}

/**
 * 创建约见
 */
export async function createMeetupApi(data: CollisionApi.CreateMeetupParams): Promise<CollisionApi.Meetup> {
  await new Promise(resolve => setTimeout(resolve, 500));

  const newMeetup: CollisionApi.Meetup = {
    id: `m${Date.now()}`,
    title: data.title,
    topic: data.topic,
    description: data.description,
    thoughts: data.thoughts,
    host: mockUsers[0]!,
    modelId: data.modelId,
    tags: data.tags || [],
    city: data.city,
    preferredTime: data.preferredTime,
    costSplit: data.costSplit,
    maxGuests: data.maxGuests,
    currentGuests: 0,
    status: 'open',
    viewCount: 0,
    interestedCount: 0,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  };

  mockMeetups.unshift(newMeetup);
  return newMeetup;
}

/**
 * 对约见表示感兴趣
 */
export async function interestMeetupApi(id: string): Promise<{ interested: boolean; interestedCount: number }> {
  await new Promise(resolve => setTimeout(resolve, 200));

  const meetup = mockMeetups.find(m => m.id === id);
  if (!meetup) throw new Error('约见不存在');

  meetup.isInterested = !meetup.isInterested;
  meetup.interestedCount += meetup.isInterested ? 1 : -1;

  return {
    interested: meetup.isInterested,
    interestedCount: meetup.interestedCount,
  };
}

/**
 * 申请约见
 */
export async function applyMeetupApi(data: CollisionApi.ApplyMeetupParams): Promise<CollisionApi.MeetupApplication> {
  await new Promise(resolve => setTimeout(resolve, 400));

  const application: CollisionApi.MeetupApplication = {
    id: `a${Date.now()}`,
    meetupId: data.meetupId,
    applicant: mockUsers[0]!,
    message: data.message,
    status: 'pending',
    createdAt: new Date().toISOString(),
  };

  mockApplications.push(application);
  return application;
}

/**
 * 获取约见申请列表
 */
export async function getMeetupApplicationsApi(meetupId: string): Promise<CollisionApi.MeetupApplication[]> {
  await new Promise(resolve => setTimeout(resolve, 200));
  return mockApplications.filter(a => a.meetupId === meetupId);
}

/**
 * 获取可用城市列表
 */
export async function getMeetupCitiesApi(): Promise<string[]> {
  await new Promise(resolve => setTimeout(resolve, 100));
  const cities = new Set(mockMeetups.map(m => m.city));
  return Array.from(cities);
}

// ==================== 付费咨询 模拟数据 ====================
const mockExperts: CollisionApi.Expert[] = [
  {
    id: 'e1',
    name: '张明远',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=expert1',
    bio: '15年互联网产品经验，前大厂产品总监',
    title: '资深产品专家',
    company: '前阿里巴巴',
    experience: 15,
    expertise: ['产品设计', '用户增长', '商业化'],
    consultCount: 128,
    rating: 4.9,
    responseRate: 0.95,
    hourlyRate: 500,
  },
  {
    id: 'e2',
    name: '李思涵',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=expert2',
    bio: '连续创业者，3次成功退出经历',
    title: '创业导师',
    company: '某知名VC合伙人',
    experience: 12,
    expertise: ['创业融资', '商业模式', '团队管理'],
    consultCount: 86,
    rating: 4.8,
    responseRate: 0.88,
    hourlyRate: 800,
  },
  {
    id: 'e3',
    name: '王建国',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=expert3',
    bio: '技术架构师，擅长大规模系统设计',
    title: '首席架构师',
    company: '某头部互联网公司',
    experience: 18,
    expertise: ['系统架构', '技术选型', '团队建设'],
    consultCount: 95,
    rating: 4.95,
    responseRate: 0.92,
    hourlyRate: 600,
  },
  {
    id: 'e4',
    name: '陈雨晴',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=expert4',
    bio: '职业规划师，帮助500+人成功转型',
    title: '职业发展顾问',
    company: '独立咨询师',
    experience: 10,
    expertise: ['职业规划', '面试辅导', '简历优化'],
    consultCount: 256,
    rating: 4.85,
    responseRate: 0.98,
    hourlyRate: 300,
  },
];

const mockConsultations: CollisionApi.Consultation[] = [
  {
    id: 'c1',
    title: '创业初期如何找到产品市场契合点（PMF）？',
    field: 'startup',
    description: '我们是一个3人技术团队，有一个B端SaaS产品的想法，已经做了MVP，但不确定如何验证产品市场契合度。',
    background: `我们团队背景：
- 3个技术合伙人，都有5年以上大厂经验
- 目标市场：中小企业数据分析
- 已完成：MVP产品开发，有10个种子用户
- 问题：用户反馈不一，不知道该往哪个方向迭代

目前的困惑：
1. 如何判断哪些用户反馈是真需求？
2. 什么指标可以说明找到了PMF？
3. 在没有找到PMF之前，要不要大规模推广？`,
    expectation: `希望得到：
1. 判断PMF的具体方法和指标
2. 早期产品迭代的优先级建议
3. 种子用户运营的经验分享
4. 是否该继续投入或及时止损的判断标准`,
    requester: mockUsers[1]!,
    modelId: '1',
    modelName: 'SWOT 分析',
    tags: ['创业', 'PMF', 'SaaS', 'B端'],
    reward: 500,
    mode: 'online',
    deadline: '2024-02-28T23:59:59Z',
    status: 'open',
    viewCount: 456,
    applicationCount: 3,
    createdAt: '2024-02-10T10:00:00Z',
    updatedAt: '2024-02-12T15:30:00Z',
  },
  {
    id: 'c2',
    title: '技术转管理，如何快速提升领导力？',
    field: 'management',
    description: '最近被提拔为技术经理，管理一个8人团队，感觉自己还在用IC思维做事，希望得到管理方面的指导。',
    background: `个人情况：
- 工作8年，一直是技术骨干
- 上个月被提拔为技术经理
- 团队8人，2个资深、4个中级、2个初级
- 目前状态：很多事情习惯自己上手，下属不太主动

具体困扰：
1. 总觉得下属做的不如自己好，忍不住想自己做
2. 开会的时候不知道该说什么
3. 和下属1:1不知道聊什么
4. 上级布置的任务不知道该怎么拆解给团队`,
    expectation: `希望获得：
1. 从IC到Manager心态转变的方法
2. 团队管理的基本框架和工具
3. 如何培养和激励下属
4. 向上管理的技巧`,
    requester: mockUsers[3]!,
    tags: ['管理', '技术转型', '领导力', '团队管理'],
    reward: 800,
    mode: 'both',
    city: '北京',
    deadline: '2024-02-25T23:59:59Z',
    status: 'open',
    viewCount: 678,
    applicationCount: 5,
    createdAt: '2024-02-08T14:00:00Z',
    updatedAt: '2024-02-11T09:00:00Z',
  },
  {
    id: 'c3',
    title: '35岁程序员，转型产品还是技术深耕？',
    field: 'career',
    description: '面临职业选择困惑，不知道是继续深耕技术还是转型产品经理，希望有经验的前辈指点。',
    background: `我的情况：
- 35岁，10年Java开发经验
- 目前在二线互联网公司做后端开发
- 技术能力中上，但对新技术学习热情减退
- 对产品有兴趣，经常会思考为什么要这样做

焦虑来源：
1. 担心年龄大了技术竞争力下降
2. 看到很多技术转产品成功的案例
3. 但也见过转型失败的例子
4. 家庭压力大，不敢轻易冒险`,
    expectation: `想了解：
1. 35岁转产品是否可行？有哪些坑？
2. 如果继续做技术，该往什么方向深耕？
3. 如何评估自己更适合哪条路？
4. 有没有折中的方案？`,
    requester: mockUsers[4]!,
    modelId: '3',
    modelName: '第一性原理思维',
    tags: ['职业规划', '35岁危机', '转型', '程序员'],
    reward: 300,
    mode: 'online',
    deadline: '2024-03-01T23:59:59Z',
    status: 'open',
    viewCount: 1234,
    applicationCount: 8,
    createdAt: '2024-02-05T09:00:00Z',
    updatedAt: '2024-02-10T20:00:00Z',
  },
  {
    id: 'c4',
    title: '如何设计一个高并发的订单系统？',
    field: 'technology',
    description: '正在主导公司新订单系统的架构设计，日订单量预计百万级，希望得到架构方面的专业指导。',
    background: `项目背景：
- 电商公司，日活UV 500万
- 现有订单系统是单体架构，性能瓶颈明显
- 需要重构为分布式架构
- 团队10人，有3人有分布式经验

技术挑战：
1. 高并发下的数据一致性
2. 分布式事务处理
3. 库存扣减的准确性
4. 订单状态机设计
5. 技术选型（消息队列、缓存、数据库）`,
    expectation: `希望获得：
1. 百万级订单系统的整体架构方案
2. 关键技术点的选型建议
3. 踩坑经验和注意事项
4. 团队分工和项目推进建议`,
    requester: mockUsers[2]!,
    tags: ['架构设计', '高并发', '分布式', '电商'],
    reward: 1000,
    mode: 'online',
    deadline: '2024-02-20T23:59:59Z',
    status: 'matched',
    selectedExpert: mockExperts[2],
    viewCount: 892,
    applicationCount: 6,
    createdAt: '2024-02-03T16:00:00Z',
    updatedAt: '2024-02-09T14:00:00Z',
  },
  {
    id: 'c5',
    title: '天使轮融资BP和路演准备指导',
    field: 'investment',
    description: '准备进行天使轮融资，需要专业指导如何准备BP和路演。',
    background: `项目情况：
- AI+教育赛道
- 产品已上线，月活2万
- 团队5人，技术为主
- 已有种子用户付费

融资需求：
- 目标金额：500万
- 用途：团队扩张 + 市场推广
- 已接触3家投资机构，但都没有下文`,
    expectation: `希望得到：
1. BP的框架和重点
2. 路演技巧和常见问题
3. 估值逻辑
4. 投资人关注点`,
    requester: mockUsers[0]!,
    tags: ['融资', 'BP', '路演', '天使轮', 'AI教育'],
    reward: 1500,
    mode: 'offline',
    city: '上海',
    deadline: '2024-02-22T23:59:59Z',
    status: 'open',
    viewCount: 567,
    applicationCount: 4,
    createdAt: '2024-02-07T11:00:00Z',
    updatedAt: '2024-02-11T16:00:00Z',
  },
];

const mockConsultationApplications: CollisionApi.ConsultationApplication[] = [
  {
    id: 'ca1',
    consultationId: 'c1',
    expert: mockExperts[1]!,
    proposal: '我有3次创业经历，其中2次成功找到PMF并完成融资。可以从实战角度分享PMF验证的方法论，包括用户访谈技巧、数据指标设计、迭代优先级判断等。',
    estimatedTime: '2小时线上沟通',
    status: 'pending',
    createdAt: '2024-02-11T10:00:00Z',
  },
  {
    id: 'ca2',
    consultationId: 'c2',
    expert: mockExperts[0]!,
    proposal: '我从技术骨干到产品总监，管理过100+人的团队。可以分享从IC到Manager的转变心得，包括：放权的艺术、1:1沟通技巧、团队激励方法、向上管理策略等。',
    estimatedTime: '1.5小时线上 + 可选线下',
    status: 'pending',
    createdAt: '2024-02-09T15:00:00Z',
  },
];

// ==================== 付费咨询 API ====================

/**
 * 获取咨询列表
 */
export async function getConsultationListApi(params: CollisionApi.ConsultationListParams = {}): Promise<CollisionApi.ConsultationListResult> {
  await new Promise(resolve => setTimeout(resolve, 400));

  let filtered = [...mockConsultations];

  // 领域筛选
  if (params.field && params.field !== 'all') {
    filtered = filtered.filter(c => c.field === params.field);
  }

  // 方式筛选
  if (params.mode && params.mode !== 'all') {
    filtered = filtered.filter(c => c.mode === params.mode || c.mode === 'both');
  }

  // 状态筛选
  if (params.status && params.status !== 'all') {
    filtered = filtered.filter(c => c.status === params.status);
  }

  // 金额筛选
  if (params.minReward !== undefined) {
    filtered = filtered.filter(c => c.reward >= params.minReward!);
  }
  if (params.maxReward !== undefined) {
    filtered = filtered.filter(c => c.reward <= params.maxReward!);
  }

  // 关键词搜索
  if (params.keyword) {
    const kw = params.keyword.toLowerCase();
    filtered = filtered.filter(c =>
      c.title.toLowerCase().includes(kw) ||
      c.description.toLowerCase().includes(kw) ||
      c.requester.name.toLowerCase().includes(kw) ||
      c.tags.some(t => t.toLowerCase().includes(kw))
    );
  }

  // 排序
  if (params.sortBy === 'reward') {
    filtered.sort((a, b) => b.reward - a.reward);
  } else if (params.sortBy === 'deadline') {
    filtered.sort((a, b) => new Date(a.deadline).getTime() - new Date(b.deadline).getTime());
  } else if (params.sortBy === 'popular') {
    filtered.sort((a, b) => b.viewCount - a.viewCount);
  } else {
    // 默认最新
    filtered.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime());
  }

  const page = params.page || 1;
  const pageSize = params.pageSize || 10;
  const start = (page - 1) * pageSize;
  const list = filtered.slice(start, start + pageSize);

  return {
    list,
    total: filtered.length,
    page,
    pageSize,
  };
}

/**
 * 获取咨询详情
 */
export async function getConsultationDetailApi(id: string): Promise<CollisionApi.Consultation | null> {
  await new Promise(resolve => setTimeout(resolve, 300));
  const consultation = mockConsultations.find(c => c.id === id);
  if (consultation) {
    consultation.viewCount += 1;
  }
  return consultation || null;
}

/**
 * 创建咨询
 */
export async function createConsultationApi(data: CollisionApi.CreateConsultationParams): Promise<CollisionApi.Consultation> {
  await new Promise(resolve => setTimeout(resolve, 500));

  const newConsultation: CollisionApi.Consultation = {
    id: `c${Date.now()}`,
    title: data.title,
    field: data.field,
    description: data.description,
    background: data.background,
    expectation: data.expectation,
    requester: mockUsers[0]!,
    modelId: data.modelId,
    tags: data.tags || [],
    reward: data.reward,
    mode: data.mode,
    city: data.city,
    deadline: data.deadline,
    status: 'open',
    viewCount: 0,
    applicationCount: 0,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  };

  mockConsultations.unshift(newConsultation);
  return newConsultation;
}

/**
 * 申请咨询（专家）
 */
export async function applyConsultationApi(data: CollisionApi.ApplyConsultationParams): Promise<CollisionApi.ConsultationApplication> {
  await new Promise(resolve => setTimeout(resolve, 400));

  const application: CollisionApi.ConsultationApplication = {
    id: `ca${Date.now()}`,
    consultationId: data.consultationId,
    expert: mockExperts[0]!,
    proposal: data.proposal,
    estimatedTime: data.estimatedTime,
    quotation: data.quotation,
    status: 'pending',
    createdAt: new Date().toISOString(),
  };

  mockConsultationApplications.push(application);

  // 更新申请数
  const consultation = mockConsultations.find(c => c.id === data.consultationId);
  if (consultation) {
    consultation.applicationCount += 1;
  }

  return application;
}

/**
 * 获取咨询申请列表
 */
export async function getConsultationApplicationsApi(consultationId: string): Promise<CollisionApi.ConsultationApplication[]> {
  await new Promise(resolve => setTimeout(resolve, 200));
  return mockConsultationApplications.filter(a => a.consultationId === consultationId);
}

/**
 * 获取推荐专家列表
 */
export async function getRecommendedExpertsApi(field?: CollisionApi.ConsultationField): Promise<CollisionApi.Expert[]> {
  await new Promise(resolve => setTimeout(resolve, 200));
  if (field) {
    return mockExperts.filter(e => e.expertise.some(exp => 
      exp.toLowerCase().includes(field.toLowerCase())
    ));
  }
  return mockExperts;
}

/**
 * 获取咨询统计数据
 */
export async function getConsultationStatsApi(): Promise<{
  totalConsultations: number;
  totalExperts: number;
  avgReward: number;
  successRate: number;
}> {
  await new Promise(resolve => setTimeout(resolve, 100));
  return {
    totalConsultations: mockConsultations.length,
    totalExperts: mockExperts.length,
    avgReward: Math.round(mockConsultations.reduce((sum, c) => sum + c.reward, 0) / mockConsultations.length),
    successRate: 0.85,
  };
}
