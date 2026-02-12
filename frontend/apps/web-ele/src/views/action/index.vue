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
  ElSlider,
  ElMessage,
  ElProgress,
} from 'element-plus';
import { updateCompletionRateApi, addFollowUpApi, updateActionStatusApi } from '#/api/action';
import type { ActionApi } from '#/api/action';

const router = useRouter();

// ==================== 状态管理 ====================
const loading = ref(false);
const actions = ref<ActionApi.Action[]>([]);
const total = ref(0);

// 分页
const currentPage = ref(1);
const pageSize = ref(10);

// 筛选
const searchKeyword = ref('');
const activeStatus = ref<ActionApi.ActionStatus | 'all'>('all');
const priorityFilter = ref<ActionApi.ActionPriority | 'all'>('all');
const topicFilter = ref<string>('');
const topicOptions = ref<{ id: string; title: string }[]>([]);
const sortBy = ref<'priority' | 'dueDate' | 'completionRate' | 'createdAt'>('priority');

// 弹窗/抽屉状态
const detailDialogVisible = ref(false);
const followUpDrawerVisible = ref(false);
const selectedAction = ref<ActionApi.Action | null>(null);
const newFollowUpContent = ref('');
const editingCompletionRate = ref(0);

// ==================== 统计数据 ====================
const stats = computed(() => {
  const all = actions.value;
  const pending = all.filter(a => a.status === 'pending').length;
  const inProgress = all.filter(a => a.status === 'in_progress').length;
  const completed = all.filter(a => a.status === 'completed').length;
  const urgent = all.filter(a => {
    if (a.status === 'completed') return false;
    if (a.priority === 'high') return true;
    if (a.dueDate && isDueSoon(a.dueDate)) return true;
    return false;
  }).length;
  return { total: all.length, pending, inProgress, completed, urgent };
});

// 状态选项
const statusTabs = [
  { id: 'all', label: '全部行动', icon: '📋' },
  { id: 'pending', label: '待执行', icon: '⏳' },
  { id: 'in_progress', label: '进行中', icon: '🚀' },
  { id: 'completed', label: '已完成', icon: '✅' },
];

// 优先级选项
const priorityOptions = [
  { value: 'all', label: '全部优先级' },
  { value: 'high', label: '🔴 高优先级' },
  { value: 'medium', label: '🟡 中优先级' },
  { value: 'low', label: '🟢 低优先级' },
];

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

const mockTopics = [
  { id: 't1', title: '产品市场定位分析' },
  { id: 't2', title: '竞品调研项目' },
  { id: 't3', title: '产品设计规划' },
  { id: 't4', title: '融资计划' },
];

// ==================== 数据获取 ====================
async function fetchActions() {
  loading.value = true;
  try {
    // 使用 Mock 数据
    let filtered = [...mockActions];
    
    // 状态筛选
    if (activeStatus.value !== 'all') {
      filtered = filtered.filter(a => a.status === activeStatus.value);
    }
    
    // 优先级筛选
    if (priorityFilter.value !== 'all') {
      filtered = filtered.filter(a => a.priority === priorityFilter.value);
    }
    
    // 课题筛选
    if (topicFilter.value) {
      filtered = filtered.filter(a => a.topicId === topicFilter.value);
    }
    
    // 关键词搜索
    if (searchKeyword.value) {
      const kw = searchKeyword.value.toLowerCase();
      filtered = filtered.filter(a => 
        a.title.toLowerCase().includes(kw) || 
        a.description.toLowerCase().includes(kw)
      );
    }
    
    // 排序
    const priorityMap: Record<ActionApi.ActionPriority, number> = { high: 3, medium: 2, low: 1 };
    filtered.sort((a, b) => {
      if (sortBy.value === 'priority') {
        return (priorityMap[b.priority] || 0) - (priorityMap[a.priority] || 0);
      }
      if (sortBy.value === 'completionRate') {
        return b.completionRate - a.completionRate;
      }
      if (sortBy.value === 'dueDate') {
        return new Date(a.dueDate || 0).getTime() - new Date(b.dueDate || 0).getTime();
      }
      return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime();
    });
    
    actions.value = filtered;
    total.value = filtered.length;
    topicOptions.value = mockTopics;
  } catch (error) {
    console.error('获取行动列表失败:', error);
    ElMessage.error('获取行动列表失败');
  } finally {
    loading.value = false;
  }
}

// ==================== 工具函数 ====================
function getPriorityStyle(priority: ActionApi.ActionPriority): string {
  const styles: Record<string, string> = {
    high: 'bg-red-100 text-red-700',
    medium: 'bg-amber-100 text-amber-700',
    low: 'bg-blue-100 text-blue-700',
  };
  return styles[priority] || 'bg-gray-100 text-gray-600';
}

function getPriorityText(priority: ActionApi.ActionPriority): string {
  const texts: Record<string, string> = {
    high: '高',
    medium: '中',
    low: '低',
  };
  return texts[priority] || priority;
}

function getPriorityIcon(priority: ActionApi.ActionPriority): string {
  const icons: Record<string, string> = {
    high: '🔴',
    medium: '🟡',
    low: '🟢',
  };
  return icons[priority] || '⚪';
}

function getStatusStyle(status: ActionApi.ActionStatus): string {
  const styles: Record<string, string> = {
    pending: 'bg-gray-100 text-gray-600',
    in_progress: 'bg-amber-100 text-amber-700',
    completed: 'bg-green-100 text-green-700',
    cancelled: 'bg-slate-100 text-slate-600',
  };
  return styles[status] || 'bg-gray-100 text-gray-600';
}

function getStatusText(status: ActionApi.ActionStatus): string {
  const texts: Record<string, string> = {
    pending: '待执行',
    in_progress: '进行中',
    completed: '已完成',
    cancelled: '已取消',
  };
  return texts[status] || status;
}

function getProgressColor(action: ActionApi.Action): string {
  if (action.status === 'completed') return '#10b981';
  if (action.priority === 'high') return '#ef4444';
  if (action.priority === 'medium') return '#f59e0b';
  return '#3b82f6';
}

function isDueSoon(dueDate: string): boolean {
  const due = new Date(dueDate);
  const now = new Date();
  const diffDays = Math.ceil((due.getTime() - now.getTime()) / (1000 * 60 * 60 * 24));
  return diffDays <= 3 && diffDays >= 0;
}

function isOverdue(dueDate: string): boolean {
  return new Date(dueDate).getTime() < new Date().getTime();
}

function getDueText(dueDate: string): string {
  const due = new Date(dueDate);
  const now = new Date();
  const diffDays = Math.ceil((due.getTime() - now.getTime()) / (1000 * 60 * 60 * 24));

  if (diffDays < 0) return '已逾期 ' + Math.abs(diffDays) + ' 天';
  if (diffDays === 0) return '今天到期';
  if (diffDays === 1) return '明天到期';
  if (diffDays <= 3) return diffDays + ' 天后到期';
  return formatDate(dueDate) + ' 到期';
}

function formatDate(dateStr: string): string {
  if (!dateStr) return '-';
  return new Date(dateStr).toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' });
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
  return formatDate(dateStr);
}

function formatDateTime(dateStr: string): string {
  if (!dateStr) return '-';
  const date = new Date(dateStr);
  return date.getFullYear() + '-' + 
    String(date.getMonth() + 1).padStart(2, '0') + '-' + 
    String(date.getDate()).padStart(2, '0') + ' ' + 
    String(date.getHours()).padStart(2, '0') + ':' + 
    String(date.getMinutes()).padStart(2, '0');
}

// 跟进记录排序
const sortedFollowUpRecords = computed(() => {
  if (!selectedAction.value) return [];
  return [...selectedAction.value.followUpRecords].reverse();
});

const latestFollowUp = computed(() => {
  if (!selectedAction.value || selectedAction.value.followUpRecords.length === 0) return null;
  return selectedAction.value.followUpRecords[selectedAction.value.followUpRecords.length - 1];
});

// ==================== 操作 ====================
function openActionDetail(action: ActionApi.Action) {
  selectedAction.value = action;
  editingCompletionRate.value = action.completionRate;
  newFollowUpContent.value = '';
  detailDialogVisible.value = true;
}

function openFollowUpDrawer() {
  newFollowUpContent.value = '';
  followUpDrawerVisible.value = true;
}

function openFollowUpFromList(action: ActionApi.Action) {
  selectedAction.value = action;
  newFollowUpContent.value = '';
  followUpDrawerVisible.value = true;
}

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

async function saveCompletionRate() {
  if (!selectedAction.value) return;
  try {
    await updateCompletionRateApi({
      id: selectedAction.value.id,
      completionRate: editingCompletionRate.value,
    });
    selectedAction.value.completionRate = editingCompletionRate.value;
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

async function markAsCompleted() {
  if (!selectedAction.value) return;
  try {
    await updateActionStatusApi({ id: selectedAction.value.id, status: 'completed' });
    selectedAction.value.status = 'completed';
    selectedAction.value.completionRate = 100;
    editingCompletionRate.value = 100;
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

async function addFollowUp() {
  if (!selectedAction.value || !newFollowUpContent.value.trim()) return;
  try {
    const record = await addFollowUpApi({
      actionId: selectedAction.value.id,
      content: newFollowUpContent.value.trim(),
    });
    selectedAction.value.followUpRecords.push(record);
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

function goToTopics() {
  router.push('/my-topics');
}

// ==================== 监听器 ====================
watch([activeStatus, priorityFilter, topicFilter, searchKeyword, sortBy], () => {
  currentPage.value = 1;
  fetchActions();
});

watch([currentPage, pageSize], () => {
  fetchActions();
});

onMounted(() => {
  fetchActions();
});
</script>

<template>
  <Page
    title="我的行动"
    description="管理你的行动清单，追踪完成进度，记录跟进内容"
    content-class="p-6 bg-gray-50"
  >
    <!-- 主内容区 -->
    <div class="flex gap-6">
      <!-- 左侧行动列表 -->
      <div class="flex-1 space-y-6">
        <!-- 筛选和搜索 -->
        <ElCard shadow="hover" class="!rounded-xl">
          <div class="flex flex-wrap items-center gap-4">
            <!-- 状态 Tab -->
            <div class="flex gap-2">
              <button
                v-for="tab in statusTabs"
                :key="tab.id"
                class="px-4 py-2 rounded-full text-sm font-medium transition-all"
                :class="[
                  activeStatus === tab.id
                    ? 'bg-emerald-100 text-emerald-700 shadow-md border border-emerald-300 font-semibold'
                    : 'bg-gray-100 text-gray-700 hover:bg-emerald-50 hover:text-emerald-600 border border-gray-200'
                ]"
                @click="activeStatus = tab.id as ActionApi.ActionStatus | 'all'"
              >
                {{ tab.icon }} {{ tab.label }}
              </button>
            </div>
            <div class="flex-1" />
            <!-- 筛选项 -->
            <div class="flex items-center gap-3">
              <ElSelect v-model="priorityFilter" placeholder="优先级" clearable class="!w-40">
                <ElOption
                  v-for="opt in priorityOptions"
                  :key="opt.value"
                  :label="opt.label"
                  :value="opt.value"
                />
              </ElSelect>
              <ElSelect v-model="topicFilter" placeholder="选择课题" clearable class="!w-40">
                <ElOption
                  v-for="topic in topicOptions"
                  :key="topic.id"
                  :label="topic.title"
                  :value="topic.id"
                />
              </ElSelect>
              <ElInput
                v-model="searchKeyword"
                placeholder="搜索行动..."
                clearable
                class="!w-44"
                @keyup.enter="fetchActions"
              />
            </div>
          </div>
        </ElCard>

        <!-- 加载状态 -->
        <div v-if="loading" class="space-y-4">
          <ElSkeleton v-for="i in 3" :key="i" :rows="3" animated class="bg-white rounded-xl p-4" />
        </div>

        <!-- 空状态 -->
        <ElCard v-else-if="actions.length === 0" shadow="hover" class="!rounded-xl">
          <ElEmpty description="暂无行动，从课题分析中创建你的第一个行动吧！">
            <template #image>
              <div class="text-6xl">✨</div>
            </template>
            <ElButton
              type="primary"
              class="!bg-emerald-600 !border-emerald-600 hover:!bg-emerald-700 !rounded-full mt-4"
              @click="goToTopics"
            >
              去创建课题
            </ElButton>
          </ElEmpty>
        </ElCard>

        <!-- 行动列表 -->
        <div v-else class="space-y-4">
          <ElCard
            v-for="action in actions"
            :key="action.id"
            shadow="hover"
            class="!rounded-xl cursor-pointer hover:shadow-lg transition-all group"
            :class="{ 
              'opacity-70': action.status === 'completed',
              '!border-l-4 !border-l-red-400': action.priority === 'high' && action.status !== 'completed'
            }"
            @click="openActionDetail(action)"
          >
            <div class="flex gap-4">
              <!-- 左侧复选框和进度环 -->
              <div class="flex-shrink-0 relative">
                <div
                  class="w-14 h-14 cursor-pointer"
                  @click.stop="toggleActionStatus(action)"
                >
                  <ElProgress
                    type="circle"
                    :percentage="action.completionRate"
                    :width="56"
                    :stroke-width="4"
                    :color="getProgressColor(action)"
                  >
                    <template #default>
                      <span v-if="action.status === 'completed'" class="text-lg text-green-500">✓</span>
                      <span v-else class="text-sm font-medium">{{ action.completionRate }}%</span>
                    </template>
                  </ElProgress>
                </div>
              </div>

              <!-- 中间内容 -->
              <div class="flex-1 min-w-0">
                <div class="flex items-start justify-between gap-4 mb-2">
                  <div class="flex items-center gap-2">
                    <span
                      class="px-2 py-0.5 rounded text-xs font-medium"
                      :class="getPriorityStyle(action.priority)"
                    >
                      {{ getPriorityIcon(action.priority) }} {{ getPriorityText(action.priority) }}
                    </span>
                    <h3
                      class="text-base font-semibold text-gray-800 group-hover:text-emerald-600 transition-colors line-clamp-1"
                      :class="{ 'line-through text-gray-400': action.status === 'completed' }"
                    >
                      {{ action.title }}
                    </h3>
                  </div>
                  <div class="flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                    <ElButton
                      type="primary"
                      size="small"
                      class="!bg-emerald-600 !border-emerald-600 !rounded-full"
                      @click.stop="openFollowUpFromList(action)"
                    >
                      记录跟进
                    </ElButton>
                  </div>
                </div>

                <p class="text-gray-500 text-sm line-clamp-1 mb-2">
                  {{ action.description }}
                </p>

                <!-- 底部元信息 -->
                <div class="flex flex-wrap items-center gap-4 text-xs text-gray-400">
                  <span class="flex items-center gap-1 text-emerald-600">
                    📁 {{ action.topicTitle }}
                  </span>
                  <span v-if="action.followUpRecords.length > 0" class="flex items-center gap-1">
                    💬 {{ action.followUpRecords.length }} 条跟进
                  </span>
                  <span
                    v-if="action.dueDate"
                    class="flex items-center gap-1"
                    :class="{
                      'text-amber-500 font-medium': isDueSoon(action.dueDate) && action.status !== 'completed',
                      'text-red-500 font-medium': isOverdue(action.dueDate) && action.status !== 'completed'
                    }"
                  >
                    ⏰ {{ getDueText(action.dueDate) }}
                  </span>
                  <span class="text-gray-300">|</span>
                  <span>{{ formatTime(action.createdAt) }}</span>
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
        <!-- 执行指南 -->
        <ElCard shadow="hover" class="!rounded-xl">
          <template #header>
            <span class="font-semibold text-gray-700">执行指南</span>
          </template>
          <div class="space-y-4">
            <div class="flex items-start gap-3">
              <div class="w-8 h-8 rounded-full bg-red-100 text-red-600 flex items-center justify-center font-bold text-sm flex-shrink-0">1</div>
              <div>
                <div class="font-medium text-gray-700 text-sm">优先高紧急</div>
                <div class="text-xs text-gray-500">先处理高优先级和即将到期的行动</div>
              </div>
            </div>
            <div class="flex items-start gap-3">
              <div class="w-8 h-8 rounded-full bg-amber-100 text-amber-600 flex items-center justify-center font-bold text-sm flex-shrink-0">2</div>
              <div>
                <div class="font-medium text-gray-700 text-sm">分解大任务</div>
                <div class="text-xs text-gray-500">把复杂行动拆成小步骤逐一完成</div>
              </div>
            </div>
            <div class="flex items-start gap-3">
              <div class="w-8 h-8 rounded-full bg-emerald-100 text-emerald-600 flex items-center justify-center font-bold text-sm flex-shrink-0">3</div>
              <div>
                <div class="font-medium text-gray-700 text-sm">定期跟进</div>
                <div class="text-xs text-gray-500">记录进展和遇到的问题</div>
              </div>
            </div>
            <div class="flex items-start gap-3">
              <div class="w-8 h-8 rounded-full bg-blue-100 text-blue-600 flex items-center justify-center font-bold text-sm flex-shrink-0">4</div>
              <div>
                <div class="font-medium text-gray-700 text-sm">及时完成</div>
                <div class="text-xs text-gray-500">完成后标记状态，保持清单清晰</div>
              </div>
            </div>
          </div>
        </ElCard>

        <!-- 快速筛选 -->
        <ElCard shadow="hover" class="!rounded-xl">
          <template #header>
            <span class="font-semibold text-gray-700">快速筛选</span>
          </template>
          <div class="space-y-2">
            <button
              class="w-full flex items-center justify-between p-3 rounded-lg hover:bg-gray-50 transition-colors text-left"
              @click="priorityFilter = 'high'; activeStatus = 'all'"
            >
              <span class="flex items-center gap-2 text-sm">🔴 高优先级</span>
              <span class="text-xs text-red-500 font-medium">{{ actions.filter(a => a.priority === 'high').length }}</span>
            </button>
            <button
              class="w-full flex items-center justify-between p-3 rounded-lg hover:bg-gray-50 transition-colors text-left"
              @click="activeStatus = 'in_progress'; priorityFilter = 'all'"
            >
              <span class="flex items-center gap-2 text-sm">🚀 进行中</span>
              <span class="text-xs text-amber-500 font-medium">{{ stats.inProgress }}</span>
            </button>
            <button
              class="w-full flex items-center justify-between p-3 rounded-lg hover:bg-gray-50 transition-colors text-left"
              @click="sortBy = 'dueDate'; activeStatus = 'all'; priorityFilter = 'all'"
            >
              <span class="flex items-center gap-2 text-sm">⏰ 即将到期</span>
              <span class="text-xs text-orange-500 font-medium">{{ actions.filter(a => a.dueDate && isDueSoon(a.dueDate) && a.status !== 'completed').length }}</span>
            </button>
          </div>
        </ElCard>

        <!-- 执行小贴士 -->
        <ElCard shadow="hover" class="!rounded-xl !bg-gradient-to-br from-emerald-50 to-emerald-100 !border-emerald-200">
          <template #header>
            <span class="font-semibold text-emerald-700">执行小贴士</span>
          </template>
          <ul class="text-sm text-emerald-800 space-y-2">
            <li class="flex items-start gap-2">
              <span class="text-emerald-500">•</span>
              每天选择1-3个最重要的行动
            </li>
            <li class="flex items-start gap-2">
              <span class="text-emerald-500">•</span>
              完成一个再开始下一个
            </li>
            <li class="flex items-start gap-2">
              <span class="text-emerald-500">•</span>
              遇阻时记录问题并寻求帮助
            </li>
            <li class="flex items-start gap-2">
              <span class="text-emerald-500">•</span>
              定期回顾复盘提升效率
            </li>
          </ul>
        </ElCard>

        <!-- 关联课题 -->
        <ElCard shadow="hover" class="!rounded-xl">
          <template #header>
            <span class="font-semibold text-gray-700">关联课题</span>
          </template>
          <div class="space-y-2">
            <div
              v-for="topic in topicOptions"
              :key="topic.id"
              class="flex items-center justify-between p-2 rounded-lg hover:bg-gray-50 cursor-pointer transition-colors"
              @click="topicFilter = topic.id"
            >
              <span class="text-sm text-gray-700 line-clamp-1">{{ topic.title }}</span>
              <span class="text-xs text-gray-400">{{ actions.filter(a => a.topicId === topic.id).length }}</span>
            </div>
          </div>
        </ElCard>
      </div>
    </div>

    <!-- 行动详情弹窗 -->
    <ElDialog
      v-model="detailDialogVisible"
      title="行动详情"
      width="700px"
    >
      <div v-if="selectedAction" class="space-y-6">
        <!-- 头部 -->
        <div class="p-4 bg-gray-50 rounded-lg">
          <div class="flex items-start justify-between mb-3">
            <div class="flex items-center gap-2">
              <span
                class="px-2 py-1 rounded text-xs font-medium"
                :class="getPriorityStyle(selectedAction.priority)"
              >
                {{ getPriorityIcon(selectedAction.priority) }} {{ getPriorityText(selectedAction.priority) }}优先级
              </span>
              <span
                class="px-2 py-1 rounded-full text-xs"
                :class="getStatusStyle(selectedAction.status)"
              >
                {{ getStatusText(selectedAction.status) }}
              </span>
            </div>
          </div>
          <h2 class="text-lg font-semibold text-gray-800 mb-2">{{ selectedAction.title }}</h2>
          <p class="text-gray-600 text-sm">{{ selectedAction.description }}</p>
          <div class="flex items-center gap-2 mt-3 text-sm text-emerald-600">
            <span>📁</span>
            <span>关联课题：{{ selectedAction.topicTitle }}</span>
          </div>
        </div>

        <!-- 指导原则 -->
        <div v-if="selectedAction.guidingPrinciple">
          <h3 class="font-medium text-gray-700 mb-2 flex items-center gap-2">
            <span>💡</span> 指导原则
          </h3>
          <div class="p-3 bg-amber-50 rounded-lg text-amber-800 text-sm">
            {{ selectedAction.guidingPrinciple }}
          </div>
        </div>

        <!-- 完成度 -->
        <div>
          <h3 class="font-medium text-gray-700 mb-3 flex items-center gap-2">
            <span>📊</span> 完成度
          </h3>
          <div class="flex items-center gap-4">
            <ElSlider v-model="editingCompletionRate" :max="100" :step="5" class="flex-1" />
            <span class="text-lg font-bold text-emerald-600 w-16 text-right">{{ editingCompletionRate }}%</span>
          </div>
          <div class="flex items-center gap-3 mt-4">
            <ElButton type="primary" class="!bg-emerald-600 !border-emerald-600" @click="saveCompletionRate">
              保存进度
            </ElButton>
            <ElButton
              v-if="editingCompletionRate === 100 && selectedAction.status !== 'completed'"
              type="success"
              @click="markAsCompleted"
            >
              标记为已完成
            </ElButton>
          </div>
        </div>

        <!-- 跟进记录 -->
        <div>
          <h3 class="font-medium text-gray-700 mb-3 flex items-center gap-2">
            <span>💬</span> 跟进记录
            <span class="text-gray-400 text-sm">({{ selectedAction.followUpRecords.length }})</span>
          </h3>
          <div v-if="latestFollowUp" class="p-3 bg-gray-50 rounded-lg mb-3">
            <div class="text-xs text-gray-500 mb-1">最新跟进 · {{ formatDateTime(latestFollowUp.createdAt) }}</div>
            <div class="text-sm text-gray-700">{{ latestFollowUp.content }}</div>
          </div>
          <ElButton type="primary" class="!bg-emerald-600 !border-emerald-600" @click="openFollowUpDrawer">
            {{ selectedAction.followUpRecords.length > 0 ? '查看全部跟进' : '添加跟进记录' }}
          </ElButton>
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
    >
      <div v-if="selectedAction" class="space-y-6">
        <!-- 行动简要信息 -->
        <div class="p-4 bg-gray-50 rounded-lg">
          <div class="flex items-center gap-2 mb-2">
            <span
              class="px-2 py-0.5 rounded text-xs"
              :class="getPriorityStyle(selectedAction.priority)"
            >
              {{ getPriorityText(selectedAction.priority) }}
            </span>
            <span
              class="px-2 py-0.5 rounded-full text-xs"
              :class="getStatusStyle(selectedAction.status)"
            >
              {{ getStatusText(selectedAction.status) }}
            </span>
          </div>
          <h3 class="font-medium text-gray-800">{{ selectedAction.title }}</h3>
          <div class="flex items-center gap-4 mt-2 text-sm text-gray-500">
            <span>完成度: {{ selectedAction.completionRate }}%</span>
          </div>
        </div>

        <!-- 添加跟进 -->
        <div>
          <h4 class="font-medium text-gray-700 mb-2">添加跟进</h4>
          <ElInput
            v-model="newFollowUpContent"
            type="textarea"
            :rows="4"
            placeholder="记录你的跟进内容、遇到的问题、下一步计划..."
            resize="none"
          />
          <ElButton
            type="primary"
            class="w-full mt-3 !bg-emerald-600 !border-emerald-600"
            :disabled="!newFollowUpContent.trim()"
            @click="addFollowUp"
          >
            添加跟进
          </ElButton>
        </div>

        <!-- 历史记录 -->
        <div>
          <h4 class="font-medium text-gray-700 mb-3">
            历史记录
            <span class="text-gray-400 text-sm">({{ selectedAction.followUpRecords.length }})</span>
          </h4>
          <div v-if="sortedFollowUpRecords.length > 0" class="space-y-4">
            <div
              v-for="(record, index) in sortedFollowUpRecords"
              :key="record.id"
              class="relative pl-6 pb-4 border-l-2 border-emerald-200 last:border-transparent"
            >
              <div class="absolute left-0 top-0 w-3 h-3 rounded-full bg-emerald-500 -translate-x-[7px]" />
              <div class="text-xs text-gray-500 mb-1">
                #{{ selectedAction.followUpRecords.length - index }} · {{ formatDateTime(record.createdAt) }}
              </div>
              <div class="text-sm text-gray-700 bg-gray-50 rounded-lg p-3">
                {{ record.content }}
              </div>
            </div>
          </div>
          <div v-else class="text-center py-8 text-gray-400">
            <div class="text-4xl mb-2">📝</div>
            <p>暂无跟进记录</p>
            <p class="text-sm">添加第一条记录开始追踪进度</p>
          </div>
        </div>
      </div>
    </ElDrawer>
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
</style>
