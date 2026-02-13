<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { Page } from '@vben/common-ui';
import {
  ElCard,
  ElButton,
  ElInput,
  ElSelect,
  ElOption,
  ElEmpty,
  ElSkeleton,
  ElPagination,
  ElMessage,
  ElDialog,
  ElProgress,
} from 'element-plus';
import { 
  getMyThinkingTopicListApi, 
  getThinkingTopicStatisticsApi,
  type ThinkingTopicApi 
} from '#/api';

const router = useRouter();

// ==================== 状态管理 ====================
const loading = ref(false);
const topics = ref<(ThinkingTopicApi.TopicInfo & { 
  recommendedActions?: string[];
  selectedModels?: string[];
  progress?: number;
})[]>([]);
const total = ref(0);

// 分页
const currentPage = ref(1);
const pageSize = ref(10);

// 筛选
const searchQuery = ref('');
const activeStatus = ref<ThinkingTopicApi.TopicStatus | 'all'>('all');
const sortBy = ref<'latest' | 'updated' | 'analysis'>('latest');

// 弹窗
const actionDialogVisible = ref(false);
const selectedTopic = ref<(ThinkingTopicApi.TopicInfo & { recommendedActions?: string[] }) | null>(null);
const completedActionsMap = ref<Record<string, boolean>>({});

// 统计数据 - 使用数字状态
const stats = computed(() => {
  return {
    total: total.value,
    inProgress: topics.value.filter(t => t.status === 0).length, // 0=进行中
    completed: topics.value.filter(t => t.status === 1).length,   // 1=已完成
    draft: topics.value.filter(t => t.status === 3).length,       // 3=草稿
  };
});

// 状态选项 - 使用数字状态
const statusTabs = [
  { id: 'all' as const, label: '全部课题', icon: '📋' },
  { id: 0, label: '进行中', icon: '⏳' },
  { id: 1, label: '已完成', icon: '✅' },
  { id: 3, label: '草稿', icon: '📝' },
  { id: 2, label: '已归档', icon: '📦' },
];

// 模型名称
const modelNames = ['SWOT分析', '5W1H', 'MECE原则', '第一性原理', '金字塔原理', '逆向思维', '奥卡姆剃刀', '二阶思维'];

// ==================== 数据获取 ====================
async function fetchTopics() {
  loading.value = true;
  try {
    const params: ThinkingTopicApi.TopicListParams = {
      page: currentPage.value,
      pageSize: pageSize.value,
      title: searchQuery.value || undefined,
    };
    if (activeStatus.value !== 'all') {
      params.status = activeStatus.value as ThinkingTopicApi.TopicStatus;
    }

    const res = await getMyThinkingTopicListApi(params);
    
    // 添加模拟数据
    topics.value = res.list.map((topic: ThinkingTopicApi.TopicInfo, index: number) => ({
      ...topic,
      recommendedActions: index % 3 === 0 ? [
        '重新评估目标用户群体，缩小范围至核心用户',
        '制定3个月内可执行的MVP功能清单',
        '寻找2-3位潜在用户进行深度访谈',
      ] : index % 3 === 1 ? [
        '整理现有数据，建立分析框架',
      ] : undefined,
      selectedModels: topic.modelName ? [topic.modelName] : modelNames.slice(index % 4, index % 4 + 1 + (index % 3)),
      progress: getProgressValue(topic.status),
    }));
    total.value = res.total;
  } catch (error) {
    console.error('获取课题列表失败:', error);
    ElMessage.error('获取课题列表失败');
  } finally {
    loading.value = false;
  }
}

// ==================== 工具函数 ====================
function getProgressValue(status: ThinkingTopicApi.TopicStatus): number {
  const map: Record<number, number> = {
    3: 15,   // 草稿
    0: 60,   // 进行中
    1: 100,  // 已完成
    2: 100,  // 已归档
  };
  return map[status] || 0;
}

function getStatusStyle(status: ThinkingTopicApi.TopicStatus): string {
  const styles: Record<number, string> = {
    3: 'bg-gray-100 text-gray-600',      // 草稿
    0: 'bg-amber-100 text-amber-700',    // 进行中
    1: 'bg-green-100 text-green-700',    // 已完成
    2: 'bg-slate-100 text-slate-600',    // 已归档
  };
  return styles[status] || 'bg-gray-100 text-gray-600';
}

function getStatusText(status: ThinkingTopicApi.TopicStatus): string {
  const texts: Record<number, string> = {
    3: '草稿',
    0: '进行中',
    1: '已完成',
    2: '已归档',
  };
  return texts[status] || '未知';
}

function getStatusIcon(status: ThinkingTopicApi.TopicStatus): string {
  const icons: Record<number, string> = {
    3: '📝',
    0: '⏳',
    1: '✅',
    2: '📦',
  };
  return icons[status] || '📋';
}

function formatTime(dateStr: string): string {
  if (!dateStr) return '';
  const date = new Date(dateStr);
  const now = new Date();
  const diff = now.getTime() - date.getTime();
  const days = Math.floor(diff / (1000 * 60 * 60 * 24));
  
  if (days === 0) return '今天';
  if (days === 1) return '昨天';
  if (days < 7) return days + '天前';
  if (days < 30) return Math.floor(days / 7) + '周前';
  return date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' });
}

function formatFullDate(dateStr: string): string {
  if (!dateStr) return '';
  return new Date(dateStr).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
}

// ==================== 导航和操作 ====================
function goToCreate() {
  router.push('/my-topics/create');
}

function goToDetail(topic: ThinkingTopicApi.TopicInfo) {
  router.push('/my-topics/' + topic.id);
}

function startAnalysis(topic: ThinkingTopicApi.TopicInfo) {
  router.push('/my-topics/' + topic.id + '?tab=analysis');
}

function viewActions(topic: ThinkingTopicApi.TopicInfo & { recommendedActions?: string[] }) {
  selectedTopic.value = topic;
  actionDialogVisible.value = true;
}

function toggleAction(topicId: number, actionIndex: number) {
  const key = topicId + '-' + actionIndex;
  completedActionsMap.value[key] = !completedActionsMap.value[key];
}

// ==================== 监听器 ====================
watch([activeStatus, searchQuery, sortBy], () => {
  currentPage.value = 1;
  fetchTopics();
});

watch([currentPage, pageSize], () => {
  fetchTopics();
});

onMounted(() => {
  fetchTopics();
});
</script>

<template>
  <Page
    title="我的课题"
    description="运用思维模型深入剖析问题本质，导向正确的行动决策"
    content-class="p-6 bg-gray-50"
  >
    <!-- 顶部统计横幅 -->
    <div class="mb-6 grid grid-cols-1 md:grid-cols-4 gap-4">
      <div class="col-span-1 md:col-span-3">
        <div class="grid grid-cols-4 gap-4">
          <!-- 统计卡片 -->
          <div class="bg-white rounded-xl p-4 border border-gray-100 shadow-sm hover:shadow-md transition-shadow">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-lg bg-purple-100 flex items-center justify-center">
                <span class="text-xl">📊</span>
              </div>
              <div>
                <div class="text-2xl font-bold text-gray-800">{{ stats.total }}</div>
                <div class="text-xs text-gray-500">全部课题</div>
              </div>
            </div>
          </div>
          <div class="bg-white rounded-xl p-4 border border-gray-100 shadow-sm hover:shadow-md transition-shadow">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-lg bg-amber-100 flex items-center justify-center">
                <span class="text-xl">⏳</span>
              </div>
              <div>
                <div class="text-2xl font-bold text-amber-600">{{ stats.inProgress }}</div>
                <div class="text-xs text-gray-500">进行中</div>
              </div>
            </div>
          </div>
          <div class="bg-white rounded-xl p-4 border border-gray-100 shadow-sm hover:shadow-md transition-shadow">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-lg bg-green-100 flex items-center justify-center">
                <span class="text-xl">✅</span>
              </div>
              <div>
                <div class="text-2xl font-bold text-green-600">{{ stats.completed }}</div>
                <div class="text-xs text-gray-500">已完成</div>
              </div>
            </div>
          </div>
          <div class="bg-white rounded-xl p-4 border border-gray-100 shadow-sm hover:shadow-md transition-shadow">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-lg bg-gray-100 flex items-center justify-center">
                <span class="text-xl">📝</span>
              </div>
              <div>
                <div class="text-2xl font-bold text-gray-600">{{ stats.draft }}</div>
                <div class="text-xs text-gray-500">草稿</div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <!-- 快速创建 -->
      <div class="bg-gradient-to-br from-purple-500 to-purple-700 rounded-xl p-4 text-white flex flex-col justify-center items-center cursor-pointer hover:shadow-lg transition-shadow" @click="goToCreate">
        <div class="text-3xl mb-2">💡</div>
        <div class="font-semibold">创建新课题</div>
        <div class="text-xs text-purple-200 mt-1">开始深度思考</div>
      </div>
    </div>

    <!-- 主内容区 -->
    <div class="flex gap-6">
      <!-- 左侧课题列表 -->
      <div class="flex-1 space-y-6">
        <!-- 筛选和搜索 -->
        <ElCard shadow="hover" class="!rounded-xl">
          <div class="flex flex-wrap items-center gap-4">
            <!-- 状态 Tab -->
            <div class="flex gap-2">
              <button
                v-for="tab in statusTabs"
                :key="String(tab.id)"
                class="px-4 py-2 rounded-full text-sm font-medium transition-all"
                :class="[
                  activeStatus === tab.id
                    ? 'bg-purple-600 text-white shadow-md'
                    : 'bg-gray-100 text-gray-600 hover:bg-purple-100 hover:text-purple-600'
                ]"
                @click="activeStatus = tab.id as ThinkingTopicApi.TopicStatus | 'all'"
              >
                {{ tab.icon }} {{ tab.label }}
              </button>
            </div>
            <div class="flex-1" />
            <!-- 搜索框 -->
            <div class="flex items-center gap-3">
              <ElInput
                v-model="searchQuery"
                placeholder="搜索课题..."
                clearable
                class="!w-48"
              />
              <ElSelect v-model="sortBy" class="!w-32">
                <ElOption label="最新创建" value="latest" />
                <ElOption label="最近更新" value="updated" />
                <ElOption label="分析最多" value="analysis" />
              </ElSelect>
            </div>
          </div>
        </ElCard>

        <!-- 加载状态 -->
        <div v-if="loading" class="space-y-4">
          <ElSkeleton v-for="i in 3" :key="i" :rows="3" animated class="bg-white rounded-xl p-4" />
        </div>

        <!-- 空状态 -->
        <ElCard v-else-if="topics.length === 0" shadow="hover" class="!rounded-xl">
          <ElEmpty description="还没有课题，创建一个开始思考吧！">
            <template #image>
              <div class="text-6xl">💭</div>
            </template>
            <ElButton
              type="primary"
              class="!bg-purple-600 !border-purple-600 hover:!bg-purple-700 !rounded-full mt-4"
              @click="goToCreate"
            >
              创建第一个课题
            </ElButton>
          </ElEmpty>
        </ElCard>

        <!-- 课题列表 -->
        <div v-else class="space-y-4">
          <ElCard
            v-for="topic in topics"
            :key="topic.id"
            shadow="hover"
            class="!rounded-xl cursor-pointer hover:shadow-lg transition-all group"
            :class="{ '!border-l-4 !border-l-green-500': topic.recommendedActions?.length }"
            @click="goToDetail(topic)"
          >
            <div class="flex gap-4">
              <!-- 左侧进度环 -->
              <div class="flex-shrink-0 w-16 h-16 relative">
                <ElProgress
                  type="circle"
                  :percentage="topic.progress || 0"
                  :width="64"
                  :stroke-width="4"
                  :color="topic.status === 1 ? '#10b981' : topic.status === 0 ? '#f59e0b' : '#9ca3af'"
                >
                  <template #default>
                    <span class="text-lg">{{ getStatusIcon(topic.status) }}</span>
                  </template>
                </ElProgress>
              </div>

              <!-- 中间内容 -->
              <div class="flex-1 min-w-0">
                <div class="flex items-start justify-between gap-4 mb-2">
                  <div class="flex items-center gap-3">
                    <h3 class="text-lg font-semibold text-gray-800 group-hover:text-purple-600 transition-colors line-clamp-1">
                      {{ topic.title }}
                    </h3>
                    <span
                      class="px-2 py-0.5 rounded-full text-xs flex-shrink-0"
                      :class="getStatusStyle(topic.status)"
                    >
                      {{ getStatusText(topic.status) }}
                    </span>
                  </div>
                  <div class="flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                    <ElButton
                      v-if="topic.status === 3"
                      type="primary"
                      size="small"
                      class="!bg-purple-600 !border-purple-600 !rounded-full"
                      @click.stop="startAnalysis(topic)"
                    >
                      开始分析
                    </ElButton>
                    <ElButton
                      v-if="topic.recommendedActions?.length"
                      type="success"
                      size="small"
                      class="!rounded-full"
                      @click.stop="viewActions(topic)"
                    >
                      查看行动
                    </ElButton>
                  </div>
                </div>

                <p class="text-gray-600 text-sm line-clamp-2 mb-3">
                  {{ topic.description }}
                </p>

                <!-- 行动建议预览 -->
                <div
                  v-if="topic.recommendedActions?.length"
                  class="mb-3 p-2 rounded-lg bg-green-50 border border-green-100"
                >
                  <div class="flex items-center gap-2 text-sm">
                    <span class="text-green-600 font-medium">行动建议:</span>
                    <span class="text-gray-700 line-clamp-1">{{ topic.recommendedActions[0] }}</span>
                    <span v-if="topic.recommendedActions.length > 1" class="text-green-600 text-xs font-medium flex-shrink-0">
                      +{{ topic.recommendedActions.length - 1 }}
                    </span>
                  </div>
                </div>

                <!-- 底部元信息 -->
                <div class="flex flex-wrap items-center gap-4 text-xs text-gray-400">
                  <span class="flex items-center gap-1">
                    <span v-if="topic.selectedModels?.length" class="text-purple-600 font-medium">
                      {{ topic.selectedModels.slice(0, 2).join(', ') }}
                      <template v-if="topic.selectedModels.length > 2">等</template>
                    </span>
                    <span v-else>未选用模型</span>
                  </span>
                  <span class="text-gray-300">|</span>
                  <span>{{ topic.analysisCount }} 次分析</span>
                  <span class="text-gray-300">|</span>
                  <span>{{ formatTime(topic.createdAt) }}</span>
                </div>
              </div>
            </div>
          </ElCard>

          <!-- 分页 -->
          <div class="flex justify-center pt-4">
            <ElPagination
              v-model:current-page="currentPage"
              :page-size="pageSize"
              :total="total"
              layout="prev, pager, next"
              background
            />
          </div>
        </div>
      </div>

      <!-- 右侧边栏 -->
      <div class="w-80 flex-shrink-0 space-y-6 hidden lg:block">
        <!-- 使用指南 -->
        <ElCard shadow="hover" class="!rounded-xl">
          <template #header>
            <span class="font-semibold text-gray-700">使用指南</span>
          </template>
          <div class="space-y-4">
            <div class="flex items-start gap-3">
              <div class="w-8 h-8 rounded-full bg-purple-100 text-purple-600 flex items-center justify-center font-bold text-sm flex-shrink-0">1</div>
              <div>
                <div class="font-medium text-gray-700 text-sm">创建课题</div>
                <div class="text-xs text-gray-500">描述你要深入思考的问题</div>
              </div>
            </div>
            <div class="flex items-start gap-3">
              <div class="w-8 h-8 rounded-full bg-purple-100 text-purple-600 flex items-center justify-center font-bold text-sm flex-shrink-0">2</div>
              <div>
                <div class="font-medium text-gray-700 text-sm">选用模型</div>
                <div class="text-xs text-gray-500">挑选适合的思维模型来分析</div>
              </div>
            </div>
            <div class="flex items-start gap-3">
              <div class="w-8 h-8 rounded-full bg-purple-100 text-purple-600 flex items-center justify-center font-bold text-sm flex-shrink-0">3</div>
              <div>
                <div class="font-medium text-gray-700 text-sm">开始分析</div>
                <div class="text-xs text-gray-500">按照模型框架进行深度思考</div>
              </div>
            </div>
            <div class="flex items-start gap-3">
              <div class="w-8 h-8 rounded-full bg-green-100 text-green-600 flex items-center justify-center font-bold text-sm flex-shrink-0">4</div>
              <div>
                <div class="font-medium text-gray-700 text-sm">获取行动</div>
                <div class="text-xs text-gray-500">得到可执行的行动建议</div>
              </div>
            </div>
          </div>
        </ElCard>

        <!-- 热门模型推荐 -->
        <ElCard shadow="hover" class="!rounded-xl">
          <template #header>
            <span class="font-semibold text-gray-700">热门模型</span>
          </template>
          <div class="flex flex-wrap gap-2">
            <span
              v-for="model in modelNames.slice(0, 6)"
              :key="model"
              class="px-3 py-1 rounded-full text-sm bg-gray-100 text-gray-600 hover:bg-purple-100 hover:text-purple-600 cursor-pointer transition-colors"
            >
              {{ model }}
            </span>
          </div>
        </ElCard>

        <!-- 思考小贴士 -->
        <ElCard shadow="hover" class="!rounded-xl !bg-gradient-to-br from-purple-50 to-purple-100 !border-purple-200">
          <template #header>
            <span class="font-semibold text-purple-700">思考小贴士</span>
          </template>
          <ul class="text-sm text-purple-800 space-y-2">
            <li class="flex items-start gap-2">
              <span class="text-purple-500">•</span>
              问题越具体，分析越有效
            </li>
            <li class="flex items-start gap-2">
              <span class="text-purple-500">•</span>
              尝试多个模型看不同角度
            </li>
            <li class="flex items-start gap-2">
              <span class="text-purple-500">•</span>
              行动建议需要及时执行
            </li>
            <li class="flex items-start gap-2">
              <span class="text-purple-500">•</span>
              定期回顾已完成的课题
            </li>
          </ul>
        </ElCard>

        <!-- 最近活动 -->
        <ElCard shadow="hover" class="!rounded-xl">
          <template #header>
            <span class="font-semibold text-gray-700">最近活动</span>
          </template>
          <div class="space-y-3">
            <div
              v-for="topic in topics.slice(0, 3)"
              :key="'recent-' + topic.id"
              class="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-50 cursor-pointer transition-colors"
              @click="goToDetail(topic)"
            >
              <span class="text-lg">{{ getStatusIcon(topic.status) }}</span>
              <div class="flex-1 min-w-0">
                <div class="text-sm text-gray-700 line-clamp-1">{{ topic.title }}</div>
                <div class="text-xs text-gray-400">{{ formatTime(topic.updatedAt) }}</div>
              </div>
            </div>
          </div>
        </ElCard>
      </div>
    </div>

    <!-- 行动建议弹窗 -->
    <ElDialog
      v-model="actionDialogVisible"
      title="行动建议"
      width="600px"
    >
      <div v-if="selectedTopic" class="space-y-4">
        <div class="p-4 bg-gray-50 rounded-lg">
          <h3 class="font-semibold text-gray-800 mb-1">{{ selectedTopic.title }}</h3>
          <p class="text-sm text-gray-500">{{ formatFullDate(selectedTopic.createdAt) }}</p>
        </div>
        
        <div class="space-y-3">
          <div
            v-for="(action, index) in selectedTopic.recommendedActions"
            :key="index"
            class="flex items-start gap-3 p-4 rounded-lg transition-colors cursor-pointer"
            :class="completedActionsMap[selectedTopic.id + '-' + index] ? 'bg-green-50' : 'bg-gray-50 hover:bg-gray-100'"
            @click="toggleAction(selectedTopic.id, index)"
          >
            <div
              class="w-6 h-6 rounded-full flex items-center justify-center flex-shrink-0"
              :class="completedActionsMap[selectedTopic.id + '-' + index] ? 'bg-green-500 text-white' : 'bg-gray-200'"
            >
              <span v-if="completedActionsMap[selectedTopic.id + '-' + index]" class="text-xs">✓</span>
              <span v-else class="text-xs text-gray-400">{{ index + 1 }}</span>
            </div>
            <span
              class="text-sm"
              :class="completedActionsMap[selectedTopic.id + '-' + index] ? 'text-green-700 line-through' : 'text-gray-700'"
            >
              {{ action }}
            </span>
          </div>
        </div>
      </div>
      <template #footer>
        <ElButton @click="actionDialogVisible = false">关闭</ElButton>
        <ElButton type="primary" class="!bg-purple-600 !border-purple-600" @click="actionDialogVisible = false">
          导出清单
        </ElButton>
      </template>
    </ElDialog>
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
