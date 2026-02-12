<script lang="ts" setup>
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';

import { Page } from '@vben/common-ui';

import {
  ElButton,
  ElCard,
  ElInput,
  ElMessage,
  ElForm,
  ElFormItem,
  ElInputNumber,
  ElTag,
} from 'element-plus';

import { createMeetupApi, type CollisionApi } from '#/api';

import { $t } from '#/locales';

const router = useRouter();

// 表单数据
const formData = ref<CollisionApi.CreateMeetupParams>({
  title: '',
  topic: 'career',
  description: '',
  thoughts: '',
  city: '',
  preferredTime: '',
  costSplit: 'aa',
  maxGuests: 3,
  tags: [],
  modelId: '',
});

// 标签输入
const tagInput = ref('');

// 提交状态
const submitting = ref(false);

// 主题选项
const topicOptions: { id: CollisionApi.MeetupTopic; name: string; icon: string; description: string }[] = [
  { id: 'career', name: '职业发展', icon: '💼', description: '职业规划、转型、成长' },
  { id: 'startup', name: '创业交流', icon: '🚀', description: '创业想法、经验分享' },
  { id: 'technology', name: '技术探讨', icon: '💻', description: '技术趋势、架构设计' },
  { id: 'investment', name: '投资理财', icon: '📈', description: '投资策略、财务规划' },
  { id: 'life', name: '生活感悟', icon: '🌟', description: '人生思考、生活方式' },
  { id: 'other', name: '其他话题', icon: '💭', description: '任何有趣的话题' },
];

// 费用选项
const costSplitOptions: { id: CollisionApi.CostSplit; name: string; icon: string; description: string }[] = [
  { id: 'host', name: '我请客', icon: '☕', description: '展示诚意，吸引更多人参与' },
  { id: 'aa', name: 'AA制', icon: '🤝', description: '公平分摊，轻松交流' },
  { id: 'guest', name: '你请客', icon: '🎁', description: '适合高价值分享者' },
];

// 热门城市
const hotCities = ['北京', '上海', '深圳', '杭州', '广州', '成都', '南京', '武汉'];

// 推荐时间
const suggestedTimes = [
  '周末下午 2-5 点',
  '工作日晚上 7-9 点',
  '周六全天',
  '周日下午',
  '时间灵活可商议',
];

// 添加标签
function addTag() {
  const tag = tagInput.value.trim();
  if (tag && formData.value.tags && formData.value.tags.length < 5 && !formData.value.tags.includes(tag)) {
    formData.value.tags.push(tag);
    tagInput.value = '';
  }
}

// 删除标签
function removeTag(index: number) {
  formData.value.tags?.splice(index, 1);
}

// 选择城市
function selectCity(city: string) {
  formData.value.city = city;
}

// 选择时间
function selectTime(time: string) {
  formData.value.preferredTime = time;
}

// 表单验证
const isFormValid = computed(() => {
  const f = formData.value;
  return f.title.trim().length >= 5 &&
    f.description.trim().length >= 20 &&
    f.thoughts.trim().length >= 50 &&
    f.city.trim() &&
    f.preferredTime.trim();
});

// 提交表单
async function handleSubmit() {
  if (!isFormValid.value) {
    ElMessage.warning('请填写完整信息');
    return;
  }

  submitting.value = true;
  try {
    const meetup = await createMeetupApi(formData.value);
    ElMessage.success('发起约见成功！');
    router.push(`/collision/meetup/${meetup.id}`);
  } catch (error) {
    console.error('发起约见失败:', error);
    ElMessage.error('发起约见失败，请重试');
  } finally {
    submitting.value = false;
  }
}

// 取消
function handleCancel() {
  router.back();
}
</script>

<template>
  <Page
    :title="$t('page.collision.meetup.createTitle')"
    description="发起一个话题，邀约志同道合的人线下交流"
  >
    <div class="max-w-3xl mx-auto">
      <ElCard shadow="hover" class="!rounded-xl">
        <ElForm label-position="top" class="space-y-6">
          <!-- 约见标题 -->
          <ElFormItem label="约见标题" required>
            <ElInput
              v-model="formData.title"
              placeholder="用一句话描述你想聊的话题"
              maxlength="50"
              show-word-limit
              class="!text-lg"
            />
            <div class="text-xs text-gray-400 mt-1">好的标题能吸引志同道合的人</div>
          </ElFormItem>

          <!-- 话题类型 -->
          <ElFormItem label="话题类型" required>
            <div class="grid grid-cols-2 sm:grid-cols-3 gap-3 w-full">
              <div
                v-for="topic in topicOptions"
                :key="topic.id"
                class="p-3 rounded-xl border-2 cursor-pointer transition-all"
                :class="[
                  formData.topic === topic.id
                    ? 'border-purple-500 bg-purple-50'
                    : 'border-gray-200 hover:border-purple-300'
                ]"
                @click="formData.topic = topic.id"
              >
                <div class="flex items-center gap-2 mb-1">
                  <span class="text-xl">{{ topic.icon }}</span>
                  <span class="font-medium">{{ topic.name }}</span>
                </div>
                <div class="text-xs text-gray-500">{{ topic.description }}</div>
              </div>
            </div>
          </ElFormItem>

          <!-- 简短描述 -->
          <ElFormItem label="简短描述" required>
            <ElInput
              v-model="formData.description"
              type="textarea"
              :rows="2"
              placeholder="简要说明你想交流什么（20-100字）"
              maxlength="100"
              show-word-limit
            />
          </ElFormItem>

          <!-- 你的思考 -->
          <ElFormItem required>
            <template #label>
              <div class="flex items-center gap-2">
                <span>你的思考</span>
                <span class="text-purple-600 text-xs bg-purple-100 px-2 py-0.5 rounded-full">核心内容</span>
              </div>
            </template>
            <ElInput
              v-model="formData.thoughts"
              type="textarea"
              :rows="8"
              placeholder="分享你对这个话题的思考和见解...

建议包含：
• 你为什么关注这个话题？
• 你目前的思考是什么？
• 你有什么困惑或想验证的观点？
• 你希望从交流中获得什么？"
              maxlength="2000"
              show-word-limit
            />
            <div class="text-xs text-gray-400 mt-1">
              有深度的思考更容易吸引对的人，至少 50 字
            </div>
          </ElFormItem>

          <!-- 城市选择 -->
          <ElFormItem label="约见城市" required>
            <ElInput
              v-model="formData.city"
              placeholder="输入城市名称"
              class="mb-2"
            />
            <div class="flex flex-wrap gap-2">
              <button
                v-for="city in hotCities"
                :key="city"
                type="button"
                class="px-3 py-1 rounded-full text-sm transition-all"
                :class="[
                  formData.city === city
                    ? 'bg-purple-600 text-white'
                    : 'bg-gray-100 text-gray-600 hover:bg-purple-100'
                ]"
                @click="selectCity(city)"
              >
                {{ city }}
              </button>
            </div>
          </ElFormItem>

          <!-- 时间偏好 -->
          <ElFormItem label="期望时间" required>
            <ElInput
              v-model="formData.preferredTime"
              placeholder="描述你方便的时间段"
              class="mb-2"
            />
            <div class="flex flex-wrap gap-2">
              <button
                v-for="time in suggestedTimes"
                :key="time"
                type="button"
                class="px-3 py-1 rounded-full text-sm transition-all"
                :class="[
                  formData.preferredTime === time
                    ? 'bg-purple-600 text-white'
                    : 'bg-gray-100 text-gray-600 hover:bg-purple-100'
                ]"
                @click="selectTime(time)"
              >
                {{ time }}
              </button>
            </div>
          </ElFormItem>

          <!-- 费用承担 -->
          <ElFormItem label="费用承担" required>
            <div class="grid grid-cols-3 gap-4 w-full">
              <div
                v-for="option in costSplitOptions"
                :key="option.id"
                class="p-4 rounded-xl border-2 cursor-pointer transition-all text-center"
                :class="[
                  formData.costSplit === option.id
                    ? 'border-purple-500 bg-purple-50'
                    : 'border-gray-200 hover:border-purple-300'
                ]"
                @click="formData.costSplit = option.id"
              >
                <div class="text-3xl mb-2">{{ option.icon }}</div>
                <div class="font-medium mb-1">{{ option.name }}</div>
                <div class="text-xs text-gray-500">{{ option.description }}</div>
              </div>
            </div>
          </ElFormItem>

          <!-- 人数限制 -->
          <ElFormItem label="约见人数">
            <div class="flex items-center gap-4">
              <ElInputNumber
                v-model="formData.maxGuests"
                :min="1"
                :max="10"
                controls-position="right"
              />
              <span class="text-gray-500">人（不含你自己）</span>
            </div>
            <div class="text-xs text-gray-400 mt-1">
              推荐 2-4 人，人少交流更深入
            </div>
          </ElFormItem>

          <!-- 标签 -->
          <ElFormItem label="标签（可选）">
            <div class="flex items-center gap-2 mb-2">
              <ElInput
                v-model="tagInput"
                placeholder="添加标签，回车确认"
                class="!w-48"
                @keyup.enter="addTag"
              />
              <ElButton @click="addTag" :disabled="!tagInput.trim()">添加</ElButton>
            </div>
            <div class="flex flex-wrap gap-2">
              <ElTag
                v-for="(tag, index) in formData.tags"
                :key="tag"
                closable
                @close="removeTag(index)"
              >
                {{ tag }}
              </ElTag>
              <span v-if="!formData.tags?.length" class="text-gray-400 text-sm">
                添加标签帮助感兴趣的人找到你
              </span>
            </div>
          </ElFormItem>

          <!-- 提交按钮 -->
          <div class="flex items-center justify-end gap-4 pt-6 border-t">
            <ElButton @click="handleCancel">取消</ElButton>
            <ElButton
              type="primary"
              :loading="submitting"
              :disabled="!isFormValid"
              class="!bg-purple-600 hover:!bg-purple-700 !px-8"
              @click="handleSubmit"
            >
              <span class="mr-2">☕</span>
              发起约见
            </ElButton>
          </div>
        </ElForm>
      </ElCard>

      <!-- 提示信息 -->
      <div class="mt-6 p-4 bg-amber-50 rounded-xl text-sm text-amber-700">
        <div class="font-medium mb-2">💡 发起约见小贴士</div>
        <ul class="space-y-1 text-amber-600">
          <li>• 写清楚你的思考，让别人了解你的深度和诚意</li>
          <li>• 选择合适的费用方式，"我请客"通常更容易吸引参与者</li>
          <li>• 选择公共场所（如连锁咖啡店）作为见面地点</li>
          <li>• 人数不宜太多，2-4人交流效果最好</li>
        </ul>
      </div>
    </div>
  </Page>
</template>

<style scoped>
:deep(.el-form-item__label) {
  font-weight: 600;
  color: #374151;
}

:deep(.el-textarea__inner) {
  font-size: 14px;
  line-height: 1.6;
}
</style>
