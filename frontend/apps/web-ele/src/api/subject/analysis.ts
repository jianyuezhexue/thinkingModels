import { requestClient } from '#/api/request';

export namespace AnalysisApi {
  /** 分析状态 */
  export type AnalysisStatus = 'pending' | 'processing' | 'completed' | 'failed';

  /** 分析输入项 */
  export interface AnalysisInput {
    key: string;
    label: string;
    value: string;
  }

  /** 分析结果项 */
  export interface AnalysisResult {
    key: string;
    label: string;
    content: string;
    suggestions?: string[];
  }

  /** 分析报告 */
  export interface Analysis {
    id: string;
    topicId: string;
    topicTitle?: string;
    modelId: string;
    modelName?: string;
    inputs: AnalysisInput[];
    results?: AnalysisResult[];
    status: AnalysisStatus;
    summary?: string;
    createdAt: string;
    updatedAt: string;
  }

  /** 创建分析参数 */
  export interface CreateAnalysisParams {
    topicId: string;
    modelId: string;
    inputs: AnalysisInput[];
  }

  /** 重新分析参数 */
  export interface ReanalyzeParams {
    analysisId: string;
    inputs?: AnalysisInput[];
  }

  /** 分析模板（模型对应的输入模板） */
  export interface AnalysisTemplate {
    modelId: string;
    modelName: string;
    fields: {
      key: string;
      label: string;
      placeholder?: string;
      required: boolean;
      type: 'text' | 'textarea' | 'select' | 'number';
      options?: { label: string; value: string }[];
    }[];
  }
}

/**
 * 获取分析模板
 */
export async function getAnalysisTemplateApi(modelId: string) {
  return requestClient.get<AnalysisApi.AnalysisTemplate>(`/subject/analysis/template/${modelId}`);
}

/**
 * 创建分析
 */
export async function createAnalysisApi(data: AnalysisApi.CreateAnalysisParams) {
  return requestClient.post<AnalysisApi.Analysis>('/subject/analysis', data);
}

/**
 * 获取分析详情
 */
export async function getAnalysisDetailApi(id: string) {
  return requestClient.get<AnalysisApi.Analysis>(`/subject/analysis/${id}`);
}

/**
 * 获取课题的所有分析
 */
export async function getTopicAnalysesApi(topicId: string) {
  return requestClient.get<AnalysisApi.Analysis[]>(`/subject/analysis/topic/${topicId}`);
}

/**
 * 获取我的分析历史
 */
export async function getMyAnalysesApi(page: number = 1, pageSize: number = 10) {
  return requestClient.get<{
    list: AnalysisApi.Analysis[];
    total: number;
    page: number;
    pageSize: number;
  }>('/subject/analysis/my', {
    params: { page, pageSize },
  });
}

/**
 * 重新分析
 */
export async function reanalyzeApi(data: AnalysisApi.ReanalyzeParams) {
  return requestClient.post<AnalysisApi.Analysis>(`/subject/analysis/${data.analysisId}/reanalyze`, {
    inputs: data.inputs,
  });
}

/**
 * 删除分析
 */
export async function deleteAnalysisApi(id: string) {
  return requestClient.delete(`/subject/analysis/${id}`);
}

/**
 * 导出分析报告
 */
export async function exportAnalysisApi(id: string, format: 'pdf' | 'docx' | 'markdown' = 'pdf') {
  return requestClient.get(`/subject/analysis/${id}/export`, {
    params: { format },
    responseType: 'blob',
  });
}
