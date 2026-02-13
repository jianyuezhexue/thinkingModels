<script lang="ts" setup>
import { onMounted, ref, watch } from 'vue';
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
  ElSkeletonItem,
  ElEmpty,
  ElTag,
  ElAvatar,
  ElTooltip,
  ElPagination,
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
const categories = [
  { id: 'all', name: '全部模型', icon: '🎯' },
  { id: 'business', name: '商业管理', icon: '💼' },
  { id: 'strategy', name: '战略规划', icon: '🎯' },
  { id: 'innovation', name: '创新思维', icon: '💡' },
  { id: 'analysis', name: '分析工具', icon: '📊' },
  { id: 'decision', name: '决策方法', icon: '⚖️' },
  { id: 'creative', name: '创意构思', icon: '🎨' },
];

// 排序选项
const sortOptions = [
  { id: 'popular', name: '最受欢迎' },
  { id: 'newest', name: '最新发布' },
  { id: 'mostAdopted', name: '最多采纳' },
  { id: 'mostLiked', name: '最多点赞' },
];

// 热门标签
const hotTags = ref([
  'SWOT分析', '金字塔原理', '六顶思考帽', '设计思维', 'PDCA循环',
  'OKR', '5Why分析', '波特五力', 'MECE原则', '费米估算'
]);

// 推荐模型
const recommendedModels = ref<ModelApi.ThinkingModel[]>([]);

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
    
    // 如果还没有推荐模型，取前3个作为推荐
    if (recommendedModels.value.length === 0 && res.list.length > 0) {
      recommendedModels.value = res.list.slice(0, 3);
    }
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
  if (num >= 10000) return (num / 10000).toFixed(1) + '万';
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K';
  return num.toString();
}

// 加载模型到我的模型库
async function handleLoad(model: ModelApi.ThinkingModel, event: Event) {
  event.stopPropagation();
  try {
    await adoptModelApi(model.id);
    ElMessage.success('已成功加载到您的模型库');
  } catch (error) {
    console.error('加载模型失败:', error);
  }
}

// 购买模型
async function handlePurchase(model: ModelApi.ThinkingModel, event: Event) {
  event.stopPropagation();
  try {
    await purchaseModelApi(model.id);
    ElMessage.success('购买成功！已添加到您的模型库');
  } catch (error) {
    console.error('购买失败:', error);
  }
}

// 引用模型
async function handleFork(model: ModelApi.ThinkingModel, event: Event) {
  event.stopPropagation();
  try {
    await forkModelApi(model.id);
    ElMessage.success('已创建副本到您的模型库');
  } catch (error) {
    console.error('引用失败:', error);
  }
}

// 获取分类图标
function getCategoryIcon(categoryId: string): string {
  const cat = categories.find(c => c.id === categoryId);
  return cat?.icon || '📁';
}

// 获取分类名称
function getCategoryName(categoryId: string): string {
  const cat = categories.find(c => c.id === categoryId);
  return cat?.name || categoryId;
}
</script>

<template>
  <Page
    description="发现、学习、应用各种强大的思维模型，提升你的思考深度和决策质量"
    title="思维模型市场"
    content-class="p-6 bg-gray-50"
  >
    <!-- 主内容区域 -->
    <div class="flex gap-6">
      <!-- 左侧主内容 -->
      <div class="flex-1 min-w-0">
        <!-- 筛选卡片 -->
        <ElCard shadow="hover" class="mb-6 !rounded-xl">
          <div class="flex flex-col gap-4">
            <!-- 分类标签 -->
            <div class="flex flex-wrap items-center gap-2">
              <button
                v-for="cat in categories"
                :key="cat.id"
                class="px-4 py-2 rounded-full text-sm font-medium transition-all"
                :class="[
                  selectedCategory === cat.id
                    ? 'bg-purple-100 text-purple-700 shadow-md border border-purple-200 font-semibold'
                    : 'bg-gray-100 text-gray-700 hover:bg-purple-50 hover:text-purple-600 border border-gray-200'
                ]"
                @click="selectedCategory = cat.id"
              >
                {{ cat.icon }} {{ cat.name }}
              </button>
            </div>

            <!-- 搜索、排序和价格筛选 -->
            <div class="flex flex-wrap items-center gap-4">
              <ElInput
                v-model="searchQuery"
                placeholder="搜索思维模型..."
                clearable
                class="flex-1 !min-w-[200px]"
                @keyup.enter="fetchModelList"
              >
                <template #prefix>
                  <svg class="h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
                  </svg>
                </template>
              </ElInput>
              <ElSelect v-model="selectedSort" class="!w-32">
                <ElOption
                  v-for="opt in sortOptions"
                  :key="opt.id"
                  :label="opt.name"
                  :value="opt.id"
                />
              </ElSelect>
              <!-- 价格筛选 -->
              <div class="flex rounded-full overflow-hidden border border-gray-200 bg-gray-100 p-1">
                <button
                  v-for="opt in [{ id: 'all', name: '全部' }, { id: 'free', name: '免费' }, { id: 'paid', name: '付费' }]"
                  :key="opt.id"
                  class="px-4 py-1.5 rounded-full text-sm font-medium transition-all"
                  :class="[
                    priceFilter === opt.id
                      ? 'bg-white text-purple-600 shadow-sm'
                      : 'text-gray-600 hover:text-purple-600'
                  ]"
                  @click="priceFilter = opt.id as 'all' | 'free' | 'paid'"
                >
                  {{ opt.name }}
                </button>
              </div>
            </div>
          </div>
        </ElCard>

        <!-- 结果统计 -->
        <div class="mb-4 flex items-center justify-between text-sm text-gray-500">
          <span>找到 <span class="font-semibold text-purple-600">{{ total }}</span> 个思维模型</span>
          <span v-if="loading" class="flex items-center gap-2">
            <svg class="animate-spin h-4 w-4 text-purple-600" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            加载中...
          </span>
        </div>

        <!-- 加载状态 -->
        <div v-if="loading" class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          <ElCard v-for="i in 6" :key="i" shadow="never" class="!rounded-xl">
            <ElSkeleton animated>
              <template #template>
                <div class="space-y-3">
                  <ElSkeletonItem variant="image" style="width: 100%; height: 140px; border-radius: 8px" />
                  <ElSkeletonItem variant="h3" style="width: 70%" />
                  <ElSkeletonItem variant="text" />
                  <ElSkeletonItem variant="text" style="width: 80%" />
                </div>
              </template>
            </ElSkeleton>
          </ElCard>
        </div>

        <!-- 空状态 -->
        <ElCard v-else-if="models.length === 0" shadow="hover" class="!rounded-xl">
          <ElEmpty description="未找到相关模型，尝试调整筛选条件">
            <template #image>
              <div class="text-6xl">🔍</div>
            </template>
            <ElButton type="primary" class="!bg-purple-600 !border-purple-600 !rounded-full mt-4" @click="selectedCategory = 'all'; priceFilter = 'all'; searchQuery = ''">
              重置筛选
            </ElButton>
          </ElEmpty>
        </ElCard>

        <!-- 模型网格 -->
        <div v-else class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          <div
            v-for="model in models"
            :key="model.id"
            class="group bg-white rounded-xl border border-gray-100 overflow-hidden cursor-pointer transition-all hover:shadow-xl hover:border-purple-200 hover:-translate-y-1"
            @click="goToDetail(model)"
          >
            <!-- 封面图 -->
            <div class="relative h-36 overflow-hidden bg-gradient-to-br from-purple-50 to-indigo-100">
              <img
                :src="model.cover || '/images/default-model-cover.svg'"
                class="h-full w-full object-cover transition-transform group-hover:scale-110"
                @error="(e) => { const img = e.target as HTMLImageElement; if (img) img.src = '/images/default-model-cover.svg'; }"
              />
              <!-- 价格标签 -->
              <span
                :class="[
                  'absolute left-3 top-3 rounded-full px-3 py-1 text-xs font-bold shadow-md',
                  model.isFree ? 'bg-green-500 text-white' : 'bg-white text-purple-600'
                ]"
              >
                {{ model.isFree ? '🎁 免费' : '💰 ¥' + model.price }}
              </span>
              <!-- 分类标签 -->
              <span class="absolute right-3 top-3 rounded-full px-2 py-1 text-xs bg-white/90 text-gray-700 shadow-sm">
                {{ getCategoryIcon(model.category) }} {{ getCategoryName(model.category) }}
              </span>
            </div>

            <!-- 内容区 -->
            <div class="p-4">
              <!-- 标题 -->
              <h3 class="font-semibold text-gray-900 group-hover:text-purple-600 transition-colors line-clamp-1 text-lg">
                {{ model.title }}
              </h3>
              <!-- 描述 -->
              <p class="mt-2 text-sm text-gray-500 line-clamp-2 h-10">
                {{ model.description }}
              </p>

              <!-- 标签 -->
              <div class="mt-3 flex flex-wrap gap-1">
                <ElTag
                  v-for="tag in model.tags.slice(0, 3)"
                  :key="tag"
                  size="small"
                  effect="plain"
                  class="!bg-purple-50 !text-purple-600 !border-purple-200"
                >
                  {{ tag }}
                </ElTag>
              </div>

              <!-- 作者 -->
              <div class="mt-3 flex items-center gap-2">
                <ElAvatar :src="model.author.avatar" :size="24" />
                <span class="text-xs text-gray-600">{{ model.author.name }}</span>
              </div>

              <!-- 统计数据 -->
              <div class="mt-3 flex items-center gap-4 border-t border-gray-100 pt-3 text-xs text-gray-500">
                <ElTooltip content="采纳数">
                  <span class="flex items-center gap-1">
                    <svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
                    </svg>
                    {{ formatNumber(model.stats.adoptions) }}
                  </span>
                </ElTooltip>
                <ElTooltip content="练习次数">
                  <span class="flex items-center gap-1">
                    <svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"/>
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                    </svg>
                    {{ formatNumber(model.stats.practices) }}
                  </span>
                </ElTooltip>
                <ElTooltip content="点赞数">
                  <span class="flex items-center gap-1 text-red-400">
                    <svg class="h-3.5 w-3.5" fill="currentColor" viewBox="0 0 24 24">
                      <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
                    </svg>
                    {{ formatNumber(model.stats.likes) }}
                  </span>
                </ElTooltip>
              </div>

              <!-- 操作按钮 -->
              <div class="mt-3 flex gap-2">
                <ElButton
                  v-if="model.isFree"
                  type="primary"
                  class="flex-1 !bg-purple-600 !border-purple-600 hover:!bg-purple-700 !rounded-full"
                  size="small"
                  @click="handleLoad(model, $event)"
                >
                  立即使用
                </ElButton>
                <ElButton
                  v-else
                  type="success"
                  class="flex-1 !rounded-full"
                  size="small"
                  @click="handlePurchase(model, $event)"
                >
                  购买 ¥{{ model.price }}
                </ElButton>
                <ElButton
                  size="small"
                  class="!rounded-full"
                  @click="handleFork(model, $event)"
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
        <div v-if="total > pageSize" class="mt-8 flex justify-center">
          <ElPagination
            v-model:current-page="currentPage"
            :page-size="pageSize"
            :total="total"
            layout="prev, pager, next"
            background
          />
        </div>
      </div>

      <!-- 右侧边栏 -->
      <div class="w-72 flex-shrink-0 space-y-6 hidden lg:block">
        <!-- 探索引导 -->
        <ElCard shadow="hover" class="!rounded-xl !bg-gradient-to-br from-purple-50 to-indigo-50 !border-purple-100">
          <div class="text-center py-4">
            <div class="w-16 h-16 mx-auto mb-4 bg-purple-100 rounded-full flex items-center justify-center">
              <span class="text-3xl">🧠</span>
            </div>
            <h3 class="text-lg font-semibold text-gray-800 mb-2">探索思维模型</h3>
            <p class="text-sm text-gray-500 mb-4">掌握思维工具，升级认知系统</p>
            <ElButton type="primary" class="w-full !bg-purple-600 !border-purple-600 hover:!bg-purple-700 !rounded-full">
              浏览全部模型
            </ElButton>
          </div>
        </ElCard>

        <!-- 热门标签 -->
        <ElCard shadow="hover" class="!rounded-xl">
          <template #header>
            <div class="flex items-center gap-2">
              <span class="text-lg">🏷️</span>
              <span class="font-semibold text-gray-700">热门标签</span>
            </div>
          </template>
          <div class="flex flex-wrap gap-2">
            <button
              v-for="tag in hotTags"
              :key="tag"
              class="px-3 py-1.5 rounded-full text-sm bg-gray-100 text-gray-600 hover:bg-purple-100 hover:text-purple-600 transition-colors"
              @click="searchQuery = tag"
            >
              {{ tag }}
            </button>
          </div>
        </ElCard>

        <!-- 推荐模型 -->
        <ElCard shadow="hover" class="!rounded-xl">
          <template #header>
            <div class="flex items-center gap-2">
              <span class="text-lg">⭐</span>
              <span class="font-semibold text-gray-700">推荐模型</span>
            </div>
          </template>
          <div class="space-y-3">
            <div
              v-for="model in recommendedModels"
              :key="model.id"
              class="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-50 cursor-pointer transition-colors"
              @click="goToDetail(model)"
            >
              <div class="w-12 h-12 rounded-lg overflow-hidden bg-purple-100 flex-shrink-0">
                <img
                  :src="model.cover"
                  class="w-full h-full object-cover"
                  @error="(e) => { const img = e.target as HTMLImageElement; if (img) img.style.display = 'none'; }"
                />
              </div>
              <div class="flex-1 min-w-0">
                <div class="text-sm font-medium text-gray-800 line-clamp-1">{{ model.title }}</div>
                <div class="text-xs text-gray-400 mt-0.5">{{ formatNumber(model.stats.adoptions) }} 人采纳</div>
              </div>
            </div>
          </div>
        </ElCard>

        <!-- 学习指南 -->
        <ElCard shadow="hover" class="!rounded-xl">
          <template #header>
            <div class="flex items-center gap-2">
              <span class="text-lg">📚</span>
              <span class="font-semibold text-gray-700">学习指南</span>
            </div>
          </template>
          <div class="space-y-4">
            <div class="flex items-start gap-3">
              <div class="w-8 h-8 rounded-full bg-purple-100 text-purple-600 flex items-center justify-center font-bold text-sm flex-shrink-0">1</div>
              <div>
                <div class="font-medium text-gray-700 text-sm">选择模型</div>
                <div class="text-xs text-gray-500">根据场景选择合适的思维模型</div>
              </div>
            </div>
            <div class="flex items-start gap-3">
              <div class="w-8 h-8 rounded-full bg-purple-100 text-purple-600 flex items-center justify-center font-bold text-sm flex-shrink-0">2</div>
              <div>
                <div class="font-medium text-gray-700 text-sm">理解原理</div>
                <div class="text-xs text-gray-500">学习模型背后的逻辑和适用条件</div>
              </div>
            </div>
            <div class="flex items-start gap-3">
              <div class="w-8 h-8 rounded-full bg-purple-100 text-purple-600 flex items-center justify-center font-bold text-sm flex-shrink-0">3</div>
              <div>
                <div class="font-medium text-gray-700 text-sm">反复练习</div>
                <div class="text-xs text-gray-500">通过实践将模型内化为思维习惯</div>
              </div>
            </div>
          </div>
        </ElCard>

        <!-- 小贴士 -->
        <ElCard shadow="hover" class="!rounded-xl !bg-gradient-to-br from-amber-50 to-orange-50 !border-amber-100">
          <template #header>
            <div class="flex items-center gap-2">
              <span class="text-lg">💡</span>
              <span class="font-semibold text-amber-700">小贴士</span>
            </div>
          </template>
          <ul class="text-sm text-amber-800 space-y-2">
            <li class="flex items-start gap-2">
              <span class="text-amber-500">•</span>
              先精通3-5个核心模型
            </li>
            <li class="flex items-start gap-2">
              <span class="text-amber-500">•</span>
              结合实际问题练习应用
            </li>
            <li class="flex items-start gap-2">
              <span class="text-amber-500">•</span>
              定期复习巩固认知
            </li>
          </ul>
        </ElCard>
      </div>
    </div>
  </Page>
</template>

<style scoped>
.line-clamp-1 {
  display: -webkit-box;
  -webkit-line-clamp: 1;
  line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
