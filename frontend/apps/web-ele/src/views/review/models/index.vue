<script lang="ts" setup>
import { onMounted, onUnmounted, ref, reactive, computed } from 'vue';

import { Page } from '@vben/common-ui';

import {
  ElButton,
  ElCard,
  ElTable,
  ElTableColumn,
  ElTag,
  ElMessage,
  ElMessageBox,
  ElDialog,
  ElForm,
  ElFormItem,
  ElInput,
  ElSelect,
  ElOption,
  ElPagination,
  ElDescriptions,
  ElDescriptionsItem,
  ElTabs,
  ElTabPane,
  ElBadge,
  ElAvatar,
  ElEmpty,
} from 'element-plus';

import {
  getThinkingModelListApi,
  reviewThinkingModelApi,
  getThinkingModelDetailApi,
} from '#/api/thinking/model';
import { getAllCategoriesApi } from '#/api/master/category';

// ===================== 类型定义 =====================

interface ModelInfo {
  id: number;
  name: string;
  description: string;
  coverImage: string;
  icon: string;
  categoryId: number;
  categoryName?: string;
  price: number;
  isFree: boolean;
  overview: string;
  difficulty: number;
  estimatedTime: number;
  status: number;
  version: string;
  isOfficial: boolean;
  author: {
    id: string;
    name: string;
    avatar: string;
  };
  stats: {
    usageCount: number;
    adoptCount: number;
    likeCount: number;
    commentCount: number;
  };
  tags: string[] | null;
  createdAt: string;
  updatedAt: string;
  submitTime?: string;
  reviewNote?: string;
  reviewerName?: string;
  reviewTime?: string;
}

// ===================== 状态定义 =====================

// 请求取消控制器（用于 loadStatusCounts）
let statusCountsAbortController: AbortController | null = null;
// 请求取消控制器（用于 loadModels）
let modelsAbortController: AbortController | null = null;

// 加载状态
const loading = ref(false);

// 当前标签
const activeTab = ref('pending');

// 模型列表
const modelList = ref<ModelInfo[]>([]);
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
});

// 搜索表单
const searchForm = reactive({
  keyword: '',
  categoryId: '',
});

// 详情弹窗
const detailDialogVisible = ref(false);
const detailLoading = ref(false);
const currentModel = ref<ModelInfo | null>(null);

// 审核弹窗
const reviewDialogVisible = ref(false);
const reviewForm = reactive({
  approved: true,
  note: '',
});
const reviewLoading = ref(false);

// 分类选项
const categoryOptions = ref<{ value: number; label: string }[]>([]);

// 状态配置
const statusConfig: Record<number, { label: string; type: 'primary' | 'success' | 'warning' | 'info' | 'danger' }> = {
  0: { label: '草稿', type: 'info' },
  1: { label: '已发布', type: 'success' },
  2: { label: '已下架', type: 'warning' },
  3: { label: '审核中', type: 'primary' },
  4: { label: '已驳回', type: 'danger' },
};

// 难度配置
const difficultyConfig: Record<number, { label: string; color: string }> = {
  1: { label: '简单', color: '#67C23A' },
  2: { label: '中等', color: '#E6A23C' },
  3: { label: '困难', color: '#F56C6C' },
};

// 各状态的数量（从后端获取）
const statusTotalCounts = reactive({
  pending: 0,
  approved: 0,
  rejected: 0,
});

// 计算各状态数量
const statusCounts = computed(() => {
  return {
    pending: activeTab.value === 'pending' ? pagination.total : statusTotalCounts.pending,
    approved: activeTab.value === 'approved' ? pagination.total : statusTotalCounts.approved,
    rejected: activeTab.value === 'rejected' ? pagination.total : statusTotalCounts.rejected,
  };
});

// 加载各状态的数量统计
async function loadStatusCounts() {
  // 取消之前的请求
  statusCountsAbortController?.abort();
  statusCountsAbortController = new AbortController();
  const signal = statusCountsAbortController.signal;

  try {
    // 并行请求各状态的模型数量
    const [pendingRes, approvedRes, rejectedRes] = await Promise.all([
      getThinkingModelListApi({ page: 1, pageSize: 1, status: 3 }, { signal }),
      getThinkingModelListApi({ page: 1, pageSize: 1, status: 1 }, { signal }),
      getThinkingModelListApi({ page: 1, pageSize: 1, status: 4 }, { signal }),
    ]);
    statusTotalCounts.pending = pendingRes.total;
    statusTotalCounts.approved = approvedRes.total;
    statusTotalCounts.rejected = rejectedRes.total;
  } catch (error: any) {
    // 如果是取消错误，静默处理
    if (error?.name === 'CanceledError' || error?.code === 'ERR_CANCELED' || signal.aborted) {
      return;
    }
    // 详细的错误信息
    const errorInfo = {
      message: error?.message || 'Unknown error',
      status: error?.status || error?.response?.status,
      data: error?.data || error?.response?.data,
      code: error?.code,
    };
    console.error('加载状态统计失败:', errorInfo);
  }
}

// ===================== 数据加载 =====================

// 加载分类列表
async function loadCategories() {
  try {
    const list = await getAllCategoriesApi();
    categoryOptions.value = list.map((item) => ({
      value: Number(item.id),
      label: item.name,
    }));
  } catch (error) {
    console.error('加载分类列表失败:', error);
  }
}

// 根据标签获取状态值
function getStatusByTab(): number | undefined {
  if (activeTab.value === 'pending') return 3; // 审核中
  if (activeTab.value === 'approved') return 1; // 已发布
  if (activeTab.value === 'rejected') return 4; // 已驳回
  return undefined;
}

// 加载模型列表
async function loadModels() {
  // 取消之前的请求
  modelsAbortController?.abort();
  modelsAbortController = new AbortController();
  const signal = modelsAbortController.signal;

  loading.value = true;
  try {
    const status = getStatusByTab();
    const params: Record<string, any> = {
      page: pagination.page,
      pageSize: pagination.pageSize,
      status,
    };

    // 搜索条件
    if (searchForm.keyword) {
      params.name = searchForm.keyword;
    }
    if (searchForm.categoryId) {
      params.categoryId = Number(searchForm.categoryId);
    }

    const res = await getThinkingModelListApi(params, { signal });
    modelList.value = res.list.map((item) => ({
      ...item,
      categoryName: categoryOptions.value.find((c) => c.value === item.categoryId)?.label || '',
      submitTime: item.updatedAt,
    }));
    pagination.total = res.total;
  } catch (error: any) {
    // 如果是取消错误，静默处理
    if (error?.name === 'CanceledError' || error?.code === 'ERR_CANCELED' || signal.aborted) {
      return;
    }
    console.error('加载模型列表失败:', error);
    ElMessage.error('加载模型列表失败');
  } finally {
    loading.value = false;
  }
}

// 查看详情
async function viewDetail(model: ModelInfo) {
  detailLoading.value = true;
  detailDialogVisible.value = true;
  try {
    const res = await getThinkingModelDetailApi(model.id);
    currentModel.value = {
      ...res,
      categoryName: categoryOptions.value.find((c) => c.value === res.categoryId)?.label || '',
      submitTime: res.updatedAt,
    };
  } catch (error) {
    console.error('加载模型详情失败:', error);
    ElMessage.error('加载模型详情失败');
    detailDialogVisible.value = false;
  } finally {
    detailLoading.value = false;
  }
}

// 打开审核弹窗
function openReviewDialog(model: ModelInfo, approved: boolean) {
  currentModel.value = model;
  reviewForm.approved = approved;
  reviewForm.note = '';
  reviewDialogVisible.value = true;
}

// 提交审核结果
async function submitReview() {
  if (!reviewForm.approved && !reviewForm.note.trim()) {
    ElMessage.warning('驳回时必须填写审核意见');
    return;
  }

  if (!currentModel.value) return;

  reviewLoading.value = true;
  try {
    await reviewThinkingModelApi({
      id: currentModel.value.id,
      approved: reviewForm.approved,
      note: reviewForm.note,
    });

    const action = reviewForm.approved ? '通过' : '驳回';
    ElMessage.success(`模型"${currentModel.value.name}"已${action}审核`);

    reviewDialogVisible.value = false;
    // 更新状态统计
    loadStatusCounts();
    loadModels();
  } catch (error) {
    console.error('审核失败:', error);
    ElMessage.error('审核失败');
  } finally {
    reviewLoading.value = false;
  }
}

// 批量审核通过
async function batchApprove() {
  try {
    await ElMessageBox.confirm('确定要批量通过所有待审核的模型吗？此功能暂未实现', '批量审核', { type: 'warning' });
    ElMessage.info('批量审核功能开发中');
  } catch {
    // 取消操作
  }
}

// 搜索
function handleSearch() {
  pagination.page = 1;
  loadModels();
}

// 重置
function handleReset() {
  searchForm.keyword = '';
  searchForm.categoryId = '';
  handleSearch();
}

// 分页
function handleSizeChange(size: number) {
  pagination.pageSize = size;
  loadModels();
}

function handleCurrentChange(page: number) {
  pagination.page = page;
  loadModels();
}

// 格式化日期
function formatDate(dateStr?: string): string {
  if (!dateStr) return '-';
  return new Date(dateStr).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  });
}

// 获取状态标签
function getStatusTag(status: number) {
  return statusConfig[status] || { label: '未知', type: 'info' };
}

// 获取难度配置
function getDifficultyConfig(difficulty: number) {
  return difficultyConfig[difficulty] || { label: '未知', color: '#909399' };
}

// 标签切换
function handleTabChange() {
  pagination.page = 1;
  loadModels();
}

onMounted(async () => {
  await loadCategories();
  loadStatusCounts();
  loadModels();
});

// 组件卸载时取消所有请求
onUnmounted(() => {
  statusCountsAbortController?.abort();
  modelsAbortController?.abort();
});
</script>

<template>
  <Page
    title="模型审核"
    description="审核用户提交的思维模型，确保内容质量和规范性"
    content-class="p-6 bg-gray-50"
  >
    <div class="space-y-6">
      <!-- 状态标签页 -->
      <ElCard shadow="hover" class="!rounded-xl">
        <ElTabs v-model="activeTab" @tab-change="handleTabChange">
          <ElTabPane name="pending">
            <template #label>
              <ElBadge :value="statusCounts.pending" :max="99" class="mr-2">
                <span class="flex items-center gap-1">
                  <span>⏳</span>
                  <span>待审核</span>
                </span>
              </ElBadge>
            </template>
          </ElTabPane>
          <ElTabPane name="approved">
            <template #label>
              <ElBadge :value="statusCounts.approved" :max="99" class="mr-2" type="success">
                <span class="flex items-center gap-1">
                  <span>✅</span>
                  <span>已通过</span>
                </span>
              </ElBadge>
            </template>
          </ElTabPane>
          <ElTabPane name="rejected">
            <template #label>
              <ElBadge :value="statusCounts.rejected" :max="99" class="mr-2" type="danger">
                <span class="flex items-center gap-1">
                  <span>❌</span>
                  <span>已驳回</span>
                </span>
              </ElBadge>
            </template>
          </ElTabPane>
        </ElTabs>

        <!-- 搜索筛选 -->
        <div class="flex flex-wrap gap-4 mt-4">
          <ElInput
            v-model="searchForm.keyword"
            placeholder="搜索模型名称、描述、作者..."
            clearable
            class="!w-64"
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <span>🔍</span>
            </template>
          </ElInput>
          <ElSelect v-model="searchForm.categoryId" placeholder="选择分类" clearable class="!w-40">
            <ElOption
              v-for="cat in categoryOptions"
              :key="cat.value"
              :label="cat.label"
              :value="cat.value"
            />
          </ElSelect>
          <ElButton type="primary" @click="handleSearch">搜索</ElButton>
          <ElButton @click="handleReset">重置</ElButton>
          <ElButton
            v-if="activeTab === 'pending'"
            type="success"
            plain
            @click="batchApprove"
          >
            批量通过
          </ElButton>
        </div>
      </ElCard>

      <!-- 模型列表 -->
      <ElCard shadow="hover" class="!rounded-xl">
        <template #header>
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <span class="text-lg">📋</span>
              <span class="font-semibold text-gray-700">
                {{ activeTab === 'pending' ? '待审核模型' : activeTab === 'approved' ? '已通过模型' : '已驳回模型' }}
              </span>
              <span class="text-sm text-gray-400">(共 {{ pagination.total }} 条)</span>
            </div>
          </div>
        </template>

        <ElTable v-loading="loading" :data="modelList" stripe class="!rounded-lg">
          <ElTableColumn label="模型信息" min-width="280">
            <template #default="{ row }">
              <div class="flex items-center gap-3">
                <div
                  class="w-16 h-12 rounded-lg bg-gradient-to-br from-blue-100 to-purple-100 flex items-center justify-center text-2xl"
                >
                  {{ row.icon || '📝' }}
                </div>
                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-2">
                    <span class="font-medium text-gray-900 truncate">{{ row.name }}</span>
                    <ElTag v-if="row.isOfficial" type="warning" size="small" effect="plain" class="!rounded-full">
                      官方
                    </ElTag>
                    <ElTag v-if="!row.isFree" type="success" size="small" effect="plain" class="!rounded-full">
                      ¥{{ row.price }}
                    </ElTag>
                  </div>
                  <div class="text-sm text-gray-500 truncate">{{ row.description }}</div>
                </div>
              </div>
            </template>
          </ElTableColumn>

          <ElTableColumn label="作者" width="120">
            <template #default="{ row }">
              <div class="flex items-center gap-2">
                <ElAvatar :size="24" class="bg-blue-500">
                  {{ row.author.name.charAt(0) }}
                </ElAvatar>
                <span class="text-sm">{{ row.author.name }}</span>
              </div>
            </template>
          </ElTableColumn>

          <ElTableColumn label="分类" width="100">
            <template #default="{ row }">
              <ElTag size="small" effect="plain" class="!rounded-full">
                {{ row.categoryName || categoryOptions.find((c) => c.value === row.categoryId)?.label }}
              </ElTag>
            </template>
          </ElTableColumn>

          <ElTableColumn label="难度" width="80" align="center">
            <template #default="{ row }">
              <span
                class="px-2 py-1 rounded-full text-xs font-medium"
                :style="{ backgroundColor: `${getDifficultyConfig(row.difficulty).color}20`, color: getDifficultyConfig(row.difficulty).color }"
              >
                {{ getDifficultyConfig(row.difficulty).label }}
              </span>
            </template>
          </ElTableColumn>

          <ElTableColumn label="用时" width="80" align="center">
            <template #default="{ row }">
              <span class="text-gray-600">{{ row.estimatedTime }}分钟</span>
            </template>
          </ElTableColumn>

          <ElTableColumn label="提交时间" width="120">
            <template #default="{ row }">
              <span class="text-sm text-gray-500">{{ formatDate(row.submitTime || row.updatedAt) }}</span>
            </template>
          </ElTableColumn>

          <ElTableColumn label="状态" width="100" align="center">
            <template #default="{ row }">
              <ElTag :type="getStatusTag(row.status).type" size="small" class="!rounded-full">
                {{ getStatusTag(row.status).label }}
              </ElTag>
            </template>
          </ElTableColumn>

          <ElTableColumn label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <div class="flex gap-2">
                <ElButton type="primary" size="small" plain @click="viewDetail(row)">
                  详情
                </ElButton>
                <template v-if="row.status === 3">
                  <ElButton type="success" size="small" plain @click="openReviewDialog(row, true)">
                    通过
                  </ElButton>
                  <ElButton type="danger" size="small" plain @click="openReviewDialog(row, false)">
                    驳回
                  </ElButton>
                </template>
                <template v-else-if="row.status === 4">
                  <ElButton type="warning" size="small" plain @click="viewDetail(row)">
                    查看原因
                  </ElButton>
                </template>
              </div>
            </template>
          </ElTableColumn>
        </ElTable>

        <ElEmpty v-if="!loading && modelList.length === 0" description="暂无数据" />

        <!-- 分页 -->
        <div v-if="modelList.length > 0" class="flex justify-end mt-4">
          <ElPagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :total="pagination.total"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </ElCard>
    </div>

    <!-- 详情弹窗 -->
    <ElDialog
      v-model="detailDialogVisible"
      title="模型详情"
      width="800px"
      destroy-on-close
    >
      <div v-loading="detailLoading">
        <template v-if="currentModel">
          <!-- 基本信息 -->
          <div class="flex gap-6 mb-6">
            <div
              class="w-32 h-24 rounded-xl bg-gradient-to-br from-blue-100 to-purple-100 flex items-center justify-center text-5xl flex-shrink-0"
            >
              {{ currentModel.icon || '📝' }}
            </div>
            <div class="flex-1">
              <div class="flex items-center gap-3 mb-2">
                <h2 class="text-xl font-bold text-gray-900">{{ currentModel.name }}</h2>
                <ElTag v-if="currentModel.isOfficial" type="warning" size="small">官方</ElTag>
                <ElTag :type="currentModel.isFree ? 'success' : 'info'" size="small">
                  {{ currentModel.isFree ? '免费' : `¥${currentModel.price}` }}
                </ElTag>
              </div>
              <p class="text-gray-600 mb-2">{{ currentModel.description }}</p>
              <div class="flex items-center gap-4 text-sm text-gray-500">
                <span>作者: {{ currentModel.author.name }}</span>
                <span>版本: {{ currentModel.version }}</span>
                <span>用时: {{ currentModel.estimatedTime }}分钟</span>
              </div>
            </div>
          </div>

          <!-- 详细信息 -->
          <ElDescriptions :column="2" border class="mb-6">
            <ElDescriptionsItem label="分类">
              {{ currentModel?.categoryName || categoryOptions.find((c) => c.value === currentModel?.categoryId)?.label }}
            </ElDescriptionsItem>
            <ElDescriptionsItem label="难度">
              <span
                class="px-2 py-1 rounded-full text-xs font-medium"
                :style="{ backgroundColor: `${getDifficultyConfig(currentModel.difficulty).color}20`, color: getDifficultyConfig(currentModel.difficulty).color }"
              >
                {{ getDifficultyConfig(currentModel.difficulty).label }}
              </span>
            </ElDescriptionsItem>
            <ElDescriptionsItem label="状态">
              <ElTag :type="getStatusTag(currentModel.status).type" size="small">
                {{ getStatusTag(currentModel.status).label }}
              </ElTag>
            </ElDescriptionsItem>
            <ElDescriptionsItem label="创建时间">
              {{ formatDate(currentModel.createdAt) }}
            </ElDescriptionsItem>
            <ElDescriptionsItem label="提交时间">
              {{ formatDate(currentModel.submitTime || currentModel.updatedAt) }}
            </ElDescriptionsItem>
            <ElDescriptionsItem label="概述" :span="2">
              {{ currentModel.overview || '-' }}
            </ElDescriptionsItem>
            <ElDescriptionsItem label="标签" :span="2">
              <div class="flex flex-wrap gap-1">
                <ElTag
                  v-for="tag in currentModel.tags"
                  :key="tag"
                  size="small"
                  effect="plain"
                  class="!rounded-full"
                >
                  {{ tag }}
                </ElTag>
                <span v-if="!currentModel.tags?.length" class="text-gray-400">-</span>
              </div>
            </ElDescriptionsItem>
            <ElDescriptionsItem label="统计数据" :span="2">
              <div class="flex gap-4">
                <span>使用: {{ currentModel.stats.usageCount }}</span>
                <span>采纳: {{ currentModel.stats.adoptCount }}</span>
                <span>点赞: {{ currentModel.stats.likeCount }}</span>
                <span>评论: {{ currentModel.stats.commentCount }}</span>
              </div>
            </ElDescriptionsItem>
            <ElDescriptionsItem v-if="currentModel.reviewNote" label="驳回原因" :span="2">
              <div class="p-3 bg-red-50 rounded-lg text-red-600">
                <div class="font-medium mb-1">审核意见:</div>
                <div>{{ currentModel.reviewNote }}</div>
                <div v-if="currentModel.reviewerName" class="mt-2 text-sm text-red-400">
                  审核人: {{ currentModel.reviewerName }} | {{ currentModel.reviewTime }}
                </div>
              </div>
            </ElDescriptionsItem>
          </ElDescriptions>
        </template>
      </div>

      <template #footer>
        <ElButton @click="detailDialogVisible = false">关闭</ElButton>
        <template v-if="currentModel?.status === 3">
          <ElButton type="success" @click="detailDialogVisible = false; openReviewDialog(currentModel!, true)">
            通过审核
          </ElButton>
          <ElButton type="danger" @click="detailDialogVisible = false; openReviewDialog(currentModel!, false)">
            驳回
          </ElButton>
        </template>
      </template>
    </ElDialog>

    <!-- 审核弹窗 -->
    <ElDialog
      v-model="reviewDialogVisible"
      :title="reviewForm.approved ? '通过审核' : '驳回模型'"
      width="500px"
      destroy-on-close
    >
      <div class="mb-4">
        <p class="text-gray-600">
          模型: <span class="font-medium text-gray-900">{{ currentModel?.name }}</span>
        </p>
        <p class="text-gray-600">
          作者: <span class="font-medium text-gray-900">{{ currentModel?.author.name }}</span>
        </p>
      </div>

      <ElForm v-if="!reviewForm.approved" label-position="top">
        <ElFormItem label="驳回原因（必填）" required>
          <ElInput
            v-model="reviewForm.note"
            type="textarea"
            :rows="4"
            placeholder="请详细说明驳回原因，帮助作者改进模型..."
          />
        </ElFormItem>
      </ElForm>

      <ElForm v-else label-position="top">
        <ElFormItem label="审核意见（选填）">
          <ElInput
            v-model="reviewForm.note"
            type="textarea"
            :rows="3"
            placeholder="可以添加一些鼓励或建议..."
          />
        </ElFormItem>
      </ElForm>

      <template #footer>
        <ElButton @click="reviewDialogVisible = false">取消</ElButton>
        <ElButton
          :type="reviewForm.approved ? 'success' : 'danger'"
          :loading="reviewLoading"
          @click="submitReview"
        >
          确认{{ reviewForm.approved ? '通过' : '驳回' }}
        </ElButton>
      </template>
    </ElDialog>
  </Page>
</template>

<style scoped>
:deep(.el-tabs__item) {
  font-size: 14px;
}

:deep(.el-badge__content) {
  font-size: 10px;
}
</style>