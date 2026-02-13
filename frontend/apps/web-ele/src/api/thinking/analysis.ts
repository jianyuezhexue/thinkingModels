import { requestClient } from '#/api/request';

export namespace ThinkingAnalysisApi {
  /** 分析状态 */
  export type AnalysisStatus = 0 | 1 | 2 | 3; // 0-待分析 1-分析中 2-已完成 3-失败

  /** 分析信息 */
  export interface AnalysisInfo {
    id: number;
    topicId: number;
    modelId: number;
    modelName: string;
    content: string;
    aiAnalysis: string;
    aiSuggestions: string;
    version: number;
    isCurrent: boolean;
    userId: number;
    status: AnalysisStatus;
    statusText: string;
    createdAt: string;
    updatedAt: string;
  }

  /** 分析详情 */
  export interface AnalysisDetail {
    id: number;
    topicId: number;
    modelId: number;
    modelName: string;
    inputContent: string;
    aiResult: string;
    userResult: string;
    conclusion: string;
    version: number;
    isCurrent: boolean;
    userId: number;
    status: AnalysisStatus;
    statusText: string;
    createdAt: string;
    updatedAt: string;
  }

  /** 创建分析参数 */
  export interface CreateAnalysisParams {
    topicId: number;
    modelId: number;
    modelName?: string;
    content?: string;
  }

  /** 更新分析参数 */
  export interface UpdateAnalysisParams {
    id: number;
    content: string;
  }

  /** 保存并AI分析参数 */
  export interface SaveWithAiParams {
    id?: number;
    topicId: number;
    modelId: number;
    modelName?: string;
    content: string;
  }

  /** 设为当前版本参数 */
  export interface SetCurrentParams {
    id: number;
  }

  /** 分析列表查询参数 */
  export interface AnalysisListParams {
    page?: number;
    pageSize?: number;
    topicId?: number;
    modelId?: number;
    userId?: number;
    isCurrent?: boolean;
    status?: AnalysisStatus;
  }

  /** 分析列表响应 */
  export interface AnalysisListResult {
    page: number;
    pageSize: number;
    total: number;
    list: AnalysisInfo[];
  }

  /** 分析历史 */
  export interface AnalysisHistory {
    topicId: string;
    modelId: string;
    modelName: string;
    versions: AnalysisInfo[];
  }

  /** 删除分析参数 */
  export interface DelAnalysisParams {
    ids: number[];
  }
}

/**
 * 获取分析列表
 */
export async function getThinkingAnalysisListApi(params: ThinkingAnalysisApi.AnalysisListParams = {}) {
  return requestClient.get<ThinkingAnalysisApi.AnalysisListResult>('/thinking/analysis/list', { params });
}

/**
 * 获取我的分析列表
 */
export async function getMyThinkingAnalysisListApi(params: ThinkingAnalysisApi.AnalysisListParams = {}) {
  return requestClient.get<ThinkingAnalysisApi.AnalysisListResult>('/thinking/analysis/my', { params });
}

/**
 * 获取分析详情
 */
export async function getThinkingAnalysisDetailApi(id: number) {
  return requestClient.get<ThinkingAnalysisApi.AnalysisDetail>(`/thinking/analysis/${id}`);
}

/**
 * 创建分析
 */
export async function createThinkingAnalysisApi(data: ThinkingAnalysisApi.CreateAnalysisParams) {
  return requestClient.post<ThinkingAnalysisApi.AnalysisInfo>('/thinking/analysis', data);
}

/**
 * 更新分析
 */
export async function updateThinkingAnalysisApi(data: ThinkingAnalysisApi.UpdateAnalysisParams) {
  return requestClient.put<ThinkingAnalysisApi.AnalysisInfo>('/thinking/analysis', data);
}

/**
 * 删除分析
 */
export async function deleteThinkingAnalysisApi(data: ThinkingAnalysisApi.DelAnalysisParams) {
  return requestClient.delete('/thinking/analysis', { data });
}

/**
 * 保存并AI分析
 */
export async function saveWithAiApi(data: ThinkingAnalysisApi.SaveWithAiParams) {
  return requestClient.post<ThinkingAnalysisApi.AnalysisInfo>('/thinking/analysis/save-with-ai', data);
}

/**
 * 获取当前版本分析
 */
export async function getCurrentAnalysisApi(topicId: number, modelId: number) {
  return requestClient.get<ThinkingAnalysisApi.AnalysisInfo>('/thinking/analysis/current', {
    params: { topicId, modelId },
  });
}

/**
 * 获取最新分析
 */
export async function getLatestAnalysisApi(topicId: number) {
  return requestClient.get<ThinkingAnalysisApi.AnalysisInfo>('/thinking/analysis/latest', {
    params: { topicId },
  });
}

/**
 * 获取课题的所有分析
 */
export async function getAnalysisByTopicApi(topicId: number) {
  return requestClient.get<ThinkingAnalysisApi.AnalysisInfo[]>(`/thinking/analysis/by-topic/${topicId}`);
}

/**
 * 获取分析历史版本
 */
export async function getAnalysisHistoryApi(topicId: number, modelId: number) {
  return requestClient.get<ThinkingAnalysisApi.AnalysisHistory>(`/thinking/analysis/history/${topicId}/${modelId}`);
}

/**
 * 设为当前版本
 */
export async function setCurrentAnalysisApi(data: ThinkingAnalysisApi.SetCurrentParams) {
  return requestClient.post('/thinking/analysis/set-current', data);
}
