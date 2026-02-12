import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    meta: {
      icon: 'lucide:settings',
      order: 99,
      title: '管理后台',
      authority: ['admin'], // 仅管理员可见
    },
    name: 'Admin',
    path: '/admin',
    redirect: '/admin/models',
    children: [
      {
        name: 'ModelManagement',
        path: 'models',
        component: () => import('#/views/admin/models/index.vue'),
        meta: {
          icon: 'lucide:brain',
          title: '模型管理',
          authority: ['admin'],
        },
      },
      {
        name: 'CategoryManagement',
        path: 'categories',
        component: () => import('#/views/admin/categories/index.vue'),
        meta: {
          icon: 'lucide:tags',
          title: '分类管理',
          authority: ['admin'],
        },
      },
    ],
  },
];

export default routes;
