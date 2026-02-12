<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { Page } from '@vben/common-ui';
import {
  ElCard,
  ElButton,
  ElInput,
  ElDatePicker,
  ElInputNumber,
  ElMessage,
} from 'element-plus';
import { createConsultationApi, type CollisionApi } from '#/api/collision';

const router = useRouter();

const form = ref<{
  title: string;
  field: CollisionApi.ConsultationField | '';
  mode: CollisionApi.ConsultationMode;
  description: string;
  background: string;
  expectation: string;
  reward: number;
  deadline: string;
}>({
  title: '',
  field: '',
  mode: 'both',
  description: '',
  background: '',
  expectation: '',
  reward: 100,
  deadline: '',
});

const submitting = ref(false);

const consultationFields = [
  { value: 'career', label: '💼 职业发展', desc: '职业规划、求职面试、晋升发展等' },
  { value: 'startup', label: '🚀 创业指导', desc: '商业模式、融资、团队搭建等' },
  { value: 'technology', label: '💻 技术架构', desc: '技术选型、架构设计、性能优化等' },
  { value: 'product', label: '📱 产品设计', desc: '产品规划、用户体验、增长策略等' },
  { value: 'investment', label: '💰 投资理财', desc: '投资策略、资产配置、风险管理等' },
  { value: 'management', label: '👥 团队管理', desc: '团队建设、组织发展、领导力等' },
  { value: 'psychology', label: '🧠 心理咨询', desc: '职场压力、人际关系、个人成长等' },
  { value: 'other', label: '📝 其他', desc: '其他类型的咨询需求' },
];

const modeOptions = [
  { value: 'online', label: '🖥️ 线上咨询', desc: '视频/语音/文字沟通' },
  { value: 'offline', label: '☕ 线下咨询', desc: '面对面沟通交流' },
  { value: 'both', label: '🔄 都可以', desc: '根据专家情况灵活安排' },
];

const rewardPresets = [50, 100, 200, 500, 1000, 2000];

const isFormValid = computed(() => {
  return (
    form.value.title.trim().length >= 5 &&
    form.value.field &&
    form.value.description.trim().length >= 20 &&
    form.value.reward >= 10 &&
    form.value.deadline
  );
});

const handleSubmit = async () => {
  if (!isFormValid.value) {
    ElMessage.warning('请完善必填信息');
    return;
  }

  submitting.value = true;
  try {
    await createConsultationApi({
      title: form.value.title,
      field: form.value.field as CollisionApi.ConsultationField,
      mode: form.value.mode,
      description: form.value.description,
      background: form.value.background,
      expectation: form.value.expectation,
      reward: form.value.reward,
      deadline: form.value.deadline,
    });
    ElMessage.success('咨询已发布成功！');
    router.push({ name: 'Collision', query: { tab: 'consultation' } });
  } catch (error) {
    ElMessage.error('发布失败，请重试');
  } finally {
    submitting.value = false;
  }
};

const goBack = () => {
  router.push({ name: 'Collision', query: { tab: 'consultation' } });
};

// 计算最早截止日期（今天之后）
const disabledDate = (date: Date) => {
  return date.getTime() < Date.now() - 8.64e7;
};
</script>

<template>
  <Page title="发布咨询" content-class="p-6 bg-gray-50">
    <!-- 返回按钮 -->
    <div class="mb-4">
      <ElButton text @click="goBack">
        ← 返回咨询列表
      </ElButton>
    </div>

    <div class="max-w-4xl mx-auto">
      <div class="flex gap-6">
        <!-- 左侧表单 -->
        <div class="flex-1 space-y-6">
          <!-- 基本信息 -->
          <ElCard shadow="hover" class="!rounded-xl">
            <template #header>
              <span class="font-semibold text-gray-700">💡 基本信息</span>
            </template>
            
            <div class="space-y-6">
              <!-- 咨询标题 -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  咨询标题 <span class="text-red-500">*</span>
                </label>
                <ElInput
                  v-model="form.title"
                  placeholder="用一句话描述您的问题，例如：如何在3年内成为技术专家？"
                  maxlength="50"
                  show-word-limit
                />
                <p class="text-xs text-gray-400 mt-1">标题需不少于5个字，简洁明了</p>
              </div>

              <!-- 咨询领域 -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  咨询领域 <span class="text-red-500">*</span>
                </label>
                <div class="grid grid-cols-2 gap-3">
                  <div
                    v-for="field in consultationFields"
                    :key="field.value"
                    class="p-3 rounded-xl border-2 cursor-pointer transition-all"
                    :class="[
                      form.field === field.value
                        ? 'border-purple-500 bg-purple-50'
                        : 'border-gray-200 hover:border-purple-300 hover:bg-gray-50'
                    ]"
                    @click="form.field = field.value as CollisionApi.ConsultationField"
                  >
                    <div class="font-medium text-gray-800">{{ field.label }}</div>
                    <div class="text-xs text-gray-500 mt-1">{{ field.desc }}</div>
                  </div>
                </div>
              </div>

              <!-- 咨询方式 -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  咨询方式
                </label>
                <div class="flex gap-4">
                  <div
                    v-for="mode in modeOptions"
                    :key="mode.value"
                    class="flex-1 p-4 rounded-xl border-2 cursor-pointer transition-all text-center"
                    :class="[
                      form.mode === mode.value
                        ? 'border-purple-500 bg-purple-50'
                        : 'border-gray-200 hover:border-purple-300'
                    ]"
                    @click="form.mode = mode.value as CollisionApi.ConsultationMode"
                  >
                    <div class="text-lg mb-1">{{ mode.label }}</div>
                    <div class="text-xs text-gray-500">{{ mode.desc }}</div>
                  </div>
                </div>
              </div>
            </div>
          </ElCard>

          <!-- 详细描述 -->
          <ElCard shadow="hover" class="!rounded-xl">
            <template #header>
              <span class="font-semibold text-gray-700">📝 详细描述</span>
            </template>
            
            <div class="space-y-6">
              <!-- 问题描述 -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  问题描述 <span class="text-red-500">*</span>
                </label>
                <ElInput
                  v-model="form.description"
                  type="textarea"
                  :rows="5"
                  placeholder="详细描述您遇到的问题或困惑..."
                  maxlength="1000"
                  show-word-limit
                />
                <p class="text-xs text-gray-400 mt-1">至少20字，描述越详细，越容易获得精准帮助</p>
              </div>

              <!-- 背景信息 -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  背景信息
                </label>
                <ElInput
                  v-model="form.background"
                  type="textarea"
                  :rows="3"
                  placeholder="可以补充您的相关背景，如工作经历、当前状况等，帮助专家更好地了解您..."
                  maxlength="500"
                  show-word-limit
                />
              </div>

              <!-- 期望获得 -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  期望获得
                </label>
                <ElInput
                  v-model="form.expectation"
                  type="textarea"
                  :rows="3"
                  placeholder="描述您期望从这次咨询中获得什么，例如具体建议、行动计划、资源推荐等..."
                  maxlength="500"
                  show-word-limit
                />
              </div>
            </div>
          </ElCard>

          <!-- 悬赏设置 -->
          <ElCard shadow="hover" class="!rounded-xl">
            <template #header>
              <span class="font-semibold text-gray-700">💰 悬赏设置</span>
            </template>
            
            <div class="space-y-6">
              <!-- 悬赏金额 -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  悬赏金额 <span class="text-red-500">*</span>
                </label>
                <div class="flex flex-wrap gap-3 mb-3">
                  <button
                    v-for="preset in rewardPresets"
                    :key="preset"
                    class="px-4 py-2 rounded-full text-sm transition-all"
                    :class="[
                      form.reward === preset
                        ? 'bg-orange-500 text-white'
                        : 'bg-gray-100 text-gray-600 hover:bg-orange-100 hover:text-orange-600'
                    ]"
                    @click="form.reward = preset"
                  >
                    ¥{{ preset }}
                  </button>
                </div>
                <div class="flex items-center gap-2">
                  <span class="text-gray-500">自定义金额:</span>
                  <ElInputNumber
                    v-model="form.reward"
                    :min="10"
                    :max="100000"
                    :step="10"
                    controls-position="right"
                  />
                  <span class="text-gray-500">元</span>
                </div>
                <p class="text-xs text-gray-400 mt-2">
                  💡 合理的悬赏金额能吸引更多优质专家申请
                </p>
              </div>

              <!-- 截止日期 -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  截止日期 <span class="text-red-500">*</span>
                </label>
                <ElDatePicker
                  v-model="form.deadline"
                  type="date"
                  placeholder="选择截止日期"
                  :disabled-date="disabledDate"
                  value-format="YYYY-MM-DD"
                  class="!w-full"
                />
                <p class="text-xs text-gray-400 mt-1">
                  建议设置7-14天的征集期
                </p>
              </div>
            </div>
          </ElCard>

          <!-- 提交按钮 -->
          <div class="flex justify-end gap-4">
            <ElButton size="large" @click="goBack">取消</ElButton>
            <ElButton
              type="primary"
              size="large"
              :loading="submitting"
              :disabled="!isFormValid"
              class="!bg-purple-600 !border-purple-600 hover:!bg-purple-700 !rounded-full !px-8"
              @click="handleSubmit"
            >
              💡 发布咨询
            </ElButton>
          </div>
        </div>

        <!-- 右侧提示 -->
        <div class="w-72 flex-shrink-0 space-y-6 hidden lg:block">
          <!-- 发布指南 -->
          <ElCard shadow="hover" class="!rounded-xl">
            <template #header>
              <span class="font-semibold text-gray-700">📋 发布指南</span>
            </template>
            <ul class="text-sm text-gray-600 space-y-3">
              <li class="flex items-start gap-2">
                <span class="text-purple-500 font-bold">1</span>
                <span>标题要简洁明了，突出核心问题</span>
              </li>
              <li class="flex items-start gap-2">
                <span class="text-purple-500 font-bold">2</span>
                <span>选择准确的咨询领域，方便专家匹配</span>
              </li>
              <li class="flex items-start gap-2">
                <span class="text-purple-500 font-bold">3</span>
                <span>详细描述问题背景和期望，提高回复质量</span>
              </li>
              <li class="flex items-start gap-2">
                <span class="text-purple-500 font-bold">4</span>
                <span>合理设置悬赏金额，吸引优质专家</span>
              </li>
              <li class="flex items-start gap-2">
                <span class="text-purple-500 font-bold">5</span>
                <span>及时查看申请，与专家保持沟通</span>
              </li>
            </ul>
          </ElCard>

          <!-- 费用说明 -->
          <ElCard shadow="hover" class="!rounded-xl">
            <template #header>
              <span class="font-semibold text-gray-700">💰 费用说明</span>
            </template>
            <div class="text-sm text-gray-600 space-y-2">
              <p>• 悬赏金额将在发布后冻结</p>
              <p>• 专家完成咨询后，悬赏金额将转给专家</p>
              <p>• 若无人应邀或取消，悬赏金额将退还</p>
              <p>• 平台不收取额外服务费</p>
            </div>
          </ElCard>

          <!-- 预览 -->
          <ElCard shadow="hover" class="!rounded-xl !bg-gradient-to-br from-purple-50 to-purple-100">
            <template #header>
              <span class="font-semibold text-purple-700">👁️ 预览效果</span>
            </template>
            <div class="space-y-3">
              <div class="flex items-center gap-2">
                <div class="w-12 h-12 rounded-lg bg-gradient-to-br from-orange-400 to-red-500 flex items-center justify-center text-white text-xs">
                  <div class="text-center">
                    <div class="text-[10px]">悬赏</div>
                    <div class="font-bold">¥{{ form.reward }}</div>
                  </div>
                </div>
                <div class="flex-1 min-w-0">
                  <div class="font-medium text-gray-800 text-sm line-clamp-1">
                    {{ form.title || '咨询标题' }}
                  </div>
                  <div class="text-xs text-gray-500">
                    {{ form.field ? consultationFields.find(f => f.value === form.field)?.label : '领域' }}
                  </div>
                </div>
              </div>
              <p class="text-xs text-gray-600 line-clamp-2">
                {{ form.description || '问题描述将显示在这里...' }}
              </p>
            </div>
          </ElCard>
        </div>
      </div>
    </div>
  </Page>
</template>
