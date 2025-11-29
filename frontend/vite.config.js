import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  base: './', // 关键：确保打包后资源路径正确
  server: {
    port: 34115 // Wails 默认前端开发端口
  }
})