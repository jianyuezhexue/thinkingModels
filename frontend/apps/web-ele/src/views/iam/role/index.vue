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
                  :style="{ backgroundColor: slotProps.row.roleColor || getRoleColor(slotProps.row.roleCode) }"
                >
                  {{ slotProps.row.roleName.charAt(0) }}
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
                {{ slotProps.row.userCount }}
              </ElTag>
            </template>
          </ElTableColumn>
          <ElTableColumn label="创建时间" prop="createTime" min-width="160" />
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
          <ElTag size="large" effect="dark" :color="currentRole?.roleColor">
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

// 角色接口
interface Role {
  id: string;
  roleName: string;
  roleCode: string;
  description: string;
  userCount: number;
  createTime: string;
  roleColor?: string;
}

// 菜单树节点
interface MenuNode {
  id: string;
  label: string;
  children?: MenuNode[];
}

// 搜索表单
const searchForm = reactive({
  roleName: '',
});

// 列表数据
const loading = ref(false);
const roleList = ref<Role[]>([]);
const tableData = ref<Role[]>([]);
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
});

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
const currentRole = ref<Role | null>(null);
const treeRef = ref<any>(null);

// 菜单树数据
const menuTree = ref<MenuNode[]>([
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

// 生成模拟数据
function generateMockData(): Role[] {
  return [
    {
      id: '1',
      roleName: '超级管理员',
      roleCode: 'admin',
      description: '系统超级管理员，拥有所有权限，可进行系统配置和用户管理',
      userCount: 3,
      createTime: '2024-01-15 10:30:00',
      roleColor: roleColorMap.admin,
    },
    {
      id: '2',
      roleName: '运营管理员',
      roleCode: 'operator',
      description: '运营人员，负责内容审核、数据分析和用户反馈处理',
      userCount: 12,
      createTime: '2024-01-15 11:20:00',
      roleColor: roleColorMap.operator,
    },
    {
      id: '3',
      roleName: '普通用户',
      roleCode: 'user',
      description: '普通注册用户，可使用基础功能创建和分享思维模型',
      userCount: 1256,
      createTime: '2024-01-16 14:20:00',
      roleColor: roleColorMap.user,
    },
    {
      id: '4',
      roleName: 'VIP会员',
      roleCode: 'vip',
      description: '付费VIP用户，享有高级功能、无限模型创建和AI辅助分析',
      userCount: 328,
      createTime: '2024-01-17 09:15:00',
      roleColor: roleColorMap.vip,
    },
    {
      id: '5',
      roleName: 'SVIP会员',
      roleCode: 'svip',
      description: '高级会员，享有VIP所有功能及团队协作、数据导出等特权',
      userCount: 86,
      createTime: '2024-02-01 16:30:00',
      roleColor: roleColorMap.svip,
    },
    {
      id: '6',
      roleName: '访客',
      roleCode: 'guest',
      description: '未注册用户，仅可浏览公开模型和基础内容',
      userCount: 2847,
      createTime: '2024-01-18 16:45:00',
      roleColor: roleColorMap.guest,
    },
    {
      id: '7',
      roleName: '测试用户',
      roleCode: 'tester',
      description: '内部测试人员，用于新功能验证和回归测试',
      userCount: 15,
      createTime: '2024-02-10 09:00:00',
      roleColor: roleColorMap.tester,
    },
    {
      id: '8',
      roleName: '内容编辑',
      roleCode: 'editor',
      description: '内容编辑人员，负责思维模型模板管理和推荐内容维护',
      userCount: 8,
      createTime: '2024-02-15 14:30:00',
      roleColor: roleColorMap.editor,
    },
  ];
}

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
    // 模拟 API 调用
    await new Promise((resolve) => setTimeout(resolve, 300));

    // 获取模拟数据
    const mockData = generateMockData();
    roleList.value = mockData;

    // 前端搜索过滤
    let filteredData = mockData;
    if (searchForm.roleName) {
      const keyword = searchForm.roleName.toLowerCase();
      filteredData = mockData.filter(
        (item) =>
          item.roleName.toLowerCase().includes(keyword) ||
          item.roleCode.toLowerCase().includes(keyword)
      );
    }

    // 分页处理
    pagination.total = filteredData.length;
    const start = (pagination.page - 1) * pagination.pageSize;
    const end = start + pagination.pageSize;
    tableData.value = filteredData.slice(start, end);
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
function handleEdit(row: Role) {
  isEdit.value = true;
  form.id = row.id;
  form.roleName = row.roleName;
  form.roleCode = row.roleCode;
  form.description = row.description;
  dialogVisible.value = true;
}

// 删除
async function handleDelete(row: Role) {
  try {
    await ElMessageBox.confirm(
      `确定要删除角色 "${row.roleName}" 吗？此操作不可恢复。`,
      '确认删除',
      { type: 'warning' }
    );
    ElMessage.success('删除成功');
    loadData();
  } catch {
    // 取消删除
  }
}

// 权限配置
function handlePermission(row: Role) {
  currentRole.value = row;
  permissionDialogVisible.value = true;
}

// 保存权限
async function handleSavePermission() {
  permissionLoading.value = true;
  try {
    await new Promise((resolve) => setTimeout(resolve, 500));
    ElMessage.success('权限配置已保存');
    permissionDialogVisible.value = false;
  } finally {
    permissionLoading.value = false;
  }
}

// 提交
async function handleSubmit() {
  await formRef.value?.validate();
  submitLoading.value = true;
  try {
    await new Promise((resolve) => setTimeout(resolve, 500));
    ElMessage.success(isEdit.value ? '更新成功' : '新增成功');
    dialogVisible.value = false;
    loadData();
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
