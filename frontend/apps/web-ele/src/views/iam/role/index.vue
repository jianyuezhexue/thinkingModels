<template>
  <Page
    description="管理系统角色和权限配置"
    title="角色管理"
    content-class="p-6 bg-gray-50"
  >
    <div class="space-y-6">
      <!-- 搜索筛选区 -->
      <ElCard shadow="hover" class="!rounded-xl">
        <template #header>
          <div class="flex items-center gap-2">
            <span class="text-lg">🔍</span>
            <span class="font-semibold text-gray-700">搜索筛选</span>
          </div>
        </template>

        <ElForm :model="searchForm" inline class="flex flex-wrap gap-4">
          <ElFormItem label="角色名称">
            <ElInput
              v-model="searchForm.roleName"
              placeholder="请输入角色名称"
              clearable
              class="!w-48"
            />
          </ElFormItem>
          <ElFormItem>
            <ElButton type="primary" @click="handleSearch">
              <template #icon>
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
                </svg>
              </template>
              搜索
            </ElButton>
            <ElButton @click="handleReset">重置</ElButton>
          </ElFormItem>
        </ElForm>
      </ElCard>

      <!-- 角色列表 -->
      <ElCard shadow="hover" class="!rounded-xl">
        <template #header>
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <span class="text-lg">🛡️</span>
              <span class="font-semibold text-gray-700">角色列表</span>
              <span class="text-sm text-gray-400">(共 {{ pagination.total }} 条)</span>
            </div>
            <ElButton type="primary" @click="handleAdd">
              <template #icon>
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
                </svg>
              </template>
              新增角色
            </ElButton>
          </div>
        </template>

        <ElTable
          v-loading="loading"
          :data="tableData"
          stripe
          class="!rounded-lg"
        >
          <ElTableColumn label="角色名称" min-width="160">
            <template #default="slotProps">
              <div v-if="slotProps?.row" class="flex items-center gap-2">
                <div
                  class="w-8 h-8 rounded-full flex items-center justify-center text-white text-sm font-bold"
                  :style="{ backgroundColor: getRoleColor(slotProps.row.roleCode) }"
                >
                  {{ slotProps.row.roleName?.charAt(0) || 'R' }}
                </div>
                <span class="font-medium">{{ slotProps.row.roleName }}</span>
              </div>
            </template>
          </ElTableColumn>
          <ElTableColumn label="角色编码" prop="roleCode" min-width="120">
            <template #default="slotProps">
              <code v-if="slotProps?.row" class="px-2 py-1 bg-gray-100 rounded text-sm">{{ slotProps.row.roleCode }}</code>
            </template>
          </ElTableColumn>
          <ElTableColumn label="描述" prop="description" min-width="220" show-overflow-tooltip />
          <ElTableColumn label="用户数" width="100" align="center">
            <template #default="slotProps">
              <ElTag v-if="slotProps?.row" type="info" effect="plain" class="!rounded-full">
                {{ slotProps.row.userCount || 0 }}
              </ElTag>
            </template>
          </ElTableColumn>
          <ElTableColumn label="创建时间" prop="createdAt" min-width="160" />
          <ElTableColumn label="操作" width="280" fixed="right">
            <template #default="slotProps">
              <div v-if="slotProps?.row" class="flex gap-2">
                <ElButton
                  type="primary"
                  size="small"
                  plain
                  @click="handleEdit(slotProps.row)"
                >
                  编辑
                </ElButton>
                <ElButton
                  type="warning"
                  size="small"
                  plain
                  @click="handlePermission(slotProps.row)"
                >
                  权限
                </ElButton>
                <ElButton
                  type="danger"
                  size="small"
                  plain
                  @click="handleDelete(slotProps.row)"
                >
                  删除
                </ElButton>
              </div>
            </template>
          </ElTableColumn>
        </ElTable>

        <!-- 分页 -->
        <div class="flex justify-end mt-4">
          <ElPagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :total="pagination.total"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </ElCard>
    </div>

    <!-- 新增/编辑对话框 -->
    <ElDialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑角色' : '新增角色'"
      width="500px"
      destroy-on-close
    >
      <ElForm
        ref="formRef"
        :model="form"
        :rules="rules"
        label-position="top"
      >
        <ElFormItem label="角色名称" prop="roleName">
          <ElInput v-model="form.roleName" placeholder="请输入角色名称" />
        </ElFormItem>
        <ElFormItem label="角色编码" prop="roleCode">
          <ElInput
            v-model="form.roleCode"
            placeholder="请输入角色编码，如：admin"
            :disabled="isEdit"
          />
        </ElFormItem>
        <ElFormItem label="描述" prop="description">
          <ElInput
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入角色描述"
          />
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="dialogVisible = false">取消</ElButton>
        <ElButton type="primary" :loading="submitLoading" @click="handleSubmit">
          确定
        </ElButton>
      </template>
    </ElDialog>

    <!-- 权限配置对话框 -->
    <ElDialog
      v-model="permissionDialogVisible"
      title="权限配置"
      width="600px"
      destroy-on-close
    >
      <ElForm label-position="top">
        <ElFormItem label="当前角色">
          <ElTag size="large" effect="dark" :color="currentRole ? getRoleColor(currentRole.roleCode) : ''">
            {{ currentRole?.roleName }}
          </ElTag>
        </ElFormItem>
        <ElFormItem label="菜单权限">
          <ElTree
            ref="treeRef"
            :data="menuTree"
            show-checkbox
            node-key="id"
            :default-expanded-keys="['1']"
            :checked-keys="checkedMenuIds"
            :props="{ label: 'label', children: 'children' }"
            class="border rounded-lg p-4"
          />
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="permissionDialogVisible = false">取消</ElButton>
        <ElButton
          type="primary"
          :loading="permissionLoading"
          @click="handleSavePermission"
        >
          保存权限
        </ElButton>
      </template>
    </ElDialog>
  </Page>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import {
  ElMessage,
  ElMessageBox,
  ElCard,
  ElForm,
  ElFormItem,
  ElInput,
  ElButton,
  ElTable,
  ElTableColumn,
  ElTag,
  ElPagination,
  ElDialog,
  ElTree,
} from 'element-plus';
import type { FormInstance, FormRules } from 'element-plus';

import { Page } from '@vben/common-ui';

import {
  getRoleListApi,
  createRoleApi,
  updateRoleApi,
  deleteRoleApi,
  updateRolePermissionApi,
} from '#/api/iam/role';
import type { RoleApi } from '#/api/iam/role';

// 角色颜色映射
const roleColorMap: Record<string, string> = {
  admin: '#7C3AED',
  operator: '#10B981',
  user: '#3B82F6',
  vip: '#F59E0B',
  svip: '#EC4899',
  guest: '#6B7280',
  tester: '#8B5CF6',
  editor: '#F97316',
};

// 获取角色颜色
function getRoleColor(roleCode: string) {
  return roleColorMap[roleCode] || '#6B7280';
}

// 搜索表单
const searchForm = reactive({
  roleName: '',
});

// 列表数据
const loading = ref(false);
const tableData = ref<RoleApi.Role[]>([]);
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
});

// 对话框
const dialogVisible = ref(false);
const isEdit = ref(false);
const formRef = ref<FormInstance>();
const submitLoading = ref(false);

// 表单
const form = reactive({
  id: '',
  roleName: '',
  roleCode: '',
  description: '',
});

// 表单校验规则
const rules: FormRules = {
  roleName: [
    { required: true, message: '请输入角色名称', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' },
  ],
  roleCode: [
    { required: true, message: '请输入角色编码', trigger: 'blur' },
    { pattern: /^[a-z0-9_]+$/, message: '只能包含小写字母、数字和下划线', trigger: 'blur' },
  ],
};

// 权限对话框
const permissionDialogVisible = ref(false);
const permissionLoading = ref(false);
const currentRole = ref<RoleApi.Role | null>(null);
const treeRef = ref<any>(null);
const checkedMenuIds = ref<string[]>([]);

// 菜单树数据
const menuTree = ref<RoleApi.MenuNode[]>([
  {
    id: '1',
    label: '系统管理',
    children: [
      { id: '1-1', label: '用户管理' },
      { id: '1-2', label: '角色管理' },
    ],
  },
  {
    id: '2',
    label: '内容管理',
    children: [
      { id: '2-1', label: '模型管理' },
      { id: '2-2', label: '分类管理' },
    ],
  },
  {
    id: '3',
    label: '数据分析',
    children: [
      { id: '3-1', label: '使用统计' },
      { id: '3-2', label: '用户分析' },
    ],
  },
]);

// 搜索
function handleSearch() {
  pagination.page = 1;
  loadData();
}

// 重置
function handleReset() {
  searchForm.roleName = '';
  handleSearch();
}

// 加载数据
async function loadData() {
  loading.value = true;
  try {
    const res = await getRoleListApi({
      page: pagination.page,
      pageSize: pagination.pageSize,
      roleName: searchForm.roleName || undefined,
    });

    tableData.value = res.list || [];
    pagination.total = res.total || 0;
  } catch (error: any) {
    ElMessage.error(error?.message || '加载失败');
    console.error('加载角色列表失败:', error);
  } finally {
    loading.value = false;
  }
}

// 新增
function handleAdd() {
  isEdit.value = false;
  form.id = '';
  form.roleName = '';
  form.roleCode = '';
  form.description = '';
  dialogVisible.value = true;
}

// 编辑
function handleEdit(row: RoleApi.Role) {
  isEdit.value = true;
  form.id = row.id;
  form.roleName = row.roleName;
  form.roleCode = row.roleCode;
  form.description = row.description;
  dialogVisible.value = true;
}

// 删除
async function handleDelete(row: RoleApi.Role) {
  try {
    await ElMessageBox.confirm(
      `确定要删除角色 "${row.roleName}" 吗？此操作不可恢复。`,
      '确认删除',
      { type: 'warning' }
    );

    await deleteRoleApi([row.id]);
    ElMessage.success('删除成功');
    loadData();
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error?.message || '删除失败');
      console.error('删除角色失败:', error);
    }
  }
}

// 权限配置
function handlePermission(row: RoleApi.Role) {
  currentRole.value = row;
  // 解析已有的菜单权限
  checkedMenuIds.value = row.menuIds ? row.menuIds.split(',').filter(Boolean) : [];
  permissionDialogVisible.value = true;
}

// 保存权限
async function handleSavePermission() {
  if (!currentRole.value) return;

  permissionLoading.value = true;
  try {
    const menuIds = treeRef.value?.getCheckedKeys() || [];
    await updateRolePermissionApi({
      id: currentRole.value.id,
      menuIds: menuIds as string[],
    });
    ElMessage.success('权限配置已保存');
    permissionDialogVisible.value = false;
    loadData();
  } catch (error: any) {
    ElMessage.error(error?.message || '保存失败');
    console.error('保存权限失败:', error);
  } finally {
    permissionLoading.value = false;
  }
}

// 提交
async function handleSubmit() {
  await formRef.value?.validate();
  submitLoading.value = true;
  try {
    if (isEdit.value) {
      // 编辑
      await updateRoleApi({
        id: form.id,
        roleName: form.roleName,
        description: form.description,
      });
      ElMessage.success('更新成功');
    } else {
      // 新增
      await createRoleApi({
        roleName: form.roleName,
        roleCode: form.roleCode,
        description: form.description,
      });
      ElMessage.success('新增成功');
    }
    dialogVisible.value = false;
    loadData();
  } catch (error: any) {
    ElMessage.error(error?.message || '操作失败');
    console.error('提交失败:', error);
  } finally {
    submitLoading.value = false;
  }
}

// 分页
function handleSizeChange(size: number) {
  pagination.pageSize = size;
  loadData();
}

function handleCurrentChange(page: number) {
  pagination.page = page;
  loadData();
}

onMounted(() => {
  loadData();
});
</script>