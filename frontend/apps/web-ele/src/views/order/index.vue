<script lang="ts" setup>
import { onMounted, ref, watch, computed } from 'vue';
import { Page } from '@vben/common-ui';
import {
  ElButton,
  ElCard,
  ElDialog,
  ElEmpty,
  ElInput,
  ElMessage,
  ElPagination,
  ElSelect,
  ElOption,
  ElSkeleton,
  ElTag,
  ElDatePicker,
  ElDescriptions,
  ElDescriptionsItem,
  ElMessageBox,
} from 'element-plus';
import type { OrderApi } from '#/api/order';
import {
  getOrderListApi,
  getOrderStatisticsApi,
  updateOrderStatusApi,
  cancelOrderApi,
  refundOrderApi,
} from '#/api/order';

// ==================== 状态管理 ====================
const loading = ref(false);
const orders = ref<OrderApi.Order[]>([]);
const total = ref(0);
const statistics = ref<OrderApi.OrderStatistics | null>(null);

// 分页
const currentPage = ref(1);
const pageSize = ref(10);

// 筛选
const searchKeyword = ref('');
const activeStatus = ref<OrderApi.OrderStatus | 'all'>('all');
const activeType = ref<OrderApi.OrderType | 'all'>('all');
const dateRange = ref<[string, string] | null>(null);
const sortBy = ref<'createdAt' | 'finalPrice'>('createdAt');
const sortOrder = ref<'asc' | 'desc'>('desc');

// 弹窗状态
const detailDialogVisible = ref(false);
const selectedOrder = ref<OrderApi.Order | null>(null);

// ==================== 统计数据 ====================
const statusTabs = [
  { id: 'all', label: '全部订单', icon: '📋' },
  { id: 'pending', label: '待支付', icon: '⏳' },
  { id: 'paid', label: '已支付', icon: '💳' },
  { id: 'completed', label: '已完成', icon: '✅' },
  { id: 'cancelled', label: '已取消', icon: '❌' },
  { id: 'refunded', label: '已退款', icon: '💸' },
];

const typeOptions = [
  { value: 'all', label: '全部类型' },
  { value: 'model_purchase', label: '模型购买' },
  { value: 'consultation', label: '付费咨询' },
  { value: 'subscription', label: '会员订阅' },
];

// Mock 数据
const mockOrders: OrderApi.Order[] = [
  {
    id: '1',
    orderNo: 'ORD202401150001',
    type: 'model_purchase',
    status: 'completed',
    userId: 'u1',
    userName: '张三',
    userAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=1',
    productId: 'm1',
    productName: 'SWOT 分析模型',
    productCover: 'https://picsum.photos/seed/swot/200/120',
    productType: '思维模型',
    originalPrice: 29.9,
    discountAmount: 0,
    finalPrice: 29.9,
    paymentMethod: 'wechat',
    paidAt: '2024-01-15T10:30:00Z',
    createdAt: '2024-01-15T10:25:00Z',
    updatedAt: '2024-01-15T10:30:00Z',
    completedAt: '2024-01-15T10:30:00Z',
  },
  {
    id: '2',
    orderNo: 'ORD202401150002',
    type: 'consultation',
    status: 'paid',
    userId: 'u2',
    userName: '李四',
    userAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=2',
    productId: 'c1',
    productName: '职业规划咨询',
    productCover: 'https://picsum.photos/seed/consult/200/120',
    productType: '付费咨询',
    originalPrice: 199,
    discountAmount: 20,
    finalPrice: 179,
    paymentMethod: 'alipay',
    paidAt: '2024-01-15T14:20:00Z',
    createdAt: '2024-01-15T14:15:00Z',
    updatedAt: '2024-01-15T14:20:00Z',
  },
  {
    id: '3',
    orderNo: 'ORD202401160001',
    type: 'model_purchase',
    status: 'pending',
    userId: 'u3',
    userName: '王五',
    userAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=3',
    productId: 'm2',
    productName: '第一性原理思考法',
    productCover: 'https://picsum.photos/seed/principle/200/120',
    productType: '思维模型',
    originalPrice: 49,
    discountAmount: 10,
    finalPrice: 39,
    createdAt: '2024-01-16T09:00:00Z',
    updatedAt: '2024-01-16T09:00:00Z',
  },
  {
    id: '4',
    orderNo: 'ORD202401160002',
    type: 'subscription',
    status: 'completed',
    userId: 'u4',
    userName: '赵六',
    userAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=4',
    productId: 's1',
    productName: '年度会员',
    productCover: 'https://picsum.photos/seed/vip/200/120',
    productType: '会员订阅',
    originalPrice: 299,
    discountAmount: 100,
    finalPrice: 199,
    paymentMethod: 'wechat',
    paidAt: '2024-01-16T11:00:00Z',
    createdAt: '2024-01-16T10:55:00Z',
    updatedAt: '2024-01-16T11:00:00Z',
    completedAt: '2024-01-16T11:00:00Z',
  },
  {
    id: '5',
    orderNo: 'ORD202401170001',
    type: 'model_purchase',
    status: 'refunded',
    userId: 'u5',
    userName: '孙七',
    userAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=5',
    productId: 'm3',
    productName: '金字塔原理',
    productCover: 'https://picsum.photos/seed/pyramid/200/120',
    productType: '思维模型',
    originalPrice: 39,
    discountAmount: 0,
    finalPrice: 39,
    paymentMethod: 'alipay',
    paidAt: '2024-01-17T08:30:00Z',
    createdAt: '2024-01-17T08:25:00Z',
    updatedAt: '2024-01-17T15:00:00Z',
    remark: '用户申请退款',
  },
  {
    id: '6',
    orderNo: 'ORD202401170002',
    type: 'consultation',
    status: 'cancelled',
    userId: 'u6',
    userName: '周八',
    userAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=6',
    productId: 'c2',
    productName: '创业咨询',
    productCover: 'https://picsum.photos/seed/startup/200/120',
    productType: '付费咨询',
    originalPrice: 299,
    discountAmount: 0,
    finalPrice: 299,
    createdAt: '2024-01-17T10:00:00Z',
    updatedAt: '2024-01-17T10:30:00Z',
    cancelledAt: '2024-01-17T10:30:00Z',
  },
];

const mockStatistics: OrderApi.OrderStatistics = {
  totalOrders: 156,
  totalAmount: 12580.5,
  pendingOrders: 12,
  pendingAmount: 890,
  completedOrders: 130,
  completedAmount: 10890.5,
  todayOrders: 8,
  todayAmount: 456,
};

// ==================== 数据获取 ====================
async function fetchOrders() {
  loading.value = true;
  try {
    // 使用 Mock 数据
    let filtered = [...mockOrders];

    // 状态筛选
    if (activeStatus.value !== 'all') {
      filtered = filtered.filter((o) => o.status === activeStatus.value);
    }

    // 类型筛选
    if (activeType.value !== 'all') {
      filtered = filtered.filter((o) => o.type === activeType.value);
    }

    // 关键词搜索
    if (searchKeyword.value) {
      const kw = searchKeyword.value.toLowerCase();
      filtered = filtered.filter(
        (o) =>
          o.orderNo.toLowerCase().includes(kw) ||
          o.userName.toLowerCase().includes(kw) ||
          o.productName.toLowerCase().includes(kw),
      );
    }

    // 排序
    filtered.sort((a, b) => {
      const aValue = sortBy.value === 'createdAt' ? new Date(a.createdAt).getTime() : a.finalPrice;
      const bValue = sortBy.value === 'createdAt' ? new Date(b.createdAt).getTime() : b.finalPrice;
      return sortOrder.value === 'desc' ? bValue - aValue : aValue - bValue;
    });

    orders.value = filtered;
    total.value = filtered.length;
    statistics.value = mockStatistics;
  } catch (error) {
    console.error('获取订单列表失败:', error);
    ElMessage.error('获取订单列表失败');
  } finally {
    loading.value = false;
  }
}

// ==================== 工具函数 ====================
function getStatusStyle(status: OrderApi.OrderStatus): string {
  const styles: Record<string, string> = {
    pending: 'bg-amber-100 text-amber-700',
    paid: 'bg-blue-100 text-blue-700',
    completed: 'bg-green-100 text-green-700',
    cancelled: 'bg-gray-100 text-gray-600',
    refunded: 'bg-purple-100 text-purple-700',
  };
  return styles[status] || 'bg-gray-100 text-gray-600';
}

function getStatusText(status: OrderApi.OrderStatus): string {
  const texts: Record<string, string> = {
    pending: '待支付',
    paid: '已支付',
    completed: '已完成',
    cancelled: '已取消',
    refunded: '已退款',
  };
  return texts[status] || status;
}

function getStatusIcon(status: OrderApi.OrderStatus): string {
  const icons: Record<string, string> = {
    pending: '⏳',
    paid: '💳',
    completed: '✅',
    cancelled: '❌',
    refunded: '💸',
  };
  return icons[status] || '📋';
}

function getTypeText(type: OrderApi.OrderType): string {
  const texts: Record<string, string> = {
    model_purchase: '模型购买',
    consultation: '付费咨询',
    subscription: '会员订阅',
  };
  return texts[type] || type;
}

function getTypeIcon(type: OrderApi.OrderType): string {
  const icons: Record<string, string> = {
    model_purchase: '🧠',
    consultation: '💬',
    subscription: '👑',
  };
  return icons[type] || '📦';
}

function getPaymentMethodText(method?: string): string {
  const texts: Record<string, string> = {
    wechat: '微信支付',
    alipay: '支付宝',
    balance: '余额支付',
  };
  return method ? texts[method] || method : '-';
}

function formatDate(dateStr: string): string {
  if (!dateStr) return '-';
  return new Date(dateStr).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: 'numeric',
  });
}

function formatDateTime(dateStr: string): string {
  if (!dateStr) return '-';
  const date = new Date(dateStr);
  return (
    date.getFullYear() +
    '-' +
    String(date.getMonth() + 1).padStart(2, '0') +
    '-' +
    String(date.getDate()).padStart(2, '0') +
    ' ' +
    String(date.getHours()).padStart(2, '0') +
    ':' +
    String(date.getMinutes()).padStart(2, '0')
  );
}

function formatPrice(price: number): string {
  return '¥' + price.toFixed(2);
}

// ==================== 操作 ====================
function openOrderDetail(order: OrderApi.Order) {
  selectedOrder.value = order;
  detailDialogVisible.value = true;
}

async function handleCompleteOrder(order: OrderApi.Order) {
  try {
    await ElMessageBox.confirm('确认将该订单标记为已完成？', '确认操作', {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'success',
    });
    await updateOrderStatusApi({ id: order.id, status: 'completed' });
    order.status = 'completed';
    order.completedAt = new Date().toISOString();
    ElMessage.success('订单已完成');
  } catch (error) {
    if (error !== 'cancel') {
      console.error('操作失败:', error);
      ElMessage.error('操作失败');
    }
  }
}

async function handleCancelOrder(order: OrderApi.Order) {
  try {
    await ElMessageBox.confirm('确认取消该订单？此操作不可撤销。', '确认取消', {
      confirmButtonText: '确认取消',
      cancelButtonText: '返回',
      type: 'warning',
    });
    await cancelOrderApi(order.id);
    order.status = 'cancelled';
    order.cancelledAt = new Date().toISOString();
    ElMessage.success('订单已取消');
  } catch (error) {
    if (error !== 'cancel') {
      console.error('操作失败:', error);
      ElMessage.error('操作失败');
    }
  }
}

async function handleRefundOrder(order: OrderApi.Order) {
  try {
    const { value: reason } = await ElMessageBox.prompt('请输入退款原因', '申请退款', {
      confirmButtonText: '确认退款',
      cancelButtonText: '取消',
      inputPlaceholder: '请输入退款原因',
      inputValidator: (value) => {
        if (!value || !value.trim()) {
          return '请输入退款原因';
        }
        return true;
      },
    });
    await refundOrderApi(order.id, reason);
    order.status = 'refunded';
    order.remark = reason;
    ElMessage.success('退款申请已提交');
  } catch (error) {
    if (error !== 'cancel') {
      console.error('操作失败:', error);
      ElMessage.error('操作失败');
    }
  }
}

// ==================== 监听器 ====================
watch([activeStatus, activeType, searchKeyword, sortBy, sortOrder], () => {
  currentPage.value = 1;
  fetchOrders();
});

watch([currentPage, pageSize], () => {
  fetchOrders();
});

onMounted(() => {
  fetchOrders();
});
</script>

<template>
  <Page
    title="订单管理"
    description="管理平台所有订单，查看订单状态和处理退款申请"
    content-class="p-6 bg-gray-50"
  >
    <!-- 主内容区 -->
    <div class="flex gap-6">
      <!-- 左侧订单列表 -->
      <div class="flex-1 space-y-6">
        <!-- 统计卡片 -->
        <div class="grid grid-cols-4 gap-4">
          <ElCard shadow="hover" class="!rounded-xl text-center">
            <div class="text-3xl font-bold text-emerald-600">{{ statistics?.totalOrders || 0 }}</div>
            <div class="text-gray-500 text-sm mt-1">总订单数</div>
          </ElCard>
          <ElCard shadow="hover" class="!rounded-xl text-center">
            <div class="text-3xl font-bold text-blue-600">{{ statistics?.pendingOrders || 0 }}</div>
            <div class="text-gray-500 text-sm mt-1">待处理订单</div>
          </ElCard>
          <ElCard shadow="hover" class="!rounded-xl text-center">
            <div class="text-3xl font-bold text-amber-600">{{ formatPrice(statistics?.todayAmount || 0) }}</div>
            <div class="text-gray-500 text-sm mt-1">今日交易额</div>
          </ElCard>
          <ElCard shadow="hover" class="!rounded-xl text-center">
            <div class="text-3xl font-bold text-purple-600">{{ formatPrice(statistics?.totalAmount || 0) }}</div>
            <div class="text-gray-500 text-sm mt-1">累计交易额</div>
          </ElCard>
        </div>

        <!-- 筛选和搜索 -->
        <ElCard shadow="hover" class="!rounded-xl">
          <div class="flex flex-wrap items-center gap-4">
            <!-- 状态 Tab -->
            <div class="flex gap-2 flex-wrap">
              <button
                v-for="tab in statusTabs"
                :key="tab.id"
                class="px-4 py-2 rounded-full text-sm font-medium transition-all"
                :class="[
                  activeStatus === tab.id
                    ? 'bg-emerald-100 text-emerald-700 shadow-md border border-emerald-300 font-semibold'
                    : 'bg-gray-100 text-gray-700 hover:bg-emerald-50 hover:text-emerald-600 border border-gray-200',
                ]"
                @click="activeStatus = tab.id as OrderApi.OrderStatus | 'all'"
              >
                {{ tab.icon }} {{ tab.label }}
              </button>
            </div>
            <div class="flex-1" />
            <!-- 筛选项 -->
            <div class="flex items-center gap-3 flex-wrap">
              <ElSelect v-model="activeType" placeholder="订单类型" clearable class="!w-36">
                <ElOption
                  v-for="opt in typeOptions"
                  :key="opt.value"
                  :label="opt.label"
                  :value="opt.value"
                />
              </ElSelect>
              <ElInput
                v-model="searchKeyword"
                placeholder="搜索订单号/用户/商品..."
                clearable
                class="!w-56"
                @keyup.enter="fetchOrders"
              />
            </div>
          </div>
        </ElCard>

        <!-- 加载状态 -->
        <div v-if="loading" class="space-y-4">
          <ElSkeleton v-for="i in 3" :key="i" :rows="3" animated class="bg-white rounded-xl p-4" />
        </div>

        <!-- 空状态 -->
        <ElCard v-else-if="orders.length === 0" shadow="hover" class="!rounded-xl">
          <ElEmpty description="暂无订单数据">
            <template #image>
              <div class="text-6xl">📋</div>
            </template>
          </ElEmpty>
        </ElCard>

        <!-- 订单列表 -->
        <div v-else class="space-y-4">
          <ElCard
            v-for="order in orders"
            :key="order.id"
            shadow="hover"
            class="!rounded-xl cursor-pointer hover:shadow-lg transition-all group"
            @click="openOrderDetail(order)"
          >
            <div class="flex gap-4">
              <!-- 商品封面 -->
              <div class="flex-shrink-0">
                <div
                  class="w-20 h-14 rounded-lg bg-gray-100 bg-cover bg-center"
                  :style="{ backgroundImage: `url(${order.productCover})` }"
                />
              </div>

              <!-- 订单信息 -->
              <div class="flex-1 min-w-0">
                <div class="flex items-start justify-between gap-4 mb-2">
                  <div class="flex items-center gap-2">
                    <span class="text-xs text-gray-400">{{ order.orderNo }}</span>
                    <ElTag size="small" :class="getStatusStyle(order.status)" class="!border-0">
                      {{ getStatusIcon(order.status) }} {{ getStatusText(order.status) }}
                    </ElTag>
                    <span class="text-xs text-gray-400">{{ getTypeIcon(order.type) }} {{ getTypeText(order.type) }}</span>
                  </div>
                  <div class="text-lg font-bold text-emerald-600">
                    {{ formatPrice(order.finalPrice) }}
                  </div>
                </div>

                <h3 class="text-base font-medium text-gray-800 group-hover:text-emerald-600 transition-colors line-clamp-1 mb-1">
                  {{ order.productName }}
                </h3>

                <!-- 底部元信息 -->
                <div class="flex flex-wrap items-center gap-4 text-xs text-gray-400">
                  <span class="flex items-center gap-1">
                    <img
                      v-if="order.userAvatar"
                      :src="order.userAvatar"
                      class="w-4 h-4 rounded-full"
                      alt="用户头像"
                    />
                    {{ order.userName }}
                  </span>
                  <span v-if="order.paymentMethod">
                    💳 {{ getPaymentMethodText(order.paymentMethod) }}
                  </span>
                  <span>{{ formatDateTime(order.createdAt) }}</span>
                  <span v-if="order.discountAmount > 0" class="text-red-400">
                    已优惠 ¥{{ order.discountAmount.toFixed(2) }}
                  </span>
                </div>
              </div>

              <!-- 操作按钮 -->
              <div class="flex-shrink-0 flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                <ElButton
                  v-if="order.status === 'paid'"
                  type="success"
                  size="small"
                  @click.stop="handleCompleteOrder(order)"
                >
                  完成订单
                </ElButton>
                <ElButton
                  v-if="order.status === 'pending'"
                  type="warning"
                  size="small"
                  @click.stop="handleCancelOrder(order)"
                >
                  取消订单
                </ElButton>
                <ElButton
                  v-if="order.status === 'paid'"
                  type="info"
                  size="small"
                  @click.stop="handleRefundOrder(order)"
                >
                  退款
                </ElButton>
              </div>
            </div>
          </ElCard>

          <!-- 分页 -->
          <div class="flex justify-center pt-4">
            <ElPagination
              v-model:current-page="currentPage"
              :page-size="pageSize"
              :total="total"
              layout="prev, pager, next"
              background
            />
          </div>
        </div>
      </div>

      <!-- 右侧边栏 -->
      <div class="w-80 flex-shrink-0 space-y-6 hidden lg:block">
        <!-- 快捷筛选 -->
        <ElCard shadow="hover" class="!rounded-xl">
          <template #header>
            <span class="font-semibold text-gray-700">快捷筛选</span>
          </template>
          <div class="space-y-2">
            <button
              class="w-full flex items-center justify-between p-3 rounded-lg hover:bg-gray-50 transition-colors text-left"
              @click="sortBy = 'createdAt'; sortOrder = 'desc'"
            >
              <span class="flex items-center gap-2 text-sm">📅 最新订单</span>
            </button>
            <button
              class="w-full flex items-center justify-between p-3 rounded-lg hover:bg-gray-50 transition-colors text-left"
              @click="sortBy = 'finalPrice'; sortOrder = 'desc'"
            >
              <span class="flex items-center gap-2 text-sm">💰 金额最高</span>
            </button>
            <button
              class="w-full flex items-center justify-between p-3 rounded-lg hover:bg-gray-50 transition-colors text-left"
              @click="activeStatus = 'pending'"
            >
              <span class="flex items-center gap-2 text-sm">⏳ 待支付订单</span>
              <span class="text-xs text-amber-500 font-medium">{{ statistics?.pendingOrders || 0 }}</span>
            </button>
          </div>
        </ElCard>

        <!-- 订单类型分布 -->
        <ElCard shadow="hover" class="!rounded-xl">
          <template #header>
            <span class="font-semibold text-gray-700">订单类型分布</span>
          </template>
          <div class="space-y-3">
            <div class="flex items-center justify-between">
              <span class="flex items-center gap-2 text-sm">
                <span>🧠</span> 模型购买
              </span>
              <span class="text-gray-600 font-medium">{{ orders.filter(o => o.type === 'model_purchase').length }}</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="flex items-center gap-2 text-sm">
                <span>💬</span> 付费咨询
              </span>
              <span class="text-gray-600 font-medium">{{ orders.filter(o => o.type === 'consultation').length }}</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="flex items-center gap-2 text-sm">
                <span>👑</span> 会员订阅
              </span>
              <span class="text-gray-600 font-medium">{{ orders.filter(o => o.type === 'subscription').length }}</span>
            </div>
          </div>
        </ElCard>

        <!-- 订单管理提示 -->
        <ElCard shadow="hover" class="!rounded-xl !bg-gradient-to-br from-blue-50 to-blue-100 !border-blue-200">
          <template #header>
            <span class="font-semibold text-blue-700">订单管理提示</span>
          </template>
          <ul class="text-sm text-blue-800 space-y-2">
            <li class="flex items-start gap-2">
              <span class="text-blue-500">•</span>
              待支付订单超时未支付将自动取消
            </li>
            <li class="flex items-start gap-2">
              <span class="text-blue-500">•</span>
              退款将在 3-5 个工作日内原路返回
            </li>
            <li class="flex items-start gap-2">
              <span class="text-blue-500">•</span>
              大额订单建议人工确认后再处理
            </li>
          </ul>
        </ElCard>
      </div>
    </div>

    <!-- 订单详情弹窗 -->
    <ElDialog v-model="detailDialogVisible" title="订单详情" width="600px">
      <div v-if="selectedOrder" class="space-y-6">
        <!-- 订单头部 -->
        <div class="p-4 bg-gray-50 rounded-lg">
          <div class="flex items-center justify-between mb-3">
            <span class="text-sm text-gray-500">{{ selectedOrder.orderNo }}</span>
            <ElTag :class="getStatusStyle(selectedOrder.status)" class="!border-0">
              {{ getStatusIcon(selectedOrder.status) }} {{ getStatusText(selectedOrder.status) }}
            </ElTag>
          </div>
          <h2 class="text-lg font-semibold text-gray-800">{{ selectedOrder.productName }}</h2>
          <p class="text-sm text-gray-500 mt-1">{{ selectedOrder.productType }}</p>
        </div>

        <!-- 订单详情 -->
        <ElDescriptions :column="2" border>
          <ElDescriptionsItem label="订单类型">
            {{ getTypeIcon(selectedOrder.type) }} {{ getTypeText(selectedOrder.type) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="支付方式">
            {{ getPaymentMethodText(selectedOrder.paymentMethod) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="下单用户">
            <span class="flex items-center gap-2">
              <img
                v-if="selectedOrder.userAvatar"
                :src="selectedOrder.userAvatar"
                class="w-5 h-5 rounded-full"
                alt="用户头像"
              />
              {{ selectedOrder.userName }}
            </span>
          </ElDescriptionsItem>
          <ElDescriptionsItem label="创建时间">
            {{ formatDateTime(selectedOrder.createdAt) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem v-if="selectedOrder.paidAt" label="支付时间">
            {{ formatDateTime(selectedOrder.paidAt) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem v-if="selectedOrder.completedAt" label="完成时间">
            {{ formatDateTime(selectedOrder.completedAt) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem v-if="selectedOrder.cancelledAt" label="取消时间">
            {{ formatDateTime(selectedOrder.cancelledAt) }}
          </ElDescriptionsItem>
        </ElDescriptions>

        <!-- 金额信息 -->
        <div class="p-4 bg-emerald-50 rounded-lg">
          <h3 class="font-medium text-emerald-700 mb-3">金额明细</h3>
          <div class="space-y-2 text-sm">
            <div class="flex justify-between">
              <span class="text-gray-600">商品原价</span>
              <span class="text-gray-800">{{ formatPrice(selectedOrder.originalPrice) }}</span>
            </div>
            <div v-if="selectedOrder.discountAmount > 0" class="flex justify-between">
              <span class="text-gray-600">优惠金额</span>
              <span class="text-red-500">-{{ formatPrice(selectedOrder.discountAmount) }}</span>
            </div>
            <div class="flex justify-between pt-2 border-t border-emerald-200">
              <span class="font-medium text-emerald-700">实付金额</span>
              <span class="text-lg font-bold text-emerald-600">{{ formatPrice(selectedOrder.finalPrice) }}</span>
            </div>
          </div>
        </div>

        <!-- 备注 -->
        <div v-if="selectedOrder.remark" class="p-4 bg-amber-50 rounded-lg">
          <h3 class="font-medium text-amber-700 mb-2">备注</h3>
          <p class="text-sm text-amber-800">{{ selectedOrder.remark }}</p>
        </div>
      </div>

      <template #footer>
        <div class="flex justify-between">
          <div class="flex gap-2">
            <ElButton
              v-if="selectedOrder?.status === 'paid'"
              type="success"
              @click="handleCompleteOrder(selectedOrder!); detailDialogVisible = false"
            >
              完成订单
            </ElButton>
            <ElButton
              v-if="selectedOrder?.status === 'paid'"
              type="info"
              @click="handleRefundOrder(selectedOrder!); detailDialogVisible = false"
            >
              申请退款
            </ElButton>
          </div>
          <ElButton @click="detailDialogVisible = false">关闭</ElButton>
        </div>
      </template>
    </ElDialog>
  </Page>
</template>

<style scoped>
.line-clamp-1 {
  display: -webkit-box;
  -webkit-line-clamp: 1;
  line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>