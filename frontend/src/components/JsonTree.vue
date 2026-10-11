<template>
  <div class="json-tree">
    <div class="tree-bar">
      <button class="tree-btn" @click="setAll(true)">展开全部</button>
      <button class="tree-btn" @click="setAll(false)">收起全部</button>
      <span class="tree-meta">{{ statsText }}</span>
    </div>
    <div class="tree-body">
      <JsonTreeRow :node="root" :depth="0" :bus="bus" />
    </div>
  </div>
</template>

<script setup>
import { computed, reactive } from 'vue';
import JsonTreeRow from './JsonTreeRow.vue';

const props = defineProps({
  root: { type: Object, required: true },
});

const bus = reactive({ v: 0, open: true });
function setAll(open) {
  bus.open = open;
  bus.v++;
}

// 轻量统计：节点总数与最大深度（树数据 ≤5MB，一次遍历可接受）
const stats = computed(() => {
  let nodes = 0;
  let maxDepth = 0;
  const walk = (n, d) => {
    nodes++;
    if (d > maxDepth) maxDepth = d;
    if (n.children) for (const c of n.children) walk(c, d + 1);
  };
  walk(props.root, 0);
  return { nodes, maxDepth };
});
const statsText = computed(() => `${stats.value.nodes} 节点 · 深度 ${stats.value.maxDepth + 1}`);
</script>

<style scoped>
.json-tree {
  height: 100%;
  min-height: 0;
  display: flex;
  flex-direction: column;
  border: 1px solid var(--border);
  border-radius: 8px;
  background-color: var(--panel);
  box-shadow: var(--shadow);
  overflow: hidden;
}

.tree-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 10px;
  border-bottom: 1px solid var(--border);
  background-color: var(--panel-alt);
  flex: none;
}
.tree-btn {
  padding: 4px 10px;
  border: 1px solid var(--border);
  border-radius: 6px;
  background-color: var(--panel);
  color: var(--muted);
  font-size: 12px;
  font-family: inherit;
  cursor: pointer;
  transition: all 0.15s ease;
}
.tree-btn:hover { color: var(--accent); border-color: var(--accent); }
.tree-meta {
  margin-left: auto;
  font-size: 12px;
  color: var(--muted);
}

.tree-body {
  flex: 1;
  min-height: 0;
  overflow: auto;
  padding: 8px 4px;
  font-family: 'JetBrains Mono', 'Cascadia Code', Consolas, 'Courier New', monospace;
}
</style>
