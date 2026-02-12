<script lang="ts" setup>
import { onMounted, ref, watch } from 'vue';

import { Page } from '@vben/common-ui';

import {
  ElButton,
  ElCard,
  ElTable,
  ElTableColumn,
  ElTag,
  ElInput,
  ElSelect,
  ElOption,
  ElMessage,
  ElMessageBox,
  ElDialog,
  ElForm,
  ElFormItem,
  ElSwitch,
  ElPagination,
} from 'element-plus';

import {
  getModelListApi,
  createModelApi,
  updateModelApi,
  deleteModelApi,
  getCategoryListApi,
  type ModelApi,
  type CategoryApi,
} from '#/api';

// 加载状态
const loading = ref(false);
const dialogLoading = ref(false);

// 模型列表
const models = ref<ModelApi.ThinkingModel[]>([]);
const total = ref(0);

// 分页
const currentPage = ref(1);
const pageSize = ref(10);

// 筛选
const searchKeyword = ref('');
const categoryFilter = ref('all');
const statusFilter = ref<'all' | 'free' | 'paid'>('all');

// 分类列表
const categories = ref<CategoryApi.Category[]>([]);

// 对话框
const dialogVisible = ref(false);
const dialogTitle = ref('创建模型');
const editingId = ref<string | null>(null);

// 表单数据
const formData = ref<ModelApi.CreateModelParams>({
  title: '',
  description: '',
  cover: '',
  categoryId: '',
  tags: [],
  isFree: true,
  price: 0,
  content: '',
});

// 表单验证规则
const rules = {
  title: [{ required: true, message: '请输入模型标题', trigger: 'blur' }],
  description: [{ required: true, message: '请输入模型描述', trigger: 'blur' }],
  categoryId: [{ required: true, message: '请选择分类', trigger: 'change' }],
};

const formRef = ref<InstanceType<typeof ElForm>>();

// 获取模型列表
async function fetchModels() {
  loading.value = true;
  try {
    const params: ModelApi.ModelListParams = {
      page: currentPage.value,
      pageSize: pageSize.value,
      keyword: searchKeyword.value || undefined,
    };

    if (categoryFilter.value !== 'all') {
      params.category = categoryFilter.value;
    }

    if (statusFilter.value === 'free') {
      params.isFree = true;
    } else if (statusFilter.value === 'paid') {
      params.isFree = false;
    }

    const res = await getModelListApi(params);
    models.value = res.list;
    total.value = res.total;
  } catch (error) {
    console.error('获取模型列表失败:', error);
    ElMessage.error('获取模型列表失败');
  } finally {
    loading.value = false;
  }
}

// 获取分类列表
async function fetchCategories() {
  try {
    const res = await getCategoryListApi();
    categories.value = res;
  } catch (error) {
    console.error('获取分类列表失败:', error);
  }
}

// 监听筛选条件变化
watch([searchKeyword, categoryFilter, statusFilter], () => {
  currentPage.value = 1;
  fetchModels();
});

// 监听分页变化
watch([currentPage, pageSize], () => {
  fetchModels();
});

// 页面加载时获取数据
onMounted(() => {
  fetchModels();
  fetchCategories();
});

// 打开创建对话框
function openCreateDialog() {
  dialogTitle.value = '创建模型';
  editingId.value = null;
  formData.value = {
    title: '',
    description: '',
    cover: '',
    categoryId: '',
    tags: [],
    isFree: true,
    price: 0,
    content: '',
  };
  dialogVisible.value = true;
}

// 打开编辑对话框
function openEditDialog(model: ModelApi.ThinkingModel) {
  dialogTitle.value = '编辑模型';
  editingId.value = model.id;
  formData.value = {
    title: model.title,
    description: model.description,
    cover: model.cover,
    categoryId: model.categoryId || '',
    tags: [...model.tags],
    isFree: model.isFree,
    price: model.price || 0,
    content: model.content || '',
  };
  dialogVisible.value = true;
}

// 提交表单
async function handleSubmit() {
  if (!formRef.value) return;

  await formRef.value.validate(async (valid) => {
    if (!valid) return;

    dialogLoading.value = true;
    try {
      if (editingId.value) {
        await updateModelApi({ ...formData.value, id: editingId.value });
        ElMessage.success('模型更新成功');
      } else {
        await createModelApi(formData.value);
        ElMessage.success('模型创建成功');
      }
      dialogVisible.value = false;
      fetchModels();
    } catch (error) {
      console.error('操作失败:', error);
      ElMessage.error('操作失败');
    } finally {
      dialogLoading.value = false;
    }
  });
}

// 删除模型
async function handleDelete(model: ModelApi.ThinkingModel) {
  try {
    await ElMessageBox.confirm(
      `确定要删除模型"${model.title}"吗？此操作不可恢复。`,
      '确认删除',
      { type: 'warning' }
    );
    await deleteModelApi(model.id);
    ElMessage.success('模型已删除');
    fetchModels();
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error);
      ElMessage.error('删除失败');
    }
  }
}

// 格式化日期
function formatDate(dateStr: string): string {
  if (!dateStr) return '-';
  const date = new Date(dateStr);
  return date.toLocaleDateString('zh-CN');
}


// 标签输入处理
const tagInput = ref('');
function handleAddTag() {
  if (tagInput.value && !formData.value.tags?.includes(tagInput.value)) {
    formData.value.tags = [...(formData.value.tags || []), tagInput.value];
    tagInput.value = '';
  }
}
function handleRemoveTag(tag: string) {
  formData.value.tags = formData.value.tags?.filter((t) => t !== tag) || [];
}
</script>

<template>
  <Page
    description="管理系统中的思维模型，包括创建、编辑、删除等操作"
    title="模型管理"
  >
    <ElCard shadow="never">
      <template #header>
        <div class="flex flex-wrap items-center justify-between gap-4">
          <div class="flex items-center gap-4">
            <h2 class="text-lg font-semibold">模型列表</h2>
            <ElSelect v-model="categoryFilter" style="width: 140px" placeholder="全部分类">
              <ElOption label="全部分类" value="all" />
              <ElOption
                v-for="cat in categories"
                :key="cat.id"
                :label="cat.name"
                :value="cat.id"
              />
            </ElSelect>
            <ElSelect v-model="statusFilter" style="width: 120px">
              <ElOption label="全部状态" value="all" />
              <ElOption label="免费" value="free" />
              <ElOption label="付费" value="paid" />
            </ElSelect>
          </div>
          <div class="flex items-center gap-3">
            <ElInput
              v-model="searchKeyword"
              placeholder="搜索模型..."
              clearable
              style="width: 220px"
              @keyup.enter="fetchModels"
            />
            <ElButton type="primary" @click="openCreateDialog">
              <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
              </svg>
              创建模型
            </ElButton>
          </div>
        </div>
      </template>

      <!-- 表格 -->
      <ElTable v-loading="loading" :data="models" stripe>
        <ElTableColumn label="模型信息" min-width="250">
          <template #default="{ row }">
            <div class="flex items-center gap-3">
              <img :src="row.cover" class="w-16 h-12 object-cover rounded" />
              <div>
                <div class="font-medium">{{ row.title }}</div>
                <div class="text-xs text-gray-500 line-clamp-1">{{ row.description }}</div>
              </div>
            </div>
          </template>
        </ElTableColumn>

        <ElTableColumn label="分类" width="120">
          <template #default="{ row }">
            <ElTag size="small">{{ row.category }}</ElTag>
          </template>
        </ElTableColumn>

        <ElTableColumn label="价格" width="100">
          <template #default="{ row }">
            <ElTag :type="row.isFree ? 'success' : 'warning'" size="small">
              {{ row.isFree ? '免费' : `¥${row.price}` }}
            </ElTag>
          </template>
        </ElTableColumn>

        <ElTableColumn label="统计" width="200">
          <template #default="{ row }">
            <div class="text-xs text-gray-500 space-y-1">
              <div>采纳: {{ row.stats.adoptions }}</div>
              <div>练习: {{ row.stats.practices }}</div>
              <div>点赞: {{ row.stats.likes }}</div>
            </div>
          </template>
        </ElTableColumn>

        <ElTableColumn label="更新日期" width="120">
          <template #default="{ row }">
            {{ formatDate(row.updatedAt) }}
          </template>
        </ElTableColumn>

        <ElTableColumn label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <div class="flex items-center gap-2">
              <ElButton type="primary" size="small" @click="openEditDialog(row)">
                编辑
              </ElButton>
              <ElButton type="danger" size="small" plain @click="handleDelete(row)">
                删除
              </ElButton>
            </div>
          </template>
        </ElTableColumn>
      </ElTable>

      <!-- 分页 -->
      <div class="flex justify-center mt-6">
        <ElPagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
        />
      </div>
    </ElCard>

    <!-- 创建/编辑对话框 -->
    <ElDialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="700px"
      :close-on-click-modal="false"
    >
      <ElForm
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-position="top"
        v-loading="dialogLoading"
      >
        <div class="grid grid-cols-2 gap-4">
          <ElFormItem label="模型标题" prop="title">
            <ElInput v-model="formData.title" placeholder="请输入模型标题" />
          </ElFormItem>

          <ElFormItem label="分类" prop="categoryId">
            <ElSelect v-model="formData.categoryId" placeholder="请选择分类" class="w-full">
              <ElOption
                v-for="cat in categories"
                :key="cat.id"
                :label="cat.name"
                :value="cat.id"
              />
            </ElSelect>
          </ElFormItem>
        </div>

        <ElFormItem label="描述" prop="description">
          <ElInput
            v-model="formData.description"
            type="textarea"
            :rows="2"
            placeholder="请输入模型描述"
          />
        </ElFormItem>

        <ElFormItem label="封面图片URL">
          <ElInput v-model="formData.cover" placeholder="https://..." />
        </ElFormItem>

        <div class="grid grid-cols-2 gap-4">
          <ElFormItem label="是否免费">
            <ElSwitch
              v-model="formData.isFree"
              active-text="免费"
              inactive-text="付费"
            />
          </ElFormItem>

          <ElFormItem v-if="!formData.isFree" label="价格（元）">
            <ElInput v-model.number="formData.price" type="number" placeholder="0.00" />
          </ElFormItem>
        </div>

        <ElFormItem label="标签">
          <div class="flex flex-wrap gap-2 mb-2">
            <ElTag
              v-for="tag in formData.tags"
              :key="tag"
              closable
              @close="handleRemoveTag(tag)"
            >
              {{ tag }}
            </ElTag>
          </div>
          <div class="flex gap-2">
            <ElInput
              v-model="tagInput"
              placeholder="输入标签后按回车"
              @keyup.enter="handleAddTag"
            />
            <ElButton @click="handleAddTag">添加</ElButton>
          </div>
        </ElFormItem>

        <ElFormItem label="详细内容（Markdown）">
          <ElInput
            v-model="formData.content"
            type="textarea"
            :rows="6"
            placeholder="## 使用指南\n\n输入详细的使用说明、步骤、示例等内容..."
          />
        </ElFormItem>
      </ElForm>

      <template #footer>
        <ElButton @click="dialogVisible = false">取消</ElButton>
        <ElButton type="primary" :loading="dialogLoading" @click="handleSubmit">
          保存
        </ElButton>
      </template>
    </ElDialog>
  </Page>
</template>

<style scoped>
.line-clamp-1 {
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
