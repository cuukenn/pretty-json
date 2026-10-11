<template>
  <div>
    <div
      class="trow"
      :class="{ clickable: container }"
      :style="{ paddingLeft: depth * 16 + 8 + 'px' }"
      @click="toggle"
    >
      <span v-if="container" class="tchev" :class="{ open: expanded }" aria-hidden="true">
        <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
          <polyline points="9 6 15 12 9 18" />
        </svg>
      </span>
      <span v-else class="tchev-sp" aria-hidden="true"></span>

      <span v-if="keyLabel" class="tkey">{{ keyLabel }}</span>
      <span v-if="keyLabel" class="tpunct">: </span>

      <template v-if="container">
        <span class="tpunct">{{ brackets[0] }}</span>
        <template v-if="!expanded">
          <span class="tpreview">… {{ node.count }} {{ node.kind === 'array' ? '项' : '键' }} </span>
          <span class="tpunct">{{ brackets[1] }}</span>
        </template>
      </template>
      <span v-else :class="'tv-' + node.kind">{{ scalarText }}</span>
    </div>

    <template v-if="container && expanded">
      <JsonTreeRow
        v-for="(child, i) in node.children"
        :key="i"
        :node="child"
        :depth="depth + 1"
        :bus="bus"
        :index="node.kind === 'array' ? i : null"
      />
      <div class="trow" :style="{ paddingLeft: depth * 16 + 8 + 'px' }">
        <span class="tchev-sp" aria-hidden="true"></span>
        <span class="tpunct">{{ brackets[1] }}</span>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue';

const props = defineProps({
  node: { type: Object, required: true },
  depth: { type: Number, required: true },
  bus: { type: Object, required: true }, // { v, open }：展开/收起全部广播
  index: { type: Number, default: null }, // 数组元素序号
});

const expanded = ref(props.depth === 0); // 默认仅展开根节点
watch(() => props.bus.v, () => { expanded.value = props.bus.open; });

const container = computed(() =>
  (props.node.kind === 'object' || props.node.kind === 'array') && props.node.count > 0
);
const brackets = computed(() => (props.node.kind === 'array' ? ['[', ']'] : ['{', '}']));
const keyLabel = computed(() =>
  props.node.key || (props.index !== null ? `[${props.index}]` : '')
);
const scalarText = computed(() => {
  switch (props.node.kind) {
    case 'string': return JSON.stringify(props.node.value);
    case 'null': return 'null';
    default: return props.node.value; // number | bool
  }
});

function toggle() {
  if (container.value) expanded.value = !expanded.value;
}
</script>

<style scoped>
.trow {
  display: flex;
  align-items: baseline;
  gap: 0;
  line-height: 1.7;
  font-size: 13px;
  white-space: pre-wrap;
  word-break: break-all;
  border-radius: 4px;
}
.trow.clickable { cursor: pointer; }
.trow.clickable:hover { background-color: var(--cm-active); }

.tchev {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 14px;
  height: 14px;
  margin-right: 4px;
  color: var(--muted);
  transition: transform 0.15s ease;
  flex: none;
  transform: rotate(0deg);
}
.tchev.open { transform: rotate(90deg); }
.tchev-sp { display: inline-block; width: 18px; flex: none; }

.tkey { color: var(--tok-key); }
.tpunct { color: var(--tok-punct); }
.tpreview { color: var(--muted); font-style: italic; font-size: 12px; }

.tv-string { color: var(--tok-string); word-break: break-all; }
.tv-number { color: var(--tok-number); }
.tv-bool { color: var(--tok-bool); }
.tv-null { color: var(--tok-null); font-style: italic; }
</style>
