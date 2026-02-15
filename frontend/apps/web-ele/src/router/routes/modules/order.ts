import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    meta: {
      icon: 'lucide:receipt',
      order: 4,
      title: '订单管理',
    },
    name: 'OrderManagement',
    path: '/order',
    component: () => import('#/views/order/index.vue'),
  },
];

export default routes;