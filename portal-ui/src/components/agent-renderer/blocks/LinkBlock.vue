<template>
  <a-card :title="block.title" :bordered="false" class="block-card">
    <p v-if="block.content" class="link-desc">{{ block.content }}</p>
    <a :href="block.url" class="link-action" :target="isExternalLink ? '_blank' : undefined" rel="noreferrer">
      {{ block.label || block.url }}
    </a>
  </a-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { AgentBlock } from '@/types/agent'

const props = defineProps<{ block: AgentBlock }>()

const isExternalLink = computed(() => /^https?:\/\//i.test(props.block.url || ''))
</script>

<style scoped>
.block-card {
  min-width: 0;
  max-width: 100%;
  border-radius: 8px;
  background: #f8fafc;
}

.link-desc {
  margin: 0 0 10px;
  color: #475569;
  line-height: 1.6;
  overflow-wrap: anywhere;
}

.link-action {
  color: #1677ff;
  font-weight: 600;
  line-height: 1.6;
  overflow-wrap: anywhere;
}
</style>
