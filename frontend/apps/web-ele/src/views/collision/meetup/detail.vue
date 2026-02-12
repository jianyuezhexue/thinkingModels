<script lang="ts" setup>
import { ref, onMounted, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { Page } from '@vben/common-ui';

import {
  ElButton,
  ElCard,
  ElMessage,
  ElSkeleton,
  ElEmpty,
  ElTag,
  ElAvatar,
  ElInput,
  ElDialog,
} from 'element-plus';

import {
  getMeetupDetailApi,
  interestMeetupApi,
  applyMeetupApi,
  type CollisionApi,
} from '#/api';

// locales can be used for i18n

const route = useRoute();
const router = useRouter();

// 加载状态
const loading = ref(true);

// 约见详情
const meetup = ref<CollisionApi.Meetup | null>(null);

// 申请列表（发起人可见）- 后续实现

// 申请对话框
const showApplyDialog = ref(false);
const applyMessage = ref('');
const applying = ref(false);

// 获取约见详情
async function fetchMeetupDetail() {
  const id = route.params.id as string;
  loading.value = true;
  try {
    meetup.value = await getMeetupDetailApi(id);
    if (!meetup.value) {
      ElMessage.error('约见不存在');
      router.push('/collision');
    }
  } catch (error) {
    console.error('获取约见详情失败:', error);
    ElMessage.error('获取约见详情失败');
  } finally {
    loading.value = false;
  }
}

// 费用标签样式
function getCostSplitStyle(costSplit: CollisionApi.CostSplit) {
  const styles: Record<CollisionApi.CostSplit, { bg: string; text: string; label: string; desc: string }> = {
    host: { bg: 'bg-green-100', text: 'text-green-700', label: '☕ 我请客', desc: '发起人请客' },
    aa: { bg: 'bg-blue-100', text: 'text-blue-700', label: '🤝 AA制', desc: '各付各的' },
    guest: { bg: 'bg-orange-100', text: 'text-orange-700', label: '🎁 你请客', desc: '参与者请客' },
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

// 主题名称
function getTopicName(topic: CollisionApi.MeetupTopic) {
  const names: Record<CollisionApi.MeetupTopic, string> = {
    career: '职业发展',
    startup: '创业交流',
    technology: '技术探讨',
    investment: '投资理财',
    life: '生活感悟',
    other: '其他话题',
  };
  return names[topic];
}

// 表示感兴趣
async function onInterest() {
  if (!meetup.value) return;
  try {
    const res = await interestMeetupApi(meetup.value.id);
    meetup.value.isInterested = res.interested;
    meetup.value.interestedCount = res.interestedCount;
    ElMessage.success(res.interested ? '已标记感兴趣' : '已取消');
  } catch (error) {
    ElMessage.error('操作失败');
  }
}

// 打开申请对话框
function openApplyDialog() {
  applyMessage.value = '';
  showApplyDialog.value = true;
}

// 提交申请
async function submitApply() {
  if (!meetup.value || !applyMessage.value.trim()) {
    ElMessage.warning('请填写申请留言');
    return;
  }

  applying.value = true;
  try {
    await applyMeetupApi({
      meetupId: meetup.value.id,
      message: applyMessage.value,
    });
    ElMessage.success('申请已提交，等待发起人确认');
    showApplyDialog.value = false;
  } catch (error) {
    ElMessage.error('申请失败，请重试');
  } finally {
    applying.value = false;
  }
}

// 分享
function onShare() {
  const url = window.location.href;
  navigator.clipboard.writeText(url).then(() => {
    ElMessage.success('链接已复制');
  }).catch(() => {
    ElMessage.info('请手动复制链接');
  });
}

// 返回列表
function goBack() {
  router.push('/collision');
}

// 格式化时间
function formatTime(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
}

// 是否可以申请
const canApply = computed(() => {
  if (!meetup.value) return false;
  return meetup.value.status === 'open' && meetup.value.currentGuests < meetup.value.maxGuests;
});

// 名额是否已满
const isFull = computed(() => {
  if (!meetup.value) return false;
  return meetup.value.currentGuests >= meetup.value.maxGuests;
});

onMounted(() => {
  fetchMeetupDetail();
});
</script>

<template>
  <Page :title="meetup?.title || '约见详情'">
    <!-- 返回按钮 -->
    <template #extra>
      <ElButton @click="goBack">← 返回列表</ElButton>
    </template>

    <!-- 加载骨架屏 -->
    <ElCard v-if="loading" shadow="hover" class="!rounded-xl">
      <ElSkeleton :rows="10" animated />
    </ElCard>

    <!-- 空状态 -->
    <ElEmpty v-else-if="!meetup" description="约见不存在或已删除" />

    <!-- 详情内容 -->
    <template v-else>
      <div class="flex gap-6">
        <!-- 左侧主内容 -->
        <div class="flex-1 min-w-0 space-y-6">
          <!-- 标题卡片 -->
          <ElCard shadow="hover" class="!rounded-xl">
            <!-- 状态标签 -->
            <div class="flex items-center justify-between mb-4">
              <span
                :class="[getStatusStyle(meetup.status).bg, getStatusStyle(meetup.status).text]"
                class="px-3 py-1 rounded-full text-sm font-medium"
              >
                {{ getStatusStyle(meetup.status).label }}
              </span>
              <span class="text-gray-400 text-sm">
                {{ formatTime(meetup.createdAt) }} 发布
              </span>
            </div>

            <!-- 标题 -->
            <h1 class="text-2xl font-bold text-gray-800 mb-4">
              {{ meetup.title }}
            </h1>

            <!-- 发起人信息 -->
            <div class="flex items-center gap-4 p-4 bg-gray-50 rounded-xl mb-4">
              <ElAvatar :src="meetup.host.avatar" :size="64" class="ring-2 ring-purple-200" />
              <div class="flex-1">
                <div class="font-semibold text-gray-800 text-lg">{{ meetup.host.name }}</div>
                <div v-if="meetup.host.bio" class="text-gray-500 text-sm mt-1">
                  {{ meetup.host.bio }}
                </div>
                <div v-if="meetup.host.interests?.length" class="flex items-center gap-2 mt-2">
                  <ElTag
                    v-for="interest in meetup.host.interests"
                    :key="interest"
                    size="small"
                    type="info"
                    round
                  >
                    {{ interest }}
                  </ElTag>
                </div>
              </div>
            </div>

            <!-- 基本信息 -->
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 mb-4">
              <div class="text-center p-3 bg-purple-50 rounded-xl">
                <div class="text-purple-600 text-xl mb-1">{{ getTopicName(meetup.topic) }}</div>
                <div class="text-gray-500 text-xs">话题类型</div>
              </div>
              <div class="text-center p-3 bg-blue-50 rounded-xl">
                <div class="text-blue-600 text-xl mb-1">📍 {{ meetup.city }}</div>
                <div class="text-gray-500 text-xs">约见城市</div>
              </div>
              <div class="text-center p-3 bg-amber-50 rounded-xl">
                <div class="text-amber-600 text-xl mb-1">🕐</div>
                <div class="text-gray-600 text-sm">{{ meetup.preferredTime }}</div>
                <div class="text-gray-500 text-xs">期望时间</div>
              </div>
              <div class="text-center p-3 bg-green-50 rounded-xl">
                <div class="text-green-600 text-xl mb-1">👥 {{ meetup.currentGuests }}/{{ meetup.maxGuests }}</div>
                <div class="text-gray-500 text-xs">已报名人数</div>
              </div>
            </div>

            <!-- 费用方式 -->
            <div
              :class="[getCostSplitStyle(meetup.costSplit).bg]"
              class="flex items-center justify-center gap-3 p-4 rounded-xl"
            >
              <span class="text-3xl">{{ getCostSplitStyle(meetup.costSplit).label.split(' ')[0] }}</span>
              <div>
                <div :class="getCostSplitStyle(meetup.costSplit).text" class="font-semibold">
                  {{ getCostSplitStyle(meetup.costSplit).label }}
                </div>
                <div class="text-sm opacity-75">{{ getCostSplitStyle(meetup.costSplit).desc }}</div>
              </div>
            </div>
          </ElCard>

          <!-- 描述 -->
          <ElCard shadow="hover" class="!rounded-xl">
            <template #header>
              <span class="font-semibold text-gray-700">📝 话题简介</span>
            </template>
            <p class="text-gray-700 leading-relaxed">
              {{ meetup.description }}
            </p>
          </ElCard>

          <!-- 发起人思考 -->
          <ElCard shadow="hover" class="!rounded-xl">
            <template #header>
              <div class="flex items-center gap-2">
                <span class="font-semibold text-gray-700">💭 发起人的思考</span>
                <span class="text-purple-600 text-xs bg-purple-100 px-2 py-0.5 rounded-full">核心内容</span>
              </div>
            </template>
            <div class="prose prose-sm max-w-none">
              <p
                v-for="(paragraph, index) in meetup.thoughts.split('\n')"
                :key="index"
                class="text-gray-700 leading-relaxed mb-3 whitespace-pre-wrap"
              >
                {{ paragraph }}
              </p>
            </div>
          </ElCard>

          <!-- 标签 -->
          <div v-if="meetup.tags.length > 0" class="flex items-center gap-2 flex-wrap">
            <ElTag v-if="meetup.modelName" type="primary" round>
              📚 {{ meetup.modelName }}
            </ElTag>
            <ElTag v-for="tag in meetup.tags" :key="tag" round>
              {{ tag }}
            </ElTag>
          </div>
        </div>

        <!-- 右侧边栏 -->
        <div class="w-80 flex-shrink-0 space-y-6 hidden lg:block">
          <!-- 操作卡片 -->
          <ElCard shadow="hover" class="!rounded-xl sticky top-4">
            <div class="space-y-4">
              <!-- 主要操作 -->
              <ElButton
                v-if="canApply"
                type="primary"
                class="!w-full !bg-purple-600 hover:!bg-purple-700 !h-12 !text-base"
                @click="openApplyDialog"
              >
                ☕ 我想参加
              </ElButton>
              <ElButton
                v-else-if="isFull"
                class="!w-full !h-12"
                disabled
              >
                👥 名额已满
              </ElButton>
              <ElButton
                v-else
                class="!w-full !h-12"
                disabled
              >
                {{ getStatusStyle(meetup.status).label }}
              </ElButton>

              <!-- 次要操作 -->
              <div class="flex gap-2">
                <ElButton
                  class="!flex-1"
                  :type="meetup.isInterested ? 'primary' : 'default'"
                  @click="onInterest"
                >
                  {{ meetup.isInterested ? '💜' : '🤍' }} {{ meetup.interestedCount }}
                </ElButton>
                <ElButton class="!flex-1" @click="onShare">
                  🔗 分享
                </ElButton>
              </div>

              <!-- 统计 -->
              <div class="flex items-center justify-around text-center pt-4 border-t">
                <div>
                  <div class="text-xl font-semibold text-gray-700">{{ meetup.viewCount }}</div>
                  <div class="text-xs text-gray-500">浏览</div>
                </div>
                <div>
                  <div class="text-xl font-semibold text-gray-700">{{ meetup.interestedCount }}</div>
                  <div class="text-xs text-gray-500">感兴趣</div>
                </div>
                <div>
                  <div class="text-xl font-semibold text-gray-700">{{ meetup.currentGuests }}</div>
                  <div class="text-xs text-gray-500">已报名</div>
                </div>
              </div>
            </div>
          </ElCard>

          <!-- 安全提示 -->
          <ElCard shadow="hover" class="!rounded-xl">
            <template #header>
              <span class="font-semibold text-gray-700">⚠️ 安全提示</span>
            </template>
            <ul class="text-sm text-gray-600 space-y-2">
              <li class="flex items-start gap-2">
                <span class="text-amber-500">•</span>
                选择公共场所见面
              </li>
              <li class="flex items-start gap-2">
                <span class="text-amber-500">•</span>
                告知朋友你的行程
              </li>
              <li class="flex items-start gap-2">
                <span class="text-amber-500">•</span>
                保护个人隐私信息
              </li>
              <li class="flex items-start gap-2">
                <span class="text-amber-500">•</span>
                注意财产安全
              </li>
            </ul>
          </ElCard>
        </div>
      </div>
    </template>

    <!-- 申请对话框 -->
    <ElDialog
      v-model="showApplyDialog"
      title="申请参加约见"
      width="500px"
    >
      <div class="space-y-4">
        <p class="text-gray-600">
          请简单介绍自己，让发起人了解你为什么想参加这次约见：
        </p>
        <ElInput
          v-model="applyMessage"
          type="textarea"
          :rows="4"
          placeholder="例如：我也在思考类似的问题，目前在尝试...，很想交流一下经验。"
          maxlength="500"
          show-word-limit
        />
      </div>
      <template #footer>
        <ElButton @click="showApplyDialog = false">取消</ElButton>
        <ElButton
          type="primary"
          :loading="applying"
          :disabled="!applyMessage.trim()"
          class="!bg-purple-600 hover:!bg-purple-700"
          @click="submitApply"
        >
          提交申请
        </ElButton>
      </template>
    </ElDialog>
  </Page>
</template>

<style scoped>
.prose p {
  margin-bottom: 0.75rem;
}
</style>
