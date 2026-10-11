<template>
  <div class="editor-host" ref="host"></div>
</template>

<script setup>
import { onBeforeUnmount, onMounted, ref, watch } from 'vue';
import { basicSetup } from 'codemirror';
import { Compartment, EditorState } from '@codemirror/state';
import { EditorView } from '@codemirror/view';
import { json } from '@codemirror/lang-json';
import { yaml } from '@codemirror/lang-yaml';
import { go } from '@codemirror/lang-go';
import { java } from '@codemirror/lang-java';
import { javascript } from '@codemirror/lang-javascript';
import { HighlightStyle, StreamLanguage, syntaxHighlighting } from '@codemirror/language';
import { toml as tomlMode } from '@codemirror/legacy-modes/mode/toml';
import { lintGutter, linter, setDiagnostics } from '@codemirror/lint';
import { tags as t } from '@lezer/highlight';

const props = defineProps({
  modelValue: { type: String, default: '' },
  readonly: { type: Boolean, default: false },
  lang: { type: String, default: 'json' }, // json | yaml | toml | go | java | ts
});
const emit = defineEmits(['update:modelValue']);

const host = ref(null);
let view = null;
let silent = false;

// 主题通过 CSS 变量实现，切换 data-theme 时编辑器颜色自动跟随
const themeExt = EditorView.theme({
  '&': { height: '100%', fontSize: '13px', backgroundColor: 'var(--panel)', color: 'var(--text)' },
  '.cm-scroller': {
    fontFamily: "'Cascadia Code', 'JetBrains Mono', Consolas, 'Courier New', monospace",
    lineHeight: '1.6',
  },
  '.cm-content': { caretColor: 'var(--accent)', padding: '8px 0' },
  '.cm-cursor, .cm-dropCursor': { borderLeftColor: 'var(--accent)' },
  '&.cm-focused .cm-selectionBackground, .cm-selectionBackground, .cm-content ::selection': {
    backgroundColor: 'var(--cm-select)',
  },
  '.cm-activeLine': { backgroundColor: 'var(--cm-active)' },
  '.cm-gutters': {
    backgroundColor: 'var(--panel)',
    color: 'var(--muted)',
    borderRight: '1px solid var(--border)',
    minWidth: '46px',
  },
  '.cm-activeLineGutter': { backgroundColor: 'var(--cm-active)' },
  '.cm-foldGutter span': { color: 'var(--muted)' },
  '.cm-lineNumbers .cm-gutterElement': { padding: '0 10px 0 14px' },
  '.cm-panels': { backgroundColor: 'var(--panel)', color: 'var(--text)', border: 'none' },
  '.cm-panels-bottom': { borderTop: '1px solid var(--border)' },
  '.cm-search': { backgroundColor: 'var(--panel)', border: '1px solid var(--border)' },
  '.cm-search button, .cm-search label': { background: 'var(--panel-alt)', color: 'var(--text)', border: '1px solid var(--border)' },
  '.cm-button': { backgroundImage: 'none', backgroundColor: 'var(--panel-alt)', color: 'var(--text)', border: '1px solid var(--border)', borderRadius: '4px' },
  '.cm-textfield': { backgroundColor: 'var(--panel)', color: 'var(--text)', border: '1px solid var(--border)' },
  '.cm-tooltip': { backgroundColor: 'var(--panel)', border: '1px solid var(--border)' },
  '.cm-tooltip-lint': { backgroundColor: 'var(--panel)' },
  '.cm-diagnostic': { borderLeftColor: 'var(--danger)', color: 'var(--text)' },
  '.cm-diagnostic-error': { borderLeftColor: 'var(--danger)' },
  '.cm-lintRange-error': {
    backgroundImage: 'url("data:image/svg+xml,%3Csvg xmlns=%27http://www.w3.org/2000/svg%27 width=%274%27 height=%272%27%3E%3Cpath d=%27M0 1 L1 0 L2 1 L3 0 L4 1%27 stroke=%27%23ef4444%27 fill=%27none%27/%3E%3C/svg%3E")',
    backgroundRepeat: 'repeat-x',
    backgroundPosition: 'bottom',
  },
});

const highlight = HighlightStyle.define([
  { tag: t.propertyName, color: 'var(--tok-key)' },
  { tag: t.string, color: 'var(--tok-string)' },
  { tag: t.number, color: 'var(--tok-number)' },
  { tag: [t.bool, t.keyword], color: 'var(--tok-bool)' },
  { tag: t.null, color: 'var(--tok-null)' },
  { tag: [t.punctuation, t.separator, t.brace, t.bracket], color: 'var(--tok-punct)' },
]);

const readonlyConf = new Compartment();
const langConf = new Compartment();

// StreamLanguage 实例复用，避免每次切换重复创建
const tomlExt = StreamLanguage.define(tomlMode);
function langExt(lang) {
  switch (lang) {
    case 'yaml':
      return yaml();
    case 'toml':
      return tomlExt;
    case 'go':
      return go();
    case 'java':
      return java();
    case 'ts':
      return javascript();
    default:
      return json();
  }
}

onMounted(() => {
  const state = EditorState.create({
    doc: props.modelValue || '',
    extensions: [
      basicSetup,
      langConf.of(langExt(props.lang)),
      themeExt,
      syntaxHighlighting(highlight),
      EditorView.lineWrapping,
      lintGutter(),
      // 编辑时自动清空旧标记；真正的诊断由后端分析后经 setDiagnostics 注入
      linter(() => [], { delay: 0 }),
      readonlyConf.of(
        props.readonly ? [EditorState.readOnly.of(true), EditorView.editable.of(false)] : []
      ),
      EditorView.updateListener.of((u) => {
        if (silent) return;
        if (u.docChanged) emit('update:modelValue', u.state.doc.toString());
      }),
    ],
  });
  view = new EditorView({ state, parent: host.value });
});

onBeforeUnmount(() => {
  view?.destroy();
  view = null;
});

watch(
  () => props.modelValue,
  (v) => {
    if (!view) return;
    const cur = view.state.doc.toString();
    const next = v ?? '';
    if (v === cur) return;
    silent = true;
    view.dispatch({ changes: { from: 0, to: cur.length, insert: next } });
    silent = false;
  }
);

watch(
  () => props.readonly,
  (ro) => {
    // state 6.7+：运行时切换扩展须用 Compartment.reconfigure()（返回 StateEffect），
    // Compartment.of() 只用于 extensions 数组
    view?.dispatch({
      effects: readonlyConf.reconfigure(ro ? [EditorState.readOnly.of(true), EditorView.editable.of(false)] : []),
    });
  }
);

watch(
  () => props.lang,
  (lang) => {
    view?.dispatch({ effects: langConf.reconfigure(langExt(lang)) });
  }
);

defineExpose({
  // 在指定行标记错误（line/col 为 1 起始，来自后端定位）
  setError(line, col, msg) {
    if (!view) return;
    let diags = [];
    if (line) {
      const l = Math.min(Math.max(line, 1), view.state.doc.lines);
      const lb = view.state.doc.line(l);
      const from = Math.min(lb.from + Math.max((col || 1) - 1, 0), lb.to);
      diags = [
        {
          from,
          to: Math.min(from + 1, lb.to),
          severity: 'error',
          message: msg || 'JSON 格式错误',
        },
      ];
    }
    view.dispatch(setDiagnostics(view.state, diags));
  },
  clearError() {
    if (!view) return;
    view.dispatch(setDiagnostics(view.state, []));
  },
  focus() {
    view?.focus();
  },
});
</script>

<style scoped>
.editor-host {
  height: 100%;
  min-height: 0;
  border: 1px solid var(--border);
  border-radius: 8px;
  overflow: hidden;
  background-color: var(--panel);
}
</style>
