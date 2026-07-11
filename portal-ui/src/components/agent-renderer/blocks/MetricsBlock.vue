<template>
  <a-card :title="block.title" :bordered="false" class="block-card">
    <div class="metric-grid">
      <div v-for="item in block.items || []" :key="item.label" class="metric-item">
        <span class="metric-label">{{ item.label }}</span>
        <strong class="metric-value">{{ item.value }}</strong>
      </div>
    </div>
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

.metric-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(128px, 1fr));
  gap: 12px;
  min-width: 0;
}

.metric-item {
  min-width: 0;
  display: grid;
  gap: 8px;
  padding: 12px;
  border: 1px solid rgba(148, 163, 184, 0.16);
  border-radius: 8px;
  background: #fff;
}

.metric-label {
  color: #64748b;
  font-size: 12px;
  line-height: 1.4;
}

.metric-value {
  min-width: 0;
  color: #0f172a;
  font-size: 20px;
  font-weight: 600;
  line-height: 1.35;
  overflow-wrap: anywhere;
  word-break: break-word;
}

@media (max-width: 520px) {
  .metric-grid {
    grid-template-columns: 1fr;
  }
}
</style>
