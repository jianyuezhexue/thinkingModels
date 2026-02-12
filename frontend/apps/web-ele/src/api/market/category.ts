import { requestClient } from '#/api/request';

export namespace CategoryApi {
  /** 分类信息 */
  export interface Category {
    id: string;
    name: string;
    description?: string;
    icon?: string;
    sort: number;
    status: number;
    modelCount?: number;
    createdAt?: string;
    updatedAt?: string;
  }

  /** 创建分类参数 */
  export interface CreateCategoryParams {
    name: string;
    description?: string;
    icon?: string;
    sort?: number;
  }

  /** 更新分类参数 */
  export interface UpdateCategoryParams extends Partial<CreateCategoryParams> {
    id: string;
  }
}

/**
 * 获取所有分类
 */
export async function getCategoryListApi() {
  return requestClient.get<CategoryApi.Category[]>('/market/category/list');
}

/**
 * 获取分类详情
 */
export async function getCategoryDetailApi(id: string) {
  return requestClient.get<CategoryApi.Category>(`/market/category/${id}`);
}

/**
 * 创建分类（管理员）
 */
export async function createCategoryApi(data: CategoryApi.CreateCategoryParams) {
  return requestClient.post<CategoryApi.Category>('/market/category', data);
}

/**
 * 更新分类（管理员）
 */
export async function updateCategoryApi(data: CategoryApi.UpdateCategoryParams) {
  return requestClient.put<CategoryApi.Category>(`/market/category/${data.id}`, data);
}

/**
 * 删除分类（管理员）
 */
export async function deleteCategoryApi(id: string) {
  return requestClient.delete(`/market/category/${id}`);
}
