<script lang="ts" setup>
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';

import { Page } from '@vben/common-ui';

import {
  ElButton,
  ElCard,
  ElInput,
  ElSelect,
  ElOption,
  ElMessage,
  ElTag,
  ElUpload,
  type UploadProps,
} from 'element-plus';

import { createDiscussionApi, getModelListApi, type CollisionApi, type ModelApi } from '#/api';
import { $t } from '#/locales';

const router = useRouter();

// 表单数据
const formData = ref<CollisionApi.CreateDiscussionParams>({
  title: '',
  content: '',
  summary: '',
  cover: '',
  category: 'inspiration',
  tags: [],
  modelId: undefined,
  status: 1, // 默认发布
});

// 标签输入
const tagInput = ref('');

// 提交状态
const submitting = ref(false);

// 模型列表 (用于关联)
const models = ref<ModelApi.ThinkingModel[]>([]);
const loadingModels = ref(false);

// 分类选项
const categories = [
  { id: 'inspiration', name: $t('page.collision.categories.inspiration'), icon: '💡', description: '分享灵光一现的创意想法' },
  { id: 'methodology', name: $t('page.collision.categories.methodology'), icon: '🧭', description: '探讨思维方法和框架' },
  { id: 'case', name: $t('page.collision.categories.case'), icon: '📋', description: '分享实际案例和分析' },
  { id: 'question', name: $t('page.collision.categories.question'), icon: '❓', description: '提出问题寻求讨论' },
  { id: 'share', name: $t('page.collision.categories.share'), icon: '📚', description: '分享学习心得和经验' },
];

// 获取模型列表
async function fetchModels() {
  loadingModels.value = true;
  try {
    const res = await getModelListApi({ pageSize: 100 });
    models.value = res.list;
  } catch (error) {
    console.error('获取模型列表失败:', error);
  } finally {
    loadingModels.value = false;
  }
}

// 页面加载时获取模型列表
fetchModels();

// 添加标签
function addTag() {
  const tag = tagInput.value.trim();
  if (!tag) return;
  if (formData.value.tags && formData.value.tags.length >= 5) {
    ElMessage.warning('最多添加5个标签');
    return;
  }
  if (formData.value.tags?.includes(tag)) {
    ElMessage.warning('标签已存在');
    return;
  }
  if (!formData.value.tags) {
    formData.value.tags = [];
  }
  formData.value.tags.push(tag);
  tagInput.value = '';
}

// 移除标签
function removeTag(tag: string) {
  if (formData.value.tags) {
    formData.value.tags = formData.value.tags.filter(t => t !== tag);
  }
}

// 处理标签输入按键
function handleTagKeydown(e: Event | KeyboardEvent) {
  if ('key' in e && e.key === 'Enter') {
    e.preventDefault();
    addTag();
  }
}

// 上传封面前的检查
const beforeCoverUpload: UploadProps['beforeUpload'] = (rawFile) => {
  if (!['image/jpeg', 'image/png', 'image/gif', 'image/webp'].includes(rawFile.type)) {
    ElMessage.error('只支持 JPG、PNG、GIF、WebP 格式的图片');
    return false;
  }
  if (rawFile.size / 1024 / 1024 > 5) {
    ElMessage.error('图片大小不能超过 5MB');
    return false;
  }
  return true;
};

// 上传封面成功
const handleCoverSuccess: UploadProps['onSuccess'] = (response) => {
  // 模拟上传成功，实际应该从 response 获取 URL
  formData.value.cover = URL.createObjectURL(response);
};

// 移除封面
function removeCover() {
  formData.value.cover = '';
}

// 表单验证
const isValid = computed(() => {
  return formData.value.title.trim() && formData.value.content.trim();
});

// 提交话题
async function handleSubmit(isDraft = false) {
  if (!formData.value.title.trim()) {
    ElMessage.warning('请输入话题标题');
    return;
  }
  if (!formData.value.content.trim()) {
    ElMessage.warning('请输入话题内容');
    return;
  }

  submitting.value = true;
  try {
    const data = {
      ...formData.value,
      status: isDraft ? 0 : 1,
      summary: formData.value.summary || formData.value.content.slice(0, 100),
    };

    await createDiscussionApi(data);
    ElMessage.success(isDraft ? '草稿保存成功' : '话题发布成功');
    router.push('/collision');
  } catch (error) {
    console.error('提交失败:', error);
    ElMessage.error(isDraft ? '保存草稿失败' : '发布失败，请重试');
  } finally {
    submitting.value = false;
  }
}

// 返回
function goBack() {
  router.push('/collision');
}
</script>

<template>
  <Page :title="$t('page.collision.create')" auto-content-height>
    <!-- 返回按钮 -->
    <div class="mb-4">
      <ElButton link @click="goBack">
        <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"/>
        </svg>
        返回话题列表
      </ElButton>
    </div>

    <div class="max-w-4xl mx-auto">
      <div class="flex gap-6">
        <!-- 主表单区 -->
        <div class="flex-1 min-w-0">
          <ElCard shadow="never">
            <!-- 标题 -->
            <div class="mb-6">
              <label class="block text-sm font-medium text-gray-700 mb-2">
                {{ $t('page.collision.form.title') }} <span class="text-red-500">*</span>
              </label>
              <ElInput
                v-model="formData.title"
                :placeholder="$t('page.collision.form.titlePlaceholder')"
                maxlength="100"
                show-word-limit
                size="large"
              />
            </div>

            <!-- 分类选择 -->
            <div class="mb-6">
              <label class="block text-sm font-medium text-gray-700 mb-2">
                {{ $t('page.collision.form.category') }} <span class="text-red-500">*</span>
              </label>
              <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-3">
                <button
                  v-for="cat in categories"
                  :key="cat.id"
                  :class="[
                    'p-4 rounded-xl border-2 transition-all text-left cursor-pointer',
                    formData.category === cat.id
                      ? 'border-purple-500 bg-purple-50'
                      : 'border-gray-200 hover:border-purple-300 hover:bg-purple-50/50'
                  ]"
                  @click="formData.category = cat.id as CollisionApi.Category"
                >
                  <div class="text-2xl mb-2">{{ cat.icon }}</div>
                  <div class="font-medium text-gray-800 text-sm">{{ cat.name }}</div>
                </button>
              </div>
            </div>

            <!-- 内容 -->
            <div class="mb-6">
              <label class="block text-sm font-medium text-gray-700 mb-2">
                {{ $t('page.collision.form.content') }} <span class="text-red-500">*</span>
              </label>
              <ElInput
                v-model="formData.content"
                type="textarea"
                :placeholder="$t('page.collision.form.contentPlaceholder')"
                :rows="12"
                maxlength="10000"
                show-word-limit
                resize="vertical"
              />
              <div class="mt-2 text-xs text-gray-400">
                支持 Markdown 语法，可以使用 **粗体**、*斜体*、`代码` 等格式
              </div>
            </div>

            <!-- 封面图 -->
            <div class="mb-6">
              <label class="block text-sm font-medium text-gray-700 mb-2">
                封面图片 <span class="text-gray-400 font-normal">(可选)</span>
              </label>
              <div v-if="formData.cover" class="relative inline-block">
                <img :src="formData.cover" class="w-64 h-40 object-cover rounded-lg" />
                <button
                  class="absolute -top-2 -right-2 w-6 h-6 bg-red-500 text-white rounded-full flex items-center justify-center hover:bg-red-600 cursor-pointer"
                  @click="removeCover"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
                  </svg>
                </button>
              </div>
              <ElUpload
                v-else
                class="cover-uploader"
                action="/api/upload"
                :show-file-list="false"
                :before-upload="beforeCoverUpload"
                :on-success="handleCoverSuccess"
              >
                <div class="w-64 h-40 border-2 border-dashed border-gray-300 rounded-lg flex flex-col items-center justify-center cursor-pointer hover:border-purple-400 hover:bg-purple-50/50 transition-colors">
                  <svg class="w-10 h-10 text-gray-400 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                  </svg>
                  <span class="text-sm text-gray-500">点击上传封面图</span>
                  <span class="text-xs text-gray-400 mt-1">支持 JPG、PNG、GIF、WebP，不超过 5MB</span>
                </div>
              </ElUpload>
            </div>

            <!-- 标签 -->
            <div class="mb-6">
              <label class="block text-sm font-medium text-gray-700 mb-2">
                {{ $t('page.collision.form.tags') }} <span class="text-gray-400 font-normal">(最多5个)</span>
              </label>
              <div class="flex flex-wrap gap-2 mb-3">
                <ElTag
                  v-for="tag in formData.tags"
                  :key="tag"
                  closable
                  effect="plain"
                  @close="removeTag(tag)"
                >
                  {{ tag }}
                </ElTag>
              </div>
              <div class="flex gap-2">
                <ElInput
                  v-model="tagInput"
                  :placeholder="$t('page.collision.form.tagsPlaceholder')"
                  class="flex-1"
                  maxlength="20"
                  @keydown="handleTagKeydown"
                />
                <ElButton @click="addTag">添加</ElButton>
              </div>
            </div>

            <!-- 关联模型 -->
            <div class="mb-6">
              <label class="block text-sm font-medium text-gray-700 mb-2">
                {{ $t('page.collision.form.relatedModel') }}
              </label>
              <ElSelect
                v-model="formData.modelId"
                :placeholder="$t('page.collision.form.relatedModelPlaceholder')"
                :loading="loadingModels"
                clearable
                filterable
                class="w-full"
              >
                <ElOption
                  v-for="model in models"
                  :key="model.id"
                  :label="model.title"
                  :value="model.id"
                >
                  <div class="flex items-center gap-2">
                    <span>📊</span>
                    <span>{{ model.title }}</span>
                  </div>
                </ElOption>
              </ElSelect>
              <div class="mt-2 text-xs text-gray-400">
                关联思维模型可以让读者更好地理解你的话题背景
              </div>
            </div>

            <!-- 摘要 -->
            <div class="mb-8">
              <label class="block text-sm font-medium text-gray-700 mb-2">
                内容摘要 <span class="text-gray-400 font-normal">(可选，不填则自动截取)</span>
              </label>
              <ElInput
                v-model="formData.summary"
                type="textarea"
                placeholder="用一两句话概括你的话题内容"
                :rows="2"
                maxlength="200"
                show-word-limit
              />
            </div>

            <!-- 操作按钮 -->
            <div class="flex items-center justify-between pt-6 border-t border-gray-100">
              <ElButton @click="goBack">取消</ElButton>
              <div class="flex gap-3">
                <ElButton
                  :loading="submitting"
                  @click="handleSubmit(true)"
                >
                  {{ $t('page.collision.form.saveDraft') }}
                </ElButton>
                <ElButton
                  type="primary"
                  :loading="submitting"
                  :disabled="!isValid"
                  @click="handleSubmit(false)"
                >
                  {{ $t('page.collision.form.submit') }}
                </ElButton>
              </div>
            </div>
          </ElCard>
        </div>

        <!-- 右侧提示栏 -->
        <div class="w-72 flex-shrink-0 space-y-6 hidden lg:block">
          <!-- 写作提示 -->
          <ElCard shadow="never" class="bg-gradient-to-br from-purple-50 to-indigo-50">
            <div class="space-y-4">
              <div class="flex items-center gap-2 text-purple-600">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"/>
                </svg>
                <span class="font-medium">写作小贴士</span>
              </div>
              <ul class="space-y-3 text-sm text-gray-600">
                <li class="flex items-start gap-2">
                  <span class="text-purple-500 mt-0.5">•</span>
                  <span>好的标题能吸引更多人参与讨论</span>
                </li>
                <li class="flex items-start gap-2">
                  <span class="text-purple-500 mt-0.5">•</span>
                  <span>分享具体案例和个人经历更有说服力</span>
                </li>
                <li class="flex items-start gap-2">
                  <span class="text-purple-500 mt-0.5">•</span>
                  <span>提出问题让读者参与思考</span>
                </li>
                <li class="flex items-start gap-2">
                  <span class="text-purple-500 mt-0.5">•</span>
                  <span>关联思维模型让内容更有深度</span>
                </li>
              </ul>
            </div>
          </ElCard>

          <!-- 热门标签推荐 -->
          <ElCard shadow="never">
            <template #header>
              <span class="font-medium">热门标签</span>
            </template>
            <div class="flex flex-wrap gap-2">
              <button
                v-for="tag in ['思维模型', '创新思维', '决策方法', '第一性原理', '认知偏差', 'AI思考', '案例分析', '经验分享']"
                :key="tag"
                :class="[
                  'px-3 py-1.5 rounded-full text-sm transition-colors cursor-pointer',
                  formData.tags?.includes(tag)
                    ? 'bg-purple-500 text-white'
                    : 'bg-gray-100 text-gray-600 hover:bg-purple-100 hover:text-purple-600'
                ]"
                @click="formData.tags?.includes(tag) ? removeTag(tag) : (formData.tags = [...(formData.tags || []), tag])"
              >
                {{ tag }}
              </button>
            </div>
          </ElCard>

          <!-- 发布须知 -->
          <ElCard shadow="never">
            <template #header>
              <span class="font-medium">发布须知</span>
            </template>
            <ul class="space-y-2 text-sm text-gray-500">
              <li class="flex items-start gap-2">
                <svg class="w-4 h-4 text-green-500 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
                </svg>
                <span>内容需与思维模型相关</span>
              </li>
              <li class="flex items-start gap-2">
                <svg class="w-4 h-4 text-green-500 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
                </svg>
                <span>禁止抄袭，鼓励原创</span>
              </li>
              <li class="flex items-start gap-2">
                <svg class="w-4 h-4 text-green-500 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
                </svg>
                <span>尊重他人，理性讨论</span>
              </li>
              <li class="flex items-start gap-2">
                <svg class="w-4 h-4 text-red-500 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
                </svg>
                <span>禁止发布广告和违规内容</span>
              </li>
            </ul>
          </ElCard>
        </div>
      </div>
    </div>
  </Page>
</template>

<style scoped>
:deep(.cover-uploader .el-upload) {
  border: none;
}
</style>
