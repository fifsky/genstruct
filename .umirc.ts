import { defineConfig } from 'umi';

export default defineConfig({
  routes: [{ path: '/', component: '@/pages/index' }],
  proxy: {
    '/api': {
      target: 'http://127.0.0.1:8080',
      changeOrigin: true,
    },
  },
});
