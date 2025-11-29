<template>
  <div class="app-container">
    <!-- 操作按钮栏 -->
    <div class="toolbar">
      <button @click="handleFormat" class="btn primary">格式化 JSON</button>
      <button @click="handleMinify" class="btn primary">压缩 JSON</button>
      <button @click="handleValidate" class="btn secondary">校验格式</button>
      <button @click="handleCopy" class="btn secondary">复制结果</button>
      <button @click="handleClear" class="btn danger">清空内容</button>
      <span class="status" :class="{ error: statusType === 'error', success: statusType === 'success' }">
        {{ statusText }}
      </span>
    </div>

    <!-- 输入输出区域（响应式布局） -->
    <div class="content">
      <div class="panel">
        <h3>输入 JSON</h3>
        <textarea
          v-model="inputValue"
          class="input-area"
          placeholder="请粘贴或输入 JSON 字符串..."
          @input="handleInput"
        ></textarea>
      </div>
      <div class="panel">
        <h3>处理结果</h3>
        <textarea
          v-model="outputValue"
          class="output-area"
          placeholder="处理后的结果将显示在这里..."
          readonly
        ></textarea>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';

// 状态管理
const inputValue = ref(''); // 输入框内容
const outputValue = ref(''); // 输出框内容
const statusText = ref(''); // 状态提示文本
const statusType = ref(''); // 状态类型（success/error）

// 等待 Wails 初始化完成（确保能调用后端方法）
onMounted(() => {
  window.addEventListener('wails:ready', () => {
    showStatus('就绪 ✅', 'success');
  });
});

// 显示状态提示（3秒后自动清除）
const showStatus = (text, type = 'info') => {
  statusText.value = text;
  statusType.value = type;
  setTimeout(() => {
    statusText.value = '';
    statusType.value = '';
  }, 3000);
};

// 格式化 JSON
const handleFormat = async () => {
  try {
    // 调用后端 FormatJSON 方法（Wails 绑定的 Go 函数）
    const result = await window.go.main.App.FormatJSON(inputValue.value);
    outputValue.value = result;
    showStatus('格式化成功 ✅', 'success');
  } catch (err) {
    showStatus(err.message, 'error');
  }
};

// 压缩 JSON
const handleMinify = async () => {
  try {
    const result = await window.go.main.App.MinifyJSON(inputValue.value);
    outputValue.value = result;
    showStatus('压缩成功 ✅', 'success');
  } catch (err) {
    showStatus(err.message, 'error');
  }
};

// 校验 JSON 格式
const handleValidate = async () => {
  const [isValid, msg] = await window.go.main.App.ValidateJSON(inputValue.value);
  showStatus(msg, isValid ? 'success' : 'error');
};

// 复制结果到剪贴板（使用浏览器原生 API）
const handleCopy = () => {
  if (!outputValue.value.trim()) {
    showStatus('无内容可复制', 'error');
    return;
  }
  navigator.clipboard.writeText(outputValue.value).then(() => {
    showStatus('复制成功 ✅', 'success');
  }).catch(() => {
    showStatus('复制失败', 'error');
  });
};

// 清空输入输出
const handleClear = () => {
  inputValue.value = '';
  outputValue.value = '';
  showStatus('已清空', 'success');
};

// 输入时自动校验（可选优化）
const handleInput = () => {
  if (inputValue.value.trim()) {
    window.go.main.App.ValidateJSON(inputValue.value).then(([isValid, msg]) => {
      statusText.value = msg;
      statusType.value = isValid ? 'success' : 'error';
    });
  } else {
    statusText.value = '';
    statusType.value = '';
  }
};
</script>

<style scoped>
/* 全局容器 */
.app-container {
  width: 100vw;
  height: 100vh;
  padding: 16px;
  box-sizing: border-box;
  font-family: "Microsoft YaHei", Arial, sans-serif;
}

/* 操作栏 */
.toolbar {
  display: flex;
  gap: 12px;
  align-items: center;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

/* 按钮样式 */
.btn {
  padding: 10px 18px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s ease;
}

.btn.primary {
  background-color: #2563eb;
  color: white;
}

.btn.primary:hover {
  background-color: #1d4ed8;
}

.btn.secondary {
  background-color: #64748b;
  color: white;
}

.btn.secondary:hover {
  background-color: #475569;
}

.btn.danger {
  background-color: #ef4444;
  color: white;
}

.btn.danger:hover {
  background-color: #dc2626;
}

/* 状态提示 */
.status {
  margin-left: auto;
  font-size: 14px;
  padding: 8px 12px;
  border-radius: 4px;
}

.status.success {
  color: #10b981;
  background-color: #ecfdf5;
}

.status.error {
  color: #ef4444;
  background-color: #fef2f2;
}

/* 输入输出区域 */
.content {
  display: flex;
  gap: 16px;
  height: calc(100vh - 70px);
  flex-wrap: wrap;
}

.panel {
  flex: 1;
  min-width: 300px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.panel h3 {
  font-size: 16px;
  color: #334155;
  font-weight: 600;
}

/* 文本域样式 */
.input-area, .output-area {
  flex: 1;
  padding: 12px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  resize: none;
  font-size: 14px;
  line-height: 1.6;
  box-sizing: border-box;
}

.input-area {
  background-color: white;
}

.output-area {
  background-color: #f8fafc;
  color: #1e293b;
}

/* 响应式适配（移动端垂直布局） */
@media (max-width: 768px) {
  .content {
    flex-direction: column;
    height: auto;
    min-height: calc(100vh - 70px);
  }

  .panel {
    height: 45vh;
  }
}
</style>