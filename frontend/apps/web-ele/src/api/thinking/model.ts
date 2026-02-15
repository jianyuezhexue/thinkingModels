import { requestClient } from '#/api/request';

export namespace ThinkingModelApi {
  /** 模型作者信息 */
  export interface ModelAuthor {
    id: string;
    name: string;
    avatar: string;
  }

  /** 模型统计数据 */
  export interface ModelStats {
    usageCount: number;
    adoptCount: number;
    likeCount: number;
    commentCount: number;
  }

  /** 思维模型信息 */
  export interface ModelInfo {
    id: number;
    name: string;
    description: string;
    coverImage: string;
    icon: string;
    categoryId: number;
    price: number;
    isFree: boolean;
    overview: string;
    difficulty: number;
    estimatedTime: number;
    status: number;
    version: string;
    isOfficial: boolean;
    sourceModelId?: number;
    author: ModelAuthor;
    stats: ModelStats;
    tags: string[] | null;
    publishTime?: string;
    createdAt: string;
    updatedAt: string;
  }

  /** 思维模型详情 */
  export interface ModelDetail extends ModelInfo {
    content: string;
  }

  /** 创建模型参数 */
  export interface CreateModelParams {
    name: string;
    description?: string;
    coverImage?: string;
    icon?: string;
    categoryId: number;
    price?: number;
    content?: string;
    overview?: string;
    difficulty?: number;
    estimatedTime?: number;
    version?: string;
    authorName?: string;
    isOfficial?: number;
  }

  /** 更新模型参数 */
  export interface UpdateModelParams {
    id: number;
    name: string;
    description?: string;
    coverImage?: string;
    icon?: string;
    categoryId: number;
    price?: number;
    content?: string;
    overview?: string;
    difficulty?: number;
    estimatedTime?: number;
    version?: string;
    authorName?: string;
  }

  /** 模型列表查询参数 */
  export interface ModelListParams {
    page?: number;
    pageSize?: number;
    name?: string;
    code?: string;
    categoryId?: number;
    status?: number;
    difficulty?: number;
    minPrice?: number;
    maxPrice?: number;
    authorId?: number;
    isOfficial?: number;
    isFree?: boolean;
    sortBy?: string;
  }

  /** 模型列表响应 */
  export interface ModelListResult {
    page: number;
    pageSize: number;
    total: number;
    list: ModelInfo[];
  }

  /** 状态数量统计响应 */
  export interface StatusCountsResult {
    pending: number;
    approved: number;
    rejected: number;
  }

  /** 发布模型参数 */
  export interface PublishModelParams {
    id: number;
  }

  /** Fork模型参数 */
  export interface ForkModelParams {
    sourceModelId: number;
    name: string;
    description?: string;
  }

  /** 删除模型参数 */
  export interface DelModelParams {
    ids: number[];
  }

  /** 审核模型参数 */
  export interface ReviewModelParams {
    id: number;
    approved: boolean;
    note?: string;
  }
}

/**
 * 获取思维模型列表
 */
export async function getThinkingModelListApi(
  params: ThinkingModelApi.ModelListParams = {},
  config?: { signal?: AbortSignal },
) {
  return requestClient.get<ThinkingModelApi.ModelListResult>('/thinkingModel/model/list', {
    params,
    ...config,
  });
}

/**
 * 获取各状态的模型数量统计
 */
export async function getModelStatusCountsApi(config?: { signal?: AbortSignal }) {
  console.log('[getModelStatusCountsApi] Starting request...');
  try {
    const result = await requestClient.get<ThinkingModelApi.StatusCountsResult>(
      '/thinkingModel/model/statusCounts',
      config,
    );
    console.log('[getModelStatusCountsApi] Response:', result);
    return result;
  } catch (error) {
    console.error('[getModelStatusCountsApi] Error:', error);
    throw error;
  }
}

/**
 * 获取我的思维模型列表
 */
export async function getMyThinkingModelListApi(
  params: ThinkingModelApi.ModelListParams = {},
  config?: { signal?: AbortSignal },
) {
  return requestClient.get<ThinkingModelApi.ModelListResult>('/thinkingModel/model/my', {
    params,
    ...config,
  });
}

/**
 * 获取思维模型详情
 */
export async function getThinkingModelDetailApi(id: number) {
  return requestClient.get<ThinkingModelApi.ModelDetail>(`/thinkingModel/model/${id}`);
}

/**
 * 创建思维模型
 */
export async function createThinkingModelApi(data: ThinkingModelApi.CreateModelParams) {
  return requestClient.post<ThinkingModelApi.ModelInfo>('/thinkingModel/model', data);
}

/**
 * 更新思维模型
 */
export async function updateThinkingModelApi(data: ThinkingModelApi.UpdateModelParams) {
  return requestClient.put<ThinkingModelApi.ModelInfo>('/thinkingModel/model', data);
}

/**
 * 删除思维模型
 */
export async function deleteThinkingModelApi(data: ThinkingModelApi.DelModelParams) {
  return requestClient.delete('/thinkingModel/model', { data });
}

/**
 * 发布思维模型
 */
export async function publishThinkingModelApi(data: ThinkingModelApi.PublishModelParams) {
  return requestClient.post<ThinkingModelApi.ModelInfo>('/thinkingModel/model/publish', data);
}

/**
 * 下架思维模型
 */
export async function unpublishThinkingModelApi(id: number) {
  return requestClient.post<ThinkingModelApi.ModelInfo>(`/thinkingModel/model/unpublish/${id}`);
}

/**
 * 引用创建思维模型
 */
export async function forkThinkingModelApi(data: ThinkingModelApi.ForkModelParams) {
  return requestClient.post<ThinkingModelApi.ModelInfo>('/thinkingModel/model/fork', data);
}

/**
 * 审核思维模型
 */
export async function reviewThinkingModelApi(data: ThinkingModelApi.ReviewModelParams) {
  return requestClient.post<ThinkingModelApi.ModelInfo>('/thinkingModel/model/review', data);
}
