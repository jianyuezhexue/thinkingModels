<script lang="ts" setup>
import { computed, onMounted, ref, watch } from 'vue';
import { useRouter } from 'vue-router';

import { Page } from '@vben/common-ui';

import {
  ElButton,
  ElCard,
  ElInput,
  ElSelect,
  ElOption,
  ElTag,
  ElMessage,
  ElEmpty,
  ElSkeleton,
  ElSkeletonItem,
  ElPagination,
} from 'element-plus';

// 类型定义
interface MyModel {
  id: string;
  title: string;
  description: string;
  cover: string;
  category: string;
  tags: string[];
  status: 'published' | 'draft' | 'under_review' | 'rejected';
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
  };
  createdAt: string;
  updatedAt: string;
}

interface StatsSummary {
  totalModels: number;
  publishedModels: number;
  totalRevenue: number;
  monthlyRevenue: number;
  totalAdoptions: number;
  totalLikes: number;
}

// 路由
const router = useRouter();

// 加载状态
const loading = ref(false);
const statsLoading = ref(false);

// 模型列表
const models = ref<MyModel[]>([]);
const total = ref(0);

// 统计数据
const stats = ref<StatsSummary>({
  totalModels: 0,
  publishedModels: 0,
  totalRevenue: 0,
  monthlyRevenue: 0,
  totalAdoptions: 0,
  totalLikes: 0,
});

// 分页
const currentPage = ref(1);
const pageSize = ref(12);

// 筛选
const searchKeyword = ref('');
const statusFilter = ref<'all' | 'published' | 'draft' | 'under_review'>('all');
const sortBy = ref<'newest' | 'popular' | 'revenue' | 'adoptions'>('newest');

// 状态选项
const statusOptions = [
  { value: 'all', label: '全部状态' },
  { value: 'published', label: '已发布', type: 'success' },
  { value: 'draft', label: '草稿', type: 'info' },
  { value: 'under_review', label: '审核中', type: 'warning' },
];

// 排序选项
const sortOptions = [
  { value: 'newest', label: '最新创建' },
  { value: 'popular', label: '最受欢迎' },
  { value: 'revenue', label: '收入最高' },
  { value: 'adoptions', label: '采纳最多' },
];

// 模拟数据 - 实际项目中应该从 API 获取
const mockStats: StatsSummary = {
  totalModels: 12,
  publishedModels: 8,
  totalRevenue: 15860,
  monthlyRevenue: 2340,
  totalAdoptions: 3456,
  totalLikes: 892,
};

const mockModels: MyModel[] = [
  {
    id: '1',
    title: 'SWOT 分析思维模型',
    description: '经典的战略分析工具，帮助分析企业或项目的优势、劣势、机会和威胁，适用于商业决策和个人发展规划。',
    cover: '/images/swot-cover.svg',
    category: 'business',
    tags: ['战略', '分析', '商业'],
    status: 'published',
    price: 29,
    isFree: false,
    stats: { adoptions: 1256, practices: 3421, likes: 328, reviews: 56 },
    revenue: { total: 8560, thisMonth: 1200, lastMonth: 980 },
    createdAt: '2024-01-15T08:00:00Z',
    updatedAt: '2024-02-10T10:30:00Z',
  },
  {
    id: '2',
    title: '第一性原理思维',
    description: '埃隆·马斯克推崇的创新思维方式，从最基本的原理出发思考问题，打破常规思维定式。',
    cover: '/images/first-principles-cover.svg',
    category: 'innovation',
    tags: ['创新', '思维', '决策'],
    status: 'published',
    price: 0,
    isFree: true,
    stats: { adoptions: 2100, practices: 5678, likes: 445, reviews: 89 },
    revenue: { total: 0, thisMonth: 0, lastMonth: 0 },
    createdAt: '2024-01-20T14:00:00Z',
    updatedAt: '2024-02-08T16:45:00Z',
  },
  {
    id: '3',
    title: '逆向思维模型',
    description: '从结果倒推过程的思维方式，帮助发现潜在风险和盲点，特别适用于项目规划和风险管理。',
    cover: '/images/reverse-thinking-cover.svg',
    category: 'strategy',
    tags: ['策略', '风险', '规划'],
    status: 'published',
    price: 19,
    isFree: false,
    stats: { adoptions: 567, practices: 1234, likes: 156, reviews: 23 },
    revenue: { total: 2340, thisMonth: 450, lastMonth: 380 },
    createdAt: '2024-02-01T09:00:00Z',
    updatedAt: '2024-02-12T11:20:00Z',
  },
  {
    id: '4',
    title: '金字塔原理',
    description: 'MECE法则和金字塔结构的表达方法，让你的思考更有逻辑，表达更清晰有力。',
    cover: '/images/pyramid-cover.svg',
    category: 'analysis',
    tags: ['逻辑', '表达', '沟通'],
    status: 'draft',
    price: 39,
    isFree: false,
    stats: { adoptions: 0, practices: 0, likes: 0, reviews: 0 },
    revenue: { total: 0, thisMonth: 0, lastMonth: 0 },
    createdAt: '2024-02-15T10:00:00Z',
    updatedAt: '2024-02-15T10:00:00Z',
  },
  {
    id: '5',
    title: '系统思维模型',
    description: '从全局视角理解复杂系统，识别反馈回路和杠杆点，适用于复杂问题分析和组织管理。',
    cover: '/images/system-thinking-cover.svg',
    category: 'analysis',
    tags: ['系统', '复杂', '管理'],
    status: 'under_review',
    price: 49,
    isFree: false,
    stats: { adoptions: 0, practices: 0, likes: 0, reviews: 0 },
    revenue: { total: 0, thisMonth: 0, lastMonth: 0 },
    createdAt: '2024-02-18T13:00:00Z',
    updatedAt: '2024-02-18T13:00:00Z',
  },
  {
    id: '6',
    title: '机会成本分析',
    description: '帮助理解决策的真实代价，评估不同选择的隐性成本，做出更明智的资源配置决策。',
    cover: '/images/opportunity-cost-cover.svg',
    category: 'decision',
    tags: ['决策', '成本', '经济'],
    status: 'published',
    price: 15,
    isFree: false,
    stats: { adoptions: 533, practices: 1456, likes: 98, reviews: 18 },
    revenue: { total: 1600, thisMonth: 280, lastMonth: 320 },
    createdAt: '2024-02-05T11:00:00Z',
    updatedAt: '2024-02-11T09:15:00Z',
  },
];

// 获取统计数据
async function fetchStats() {
  statsLoading.value = true;
  try {
    // TODO: 替换为实际 API 调用
    // const res = await getMyModelStatsApi();
    await new Promise(resolve => setTimeout(resolve, 500));
    stats.value = mockStats;
  } catch (error) {
    console.error('获取统计数据失败:', error);
  } finally {
    statsLoading.value = false;
  }
}

// 获取模型列表
async function fetchModels() {
  loading.value = true;
  try {
    // TODO: 替换为实际 API 调用
    // const params = {
    //   page: currentPage.value,
    //   pageSize: pageSize.value,
    //   keyword: searchKeyword.value || undefined,
    //   status: statusFilter.value === 'all' ? undefined : statusFilter.value,
    //   sortBy: sortBy.value,
    // };
    // const res = await getMyModelsApi(params);
    await new Promise(resolve => setTimeout(resolve, 800));

    // 模拟筛选和排序
    let filtered = [...mockModels];
    if (statusFilter.value !== 'all') {
      filtered = filtered.filter(m => m.status === statusFilter.value);
    }
    if (searchKeyword.value) {
      const kw = searchKeyword.value.toLowerCase();
      filtered = filtered.filter(m =>
        m.title.toLowerCase().includes(kw) ||
        m.description.toLowerCase().includes(kw)
      );
    }

    // 模拟分页
    const start = (currentPage.value - 1) * pageSize.value;
    const end = start + pageSize.value;
    models.value = filtered.slice(start, end);
    total.value = filtered.length;
  } catch (error) {
    console.error('获取模型列表失败:', error);
    ElMessage.error('获取模型列表失败');
  } finally {
    loading.value = false;
  }
}

// 监听筛选条件变化
watch([searchKeyword, statusFilter, sortBy], () => {
  currentPage.value = 1;
  fetchModels();
});

// 监听分页变化
watch([currentPage, pageSize], () => {
  fetchModels();
});

// 页面加载时获取数据
onMounted(() => {
  fetchStats();
  fetchModels();
});

// 跳转到创建模型页面
function goToCreate() {
  router.push('/my-models/create');
}

// 跳转到模型详情
function goToDetail(model: MyModel) {
  router.push(`/my-models/${model.id}`);
}

// 编辑模型
function handleEdit(model: MyModel, event: Event) {
  event.stopPropagation();
  router.push(`/my-models/${model.id}?edit=true`);
}

// 删除模型
async function handleDelete(model: MyModel, event: Event) {
  event.stopPropagation();
  try {
    // TODO: 调用删除 API
    // await deleteModelApi(model.id);
    ElMessage.success('模型已删除');
    fetchModels();
    fetchStats();
  } catch (error) {
    console.error('删除失败:', error);
    ElMessage.error('删除失败');
  }
}

// 发布模型
async function handlePublish(model: MyModel, event: Event) {
  event.stopPropagation();
  try {
    // TODO: 调用发布 API
    // await publishModelApi(model.id);
    ElMessage.success('模型已提交审核');
    fetchModels();
  } catch (error) {
    console.error('发布失败:', error);
    ElMessage.error('发布失败');
  }
}

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

// 获取状态标签文本
function getStatusText(status: string): string {
  const map: Record<string, string> = {
    published: '已发布',
    draft: '草稿',
    under_review: '审核中',
    rejected: '已驳回',
  };
  return map[status] || status;
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
</script>

<template>
  <Page
    description="管理你创建的思维模型，跟踪收入和用户反馈"
    title="我的模型"
  >
    <!-- 统计概览卡片 -->
    <div class="mb-6 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <!-- 总模型数 -->
      <ElCard shadow="never" class="stat-card">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-500">总模型数</p>
            <p class="mt-1 text-2xl font-bold text-gray-900">
              {{ statsLoading ? '-' : stats.totalModels }}
            </p>
            <p class="mt-1 text-xs text-gray-400">
              {{ stats.publishedModels }} 个已发布
            </p>
          </div>
          <div class="rounded-lg bg-purple-50 p-3">
            <svg class="h-6 w-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"/>
            </svg>
          </div>
        </div>
      </ElCard>

      <!-- 总收入 -->
      <ElCard shadow="never" class="stat-card">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-500">累计收入</p>
            <p class="mt-1 text-2xl font-bold text-green-600">
              {{ statsLoading ? '-' : formatMoney(stats.totalRevenue) }}
            </p>
            <p class="mt-1 text-xs" :class="stats.monthlyRevenue > 0 ? 'text-green-500' : 'text-gray-400'">
              本月 +{{ formatMoney(stats.monthlyRevenue) }}
            </p>
          </div>
          <div class="rounded-lg bg-green-50 p-3">
            <svg class="h-6 w-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
          </div>
        </div>
      </ElCard>

      <!-- 总采纳数 -->
      <ElCard shadow="never" class="stat-card">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-500">被采纳次数</p>
            <p class="mt-1 text-2xl font-bold text-blue-600">
              {{ statsLoading ? '-' : formatNumber(stats.totalAdoptions) }}
            </p>
            <p class="mt-1 text-xs text-blue-400">
              帮助了 {{ Math.floor(stats.totalAdoptions * 0.7) }}+ 位思考者
            </p>
          </div>
          <div class="rounded-lg bg-blue-50 p-3">
            <svg class="h-6 w-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
            </svg>
          </div>
        </div>
      </ElCard>

      <!-- 总点赞数 -->
      <ElCard shadow="never" class="stat-card">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-500">获得点赞</p>
            <p class="mt-1 text-2xl font-bold text-pink-600">
              {{ statsLoading ? '-' : formatNumber(stats.totalLikes) }}
            </p>
            <p class="mt-1 text-xs text-pink-400">
              好评率 98.5%
            </p>
          </div>
          <div class="rounded-lg bg-pink-50 p-3">
            <svg class="h-6 w-6 text-pink-600" fill="currentColor" viewBox="0 0 24 24">
              <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
            </svg>
          </div>
        </div>
      </ElCard>
    </div>

    <!-- 模型列表 -->
    <ElCard shadow="never">
      <template #header>
        <div class="flex flex-wrap items-center justify-between gap-4">
          <div class="flex items-center gap-4">
            <h2 class="text-lg font-semibold">模型列表</h2>
            <ElSelect v-model="statusFilter" style="width: 140px">
              <ElOption
                v-for="opt in statusOptions"
                :key="opt.value"
                :label="opt.label"
                :value="opt.value"
              />
            </ElSelect>
            <ElSelect v-model="sortBy" style="width: 140px">
              <ElOption
                v-for="opt in sortOptions"
                :key="opt.value"
                :label="opt.label"
                :value="opt.value"
              />
            </ElSelect>
          </div>
          <div class="flex items-center gap-3">
            <ElInput
              v-model="searchKeyword"
              placeholder="搜索模型..."
              clearable
              style="width: 220px"
              @keyup.enter="fetchModels"
            >
              <template #prefix>
                <svg class="h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
                </svg>
              </template>
            </ElInput>
            <ElButton type="primary" @click="goToCreate">
              <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
              </svg>
              创建模型
            </ElButton>
          </div>
        </div>
      </template>

      <!-- 加载状态 -->
      <div v-if="loading" class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
        <ElSkeleton v-for="i in 8" :key="i" animated>
          <template #template>
            <div class="space-y-3">
              <ElSkeletonItem variant="image" style="width: 100%; height: 160px" />
              <ElSkeletonItem variant="p" style="width: 70%" />
              <ElSkeletonItem variant="text" style="width: 50%" />
            </div>
          </template>
        </ElSkeleton>
      </div>

      <!-- 空状态 -->
      <ElEmpty v-else-if="models.length === 0" description="暂无模型，创建你的第一个思维模型吧！">
        <ElButton type="primary" @click="goToCreate">创建模型</ElButton>
      </ElEmpty>

      <!-- 模型网格 -->
      <div v-else class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
        <div
          v-for="model in models"
          :key="model.id"
          class="group cursor-pointer overflow-hidden rounded-lg bg-white shadow-sm transition-all hover:shadow-lg border border-gray-100"
          @click="goToDetail(model)"
        >
          <!-- 封面图 -->
          <div class="relative h-40 overflow-hidden">
            <img
              :src="model.cover || '/images/default-model-cover.svg'"
              class="h-full w-full object-cover transition-transform group-hover:scale-105"
              @error="(e) => { const img = e.target as HTMLImageElement; if (img) img.src = '/images/default-model-cover.svg'; }"
            />
            <!-- 状态标签 -->
            <ElTag
              :type="getStatusType(model.status)"
              size="small"
              class="absolute left-2 top-2"
              effect="dark"
            >
              {{ getStatusText(model.status) }}
            </ElTag>
            <!-- 价格标签 -->
            <span
              :class="[
                'absolute right-2 top-2 rounded-full px-2 py-0.5 text-xs font-bold',
                model.isFree ? 'bg-green-500 text-white' : 'bg-white text-purple-600'
              ]"
            >
              {{ model.isFree ? '免费' : '¥' + model.price }}
            </span>
          </div>

          <!-- 内容 -->
          <div class="p-4">
            <h3 class="font-semibold text-gray-900 group-hover:text-purple-600 transition-colors">{{ model.title }}</h3>
            <p class="mt-1 line-clamp-2 h-10 text-sm text-gray-500">{{ model.description }}</p>

            <!-- 标签 -->
            <div class="mt-2 flex flex-wrap gap-1">
              <span
                v-for="tag in model.tags.slice(0, 2)"
                :key="tag"
                class="rounded px-2 py-0.5 text-xs"
                :style="{ backgroundColor: 'var(--el-color-primary-light-9)', color: 'var(--el-color-primary)' }"
              >
                {{ tag }}
              </span>
            </div>

            <!-- 收入展示 (仅已发布且付费模型) -->
            <div v-if="model.status === 'published' && !model.isFree" class="mt-3 flex items-center gap-3 rounded-lg bg-gradient-to-r from-purple-50 to-pink-50 px-3 py-2">
              <div class="flex items-center gap-1 text-sm font-medium text-purple-700">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
                {{ formatMoney(model.revenue.total) }}
              </div>
              <div class="text-xs text-gray-400">
                本月 +{{ formatMoney(model.revenue.thisMonth) }}
              </div>
            </div>

            <!-- 统计数据 -->
            <div class="mt-3 flex items-center gap-3 border-t border-gray-100 pt-3 text-xs text-gray-500">
              <span class="flex items-center gap-1" title="采纳">
                <svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
                </svg>
                {{ formatNumber(model.stats.adoptions) }}
              </span>
              <span class="flex items-center gap-1" title="练习">
                <svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"/>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
                {{ formatNumber(model.stats.practices) }}
              </span>
              <span class="flex items-center gap-1" :style="{ color: 'var(--el-color-primary)' }" title="点赞">
                <svg class="h-3.5 w-3.5" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
                </svg>
                {{ formatNumber(model.stats.likes) }}
              </span>
            </div>

            <!-- 操作按钮 -->
            <div class="mt-3 flex gap-2">
              <ElButton
                v-if="model.status === 'draft'"
                type="primary"
                class="flex-1"
                size="small"
                @click.stop="handlePublish(model, $event)"
              >
                提交审核
              </ElButton>
              <ElButton
                v-else-if="model.status === 'published'"
                type="primary"
                class="flex-1"
                size="small"
                @click.stop="handleEdit(model, $event)"
              >
                编辑模型
              </ElButton>
              <ElButton
                v-else
                type="info"
                class="flex-1"
                size="small"
                disabled
              >
                审核中...
              </ElButton>
              <ElButton
                size="small"
                @click.stop="handleDelete(model, $event)"
              >
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                </svg>
              </ElButton>
            </div>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <div v-if="total > pageSize" class="mt-8 flex justify-center">
        <ElPagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[8, 12, 24, 48]"
          layout="total, sizes, prev, pager, next"
        />
      </div>
    </ElCard>
  </Page>
</template>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.stat-card :deep(.el-card__body) {
  padding: 16px;
}
</style>
