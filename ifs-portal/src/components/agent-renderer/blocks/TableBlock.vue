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
  min-width: 0;
  max-width: 100%;
  border-radius: 8px;
  background: #f8fafc;
}

:deep(.ant-card-body) {
  overflow-x: auto;
}

:deep(.vxe-table) {
  min-width: 100%;
}
</style>
