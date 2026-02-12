<script lang="ts" setup>
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router';

import { Page } from '@vben/common-ui';

import {
  ElButton,
  ElCard,
  ElInput,
  ElSelect,
  ElOption,
  ElForm,
  ElFormItem,
  ElSwitch,
  ElTag,
  ElMessage,
  ElSteps,
  ElStep,
  ElUpload,
} from 'element-plus';
import type { UploadProps, UploadFile } from 'element-plus';

// 路由
const router = useRouter();

// 当前步骤
const currentStep = ref(0);

// 表单数据
const form = reactive({
  title: '',
  description: '',
  category: '',
  tags: [] as string[],
  cover: '',
  isFree: true,
  price: 0,
  content: {
    overview: '',
    steps: [
      { title: '', description: '' },
    ],
    examples: [
      { title: '', content: '' },
    ],
    tips: [] as string[],
  },
});

// 表单校验规则
const rules = {
  title: [
    { required: true, message: '请输入模型名称', trigger: 'blur' },
    { min: 2, max: 50, message: '名称长度在 2 到 50 个字符', trigger: 'blur' },
  ],
  description: [
    { required: true, message: '请输入模型描述', trigger: 'blur' },
    { min: 10, max: 500, message: '描述长度在 10 到 500 个字符', trigger: 'blur' },
  ],
  category: [
    { required: true, message: '请选择模型分类', trigger: 'change' },
  ],
};

// 分类选项
const categories = [
  { value: 'business', label: '商业管理' },
  { value: 'strategy', label: '战略规划' },
  { value: 'innovation', label: '创新思维' },
  { value: 'analysis', label: '分析工具' },
  { value: 'decision', label: '决策方法' },
  { value: 'creative', label: '创意构思' },
  { value: 'psychology', label: '心理学' },
  { value: 'communication', label: '沟通表达' },
];

// 标签输入
const tagInput = ref('');
const tagInputVisible = ref(false);

// 表单引用
const formRef = ref();

// 添加标签
function handleAddTag() {
  if (tagInput.value && !form.tags.includes(tagInput.value)) {
    form.tags.push(tagInput.value);
  }
  tagInput.value = '';
  tagInputVisible.value = false;
}

// 删除标签
function handleRemoveTag(tag: string) {
  form.tags = form.tags.filter(t => t !== tag);
}

// 添加步骤
function addStep() {
  form.content.steps.push({ title: '', description: '' });
}

// 删除步骤
function removeStep(index: number) {
  form.content.steps.splice(index, 1);
}

// 添加案例
function addExample() {
  form.content.examples.push({ title: '', content: '' });
}

// 删除案例
function removeExample(index: number) {
  form.content.examples.splice(index, 1);
}

// 上一步
function prevStep() {
  if (currentStep.value > 0) {
    currentStep.value--;
  }
}

// 下一步
async function nextStep() {
  if (currentStep.value === 0) {
    await formRef.value?.validate();
  }
  if (currentStep.value < 2) {
    currentStep.value++;
  }
}

// 提交表单
async function handleSubmit() {
  try {
    // TODO: 调用创建 API
    // await createModelApi(form);
    ElMessage.success('模型创建成功！');
    router.push('/my-models');
  } catch (error) {
    console.error('创建失败:', error);
    ElMessage.error('创建失败，请重试');
  }
}

// 保存草稿
async function handleSaveDraft() {
  try {
    // TODO: 调用保存草稿 API
    // await saveModelDraftApi(form);
    ElMessage.success('草稿已保存');
    router.push('/my-models');
  } catch (error) {
    console.error('保存失败:', error);
    ElMessage.error('保存失败');
  }
}

// 取消
function handleCancel() {
  router.back();
}

// 封面上传
const handleCoverChange: UploadProps['onChange'] = (uploadFile: UploadFile) => {
  if (uploadFile.raw) {
    const reader = new FileReader();
    reader.onload = (e) => {
      form.cover = e.target?.result as string;
    };
    reader.readAsDataURL(uploadFile.raw);
  }
};
</script>

<template>
  <Page
    description="创建并分享你的思维模型，帮助他人更好地思考"
    title="创建思维模型"
  >
    <ElCard shadow="never">
      <!-- 步骤条 -->
      <ElSteps :active="currentStep" finish-status="success" class="mb-8">
        <ElStep title="基本信息" description="填写模型的基本信息" />
        <ElStep title="详细内容" description="添加模型的使用步骤和案例" />
        <ElStep title="发布设置" description="设置价格和发布选项" />
      </ElSteps>

      <!-- 步骤 1: 基本信息 -->
      <div v-if="currentStep === 0">
        <ElForm
          ref="formRef"
          :model="form"
          :rules="rules"
          label-position="top"
          style="max-width: 800px"
        >
          <ElFormItem label="模型封面" prop="cover">
            <ElUpload
              class="cover-uploader"
              :auto-upload="false"
              :on-change="handleCoverChange"
              :show-file-list="false"
              accept="image/*"
            >
              <div
                v-if="form.cover"
                class="cover-preview"
                :style="{ backgroundImage: `url(${form.cover})` }"
              >
                <div class="cover-overlay">
                  <span>更换封面</span>
                </div>
              </div>
              <div v-else class="cover-placeholder">
                <svg class="h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                </svg>
                <p class="mt-2 text-sm text-gray-500">点击上传封面图片</p>
                <p class="text-xs text-gray-400">建议尺寸 800x400，支持 JPG、PNG</p>
              </div>
            </ElUpload>
          </ElFormItem>

          <ElFormItem label="模型名称" prop="title">
            <ElInput
              v-model="form.title"
              placeholder="给你的思维模型起个名字，例如：SWOT 分析模型"
              maxlength="50"
              show-word-limit
            />
          </ElFormItem>

          <ElFormItem label="模型描述" prop="description">
            <ElInput
              v-model="form.description"
              type="textarea"
              :rows="4"
              placeholder="简要描述这个思维模型的用途、适用场景和价值..."
              maxlength="500"
              show-word-limit
            />
          </ElFormItem>

          <ElFormItem label="模型分类" prop="category">
            <ElSelect v-model="form.category" placeholder="选择模型分类" style="width: 100%">
              <ElOption
                v-for="cat in categories"
                :key="cat.value"
                :label="cat.label"
                :value="cat.value"
              />
            </ElSelect>
          </ElFormItem>

          <ElFormItem label="标签">
            <div class="flex flex-wrap gap-2">
              <ElTag
                v-for="tag in form.tags"
                :key="tag"
                closable
                @close="handleRemoveTag(tag)"
              >
                {{ tag }}
              </ElTag>
              <ElInput
                v-if="tagInputVisible"
                v-model="tagInput"
                size="small"
                style="width: 100px"
                @keyup.enter="handleAddTag"
                @blur="handleAddTag"
                v-focus
              />
              <ElButton
                v-else
                size="small"
                @click="tagInputVisible = true"
              >
                + 添加标签
              </ElButton>
            </div>
            <p class="text-xs text-gray-400 mt-1">标签帮助用户更好地发现你的模型，建议 3-5 个</p>
          </ElFormItem>
        </ElForm>
      </div>

      <!-- 步骤 2: 详细内容 -->
      <div v-else-if="currentStep === 1">
        <div style="max-width: 800px">
          <div class="mb-6">
            <h3 class="text-lg font-medium text-gray-900 mb-2">模型概述</h3>
            <p class="text-sm text-gray-500 mb-4">介绍这个思维模型的背景、原理和核心价值</p>
            <ElInput
              v-model="form.content.overview"
              type="textarea"
              :rows="6"
              placeholder="详细介绍这个思维模型..."
            />
          </div>

          <div class="mb-6">
            <div class="flex items-center justify-between mb-4">
              <div>
                <h3 class="text-lg font-medium text-gray-900">使用步骤</h3>
                <p class="text-sm text-gray-500">按顺序列出使用这个模型的步骤</p>
              </div>
              <ElButton type="primary" plain size="small" @click="addStep">
                + 添加步骤
              </ElButton>
            </div>
            <div class="space-y-4">
              <div
                v-for="(step, index) in form.content.steps"
                :key="index"
                class="rounded-lg border border-gray-200 p-4"
              >
                <div class="flex items-start gap-4">
                  <div class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-purple-100 text-sm font-medium text-purple-600">
                    {{ index + 1 }}
                  </div>
                  <div class="flex-1 space-y-3">
                    <ElInput
                      v-model="step.title"
                      placeholder="步骤标题"
                    />
                    <ElInput
                      v-model="step.description"
                      type="textarea"
                      :rows="2"
                      placeholder="步骤说明..."
                    />
                  </div>
                  <ElButton
                    v-if="form.content.steps.length > 1"
                    type="danger"
                    plain
                    size="small"
                    @click="removeStep(index)"
                  >
                    删除
                  </ElButton>
                </div>
              </div>
            </div>
          </div>

          <div class="mb-6">
            <div class="flex items-center justify-between mb-4">
              <div>
                <h3 class="text-lg font-medium text-gray-900">实践案例</h3>
                <p class="text-sm text-gray-500">提供真实或假设的应用案例</p>
              </div>
              <ElButton type="primary" plain size="small" @click="addExample">
                + 添加案例
              </ElButton>
            </div>
            <div class="space-y-4">
              <div
                v-for="(example, index) in form.content.examples"
                :key="index"
                class="rounded-lg border border-gray-200 p-4"
              >
                <div class="flex items-start gap-4">
                  <div class="flex-1 space-y-3">
                    <ElInput
                      v-model="example.title"
                      placeholder="案例标题"
                    />
                    <ElInput
                      v-model="example.content"
                      type="textarea"
                      :rows="3"
                      placeholder="案例详细内容..."
                    />
                  </div>
                  <ElButton
                    v-if="form.content.examples.length > 1"
                    type="danger"
                    plain
                    size="small"
                    @click="removeExample(index)"
                  >
                    删除
                  </ElButton>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 步骤 3: 发布设置 -->
      <div v-else-if="currentStep === 2">
        <div style="max-width: 600px">
          <ElForm label-position="top">
            <ElFormItem label="是否免费">
              <ElSwitch
                v-model="form.isFree"
                active-text="免费"
                inactive-text="付费"
              />
              <p class="text-xs text-gray-400 mt-2">
                {{ form.isFree ? '用户可免费使用你的模型，有助于获得更多曝光和反馈' : '设置合理的价格，让你的知识产生价值' }}
              </p>
            </ElFormItem>

            <ElFormItem v-if="!form.isFree" label="价格 (元)">
              <ElInput
                v-model.number="form.price"
                type="number"
                min="1"
                max="9999"
                placeholder="输入价格"
              >
                <template #prefix>¥</template>
              </ElInput>
              <p class="text-xs text-gray-400 mt-2">
                建议定价：简单模型 9-29 元，复杂模型 39-99 元。平台将抽取 20% 作为服务费。
              </p>
            </ElFormItem>
          </ElForm>

          <div class="mt-8 rounded-lg bg-gray-50 p-4">
            <h4 class="font-medium text-gray-900 mb-2">发布须知</h4>
            <ul class="text-sm text-gray-500 space-y-1">
              <li>• 模型提交后将进入审核流程，通常在 24 小时内完成</li>
              <li>• 确保内容原创或已获得授权，禁止抄袭</li>
              <li>• 模型一经发布，可被其他用户采纳、练习和评价</li>
              <li>• 付费模型的收入将在用户确认收货后结算</li>
            </ul>
          </div>
        </div>
      </div>

      <!-- 底部按钮 -->
      <div class="mt-8 flex justify-between">
        <div>
          <ElButton v-if="currentStep === 0" @click="handleCancel">
            取消
          </ElButton>
          <ElButton v-else @click="prevStep">
            上一步
          </ElButton>
        </div>
        <div class="flex gap-3">
          <ElButton v-if="currentStep < 2" type="primary" @click="nextStep">
            下一步
          </ElButton>
          <template v-else>
            <ElButton @click="handleSaveDraft">
              保存草稿
            </ElButton>
            <ElButton type="primary" @click="handleSubmit">
              提交审核
            </ElButton>
          </template>
        </div>
      </div>
    </ElCard>
  </Page>
</template>

<style scoped>
.cover-uploader :deep(.el-upload) {
  display: block;
}

.cover-preview {
  width: 320px;
  height: 180px;
  background-size: cover;
  background-position: center;
  border-radius: 8px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
}

.cover-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s;
}

.cover-preview:hover .cover-overlay {
  opacity: 1;
}

.cover-overlay span {
  color: white;
  font-size: 14px;
}

.cover-placeholder {
  width: 320px;
  height: 180px;
  border: 2px dashed #d9d9d9;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: border-color 0.3s;
}

.cover-placeholder:hover {
  border-color: var(--el-color-primary);
}
</style>
