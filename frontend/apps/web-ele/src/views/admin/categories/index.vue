<script lang="ts" setup>
import { onMounted, ref } from 'vue';

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
  ElInputNumber,
} from 'element-plus';

import {
  getCategoryListApi,
  createCategoryApi,
  updateCategoryApi,
  deleteCategoryApi,
  type CategoryApi,
} from '#/api';

// 加载状态
const loading = ref(false);
const dialogLoading = ref(false);

// 分类列表
const categories = ref<CategoryApi.Category[]>([]);

// 对话框
const dialogVisible = ref(false);
const dialogTitle = ref('创建分类');
const editingId = ref<string | null>(null);

// 表单数据
const formData = ref<CategoryApi.CreateCategoryParams>({
  name: '',
  description: '',
  icon: '',
  sort: 0,
});

// 表单验证规则
const rules = {
  name: [{ required: true, message: '请输入分类名称', trigger: 'blur' }],
};

const formRef = ref<InstanceType<typeof ElForm>>();

// 获取分类列表
async function fetchCategories() {
  loading.value = true;
  try {
    const res = await getCategoryListApi();
    categories.value = res.sort((a, b) => a.sort - b.sort);
  } catch (error) {
    console.error('获取分类列表失败:', error);
    ElMessage.error('获取分类列表失败');
  } finally {
    loading.value = false;
  }
}

// 页面加载时获取数据
onMounted(() => {
  fetchCategories();
});

// 打开创建对话框
function openCreateDialog() {
  dialogTitle.value = '创建分类';
  editingId.value = null;
  formData.value = {
    name: '',
    description: '',
    icon: '',
    sort: categories.value.length,
  };
  dialogVisible.value = true;
}

// 打开编辑对话框
function openEditDialog(category: CategoryApi.Category) {
  dialogTitle.value = '编辑分类';
  editingId.value = category.id;
  formData.value = {
    name: category.name,
    description: category.description || '',
    icon: category.icon || '',
    sort: category.sort,
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
        await updateCategoryApi({ ...formData.value, id: editingId.value });
        ElMessage.success('分类更新成功');
      } else {
        await createCategoryApi(formData.value);
        ElMessage.success('分类创建成功');
      }
      dialogVisible.value = false;
      fetchCategories();
    } catch (error) {
      console.error('操作失败:', error);
      ElMessage.error('操作失败');
    } finally {
      dialogLoading.value = false;
    }
  });
}

// 删除分类
async function handleDelete(category: CategoryApi.Category) {
  try {
    await ElMessageBox.confirm(
      `确定要删除分类"${category.name}"吗？该分类下的模型需要重新分配分类。`,
      '确认删除',
      { type: 'warning' }
    );
    await deleteCategoryApi(category.id);
    ElMessage.success('分类已删除');
    fetchCategories();
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error);
      ElMessage.error('删除失败');
    }
  }
}

// 格式化日期
function formatDate(dateStr?: string): string {
  if (!dateStr) return '-';
  const date = new Date(dateStr);
  return date.toLocaleDateString('zh-CN');
}
</script>

<template>
  <Page
    description="管理思维模型的分类，用于组织和筛选模型"
    title="分类管理"
  >
    <ElCard shadow="never">
      <template #header>
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">分类列表</h2>
          <ElButton type="primary" @click="openCreateDialog">
            <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
            </svg>
            创建分类
          </ElButton>
        </div>
      </template>

      <!-- 表格 -->
      <ElTable v-loading="loading" :data="categories" stripe>
        <ElTableColumn label="排序" width="80">
          <template #default="{ row }">
            <span class="text-gray-500">{{ row.sort }}</span>
          </template>
        </ElTableColumn>

        <ElTableColumn label="分类名称" min-width="150">
          <template #default="{ row }">
            <div class="font-medium">{{ row.name }}</div>
          </template>
        </ElTableColumn>

        <ElTableColumn label="图标" width="100">
          <template #default="{ row }">
            <span v-if="row.icon" class="text-xl">{{ row.icon }}</span>
            <span v-else class="text-gray-400">-</span>
          </template>
        </ElTableColumn>

        <ElTableColumn label="描述" min-width="200">
          <template #default="{ row }">
            <span class="text-gray-500">{{ row.description || '-' }}</span>
          </template>
        </ElTableColumn>

        <ElTableColumn label="状态" width="100">
          <template #default="{ row }">
            <ElTag :type="row.status === 1 ? 'success' : 'info'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </ElTag>
          </template>
        </ElTableColumn>

        <ElTableColumn label="模型数" width="100">
          <template #default="{ row }">
            <span class="text-gray-500">{{ row.modelCount || 0 }}</span>
          </template>
        </ElTableColumn>

        <ElTableColumn label="创建日期" width="120">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
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
    </ElCard>

    <!-- 创建/编辑对话框 -->
    <ElDialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="500px"
      :close-on-click-modal="false"
    >
      <ElForm
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="80px"
        v-loading="dialogLoading"
      >
        <ElFormItem label="名称" prop="name">
          <ElInput v-model="formData.name" placeholder="请输入分类名称" />
        </ElFormItem>

        <ElFormItem label="描述">
          <ElInput
            v-model="formData.description"
            type="textarea"
            :rows="2"
            placeholder="请输入分类描述（可选）"
          />
        </ElFormItem>

        <ElFormItem label="图标">
          <ElInput v-model="formData.icon" placeholder="emoji 或图标类名（可选）" />
        </ElFormItem>

        <ElFormItem label="排序">
          <ElInputNumber v-model="formData.sort" :min="0" :max="100" />
          <span class="text-xs text-gray-400 ml-2">数字越小排序越靠前</span>
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
