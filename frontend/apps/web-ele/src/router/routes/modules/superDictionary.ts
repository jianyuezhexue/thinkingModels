import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    meta: {
      icon: 'lucide:book-key',
      order: 5,
      title: '数据字典',
    },
    name: 'SuperDictionary',
    path: '/superDictionary',
    component: () => import('#/views/admin/superDictionary/index.vue'),
  },
];

export default routes;