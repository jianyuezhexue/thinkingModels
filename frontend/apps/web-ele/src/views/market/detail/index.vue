<template>
  <Page title="模型详情" description="深入了解思维模型，开始你的思考之旅" content-class="p-6 bg-gray-50">
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
      <!-- 返回按钮 -->
      <div class="mb-4">
        <button
          class="flex items-center gap-2 text-gray-600 hover:text-purple-600 transition-colors"
          @click="goBack"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
          </svg>
          <span>返回模型市场</span>
        </button>
      </div>

      <!-- 内容 -->
      <div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <!-- 左侧：模型信息 -->
        <div class="lg:col-span-2 space-y-6">
          <!-- 封面和基本信息 -->
          <ElCard shadow="hover" class="!rounded-xl overflow-hidden">
            <!-- 封面 -->
            <div class="relative h-64 w-full overflow-hidden rounded-xl -mt-5 -mx-5 mb-4" style="width: calc(100% + 40px);">
              <img
                v-if="model.coverImage"
                :src="model.coverImage"
                class="h-full w-full object-cover"
                @error="(e) => { const img = e.target as HTMLImageElement; if (img) img.style.display = 'none'; }"
              />
              <div v-else class="h-full w-full bg-gradient-to-br from-purple-100 to-indigo-100 flex items-center justify-center">
                <span class="text-6xl">{{ model.icon || '📝' }}</span>
              </div>
              <div class="absolute inset-0 bg-gradient-to-t from-black/70 via-black/20 to-transparent"></div>
              <!-- 封面上的信息 -->
              <div class="absolute bottom-6 left-6 right-6">
                <div class="flex items-center gap-3 mb-3">
                  <span
                    :class="[
                      'px-4 py-1.5 rounded-full text-sm font-bold shadow-lg',
                      model.isFree ? 'bg-green-500 text-white' : 'bg-white text-purple-600'
                    ]"
                  >
                    {{ model.isFree ? '🎁 免费' : '💰 ¥' + model.price }}
                  </span>
                  <span class="px-3 py-1 rounded-full text-xs bg-white/90 text-gray-700">
                    {{ model.categoryName || '未分类' }}
                  </span>
                  <span v-if="model.isOfficial" class="px-3 py-1 rounded-full text-xs bg-amber-500 text-white">
                    官方
                  </span>
                </div>
                <h1 class="text-2xl font-bold text-white drop-shadow-lg">{{ model.name }}</h1>
              </div>
            </div>

            <!-- 描述 -->
            <p class="text-gray-600 leading-relaxed">{{ model.description }}</p>

            <!-- 作者信息 -->
            <div class="mt-5 flex items-center justify-between p-4 bg-gray-50 rounded-xl">
              <div class="flex items-center gap-3">
                <ElAvatar :size="48" class="bg-blue-500 ring-2 ring-white shadow-md">
                  {{ model.author?.name?.charAt(0) || '?' }}
                </ElAvatar>
                <div>
                  <div class="font-semibold text-gray-800">{{ model.author?.name || '未知作者' }}</div>
                  <div class="text-sm text-gray-500">模型创建者</div>
                </div>
              </div>
              <div class="text-sm text-gray-400">
                更新于 {{ formatDate(model.updatedAt) }}
              </div>
            </div>

            <!-- 标签 -->
            <div v-if="model.tags?.length" class="mt-4 flex flex-wrap gap-2">
              <ElTag
                v-for="tag in model.tags"
                :key="tag"
                effect="plain"
                class="!bg-purple-50 !text-purple-600 !border-purple-200 !rounded-full"
              >
                {{ tag }}
              </ElTag>
            </div>

            <!-- 统计数据 -->
            <div class="mt-6 grid grid-cols-4 gap-4">
              <div class="text-center p-3 rounded-xl bg-purple-50 hover:bg-purple-100 transition-colors">
                <div class="text-xl font-bold text-purple-600">{{ formatNumber(model.stats.usageCount) }}</div>
                <div class="text-xs text-gray-600 mt-1">使用</div>
              </div>
              <div class="text-center p-3 rounded-xl bg-blue-50 hover:bg-blue-100 transition-colors">
                <div class="text-xl font-bold text-blue-600">{{ formatNumber(model.stats.adoptCount) }}</div>
                <div class="text-xs text-gray-600 mt-1">采纳</div>
              </div>
              <div class="text-center p-3 rounded-xl bg-red-50 hover:bg-red-100 transition-colors">
                <div class="text-xl font-bold text-red-500">{{ formatNumber(model.stats.likeCount) }}</div>
                <div class="text-xs text-gray-600 mt-1">点赞</div>
              </div>
              <div class="text-center p-3 rounded-xl bg-green-50 hover:bg-green-100 transition-colors">
                <div class="text-xl font-bold text-green-600">{{ formatNumber(model.stats.commentCount) }}</div>
                <div class="text-xs text-gray-600 mt-1">评论</div>
              </div>
            </div>
          </ElCard>

          <!-- Tabs -->
          <ElCard shadow="hover" class="!rounded-xl">
            <ElTabs v-model="activeTab" class="model-detail-tabs">
              <!-- 使用指南 Tab -->
              <ElTabPane label="📖 模型概述" name="guide">
                <div class="py-4">
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
                  <div v-if="isContentMarkdown" class="markdown-content prose prose-slate max-w-none"
                    v-html="renderMarkdown(model.content || '')">
                  </div>
                  <!-- 纯文本渲染 -->
                  <div v-else-if="model.content" class="prose max-w-none">
                    <p class="text-gray-600 leading-relaxed whitespace-pre-line">
                      {{ model.content }}
                    </p>
                  </div>
                  <!-- 无内容时显示默认提示 -->
                  <div v-else class="text-gray-500 text-center py-8">
                    暂无概述内容
                  </div>
                </div>
              </ElTabPane>

              <!-- 使用步骤 Tab -->
              <ElTabPane label="📝 使用步骤" name="steps">
                <div class="py-4">
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
                  <ElEmpty v-else description="暂无使用步骤">
                    <template #image>
                      <div class="text-5xl">📋</div>
                    </template>
                  </ElEmpty>
                </div>
              </ElTabPane>

              <!-- 案例 Tab -->
              <ElTabPane label="💡 实践案例" name="examples">
                <div class="py-4">
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
                  <ElEmpty v-else description="暂无实践案例">
                    <template #image>
                      <div class="text-5xl">💡</div>
                    </template>
                  </ElEmpty>
                </div>
              </ElTabPane>

              <!-- 讨论 Tab -->
              <ElTabPane label="💬 讨论" name="discussion">
                <div class="py-4 space-y-6">
                  <!-- 发表评论 -->
                  <div class="rounded-xl bg-gradient-to-r from-purple-50 to-indigo-50 p-5">
                    <h4 class="mb-3 font-semibold text-gray-800">参与讨论</h4>
                    <ElInput
                      v-model="newComment"
                      type="textarea"
                      :rows="3"
                      placeholder="分享你的想法、疑问或经验..."
                      class="mb-3"
                    />
                    <div class="flex items-center justify-between">
                      <ElButton type="primary" class="!bg-purple-600 !border-purple-600 !rounded-full" @click="handleSubmitComment">
                        发布评论
                      </ElButton>
                      <span class="text-sm text-gray-500">{{ comments.length }} 条讨论</span>
                    </div>
                  </div>

                  <!-- 评论列表 -->
                  <div v-if="comments.length > 0" class="space-y-4">
                    <div
                      v-for="comment in comments"
                      :key="comment.id"
                      class="rounded-xl border border-gray-100 p-5 hover:border-purple-200 hover:shadow-sm transition-all"
                    >
                      <div class="flex items-start gap-4">
                        <ElAvatar :size="44" class="flex-shrink-0 bg-blue-500">
                          {{ comment.author.name?.charAt(0) || '?' }}
                        </ElAvatar>
                        <div class="flex-1">
                          <div class="flex items-center gap-2 mb-2">
                            <span class="font-semibold text-gray-800">{{ comment.author.name }}</span>
                            <span class="text-xs text-gray-400">{{ comment.createdAt }}</span>
                          </div>
                          <p class="text-gray-700 leading-relaxed">{{ comment.content }}</p>
                          <div class="mt-3 flex items-center gap-4 text-sm text-gray-500">
                            <button class="flex items-center gap-1 hover:text-purple-600 transition-colors">
                              👍 {{ comment.likes }}
                            </button>
                            <button class="hover:text-purple-600 transition-colors">回复</button>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                  <ElEmpty v-else description="暂无讨论，来发表第一条评论吧！">
                    <template #image>
                      <div class="text-5xl">💭</div>
                    </template>
                  </ElEmpty>
                </div>
              </ElTabPane>
            </ElTabs>
          </ElCard>
        </div>

        <!-- 右侧：操作和推荐 -->
        <div class="space-y-6">
          <!-- 模型信息卡片 -->
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
            </div>
          </ElCard>

          <!-- 操作按钮 -->
          <ElCard shadow="hover" class="!rounded-xl">
            <div class="space-y-3">
              <ElButton
                v-if="model.isFree"
                type="primary"
                size="large"
                class="w-full !bg-purple-600 !border-purple-600 hover:!bg-purple-700 !rounded-full !h-14 !text-base"
                @click="handleFork"
              >
                <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
                </svg>
                引用创建副本
              </ElButton>
              <ElButton
                v-else
                type="success"
                size="large"
                class="w-full !rounded-full !h-16 !text-base !font-semibold"
                @click="handlePurchase"
              >
                <svg class="h-6 w-6 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z"/>
                </svg>
                立即购买 · ¥{{ model.price }}
              </ElButton>
              <ElButton
                size="large"
                class="w-full !rounded-full !h-12"
                @click="handleLike"
              >
                ❤️ 点赞 ({{ formatNumber(model.stats.likeCount) }})
              </ElButton>
            </div>

            <div class="my-4 border-t border-gray-100"></div>

            <ElButton
              type="warning"
              size="large"
              class="w-full !rounded-full !h-12"
              @click="goToCreateTopic"
            >
              <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"/>
              </svg>
              用此模型分析课题
            </ElButton>
          </ElCard>

          <!-- 快速导航 -->
          <ElCard shadow="hover" class="!rounded-xl">
            <template #header>
              <div class="flex items-center gap-2">
                <span class="text-lg">🧭</span>
                <span class="font-semibold text-gray-700">快速导航</span>
              </div>
            </template>
            <div class="space-y-3">
              <div
                class="flex items-center justify-between p-3 rounded-lg hover:bg-gray-50 cursor-pointer transition-colors"
                @click="router.push('/my-models')"
              >
                <span class="text-sm text-gray-600">我的模型库</span>
                <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
                </svg>
              </div>
              <div
                class="flex items-center justify-between p-3 rounded-lg hover:bg-gray-50 cursor-pointer transition-colors"
                @click="goToCreateTopic"
              >
                <span class="text-sm text-gray-600">创建新课题</span>
                <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
                </svg>
              </div>
              <div
                class="flex items-center justify-between p-3 rounded-lg hover:bg-gray-50 cursor-pointer transition-colors"
                @click="router.push(`/market?category=${model.categoryId}`)"
              >
                <span class="text-sm text-gray-600">同类模型</span>
                <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
                </svg>
              </div>
            </div>
          </ElCard>

          <!-- 学习小贴士 -->
          <ElCard shadow="hover" class="!rounded-xl !bg-gradient-to-br from-amber-50 to-orange-50 !border-amber-100">
            <template #header>
              <div class="flex items-center gap-2">
                <span class="text-lg">💡</span>
                <span class="font-semibold text-amber-700">学习小贴士</span>
              </div>
            </template>
            <ul class="text-sm text-amber-800 space-y-2">
              <li class="flex items-start gap-2">
                <span class="text-amber-500">•</span>
                先理解原理再动手实践
              </li>
              <li class="flex items-start gap-2">
                <span class="text-amber-500">•</span>
                结合实际问题反复练习
              </li>
              <li class="flex items-start gap-2">
                <span class="text-amber-500">•</span>
                记录思考过程便于复盘
              </li>
            </ul>
          </ElCard>
        </div>
      </div>
    </template>

    <!-- 空状态 -->
    <ElCard v-else shadow="hover" class="!rounded-xl">
      <ElEmpty description="模型不存在或已被下架">
        <template #image>
          <div class="text-6xl">🔍</div>
        </template>
        <ElButton type="primary" class="!bg-purple-600 !border-purple-600 !rounded-full mt-4" @click="goBack">
          返回市场
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
  breaks: true,
  gfm: true,
});

import {
  ElButton,
  ElCard,
  ElTabs,
  ElTabPane,
  ElInput,
  ElAvatar,
  ElTag,
  ElEmpty,
  ElMessage,
  ElSkeleton,
  ElSkeletonItem,
} from 'element-plus';

import {
  getThinkingModelDetailApi,
  forkThinkingModelApi,
  likeThinkingModelApi,
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

// 当前激活的Tab
const activeTab = ref('guide');

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
  status: number;
  price: number;
  isFree: boolean;
  isOfficial: boolean;
  content: string;
  usageGuide: string;
  examples: string;
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

// 评论相关
const newComment = ref('');
const comments = ref<Array<{
  id: string;
  author: { name: string };
  content: string;
  createdAt: string;
  likes: number;
}>>([]);

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

// 解析视频 URL
function parseVideoUrl(url: string): VideoEmbed | null {
  if (!url || !url.trim()) return null;

  const bilibiliMatch = url.match(/bilibili\.com\/video\/(BV[a-zA-Z0-9]+)/);
  if (bilibiliMatch) {
    const bvid = bilibiliMatch[1];
    const pageMatch = url.match(/[?&]p=(\d+)/);
    const page = pageMatch ? pageMatch[1] : '1';
    return {
      platform: 'bilibili',
      embedUrl: `//player.bilibili.com/player.html?bvid=${bvid}&page=${page}&high_quality=1&danmaku=0`,
      title: 'Bilibili 视频',
    };
  }

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

function hasValidVideo(videoUrl: string): boolean {
  return parseVideoUrl(videoUrl) !== null;
}

function getVideoEmbed(videoUrl: string): VideoEmbed | null {
  return parseVideoUrl(videoUrl);
}

// 判断内容是否为 Markdown 格式
function isMarkdown(content: string): boolean {
  if (!content || !content.trim()) return false;

  const markdownPatterns = [
    /^#{1,6}\s+.+/m,
    /\*\*.+?\*\*/,
    /\*.+?\*/,
    /^\s*[-*+]\s+.+/m,
    /^\s*\d+\.\s+.+/m,
    /\[.+?\]\(.+?\)/,
    /!\[.*?\]\(.+?\)/,
    /^```/m,
    /`.+?`/,
    /^>\s+.+/m,
    /^\s*---+\s*$/m,
    /^\|.+\|/m,
  ];

  let matchCount = 0;
  for (const pattern of markdownPatterns) {
    if (pattern.test(content)) {
      matchCount++;
      if (matchCount >= 2) return true;
    }
  }

  return false;
}

// 渲染 Markdown
function renderMarkdown(content: string): string {
  if (!content) return '';
  try {
    const rawHtml = marked.parse(content) as string;
    return DOMPurify.sanitize(rawHtml);
  } catch {
    return content;
  }
}

// 判断当前模型内容是否为 Markdown
const isContentMarkdown = computed(() => {
  return model.value ? isMarkdown(model.value.content) : false;
});

// 引用创建副本
async function handleFork() {
  if (!model.value) return;
  try {
    await forkThinkingModelApi({
      sourceModelId: model.value.id,
      name: `${model.value.name} (副本)`,
      description: model.value.description,
    });
    ElMessage.success('已创建副本到您的模型库');
    router.push('/my-models');
  } catch (error) {
    console.error('创建副本失败:', error);
    ElMessage.error('创建副本失败');
  }
}

// 购买模型
function handlePurchase() {
  ElMessage.info('购买功能开发中，敬请期待');
}

// 点赞模型
async function handleLike() {
  if (!model.value) return;
  try {
    const res = await likeThinkingModelApi(model.value.id);
    // 更新点赞数
    model.value.stats.likeCount = res.likeCount;
    if (res.liked) {
      ElMessage.success('点赞成功');
    } else {
      ElMessage.info('今天已经点赞过了');
    }
  } catch (error) {
    console.error('点赞失败:', error);
    ElMessage.error('点赞失败');
  }
}

// 发表评论
function handleSubmitComment() {
  if (!newComment.value.trim()) {
    ElMessage.warning('请输入评论内容');
    return;
  }
  comments.value.unshift({
    id: Date.now().toString(),
    author: { name: '我' },
    content: newComment.value,
    createdAt: new Date().toLocaleString('zh-CN'),
    likes: 0,
  });
  newComment.value = '';
  ElMessage.success('评论已发布');
}

// 返回市场
function goBack() {
  router.push('/market');
}

// 跳转到创建课题页面
function goToCreateTopic() {
  router.push('/my-topics/create');
}
</script>

<style scoped>
.model-detail-tabs :deep(.el-tabs__header) {
  margin-bottom: 0;
}

.model-detail-tabs :deep(.el-tabs__nav-wrap::after) {
  display: none;
}

.model-detail-tabs :deep(.el-tabs__item) {
  padding: 0 20px;
  height: 48px;
  line-height: 48px;
  font-weight: 500;
}

.model-detail-tabs :deep(.el-tabs__item.is-active) {
  color: #7c3aed;
}

.model-detail-tabs :deep(.el-tabs__active-bar) {
  background-color: #7c3aed;
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