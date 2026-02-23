import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  // TODO: MVP 后期添加用户管理功能
  // {
  //   meta: {
  //     icon: 'lucide:users',
  //     order: 98,
  //     title: '用户管理',
  //   },
  //   name: 'UserManagement',
  //   path: '/iam/user',
  //   component: () => import('#/views/iam/user/index.vue'),
  // },
  {
    meta: {
      icon: 'lucide:user-cog',
      order: 99,
      title: '角色管理',
    },
    name: 'RoleManagement',
    path: '/iam/role',
    component: () => import('#/views/iam/role/index.vue'),
  },
];

export default routes;
