<template>
  <div ref="wrapperRef" class="iframe-wrapper" :style="{ height }">
    <a-spin :spinning="loading" class="iframe-spin">
      <iframe
        :src="src"
        frameborder="0"
        class="iframe-content"
        scrolling="auto"
        @load="handleLoad"
      />
    </a-spin>
  </div>
</template>

<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, ref } from "vue";

const props = defineProps({
  src: {
    type: String,
    required: true
  }
});

const height = ref("");
const loading = ref(true);
const wrapperRef = ref<HTMLElement>();
let resizeObserver: ResizeObserver | undefined;
let resizeFrame = 0;

function syncHeight() {
  cancelAnimationFrame(resizeFrame);
  resizeFrame = requestAnimationFrame(() => {
    const top = wrapperRef.value?.getBoundingClientRect().top || 0;
    const viewportHeight = window.visualViewport?.height || document.documentElement.clientHeight;
    height.value = `${Math.max(Math.floor(viewportHeight - top), 320)}px`;
  });
}

function handleLoad() {
  loading.value = false;
}

onMounted(() => {
  nextTick(syncHeight);
  window.addEventListener("resize", syncHeight);
  window.visualViewport?.addEventListener("resize", syncHeight);
  resizeObserver = new ResizeObserver(syncHeight);
  if (wrapperRef.value?.parentElement) resizeObserver.observe(wrapperRef.value.parentElement);
});

onBeforeUnmount(() => {
  window.removeEventListener("resize", syncHeight);
  window.visualViewport?.removeEventListener("resize", syncHeight);
  resizeObserver?.disconnect();
  cancelAnimationFrame(resizeFrame);
});
</script>

<style scoped>
.iframe-wrapper {
  width: 100%;
  min-width: 0;
  overflow: hidden;
}

.iframe-spin {
  display: block;
  width: 100%;
  height: 100%;
}

.iframe-spin :deep(.ant-spin-container) {
  width: 100%;
  height: 100%;
}

.iframe-content {
  display: block;
  width: 100%;
  height: 100%;
  border: 0;
}
</style>
