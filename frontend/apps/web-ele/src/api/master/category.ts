import { requestClient } from '#/api/request';

export namespace MasterCategoryApi {
  /** 分类信息 */
  export interface Category {
    id: string;
    name: string;
    description?: string;
    icon?: string;
    heat: number;
    createdAt?: string;
    updatedAt?: string;
  }

  /** 创建分类参数 */
  export interface CreateCategoryParams {
    name: string;
    description?: string;
    icon?: string;
    heat?: number;
  }

  /** 更新分类参数 */
  export interface UpdateCategoryParams extends Partial<CreateCategoryParams> {
    id: string;
  }
}

/**
 * 获取所有分类（按热度降序）
 */
export async function getAllCategoriesApi() {
  const res = await requestClient.get<{ list: MasterCategoryApi.Category[] }>('/thinkingModel/category/all');
  return res.list;
}

/**
 * 获取分类列表（分页）
 */
export async function getCategoryListApi(params: { page?: number; pageSize?: number }) {
  return requestClient.get<{
    list: MasterCategoryApi.Category[];
    total: number;
    page: number;
    pageSize: number;
  }>('/thinkingModel/category/list', { params });
}

/**
 * 获取分类详情
 */
export async function getCategoryDetailApi(id: string) {
  return requestClient.get<MasterCategoryApi.Category>(`/thinkingModel/category/${id}`);
}

/**
 * 创建分类
 */
export async function createCategoryApi(data: MasterCategoryApi.CreateCategoryParams) {
  return requestClient.post<MasterCategoryApi.Category>('/thinkingModel/category', data);
}

/**
 * 更新分类
 */
export async function updateCategoryApi(data: MasterCategoryApi.UpdateCategoryParams) {
  return requestClient.put<MasterCategoryApi.Category>('/thinkingModel/category', data);
}

/**
 * 删除分类
 */
export async function deleteCategoryApi(id: string) {
  return requestClient.delete('/thinkingModel/category', { params: { ids: [id] } });
}

/**
 * 增加分类热度
 */
export async function increaseHeatApi(id: string, delta: number) {
  return requestClient.post<MasterCategoryApi.Category>('/thinkingModel/category/increaseHeat', {
    id,
    delta,
  });
}
