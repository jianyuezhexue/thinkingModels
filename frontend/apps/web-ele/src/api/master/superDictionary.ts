import { requestClient } from '#/api/request';

export namespace SuperDictionaryApi {
  /** 超级字典实体 */
  export interface SuperDictionary {
    id: number;
    parentId: number;
    dictValue: string;
    dictName: string;
    level: number;
    levelName: string;
    description: string;
    eval: string;
    extSchema: string;
    extJson: string;
    createdAt: string;
    updatedAt: string;
  }

  /** 树形节点 */
  export interface TreeNode extends SuperDictionary {
    children: TreeNode[];
  }

  /** 创建参数 */
  export interface CreateParams {
    parentId?: number;
    dictValue: string;
    dictName: string;
    level?: number;
    levelName?: string;
    description?: string;
    eval?: string;
    extSchema?: string;
    extJson?: string;
  }

  /** 更新参数 */
  export interface UpdateParams extends Partial<CreateParams> {
    id: number;
  }

  /** 列表查询参数 */
  export interface ListParams {
    page?: number;
    pageSize?: number;
    id?: number;
    ids?: number[];
    parentId?: number;
    dictValue?: string;
    dictName?: string;
    level?: number;
    levelName?: string;
    description?: string;
  }

  /** 列表结果 */
  export interface ListResult {
    page: number;
    pageSize: number;
    total: number;
    list: SuperDictionary[];
  }
}

/**
 * 获取字典列表（分页）
 */
export async function getSuperDictionaryListApi(params: SuperDictionaryApi.ListParams = {}) {
  return requestClient.post<SuperDictionaryApi.ListResult>('/master/superDictionary/list', params);
}

/**
 * 获取字典详情
 */
export async function getSuperDictionaryDetailApi(id: number) {
  return requestClient.post<SuperDictionaryApi.SuperDictionary>(`/master/superDictionary/${id}`);
}

/**
 * 创建字典
 */
export async function createSuperDictionaryApi(data: SuperDictionaryApi.CreateParams) {
  return requestClient.post<SuperDictionaryApi.SuperDictionary>('/master/superDictionary', data);
}

/**
 * 更新字典
 */
export async function updateSuperDictionaryApi(data: SuperDictionaryApi.UpdateParams) {
  return requestClient.put<SuperDictionaryApi.SuperDictionary>('/master/superDictionary', data);
}

/**
 * 批量删除字典
 */
export async function deleteSuperDictionaryApi(ids: number[]) {
  return requestClient.delete('/master/superDictionary', { data: { ids } });
}

/**
 * 获取树形结构
 */
export async function getSuperDictionaryTreeApi(parentId: number = 0) {
  return requestClient.get<SuperDictionaryApi.TreeNode[]>('/master/superDictionary/tree', {
    params: { parentId },
  });
}

/**
 * 获取子节点列表
 */
export async function getSuperDictionaryChildrenApi(parentId: number) {
  return requestClient.post<SuperDictionaryApi.SuperDictionary[]>('/master/superDictionary/children', {
    parentId,
  });
}