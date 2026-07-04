<template>
  <a-card :title="block.title" :bordered="false" class="block-card">
    <vxe-table :data="block.data || []" border stripe size="small" max-height="360">
      <vxe-column
        v-for="column in block.columns || []"
        :key="column.field"
        :field="column.field"
        :title="column.label"
        min-width="140"
      />
    </vxe-table>
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
