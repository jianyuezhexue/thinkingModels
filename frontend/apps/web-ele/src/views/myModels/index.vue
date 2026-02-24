<template>
  <Page description="管理你创建的思维模型，跟踪收入和用户反馈" title="我的模型" content-class="p-6 bg-gray-50">
    <!-- 顶部统计卡片 -->
    <div class="mb-6 grid grid-cols-2 md:grid-cols-4 gap-4">
      <div class="bg-white rounded-xl p-5 border border-gray-100 shadow-sm hover:shadow-md transition-shadow">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-500">总模型数</p>
            <p class="mt-1 text-2xl font-bold text-gray-900">{{ stats.totalModels }}</p>
            <p class="mt-1 text-xs text-purple-500">{{ stats.publishedModels }} 个已发布</p>
          </div>
          <div class="w-12 h-12 rounded-xl bg-purple-100 flex items-center justify-center">
            <span class="text-2xl">📦</span>
          </div>
        </div>
      </div>
      <div class="bg-white rounded-xl p-5 border border-gray-100 shadow-sm hover:shadow-md transition-shadow">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-500">审核中</p>
            <p class="mt-1 text-2xl font-bold text-amber-600">{{ stats.underReviewModels }}</p>
            <p class="mt-1 text-xs text-amber-500">等待审核</p>
          </div>
          <div class="w-12 h-12 rounded-xl bg-amber-100 flex items-center justify-center">
            <span class="text-2xl">⏳</span>
          </div>
        </div>
      </div>
      <div class="bg-white rounded-xl p-5 border border-gray-100 shadow-sm hover:shadow-md transition-shadow">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-500">被采纳次数</p>
            <p class="mt-1 text-2xl font-bold text-blue-600">{{ formatNumber(stats.totalAdoptions) }}</p>
            <p class="mt-1 text-xs text-blue-500">帮助了 {{ Math.floor(stats.totalAdoptions * 0.7) }}+ 人</p>
          </div>
          <div class="w-12 h-12 rounded-xl bg-blue-100 flex items-center justify-center">
            <span class="text-2xl">✅</span>
          </div>
        </div>
      </div>
      <div class="bg-white rounded-xl p-5 border border-gray-100 shadow-sm hover:shadow-md transition-shadow">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-500">获得点赞</p>
            <p class="mt-1 text-2xl font-bold text-red-500">{{ formatNumber(stats.totalLikes) }}</p>
            <p class="mt-1 text-xs text-red-400">好评率 98.5%</p>
          </div>
          <div class="w-12 h-12 rounded-xl bg-red-100 flex items-center justify-center">
            <span class="text-2xl">❤️</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 主内容区 -->
    <div class="flex gap-6">
      <!-- 左侧列表 -->
      <div class="flex-1 min-w-0 space-y-6">
        <!-- 筛选卡片 -->
        <ElCard shadow="hover" class="!rounded-xl">
          <div class="flex flex-col gap-4">
            <!-- 状态 Tab -->
            <div class="flex flex-wrap items-center gap-2">
              <button v-for="tab in statusTabs" :key="tab.id"
                class="px-4 py-2 rounded-full text-sm font-medium transition-all" :class="[
                  activeStatus === tab.id
                    ? 'bg-purple-100 text-purple-700 shadow-md border border-purple-200 font-semibold'
                    : 'bg-gray-100 text-gray-700 hover:bg-purple-50 hover:text-purple-600 border border-gray-200'
                ]" @click="activeStatus = tab.id as any">
                {{ tab.icon }} {{ tab.label }}
              </button>
              <div class="flex-1" />
              <ElButton type="primary" class="!bg-purple-600 !border-purple-600 hover:!bg-purple-700 !rounded-full"
                @click="goToCreate">
                <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                </svg>
                创建模型
              </ElButton>
            </div>

            <!-- 搜索和排序 -->
            <div class="flex items-center gap-3">
              <ElInput v-model="searchKeyword" placeholder="搜索模型..." clearable class="flex-1"
                @keyup.enter="fetchModels">
                <template #prefix>
                  <svg class="h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                  </svg>
                </template>
              </ElInput>
              <ElSelect v-model="sortBy" class="!w-32">
                <ElOption v-for="opt in sortOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
              </ElSelect>
            </div>
          </div>
        </ElCard>

        <!-- 结果统计 -->
        <div class="flex items-center justify-between text-sm text-gray-500">
          <span>共 <span class="font-semibold text-purple-600">{{ total }}</span> 个模型</span>
        </div>

        <!-- 加载状态 -->
        <div v-if="loading" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <ElCard v-for="i in 4" :key="i" shadow="never" class="!rounded-xl">
            <ElSkeleton animated>
              <template #template>
                <div class="space-y-3">
                  <ElSkeletonItem variant="image" style="width: 100%; height: 120px; border-radius: 8px" />
                  <ElSkeletonItem variant="h3" style="width: 70%" />
                  <ElSkeletonItem variant="text" />
                </div>
              </template>
            </ElSkeleton>
          </ElCard>
        </div>

        <!-- 空状态 -->
        <ElCard v-else-if="models.length === 0" shadow="hover" class="!rounded-xl">
          <ElEmpty description="还没有模型，创建你的第一个思维模型吧！">
            <template #image>
              <div class="text-6xl">🧠</div>
            </template>
            <ElButton type="primary" class="!bg-purple-600 !border-purple-600 !rounded-full mt-4" @click="goToCreate">
              创建模型
            </ElButton>
          </ElEmpty>
        </ElCard>

        <!-- 模型列表 -->
        <div v-else class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <div v-for="model in models" :key="model.id"
            class="group bg-white rounded-xl border border-gray-100 overflow-hidden cursor-pointer transition-all hover:shadow-xl hover:border-purple-200 hover:-translate-y-1"
            @click="goToDetail(model)">
            <!-- 封面 -->
            <div class="relative h-32 overflow-hidden bg-gradient-to-br from-purple-100 to-indigo-100">
              <img v-if="model.coverImage" :src="model.coverImage"
                class="h-full w-full object-cover transition-transform group-hover:scale-110"
                @error="(e) => { const img = e.target as HTMLImageElement; if (img) img.style.display = 'none'; }" />
              <div v-else class="h-full w-full flex items-center justify-center text-4xl">
                {{ model.icon || '📝' }}
              </div>
              <!-- 状态标签 -->
              <span :class="[
                'absolute left-3 top-3 rounded-full px-2.5 py-1 text-xs font-medium',
                getStatusStyle(model.status).bg,
                getStatusStyle(model.status).text
              ]">
                {{ getStatusStyle(model.status).label }}
              </span>
              <!-- 价格标签 -->
              <span :class="[
                'absolute right-3 top-3 rounded-full px-2.5 py-1 text-xs font-bold shadow-md',
                model.isFree ? 'bg-green-500 text-white' : 'bg-white text-purple-600'
              ]">
                {{ model.isFree ? '免费' : '¥' + model.price }}
              </span>
            </div>

            <!-- 内容 -->
            <div class="p-4">
              <h3 class="font-semibold text-gray-900 group-hover:text-purple-600 transition-colors line-clamp-1">
                {{ model.name }}
              </h3>
              <p class="mt-1 text-sm text-gray-500 line-clamp-2 h-10">
                {{ model.description }}
              </p>

              <!-- 标签 -->
              <div v-if="model.tags?.length" class="mt-2 flex flex-wrap gap-1">
                <ElTag v-for="tag in model.tags.slice(0, 4)" :key="tag" size="small" effect="plain"
                  class="!bg-purple-50 !text-purple-600 !border-purple-200 !rounded-full">
                  {{ tag }}
                </ElTag>
              </div>

              <!-- 统计 -->
              <div v-if="model.status === 1"
                class="mt-3 flex items-center gap-4 text-xs text-gray-500 border-t border-gray-100 pt-3">
                <span class="flex items-center gap-1">
                  <svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                  {{ formatNumber(model.stats.adoptCount) }}
                </span>
                <span class="flex items-center gap-1 text-red-400">
                  <svg class="h-3.5 w-3.5" fill="currentColor" viewBox="0 0 24 24">
                    <path
                      d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z" />
                  </svg>
                  {{ formatNumber(model.stats.likeCount) }}
                </span>
              </div>

              <!-- 操作按钮 -->
              <div class="mt-3 flex gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                <ElButton size="small" class="flex-1 !rounded-full" @click="handleEdit(model, $event)">
                  编辑
                </ElButton>
                <ElButton v-if="model.status === 0" type="primary" size="small"
                  class="flex-1 !bg-purple-600 !border-purple-600 !rounded-full" @click="handlePublish(model, $event)">
                  发布
                </ElButton>
                <ElButton type="danger" size="small" plain class="!rounded-full" @click="handleDelete(model, $event)">
                  <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </ElButton>
              </div>
            </div>
          </div>
        </div>

        <!-- 分页 -->
        <div v-if="total > 0" class="flex justify-center pt-4">
          <ElPagination v-model:current-page="currentPage" v-model:page-size="pageSize" :total="total"
            :page-sizes="[12, 24, 48]" layout="total, sizes, prev, pager, next, jumper" background
            @size-change="handleSizeChange" @current-change="handleCurrentChange" />
        </div>
      </div>

      <!-- 右侧边栏 -->
      <div class="w-80 flex-shrink-0 space-y-6 hidden lg:block">
        <!-- 创建引导 -->
        <ElCard shadow="hover" class="!rounded-xl !bg-gradient-to-br from-purple-50 to-indigo-50 !border-purple-100">
          <div class="text-center py-4">
            <div class="w-16 h-16 mx-auto mb-4 bg-purple-100 rounded-full flex items-center justify-center">
              <span class="text-3xl">💡</span>
            </div>
            <h3 class="text-lg font-semibold text-gray-800 mb-2">分享你的思维模型</h3>
            <p class="text-sm text-gray-500 mb-4">创建并分享你的思维模型，帮助他人提升思考能力</p>
            <ElButton type="primary" class="w-full !bg-purple-600 !border-purple-600 hover:!bg-purple-700 !rounded-full"
              @click="goToCreate">
              创建新模型
            </ElButton>
          </div>
        </ElCard>

        <!-- 热门模型 -->
        <ElCard shadow="hover" class="!rounded-xl">
          <template #header>
            <div class="flex items-center gap-2">
              <span class="text-lg">🔥</span>
              <span class="font-semibold text-gray-700">热门模型</span>
            </div>
          </template>
          <div class="space-y-3">
            <template v-if="models.filter(m => m.status === 1).length > 0">
              <div v-for="model in models.filter(m => m.status === 1).slice(0, 3)" :key="model.id"
                class="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-50 cursor-pointer transition-colors"
                @click="goToDetail(model)">
                <div
                  class="w-10 h-10 rounded-lg overflow-hidden bg-purple-100 flex items-center justify-center flex-shrink-0 text-xl">
                  {{ model.icon || '📝' }}
                </div>
                <div class="flex-1 min-w-0">
                  <div class="text-sm font-medium text-gray-800 line-clamp-1">{{ model.name }}</div>
                  <div class="text-xs text-gray-400">{{ formatNumber(model.stats.adoptCount) }} 采纳</div>
                </div>
              </div>
            </template>
            <div v-else class="text-center text-gray-400 py-4">
              暂无已发布的模型
            </div>
          </div>
        </ElCard>

        <!-- 创作指南 -->
        <ElCard shadow="hover" class="!rounded-xl">
          <template #header>
            <div class="flex items-center gap-2">
              <span class="text-lg">📚</span>
              <span class="font-semibold text-gray-700">创作指南</span>
            </div>
          </template>
          <div class="space-y-4">
            <div class="flex items-start gap-3">
              <div
                class="w-8 h-8 rounded-full bg-purple-100 text-purple-600 flex items-center justify-center font-bold text-sm flex-shrink-0">
                1</div>
              <div>
                <div class="font-medium text-gray-700 text-sm">明确模型用途</div>
                <div class="text-xs text-gray-500">确定模型解决什么问题</div>
              </div>
            </div>
            <div class="flex items-start gap-3">
              <div
                class="w-8 h-8 rounded-full bg-purple-100 text-purple-600 flex items-center justify-center font-bold text-sm flex-shrink-0">
                2</div>
              <div>
                <div class="font-medium text-gray-700 text-sm">清晰的使用步骤</div>
                <div class="text-xs text-gray-500">让用户容易上手</div>
              </div>
            </div>
            <div class="flex items-start gap-3">
              <div
                class="w-8 h-8 rounded-full bg-purple-100 text-purple-600 flex items-center justify-center font-bold text-sm flex-shrink-0">
                3</div>
              <div>
                <div class="font-medium text-gray-700 text-sm">丰富的案例</div>
                <div class="text-xs text-gray-500">通过实例帮助理解</div>
              </div>
            </div>
          </div>
        </ElCard>

        <!-- 小贴士 -->
        <ElCard shadow="hover" class="!rounded-xl !bg-gradient-to-br from-amber-50 to-orange-50 !border-amber-100">
          <template #header>
            <div class="flex items-center gap-2">
              <span class="text-lg">💡</span>
              <span class="font-semibold text-amber-700">创作小贴士</span>
            </div>
          </template>
          <ul class="text-sm text-amber-800 space-y-2">
            <li class="flex items-start gap-2">
              <span class="text-amber-500">•</span>
              优质封面图能提升50%点击率
            </li>
            <li class="flex items-start gap-2">
              <span class="text-amber-500">•</span>
              详细的使用步骤更受欢迎
            </li>
            <li class="flex items-start gap-2">
              <span class="text-amber-500">•</span>
              定期更新保持模型活力
            </li>
          </ul>
        </ElCard>
      </div>
    </div>
  </Page>
</template>

<script lang="ts" setup>
import { onMounted, ref, watch, onUnmounted } from 'vue';
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
  ElMessageBox,
} from 'element-plus';

import {
  getMyThinkingModelListApi,
  getModelStatusCountsApi,
  publishThinkingModelApi,
  deleteThinkingModelApi,
} from '#/api/thinking/model';
import { getAllCategoriesApi } from '#/api/master/category';

// 类型定义
interface MyModel {
  id: number;
  name: string;
  description: string;
  coverImage: string;
  icon: string;
  categoryId: number;
  categoryName?: string;
  tags: string[] | null;
  status: number; // 0=草稿, 1=已发布, 2=已下架, 3=审核中, 4=已驳回
  price: number;
  isFree: boolean;
  stats: {
    usageCount: number;
    adoptCount: number;
    likeCount: number;
    commentCount: number;
  };
  author: {
    id: string;
    name: string;
    avatar: string;
  };
  createdAt: string;
  updatedAt: string;
}

interface StatsSummary {
  totalModels: number;
  publishedModels: number;
  draftModels: number;
  underReviewModels: number;
  totalAdoptions: number;
  totalLikes: number;
}

// 路由
const router = useRouter();

// 请求取消控制器
let abortController: AbortController | null = null;

// 加载状态
const loading = ref(false);

// 模型列表
const models = ref<MyModel[]>([]);
const total = ref(0);

// 统计数据
const stats = ref<StatsSummary>({
  totalModels: 0,
  publishedModels: 0,
  draftModels: 0,
  underReviewModels: 0,
  totalAdoptions: 0,
  totalLikes: 0,
});

// 分页
const currentPage = ref(1);
const pageSize = ref(12);

// 筛选
const searchKeyword = ref('');
const activeStatus = ref<number | 'all'>('all');
const sortBy = ref<'newest' | 'popular' | 'adoptions'>('newest');

// 分类选项
const categoryOptions = ref<{ value: number; label: string }[]>([]);
const selectedCategory = ref<number | ''>('');

// 状态 Tab（后端状态值：0=草稿, 1=已发布, 2=已下架, 3=审核中, 4=已驳回）
const statusTabs = [
  { id: 'all' as const, status: undefined, label: '全部模型', icon: '📦' },
  { id: 1, status: 1, label: '已发布', icon: '✅' },
  { id: 0, status: 0, label: '草稿箱', icon: '📝' },
  { id: 3, status: 3, label: '审核中', icon: '⏳' },
];

// 排序选项
const sortOptions = [
  { value: 'newest', label: '最新创建' },
  { value: 'popular', label: '最受欢迎' },
  { value: 'adoptions', label: '采纳最多' },
];

// 加载分类列表
async function loadCategories() {
  try {
    const list = await getAllCategoriesApi();
    categoryOptions.value = list.map((item) => ({
      value: Number(item.id),
      label: item.name,
    }));
  } catch (error) {
    console.error('加载分类列表失败:', error);
  }
}

// 加载统计数据
async function loadStats() {
  try {
    const counts = await getModelStatusCountsApi();
    const allModels = models.value;
    stats.value = {
      totalModels: counts.pending + counts.approved + counts.rejected,
      publishedModels: counts.approved,
      draftModels: allModels.filter(m => m.status === 0).length,
      underReviewModels: counts.pending,
      totalAdoptions: allModels.reduce((sum, m) => sum + m.stats.adoptCount, 0),
      totalLikes: allModels.reduce((sum, m) => sum + m.stats.likeCount, 0),
    };
  } catch (error) {
    console.error('加载统计数据失败:', error);
  }
}

// 获取模型列表
async function fetchModels() {
  // 取消之前的请求
  abortController?.abort();
  abortController = new AbortController();
  const signal = abortController.signal;

  loading.value = true;
  try {
    const params: Record<string, any> = {
      page: currentPage.value,
      pageSize: pageSize.value,
    };

    // 状态筛选
    if (activeStatus.value !== 'all') {
      params.status = activeStatus.value;
    }

    // 关键字搜索
    if (searchKeyword.value) {
      params.name = searchKeyword.value;
    }

    // 分类筛选
    if (selectedCategory.value !== '') {
      params.categoryId = selectedCategory.value;
    }

    // 排序
    if (sortBy.value === 'popular') {
      params.sortBy = 'likeCount';
    } else if (sortBy.value === 'adoptions') {
      params.sortBy = 'adoptCount';
    }

    const res = await getMyThinkingModelListApi(params, { signal });
    models.value = res.list.map((item) => ({
      ...item,
      categoryName: categoryOptions.value.find((c) => c.value === item.categoryId)?.label || '',
    }));
    total.value = res.total;
  } catch (error: any) {
    // 如果是取消错误，静默处理
    if (error?.name === 'CanceledError' || error?.code === 'ERR_CANCELED' || signal.aborted) {
      return;
    }
    console.error('获取模型列表失败:', error);
    ElMessage.error('获取模型列表失败');
  } finally {
    loading.value = false;
  }
}


// 监听变化
watch([searchKeyword, activeStatus, sortBy, selectedCategory], () => {
  currentPage.value = 1;
  fetchModels();
});

watch([currentPage, pageSize], () => {
  fetchModels();
});

// 分页处理
function handleSizeChange(size: number) {
  pageSize.value = size;
  fetchModels();
}

function handleCurrentChange(page: number) {
  currentPage.value = page;
  fetchModels();
}

onMounted(async () => {
  await loadCategories();
  await fetchModels();
  loadStats();
});

onUnmounted(() => {
  abortController?.abort();
});

// 操作函数
function goToCreate() {
  router.push('/my-models/create');
}

function goToDetail(model: MyModel) {
  router.push(`/my-models/${model.id}`);
}

function handleEdit(model: MyModel, event: Event) {
  event.stopPropagation();
  router.push(`/my-models/create?id=${model.id}`);
}

async function handleDelete(model: MyModel, event: Event) {
  event.stopPropagation();
  try {
    await ElMessageBox.confirm(
      `确定要删除模型「${model.name}」吗？此操作不可恢复。`,
      '删除确认',
      { type: 'warning' }
    );
    await deleteThinkingModelApi({ ids: [model.id] });
    ElMessage.success('模型已删除');
    fetchModels();
    loadStats();
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('删除失败:', error);
      ElMessage.error('删除失败');
    }
  }
}

async function handlePublish(m: MyModel, event: Event) {
  event.stopPropagation();
  try {
    await ElMessageBox.confirm(
      `提交「${m.name}」审核后，模型将在审核通过后发布到市场。确定提交吗？`,
      '提交审核',
      { type: 'info' }
    );
    await publishThinkingModelApi({ id: m.id });
    ElMessage.success('模型已提交审核');
    fetchModels();
    loadStats();
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('提交审核失败:', error);
      ElMessage.error('提交审核失败');
    }
  }
}

// 工具函数 - 状态映射：0=草稿, 1=已发布, 2=已下架, 3=审核中, 4=已驳回
function getStatusStyle(status: number): { bg: string; text: string; label: string } {
  const styles: Record<number, { bg: string; text: string; label: string }> = {
    0: { bg: 'bg-gray-100', text: 'text-gray-600', label: '草稿' },
    1: { bg: 'bg-green-100', text: 'text-green-700', label: '已发布' },
    2: { bg: 'bg-yellow-100', text: 'text-yellow-700', label: '已下架' },
    3: { bg: 'bg-amber-100', text: 'text-amber-700', label: '审核中' },
    4: { bg: 'bg-red-100', text: 'text-red-700', label: '已驳回' },
  };
  return styles[status] || { bg: 'bg-gray-100', text: 'text-gray-600', label: '未知' };
}

function formatNumber(num: number): string {
  if (num >= 10000) return (num / 10000).toFixed(1) + '万';
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K';
  return num.toString();
}


</script>

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
