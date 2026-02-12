import { requestClient } from '#/api/request';

export namespace ActionApi {
  /** 行动状态 */
  export type ActionStatus = 'pending' | 'in_progress' | 'completed' | 'cancelled';

  /** 行动优先级 */
  export type ActionPriority = 'high' | 'medium' | 'low';

  /** 跟进记录 */
  export interface FollowUpRecord {
    id: string;
    content: string;
    createdAt: string;
  }

  /** 行动信息 */
  export interface Action {
    id: string;
    title: string;
    description: string;
    topicId: string;
    topicTitle: string;
    guidingPrinciple?: string; // 指导原则
    completionRate: number; // 完成度 0-100
    status: ActionStatus;
    priority: ActionPriority; // 优先级
    followUpRecords: FollowUpRecord[];
    createdAt: string;
    updatedAt: string;
    dueDate?: string;
  }

  /** 行动列表查询参数 */
  export interface ActionListParams {
    page?: number;
    pageSize?: number;
    status?: ActionStatus;
    priority?: ActionPriority;
    topicId?: string;
    keyword?: string;
    sortBy?: 'createdAt' | 'dueDate' | 'completionRate' | 'priority';
  }

  /** 行动列表响应 */
  export interface ActionListResult {
    list: Action[];
    total: number;
    page: number;
    pageSize: number;
  }

  /** 更新完成度参数 */
  export interface UpdateCompletionRateParams {
    id: string;
    completionRate: number;
  }

  /** 添加跟进记录参数 */
  export interface AddFollowUpParams {
    actionId: string;
    content: string;
  }

  /** 更新行动状态参数 */
  export interface UpdateStatusParams {
    id: string;
    status: ActionStatus;
  }
}

/**
 * 获取行动列表
 */
export async function getActionListApi(params: ActionApi.ActionListParams = {}) {
  return requestClient.get<ActionApi.ActionListResult>('/action/list', {
    params,
  });
}

/**
 * 获取行动详情
 */
export async function getActionDetailApi(id: string) {
  return requestClient.get<ActionApi.Action>(`/action/${id}`);
}

/**
 * 更新行动完成度
 */
export async function updateCompletionRateApi(data: ActionApi.UpdateCompletionRateParams) {
  return requestClient.put<ActionApi.Action>(`/action/${data.id}/completion`, {
    completionRate: data.completionRate,
  });
}

/**
 * 添加跟进记录
 */
export async function addFollowUpApi(data: ActionApi.AddFollowUpParams) {
  return requestClient.post<ActionApi.FollowUpRecord>(`/action/${data.actionId}/follow-up`, {
    content: data.content,
  });
}

/**
 * 更新行动状态
 */
export async function updateActionStatusApi(data: ActionApi.UpdateStatusParams) {
  return requestClient.put<ActionApi.Action>(`/action/${data.id}/status`, {
    status: data.status,
  });
}

/**
 * 获取用户的所有行动
 */
export async function getMyActionsApi(params: Omit<ActionApi.ActionListParams, 'userId'> = {}) {
  return requestClient.get<ActionApi.ActionListResult>('/action/my', {
    params,
  });
}
