# Frontend Code Standards

## Project Overview

This is a Vue 3 + Element Plus + VBEN Admin frontend project using pnpm monorepo architecture.

## Tech Stack

- **Framework**: Vue 3 (Composition API with `<script setup>`)
- **UI Library**: Element Plus
- **State Management**: Pinia
- **Router**: Vue Router 4
- **Build Tool**: Vite
- **Package Manager**: pnpm
- **Styling**: Tailwind CSS + PostCSS
- **Language**: TypeScript

## Directory Structure

```
frontend/
├── apps/
│   └── web-ele/              # Main application
│       ├── src/
│       │   ├── api/          # API layer
│       │   │   ├── core/     # Core APIs (auth, user, menu)
│       │   │   ├── market/   # Market-related APIs
│       │   │   └── subject/  # Subject/topic APIs
│       │   ├── views/        # Page views
│       │   │   ├── market/   # Market module pages
│       │   │   ├── subject/  # Subject module pages
│       │   │   ├── my-models/# My models pages
│       │   │   └── _core/    # Core pages (login, errors)
│       │   ├── router/       # Router configuration
│       │   │   └── routes/
│       │   │       └── modules/  # Route modules (auto-loaded)
│       │   ├── locales/      # i18n translations
│       │   │   └── langs/
│       │   │       ├── zh-CN/
│       │   │       └── en-US/
│       │   ├── layouts/      # Layout components
│       │   ├── adapter/      # Component library adapter
│       │   └── store/        # Pinia stores
├── packages/                 # Shared packages
├── internal/                 # Internal tools
└── docs/                     # Documentation
```

## Naming Conventions

### Files

| Type | Pattern | Example |
|------|---------|---------|
| Views | `PascalCase/index.vue` | `market/index.vue` |
| Components | `PascalCase.vue` | `ModelCard.vue` |
| Composables | `camelCase.ts` | `useModelList.ts` |
| Utils | `camelCase.ts` | `formatDate.ts` |
| APIs | `camelCase.ts` | `model.ts` |
| Routes | `kebab-case.ts` | `my-models.ts` |

### Code

| Type | Pattern | Example |
|------|---------|---------|
| Components | `PascalCase` | `ModelList` |
| Props | `camelCase` | `modelId` |
| Events | `camelCase` | `onUpdate` |
| Functions | `camelCase` | `fetchModels()` |
| Constants | `UPPER_SNAKE_CASE` | `MAX_PAGE_SIZE` |
| Types/Interfaces | `PascalCase` | `ModelApi` |
| Enums | `PascalCase` | `ModelStatus` |

## Vue Component Structure

```vue
<script lang="ts" setup>
// 1. Imports (grouped: Vue, 3rd-party, local)
import { computed, onMounted, ref, watch } from 'vue';
import { useRouter } from 'vue-router';

import { Page } from '@vben/common-ui';
import { ElButton, ElCard } from 'element-plus';

import { getModelListApi } from '#/api';

// 2. Types/Interfaces
interface Model {
  id: string;
  title: string;
}

// 3. Props & Emits
const props = defineProps<{
  modelId: string;
}>();

const emit = defineEmits<{
  (e: 'update', id: string): void;
}>();

// 4. Reactive State
const loading = ref(false);
const models = ref<Model[]>([]);

// 5. Computed
const totalCount = computed(() => models.value.length);

// 6. Methods
async function fetchModels() {
  loading.value = true;
  try {
    const res = await getModelListApi();
    models.value = res.list;
  } finally {
    loading.value = false;
  }
}

// 7. Lifecycle
onMounted(() => {
  fetchModels();
});
</script>

<template>
  <Page title="Page Title" description="Page description">
    <ElCard>
      <!-- Content -->
    </ElCard>
  </Page>
</template>

<style scoped>
/* Component-specific styles */
</style>
```

## API Layer Standards

### File Structure

```typescript
// api/market/model.ts
import { requestClient } from '#/api/request';

// Namespace for types
export namespace ModelApi {
  export interface ThinkingModel {
    id: string;
    title: string;
    description: string;
  }

  export interface ModelListParams {
    page?: number;
    pageSize?: number;
    keyword?: string;
    category?: string;
  }
}

// API functions
export async function getModelListApi(params?: ModelApi.ModelListParams) {
  return requestClient.get<{
    list: ModelApi.ThinkingModel[];
    total: number;
  }>('/market/models', { params });
}

export async function getModelDetailApi(id: string) {
  return requestClient.get<ModelApi.ThinkingModel>(`/market/models/${id}`);
}

export async function createModelApi(data: Partial<ModelApi.ThinkingModel>) {
  return requestClient.post<ModelApi.ThinkingModel>('/market/models', data);
}
```

### Export Pattern

```typescript
// api/index.ts
export * from './market/model';
export * from './market/category';
export * from './subject/topic';
```

## Route Configuration

### Module Routes (Auto-loaded)

```typescript
// router/routes/modules/my-models.ts
import type { RouteRecordRaw } from 'vue-router';
import { $t } from '#/locales';

const routes: RouteRecordRaw[] = [
  {
    meta: {
      icon: 'lucide:box',        // Icon from lucide
      order: 2,                   // Menu order
      title: $t('page.models.myModels'), // i18n key
    },
    name: 'MyModels',
    path: '/my-models',
    component: () => import('#/views/my-models/index.vue'),
  },
  // Hidden child routes
  {
    name: 'MyModelDetail',
    path: '/my-models/:id',
    component: () => import('#/views/my-models/detail/index.vue'),
    meta: {
      hideInMenu: true,
      title: $t('page.models.detail'),
    },
  },
];

export default routes;
```

## Internationalization (i18n)

### File Structure

```
locales/langs/
├── zh-CN/
│   ├── page.json     # Page titles, menu items
│   └── demos.json    # Demo content
└── en-US/
    ├── page.json
    └── demos.json
```

### Usage Pattern

```typescript
// In routes
import { $t } from '#/locales';

meta: {
  title: $t('page.models.market'),
}
```

```vue
<!-- In components (if needed) -->
<template>
  <span>{{ $t('page.models.create') }}</span>
</template>
```

## Styling Standards

### Tailwind CSS First

```vue
<template>
  <!-- Good: Use Tailwind utility classes -->
  <div class="flex items-center gap-4 p-4 bg-white rounded-lg shadow-sm">
    <span class="text-sm font-medium text-gray-900">Title</span>
  </div>

  <!-- Avoid: Custom CSS when Tailwind can do it -->
  <div class="custom-card">...</div>
</template>
```

### Element Plus Component Styling

```vue
<style scoped>
/* Override Element Plus variables if needed */
:deep(.el-button--primary) {
  --el-button-bg-color: var(--custom-primary);
}

/* Component-specific custom styles */
.custom-class {
  /* Only when Tailwind can't handle it */
}
</style>
```

### Common Utility Classes

```css
/* Use these standard utilities */
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
```

## TypeScript Standards

### Type Definitions

```typescript
// Prefer interfaces for object shapes
interface Model {
  id: string;
  title: string;
}

// Use type for unions/intersections
type Status = 'draft' | 'published' | 'archived';

// Export namespace for API types
export namespace ModelApi {
  export interface CreateParams {
    title: string;
  }
}
```

### Props Typing

```vue
<script setup lang="ts">
// Good: Explicit type annotation
const props = defineProps<{
  modelId: string;
  model?: Model;
  loading?: boolean;
}>();

// Provide defaults if needed
const props = withDefaults(defineProps<{
  loading?: boolean;
}>(), {
  loading: false,
});
</script>
```

## Error Handling

```typescript
// API calls with error handling
async function fetchData() {
  loading.value = true;
  try {
    const res = await getDataApi();
    data.value = res;
  } catch (error) {
    console.error('Failed to fetch data:', error);
    ElMessage.error('获取数据失败');
  } finally {
    loading.value = false;
  }
}
```

## Best Practices

1. **Composition API**: Always use `<script setup>` syntax
2. **Import Order**: Vue core → 3rd party → @vben packages → local (#/)
3. **Reactivity**: Use `ref` for primitives, `reactive` for objects
4. **API Calls**: Centralize in `api/` directory, use namespaces
5. **Types**: Always define types, avoid `any`
6. **Components**: Use Element Plus components, avoid custom when possible
7. **Icons**: Use `lucide` icons (via `icon: 'lucide:name'`)
8. **Comments**: Minimal comments, code should be self-documenting
9. **Performance**: Use `computed` for derived state, `watch` for side effects
10. **Accessibility**: Use semantic HTML, proper ARIA labels

## Adding a New Module

1. **Create route file**: `router/routes/modules/{module}.ts`
2. **Add translations**: Update `locales/langs/*/page.json`
3. **Create views**: `views/{module}/index.vue`
4. **Create APIs**: `api/{module}/{entity}.ts`
5. **Export APIs**: Add to `api/index.ts`

## Code Review Checklist

- [ ] TypeScript types are properly defined
- [ ] No `console.log` in production code
- [ ] API errors are handled with user-friendly messages
- [ ] i18n keys are used for all user-facing text
- [ ] Tailwind classes are used for styling (no custom CSS unless necessary)
- [ ] Component follows the standard structure
- [ ] Props are properly typed with defaults if optional
- [ ] File and function names follow naming conventions
