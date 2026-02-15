import { requestClient } from '#/api/request';

export namespace OrderApi {
  /** 订单状态 */
  export type OrderStatus = 'pending' | 'paid' | 'completed' | 'cancelled' | 'refunded';

  /** 订单类型 */
  export type OrderType = 'model_purchase' | 'consultation' | 'subscription';

  /** 订单信息 */
  export interface Order {
    id: string;
    orderNo: string;
    type: OrderType;
    status: OrderStatus;
    userId: string;
    userName: string;
    userAvatar?: string;
    // 商品信息
    productId: string;
    productName: string;
    productCover?: string;
    productType: string;
    // 金额信息
    originalPrice: number;
    discountAmount: number;
    finalPrice: number;
    // 支付信息
    paymentMethod?: string;
    paidAt?: string;
    // 时间信息
    createdAt: string;
    updatedAt: string;
    completedAt?: string;
    cancelledAt?: string;
    // 其他信息
    remark?: string;
  }

  /** 订单列表查询参数 */
  export interface OrderListParams {
    page?: number;
    pageSize?: number;
    status?: OrderStatus;
    type?: OrderType;
    orderNo?: string;
    userName?: string;
    startDate?: string;
    endDate?: string;
    sortBy?: 'createdAt' | 'finalPrice';
    sortOrder?: 'asc' | 'desc';
  }

  /** 订单列表响应 */
  export interface OrderListResult {
    list: Order[];
    total: number;
    page: number;
    pageSize: number;
  }

  /** 订单统计数据 */
  export interface OrderStatistics {
    totalOrders: number;
    totalAmount: number;
    pendingOrders: number;
    pendingAmount: number;
    completedOrders: number;
    completedAmount: number;
    todayOrders: number;
    todayAmount: number;
  }
}

/**
 * 获取订单列表
 */
export async function getOrderListApi(params: OrderApi.OrderListParams = {}) {
  return requestClient.get<OrderApi.OrderListResult>('/order/list', {
    params,
  });
}

/**
 * 获取订单详情
 */
export async function getOrderDetailApi(id: string) {
  return requestClient.get<OrderApi.Order>(`/order/${id}`);
}

/**
 * 获取订单统计
 */
export async function getOrderStatisticsApi() {
  return requestClient.get<OrderApi.OrderStatistics>('/order/statistics');
}

/**
 * 更新订单状态
 */
export async function updateOrderStatusApi(data: { id: string; status: OrderApi.OrderStatus }) {
  return requestClient.put<OrderApi.Order>(`/order/${data.id}/status`, {
    status: data.status,
  });
}

/**
 * 取消订单
 */
export async function cancelOrderApi(id: string) {
  return requestClient.put<OrderApi.Order>(`/order/${id}/cancel`);
}

/**
 * 订单退款
 */
export async function refundOrderApi(id: string, reason: string) {
  return requestClient.put<OrderApi.Order>(`/order/${id}/refund`, {
    reason,
  });
}