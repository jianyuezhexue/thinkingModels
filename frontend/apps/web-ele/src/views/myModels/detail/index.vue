<template>
  <Page description="管理模型详情、查看数据分析" title="模型详情" content-class="p-6 bg-gray-50">
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
        <button class="flex items-center gap-2 text-gray-600 hover:text-purple-600 transition-colors" @click="goBack">
          <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
          <span>返回列表</span>
        </button>
        <div class="flex items-center gap-3">
          <ElButton v-if="model.status === 0" type="primary"
            class="!bg-purple-600 !border-purple-600 hover:!bg-purple-700 !rounded-full" @click="handlePublish">
            <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
            </svg>
            提交审核
          </ElButton>
          <!-- <ElButton v-if="model.status === 1" plain class="!rounded-full" @click="handleUnpublish">
            下架模型
          </ElButton> -->
          <ElButton class="!rounded-full" @click="handleEdit">
            <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>
            编辑
          </ElButton>
          <!-- <ElButton type="danger" plain class="!rounded-full" @click="handleDelete">
            <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
          </ElButton> -->
        </div>
      </div>

      <div class="flex gap-6">
        <!-- 左侧主内容 -->
        <div class="flex-1 min-w-0 space-y-6">
          <!-- 模型头卡 -->
          <ElCard shadow="hover" class="!rounded-xl overflow-hidden">
            <div class="flex flex-col lg:flex-row gap-6">
              <!-- 封面 -->
              <div
                class="relative w-full lg:w-72 h-48 rounded-xl overflow-hidden bg-gradient-to-br from-purple-100 to-indigo-100 flex-shrink-0 flex items-center justify-center">
                <img v-if="model.coverImage" :src="model.coverImage" class="w-full h-full object-cover"
                  @error="(e) => { const img = e.target as HTMLImageElement; if (img) img.style.display = 'none'; }" />
                <div v-else class="text-5xl">{{ model.icon || '📝' }}</div>
                <!-- 状态标签 -->
                <span :class="[
                  'absolute left-3 top-3 rounded-full px-3 py-1 text-sm font-medium',
                  getStatusStyle(model.status).bg,
                  getStatusStyle(model.status).text
                ]">
                  {{ getStatusStyle(model.status).label }}
                </span>
                <!-- 价格标签 -->
                <span :class="[
                  'absolute right-3 top-3 rounded-full px-3 py-1 text-sm font-bold shadow-lg',
                  model.isFree ? 'bg-green-500 text-white' : 'bg-white text-purple-600'
                ]">
                  {{ model.isFree ? '免费' : '¥' + model.price }}
                </span>
              </div>

              <!-- 信息 -->
              <div class="flex-1">
                <h1 class="text-2xl font-bold text-gray-900 mb-2">{{ model.name }}</h1>
                <p class="text-gray-500 mb-4 leading-relaxed">{{ model.description }}</p>

                <!-- 标签 -->
                <div v-if="model.tags?.length" class="flex flex-wrap gap-2 mb-4">
                  <ElTag v-for="tag in model.tags" :key="tag" effect="plain"
                    class="!bg-purple-50 !text-purple-600 !border-purple-200 !rounded-full">
                    {{ tag }}
                  </ElTag>
                </div>

                <!-- 统计栏 -->
                <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 pt-4 border-t border-gray-100">
                  <div class="text-center">
                    <div class="text-2xl font-bold text-purple-600">{{ formatNumber(model.stats.usageCount) }}</div>
                    <div class="text-sm text-gray-500">使用次数</div>
                  </div>
                  <div class="text-center">
                    <div class="text-2xl font-bold text-blue-600">{{ formatNumber(model.stats.adoptCount) }}</div>
                    <div class="text-sm text-gray-500">被采纳</div>
                  </div>
                  <div class="text-center">
                    <div class="text-2xl font-bold text-red-500">{{ formatNumber(model.stats.likeCount) }}</div>
                    <div class="text-sm text-gray-500">获赞</div>
                  </div>
                  <div class="text-center">
                    <div class="text-2xl font-bold text-green-600">{{ formatNumber(model.stats.commentCount) }}</div>
                    <div class="text-sm text-gray-500">评论数</div>
                  </div>
                </div>
              </div>
            </div>
          </ElCard>

          <!-- Tab 导航 -->
          <div class="flex flex-wrap gap-2">
            <button v-for="tab in tabs" :key="tab.id"
              class="px-5 py-2.5 rounded-full text-sm font-medium transition-all" :class="[
                activeTab === tab.id
                  ? 'bg-purple-100 text-purple-700 shadow-md border border-purple-200 font-semibold'
                  : 'bg-white text-gray-600 hover:bg-purple-50 hover:text-purple-600 border border-gray-200'
              ]" @click="activeTab = tab.id as any">
              {{ tab.icon }} {{ tab.label }}
            </button>
          </div>

          <!-- Tab 内容 -->
          <!-- 概述 -->
          <ElCard v-if="activeTab === 'content'" shadow="hover" class="!rounded-xl">
            <template #header>
              <div class="flex items-center gap-2">
                <span class="text-lg">📋</span>
                <span class="font-semibold text-gray-700">模型概述</span>
              </div>
            </template>

            <!-- 视频展示区域 -->
            <div v-if="model.videoUrl && hasValidVideo(model.videoUrl)" class="mb-6">
              <div class="aspect-video w-full rounded-xl overflow-hidden bg-gray-100">
                <iframe :src="getVideoEmbed(model.videoUrl)?.embedUrl" :title="getVideoEmbed(model.videoUrl)?.title"
                  class="w-full h-full" frameborder="0"
                  allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                  allowfullscreen></iframe>
              </div>
              <div class="mt-2 flex items-center gap-2 text-sm text-gray-500">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
                </svg>
                <span>{{ getVideoEmbed(model.videoUrl)?.platform === 'bilibili' ? 'Bilibili' : 'YouTube' }} 视频教程</span>
              </div>
            </div>

            <!-- Markdown 渲染 -->
            <div v-if="isContentMarkdown" class="markdown-content prose prose-slate max-w-none" v-html="renderMarkdown(model.content || '')">
            </div>
            <!-- 纯文本渲染 -->
            <div v-else class="prose max-w-none">
              <p class="text-gray-600 leading-relaxed whitespace-pre-line">
                {{ model.content || '暂无概述' }}
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
            <div v-if="parsedSteps.length" class="space-y-4">
              <div v-for="(step, index) in parsedSteps" :key="index"
                class="flex gap-4 p-4 rounded-xl bg-gradient-to-r from-purple-50 to-indigo-50 border border-purple-100">
                <div
                  class="w-10 h-10 rounded-full bg-purple-600 text-white flex items-center justify-center font-bold flex-shrink-0">
                  {{ index + 1 }}
                </div>
                <div class="flex-1">
                  <h4 class="font-semibold text-gray-800 mb-2">{{ step.title }}</h4>
                  <p class="text-sm text-gray-600 leading-relaxed">{{ step.description }}</p>
                </div>
              </div>
            </div>
            <div v-else class="text-center text-gray-400 py-8">
              暂无使用步骤
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
            <div v-if="parsedExamples.length" class="space-y-6">
              <div v-for="(example, index) in parsedExamples" :key="index"
                class="p-5 rounded-xl bg-gradient-to-br from-gray-50 to-slate-50 border border-gray-100">
                <h4 class="font-semibold text-gray-800 mb-3 flex items-center gap-2">
                  <span
                    class="w-6 h-6 rounded-full bg-amber-100 text-amber-600 flex items-center justify-center text-sm">
                    {{ index + 1 }}
                  </span>
                  {{ example.title }}
                </h4>
                <p class="text-sm text-gray-600 leading-relaxed whitespace-pre-line">{{ example.content }}</p>
              </div>
            </div>
            <div v-else class="text-center text-gray-400 py-8">
              暂无实践案例
            </div>
          </ElCard>

          <!-- 用户反馈 -->
          <ElCard v-if="activeTab === 'feedback'" shadow="hover" class="!rounded-xl">
            <template #header>
              <div class="flex items-center gap-2">
                <span class="text-lg">💬</span>
                <span class="font-semibold text-gray-700">用户反馈</span>
              </div>
            </template>
            <div class="text-center text-gray-400 py-8">
              暂无用户反馈
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
                  <div class="text-2xl font-bold text-purple-600 mb-1">{{ formatNumber(model.stats.usageCount) }}</div>
                  <div class="text-sm text-gray-500">总使用</div>
                </div>
                <div class="p-4 rounded-xl bg-blue-50 border border-blue-100 text-center">
                  <div class="text-2xl font-bold text-blue-600 mb-1">{{ formatNumber(model.stats.adoptCount) }}</div>
                  <div class="text-sm text-gray-500">采纳数</div>
                </div>
                <div class="p-4 rounded-xl bg-red-50 border border-red-100 text-center">
                  <div class="text-2xl font-bold text-red-500 mb-1">{{ formatNumber(model.stats.likeCount) }}</div>
                  <div class="text-sm text-gray-500">点赞数</div>
                </div>
                <div class="p-4 rounded-xl bg-green-50 border border-green-100 text-center">
                  <div class="text-2xl font-bold text-green-600 mb-1">{{ formatNumber(model.stats.commentCount) }}</div>
                  <div class="text-sm text-gray-500">评论数</div>
                </div>
              </div>
            </div>
          </ElCard>
        </div>

        <!-- 右侧边栏 -->
        <div class="w-80 flex-shrink-0 space-y-6 hidden lg:block">
          <!-- 模型信息 -->
          <ElCard shadow="hover" class="!rounded-xl">
            <template #header>
              <div class="flex items-center gap-2">
                <span class="text-lg">ℹ️</span>
                <span class="font-semibold text-gray-700">模型信息</span>
              </div>
            </template>
            <div class="space-y-3 text-sm">
              <div class="flex items-center justify-between">
                <span class="text-gray-500">分类</span>
                <span class="text-gray-800">{{ model.categoryName || '-' }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-gray-500">难度</span>
                <span class="px-2 py-0.5 rounded-full text-xs font-medium"
                  :style="{ backgroundColor: `${getDifficultyConfig(model.difficulty).color}20`, color: getDifficultyConfig(model.difficulty).color }">
                  {{ getDifficultyConfig(model.difficulty).label }}
                </span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-gray-500">预计用时</span>
                <span class="text-gray-800">{{ model.estimatedTime }} 分钟</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-gray-500">版本</span>
                <span class="text-gray-800">v{{ model.version || '1.0' }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-gray-500">创建时间</span>
                <span class="text-gray-800">{{ formatDate(model.createdAt) }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-gray-500">最后更新</span>
                <span class="text-gray-800">{{ formatDate(model.updatedAt) }}</span>
              </div>
            </div>
          </ElCard>

          <!-- 驳回原因（仅驳回状态显示） -->
          <ElCard v-if="model.status === 4 && model.reviewNote" shadow="hover"
            class="!rounded-xl !bg-red-50 !border-red-100">
            <template #header>
              <div class="flex items-center gap-2">
                <span class="text-lg">⚠️</span>
                <span class="font-semibold text-red-700">驳回原因</span>
              </div>
            </template>
            <p class="text-sm text-red-600">{{ model.reviewNote }}</p>
          </ElCard>

          <!-- 作者信息 -->
          <ElCard shadow="hover" class="!rounded-xl">
            <template #header>
              <div class="flex items-center gap-2">
                <span class="text-lg">👤</span>
                <span class="font-semibold text-gray-700">作者信息</span>
              </div>
            </template>
            <div class="flex items-center gap-3">
              <ElAvatar :size="48" class="bg-blue-500">
                {{ model.author?.name?.charAt(0) || '?' }}
              </ElAvatar>
              <div>
                <div class="font-medium text-gray-800">{{ model.author?.name || '未知作者' }}</div>
                <div class="text-xs text-gray-400">模型创建者</div>
              </div>
            </div>
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

<script lang="ts" setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { Page } from '@vben/common-ui';

import { marked } from 'marked';
import DOMPurify from 'dompurify';

// 配置 marked 选项
marked.setOptions({
  breaks: true, // 支持 GFM 换行
  gfm: true, // 启用 GitHub Flavored Markdown
});

import {
  ElButton,
  ElCard,
  ElTag,
  ElMessage,
  ElEmpty,
  ElSkeleton,
  ElSkeletonItem,
  ElMessageBox,
  ElAvatar,
} from 'element-plus';

import {
  getThinkingModelDetailApi,
  deleteThinkingModelApi,
  publishThinkingModelApi,
  unpublishThinkingModelApi,
} from '#/api/thinking/model';
import { getAllCategoriesApi } from '#/api/master/category';

// 路由
const route = useRoute();
const router = useRouter();
const modelId = computed(() => Number(route.params.id));

// 请求取消控制器
let abortController: AbortController | null = null;

// 加载状态
const loading = ref(true);

// 当前选中 Tab
const activeTab = ref<'content' | 'steps' | 'examples' | 'feedback' | 'analytics'>('content');

// Tab 配置
const tabs = [
  { id: 'content', label: '模型概述', icon: '📋' },
  { id: 'steps', label: '使用步骤', icon: '📝' },
  { id: 'examples', label: '实践案例', icon: '💡' },
  { id: 'feedback', label: '用户反馈', icon: '💬' },
  { id: 'analytics', label: '数据分析', icon: '📊' },
];

// 分类选项
const categoryOptions = ref<{ value: number; label: string }[]>([]);

// 模型数据类型
interface ModelDetail {
  id: number;
  name: string;
  description: string;
  coverImage: string;
  videoUrl: string;
  icon: string;
  categoryId: number;
  categoryName?: string;
  tags: string[] | null;
  status: number; // 0=草稿, 1=已发布, 2=已下架, 3=审核中, 4=已驳回
  price: number;
  isFree: boolean;
  content: string;
  usageGuide: string;  // 使用指南（JSON 字符串）
  examples: string;    // 案例（JSON 字符串）
  difficulty: number;
  estimatedTime: number;
  version: string;
  author: {
    id: string;
    name: string;
    avatar: string;
  };
  stats: {
    usageCount: number;
    adoptCount: number;
    likeCount: number;
    commentCount: number;
  };
  createdAt: string;
  updatedAt: string;
  reviewNote?: string;
}

// 解析后的步骤和案例
interface Step {
  title: string;
  description: string;
}

interface Example {
  title: string;
  content: string;
}

const model = ref<ModelDetail | null>(null);
const parsedSteps = ref<Step[]>([]);
const parsedExamples = ref<Example[]>([]);

// 解析 JSON 字符串数组
function parseJsonArray<T>(jsonStr: string): T[] {
  if (!jsonStr) return [];
  try {
    const parsed = JSON.parse(jsonStr);
    return Array.isArray(parsed) ? parsed : [];
  } catch {
    return [];
  }
}

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

// 获取模型详情
async function fetchModelDetail() {
  abortController?.abort();
  abortController = new AbortController();
  const signal = abortController.signal;

  loading.value = true;
  try {
    const res = await getThinkingModelDetailApi(modelId.value);
    model.value = {
      ...res,
      categoryName: categoryOptions.value.find((c) => c.value === res.categoryId)?.label || '',
    };
    // 解析使用步骤和案例
    parsedSteps.value = parseJsonArray<Step>(res.usageGuide || '');
    parsedExamples.value = parseJsonArray<Example>(res.examples || '');
  } catch (error: any) {
    if (error?.name === 'CanceledError' || error?.code === 'ERR_CANCELED' || signal.aborted) {
      return;
    }
    console.error('获取模型详情失败:', error);
    ElMessage.error('获取模型详情失败');
  } finally {
    loading.value = false;
  }
}

onMounted(async () => {
  await loadCategories();
  await fetchModelDetail();
});

onUnmounted(() => {
  abortController?.abort();
});

// 操作函数
function goBack() {
  router.push('/my-models');
}

function handleEdit() {
  router.push(`/my-models/create?id=${modelId.value}`);
}

async function handleDelete() {
  if (!model.value) return;
  try {
    await ElMessageBox.confirm(
      `确定要删除模型「${model.value.name}」吗？此操作不可恢复。`,
      '删除确认',
      { type: 'warning' }
    );
    await deleteThinkingModelApi({ ids: [model.value.id] });
    ElMessage.success('模型已删除');
    router.push('/my-models');
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('删除失败:', error);
      ElMessage.error('删除失败');
    }
  }
}

async function handlePublish() {
  if (!model.value) return;
  try {
    await ElMessageBox.confirm(
      '提交审核后，模型将在审核通过后发布到市场。确定提交吗？',
      '提交审核',
      { type: 'info' }
    );
    await publishThinkingModelApi({ id: model.value.id });
    ElMessage.success('模型已提交审核');
    fetchModelDetail();
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('提交审核失败:', error);
      ElMessage.error('提交审核失败');
    }
  }
}

async function handleUnpublish() {
  if (!model.value) return;
  try {
    await ElMessageBox.confirm(
      '下架后用户将无法继续购买此模型。确定下架吗？',
      '确认下架',
      { type: 'warning' }
    );
    await unpublishThinkingModelApi(model.value.id);
    ElMessage.success('模型已下架');
    fetchModelDetail();
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('下架失败:', error);
      ElMessage.error('下架失败');
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

function formatDate(dateStr: string): string {
  if (!dateStr) return '-';
  const date = new Date(dateStr);
  return date.toLocaleDateString('zh-CN', { year: 'numeric', month: 'short', day: 'numeric' });
}

// 获取难度配置
function getDifficultyConfig(difficulty: number) {
  const config: Record<number, { label: string; color: string }> = {
    1: { label: '简单', color: '#67C23A' },
    2: { label: '中等', color: '#E6A23C' },
    3: { label: '困难', color: '#F56C6C' },
  };
  return config[difficulty] || { label: '未知', color: '#909399' };
}

// 视频平台类型
type VideoPlatform = 'bilibili' | 'youtube' | 'unknown';

// 视频嵌入信息
interface VideoEmbed {
  platform: VideoPlatform;
  embedUrl: string;
  title: string;
}

// 解析视频 URL，返回嵌入信息
function parseVideoUrl(url: string): VideoEmbed | null {
  if (!url || !url.trim()) return null;

  // Bilibili 视频匹配
  // 支持: https://www.bilibili.com/video/BV1xxx 或 https://www.bilibili.com/video/BV1xxx?p=1
  const bilibiliMatch = url.match(/bilibili\.com\/video\/(BV[a-zA-Z0-9]+)/);
  if (bilibiliMatch) {
    const bvid = bilibiliMatch[1];
    // 提取分P参数
    const pageMatch = url.match(/[?&]p=(\d+)/);
    const page = pageMatch ? pageMatch[1] : '1';
    return {
      platform: 'bilibili',
      embedUrl: `//player.bilibili.com/player.html?bvid=${bvid}&page=${page}&high_quality=1&danmaku=0`,
      title: 'Bilibili 视频',
    };
  }

  // YouTube 视频匹配
  // 支持: https://www.youtube.com/watch?v=xxx 或 https://youtu.be/xxx
  const youtubeMatch = url.match(/(?:youtube\.com\/watch\?v=|youtu\.be\/)([a-zA-Z0-9_-]+)/);
  if (youtubeMatch) {
    const videoId = youtubeMatch[1];
    return {
      platform: 'youtube',
      embedUrl: `https://www.youtube.com/embed/${videoId}`,
      title: 'YouTube 视频',
    };
  }

  return null;
}

// 判断是否有有效视频
function hasValidVideo(videoUrl: string): boolean {
  return parseVideoUrl(videoUrl) !== null;
}

// 获取视频嵌入信息
function getVideoEmbed(videoUrl: string): VideoEmbed | null {
  return parseVideoUrl(videoUrl);
}

// 判断内容是否为 Markdown 格式
function isMarkdown(content: string): boolean {
  if (!content || !content.trim()) return false;

  // 常见的 Markdown 特征标记
  const markdownPatterns = [
    /^#{1,6}\s+.+/m,           // 标题 # ## ### 等
    /\*\*.+?\*\*/,             // 粗体 **text**
    /\*.+?\*/,                 // 斜体 *text*
    /^\s*[-*+]\s+.+/m,         // 无序列表
    /^\s*\d+\.\s+.+/m,         // 有序列表
    /\[.+?\]\(.+?\)/,          // 链接 [text](url)
    /!\[.*?\]\(.+?\)/,         // 图片 ![alt](url)
    /^```/m,                   // 代码块
    /`.+?`/,                   // 行内代码
    /^>\s+.+/m,                // 引用块
    /^\s*---+\s*$/m,           // 分隔线
    /^\|.+\|/m,                // 表格
  ];

  // 检测到 2 个或以上 Markdown 特征则认为是 Markdown
  let matchCount = 0;
  for (const pattern of markdownPatterns) {
    if (pattern.test(content)) {
      matchCount++;
      if (matchCount >= 2) return true;
    }
  }

  return false;
}

// 渲染 Markdown 内容为安全的 HTML
function renderMarkdown(content: string): string {
  if (!content) return '';

  try {
    // 解析 Markdown 为 HTML
    const rawHtml = marked.parse(content) as string;
    // 使用 DOMPurify 清理 XSS 风险
    return DOMPurify.sanitize(rawHtml);
  } catch {
    return content;
  }
}

// 判断当前模型内容是否为 Markdown
const isContentMarkdown = computed(() => {
  return model.value ? isMarkdown(model.value.content) : false;
});
</script>

<style scoped>
.prose {
  color: #374151;
}

/* Markdown 渲染样式 */
.markdown-content {
  color: #374151;
  line-height: 1.75;
}

.markdown-content :deep(h1) {
  font-size: 1.875rem;
  font-weight: 700;
  margin-bottom: 1rem;
  margin-top: 2rem;
  color: #111827;
}

.markdown-content :deep(h2) {
  font-size: 1.5rem;
  font-weight: 600;
  margin-bottom: 0.75rem;
  margin-top: 1.75rem;
  color: #111827;
  border-bottom: 1px solid #e5e7eb;
  padding-bottom: 0.5rem;
}

.markdown-content :deep(h3) {
  font-size: 1.25rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
  margin-top: 1.5rem;
  color: #111827;
}

.markdown-content :deep(h4),
.markdown-content :deep(h5),
.markdown-content :deep(h6) {
  font-size: 1.125rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
  margin-top: 1.25rem;
  color: #374151;
}

.markdown-content :deep(p) {
  margin-bottom: 1rem;
}

.markdown-content :deep(ul),
.markdown-content :deep(ol) {
  margin-bottom: 1rem;
  padding-left: 1.5rem;
}

.markdown-content :deep(ul) {
  list-style-type: disc;
}

.markdown-content :deep(ol) {
  list-style-type: decimal;
}

.markdown-content :deep(li) {
  margin-bottom: 0.5rem;
}

.markdown-content :deep(a) {
  color: #7c3aed;
  text-decoration: underline;
  transition: color 0.2s;
}

.markdown-content :deep(a:hover) {
  color: #6d28d9;
}

.markdown-content :deep(blockquote) {
  border-left: 4px solid #7c3aed;
  padding-left: 1rem;
  margin: 1rem 0;
  color: #6b7280;
  background-color: #f9fafb;
  padding: 1rem;
  border-radius: 0 0.5rem 0.5rem 0;
}

.markdown-content :deep(code) {
  background-color: #f3f4f6;
  padding: 0.125rem 0.375rem;
  border-radius: 0.25rem;
  font-size: 0.875em;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  color: #dc2626;
}

.markdown-content :deep(pre) {
  background-color: #1f2937;
  color: #f9fafb;
  padding: 1rem;
  border-radius: 0.5rem;
  overflow-x: auto;
  margin: 1rem 0;
}

.markdown-content :deep(pre code) {
  background-color: transparent;
  padding: 0;
  color: inherit;
  font-size: 0.875rem;
}

.markdown-content :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin: 1rem 0;
}

.markdown-content :deep(th),
.markdown-content :deep(td) {
  border: 1px solid #e5e7eb;
  padding: 0.75rem 1rem;
  text-align: left;
}

.markdown-content :deep(th) {
  background-color: #f9fafb;
  font-weight: 600;
}

.markdown-content :deep(tr:nth-child(even)) {
  background-color: #f9fafb;
}

.markdown-content :deep(hr) {
  border: none;
  border-top: 2px solid #e5e7eb;
  margin: 2rem 0;
}

.markdown-content :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 0.5rem;
  margin: 1rem 0;
}

.markdown-content :deep(strong) {
  font-weight: 600;
  color: #111827;
}

.markdown-content :deep(em) {
  font-style: italic;
}
</style>