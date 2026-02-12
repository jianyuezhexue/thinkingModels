<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { Page } from '@vben/common-ui';
import {
  ElButton,
  ElCard,
  ElMessage,
  ElSkeleton,
  ElDialog,
  ElSelect,
  ElOption,
  ElInput,
  ElForm,
  ElFormItem,
  ElEmpty,
  ElProgress,
} from 'element-plus';
import {
  getTopicDetailApi,
  selectModelForTopicApi,
  getModelListApi,
  getTopicAnalysesApi,
  getAnalysisTemplateApi,
  createAnalysisApi,
  type TopicApi,
  type ModelApi,
  type AnalysisApi,
} from '#/api';

// 路由
const route = useRoute();
const router = useRouter();
const topicId = computed(() => route.params.id as string);

// 加载状态
const loading = ref(true);

// 课题数据
const topic = ref<TopicApi.Topic | null>(null);

// 分析列表
const analyses = ref<AnalysisApi.Analysis[]>([]);

// 可用模型列表
const availableModels = ref<ModelApi.ThinkingModel[]>([]);

// 对话框显示状态
const showModelDialog = ref(false);
const showAnalysisDialog = ref(false);

// 选中的模型
const selectedModelId = ref('');

// 分析模板
const analysisTemplate = ref<AnalysisApi.AnalysisTemplate | null>(null);

// 分析表单数据
const analysisForm = ref<Record<string, string>>({});

// 分析提交状态
const submittingAnalysis = ref(false);

// 展开的分析 ID
const expandedAnalysisId = ref<string | null>(null);

// ==================== 数据获取 ====================
async function fetchTopicDetail() {
  loading.value = true;
  try {
    const res = await getTopicDetailApi(topicId.value);
    topic.value = res;
    fetchAnalyses();
  } catch (error) {
    console.error('获取课题详情失败:', error);
    ElMessage.error('获取课题详情失败');
  } finally {
    loading.value = false;
  }
}

async function fetchAnalyses() {
  try {
    const res = await getTopicAnalysesApi(topicId.value);
    analyses.value = res;
  } catch (error) {
    console.error('获取分析记录失败:', error);
  }
}

async function fetchAvailableModels() {
  try {
    const res = await getModelListApi({ pageSize: 100 });
    availableModels.value = res.list;
  } catch (error) {
    console.error('获取模型列表失败:', error);
  }
}

onMounted(() => {
  fetchTopicDetail();
  fetchAvailableModels();
});

// ==================== 工具函数 ====================
function getStatusStyle(status: TopicApi.TopicStatus): string {
  const styles: Record<string, string> = {
    draft: 'bg-gray-100 text-gray-600',
    in_progress: 'bg-amber-100 text-amber-700',
    completed: 'bg-green-100 text-green-700',
    archived: 'bg-slate-100 text-slate-600',
  };
  return styles[status] || 'bg-gray-100 text-gray-600';
}

function getStatusText(status: TopicApi.TopicStatus): string {
  const texts: Record<string, string> = {
    draft: '草稿',
    in_progress: '进行中',
    completed: '已完成',
    archived: '已归档',
  };
  return texts[status] || status;
}

function getStatusIcon(status: TopicApi.TopicStatus): string {
  const icons: Record<string, string> = {
    draft: '📝',
    in_progress: '⏳',
    completed: '✅',
    archived: '📦',
  };
  return icons[status] || '📋';
}

function getProgressValue(status: TopicApi.TopicStatus): number {
  const map: Record<string, number> = {
    draft: 15,
    in_progress: 60,
    completed: 100,
    archived: 100,
  };
  return map[status] || 0;
}

function formatDate(dateStr: string): string {
  if (!dateStr) return '-';
  return new Date(dateStr).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
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

// ==================== 操作 ====================
async function openModelDialog() {
  await fetchAvailableModels();
  selectedModelId.value = topic.value?.modelId || '';
  showModelDialog.value = true;
}

async function handleSelectModel() {
  if (!selectedModelId.value) {
    ElMessage.warning('请选择一个思维模型');
    return;
  }

  try {
    await selectModelForTopicApi({
      topicId: topicId.value,
      modelId: selectedModelId.value,
    });
    ElMessage.success('模型选用成功');
    showModelDialog.value = false;
    fetchTopicDetail();
  } catch (error) {
    console.error('选用模型失败:', error);
    ElMessage.error('选用模型失败');
  }
}

async function openAnalysisDialog() {
  if (!topic.value?.modelId) {
    ElMessage.warning('请先选用一个思维模型');
    return;
  }

  try {
    const template = await getAnalysisTemplateApi(topic.value.modelId);
    analysisTemplate.value = template;
    analysisForm.value = {};
    template.fields.forEach((field) => {
      analysisForm.value[field.key] = '';
    });
    showAnalysisDialog.value = true;
  } catch (error) {
    console.error('获取分析模板失败:', error);
    ElMessage.error('获取分析模板失败');
  }
}

async function handleSubmitAnalysis() {
  if (!topic.value?.modelId || !analysisTemplate.value) return;

  const emptyFields = analysisTemplate.value.fields.filter(
    (f) => f.required && !analysisForm.value[f.key]?.trim()
  );
  if (emptyFields.length > 0) {
    ElMessage.warning('请填写: ' + emptyFields.map((f) => f.label).join(', '));
    return;
  }

  submittingAnalysis.value = true;
  try {
    const inputs: AnalysisApi.AnalysisInput[] = analysisTemplate.value.fields.map((field) => ({
      key: field.key,
      label: field.label,
      value: analysisForm.value[field.key] ?? '',
    }));

    await createAnalysisApi({
      topicId: topicId.value,
      modelId: topic.value.modelId,
      inputs,
    });

    ElMessage.success('分析提交成功，正在生成结果...');
    showAnalysisDialog.value = false;
    fetchAnalyses();
  } catch (error) {
    console.error('提交分析失败:', error);
    ElMessage.error('提交分析失败');
  } finally {
    submittingAnalysis.value = false;
  }
}

function toggleAnalysis(id: string) {
  expandedAnalysisId.value = expandedAnalysisId.value === id ? null : id;
}

function goBack() {
  router.push('/my-topics');
}

function goToMarket() {
  router.push('/market');
}
</script>

<template>
  <Page
    title="课题详情"
    description="查看课题详情，使用思维模型进行深度分析"
    content-class="p-6 bg-gray-50"
  >
    <!-- 加载状态 -->
    <div v-if="loading" class="space-y-6">
      <ElCard shadow="hover" class="!rounded-xl">
        <ElSkeleton animated :rows="4" />
      </ElCard>
    </div>

    <!-- 内容 -->
    <div v-else-if="topic" class="space-y-6">
      <!-- 顶部导航 -->
      <div class="flex items-center justify-between">
        <button
          class="flex items-center gap-2 text-gray-500 hover:text-purple-600 transition-colors"
          @click="goBack"
        >
          <span class="text-lg">←</span>
          <span>返回课题列表</span>
        </button>
        <div class="flex items-center gap-3">
          <ElButton
            type="primary"
            class="!bg-purple-600 !border-purple-600 hover:!bg-purple-700 !rounded-full"
            @click="openAnalysisDialog"
            :disabled="!topic.modelId"
          >
            开始新分析
          </ElButton>
        </div>
      </div>

      <!-- 课题信息卡片 -->
      <ElCard shadow="hover" class="!rounded-xl overflow-hidden">
        <div class="flex gap-6">
          <!-- 左侧进度环 -->
          <div class="flex-shrink-0">
            <ElProgress
              type="circle"
              :percentage="getProgressValue(topic.status)"
              :width="100"
              :stroke-width="6"
              :color="topic.status === 'completed' ? '#10b981' : topic.status === 'in_progress' ? '#f59e0b' : '#9ca3af'"
            >
              <template #default>
                <span class="text-3xl">{{ getStatusIcon(topic.status) }}</span>
              </template>
            </ElProgress>
          </div>

          <!-- 右侧内容 -->
          <div class="flex-1 min-w-0">
            <div class="flex items-start justify-between mb-4">
              <div>
                <div class="flex items-center gap-3 mb-2">
                  <h1 class="text-2xl font-bold text-gray-800">{{ topic.title }}</h1>
                  <span
                    class="px-3 py-1 rounded-full text-sm"
                    :class="getStatusStyle(topic.status)"
                  >
                    {{ getStatusText(topic.status) }}
                  </span>
                </div>
                <p class="text-gray-600 whitespace-pre-wrap">{{ topic.description }}</p>
              </div>
            </div>

            <!-- 元信息 -->
            <div class="flex flex-wrap items-center gap-6 text-sm text-gray-500 pt-4 border-t border-gray-100">
              <span class="flex items-center gap-2">
                <span class="w-8 h-8 rounded-lg bg-purple-100 flex items-center justify-center text-purple-600">🧠</span>
                <span>{{ topic.modelName || '未选用模型' }}</span>
              </span>
              <span class="flex items-center gap-2">
                <span class="w-8 h-8 rounded-lg bg-blue-100 flex items-center justify-center text-blue-600">📊</span>
                <span>{{ topic.analysisCount }} 次分析</span>
              </span>
              <span class="flex items-center gap-2">
                <span class="w-8 h-8 rounded-lg bg-green-100 flex items-center justify-center text-green-600">📅</span>
                <span>创建于 {{ formatDate(topic.createdAt) }}</span>
              </span>
            </div>
          </div>
        </div>
      </ElCard>

      <!-- 主内容区 -->
      <div class="flex gap-6">
        <!-- 左侧分析记录 -->
        <div class="flex-1 space-y-4">
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-semibold text-gray-800">分析记录</h2>
            <span class="text-sm text-gray-500">共 {{ analyses.length }} 条记录</span>
          </div>

          <!-- 空状态 -->
          <ElCard v-if="analyses.length === 0" shadow="hover" class="!rounded-xl">
            <ElEmpty description="暂无分析记录">
              <template #image>
                <div class="text-6xl">📝</div>
              </template>
              <div class="text-center">
                <p class="text-gray-500 mb-4">选用思维模型后，开始你的第一次分析吧！</p>
                <ElButton
                  v-if="topic.modelId"
                  type="primary"
                  class="!bg-purple-600 !border-purple-600 !rounded-full"
                  @click="openAnalysisDialog"
                >
                  开始分析
                </ElButton>
                <ElButton
                  v-else
                  type="primary"
                  class="!bg-purple-600 !border-purple-600 !rounded-full"
                  @click="openModelDialog"
                >
                  选用模型
                </ElButton>
              </div>
            </ElEmpty>
          </ElCard>

          <!-- 分析列表 -->
          <div v-else class="space-y-4">
            <ElCard
              v-for="analysis in analyses"
              :key="analysis.id"
              shadow="hover"
              class="!rounded-xl cursor-pointer transition-all"
              :class="{ '!border-purple-300': expandedAnalysisId === analysis.id }"
              @click="toggleAnalysis(analysis.id)"
            >
              <!-- 折叠头部 -->
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 rounded-lg bg-purple-100 flex items-center justify-center">
                    <span class="text-xl">🧠</span>
                  </div>
                  <div>
                    <div class="font-medium text-gray-800">{{ analysis.modelName }}</div>
                    <div class="text-xs text-gray-500">{{ formatTime(analysis.createdAt) }}</div>
                  </div>
                </div>
                <div class="flex items-center gap-3">
                  <span
                    v-if="analysis.status === 'completed'"
                    class="px-2 py-1 rounded-full text-xs bg-green-100 text-green-700"
                  >
                    已完成
                  </span>
                  <span
                    v-else-if="analysis.status === 'processing'"
                    class="px-2 py-1 rounded-full text-xs bg-amber-100 text-amber-700"
                  >
                    分析中
                  </span>
                  <span
                    v-else-if="analysis.status === 'failed'"
                    class="px-2 py-1 rounded-full text-xs bg-red-100 text-red-700"
                  >
                    失败
                  </span>
                  <span
                    class="text-gray-400 transition-transform"
                    :class="{ 'rotate-180': expandedAnalysisId === analysis.id }"
                  >
                    ▼
                  </span>
                </div>
              </div>

              <!-- 展开内容 -->
              <div
                v-if="expandedAnalysisId === analysis.id"
                class="mt-4 pt-4 border-t border-gray-100 space-y-4"
                @click.stop
              >
                <!-- 输入内容 -->
                <div class="bg-gray-50 rounded-lg p-4">
                  <h4 class="font-medium text-sm text-gray-600 mb-3 flex items-center gap-2">
                    <span>📥</span> 你的输入
                  </h4>
                  <div class="space-y-3">
                    <div
                      v-for="input in analysis.inputs"
                      :key="input.key"
                      class="bg-white rounded-lg p-3"
                    >
                      <div class="text-xs text-gray-500 mb-1">{{ input.label }}</div>
                      <div class="text-sm text-gray-700">{{ input.value }}</div>
                    </div>
                  </div>
                </div>

                <!-- 分析结果 -->
                <div v-if="analysis.results && analysis.results.length > 0">
                  <h4 class="font-medium text-sm text-purple-600 mb-3 flex items-center gap-2">
                    <span>✨</span> 分析结果
                  </h4>
                  <div class="space-y-3">
                    <div
                      v-for="result in analysis.results"
                      :key="result.key"
                      class="bg-purple-50 border border-purple-100 rounded-lg p-4"
                    >
                      <h5 class="font-medium text-sm text-purple-700 mb-2">{{ result.label }}</h5>
                      <p class="text-sm text-gray-700 whitespace-pre-wrap">{{ result.content }}</p>
                    </div>
                  </div>
                </div>

                <!-- 分析中状态 -->
                <div v-else-if="analysis.status === 'processing'" class="text-center py-6">
                  <div class="inline-block animate-spin text-3xl mb-2">⏳</div>
                  <div class="text-purple-600 text-sm">分析进行中，请稍候...</div>
                </div>

                <!-- 失败状态 -->
                <div v-else-if="analysis.status === 'failed'" class="text-center py-6">
                  <div class="text-3xl mb-2">❌</div>
                  <div class="text-red-500 text-sm">分析失败，请重试</div>
                </div>
              </div>
            </ElCard>
          </div>
        </div>

        <!-- 右侧边栏 -->
        <div class="w-80 flex-shrink-0 space-y-6 hidden lg:block">
          <!-- 思维模型卡片 -->
          <ElCard shadow="hover" class="!rounded-xl">
            <template #header>
              <span class="font-semibold text-gray-700">当前模型</span>
            </template>

            <div v-if="topic.modelId" class="text-center">
              <div class="w-16 h-16 rounded-xl bg-purple-100 flex items-center justify-center mx-auto mb-3">
                <span class="text-3xl">🧠</span>
              </div>
              <div class="font-semibold text-lg text-gray-800 mb-1">{{ topic.modelName }}</div>
              <div class="text-sm text-gray-500 mb-4">已选用</div>
              <ElButton
                class="w-full !rounded-full"
                @click="openModelDialog"
              >
                更换模型
              </ElButton>
            </div>

            <div v-else class="text-center py-4">
              <div class="w-16 h-16 rounded-xl bg-gray-100 flex items-center justify-center mx-auto mb-3">
                <span class="text-3xl">🤔</span>
              </div>
              <div class="text-gray-500 mb-4">还没有选用模型</div>
              <ElButton
                type="primary"
                class="w-full !bg-purple-600 !border-purple-600 !rounded-full mb-2"
                @click="openModelDialog"
              >
                选用模型
              </ElButton>
              <button
                class="text-sm text-purple-600 hover:text-purple-700"
                @click="goToMarket"
              >
                去市场浏览 →
              </button>
            </div>
          </ElCard>

          <!-- 快速操作 -->
          <ElCard shadow="hover" class="!rounded-xl">
            <template #header>
              <span class="font-semibold text-gray-700">快速操作</span>
            </template>
            <div class="space-y-3">
              <button
                class="w-full flex items-center gap-3 p-3 rounded-lg bg-purple-50 hover:bg-purple-100 text-purple-700 transition-colors"
                :class="{ 'opacity-50 cursor-not-allowed': !topic.modelId }"
                :disabled="!topic.modelId"
                @click="openAnalysisDialog"
              >
                <span class="text-lg">📝</span>
                <span class="font-medium">开始分析</span>
              </button>
              <button
                class="w-full flex items-center gap-3 p-3 rounded-lg bg-gray-50 hover:bg-gray-100 text-gray-700 transition-colors"
                @click="goToMarket"
              >
                <span class="text-lg">🛒</span>
                <span class="font-medium">浏览更多模型</span>
              </button>
              <button
                class="w-full flex items-center gap-3 p-3 rounded-lg bg-gray-50 hover:bg-gray-100 text-gray-700 transition-colors"
                @click="goBack"
              >
                <span class="text-lg">📋</span>
                <span class="font-medium">返回课题列表</span>
              </button>
            </div>
          </ElCard>

          <!-- 分析提示 -->
          <ElCard shadow="hover" class="!rounded-xl !bg-gradient-to-br from-purple-50 to-purple-100 !border-purple-200">
            <template #header>
              <span class="font-semibold text-purple-700">分析提示</span>
            </template>
            <ul class="text-sm text-purple-800 space-y-2">
              <li class="flex items-start gap-2">
                <span class="text-purple-500">•</span>
                详细描述问题能获得更好的分析
              </li>
              <li class="flex items-start gap-2">
                <span class="text-purple-500">•</span>
                尝试不同模型获得多角度见解
              </li>
              <li class="flex items-start gap-2">
                <span class="text-purple-500">•</span>
                分析结果可导出和分享
              </li>
            </ul>
          </ElCard>
        </div>
      </div>
    </div>

    <!-- 错误状态 -->
    <ElCard v-else shadow="hover" class="!rounded-xl">
      <ElEmpty description="课题不存在或已被删除">
        <template #image>
          <div class="text-6xl">😕</div>
        </template>
        <ElButton
          type="primary"
          class="!bg-purple-600 !border-purple-600 !rounded-full"
          @click="goBack"
        >
          返回列表
        </ElButton>
      </ElEmpty>
    </ElCard>

    <!-- 选用模型对话框 -->
    <ElDialog
      v-model="showModelDialog"
      title="选用思维模型"
      width="600px"
    >
      <div class="mb-4 p-4 bg-purple-50 rounded-lg">
        <div class="flex items-start gap-3">
          <span class="text-2xl">💡</span>
          <div>
            <p class="font-medium text-purple-900">选择适合的思维模型</p>
            <p class="text-sm text-purple-700 mt-1">不同的思维模型适用于不同类型的问题分析。</p>
          </div>
        </div>
      </div>

      <ElSelect
        v-model="selectedModelId"
        placeholder="请选择思维模型"
        class="w-full"
        size="large"
      >
        <ElOption
          v-for="model in availableModels"
          :key="model.id"
          :label="model.title"
          :value="model.id"
        >
          <div class="flex items-center justify-between py-1">
            <div class="flex items-center gap-2">
              <span class="text-lg">🧠</span>
              <span>{{ model.title }}</span>
            </div>
            <span
              class="px-2 py-0.5 rounded-full text-xs"
              :class="model.isFree ? 'bg-green-100 text-green-700' : 'bg-amber-100 text-amber-700'"
            >
              {{ model.isFree ? '免费' : '¥' + model.price }}
            </span>
          </div>
        </ElOption>
      </ElSelect>

      <template #footer>
        <ElButton @click="showModelDialog = false">取消</ElButton>
        <ElButton
          type="primary"
          class="!bg-purple-600 !border-purple-600"
          @click="handleSelectModel"
        >
          确定选用
        </ElButton>
      </template>
    </ElDialog>

    <!-- 分析对话框 -->
    <ElDialog
      v-model="showAnalysisDialog"
      title="开始分析"
      width="700px"
    >
      <div v-if="analysisTemplate">
        <div class="mb-4 p-4 bg-purple-50 rounded-lg">
          <div class="flex items-start gap-3">
            <span class="text-2xl">🧠</span>
            <div>
              <p class="font-medium text-purple-900">{{ analysisTemplate.modelName }}</p>
              <p class="text-sm text-purple-700 mt-1">请根据提示填写各项内容，系统将根据你的输入生成分析结果。</p>
            </div>
          </div>
        </div>

        <ElForm label-position="top">
          <ElFormItem
            v-for="field in analysisTemplate.fields"
            :key="field.key"
            :label="field.label"
            :required="field.required"
          >
            <ElInput
              v-if="field.type === 'textarea'"
              v-model="analysisForm[field.key]"
              type="textarea"
              :rows="4"
              :placeholder="field.placeholder || '请输入' + field.label"
            />
            <ElInput
              v-else
              v-model="analysisForm[field.key]"
              :placeholder="field.placeholder || '请输入' + field.label"
            />
          </ElFormItem>
        </ElForm>
      </div>

      <template #footer>
        <ElButton @click="showAnalysisDialog = false">取消</ElButton>
        <ElButton
          type="primary"
          class="!bg-purple-600 !border-purple-600"
          :loading="submittingAnalysis"
          @click="handleSubmitAnalysis"
        >
          提交分析
        </ElButton>
      </template>
    </ElDialog>
  </Page>
</template>
