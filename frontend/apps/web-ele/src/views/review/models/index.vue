<script lang="ts" setup>
import { onMounted, ref, reactive, computed } from 'vue';

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
  tags: string[];
  createdAt: string;
  updatedAt: string;
  submitTime?: string;
  reviewNote?: string;
  reviewerName?: string;
  reviewTime?: string;
}

// ===================== 状态定义 =====================

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
const categoryOptions = [
  { value: 1, label: '商业管理' },
  { value: 2, label: '战略规划' },
  { value: 3, label: '创新思维' },
  { value: 4, label: '分析工具' },
  { value: 5, label: '决策方法' },
  { value: 6, label: '创意构思' },
];

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

// 计算各状态数量
const statusCounts = computed(() => {
  return {
    pending: modelList.value.filter((m) => m.status === 3).length || 5,
    approved: modelList.value.filter((m) => m.status === 1).length || 12,
    rejected: modelList.value.filter((m) => m.status === 4).length || 2,
  };
});

// ===================== 数据加载 =====================

// 模拟数据
function generateMockData(): ModelInfo[] {
  const mockModels: ModelInfo[] = [
    {
      id: 1,
      name: 'SWOT 分析模型',
      description: '系统分析优势、劣势、机会与威胁的经典战略框架，适用于企业战略规划和个人决策分析。',
      coverImage: 'https://images.unsplash.com/photo-1454165804606-c3d57bc86b40?w=800&h=400&fit=crop',
      icon: '📊',
      categoryId: 2,
      categoryName: '战略规划',
      price: 0,
      isFree: true,
      overview: 'SWOT 分析是一种战略规划工具，用于帮助识别优势、劣势、机会和威胁。',
      difficulty: 2,
      estimatedTime: 30,
      status: 3,
      version: '1.0.0',
      isOfficial: false,
      author: { id: '1', name: '张三', avatar: '' },
      stats: { usageCount: 1250, adoptCount: 890, likeCount: 456, commentCount: 78 },
      tags: ['战略', '分析', '决策'],
      createdAt: '2024-01-15 10:30:00',
      updatedAt: '2024-01-20 14:20:00',
      submitTime: '2024-01-20 15:00:00',
    },
    {
      id: 2,
      name: '第一性原理思维',
      description: '像马斯克一样回归本质，打破常规的创新思考方式，深入问题本质找到创新解决方案。',
      coverImage: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=800&h=400&fit=crop',
      icon: '💡',
      categoryId: 3,
      categoryName: '创新思维',
      price: 99,
      isFree: false,
      overview: '第一性原理是一种从最基本的事实出发，逐步推导出解决方案的思维方式。',
      difficulty: 3,
      estimatedTime: 45,
      status: 3,
      version: '1.0.0',
      isOfficial: false,
      author: { id: '2', name: '李四', avatar: '' },
      stats: { usageCount: 890, adoptCount: 650, likeCount: 320, commentCount: 45 },
      tags: ['创新', '思维', '问题解决'],
      createdAt: '2024-01-18 09:15:00',
      updatedAt: '2024-01-22 11:30:00',
      submitTime: '2024-01-22 12:00:00',
    },
    {
      id: 3,
      name: 'PDCA 循环',
      description: '计划-执行-检查-行动的持续改进方法论，适用于质量管理、流程优化等场景。',
      coverImage: 'https://images.unsplash.com/photo-1552664730-d307ca884978?w=800&h=400&fit=crop',
      icon: '🔄',
      categoryId: 1,
      categoryName: '商业管理',
      price: 0,
      isFree: true,
      overview: 'PDCA 循环是一种持续改进的管理方法，通过计划、执行、检查、行动四个阶段不断优化。',
      difficulty: 1,
      estimatedTime: 20,
      status: 4,
      version: '1.0.0',
      isOfficial: false,
      author: { id: '3', name: '王五', avatar: '' },
      stats: { usageCount: 560, adoptCount: 400, likeCount: 180, commentCount: 25 },
      tags: ['管理', '改进', '质量'],
      createdAt: '2024-01-10 14:00:00',
      updatedAt: '2024-01-19 16:45:00',
      submitTime: '2024-01-19 17:00:00',
      reviewNote: '模型描述不够详细，请补充具体应用场景和案例分析。另外，概述部分需要更清晰地说明使用步骤。',
      reviewerName: '审核员A',
      reviewTime: '2024-01-20 10:00:00',
    },
    {
      id: 4,
      name: '六顶思考帽',
      description: '德博诺的经典创意工具，全方位激发团队创造力，从不同角度思考问题。',
      coverImage: 'https://images.unsplash.com/photo-1517245386807-bb43f82c33c4?w=800&h=400&fit=crop',
      icon: '🎩',
      categoryId: 3,
      categoryName: '创新思维',
      price: 49,
      isFree: false,
      overview: '六顶思考帽是一种平行思维工具，通过六种不同颜色的帽子代表六种思维方式。',
      difficulty: 2,
      estimatedTime: 35,
      status: 1,
      version: '1.0.0',
      isOfficial: true,
      author: { id: '4', name: '赵六', avatar: '' },
      stats: { usageCount: 2100, adoptCount: 1800, likeCount: 890, commentCount: 156 },
      tags: ['创意', '团队', '思维'],
      createdAt: '2024-01-05 08:00:00',
      updatedAt: '2024-01-15 09:30:00',
      submitTime: '2024-01-15 10:00:00',
    },
    {
      id: 5,
      name: '5W1H 分析模型',
      description: '六个基本问题全面剖析，确保思考无遗漏，适用于问题分析、需求梳理等场景。',
      coverImage: 'https://images.unsplash.com/photo-1450101499163-c8848c66ca85?w=800&h=400&fit=crop',
      icon: '❓',
      categoryId: 4,
      categoryName: '分析工具',
      price: 0,
      isFree: true,
      overview: '5W1H 分析法通过 What、Why、Who、When、Where、How 六个问题全面分析问题。',
      difficulty: 1,
      estimatedTime: 15,
      status: 3,
      version: '1.0.0',
      isOfficial: false,
      author: { id: '5', name: '钱七', avatar: '' },
      stats: { usageCount: 3400, adoptCount: 2800, likeCount: 1200, commentCount: 234 },
      tags: ['分析', '问题', '全面'],
      createdAt: '2024-01-12 11:20:00',
      updatedAt: '2024-01-23 13:45:00',
      submitTime: '2024-01-23 14:00:00',
    },
  ];

  return mockModels;
}

// 加载模型列表
async function loadModels() {
  loading.value = true;
  try {
    // 模拟 API 调用
    await new Promise((resolve) => setTimeout(resolve, 500));

    // 根据标签过滤
    let filteredData = generateMockData();
    if (activeTab.value === 'pending') {
      filteredData = filteredData.filter((m) => m.status === 3);
    } else if (activeTab.value === 'approved') {
      filteredData = filteredData.filter((m) => m.status === 1);
    } else if (activeTab.value === 'rejected') {
      filteredData = filteredData.filter((m) => m.status === 4);
    }

    // 搜索过滤
    if (searchForm.keyword) {
      const keyword = searchForm.keyword.toLowerCase();
      filteredData = filteredData.filter(
        (m) =>
          m.name.toLowerCase().includes(keyword) ||
          m.description.toLowerCase().includes(keyword) ||
          m.author.name.toLowerCase().includes(keyword)
      );
    }
    if (searchForm.categoryId) {
      filteredData = filteredData.filter((m) => m.categoryId === parseInt(searchForm.categoryId));
    }

    modelList.value = filteredData;
    pagination.total = filteredData.length;
  } catch (error) {
    console.error('加载模型列表失败:', error);
    ElMessage.error('加载模型列表失败');
  } finally {
    loading.value = false;
  }
}

// 查看详情
async function viewDetail(model: ModelInfo) {
  detailLoading.value = true;
  currentModel.value = model;
  detailDialogVisible.value = true;
  detailLoading.value = false;
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

  reviewLoading.value = true;
  try {
    // 模拟 API 调用
    await new Promise((resolve) => setTimeout(resolve, 500));

    const action = reviewForm.approved ? '通过' : '驳回';
    ElMessage.success(`模型"${currentModel.value?.name}"已${action}审核`);

    // 更新列表中的状态
    if (currentModel.value) {
      const index = modelList.value.findIndex((m) => m.id === currentModel.value!.id);
      if (index > -1) {
        modelList.value[index].status = reviewForm.approved ? 1 : 4;
        if (!reviewForm.approved) {
          modelList.value[index].reviewNote = reviewForm.note;
        }
      }
    }

    reviewDialogVisible.value = false;
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
    await ElMessageBox.confirm('确定要批量通过所有待审核的模型吗？', '批量审核', { type: 'warning' });
    ElMessage.success('批量审核通过成功');
    loadModels();
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

onMounted(() => {
  loadModels();
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
            <ElDescriptionsItem label="模型编码">
              <code class="px-2 py-1 bg-gray-100 rounded text-sm">{{ currentModel.code }}</code>
            </ElDescriptionsItem>
            <ElDescriptionsItem label="分类">
              {{ currentModel.categoryName || categoryOptions.find((c) => c.value === currentModel.categoryId)?.label }}
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