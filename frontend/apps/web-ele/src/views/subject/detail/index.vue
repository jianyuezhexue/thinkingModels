<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { Page } from '@vben/common-ui';

import {
  ElButton,
  ElCard,
  ElTag,
  ElMessage,
  ElSkeleton,
  ElDialog,
  ElSelect,
  ElOption,
  ElInput,
  ElForm,
  ElFormItem,
  ElEmpty,
  ElDivider,
  ElCollapse,
  ElCollapseItem,
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

// 获取课题详情
async function fetchTopicDetail() {
  loading.value = true;
  try {
    const res = await getTopicDetailApi(topicId.value);
    topic.value = res;
    // 获取该课题的分析记录
    fetchAnalyses();
  } catch (error) {
    console.error('获取课题详情失败:', error);
    ElMessage.error('获取课题详情失败');
  } finally {
    loading.value = false;
  }
}

// 获取分析记录
async function fetchAnalyses() {
  try {
    const res = await getTopicAnalysesApi(topicId.value);
    analyses.value = res;
  } catch (error) {
    console.error('获取分析记录失败:', error);
  }
}

// 获取可用模型列表
async function fetchAvailableModels() {
  try {
    const res = await getModelListApi({ pageSize: 100 });
    availableModels.value = res.list;
  } catch (error) {
    console.error('获取模型列表失败:', error);
  }
}

// 页面加载时获取数据
onMounted(() => {
  fetchTopicDetail();
  fetchAvailableModels();
});

// 打开选用模型对话框
async function openModelDialog() {
  await fetchAvailableModels();
  selectedModelId.value = topic.value?.modelId || '';
  showModelDialog.value = true;
}

// 选用模型
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

// 打开分析对话框
async function openAnalysisDialog() {
  if (!topic.value?.modelId) {
    ElMessage.warning('请先选用一个思维模型');
    return;
  }

  try {
    const template = await getAnalysisTemplateApi(topic.value.modelId);
    analysisTemplate.value = template;
    // 初始化表单
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

// 提交分析
async function handleSubmitAnalysis() {
  if (!topic.value?.modelId || !analysisTemplate.value) return;

  // 验证必填项
  const emptyFields = analysisTemplate.value.fields.filter(
    (f) => f.required && !analysisForm.value[f.key]?.trim()
  );
  if (emptyFields.length > 0) {
    ElMessage.warning(`请填写: ${emptyFields.map((f) => f.label).join(', ')}`);
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

// 返回列表
function goBack() {
  router.push('/my-topics');
}

// 去市场选模型
function goToMarket() {
  router.push('/market');
}
</script>

<template>
  <Page title="课题详情" description="查看课题详情，使用思维模型进行深度分析">
    <!-- 加载状态 -->
    <div v-if="loading" class="space-y-6">
      <ElCard shadow="never">
        <ElSkeleton animated :rows="3" />
      </ElCard>
      <ElCard shadow="never">
        <ElSkeleton animated :rows="5" />
      </ElCard>
    </div>

    <!-- 内容 -->
    <div v-else-if="topic" class="grid grid-cols-1 gap-6 lg:grid-cols-3">
      <!-- 左侧：课题信息 -->
      <div class="lg:col-span-2 space-y-6">
        <!-- 基本信息 -->
        <ElCard shadow="never">
          <template #header>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-3">
                <h2 class="text-xl font-bold">{{ topic.title }}</h2>
                <ElTag :type="getStatusType(topic.status)">
                  {{ getStatusText(topic.status) }}
                </ElTag>
              </div>
              <ElButton @click="goBack">← 返回列表</ElButton>
            </div>
          </template>

          <div class="space-y-4">
            <div>
              <h3 class="text-sm font-medium text-gray-500 mb-2">课题描述</h3>
              <p class="text-gray-700 whitespace-pre-wrap">{{ topic.description }}</p>
            </div>

            <ElDivider />

            <div class="flex flex-wrap items-center gap-6 text-sm text-gray-500">
              <span class="flex items-center gap-1">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"/>
                </svg>
                当前模型：{{ topic.modelName || '未选用' }}
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
            </div>
          </div>
        </ElCard>

        <!-- 分析记录 -->
        <ElCard shadow="never">
          <template #header>
            <div class="flex items-center justify-between">
              <h3 class="font-medium">分析记录</h3>
              <ElButton type="primary" @click="openAnalysisDialog">
                <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>
                </svg>
                开始新分析
              </ElButton>
            </div>
          </template>

          <ElEmpty v-if="analyses.length === 0" description="暂无分析记录，点击上方按钮开始分析">
            <ElButton type="primary" @click="openAnalysisDialog">开始分析</ElButton>
          </ElEmpty>

          <ElCollapse v-else>
            <ElCollapseItem
              v-for="analysis in analyses"
              :key="analysis.id"
              :title="`${analysis.modelName} - ${formatDate(analysis.createdAt)}`"
            >
              <div class="space-y-4">
                <!-- 输入内容 -->
                <div class="bg-gray-50 rounded-lg p-4">
                  <h4 class="font-medium text-sm text-gray-600 mb-2">你的输入</h4>
                  <div class="space-y-2">
                    <div v-for="input in analysis.inputs" :key="input.key">
                      <span class="text-xs text-gray-500">{{ input.label }}:</span>
                      <p class="text-sm text-gray-700">{{ input.value }}</p>
                    </div>
                  </div>
                </div>

                <!-- 分析结果 -->
                <div v-if="analysis.results && analysis.results.length > 0">
                  <h4 class="font-medium text-sm text-purple-600 mb-2">分析结果</h4>
                  <div class="space-y-3">
                    <div
                      v-for="result in analysis.results"
                      :key="result.key"
                      class="border border-purple-100 rounded-lg p-3"
                    >
                      <h5 class="font-medium text-sm mb-1">{{ result.label }}</h5>
                      <p class="text-sm text-gray-700">{{ result.content }}</p>
                    </div>
                  </div>
                </div>

                <!-- 分析中状态 -->
                <div v-else-if="analysis.status === 'processing'" class="text-center py-4">
                  <div class="text-purple-600">
                    <svg class="h-6 w-6 mx-auto mb-2 animate-spin" fill="none" viewBox="0 0 24 24">
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    <span class="text-sm">分析进行中，请稍候...</span>
                  </div>
                </div>

                <!-- 失败状态 -->
                <div v-else-if="analysis.status === 'failed'" class="text-center py-4 text-red-500">
                  <span class="text-sm">分析失败，请重试</span>
                </div>

                <!-- 等待状态 -->
                <div v-else class="text-center py-4 text-gray-400">
                  <span class="text-sm">等待分析...</span>
                </div>
              </div>
            </ElCollapseItem>
          </ElCollapse>
        </ElCard>
      </div>

      <!-- 右侧：操作面板 -->
      <div class="space-y-6">
        <!-- 思维模型 -->
        <ElCard shadow="never">
          <template #header>
            <h3 class="font-medium">思维模型</h3>
          </template>

          <div v-if="topic.modelId" class="text-center py-4">
            <div class="text-2xl mb-2">🧠</div>
            <div class="font-medium text-lg mb-1">{{ topic.modelName }}</div>
            <div class="text-sm text-gray-500 mb-4">已选用</div>
            <ElButton type="primary" class="w-full" @click="openModelDialog">
              更换模型
            </ElButton>
          </div>

          <div v-else class="text-center py-6">
            <div class="text-4xl mb-3">🤔</div>
            <div class="text-gray-500 mb-4">还没有选用思维模型</div>
            <ElButton type="primary" class="w-full mb-2" @click="openModelDialog">
              选用模型
            </ElButton>
            <ElButton link type="primary" @click="goToMarket">
              去市场浏览 →
            </ElButton>
          </div>
        </ElCard>

        <!-- 快速操作 -->
        <ElCard shadow="never">
          <template #header>
            <h3 class="font-medium">快速操作</h3>
          </template>

          <div class="space-y-3">
            <ElButton class="w-full" @click="openAnalysisDialog" :disabled="!topic.modelId">
              <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>
              </svg>
              开始分析
            </ElButton>
            <ElButton class="w-full" @click="goToMarket">
              <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z"/>
              </svg>
              浏览更多模型
            </ElButton>
            <ElButton class="w-full" @click="goBack">
              <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 10h16M4 14h16M4 18h16"/>
              </svg>
              返回课题列表
            </ElButton>
          </div>
        </ElCard>
      </div>
    </div>

    <!-- 错误状态 -->
    <ElEmpty v-else description="课题不存在或已被删除">
      <ElButton @click="goBack">返回列表</ElButton>
    </ElEmpty>

    <!-- 选用模型对话框 -->
    <ElDialog
      v-model="showModelDialog"
      title="选用思维模型"
      width="600px"
    >
      <p class="text-gray-500 mb-4">选择一个适合本课题的思维模型，这将帮助你更系统地分析问题。</p>
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
          <div class="flex items-center justify-between">
            <span>{{ model.title }}</span>
            <ElTag size="small" :type="model.isFree ? 'success' : 'warning'">
              {{ model.isFree ? '免费' : `¥${model.price}` }}
            </ElTag>
          </div>
        </ElOption>
      </ElSelect>

      <template #footer>
        <ElButton @click="showModelDialog = false">取消</ElButton>
        <ElButton type="primary" @click="handleSelectModel">确定</ElButton>
      </template>
    </ElDialog>

    <!-- 分析对话框 -->
    <ElDialog
      v-model="showAnalysisDialog"
      title="使用思维模型进行分析"
      width="700px"
    >
      <div v-if="analysisTemplate">
        <div class="bg-blue-50 rounded-lg p-4 mb-4">
          <div class="flex items-start gap-3">
            <svg class="h-5 w-5 text-blue-500 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
            <div>
              <p class="font-medium text-blue-900">{{ analysisTemplate.modelName }}</p>
              <p class="text-sm text-blue-700 mt-1">请根据提示填写各项内容，系统将根据你的输入生成分析结果。</p>
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
              :placeholder="field.placeholder || `请输入${field.label}`"
            />
            <ElInput
              v-else
              v-model="analysisForm[field.key]"
              :placeholder="field.placeholder || `请输入${field.label}`"
            />
          </ElFormItem>
        </ElForm>
      </div>

      <template #footer>
        <ElButton @click="showAnalysisDialog = false">取消</ElButton>
        <ElButton
          type="primary"
          :loading="submittingAnalysis"
          @click="handleSubmitAnalysis"
        >
          提交分析
        </ElButton>
      </template>
    </ElDialog>
  </Page>
</template>
