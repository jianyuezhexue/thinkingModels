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
  ElEmpty,
  ElTag,
  ElAvatar,
} from 'element-plus';

import {
  getMeetupListApi,
  interestMeetupApi,
  getMeetupCitiesApi,
  type CollisionApi,
} from '#/api';

import { $t } from '#/locales';

// 加载状态
const loading = ref(false);

// 约见列表数据
const meetups = ref<CollisionApi.Meetup[]>([]);
const total = ref(0);

// 筛选状态
const searchQuery = ref('');
const selectedTopic = ref<CollisionApi.MeetupTopic | 'all'>('all');
const selectedCity = ref('');
const selectedCostSplit = ref<CollisionApi.CostSplit | 'all'>('all');
const selectedSort = ref<CollisionApi.MeetupListParams['sortBy']>('latest');

// 分页
const currentPage = ref(1);
const pageSize = ref(10);

// 可用城市
const cities = ref<string[]>([]);

// 主题分类
const topics = [
  { id: 'all', name: '全部主题', icon: '🎯' },
  { id: 'career', name: '职业发展', icon: '💼' },
  { id: 'startup', name: '创业交流', icon: '🚀' },
  { id: 'technology', name: '技术探讨', icon: '💻' },
  { id: 'investment', name: '投资理财', icon: '📈' },
  { id: 'life', name: '生活感悟', icon: '🌟' },
  { id: 'other', name: '其他话题', icon: '💭' },
];

// 费用选项
const costSplitOptions = [
  { id: 'all', name: '不限', description: '' },
  { id: 'host', name: '我请客', description: '发起人请客' },
  { id: 'aa', name: 'AA制', description: '各付各的' },
  { id: 'guest', name: '你请客', description: '参与者请客' },
];

// 排序选项
const sortOptions = [
  { id: 'latest', name: '最新发布' },
  { id: 'popular', name: '最受关注' },
  { id: 'soonest', name: '即将开始' },
];

const router = useRouter();

// 获取约见列表
async function fetchMeetupList() {
  loading.value = true;
  try {
    const params: CollisionApi.MeetupListParams = {
      page: currentPage.value,
      pageSize: pageSize.value,
      sortBy: selectedSort.value,
      keyword: searchQuery.value || undefined,
      topic: selectedTopic.value,
      city: selectedCity.value || undefined,
      costSplit: selectedCostSplit.value,
      status: 'open',
    };

    const res = await getMeetupListApi(params);
    meetups.value = res.list;
    total.value = res.total;
  } catch (error) {
    console.error('获取约见列表失败:', error);
    ElMessage.error('获取约见列表失败');
  } finally {
    loading.value = false;
  }
}

// 获取城市列表
async function fetchCities() {
  try {
    cities.value = await getMeetupCitiesApi();
  } catch (error) {
    console.error('获取城市列表失败:', error);
  }
}

// 主题切换
function onTopicChange(topicId: string) {
  selectedTopic.value = topicId as CollisionApi.MeetupTopic | 'all';
  currentPage.value = 1;
  fetchMeetupList();
}

// 搜索
function onSearch() {
  currentPage.value = 1;
  fetchMeetupList();
}

// 表示感兴趣
async function onInterest(meetup: CollisionApi.Meetup) {
  try {
    const res = await interestMeetupApi(meetup.id);
    meetup.isInterested = res.interested;
    meetup.interestedCount = res.interestedCount;
    ElMessage.success(res.interested ? '已标记感兴趣' : '已取消');
  } catch (error) {
    ElMessage.error('操作失败');
  }
}

// 查看详情
function viewDetail(meetup: CollisionApi.Meetup) {
  router.push(`/collision/meetup/${meetup.id}`);
}

// 发起约见
function createMeetup() {
  router.push('/collision/meetup/create');
}

// 费用标签样式
function getCostSplitStyle(costSplit: CollisionApi.CostSplit) {
  const styles: Record<CollisionApi.CostSplit, { bg: string; text: string; label: string }> = {
    host: { bg: 'bg-green-100', text: 'text-green-700', label: '☕ 我请客' },
    aa: { bg: 'bg-blue-100', text: 'text-blue-700', label: '🤝 AA制' },
    guest: { bg: 'bg-orange-100', text: 'text-orange-700', label: '🎁 你请客' },
  };
  return styles[costSplit];
}

// 状态标签样式
function getStatusStyle(status: CollisionApi.MeetupStatus) {
  const styles: Record<CollisionApi.MeetupStatus, { bg: string; text: string; label: string }> = {
    open: { bg: 'bg-emerald-100', text: 'text-emerald-700', label: '🟢 招募中' },
    pending: { bg: 'bg-amber-100', text: 'text-amber-700', label: '⏳ 待确认' },
    confirmed: { bg: 'bg-purple-100', text: 'text-purple-700', label: '✅ 已确认' },
    completed: { bg: 'bg-gray-100', text: 'text-gray-700', label: '✔️ 已完成' },
    cancelled: { bg: 'bg-red-100', text: 'text-red-700', label: '❌ 已取消' },
  };
  return styles[status];
}

// 格式化时间
function formatTime(dateStr: string) {
  const date = new Date(dateStr);
  const now = new Date();
  const diff = now.getTime() - date.getTime();
  const days = Math.floor(diff / (1000 * 60 * 60 * 24));
  
  if (days === 0) return '今天';
  if (days === 1) return '昨天';
  if (days < 7) return `${days}天前`;
  return date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' });
}

// 监听筛选变化
watch([selectedCity, selectedCostSplit, selectedSort], () => {
  currentPage.value = 1;
  fetchMeetupList();
});

onMounted(() => {
  fetchMeetupList();
  fetchCities();
});
</script>

<template>
  <Page
    :title="$t('page.collision.meetup.title')"
    description="发布话题和思考，邀约志同道合的人线下交流"
  >
    <div class="flex gap-6">
      <!-- 左侧主内容 -->
      <div class="flex-1 min-w-0">
        <!-- 顶部操作栏 -->
        <div class="mb-6 flex items-center justify-between gap-4 flex-wrap">
          <ElInput
            v-model="searchQuery"
            placeholder="搜索约见话题..."
            class="!w-80"
            clearable
            @keyup.enter="onSearch"
          >
            <template #prefix>
              <span class="text-gray-400">🔍</span>
            </template>
          </ElInput>
          
          <ElButton type="primary" @click="createMeetup" class="!bg-purple-600 hover:!bg-purple-700">
            <span class="mr-2">☕</span>
            发起约见
          </ElButton>
        </div>

        <!-- 主题筛选 -->
        <div class="mb-6 flex items-center gap-2 flex-wrap">
          <button
            v-for="topic in topics"
            :key="topic.id"
            class="px-4 py-2 rounded-full text-sm font-medium transition-all"
            :class="[
              selectedTopic === topic.id
                ? 'bg-purple-600 text-white shadow-md'
                : 'bg-white text-gray-600 hover:bg-purple-50 hover:text-purple-600 border'
            ]"
            @click="onTopicChange(topic.id)"
          >
            <span class="mr-1">{{ topic.icon }}</span>
            {{ topic.name }}
          </button>
        </div>

        <!-- 二级筛选 -->
        <div class="mb-6 flex items-center gap-4 flex-wrap">
          <ElSelect v-model="selectedCity" placeholder="选择城市" clearable class="!w-32">
            <ElOption label="全部城市" value="" />
            <ElOption v-for="city in cities" :key="city" :label="city" :value="city" />
          </ElSelect>

          <ElSelect v-model="selectedCostSplit" placeholder="费用方式" class="!w-32">
            <ElOption
              v-for="opt in costSplitOptions"
              :key="opt.id"
              :label="opt.name"
              :value="opt.id"
            />
          </ElSelect>

          <ElSelect v-model="selectedSort" class="!w-32">
            <ElOption
              v-for="opt in sortOptions"
              :key="opt.id"
              :label="opt.name"
              :value="opt.id"
            />
          </ElSelect>
        </div>

        <!-- 加载骨架屏 -->
        <div v-if="loading" class="space-y-4">
          <ElCard v-for="i in 3" :key="i" shadow="hover" class="!rounded-xl">
            <ElSkeleton :rows="4" animated />
          </ElCard>
        </div>

        <!-- 空状态 -->
        <ElEmpty v-else-if="meetups.length === 0" description="暂无约见，来发起第一个吧！" />

        <!-- 约见列表 -->
        <div v-else class="space-y-4">
          <ElCard
            v-for="meetup in meetups"
            :key="meetup.id"
            shadow="hover"
            class="!rounded-xl cursor-pointer transition-all hover:shadow-lg"
            @click="viewDetail(meetup)"
          >
            <div class="flex gap-4">
              <!-- 发起人头像 -->
              <div class="flex-shrink-0">
                <ElAvatar :src="meetup.host.avatar" :size="56" class="ring-2 ring-purple-200" />
              </div>

              <!-- 内容区 -->
              <div class="flex-1 min-w-0">
                <!-- 标题行 -->
                <div class="flex items-start justify-between gap-4 mb-2">
                  <h3 class="text-lg font-semibold text-gray-800 line-clamp-1">
                    {{ meetup.title }}
                  </h3>
                  <div class="flex items-center gap-2 flex-shrink-0">
                    <span 
                      :class="[getStatusStyle(meetup.status).bg, getStatusStyle(meetup.status).text]"
                      class="px-2 py-0.5 rounded-full text-xs font-medium"
                    >
                      {{ getStatusStyle(meetup.status).label }}
                    </span>
                  </div>
                </div>

                <!-- 发起人信息 -->
                <div class="flex items-center gap-3 mb-3 text-sm text-gray-500">
                  <span class="font-medium text-gray-700">{{ meetup.host.name }}</span>
                  <span v-if="meetup.host.bio" class="text-gray-400">{{ meetup.host.bio }}</span>
                </div>

                <!-- 描述 -->
                <p class="text-gray-600 mb-3 line-clamp-2">
                  {{ meetup.description }}
                </p>

                <!-- 标签 -->
                <div class="flex items-center gap-2 mb-3 flex-wrap">
                  <span 
                    :class="[getCostSplitStyle(meetup.costSplit).bg, getCostSplitStyle(meetup.costSplit).text]"
                    class="px-2 py-0.5 rounded-full text-xs font-medium"
                  >
                    {{ getCostSplitStyle(meetup.costSplit).label }}
                  </span>
                  <ElTag v-if="meetup.modelName" type="info" size="small" round>
                    {{ meetup.modelName }}
                  </ElTag>
                  <ElTag v-for="tag in meetup.tags.slice(0, 3)" :key="tag" size="small" round>
                    {{ tag }}
                  </ElTag>
                </div>

                <!-- 底部信息 -->
                <div class="flex items-center justify-between text-sm text-gray-500">
                  <div class="flex items-center gap-4">
                    <span>📍 {{ meetup.city }}</span>
                    <span>🕐 {{ meetup.preferredTime }}</span>
                    <span>👥 {{ meetup.currentGuests }}/{{ meetup.maxGuests }} 人</span>
                  </div>
                  <div class="flex items-center gap-4">
                    <button
                      class="flex items-center gap-1 hover:text-purple-600 transition-colors"
                      :class="{ 'text-purple-600': meetup.isInterested }"
                      @click.stop="onInterest(meetup)"
                    >
                      <span>{{ meetup.isInterested ? '💜' : '🤍' }}</span>
                      {{ meetup.interestedCount }}
                    </button>
                    <span class="text-gray-400">{{ formatTime(meetup.createdAt) }}</span>
                  </div>
                </div>
              </div>
            </div>
          </ElCard>
        </div>
      </div>

      <!-- 右侧边栏 -->
      <div class="w-80 flex-shrink-0 space-y-6 hidden lg:block">
        <!-- 约见须知 -->
        <ElCard shadow="hover" class="!rounded-xl">
          <template #header>
            <span class="font-semibold text-gray-700">☕ 约见须知</span>
          </template>
          <ul class="text-sm text-gray-600 space-y-2">
            <li class="flex items-start gap-2">
              <span class="text-purple-500">•</span>
              选择公共场所（如咖啡厅）见面
            </li>
            <li class="flex items-start gap-2">
              <span class="text-purple-500">•</span>
              提前明确话题和交流预期
            </li>
            <li class="flex items-start gap-2">
              <span class="text-purple-500">•</span>
              尊重彼此时间，准时赴约
            </li>
            <li class="flex items-start gap-2">
              <span class="text-purple-500">•</span>
              保持真诚开放的交流态度
            </li>
            <li class="flex items-start gap-2">
              <span class="text-purple-500">•</span>
              注意个人信息和财产安全
            </li>
          </ul>
        </ElCard>

        <!-- 费用说明 -->
        <ElCard shadow="hover" class="!rounded-xl">
          <template #header>
            <span class="font-semibold text-gray-700">💰 费用说明</span>
          </template>
          <div class="space-y-3">
            <div class="flex items-center gap-3 p-2 rounded-lg bg-green-50">
              <span class="text-2xl">☕</span>
              <div>
                <div class="font-medium text-green-700">我请客</div>
                <div class="text-xs text-green-600">发起人承担费用</div>
              </div>
            </div>
            <div class="flex items-center gap-3 p-2 rounded-lg bg-blue-50">
              <span class="text-2xl">🤝</span>
              <div>
                <div class="font-medium text-blue-700">AA制</div>
                <div class="text-xs text-blue-600">各付各的，轻松交流</div>
              </div>
            </div>
            <div class="flex items-center gap-3 p-2 rounded-lg bg-orange-50">
              <span class="text-2xl">🎁</span>
              <div>
                <div class="font-medium text-orange-700">你请客</div>
                <div class="text-xs text-orange-600">参与者承担费用</div>
              </div>
            </div>
          </div>
        </ElCard>

        <!-- 热门城市 -->
        <ElCard shadow="hover" class="!rounded-xl">
          <template #header>
            <span class="font-semibold text-gray-700">🏙️ 热门城市</span>
          </template>
          <div class="flex flex-wrap gap-2">
            <button
              v-for="city in cities"
              :key="city"
              class="px-3 py-1 rounded-full text-sm transition-all"
              :class="[
                selectedCity === city
                  ? 'bg-purple-600 text-white'
                  : 'bg-gray-100 text-gray-600 hover:bg-purple-100 hover:text-purple-600'
              ]"
              @click="selectedCity = selectedCity === city ? '' : city"
            >
              {{ city }}
            </button>
          </div>
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
