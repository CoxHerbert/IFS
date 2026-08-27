<template>
  <div class="app-container agent-config-page">
    <a-alert
      type="info"
      show-icon
      message="本地模型服务地址"
      description="如果后台部署在阿里云，本地 Ollama 需要通过 VPN、内网穿透、专线或反向代理暴露一个后台服务器可访问的 HTTP 地址。"
    />

    <a-form ref="formRef" :model="form" :rules="rules" layout="vertical" class="config-form">
      <a-row :gutter="16">
        <a-col :xs="24" :md="12">
          <a-form-item label="Ollama Base URL" name="baseUrl">
            <a-input v-model:value="form.baseUrl" placeholder="http://localhost:11434" />
          </a-form-item>
        </a-col>
        <a-col :xs="24" :md="6">
          <a-form-item label="默认模型" name="defaultModel">
            <a-select v-model:value="form.defaultModel" :options="modelSelectOptions" />
          </a-form-item>
        </a-col>
        <a-col :xs="24" :md="6">
          <a-form-item label="超时时间（秒）" name="timeout">
            <a-input-number v-model:value="form.timeout" :min="5" :max="600" style="width: 100%" />
          </a-form-item>
        </a-col>
      </a-row>

      <div class="section-head">
        <strong>可选模型</strong>
        <a-button type="primary" ghost @click="addModel">新增模型</a-button>
      </div>

      <a-table :data-source="form.models" :columns="columns" :pagination="false" row-key="value" bordered size="small">
        <template #bodyCell="{ column, record, index }">
          <template v-if="column.key === 'label'">
            <a-input v-model:value="record.label" placeholder="Qwen 2.5 7B" />
          </template>
          <template v-else-if="column.key === 'value'">
            <a-input v-model:value="record.value" placeholder="qwen2.5:7b" @blur="syncDefaultModel" />
          </template>
          <template v-else-if="column.key === 'description'">
            <a-input v-model:value="record.description" placeholder="模型说明" />
          </template>
          <template v-else-if="column.key === 'default'">
            <a-radio :checked="record.default" @change="setDefault(index)" />
          </template>
          <template v-else-if="column.key === 'operation'">
            <a-button type="link" danger @click="removeModel(index)">删除</a-button>
          </template>
        </template>
      </a-table>

      <div class="form-actions">
        <a-button @click="loadConfig">重置</a-button>
        <a-button :loading="testing" @click="testConfig">测试连接</a-button>
        <a-button type="primary" :loading="saving" @click="saveConfig">保存配置</a-button>
      </div>
    </a-form>
  </div>
</template>

<script setup>
import { computed, getCurrentInstance, onMounted, reactive, ref } from 'vue'
import { getAgentOllamaConfig, testAgentOllamaConfig, updateAgentOllamaConfig } from '@/api/agent/config'

const { proxy } = getCurrentInstance()
const formRef = ref()
const saving = ref(false)
const testing = ref(false)

const form = reactive({
  baseUrl: 'http://localhost:11434',
  defaultModel: 'qwen2.5:7b',
  timeout: 90,
  models: []
})

const rules = {
  baseUrl: [{ required: true, message: 'Ollama Base URL 不能为空', trigger: 'blur' }],
  defaultModel: [{ required: true, message: '默认模型不能为空', trigger: 'change' }],
  timeout: [{ required: true, message: '超时时间不能为空', trigger: 'change' }]
}

const columns = [
  { title: '显示名称', dataIndex: 'label', key: 'label', width: 220 },
  { title: '模型标识', dataIndex: 'value', key: 'value', width: 220 },
  { title: '说明', dataIndex: 'description', key: 'description' },
  { title: '默认', dataIndex: 'default', key: 'default', width: 80, align: 'center' },
  { title: '操作', key: 'operation', width: 90, align: 'center' }
]

const modelSelectOptions = computed(() =>
  form.models
    .filter((item) => item.value)
    .map((item) => ({ label: item.label || item.value, value: item.value }))
)

function applyConfig(data) {
  form.baseUrl = data?.baseUrl || 'http://localhost:11434'
  form.defaultModel = data?.defaultModel || 'qwen2.5:7b'
  form.timeout = data?.timeout || 90
  form.models = Array.isArray(data?.models) ? data.models.map((item) => ({ ...item })) : []
  if (!form.models.length) {
    addModel()
  }
  syncDefaultModel()
}

async function loadConfig() {
  const response = await getAgentOllamaConfig()
  applyConfig(response.data || response)
}

function addModel() {
  form.models.push({ label: '', value: '', description: '', default: form.models.length === 0 })
}

function removeModel(index) {
  form.models.splice(index, 1)
  if (!form.models.length) {
    addModel()
  }
  syncDefaultModel()
}

function setDefault(index) {
  form.models.forEach((item, itemIndex) => {
    item.default = itemIndex === index
  })
  form.defaultModel = form.models[index]?.value || form.defaultModel
}

function syncDefaultModel() {
  const defaultItem = form.models.find((item) => item.default && item.value)
  if (defaultItem) {
    form.defaultModel = defaultItem.value
  } else if (!form.models.some((item) => item.value === form.defaultModel)) {
    const first = form.models.find((item) => item.value)
    form.defaultModel = first?.value || ''
    if (first) {
      first.default = true
    }
  }
}

async function saveConfig() {
  await formRef.value?.validate()
  syncDefaultModel()
  saving.value = true
  try {
    const payload = {
      baseUrl: form.baseUrl,
      defaultModel: form.defaultModel,
      timeout: form.timeout,
      models: form.models.filter((item) => item.value)
    }
    const response = await updateAgentOllamaConfig(payload)
    applyConfig(response.data || payload)
    proxy?.$modal?.msgSuccess?.('保存成功')
  } finally {
    saving.value = false
  }
}

async function testConfig() {
  await formRef.value?.validate()
  syncDefaultModel()
  testing.value = true
  try {
    await testAgentOllamaConfig({
      baseUrl: form.baseUrl,
      defaultModel: form.defaultModel,
      timeout: form.timeout,
      models: form.models.filter((item) => item.value)
    })
    proxy?.$modal?.msgSuccess?.('连接成功')
  } finally {
    testing.value = false
  }
}

onMounted(loadConfig)
</script>

<style scoped>
.agent-config-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.config-form {
  max-width: 1120px;
}

.section-head,
.form-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin: 8px 0 12px;
}

.form-actions {
  justify-content: flex-end;
  gap: 8px;
  margin-top: 16px;
}
</style>
