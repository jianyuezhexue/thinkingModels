<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { Page } from '@vben/common-ui';

import {
  ElButton,
  ElCard,
  ElTag,
  ElMessage,
  ElTabs,
  ElTabPane,
  ElEmpty,
  ElSkeleton,
  ElSkeletonItem,
  ElDialog,
  ElInput,
} from 'element-plus';

// 路由
const route = useRoute();
const router = useRouter();
const modelId = computed(() => route.params.id as string);
const isEditMode = computed(() => route.query.edit === 'true');

// 加载状态
const loading = ref(true);

// 模型数据
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

// 删除确认对话框
const deleteDialogVisible = ref(false);

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
    overview: 'SWOT 分析是一种战略规划工具，用于评估企业、项目或个人的优势（Strengths）、劣势（Weaknesses）、机会（Opportunities）和威胁（Threats）。这种方法帮助决策者全面了解内外部环境，制定更有效的战略。',
    steps: [
      {
        title: '识别优势 (Strengths)',
        description: '列出你或你的组织相对于竞争对手的优势。包括资源、能力、经验、品牌等内部因素。',
      },
      {
        title: '识别劣势 (Weaknesses)',
        description: '诚实地列出需要改进的领域。这些是你相对于竞争对手的不足之处。',
      },
      {
        title: '发现机会 (Opportunities)',
        description: '分析外部环境中的有利因素。包括市场趋势、政策变化、技术发展等。',
      },
      {
        title: '识别威胁 (Threats)',
        description: '评估可能对你产生负面影响的外部因素。包括竞争、经济环境、法规变化等。',
      },
    ],
    examples: [
      {
        title: '案例：某电商平台的 SWOT 分析',
        content: '优势：拥有庞大的用户基础、完善的物流体系、强大的技术团队。劣势：运营成本较高、对第三方商家管控力有限。机会：下沉市场增长、跨境电商发展、直播带货兴起。威胁：竞争对手价格战、监管政策趋严、用户获取成本上升。',
      },
      {
        title: '案例：个人职业发展的 SWOT 分析',
        content: '优势：专业技能扎实、沟通能力强、学习能力快。劣势：管理经验不足、人脉资源有限。机会：行业快速发展、新技术带来新岗位、公司有晋升机会。威胁：行业竞争加剧、技术更新快、AI 可能替代部分工作。',
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
    // TODO: 调用实际 API
    // const res = await getMyModelDetailApi(modelId.value);
    await new Promise(resolve => setTimeout(resolve, 800));
    model.value = mockModelDetail;
  } catch (error) {
    console.error('获取模型详情失败:', error);
    ElMessage.error('获取模型详情失败');
  } finally {
    loading.value = false;
  }
}

// 页面加载
onMounted(() => {
  fetchModelDetail();
});

// 返回列表
function goBack() {
  router.push('/my-models');
}

// 编辑模型
function handleEdit() {
  router.push(`/my-models/${modelId.value}?edit=true`);
}

// 删除模型
async function handleDelete() {
  deleteDialogVisible.value = true;
}

// 确认删除
async function confirmDelete() {
  try {
    // await deleteModelApi(modelId.value);
    ElMessage.success('模型已删除');
    deleteDialogVisible.value = false;
    router.push('/my-models');
  } catch (error) {
    console.error('删除失败:', error);
    ElMessage.error('删除失败');
  }
}

// 发布模型
async function handlePublish() {
  try {
    // await publishModelApi(modelId.value);
    ElMessage.success('模型已提交审核');
    fetchModelDetail();
  } catch (error) {
    console.error('发布失败:', error);
    ElMessage.error('发布失败');
  }
}

// 格式化金额
function formatMoney(amount: number): string {
  if (amount >= 10000) {
    return '¥' + (amount / 10000).toFixed(1) + '万';
  }
  return '¥' + amount.toLocaleString();
}

// 格式化数字
function formatNumber(num: number): string {
  if (num >= 10000) return (num / 10000).toFixed(1) + '万';
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K';
  return num.toString();
}

// 格式化日期
function formatDate(dateStr: string): string {
  if (!dateStr) return '-';
  const date = new Date(dateStr);
  return date.toLocaleDateString('zh-CN');
}

// 计算收入趋势
const revenueTrend = computed(() => {
  if (!model.value) return 0;
  const { thisMonth, lastMonth } = model.value.revenue;
  if (lastMonth === 0) return 100;
  return Math.round(((thisMonth - lastMonth) / lastMonth) * 100);
});

// 获取状态标签类型
function getStatusType(status: string): any {
  const map: Record<string, any> = {
    published: 'success',
    draft: 'info',
    under_review: 'warning',
    rejected: 'danger',
  };
  return map[status] || '';
}
</script>

<template>
  <Page
    :description="model?.description || '模型详情'"
    :title="model?.title || '模型详情'"
  >
    <!-- 加载状态 -->
    <div v-if="loading" class="space-y-6">
      <ElSkeleton animated>
        <template #template>
          <div class="flex gap-6">
            <ElSkeletonItem variant="image" style="width: 300px; height: 200px" />
            <div class="flex-1 space-y-4">
              <ElSkeletonItem variant="p" style="width: 60%" />
              <ElSkeletonItem variant="text" style="width: 80%" />
              <ElSkeletonItem variant="text" style="width: 40%" />
            </div>
          </div>
        </template>
      </ElSkeleton>
    </div>

    <!-- 模型详情 -->
    <template v-else-if="model">
      <!-- 顶部操作栏 -->
      <div class="mb-6 flex items-center justify-between">
        <ElButton @click="goBack">
          <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"/>
          </svg>
          返回列表
        </ElButton>
        <div class="flex gap-3">
          <ElButton v-if="model.status === 'draft'" type="primary" @click="handlePublish">
            提交审核
          </ElButton>
          <ElButton type="primary" plain @click="handleEdit">
            <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
            </svg>
            编辑
          </ElButton>
          <ElButton type="danger" plain @click="handleDelete">
            <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
            </svg>
            删除
          </ElButton>
        </div>
      </div>

      <!-- 基本信息卡片 -->
      <ElCard shadow="never" class="mb-6">
        <div class="flex flex-col lg:flex-row gap-6">
          <!-- 封面图 -->
          <div class="shrink-0">
            <img
              :src="model.cover || '/images/default-model-cover.svg'"
              class="w-full lg:w-80 h-48 object-cover rounded-lg"
              @error="(e) => { const img = e.target as HTMLImageElement; if (img) img.src = '/images/default-model-cover.svg'; }"
            />
          </div>

          <!-- 信息区 -->
          <div class="flex-1">
            <div class="flex items-start justify-between">
              <div>
                <div class="flex items-center gap-3 mb-2">
                  <h1 class="text-2xl font-bold text-gray-900">{{ model.title }}</h1>
                  <ElTag :type="getStatusType(model.status)">{{ model.statusText }}</ElTag>
                  <span
                    :class="[
                      'rounded-full px-2 py-0.5 text-xs font-bold',
                      model.isFree ? 'bg-green-100 text-green-700' : 'bg-purple-100 text-purple-700'
                    ]"
                  >
                    {{ model.isFree ? '免费' : '¥' + model.price }}
                  </span>
                </div>
                <p class="text-gray-500 mb-4">{{ model.description }}</p>
                <div class="flex flex-wrap gap-2 mb-4">
                  <span
                    v-for="tag in model.tags"
                    :key="tag"
                    class="rounded px-2 py-0.5 text-xs"
                    :style="{ backgroundColor: 'var(--el-color-primary-light-9)', color: 'var(--el-color-primary)' }"
                  >
                    {{ tag }}
                  </span>
                </div>
              </div>
            </div>

            <div class="grid grid-cols-2 md:grid-cols-4 gap-4 pt-4 border-t border-gray-100">
              <div>
                <p class="text-sm text-gray-500">被采纳</p>
                <p class="text-lg font-semibold text-gray-900">{{ formatNumber(model.stats.adoptions) }}</p>
              </div>
              <div>
                <p class="text-sm text-gray-500">练习次数</p>
                <p class="text-lg font-semibold text-gray-900">{{ formatNumber(model.stats.practices) }}</p>
              </div>
              <div>
                <p class="text-sm text-gray-500">获得点赞</p>
                <p class="text-lg font-semibold text-gray-900">{{ formatNumber(model.stats.likes) }}</p>
              </div>
              <div v-if="!model.isFree">
                <p class="text-sm text-gray-500">累计收入</p>
                <p class="text-lg font-semibold text-green-600">{{ formatMoney(model.revenue.total) }}</p>
              </div>
            </div>
          </div>
        </div>
      </ElCard>

      <!-- 收入分析 (仅付费模型) -->
      <ElCard v-if="!model.isFree" shadow="never" class="mb-6">
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold">收入分析</h3>
            <div class="flex items-center gap-2 text-sm">
              <span class="text-gray-500">本月收入:</span>
              <span class="font-semibold text-green-600">{{ formatMoney(model.revenue.thisMonth) }}</span>
              <span
                :class="[
                  'text-xs px-2 py-0.5 rounded',
                  revenueTrend >= 0 ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'
                ]"
              >
                {{ revenueTrend >= 0 ? '↑' : '↓' }} {{ Math.abs(revenueTrend) }}%
              </span>
            </div>
          </div>
        </template>
        <div class="grid md:grid-cols-3 gap-6">
          <div class="bg-gradient-to-br from-purple-50 to-pink-50 rounded-lg p-4">
            <p class="text-sm text-gray-600 mb-1">累计收入</p>
            <p class="text-2xl font-bold text-purple-700">{{ formatMoney(model.revenue.total) }}</p>
            <p class="text-xs text-gray-400 mt-2">扣除平台服务费后</p>
          </div>
          <div class="bg-gradient-to-br from-blue-50 to-cyan-50 rounded-lg p-4">
            <p class="text-sm text-gray-600 mb-1">本月收入</p>
            <p class="text-2xl font-bold text-blue-700">{{ formatMoney(model.revenue.thisMonth) }}</p>
            <p class="text-xs text-gray-400 mt-2">较上月 {{ revenueTrend >= 0 ? '增长' : '下降' }} {{ Math.abs(revenueTrend) }}%</p>
          </div>
          <div class="bg-gradient-to-br from-green-50 to-emerald-50 rounded-lg p-4">
            <p class="text-sm text-gray-600 mb-1">付费采纳</p>
            <p class="text-2xl font-bold text-green-700">{{ formatNumber(Math.floor(model.stats.adoptions * 0.6)) }}</p>
            <p class="text-xs text-gray-400 mt-2">转化率 60%</p>
          </div>
        </div>

        <!-- 收入趋势图占位 -->
        <div class="mt-6 p-4 bg-gray-50 rounded-lg">
          <p class="text-sm text-gray-500 mb-4">近 6 个月收入趋势</p>
          <div class="flex items-end gap-4 h-32">
            <div
              v-for="(item, index) in model.revenue.history"
              :key="index"
              class="flex-1 flex flex-col items-center gap-2"
            >
              <div
                class="w-full bg-purple-500 rounded-t transition-all duration-500"
                :style="{ height: `${(item.amount / 1500) * 100}%` }"
              />
              <span class="text-xs text-gray-500">{{ item.month.slice(5) }}月</span>
            </div>
          </div>
        </div>
      </ElCard>

      <!-- 详细内容 -->
      <ElCard shadow="never">
        <ElTabs>
          <ElTabPane label="模型内容">
            <div class="space-y-6">
              <!-- 概述 -->
              <div>
                <h3 class="text-lg font-semibold text-gray-900 mb-3">模型概述</h3>
                <p class="text-gray-600 leading-relaxed">{{ model.content.overview }}</p>
              </div>

              <!-- 使用步骤 -->
              <div>
                <h3 class="text-lg font-semibold text-gray-900 mb-4">使用步骤</h3>
                <div class="space-y-4">
                  <div
                    v-for="(step, index) in model.content.steps"
                    :key="index"
                    class="flex gap-4 p-4 rounded-lg bg-gray-50"
                  >
                    <div class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-purple-100 text-sm font-medium text-purple-600">
                      {{ index + 1 }}
                    </div>
                    <div>
                      <h4 class="font-medium text-gray-900">{{ step.title }}</h4>
                      <p class="text-sm text-gray-500 mt-1">{{ step.description }}</p>
                    </div>
                  </div>
                </div>
              </div>

              <!-- 实践案例 -->
              <div>
                <h3 class="text-lg font-semibold text-gray-900 mb-4">实践案例</h3>
                <div class="space-y-4">
                  <ElCard
                    v-for="(example, index) in model.content.examples"
                    :key="index"
                    shadow="never"
                    class="bg-gray-50"
                  >
                    <template #header>
                      <span class="font-medium">{{ example.title }}</span>
                    </template>
                    <p class="text-gray-600 text-sm">{{ example.content }}</p>
                  </ElCard>
                </div>
              </div>
            </div>
          </ElTabPane>

          <ElTabPane label="用户反馈">
            <div class="text-center py-12">
              <svg class="h-16 w-16 mx-auto text-gray-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/>
              </svg>
              <p class="text-gray-500">用户反馈功能即将上线</p>
            </div>
          </ElTabPane>

          <ElTabPane label="数据分析">
            <div class="text-center py-12">
              <svg class="h-16 w-16 mx-auto text-gray-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
              </svg>
              <p class="text-gray-500">详细数据分析功能即将上线</p>
            </div>
          </ElTabPane>
        </ElTabs>
      </ElCard>
    </template>

    <!-- 删除确认对话框 -->
    <ElDialog
      v-model="deleteDialogVisible"
      title="确认删除"
      width="400px"
    >
      <p class="text-gray-600">
        确定要删除模型 <strong>{{ model?.title }}</strong> 吗？此操作不可恢复。
      </p>
      <template #footer>
        <ElButton @click="deleteDialogVisible = false">取消</ElButton>
        <ElButton type="danger" @click="confirmDelete">确认删除</ElButton>
      </template>
    </ElDialog>
  </Page>
</template>
