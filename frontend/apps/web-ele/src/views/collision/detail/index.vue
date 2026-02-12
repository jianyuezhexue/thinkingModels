<script lang="ts" setup>
import { computed, onMounted, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { Page } from '@vben/common-ui';

import {
  ElButton,
  ElCard,
  ElAvatar,
  ElTag,
  ElMessage,
  ElSkeleton,
  ElSkeletonItem,
  ElInput,
  ElEmpty,
  ElDivider,
} from 'element-plus';

import {
  getDiscussionDetailApi,
  getCommentListApi,
  createCommentApi,
  likeDiscussionApi,
  favoriteDiscussionApi,
  likeCommentApi,
  type CollisionApi,
} from '#/api';

import { $t } from '#/locales';

const route = useRoute();
const router = useRouter();

// 加载状态
const loading = ref(false);
const commentLoading = ref(false);
const submitting = ref(false);

// 话题详情
const discussion = ref<CollisionApi.Discussion | null>(null);

// 评论列表
const comments = ref<CollisionApi.Comment[]>([]);
const commentTotal = ref(0);
const commentPage = ref(1);
const commentPageSize = ref(10);

// 评论输入
const commentContent = ref('');
const replyTo = ref<CollisionApi.Comment | null>(null);
const replyParentId = ref<string | null>(null);

// 获取话题 ID
const discussionId = computed(() => route.params.id as string);

// 获取话题详情
async function fetchDiscussion() {
  loading.value = true;
  try {
    discussion.value = await getDiscussionDetailApi(discussionId.value);
    if (!discussion.value) {
      ElMessage.error('话题不存在');
      router.push('/collision');
    }
  } catch (error) {
    console.error('获取话题详情失败:', error);
    ElMessage.error('获取话题详情失败');
  } finally {
    loading.value = false;
  }
}

// 获取评论列表
async function fetchComments() {
  commentLoading.value = true;
  try {
    const res = await getCommentListApi({
      discussionId: discussionId.value,
      page: commentPage.value,
      pageSize: commentPageSize.value,
    });
    if (commentPage.value === 1) {
      comments.value = res.list;
    } else {
      comments.value.push(...res.list);
    }
    commentTotal.value = res.total;
  } catch (error) {
    console.error('获取评论列表失败:', error);
  } finally {
    commentLoading.value = false;
  }
}

// 加载更多评论
function loadMoreComments() {
  commentPage.value++;
  fetchComments();
}

// 发表评论
async function submitComment() {
  if (!commentContent.value.trim()) {
    ElMessage.warning('请输入评论内容');
    return;
  }

  submitting.value = true;
  try {
    const newComment = await createCommentApi({
      discussionId: discussionId.value,
      parentId: replyParentId.value || undefined,
      replyToId: replyTo.value?.user.id,
      content: commentContent.value,
    });

    // 如果是回复，添加到对应评论的 replies 中
    if (replyParentId.value) {
      const parentComment = comments.value.find(c => c.id === replyParentId.value);
      if (parentComment) {
        if (!parentComment.replies) parentComment.replies = [];
        parentComment.replies.push(newComment);
        parentComment.replyCount++;
      }
    } else {
      // 否则添加到评论列表顶部
      comments.value.unshift(newComment);
      commentTotal.value++;
    }

    // 更新话题评论数
    if (discussion.value) {
      discussion.value.commentCount++;
    }

    // 清空输入
    commentContent.value = '';
    replyTo.value = null;
    replyParentId.value = null;

    ElMessage.success('评论发表成功');
  } catch (error) {
    console.error('发表评论失败:', error);
    ElMessage.error('发表评论失败');
  } finally {
    submitting.value = false;
  }
}

// 回复评论
function handleReply(comment: CollisionApi.Comment, parentId?: string) {
  replyTo.value = comment;
  replyParentId.value = parentId || comment.id;
  // 滚动到评论框
  document.getElementById('comment-input')?.scrollIntoView({ behavior: 'smooth' });
}

// 取消回复
function cancelReply() {
  replyTo.value = null;
  replyParentId.value = null;
}

// 点赞话题
async function handleLike() {
  if (!discussion.value) return;
  try {
    const res = await likeDiscussionApi(discussion.value.id);
    discussion.value.isLiked = res.liked;
    discussion.value.likeCount = res.likeCount;
  } catch (error) {
    console.error('点赞失败:', error);
  }
}

// 收藏话题
async function handleFavorite() {
  if (!discussion.value) return;
  try {
    const res = await favoriteDiscussionApi(discussion.value.id);
    discussion.value.isFavorited = res.favorited;
    discussion.value.favoriteCount = res.favoriteCount;
  } catch (error) {
    console.error('收藏失败:', error);
  }
}

// 点赞评论
async function handleLikeComment(comment: CollisionApi.Comment) {
  try {
    const res = await likeCommentApi(comment.id);
    comment.isLiked = res.liked;
    comment.likeCount = res.likeCount;
  } catch (error) {
    console.error('点赞评论失败:', error);
  }
}

// 分享
function handleShare() {
  if (navigator.share) {
    navigator.share({
      title: discussion.value?.title,
      text: discussion.value?.summary,
      url: window.location.href,
    });
  } else {
    navigator.clipboard.writeText(window.location.href);
    ElMessage.success('链接已复制到剪贴板');
  }
}

// 返回列表
function goBack() {
  router.push('/collision');
}

// 跳转到模型详情
function goToModel() {
  if (discussion.value?.modelId) {
    router.push(`/market/${discussion.value.modelId}`);
  }
}

// 格式化数字
function formatNumber(num: number): string {
  if (num >= 10000) return (num / 10000).toFixed(1) + '万';
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K';
  return num.toString();
}

// 格式化时间
function formatTime(time: string): string {
  const date = new Date(time);
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
}

// 格式化相对时间
function formatRelativeTime(time: string): string {
  const now = new Date();
  const date = new Date(time);
  const diff = now.getTime() - date.getTime();
  const minutes = Math.floor(diff / 60000);
  const hours = Math.floor(diff / 3600000);
  const days = Math.floor(diff / 86400000);

  if (minutes < 1) return '刚刚';
  if (minutes < 60) return `${minutes}分钟前`;
  if (hours < 24) return `${hours}小时前`;
  if (days < 7) return `${days}天前`;
  return date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' });
}

// 获取分类颜色
function getCategoryColor(category: CollisionApi.Category): string {
  const colors: Record<CollisionApi.Category, string> = {
    inspiration: 'bg-amber-100 text-amber-700',
    methodology: 'bg-blue-100 text-blue-700',
    case: 'bg-green-100 text-green-700',
    question: 'bg-purple-100 text-purple-700',
    share: 'bg-pink-100 text-pink-700',
  };
  return colors[category] || 'bg-gray-100 text-gray-700';
}

// 获取分类名称
function getCategoryName(category: CollisionApi.Category): string {
  const names: Record<CollisionApi.Category, string> = {
    inspiration: '灵感创意',
    methodology: '方法论',
    case: '案例分析',
    question: '问题讨论',
    share: '经验分享',
  };
  return names[category] || category;
}

// 是否还有更多评论
const hasMoreComments = computed(() => comments.value.length < commentTotal.value);

// 页面加载
onMounted(() => {
  fetchDiscussion();
  fetchComments();
});

// 监听路由变化
watch(() => route.params.id, () => {
  commentPage.value = 1;
  fetchDiscussion();
  fetchComments();
});
</script>

<template>
  <Page auto-content-height>
    <!-- 返回按钮 -->
    <div class="mb-4">
      <ElButton link @click="goBack">
        <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"/>
        </svg>
        返回话题列表
      </ElButton>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="max-w-4xl mx-auto">
      <ElCard shadow="never">
        <ElSkeleton animated>
          <template #template>
            <div class="space-y-4">
              <ElSkeletonItem variant="h1" style="width: 70%" />
              <div class="flex items-center gap-3">
                <ElSkeletonItem variant="circle" style="width: 40px; height: 40px" />
                <ElSkeletonItem variant="text" style="width: 150px" />
              </div>
              <ElSkeletonItem variant="text" />
              <ElSkeletonItem variant="text" />
              <ElSkeletonItem variant="text" style="width: 80%" />
            </div>
          </template>
        </ElSkeleton>
      </ElCard>
    </div>

    <!-- 话题内容 -->
    <div v-else-if="discussion" class="max-w-4xl mx-auto">
      <div class="flex gap-6">
        <!-- 主内容 -->
        <div class="flex-1 min-w-0">
          <!-- 话题卡片 -->
          <ElCard shadow="never" class="mb-6">
            <!-- 标记 -->
            <div v-if="discussion.isTop || discussion.isFeatured" class="flex gap-2 mb-4">
              <span v-if="discussion.isTop" class="px-2 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-600">
                置顶
              </span>
              <span v-if="discussion.isFeatured" class="px-2 py-0.5 rounded-full text-xs font-medium bg-purple-100 text-purple-600">
                精选
              </span>
            </div>

            <!-- 标题 -->
            <h1 class="text-2xl font-bold text-gray-900 mb-4">{{ discussion.title }}</h1>

            <!-- 作者信息 -->
            <div class="flex items-center justify-between mb-6 pb-6 border-b border-gray-100">
              <div class="flex items-center gap-3">
                <ElAvatar :src="discussion.user.avatar" :size="44" />
                <div>
                  <div class="font-medium text-gray-800">{{ discussion.user.name }}</div>
                  <div class="text-sm text-gray-400">{{ formatTime(discussion.publishTime) }}</div>
                </div>
              </div>

              <!-- 分类和模型 -->
              <div class="flex items-center gap-2">
                <span :class="['px-3 py-1 rounded-full text-sm font-medium', getCategoryColor(discussion.category)]">
                  {{ getCategoryName(discussion.category) }}
                </span>
                <button
                  v-if="discussion.modelName"
                  class="px-3 py-1 rounded-full text-sm font-medium bg-blue-50 text-blue-600 hover:bg-blue-100 transition-colors cursor-pointer"
                  @click="goToModel"
                >
                  📊 {{ discussion.modelName }}
                </button>
              </div>
            </div>

            <!-- 封面图 -->
            <div v-if="discussion.cover" class="mb-6 rounded-xl overflow-hidden">
              <img :src="discussion.cover" class="w-full max-h-96 object-cover" />
            </div>

            <!-- 内容 -->
            <div class="prose prose-lg max-w-none mb-6 text-gray-700 leading-relaxed whitespace-pre-wrap">
              {{ discussion.content }}
            </div>

            <!-- 标签 -->
            <div class="flex flex-wrap gap-2 mb-6">
              <ElTag v-for="tag in discussion.tags" :key="tag" effect="plain">
                {{ tag }}
              </ElTag>
            </div>

            <!-- 操作栏 -->
            <div class="flex items-center justify-between pt-6 border-t border-gray-100">
              <div class="flex items-center gap-6">
                <!-- 点赞 -->
                <button
                  :class="[
                    'flex items-center gap-2 px-4 py-2 rounded-full transition-all cursor-pointer',
                    discussion.isLiked
                      ? 'bg-red-50 text-red-500'
                      : 'bg-gray-50 text-gray-600 hover:bg-red-50 hover:text-red-500'
                  ]"
                  @click="handleLike"
                >
                  <svg class="w-5 h-5" :fill="discussion.isLiked ? 'currentColor' : 'none'" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"/>
                  </svg>
                  <span class="font-medium">{{ formatNumber(discussion.likeCount) }}</span>
                </button>

                <!-- 收藏 -->
                <button
                  :class="[
                    'flex items-center gap-2 px-4 py-2 rounded-full transition-all cursor-pointer',
                    discussion.isFavorited
                      ? 'bg-yellow-50 text-yellow-500'
                      : 'bg-gray-50 text-gray-600 hover:bg-yellow-50 hover:text-yellow-500'
                  ]"
                  @click="handleFavorite"
                >
                  <svg class="w-5 h-5" :fill="discussion.isFavorited ? 'currentColor' : 'none'" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z"/>
                  </svg>
                  <span class="font-medium">{{ formatNumber(discussion.favoriteCount) }}</span>
                </button>

                <!-- 分享 -->
                <button
                  class="flex items-center gap-2 px-4 py-2 rounded-full bg-gray-50 text-gray-600 hover:bg-purple-50 hover:text-purple-500 transition-all cursor-pointer"
                  @click="handleShare"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z"/>
                  </svg>
                  <span class="font-medium">分享</span>
                </button>
              </div>

              <!-- 统计 -->
              <div class="flex items-center gap-4 text-sm text-gray-400">
                <span class="flex items-center gap-1">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                  </svg>
                  {{ formatNumber(discussion.viewCount) }} 浏览
                </span>
                <span class="flex items-center gap-1">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/>
                  </svg>
                  {{ formatNumber(discussion.commentCount) }} 评论
                </span>
              </div>
            </div>
          </ElCard>

          <!-- 评论区 -->
          <ElCard shadow="never">
            <template #header>
              <div class="flex items-center gap-2">
                <svg class="w-5 h-5 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/>
                </svg>
                <span class="font-medium">评论 ({{ commentTotal }})</span>
              </div>
            </template>

            <!-- 评论输入框 -->
            <div id="comment-input" class="mb-6">
              <!-- 回复提示 -->
              <div v-if="replyTo" class="mb-2 flex items-center gap-2 text-sm text-gray-500">
                <span>回复 @{{ replyTo.user.name }}</span>
                <button
                  class="text-purple-500 hover:text-purple-600 cursor-pointer"
                  @click="cancelReply"
                >
                  取消
                </button>
              </div>

              <div class="flex gap-3">
                <ElAvatar
                  src="https://api.dicebear.com/7.x/avataaars/svg?seed=me"
                  :size="40"
                />
                <div class="flex-1">
                  <ElInput
                    v-model="commentContent"
                    type="textarea"
                    :rows="3"
                    :placeholder="replyTo ? `回复 @${replyTo.user.name}...` : $t('page.collision.comment.placeholder')"
                    resize="none"
                  />
                  <div class="mt-2 flex justify-end">
                    <ElButton
                      type="primary"
                      :loading="submitting"
                      :disabled="!commentContent.trim()"
                      @click="submitComment"
                    >
                      {{ $t('page.collision.comment.submit') }}
                    </ElButton>
                  </div>
                </div>
              </div>
            </div>

            <ElDivider />

            <!-- 评论列表 -->
            <div v-if="commentLoading && commentPage === 1" class="space-y-4">
              <div v-for="i in 3" :key="i" class="flex gap-3">
                <ElSkeleton animated>
                  <template #template>
                    <div class="flex gap-3 w-full">
                      <ElSkeletonItem variant="circle" style="width: 40px; height: 40px; flex-shrink: 0" />
                      <div class="flex-1 space-y-2">
                        <ElSkeletonItem variant="text" style="width: 120px" />
                        <ElSkeletonItem variant="text" />
                        <ElSkeletonItem variant="text" style="width: 80%" />
                      </div>
                    </div>
                  </template>
                </ElSkeleton>
              </div>
            </div>

            <ElEmpty
              v-else-if="comments.length === 0"
              :description="$t('page.collision.comment.noComments')"
            />

            <div v-else class="space-y-6">
              <!-- 评论项 -->
              <div v-for="comment in comments" :key="comment.id" class="group">
                <div class="flex gap-3">
                  <ElAvatar :src="comment.user.avatar" :size="40" class="flex-shrink-0" />
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center gap-2 mb-1">
                      <span class="font-medium text-gray-800">{{ comment.user.name }}</span>
                      <span class="text-sm text-gray-400">{{ formatRelativeTime(comment.createdAt) }}</span>
                    </div>
                    <p class="text-gray-700 mb-2 whitespace-pre-wrap">{{ comment.content }}</p>
                    <div class="flex items-center gap-4">
                      <button
                        :class="[
                          'flex items-center gap-1 text-sm transition-colors cursor-pointer',
                          comment.isLiked ? 'text-red-500' : 'text-gray-400 hover:text-red-500'
                        ]"
                        @click="handleLikeComment(comment)"
                      >
                        <svg class="w-4 h-4" :fill="comment.isLiked ? 'currentColor' : 'none'" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"/>
                        </svg>
                        {{ comment.likeCount || '' }}
                      </button>
                      <button
                        class="flex items-center gap-1 text-sm text-gray-400 hover:text-purple-500 transition-colors cursor-pointer"
                        @click="handleReply(comment)"
                      >
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6"/>
                        </svg>
                        {{ $t('page.collision.comment.reply') }}
                      </button>
                    </div>

                    <!-- 子回复 -->
                    <div v-if="comment.replies && comment.replies.length > 0" class="mt-4 pl-4 border-l-2 border-gray-100 space-y-4">
                      <div v-for="reply in comment.replies" :key="reply.id" class="flex gap-3">
                        <ElAvatar :src="reply.user.avatar" :size="32" class="flex-shrink-0" />
                        <div class="flex-1 min-w-0">
                          <div class="flex items-center gap-2 mb-1">
                            <span class="font-medium text-gray-800 text-sm">{{ reply.user.name }}</span>
                            <span v-if="reply.replyToName" class="text-sm text-gray-400">
                              回复 <span class="text-purple-500">@{{ reply.replyToName }}</span>
                            </span>
                            <span class="text-xs text-gray-400">{{ formatRelativeTime(reply.createdAt) }}</span>
                          </div>
                          <p class="text-gray-700 text-sm mb-2 whitespace-pre-wrap">{{ reply.content }}</p>
                          <div class="flex items-center gap-4">
                            <button
                              :class="[
                                'flex items-center gap-1 text-xs transition-colors cursor-pointer',
                                reply.isLiked ? 'text-red-500' : 'text-gray-400 hover:text-red-500'
                              ]"
                              @click="handleLikeComment(reply)"
                            >
                              <svg class="w-3.5 h-3.5" :fill="reply.isLiked ? 'currentColor' : 'none'" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"/>
                              </svg>
                              {{ reply.likeCount || '' }}
                            </button>
                            <button
                              class="flex items-center gap-1 text-xs text-gray-400 hover:text-purple-500 transition-colors cursor-pointer"
                              @click="handleReply(reply, comment.id)"
                            >
                              <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6"/>
                              </svg>
                              {{ $t('page.collision.comment.reply') }}
                            </button>
                          </div>
                        </div>
                      </div>

                      <!-- 查看更多回复 -->
                      <button
                        v-if="comment.replyCount > (comment.replies?.length || 0)"
                        class="text-sm text-purple-500 hover:text-purple-600 cursor-pointer"
                      >
                        查看更多 {{ comment.replyCount - (comment.replies?.length || 0) }} 条回复
                      </button>
                    </div>
                  </div>
                </div>
              </div>

              <!-- 加载更多 -->
              <div v-if="hasMoreComments" class="text-center">
                <ElButton
                  :loading="commentLoading"
                  @click="loadMoreComments"
                >
                  {{ $t('page.collision.comment.loadMore') }}
                </ElButton>
              </div>
            </div>
          </ElCard>
        </div>

        <!-- 右侧作者卡片 (大屏显示) -->
        <div class="w-72 flex-shrink-0 space-y-6 hidden xl:block">
          <!-- 作者信息 -->
          <ElCard shadow="never">
            <div class="text-center">
              <ElAvatar :src="discussion.user.avatar" :size="80" class="mb-4" />
              <h3 class="text-lg font-semibold text-gray-800 mb-1">{{ discussion.user.name }}</h3>
              <p class="text-sm text-gray-500 mb-4">思维探索爱好者</p>
              <ElButton type="primary" class="w-full">关注</ElButton>
            </div>
          </ElCard>

          <!-- 相关话题 -->
          <ElCard shadow="never">
            <template #header>
              <span class="font-medium">相关话题</span>
            </template>
            <div class="space-y-3">
              <div class="text-center text-sm text-gray-400 py-4">
                暂无相关话题
              </div>
            </div>
          </ElCard>
        </div>
      </div>
    </div>
  </Page>
</template>

<style scoped>
.prose {
  line-height: 1.8;
}

.prose p {
  margin-bottom: 1em;
}
</style>
