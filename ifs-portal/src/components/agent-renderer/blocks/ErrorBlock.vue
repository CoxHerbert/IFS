<template>
  <a-alert
    :message="block.title || '错误'"
    :description="block.content || '无法展示该结果。'"
    type="error"
    show-icon
  />
</template>

<script setup lang="ts">
import type { AgentBlock } from '@/types/agent'

defineProps<{ block: AgentBlock }>()
</script>
