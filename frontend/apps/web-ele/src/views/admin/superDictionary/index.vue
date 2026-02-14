<script lang="ts" setup>
import { onMounted, ref, watch, reactive, onUnmounted } from 'vue';

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
  ElTree,
  ElDescriptions,
  ElDescriptionsItem,
  ElEmpty,
  ElTooltip,
  ElPagination,
} from 'element-plus';
import type { FormInstance, FormRules } from 'element-plus';

import {
  getSuperDictionaryListApi,
  getSuperDictionaryTreeApi,
  getSuperDictionaryDetailApi,
  createSuperDictionaryApi,
  updateSuperDictionaryApi,
  deleteSuperDictionaryApi,
  type SuperDictionaryApi,
} from '#/api/master/superDictionary';

// ===================== 状态定义 =====================

// 视图模式
const viewMode = ref<'tree' | 'list'>('tree');

// 加载状态
const loading = ref(false);
const treeLoading = ref(false);
const detailLoading = ref(false);
const submitLoading = ref(false);

// 树形数据
const treeData = ref<SuperDictionaryApi.TreeNode[]>([]);
const selectedNode = ref<SuperDictionaryApi.TreeNode | null>(null);
const expandedKeys = ref<number[]>([]);

// 列表数据
const tableData = ref<SuperDictionaryApi.SuperDictionary[]>([]);
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
});

// 弹窗
const dialogVisible = ref(false);
const dialogType = ref<'create' | 'edit' | 'createChild'>('create');
const formRef = ref<FormInstance>();

// 详情弹窗
const detailDialogVisible = ref(false);
const detailData = ref<SuperDictionaryApi.SuperDictionary | null>(null);

// 搜索表单
const searchForm = reactive({
  dictName: '',
  dictValue: '',
});

// 表单数据
const formData = reactive({
  id: 0,
  parentId: 0,
  dictValue: '',
  dictName: '',
  level: 1,
  levelName: '',
  description: '',
  eval: '',
  extSchema: '',
  extJson: '',
});

// 表单验证规则
const rules: FormRules = {
  dictValue: [
    { required: true, message: '请输入字典值', trigger: 'blur' },
    { min: 1, max: 100, message: '长度在 1 到 100 个字符', trigger: 'blur' },
  ],
  dictName: [
    { required: true, message: '请输入字典名称', trigger: 'blur' },
    { min: 1, max: 100, message: '长度在 1 到 100 个字符', trigger: 'blur' },
  ],
};

// ===================== 拖拽调整宽度 =====================

const leftWidth = ref(350); // 左侧默认宽度
const isDragging = ref(false);
const containerRef = ref<HTMLElement | null>(null);

const MIN_LEFT_WIDTH = 200;
const MAX_LEFT_WIDTH = 600;

function startDrag(e: MouseEvent) {
  isDragging.value = true;
  document.body.style.cursor = 'col-resize';
  document.body.style.userSelect = 'none';
  e.preventDefault();
}

function onDrag(e: MouseEvent) {
  if (!isDragging.value || !containerRef.value) return;

  const containerRect = containerRef.value.getBoundingClientRect();
  const newWidth = e.clientX - containerRect.left;

  leftWidth.value = Math.min(MAX_LEFT_WIDTH, Math.max(MIN_LEFT_WIDTH, newWidth));
}

function stopDrag() {
  isDragging.value = false;
  document.body.style.cursor = '';
  document.body.style.userSelect = '';
}

onMounted(() => {
  loadTreeData();
  loadListData();

  // 添加全局拖拽事件
  document.addEventListener('mousemove', onDrag);
  document.addEventListener('mouseup', stopDrag);
});

onUnmounted(() => {
  // 移除全局拖拽事件
  document.removeEventListener('mousemove', onDrag);
  document.removeEventListener('mouseup', stopDrag);
});

// ===================== 数据加载 =====================

// 加载树形数据
async function loadTreeData() {
  treeLoading.value = true;
  try {
    const res = await getSuperDictionaryTreeApi(0);
    treeData.value = res || [];
    // 默认展开第一层
    if (res && res.length > 0) {
      expandedKeys.value = res.map((item) => item.id);
    }
  } catch (error) {
    console.error('加载树形数据失败:', error);
    ElMessage.error('加载树形数据失败');
  } finally {
    treeLoading.value = false;
  }
}

// 加载列表数据
async function loadListData() {
  loading.value = true;
  try {
    const params: SuperDictionaryApi.ListParams = {
      page: pagination.page,
      pageSize: pagination.pageSize,
      ...searchForm,
    };
    const res = await getSuperDictionaryListApi(params);
    tableData.value = res.list || [];
    pagination.total = res.total || 0;
  } catch (error) {
    console.error('加载列表数据失败:', error);
    ElMessage.error('加载列表数据失败');
  } finally {
    loading.value = false;
  }
}

// 加载详情
async function loadDetail(id: number) {
  detailLoading.value = true;
  try {
    const res = await getSuperDictionaryDetailApi(id);
    detailData.value = res;
    detailDialogVisible.value = true;
  } catch (error) {
    console.error('加载详情失败:', error);
    ElMessage.error('加载详情失败');
  } finally {
    detailLoading.value = false;
  }
}

// ===================== 树形操作 =====================

// 树形配置
const treeProps = {
  label: 'dictName',
  children: 'children',
  value: 'id',
};

// 点击树节点
function handleNodeClick(data: SuperDictionaryApi.TreeNode) {
  selectedNode.value = data;
}

// ===================== 弹窗操作 =====================

// 打开新增弹窗
function openCreateDialog() {
  dialogType.value = 'create';
  resetForm();
  formData.parentId = 0;
  formData.level = 1;
  dialogVisible.value = true;
}

// 打开编辑弹窗
function openEditDialog(row: SuperDictionaryApi.SuperDictionary) {
  dialogType.value = 'edit';
  Object.assign(formData, {
    id: row.id,
    parentId: row.parentId,
    dictValue: row.dictValue,
    dictName: row.dictName,
    level: row.level,
    levelName: row.levelName,
    description: row.description,
    eval: row.eval,
    extSchema: row.extSchema,
    extJson: row.extJson,
  });
  dialogVisible.value = true;
}

// 打开新增子节点弹窗
function openCreateChildDialog(parent: SuperDictionaryApi.SuperDictionary) {
  dialogType.value = 'createChild';
  resetForm();
  formData.parentId = parent.id;
  formData.level = parent.level + 1;
  formData.levelName = `${parent.dictName} - 子项`;
  dialogVisible.value = true;
}

// 重置表单
function resetForm() {
  Object.assign(formData, {
    id: 0,
    parentId: 0,
    dictValue: '',
    dictName: '',
    level: 1,
    levelName: '',
    description: '',
    eval: '',
    extSchema: '',
    extJson: '',
  });
  formRef.value?.resetFields();
}

// 提交表单
async function handleSubmit() {
  if (!formRef.value) return;

  await formRef.value.validate(async (valid) => {
    if (!valid) return;

    submitLoading.value = true;
    try {
      if (dialogType.value === 'create' || dialogType.value === 'createChild') {
        await createSuperDictionaryApi({
          parentId: formData.parentId,
          dictValue: formData.dictValue,
          dictName: formData.dictName,
          level: formData.level,
          levelName: formData.levelName,
          description: formData.description,
          eval: formData.eval,
          extSchema: formData.extSchema,
          extJson: formData.extJson,
        });
        ElMessage.success('创建成功');
      } else {
        await updateSuperDictionaryApi({
          id: formData.id,
          parentId: formData.parentId,
          dictValue: formData.dictValue,
          dictName: formData.dictName,
          level: formData.level,
          levelName: formData.levelName,
          description: formData.description,
          eval: formData.eval,
          extSchema: formData.extSchema,
          extJson: formData.extJson,
        });
        ElMessage.success('更新成功');
      }
      dialogVisible.value = false;
      refreshData();
    } catch (error) {
      console.error('操作失败:', error);
      ElMessage.error('操作失败');
    } finally {
      submitLoading.value = false;
    }
  });
}

// ===================== 删除操作 =====================

async function handleDelete(row: SuperDictionaryApi.SuperDictionary) {
  try {
    await ElMessageBox.confirm(
      `确定要删除字典"${row.dictName}"吗？如果有子节点将一并删除，此操作不可恢复。`,
      '确认删除',
      { type: 'warning' }
    );
    await deleteSuperDictionaryApi([row.id]);
    ElMessage.success('删除成功');
    refreshData();
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error);
      ElMessage.error('删除失败');
    }
  }
}

// 批量删除
const selectedRows = ref<SuperDictionaryApi.SuperDictionary[]>([]);

function handleSelectionChange(rows: SuperDictionaryApi.SuperDictionary[]) {
  selectedRows.value = rows;
}

async function handleBatchDelete() {
  if (selectedRows.value.length === 0) {
    ElMessage.warning('请选择要删除的字典');
    return;
  }

  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedRows.value.length} 个字典吗？如果有子节点将一并删除，此操作不可恢复。`,
      '确认删除',
      { type: 'warning' }
    );
    const ids = selectedRows.value.map((row) => row.id);
    await deleteSuperDictionaryApi(ids);
    ElMessage.success('批量删除成功');
    refreshData();
  } catch (error) {
    if (error !== 'cancel') {
      console.error('批量删除失败:', error);
      ElMessage.error('批量删除失败');
    }
  }
}

// ===================== 搜索与分页 =====================

function handleSearch() {
  pagination.page = 1;
  loadListData();
}

function handleReset() {
  searchForm.dictName = '';
  searchForm.dictValue = '';
  handleSearch();
}

function handleSizeChange(size: number) {
  pagination.pageSize = size;
  loadListData();
}

function handleCurrentChange(page: number) {
  pagination.page = page;
  loadListData();
}

// 刷新数据
function refreshData() {
  if (viewMode.value === 'tree') {
    loadTreeData();
  } else {
    loadListData();
  }
  selectedNode.value = null;
}

// ===================== 工具函数 =====================

// 格式化日期
function formatDate(dateStr?: string): string {
  if (!dateStr) return '-';
  const date = new Date(dateStr);
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  });
}

// 获取层级标签类型
function getLevelTagType(level: number): 'success' | 'warning' | 'info' | 'danger' {
  const types: Record<number, 'success' | 'warning' | 'info' | 'danger'> = {
    1: 'danger',
    2: 'warning',
    3: 'success',
    4: 'info',
  };
  return types[level] || 'info';
}

// 获取层级名称
function getLevelName(level: number): string {
  const names: Record<number, string> = {
    1: '一级',
    2: '二级',
    3: '三级',
    4: '四级',
  };
  return names[level] || `${level}级`;
}

// ===================== 生命周期 =====================

// 监听视图模式切换
watch(viewMode, () => {
  refreshData();
});
</script>

<template>
  <Page
    title="超级字典管理"
    description="管理系统层级化字典数据，支持树形结构和列表视图"
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
          <ElFormItem label="字典名称">
            <ElInput
              v-model="searchForm.dictName"
              placeholder="请输入字典名称"
              clearable
              class="!w-48"
            />
          </ElFormItem>
          <ElFormItem label="字典值">
            <ElInput
              v-model="searchForm.dictValue"
              placeholder="请输入字典值"
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

      <!-- 主内容区 -->
      <div ref="containerRef" class="flex gap-0 relative">
        <!-- 左侧：树形结构 -->
        <div
          class="shrink-0"
          :style="{ width: leftWidth + 'px' }"
        >
          <ElCard shadow="hover" class="!rounded-xl h-full">
            <template #header>
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <span class="text-lg">🌳</span>
                  <span class="font-semibold text-gray-700">字典树</span>
                </div>
                <ElButton type="primary" size="small" @click="openCreateDialog">
                  <template #icon>
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
                    </svg>
                  </template>
                  新增根节点
                </ElButton>
              </div>
            </template>

            <div v-loading="treeLoading" class="min-h-[400px]">
              <ElEmpty v-if="treeData.length === 0 && !treeLoading" description="暂无字典数据">
                <ElButton type="primary" @click="openCreateDialog">创建第一个字典</ElButton>
              </ElEmpty>

              <ElTree
                v-else
                :data="treeData"
                :props="treeProps"
                node-key="id"
                :expand-on-click-node="false"
                :default-expanded-keys="expandedKeys"
                highlight-current
                class="!bg-transparent"
                @node-click="handleNodeClick"
              >
                <template #default="{ data }">
                  <div class="flex items-center justify-between w-full pr-2 group">
                    <div class="flex items-center gap-2 flex-1 min-w-0">
                      <ElTag
                        :type="getLevelTagType(data.level)"
                        size="small"
                        effect="plain"
                        class="!rounded-full shrink-0"
                      >
                        {{ getLevelName(data.level) }}
                      </ElTag>
                      <span class="truncate font-medium" :title="data.dictName">
                        {{ data.dictName }}
                      </span>
                      <span class="text-xs text-gray-400 truncate" :title="data.dictValue">
                        ({{ data.dictValue }})
                      </span>
                    </div>
                    <div class="opacity-0 group-hover:opacity-100 transition-opacity flex gap-1">
                      <ElTooltip content="新增子节点" placement="top">
                        <ElButton
                          type="primary"
                          size="small"
                          circle
                          @click.stop="openCreateChildDialog(data)"
                        >
                          <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
                          </svg>
                        </ElButton>
                      </ElTooltip>
                      <ElTooltip content="编辑" placement="top">
                        <ElButton
                          type="warning"
                          size="small"
                          circle
                          @click.stop="openEditDialog(data)"
                        >
                          <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h10a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
                          </svg>
                        </ElButton>
                      </ElTooltip>
                      <ElTooltip content="删除" placement="top">
                        <ElButton
                          type="danger"
                          size="small"
                          circle
                          @click.stop="handleDelete(data)"
                        >
                          <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                          </svg>
                        </ElButton>
                      </ElTooltip>
                    </div>
                  </div>
                </template>
              </ElTree>
            </div>
          </ElCard>
        </div>

        <!-- 可拖拽分隔条 -->
        <div
          class="w-1 bg-gray-200 hover:bg-blue-400 cursor-col-resize transition-colors flex items-center justify-center relative group"
          :class="{ 'bg-blue-500': isDragging }"
          @mousedown="startDrag"
        >
          <div class="absolute w-1 h-10 bg-gray-300 rounded-full group-hover:bg-blue-400 transition-colors" :class="{ 'bg-blue-500': isDragging }"></div>
        </div>

        <!-- 右侧：详情与列表 -->
        <div class="flex-1 min-w-0 space-y-6 pl-6">
          <!-- 选中节点详情 -->
          <ElCard v-if="selectedNode" shadow="hover" class="!rounded-xl">
            <template #header>
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <span class="text-lg">📋</span>
                  <span class="font-semibold text-gray-700">字典详情</span>
                </div>
                <div class="flex gap-2">
                  <ElButton type="primary" size="small" @click="openEditDialog(selectedNode)">
                    编辑
                  </ElButton>
                  <ElButton type="success" size="small" @click="openCreateChildDialog(selectedNode)">
                    新增子节点
                  </ElButton>
                </div>
              </div>
            </template>

            <ElDescriptions :column="2" border class="!rounded-lg overflow-hidden">
              <ElDescriptionsItem label="字典名称">
                <span class="font-medium">{{ selectedNode.dictName }}</span>
              </ElDescriptionsItem>
              <ElDescriptionsItem label="字典值">
                <code class="px-2 py-1 bg-gray-100 rounded text-sm">{{ selectedNode.dictValue }}</code>
              </ElDescriptionsItem>
              <ElDescriptionsItem label="层级">
                <ElTag :type="getLevelTagType(selectedNode.level)" effect="plain" class="!rounded-full">
                  {{ getLevelName(selectedNode.level) }}
                </ElTag>
              </ElDescriptionsItem>
              <ElDescriptionsItem label="层级名称">
                {{ selectedNode.levelName || '-' }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="父级ID">
                {{ selectedNode.parentId || '根节点' }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="ID">
                {{ selectedNode.id }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="描述" :span="2">
                {{ selectedNode.description || '-' }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="创建时间">
                {{ formatDate(selectedNode.createdAt) }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="更新时间">
                {{ formatDate(selectedNode.updatedAt) }}
              </ElDescriptionsItem>
              <ElDescriptionsItem v-if="selectedNode.eval" label="表达式" :span="2">
                <code class="px-2 py-1 bg-amber-50 text-amber-600 rounded text-sm">{{ selectedNode.eval }}</code>
              </ElDescriptionsItem>
            </ElDescriptions>

            <!-- 扩展数据 -->
            <div v-if="selectedNode.extJson" class="mt-4">
              <div class="text-sm font-semibold text-gray-600 mb-2">扩展数据</div>
              <pre class="bg-gray-50 p-3 rounded-lg text-sm overflow-auto max-h-32">{{ selectedNode.extJson }}</pre>
            </div>

            <!-- 子节点列表 -->
            <div v-if="selectedNode.children && selectedNode.children.length > 0" class="mt-4">
              <div class="text-sm font-semibold text-gray-600 mb-2">
                子节点 ({{ selectedNode.children.length }})
              </div>
              <div class="flex flex-wrap gap-2">
                <ElTag
                  v-for="child in selectedNode.children"
                  :key="child.id"
                  effect="plain"
                  class="!rounded-full cursor-pointer hover:!bg-blue-50"
                  @click="selectedNode = child"
                >
                  {{ child.dictName }}
                </ElTag>
              </div>
            </div>
          </ElCard>

          <!-- 字典列表 -->
          <ElCard shadow="hover" class="!rounded-xl">
            <template #header>
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <span class="text-lg">📊</span>
                  <span class="font-semibold text-gray-700">字典列表</span>
                  <span class="text-sm text-gray-400">(共 {{ pagination.total }} 条)</span>
                </div>
                <div class="flex gap-2">
                  <ElButton
                    type="danger"
                    size="small"
                    plain
                    :disabled="selectedRows.length === 0"
                    @click="handleBatchDelete"
                  >
                    批量删除 ({{ selectedRows.length }})
                  </ElButton>
                  <ElButton type="primary" size="small" @click="openCreateDialog">
                    <template #icon>
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
                      </svg>
                    </template>
                    新增字典
                  </ElButton>
                </div>
              </div>
            </template>

            <ElTable
              v-loading="loading"
              :data="tableData"
              stripe
              class="!rounded-lg"
              @selection-change="handleSelectionChange"
            >
              <ElTableColumn type="selection" width="50" />
              <ElTableColumn label="字典名称" min-width="160">
                <template #default="{ row }">
                  <div class="flex items-center gap-2">
                    <ElTag
                      :type="getLevelTagType(row.level)"
                      size="small"
                      effect="plain"
                      class="!rounded-full"
                    >
                      {{ getLevelName(row.level) }}
                    </ElTag>
                    <span class="font-medium">{{ row.dictName }}</span>
                  </div>
                </template>
              </ElTableColumn>
              <ElTableColumn label="字典值" prop="dictValue" min-width="120">
                <template #default="{ row }">
                  <code class="px-2 py-1 bg-gray-100 rounded text-sm">{{ row.dictValue }}</code>
                </template>
              </ElTableColumn>
              <ElTableColumn label="层级名称" prop="levelName" min-width="120">
                <template #default="{ row }">
                  <span class="text-gray-500">{{ row.levelName || '-' }}</span>
                </template>
              </ElTableColumn>
              <ElTableColumn label="描述" prop="description" min-width="200" show-overflow-tooltip />
              <ElTableColumn label="创建时间" prop="createdAt" min-width="160">
                <template #default="{ row }">
                  <span class="text-gray-500 text-sm">{{ formatDate(row.createdAt) }}</span>
                </template>
              </ElTableColumn>
              <ElTableColumn label="操作" width="200" fixed="right">
                <template #default="{ row }">
                  <div class="flex gap-2">
                    <ElButton type="primary" size="small" plain @click="loadDetail(row.id)">
                      详情
                    </ElButton>
                    <ElButton type="warning" size="small" plain @click="openEditDialog(row)">
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
      </div>
    </div>

    <!-- 新增/编辑对话框 -->
    <ElDialog
      v-model="dialogVisible"
      :title="dialogType === 'create' ? '新增根字典' : dialogType === 'createChild' ? '新增子字典' : '编辑字典'"
      width="600px"
      destroy-on-close
      :close-on-click-modal="false"
    >
      <ElForm
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-position="top"
      >
        <div class="grid grid-cols-2 gap-4">
          <ElFormItem label="字典值" prop="dictValue">
            <ElInput v-model="formData.dictValue" placeholder="请输入字典值（如：status_active）" />
          </ElFormItem>
          <ElFormItem label="字典名称" prop="dictName">
            <ElInput v-model="formData.dictName" placeholder="请输入字典名称（如：激活状态）" />
          </ElFormItem>
        </div>

        <div class="grid grid-cols-2 gap-4">
          <ElFormItem label="层级">
            <ElInputNumber v-model="formData.level" :min="1" :max="10" class="!w-full" disabled />
          </ElFormItem>
          <ElFormItem label="层级名称">
            <ElInput v-model="formData.levelName" placeholder="如：一级分类" />
          </ElFormItem>
        </div>

        <ElFormItem label="描述">
          <ElInput
            v-model="formData.description"
            type="textarea"
            :rows="3"
            placeholder="请输入字典描述"
          />
        </ElFormItem>

        <ElFormItem label="表达式 (Eval)">
          <ElInput
            v-model="formData.eval"
            placeholder="可选：用于条件计算的表达式"
          />
          <template #extra>
            <span class="text-xs text-gray-400">用于动态计算或条件判断的表达式</span>
          </template>
        </ElFormItem>

        <ElFormItem label="扩展 Schema">
          <ElInput
            v-model="formData.extSchema"
            type="textarea"
            :rows="2"
            placeholder="可选：JSON Schema 定义"
          />
        </ElFormItem>

        <ElFormItem label="扩展数据 (JSON)">
          <ElInput
            v-model="formData.extJson"
            type="textarea"
            :rows="3"
            placeholder="可选：扩展数据的 JSON 格式"
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

    <!-- 详情对话框 -->
    <ElDialog
      v-model="detailDialogVisible"
      title="字典详情"
      width="700px"
      destroy-on-close
    >
      <div v-loading="detailLoading">
        <ElDescriptions v-if="detailData" :column="2" border>
          <ElDescriptionsItem label="ID">
            {{ detailData.id }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="父级ID">
            {{ detailData.parentId || '根节点' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="字典名称">
            <span class="font-medium">{{ detailData.dictName }}</span>
          </ElDescriptionsItem>
          <ElDescriptionsItem label="字典值">
            <code class="px-2 py-1 bg-gray-100 rounded text-sm">{{ detailData.dictValue }}</code>
          </ElDescriptionsItem>
          <ElDescriptionsItem label="层级">
            <ElTag :type="getLevelTagType(detailData.level)" effect="plain" class="!rounded-full">
              {{ getLevelName(detailData.level) }}
            </ElTag>
          </ElDescriptionsItem>
          <ElDescriptionsItem label="层级名称">
            {{ detailData.levelName || '-' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="描述" :span="2">
            {{ detailData.description || '-' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="表达式">
            <code v-if="detailData.eval" class="px-2 py-1 bg-amber-50 text-amber-600 rounded text-sm">
              {{ detailData.eval }}
            </code>
            <span v-else class="text-gray-400">-</span>
          </ElDescriptionsItem>
          <ElDescriptionsItem label="创建时间">
            {{ formatDate(detailData.createdAt) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="更新时间">
            {{ formatDate(detailData.updatedAt) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem v-if="detailData.extSchema" label="扩展Schema" :span="2">
            <pre class="bg-gray-50 p-3 rounded text-sm overflow-auto max-h-32">{{ detailData.extSchema }}</pre>
          </ElDescriptionsItem>
          <ElDescriptionsItem v-if="detailData.extJson" label="扩展数据" :span="2">
            <pre class="bg-gray-50 p-3 rounded text-sm overflow-auto max-h-32">{{ detailData.extJson }}</pre>
          </ElDescriptionsItem>
        </ElDescriptions>
      </div>

      <template #footer>
        <ElButton @click="detailDialogVisible = false">关闭</ElButton>
        <ElButton v-if="detailData" type="primary" @click="openEditDialog(detailData); detailDialogVisible = false">
          编辑
        </ElButton>
      </template>
    </ElDialog>
  </Page>
</template>

<style scoped>
:deep(.el-tree-node__content) {
  height: auto;
  padding: 8px 0;
}

:deep(.el-tree-node__content:hover) {
  background-color: #f5f7fa;
}

:deep(.el-tree-node.is-current > .el-tree-node__content) {
  background-color: #ecf5ff;
}

/* 拖拽时禁止文字选中 */
.dragging {
  user-select: none;
}
</style>