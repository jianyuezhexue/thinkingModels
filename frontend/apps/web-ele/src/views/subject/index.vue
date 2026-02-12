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
} from 'element-plus';

import { getMyTopicsApi, deleteTopicApi, type TopicApi } from '#/api';

// 路由
const router = useRouter();

// 加载状态
const loading = ref(false);

// 课题列表
const topics = ref<TopicApi.Topic[]>([]);
const total = ref(0);

// 分页
const currentPage = ref(1);
const pageSize = ref(10);

// 筛选
const searchKeyword = ref('');
const statusFilter = ref<TopicApi.TopicStatus | 'all'>('all');

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
    topics.value = res.list;
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

<template>
  <Page
    description="管理你的思考课题，使用思维模型进行深入分析"
    title="我的课题"
  >
    <ElCard shadow="never">
      <template #header>
        <div class="flex flex-wrap items-center justify-between gap-4">
          <div class="flex items-center gap-4">
            <h2 class="text-lg font-semibold">课题列表</h2>
            <ElSelect v-model="statusFilter" style="width: 140px">
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
              @keyup.enter="fetchTopics"
            >
              <template #prefix>
                <svg class="h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
                </svg>
              </template>
            </ElInput>
            <ElButton type="primary" @click="goToCreateTopic">
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
        <ElButton type="primary" @click="goToCreateTopic">创建课题</ElButton>
      </ElEmpty>

      <!-- 课题列表 -->
      <div v-else class="space-y-4">
        <div
          v-for="topic in topics"
          :key="topic.id"
          class="group rounded-lg border border-gray-200 bg-white p-4 transition-all hover:border-purple-300 hover:shadow-md cursor-pointer"
          @click="goToDetail(topic)"
        >
          <div class="flex items-start justify-between gap-4">
            <div class="flex-1 min-w-0">
              <!-- 标题和状态 -->
              <div class="flex items-center gap-3 mb-2">
                <h3 class="font-semibold text-gray-900 group-hover:text-purple-600 transition-colors truncate">
                  {{ topic.title }}
                </h3>
                <ElTag :type="getStatusType(topic.status)" size="small">
                  {{ getStatusText(topic.status) }}
                </ElTag>
              </div>

              <!-- 描述 -->
              <p class="text-sm text-gray-500 line-clamp-2 mb-3">{{ topic.description }}</p>

              <!-- 元信息 -->
              <div class="flex flex-wrap items-center gap-4 text-xs text-gray-400">
                <span class="flex items-center gap-1">
                  <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"/>
                  </svg>
                  {{ topic.modelName || '未选用模型' }}
                </span>
                <span class="flex items-center gap-1">
                  <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
                  </svg>
                  {{ topic.analysisCount }} 次分析
                </span>
                <span class="flex items-center gap-1">
                  <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                  </svg>
                  创建于 {{ formatDate(topic.createdAt) }}
                </span>
                <span class="flex items-center gap-1">
                  <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                  </svg>
                  更新于 {{ formatDate(topic.updatedAt) }}
                </span>
              </div>
            </div>

            <!-- 操作按钮 -->
            <div class="flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
              <ElButton
                type="primary"
                size="small"
                @click.stop="goToDetail(topic)"
              >
                查看详情
              </ElButton>
              <ElButton
                type="danger"
                size="small"
                plain
                @click.stop="handleDelete(topic)"
              >
                删除
              </ElButton>
            </div>
          </div>
        </div>

        <!-- 分页 -->
        <div class="flex justify-center pt-4">
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
  </Page>
</template>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
