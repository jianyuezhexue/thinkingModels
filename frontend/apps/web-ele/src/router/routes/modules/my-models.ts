import type { RouteRecordRaw } from 'vue-router';

import { $t } from '#/locales';

const routes: RouteRecordRaw[] = [
  {
    meta: {
      icon: 'lucide:box',
      order: 2,
      title: $t('page.models.myModels'),
    },
    name: 'MyModels',
    path: '/my-models',
    component: () => import('#/views/my-models/index.vue'),
  },
  // 模型详情/编辑页 - 隐藏菜单
  {
    name: 'MyModelDetail',
    path: '/my-models/:id',
    component: () => import('#/views/my-models/detail/index.vue'),
    meta: {
      hideInMenu: true,
      keepAlive: true,
      title: $t('page.models.detail'),
    },
  },
  // 创建模型页面
  {
    name: 'CreateModel',
    path: '/my-models/create',
    component: () => import('#/views/my-models/create/index.vue'),
    meta: {
      hideInMenu: true,
      title: $t('page.models.create'),
    },
  },
];

export default routes;
