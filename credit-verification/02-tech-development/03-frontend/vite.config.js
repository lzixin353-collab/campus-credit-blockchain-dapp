import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import path from 'path'

export default defineConfig({
  plugins: [
    vue(),
    // Element Plus自动导入（适配2.3.8）
    AutoImport({
      resolvers: [ElementPlusResolver()],
    }),
    Components({
      resolvers: [ElementPlusResolver()],
    }),
  ],
  // 路径别名（避免导入路径错误）
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src'),
    },
  },
  // 开发服务器（前端 8081，/api 代理到后端 8080，保留 /api 前缀与后端路由一致）
  server: {
    port: 8081,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
    open: true,
    // 关闭HMR覆盖层（可选，避免报错干扰）
    hmr: {
      overlay: false
    }
  },
  // 版本兼容：锁定依赖构建版本
  optimizeDeps: {
    include: ['vue', 'element-plus', 'web3@4.1.2', 'axios@1.4.0']
  }
})