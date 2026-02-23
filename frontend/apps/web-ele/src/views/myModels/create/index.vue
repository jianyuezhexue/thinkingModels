<template>
  <Page
    :description="isEdit ? '修改模型内容和设置' : '创建并分享你的思维模型'"
    :title="isEdit ? '编辑模型' : '创建思维模型'"
    content-class="p-6 bg-gray-50"
  >
    <div class="flex gap-6">
      <!-- 左侧主表单 -->
      <div class="flex-1 min-w-0 space-y-6">
        <!-- 步骤导航 -->
        <ElCard shadow="hover" class="!rounded-xl">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-4">
              <button
                v-for="(step, index) in steps"
                :key="step.id"
                class="flex items-center gap-3 px-4 py-3 rounded-xl transition-all"
                :class="[
                  currentStep === index
                    ? 'bg-purple-100 border-2 border-purple-300'
                    : index < currentStep
                    ? 'bg-green-50 border border-green-200'
                    : 'bg-gray-50 border border-gray-200 hover:border-purple-200'
                ]"
                @click="goToStep(index)"
              >
                <div
                  class="w-10 h-10 rounded-full flex items-center justify-center text-lg"
                  :class="[
                    currentStep === index
                      ? 'bg-purple-600 text-white'
                      : index < currentStep
                      ? 'bg-green-500 text-white'
                      : 'bg-gray-200 text-gray-500'
                  ]"
                >
                  <span v-if="index < currentStep">✓</span>
                  <span v-else>{{ step.icon }}</span>
                </div>
                <div class="text-left">
                  <div class="font-semibold" :class="currentStep === index ? 'text-purple-700' : 'text-gray-700'">
                    {{ step.label }}
                  </div>
                  <div class="text-xs text-gray-400">{{ step.description }}</div>
                </div>
              </button>
            </div>
            <div class="text-right">
              <div class="text-sm text-gray-500">完成度</div>
              <div class="text-2xl font-bold text-purple-600">{{ formProgress }}%</div>
            </div>
          </div>
        </ElCard>

        <!-- 步骤 1: 基本信息 -->
        <ElCard v-if="currentStep === 0" shadow="hover" class="!rounded-xl">
          <template #header>
            <div class="flex items-center gap-2">
              <span class="text-lg">📋</span>
              <span class="font-semibold text-gray-700">基本信息</span>
            </div>
          </template>

          <ElForm ref="formRef" :model="form" :rules="rules" label-position="top" class="max-w-2xl">
            <!-- 封面上传 -->
            <ElFormItem label="模型封面">
              <div class="w-full space-y-3">
                <!-- 已选封面预览 -->
                <div
                  v-if="form.cover"
                  class="relative w-full h-48 rounded-xl overflow-hidden group"
                >
                  <img :src="form.cover" class="w-full h-full object-cover" />
                  <div class="absolute inset-0 bg-black/50 flex items-center justify-center gap-3 opacity-0 group-hover:opacity-100 transition-opacity">
                    <button
                      type="button"
                      class="px-4 py-2 bg-white text-gray-700 rounded-full text-sm font-medium hover:bg-purple-50 transition-colors"
                      @click="openCoverDialog"
                    >
                      更换图片
                    </button>
                    <button
                      type="button"
                      class="px-4 py-2 bg-white/20 text-white rounded-full text-sm font-medium hover:bg-white/30 transition-colors"
                      @click="form.cover = ''"
                    >
                      移除
                    </button>
                  </div>
                </div>

                <!-- 未选封面时的选择区域 -->
                <div v-else class="w-full">
                  <!-- 主要区域：选择预设图片 -->
                  <div
                    class="w-full h-44 border-2 border-dashed border-purple-200 rounded-xl flex flex-col items-center justify-center cursor-pointer hover:border-purple-400 hover:bg-purple-50 transition-colors bg-purple-50/30 mb-3"
                    @click="openCoverDialog"
                  >
                    <div class="w-14 h-14 rounded-full bg-purple-100 flex items-center justify-center mb-3">
                      <svg class="h-7 w-7 text-purple-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                      </svg>
                    </div>
                    <p class="text-base text-purple-600 font-medium">点击选择预设图片</p>
                    <p class="text-xs text-purple-400 mt-1">24张高质量免费图片可选</p>
                  </div>

                  <!-- 次要区域：上传自定义图片 -->
                  <!-- <ElUpload
                    class="w-full"
                    :auto-upload="false"
                    :on-change="handleCoverChange"
                    :show-file-list="false"
                    accept="image/*"
                  >
                    <div class="w-full py-3 border border-gray-200 rounded-lg flex items-center justify-center gap-2 cursor-pointer hover:border-purple-300 hover:bg-gray-50 transition-colors">
                      <svg class="h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"/>
                      </svg>
                      <span class="text-sm text-gray-500">或上传自定义图片</span>
                    </div>
                  </ElUpload> -->
                </div>
              </div>
            </ElFormItem>

            <!-- 视频地址 -->
            <ElFormItem label="视频地址">
              <ElInput
                v-model="form.videoUrl"
                placeholder="输入视频链接（可选），支持 YouTube、Bilibili 等平台"
                maxlength="500"
                show-word-limit
                class="!rounded-lg"
              >
                <template #prefix>
                  <svg class="h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"/>
                  </svg>
                </template>
              </ElInput>
              <div class="text-xs text-gray-400 mt-1">添加视频教程链接，帮助用户更好地理解和使用模型</div>
            </ElFormItem>

            <!-- 模型名称 -->
            <ElFormItem label="模型名称" prop="title">
              <ElInput
                v-model="form.title"
                placeholder="给你的思维模型起个名字，例如：SWOT 分析模型"
                maxlength="50"
                show-word-limit
                class="!rounded-lg"
              />
            </ElFormItem>

            <!-- 模型描述 -->
            <ElFormItem label="模型描述" prop="description">
              <ElInput
                v-model="form.description"
                type="textarea"
                :rows="4"
                placeholder="简要描述这个思维模型的用途、适用场景和价值..."
                maxlength="500"
                show-word-limit
              />
            </ElFormItem>

            <!-- 模型分类 -->
            <ElFormItem prop="category">
              <template #label>
                <span>模型分类 <span class="text-xs text-gray-400">(已选择 {{ form.category.length }} 个)</span></span>
              </template>
              <div v-if="categoriesLoading" class="text-center py-4">
                <span class="text-gray-500">加载分类中...</span>
              </div>
              <div v-else-if="categoriesError" class="text-center py-4">
                <span class="text-red-500">{{ categoriesError }}</span>
                <ElButton size="small" @click="loadCategories">重试</ElButton>
              </div>
              <div v-else class="grid grid-cols-4 gap-3">
                <button
                  v-for="cat in categories"
                  :key="cat.value"
                  type="button"
                  class="p-3 rounded-lg border-2 text-center transition-all"
                  :class="[
                    form.category.includes(cat.value)
                      ? 'border-purple-500 bg-purple-50 text-purple-700'
                      : 'border-gray-200 hover:border-purple-300 text-gray-600'
                  ]"
                  @click="toggleCategory(cat.value)"
                >
                  <div class="text-xl mb-1">{{ cat.icon }}</div>
                  <div class="text-sm font-medium">{{ cat.label }}</div>
                </button>
              </div>
            </ElFormItem>

            <!-- 标签 -->
            <ElFormItem label="标签">
              <div class="space-y-3">
                <div class="flex flex-wrap gap-2">
                  <ElTag
                    v-for="tag in form.tags"
                    :key="tag"
                    closable
                    effect="plain"
                    class="!bg-purple-50 !text-purple-600 !border-purple-200 !rounded-full"
                    @close="handleRemoveTag(tag)"
                  >
                    {{ tag }}
                  </ElTag>
                  <ElInput
                    v-if="tagInputVisible"
                    v-model="tagInput"
                    size="small"
                    class="!w-24"
                    @keyup.enter="handleAddTag"
                    @blur="handleAddTag"
                  />
                  <ElButton
                    v-else-if="form.tags.length < 5"
                    size="small"
                    class="!rounded-full"
                    @click="tagInputVisible = true"
                  >
                    + 添加标签
                  </ElButton>
                </div>
                <div class="flex flex-wrap gap-2">
                  <span class="text-xs text-gray-400 mr-2">推荐：</span>
                  <button
                    v-for="tag in suggestedTags.filter(t => !form.tags.includes(t))"
                    :key="tag"
                    type="button"
                    class="px-2 py-0.5 text-xs bg-gray-100 text-gray-500 rounded-full hover:bg-purple-100 hover:text-purple-600 transition-colors"
                    @click="addSuggestedTag(tag)"
                  >
                    + {{ tag }}
                  </button>
                </div>
              </div>
            </ElFormItem>
          </ElForm>
        </ElCard>

        <!-- 步骤 2: 内容编辑 -->
        <template v-if="currentStep === 1">
          <!-- 模型概述 -->
          <ElCard shadow="hover" class="!rounded-xl">
            <template #header>
              <div class="flex items-center gap-2">
                <span class="text-lg">📖</span>
                <span class="font-semibold text-gray-700">模型概述</span>
              </div>
            </template>
            <p class="text-sm text-gray-500 mb-4">介绍这个思维模型的背景、原理和核心价值</p>
            <ElInput
              v-model="form.content.overview"
              type="textarea"
              :rows="6"
              placeholder="详细介绍这个思维模型的背景、理论基础、核心原理和适用场景..."
            />
          </ElCard>

          <!-- 使用步骤 -->
          <ElCard shadow="hover" class="!rounded-xl">
            <template #header>
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <span class="text-lg">📝</span>
                  <span class="font-semibold text-gray-700">使用步骤</span>
                </div>
                <ElButton type="primary" plain size="small" class="!rounded-full" @click="addStep">
                  + 添加步骤
                </ElButton>
              </div>
            </template>
            <p class="text-sm text-gray-500 mb-4">按顺序列出使用这个模型的详细步骤</p>
            <div class="space-y-4">
              <div
                v-for="(step, index) in form.content.steps"
                :key="index"
                class="p-4 rounded-xl bg-gradient-to-r from-purple-50 to-indigo-50 border border-purple-100"
              >
                <div class="flex items-start gap-4">
                  <div class="w-10 h-10 rounded-full bg-purple-600 text-white flex items-center justify-center font-bold flex-shrink-0">
                    {{ index + 1 }}
                  </div>
                  <div class="flex-1 space-y-3">
                    <ElInput
                      v-model="step.title"
                      placeholder="步骤标题，例如：识别优势 (Strengths)"
                    />
                    <ElInput
                      v-model="step.description"
                      type="textarea"
                      :rows="3"
                      placeholder="详细说明这个步骤的操作方法和注意事项..."
                    />
                  </div>
                  <ElButton
                    v-if="form.content.steps.length > 1"
                    type="danger"
                    plain
                    size="small"
                    class="!rounded-full"
                    @click="removeStep(index)"
                  >
                    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                    </svg>
                  </ElButton>
                </div>
              </div>
            </div>
          </ElCard>

          <!-- 实践案例 -->
          <ElCard shadow="hover" class="!rounded-xl">
            <template #header>
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <span class="text-lg">💡</span>
                  <span class="font-semibold text-gray-700">实践案例</span>
                </div>
                <ElButton type="primary" plain size="small" class="!rounded-full" @click="addExample">
                  + 添加案例
                </ElButton>
              </div>
            </template>
            <p class="text-sm text-gray-500 mb-4">提供真实或假设的应用案例，帮助用户理解</p>
            <div class="space-y-4">
              <div
                v-for="(example, index) in form.content.examples"
                :key="index"
                class="p-4 rounded-xl bg-gradient-to-r from-amber-50 to-orange-50 border border-amber-100"
              >
                <div class="flex items-start gap-4">
                  <div class="w-10 h-10 rounded-full bg-amber-500 text-white flex items-center justify-center font-bold flex-shrink-0">
                    {{ index + 1 }}
                  </div>
                  <div class="flex-1 space-y-3">
                    <ElInput
                      v-model="example.title"
                      placeholder="案例标题，例如：某电商平台的 SWOT 分析"
                    />
                    <ElInput
                      v-model="example.content"
                      type="textarea"
                      :rows="4"
                      placeholder="详细描述这个案例的背景、分析过程和结论..."
                    />
                  </div>
                  <ElButton
                    v-if="form.content.examples.length > 1"
                    type="danger"
                    plain
                    size="small"
                    class="!rounded-full"
                    @click="removeExample(index)"
                  >
                    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                    </svg>
                  </ElButton>
                </div>
              </div>
            </div>
          </ElCard>
        </template>

        <!-- 步骤 3: 发布设置 -->
        <ElCard v-if="currentStep === 2" shadow="hover" class="!rounded-xl">
          <template #header>
            <div class="flex items-center gap-2">
              <span class="text-lg">🚀</span>
              <span class="font-semibold text-gray-700">发布设置</span>
            </div>
          </template>

          <div class="max-w-xl space-y-6">
            <!-- 定价设置 -->
            <div class="p-6 rounded-xl bg-gradient-to-r from-purple-50 to-indigo-50 border border-purple-100">
              <h4 class="font-semibold text-gray-800 mb-4">💰 定价设置</h4>
              <div class="flex items-center gap-4 mb-4">
                <button
                  type="button"
                  class="flex-1 p-4 rounded-xl border-2 transition-all text-center"
                  :class="form.isFree ? 'border-green-500 bg-green-50' : 'border-gray-200 hover:border-green-300'"
                  @click="form.isFree = true"
                >
                  <div class="text-2xl mb-1">🆓</div>
                  <div class="font-semibold" :class="form.isFree ? 'text-green-700' : 'text-gray-600'">免费</div>
                  <div class="text-xs text-gray-400">获得更多曝光</div>
                </button>
                <button
                  type="button"
                  class="flex-1 p-4 rounded-xl border-2 transition-all text-center"
                  :class="!form.isFree ? 'border-purple-500 bg-purple-50' : 'border-gray-200 hover:border-purple-300'"
                  @click="form.isFree = false"
                >
                  <div class="text-2xl mb-1">💎</div>
                  <div class="font-semibold" :class="!form.isFree ? 'text-purple-700' : 'text-gray-600'">付费</div>
                  <div class="text-xs text-gray-400">知识变现</div>
                </button>
              </div>

              <div v-if="!form.isFree" class="space-y-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">设置价格</label>
                  <div class="flex items-center gap-3">
                    <ElInput
                      v-model.number="form.price"
                      type="number"
                      min="1"
                      max="999"
                      class="!w-32"
                    >
                      <template #prefix>¥</template>
                    </ElInput>
                    <div class="flex gap-2">
                      <button
                        v-for="price in [9, 19, 29, 49, 99]"
                        :key="price"
                        type="button"
                        class="px-3 py-1.5 text-sm rounded-full transition-colors"
                        :class="form.price === price ? 'bg-purple-600 text-white' : 'bg-gray-100 text-gray-600 hover:bg-purple-100'"
                        @click="form.price = price"
                      >
                        ¥{{ price }}
                      </button>
                    </div>
                  </div>
                </div>
                <div class="p-3 bg-white rounded-lg text-sm text-gray-500">
                  <div class="flex items-center justify-between mb-1">
                    <span>售价</span>
                    <span>¥{{ form.price }}</span>
                  </div>
                  <div class="flex items-center justify-between mb-1">
                    <span>平台服务费 (20%)</span>
                    <span class="text-red-500">-¥{{ (form.price * 0.2).toFixed(2) }}</span>
                  </div>
                  <div class="flex items-center justify-between pt-2 border-t border-gray-100 font-semibold">
                    <span>预计收入</span>
                    <span class="text-green-600">¥{{ (form.price * 0.8).toFixed(2) }}</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- 发布须知 -->
            <div class="p-6 rounded-xl bg-gradient-to-r from-amber-50 to-orange-50 border border-amber-100">
              <h4 class="font-semibold text-amber-800 mb-3">📋 发布须知</h4>
              <ul class="text-sm text-amber-700 space-y-2">
                <li class="flex items-start gap-2">
                  <span class="text-amber-500 mt-0.5">•</span>
                  <span>模型提交后将进入审核流程，通常在 24 小时内完成</span>
                </li>
                <li class="flex items-start gap-2">
                  <span class="text-amber-500 mt-0.5">•</span>
                  <span>确保内容原创或已获得授权，禁止抄袭</span>
                </li>
                <li class="flex items-start gap-2">
                  <span class="text-amber-500 mt-0.5">•</span>
                  <span>模型一经发布，可被其他用户采纳、练习和评价</span>
                </li>
                <li class="flex items-start gap-2">
                  <span class="text-amber-500 mt-0.5">•</span>
                  <span>付费模型的收入将在用户购买后 T+7 日结算</span>
                </li>
              </ul>
            </div>
          </div>
        </ElCard>

        <!-- 底部操作栏 -->
        <ElCard shadow="hover" class="!rounded-xl">
          <div class="flex items-center justify-between">
            <div>
              <ElButton v-if="currentStep === 0" class="!rounded-full" @click="handleCancel">
                取消
              </ElButton>
              <ElButton v-else class="!rounded-full" @click="prevStep">
                <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
                </svg>
                上一步
              </ElButton>
            </div>
            <div class="flex items-center gap-3">
              <ElButton class="!rounded-full" @click="handleSaveDraft">
                <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-3m-1 4l-3 3m0 0l-3-3m3 3V4"/>
                </svg>
                保存草稿
              </ElButton>
              <ElButton
                v-if="currentStep < 2"
                type="primary"
                class="!bg-purple-600 !border-purple-600 hover:!bg-purple-700 !rounded-full"
                @click="nextStep"
              >
                下一步
                <svg class="h-4 w-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
                </svg>
              </ElButton>
              <ElButton
                v-else
                type="primary"
                class="!bg-purple-600 !border-purple-600 hover:!bg-purple-700 !rounded-full"
                @click="handleSubmit"
                :loading="submitLoading"
              >
                <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"/>
                </svg>
                提交审核
              </ElButton>
            </div>
          </div>
        </ElCard>
      </div>

      <!-- 右侧边栏 -->
      <div class="w-80 flex-shrink-0 space-y-6 hidden lg:block">
        <!-- 预览卡片 -->
        <ElCard shadow="hover" class="!rounded-xl">
          <template #header>
            <div class="flex items-center gap-2">
              <span class="text-lg">👁️</span>
              <span class="font-semibold text-gray-700">模型预览</span>
            </div>
          </template>
          <div class="space-y-4">
            <div class="h-32 rounded-lg overflow-hidden bg-gradient-to-br from-purple-100 to-indigo-100">
              <img
                v-if="form.cover"
                :src="form.cover"
                class="w-full h-full object-cover"
              />
              <div v-else class="w-full h-full flex items-center justify-center text-4xl">
                🖼️
              </div>
            </div>
            <div>
              <h4 class="font-semibold text-gray-800 line-clamp-1">
                {{ form.title || '模型名称' }}
              </h4>
              <p class="text-sm text-gray-500 mt-1 line-clamp-2">
                {{ form.description || '模型描述将显示在这里...' }}
              </p>
            </div>
            <div class="flex flex-wrap gap-1">
              <span
                v-for="tag in form.tags.slice(0, 3)"
                :key="tag"
                class="px-2 py-0.5 text-xs bg-purple-100 text-purple-600 rounded-full"
              >
                {{ tag }}
              </span>
              <span v-if="form.tags.length === 0" class="text-xs text-gray-400">暂无标签</span>
            </div>
            <div class="flex items-center justify-between pt-3 border-t border-gray-100">
              <span class="text-sm text-gray-500">
                {{ form.category.length > 0 ? form.category.map(v => categories.find(c => c.value === v)?.label).filter(Boolean).join('、') : '未选择分类' }}
              </span>
              <span
                :class="[
                  'px-2 py-0.5 text-sm font-bold rounded-full',
                  form.isFree ? 'bg-green-100 text-green-600' : 'bg-purple-100 text-purple-600'
                ]"
              >
                {{ form.isFree ? '免费' : '¥' + form.price }}
              </span>
            </div>
          </div>
        </ElCard>

        <!-- 创作指南 -->
        <ElCard shadow="hover" class="!rounded-xl !bg-gradient-to-br from-purple-50 to-indigo-50 !border-purple-100">
          <template #header>
            <div class="flex items-center gap-2">
              <span class="text-lg">📚</span>
              <span class="font-semibold text-purple-700">创作指南</span>
            </div>
          </template>
          <div class="space-y-4">
            <div class="flex items-start gap-3">
              <div class="w-8 h-8 rounded-full bg-purple-200 text-purple-700 flex items-center justify-center font-bold text-sm flex-shrink-0">1</div>
              <div>
                <div class="font-medium text-gray-700 text-sm">明确用途</div>
                <div class="text-xs text-gray-500">确定模型解决什么问题</div>
              </div>
            </div>
            <div class="flex items-start gap-3">
              <div class="w-8 h-8 rounded-full bg-purple-200 text-purple-700 flex items-center justify-center font-bold text-sm flex-shrink-0">2</div>
              <div>
                <div class="font-medium text-gray-700 text-sm">清晰步骤</div>
                <div class="text-xs text-gray-500">让用户容易上手</div>
              </div>
            </div>
            <div class="flex items-start gap-3">
              <div class="w-8 h-8 rounded-full bg-purple-200 text-purple-700 flex items-center justify-center font-bold text-sm flex-shrink-0">3</div>
              <div>
                <div class="font-medium text-gray-700 text-sm">丰富案例</div>
                <div class="text-xs text-gray-500">通过实例帮助理解</div>
              </div>
            </div>
          </div>
        </ElCard>

        <!-- 小贴士 -->
        <ElCard shadow="hover" class="!rounded-xl !bg-gradient-to-br from-amber-50 to-orange-50 !border-amber-100">
          <template #header>
            <div class="flex items-center gap-2">
              <span class="text-lg">💡</span>
              <span class="font-semibold text-amber-700">小贴士</span>
            </div>
          </template>
          <ul class="text-sm text-amber-800 space-y-2">
            <li class="flex items-start gap-2">
              <span class="text-amber-500">•</span>
              优质封面图能提升50%点击率
            </li>
            <li class="flex items-start gap-2">
              <span class="text-amber-500">•</span>
              3-5个步骤最易被用户接受
            </li>
            <li class="flex items-start gap-2">
              <span class="text-amber-500">•</span>
              真实案例更有说服力
            </li>
            <li class="flex items-start gap-2">
              <span class="text-amber-500">•</span>
              ¥19-49 是最佳定价区间
            </li>
          </ul>
        </ElCard>
      </div>
    </div>

    <!-- 封面图片选择对话框 -->
    <ElDialog
      v-model="coverDialogVisible"
      title="选择预设封面"
      width="900"
      destroy-on-close
      :close-on-click-modal="true"
    >
      <div class="space-y-4">
        <p class="text-sm text-gray-500">精选 24 张高质量图片来自 Unsplash，可免费使用</p>
        <div class="grid grid-cols-4 gap-3 max-h-[480px] overflow-y-auto pr-2">
          <div
            v-for="cover in presetCovers"
            :key="cover.id"
            class="relative h-28 rounded-lg overflow-hidden cursor-pointer group ring-2 ring-transparent hover:ring-purple-400 transition-all"
            :class="{ 'ring-purple-500': form.cover === cover.url }"
            @click="selectPresetCover(cover.url)"
          >
            <img
              :src="cover.url"
              :alt="cover.label"
              class="w-full h-full object-cover"
              loading="lazy"
            />
            <div class="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent flex items-end p-2 opacity-0 group-hover:opacity-100 transition-opacity">
              <span class="text-white text-xs font-medium">{{ cover.label }}</span>
            </div>
            <div
              v-if="form.cover === cover.url"
              class="absolute top-2 right-2 w-5 h-5 rounded-full bg-purple-500 flex items-center justify-center"
            >
              <svg class="w-3 h-3 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
              </svg>
            </div>
          </div>
        </div>
      </div>
    </ElDialog>
  </Page>
</template>

<script lang="ts" setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { Page } from '@vben/common-ui';
import { requestClient } from '#/api/request';

import {
  ElButton,
  ElCard,
  ElInput,
  ElForm,
  ElFormItem,
  ElTag,
  ElMessage,
  // ElUpload,
  ElMessageBox,
  ElDialog,
} from 'element-plus';
import type { FormInstance } from 'element-plus';

// API
import { getAllCategoriesApi, type MasterCategoryApi } from '#/api/master/category';
import { getThinkingModelDetailApi } from '#/api/thinking/model';

// Store
import { useResourcesStore } from '#/store/resources';

// 路由
const route = useRoute();
const router = useRouter();
const editId = computed(() => route.query.id as string | undefined);
const isEdit = computed(() => !!editId.value);

// Store
const resourcesStore = useResourcesStore();

// 当前步骤
const currentStep = ref(0);
const steps = [
  { id: 0, label: '基本信息', icon: '📋', description: '模型名称、分类、简介' },
  { id: 1, label: '内容编辑', icon: '📝', description: '使用步骤和案例' },
  { id: 2, label: '发布设置', icon: '🚀', description: '定价和发布选项' },
];

// 表单数据
const form = reactive({
  title: '',
  description: '',
  category: [] as string[],
  tags: [] as string[],
  cover: '',
  videoUrl: '',
  isFree: true,
  price: 29,
  content: {
    overview: '',
    steps: [
      { title: '', description: '' },
    ],
    examples: [
      { title: '', content: '' },
    ],
  },
});

// 表单引用
const formRef = ref<FormInstance>();

// 表单校验规则
const rules = {
  title: [
    { required: true, message: '请输入模型名称', trigger: 'blur' },
    { min: 2, max: 50, message: '名称长度在 2 到 50 个字符', trigger: 'blur' },
  ],
  description: [
    { required: true, message: '请输入模型描述', trigger: 'blur' },
    { min: 10, max: 500, message: '描述长度在 10 到 500 个字符', trigger: 'blur' },
  ],
  category: [
    { required: true, message: '请至少选择一个模型分类', trigger: 'change' },
    { type: 'array' as const, min: 1, message: '请至少选择一个分类', trigger: 'change' },
  ],
};

// 从 Store 获取静态资源
const suggestedTags = computed(() => resourcesStore.suggestedTags);
const presetCovers = computed(() => resourcesStore.presetCovers);

// 分类数据（从后端直接获取）
const categoriesData = ref<MasterCategoryApi.Category[]>([]);
const categoriesLoading = ref(false);
const categoriesError = ref<string | null>(null);

// 分类选项（转换后的格式）
const categories = computed(() =>
  categoriesData.value.map((cat) => ({
    value: String(cat.id),
    label: cat.name,
    icon: cat.icon || '📦',
  })),
);

// 加载分类列表
async function loadCategories() {
  if (categoriesData.value.length > 0) return;

  categoriesLoading.value = true;
  categoriesError.value = null;

  try {
    const res = await getAllCategoriesApi();
    categoriesData.value = res;
  } catch (error: any) {
    categoriesError.value = error.message || '加载分类失败';
    console.error('加载分类失败:', error);
  } finally {
    categoriesLoading.value = false;
  }
}

// 分类选择切换
function toggleCategory(value: string) {
  const index = form.category.indexOf(value);
  if (index > -1) {
    form.category.splice(index, 1);
  } else {
    form.category.push(value);
  }
}
const submitLoading = ref(false);

// 图片选择对话框
const coverDialogVisible = ref(false);

// 打开图片选择对话框
function openCoverDialog() {
  coverDialogVisible.value = true;
}

// 选择预设封面
function selectPresetCover(url: string) {
  form.cover = url;
  coverDialogVisible.value = false;
  ElMessage.success('封面已选择');
}

// 标签输入
const tagInput = ref('');
const tagInputVisible = ref(false);

// 加载状态
const pageLoading = ref(false);

// 请求取消控制器
let abortController: AbortController | null = null;

// 解析内容
function parseContent(content: string): { overview: string; steps: { title: string; description: string }[]; examples: { title: string; content: string }[] } {
  try {
    const parsed = JSON.parse(content);
    return {
      overview: parsed.overview || '',
      steps: parsed.steps || [],
      examples: parsed.examples || [],
    };
  } catch {
    return {
      overview: '',
      steps: [],
      examples: [],
    };
  }
}

// 加载数据
onMounted(async () => {
  // 加载分类数据
  await loadCategories();

  if (isEdit.value && editId.value) {
    // 加载编辑数据
    pageLoading.value = true;
    abortController = new AbortController();
    try {
      const res = await getThinkingModelDetailApi(Number(editId.value));
      form.title = res.name;
      form.description = res.description || '';
      form.cover = res.coverImage || '';
      form.videoUrl = res.videoUrl || '';
      form.category = res.categoryId ? [String(res.categoryId)] : [];
      form.isFree = res.isFree;
      form.price = res.price || 29;

      // 解析内容
      const content = parseContent(res.content || '');
      form.content.overview = content.overview || res.overview || '';
      form.content.steps = content.steps.length > 0 ? content.steps : [{ title: '', description: '' }];
      form.content.examples = content.examples.length > 0 ? content.examples : [{ title: '', content: '' }];

      // 标签
      if (res.tags && Array.isArray(res.tags)) {
        form.tags = res.tags.slice(0, 5);
      }
    } catch (error) {
      console.error('加载模型数据失败:', error);
      ElMessage.error('加载模型数据失败');
      router.push('/my-models');
    } finally {
      pageLoading.value = false;
    }
  }
});

onUnmounted(() => {
  abortController?.abort();
});

// 标签操作
function handleAddTag() {
  const tag = tagInput.value.trim();
  if (tag && !form.tags.includes(tag) && form.tags.length < 5) {
    form.tags.push(tag);
  }
  tagInput.value = '';
  tagInputVisible.value = false;
}

function handleRemoveTag(tag: string) {
  form.tags = form.tags.filter(t => t !== tag);
}

function addSuggestedTag(tag: string) {
  if (!form.tags.includes(tag) && form.tags.length < 5) {
    form.tags.push(tag);
  }
}

// 步骤操作
function addStep() {
  form.content.steps.push({ title: '', description: '' });
}

function removeStep(index: number) {
  if (form.content.steps.length > 1) {
    form.content.steps.splice(index, 1);
  }
}

// 案例操作
function addExample() {
  form.content.examples.push({ title: '', content: '' });
}

function removeExample(index: number) {
  if (form.content.examples.length > 1) {
    form.content.examples.splice(index, 1);
  }
}

// 步骤导航
async function goToStep(step: number) {
  if (step < currentStep.value) {
    currentStep.value = step;
    return;
  }
  
  // 验证当前步骤
  if (currentStep.value === 0) {
    try {
      await formRef.value?.validate();
    } catch {
      ElMessage.warning('请完善基本信息');
      return;
    }
  }
  
  if (currentStep.value === 1) {
    if (!form.content.overview.trim()) {
      ElMessage.warning('请填写模型概述');
      return;
    }
    const hasValidStep = form.content.steps.some(s => s.title.trim() && s.description.trim());
    if (!hasValidStep) {
      ElMessage.warning('请至少添加一个完整的使用步骤');
      return;
    }
  }
  
  currentStep.value = step;
}

function prevStep() {
  if (currentStep.value > 0) {
    currentStep.value--;
  }
}

async function nextStep() {
  await goToStep(currentStep.value + 1);
}

// 提交
async function handleSubmit() {
  try {
    await ElMessageBox.confirm(
      '提交后模型将进入审核流程，审核通过后将发布到市场。确定提交吗？',
      '提交审核',
      { type: 'info' }
    );

    submitLoading.value = true;

    // 构建提交数据（字段名使用驼峰格式与后端 Go struct 对应）
    const submitData = {
      name: form.title,
      description: form.description,
      categoryId: form.category[0] ? parseInt(form.category[0]) : 1,
      difficulty: 2, // 默认难度：1=简单 2=中等 3=困难
      estimatedTime: 30, // 默认预估时间（分钟）
      coverImage: form.cover,
      videoUrl: form.videoUrl,
      content: JSON.stringify({
        overview: form.content.overview,
        steps: form.content.steps.filter((s: any) => s.title.trim() || s.description.trim()),
        examples: form.content.examples.filter((e: any) => e.title.trim() || e.content.trim()),
        tags: form.tags,
      }),
      overview: form.content.overview,
      price: form.isFree ? 0 : form.price,
    };

    if (isEdit.value && editId.value) {
      // 编辑模式 - 更新模型并提交审核
      await requestClient.put('/thinkingModel/model', { ...submitData, id: Number(editId.value) });
      await requestClient.post('/thinkingModel/model/publish', { id: Number(editId.value) });
      ElMessage.success('模型已更新并提交审核');
    } else {
      // 创建模式 - 创建模型
      const res = await requestClient.post('/thinkingModel/model', submitData);
      // 创建成功后提交审核
      if (res && res.id) {
        await requestClient.post('/thinkingModel/model/publish', { id: res.id });
      }
      ElMessage.success('模型已创建并提交审核');
    }

    router.push('/my-models');
  } catch (error: any) {
    if (error !== 'cancel') {
      // 不是用户取消，显示错误信息
      const errorMsg = error?.response?.data?.msg || error?.response?.data?.message || error?.message || '提交失败，请重试';
      ElMessage.error(errorMsg);
    }
  } finally {
    submitLoading.value = false;
  }
}

// 保存草稿
async function handleSaveDraft() {
  try {
    submitLoading.value = true;

    // 构建草稿数据（字段名使用驼峰格式与后端 Go struct 对应）
    const draftData = {
      name: form.title || '未命名草稿',
      description: form.description,
      categoryId: form.category[0] ? parseInt(form.category[0]) : 1,
      difficulty: 2,
      estimatedTime: 30,
      coverImage: form.cover,
      videoUrl: form.videoUrl,
      content: JSON.stringify({
        overview: form.content.overview,
        steps: form.content.steps.filter((s: any) => s.title.trim() || s.description.trim()),
        examples: form.content.examples.filter((e: any) => e.title.trim() || e.content.trim()),
        tags: form.tags,
      }),
      overview: form.content.overview,
      price: form.isFree ? 0 : form.price,
    };

    if (isEdit.value && editId.value) {
      // 编辑模式 - 更新草稿
      await requestClient.put('/thinkingModel/model', { ...draftData, id: Number(editId.value) });
    } else {
      // 创建模式 - 创建草稿
      await requestClient.post('/thinkingModel/model', draftData);
    }

    ElMessage.success('草稿已保存');
    router.push('/my-models');
  } catch (error: any) {
    const errorMsg = error?.response?.data?.msg || error?.response?.data?.message || error?.message || '保存失败，请重试';
    ElMessage.error(errorMsg);
  } finally {
    submitLoading.value = false;
  }
}

// 取消
function handleCancel() {
  router.back();
}

// 计算完成进度
const formProgress = computed(() => {
  let filled = 0;
  let total = 8;
  
  if (form.title) filled++;
  if (form.description) filled++;
  if (form.category.length > 0) filled++;
  if (form.tags.length > 0) filled++;
  if (form.content.overview) filled++;
  if (form.content.steps.some(s => s.title && s.description)) filled++;
  if (form.content.examples.some(e => e.title && e.content)) filled++;
  if (form.cover) filled++;
  
  return Math.round((filled / total) * 100);
});
</script>

<style scoped>
.line-clamp-1 {
  display: -webkit-box;
  -webkit-line-clamp: 1;
  line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
