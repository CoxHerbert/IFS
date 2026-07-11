<template>
  <a-card :title="block.title" :bordered="false" class="block-card">
    <pre>{{ block.content }}</pre>
  </a-card>
</template>

<script setup lang="ts">
import type { AgentBlock } from '@/types/agent'

defineProps<{ block: AgentBlock }>()
</script>

<style scoped>
.block-card {
  min-width: 0;
  max-width: 100%;
  border-radius: 8px;
  background: #f8fafc;
}

pre {
  margin: 0;
  white-space: pre-wrap;
  overflow-wrap: anywhere;
  word-break: break-word;
  font-family: inherit;
  line-height: 1.7;
}
</style>
