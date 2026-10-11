<template>
  <div class="app-container">
    <!-- 顶部工具栏 -->
    <header class="toolbar">
      <div class="brand">
        <span class="logo">{}</span>
        <span class="name">pretty-json</span>
      </div>

      <div class="btn-group">
        <button @click="handleFormat" class="btn primary" title="Ctrl+Enter">格式化</button>
        <button @click="handleMinify" class="btn primary" title="Ctrl+Shift+M">压缩</button>
        <button @click="handleValidate" class="btn secondary" title="Ctrl+Shift+K">校验</button>
      </div>
      <div class="btn-divider"></div>
      <div class="btn-group">
        <button @click="handleOpen" class="btn secondary" title="Ctrl+O">打开</button>
        <button @click="handleSave" class="btn secondary" title="Ctrl+S">保存</button>
      </div>
      <div class="btn-divider"></div>
      <div class="btn-group">
        <button @click="handleCopy" class="btn secondary">复制结果</button>
        <button @click="handleClear" class="btn danger">清空</button>
      </div>

      <div class="spacer"></div>
      <span class="shortcut-hint">Ctrl+Enter 格式化 · Ctrl+Shift+M 压缩 · Ctrl+Shift+K 校验</span>
      <button @click="toggleTheme" class="btn icon" :title="theme === 'dark' ? '切换亮色主题' : '切换暗色主题'">
        {{ theme === 'dark' ? '☀️' : '🌙' }}
      </button>
    </header>

    <!-- M2 工具面板：格式互转 / JSONPath / 代码生成 -->
    <ToolPanel :input="inputValue" @result="onToolResult" @error="(m) => toast(m, 'error')" />

    <!-- 输入输出区域（可拖拽分隔） -->
    <main class="content" ref="contentRef">
      <section class="panel" :style="{ width: leftPct + '%' }">
        <div class="panel-head">
          <h3>输入 {{ inputLangLabel }}</h3>
          <span class="hint">实时校验 · 错误行列定位</span>
        </div>
        <JsonEditor ref="inputEditor" v-model="inputValue" :lang="inputLang" />
      </section>

      <div class="splitter" @mousedown="startDrag" title="拖拽调整宽度">
        <div class="splitter-bar"></div>
      </div>

      <section class="panel grow">
        <div class="panel-head">
          <h3>处理结果 <span class="lang-badge">{{ outputLangLabel }}</span></h3>
          <span class="hint">保留键顺序 · 大数字保精度</span>
        </div>
        <JsonEditor v-model="outputValue" :readonly="true" :lang="outputLang" />
      </section>
    </main>

    <!-- 状态栏 -->
    <StatusBar :input="inputValue" :result="analyzeResult" />

    <!-- Toast 通知 -->
    <ToastHost :toasts="toasts" />
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue';
import JsonEditor from './components/JsonEditor.vue';
import ToolPanel from './components/ToolPanel.vue';
import StatusBar from './components/StatusBar.vue';
import ToastHost from './components/ToastHost.vue';

// ---------- 状态 ----------
const inputValue = ref('');
const outputValue = ref('');
const analyzeResult = ref(null); // AnalyzeResult：{ valid, message, line, col, stats }
const inputEditor = ref(null);

// ---------- 编辑器语言（json | yaml | toml | go | java | ts） ----------
const inputLang = ref('json');
const outputLang = ref('json');
const currentFileName = ref('');

const langLabels = { json: 'JSON', yaml: 'YAML', toml: 'TOML', go: 'Go', java: 'Java', ts: 'TypeScript' };
const inputLangLabel = computed(() => langLabels[inputLang.value] || 'JSON');
const outputLangLabel = computed(() => langLabels[outputLang.value] || 'JSON');

// ---------- 主题 ----------
const theme = ref('light');
function applyTheme() {
  document.documentElement.dataset.theme = theme.value;
}
function toggleTheme() {
  theme.value = theme.value === 'dark' ? 'light' : 'dark';
}
watch(theme, (t) => {
  localStorage.setItem('pj-theme', t);
  applyTheme();
});

// ---------- Toast ----------
let toastSeq = 0;
const toasts = ref([]);
function toast(text, type = 'info') {
  const id = ++toastSeq;
  toasts.value.push({ id, type, text });
  setTimeout(() => {
    toasts.value = toasts.value.filter((t) => t.id !== id);
  }, 3000);
}

// ---------- 实时校验（300ms 防抖 + 竞态保护） ----------
let analyzeTimer = null;
let analyzeSeq = 0;

function refreshAnalyze(immediate = false) {
  if (analyzeTimer) clearTimeout(analyzeTimer);
  const run = async () => {
    if (!inputValue.value.trim() || inputLang.value !== 'json') {
      analyzeResult.value = null;
      inputEditor.value?.clearError();
      return;
    }
    const seq = ++analyzeSeq;
    try {
      const r = await window.go.main.App.AnalyzeJSON(inputValue.value);
      if (seq !== analyzeSeq) return; // 已有更新的请求，丢弃旧结果
      analyzeResult.value = r;
      if (r.valid) inputEditor.value?.clearError();
      else inputEditor.value?.setError(r.line, r.col, r.message);
    } catch {
      /* Wails 未就绪时忽略 */
    }
  };
  if (immediate) run();
  else analyzeTimer = setTimeout(run, 300);
}

watch(inputValue, () => refreshAnalyze());

// ---------- 操作 ----------
async function handleFormat() {
  if (!inputValue.value.trim()) {
    toast('请输入JSON内容', 'error');
    return;
  }
  try {
    outputValue.value = await window.go.main.App.FormatJSON(inputValue.value);
    outputLang.value = 'json';
    toast('格式化成功 ✅', 'success');
  } catch (err) {
    toast(err?.message || String(err), 'error');
  }
}

async function handleMinify() {
  if (!inputValue.value.trim()) {
    toast('请输入JSON内容', 'error');
    return;
  }
  try {
    outputValue.value = await window.go.main.App.MinifyJSON(inputValue.value);
    outputLang.value = 'json';
    toast('压缩成功 ✅', 'success');
  } catch (err) {
    toast(err?.message || String(err), 'error');
  }
}

async function handleValidate() {
  try {
    const r = await window.go.main.App.AnalyzeJSON(inputValue.value);
    analyzeResult.value = r;
    if (r.valid) {
      inputEditor.value?.clearError();
      toast('JSON格式正确 ✅', 'success');
    } else {
      inputEditor.value?.setError(r.line, r.col, r.message);
      toast(r.message, 'error');
    }
  } catch (err) {
    toast(err?.message || String(err), 'error');
  }
}

function handleCopy() {
  if (!outputValue.value.trim()) {
    toast('无内容可复制', 'error');
    return;
  }
  navigator.clipboard
    .writeText(outputValue.value)
    .then(() => toast('复制成功 ✅', 'success'))
    .catch(() => toast('复制失败', 'error'));
}

function handleClear() {
  inputValue.value = '';
  outputValue.value = '';
  analyzeResult.value = null;
  inputLang.value = 'json';
  outputLang.value = 'json';
  currentFileName.value = '';
  inputEditor.value?.clearError();
  toast('已清空', 'success');
}

// ---------- 文件打开 / 保存 / 拖放 ----------
function extToLang(name) {
  const m = /\.([a-z0-9]+)$/i.exec(name || '');
  const ext = m?.[1]?.toLowerCase();
  if (ext === 'yaml' || ext === 'yml') return 'yaml';
  if (ext === 'toml') return 'toml';
  return 'json';
}

function applyFile(fd) {
  if (!fd) return;
  inputValue.value = fd.content ?? '';
  currentFileName.value = fd.name || '';
  inputLang.value = extToLang(fd.name);
  refreshAnalyze(true);
  toast(`已打开 ${fd.name} ✅`, 'success');
}

async function handleOpen() {
  try {
    const fd = await window.go.main.App.OpenJSONFile();
    if (!fd) return; // 用户取消
    applyFile(fd);
  } catch (err) {
    toast(err?.message || String(err), 'error');
  }
}

const saveExtMap = { json: 'json', yaml: 'yaml', toml: 'toml', go: 'go', ts: 'ts', java: 'java' };

async function handleSave() {
  const content = outputValue.value || inputValue.value;
  if (!content.trim()) {
    toast('无内容可保存', 'error');
    return;
  }
  const base = currentFileName.value
    ? currentFileName.value.replace(/\.[^.]+$/, '')
    : 'pretty-json';
  const defaultName = `${base || 'pretty-json'}.${saveExtMap[outputLang.value] || 'json'}`;
  try {
    const p = await window.go.main.App.SaveTextFile(content, defaultName);
    if (p) toast(`已保存 ${p} ✅`, 'success');
  } catch (err) {
    toast(err?.message || String(err), 'error');
  }
}

// ToolPanel 结果统一落到输出编辑器，并切换其高亮语言
function onToolResult({ text, lang, note }) {
  outputValue.value = text;
  outputLang.value = lang || 'json';
  if (note) toast(note, 'success');
}

// ---------- 分隔条拖拽 ----------
const contentRef = ref(null);
const leftPct = ref(50);

function startDrag() {
  const onMove = (ev) => {
    if (!contentRef.value) return;
    const rect = contentRef.value.getBoundingClientRect();
    let pct = ((ev.clientX - rect.left) / rect.width) * 100;
    pct = Math.min(80, Math.max(20, pct));
    leftPct.value = pct;
  };
  const onUp = () => {
    document.body.style.userSelect = '';
    window.removeEventListener('mousemove', onMove);
    window.removeEventListener('mouseup', onUp);
  };
  document.body.style.userSelect = 'none';
  window.addEventListener('mousemove', onMove);
  window.addEventListener('mouseup', onUp);
}

// ---------- 快捷键 ----------
function onKeydown(e) {
  if (!(e.ctrlKey || e.metaKey)) return;
  const k = e.key.toLowerCase();
  if (k === 'enter') {
    e.preventDefault();
    handleFormat();
  } else if (e.shiftKey && k === 'm') {
    e.preventDefault();
    handleMinify();
  } else if (e.shiftKey && k === 'k') {
    e.preventDefault();
    handleValidate();
  } else if (k === 'o') {
    e.preventDefault();
    handleOpen();
  } else if (k === 's') {
    e.preventDefault();
    handleSave();
  }
}

// ---------- 生命周期 ----------
onMounted(() => {
  const saved = localStorage.getItem('pj-theme');
  theme.value = saved || (window.matchMedia?.('(prefers-color-scheme: dark)').matches ? 'dark' : 'light');
  applyTheme();
  window.addEventListener('keydown', onKeydown);
  window.addEventListener('wails:ready', () => {
    toast('就绪 ✅', 'success');
    refreshAnalyze(true);
  });
  // 拖放文件：后端 OnFileDrop 推送 file:drop 事件
  if (window.runtime?.EventsOn) {
    window.runtime.EventsOn('file:drop', applyFile);
    window.runtime.EventsOn('file:drop:error', (msg) => toast(msg, 'error'));
  }
});

onBeforeUnmount(() => {
  window.removeEventListener('keydown', onKeydown);
  if (analyzeTimer) clearTimeout(analyzeTimer);
});
</script>

<style scoped>
.app-container {
  display: flex;
  flex-direction: column;
  width: 100vw;
  height: 100vh;
  background-color: var(--bg);
}

/* ---------- 工具栏 ---------- */
.toolbar {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 16px;
  background-color: var(--panel);
  border-bottom: 1px solid var(--border);
  flex-wrap: wrap;
}

.brand {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-right: 10px;
  user-select: none;
}

.brand .logo {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  border-radius: 8px;
  background: linear-gradient(135deg, #2563eb, #7c3aed);
  color: #fff;
  font-weight: 700;
  font-size: 15px;
  font-family: Consolas, monospace;
}

.brand .name {
  font-size: 15px;
  font-weight: 600;
  color: var(--text);
}

.btn-group {
  display: flex;
  gap: 8px;
}

.btn-divider {
  width: 1px;
  height: 22px;
  background-color: var(--border);
  margin: 0 4px;
}

.btn {
  padding: 8px 16px;
  border: 1px solid transparent;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
  font-family: inherit;
  transition: all 0.15s ease;
}

.btn.primary {
  background-color: var(--accent);
  color: #fff;
}
.btn.primary:hover { background-color: var(--accent-hover); }

.btn.secondary {
  background-color: var(--panel-alt);
  color: var(--text);
  border-color: var(--border);
}
.btn.secondary:hover { background-color: var(--border); }

.btn.danger {
  background-color: transparent;
  color: var(--danger);
  border-color: var(--danger);
}
.btn.danger:hover { background-color: var(--danger); color: #fff; }

.btn.icon {
  background-color: var(--panel-alt);
  border-color: var(--border);
  font-size: 15px;
  padding: 6px 12px;
}

.spacer { flex: 1; }

.shortcut-hint {
  font-size: 12px;
  color: var(--muted);
  user-select: none;
}

/* ---------- 主内容区 ---------- */
.content {
  display: flex;
  flex: 1;
  min-height: 0;
  padding: 12px;
  gap: 0;
}

.panel {
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-width: 240px;
  min-height: 0;
}

.panel.grow {
  flex: 1;
  min-width: 0;
}

.panel-head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
}

.panel-head h3 {
  font-size: 14px;
  color: var(--text);
  font-weight: 600;
}

.lang-badge {
  display: inline-block;
  margin-left: 6px;
  padding: 1px 8px;
  border-radius: 10px;
  background-color: var(--panel-alt);
  border: 1px solid var(--border);
  color: var(--accent);
  font-size: 11px;
  font-weight: 600;
  vertical-align: 1px;
  user-select: none;
}

.panel-head .hint {
  font-size: 12px;
  color: var(--muted);
}

.splitter {
  width: 14px;
  align-self: stretch;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: col-resize;
  flex-shrink: 0;
}

.splitter-bar {
  width: 4px;
  height: 48px;
  border-radius: 2px;
  background-color: var(--border);
  transition: background-color 0.15s ease;
}

.splitter:hover .splitter-bar {
  background-color: var(--accent);
}

@media (max-width: 768px) {
  .content {
    flex-direction: column;
  }
  .panel {
    width: 100% !important;
    flex: 1;
  }
  .splitter {
    display: none;
  }
  .shortcut-hint {
    display: none;
  }
}
</style>
