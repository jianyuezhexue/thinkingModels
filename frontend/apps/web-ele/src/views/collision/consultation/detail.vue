<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { Page } from '@vben/common-ui';
import {
  ElCard,
  ElButton,
  ElAvatar,
  ElDialog,
  ElInput,
  ElMessage,
  ElEmpty,
  ElSkeleton,
} from 'element-plus';
import {
  getConsultationDetailApi,
  applyConsultationApi,
  type CollisionApi,
} from '#/api/collision';

const route = useRoute();
const router = useRouter();

const consultation = ref<CollisionApi.Consultation | null>(null);
const loading = ref(true);
const applyDialogVisible = ref(false);
const applyForm = ref({
  proposal: '',
  estimatedTime: '',
  quotation: 0,
});
const applying = ref(false);

const consultationId = computed(() => String(route.params.id));

const getStatusStyle = (status: string) => {
  const styles: Record<string, string> = {
    open: 'bg-green-100 text-green-700',
    matched: 'bg-blue-100 text-blue-700',
    inProgress: 'bg-purple-100 text-purple-700',
    completed: 'bg-gray-100 text-gray-700',
    cancelled: 'bg-red-100 text-red-700',
    expired: 'bg-gray-100 text-gray-500',
  };
  return styles[status] || 'bg-gray-100 text-gray-600';
};

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    open: '征集中',
    matched: '已匹配',
    inProgress: '进行中',
    completed: '已完成',
    cancelled: '已取消',
    expired: '已过期',
  };
  return texts[status] || status;
};

const getModeStyle = (mode: string) => {
  const styles: Record<string, string> = {
    online: 'bg-blue-100 text-blue-700',
    offline: 'bg-green-100 text-green-700',
    both: 'bg-purple-100 text-purple-700',
  };
  return styles[mode] || 'bg-gray-100 text-gray-600';
};

const getModeText = (mode: string) => {
  const texts: Record<string, string> = {
    online: '🖥️ 线上咨询',
    offline: '☕ 线下咨询',
    both: '🔄 线上/线下均可',
  };
  return texts[mode] || mode;
};

const getFieldName = (field: string) => {
  const fields: Record<string, string> = {
    career: '职业发展',
    startup: '创业指导',
    technology: '技术架构',
    product: '产品设计',
    investment: '投资理财',
    management: '团队管理',
    psychology: '心理咨询',
    other: '其他',
  };
  return fields[field] || field;
};

const formatDate = (dateStr?: string) => {
  if (!dateStr) return '';
  const date = new Date(dateStr);
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
};

const formatDateTime = (dateStr?: string) => {
  if (!dateStr) return '';
  const date = new Date(dateStr);
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  });
};

const daysUntilDeadline = computed(() => {
  if (!consultation.value?.deadline) return 0;
  const deadline = new Date(consultation.value.deadline);
  const now = new Date();
  const diff = deadline.getTime() - now.getTime();
  return Math.ceil(diff / (1000 * 60 * 60 * 24));
});

const fetchConsultationDetail = async () => {
  loading.value = true;
  try {
    const res = await getConsultationDetailApi(consultationId.value);
    consultation.value = res;
  } catch (error) {
    console.error('Failed to fetch consultation detail:', error);
  } finally {
    loading.value = false;
  }
};

const openApplyDialog = () => {
  applyDialogVisible.value = true;
};

const submitApplication = async () => {
  if (!applyForm.value.proposal.trim()) {
    ElMessage.warning('请填写申请说明');
    return;
  }
  if (!applyForm.value.estimatedTime.trim()) {
    ElMessage.warning('请填写预计时间');
    return;
  }

  applying.value = true;
  try {
    await applyConsultationApi({
      consultationId: consultationId.value,
      proposal: applyForm.value.proposal,
      estimatedTime: applyForm.value.estimatedTime,
      quotation: applyForm.value.quotation || undefined,
    });
    ElMessage.success('申请已提交，请等待对方回复');
    applyDialogVisible.value = false;
    applyForm.value = { proposal: '', estimatedTime: '', quotation: 0 };
    fetchConsultationDetail();
  } catch (error) {
    ElMessage.error('申请提交失败');
  } finally {
    applying.value = false;
  }
};

const goBack = () => {
  router.push({ name: 'Collision', query: { tab: 'consultation' } });
};

onMounted(() => {
  fetchConsultationDetail();
});
</script>

<template>
  <Page title="咨询详情" content-class="p-6 bg-gray-50">
    <!-- 返回按钮 -->
    <div class="mb-4">
      <ElButton text @click="goBack">
        ← 返回咨询列表
      </ElButton>
    </div>

    <ElSkeleton v-if="loading" :rows="10" animated />

    <template v-else-if="consultation">
      <div class="flex gap-6">
        <!-- 左侧主内容 -->
        <div class="flex-1 space-y-6">
          <!-- 咨询标题卡片 -->
          <ElCard shadow="hover" class="!rounded-xl">
            <div class="flex items-start gap-4">
              <!-- 悬赏金额 -->
              <div class="flex-shrink-0 w-24 h-24 rounded-xl bg-gradient-to-br from-orange-400 to-red-500 flex flex-col items-center justify-center text-white shadow-lg">
                <span class="text-xs opacity-80">悬赏金额</span>
                <span class="text-2xl font-bold">¥{{ consultation.reward }}</span>
              </div>
              
              <div class="flex-1">
                <div class="flex items-center gap-3 mb-2">
                  <span
                    class="px-3 py-1 rounded-full text-sm"
                    :class="getStatusStyle(consultation.status)"
                  >
                    {{ getStatusText(consultation.status) }}
                  </span>
                  <span
                    class="px-3 py-1 rounded-full text-sm"
                    :class="getModeStyle(consultation.mode)"
                  >
                    {{ getModeText(consultation.mode) }}
                  </span>
                  <span class="px-3 py-1 rounded-full text-sm bg-gray-100 text-gray-600">
                    {{ getFieldName(consultation.field) }}
                  </span>
                </div>
                <h1 class="text-2xl font-bold text-gray-800 mb-3">
                  {{ consultation.title }}
                </h1>
                <div class="flex items-center gap-6 text-sm text-gray-500">
                  <span>⏰ 截止日期: {{ formatDate(consultation.deadline) }}</span>
                  <span
                    v-if="daysUntilDeadline > 0"
                    class="text-orange-500 font-medium"
                  >
                    还剩 {{ daysUntilDeadline }} 天
                  </span>
                  <span v-else class="text-red-500 font-medium">已截止</span>
                </div>
              </div>
            </div>
          </ElCard>

          <!-- 问题描述 -->
          <ElCard shadow="hover" class="!rounded-xl">
            <template #header>
              <span class="font-semibold text-gray-700">💭 问题描述</span>
            </template>
            <p class="text-gray-700 whitespace-pre-wrap leading-relaxed">
              {{ consultation.description }}
            </p>
          </ElCard>

          <!-- 背景信息 -->
          <ElCard v-if="consultation.background" shadow="hover" class="!rounded-xl">
            <template #header>
              <span class="font-semibold text-gray-700">📋 背景信息</span>
            </template>
            <p class="text-gray-700 whitespace-pre-wrap leading-relaxed">
              {{ consultation.background }}
            </p>
          </ElCard>

          <!-- 期望获得 -->
          <ElCard v-if="consultation.expectation" shadow="hover" class="!rounded-xl">
            <template #header>
              <span class="font-semibold text-gray-700">🎯 期望获得</span>
            </template>
            <p class="text-gray-700 whitespace-pre-wrap leading-relaxed">
              {{ consultation.expectation }}
            </p>
          </ElCard>

          <!-- 专家申请列表 -->
          <ElCard shadow="hover" class="!rounded-xl">
            <template #header>
              <div class="flex items-center justify-between">
                <span class="font-semibold text-gray-700">👨‍💼 专家申请 ({{ consultation.applicationCount || 0 }})</span>
              </div>
            </template>
            
            <ElEmpty
              v-if="!consultation.applicationCount"
              description="暂无专家申请，快来成为第一个！"
            />
            
            <div v-else class="text-center py-4">
              <p class="text-gray-600">已有 {{ consultation.applicationCount }} 位专家申请</p>
              <p class="text-sm text-gray-400 mt-2">发布者可在消息中心查看详情</p>
            </div>
          </ElCard>
        </div>

        <!-- 右侧边栏 -->
        <div class="w-80 flex-shrink-0 space-y-6 hidden lg:block">
          <!-- 发布者信息 -->
          <ElCard shadow="hover" class="!rounded-xl">
            <template #header>
              <span class="font-semibold text-gray-700">👤 发布者</span>
            </template>
            <div class="flex flex-col items-center text-center">
              <ElAvatar :size="80" :src="consultation.requester?.avatar" class="mb-3">
                {{ consultation.requester?.name?.charAt(0) }}
              </ElAvatar>
              <h3 class="font-semibold text-gray-800 mb-1">
                {{ consultation.requester?.name }}
              </h3>
              <p class="text-sm text-gray-500 mb-4">
                发布于 {{ formatDateTime(consultation.createdAt) }}
              </p>
            </div>
          </ElCard>

          <!-- 申请咨询按钮 -->
          <ElCard
            v-if="consultation.status === 'open'"
            shadow="hover"
            class="!rounded-xl !bg-gradient-to-br from-purple-50 to-purple-100"
          >
            <div class="text-center">
              <div class="text-4xl mb-3">🌟</div>
              <h3 class="font-semibold text-purple-700 mb-2">成为咨询专家</h3>
              <p class="text-sm text-purple-600 mb-4">
                如果您有相关经验，可以申请成为本次咨询的专家
              </p>
              <ElButton
                type="primary"
                size="large"
                class="w-full !bg-purple-600 !border-purple-600 hover:!bg-purple-700 !rounded-full"
                @click="openApplyDialog"
              >
                💼 申请咨询
              </ElButton>
            </div>
          </ElCard>

          <!-- 咨询统计 -->
          <ElCard shadow="hover" class="!rounded-xl">
            <template #header>
              <span class="font-semibold text-gray-700">📊 统计数据</span>
            </template>
            <div class="grid grid-cols-2 gap-4">
              <div class="text-center p-3 bg-gray-50 rounded-lg">
                <div class="text-2xl font-bold text-purple-600">{{ consultation.viewCount || 0 }}</div>
                <div class="text-xs text-gray-500">浏览量</div>
              </div>
              <div class="text-center p-3 bg-gray-50 rounded-lg">
                <div class="text-2xl font-bold text-orange-500">{{ consultation.applicationCount || 0 }}</div>
                <div class="text-xs text-gray-500">申请数</div>
              </div>
            </div>
          </ElCard>

          <!-- 温馨提示 -->
          <ElCard shadow="hover" class="!rounded-xl">
            <template #header>
              <span class="font-semibold text-gray-700">💡 温馨提示</span>
            </template>
            <ul class="text-sm text-gray-600 space-y-2">
              <li class="flex items-start gap-2">
                <span class="text-purple-500">•</span>
                请认真阅读咨询需求
              </li>
              <li class="flex items-start gap-2">
                <span class="text-purple-500">•</span>
                确保有相关专业能力
              </li>
              <li class="flex items-start gap-2">
                <span class="text-purple-500">•</span>
                申请后请保持沟通畅通
              </li>
              <li class="flex items-start gap-2">
                <span class="text-purple-500">•</span>
                咨询完成后请确认结算
              </li>
            </ul>
          </ElCard>
        </div>
      </div>
    </template>

    <ElEmpty v-else description="咨询不存在或已被删除" />

    <!-- 申请对话框 -->
    <ElDialog
      v-model="applyDialogVisible"
      title="申请成为咨询专家"
      width="600px"
      :close-on-click-modal="false"
    >
      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            申请说明 <span class="text-red-500">*</span>
          </label>
          <ElInput
            v-model="applyForm.proposal"
            type="textarea"
            :rows="4"
            placeholder="介绍您的相关背景、解决思路和专业能力..."
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            预计时间 <span class="text-red-500">*</span>
          </label>
          <ElInput
            v-model="applyForm.estimatedTime"
            placeholder="例如：1-2小时线上沟通"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            报价（可选）
          </label>
          <ElInput
            v-model.number="applyForm.quotation"
            type="number"
            placeholder="如果您有额外报价需求，请填写金额"
          />
        </div>
      </div>
      <template #footer>
        <ElButton @click="applyDialogVisible = false">取消</ElButton>
        <ElButton
          type="primary"
          :loading="applying"
          class="!bg-purple-600 !border-purple-600"
          @click="submitApplication"
        >
          提交申请
        </ElButton>
      </template>
    </ElDialog>
  </Page>
</template>
