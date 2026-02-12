<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { Page } from '@vben/common-ui';

import {
  ElButton,
  ElCard,
  ElTabs,
  ElTabPane,
  ElInput,
  ElAvatar,
  ElTag,
  ElEmpty,
  ElDivider,
  ElMessage,
  ElSkeleton,
  ElSkeletonItem,
} from 'element-plus';

import {
  getModelDetailApi,
  getRecommendedModelsApi,
  adoptModelApi,
  purchaseModelApi,
  forkModelApi,
  likeModelApi,
  type ModelApi,
} from '#/api';

// 路由
const route = useRoute();
const router = useRouter();
const modelId = computed(() => route.params.id as string);

// 加载状态
const loading = ref(true);

// 模型数据
const model = ref<ModelApi.ThinkingModel | null>(null);

// 相关推荐
const relatedModels = ref<ModelApi.ThinkingModel[]>([]);

// 当前激活的Tab
const activeTab = ref('guide');

// 练习记录（模拟数据）
const practiceRecords = ref([
  {
    id: 'p1',
    user: { id: 'u20', name: '王思维', avatar: 'https://avatar.vercel.sh/wangsw.svg?text=WS', isCertified: true },
    topicTitle: '如何用SWOT分析评估新产品上线风险',
    summary: '通过SWOT分析框架，我系统评估了我们团队新开发的AI助手产品的市场前景。优势在于技术领先，劣势是品牌知名度不足，机会是市场需求增长，威胁是巨头竞争激烈...',
    createdAt: '2024-02-18 15:30',
    views: 1256,
    likes: 89,
    comments: 23,
    isExcellent: true,
  },
  {
    id: 'p2',
    user: { id: 'u21', name: '李策划', avatar: 'https://avatar.vercel.sh/lich.svg?text=LC', isCertified: false },
    topicTitle: 'SWOT分析在创业计划书中的应用',
    summary: '在准备融资计划书时，使用SWOT分析帮助投资人快速理解我们的商业模式。将内部资源能力与外部环境结合分析，让融资路演更有说服力...',
    createdAt: '2024-02-16 09:15',
    views: 892,
    likes: 56,
    comments: 12,
    isExcellent: true,
  },
  {
    id: 'p3',
    user: { id: 'u22', name: '张产品', avatar: '', isCertified: true },
    topicTitle: '产品迭代中的SWOT实战案例',
    summary: '分享一个真实案例：我们在做产品迭代决策时，通过SWOT分析发现了被忽视的技术债务风险，及时调整了优先级，避免了后期大规模重构...',
    createdAt: '2024-02-14 16:45',
    views: 2341,
    likes: 178,
    comments: 45,
    isExcellent: false,
  },
  {
    id: 'p4',
    user: { id: 'u23', name: '陈决策', avatar: 'https://avatar.vercel.sh/chenjc.svg?text=CJ', isCertified: false },
    topicTitle: '个人职业发展的SWOT分析',
    summary: '用SWOT分析自己的职业发展路径，识别了需要提升的技能和可以把握的行业机会。这种方法不仅适用于企业，对个人规划也很有帮助...',
    createdAt: '2024-02-12 11:20',
    views: 567,
    likes: 34,
    comments: 8,
    isExcellent: false,
  },
]);

// 查看练习详情
function viewPracticeDetail(id: string) {
  router.push(`/practices/${id}`);
}

// 评论相关（模拟数据）
const newComment = ref('');
const comments = ref([
  {
    id: 'c1',
    author: { id: 'u10', name: '李思考', avatar: 'https://avatar.vercel.sh/lisk.svg?text=LS' },
    content: '这个模型在实际工作中非常有用，特别是在做产品规划的时候。建议大家多练习！',
    createdAt: '2024-02-15 14:30',
    likes: 23,
  },
  {
    id: 'c2',
    author: { id: 'u12', name: '赵分析', avatar: 'https://avatar.vercel.sh/zhaofx.svg?text=ZF' },
    content: '有没有人可以分享一下如何在团队会议中引导大家使用这个思维模型？',
    createdAt: '2024-02-14 09:15',
    likes: 15,
  },
]);

// 获取模型详情
async function fetchModelDetail() {
  loading.value = true;
  try {
    const res = await getModelDetailApi(modelId.value);
    model.value = res;
    // 获取相关推荐
    fetchRelatedModels(res.category);
  } catch (error) {
    console.error('获取模型详情失败:', error);
    ElMessage.error('获取模型详情失败');
  } finally {
    loading.value = false;
  }
}

// 获取相关推荐
async function fetchRelatedModels(category: string) {
  try {
    const res = await getRecommendedModelsApi(category, 4);
    relatedModels.value = res.filter((m) => m.id !== modelId.value).slice(0, 3);
  } catch (error) {
    console.error('获取推荐模型失败:', error);
  }
}

// 页面加载时获取数据
onMounted(() => {
  fetchModelDetail();
});

// 格式化数字
function formatNumber(num: number): string {
  if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M';
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K';
  return num.toString();
}

// 加载模型
async function handleLoad() {
  if (!model.value) return;
  try {
    await adoptModelApi(model.value.id);
    ElMessage.success('已成功加载到您的模型库');
  } catch (error) {
    console.error('加载失败:', error);
  }
}

// 购买模型
async function handlePurchase() {
  if (!model.value) return;
  try {
    await purchaseModelApi(model.value.id);
    ElMessage.success('购买成功！已添加到您的模型库');
  } catch (error) {
    console.error('购买失败:', error);
  }
}

// 引用模型
async function handleFork() {
  if (!model.value) return;
  try {
    await forkModelApi(model.value.id);
    ElMessage.success('已创建副本到您的模型库');
  } catch (error) {
    console.error('引用失败:', error);
  }
}

// 点赞模型
async function handleLike() {
  if (!model.value) return;
  try {
    await likeModelApi(model.value.id);
    model.value.stats.likes++;
    ElMessage.success('已点赞');
  } catch (error) {
    console.error('点赞失败:', error);
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
    author: { id: 'me', name: '我', avatar: 'https://avatar.vercel.sh/me.svg?text=ME' },
    content: newComment.value,
    createdAt: new Date().toLocaleString('zh-CN'),
    likes: 0,
  });
  newComment.value = '';
  ElMessage.success('评论已发布');
}

// 跳转到相关模型
function goToRelatedModel(id: string) {
  router.push(`/market/${id}`);
}

// 返回市场
function goBack() {
  router.push('/market');
}

// 跳转到创建课题页面
function goToCreateTopic() {
  router.push('/my-topics/create');
}

// 格式化 Markdown 内容为 HTML
const formattedContent = computed(() => {
  if (!model.value?.content) return '';
  return model.value.content
    .replace(/\n/g, '<br>')
    .replace(/## (.*)/g, '<h2 class="text-xl font-bold mt-6 mb-3">$1</h2>')
    .replace(/### (.*)/g, '<h3 class="text-lg font-semibold mt-4 mb-2">$1</h3>')
    .replace(/\d\. \*\*(.*)\*\*/g, '<strong>$1</strong>');
});
</script>

<template>
  <Page title="模型详情" description="深入了解思维模型，开始你的思考之旅">
    <!-- 加载状态 -->
    <div v-if="loading" class="grid grid-cols-1 gap-6 lg:grid-cols-3">
      <div class="lg:col-span-2 space-y-6">
        <ElCard shadow="never">
          <ElSkeleton animated>
            <template #template>
              <ElSkeletonItem variant="image" style="width: 100%; height: 256px" />
              <div class="mt-4 space-y-3">
                <ElSkeletonItem variant="p" style="width: 50%" />
                <ElSkeletonItem variant="text" style="width: 30%" />
              </div>
            </template>
          </ElSkeleton>
        </ElCard>
      </div>
      <div class="space-y-6">
        <ElCard shadow="never">
          <ElSkeleton :rows="3" animated />
        </ElCard>
      </div>
    </div>

    <!-- 内容 -->
    <div v-else-if="model" class="grid grid-cols-1 gap-6 lg:grid-cols-3">
      <!-- 左侧：模型信息 -->
      <div class="lg:col-span-2 space-y-6">
        <!-- 封面和基本信息 -->
        <ElCard shadow="never">
          <div class="relative h-64 w-full overflow-hidden rounded-lg">
            <img
              :src="model.cover || '/images/default-model-cover.svg'"
              class="h-full w-full object-cover"
              @error="(e) => { const img = e.target as HTMLImageElement; if (img) img.src = '/images/default-model-cover.svg'; }"
            />
            <div class="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent"></div>
            <div class="absolute bottom-4 left-4 right-4">
              <div class="flex items-center gap-2">
                <ElTag
                  :type="model.isFree ? 'success' : 'primary'"
                  effect="dark"
                  size="large"
                >
                  {{ model.isFree ? '免费' : `¥${model.price}` }}
                </ElTag>
                <ElTag type="info" effect="plain">{{ model.category }}</ElTag>
              </div>
            </div>
          </div>

          <div class="mt-4">
            <h1 class="text-2xl font-bold text-gray-900">{{ model.title }}</h1>
            <p class="mt-2 text-gray-600">{{ model.description }}</p>

            <div class="mt-4 flex items-center justify-between">
              <div class="flex items-center gap-3">
                <ElAvatar
                  :size="48"
                  :src="model.author.avatar || '/images/default-avatar.svg'"
                  @error="(e) => { const img = e.target as HTMLImageElement; if (img) img.src = '/images/default-avatar.svg'; }"
                />
                <div>
                  <div class="font-medium">{{ model.author.name }}</div>
                  <div class="text-sm text-gray-500">{{ model.author.bio || '资深思维模型专家' }}</div>
                </div>
              </div>
              <div class="text-sm text-gray-400">
                更新于 {{ model.updatedAt }}
              </div>
            </div>

            <div class="mt-4 flex flex-wrap gap-2">
              <ElTag
                v-for="tag in model.tags"
                :key="tag"
                type="primary"
                effect="light"
              >
                {{ tag }}
              </ElTag>
            </div>

            <!-- 统计数据 -->
            <ElDivider />
            <div class="grid grid-cols-5 gap-4 text-center">
              <div>
                <div class="text-lg font-bold">{{ formatNumber(model.stats.adoptions) }}</div>
                <div class="text-xs text-gray-500">采纳</div>
              </div>
              <div>
                <div class="text-lg font-bold">{{ formatNumber(model.stats.practices) }}</div>
                <div class="text-xs text-gray-500">练习</div>
              </div>
              <div>
                <div class="text-lg font-bold">{{ formatNumber(model.stats.discussions) }}</div>
                <div class="text-xs text-gray-500">讨论</div>
              </div>
              <div>
                <div class="text-lg font-bold">{{ formatNumber(model.stats.forks) }}</div>
                <div class="text-xs text-gray-500">引用</div>
              </div>
              <div>
                <div class="text-lg font-bold text-purple-600">{{ formatNumber(model.stats.likes) }}</div>
                <div class="text-xs text-gray-500">点赞</div>
              </div>
            </div>
          </div>
        </ElCard>

        <!-- Tabs -->
        <ElCard shadow="never">
          <ElTabs v-model="activeTab" type="border-card">
            <!-- 使用指南 Tab -->
            <ElTabPane label="使用指南" name="guide">
              <div class="prose max-w-none">
                <div
                  v-if="model.content"
                  class="text-gray-700 leading-relaxed"
                  v-html="formattedContent"
                ></div>
                <div v-else class="text-gray-500">
                  <h2 class="text-xl font-bold mt-6 mb-3">什么是{{ model.title }}？</h2>
                  <p class="mb-4">{{ model.description }}</p>
                  <h2 class="text-xl font-bold mt-6 mb-3">如何使用</h2>
                  <ol class="list-decimal list-inside space-y-2">
                    <li>理解模型的核心概念和原理</li>
                    <li>阅读示例，学习如何应用</li>
                    <li>在实际问题中尝试使用</li>
                    <li>记录你的思考过程</li>
                    <li>与他人分享和讨论</li>
                  </ol>
                  <h2 class="text-xl font-bold mt-6 mb-3">应用场景</h2>
                  <ul class="list-disc list-inside space-y-2">
                    <li>商业决策</li>
                    <li>产品规划</li>
                    <li>问题解决</li>
                    <li>创新思考</li>
                  </ul>
                </div>
              </div>
            </ElTabPane>

            <!-- 讨论 Tab -->
            <ElTabPane label="讨论" name="discussion">
              <div class="space-y-6">
                <!-- 发表评论 -->
                <div class="rounded-lg bg-gray-50 p-4">
                  <h4 class="mb-3 font-medium">参与讨论</h4>
                  <ElInput
                    v-model="newComment"
                    type="textarea"
                    :rows="3"
                    placeholder="分享你的想法、疑问或经验..."
                    class="mb-3"
                  />
                  <div class="flex items-center justify-between">
                    <ElButton type="primary" @click="handleSubmitComment">
                      发布评论
                    </ElButton>
                    <span class="text-xs text-gray-500">{{ comments.length }} 条讨论</span>
                  </div>
                </div>

                <!-- 评论列表 -->
                <div v-if="comments.length > 0" class="space-y-4">
                  <div
                    v-for="comment in comments"
                    :key="comment.id"
                    class="rounded-lg border border-gray-100 p-4"
                  >
                    <div class="flex items-start gap-3">
                      <ElAvatar
                        :size="40"
                        :src="comment.author.avatar || '/images/default-avatar.svg'"
                        @error="(e) => { const img = e.target as HTMLImageElement; if (img) img.src = '/images/default-avatar.svg'; }"
                      />
                      <div class="flex-1">
                        <div class="flex items-center gap-2 mb-1">
                          <span class="font-medium">{{ comment.author.name }}</span>
                          <span class="text-xs text-gray-400">{{ comment.createdAt }}</span>
                        </div>
                        <p class="text-sm text-gray-700">{{ comment.content }}</p>
                        <div class="mt-2 flex items-center gap-4 text-xs text-gray-500">
                          <span class="cursor-pointer hover:text-purple-600">👍 {{ comment.likes }}</span>
                          <span class="cursor-pointer hover:text-purple-600">回复</span>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
                <ElEmpty v-else description="暂无讨论，来发表第一条评论吧！" />
              </div>
            </ElTabPane>

            <!-- 练习记录 Tab -->
            <ElTabPane label="练习记录" name="practices">
              <div class="space-y-6">
                <!-- 练习统计 -->
                <div class="grid grid-cols-3 gap-4">
                  <div class="rounded-lg bg-purple-50 p-4 text-center">
                    <div class="text-2xl font-bold text-purple-600">{{ formatNumber(model.stats.practices) }}</div>
                    <div class="text-sm text-gray-600">总练习次数</div>
                  </div>
                  <div class="rounded-lg bg-blue-50 p-4 text-center">
                    <div class="text-2xl font-bold text-blue-600">{{ practiceRecords.length }}</div>
                    <div class="text-sm text-gray-600">公开练习</div>
                  </div>
                  <div class="rounded-lg bg-green-50 p-4 text-center">
                    <div class="text-2xl font-bold text-green-600">{{ practiceRecords.filter(p => p.isExcellent).length }}</div>
                    <div class="text-sm text-gray-600">优秀练习</div>
                  </div>
                </div>

                <!-- 练习列表 -->
                <div v-if="practiceRecords.length > 0" class="space-y-4">
                  <div
                    v-for="record in practiceRecords"
                    :key="record.id"
                    class="rounded-lg border border-gray-100 p-4 transition-all hover:border-purple-200 hover:shadow-sm"
                  >
                    <div class="flex items-start gap-4">
                      <ElAvatar
                        :size="48"
                        :src="record.user.avatar || '/images/default-avatar.svg'"
                        @error="(e) => { const img = e.target as HTMLImageElement; if (img) img.src = '/images/default-avatar.svg'; }"
                      />
                      <div class="flex-1">
                        <div class="flex items-center justify-between">
                          <div class="flex items-center gap-2">
                            <span class="font-medium">{{ record.user.name }}</span>
                            <ElTag v-if="record.isExcellent" type="success" size="small">优秀</ElTag>
                            <ElTag v-if="record.user.isCertified" type="warning" size="small">认证用户</ElTag>
                          </div>
                          <span class="text-xs text-gray-400">{{ record.createdAt }}</span>
                        </div>

                        <div class="mt-2">
                          <h4 class="font-medium text-gray-900">{{ record.topicTitle }}</h4>
                          <p class="mt-1 line-clamp-3 text-sm text-gray-600">{{ record.summary }}</p>
                        </div>

                        <div class="mt-3 flex items-center gap-4 text-xs text-gray-500">
                          <span class="flex items-center gap-1">
                            <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                            </svg>
                            {{ formatNumber(record.views) }} 浏览
                          </span>
                          <span class="flex items-center gap-1">
                            <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"/>
                            </svg>
                            {{ formatNumber(record.likes) }} 点赞
                          </span>
                          <span class="flex items-center gap-1">
                            <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/>
                            </svg>
                            {{ formatNumber(record.comments) }} 评论
                          </span>
                        </div>

                        <div class="mt-3">
                          <ElButton link type="primary" size="small" @click="viewPracticeDetail(record.id)">
                            查看详情 →
                          </ElButton>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                <ElEmpty v-else description="暂无练习记录，成为第一个练习者吧！">
                  <template #extra>
                    <ElButton type="primary" @click="goToCreateTopic">开始练习</ElButton>
                  </template>
                </ElEmpty>
              </div>
            </ElTabPane>

            <!-- 版本历史 Tab -->
            <ElTabPane label="版本历史" name="versions">
              <div class="space-y-4">
                <div class="flex items-center justify-between rounded-lg border border-gray-100 p-4">
                  <div>
                    <div class="font-medium">v2.0 当前版本</div>
                    <div class="text-sm text-gray-500">新增更多实战案例，优化使用说明</div>
                  </div>
                  <span class="text-xs text-gray-400">{{ model.updatedAt }}</span>
                </div>
                <div class="flex items-center justify-between rounded-lg border border-gray-100 p-4 bg-gray-50">
                  <div>
                    <div class="font-medium">v1.0</div>
                    <div class="text-sm text-gray-500">初始版本发布</div>
                  </div>
                  <span class="text-xs text-gray-400">2023-12-15</span>
                </div>
              </div>
            </ElTabPane>
          </ElTabs>
        </ElCard>
      </div>

      <!-- 右侧：操作和推荐 -->
      <div class="space-y-6">
        <!-- 操作按钮 -->
        <ElCard shadow="never">
          <div class="space-y-3">
            <ElButton
              v-if="model.isFree"
              type="primary"
              size="large"
              class="w-full"
              @click="handleLoad"
            >
              <svg class="h-5 w-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
              </svg>
              加载到我的模型
            </ElButton>
            <ElButton
              v-else
              type="success"
              size="large"
              class="w-full"
              @click="handlePurchase"
            >
              <svg class="h-5 w-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z"/>
              </svg>
              购买 ¥{{ model.price }}
            </ElButton>
            <ElButton
              size="large"
              class="w-full"
              @click="handleFork"
            >
              <svg class="h-5 w-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7v8a2 2 0 002 2h6M8 7V5a2 2 0 012-2h4.586a1 1 0 01.707.293l4.414 4.414a1 1 0 01.293.707V15a2 2 0 01-2 2h-2M8 7H6a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2v-2"/>
              </svg>
              引用创建副本
            </ElButton>
            <ElButton
              size="large"
              class="w-full"
              @click="handleLike"
            >
              👍 点赞 ({{ formatNumber(model.stats.likes) }})
            </ElButton>
            <ElDivider />
            <ElButton
              type="warning"
              size="large"
              class="w-full"
              @click="goToCreateTopic"
            >
              <svg class="h-5 w-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"/>
              </svg>
              用此模型分析课题
            </ElButton>
            <ElButton
              size="large"
              class="w-full"
              @click="goBack"
            >
              ← 返回市场
            </ElButton>
          </div>
        </ElCard>

        <!-- 相关模型推荐 -->
        <ElCard v-if="relatedModels.length > 0" shadow="never" header="相关模型推荐">
          <div class="space-y-4">
            <div
              v-for="related in relatedModels"
              :key="related.id"
              class="cursor-pointer rounded-lg border border-gray-100 p-3 transition-colors hover:border-purple-300"
              @click="goToRelatedModel(related.id)"
            >
              <div class="font-medium text-sm">{{ related.title }}</div>
              <div class="mt-1 flex items-center justify-between text-xs text-gray-500">
                <span>{{ related.category }}</span>
                <span>{{ formatNumber(related.stats.adoptions) }} 采纳</span>
              </div>
            </div>
          </div>
        </ElCard>

        <!-- 快速导航 -->
        <ElCard shadow="never" header="快速导航">
          <div class="space-y-3 text-sm">
            <div class="flex items-center justify-between">
              <span class="text-gray-500">我的模型库</span>
              <ElButton link type="primary" @click="router.push('/my-topics')">查看 →</ElButton>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-gray-500">创建新课题</span>
              <ElButton link type="primary" @click="goToCreateTopic">创建 →</ElButton>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-gray-500">同类模型</span>
              <ElButton link type="primary" @click="router.push(`/market?category=${model.category}`)">查看 →</ElButton>
            </div>
          </div>
        </ElCard>
      </div>
    </div>

    <!-- 错误状态 -->
    <ElEmpty v-else description="模型不存在或已被删除" />
  </Page>
</template>

<style scoped>
:deep(.el-tabs__content) {
  padding: 20px 0;
}

.prose h2 {
  color: var(--el-text-color-primary);
}

.prose p {
  margin-bottom: 1rem;
}

.line-clamp-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
