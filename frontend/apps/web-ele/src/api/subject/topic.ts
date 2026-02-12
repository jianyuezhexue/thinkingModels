import { requestClient } from '#/api/request';

export namespace TopicApi {
  /** 课题状态 */
  export type TopicStatus = 'draft' | 'in_progress' | 'completed' | 'archived';

  /** 课题信息 */
  export interface Topic {
    id: string;
    title: string;
    description: string;
    userId: string;
    status: TopicStatus;
    modelId?: string;
    modelName?: string;
    analysisCount: number;
    createdAt: string;
    updatedAt: string;
  }

  /** 课题列表查询参数 */
  export interface TopicListParams {
    page?: number;
    pageSize?: number;
    status?: TopicStatus;
    keyword?: string;
  }

  /** 课题列表响应 */
  export interface TopicListResult {
    list: Topic[];
    total: number;
    page: number;
    pageSize: number;
  }

  /** 创建课题参数 */
  export interface CreateTopicParams {
    title: string;
    description: string;
  }

  /** 更新课题参数 */
  export interface UpdateTopicParams {
    id: string;
    title?: string;
    description?: string;
    status?: TopicStatus;
  }

  /** 选用模型参数 */
  export interface SelectModelParams {
    topicId: string;
    modelId: string;
  }
}

/**
 * 获取课题列表
 */
export async function getTopicListApi(params: TopicApi.TopicListParams = {}) {
  return requestClient.get<TopicApi.TopicListResult>('/subject/topic/list', {
    params,
  });
}

/**
 * 获取课题详情
 */
export async function getTopicDetailApi(id: string) {
  return requestClient.get<TopicApi.Topic>(`/subject/topic/${id}`);
}

/**
 * 创建课题
 */
export async function createTopicApi(data: TopicApi.CreateTopicParams) {
  return requestClient.post<TopicApi.Topic>('/subject/topic', data);
}

/**
 * 更新课题
 */
export async function updateTopicApi(data: TopicApi.UpdateTopicParams) {
  return requestClient.put<TopicApi.Topic>(`/subject/topic/${data.id}`, data);
}

/**
 * 删除课题
 */
export async function deleteTopicApi(id: string) {
  return requestClient.delete(`/subject/topic/${id}`);
}

/**
 * 为课题选用思维模型
 */
export async function selectModelForTopicApi(params: TopicApi.SelectModelParams) {
  return requestClient.post<TopicApi.Topic>(`/subject/topic/${params.topicId}/model`, {
    modelId: params.modelId,
  });
}

/**
 * 获取用户的所有课题
 */
export async function getMyTopicsApi(params: Omit<TopicApi.TopicListParams, 'userId'> = {}) {
  return requestClient.get<TopicApi.TopicListResult>('/subject/topic/my', {
    params,
  });
}
