<script lang="ts" setup>
import { ref, reactive, computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { Page } from '@vben/common-ui';

import {
  ElButton,
  ElCard,
  ElInput,
  ElForm,
  ElFormItem,
  ElTag,
  ElMessage,
} from 'element-plus';
import type { FormInstance } from 'element-plus';

import { createTopicApi, type TopicApi } from '#/api';

// 路由
const route = useRoute();
const router = useRouter();
const editId = computed(() => route.query.id as string | undefined);
const isEdit = computed(() => !!editId.value);

// 当前步骤
const currentStep = ref(0);
const steps = [
  { id: 0, label: '基本信息', icon: '📋', description: '课题标题与背景' },
  { id: 1, label: '详细描述', icon: '📝', description: '目标与约束条件' },
  { id: 2, label: '确认创建', icon: '🚀', description: '检查并提交' },
];

// 表单数据
const form = reactive({
  title: '',
  description: '',
  background: '',
  goal: '',
  constraints: '',
  tags: [] as string[],
});

// 表单引用
const formRef = ref<FormInstance>();

// 表单校验规则
const rules = {
  title: [
    { required: true, message: '请输入课题标题', trigger: 'blur' },
    { min: 2, max: 100, message: '标题长度在 2 到 100 个字符', trigger: 'blur' },
  ],
};

// 推荐标签
const suggestedTags = ['商业决策', '产品规划', '团队管理', '职业发展', '创新思维', '问题解决', '战略分析', '个人成长'];

// 标签输入
const tagInput = ref('');
const tagInputVisible = ref(false);

// 提交状态
const submitting = ref(false);

// 加载编辑数据
onMounted(async () => {
  if (isEdit.value) {
    // 模拟加载编辑数据
    await new Promise(resolve => setTimeout(resolve, 500));
    form.title = '如何提高团队创新能力';
    form.background = '团队目前的创新速度跟不上市场需求';
    form.goal = '在3个月内将新想法的落地速度提升50%';
    form.constraints = '预算有限，需要利用现有资源';
    form.tags = ['团队管理', '创新思维'];
  }
});

// 标签操作
function handleAddTag() {
  const tag = tagInput.value.trim();
  if (tag && !form.tags.includes(tag) && form.tags.length < 5) {
    form.tags.push(tag);
  }
  tagInput.value = '';
  tagInputVisible.value = false;
}

function handleRemoveTag(tag: string) {
  form.tags = form.tags.filter(t => t !== tag);
}

function addSuggestedTag(tag: string) {
  if (!form.tags.includes(tag) && form.tags.length < 5) {
    form.tags.push(tag);
  }
}

// 步骤导航
async function goToStep(step: number) {
  if (step < currentStep.value) {
    currentStep.value = step;
    return;
  }
  
  // 验证当前步骤
  if (currentStep.value === 0) {
    try {
      await formRef.value?.validateField(['title']);
    } catch {
      ElMessage.warning('请填写课题标题');
      return;
    }
  }
  
  currentStep.value = step;
}

function prevStep() {
  if (currentStep.value > 0) {
    currentStep.value--;
  }
}

async function nextStep() {
  await goToStep(currentStep.value + 1);
}

// 提交
async function handleSubmit() {
  submitting.value = true;
  try {
    // 组装描述
    const description = [
      form.background && `背景：${form.background}`,
      form.goal && `目标：${form.goal}`,
      form.constraints && `约束：${form.constraints}`,
    ].filter(Boolean).join('\n\n');
    
    const params: TopicApi.CreateTopicParams = {
      title: form.title,
      description: description || form.title,
    };
    
    const res = await createTopicApi(params);
    ElMessage.success(isEdit.value ? '课题已更新' : '课题创建成功！');
    router.push(`/my-topics/${res.id}`);
  } catch (error) {
    console.error('提交失败:', error);
    ElMessage.error('提交失败，请重试');
  } finally {
    submitting.value = false;
  }
}

// 保存草稿
async function handleSaveDraft() {
  ElMessage.success('草稿已保存');
}

// 取消
function handleCancel() {
  router.back();
}

// 计算完成进度
const formProgress = computed(() => {
  let filled = 0;
  let total = 5;
  
  if (form.title) filled++;
  if (form.background) filled++;
  if (form.goal) filled++;
  if (form.constraints) filled++;
  if (form.tags.length > 0) filled++;
  
  return Math.round((filled / total) * 100);
});

// 生成完整描述预览
const fullDescription = computed(() => {
  const parts = [];
  if (form.background) parts.push(`【背景】\n${form.background}`);
  if (form.goal) parts.push(`【目标】\n${form.goal}`);
  if (form.constraints) parts.push(`【约束条件】\n${form.constraints}`);
  return parts.join('\n\n') || '暂无描述...';
});
</script>

<template>
  <Page
    :description="isEdit ? '修改课题内容' : '创建一个新课题，开始你的深度思考之旅'"
    :title="isEdit ? '编辑课题' : '创建课题'"
    content-class="p-6 bg-gray-50"
  >
    <div class="flex gap-6">
      <!-- 左侧主表单 -->
      <div class="flex-1 min-w-0 space-y-6">
        <!-- 步骤导航 -->
        <ElCard shadow="hover" class="!rounded-xl">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-4">
              <button
                v-for="(step, index) in steps"
                :key="step.id"
                class="flex items-center gap-3 px-4 py-3 rounded-xl transition-all"
                :class="[
                  currentStep === index
                    ? 'bg-emerald-100 border-2 border-emerald-300'
                    : index < currentStep
                    ? 'bg-green-50 border border-green-200'
                    : 'bg-gray-50 border border-gray-200 hover:border-emerald-200'
                ]"
                @click="goToStep(index)"
              >
                <div
                  class="w-10 h-10 rounded-full flex items-center justify-center text-lg"
                  :class="[
                    currentStep === index
                      ? 'bg-emerald-600 text-white'
                      : index < currentStep
                      ? 'bg-green-500 text-white'
                      : 'bg-gray-200 text-gray-500'
                  ]"
                >
                  <span v-if="index < currentStep">✓</span>
                  <span v-else>{{ step.icon }}</span>
                </div>
                <div class="text-left">
                  <div class="font-semibold" :class="currentStep === index ? 'text-emerald-700' : 'text-gray-700'">
                    {{ step.label }}
                  </div>
                  <div class="text-xs text-gray-400">{{ step.description }}</div>
                </div>
              </button>
            </div>
            <div class="text-right">
              <div class="text-sm text-gray-500">完成度</div>
              <div class="text-2xl font-bold text-emerald-600">{{ formProgress }}%</div>
            </div>
          </div>
        </ElCard>

        <!-- 步骤 1: 基本信息 -->
        <ElCard v-if="currentStep === 0" shadow="hover" class="!rounded-xl">
          <template #header>
            <div class="flex items-center gap-2">
              <span class="text-lg">📋</span>
              <span class="font-semibold text-gray-700">基本信息</span>
            </div>
          </template>
          
          <ElForm ref="formRef" :model="form" :rules="rules" label-position="top" class="max-w-2xl">
            <!-- 课题标题 -->
            <ElFormItem label="课题标题" prop="title">
              <ElInput
                v-model="form.title"
                placeholder="用一句话概括你要思考的问题，例如：如何提高团队的创新能力？"
                maxlength="100"
                show-word-limit
                size="large"
                class="!rounded-lg"
              />
            </ElFormItem>

            <!-- 背景说明 -->
            <ElFormItem label="背景说明">
              <ElInput
                v-model="form.background"
                type="textarea"
                :rows="4"
                placeholder="描述这个课题产生的原因和背景...&#10;例如：团队目前的创新速度跟不上市场需求，需要找到提升方法"
              />
            </ElFormItem>

            <!-- 标签 -->
            <ElFormItem label="课题标签">
              <div class="space-y-3">
                <div class="flex flex-wrap gap-2">
                  <ElTag
                    v-for="tag in form.tags"
                    :key="tag"
                    closable
                    effect="plain"
                    class="!bg-emerald-50 !text-emerald-600 !border-emerald-200 !rounded-full"
                    @close="handleRemoveTag(tag)"
                  >
                    {{ tag }}
                  </ElTag>
                  <ElInput
                    v-if="tagInputVisible"
                    v-model="tagInput"
                    size="small"
                    class="!w-24"
                    @keyup.enter="handleAddTag"
                    @blur="handleAddTag"
                  />
                  <ElButton
                    v-else-if="form.tags.length < 5"
                    size="small"
                    class="!rounded-full"
                    @click="tagInputVisible = true"
                  >
                    + 添加标签
                  </ElButton>
                </div>
                <div class="flex flex-wrap gap-2">
                  <span class="text-xs text-gray-400 mr-2">推荐：</span>
                  <button
                    v-for="tag in suggestedTags.filter(t => !form.tags.includes(t))"
                    :key="tag"
                    type="button"
                    class="px-2 py-0.5 text-xs bg-gray-100 text-gray-500 rounded-full hover:bg-emerald-100 hover:text-emerald-600 transition-colors"
                    @click="addSuggestedTag(tag)"
                  >
                    + {{ tag }}
                  </button>
                </div>
              </div>
            </ElFormItem>
          </ElForm>
        </ElCard>

        <!-- 步骤 2: 详细描述 -->
        <template v-if="currentStep === 1">
          <ElCard shadow="hover" class="!rounded-xl">
            <template #header>
              <div class="flex items-center gap-2">
                <span class="text-lg">🎯</span>
                <span class="font-semibold text-gray-700">目标设定</span>
              </div>
            </template>
            <p class="text-sm text-gray-500 mb-4">明确你想要达成的具体目标</p>
            <ElInput
              v-model="form.goal"
              type="textarea"
              :rows="4"
              placeholder="描述你期望的结果...&#10;例如：在3个月内将新想法的落地速度提升50%"
            />
          </ElCard>

          <ElCard shadow="hover" class="!rounded-xl">
            <template #header>
              <div class="flex items-center gap-2">
                <span class="text-lg">⚠️</span>
                <span class="font-semibold text-gray-700">约束条件</span>
              </div>
            </template>
            <p class="text-sm text-gray-500 mb-4">列出需要考虑的限制因素</p>
            <ElInput
              v-model="form.constraints"
              type="textarea"
              :rows="4"
              placeholder="描述资源、时间、预算等限制...&#10;例如：预算有限，需要利用现有资源"
            />
          </ElCard>
        </template>

        <!-- 步骤 3: 确认创建 -->
        <ElCard v-if="currentStep === 2" shadow="hover" class="!rounded-xl">
          <template #header>
            <div class="flex items-center gap-2">
              <span class="text-lg">🚀</span>
              <span class="font-semibold text-gray-700">确认课题信息</span>
            </div>
          </template>

          <div class="space-y-6">
            <!-- 预览 -->
            <div class="p-6 rounded-xl bg-gradient-to-r from-emerald-50 to-teal-50 border border-emerald-100">
              <h3 class="text-xl font-bold text-gray-800 mb-4">{{ form.title || '未设置标题' }}</h3>
              
              <div class="flex flex-wrap gap-2 mb-4">
                <ElTag
                  v-for="tag in form.tags"
                  :key="tag"
                  effect="plain"
                  class="!bg-white !text-emerald-600 !border-emerald-200 !rounded-full"
                >
                  {{ tag }}
                </ElTag>
                <span v-if="form.tags.length === 0" class="text-sm text-gray-400">暂无标签</span>
              </div>
              
              <div class="text-gray-600 whitespace-pre-line leading-relaxed">
                {{ fullDescription }}
              </div>
            </div>

            <!-- 提示 -->
            <div class="p-5 rounded-xl bg-gradient-to-r from-blue-50 to-indigo-50 border border-blue-100">
              <h4 class="font-semibold text-blue-800 mb-3 flex items-center gap-2">
                <span>💡</span> 创建后你可以
              </h4>
              <ul class="text-sm text-blue-700 space-y-2">
                <li class="flex items-start gap-2">
                  <span class="text-blue-500">•</span>
                  <span>选用合适的思维模型进行深入分析</span>
                </li>
                <li class="flex items-start gap-2">
                  <span class="text-blue-500">•</span>
                  <span>多次分析，对比不同思路和视角</span>
                </li>
                <li class="flex items-start gap-2">
                  <span class="text-blue-500">•</span>
                  <span>保存分析结果，形成完整的思考记录</span>
                </li>
              </ul>
            </div>
          </div>
        </ElCard>

        <!-- 底部操作栏 -->
        <ElCard shadow="hover" class="!rounded-xl">
          <div class="flex items-center justify-between">
            <div>
              <ElButton v-if="currentStep === 0" class="!rounded-full" @click="handleCancel">
                取消
              </ElButton>
              <ElButton v-else class="!rounded-full" @click="prevStep">
                <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
                </svg>
                上一步
              </ElButton>
            </div>
            <div class="flex items-center gap-3">
              <ElButton class="!rounded-full" @click="handleSaveDraft">
                <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-3m-1 4l-3 3m0 0l-3-3m3 3V4"/>
                </svg>
                保存草稿
              </ElButton>
              <ElButton
                v-if="currentStep < 2"
                type="primary"
                class="!bg-emerald-600 !border-emerald-600 hover:!bg-emerald-700 !rounded-full"
                @click="nextStep"
              >
                下一步
                <svg class="h-4 w-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
                </svg>
              </ElButton>
              <ElButton
                v-else
                type="primary"
                :loading="submitting"
                class="!bg-emerald-600 !border-emerald-600 hover:!bg-emerald-700 !rounded-full"
                @click="handleSubmit"
              >
                <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
                </svg>
                {{ isEdit ? '保存修改' : '创建课题' }}
              </ElButton>
            </div>
          </div>
        </ElCard>
      </div>

      <!-- 右侧边栏 -->
      <div class="w-80 flex-shrink-0 space-y-6 hidden lg:block">
        <!-- 预览卡片 -->
        <ElCard shadow="hover" class="!rounded-xl">
          <template #header>
            <div class="flex items-center gap-2">
              <span class="text-lg">👁️</span>
              <span class="font-semibold text-gray-700">课题预览</span>
            </div>
          </template>
          <div class="space-y-4">
            <div class="h-24 rounded-lg overflow-hidden bg-gradient-to-br from-emerald-100 to-teal-100 flex items-center justify-center">
              <span class="text-4xl">🧩</span>
            </div>
            <div>
              <h4 class="font-semibold text-gray-800 line-clamp-2">
                {{ form.title || '课题标题' }}
              </h4>
              <p class="text-sm text-gray-500 mt-2 line-clamp-3">
                {{ form.background || '课题背景将显示在这里...' }}
              </p>
            </div>
            <div class="flex flex-wrap gap-1">
              <span
                v-for="tag in form.tags.slice(0, 3)"
                :key="tag"
                class="px-2 py-0.5 text-xs bg-emerald-100 text-emerald-600 rounded-full"
              >
                {{ tag }}
              </span>
              <span v-if="form.tags.length === 0" class="text-xs text-gray-400">暂无标签</span>
            </div>
          </div>
        </ElCard>

        <!-- 创作指南 -->
        <ElCard shadow="hover" class="!rounded-xl !bg-gradient-to-br from-emerald-50 to-teal-50 !border-emerald-100">
          <template #header>
            <div class="flex items-center gap-2">
              <span class="text-lg">📚</span>
              <span class="font-semibold text-emerald-700">创作指南</span>
            </div>
          </template>
          <div class="space-y-4">
            <div class="flex items-start gap-3">
              <div class="w-8 h-8 rounded-full bg-emerald-200 text-emerald-700 flex items-center justify-center font-bold text-sm flex-shrink-0">1</div>
              <div>
                <div class="font-medium text-gray-700 text-sm">明确问题</div>
                <div class="text-xs text-gray-500">用简洁的语言描述核心问题</div>
              </div>
            </div>
            <div class="flex items-start gap-3">
              <div class="w-8 h-8 rounded-full bg-emerald-200 text-emerald-700 flex items-center justify-center font-bold text-sm flex-shrink-0">2</div>
              <div>
                <div class="font-medium text-gray-700 text-sm">充分背景</div>
                <div class="text-xs text-gray-500">提供足够的上下文信息</div>
              </div>
            </div>
            <div class="flex items-start gap-3">
              <div class="w-8 h-8 rounded-full bg-emerald-200 text-emerald-700 flex items-center justify-center font-bold text-sm flex-shrink-0">3</div>
              <div>
                <div class="font-medium text-gray-700 text-sm">设定目标</div>
                <div class="text-xs text-gray-500">清晰的目标让分析更有方向</div>
              </div>
            </div>
          </div>
        </ElCard>

        <!-- 小贴士 -->
        <ElCard shadow="hover" class="!rounded-xl !bg-gradient-to-br from-amber-50 to-orange-50 !border-amber-100">
          <template #header>
            <div class="flex items-center gap-2">
              <span class="text-lg">💡</span>
              <span class="font-semibold text-amber-700">小贴士</span>
            </div>
          </template>
          <ul class="text-sm text-amber-800 space-y-2">
            <li class="flex items-start gap-2">
              <span class="text-amber-500">•</span>
              好的问题是成功分析的一半
            </li>
            <li class="flex items-start gap-2">
              <span class="text-amber-500">•</span>
              背景信息越丰富，分析越精准
            </li>
            <li class="flex items-start gap-2">
              <span class="text-amber-500">•</span>
              可以多次分析同一课题
            </li>
            <li class="flex items-start gap-2">
              <span class="text-amber-500">•</span>
              尝试不同模型获得多元视角
            </li>
          </ul>
        </ElCard>

        <!-- 推荐模型 -->
        <ElCard shadow="hover" class="!rounded-xl">
          <template #header>
            <div class="flex items-center gap-2">
              <span class="text-lg">🧠</span>
              <span class="font-semibold text-gray-700">推荐模型</span>
            </div>
          </template>
          <div class="space-y-3">
            <div
              class="p-3 rounded-lg hover:bg-gray-50 cursor-pointer transition-colors border border-gray-100"
              @click="router.push('/market/1')"
            >
              <div class="font-medium text-gray-800 text-sm">SWOT 分析</div>
              <div class="text-xs text-gray-400 mt-1">战略决策的经典工具</div>
            </div>
            <div
              class="p-3 rounded-lg hover:bg-gray-50 cursor-pointer transition-colors border border-gray-100"
              @click="router.push('/market/2')"
            >
              <div class="font-medium text-gray-800 text-sm">第一性原理</div>
              <div class="text-xs text-gray-400 mt-1">回归本质思考问题</div>
            </div>
            <div
              class="p-3 rounded-lg hover:bg-gray-50 cursor-pointer transition-colors border border-gray-100"
              @click="router.push('/market/3')"
            >
              <div class="font-medium text-gray-800 text-sm">5W1H 分析</div>
              <div class="text-xs text-gray-400 mt-1">全面剖析问题各方面</div>
            </div>
          </div>
        </ElCard>
      </div>
    </div>
  </Page>
</template>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.line-clamp-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
