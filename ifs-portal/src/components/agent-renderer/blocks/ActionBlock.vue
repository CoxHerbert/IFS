<template>
  <a-card :bordered="false" class="block-card">
    <a-button type="primary" :loading="executing" @click="handleExecute">
      {{ block.label || '执行操作' }}
    </a-button>
  </a-card>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { message as antMessage } from 'ant-design-vue'
import type { AgentBlock, AgentResult } from '@/types/agent'
import { executeAgentAction } from '@/api/agentForm'

const props = defineProps<{
  block: AgentBlock
  sessionId?: number
}>()

const emit = defineEmits<{
  (event: 'submitted', result: AgentResult): void
}>()

const executing = ref(false)

async function handleExecute() {
  if (!props.sessionId) {
    antMessage.warning('缺少会话 ID')
    return
  }
  executing.value = true
  try {
    const result = await executeAgentAction({
      sessionId: props.sessionId,
      actionCode: props.block.actionCode || '',
      payload: props.block.payload || {},
    })
    emit('submitted', result)
  } catch (error) {
    antMessage.error(error instanceof Error ? error.message : '执行失败')
  } finally {
    executing.value = false
  }
}
</script>

<style scoped>
.block-card {
  border-radius: 8px;
  background: #f8fafc;
}
</style>
