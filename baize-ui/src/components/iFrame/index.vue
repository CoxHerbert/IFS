<template>
  <div class="iframe-wrapper" :style="{ height }">
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
import { onBeforeUnmount, onMounted, ref } from "vue";

const props = defineProps({
  src: {
    type: String,
    required: true
  }
});

const height = ref("");
const loading = ref(true);

function syncHeight() {
  height.value = `${document.documentElement.clientHeight - 94.5}px`;
}

function handleLoad() {
  loading.value = false;
}

onMounted(() => {
  syncHeight();
  window.addEventListener("resize", syncHeight);
});

onBeforeUnmount(() => {
  window.removeEventListener("resize", syncHeight);
});
</script>

<style scoped>
.iframe-wrapper {
  width: 100%;
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
  width: 100%;
  height: 100%;
  border: 0;
}
</style>
