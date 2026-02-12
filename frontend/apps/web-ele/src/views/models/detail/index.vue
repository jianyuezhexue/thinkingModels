<script lang="ts" setup>
import { computed, ref } from 'vue';
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
} from 'element-plus';

// Types
interface ModelAuthor {
  id: string;
  name: string;
  avatar: string;
  bio?: string;
}

interface Comment {
  id: string;
  author: ModelAuthor;
  content: string;
  createdAt: string;
  likes: number;
  replies?: Comment[];
}

interface PracticeRecord {
  id: string;
  title: string;
  content: string;
  createdAt: string;
  isPublic: boolean;
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
  content?: string; // 模型详细内容/使用指南
}

const route = useRoute();
const router = useRouter();
const modelId = computed(() => route.params.id as string);

// 当前激活的Tab
const activeTab = ref('practice');

// 从市场页面传递过来的模型数据（实际项目中可以从store或API获取）
const model = ref<ThinkingModel>({
  id: modelId.value,
  title: '第一性原理思维',
  description: '像马斯克一样回归本质，打破常规的创新思考方式',
  cover: 'https://images.unsplash.com/photo-1507413245164-6160d8298b31?w=800&h=400&fit=crop',
  author: {
    id: 'u3',
    name: '王创新',
    avatar: 'https://avatar.vercel.sh/wangcx.svg?text=WC',
    bio: '资深产品创新专家，前字节跳动产品经理',
  },
  isFree: true,
  category: 'innovation',
  tags: ['创新思维', '底层逻辑', '马斯克', '第一性原理'],
  stats: { adoptions: 15230, practices: 48200, discussions: 1234, forks: 3456, likes: 10234 },
  updatedAt: '2024-01-20',
  content: `
## 什么是第一性原理思维？

第一性原理思维是一种回归事物本质的思考方式，由亚里士多德提出，被埃隆·马斯克广泛应用。

### 核心步骤：
1. **识别并质疑现有假设** - 打破常规认知
2. **拆解问题到基本要素** - 找到最基本的真理
3. **从基础重新构建解决方案** - 基于本质创造新方案

### 应用场景：
- 产品创新
- 商业模式设计
- 技术突破
- 个人成长
  `,
});

// 练习相关
const newPractice = ref({
  title: '',
  content: '',
  isPublic: true,
});
const myPractices = ref<PracticeRecord[]>([
  {
    id: 'p1',
    title: '电池成本优化分析',
    content: '运用第一性原理分析电动车电池成本：原材料（钴、锂、镍）的市场价格是多少？加工成本如何？通过重新设计电池结构和供应链，可以降低成本60%...',
    createdAt: '2024-02-10',
    isPublic: true,
  },
]);

// 讨论相关
const newComment = ref('');
const comments = ref<Comment[]>([
  {
    id: 'c1',
    author: { id: 'u10', name: '李思考', avatar: 'https://avatar.vercel.sh/lisk.svg?text=LS' },
    content: '这个模型在实际工作中非常有用，特别是在做产品规划的时候。建议大家多练习！',
    createdAt: '2024-02-15 14:30',
    likes: 23,
    replies: [
      {
        id: 'c1-1',
        author: { id: 'u11', name: '张思维', avatar: 'https://avatar.vercel.sh/zhangsw.svg?text=ZS' },
        content: '同意！我在设计新产品时用这个方法，确实能找到差异化的突破口。',
        createdAt: '2024-02-15 15:20',
        likes: 8,
      },
    ],
  },
  {
    id: 'c2',
    author: { id: 'u12', name: '赵分析', avatar: 'https://avatar.vercel.sh/zhaofx.svg?text=ZF' },
    content: '有没有人可以分享一下如何在团队会议中引导大家使用第一性原理？',
    createdAt: '2024-02-14 09:15',
    likes: 15,
  },
]);

// 相关模型推荐
const relatedModels = ref([
  { id: '1', title: 'SWOT 分析模型', category: 'strategy', adoptions: 12580 },
  { id: '5', title: '决策矩阵', category: 'decision', adoptions: 4560 },
  { id: '6', title: '六顶思考帽', category: 'creative', adoptions: 7230 },
]);

// Actions
function formatNumber(num: number): string {
  if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M';
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K';
  return num.toString();
}

function handleSubmitPractice() {
  if (!newPractice.value.title.trim() || !newPractice.value.content.trim()) {
    ElMessage.warning('请填写标题和内容');
    return;
  }
  myPractices.value.unshift({
    id: Date.now().toString(),
    title: newPractice.value.title,
    content: newPractice.value.content,
    createdAt: new Date().toISOString().split('T')[0] || '',
    isPublic: newPractice.value.isPublic,
  });
  newPractice.value = { title: '', content: '', isPublic: true };
  ElMessage.success('练习记录已保存');
}

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

function handleLoad() {
  ElMessage.success('已加载到我的模型');
}

function handlePurchase() {
  ElMessage.info('跳转到购买页面...');
}

function handleFork() {
  ElMessage.success('已创建副本');
}

function handleLike() {
  model.value.stats.likes++;
  ElMessage.success('已点赞');
}

function goToRelatedModel(id: string) {
  router.push(`/models/detail/${id}`);
}

function goBack() {
  router.push('/models');
}
</script>

<template>
  <Page
    :description="model.description"
    :title="model.title"
  >
    <div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
      <!-- 左侧：模型信息 -->
      <div class="lg:col-span-2 space-y-6">
        <!-- 封面和基本信息 -->
        <ElCard shadow="never">
          <div class="relative h-64 w-full overflow-hidden rounded-lg">
            <img :src="model.cover" class="h-full w-full object-cover" />
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
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-3">
                <ElAvatar :size="48" :src="model.author.avatar" />
                <div>
                  <div class="font-medium">{{ model.author.name }}</div>
                  <div class="text-sm text-gray-500">{{ model.author.bio }}</div>
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

        <!-- 练习和讨论 Tabs -->
        <ElCard shadow="never">
          <ElTabs v-model="activeTab" type="border-card">
            <!-- 练习 Tab -->
            <ElTabPane label="练习" name="practice">
              <div class="space-y-6">
                <!-- 新建练习 -->
                <div class="rounded-lg bg-gray-50 p-4">
                  <h4 class="mb-3 font-medium">记录新练习</h4>
                  <ElInput
                    v-model="newPractice.title"
                    placeholder="练习标题"
                    class="mb-3"
                  />
                  <ElInput
                    v-model="newPractice.content"
                    type="textarea"
                    :rows="4"
                    placeholder="描述你使用这个思维模型的过程、思考和收获..."
                    class="mb-3"
                  />
                  <div class="flex items-center justify-between">
                    <ElButton type="primary" @click="handleSubmitPractice">
                      保存练习
                    </ElButton>
                  </div>
                </div>

                <!-- 练习列表 -->
                <div v-if="myPractices.length > 0" class="space-y-4">
                  <h4 class="font-medium">我的练习记录</h4>
                  <div
                    v-for="practice in myPractices"
                    :key="practice.id"
                    class="rounded-lg border border-gray-100 p-4"
                  >
                    <div class="flex items-center justify-between mb-2">
                      <h5 class="font-medium">{{ practice.title }}</h5>
                      <span class="text-xs text-gray-400">{{ practice.createdAt }}</span>
                    </div>
                    <p class="text-sm text-gray-600">{{ practice.content }}</p>
                  </div>
                </div>
                <ElEmpty v-else description="暂无练习记录，开始你的第一次练习吧！" />
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
                  <ElButton type="primary" @click="handleSubmitComment">
                    发布评论
                  </ElButton>
                </div>

                <!-- 评论列表 -->
                <div v-if="comments.length > 0" class="space-y-4">
                  <h4 class="font-medium">全部讨论 ({{ comments.length }})</h4>
                  <div
                    v-for="comment in comments"
                    :key="comment.id"
                    class="rounded-lg border border-gray-100 p-4"
                  >
                    <div class="flex items-start gap-3">
                      <ElAvatar :size="40" :src="comment.author.avatar" />
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

                        <!-- 回复列表 -->
                        <div v-if="comment.replies?.length" class="mt-3 space-y-3">
                          <div
                            v-for="reply in comment.replies"
                            :key="reply.id"
                            class="flex items-start gap-2 rounded bg-gray-50 p-3"
                          >
                            <ElAvatar :size="32" :src="reply.author.avatar" />
                            <div class="flex-1">
                              <div class="flex items-center gap-2 mb-1">
                                <span class="font-medium text-sm">{{ reply.author.name }}</span>
                                <span class="text-xs text-gray-400">{{ reply.createdAt }}</span>
                              </div>
                              <p class="text-sm text-gray-600">{{ reply.content }}</p>
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
                <ElEmpty v-else description="暂无讨论，来发表第一条评论吧！" />
              </div>
            </ElTabPane>

            <!-- 使用指南 Tab -->
            <ElTabPane label="使用指南" name="guide">
              <div class="prose max-w-none">
                <div v-html="model.content?.replace(/\n/g, '<br>').replace(/## (.*)/g, '<h2>$1</h2>').replace(/### (.*)/g, '<h3>$1</h3>').replace(/\d\. \*\*(.*)\*\*/g, '<strong>$1</strong>')"></div>
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
                  <span class="text-xs text-gray-400">2024-01-20</span>
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

      <!-- 右侧：操作和相关推荐 -->
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
              加载到我的模型
            </ElButton>
            <ElButton
              v-else
              type="success"
              size="large"
              class="w-full"
              @click="handlePurchase"
            >
              购买 ¥{{ model.price }}
            </ElButton>
            <ElButton
              size="large"
              class="w-full"
              @click="handleFork"
            >
              引用创建副本
            </ElButton>
            <ElButton
              size="large"
              class="w-full"
              @click="handleLike"
            >
              👍 点赞 ({{ formatNumber(model.stats.likes) }})
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
        <ElCard shadow="never" header="相关模型推荐">
          <div class="space-y-4">
            <div
              v-for="related in relatedModels"
              :key="related.id"
              class="cursor-pointer rounded-lg border border-gray-100 p-3 transition-colors hover:border-purple-300"
              @click="goToRelatedModel(related.id)"
            >
              <div class="font-medium">{{ related.title }}</div>
              <div class="mt-1 flex items-center justify-between text-xs text-gray-500">
                <span>{{ related.category }}</span>
                <span>{{ formatNumber(related.adoptions) }} 采纳</span>
              </div>
            </div>
          </div>
        </ElCard>

        <!-- 快速导航 -->
        <ElCard shadow="never" header="快速导航">
          <div class="space-y-2 text-sm">
            <div class="flex items-center justify-between">
              <span class="text-gray-500">作者主页</span>
              <ElButton link type="primary">查看 →</ElButton>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-gray-500">同类模型</span>
              <ElButton link type="primary">查看 →</ElButton>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-gray-500">热门讨论</span>
              <ElButton link type="primary">查看 →</ElButton>
            </div>
          </div>
        </ElCard>
      </div>
    </div>
  </Page>
</template>

<style scoped>
:deep(.el-tabs__content) {
  padding: 20px 0;
}
</style>
