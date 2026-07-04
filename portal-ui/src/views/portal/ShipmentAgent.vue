<template>
  <main class="shipment-agent-page">
    <section class="page-head">
      <div>
        <h1>出货计划智能分析</h1>
        <p>上传 Excel/CSV 后由服务端解析、计算 CBM，并交给本地模型生成专业建议。</p>
      </div>
      <a-button type="primary" :disabled="!selectedFile" :loading="analyzing" @click="handleAnalyze">
        提交 Agent 分析
      </a-button>
    </section>

    <ExcelUploader @select="handleFile" />

    <a-alert v-if="selectedFile" type="info" show-icon :message="`已选择文件：${selectedFile.name}`" />
    <a-alert v-if="errorMessage" type="error" show-icon :message="errorMessage" />

    <AgentResultRenderer v-if="agentResult" :result="agentResult" />
  </main>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { message as antMessage } from 'ant-design-vue'
import AgentResultRenderer from '@/components/agent-renderer/AgentResultRenderer.vue'
import ExcelUploader from '@/components/shipment/ExcelUploader.vue'
import { analyzeShipment } from '@/api/shipmentAgent'
import type { AgentResult } from '@/types/agent'

const selectedFile = ref<File>()
const agentResult = ref<AgentResult>()
const errorMessage = ref('')
const analyzing = ref(false)

function handleFile(file: File) {
  selectedFile.value = file
  agentResult.value = undefined
  errorMessage.value = ''
  antMessage.success(`已选择 ${file.name}`)
}

async function handleAnalyze() {
  if (!selectedFile.value) {
    errorMessage.value = '请先选择 Excel 或 CSV 文件'
    return
  }
  analyzing.value = true
  errorMessage.value = ''
  try {
    agentResult.value = await analyzeShipment(selectedFile.value)
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Agent 分析失败'
  } finally {
    analyzing.value = false
  }
}
</script>

<style scoped>
.shipment-agent-page {
  width: min(1180px, calc(100% - 40px));
  margin: 28px auto;
  display: grid;
  gap: 18px;
}

.page-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
}

.page-head h1 {
  margin: 0;
  font-size: 30px;
}

.page-head p {
  margin: 8px 0 0;
  color: #64748b;
}
</style>
