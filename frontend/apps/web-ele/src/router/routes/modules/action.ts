import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    meta: {
      icon: 'lucide:check-circle',
      order: 3,
      title: '我的行动',
    },
    name: 'MyActions',
    path: '/my-actions',
    component: () => import('#/views/action/index.vue'),
  },
];

export default routes;
