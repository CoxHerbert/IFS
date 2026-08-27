<template>
  <a-card :title="block.title" :bordered="false" class="block-card">
    <a-button type="primary" :href="block.url" :disabled="!block.url" download>
      下载{{ block.name ? `：${block.name}` : '' }}
    </a-button>
  </a-card>
</template>

<script setup lang="ts">
import type { AgentBlock } from '@/types/agent'

defineProps<{ block: AgentBlock }>()
</script>

<style scoped>
.block-card {
  border-radius: 8px;
  background: #f8fafc;
}
</style>
