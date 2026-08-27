<template>
  <a-card :title="block.title" :bordered="false" class="block-card">
    <a-form ref="formRef" :model="formState" layout="vertical">
      <a-form-item
        v-for="field in block.fields || []"
        :key="field.field"
        :label="field.label"
        :name="field.field"
        :rules="field.required ? [{ required: true, message: `请填写${field.label}` }] : []"
      >
        <a-input
          v-if="field.component === 'input'"
          v-model:value="formState[field.field]"
          :placeholder="field.placeholder"
        />
        <a-textarea
          v-else-if="field.component === 'textarea'"
          v-model:value="formState[field.field]"
          :placeholder="field.placeholder"
          :auto-size="{ minRows: 2, maxRows: 5 }"
        />
        <a-input-number
          v-else-if="field.component === 'number'"
          v-model:value="formState[field.field]"
          :placeholder="field.placeholder"
          style="width: 100%"
        />
        <a-select
          v-else-if="field.component === 'select'"
          v-model:value="formState[field.field]"
          :placeholder="field.placeholder || `请选择${field.label}`"
        >
          <a-select-option v-for="option in field.options || []" :key="String(option.value)" :value="option.value">
            {{ option.label }}
          </a-select-option>
        </a-select>
        <a-date-picker
          v-else-if="field.component === 'date'"
          v-model:value="formState[field.field]"
          value-format="YYYY-MM-DD"
          style="width: 100%"
        />
        <a-upload v-else-if="field.component === 'upload'" :before-upload="handleFile(field.field)" :max-count="1">
          <a-button>选择文件</a-button>
        </a-upload>
        <a-input v-else v-model:value="formState[field.field]" :placeholder="field.placeholder" />
      </a-form-item>

      <a-button type="primary" :loading="submitting" @click="handleSubmit">提交</a-button>
    </a-form>
  </a-card>
</template>

<script setup lang="ts">
import { reactive, ref, watch } from 'vue'
import { message as antMessage } from 'ant-design-vue'
import type { FormInstance } from 'ant-design-vue'
import type { AgentBlock, AgentResult } from '@/types/agent'
import { submitAgentForm } from '@/api/agentForm'

const props = defineProps<{
  block: AgentBlock
  sessionId?: number
}>()

const emit = defineEmits<{
  (event: 'submitted', result: AgentResult): void
}>()

const formRef = ref<FormInstance>()
const submitting = ref(false)
const formState = reactive<Record<string, unknown>>({})

function handleFile(field: string) {
  return (file: File) => {
    formState[field] = file
    return false
  }
}

watch(
  () => props.block,
  (block) => {
    Object.keys(formState).forEach((key) => delete formState[key])
    Object.assign(formState, block.initialValues || {})
  },
  { immediate: true },
)

async function handleSubmit() {
  if (!props.sessionId) {
    antMessage.warning('缺少会话 ID')
    return
  }
  await formRef.value?.validate()
  submitting.value = true
  try {
    const result = await submitAgentForm({
      sessionId: props.sessionId,
      formCode: props.block.formCode || '',
      submitApi: props.block.submitApi,
      values: { ...formState },
    })
    emit('submitted', result)
  } catch (error) {
    antMessage.error(error instanceof Error ? error.message : '提交失败')
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.block-card {
  border-radius: 8px;
  background: #f8fafc;
}
</style>
