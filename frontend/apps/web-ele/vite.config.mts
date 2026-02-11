// @ts-nocheck
import { defineConfig } from '@vben/vite-config';
import ElementPlus from 'unplugin-element-plus/vite';

export default defineConfig(async () => {
  return {
    application: {},
    vite: {
      plugins: [
        ElementPlus({
          format: 'esm',
        }),
      ],
      server: {
        // 前端直接连接后端，不使用代理
      },
    },
  };
});
