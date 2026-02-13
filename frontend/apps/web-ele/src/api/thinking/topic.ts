import { requestClient } from '#/api/request';

export namespace ThinkingTopicApi {
  /** 课题状态 */
  export type TopicStatus = 0 | 1 | 2 | 3; // 0-进行中 1-已完成 2-已归档 3-草稿

  /** 课题信息 */
  export interface TopicInfo {
    id: number;
    title: string;
    description: string;
    background: string;
    goal: string;
    constraints: string;
    status: TopicStatus;
    statusText: string;
    userId: number;
    modelId: number;
    modelName: string;
    priority: number;
    priorityText: string;
    tags: string[] | null;
    analysisCount: number;
    actionCount: number;
    deadline?: string;
    completeTime?: string;
    createdAt: string;
    updatedAt: string;
  }

  /** 课题详情 */
  export interface TopicDetail extends TopicInfo {}

  /** 创建课题参数 */
  export interface CreateTopicParams {
    title: string;
    description?: string;
    background?: string;
    goal?: string;
    constraints?: string;
    priority?: number;
    tags?: string;
    deadline?: string;
  }

  /** 更新课题参数 */
  export interface UpdateTopicParams {
    id: number;
    title: string;
    description?: string;
    background?: string;
    goal?: string;
    constraints?: string;
    priority?: number;
    tags?: string;
    deadline?: string;
  }

  /** 课题列表查询参数 */
  export interface TopicListParams {
    page?: number;
    pageSize?: number;
    title?: string;
    userId?: number;
    modelId?: number;
    status?: TopicStatus;
    priority?: number;
  }

  /** 课题列表响应 */
  export interface TopicListResult {
    page: number;
    pageSize: number;
    total: number;
    list: TopicInfo[];
  }

  /** 更新课题状态参数 */
  export interface UpdateTopicStatusParams {
    id: number;
    status: TopicStatus;
  }

  /** 为课题选用模型参数 */
  export interface SelectModelParams {
    topicId: number;
    modelId: number;
    modelName?: string;
  }

  /** 课题统计 */
  export interface TopicStatistics {
    total: number;
    draftCount: number;
    progressCount: number;
    completeCount: number;
    archiveCount: number;
  }

  /** 删除课题参数 */
  export interface DelTopicParams {
    ids: number[];
  }
}

/**
 * 获取课题列表
 */
export async function getThinkingTopicListApi(params: ThinkingTopicApi.TopicListParams = {}) {
  return requestClient.get<ThinkingTopicApi.TopicListResult>('/thinking/topic/list', { params });
}

/**
 * 获取我的课题列表
 */
export async function getMyThinkingTopicListApi(params: ThinkingTopicApi.TopicListParams = {}) {
  return requestClient.get<ThinkingTopicApi.TopicListResult>('/thinking/topic/my', { params });
}

/**
 * 获取课题详情
 */
export async function getThinkingTopicDetailApi(id: number) {
  return requestClient.get<ThinkingTopicApi.TopicDetail>(`/thinking/topic/${id}`);
}

/**
 * 创建课题
 */
export async function createThinkingTopicApi(data: ThinkingTopicApi.CreateTopicParams) {
  return requestClient.post<ThinkingTopicApi.TopicInfo>('/thinking/topic', data);
}

/**
 * 更新课题
 */
export async function updateThinkingTopicApi(data: ThinkingTopicApi.UpdateTopicParams) {
  return requestClient.put<ThinkingTopicApi.TopicInfo>('/thinking/topic', data);
}

/**
 * 删除课题
 */
export async function deleteThinkingTopicApi(data: ThinkingTopicApi.DelTopicParams) {
  return requestClient.delete('/thinking/topic', { data });
}

/**
 * 更新课题状态
 */
export async function updateThinkingTopicStatusApi(data: ThinkingTopicApi.UpdateTopicStatusParams) {
  return requestClient.post('/thinking/topic/status', data);
}

/**
 * 为课题选用模型
 */
export async function selectModelForTopicApi(data: ThinkingTopicApi.SelectModelParams) {
  return requestClient.post('/thinking/topic/select-model', data);
}

/**
 * 移除课题模型
 */
export async function removeModelFromTopicApi(topicId: number) {
  return requestClient.post(`/thinking/topic/remove-model/${topicId}`);
}

/**
 * 完成课题
 */
export async function completeThinkingTopicApi(id: number) {
  return requestClient.post(`/thinking/topic/complete/${id}`);
}

/**
 * 归档课题
 */
export async function archiveThinkingTopicApi(id: number) {
  return requestClient.post(`/thinking/topic/archive/${id}`);
}

/**
 * 重新打开课题
 */
export async function reopenThinkingTopicApi(id: number) {
  return requestClient.post(`/thinking/topic/reopen/${id}`);
}

/**
 * 获取课题统计
 */
export async function getThinkingTopicStatisticsApi() {
  return requestClient.get<ThinkingTopicApi.TopicStatistics>('/thinking/topic/statistics');
}
