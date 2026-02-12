import { requestClient } from '#/api/request';

export namespace ModelApi {
  /** 模型作者 */
  export interface ModelAuthor {
    id: string;
    name: string;
    avatar: string;
    bio?: string;
  }

  /** 模型统计 */
  export interface ModelStats {
    adoptions: number;
    practices: number;
    discussions: number;
    forks: number;
    likes: number;
  }

  /** 思维模型 */
  export interface ThinkingModel {
    id: string;
    title: string;
    description: string;
    cover: string;
    author: ModelAuthor;
    isFree: boolean;
    price?: number;
    category: string;
    categoryId?: string;
    tags: string[];
    stats: ModelStats;
    updatedAt: string;
    content?: string;
    status?: number;
  }

  /** 模型列表查询参数 */
  export interface ModelListParams {
    page?: number;
    pageSize?: number;
    category?: string;
    keyword?: string;
    isFree?: boolean;
    sortBy?: 'popular' | 'newest' | 'mostAdopted' | 'mostLiked';
  }

  /** 模型列表响应 */
  export interface ModelListResult {
    list: ThinkingModel[];
    total: number;
    page: number;
    pageSize: number;
  }

  /** 创建模型参数 */
  export interface CreateModelParams {
    title: string;
    description: string;
    cover?: string;
    categoryId: string;
    tags?: string[];
    isFree: boolean;
    price?: number;
    content?: string;
  }

  /** 更新模型参数 */
  export interface UpdateModelParams extends Partial<CreateModelParams> {
    id: string;
  }

  /** 模型操作响应 */
  export interface ModelOperationResult {
    success: boolean;
    message?: string;
    data?: ThinkingModel;
  }
}

/**
 * 获取模型列表
 */
export async function getModelListApi(params: ModelApi.ModelListParams = {}) {
  return requestClient.get<ModelApi.ModelListResult>('/market/model/list', {
    params,
  });
}

/**
 * 获取模型详情
 */
export async function getModelDetailApi(id: string) {
  return requestClient.get<ModelApi.ThinkingModel>(`/market/model/${id}`);
}

/**
 * 创建模型（管理员）
 */
export async function createModelApi(data: ModelApi.CreateModelParams) {
  return requestClient.post<ModelApi.ThinkingModel>('/market/model', data);
}

/**
 * 更新模型（管理员）
 */
export async function updateModelApi(data: ModelApi.UpdateModelParams) {
  return requestClient.put<ModelApi.ThinkingModel>(`/market/model/${data.id}`, data);
}

/**
 * 删除模型（管理员）
 */
export async function deleteModelApi(id: string) {
  return requestClient.delete(`/market/model/${id}`);
}

/**
 * 采纳/加载模型到我的模型库
 */
export async function adoptModelApi(id: string) {
  return requestClient.post<ModelApi.ModelOperationResult>(`/market/model/${id}/adopt`);
}

/**
 * 购买模型
 */
export async function purchaseModelApi(id: string) {
  return requestClient.post<ModelApi.ModelOperationResult>(`/market/model/${id}/purchase`);
}

/**
 * 引用/复刻模型
 */
export async function forkModelApi(id: string) {
  return requestClient.post<ModelApi.ThinkingModel>(`/market/model/${id}/fork`);
}

/**
 * 点赞模型
 */
export async function likeModelApi(id: string) {
  return requestClient.post<{ liked: boolean; likes: number }>(`/market/model/${id}/like`);
}

/**
 * 获取推荐模型
 */
export async function getRecommendedModelsApi(category?: string, limit: number = 4) {
  return requestClient.get<ModelApi.ThinkingModel[]>('/market/model/recommended', {
    params: { category, limit },
  });
}
