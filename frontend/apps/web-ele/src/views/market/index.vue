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
  ElMessage,
  ElSkeleton,
  ElEmpty,
} from 'element-plus';

import { getModelListApi, adoptModelApi, purchaseModelApi, forkModelApi, type ModelApi } from '#/api';

// 加载状态
const loading = ref(false);

// 模型列表数据
const models = ref<ModelApi.ThinkingModel[]>([]);
const total = ref(0);

// 筛选状态
const searchQuery = ref('');
const selectedCategory = ref('all');
const selectedSort = ref<ModelApi.ModelListParams['sortBy']>('popular');
const priceFilter = ref<'all' | 'free' | 'paid'>('all');

// 分页
const currentPage = ref(1);
const pageSize = ref(12);

// 分类列表
const categories = ref<{ id: string; name: string }[]>([
  { id: 'all', name: '全部' },
  { id: 'business', name: '商业管理' },
  { id: 'strategy', name: '战略规划' },
  { id: 'innovation', name: '创新思维' },
  { id: 'analysis', name: '分析工具' },
  { id: 'decision', name: '决策方法' },
  { id: 'creative', name: '创意构思' },
]);

// 排序选项
const sortOptions = [
  { id: 'popular', name: '最受欢迎' },
  { id: 'newest', name: '最新发布' },
  { id: 'mostAdopted', name: '最多采纳' },
  { id: 'mostLiked', name: '最多点赞' },
];

const router = useRouter();

// 获取模型列表
async function fetchModelList() {
  loading.value = true;
  try {
    const params: ModelApi.ModelListParams = {
      page: currentPage.value,
      pageSize: pageSize.value,
      sortBy: selectedSort.value,
      keyword: searchQuery.value || undefined,
    };

    if (selectedCategory.value !== 'all') {
      params.category = selectedCategory.value;
    }

    if (priceFilter.value === 'free') {
      params.isFree = true;
    } else if (priceFilter.value === 'paid') {
      params.isFree = false;
    }

    const res = await getModelListApi(params);
    models.value = res.list;
    total.value = res.total;
  } catch (error) {
    console.error('获取模型列表失败:', error);
    ElMessage.error('获取模型列表失败');
  } finally {
    loading.value = false;
  }
}

// 监听筛选条件变化，重置页码并重新加载
watch([searchQuery, selectedCategory, selectedSort, priceFilter], () => {
  currentPage.value = 1;
  fetchModelList();
});

// 监听页码变化
watch(currentPage, () => {
  fetchModelList();
});

// 页面加载时获取数据
onMounted(() => {
  fetchModelList();
});

// 跳转到详情页
function goToDetail(model: ModelApi.ThinkingModel) {
  router.push(`/market/${model.id}`);
}

// 格式化数字
function formatNumber(num: number): string {
  if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M';
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K';
  return num.toString();
}

// 加载模型到我的模型库
async function handleLoad(model: ModelApi.ThinkingModel) {
  try {
    await adoptModelApi(model.id);
    ElMessage.success('已成功加载到您的模型库');
  } catch (error) {
    console.error('加载模型失败:', error);
  }
}

// 购买模型
async function handlePurchase(model: ModelApi.ThinkingModel) {
  try {
    await purchaseModelApi(model.id);
    ElMessage.success('购买成功！已添加到您的模型库');
  } catch (error) {
    console.error('购买失败:', error);
  }
}

// 引用模型
async function handleFork(model: ModelApi.ThinkingModel) {
  try {
    await forkModelApi(model.id);
    ElMessage.success('已创建副本到您的模型库');
  } catch (error) {
    console.error('引用失败:', error);
  }
}

// 计算总页数
const totalPages = computed(() => Math.ceil(total.value / pageSize.value));

// 分页显示的页码
const displayPages = computed(() => {
  const pages: number[] = [];
  const maxVisible = 7;
  const half = Math.floor(maxVisible / 2);

  let start = Math.max(1, currentPage.value - half);
  let end = Math.min(totalPages.value, start + maxVisible - 1);

  if (end - start + 1 < maxVisible) {
    start = Math.max(1, end - maxVisible + 1);
  }

  for (let i = start; i <= end; i++) {
    pages.push(i);
  }

  return pages;
});
</script>

<template>
  <Page
    description="发现、学习、应用各种强大的思维模型，提升你的思考深度和决策质量"
    title="思维模型市场"
  >
    <ElCard class="mb-4" shadow="never">
      <template #header>
        <div class="flex flex-wrap items-center justify-between gap-4">
          <!-- 分类筛选 -->
          <div class="flex flex-wrap gap-2">
            <ElButton
              v-for="cat in categories"
              :key="cat.id"
              :type="selectedCategory === cat.id ? 'primary' : 'default'"
              size="small"
              @click="selectedCategory = cat.id"
            >
              {{ cat.name }}
            </ElButton>
          </div>

          <!-- 搜索和控制 -->
          <div class="flex items-center gap-3">
            <ElInput
              v-model="searchQuery"
              placeholder="搜索思维模型..."
              clearable
              style="width: 220px"
              @keyup.enter="fetchModelList"
            >
              <template #prefix>
                <svg class="h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
                </svg>
              </template>
            </ElInput>
            <ElSelect v-model="selectedSort" style="width: 130px">
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
                size="small"
                @click="priceFilter = 'all'"
              >
                全部
              </ElButton>
              <ElButton
                :type="priceFilter === 'free' ? 'primary' : 'default'"
                size="small"
                @click="priceFilter = 'free'"
              >
                免费
              </ElButton>
              <ElButton
                :type="priceFilter === 'paid' ? 'primary' : 'default'"
                size="small"
                @click="priceFilter = 'paid'"
              >
                付费
              </ElButton>
            </div>
          </div>
        </div>
      </template>

      <!-- 结果统计 -->
      <div class="mb-4 flex items-center justify-between text-sm text-gray-500">
        <span>找到 {{ total }} 个思维模型</span>
        <span v-if="loading">加载中...</span>
      </div>

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
      <ElEmpty v-else-if="models.length === 0" description="未找到相关模型，尝试调整筛选条件" />

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
            <img :src="model.cover" class="h-full w-full object-cover transition-transform group-hover:scale-105" />
            <span
              :class="[
                'absolute left-2 top-2 rounded-full px-2 py-0.5 text-xs font-bold',
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

            <!-- 作者 -->
            <div class="mt-3 flex items-center gap-2">
              <img :src="model.author.avatar" class="h-5 w-5 rounded-full" />
              <span class="text-xs text-gray-600">{{ model.author.name }}</span>
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
              <span class="flex items-center gap-1" title="讨论">
                <svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/>
                </svg>
                {{ formatNumber(model.stats.discussions) }}
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
                v-if="model.isFree"
                type="primary"
                class="flex-1"
                size="small"
                @click.stop="handleLoad(model)"
              >
                立即使用
              </ElButton>
              <ElButton
                v-else
                type="success"
                class="flex-1"
                size="small"
                @click.stop="handlePurchase(model)"
              >
                购买 ¥{{ model.price }}
              </ElButton>
              <ElButton
                size="small"
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

      <!-- 分页 -->
      <div v-if="totalPages > 1" class="mt-8 flex justify-center gap-2">
        <ElButton
          :disabled="currentPage === 1"
          size="small"
          @click="currentPage--"
        >
          上一页
        </ElButton>
        <ElButton
          v-for="page in displayPages"
          :key="page"
          :type="currentPage === page ? 'primary' : 'default'"
          size="small"
          @click="currentPage = page"
        >
          {{ page }}
        </ElButton>
        <ElButton
          :disabled="currentPage === totalPages"
          size="small"
          @click="currentPage++"
        >
          下一页
        </ElButton>
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
</style>
