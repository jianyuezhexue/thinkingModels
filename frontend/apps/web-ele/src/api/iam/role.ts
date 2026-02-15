import { requestClient } from '#/api/request';

export namespace RoleApi {
  /** 角色信息 */
  export interface Role {
    id: string;
    roleName: string;
    roleCode: string;
    description: string;
    status: number;
    sort: number;
    menuIds: string;
    userCount: number;
    createdAt: string;
    updatedAt: string;
  }

  /** 创建角色参数 */
  export interface CreateRoleParams {
    roleName: string;
    roleCode: string;
    description?: string;
    status?: number;
    sort?: number;
    menuIds?: string;
  }

  /** 更新角色参数 */
  export interface UpdateRoleParams {
    id: string;
    roleName: string;
    description?: string;
    status?: number;
    sort?: number;
    menuIds?: string;
  }

  /** 更新角色权限参数 */
  export interface UpdateRolePermissionParams {
    id: string;
    menuIds: string[];
  }

  /** 搜索角色参数 */
  export interface SearchRoleParams {
    page?: number;
    pageSize?: number;
    roleName?: string;
    roleCode?: string;
    status?: number;
  }

  /** 角色列表响应 */
  export interface ListRoleResponse {
    list: Role[];
    total: number;
    page: number;
    pageSize: number;
  }

  /** 菜单树节点 */
  export interface MenuNode {
    id: string;
    label: string;
    children?: MenuNode[];
  }
}

/**
 * 获取角色列表（分页）
 */
export async function getRoleListApi(params: RoleApi.SearchRoleParams) {
  return requestClient.post<RoleApi.ListRoleResponse>('/role/list', params);
}

/**
 * 获取所有角色（用于下拉选择）
 */
export async function getAllRolesApi() {
  return requestClient.get<RoleApi.Role[]>('/role/all');
}

/**
 * 获取角色详情
 */
export async function getRoleDetailApi(id: string) {
  return requestClient.get<RoleApi.Role>(`/role/${id}`);
}

/**
 * 创建角色
 */
export async function createRoleApi(data: RoleApi.CreateRoleParams) {
  return requestClient.post<RoleApi.Role>('/role', data);
}

/**
 * 更新角色
 */
export async function updateRoleApi(data: RoleApi.UpdateRoleParams) {
  return requestClient.put<RoleApi.Role>('/role', data);
}

/**
 * 删除角色
 */
export async function deleteRoleApi(ids: string[]) {
  return requestClient.delete('/role', { data: { ids } });
}

/**
 * 更新角色权限
 */
export async function updateRolePermissionApi(data: RoleApi.UpdateRolePermissionParams) {
  return requestClient.put<RoleApi.Role>('/role/permission', data);
}