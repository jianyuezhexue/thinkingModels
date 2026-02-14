import { ref, computed } from 'vue';

import { defineStore } from 'pinia';

/**
 * 资源数据 Store
 * 集中管理可复用的静态资源（标签、封面图片等）
 */
export const useResourcesStore = defineStore('resources', () => {
  // ==================== 推荐标签 ====================
  /** 推荐标签列表 */
  const suggestedTags = ref<string[]>([
    '战略',
    '分析',
    '思维',
    '创新',
    '管理',
    '决策',
    '效率',
    '逻辑',
    '沟通',
    '规划',
  ]);

  // ==================== 预设封面图片 ====================
  /** 预设封面图片分类 */
  type CoverCategory =
    | 'business'
    | 'strategy'
    | 'analysis'
    | 'decision'
    | 'creative'
    | 'innovation';

  /** 封面图片项 */
  interface PresetCover {
    id: string;
    url: string;
    label: string;
    category: CoverCategory;
  }

  /** 预设封面图片列表 */
  const presetCovers = ref<PresetCover[]>([
    // 商业与战略
    {
      id: '1',
      url: 'https://images.unsplash.com/photo-1454165804606-c3d57bc86b40?w=800&h=400&fit=crop',
      label: '商业会议',
      category: 'business',
    },
    {
      id: '2',
      url: 'https://images.unsplash.com/photo-1552664730-d307ca884978?w=800&h=400&fit=crop',
      label: '团队协作',
      category: 'business',
    },
    {
      id: '3',
      url: 'https://images.unsplash.com/photo-1531403009284-440f080d1e12?w=800&h=400&fit=crop',
      label: '战略规划',
      category: 'strategy',
    },
    {
      id: '4',
      url: 'https://images.unsplash.com/photo-1542744173-8e7e53415bb0?w=800&h=400&fit=crop',
      label: '商务演示',
      category: 'business',
    },
    // 数据与分析
    {
      id: '5',
      url: 'https://images.unsplash.com/photo-1460925895917-afdab827c52f?w=800&h=400&fit=crop',
      label: '数据分析',
      category: 'analysis',
    },
    {
      id: '6',
      url: 'https://images.unsplash.com/photo-1551288049-bebda4e38f71?w=800&h=400&fit=crop',
      label: '图表可视化',
      category: 'analysis',
    },
    {
      id: '7',
      url: 'https://images.unsplash.com/photo-1516321318423-f06f85e504b3?w=800&h=400&fit=crop',
      label: '决策分析',
      category: 'decision',
    },
    {
      id: '8',
      url: 'https://images.unsplash.com/photo-1504868584819-f8e8b4b6d7e3?w=800&h=400&fit=crop',
      label: '数据仪表盘',
      category: 'analysis',
    },
    // 创意与创新
    {
      id: '9',
      url: 'https://images.unsplash.com/photo-1507925921958-8a62f3d1a50d?w=800&h=400&fit=crop',
      label: '创意笔记',
      category: 'creative',
    },
    {
      id: '10',
      url: 'https://images.unsplash.com/photo-1512758017271-d7b84c2113f1?w=800&h=400&fit=crop',
      label: '灵感创意',
      category: 'creative',
    },
    {
      id: '11',
      url: 'https://images.unsplash.com/photo-1517245386807-bb43f82c33c4?w=800&h=400&fit=crop',
      label: '头脑风暴',
      category: 'innovation',
    },
    {
      id: '12',
      url: 'https://images.unsplash.com/photo-1556761175-5973dc0f32e7?w=800&h=400&fit=crop',
      label: '创新团队',
      category: 'innovation',
    },
    // 心理与决策
    {
      id: '13',
      url: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=800&h=400&fit=crop',
      label: '专注思考',
      category: 'decision',
    },
    {
      id: '14',
      url: 'https://images.unsplash.com/photo-1517842645767-c639042777db?w=800&h=400&fit=crop',
      label: '深度阅读',
      category: 'strategy',
    },
    {
      id: '15',
      url: 'https://images.unsplash.com/photo-1434030216411-0b793f4b4173?w=800&h=400&fit=crop',
      label: '学习成长',
      category: 'strategy',
    },
    {
      id: '16',
      url: 'https://images.unsplash.com/photo-1484480974693-6ca0a78fb36b?w=800&h=400&fit=crop',
      label: '目标规划',
      category: 'strategy',
    },
  ]);

  /** 按分类分组的封面图片 */
  const coversByCategory = computed(() => {
    const grouped: Partial<Record<CoverCategory, PresetCover[]>> = {};
    presetCovers.value.forEach((cover) => {
      if (!grouped[cover.category]) {
        grouped[cover.category] = [];
      }
      grouped[cover.category]!.push(cover);
    });
    return grouped;
  });

  /** 根据分类获取封面 */
  function getCoversByCategory(category: CoverCategory): PresetCover[] {
    return presetCovers.value.filter((cover) => cover.category === category);
  }

  /** 随机获取封面 */
  function getRandomCover(): PresetCover | null {
    if (presetCovers.value.length === 0) return null;
    const randomIndex = Math.floor(Math.random() * presetCovers.value.length);
    return presetCovers.value[randomIndex]!;
  }

  return {
    // 推荐标签
    suggestedTags,

    // 预设封面
    presetCovers,
    coversByCategory,
    getCoversByCategory,
    getRandomCover,
  };
});
