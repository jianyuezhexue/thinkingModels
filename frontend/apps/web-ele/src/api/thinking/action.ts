import { requestClient } from '#/api/request';

export namespace ThinkingActionApi {
  /** 行动状态 */
  export type ActionStatus = 0 | 1 | 2 | 3; // 0-待开始 1-进行中 2-已完成 3-已取消

  /** 行动信息 */
  export interface ActionInfo {
    id: number;
    title: string;
    description: string;
    userId: number;
    topicId: number;
    topicTitle: string;
    analysisId: number;
    priority: number;
    priorityText: string;
    status: ActionStatus;
    statusText: string;
    progress: number;
    deadline?: string;
    completedAt?: string;
    guidePrinciple: string;
    followUpCount: number;
    isOverdue: boolean;
    createdAt: string;
    updatedAt: string;
  }

  /** 行动详情 */
  export interface ActionDetail extends ActionInfo {}

  /** 创建行动参数 */
  export interface CreateActionParams {
    title: string;
    description?: string;
    topicId?: number;
    topicTitle?: string;
    analysisId?: number;
    priority?: number;
    deadline?: string;
    guidePrinciple?: string;
  }

  /** 更新行动参数 */
  export interface UpdateActionParams {
    id: number;
    title: string;
    description?: string;
    priority?: number;
    deadline?: string;
    guidePrinciple?: string;
  }

  /** 更新行动进度参数 */
  export interface UpdateProgressParams {
    id: number;
    progress: number;
  }

  /** 从分析创建行动参数 */
  export interface CreateFromAnalysisParams {
    analysisId: number;
    topicId?: number;
    topicTitle?: string;
    actions: CreateActionBatch[];
  }

  /** 批量创建行动项 */
  export interface CreateActionBatch {
    title: string;
    description?: string;
    priority?: number;
    deadline?: string;
    guidePrinciple?: string;
  }

  /** 行动列表查询参数 */
  export interface ActionListParams {
    page?: number;
    pageSize?: number;
    title?: string;
    userId?: number;
    topicId?: number;
    analysisId?: number;
    status?: ActionStatus;
    priority?: number;
    isOverdue?: boolean;
  }

  /** 行动列表响应 */
  export interface ActionListResult {
    page: number;
    pageSize: number;
    total: number;
    list: ActionInfo[];
  }

  /** 行动统计 */
  export interface ActionStatistics {
    total: number;
    pendingCount: number;
    progressCount: number;
    completeCount: number;
    cancelCount: number;
    overdueCount: number;
  }

  /** 删除行动参数 */
  export interface DelActionParams {
    ids: number[];
  }
}

/**
 * 获取行动列表
 */
export async function getThinkingActionListApi(params: ThinkingActionApi.ActionListParams = {}) {
  return requestClient.get<ThinkingActionApi.ActionListResult>('/thinking/action/list', { params });
}

/**
 * 获取我的行动列表
 */
export async function getMyThinkingActionListApi(params: ThinkingActionApi.ActionListParams = {}) {
  return requestClient.get<ThinkingActionApi.ActionListResult>('/thinking/action/my', { params });
}

/**
 * 获取行动详情
 */
export async function getThinkingActionDetailApi(id: number) {
  return requestClient.get<ThinkingActionApi.ActionDetail>(`/thinking/action/${id}`);
}

/**
 * 创建行动
 */
export async function createThinkingActionApi(data: ThinkingActionApi.CreateActionParams) {
  return requestClient.post<ThinkingActionApi.ActionInfo>('/thinking/action', data);
}

/**
 * 更新行动
 */
export async function updateThinkingActionApi(data: ThinkingActionApi.UpdateActionParams) {
  return requestClient.put<ThinkingActionApi.ActionInfo>('/thinking/action', data);
}

/**
 * 删除行动
 */
export async function deleteThinkingActionApi(data: ThinkingActionApi.DelActionParams) {
  return requestClient.delete('/thinking/action', { data });
}

/**
 * 从分析创建行动
 */
export async function createActionsFromAnalysisApi(data: ThinkingActionApi.CreateFromAnalysisParams) {
  return requestClient.post<ThinkingActionApi.ActionInfo[]>('/thinking/action/from-analysis', data);
}

/**
 * 获取课题的行动列表
 */
export async function getActionsByTopicApi(topicId: number) {
  return requestClient.get<ThinkingActionApi.ActionInfo[]>(`/thinking/action/by-topic/${topicId}`);
}

/**
 * 获取分析的行动列表
 */
export async function getActionsByAnalysisApi(analysisId: number) {
  return requestClient.get<ThinkingActionApi.ActionInfo[]>(`/thinking/action/by-analysis/${analysisId}`);
}

/**
 * 更新行动进度
 */
export async function updateActionProgressApi(data: ThinkingActionApi.UpdateProgressParams) {
  return requestClient.post('/thinking/action/progress', data);
}

/**
 * 完成行动
 */
export async function completeThinkingActionApi(id: number) {
  return requestClient.post(`/thinking/action/complete/${id}`);
}

/**
 * 取消行动
 */
export async function cancelThinkingActionApi(id: number) {
  return requestClient.post(`/thinking/action/cancel/${id}`);
}

/**
 * 获取行动统计
 */
export async function getThinkingActionStatisticsApi() {
  return requestClient.get<ThinkingActionApi.ActionStatistics>('/thinking/action/statistics');
}
