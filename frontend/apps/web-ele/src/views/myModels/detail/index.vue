<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { Page } from '@vben/common-ui';

import {
  ElButton,
  ElCard,
  ElTag,
  ElMessage,
  ElEmpty,
  ElSkeleton,
  ElSkeletonItem,
  ElMessageBox,
  ElProgress,
  ElAvatar,
} from 'element-plus';

// 路由
const route = useRoute();
const router = useRouter();
const modelId = computed(() => route.params.id as string);

// 加载状态
const loading = ref(true);

// 当前选中 Tab
const activeTab = ref<'overview' | 'steps' | 'examples' | 'feedback' | 'analytics'>('overview');

// Tab 配置
const tabs = [
  { id: 'overview', label: '模型概述', icon: '📋' },
  { id: 'steps', label: '使用步骤', icon: '📝' },
  { id: 'examples', label: '实践案例', icon: '💡' },
  { id: 'feedback', label: '用户反馈', icon: '💬' },
  { id: 'analytics', label: '数据分析', icon: '📊' },
];

// 模型数据类型
interface ModelDetail {
  id: string;
  title: string;
  description: string;
  cover: string;
  category: string;
  categoryName: string;
  tags: string[];
  status: 'published' | 'draft' | 'under_review' | 'rejected';
  statusText: string;
  price: number;
  isFree: boolean;
  stats: {
    adoptions: number;
    practices: number;
    likes: number;
    reviews: number;
    views: number;
  };
  revenue: {
    total: number;
    thisMonth: number;
    lastMonth: number;
    history: { month: string; amount: number }[];
  };
  content: {
    overview: string;
    steps: { title: string; description: string }[];
    examples: { title: string; content: string }[];
  };
  createdAt: string;
  updatedAt: string;
  rejectReason?: string;
}

const model = ref<ModelDetail | null>(null);

// 模拟反馈数据
const mockFeedbacks = [
  { id: '1', user: '张同学', avatar: '', content: '非常实用的思维模型，帮助我整理了很多思路！', rating: 5, date: '2024-02-18' },
  { id: '2', user: '李经理', avatar: '', content: '用来做项目决策分析很有帮助，推荐给团队了。', rating: 5, date: '2024-02-15' },
  { id: '3', user: '王创业者', avatar: '', content: '结构清晰，案例丰富，对商业分析很有帮助。', rating: 4, date: '2024-02-10' },
];

// 模拟数据
const mockModelDetail: ModelDetail = {
  id: '1',
  title: 'SWOT 分析思维模型',
  description: '经典的战略分析工具，帮助分析企业或项目的优势、劣势、机会和威胁，适用于商业决策和个人发展规划。',
  cover: '/images/swot-cover.svg',
  category: 'business',
  categoryName: '商业管理',
  tags: ['战略', '分析', '商业', '管理'],
  status: 'published',
  statusText: '已发布',
  price: 29,
  isFree: false,
  stats: {
    adoptions: 1256,
    practices: 3421,
    likes: 328,
    reviews: 56,
    views: 8900,
  },
  revenue: {
    total: 8560,
    thisMonth: 1200,
    lastMonth: 980,
    history: [
      { month: '2023-09', amount: 450 },
      { month: '2023-10', amount: 520 },
      { month: '2023-11', amount: 680 },
      { month: '2023-12', amount: 890 },
      { month: '2024-01', amount: 1120 },
      { month: '2024-02', amount: 1200 },
    ],
  },
  content: {
    overview: 'SWOT 分析是一种战略规划工具，用于评估企业、项目或个人的优势（Strengths）、劣势（Weaknesses）、机会（Opportunities）和威胁（Threats）。这种方法帮助决策者全面了解内外部环境，制定更有效的战略。\n\n通过系统分析内部优势与劣势，以及外部机会与威胁，SWOT 分析能够帮助决策者：\n1. 更清晰地了解当前状况\n2. 发现潜在的战略方向\n3. 识别需要改进的领域\n4. 为未来的决策提供依据',
    steps: [
      {
        title: '识别优势 (Strengths)',
        description: '列出你或你的组织相对于竞争对手的优势。包括资源、能力、经验、品牌等内部因素。与竞争对手相比，你有什么独特的优势？你有哪些其他人难以复制的资源或能力？',
      },
      {
        title: '识别劣势 (Weaknesses)',
        description: '诚实地列出需要改进的领域。这些是你相对于竞争对手的不足之处。有哪些领域你需要改进？你缺乏哪些资源或能力？',
      },
      {
        title: '发现机会 (Opportunities)',
        description: '分析外部环境中的有利因素。包括市场趋势、政策变化、技术发展等。外部环境中有哪些有利的变化？有哪些未被满足的市场需求？',
      },
      {
        title: '识别威胁 (Threats)',
        description: '评估可能对你产生负面影响的外部因素。包括竞争、经济环境、法规变化等。有哪些外部因素可能对你造成威胁？',
      },
      {
        title: '整合分析与制定策略',
        description: '将 SWOT 四个要素进行交叉分析，制定 SO、WO、ST、WT 四种策略组合，形成全面的行动方案。',
      },
    ],
    examples: [
      {
        title: '案例：某电商平台的 SWOT 分析',
        content: '优势：拥有庞大的用户基础（5亿+注册用户）、完善的物流体系、强大的技术团队、深厚的品牌认知。劣势：运营成本较高、对第三方商家管控力有限、用户活跃度有下降趋势。机会：下沉市场增长潜力大、跨境电商政策利好、直播带货兴起带来新增长点。威胁：竞争对手价格战激烈、监管政策趋严、用户获取成本持续上升。\n\n策略建议：\n1. SO策略：利用用户基础优势，加速拓展下沉市场\n2. WO策略：通过直播带货提升用户活跃度\n3. ST策略：强化品牌优势，避开价格战\n4. WT策略：降本增效，优化商家管理体系',
      },
      {
        title: '案例：个人职业发展的 SWOT 分析',
        content: '优势：专业技能扎实（熟练掌握3门编程语言）、沟通能力强、学习能力快、具有项目管理经验。劣势：管理经验不足、行业人脉资源有限、英语口语能力一般。机会：AI行业快速发展、公司有内部晋升机会、远程工作成为趋势。威胁：35岁职场焦虑、技术更新换代快、AI可能替代部分工作。\n\n策略建议：\n1. SO策略：利用技术优势切入AI领域\n2. WO策略：寻找mentor积累管理经验\n3. ST策略：持续学习保持技术竞争力\n4. WT策略：建立个人品牌，拓展人脉网络',
      },
    ],
  },
  createdAt: '2024-01-15T08:00:00Z',
  updatedAt: '2024-02-10T10:30:00Z',
};

// 获取模型详情
async function fetchModelDetail() {
  loading.value = true;
  try {
    await new Promise(resolve => setTimeout(resolve, 600));
    model.value = mockModelDetail;
  } catch (error) {
    console.error('获取模型详情失败:', error);
    ElMessage.error('获取模型详情失败');
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  fetchModelDetail();
});

// 操作函数
function goBack() {
  router.push('/my-models');
}

function handleEdit() {
  router.push(`/my-models/create?id=${modelId.value}`);
}

async function handleDelete() {
  try {
    await ElMessageBox.confirm(
      `确定要删除模型「${model.value?.title}」吗？此操作不可恢复。`,
      '删除确认',
      { type: 'warning' }
    );
    ElMessage.success('模型已删除');
    router.push('/my-models');
  } catch {
    // 用户取消
  }
}

async function handlePublish() {
  try {
    await ElMessageBox.confirm(
      '提交审核后，模型将在审核通过后发布到市场。确定提交吗？',
      '提交审核',
      { type: 'info' }
    );
    ElMessage.success('模型已提交审核');
    fetchModelDetail();
  } catch {
    // 用户取消
  }
}

async function handleUnpublish() {
  try {
    await ElMessageBox.confirm(
      '下架后用户将无法继续购买此模型。确定下架吗？',
      '确认下架',
      { type: 'warning' }
    );
    ElMessage.success('模型已下架');
    fetchModelDetail();
  } catch {
    // 用户取消
  }
}

// 工具函数
function getStatusStyle(status: string): { bg: string; text: string; label: string } {
  const styles: Record<string, { bg: string; text: string; label: string }> = {
    published: { bg: 'bg-green-100', text: 'text-green-700', label: '已发布' },
    draft: { bg: 'bg-gray-100', text: 'text-gray-600', label: '草稿' },
    under_review: { bg: 'bg-amber-100', text: 'text-amber-700', label: '审核中' },
    rejected: { bg: 'bg-red-100', text: 'text-red-700', label: '已驳回' },
  };
  return styles[status] || { bg: 'bg-gray-100', text: 'text-gray-600', label: status };
}

function formatMoney(amount: number): string {
  if (amount >= 10000) return '¥' + (amount / 10000).toFixed(1) + '万';
  return '¥' + amount.toLocaleString();
}

function formatNumber(num: number): string {
  if (num >= 10000) return (num / 10000).toFixed(1) + '万';
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K';
  return num.toString();
}

function formatDate(dateStr: string): string {
  if (!dateStr) return '-';
  const date = new Date(dateStr);
  return date.toLocaleDateString('zh-CN', { year: 'numeric', month: 'short', day: 'numeric' });
}

// 收入趋势
const revenueTrend = computed(() => {
  if (!model.value) return 0;
  const { thisMonth, lastMonth } = model.value.revenue;
  if (lastMonth === 0) return thisMonth > 0 ? 100 : 0;
  return Math.round(((thisMonth - lastMonth) / lastMonth) * 100);
});

// 最大月收入（用于图表）
const maxMonthRevenue = computed(() => {
  if (!model.value) return 1500;
  return Math.max(...model.value.revenue.history.map(h => h.amount), 1);
});
</script>

<template>
  <Page
    description="管理模型详情、查看数据分析"
    title="模型详情"
    content-class="p-6 bg-gray-50"
  >
    <!-- 加载状态 -->
    <div v-if="loading" class="space-y-6">
      <ElCard shadow="hover" class="!rounded-xl">
        <ElSkeleton animated>
          <template #template>
            <div class="flex gap-6">
              <ElSkeletonItem variant="image" style="width: 280px; height: 180px; border-radius: 12px" />
              <div class="flex-1 space-y-4">
                <ElSkeletonItem variant="h1" style="width: 60%" />
                <ElSkeletonItem variant="text" style="width: 80%" />
                <ElSkeletonItem variant="text" style="width: 40%" />
              </div>
            </div>
          </template>
        </ElSkeleton>
      </ElCard>
    </div>

    <!-- 模型详情 -->
    <template v-else-if="model">
      <!-- 顶部返回和操作 -->
      <div class="mb-6 flex items-center justify-between">
        <button
          class="flex items-center gap-2 text-gray-600 hover:text-purple-600 transition-colors"
          @click="goBack"
        >
          <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"/>
          </svg>
          <span>返回列表</span>
        </button>
        <div class="flex items-center gap-3">
          <ElButton
            v-if="model.status === 'draft'"
            type="primary"
            class="!bg-purple-600 !border-purple-600 hover:!bg-purple-700 !rounded-full"
            @click="handlePublish"
          >
            <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"/>
            </svg>
            提交审核
          </ElButton>
          <ElButton
            v-if="model.status === 'published'"
            plain
            class="!rounded-full"
            @click="handleUnpublish"
          >
            下架模型
          </ElButton>
          <ElButton class="!rounded-full" @click="handleEdit">
            <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
            </svg>
            编辑
          </ElButton>
          <ElButton type="danger" plain class="!rounded-full" @click="handleDelete">
            <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
            </svg>
          </ElButton>
        </div>
      </div>

      <div class="flex gap-6">
        <!-- 左侧主内容 -->
        <div class="flex-1 min-w-0 space-y-6">
          <!-- 模型头卡 -->
          <ElCard shadow="hover" class="!rounded-xl overflow-hidden">
            <div class="flex flex-col lg:flex-row gap-6">
              <!-- 封面 -->
              <div class="relative w-full lg:w-72 h-48 rounded-xl overflow-hidden bg-gradient-to-br from-purple-100 to-indigo-100 flex-shrink-0">
                <img
                  :src="model.cover || '/images/default-model-cover.svg'"
                  class="w-full h-full object-cover"
                  @error="(e) => { const img = e.target as HTMLImageElement; if (img) img.src = '/images/default-model-cover.svg'; }"
                />
                <!-- 状态标签 -->
                <span
                  :class="[
                    'absolute left-3 top-3 rounded-full px-3 py-1 text-sm font-medium',
                    getStatusStyle(model.status).bg,
                    getStatusStyle(model.status).text
                  ]"
                >
                  {{ getStatusStyle(model.status).label }}
                </span>
                <!-- 价格标签 -->
                <span
                  :class="[
                    'absolute right-3 top-3 rounded-full px-3 py-1 text-sm font-bold shadow-lg',
                    model.isFree ? 'bg-green-500 text-white' : 'bg-white text-purple-600'
                  ]"
                >
                  {{ model.isFree ? '免费' : '¥' + model.price }}
                </span>
              </div>

              <!-- 信息 -->
              <div class="flex-1">
                <h1 class="text-2xl font-bold text-gray-900 mb-2">{{ model.title }}</h1>
                <p class="text-gray-500 mb-4 leading-relaxed">{{ model.description }}</p>
                
                <!-- 标签 -->
                <div class="flex flex-wrap gap-2 mb-4">
                  <ElTag
                    v-for="tag in model.tags"
                    :key="tag"
                    effect="plain"
                    class="!bg-purple-50 !text-purple-600 !border-purple-200 !rounded-full"
                  >
                    {{ tag }}
                  </ElTag>
                </div>

                <!-- 统计栏 -->
                <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 pt-4 border-t border-gray-100">
                  <div class="text-center">
                    <div class="text-2xl font-bold text-purple-600">{{ formatNumber(model.stats.views) }}</div>
                    <div class="text-sm text-gray-500">浏览量</div>
                  </div>
                  <div class="text-center">
                    <div class="text-2xl font-bold text-blue-600">{{ formatNumber(model.stats.adoptions) }}</div>
                    <div class="text-sm text-gray-500">被采纳</div>
                  </div>
                  <div class="text-center">
                    <div class="text-2xl font-bold text-red-500">{{ formatNumber(model.stats.likes) }}</div>
                    <div class="text-sm text-gray-500">获赞</div>
                  </div>
                  <div class="text-center">
                    <div class="text-2xl font-bold text-green-600">{{ formatNumber(model.stats.reviews) }}</div>
                    <div class="text-sm text-gray-500">评价数</div>
                  </div>
                </div>
              </div>
            </div>
          </ElCard>

          <!-- Tab 导航 -->
          <div class="flex flex-wrap gap-2">
            <button
              v-for="tab in tabs"
              :key="tab.id"
              class="px-5 py-2.5 rounded-full text-sm font-medium transition-all"
              :class="[
                activeTab === tab.id
                  ? 'bg-purple-100 text-purple-700 shadow-md border border-purple-200 font-semibold'
                  : 'bg-white text-gray-600 hover:bg-purple-50 hover:text-purple-600 border border-gray-200'
              ]"
              @click="activeTab = tab.id as any"
            >
              {{ tab.icon }} {{ tab.label }}
            </button>
          </div>

          <!-- Tab 内容 -->
          <!-- 概述 -->
          <ElCard v-if="activeTab === 'overview'" shadow="hover" class="!rounded-xl">
            <template #header>
              <div class="flex items-center gap-2">
                <span class="text-lg">📋</span>
                <span class="font-semibold text-gray-700">模型概述</span>
              </div>
            </template>
            <div class="prose max-w-none">
              <p class="text-gray-600 leading-relaxed whitespace-pre-line">
                {{ model.content.overview }}
              </p>
            </div>
          </ElCard>

          <!-- 使用步骤 -->
          <ElCard v-if="activeTab === 'steps'" shadow="hover" class="!rounded-xl">
            <template #header>
              <div class="flex items-center gap-2">
                <span class="text-lg">📝</span>
                <span class="font-semibold text-gray-700">使用步骤</span>
              </div>
            </template>
            <div class="space-y-4">
              <div
                v-for="(step, index) in model.content.steps"
                :key="index"
                class="flex gap-4 p-4 rounded-xl bg-gradient-to-r from-purple-50 to-indigo-50 border border-purple-100"
              >
                <div class="w-10 h-10 rounded-full bg-purple-600 text-white flex items-center justify-center font-bold flex-shrink-0">
                  {{ index + 1 }}
                </div>
                <div class="flex-1">
                  <h4 class="font-semibold text-gray-800 mb-2">{{ step.title }}</h4>
                  <p class="text-sm text-gray-600 leading-relaxed">{{ step.description }}</p>
                </div>
              </div>
            </div>
          </ElCard>

          <!-- 案例 -->
          <ElCard v-if="activeTab === 'examples'" shadow="hover" class="!rounded-xl">
            <template #header>
              <div class="flex items-center gap-2">
                <span class="text-lg">💡</span>
                <span class="font-semibold text-gray-700">实践案例</span>
              </div>
            </template>
            <div class="space-y-6">
              <div
                v-for="(example, index) in model.content.examples"
                :key="index"
                class="p-5 rounded-xl bg-gradient-to-br from-gray-50 to-slate-50 border border-gray-100"
              >
                <h4 class="font-semibold text-gray-800 mb-3 flex items-center gap-2">
                  <span class="w-6 h-6 rounded-full bg-amber-100 text-amber-600 flex items-center justify-center text-sm">
                    {{ index + 1 }}
                  </span>
                  {{ example.title }}
                </h4>
                <p class="text-sm text-gray-600 leading-relaxed whitespace-pre-line">{{ example.content }}</p>
              </div>
            </div>
          </ElCard>

          <!-- 用户反馈 -->
          <ElCard v-if="activeTab === 'feedback'" shadow="hover" class="!rounded-xl">
            <template #header>
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <span class="text-lg">💬</span>
                  <span class="font-semibold text-gray-700">用户反馈</span>
                </div>
                <span class="text-sm text-gray-400">共 {{ model.stats.reviews }} 条评价</span>
              </div>
            </template>
            <div class="space-y-4">
              <div
                v-for="feedback in mockFeedbacks"
                :key="feedback.id"
                class="p-4 rounded-xl bg-gray-50 border border-gray-100"
              >
                <div class="flex items-start gap-3">
                  <ElAvatar size="small" :style="{ backgroundColor: '#7c3aed' }">
                    {{ feedback.user.charAt(0) }}
                  </ElAvatar>
                  <div class="flex-1">
                    <div class="flex items-center justify-between mb-1">
                      <span class="font-medium text-gray-800">{{ feedback.user }}</span>
                      <div class="flex items-center gap-1">
                        <span v-for="i in feedback.rating" :key="i" class="text-amber-400 text-sm">★</span>
                        <span v-for="i in (5 - feedback.rating)" :key="'e' + i" class="text-gray-300 text-sm">★</span>
                      </div>
                    </div>
                    <p class="text-sm text-gray-600">{{ feedback.content }}</p>
                    <p class="text-xs text-gray-400 mt-2">{{ feedback.date }}</p>
                  </div>
                </div>
              </div>
            </div>
          </ElCard>

          <!-- 数据分析 -->
          <ElCard v-if="activeTab === 'analytics'" shadow="hover" class="!rounded-xl">
            <template #header>
              <div class="flex items-center gap-2">
                <span class="text-lg">📊</span>
                <span class="font-semibold text-gray-700">数据分析</span>
              </div>
            </template>
            <div class="space-y-6">
              <!-- 关键指标 -->
              <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div class="p-4 rounded-xl bg-purple-50 border border-purple-100 text-center">
                  <div class="text-2xl font-bold text-purple-600 mb-1">{{ formatNumber(model.stats.views) }}</div>
                  <div class="text-sm text-gray-500">总浏览</div>
                </div>
                <div class="p-4 rounded-xl bg-blue-50 border border-blue-100 text-center">
                  <div class="text-2xl font-bold text-blue-600 mb-1">{{ ((model.stats.adoptions / model.stats.views) * 100).toFixed(1) }}%</div>
                  <div class="text-sm text-gray-500">采纳率</div>
                </div>
                <div class="p-4 rounded-xl bg-green-50 border border-green-100 text-center">
                  <div class="text-2xl font-bold text-green-600 mb-1">{{ ((model.stats.likes / model.stats.adoptions) * 100).toFixed(1) }}%</div>
                  <div class="text-sm text-gray-500">好评率</div>
                </div>
                <div class="p-4 rounded-xl bg-amber-50 border border-amber-100 text-center">
                  <div class="text-2xl font-bold text-amber-600 mb-1">{{ (model.stats.practices / model.stats.adoptions).toFixed(1) }}</div>
                  <div class="text-sm text-gray-500">人均练习</div>
                </div>
              </div>

              <!-- 趋势提示 -->
              <div class="p-4 rounded-xl bg-gradient-to-r from-purple-50 to-indigo-50 border border-purple-100">
                <h4 class="font-medium text-gray-800 mb-2">💡 数据洞察</h4>
                <ul class="text-sm text-gray-600 space-y-1">
                  <li>• 本模型采纳率高于平均水平 35%</li>
                  <li>• 用户平均练习次数持续增长</li>
                  <li>• 建议增加更多实战案例提升互动</li>
                </ul>
              </div>
            </div>
          </ElCard>
        </div>

        <!-- 右侧边栏 -->
        <div class="w-80 flex-shrink-0 space-y-6 hidden lg:block">
          <!-- 收入概览 -->
          <ElCard v-if="!model.isFree" shadow="hover" class="!rounded-xl !bg-gradient-to-br from-green-50 to-emerald-50 !border-green-100">
            <template #header>
              <div class="flex items-center gap-2">
                <span class="text-lg">💰</span>
                <span class="font-semibold text-gray-700">收入概览</span>
              </div>
            </template>
            <div class="space-y-4">
              <div class="text-center py-3">
                <div class="text-3xl font-bold text-green-600">{{ formatMoney(model.revenue.total) }}</div>
                <div class="text-sm text-gray-500">累计收入</div>
              </div>
              <div class="flex items-center justify-between p-3 bg-white rounded-lg">
                <span class="text-sm text-gray-500">本月收入</span>
                <div class="flex items-center gap-2">
                  <span class="font-semibold text-green-600">{{ formatMoney(model.revenue.thisMonth) }}</span>
                  <span
                    :class="[
                      'text-xs px-2 py-0.5 rounded-full',
                      revenueTrend >= 0 ? 'bg-green-100 text-green-600' : 'bg-red-100 text-red-600'
                    ]"
                  >
                    {{ revenueTrend >= 0 ? '+' : '' }}{{ revenueTrend }}%
                  </span>
                </div>
              </div>
              
              <!-- 收入趋势图 -->
              <div class="pt-4">
                <div class="text-sm text-gray-500 mb-3">近6个月趋势</div>
                <div class="flex items-end gap-2 h-24">
                  <div
                    v-for="(item, index) in model.revenue.history"
                    :key="index"
                    class="flex-1 flex flex-col items-center gap-1"
                  >
                    <div
                      class="w-full bg-gradient-to-t from-green-500 to-emerald-400 rounded-t transition-all duration-500"
                      :style="{ height: `${(item.amount / maxMonthRevenue) * 100}%` }"
                    />
                    <span class="text-xs text-gray-400">{{ item.month.slice(5) }}</span>
                  </div>
                </div>
              </div>
            </div>
          </ElCard>

          <!-- 模型表现 -->
          <ElCard shadow="hover" class="!rounded-xl">
            <template #header>
              <div class="flex items-center gap-2">
                <span class="text-lg">📈</span>
                <span class="font-semibold text-gray-700">模型表现</span>
              </div>
            </template>
            <div class="space-y-4">
              <div>
                <div class="flex items-center justify-between text-sm mb-2">
                  <span class="text-gray-500">采纳率</span>
                  <span class="text-purple-600 font-medium">{{ ((model.stats.adoptions / model.stats.views) * 100).toFixed(1) }}%</span>
                </div>
                <ElProgress :percentage="Math.min(100, Math.round((model.stats.adoptions / model.stats.views) * 100))" :stroke-width="8" color="#7c3aed" :show-text="false" />
              </div>
              <div>
                <div class="flex items-center justify-between text-sm mb-2">
                  <span class="text-gray-500">好评率</span>
                  <span class="text-green-600 font-medium">{{ ((model.stats.likes / model.stats.adoptions) * 100).toFixed(1) }}%</span>
                </div>
                <ElProgress :percentage="Math.min(100, Math.round((model.stats.likes / model.stats.adoptions) * 100))" :stroke-width="8" color="#10b981" :show-text="false" />
              </div>
              <div>
                <div class="flex items-center justify-between text-sm mb-2">
                  <span class="text-gray-500">复购率</span>
                  <span class="text-blue-600 font-medium">45%</span>
                </div>
                <ElProgress :percentage="45" :stroke-width="8" color="#3b82f6" :show-text="false" />
              </div>
            </div>
          </ElCard>

          <!-- 基本信息 -->
          <ElCard shadow="hover" class="!rounded-xl">
            <template #header>
              <div class="flex items-center gap-2">
                <span class="text-lg">ℹ️</span>
                <span class="font-semibold text-gray-700">基本信息</span>
              </div>
            </template>
            <div class="space-y-3 text-sm">
              <div class="flex items-center justify-between">
                <span class="text-gray-500">分类</span>
                <span class="text-gray-800">{{ model.categoryName }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-gray-500">创建时间</span>
                <span class="text-gray-800">{{ formatDate(model.createdAt) }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-gray-500">最后更新</span>
                <span class="text-gray-800">{{ formatDate(model.updatedAt) }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-gray-500">模型 ID</span>
                <span class="text-gray-400 font-mono text-xs">{{ model.id }}</span>
              </div>
            </div>
          </ElCard>

          <!-- 操作提示 -->
          <ElCard shadow="hover" class="!rounded-xl !bg-gradient-to-br from-amber-50 to-orange-50 !border-amber-100">
            <template #header>
              <div class="flex items-center gap-2">
                <span class="text-lg">💡</span>
                <span class="font-semibold text-amber-700">优化建议</span>
              </div>
            </template>
            <ul class="text-sm text-amber-800 space-y-2">
              <li class="flex items-start gap-2">
                <span class="text-amber-500 mt-1">•</span>
                <span>添加更多实战案例可提升采纳率</span>
              </li>
              <li class="flex items-start gap-2">
                <span class="text-amber-500 mt-1">•</span>
                <span>定期更新内容保持模型活力</span>
              </li>
              <li class="flex items-start gap-2">
                <span class="text-amber-500 mt-1">•</span>
                <span>回复用户反馈可提升好评率</span>
              </li>
            </ul>
          </ElCard>
        </div>
      </div>
    </template>

    <!-- 空状态 -->
    <ElCard v-else shadow="hover" class="!rounded-xl">
      <ElEmpty description="模型不存在或已被删除">
        <template #image>
          <div class="text-6xl">🔍</div>
        </template>
        <ElButton type="primary" class="!bg-purple-600 !border-purple-600 !rounded-full mt-4" @click="goBack">
          返回列表
        </ElButton>
      </ElEmpty>
    </ElCard>
  </Page>
</template>

<style scoped>
.prose {
  color: #374151;
}
</style>
