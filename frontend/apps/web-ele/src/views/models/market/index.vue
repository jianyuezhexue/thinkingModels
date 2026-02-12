<script lang="ts" setup>
import { computed, ref, watch } from 'vue';
import { useRouter } from 'vue-router';

import { Page } from '@vben/common-ui';

import { ElButton, ElCard, ElInput, ElSelect, ElOption } from 'element-plus';

// Types
interface ModelAuthor {
  id: string;
  name: string;
  avatar: string;
}

interface ThinkingModel {
  id: string;
  title: string;
  description: string;
  cover: string;
  author: ModelAuthor;
  isFree: boolean;
  price?: number;
  category: string;
  tags: string[];
  stats: {
    adoptions: number;
    practices: number;
    discussions: number;
    forks: number;
    likes: number;
  };
  updatedAt: string;
}

// Filter states
const searchQuery = ref('');
const selectedCategory = ref('all');
const selectedSort = ref('popular');
const priceFilter = ref<'all' | 'free' | 'paid'>('all');

// Pagination
const currentPage = ref(1);
const pageSize = ref(8);

// Detail modal
const showDetailModal = ref(false);
const selectedModel = ref<ThinkingModel | null>(null);

// Categories
const categories = [
  { id: 'all', name: '全部' },
  { id: 'business', name: '商业管理' },
  { id: 'strategy', name: '战略规划' },
  { id: 'innovation', name: '创新思维' },
  { id: 'analysis', name: '分析工具' },
  { id: 'decision', name: '决策方法' },
  { id: 'creative', name: '创意构思' },
];

// Sort options
const sortOptions = [
  { id: 'popular', name: '最受欢迎' },
  { id: 'newest', name: '最新发布' },
  { id: 'mostAdopted', name: '最多采纳' },
  { id: 'mostLiked', name: '最多点赞' },
];

// Mock data
const models = ref<ThinkingModel[]>([
  {
    id: '1',
    title: 'SWOT 分析模型',
    description: '系统分析优势、劣势、机会与威胁的经典战略框架',
    cover: 'https://images.unsplash.com/photo-1552664730-d307ca884978?w=400&h=200&fit=crop',
    author: { id: 'u1', name: '张明远', avatar: 'https://avatar.vercel.sh/zhangmy.svg?text=ZM' },
    isFree: true,
    category: 'strategy',
    tags: ['战略分析', '商业思维', '经典模型'],
    stats: { adoptions: 12580, practices: 45620, discussions: 892, forks: 2341, likes: 8934 },
    updatedAt: '2024-01-15',
  },
  {
    id: '2',
    title: 'PDCA 循环',
    description: '计划-执行-检查-行动的持续改进方法论',
    cover: 'https://images.unsplash.com/photo-1512758017271-d7b84c2113f1?w=400&h=200&fit=crop',
    author: { id: 'u2', name: '李思维', avatar: 'https://avatar.vercel.sh/lisiwei.svg?text=LS' },
    isFree: false,
    price: 29.9,
    category: 'business',
    tags: ['持续改进', '管理方法', '质量控制'],
    stats: { adoptions: 8920, practices: 32150, discussions: 567, forks: 1234, likes: 6234 },
    updatedAt: '2024-02-10',
  },
  {
    id: '3',
    title: '第一性原理思维',
    description: '像马斯克一样回归本质，打破常规的创新思考方式',
    cover: 'https://images.unsplash.com/photo-1507413245164-6160d8298b31?w=400&h=200&fit=crop',
    author: { id: 'u3', name: '王创新', avatar: 'https://avatar.vercel.sh/wangcx.svg?text=WC' },
    isFree: true,
    category: 'innovation',
    tags: ['创新思维', '底层逻辑', '马斯克'],
    stats: { adoptions: 15230, practices: 48200, discussions: 1234, forks: 3456, likes: 10234 },
    updatedAt: '2024-01-20',
  },
  {
    id: '4',
    title: '5W1H 分析模型',
    description: '六个基本问题全面剖析，确保思考无遗漏',
    cover: 'https://images.unsplash.com/photo-1454165804606-c3d57bc86b40?w=400&h=200&fit=crop',
    author: { id: 'u4', name: '赵分析', avatar: 'https://avatar.vercel.sh/zhaofx.svg?text=ZF' },
    isFree: true,
    category: 'analysis',
    tags: ['问题分析', '新闻写作', '全面思考'],
    stats: { adoptions: 6780, practices: 23450, discussions: 345, forks: 890, likes: 4567 },
    updatedAt: '2024-02-05',
  },
  {
    id: '5',
    title: '决策矩阵',
    description: '多维度量化评估，让复杂决策变得简单清晰',
    cover: 'https://images.unsplash.com/photo-1551288049-bebda4e38f71?w=400&h=200&fit=crop',
    author: { id: 'u5', name: '陈决策', avatar: 'https://avatar.vercel.sh/chenjc.svg?text=CJ' },
    isFree: false,
    price: 19.9,
    category: 'decision',
    tags: ['决策分析', '多维度评估', '优先级'],
    stats: { adoptions: 4560, practices: 18230, discussions: 278, forks: 567, likes: 3456 },
    updatedAt: '2024-02-12',
  },
  {
    id: '6',
    title: '六顶思考帽',
    description: '德博诺的经典创意工具，全方位激发团队创造力',
    cover: 'https://images.unsplash.com/photo-1517048676732-d65bc937f952?w=400&h=200&fit=crop',
    author: { id: 'u6', name: '刘思维', avatar: 'https://avatar.vercel.sh/liusw.svg?text=LS' },
    isFree: false,
    price: 39.9,
    category: 'creative',
    tags: ['创意思考', '团队讨论', '全面视角'],
    stats: { adoptions: 7230, practices: 28900, discussions: 456, forks: 1234, likes: 5678 },
    updatedAt: '2024-01-28',
  },
  {
    id: '7',
    title: 'OKR 目标管理',
    description: '谷歌都在用的目标与关键成果管理法',
    cover: 'https://images.unsplash.com/photo-1460925895917-afdab827c52f?w=400&h=200&fit=crop',
    author: { id: 'u7', name: '孙目标', avatar: 'https://avatar.vercel.sh/sunmb.svg?text=SM' },
    isFree: true,
    category: 'business',
    tags: ['目标管理', '绩效考核', '谷歌方法'],
    stats: { adoptions: 9870, practices: 35670, discussions: 678, forks: 1890, likes: 7234 },
    updatedAt: '2024-02-08',
  },
  {
    id: '8',
    title: 'ABCDE 优先级法',
    description: '艾森豪威尔矩阵升级版，高效管理时间和任务',
    cover: 'https://images.unsplash.com/photo-1484480974693-6ca0a78fb36b?w=400&h=200&fit=crop',
    author: { id: 'u8', name: '周优先', avatar: 'https://avatar.vercel.sh/zhouyx.svg?text=ZY' },
    isFree: true,
    category: 'decision',
    tags: ['时间管理', '优先级', '效率提升'],
    stats: { adoptions: 8340, practices: 29870, discussions: 423, forks: 1023, likes: 5123 },
    updatedAt: '2024-01-30',
  },
]);

// Filtered models
const filteredModels = computed(() => {
  let result = models.value;

  if (selectedCategory.value !== 'all') {
    result = result.filter((m) => m.category === selectedCategory.value);
  }

  if (priceFilter.value === 'free') {
    result = result.filter((m) => m.isFree);
  } else if (priceFilter.value === 'paid') {
    result = result.filter((m) => !m.isFree);
  }

  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase();
    result = result.filter(
      (m) =>
        m.title.toLowerCase().includes(query) ||
        m.description.toLowerCase().includes(query) ||
        m.tags.some((t) => t.toLowerCase().includes(query))
    );
  }

  switch (selectedSort.value) {
    case 'popular':
      result = [...result].sort((a, b) => b.stats.adoptions - a.stats.adoptions);
      break;
    case 'newest':
      result = [...result].sort((a, b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime());
      break;
    case 'mostAdopted':
      result = [...result].sort((a, b) => b.stats.adoptions - a.stats.adoptions);
      break;
    case 'mostLiked':
      result = [...result].sort((a, b) => b.stats.likes - a.stats.likes);
      break;
  }

  return result;
});

// Paginated models
const paginatedModels = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value;
  return filteredModels.value.slice(start, start + pageSize.value);
});

const totalPages = computed(() => Math.ceil(filteredModels.value.length / pageSize.value));

// Reset page when filters change
watch([searchQuery, selectedCategory, selectedSort, priceFilter], () => {
  currentPage.value = 1;
});

const router = useRouter();

// Actions
function goToDetail(model: ThinkingModel) {
  router.push(`/models/detail/${model.id}`);
}

function openDetailModal(model: ThinkingModel) {
  selectedModel.value = model;
  showDetailModal.value = true;
}

function formatNumber(num: number): string {
  if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M';
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K';
  return num.toString();
}

function handleLoad(model: ThinkingModel) {
  console.log('Load model:', model.id);
}

function handlePurchase(model: ThinkingModel) {
  console.log('Purchase model:', model.id);
}

function handleFork(model: ThinkingModel) {
  console.log('Fork model:', model.id);
}
</script>

<template>
  <Page
    description="发现、学习、应用各种强大的思维模型，提升你的思考深度和决策质量"
    title="思维模型市场"
  >
    <ElCard class="mb-4" shadow="never">
      <template #header>
        <div class="flex flex-wrap items-center justify-between gap-4">
          <!-- Category -->
          <div class="flex flex-wrap gap-2">
            <ElButton
              v-for="cat in categories"
              :key="cat.id"
              :type="selectedCategory === cat.id ? 'primary' : 'default'"
              @click="selectedCategory = cat.id"
            >
              {{ cat.name }}
            </ElButton>
          </div>

          <!-- Controls -->
          <div class="flex items-center gap-3">
            <ElInput
              v-model="searchQuery"
              placeholder="搜索模型..."
              clearable
              style="width: 200px"
            />
            <ElSelect v-model="selectedSort" style="width: 120px">
              <ElOption
                v-for="opt in sortOptions"
                :key="opt.id"
                :label="opt.name"
                :value="opt.id"
              />
            </ElSelect>
            <div class="flex rounded overflow-hidden border border-gray-300">
              <ElButton
                :type="priceFilter === 'all' ? 'primary' : 'default'"
                @click="priceFilter = 'all'"
              >
                全部
              </ElButton>
              <ElButton
                :type="priceFilter === 'free' ? 'primary' : 'default'"
                @click="priceFilter = 'free'"
              >
                免费
              </ElButton>
              <ElButton
                :type="priceFilter === 'paid' ? 'primary' : 'default'"
                @click="priceFilter = 'paid'"
              >
                付费
              </ElButton>
            </div>
          </div>
        </div>
      </template>

      <!-- Results Count -->
      <div class="mb-4 text-sm text-gray-500">
        找到 {{ filteredModels.length }} 个思维模型
      </div>

      <!-- Grid -->
      <div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
        <div
          v-for="model in paginatedModels"
          :key="model.id"
          class="cursor-pointer overflow-hidden rounded-lg bg-white shadow-sm transition-shadow hover:shadow-md border border-gray-100"
          @click="goToDetail(model)"
        >
          <div class="relative h-40 overflow-hidden">
            <img :src="model.cover" class="h-full w-full object-cover" />
            <span
              :class="[
                'absolute left-2 top-2 rounded-full px-2 py-0.5 text-xs font-bold',
                model.isFree ? 'bg-green-500 text-white' : 'bg-white text-purple-600'
              ]"
            >
              {{ model.isFree ? '免费' : '¥' + model.price }}
            </span>
          </div>
          <div class="p-4">
            <h3 class="font-semibold text-gray-900">{{ model.title }}</h3>
            <p class="mt-1 line-clamp-2 h-10 text-sm text-gray-500">{{ model.description }}</p>
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
            <!-- Author & Stats -->
            <div class="mt-3 flex items-center gap-2">
              <img :src="model.author.avatar" class="h-5 w-5 rounded-full" />
              <span class="text-xs text-gray-600">{{ model.author.name }}</span>
            </div>

            <!-- Stats Row -->
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
              <span class="flex items-center gap-1" title="讨论">
                <svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/>
                </svg>
                {{ formatNumber(model.stats.discussions) }}
              </span>
              <span class="flex items-center gap-1" title="引用">
                <svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7v8a2 2 0 002 2h6M8 7V5a2 2 0 012-2h4.586a1 1 0 01.707.293l4.414 4.414a1 1 0 01.293.707V15a2 2 0 01-2 2h-2M8 7H6a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2v-2"/>
                </svg>
                {{ formatNumber(model.stats.forks) }}
              </span>
              <span class="flex items-center gap-1" :style="{ color: 'var(--el-color-primary)' }" title="点赞">
                <svg class="h-3.5 w-3.5" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
                </svg>
                {{ formatNumber(model.stats.likes) }}
              </span>
            </div>

            <!-- Action Buttons -->
            <div class="mt-3 flex gap-2">
              <ElButton
                v-if="model.isFree"
                type="primary"
                class="flex-1"
                @click.stop="handleLoad(model)"
              >
                加载到模型
              </ElButton>
              <ElButton
                v-else
                type="success"
                class="flex-1"
                @click.stop="handlePurchase(model)"
              >
                购买 ¥{{ model.price }}
              </ElButton>
              <ElButton
                @click.stop="handleFork(model)"
                title="引用创建副本"
              >
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7v8a2 2 0 002 2h6M8 7V5a2 2 0 012-2h4.586a1 1 0 01.707.293l4.414 4.414a1 1 0 01.293.707V15a2 2 0 01-2 2h-2M8 7H6a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2v-2"/>
                </svg>
              </ElButton>
            </div>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="mt-8 flex justify-center gap-2">
        <ElButton
          v-for="page in totalPages"
          :key="page"
          :type="currentPage === page ? 'primary' : 'default'"
          @click="currentPage = page"
        >
          {{ page }}
        </ElButton>
      </div>
    </ElCard>

    <!-- Modal -->
    <div
      v-if="showDetailModal && selectedModel"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4"
      @click="showDetailModal = false"
    >
      <div class="w-full max-w-2xl max-h-[90vh] overflow-y-auto rounded-xl bg-white shadow-2xl" @click.stop>
        <!-- Cover Image -->
        <div class="relative h-56 w-full overflow-hidden">
          <img :src="selectedModel.cover" class="h-full w-full object-cover" />
          <div class="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent"></div>
          <ElButton
            circle
            class="absolute right-4 top-4"
            @click="showDetailModal = false"
          >
            <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </ElButton>
        </div>

        <div class="p-6">
          <!-- Header -->
          <div class="mb-4">
            <div class="flex items-start justify-between gap-4">
              <h2 class="text-2xl font-bold text-gray-900">{{ selectedModel.title }}</h2>
              <ElTag
                :type="selectedModel.isFree ? 'success' : 'primary'"
                size="large"
                effect="dark"
              >
                {{ selectedModel.isFree ? '免费' : '¥' + selectedModel.price }}
              </ElTag>
            </div>
            <p class="mt-2 text-gray-600">{{ selectedModel.description }}</p>
          </div>

          <!-- Tags -->
          <div class="mb-4 flex flex-wrap gap-2">
            <ElTag
              v-for="tag in selectedModel.tags"
              :key="tag"
              type="primary"
              effect="light"
            >
              {{ tag }}
            </ElTag>
          </div>

          <!-- Author & Date -->
          <div class="mb-6 flex items-center justify-between border-b border-gray-100 pb-4">
            <div class="flex items-center gap-3">
              <img :src="selectedModel.author.avatar" class="h-10 w-10 rounded-full" />
              <div>
                <div class="font-medium text-gray-900">{{ selectedModel.author.name }}</div>
                <div class="text-sm text-gray-500">作者</div>
              </div>
            </div>
            <div class="text-sm text-gray-500">
              更新于 {{ selectedModel.updatedAt }}
            </div>
          </div>

          <!-- Stats Grid -->
          <div class="mb-6 grid grid-cols-5 gap-4 rounded-lg bg-gray-50 p-4">
            <div class="text-center">
              <div class="flex items-center justify-center gap-1 text-gray-500">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
                </svg>
                <span class="text-xs">采纳</span>
              </div>
              <div class="mt-1 text-lg font-bold text-gray-900">{{ formatNumber(selectedModel.stats.adoptions) }}</div>
            </div>
            <div class="text-center">
              <div class="flex items-center justify-center gap-1 text-gray-500">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"/>
                </svg>
                <span class="text-xs">练习</span>
              </div>
              <div class="mt-1 text-lg font-bold text-gray-900">{{ formatNumber(selectedModel.stats.practices) }}</div>
            </div>
            <div class="text-center">
              <div class="flex items-center justify-center gap-1 text-gray-500">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8"/>
                </svg>
                <span class="text-xs">讨论</span>
              </div>
              <div class="mt-1 text-lg font-bold text-gray-900">{{ formatNumber(selectedModel.stats.discussions) }}</div>
            </div>
            <div class="text-center">
              <div class="flex items-center justify-center gap-1 text-gray-500">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7v8a2 2 0 002 2h6M8 7V5a2 2 0 012-2h4.586"/>
                </svg>
                <span class="text-xs">引用</span>
              </div>
              <div class="mt-1 text-lg font-bold text-gray-900">{{ formatNumber(selectedModel.stats.forks) }}</div>
            </div>
            <div class="text-center">
              <div class="flex items-center justify-center gap-1" :style="{ color: 'var(--el-color-primary)' }">
                <svg class="h-4 w-4" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3"/>
                </svg>
                <span class="text-xs">点赞</span>
              </div>
              <div class="mt-1 text-lg font-bold" :style="{ color: 'var(--el-color-primary)' }">{{ formatNumber(selectedModel.stats.likes) }}</div>
            </div>
          </div>

          <!-- Action Buttons -->
          <div class="flex gap-3">
            <ElButton
              v-if="selectedModel.isFree"
              type="primary"
              size="large"
              class="flex-1"
              @click="handleLoad(selectedModel); showDetailModal = false"
            >
              加载到我的模型
            </ElButton>
            <ElButton
              v-else
              type="success"
              size="large"
              class="flex-1"
              @click="handlePurchase(selectedModel); showDetailModal = false"
            >
              购买 ¥{{ selectedModel.price }}
            </ElButton>
            <ElButton
              size="large"
              @click="handleFork(selectedModel); showDetailModal = false"
            >
              <svg class="h-5 w-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7v8a2 2 0 002 2h6M8 7V5a2 2 0 012-2h4.586a1 1 0 01.707.293l4.414 4.414a1 1 0 01.293.707V15a2 2 0 01-2 2h-2M8 7H6a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2v-2"/>
              </svg>
              引用创建副本
            </ElButton>
          </div>
        </div>
      </div>
    </div>
  </Page>
</template>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
