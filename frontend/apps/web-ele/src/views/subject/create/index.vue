<script lang="ts" setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';

import { Page } from '@vben/common-ui';

import {
  ElButton,
  ElCard,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElSteps,
  ElStep,
} from 'element-plus';

import { createTopicApi, type TopicApi } from '#/api';

// 路由
const router = useRouter();

// 当前步骤
const currentStep = ref(0);

// 表单数据
const formData = ref<TopicApi.CreateTopicParams>({
  title: '',
  description: '',
});

// 表单验证规则
const rules = {
  title: [
    { required: true, message: '请输入课题标题', trigger: 'blur' },
    { min: 2, max: 100, message: '标题长度在 2 到 100 个字符之间', trigger: 'blur' },
  ],
  description: [
    { required: true, message: '请输入课题描述', trigger: 'blur' },
    { min: 10, max: 2000, message: '描述长度在 10 到 2000 个字符之间', trigger: 'blur' },
  ],
};

// 表单引用
const formRef = ref<InstanceType<typeof ElForm>>();

// 提交加载状态
const submitting = ref(false);

// 创建课题
async function handleSubmit() {
  if (!formRef.value) return;

  await formRef.value.validate(async (valid) => {
    if (!valid) {
      ElMessage.warning('请填写完整信息');
      return;
    }

    submitting.value = true;
    try {
      const res = await createTopicApi(formData.value);
      ElMessage.success('课题创建成功！');
      // 跳转到课题详情页
      router.push(`/my-topics/${res.id}`);
    } catch (error) {
      console.error('创建课题失败:', error);
      ElMessage.error('创建课题失败，请重试');
    } finally {
      submitting.value = false;
    }
  });
}

// 下一步
async function nextStep() {
  if (!formRef.value) return;

  if (currentStep.value === 0) {
    await formRef.value.validateField(['title'], (valid) => {
      if (valid) {
        currentStep.value++;
      }
    });
  }
}

// 上一步
function prevStep() {
  currentStep.value--;
}

// 取消创建
function handleCancel() {
  router.push('/my-topics');
}

// 步骤列表
const steps = [
  { title: '基本信息', description: '填写课题标题' },
  { title: '详细描述', description: '描述课题背景和目标' },
  { title: '确认创建', description: '检查并提交' },
];
</script>

<template>
  <Page
    description="创建一个新课题，开始你的深度思考之旅"
    title="创建课题"
  >
    <div class="max-w-4xl mx-auto">
      <!-- 步骤条 -->
      <ElCard class="mb-6" shadow="never">
        <ElSteps :active="currentStep" finish-status="success" align-center>
          <ElStep
            v-for="step in steps"
            :key="step.title"
            :title="step.title"
            :description="step.description"
          />
        </ElSteps>
      </ElCard>

      <!-- 表单内容 -->
      <ElCard shadow="never">
        <ElForm
          ref="formRef"
          :model="formData"
          :rules="rules"
          label-position="top"
          class="max-w-2xl mx-auto"
        >
          <!-- 第一步：基本信息 -->
          <div v-show="currentStep === 0">
            <h3 class="text-lg font-medium mb-4">给你的课题起个名字</h3>
            <p class="text-gray-500 mb-6">一个好的标题能帮助你更好地聚焦思考的方向。</p>

            <ElFormItem label="课题标题" prop="title">
              <ElInput
                v-model="formData.title"
                placeholder="例如：如何提高团队的创新能力？"
                size="large"
              />
            </ElFormItem>

            <div class="flex justify-end gap-3 mt-8">
              <ElButton @click="handleCancel">取消</ElButton>
              <ElButton type="primary" @click="nextStep">下一步</ElButton>
            </div>
          </div>

          <!-- 第二步：详细描述 -->
          <div v-show="currentStep === 1">
            <h3 class="text-lg font-medium mb-4">详细描述你的课题</h3>
            <p class="text-gray-500 mb-6">背景信息越充分，思维模型的分析效果越好。</p>

            <ElFormItem label="课题描述" prop="description">
              <ElInput
                v-model="formData.description"
                type="textarea"
                :rows="8"
                placeholder="描述课题的背景、目标、相关约束条件等...&#10;&#10;例如：&#10;- 背景：团队目前的创新速度跟不上市场需求&#10;- 目标：在3个月内将新想法的落地速度提升50%&#10;- 约束：预算有限，需要利用现有资源"
              />
            </ElFormItem>

            <div class="flex justify-end gap-3 mt-8">
              <ElButton @click="prevStep">上一步</ElButton>
              <ElButton type="primary" @click="currentStep++">下一步</ElButton>
            </div>
          </div>

          <!-- 第三步：确认 -->
          <div v-show="currentStep === 2">
            <h3 class="text-lg font-medium mb-4">确认课题信息</h3>
            <p class="text-gray-500 mb-6">确认无误后，点击创建按钮。</p>

            <div class="bg-gray-50 rounded-lg p-6 mb-6">
              <div class="mb-4">
                <label class="text-sm text-gray-500 block mb-1">课题标题</label>
                <div class="font-medium text-lg">{{ formData.title }}</div>
              </div>
              <div>
                <label class="text-sm text-gray-500 block mb-1">课题描述</label>
                <div class="text-gray-700 whitespace-pre-wrap">{{ formData.description }}</div>
              </div>
            </div>

            <div class="bg-blue-50 rounded-lg p-4 mb-6">
              <div class="flex items-start gap-3">
                <svg class="h-5 w-5 text-blue-500 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
                <div class="text-sm text-blue-700">
                  <p class="font-medium mb-1">创建后你可以：</p>
                  <ul class="list-disc list-inside space-y-1">
                    <li>选用合适的思维模型进行分析</li>
                    <li>多次分析，对比不同思路</li>
                    <li>保存分析结果，形成完整的思考记录</li>
                  </ul>
                </div>
              </div>
            </div>

            <div class="flex justify-end gap-3 mt-8">
              <ElButton @click="prevStep">上一步</ElButton>
              <ElButton type="primary" :loading="submitting" @click="handleSubmit">
                创建课题
              </ElButton>
            </div>
          </div>
        </ElForm>
      </ElCard>

      <!-- 提示信息 -->
      <div class="mt-6 text-center text-sm text-gray-400">
        <p>创建课题后，你可以在详情页选用思维模型进行深入分析</p>
      </div>
    </div>
  </Page>
</template>
