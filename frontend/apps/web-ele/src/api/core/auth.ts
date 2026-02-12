import { baseRequestClient, requestClient } from '#/api/request';

export namespace AuthApi {
  /** 登录接口参数 */
  export interface LoginParams {
    password?: string;
    username?: string;
  }

  /** 用户信息 */
  export interface UserInfo {
    id: number;
    username: string;
    nickname: string;
    email: string;
    phone: string;
    avatar: string;
    status: number;
    lastLoginTime: string;
    createdAt: string;
  }

  /** 登录接口返回值 */
  export interface LoginResult {
    accessToken: string;
    refreshToken: string;
    expiresIn: number;
    userInfo: UserInfo;
  }

  export interface RefreshTokenResult {
    data: string;
    status: number;
  }

  /** 注册接口参数 */
  export interface RegisterParams {
    username: string;
    password: string;
    nickname?: string;
    email?: string;
    phone?: string;
  }

  /** 注册接口返回值 */
  export interface RegisterResult {
    id: number;
    username: string;
    nickname: string;
    email: string;
    phone: string;
    status: number;
    createdAt: string;
  }
}

/**
 * 登录
 */
export async function loginApi(data: AuthApi.LoginParams) {
  return requestClient.post<AuthApi.LoginResult>('/auth/login', data);
}

/**
 * 刷新accessToken
 */
export async function refreshTokenApi() {
  return baseRequestClient.post<AuthApi.RefreshTokenResult>('/auth/refresh', {
    withCredentials: true,
  });
}

/**
 * 退出登录
 */
export async function logoutApi() {
  return baseRequestClient.post('/auth/logout', {
    withCredentials: true,
  });
}

/**
 * 获取用户权限码
 */
export async function getAccessCodesApi() {
  return requestClient.get<string[]>('/auth/codes');
}

/**
 * 注册
 */
export async function registerApi(data: AuthApi.RegisterParams) {
  return requestClient.post<AuthApi.RegisterResult>('/auth/register', data);
}
