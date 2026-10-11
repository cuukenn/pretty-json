<template>
  <TransitionGroup name="toast" tag="div" class="toast-host">
    <div v-for="t in toasts" :key="t.id" class="toast" :class="t.type">
      {{ t.text }}
    </div>
  </TransitionGroup>
</template>

<script setup>
defineProps({
  toasts: { type: Array, default: () => [] }, // [{ id, type, text }]
});
</script>

<style scoped>
.toast-host {
  position: fixed;
  right: 20px;
  bottom: 46px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  z-index: 1000;
  pointer-events: none;
}

.toast {
  padding: 10px 16px;
  border-radius: 8px;
  font-size: 13px;
  color: #fff;
  background-color: #334155;
  box-shadow: var(--shadow);
  max-width: 420px;
  word-break: break-all;
}

.toast.success { background-color: #059669; }
.toast.error { background-color: #dc2626; }

.toast-enter-active, .toast-leave-active {
  transition: all 0.25s ease;
}
.toast-enter-from, .toast-leave-to {
  opacity: 0;
  transform: translateY(8px);
}
.toast-leave-active {
  position: absolute;
  right: 0;
}
</style>
