import type { RouteRecordRaw } from 'vue-router';

import { $t } from '#/locales';

const routes: RouteRecordRaw[] = [
  {
    meta: {
      icon: 'lucide:store',
      order: 1,
      title: $t('page.models.market'),
    },
    name: 'Market',
    path: '/market',
    component: () => import('#/views/market/index.vue'),
  },
  // 模型详情页 - 隐藏菜单
  {
    name: 'MarketDetail',
    path: '/market/:id',
    component: () => import('#/views/market/detail/index.vue'),
    meta: {
      hideInMenu: true,
      keepAlive: true,
      title: $t('page.models.detail'),
    },
  },
];

export default routes;
