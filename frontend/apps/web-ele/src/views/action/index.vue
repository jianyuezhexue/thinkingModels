<template>
  <Page
    description="管理你的行动清单，追踪完成进度，记录跟进内容"
    title="我的行动"
  >
    <!-- 行动列表 -->
    <ElCard shadow="never" class="action-list-card">
      <template #header>
        <div class="flex flex-wrap items-center justify-between gap-4">
          <div class="flex items-center gap-4">
            <h2 class="text-lg font-semibold">行动清单</h2>
            <ElSelect v-model="statusFilter" style="width: 140px" size="small">
              <ElOption
                v-for="opt in statusOptions"
                :key="opt.value"
                :label="opt.label"
                :value="opt.value"
              />
            </ElSelect>
          </div>
          <div class="flex items-center gap-3">
            <ElSelect
              v-model="priorityFilter"
              placeholder="优先级"
              clearable
              style="width: 120px"
              size="small"
            >
              <ElOption
                v-for="opt in priorityOptions"
                :key="opt.value"
                :label="opt.label"
                :value="opt.value"
              />
            </ElSelect>
            <ElSelect
              v-model="topicFilter"
              placeholder="选择课题"
              clearable
              style="width: 160px"
              size="small"
            >
              <ElOption
                v-for="topic in topicOptions"
                :key="topic.id"
                :label="topic.title"
                :value="topic.id"
              />
            </ElSelect>
            <ElSelect
              v-model="sortBy"
              style="width: 120px"
              size="small"
            >
              <ElOption
                v-for="opt in sortOptions"
                :key="opt.value"
                :label="opt.label"
                :value="opt.value"
              />
            </ElSelect>
            <ElInput
              v-model="searchKeyword"
              placeholder="搜索行动..."
              clearable
              style="width: 180px"
              size="small"
              @keyup.enter="fetchActions"
            >
              <template #prefix>
                <svg class="h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
                </svg>
              </template>
            </ElInput>
          </div>
        </div>
      </template>

      <!-- 加载状态 -->
      <div v-if="loading" class="space-y-4">
        <ElSkeleton v-for="i in 3" :key="i" animated>
          <template #template>
            <div class="flex items-center gap-4 py-4">
              <ElSkeletonItem variant="text" style="width: 40px" />
              <ElSkeletonItem variant="text" style="width: 300px" />
              <ElSkeletonItem variant="text" style="width: 120px" />
              <ElSkeletonItem variant="text" style="width: 100px" />
            </div>
          </template>
        </ElSkeleton>
      </div>

      <!-- 空状态 -->
      <ElEmpty v-else-if="actions.length === 0" description="暂无行动，从课题分析中创建你的第一个行动吧！">
        <template #image>
          <div class="empty-illustration">
            <svg class="h-24 w-24 text-emerald-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
          </div>
        </template>
        <ElButton type="primary" @click="goToTopics" class="btn-primary-gradient mt-4">
          去创建课题
        </ElButton>
      </ElEmpty>

      <!-- 行动列表 -->
      <div v-else class="action-list">
        <div
          v-for="action in actions"
          :key="action.id"
          class="action-card"
          :class="{ 'action-completed': action.status === 'completed' }"
          @click="openActionDetail(action)"
        >
          <!-- 复选框 -->
          <div class="action-checkbox" @click.stop="toggleActionStatus(action)">
            <div class="checkbox-circle" :class="{ checked: action.status === 'completed' }">
              <svg v-if="action.status === 'completed'" class="h-4 w-4" fill="currentColor" viewBox="0 0 24 24">
                <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/>
              </svg>
            </div>
          </div>

          <!-- 行动内容 -->
          <div class="action-content">
            <div class="action-header">
              <div class="action-title-wrapper">
                <ElTag
                  :type="getPriorityType(action.priority)"
                  size="small"
                  effect="light"
                  class="priority-tag"
                >
                  {{ getPriorityText(action.priority) }}
                </ElTag>
                <h3 class="action-title" :class="{ completed: action.status === 'completed' }">
                  {{ action.title }}
                </h3>
              </div>
              <ElTag :type="getStatusType(action.status)" size="small" effect="light">
                {{ getStatusText(action.status) }}
              </ElTag>
            </div>

            <p class="action-desc">{{ action.description }}</p>

            <!-- 进度条 -->
            <div class="action-progress-row">
              <div class="progress-bar">
                <div class="progress-fill" :style="{ width: `${action.completionRate}%` }" :class="`progress-${action.priority}`" />
              </div>
              <span class="progress-text">{{ action.completionRate }}%</span>
            </div>

            <!-- 元信息 -->
            <div class="action-meta">
              <span class="meta-item">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"/>
                </svg>
                {{ action.topicTitle }}
              </span>
              <span class="meta-item" v-if="action.followUpRecords.length > 0">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z"/>
                </svg>
                {{ action.followUpRecords.length }} 条跟进
              </span>
              <span class="meta-item" v-if="action.dueDate" :class="{ 'due-soon': isDueSoon(action.dueDate), 'due-overdue': isOverdue(action.dueDate) }">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
                {{ getDueText(action.dueDate) }}
              </span>
              <span class="meta-item" v-else>
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                </svg>
                {{ formatDate(action.createdAt) }}
              </span>
            </div>
          </div>

          <!-- 快捷操作 -->
          <div class="action-quick-actions">
            <ElButton type="primary" size="small" text @click.stop="openFollowUpFromList(action)">
              <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
              </svg>
              记录跟进
            </ElButton>
          </div>
        </div>

        <!-- 分页 -->
        <div class="flex justify-center pt-6">
          <ElPagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :total="total"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next"
          />
        </div>
      </div>
    </ElCard>

    <!-- 行动详情弹窗 -->
    <ElDialog
      v-model="detailDialogVisible"
      title="行动详情"
      width="700px"
      class="action-detail-dialog"
      :close-on-click-modal="false"
    >
      <div v-if="selectedAction" class="action-detail-content">
        <!-- 头部信息 -->
        <div class="detail-header">
          <div class="detail-title-row">
            <h2 class="detail-title">{{ selectedAction.title }}</h2>
            <ElTag :type="getStatusType(selectedAction.status)" size="small" effect="light">
              {{ getStatusText(selectedAction.status) }}
            </ElTag>
          </div>
          <p class="detail-desc">{{ selectedAction.description }}</p>
          <div class="detail-topic">
            <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"/>
            </svg>
            <span>关联课题：{{ selectedAction.topicTitle }}</span>
          </div>
        </div>

        <!-- 指导原则 -->
        <div class="detail-section">
          <div class="section-title">
            <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"/>
            </svg>
            指导原则
          </div>
          <div class="guiding-principle-box">
            {{ selectedAction.guidingPrinciple || '暂无指导原则' }}
          </div>
        </div>

        <!-- 完成度 -->
        <div class="detail-section">
          <div class="section-title">
            <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
            </svg>
            完成度
          </div>
          <div class="completion-control">
            <ElSlider v-model="editingCompletionRate" :max="100" :step="5" show-stops />
            <span class="completion-value">{{ editingCompletionRate }}%</span>
          </div>
          <div class="completion-actions">
            <ElButton type="primary" size="small" @click="saveCompletionRate">
              保存进度
            </ElButton>
            <ElButton
              v-if="editingCompletionRate === 100 && selectedAction.status !== 'completed'"
              type="success"
              size="small"
              @click="markAsCompleted"
            >
              标记为已完成
            </ElButton>
          </div>
        </div>

        <!-- 跟进记录入口 -->
        <div class="detail-section">
          <div class="section-title">
            <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z"/>
            </svg>
            跟进记录
            <span class="record-count">({{ selectedAction.followUpRecords.length }})</span>
          </div>
          <div class="follow-up-summary">
            <div class="follow-up-preview" v-if="latestFollowUp">
              <div class="latest-record">
                <span class="latest-label">最新跟进：</span>
                <span class="latest-content">{{ latestFollowUp.content }}</span>
                <span class="latest-time">{{ formatDateTime(latestFollowUp.createdAt) }}</span>
              </div>
            </div>
            <div v-else class="no-follow-up-inline">
              暂无跟进记录
            </div>
            <ElButton type="primary" @click="openFollowUpDrawer" class="mt-3">
              <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
              </svg>
              {{ selectedAction.followUpRecords.length > 0 ? '查看全部' : '添加跟进' }}
            </ElButton>
          </div>
        </div>
      </div>

      <template #footer>
        <ElButton @click="detailDialogVisible = false">关闭</ElButton>
      </template>
    </ElDialog>

    <!-- 跟进记录抽屉 -->
    <ElDrawer
      v-model="followUpDrawerVisible"
      title="跟进记录"
      size="450px"
      :with-header="true"
      class="follow-up-drawer"
    >
      <div v-if="selectedAction" class="drawer-content">
        <!-- 行动简要信息 -->
        <div class="drawer-action-header">
          <div class="drawer-action-title">{{ selectedAction.title }}</div>
          <div class="drawer-action-meta">
            <ElTag :type="getStatusType(selectedAction.status)" size="small" effect="light">
              {{ getStatusText(selectedAction.status) }}
            </ElTag>
            <span class="drawer-progress">{{ selectedAction.completionRate }}%</span>
          </div>
        </div>

        <!-- 新增跟进 -->
        <div class="drawer-add-section">
          <div class="drawer-section-title">添加跟进</div>
          <ElInput
            v-model="newFollowUpContent"
            type="textarea"
            :rows="4"
            placeholder="记录你的跟进内容、遇到的问题、下一步计划..."
            resize="none"
          />
          <ElButton
            type="primary"
            class="mt-3 w-full"
            :disabled="!newFollowUpContent.trim()"
            @click="addFollowUp"
          >
            添加跟进
          </ElButton>
        </div>

        <!-- 跟进记录时间轴 -->
        <div class="drawer-timeline-section">
          <div class="drawer-section-title">
            历史记录
            <span class="record-count">({{ selectedAction.followUpRecords.length }})</span>
          </div>
          <div v-if="selectedAction.followUpRecords.length > 0" class="timeline-list">
            <div
              v-for="(record, index) in sortedFollowUpRecords"
              :key="record.id"
              class="timeline-item"
            >
              <div class="timeline-marker"></div>
              <div class="timeline-content">
                <div class="timeline-header">
                  <span class="timeline-index">#{{ selectedAction.followUpRecords.length - index }}</span>
                  <span class="timeline-time">{{ formatDateTime(record.createdAt) }}</span>
                </div>
                <div class="timeline-body">{{ record.content }}</div>
              </div>
            </div>
          </div>
          <div v-else class="timeline-empty">
            <svg class="h-12 w-12 text-gray-300 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
            </svg>
            <p>暂无跟进记录</p>
            <p class="text-sm text-gray-400">添加第一条记录开始追踪进度</p>
          </div>
        </div>
      </div>
    </ElDrawer>
  </Page>
</template>

<script lang="ts" setup>
import { onMounted, ref, watch, computed } from 'vue';
import { useRouter } from 'vue-router';

import { Page } from '@vben/common-ui';

import {
  ElButton,
  ElCard,
  ElDialog,
  ElDrawer,
  ElEmpty,
  ElInput,
  ElOption,
  ElPagination,
  ElSelect,
  ElSkeleton,
  ElSkeletonItem,
  ElSlider,
  ElTag,
  ElMessage,
} from 'element-plus';

import { getMyActionsApi, updateCompletionRateApi, addFollowUpApi, updateActionStatusApi } from '#/api/action';
import { getMyTopicsApi } from '#/api';
import type { ActionApi } from '#/api/action';

// 路由
const router = useRouter();

// 加载状态
const loading = ref(false);

// 行动列表
const actions = ref<ActionApi.Action[]>([]);
const total = ref(0);

// 分页
const currentPage = ref(1);
const pageSize = ref(10);

// 筛选
const searchKeyword = ref('');
const statusFilter = ref<ActionApi.ActionStatus | 'all'>('all');
const topicFilter = ref<string>('');
const topicOptions = ref<{ id: string; title: string }[]>([]);
const priorityFilter = ref<ActionApi.ActionPriority | 'all'>('all');
const sortBy = ref<'createdAt' | 'dueDate' | 'completionRate' | 'priority'>('priority');

// 弹窗/抽屉状态
const detailDialogVisible = ref(false);
const followUpDrawerVisible = ref(false);
const selectedAction = ref<ActionApi.Action | null>(null);
const newFollowUpContent = ref('');
const editingCompletionRate = ref(0);

// 统计
const totalActions = computed(() => actions.value.length);
const pendingActions = computed(() => actions.value.filter(a => a.status === 'pending').length);
const inProgressActions = computed(() => actions.value.filter(a => a.status === 'in_progress').length);
const completedActions = computed(() => actions.value.filter(a => a.status === 'completed').length);
const urgentActions = computed(() => actions.value.filter(a => {
  // 高优先级未完成，或即将到期（3天内）未完成
  if (a.status === 'completed') return false;
  if (a.priority === 'high') return true;
  if (a.dueDate && isDueSoon(a.dueDate)) return true;
  return false;
}).length);

// 跟进记录排序（最新的在前面）
const sortedFollowUpRecords = computed(() => {
  if (!selectedAction.value) return [];
  return [...selectedAction.value.followUpRecords].reverse();
});

// 获取最新的跟进记录
const latestFollowUp = computed(() => {
  if (!selectedAction.value || selectedAction.value.followUpRecords.length === 0) return null;
  return selectedAction.value.followUpRecords[selectedAction.value.followUpRecords.length - 1];
});

// 状态选项
const statusOptions = [
  { value: 'all', label: '全部状态', type: '' },
  { value: 'pending', label: '待执行', type: 'info' },
  { value: 'in_progress', label: '进行中', type: 'warning' },
  { value: 'completed', label: '已完成', type: 'success' },
  { value: 'cancelled', label: '已取消', type: '' },
];

// 优先级选项
const priorityOptions = [
  { value: 'all', label: '全部优先级' },
  { value: 'high', label: '高优先级' },
  { value: 'medium', label: '中优先级' },
  { value: 'low', label: '低优先级' },
];

// 排序选项
const sortOptions = [
  { value: 'priority', label: '按优先级' },
  { value: 'dueDate', label: '按截止日期' },
  { value: 'completionRate', label: '按完成度' },
  { value: 'createdAt', label: '按创建时间' },
];

// 优先级映射
const priorityMap: Record<ActionApi.ActionPriority, number> = {
  high: 3,
  medium: 2,
  low: 1,
};

// 获取行动列表
async function fetchActions() {
  loading.value = true;
  try {
    const params: ActionApi.ActionListParams = {
      page: currentPage.value,
      pageSize: pageSize.value,
      keyword: searchKeyword.value || undefined,
      sortBy: sortBy.value,
    };

    if (statusFilter.value !== 'all') {
      params.status = statusFilter.value;
    }

    if (topicFilter.value) {
      params.topicId = topicFilter.value;
    }

    if (priorityFilter.value !== 'all') {
      params.priority = priorityFilter.value;
    }

    const res = await getMyActionsApi(params);

    // Mock 数据，添加更多详细信息
    const priorities: ActionApi.ActionPriority[] = ['high', 'medium', 'low'];
    const principles = [
      'SMART原则：具体的、可衡量的、可达成的、相关的、有时限的',
      '先做重要不紧急的事，避免总是救火',
      '设定检查点，每完成25%回顾一次方向是否正确',
      '遇到困难时，先拆解为更小可执行的动作',
    ];

    actions.value = res.list.map((action: ActionApi.Action, index: number) => {
      const dueDate = new Date();
      dueDate.setDate(dueDate.getDate() + (index % 7) - 2);

      const principle = action.guidingPrinciple || principles[index % 4];

      return {
        ...action,
        priority: action.priority || priorities[index % 3],
        guidingPrinciple: principle,
        completionRate: action.completionRate ?? (index % 5) * 20,
        followUpRecords: action.followUpRecords || [],
        dueDate: action.dueDate || dueDate.toISOString(),
      };
    });

    // 前端排序（Mock 阶段）
    actions.value.sort((a, b) => {
      if (sortBy.value === 'priority') {
        return priorityMap[b.priority] - priorityMap[a.priority];
      }
      if (sortBy.value === 'completionRate') {
        return b.completionRate - a.completionRate;
      }
      if (sortBy.value === 'dueDate') {
        return new Date(a.dueDate || 0).getTime() - new Date(b.dueDate || 0).getTime();
      }
      return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime();
    });

    total.value = res.total;
  } catch (error) {
    console.error('获取行动列表失败:', error);
    ElMessage.error('获取行动列表失败');
  } finally {
    loading.value = false;
  }
}

// 获取课题选项
async function fetchTopicOptions() {
  try {
    const res = await getMyTopicsApi({ pageSize: 100 });
    topicOptions.value = res.list.map(topic => ({
      id: topic.id,
      title: topic.title,
    }));
  } catch (error) {
    console.error('获取课题列表失败:', error);
  }
}

// 监听筛选条件变化
watch([searchKeyword, statusFilter, topicFilter, priorityFilter, sortBy], () => {
  currentPage.value = 1;
  fetchActions();
});

// 监听分页变化
watch([currentPage, pageSize], () => {
  fetchActions();
});

// Mock 数据
const mockActions: ActionApi.Action[] = [
  {
    id: '1',
    title: '重新评估目标用户群体，缩小范围至核心用户',
    description: '通过用户调研和数据分析，明确产品的核心用户画像，避免资源分散',
    topicId: 't1',
    topicTitle: '产品市场定位分析',
    guidingPrinciple: 'SMART原则：具体的、可衡量的、可达成的、相关的、有时限的',
    completionRate: 75,
    status: 'in_progress',
    priority: 'high',
    followUpRecords: [
      { id: 'f1', content: '已完成20份用户问卷收集', createdAt: '2024-01-15T10:30:00Z' },
      { id: 'f2', content: '初步筛选出3个潜在用户群体', createdAt: '2024-01-16T14:20:00Z' },
    ],
    createdAt: '2024-01-10T08:00:00Z',
    updatedAt: '2024-01-16T14:20:00Z',
    dueDate: '2024-01-20T23:59:59Z',
  },
  {
    id: '2',
    title: '制定3个月内可执行的MVP功能清单',
    description: '基于核心用户需求，确定MVP版本的功能范围和优先级',
    topicId: 't1',
    topicTitle: '产品市场定位分析',
    guidingPrinciple: '先做重要不紧急的事，避免总是救火',
    completionRate: 30,
    status: 'in_progress',
    priority: 'high',
    followUpRecords: [
      { id: 'f3', content: '已列出15个候选功能', createdAt: '2024-01-14T09:00:00Z' },
    ],
    createdAt: '2024-01-12T10:00:00Z',
    updatedAt: '2024-01-14T09:00:00Z',
    dueDate: '2024-01-25T23:59:59Z',
  },
  {
    id: '3',
    title: '寻找2-3位潜在用户进行深度访谈',
    description: '通过一对一访谈深入了解用户痛点和需求',
    topicId: 't1',
    topicTitle: '产品市场定位分析',
    guidingPrinciple: '遇到困难时，先拆解为更小可执行的动作',
    completionRate: 0,
    status: 'pending',
    priority: 'medium',
    followUpRecords: [],
    createdAt: '2024-01-15T08:00:00Z',
    updatedAt: '2024-01-15T08:00:00Z',
    dueDate: '2024-01-30T23:59:59Z',
  },
  {
    id: '4',
    title: '完成竞品分析报告',
    description: '分析主要竞争对手的产品特点、定价策略和市场份额',
    topicId: 't2',
    topicTitle: '竞品调研项目',
    guidingPrinciple: '设定检查点，每完成25%回顾一次方向是否正确',
    completionRate: 100,
    status: 'completed',
    priority: 'medium',
    followUpRecords: [
      { id: 'f4', content: '确定5家主要竞品', createdAt: '2024-01-08T10:00:00Z' },
      { id: 'f5', content: '完成功能对比表格', createdAt: '2024-01-10T15:30:00Z' },
      { id: 'f6', content: '撰写分析报告', createdAt: '2024-01-12T16:00:00Z' },
    ],
    createdAt: '2024-01-05T08:00:00Z',
    updatedAt: '2024-01-12T16:00:00Z',
    dueDate: '2024-01-15T23:59:59Z',
  },
  {
    id: '5',
    title: '设计产品原型图',
    description: '基于MVP功能清单，设计核心流程的原型图',
    topicId: 't3',
    topicTitle: '产品设计规划',
    guidingPrinciple: '先做重要不紧急的事，避免总是救火',
    completionRate: 45,
    status: 'in_progress',
    priority: 'low',
    followUpRecords: [
      { id: 'f7', content: '完成首页原型设计', createdAt: '2024-01-16T11:00:00Z' },
    ],
    createdAt: '2024-01-14T08:00:00Z',
    updatedAt: '2024-01-16T11:00:00Z',
    dueDate: '2024-02-05T23:59:59Z',
  },
  {
    id: '6',
    title: '准备投资人路演PPT',
    description: '整理商业模式、市场分析和财务预测，制作路演材料',
    topicId: 't4',
    topicTitle: '融资计划',
    guidingPrinciple: 'SMART原则：具体的、可衡量的、可达成的、相关的、有时限的',
    completionRate: 10,
    status: 'pending',
    priority: 'high',
    followUpRecords: [],
    createdAt: '2024-01-18T08:00:00Z',
    updatedAt: '2024-01-18T08:00:00Z',
    dueDate: '2024-01-22T23:59:59Z',
  },
];

// Mock 课题选项
const mockTopics = [
  { id: 't1', title: '产品市场定位分析' },
  { id: 't2', title: '竞品调研项目' },
  { id: 't3', title: '产品设计规划' },
  { id: 't4', title: '融资计划' },
];

// 页面加载时获取数据
onMounted(() => {
  // 使用 Mock 数据
  actions.value = mockActions;
  total.value = mockActions.length;
  topicOptions.value = mockTopics;
  // fetchActions();
  // fetchTopicOptions();
});

// 打开行动详情
function openActionDetail(action: ActionApi.Action) {
  selectedAction.value = action;
  editingCompletionRate.value = action.completionRate;
  newFollowUpContent.value = '';
  detailDialogVisible.value = true;
}

// 打开跟进记录抽屉
function openFollowUpDrawer() {
  newFollowUpContent.value = '';
  followUpDrawerVisible.value = true;
}

// 从列表直接打开跟进抽屉
function openFollowUpFromList(action: ActionApi.Action) {
  selectedAction.value = action;
  newFollowUpContent.value = '';
  followUpDrawerVisible.value = true;
}

// 切换行动状态
async function toggleActionStatus(action: ActionApi.Action) {
  const newStatus = action.status === 'completed' ? 'pending' : 'completed';
  try {
    await updateActionStatusApi({ id: action.id, status: newStatus });
    action.status = newStatus;
    if (newStatus === 'completed') {
      action.completionRate = 100;
    }
    ElMessage.success(newStatus === 'completed' ? '行动已完成' : '已恢复为待执行');
  } catch (error) {
    console.error('更新状态失败:', error);
    ElMessage.error('操作失败');
  }
}

// 保存完成度
async function saveCompletionRate() {
  if (!selectedAction.value) return;
  try {
    await updateCompletionRateApi({
      id: selectedAction.value.id,
      completionRate: editingCompletionRate.value,
    });
    selectedAction.value.completionRate = editingCompletionRate.value;
    // 更新列表中的数据
    const actionInList = actions.value.find(a => a.id === selectedAction.value!.id);
    if (actionInList) {
      actionInList.completionRate = editingCompletionRate.value;
    }
    ElMessage.success('完成度已更新');
  } catch (error) {
    console.error('更新完成度失败:', error);
    ElMessage.error('保存失败');
  }
}

// 标记为已完成
async function markAsCompleted() {
  if (!selectedAction.value) return;
  try {
    await updateActionStatusApi({ id: selectedAction.value.id, status: 'completed' });
    selectedAction.value.status = 'completed';
    selectedAction.value.completionRate = 100;
    editingCompletionRate.value = 100;
    // 更新列表中的数据
    const actionInList = actions.value.find(a => a.id === selectedAction.value?.id);
    if (actionInList) {
      actionInList.status = 'completed';
      actionInList.completionRate = 100;
    }
    ElMessage.success('行动已标记为完成');
  } catch (error) {
    console.error('更新状态失败:', error);
    ElMessage.error('操作失败');
  }
}

// 添加跟进记录
async function addFollowUp() {
  if (!selectedAction.value || !newFollowUpContent.value.trim()) return;
  try {
    const record = await addFollowUpApi({
      actionId: selectedAction.value.id,
      content: newFollowUpContent.value.trim(),
    });
    selectedAction.value.followUpRecords.push(record);
    // 更新列表中的数据
    const actionInList = actions.value.find(a => a.id === selectedAction.value!.id);
    if (actionInList) {
      actionInList.followUpRecords.push(record);
    }
    newFollowUpContent.value = '';
    ElMessage.success('跟进记录已添加');
  } catch (error) {
    console.error('添加跟进记录失败:', error);
    ElMessage.error('添加失败');
  }
}

// 获取状态标签类型
function getStatusType(status: ActionApi.ActionStatus): any {
  const map: Record<ActionApi.ActionStatus, any> = {
    pending: 'info',
    in_progress: 'warning',
    completed: 'success',
    cancelled: '',
  };
  return map[status] || '';
}

// 获取状态标签文本
function getStatusText(status: ActionApi.ActionStatus): string {
  const map: Record<ActionApi.ActionStatus, string> = {
    pending: '待执行',
    in_progress: '进行中',
    completed: '已完成',
    cancelled: '已取消',
  };
  return map[status] || status;
}

// 获取优先级标签类型
function getPriorityType(priority: ActionApi.ActionPriority): any {
  const map: Record<ActionApi.ActionPriority, any> = {
    high: 'danger',
    medium: 'warning',
    low: 'info',
  };
  return map[priority] || 'info';
}

// 获取优先级标签文本
function getPriorityText(priority: ActionApi.ActionPriority): string {
  const map: Record<ActionApi.ActionPriority, string> = {
    high: '高',
    medium: '中',
    low: '低',
  };
  return map[priority] || priority;
}

// 检查是否即将到期（3天内）
function isDueSoon(dueDate: string): boolean {
  const due = new Date(dueDate);
  const now = new Date();
  const diffDays = Math.ceil((due.getTime() - now.getTime()) / (1000 * 60 * 60 * 24));
  return diffDays <= 3 && diffDays >= 0;
}

// 检查是否已逾期
function isOverdue(dueDate: string): boolean {
  const due = new Date(dueDate);
  const now = new Date();
  return due.getTime() < now.getTime();
}

// 获取截止日期显示文本
function getDueText(dueDate: string): string {
  const due = new Date(dueDate);
  const now = new Date();
  const diffDays = Math.ceil((due.getTime() - now.getTime()) / (1000 * 60 * 60 * 24));

  if (diffDays < 0) {
    return `已逾期 ${Math.abs(diffDays)} 天`;
  } else if (diffDays === 0) {
    return '今天到期';
  } else if (diffDays === 1) {
    return '明天到期';
  } else if (diffDays <= 3) {
    return `${diffDays} 天后到期`;
  }
  return `截止 ${formatDate(dueDate)}`;
}

// 格式化日期
function formatDate(dateStr: string): string {
  if (!dateStr) return '-';
  const date = new Date(dateStr);
  return date.toLocaleDateString('zh-CN');
}

// 格式化日期时间
function formatDateTime(dateStr: string): string {
  if (!dateStr) return '-';
  const date = new Date(dateStr);
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`;
}

// 跳转到课题页面
function goToTopics() {
  router.push('/my-topics');
}
</script>

<style scoped>
/* 统计卡片 */
.stat-card {
  background: white;
  border-radius: 12px;
  padding: 16px;
  border: 1px solid #e5e7eb;
  transition: all 0.3s ease;
}

.stat-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  transform: translateY(-2px);
}

.stat-card-primary {
  border-left: 4px solid #6366f1;
}

.stat-card-success {
  border-left: 4px solid #10b981;
}

.stat-card-warning {
  border-left: 4px solid #f59e0b;
}

.stat-card-info {
  border-left: 4px solid #3b82f6;
}

.stat-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 12px;
}

.stat-card-primary .stat-icon {
  background: #eef2ff;
  color: #6366f1;
}

.stat-card-success .stat-icon {
  background: #ecfdf5;
  color: #10b981;
}

.stat-card-warning .stat-icon {
  background: #fffbeb;
  color: #f59e0b;
}

.stat-card-info .stat-icon {
  background: #eff6ff;
  color: #3b82f6;
}

.stat-card-danger {
  border-left: 4px solid #ef4444;
}

.stat-card-danger .stat-icon {
  background: #fef2f2;
  color: #ef4444;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #1f2937;
  line-height: 1;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 13px;
  color: #6b7280;
}

/* 按钮样式 */
.btn-primary-gradient {
  background: linear-gradient(135deg, #6366f1 0%, #818cf8 100%);
  border: none;
  transition: all 0.3s ease;
}

.btn-primary-gradient:hover {
  background: linear-gradient(135deg, #4f46e5 0%, #6366f1 100%);
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.3);
}

/* 行动列表卡片 */
.action-list-card {
  width: 100%;
}

.action-list-card :deep(.el-card__body) {
  padding: 20px;
}

/* 行动列表 */
.action-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* 行动卡片 */
.action-card {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  padding: 16px 20px;
  background: white;
  border-radius: 12px;
  border: 1px solid #e5e7eb;
  cursor: pointer;
  transition: all 0.3s ease;
}

.action-card:hover {
  border-color: #c4b5fd;
  box-shadow: 0 4px 16px rgba(99, 102, 241, 0.12);
}

.action-completed {
  background: #f9fafb;
  opacity: 0.8;
}

/* 复选框 */
.action-checkbox {
  flex-shrink: 0;
  padding-top: 4px;
}

.checkbox-circle {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  border: 2px solid #d1d5db;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
  color: white;
}

.checkbox-circle:hover {
  border-color: #6366f1;
}

.checkbox-circle.checked {
  background: #10b981;
  border-color: #10b981;
}

/* 行动内容 */
.action-content {
  flex: 1 1 auto;
  min-width: 0;
}

.action-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 6px;
}

.action-title-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  min-width: 0;
}

.priority-tag {
  flex-shrink: 0;
  font-size: 11px;
  padding: 0 6px;
  height: 18px;
  line-height: 16px;
}

.action-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.action-title.completed {
  text-decoration: line-through;
  color: #9ca3af;
}

.action-desc {
  font-size: 13px;
  color: #6b7280;
  margin-bottom: 10px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
}

/* 进度条行 */
.action-progress-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 10px;
}

.progress-bar {
  flex: 1;
  height: 6px;
  background: #e5e7eb;
  border-radius: 3px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #6366f1 0%, #818cf8 100%);
  border-radius: 3px;
  transition: width 0.3s ease;
}

.action-completed .progress-fill {
  background: #10b981;
}

.progress-high {
  background: linear-gradient(90deg, #ef4444 0%, #f87171 100%);
}

.progress-medium {
  background: linear-gradient(90deg, #f59e0b 0%, #fbbf24 100%);
}

.progress-low {
  background: linear-gradient(90deg, #3b82f6 0%, #60a5fa 100%);
}

.progress-text {
  font-size: 12px;
  font-weight: 600;
  color: #6366f1;
  min-width: 36px;
}

.action-completed .progress-text {
  color: #10b981;
}

/* 元信息 */
.action-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  font-size: 12px;
  color: #9ca3af;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.meta-item.due-soon {
  color: #f59e0b;
  font-weight: 500;
}

.meta-item.due-overdue {
  color: #ef4444;
  font-weight: 500;
}

/* 快捷操作 */
.action-quick-actions {
  flex-shrink: 0;
  padding-top: 4px;
}

/* 空状态 */
.empty-illustration {
  display: flex;
  justify-content: center;
  margin-bottom: 16px;
}

/* 弹窗样式 */
.action-detail-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #eef2ff 0%, #e0e7ff 100%);
  margin-right: 0;
  padding: 20px;
  border-bottom: 1px solid #c4b5fd;
}

.action-detail-dialog :deep(.el-dialog__title) {
  color: #4f46e5;
  font-weight: 600;
}

.action-detail-content {
  padding: 8px 0;
}

/* 详情头部 */
.detail-header {
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px dashed #e5e7eb;
}

.detail-title-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.detail-title {
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.detail-desc {
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 12px;
  line-height: 1.5;
}

.detail-topic {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #6366f1;
  font-weight: 500;
}

/* 详情区块 */
.detail-section {
  margin-bottom: 24px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 15px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #f3f4f6;
}

.section-title svg {
  color: #6366f1;
}

.record-count {
  font-size: 13px;
  color: #9ca3af;
  font-weight: 400;
}

/* 指导原则 */
.guiding-principle-box {
  background: linear-gradient(135deg, #f5f3ff 0%, #ede9fe 100%);
  border-left: 4px solid #6366f1;
  padding: 16px 20px;
  border-radius: 0 12px 12px 0;
  font-size: 14px;
  color: #4c1d95;
  line-height: 1.6;
}

/* 完成度控制 */
.completion-control {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 12px;
}

.completion-control :deep(.el-slider) {
  flex: 1;
}

.completion-value {
  font-size: 20px;
  font-weight: 700;
  color: #6366f1;
  min-width: 50px;
}

.completion-actions {
  display: flex;
  gap: 12px;
}

/* 跟进记录 */
.add-follow-up {
  background: #f9fafb;
  padding: 16px;
  border-radius: 12px;
  margin-bottom: 16px;
}

.follow-up-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.follow-up-item {
  display: flex;
  gap: 12px;
  padding: 14px 16px;
  background: #f9fafb;
  border-radius: 10px;
  border-left: 3px solid #10b981;
}

.follow-up-index {
  flex-shrink: 0;
  width: 28px;
  height: 28px;
  background: #10b981;
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
}

.follow-up-body {
  flex: 1;
}

.follow-up-content {
  font-size: 14px;
  color: #374151;
  line-height: 1.6;
  margin-bottom: 6px;
}

.follow-up-time {
  font-size: 12px;
  color: #9ca3af;
}

.no-follow-up {
  text-align: center;
  padding: 32px;
  color: #9ca3af;
  font-size: 14px;
  background: #f9fafb;
  border-radius: 12px;
}

/* 跟进摘要（详情弹窗中） */
.follow-up-summary {
  padding: 16px;
  background: #f9fafb;
  border-radius: 12px;
}

.follow-up-preview {
  margin-bottom: 12px;
}

.latest-record {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.latest-label {
  font-size: 13px;
  color: #6366f1;
  font-weight: 500;
}

.latest-content {
  font-size: 14px;
  color: #374151;
  line-height: 1.5;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.latest-time {
  font-size: 12px;
  color: #9ca3af;
}

.no-follow-up-inline {
  font-size: 14px;
  color: #9ca3af;
  padding: 12px 0;
  text-align: center;
}

/* 跟进抽屉样式 */
.follow-up-drawer :deep(.el-drawer__header) {
  background: linear-gradient(135deg, #ecfdf5 0%, #d1fae5 100%);
  margin-bottom: 0;
  padding: 20px;
  border-bottom: 1px solid #a7f3d0;
}

.follow-up-drawer :deep(.el-drawer__title) {
  color: #059669;
  font-weight: 600;
  font-size: 16px;
}

.drawer-content {
  padding: 20px;
  height: 100%;
  overflow-y: auto;
}

/* 抽屉头部行动信息 */
.drawer-action-header {
  background: linear-gradient(135deg, #f5f3ff 0%, #ede9fe 100%);
  border-left: 4px solid #6366f1;
  padding: 16px 20px;
  border-radius: 0 12px 12px 0;
  margin-bottom: 24px;
}

.drawer-action-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 10px;
  line-height: 1.5;
}

.drawer-action-meta {
  display: flex;
  align-items: center;
  gap: 12px;
}

.drawer-progress {
  font-size: 14px;
  font-weight: 700;
  color: #6366f1;
}

/* 抽屉区块标题 */
.drawer-section-title {
  font-size: 14px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 12px;
  display: flex;
  align-items: center;
  gap: 6px;
}

/* 新增跟进区域 */
.drawer-add-section {
  background: #f9fafb;
  padding: 16px;
  border-radius: 12px;
  margin-bottom: 24px;
}

/* 时间轴区域 */
.drawer-timeline-section {
  padding-bottom: 20px;
}

/* 时间轴样式 */
.timeline-list {
  position: relative;
  padding-left: 8px;
}

.timeline-list::before {
  content: '';
  position: absolute;
  left: 16px;
  top: 0;
  bottom: 0;
  width: 2px;
  background: linear-gradient(180deg, #10b981 0%, #34d399 100%);
}

.timeline-item {
  position: relative;
  padding-left: 32px;
  padding-bottom: 20px;
}

.timeline-item:last-child {
  padding-bottom: 0;
}

.timeline-marker {
  position: absolute;
  left: -8px;
  top: 4px;
  width: 16px;
  height: 16px;
  background: #10b981;
  border-radius: 50%;
  border: 3px solid white;
  box-shadow: 0 0 0 2px #10b981;
}

.timeline-content {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  padding: 14px 16px;
}

.timeline-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.timeline-index {
  font-size: 12px;
  font-weight: 600;
  color: #10b981;
  background: #ecfdf5;
  padding: 2px 8px;
  border-radius: 10px;
}

.timeline-time {
  font-size: 12px;
  color: #9ca3af;
}

.timeline-body {
  font-size: 14px;
  color: #374151;
  line-height: 1.6;
}

.timeline-empty {
  text-align: center;
  padding: 40px 20px;
  color: #9ca3af;
}

.timeline-empty p {
  margin: 0;
}

.timeline-empty p:last-child {
  margin-top: 4px;
}

/* 响应式 */
@media (max-width: 768px) {
  .action-card {
    flex-direction: column;
    gap: 12px;
  }

  .action-quick-actions {
    padding-top: 0;
    width: 100%;
  }

  .completion-control {
    flex-direction: column;
    gap: 12px;
  }

  .completion-control :deep(.el-slider) {
    width: 100%;
  }
}
</style>
