import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    meta: {
      icon: 'lucide:folder-open',
      order: 2,
      title: '我的课题',
    },
    name: 'MyTopics',
    path: '/my-topics',
    component: () => import('#/views/subject/index.vue'),
  },
  // 创建课题页 - 隐藏菜单
  {
    name: 'CreateTopic',
    path: '/my-topics/create',
    component: () => import('#/views/subject/create/index.vue'),
    meta: {
      hideInMenu: true,
      title: '创建课题',
    },
  },
  // 课题详情页 - 隐藏菜单
  {
    name: 'TopicDetail',
    path: '/my-topics/:id',
    component: () => import('#/views/subject/detail/index.vue'),
    meta: {
      hideInMenu: true,
      title: '课题详情',
    },
  },
];

export default routes;
