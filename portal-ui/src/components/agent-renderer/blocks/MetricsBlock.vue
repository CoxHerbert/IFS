<template>
  <a-card :title="block.title" :bordered="false" class="block-card">
    <a-row :gutter="[12, 12]">
      <a-col v-for="item in block.items || []" :key="item.label" :xs="24" :sm="12" :md="8">
        <a-statistic :title="item.label" :value="String(item.value)" />
      </a-col>
    </a-row>
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
