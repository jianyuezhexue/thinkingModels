<template>
  <Page
    description="管理你的思考课题，运用思维模型深入剖析问题本质，并导向正确的行动决策"
    title="我的课题"
  >
    <!-- 课题列表 -->
    <ElCard shadow="never" class="topic-list-card">
      <template #header>
        <div class="flex flex-wrap items-center justify-between gap-4">
          <div class="flex items-center gap-4">
            <h2 class="text-lg font-semibold">课题列表</h2>
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
            <ElInput
              v-model="searchKeyword"
              placeholder="搜索课题..."
              clearable
              style="width: 220px"
              size="small"
              @keyup.enter="fetchTopics"
            >
              <template #prefix>
                <svg class="h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
                </svg>
              </template>
            </ElInput>
            <ElButton type="primary" @click="goToCreateTopic" class="btn-primary-gradient">
              <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
              </svg>
              创建课题
            </ElButton>
          </div>
        </div>
      </template>

      <!-- 加载状态 -->
      <div v-if="loading" class="space-y-4">
        <ElSkeleton v-for="i in 3" :key="i" animated>
          <template #template>
            <div class="flex items-center gap-4">
              <ElSkeletonItem variant="text" style="width: 200px" />
              <ElSkeletonItem variant="text" style="width: 100px" />
              <ElSkeletonItem variant="text" style="width: 80px" />
              <ElSkeletonItem variant="text" style="width: 120px" />
            </div>
          </template>
        </ElSkeleton>
      </div>

      <!-- 空状态 -->
      <ElEmpty v-else-if="topics.length === 0" description="暂无课题，创建一个开始思考吧！">
        <template #image>
          <div class="empty-illustration">
            <svg class="h-24 w-24 text-purple-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"/>
            </svg>
          </div>
        </template>
        <ElButton type="primary" @click="goToCreateTopic" class="btn-primary-gradient mt-4">
          创建第一个课题
        </ElButton>
      </ElEmpty>

      <!-- 课题列表 -->
      <div v-else class="topic-list">
        <div
          v-for="topic in topics"
          :key="topic.id"
          class="topic-card group"
          :class="{ 'topic-has-action': topic.recommendedActions && topic.recommendedActions.length > 0 }"
          @click="goToDetail(topic)"
        >
          <!-- 左侧：课题信息 -->
          <div class="topic-main">
            <div class="topic-header">
              <h3 class="topic-title">{{ topic.title }}</h3>
              <div class="topic-badges">
                <ElTag :type="getStatusType(topic.status)" size="small" effect="light">
                  {{ getStatusText(topic.status) }}
                </ElTag>
                <ElTag
                  v-if="topic.recommendedActions && topic.recommendedActions.length > 0"
                  type="success"
                  size="small"
                  effect="light"
                  class="action-badge"
                >
                  <svg class="h-3 w-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                  </svg>
                  {{ topic.recommendedActions.length }} 个行动建议
                </ElTag>
              </div>
            </div>

            <p class="topic-desc">{{ topic.description }}</p>

            <!-- 行动建议预览（如果有） -->
            <div v-if="topic.recommendedActions && topic.recommendedActions.length > 0" class="action-preview">
              <span class="action-label-inline">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
                </svg>
                行动：
              </span>
              <span class="action-item-inline">{{ topic.recommendedActions[0] }}</span>
              <span v-if="topic.recommendedActions.length > 1" class="action-more-inline">
                +{{ topic.recommendedActions.length - 1 }}
              </span>
            </div>

            <!-- 元信息 -->
            <div class="topic-meta">
              <span class="meta-item models-wrapper">
                <svg class="h-4 w-4 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"/>
                </svg>
                <template v-if="topic.selectedModels && topic.selectedModels.length > 0">
                  <span class="models-list">{{ topic.selectedModels.join('、') }}</span>
                </template>
                <template v-else>未选用模型</template>
              </span>
              <span class="meta-item">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
                </svg>
                {{ topic.analysisCount }} 次分析
              </span>
              <span class="meta-item">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                </svg>
                {{ formatDate(topic.createdAt) }}
              </span>
            </div>
          </div>

          <!-- 右侧：操作区 -->
          <div class="topic-actions">
            <div class="action-buttons">
              <ElButton
                v-if="topic.recommendedActions && topic.recommendedActions.length > 0"
                type="success"
                size="small"
                class="btn-action"
                @click.stop="viewActions(topic)"
              >
                <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
                </svg>
                查看行动
              </ElButton>
              <ElButton
                v-else-if="topic.status === 'draft'"
                type="primary"
                size="small"
                class="btn-analyze"
                @click.stop="startAnalysis(topic)"
              >
                <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"/>
                </svg>
                开始分析
              </ElButton>
              <ElButton
                v-else
                type="primary"
                size="small"
                plain
                @click.stop="goToDetail(topic)"
              >
                查看详情
              </ElButton>
            </div>

            <!-- 进度指示器 -->
            <div class="topic-progress">
              <div class="progress-ring" :class="`progress-${topic.status}`">
                <svg viewBox="0 0 36 36">
                  <path
                    class="progress-ring-bg"
                    d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
                  />
                  <path
                    class="progress-ring-fill"
                    :stroke-dasharray="`${getProgress(topic)}, 100`"
                    d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
                  />
                </svg>
                <div class="progress-text">{{ getProgress(topic) }}%</div>
              </div>
              <div class="progress-label">{{ getProgressLabel(topic) }}</div>
            </div>
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

    <!-- 行动建议弹窗 -->
    <ElDialog
      v-model="actionDialogVisible"
      title="行动建议"
      width="600px"
      class="action-dialog"
    >
      <div v-if="selectedTopic" class="action-dialog-content">
        <div class="dialog-topic-title">{{ selectedTopic.title }}</div>
        <div class="action-list">
          <div
            v-for="(action, index) in selectedTopic.recommendedActions"
            :key="index"
            class="action-list-item"
            :class="{ 'action-completed': completedActionsMap[`${selectedTopic.id}-${index}`] }"
          >
            <div class="action-checkbox" @click="toggleAction(selectedTopic.id, index)">
              <svg v-if="completedActionsMap[`${selectedTopic.id}-${index}`]" class="h-5 w-5 text-emerald-500" fill="currentColor" viewBox="0 0 24 24">
                <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/>
              </svg>
              <svg v-else class="h-5 w-5 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <circle cx="12" cy="12" r="10" stroke-width="2"/>
              </svg>
            </div>
            <div class="action-text">{{ action }}</div>
          </div>
        </div>
      </div>
      <template #footer>
        <ElButton @click="actionDialogVisible = false">关闭</ElButton>
        <ElButton type="primary" @click="exportActions">导出行动清单</ElButton>
      </template>
    </ElDialog>
  </Page>
</template>

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
  ElTag,
  ElMessage,
  ElEmpty,
  ElSkeleton,
  ElSkeletonItem,
  ElPagination,
  ElDialog,
} from 'element-plus';

import { getMyTopicsApi, deleteTopicApi, type TopicApi } from '#/api';

// 路由
const router = useRouter();

// 加载状态
const loading = ref(false);

// 课题列表
const topics = ref<(TopicApi.Topic & { recommendedActions?: string[]; selectedModels?: string[] })[]>([]);
const total = ref(0);

// 分页
const currentPage = ref(1);
const pageSize = ref(10);

// 筛选
const searchKeyword = ref('');
const statusFilter = ref<TopicApi.TopicStatus | 'all'>('all');

// 弹窗状态
const actionDialogVisible = ref(false);
const selectedTopic = ref<(TopicApi.Topic & { recommendedActions?: string[] }) | null>(null);
const completedActionsMap = ref<Record<string, boolean>>({});

// 状态选项
const statusOptions = [
  { value: 'all', label: '全部状态' },
  { value: 'draft', label: '草稿', type: 'info' },
  { value: 'in_progress', label: '进行中', type: 'warning' },
  { value: 'completed', label: '已完成', type: 'success' },
  { value: 'archived', label: '已归档', type: '' },
];

// 获取课题列表
async function fetchTopics() {
  loading.value = true;
  try {
    const params: TopicApi.TopicListParams = {
      page: currentPage.value,
      pageSize: pageSize.value,
      keyword: searchKeyword.value || undefined,
    };

    if (statusFilter.value !== 'all') {
      params.status = statusFilter.value;
    }

    const res = await getMyTopicsApi(params);

    // 模拟模型列表
    const modelNames = ['SWOT分析', '5W1H', 'MECE原则', '第一性原理', '金字塔原理', '逆向思维', '奥卡姆剃刀', '二阶思维'];

    // 模拟添加行动建议数据和选用模型
    topics.value = res.list.map((topic: TopicApi.Topic, index: number) => ({
      ...topic,
      recommendedActions: index % 3 === 0 ? [
        '重新评估目标用户群体，缩小范围至核心用户',
        '制定3个月内可执行的MVP功能清单',
        '寻找2-3位潜在用户进行深度访谈',
      ] : undefined,
      // 随机选用1-3个模型
      selectedModels: modelNames.slice(index % 4, index % 4 + 1 + (index % 3)),
    }));

    total.value = res.total;
  } catch (error) {
    console.error('获取课题列表失败:', error);
    ElMessage.error('获取课题列表失败');
  } finally {
    loading.value = false;
  }
}

// 监听筛选条件变化
watch([searchKeyword, statusFilter], () => {
  currentPage.value = 1;
  fetchTopics();
});

// 监听分页变化
watch([currentPage, pageSize], () => {
  fetchTopics();
});

// 页面加载时获取数据
onMounted(() => {
  fetchTopics();
});

// 跳转到创建课题页面
function goToCreateTopic() {
  router.push('/my-topics/create');
}

// 跳转到课题详情
function goToDetail(topic: TopicApi.Topic) {
  router.push(`/my-topics/${topic.id}`);
}

// 开始分析
function startAnalysis(topic: TopicApi.Topic) {
  router.push(`/my-topics/${topic.id}?tab=analysis`);
}

// 查看行动建议
function viewActions(topic: TopicApi.Topic & { recommendedActions?: string[] }) {
  selectedTopic.value = topic;
  actionDialogVisible.value = true;
}

// 切换行动完成状态
function toggleAction(topicId: string, actionIndex: number) {
  const key = `${topicId}-${actionIndex}`;
  completedActionsMap.value[key] = !completedActionsMap.value[key];
}

// 导出行动清单
function exportActions() {
  ElMessage.success('行动清单已导出');
  actionDialogVisible.value = false;
}

// 获取进度
function getProgress(topic: TopicApi.Topic): number {
  const map: Record<string, number> = {
    draft: 25,
    in_progress: 60,
    completed: 100,
    archived: 100,
  };
  return map[topic.status] || 0;
}

// 获取进度标签
function getProgressLabel(topic: TopicApi.Topic): string {
  const map: Record<string, string> = {
    draft: '待分析',
    in_progress: '分析中',
    completed: '已完成',
    archived: '已归档',
  };
  return map[topic.status] || '未知';
}

// 删除课题
async function handleDelete(topic: TopicApi.Topic) {
  try {
    await deleteTopicApi(topic.id);
    ElMessage.success('课题已删除');
    fetchTopics();
  } catch (error) {
    console.error('删除课题失败:', error);
    ElMessage.error('删除失败');
  }
}

// 获取状态标签类型
function getStatusType(status: TopicApi.TopicStatus): any {
  const map: Record<TopicApi.TopicStatus, any> = {
    draft: 'info',
    in_progress: 'warning',
    completed: 'success',
    archived: '',
  };
  return map[status] || '';
}

// 获取状态标签文本
function getStatusText(status: TopicApi.TopicStatus): string {
  const map: Record<TopicApi.TopicStatus, string> = {
    draft: '草稿',
    in_progress: '进行中',
    completed: '已完成',
    archived: '已归档',
  };
  return map[status] || status;
}

// 格式化日期
function formatDate(dateStr: string): string {
  if (!dateStr) return '-';
  const date = new Date(dateStr);
  return date.toLocaleDateString('zh-CN');
}
</script>

<style scoped>
/* 工作流程卡片 */
.workflow-card {
  background: linear-gradient(135deg, #f5f3ff 0%, #ede9fe 100%);
  border: 1px solid #ddd6fe;
}

.workflow-icon {
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, #6366f1 0%, #818cf8 100%);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 工作流步骤 */
.workflow-steps {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 0;
  gap: 16px;
}

.step-item {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
  padding: 16px;
  background: white;
  border-radius: 12px;
  border: 2px solid #e5e7eb;
  transition: all 0.3s ease;
}

.step-item.step-active {
  border-color: #6366f1;
  background: linear-gradient(135deg, #fff 0%, #f5f3ff 100%);
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.15);
}

.step-number {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #e5e7eb;
  color: #6b7280;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 14px;
  transition: all 0.3s ease;
}

.step-active .step-number {
  background: linear-gradient(135deg, #6366f1 0%, #818cf8 100%);
  color: white;
}

.step-number-action {
  background: linear-gradient(135deg, #10b981 0%, #34d399 100%) !important;
  color: white !important;
}

.step-content {
  flex: 1;
}

.step-title {
  font-weight: 600;
  font-size: 14px;
  color: #374151;
  margin-bottom: 2px;
}

.step-active .step-title {
  color: #6366f1;
}

.step-desc {
  font-size: 12px;
  color: #9ca3af;
}

.step-arrow {
  color: #d1d5db;
}

.step-badge {
  color: #10b981;
}

/* 工作流进度条 */
.workflow-progress {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px dashed #ddd6fe;
}

.progress-track {
  height: 6px;
  background: #e5e7eb;
  border-radius: 3px;
  overflow: hidden;
  margin-bottom: 12px;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #6366f1 0%, #818cf8 50%, #10b981 100%);
  border-radius: 3px;
  transition: width 0.5s ease;
}

.progress-labels {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #9ca3af;
}

.progress-labels .text-active {
  color: #6366f1;
  font-weight: 500;
}

.progress-labels .text-action {
  color: #10b981;
}

/* 统计卡片 */
.stat-card {
  background: white;
  border-radius: 12px;
  padding: 16px;
  border: 1px solid #e5e7eb;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
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
  margin-bottom: 4px;
}

.stat-trend {
  font-size: 11px;
  color: #10b981;
  font-weight: 500;
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

/* 课题列表卡片 */
.topic-list-card {
  width: 100%;
}

.topic-list-card :deep(.el-card__body) {
  padding: 20px;
}

/* 课题列表 */
.topic-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  width: 100%;
}

.topic-card {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  padding: 20px 24px;
  background: white;
  border-radius: 12px;
  border: 1px solid #e5e7eb;
  cursor: pointer;
  transition: all 0.3s ease;
}

.topic-card:hover {
  border-color: #c4b5fd;
  box-shadow: 0 4px 16px rgba(99, 102, 241, 0.12);
}

.topic-has-action {
  border-left: 4px solid #10b981;
}

.topic-main {
  flex: 1 1 auto;
  min-width: 0;
  width: calc(100% - 120px);
}

.topic-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.topic-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.topic-card:hover .topic-title {
  color: #6366f1;
}

.topic-badges {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.action-badge {
  background: #ecfdf5 !important;
  color: #059669 !important;
  border-color: #6ee7b7 !important;
  white-space: nowrap !important;
  flex-shrink: 0;
}

.action-badge :deep(.el-tag__content) {
  display: inline-flex !important;
  align-items: center;
  white-space: nowrap !important;
}

.topic-desc {
  font-size: 13px;
  color: #6b7280;
  margin-bottom: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

/* 行动建议预览 - 单行内联布局 */
.action-preview {
  display: block;
  padding: 8px 14px;
  background: linear-gradient(135deg, #ecfdf5 0%, #d1fae5 100%);
  border-radius: 8px;
  margin-bottom: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  line-height: 1.8;
}

.action-label-inline {
  display: inline;
  font-size: 13px;
  color: #059669;
  font-weight: 600;
  margin-right: 8px;
}

.action-label-inline svg {
  display: inline-block;
  vertical-align: middle;
  margin-right: 4px;
  margin-top: -2px;
}

.action-item-inline {
  display: inline;
  font-size: 13px;
  color: #374151;
  background: white;
  padding: 3px 10px;
  border-radius: 4px;
  border: 1px solid #a7f3d0;
  margin-right: 8px;
}

.action-more-inline {
  display: inline;
  font-size: 12px;
  color: #059669;
  font-weight: 500;
}

/* 元信息 */
.topic-meta {
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

.meta-item.models-wrapper {
  max-width: 300px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.models-list {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: #6366f1;
  font-weight: 500;
}

/* 操作区 */
.topic-actions {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  width: 100px;
  flex: 0 0 100px;
}

.action-buttons {
  display: flex;
  gap: 8px;
}

.btn-action {
  background: linear-gradient(135deg, #10b981 0%, #34d399 100%);
  border: none;
  color: white;
}

.btn-action:hover {
  background: linear-gradient(135deg, #059669 0%, #10b981 100%);
  color: white;
}

.btn-analyze {
  background: linear-gradient(135deg, #6366f1 0%, #818cf8 100%);
  border: none;
  color: white;
}

.btn-analyze:hover {
  background: linear-gradient(135deg, #4f46e5 0%, #6366f1 100%);
  color: white;
}

/* 进度环 */
.topic-progress {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
}

.progress-ring {
  position: relative;
  width: 48px;
  height: 48px;
}

.progress-ring svg {
  width: 100%;
  height: 100%;
  transform: rotate(-90deg);
}

.progress-ring-bg {
  fill: none;
  stroke: #e5e7eb;
  stroke-width: 3;
}

.progress-ring-fill {
  fill: none;
  stroke-width: 3;
  stroke-linecap: round;
  transition: stroke-dasharray 0.5s ease;
}

.progress-draft .progress-ring-fill {
  stroke: #9ca3af;
}

.progress-in_progress .progress-ring-fill {
  stroke: #f59e0b;
}

.progress-completed .progress-ring-fill {
  stroke: #10b981;
}

.progress-archived .progress-ring-fill {
  stroke: #6b7280;
}

.progress-text {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 11px;
  font-weight: 600;
  color: #374151;
}

.progress-label {
  font-size: 11px;
  color: #9ca3af;
}

/* 空状态 */
.empty-illustration {
  display: flex;
  justify-content: center;
  margin-bottom: 16px;
}

/* 弹窗样式 */
:deep(.action-dialog .el-dialog__header) {
  background: linear-gradient(135deg, #ecfdf5 0%, #d1fae5 100%);
  margin-right: 0;
  padding: 20px;
  border-bottom: 1px solid #a7f3d0;
}

:deep(.action-dialog .el-dialog__title) {
  color: #059669;
  font-weight: 600;
}

.action-dialog-content {
  padding: 8px 0;
}

.dialog-topic-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 1px dashed #e5e7eb;
}

.action-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.action-list-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px;
  background: #f9fafb;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.action-list-item:hover {
  background: #f3f4f6;
}

.action-completed {
  background: #ecfdf5 !important;
  opacity: 0.7;
}

.action-completed .action-text {
  text-decoration: line-through;
  color: #9ca3af;
}

.action-checkbox {
  cursor: pointer;
  flex-shrink: 0;
  margin-top: 2px;
}

.action-text {
  font-size: 14px;
  color: #374151;
  line-height: 1.5;
}

/* 响应式 */
@media (max-width: 768px) {
  .workflow-steps {
    flex-direction: column;
    gap: 12px;
  }

  .step-arrow {
    transform: rotate(90deg);
  }

  .topic-card {
    flex-direction: column;
  }

  .topic-actions {
    flex-direction: row;
    justify-content: space-between;
    width: 100%;
  }
}
</style>
