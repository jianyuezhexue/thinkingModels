import type { RouteRecordRaw } from 'vue-router';

import { $t } from '#/locales';

const routes: RouteRecordRaw[] = [
  {
    meta: {
      icon: 'lucide:brain',
      order: 1,
      title: $t('page.models.menuTitle'),
    },
    name: 'Models',
    path: '/models',
    component: () => import('#/views/models/market/index.vue'),
  },
  // 隐藏路由 - 模型详情页（不在菜单中显示）
  {
    name: 'ModelDetail',
    path: '/models/detail/:id',
    component: () => import('#/views/models/detail/index.vue'),
    meta: {
      hideInMenu: true,
      keepAlive: true,
      title: $t('page.models.detail'),
    },
  },
];

export default routes;
