<template>
  <div class="toolstrip">
    <div class="tabs">
      <button
        v-for="t in tabs"
        :key="t.id"
        class="tab"
        :class="{ active: active === t.id }"
        @click="active = t.id"
      >
        {{ t.label }}
      </button>
    </div>

    <div class="controls">
      <!-- 格式互转 -->
      <template v-if="active === 'convert'">
        <button class="tbtn" :disabled="busy" @click="convert('JSONToYAML', 'yaml', 'JSON → YAML 完成 ✅')">JSON → YAML</button>
        <button class="tbtn" :disabled="busy" @click="convert('YAMLToJSON', 'json', 'YAML → JSON 完成 ✅')">YAML → JSON</button>
        <button class="tbtn" :disabled="busy" @click="convert('JSONToTOML', 'toml', 'JSON → TOML 完成 ✅')">JSON → TOML</button>
        <button class="tbtn" :disabled="busy" @click="convert('TOMLToJSON', 'json', 'TOML → JSON 完成 ✅')">TOML → JSON</button>
      </template>

      <!-- JSONPath -->
      <template v-else-if="active === 'jsonpath'">
        <input
          v-model="path"
          class="path-input"
          placeholder="JSONPath，如 $.store.book[*].title"
          spellcheck="false"
          @keyup.enter="queryPath"
        />
        <button class="tbtn primary" :disabled="busy" @click="queryPath">查询</button>
        <span class="tp-hint">支持 .name · ..name · ['k'] · [n] · [*] · [a:b:step]</span>
      </template>

      <!-- 代码生成 -->
      <template v-else>
        <select v-model="codeLang" class="select">
          <option value="go">Go</option>
          <option value="typescript">TypeScript</option>
          <option value="java">Java</option>
        </select>
        <label v-if="codeLang === 'java'" class="chk" title="用 Lombok @Data 注解替代显式 getter/setter">
          <input type="checkbox" v-model="javaLombok" />
          Lombok @Data
        </label>
        <button class="tbtn primary" :disabled="busy" @click="genCode">生成类型定义</button>
        <span class="tp-hint">自动推断结构体 / 接口 / 类</span>
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';

const props = defineProps({
  input: { type: String, default: '' },
});
const emit = defineEmits(['result', 'error']);

const tabs = [
  { id: 'convert', label: '格式互转' },
  { id: 'jsonpath', label: 'JSONPath' },
  { id: 'codegen', label: '代码生成' },
];
const active = ref('convert');
const busy = ref(false);
const path = ref('');
const codeLang = ref('go');
const javaLombok = ref(false);

function api() {
  const a = window.go?.main?.App;
  if (!a) {
    emit('error', '后端未就绪');
    return null;
  }
  return a;
}

// 通用调用：成功 emit result({text, lang, note})，失败 emit error(message)
async function call(fn, args, lang, note) {
  const a = api();
  if (!a) return;
  if (!props.input.trim()) {
    emit('error', '请输入内容');
    return;
  }
  busy.value = true;
  try {
    const out = await a[fn](...args);
    emit('result', { text: out, lang, note });
  } catch (err) {
    emit('error', err?.message || String(err));
  } finally {
    busy.value = false;
  }
}

function convert(method, lang, note) {
  call(method, [props.input], lang, note);
}

async function queryPath() {
  const a = api();
  if (!a) return;
  if (!props.input.trim()) {
    emit('error', '请输入内容');
    return;
  }
  if (!path.value.trim()) {
    emit('error', '请输入 JSONPath 表达式');
    return;
  }
  busy.value = true;
  try {
    const res = await a.JSONPathQuery(props.input, path.value.trim());
    if (res.count === 0) {
      emit('error', '没有匹配结果');
    } else {
      emit('result', { text: res.output, lang: 'json', note: `JSONPath 匹配 ${res.count} 处 ✅` });
    }
  } catch (err) {
    emit('error', err?.message || String(err));
  } finally {
    busy.value = false;
  }
}

const codeLangMap = { go: 'go', typescript: 'ts', java: 'java' };

async function genCode() {
  const a = api();
  if (!a) return;
  if (!props.input.trim()) {
    emit('error', '请输入JSON内容');
    return;
  }
  busy.value = true;
  try {
    const opts = codeLang.value === 'java' ? { javaLombok: javaLombok.value } : {};
    const out = await a.GenerateCode(props.input, codeLang.value, opts);
    emit('result', {
      text: out,
      lang: codeLangMap[codeLang.value],
      note: `${codeLang.value === 'typescript' ? 'TypeScript' : codeLang.value} 类型定义已生成 ✅`,
    });
  } catch (err) {
    emit('error', err?.message || String(err));
  } finally {
    busy.value = false;
  }
}
</script>

<style scoped>
.toolstrip {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 6px 16px;
  background-color: var(--panel);
  border-bottom: 1px solid var(--border);
  flex-wrap: wrap;
}

.tabs {
  display: flex;
  gap: 4px;
}

.tab {
  padding: 6px 14px;
  border: 1px solid transparent;
  border-radius: 6px;
  background: transparent;
  color: var(--muted);
  font-size: 13px;
  font-family: inherit;
  cursor: pointer;
  transition: all 0.15s ease;
}
.tab:hover { color: var(--text); }
.tab.active {
  color: var(--accent);
  background-color: var(--panel-alt);
  border-color: var(--border);
  font-weight: 600;
}

.controls {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  min-width: 0;
}

.tbtn {
  padding: 6px 12px;
  border: 1px solid var(--border);
  border-radius: 6px;
  background-color: var(--panel-alt);
  color: var(--text);
  font-size: 12.5px;
  font-family: inherit;
  cursor: pointer;
  transition: all 0.15s ease;
  white-space: nowrap;
}
.tbtn:hover:not(:disabled) { background-color: var(--border); }
.tbtn:disabled { opacity: 0.5; cursor: wait; }
.tbtn.primary {
  background-color: var(--accent);
  border-color: var(--accent);
  color: #fff;
}
.tbtn.primary:hover:not(:disabled) { background-color: var(--accent-hover); }

.path-input {
  flex: 1;
  min-width: 180px;
  max-width: 420px;
  padding: 7px 10px;
  border: 1px solid var(--border);
  border-radius: 6px;
  background-color: var(--bg);
  color: var(--text);
  font-size: 13px;
  font-family: 'Cascadia Code', 'JetBrains Mono', Consolas, monospace;
}
.path-input:focus { outline: none; border-color: var(--accent); }

.select {
  padding: 6px 10px;
  border: 1px solid var(--border);
  border-radius: 6px;
  background-color: var(--panel-alt);
  color: var(--text);
  font-size: 13px;
  font-family: inherit;
  cursor: pointer;
}

.tp-hint {
  font-size: 12px;
  color: var(--muted);
  user-select: none;
}

.chk {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 12.5px;
  color: var(--text);
  cursor: pointer;
  user-select: none;
  white-space: nowrap;
}
.chk input { accent-color: var(--accent); cursor: pointer; }

@media (max-width: 900px) {
  .tp-hint { display: none; }
}
</style>
