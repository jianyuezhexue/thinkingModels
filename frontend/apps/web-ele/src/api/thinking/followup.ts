import { requestClient } from '#/api/request';

export namespace ThinkingFollowUpApi {
  /** 跟进信息 */
  export interface FollowUpInfo {
    id: number;
    actionId: number;
    userId: number;
    userName: string;
    content: string;
    progressBefore: number;
    progressAfter: number;
    progressDelta: number;
    createdAt: string;
    updatedAt: string;
  }

  /** 创建跟进参数 */
  export interface CreateFollowUpParams {
    actionId: number;
    content: string;
    progressBefore?: number;
    progressAfter?: number;
  }

  /** 更新跟进参数 */
  export interface UpdateFollowUpParams {
    id: number;
    content: string;
  }

  /** 跟进列表查询参数 */
  export interface FollowUpListParams {
    page?: number;
    pageSize?: number;
    actionId?: number;
    userId?: number;
  }

  /** 跟进列表响应 */
  export interface FollowUpListResult {
    page: number;
    pageSize: number;
    total: number;
    list: FollowUpInfo[];
  }

  /** 删除跟进参数 */
  export interface DelFollowUpParams {
    ids: number[];
  }
}

/**
 * 获取行动的跟进列表
 */
export async function getFollowUpsByActionApi(actionId: number) {
  return requestClient.get<ThinkingFollowUpApi.FollowUpInfo[]>(`/thinking/followup/by-action/${actionId}`);
}

/**
 * 获取跟进详情
 */
export async function getThinkingFollowUpDetailApi(id: number) {
  return requestClient.get<ThinkingFollowUpApi.FollowUpInfo>(`/thinking/followup/${id}`);
}

/**
 * 创建跟进
 */
export async function createThinkingFollowUpApi(data: ThinkingFollowUpApi.CreateFollowUpParams) {
  return requestClient.post<ThinkingFollowUpApi.FollowUpInfo>('/thinking/followup', data);
}

/**
 * 更新跟进
 */
export async function updateThinkingFollowUpApi(data: ThinkingFollowUpApi.UpdateFollowUpParams) {
  return requestClient.put<ThinkingFollowUpApi.FollowUpInfo>('/thinking/followup', data);
}

/**
 * 删除跟进
 */
export async function deleteThinkingFollowUpApi(data: ThinkingFollowUpApi.DelFollowUpParams) {
  return requestClient.delete('/thinking/followup', { data });
}
