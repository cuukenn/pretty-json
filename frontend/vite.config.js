import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  base: './', // 关键：确保打包后资源路径正确
  server: {
    port: 34115 // Wails 默认前端开发端口
  },
  // 拆包：Vue 运行时与 CodeMirror 分离，避免单个超大 chunk（并行加载 + 缓存友好）
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          vue: ['vue'],
          codemirror: [
            'codemirror',
            '@codemirror/state',
            '@codemirror/view',
            '@codemirror/language',
            '@codemirror/lint',
            '@codemirror/legacy-modes/mode/toml',
            '@lezer/highlight'
          ],
          'codemirror-langs': [
            '@codemirror/lang-json',
            '@codemirror/lang-yaml',
            '@codemirror/lang-go',
            '@codemirror/lang-java',
            '@codemirror/lang-javascript'
          ]
        }
      }
    }
  },
  // 把所有 CodeMirror 相关包锁定进首次预打包：
  // wails dev 代理会阻断 Vite HMR 的 full-reload，若会话中途发现新依赖触发二次预打包，
  // 页面会混用新旧两套 @codemirror/state 副本（Compartment 跨实例失效 → effect.is 报错）
  optimizeDeps: {
    include: [
      'vue',
      'codemirror',
      '@codemirror/state',
      '@codemirror/view',
      '@codemirror/language',
      '@codemirror/lint',
      '@codemirror/lang-json',
      '@codemirror/lang-yaml',
      '@codemirror/lang-go',
      '@codemirror/lang-java',
      '@codemirror/lang-javascript',
      '@codemirror/legacy-modes/mode/toml',
      '@lezer/highlight'
    ]
  }
})