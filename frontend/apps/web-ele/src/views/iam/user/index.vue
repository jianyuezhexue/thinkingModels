<template>
  <Page
    description="管理系统用户信息和权限"
    title="用户管理"
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
          <ElFormItem label="用户名">
            <ElInput
              v-model="searchForm.username"
              placeholder="请输入用户名"
              clearable
              class="!w-48"
            />
          </ElFormItem>
          <ElFormItem label="手机号">
            <ElInput
              v-model="searchForm.phone"
              placeholder="请输入手机号"
              clearable
              class="!w-48"
            />
          </ElFormItem>
          <ElFormItem label="角色">
            <ElSelect
              v-model="searchForm.role"
              placeholder="请选择角色"
              clearable
              class="!w-40"
            >
              <ElOption
                v-for="role in roleOptions"
                :key="role.value"
                :label="role.label"
                :value="role.value"
              />
            </ElSelect>
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

      <!-- 用户列表 -->
      <ElCard shadow="hover" class="!rounded-xl">
        <template #header>
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <span class="text-lg">👥</span>
              <span class="font-semibold text-gray-700">用户列表</span>
              <span class="text-sm text-gray-400">(共 {{ pagination.total }} 条)</span>
            </div>
            <ElButton type="primary" @click="handleAdd">
              <template #icon>
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
                </svg>
              </template>
              新增用户
            </ElButton>
          </div>
        </template>

        <ElTable
          v-loading="loading"
          :data="userList"
          stripe
          class="!rounded-lg"
        >
          <ElTableColumn label="头像" width="80" align="center">
            <template #default="scope">
              <ElAvatar :src="scope.row.avatar" :size="40" class="bg-purple-100">
                {{ scope.row.username?.charAt(0)?.toUpperCase() }}
              </ElAvatar>
            </template>
          </ElTableColumn>
          <ElTableColumn label="用户名" prop="username" min-width="120" />
          <ElTableColumn label="手机号" prop="phone" min-width="120" />
          <ElTableColumn label="角色" min-width="120">
            <template #default="scope">
              <ElTag
                :type="scope.row.role === 'admin' ? 'danger' : scope.row.role === 'vip' ? 'warning' : 'info'"
                effect="light"
                class="!rounded-full"
              >
                {{ getRoleLabel(scope.row.role) }}
              </ElTag>
            </template>
          </ElTableColumn>
          <ElTableColumn label="充值额度" width="120" align="center">
            <template #default="scope">
              <div class="flex flex-col items-center">
                <span class="font-semibold text-green-600">
                  ¥{{ formatBalance(scope.row.balance) }}
                </span>
                <span class="text-xs text-gray-400">剩余额度</span>
              </div>
            </template>
          </ElTableColumn>
          <ElTableColumn label="模型数量" width="100" align="center">
            <template #default="scope">
              <div class="flex flex-col items-center">
                <span class="font-semibold text-blue-600">{{ scope.row.modelCount || 0 }}</span>
                <span class="text-xs text-gray-400">个模型</span>
              </div>
            </template>
          </ElTableColumn>
          <ElTableColumn label="状态" width="100" align="center">
            <template #default="scope">
              <ElTag
                :type="scope.row.status === 1 ? 'success' : 'danger'"
                effect="light"
                class="!rounded-full"
              >
                {{ scope.row.status === 1 ? '正常' : '禁用' }}
              </ElTag>
            </template>
          </ElTableColumn>
          <ElTableColumn label="创建时间" prop="createTime" min-width="160" />
          <ElTableColumn label="操作" width="120" fixed="right">
            <template #default="scope">
              <div class="flex gap-2">
                <ElButton type="primary" size="small" plain @click="handleEdit(scope.row)">
                  编辑
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
      :title="isEdit ? '编辑用户' : '新增用户'"
      width="500px"
      destroy-on-close
    >
      <ElForm
        ref="formRef"
        :model="form"
        :rules="rules"
        label-position="top"
      >
        <ElFormItem label="用户名" prop="username">
          <ElInput v-model="form.username" placeholder="请输入用户名" />
        </ElFormItem>
        <ElFormItem label="手机号" prop="phone">
          <ElInput v-model="form.phone" placeholder="请输入手机号" />
        </ElFormItem>
        <ElFormItem label="角色" prop="role">
          <ElSelect v-model="form.role" placeholder="请选择角色" class="w-full">
            <ElOption
              v-for="role in roleOptions"
              :key="role.value"
              :label="role.label"
              :value="role.value"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="状态" prop="status">
          <ElRadioGroup v-model="form.status">
            <ElRadio :label="1">正常</ElRadio>
            <ElRadio :label="0">禁用</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="dialogVisible = false">取消</ElButton>
        <ElButton type="primary" :loading="submitLoading" @click="handleSubmit">
          确定
        </ElButton>
      </template>
    </ElDialog>
  </Page>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import {
  ElAvatar,
  ElButton,
  ElCard,
  ElDialog,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElMessageBox,
  ElOption,
  ElPagination,
  ElRadio,
  ElRadioGroup,
  ElSelect,
  ElTable,
  ElTableColumn,
  ElTag,
} from 'element-plus';
import type { FormInstance, FormRules } from 'element-plus';

import { Page } from '@vben/common-ui';

// 用户接口
interface User {
  id: string;
  username: string;
  phone: string;
  role: string;
  status: number;
  avatar?: string;
  createTime: string;
  balance: number;
  modelCount: number;
}

// 搜索表单
const searchForm = reactive({
  username: '',
  phone: '',
  role: '',
});

// 角色选项
const roleOptions = [
  { value: 'admin', label: '管理员' },
  { value: 'user', label: '普通用户' },
  { value: 'vip', label: 'VIP用户' },
];

// 获取角色显示文本
function getRoleLabel(role: string) {
  const map: Record<string, string> = {
    admin: '管理员',
    user: '普通用户',
    vip: 'VIP用户',
  };
  return map[role] || role;
}

// 格式化金额
function formatBalance(balance: number): string {
  if (balance === undefined || balance === null) return '0.00';
  return balance.toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  });
}

// 列表数据
const loading = ref(false);
const userList = ref<User[]>([]);
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
  username: '',
  phone: '',
  role: '',
  status: 1,
});

// 表单校验规则
const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' },
  ],
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '手机号格式不正确', trigger: 'blur' },
  ],
  role: [{ required: true, message: '请选择角色', trigger: 'change' }],
};

// 搜索
function handleSearch() {
  pagination.page = 1;
  loadData();
}

// 重置
function handleReset() {
  searchForm.username = '';
  searchForm.phone = '';
  searchForm.role = '';
  handleSearch();
}

// 加载数据
async function loadData() {
  loading.value = true;
  try {
    // 模拟 API 调用
    await new Promise((resolve) => setTimeout(resolve, 500));

    // 模拟数据
    const mockData: User[] = [
      {
        id: '1',
        username: '张三',
        phone: '13800138001',
        role: 'admin',
        status: 1,
        avatar: '',
        createTime: '2024-01-15 10:30:00',
        balance: 9999.99,
        modelCount: 156,
      },
      {
        id: '2',
        username: '李四',
        phone: '13800138002',
        role: 'user',
        status: 1,
        avatar: '',
        createTime: '2024-01-16 14:20:00',
        balance: 0,
        modelCount: 5,
      },
      {
        id: '3',
        username: '王五',
        phone: '13800138003',
        role: 'vip',
        status: 0,
        avatar: '',
        createTime: '2024-01-17 09:15:00',
        balance: 299.5,
        modelCount: 42,
      },
      {
        id: '4',
        username: '赵六',
        phone: '13800138004',
        role: 'vip',
        status: 1,
        avatar: '',
        createTime: '2024-01-18 11:30:00',
        balance: 1999,
        modelCount: 89,
      },
      {
        id: '5',
        username: '钱七',
        phone: '13800138005',
        role: 'user',
        status: 1,
        avatar: '',
        createTime: '2024-01-19 16:45:00',
        balance: 50,
        modelCount: 12,
      },
    ];

    userList.value = mockData;
    pagination.total = 100;
  } finally {
    loading.value = false;
  }
}

// 新增
function handleAdd() {
  isEdit.value = false;
  form.id = '';
  form.username = '';
  form.phone = '';
  form.role = '';
  form.status = 1;
  dialogVisible.value = true;
}

// 编辑
function handleEdit(row: User) {
  isEdit.value = true;
  form.id = row.id;
  form.username = row.username;
  form.phone = row.phone;
  form.role = row.role;
  form.status = row.status;
  dialogVisible.value = true;
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
