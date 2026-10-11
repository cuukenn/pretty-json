<template>
  <div class="statusbar">
    <div class="left">
      <span class="state" :class="stateClass">
        <i class="dot"></i>{{ stateText }}
      </span>
    </div>
    <div class="right">
      <span title="字符数">{{ chars }} 字符</span>
      <span class="sep">·</span>
      <span title="行数">{{ lines }} 行</span>
      <span class="sep">·</span>
      <span title="UTF-8 字节数">{{ bytesText }}</span>
      <template v-if="stats">
        <span class="sep">·</span>
        <span title="对象键总数">键 {{ stats.keys }}</span>
        <span class="sep">·</span>
        <span title="节点总数">节点 {{ stats.nodes }}</span>
        <span class="sep">·</span>
        <span title="最大嵌套深度">深度 {{ stats.depth }}</span>
        <span class="sep">·</span>
        <span title="对象/数组个数">对象 {{ stats.objects }} · 数组 {{ stats.arrays }}</span>
      </template>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  input: { type: String, default: '' },
  result: { type: Object, default: null }, // AnalyzeResult
});

const chars = computed(() => props.input.length);
const lines = computed(() => (props.input ? props.input.split('\n').length : 0));
const bytesText = computed(() => {
  const n = new TextEncoder().encode(props.input).length;
  if (n < 1024) return `${n} B`;
  if (n < 1024 * 1024) return `${(n / 1024).toFixed(1)} KB`;
  return `${(n / 1024 / 1024).toFixed(2)} MB`;
});
const stats = computed(() => props.result?.stats ?? null);

const stateClass = computed(() => {
  if (!props.result || !props.input.trim()) return 'idle';
  return props.result.valid ? 'ok' : 'err';
});
const stateText = computed(() => {
  if (!props.input.trim()) return '就绪';
  if (!props.result) return '校验中…';
  if (props.result.valid) return '合法 JSON';
  if (props.result.line) return `错误：第 ${props.result.line} 行第 ${props.result.col} 列`;
  return `错误：${props.result.message}`;
});
</script>

<style scoped>
.statusbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 30px;
  padding: 0 16px;
  font-size: 12px;
  color: var(--muted);
  background-color: var(--panel);
  border-top: 1px solid var(--border);
  user-select: none;
}

.state {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.state .dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: var(--muted);
}

.state.ok .dot { background-color: var(--success); }
.state.ok { color: var(--success); }
.state.err .dot { background-color: var(--danger); }
.state.err { color: var(--danger); }

.right {
  display: flex;
  align-items: center;
  gap: 6px;
}

.sep {
  color: var(--border);
}
</style>
