<script lang="ts" setup>
import { computed, onMounted, ref, watch } from 'vue';
import { useRouter, useRoute } from 'vue-router';

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

} from 'element-plus';

import {
  getDiscussionListApi,
  likeDiscussionApi,
  favoriteDiscussionApi,
  getHotTagsApi,
  getActiveUsersApi,
  getMeetupListApi,
  interestMeetupApi,
  getMeetupCitiesApi,
  getConsultationListApi,
  getRecommendedExpertsApi,
  type CollisionApi,
} from '#/api';

import { $t } from '#/locales';

const route = useRoute();

// 当前激活的 Tab
const activeTab = ref<'meetup' | 'discussion' | 'consultation'>('meetup');

// 加载状态
const loading = ref(false);

// ==================== 话题讨论相关 ====================
// 话题列表数据
const discussions = ref<CollisionApi.Discussion[]>([]);
const total = ref(0);

// 筛选状态
const searchQuery = ref('');
const selectedCategory = ref<CollisionApi.Category | 'all'>('all');
const selectedSort = ref<CollisionApi.DiscussionListParams['sortBy']>('latest');

// 分页
const currentPage = ref(1);
const pageSize = ref(10);

// 热门标签
const hotTags = ref<string[]>([]);

// ==================== 找人聊聊相关 ====================
// 约见列表数据
const meetups = ref<CollisionApi.Meetup[]>([]);
const meetupTotal = ref(0);

// 约见筛选状态
const meetupSearchQuery = ref('');
const selectedTopic = ref<CollisionApi.MeetupTopic | 'all'>('all');
const selectedCity = ref('');
const selectedCostSplit = ref<CollisionApi.CostSplit | 'all'>('all');
const meetupSort = ref<CollisionApi.MeetupListParams['sortBy']>('latest');

// 约见分页
const meetupPage = ref(1);

// 可用城市
const cities = ref<string[]>([]);

// 主题列表
const meetupTopics = [
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
  { id: 'all', name: '不限' },
  { id: 'host', name: '我请客' },
  { id: 'aa', name: 'AA制' },
  { id: 'guest', name: '你请客' },
];

// ==================== 付费咨询相关 ====================
// 咨询列表数据
const consultations = ref<CollisionApi.Consultation[]>([]);
const consultationTotal = ref(0);

// 咨询筛选状态
const consultationSearchQuery = ref('');
const selectedField = ref<CollisionApi.ConsultationField | 'all'>('all');
const selectedMode = ref<CollisionApi.ConsultationMode | 'all'>('all');
const consultationSort = ref<CollisionApi.ConsultationListParams['sortBy']>('latest');

// 咨询分页
const consultationPage = ref(1);

// 推荐专家
const recommendedExperts = ref<CollisionApi.Expert[]>([]);

// 领域列表
const consultationFields = [
  { id: 'all', name: '全部领域', icon: '🎯' },
  { id: 'career', name: '职业规划', icon: '💼' },
  { id: 'startup', name: '创业咨询', icon: '🚀' },
  { id: 'technology', name: '技术架构', icon: '💻' },
  { id: 'product', name: '产品设计', icon: '📱' },
  { id: 'investment', name: '投资理财', icon: '📈' },
  { id: 'management', name: '团队管理', icon: '👥' },
  { id: 'psychology', name: '心理咨询', icon: '🧠' },
  { id: 'other', name: '其他领域', icon: '💭' },
];

// 咨询方式选项
const consultationModeOptions = [
  { id: 'all', name: '不限' },
  { id: 'online', name: '线上' },
  { id: 'offline', name: '线下' },
  { id: 'both', name: '均可' },
];

// 活跃用户
const activeUsers = ref<CollisionApi.UserInfo[]>([]);

// 分类列表
const categories = [
  { id: 'all', name: $t('page.collision.categories.all'), icon: 'lucide:layout-grid' },
  { id: 'inspiration', name: $t('page.collision.categories.inspiration'), icon: 'lucide:lightbulb' },
  { id: 'methodology', name: $t('page.collision.categories.methodology'), icon: 'lucide:compass' },
  { id: 'case', name: $t('page.collision.categories.case'), icon: 'lucide:file-text' },
  { id: 'question', name: $t('page.collision.categories.question'), icon: 'lucide:help-circle' },
  { id: 'share', name: $t('page.collision.categories.share'), icon: 'lucide:share-2' },
];

// 排序选项
const sortOptions = [
  { id: 'latest', name: $t('page.collision.sort.latest') },
  { id: 'popular', name: $t('page.collision.sort.popular') },
  { id: 'mostCommented', name: $t('page.collision.sort.mostCommented') },
  { id: 'mostLiked', name: $t('page.collision.sort.mostLiked') },
];

const router = useRouter();

// 获取话题列表
async function fetchDiscussionList() {
  loading.value = true;
  try {
    const params: CollisionApi.DiscussionListParams = {
      page: currentPage.value,
      pageSize: pageSize.value,
      sortBy: selectedSort.value,
      keyword: searchQuery.value || undefined,
      category: selectedCategory.value,
    };

    const res = await getDiscussionListApi(params);
    discussions.value = res.list;
    total.value = res.total;
  } catch (error) {
    console.error('获取话题列表失败:', error);
    ElMessage.error('获取话题列表失败');
  } finally {
    loading.value = false;
  }
}

// 获取约见列表
async function fetchMeetupList() {
  loading.value = true;
  try {
    const params: CollisionApi.MeetupListParams = {
      page: meetupPage.value,
      pageSize: pageSize.value,
      sortBy: meetupSort.value,
      keyword: meetupSearchQuery.value || undefined,
      topic: selectedTopic.value,
      city: selectedCity.value || undefined,
      costSplit: selectedCostSplit.value,
      status: 'open',
    };

    const res = await getMeetupListApi(params);
    meetups.value = res.list;
    meetupTotal.value = res.total;
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

// 获取热门标签和活跃用户
async function fetchSidebarData() {
  try {
    const [tags, users] = await Promise.all([
      getHotTagsApi(),
      getActiveUsersApi(),
    ]);
    hotTags.value = tags;
    activeUsers.value = users;
  } catch (error) {
    console.error('获取侧边栏数据失败:', error);
  }
}

// 获取咨询列表
async function fetchConsultationList() {
  loading.value = true;
  try {
    const params: CollisionApi.ConsultationListParams = {
      page: consultationPage.value,
      pageSize: pageSize.value,
      sortBy: consultationSort.value,
      keyword: consultationSearchQuery.value || undefined,
      field: selectedField.value,
      mode: selectedMode.value,
      status: 'open',
    };

    const res = await getConsultationListApi(params);
    consultations.value = res.list;
    consultationTotal.value = res.total;
  } catch (error) {
    console.error('获取咨询列表失败:', error);
    ElMessage.error('获取咨询列表失败');
  } finally {
    loading.value = false;
  }
}

// 获取推荐专家
async function fetchRecommendedExperts() {
  try {
    recommendedExperts.value = await getRecommendedExpertsApi();
  } catch (error) {
    console.error('获取推荐专家失败:', error);
  }
}

// 监听话题筛选条件变化
watch([searchQuery, selectedCategory, selectedSort], () => {
  currentPage.value = 1;
  if (activeTab.value === 'discussion') {
    fetchDiscussionList();
  }
});

// 监听约见筛选条件变化
watch([meetupSearchQuery, selectedTopic, selectedCity, selectedCostSplit, meetupSort], () => {
  meetupPage.value = 1;
  if (activeTab.value === 'meetup') {
    fetchMeetupList();
  }
});

// 监听咨询筛选条件变化
watch([consultationSearchQuery, selectedField, selectedMode, consultationSort], () => {
  consultationPage.value = 1;
  if (activeTab.value === 'consultation') {
    fetchConsultationList();
  }
});

// Tab 切换时加载对应数据
watch(activeTab, (tab) => {
  if (tab === 'discussion') {
    fetchDiscussionList();
    fetchSidebarData();
  } else if (tab === 'meetup') {
    fetchMeetupList();
    fetchCities();
  } else if (tab === 'consultation') {
    fetchConsultationList();
    fetchRecommendedExperts();
  }
});

// 页面加载时获取数据
onMounted(() => {
  // 检查 URL 参数来确定初始 Tab
  const tab = route.query.tab as string;
  if (tab === 'discussion') {
    activeTab.value = 'discussion';
    fetchDiscussionList();
    fetchSidebarData();
  } else if (tab === 'consultation') {
    activeTab.value = 'consultation';
    fetchConsultationList();
    fetchRecommendedExperts();
  } else {
    // 默认显示找人聊聊
    activeTab.value = 'meetup';
    fetchMeetupList();
    fetchCities();
  }
});

// 跳转到详情页
function goToDetail(discussion: CollisionApi.Discussion) {
  router.push(`/collision/${discussion.id}`);
}

// 跳转到约见详情页
function goToMeetupDetail(meetup: CollisionApi.Meetup) {
  router.push(`/collision/meetup/${meetup.id}`);
}

// 跳转到创建页
function goToCreate() {
  router.push('/collision/create');
}

// 跳转到创建约见页
function goToCreateMeetup() {
  router.push('/collision/meetup/create');
}

// 跳转到咨询详情页
function goToConsultationDetail(consultationId: string) {
  router.push(`/collision/consultation/${consultationId}`);
}

// 跳转到创建咨询页
function goToCreateConsultation() {
  router.push('/collision/consultation/create');
}

// 对约见表示感兴趣
async function onMeetupInterest(meetup: CollisionApi.Meetup, event: Event) {
  event.stopPropagation();
  try {
    const res = await interestMeetupApi(meetup.id);
    meetup.isInterested = res.interested;
    meetup.interestedCount = res.interestedCount;
    ElMessage.success(res.interested ? '已标记感兴趣' : '已取消');
  } catch (error) {
    ElMessage.error('操作失败');
  }
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
function getMeetupStatusStyle(status: CollisionApi.MeetupStatus) {
  const styles: Record<CollisionApi.MeetupStatus, { bg: string; text: string; label: string }> = {
    open: { bg: 'bg-emerald-100', text: 'text-emerald-700', label: '🟢 招募中' },
    pending: { bg: 'bg-amber-100', text: 'text-amber-700', label: '⏳ 待确认' },
    confirmed: { bg: 'bg-purple-100', text: 'text-purple-700', label: '✅ 已确认' },
    completed: { bg: 'bg-gray-100', text: 'text-gray-700', label: '✔️ 已完成' },
    cancelled: { bg: 'bg-red-100', text: 'text-red-700', label: '❌ 已取消' },
  };
  return styles[status];
}

// 咨询状态标签样式
function getConsultationStatusStyle(status: CollisionApi.ConsultationStatus) {
  const styles: Record<CollisionApi.ConsultationStatus, { bg: string; text: string; label: string }> = {
    open: { bg: 'bg-emerald-100', text: 'text-emerald-700', label: '🟢 招募中' },
    matched: { bg: 'bg-blue-100', text: 'text-blue-700', label: '🤝 已匹配' },
    inProgress: { bg: 'bg-purple-100', text: 'text-purple-700', label: '⏳ 进行中' },
    completed: { bg: 'bg-gray-100', text: 'text-gray-700', label: '✔️ 已完成' },
    cancelled: { bg: 'bg-red-100', text: 'text-red-700', label: '❌ 已取消' },
    expired: { bg: 'bg-gray-100', text: 'text-gray-500', label: '⏰ 已过期' },
  };
  return styles[status];
}

// 咨询方式标签样式
function getConsultationModeStyle(mode: CollisionApi.ConsultationMode) {
  const styles: Record<CollisionApi.ConsultationMode, { bg: string; text: string; label: string }> = {
    online: { bg: 'bg-blue-100', text: 'text-blue-700', label: '💻 线上' },
    offline: { bg: 'bg-orange-100', text: 'text-orange-700', label: '🤝 线下' },
    both: { bg: 'bg-purple-100', text: 'text-purple-700', label: '🔄 均可' },
  };
  return styles[mode];
}

// 获取领域名称
function getFieldName(field: CollisionApi.ConsultationField) {
  const fieldItem = consultationFields.find(f => f.id === field);
  return fieldItem?.name || field;
}

// 格式化悬赏金额
function formatReward(reward: number) {
  if (reward >= 1000) {
    return `¥${(reward / 1000).toFixed(reward % 1000 === 0 ? 0 : 1)}k`;
  }
  return `¥${reward}`;
}

// 格式化截止日期
function formatDeadline(deadline: string) {
  const date = new Date(deadline);
  const now = new Date();
  const diff = date.getTime() - now.getTime();
  const days = Math.ceil(diff / (1000 * 60 * 60 * 24));
  
  if (days < 0) return '已过期';
  if (days === 0) return '今天截止';
  if (days === 1) return '明天截止';
  if (days <= 7) return `${days}天后截止`;
  return date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' }) + ' 截止';
}

// 点赞
async function handleLike(discussion: CollisionApi.Discussion, event: Event) {
  event.stopPropagation();
  try {
    const res = await likeDiscussionApi(discussion.id);
    discussion.isLiked = res.liked;
    discussion.likeCount = res.likeCount;
  } catch (error) {
    console.error('点赞失败:', error);
  }
}

// 收藏
async function handleFavorite(discussion: CollisionApi.Discussion, event: Event) {
  event.stopPropagation();
  try {
    const res = await favoriteDiscussionApi(discussion.id);
    discussion.isFavorited = res.favorited;
    discussion.favoriteCount = res.favoriteCount;
  } catch (error) {
    console.error('收藏失败:', error);
  }
}

// 格式化数字
function formatNumber(num: number): string {
  if (num >= 10000) return (num / 10000).toFixed(1) + '万';
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K';
  return num.toString();
}

// 格式化时间
function formatTime(time: string): string {
  const now = new Date();
  const date = new Date(time);
  const diff = now.getTime() - date.getTime();
  const minutes = Math.floor(diff / 60000);
  const hours = Math.floor(diff / 3600000);
  const days = Math.floor(diff / 86400000);

  if (minutes < 1) return '刚刚';
  if (minutes < 60) return `${minutes}分钟前`;
  if (hours < 24) return `${hours}小时前`;
  if (days < 7) return `${days}天前`;
  return date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' });
}

// 获取分类颜色
function getCategoryColor(category: CollisionApi.Category): string {
  const colors: Record<CollisionApi.Category, string> = {
    inspiration: 'bg-amber-100 text-amber-700',
    methodology: 'bg-blue-100 text-blue-700',
    case: 'bg-green-100 text-green-700',
    question: 'bg-purple-100 text-purple-700',
    share: 'bg-pink-100 text-pink-700',
  };
  return colors[category] || 'bg-gray-100 text-gray-700';
}

// 获取分类名称
function getCategoryName(category: CollisionApi.Category): string {
  const cat = categories.find(c => c.id === category);
  return cat?.name || category;
}

// 计算总页数
const totalPages = computed(() => Math.ceil(total.value / pageSize.value));

// 分页显示的页码
const displayPages = computed(() => {
  const pages: number[] = [];
  const maxVisible = 5;
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
    :description="$t('page.collision.subtitle')"
    :title="$t('page.collision.title')"
  >
    <!-- 顶部 Tab 切换 -->
    <div class="mb-6">
      <div class="flex items-center gap-2 p-1 bg-gray-100 rounded-xl w-fit">
        <button
          class="px-6 py-2.5 rounded-lg text-sm font-medium transition-all"
          :class="[
            activeTab === 'meetup'
              ? 'bg-white text-purple-600 shadow-sm'
              : 'text-gray-600 hover:text-purple-600'
          ]"
          @click="activeTab = 'meetup'"
        >
          ☕ 找人聊聊
        </button>
        <button
          class="px-6 py-2.5 rounded-lg text-sm font-medium transition-all"
          :class="[
            activeTab === 'discussion'
              ? 'bg-white text-purple-600 shadow-sm'
              : 'text-gray-600 hover:text-purple-600'
          ]"
          @click="activeTab = 'discussion'"
        >
          💬 话题讨论
        </button>
        <button
          class="px-6 py-2.5 rounded-lg text-sm font-medium transition-all"
          :class="[
            activeTab === 'consultation'
              ? 'bg-white text-purple-600 shadow-sm'
              : 'text-gray-600 hover:text-purple-600'
          ]"
          @click="activeTab = 'consultation'"
        >
          💰 付费咨询
        </button>
      </div>
    </div>

    <!-- 话题讨论内容 -->
    <div v-show="activeTab === 'discussion'" class="flex gap-6">
      <!-- 主内容区 -->
      <div class="flex-1 min-w-0">
        <!-- 筛选和搜索栏 -->
        <ElCard class="mb-6" shadow="never">
          <div class="flex flex-col gap-4">
            <!-- 分类标签 -->
            <div class="flex flex-wrap items-center gap-2">
              <button
                v-for="cat in categories"
                :key="cat.id"
                :class="[
                  'px-4 py-2 rounded-full text-sm font-medium transition-all cursor-pointer',
                  selectedCategory === cat.id
                    ? 'bg-purple-600 text-white shadow-md'
                    : 'bg-gray-100 text-gray-600 hover:bg-purple-100 hover:text-purple-600'
                ]"
                @click="selectedCategory = cat.id as CollisionApi.Category | 'all'"
              >
                {{ cat.name }}
              </button>
            </div>

            <!-- 搜索和排序 -->
            <div class="flex items-center gap-4">
              <ElInput
                v-model="searchQuery"
                :placeholder="$t('page.collision.searchPlaceholder')"
                clearable
                class="flex-1"
                @keyup.enter="fetchDiscussionList"
              >
                <template #prefix>
                  <svg class="h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
                  </svg>
                </template>
              </ElInput>
              <ElSelect v-model="selectedSort" style="width: 140px">
                <ElOption
                  v-for="opt in sortOptions"
                  :key="opt.id"
                  :label="opt.name"
                  :value="opt.id"
                />
              </ElSelect>
              <ElButton type="primary" @click="goToCreate">
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
                </svg>
                {{ $t('page.collision.action.newTopic') }}
              </ElButton>
            </div>
          </div>
        </ElCard>

        <!-- 结果统计 -->
        <div class="mb-4 flex items-center justify-between text-sm text-gray-500">
          <span>共 {{ total }} 个话题</span>
          <span v-if="loading">加载中...</span>
        </div>

        <!-- 加载状态 -->
        <div v-if="loading" class="space-y-4">
          <ElCard v-for="i in 3" :key="i" shadow="never">
            <ElSkeleton animated>
              <template #template>
                <div class="flex gap-4">
                  <ElSkeletonItem variant="circle" style="width: 48px; height: 48px" />
                  <div class="flex-1 space-y-3">
                    <ElSkeletonItem variant="h3" style="width: 60%" />
                    <ElSkeletonItem variant="text" />
                    <ElSkeletonItem variant="text" style="width: 80%" />
                  </div>
                </div>
              </template>
            </ElSkeleton>
          </ElCard>
        </div>

        <!-- 空状态 -->
        <ElEmpty
          v-else-if="discussions.length === 0"
          :description="$t('page.collision.empty.description')"
        >
          <ElButton type="primary" @click="goToCreate">
            {{ $t('page.collision.action.newTopic') }}
          </ElButton>
        </ElEmpty>

        <!-- 话题列表 -->
        <div v-else class="space-y-4">
          <div
            v-for="discussion in discussions"
            :key="discussion.id"
            class="group bg-white rounded-xl border border-gray-100 p-5 cursor-pointer transition-all hover:shadow-lg hover:border-purple-200"
            @click="goToDetail(discussion)"
          >
            <!-- 顶部标记 -->
            <div v-if="discussion.isTop || discussion.isFeatured" class="flex gap-2 mb-3">
              <span v-if="discussion.isTop" class="px-2 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-600">
                {{ $t('page.collision.badge.top') }}
              </span>
              <span v-if="discussion.isFeatured" class="px-2 py-0.5 rounded-full text-xs font-medium bg-purple-100 text-purple-600">
                {{ $t('page.collision.badge.featured') }}
              </span>
            </div>

            <div class="flex gap-4">
              <!-- 用户头像 -->
              <ElAvatar :src="discussion.user.avatar" :size="48" class="flex-shrink-0" />

              <!-- 内容区 -->
              <div class="flex-1 min-w-0">
                <!-- 标题 -->
                <h3 class="text-lg font-semibold text-gray-900 group-hover:text-purple-600 transition-colors line-clamp-1">
                  {{ discussion.title }}
                </h3>

                <!-- 摘要 -->
                <p class="mt-2 text-sm text-gray-500 line-clamp-2">
                  {{ discussion.summary || discussion.content }}
                </p>

                <!-- 标签和模型 -->
                <div class="mt-3 flex flex-wrap items-center gap-2">
                  <span :class="['px-2 py-0.5 rounded text-xs font-medium', getCategoryColor(discussion.category)]">
                    {{ getCategoryName(discussion.category) }}
                  </span>
                  <span
                    v-if="discussion.modelName"
                    class="px-2 py-0.5 rounded text-xs font-medium bg-blue-50 text-blue-600"
                  >
                    📊 {{ discussion.modelName }}
                  </span>
                  <ElTag
                    v-for="tag in discussion.tags.slice(0, 3)"
                    :key="tag"
                    size="small"
                    type="info"
                    effect="plain"
                  >
                    {{ tag }}
                  </ElTag>
                </div>

                <!-- 底部信息栏 -->
                <div class="mt-4 flex items-center justify-between">
                  <!-- 作者和时间 -->
                  <div class="flex items-center gap-2 text-sm text-gray-400">
                    <span class="font-medium text-gray-600">{{ discussion.user.name }}</span>
                    <span>·</span>
                    <span>{{ formatTime(discussion.publishTime) }}</span>
                  </div>

                  <!-- 统计和操作 -->
                  <div class="flex items-center gap-4">
                    <!-- 浏览量 -->
                    <ElTooltip :content="$t('page.collision.stats.views')">
                      <span class="flex items-center gap-1 text-sm text-gray-400">
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                        </svg>
                        {{ formatNumber(discussion.viewCount) }}
                      </span>
                    </ElTooltip>

                    <!-- 评论 -->
                    <ElTooltip :content="$t('page.collision.stats.comments')">
                      <span class="flex items-center gap-1 text-sm text-gray-400">
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/>
                        </svg>
                        {{ formatNumber(discussion.commentCount) }}
                      </span>
                    </ElTooltip>

                    <!-- 点赞按钮 -->
                    <button
                      :class="[
                        'flex items-center gap-1 text-sm transition-colors',
                        discussion.isLiked ? 'text-red-500' : 'text-gray-400 hover:text-red-500'
                      ]"
                      @click="handleLike(discussion, $event)"
                    >
                      <svg class="w-4 h-4" :fill="discussion.isLiked ? 'currentColor' : 'none'" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"/>
                      </svg>
                      {{ formatNumber(discussion.likeCount) }}
                    </button>

                    <!-- 收藏按钮 -->
                    <button
                      :class="[
                        'flex items-center gap-1 text-sm transition-colors',
                        discussion.isFavorited ? 'text-yellow-500' : 'text-gray-400 hover:text-yellow-500'
                      ]"
                      @click="handleFavorite(discussion, $event)"
                    >
                      <svg class="w-4 h-4" :fill="discussion.isFavorited ? 'currentColor' : 'none'" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z"/>
                      </svg>
                      {{ formatNumber(discussion.favoriteCount) }}
                    </button>
                  </div>
                </div>
              </div>

              <!-- 封面图 -->
              <div
                v-if="discussion.cover"
                class="flex-shrink-0 w-32 h-24 rounded-lg overflow-hidden"
              >
                <img
                  :src="discussion.cover"
                  class="w-full h-full object-cover"
                  @error="(e) => { const img = e.target as HTMLImageElement; img.style.display = 'none'; }"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- 分页 -->
        <div v-if="totalPages > 1" class="mt-8 flex justify-center">
          <div class="flex items-center gap-2">
            <button
              :disabled="currentPage === 1"
              class="px-3 py-2 rounded-lg text-sm font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed hover:bg-purple-50 hover:text-purple-600"
              @click="currentPage--"
            >
              上一页
            </button>
            <button
              v-for="page in displayPages"
              :key="page"
              :class="[
                'w-10 h-10 rounded-lg text-sm font-medium transition-colors',
                currentPage === page
                  ? 'bg-purple-600 text-white'
                  : 'hover:bg-purple-50 hover:text-purple-600'
              ]"
              @click="currentPage = page"
            >
              {{ page }}
            </button>
            <button
              :disabled="currentPage === totalPages"
              class="px-3 py-2 rounded-lg text-sm font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed hover:bg-purple-50 hover:text-purple-600"
              @click="currentPage++"
            >
              下一页
            </button>
          </div>
        </div>
      </div>

      <!-- 右侧边栏 -->
      <div class="w-72 flex-shrink-0 space-y-6 hidden lg:block">
        <!-- 发起话题卡片 -->
        <ElCard shadow="never" class="bg-gradient-to-br from-purple-50 to-indigo-50">
          <div class="text-center py-4">
            <div class="w-16 h-16 mx-auto mb-4 bg-purple-100 rounded-full flex items-center justify-center">
              <svg class="w-8 h-8 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"/>
              </svg>
            </div>
            <h3 class="text-lg font-semibold text-gray-800 mb-2">有想法？来碰撞！</h3>
            <p class="text-sm text-gray-500 mb-4">分享你的灵感和思考，与更多人交流</p>
            <ElButton type="primary" class="w-full" @click="goToCreate">
              {{ $t('page.collision.action.newTopic') }}
            </ElButton>
          </div>
        </ElCard>

        <!-- 热门标签 -->
        <ElCard shadow="never">
          <template #header>
            <div class="flex items-center gap-2">
              <svg class="w-5 h-5 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"/>
              </svg>
              <span class="font-medium">热门标签</span>
            </div>
          </template>
          <div class="flex flex-wrap gap-2">
            <button
              v-for="tag in hotTags"
              :key="tag"
              class="px-3 py-1.5 rounded-full text-sm bg-gray-100 text-gray-600 hover:bg-purple-100 hover:text-purple-600 transition-colors cursor-pointer"
              @click="searchQuery = tag"
            >
              {{ tag }}
            </button>
          </div>
        </ElCard>

        <!-- 活跃用户 -->
        <ElCard shadow="never">
          <template #header>
            <div class="flex items-center gap-2">
              <svg class="w-5 h-5 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"/>
              </svg>
              <span class="font-medium">活跃用户</span>
            </div>
          </template>
          <div class="space-y-3">
            <div
              v-for="user in activeUsers"
              :key="user.id"
              class="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-50 transition-colors cursor-pointer"
            >
              <ElAvatar :src="user.avatar" :size="36" />
              <span class="text-sm font-medium text-gray-700">{{ user.name }}</span>
            </div>
          </div>
        </ElCard>

        <!-- 话题规则 -->
        <ElCard shadow="never">
          <template #header>
            <div class="flex items-center gap-2">
              <svg class="w-5 h-5 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
              <span class="font-medium">社区规则</span>
            </div>
          </template>
          <ul class="space-y-2 text-sm text-gray-500">
            <li class="flex items-start gap-2">
              <span class="text-purple-500">•</span>
              <span>尊重他人，理性讨论</span>
            </li>
            <li class="flex items-start gap-2">
              <span class="text-purple-500">•</span>
              <span>分享原创内容，禁止抄袭</span>
            </li>
            <li class="flex items-start gap-2">
              <span class="text-purple-500">•</span>
              <span>话题需与思维模型相关</span>
            </li>
            <li class="flex items-start gap-2">
              <span class="text-purple-500">•</span>
              <span>禁止发布广告和恶意内容</span>
            </li>
          </ul>
        </ElCard>
      </div>
    </div>

    <!-- 找人聊聊内容 -->
    <div v-show="activeTab === 'meetup'" class="flex gap-6">
      <!-- 主内容区 -->
      <div class="flex-1 min-w-0">
        <!-- 顶部操作栏 -->
        <div class="mb-6 flex items-center justify-between gap-4 flex-wrap">
          <ElInput
            v-model="meetupSearchQuery"
            placeholder="搜索约见话题..."
            class="!w-80"
            clearable
            @keyup.enter="fetchMeetupList"
          >
            <template #prefix>
              <span class="text-gray-400">🔍</span>
            </template>
          </ElInput>
          
          <ElButton type="primary" @click="goToCreateMeetup" class="!bg-purple-600 hover:!bg-purple-700">
            <span class="mr-2">☕</span>
            发起约见
          </ElButton>
        </div>

        <!-- 主题筛选 -->
        <div class="mb-6 flex items-center gap-2 flex-wrap">
          <button
            v-for="topic in meetupTopics"
            :key="topic.id"
            class="px-4 py-2 rounded-full text-sm font-medium transition-all"
            :class="[
              selectedTopic === topic.id
                ? 'bg-purple-600 text-white shadow-md'
                : 'bg-white text-gray-600 hover:bg-purple-50 hover:text-purple-600 border'
            ]"
            @click="selectedTopic = topic.id as CollisionApi.MeetupTopic | 'all'"
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

          <ElSelect v-model="meetupSort" class="!w-32">
            <ElOption label="最新发布" value="latest" />
            <ElOption label="最受关注" value="popular" />
            <ElOption label="即将开始" value="soonest" />
          </ElSelect>
        </div>

        <!-- 结果统计 -->
        <div class="mb-4 flex items-center justify-between text-sm text-gray-500">
          <span>共 {{ meetupTotal }} 个约见</span>
          <span v-if="loading">加载中...</span>
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
            @click="goToMeetupDetail(meetup)"
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
                      :class="[getMeetupStatusStyle(meetup.status).bg, getMeetupStatusStyle(meetup.status).text]"
                      class="px-2 py-0.5 rounded-full text-xs font-medium"
                    >
                      {{ getMeetupStatusStyle(meetup.status).label }}
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
                      @click.stop="onMeetupInterest(meetup, $event)"
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

    <!-- 付费咨询内容 -->
    <div v-show="activeTab === 'consultation'" class="flex gap-6">
      <!-- 左侧主内容 -->
      <div class="flex-1 space-y-6">
        <!-- 筛选和搜索 -->
        <ElCard shadow="hover" class="!rounded-xl">
          <div class="flex flex-wrap items-center gap-4">
            <div class="flex-1 min-w-64">
              <ElInput
                v-model="consultationSearchQuery"
                placeholder="搜索咨询话题..."
                prefix-icon="Search"
                clearable
                class="!rounded-full"
              />
            </div>
            <ElSelect
              v-model="selectedField"
              placeholder="咨询领域"
              clearable
              class="w-32"
            >
              <ElOption
                v-for="field in consultationFields"
                :key="field.id"
                :label="field.name"
                :value="field.id"
              />
            </ElSelect>
            <ElSelect
              v-model="selectedMode"
              placeholder="咨询方式"
              clearable
              class="w-32"
            >
              <ElOption
                v-for="mode in consultationModeOptions"
                :key="mode.id"
                :label="mode.name"
                :value="mode.id"
              />
            </ElSelect>
            <ElSelect v-model="consultationSort" class="w-32">
              <ElOption label="最新发布" value="newest" />
              <ElOption label="悬赏最高" value="reward" />
              <ElOption label="即将截止" value="deadline" />
              <ElOption label="申请最多" value="applications" />
            </ElSelect>
            <ElButton
              type="primary"
              class="!bg-purple-600 !border-purple-600 hover:!bg-purple-700 !rounded-full"
              @click="goToCreateConsultation"
            >
              💡 发布咨询
            </ElButton>
          </div>
        </ElCard>

        <!-- 咨询列表 -->
        <div class="space-y-4">
          <ElCard
            v-for="consultation in consultations"
            :key="consultation.id"
            shadow="hover"
            class="!rounded-xl cursor-pointer hover:shadow-lg transition-shadow"
            @click="goToConsultationDetail(consultation.id)"
          >
            <div class="flex gap-4">
              <!-- 悬赏金额 -->
              <div class="flex-shrink-0 w-20 h-20 rounded-xl bg-gradient-to-br from-orange-400 to-red-500 flex flex-col items-center justify-center text-white">
                <span class="text-xs">悬赏</span>
                <span class="text-xl font-bold">{{ formatReward(consultation.reward) }}</span>
              </div>
              
              <!-- 主要内容 -->
              <div class="flex-1 min-w-0">
                <div class="flex items-start justify-between gap-4 mb-2">
                  <h3 class="text-lg font-semibold text-gray-800 line-clamp-1">
                    {{ consultation.title }}
                  </h3>
                  <span
                    class="px-2 py-0.5 rounded-full text-xs flex-shrink-0"
                    :class="getConsultationStatusStyle(consultation.status)"
                  >
                    {{ consultation.status === 'open' ? '征集中' : consultation.status === 'matched' ? '已匹配' : consultation.status === 'inProgress' ? '进行中' : consultation.status === 'completed' ? '已完成' : consultation.status === 'cancelled' ? '已取消' : '已过期' }}
                  </span>
                </div>
                
                <p class="text-gray-600 text-sm line-clamp-2 mb-3">
                  {{ consultation.description }}
                </p>
                
                <div class="flex flex-wrap items-center gap-3 text-sm">
                  <span
                    class="px-2 py-0.5 rounded-full text-xs"
                    :class="getConsultationModeStyle(consultation.mode)"
                  >
                    {{ consultation.mode === 'online' ? '🖥️ 线上' : consultation.mode === 'offline' ? '☕ 线下' : '🔄 都可以' }}
                  </span>
                  <span class="px-2 py-0.5 rounded-full text-xs bg-gray-100 text-gray-600">
                    {{ getFieldName(consultation.field) }}
                  </span>
                  <span class="text-gray-400">|</span>
                  <span class="text-gray-500">
                    ⏰ 截止: {{ formatDeadline(consultation.deadline) }}
                  </span>
                  <span class="text-gray-400">|</span>
                  <span class="text-purple-600">
                    📝 {{ consultation.applicationCount || 0 }}人申请
                  </span>
                </div>
              </div>
            </div>
            
            <!-- 底部信息 -->
            <div class="flex items-center justify-between mt-4 pt-4 border-t border-gray-100">
              <div class="flex items-center gap-3">
                <ElAvatar :size="32" :src="consultation.requester?.avatar">
                  {{ consultation.requester?.name?.charAt(0) }}
                </ElAvatar>
                <div>
                  <span class="text-sm text-gray-700">{{ consultation.requester?.name }}</span>
                  <span class="text-xs text-gray-400 ml-2">发布于 {{ formatTime(consultation.createdAt) }}</span>
                </div>
              </div>
              <div class="flex items-center gap-2">
                <span class="text-sm text-gray-500">
                  👁️ {{ consultation.viewCount || 0 }}
                </span>
              </div>
            </div>
          </ElCard>
        </div>

        <!-- 分页 -->
        <div class="flex justify-center mt-6">
          <ElPagination
            v-model:current-page="currentPage"
            :page-size="pageSize"
            :total="consultationTotal"
            layout="prev, pager, next"
            background
          />
        </div>
      </div>

      <!-- 右侧边栏 -->
      <div class="w-80 flex-shrink-0 space-y-6 hidden lg:block">
        <!-- 推荐专家 -->
        <ElCard shadow="hover" class="!rounded-xl">
          <template #header>
            <span class="font-semibold text-gray-700">🌟 推荐专家</span>
          </template>
          <div class="space-y-4">
            <div
              v-for="expert in recommendedExperts"
              :key="expert.id"
              class="flex items-start gap-3 p-2 rounded-lg hover:bg-gray-50 cursor-pointer transition-colors"
            >
              <ElAvatar :size="48" :src="expert.avatar">
                {{ expert.name?.charAt(0) }}
              </ElAvatar>
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2">
                  <span class="font-medium text-gray-800">{{ expert.name }}</span>
                  <span class="text-xs text-purple-600 bg-purple-50 px-1.5 py-0.5 rounded">
                    {{ expert.title }}
                  </span>
                </div>
                <div class="text-xs text-gray-500 mt-0.5">{{ expert.company }}</div>
                <div class="flex items-center gap-2 mt-1 text-xs text-gray-400">
                  <span>⭐ {{ expert.rating?.toFixed(1) }}</span>
                  <span>|</span>
                  <span>📋 {{ expert.consultCount }}次咨询</span>
                </div>
                <div class="flex flex-wrap gap-1 mt-1">
                  <span
                    v-for="tag in expert.expertise?.slice(0, 2)"
                    :key="tag"
                    class="text-xs bg-gray-100 text-gray-600 px-1.5 py-0.5 rounded"
                  >
                    {{ tag }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </ElCard>

        <!-- 咨询须知 -->
        <ElCard shadow="hover" class="!rounded-xl">
          <template #header>
            <span class="font-semibold text-gray-700">📋 咨询须知</span>
          </template>
          <ul class="text-sm text-gray-600 space-y-2">
            <li class="flex items-start gap-2">
              <span class="text-purple-500">1.</span>
              清晰描述问题背景和期望
            </li>
            <li class="flex items-start gap-2">
              <span class="text-purple-500">2.</span>
              合理设置悬赏金额
            </li>
            <li class="flex items-start gap-2">
              <span class="text-purple-500">3.</span>
              选择合适的咨询方式
            </li>
            <li class="flex items-start gap-2">
              <span class="text-purple-500">4.</span>
              咨询完成后及时确认
            </li>
            <li class="flex items-start gap-2">
              <span class="text-purple-500">5.</span>
              对专家服务进行评价
            </li>
          </ul>
        </ElCard>

        <!-- 热门领域 -->
        <ElCard shadow="hover" class="!rounded-xl">
          <template #header>
            <span class="font-semibold text-gray-700">🔥 热门领域</span>
          </template>
          <div class="flex flex-wrap gap-2">
            <button
              v-for="field in consultationFields"
              :key="field.id"
              class="px-3 py-1 rounded-full text-sm transition-all"
              :class="[
                selectedField === field.id
                  ? 'bg-purple-600 text-white'
                  : 'bg-gray-100 text-gray-600 hover:bg-purple-100 hover:text-purple-600'
              ]"
              @click="selectedField = selectedField === field.id ? 'all' : (field.id as CollisionApi.ConsultationField | 'all')"
            >
              {{ field.name }}
            </button>
          </div>
        </ElCard>

        <!-- 悬赏排行 -->
        <ElCard shadow="hover" class="!rounded-xl">
          <template #header>
            <span class="font-semibold text-gray-700">💰 高额悬赏</span>
          </template>
          <div class="space-y-3">
            <div
              v-for="(consultation, index) in consultations.slice(0, 3)"
              :key="consultation.id"
              class="flex items-center gap-3 cursor-pointer hover:bg-gray-50 p-2 rounded-lg transition-colors"
              @click="goToConsultationDetail(consultation.id)"
            >
              <span
                class="w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold"
                :class="[
                  index === 0 ? 'bg-yellow-100 text-yellow-600' :
                  index === 1 ? 'bg-gray-100 text-gray-600' :
                  'bg-orange-50 text-orange-600'
                ]"
              >
                {{ index + 1 }}
              </span>
              <div class="flex-1 min-w-0">
                <div class="text-sm text-gray-700 line-clamp-1">{{ consultation.title }}</div>
              </div>
              <span class="text-sm font-semibold text-orange-500">
                ¥{{ consultation.reward }}
              </span>
            </div>
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
