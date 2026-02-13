import type { RouteRecordRaw } from 'vue-router';

// import { $t } from '#/locales';

// MVP阶段暂时隐藏思维碰撞功能，聚焦核心功能
const routes: RouteRecordRaw[] = [
  // {
  //   meta: {
  //     icon: 'lucide:lightbulb',
  //     order: 1.5, // 放在模型市场下面
  //     title: $t('page.collision.title'),
  //   },
  //   name: 'Collision',
  //   path: '/collision',
  //   component: () => import('#/views/collision/index.vue'),
  // },
  // // 创建话题页 - 隐藏菜单 (必须在 :id 路由之前)
  // {
  //   name: 'CreateCollision',
  //   path: '/collision/create',
  //   component: () => import('#/views/collision/create/index.vue'),
  //   meta: {
  //     hideInMenu: true,
  //     title: $t('page.collision.create'),
  //   },
  // },
  // // 约见列表页 - 独立页面
  // {
  //   name: 'MeetupList',
  //   path: '/collision/meetup',
  //   component: () => import('#/views/collision/meetup/index.vue'),
  //   meta: {
  //     hideInMenu: true,
  //     title: $t('page.collision.meetup.title'),
  //   },
  // },
  // // 创建约见页
  // {
  //   name: 'CreateMeetup',
  //   path: '/collision/meetup/create',
  //   component: () => import('#/views/collision/meetup/create.vue'),
  //   meta: {
  //     hideInMenu: true,
  //     title: $t('page.collision.meetup.createTitle'),
  //   },
  // },
  // // 约见详情页
  // {
  //   name: 'MeetupDetail',
  //   path: '/collision/meetup/:id',
  //   component: () => import('#/views/collision/meetup/detail.vue'),
  //   meta: {
  //     hideInMenu: true,
  //     title: $t('page.collision.meetup.detail'),
  //   },
  // },
  // // 创建咨询页
  // {
  //   name: 'CreateConsultation',
  //   path: '/collision/consultation/create',
  //   component: () => import('#/views/collision/consultation/create.vue'),
  //   meta: {
  //     hideInMenu: true,
  //     title: $t('page.collision.consultation.createTitle'),
  //   },
  // },
  // // 咨询详情页
  // {
  //   name: 'ConsultationDetail',
  //   path: '/collision/consultation/:id',
  //   component: () => import('#/views/collision/consultation/detail.vue'),
  //   meta: {
  //     hideInMenu: true,
  //     title: $t('page.collision.consultation.detail'),
  //   },
  // },
  // // 话题详情页 - 隐藏菜单 (放在最后)
  // {
  //   name: 'CollisionDetail',
  //   path: '/collision/:id',
  //   component: () => import('#/views/collision/detail/index.vue'),
  //   meta: {
  //     hideInMenu: true,
  //     keepAlive: true,
  //     title: $t('page.collision.detail'),
  //   },
  // },
];

export default routes;
