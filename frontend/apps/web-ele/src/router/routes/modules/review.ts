import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    meta: {
      icon: 'lucide:clipboard-check',
      order: 6,
      title: '模型审核',
    },
    name: 'ModelReview',
    path: '/review/models',
    component: () => import('#/views/review/models/index.vue'),
  },
];

export default routes;