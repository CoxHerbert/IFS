<template>
  <section v-if="result" class="agent-result">
    <div class="result-head">
      <strong>{{ result.title }}</strong>
      <span>{{ result.summary }}</span>
    </div>
    <component
      :is="getBlockComponent(block.type)"
      v-for="(block, index) in result.blocks"
      :key="`${block.type}-${block.title || index}`"
      :block="block"
      :session-id="sessionId"
      @submitted="handleSubmitted"
    />
  </section>
</template>

<script setup lang="ts">
import type { Component } from 'vue'
import type { AgentResult } from '@/types/agent'
import MetricsBlock from './blocks/MetricsBlock.vue'
import TableBlock from './blocks/TableBlock.vue'
import MarkdownBlock from './blocks/MarkdownBlock.vue'
import FileBlock from './blocks/FileBlock.vue'
import ErrorBlock from './blocks/ErrorBlock.vue'
import FormBlock from './blocks/FormBlock.vue'
import ActionBlock from './blocks/ActionBlock.vue'
import LinkBlock from './blocks/LinkBlock.vue'

defineProps<{
  result: AgentResult | undefined
  sessionId?: number
}>()

const emit = defineEmits<{
  (event: 'submitted', result: AgentResult): void
}>()

function getBlockComponent(type: string): Component {
  const blockMap: Record<string, Component> = {
    metrics: MetricsBlock,
    table: TableBlock,
    markdown: MarkdownBlock,
    file: FileBlock,
    link: LinkBlock,
    error: ErrorBlock,
    form: FormBlock,
    action: ActionBlock,
    summary: MarkdownBlock,
  }
  return blockMap[type] || ErrorBlock
}

function handleSubmitted(result: AgentResult) {
  emit('submitted', result)
}
</script>

<style scoped>
.agent-result {
  display: grid;
  gap: 12px;
  min-width: 0;
  width: 100%;
}

.result-head {
  display: grid;
  gap: 6px;
  min-width: 0;
}

.result-head strong {
  overflow-wrap: anywhere;
}

.result-head span {
  color: #64748b;
  line-height: 1.6;
  overflow-wrap: anywhere;
}

:deep(.ant-card) {
  min-width: 0;
  max-width: 100%;
}

:deep(.ant-card-body) {
  min-width: 0;
  overflow-x: auto;
}
</style>
