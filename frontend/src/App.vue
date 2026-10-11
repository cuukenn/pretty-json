<template>
  <div class="app-container">
    <!-- 顶部工具栏 -->
    <header class="toolbar">
      <div class="brand">
        <span class="logo">{}</span>
        <span class="name">pretty-json</span>
      </div>

      <div class="btn-group">
        <button @click="handleFormat" class="btn primary" title="格式化 (Ctrl+Enter)" aria-label="格式化">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
            <path d="M8 3H7a2 2 0 0 0-2 2v5a2 2 0 0 1-2 2 2 2 0 0 1 2 2v5a2 2 0 0 0 2 2h1" />
            <path d="M16 3h1a2 2 0 0 1 2 2v5a2 2 0 0 0 2 2 2 2 0 0 0-2 2v5a2 2 0 0 1-2 2h-1" />
          </svg>
        </button>
        <button @click="handleMinify" class="btn primary" title="压缩 (Ctrl+Shift+M)" aria-label="压缩">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
            <polyline points="4 14 10 14 10 20" />
            <polyline points="20 10 14 10 14 4" />
            <line x1="14" y1="10" x2="21" y2="3" />
            <line x1="3" y1="21" x2="10" y2="14" />
          </svg>
        </button>
        <button @click="handleValidate" class="btn secondary" title="校验 (Ctrl+Shift+K)" aria-label="校验">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
            <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14" />
            <polyline points="22 4 12 14.01 9 11.01" />
          </svg>
        </button>
      </div>
      <div class="btn-divider"></div>
      <div class="btn-group">
        <button @click="handleOpen" class="btn secondary" title="打开文件 (Ctrl+O)" aria-label="打开文件">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
            <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z" />
          </svg>
        </button>
        <button @click="handleSave" class="btn secondary" title="保存 (Ctrl+S)" aria-label="保存">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
            <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z" />
            <polyline points="17 21 17 13 7 13 7 21" />
            <polyline points="7 3 7 8 15 8" />
          </svg>
        </button>
      </div>
      <div class="btn-divider"></div>
      <div class="btn-group">
        <button @click="handleCopy" class="btn secondary" title="复制结果" aria-label="复制结果">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
            <rect x="9" y="9" width="13" height="13" rx="2" ry="2" />
            <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
          </svg>
        </button>
        <button @click="openDiff" class="btn secondary" title="对比输入与输出 (Ctrl+D)" aria-label="对比">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
            <circle cx="18" cy="18" r="3" />
            <circle cx="6" cy="6" r="3" />
            <path d="M13 6h3a2 2 0 0 1 2 2v7" />
            <line x1="6" y1="9" x2="6" y2="21" />
          </svg>
        </button>
        <button @click="handleClear" class="btn danger" title="清空输入与输出" aria-label="清空">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
            <polyline points="3 6 5 6 21 6" />
            <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
            <line x1="10" y1="11" x2="10" y2="17" />
            <line x1="14" y1="11" x2="14" y2="17" />
          </svg>
        </button>
      </div>

      <div class="spacer"></div>
      <span class="shortcut-hint">Ctrl+Enter 格式化 · Ctrl+Shift+M 压缩 · Ctrl+Shift+K 校验 · F1 帮助</span>
      <button @click="showHelp = true" class="btn icon" title="快捷键帮助 (F1)" aria-label="帮助">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24"
          fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
          <circle cx="12" cy="12" r="10" />
          <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3" />
          <line x1="12" y1="17" x2="12.01" y2="17" />
        </svg>
      </button>
      <div class="pop-anchor">
        <button @click="showHistory = !showHistory" class="btn icon" title="输入历史" aria-label="输入历史" aria-haspopup="true">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24"
            fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
            <circle cx="12" cy="12" r="10" />
            <polyline points="12 6 12 12 16 14" />
          </svg>
        </button>
        <div v-if="showHistory" class="pop-backdrop" @click="showHistory = false"></div>
        <div v-if="showHistory" class="popover" role="menu">
          <div class="pop-head">
            <span>最近输入（上限 {{ HISTORY_MAX }} 条）</span>
            <button class="pop-clear" @click="clearHistory" :disabled="!history.length">清空</button>
          </div>
          <ul class="pop-list">
            <li v-for="(h, i) in history" :key="h.ts">
              <button class="pop-item" @click="restoreHistory(h)">
                <span class="pop-item-name">{{ h.name || '手动输入' }}</span>
                <span class="pop-item-meta">{{ fmtTs(h.ts) }} · {{ fmtChars(h.text.length) }}</span>
                <span class="pop-item-prev">{{ h.text.slice(0, 80) }}</span>
              </button>
              <button class="pop-del" @click="deleteHistory(i)" title="删除此条" aria-label="删除此条">✕</button>
            </li>
            <li v-if="!history.length" class="pop-empty">暂无历史 · 打开文件或格式化成功后自动记录</li>
          </ul>
        </div>
      </div>
      <button @click="cycleTheme" class="btn icon" :title="`主题：${themeLabel}（点击切换）`" aria-label="切换主题">
        <svg v-if="themePref === 'system'" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24"
          fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
          <rect x="2" y="3" width="20" height="14" rx="2" ry="2" />
          <line x1="8" y1="21" x2="16" y2="21" />
          <line x1="12" y1="17" x2="12" y2="21" />
        </svg>
        <svg v-else-if="theme === 'dark'" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24"
          fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
          <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
        </svg>
        <svg v-else xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24"
          fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
          <circle cx="12" cy="12" r="4" />
          <path d="M12 2v2m0 16v2M4.93 4.93l1.41 1.41m11.32 11.32 1.41 1.41M2 12h2m16 0h2M4.93 19.07l1.41-1.41M17.66 6.34l1.41-1.41" />
        </svg>
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
        <div class="editor-wrap">
          <JsonEditor ref="inputEditor" v-model="inputValue" :lang="inputLang" />
          <div v-if="!inputValue" class="empty-state">
            <div class="empty-logo" aria-hidden="true">{ }</div>
            <p>粘贴 / 输入 JSON，或拖入文件打开</p>
            <button class="empty-btn" @click="loadSample">加载示例 JSON</button>
          </div>
        </div>
      </section>

      <div class="splitter" @mousedown="startDrag" title="拖拽调整宽度">
        <div class="splitter-bar"></div>
      </div>

      <section class="panel grow">
        <div class="panel-head">
          <h3>处理结果 <span class="lang-badge">{{ outputLangLabel }}</span></h3>
          <div class="head-right">
            <span class="hint">保留键顺序 · 大数字保精度</span>
            <div class="view-switch" role="group" aria-label="输出视图切换">
              <button :class="{ on: outputView === 'text' }" @click="outputView = 'text'">文本</button>
              <button
                :class="{ on: outputView === 'tree' }"
                :disabled="!canTreeView"
                :title="canTreeView ? '树视图（可折叠浏览）' : '文档过大或为空，树视图不可用'"
                @click="switchTree"
              >树</button>
            </div>
          </div>
        </div>
        <JsonEditor v-show="outputView === 'text'" v-model="outputValue" :readonly="true" :lang="outputLang" />
        <JsonTree v-if="outputView === 'tree' && treeRoot" :root="treeRoot" />
        <div v-if="outputView === 'tree' && !treeRoot" class="tree-empty">树视图不可用（JSON 无效或正在构建）</div>
      </section>
    </main>

    <!-- 状态栏 -->
    <StatusBar :input="inputValue" :result="analyzeResult" />

    <!-- Toast 通知 -->
    <ToastHost :toasts="toasts" />

    <!-- 快捷键帮助 -->
    <div v-if="showHelp" class="modal-backdrop" @click.self="showHelp = false">
      <div class="modal" role="dialog" aria-label="快捷键帮助">
        <div class="modal-head">
          <h3>快捷键</h3>
          <button class="modal-x" @click="showHelp = false" aria-label="关闭">✕</button>
        </div>
        <dl class="kb-list">
          <dt>Ctrl+Enter</dt><dd>格式化</dd>
          <dt>Ctrl+Shift+M</dt><dd>压缩</dd>
          <dt>Ctrl+Shift+K</dt><dd>校验</dd>
          <dt>Ctrl+O</dt><dd>打开文件</dd>
          <dt>Ctrl+S</dt><dd>另存为</dd>
          <dt>Ctrl+D</dt><dd>对比输入与输出</dd>
          <dt>拖入文件</dt><dd>打开 JSON / YAML / TOML（≤20MB）</dd>
          <dt>F1 / Ctrl+/</dt><dd>开关本帮助</dd>
          <dt>Esc</dt><dd>关闭弹层 / 帮助</dd>
        </dl>
      </div>
    </div>

    <!-- 输入/输出对比 -->
    <div v-if="showDiff" class="modal-backdrop" @click.self="showDiff = false">
      <div class="modal wide" role="dialog" aria-label="输入与输出对比">
        <div class="modal-head">
          <h3>对比：输入 → 输出
            <span class="diff-legend"><i class="lg del"></i>仅输入<i class="lg add"></i>仅输出</span>
          </h3>
          <button class="modal-x" @click="showDiff = false" aria-label="关闭">✕</button>
        </div>
        <div class="diff-body">
          <div v-for="(r, i) in diffRows" :key="i" class="dline" :class="r.t">
            <span class="dno">{{ r.no }}</span>
            <span class="dmark">{{ r.t === 'add' ? '+' : r.t === 'del' ? '-' : ' ' }}</span>
            <span class="dtxt">{{ r.text }}</span>
          </div>
          <div v-if="!diffRows.length" class="diff-empty">两侧内容完全一致</div>
          <div v-if="diffTruncated" class="diff-more">差异过大，仅显示前 {{ DIFF_ROWS_MAX }} 行</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue';
import JsonEditor from './components/JsonEditor.vue';
import JsonTree from './components/JsonTree.vue';
import ToolPanel from './components/ToolPanel.vue';
import StatusBar from './components/StatusBar.vue';
import ToastHost from './components/ToastHost.vue';
import { diffLines } from 'diff';

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

// ---------- 输出树视图 ----------
const outputView = ref('text'); // text | tree
const treeRoot = ref(null);
const TREE_LIMIT = 1000000; // 前端树视图字符上限（后端另有 5MB 硬上限）
const treeSource = computed(() => outputValue.value || inputValue.value);
const canTreeView = computed(() => !!treeSource.value.trim() && treeSource.value.length <= TREE_LIMIT);

let treeTimer = null;
let treeSeq = 0;
async function refreshTree(immediate = false) {
  if (outputView.value !== 'tree') return;
  if (treeTimer) clearTimeout(treeTimer);
  const run = async () => {
    const src = treeSource.value;
    if (!src.trim() || src.length > TREE_LIMIT) {
      treeRoot.value = null;
      return;
    }
    const seq = ++treeSeq;
    try {
      const t = await window.go.main.App.BuildJSONTree(src);
      if (seq === treeSeq) treeRoot.value = t;
    } catch (err) {
      if (seq !== treeSeq) return;
      treeRoot.value = null;
      outputView.value = 'text';
      toast(err?.message || '树视图构建失败', 'error');
    }
  };
  if (immediate) run();
  else treeTimer = setTimeout(run, 300);
}
function switchTree() {
  if (!canTreeView.value) {
    toast('文档过大或为空，树视图不可用', 'error');
    return;
  }
  outputView.value = 'tree';
  refreshTree(true);
}
watch([outputView, treeSource], () => refreshTree());

// ---------- 输入 / 输出对比 ----------
const showDiff = ref(false);
const diffRows = ref([]);
const diffTruncated = ref(false);
const DIFF_LIMIT = 200000; // 参与对比的单侧字符上限
const DIFF_ROWS_MAX = 4000;

function openDiff() {
  const a = inputValue.value;
  const b = outputValue.value;
  if (!a.trim() && !b.trim()) {
    toast('两侧均为空，无可对比内容', 'error');
    return;
  }
  if (a.length > DIFF_LIMIT || b.length > DIFF_LIMIT) {
    toast('文档过大，对比上限为 20 万字符', 'error');
    return;
  }
  const rows = [];
  let oldNo = 1;
  let newNo = 1;
  let truncated = false;
  for (const part of diffLines(a, b)) {
    const lines = part.value.split('\n');
    if (lines.length && lines[lines.length - 1] === '') lines.pop(); // 末尾换行不产出行
    for (const text of lines) {
      if (rows.length >= DIFF_ROWS_MAX) { truncated = true; break; }
      if (part.added) rows.push({ t: 'add', no: String(newNo++), text });
      else if (part.removed) rows.push({ t: 'del', no: String(oldNo++), text });
      else { rows.push({ t: 'ctx', no: String(oldNo), text }); oldNo++; newNo++; }
    }
    if (truncated) break;
  }
  diffRows.value = rows;
  diffTruncated.value = truncated;
  showDiff.value = true;
}

// ---------- 主题（浅色 → 深色 → 跟随系统 三态循环） ----------
const themePref = ref('system'); // dark | light | system
const sysDark = ref(false);
const theme = computed(() => (themePref.value === 'system' ? (sysDark.value ? 'dark' : 'light') : themePref.value));
const themeLabel = computed(() => ({ dark: '深色', light: '浅色', system: '跟随系统' })[themePref.value]);
watch(theme, (t) => { document.documentElement.dataset.theme = t; }, { immediate: true });
watch(themePref, (t) => localStorage.setItem('pj-theme', t));
function cycleTheme() {
  themePref.value = themePref.value === 'light' ? 'dark' : themePref.value === 'dark' ? 'system' : 'light';
}
function onSysThemeChange(e) { sysDark.value = e.matches; }

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

// ---------- 输入历史（localStorage，最多 10 条） ----------
const HISTORY_KEY = 'pj-history';
const HISTORY_MAX = 10;
const HISTORY_ITEM_MAX = 512 * 1024; // 单条上限 512KB，防止撑爆存储
const history = ref([]);
const showHistory = ref(false);
const showHelp = ref(false);

function loadHistory() {
  try {
    const raw = localStorage.getItem(HISTORY_KEY);
    const arr = raw ? JSON.parse(raw) : [];
    history.value = Array.isArray(arr) ? arr.slice(0, HISTORY_MAX) : [];
  } catch {
    history.value = [];
  }
}
function saveHistory() {
  try {
    localStorage.setItem(HISTORY_KEY, JSON.stringify(history.value));
  } catch {
    /* 存储满时静默降级 */
  }
}
function pushHistory(text, name) {
  if (!text || text.length > HISTORY_ITEM_MAX) return;
  if (history.value[0] && history.value[0].text === text) return; // 连续相同不重复记录
  history.value.unshift({ text, name: name || '', ts: Date.now() });
  if (history.value.length > HISTORY_MAX) history.value.length = HISTORY_MAX;
  saveHistory();
}
function restoreHistory(h) {
  inputValue.value = h.text;
  currentFileName.value = h.name || '';
  inputLang.value = extToLang(h.name);
  showHistory.value = false;
  refreshAnalyze(true);
  toast('已恢复历史内容 ✅', 'success');
}
function deleteHistory(i) {
  history.value.splice(i, 1);
  saveHistory();
}
function clearHistory() {
  history.value = [];
  saveHistory();
}
function fmtTs(ts) {
  return new Date(ts).toLocaleString('zh-CN', { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' });
}
function fmtChars(n) {
  return n >= 1048576 ? (n / 1048576).toFixed(1) + ' MB' : n >= 1024 ? (n / 1024).toFixed(1) + ' KB' : n + ' B';
}

// ---------- 示例 JSON（空状态引导） ----------
const sampleJson = `{
  "project": "pretty-json",
  "version": "1.0.0",
  "stars": 9007199254740993,
  "rating": 4.85,
  "private": false,
  "homepage": null,
  "keywords": ["json", "formatter", "jsonpath", "codegen"],
  "author": { "name": "cuukenn", "location": "上海" },
  "releases": [
    { "tag": "v1.0", "notes": "格式化 / 压缩 / 校验" },
    { "tag": "v1.1", "notes": "YAML · TOML · JSONPath · 代码生成" }
  ]
}`;
function loadSample() {
  inputValue.value = sampleJson;
  inputLang.value = 'json';
  refreshAnalyze(true);
}

// ---------- 实时校验（300ms 防抖 + 竞态保护 + 大文档守卫） ----------
let analyzeTimer = null;
let analyzeSeq = 0;
const AUTO_ANALYZE_LIMIT = 1000000; // 字符数：超过则暂停实时校验，手动「校验」仍可用

function refreshAnalyze(immediate = false) {
  if (analyzeTimer) clearTimeout(analyzeTimer);
  const run = async () => {
    if (!inputValue.value.trim() || inputLang.value !== 'json') {
      analyzeResult.value = null;
      inputEditor.value?.clearError();
      return;
    }
    if (!immediate && inputValue.value.length > AUTO_ANALYZE_LIMIT) {
      analyzeResult.value = { valid: false, paused: true };
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
    pushHistory(inputValue.value, currentFileName.value);
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
    pushHistory(inputValue.value, currentFileName.value);
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
  treeRoot.value = null;
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
  pushHistory(fd.content ?? '', fd.name);
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
  let raf = 0;
  let pendingX = 0;
  const apply = () => {
    raf = 0;
    if (!contentRef.value) return;
    const rect = contentRef.value.getBoundingClientRect();
    let pct = ((pendingX - rect.left) / rect.width) * 100;
    pct = Math.min(80, Math.max(20, pct));
    leftPct.value = pct;
  };
  const onMove = (ev) => {
    pendingX = ev.clientX;
    if (!raf) raf = requestAnimationFrame(apply); // rAF 节流：每帧至多一次布局
  };
  const onUp = () => {
    if (raf) cancelAnimationFrame(raf);
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
  if (e.key === 'Escape') {
    showHelp.value = false;
    showHistory.value = false;
    showDiff.value = false;
    return;
  }
  if (e.key === 'F1') {
    e.preventDefault();
    showHelp.value = !showHelp.value;
    return;
  }
  if ((e.ctrlKey || e.metaKey) && e.key === '/') {
    e.preventDefault();
    showHelp.value = !showHelp.value;
    return;
  }
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
  } else if (k === 'd') {
    e.preventDefault();
    openDiff();
  }
}

// ---------- 生命周期 ----------
let sysMq = null;
onMounted(() => {
  const saved = localStorage.getItem('pj-theme');
  themePref.value = saved === 'dark' || saved === 'light' ? saved : 'system';
  sysMq = window.matchMedia?.('(prefers-color-scheme: dark)');
  if (sysMq) {
    sysDark.value = sysMq.matches;
    sysMq.addEventListener('change', onSysThemeChange);
  }
  loadHistory();
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
  if (sysMq) sysMq.removeEventListener('change', onSysThemeChange);
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
  box-shadow: var(--shadow);
  position: relative;
  z-index: 5;
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
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  padding: 0;
  border: 1px solid transparent;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
  font-family: inherit;
  transition: all 0.15s ease;
}
.btn svg { display: block; }
.btn:focus-visible { outline: 2px solid var(--accent); outline-offset: 1px; }

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
}
.btn.icon svg { display: block; }

/* ---------- 输入历史弹层 ---------- */
.pop-anchor { position: relative; }
.pop-backdrop { position: fixed; inset: 0; z-index: 50; }
.popover {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  width: 330px;
  max-width: 80vw;
  background-color: var(--panel);
  border: 1px solid var(--border);
  border-radius: 10px;
  box-shadow: var(--shadow);
  z-index: 60;
  padding: 8px;
}
.pop-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 4px 8px 8px;
  font-size: 12px;
  color: var(--muted);
}
.pop-clear {
  background: none;
  border: none;
  color: var(--danger);
  font-size: 12px;
  font-family: inherit;
  cursor: pointer;
}
.pop-clear:disabled { opacity: 0.4; cursor: default; }
.pop-list { list-style: none; margin: 0; padding: 0; max-height: 320px; overflow-y: auto; }
.pop-list li { display: flex; align-items: stretch; gap: 4px; }
.pop-item {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
  text-align: left;
  background: none;
  border: 1px solid transparent;
  border-radius: 6px;
  padding: 6px 8px;
  cursor: pointer;
  color: var(--text);
  font-family: inherit;
}
.pop-item:hover { background-color: var(--panel-alt); border-color: var(--border); }
.pop-item-name { font-size: 12.5px; font-weight: 600; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.pop-item-meta { font-size: 11px; color: var(--muted); }
.pop-item-prev {
  font-size: 11.5px;
  color: var(--muted);
  font-family: 'Cascadia Code', 'JetBrains Mono', Consolas, monospace;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.pop-del {
  background: none;
  border: none;
  color: var(--muted);
  font-size: 13px;
  cursor: pointer;
  padding: 0 6px;
}
.pop-del:hover { color: var(--danger); }
.pop-empty { font-size: 12px; color: var(--muted); padding: 10px 8px; }

/* ---------- 快捷键帮助模态 ---------- */
.modal-backdrop {
  position: fixed;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
}
.modal {
  width: 440px;
  max-width: calc(100vw - 40px);
  background-color: var(--panel);
  border: 1px solid var(--border);
  border-radius: 10px;
  box-shadow: var(--shadow);
}
.modal-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border);
}
.modal-head h3 { font-size: 14px; color: var(--text); }
.modal-x {
  background: none;
  border: none;
  color: var(--muted);
  font-size: 15px;
  cursor: pointer;
  line-height: 1;
}
.modal-x:hover { color: var(--text); }
.kb-list {
  display: grid;
  grid-template-columns: auto 1fr;
  gap: 8px 16px;
  padding: 14px 16px;
  margin: 0;
}
.kb-list dt {
  font-family: 'Cascadia Code', Consolas, monospace;
  font-size: 12px;
  background-color: var(--panel-alt);
  border: 1px solid var(--border);
  border-radius: 4px;
  padding: 2px 8px;
  justify-self: start;
  color: var(--text);
}
.kb-list dd { margin: 0; font-size: 13px; color: var(--text); align-self: center; }

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

/* ---------- 输入区空状态 ---------- */
.editor-wrap {
  position: relative;
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
}
.empty-state {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  pointer-events: none;
  color: var(--muted);
}
.empty-logo {
  font-size: 30px;
  font-family: Consolas, monospace;
  font-weight: 700;
  opacity: 0.3;
}
.empty-state p { font-size: 13px; margin: 0; }
.empty-btn {
  pointer-events: auto;
  padding: 7px 16px;
  border: 1px solid var(--accent);
  border-radius: 6px;
  background: transparent;
  color: var(--accent);
  font-size: 13px;
  font-family: inherit;
  cursor: pointer;
  transition: all 0.15s ease;
}
.empty-btn:hover { background: var(--accent); color: #fff; }

/* ---------- 输出视图切换 ---------- */
.head-right {
  display: flex;
  align-items: center;
  gap: 10px;
}
.view-switch {
  display: inline-flex;
  border: 1px solid var(--border);
  border-radius: 6px;
  overflow: hidden;
}
.view-switch button {
  padding: 4px 12px;
  font-size: 12px;
  font-family: inherit;
  border: none;
  background-color: var(--panel);
  color: var(--muted);
  cursor: pointer;
  transition: all 0.15s ease;
}
.view-switch button + button { border-left: 1px solid var(--border); }
.view-switch button.on { background-color: var(--accent); color: #fff; }
.view-switch button:disabled { opacity: 0.45; cursor: not-allowed; }
.tree-empty {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--muted);
  font-size: 13px;
  border: 1px dashed var(--border);
  border-radius: 8px;
}

/* ---------- 输入/输出对比模态 ---------- */
.modal.wide {
  width: 920px;
  max-width: calc(100vw - 40px);
}
.diff-legend {
  font-size: 12px;
  font-weight: 400;
  color: var(--muted);
  margin-left: 10px;
}
.diff-legend .lg {
  display: inline-block;
  width: 10px;
  height: 10px;
  border-radius: 2px;
  margin: 0 4px 0 10px;
  vertical-align: -1px;
}
.lg.del { background-color: rgba(239, 68, 68, 0.4); }
.lg.add { background-color: rgba(16, 185, 129, 0.4); }
.diff-body {
  max-height: 65vh;
  overflow: auto;
  margin: 12px 16px 16px;
  border: 1px solid var(--border);
  border-radius: 6px;
  font-family: 'JetBrains Mono', 'Cascadia Code', Consolas, 'Courier New', monospace;
  font-size: 12px;
}
.dline { display: flex; line-height: 1.6; }
.dline.add { background-color: rgba(16, 185, 129, 0.12); }
.dline.del { background-color: rgba(239, 68, 68, 0.12); }
.dno {
  flex: none;
  width: 46px;
  text-align: right;
  padding-right: 10px;
  color: var(--muted);
  user-select: none;
}
.dmark { flex: none; width: 16px; color: var(--muted); user-select: none; }
.dline.add .dmark { color: var(--success); }
.dline.del .dmark { color: var(--danger); }
.dtxt { flex: 1; white-space: pre-wrap; word-break: break-all; padding-right: 8px; color: var(--text); }
.diff-empty, .diff-more { padding: 16px; text-align: center; color: var(--muted); }

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
