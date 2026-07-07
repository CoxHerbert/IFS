<template>
  <div class="app-container agent-page">
    <div class="session-panel">
      <div class="panel-head">
        <strong>IFS 智能助手</strong>
        <a-button type="primary" size="small" @click="handleCreateSession">新建</a-button>
      </div>

      <div class="model-select-wrap">
        <span>当前模型</span>
        <a-select v-model:value="selectedModel" size="small" style="width: 100%">
          <a-select-option v-for="item in models" :key="item.value" :value="item.value" :label="item.label">
            <span>{{ item.label }}</span>
            <small class="model-option-desc">{{ item.description }}</small>
          </a-select-option>
        </a-select>
      </div>

      <div class="session-list">
        <div v-for="item in sessions" :key="item.id" class="session-item"
          :class="{ active: item.id === activeSessionId }" @click="openSession(item.id)">
          <div class="session-main">
            <a-input v-if="editingSessionId === item.id" ref="sessionTitleInputRef" v-model:value="editingTitle"
              size="small" :maxlength="80" class="session-title-input" @click.stop
              @pressEnter="submitRenameSession(item)" @blur="submitRenameSession(item)"
              @keydown.esc.stop.prevent="cancelRenameSession" />
            <strong v-else class="session-title" @click.stop="startRenameSession(item)">
              {{ item.title }}
            </strong>
            <span>{{ item.updatedAt || item.modelName }}</span>
          </div>

          <a-popconfirm title="确定删除这个对话吗？" ok-text="确定" cancel-text="取消" @confirm="handleDeleteSession(item.id)">
            <a-button type="link" size="small" class="delete-btn" @click.stop>
              删除
            </a-button>
          </a-popconfirm>
        </div>
      </div>
    </div>

    <div class="chat-panel">
      <div ref="messageListRef" class="message-list">
        <a-empty v-if="!messages.length" description="开始一段对话" />

        <div v-for="item in messages" :key="item.id" class="message" :class="item.role">
          <div class="bubble">
            <template v-if="item.role === 'assistant' && item.blockResult">
              <div class="result-title">
                <strong>{{ item.blockResult.title }}</strong>
                <span>{{ item.blockResult.summary }}</span>
              </div>

              <div v-for="(block, index) in item.blockResult.blocks" :key="block.type + index" class="block">
                <div v-if="block.type === 'metrics'" class="metric-grid">
                  <div v-for="metric in block.items || []" :key="metric.label" class="metric">
                    <small>{{ metric.label }}</small>
                    <strong>{{ metric.value }}</strong>
                  </div>
                </div>

                <a-table v-else-if="block.type === 'table'" :data-source="block.data || []"
                  :columns="getBlockTableColumns(block)" :pagination="false"
                  :row-key="(record, rowIndex) => record.id || record.key || rowIndex" bordered size="small" />

                <a-alert v-else-if="block.type === 'error'" type="error" :message="block.title || '错误'"
                  :description="block.content" show-icon />

                <div v-else-if="block.type === 'form'" class="dynamic-form">
                  <strong class="block-title">{{ block.title }}</strong>

                  <a-form :model="getFormState(block)" :label-col="{ style: { width: '96px' } }" class="form-body">
                    <a-form-item v-for="field in block.fields || []" :key="field.field" :label="field.label"
                      :required="field.required">
                      <a-input v-if="field.component === 'input'" v-model:value="getFormState(block)[field.field]"
                        :placeholder="field.placeholder" />

                      <a-textarea v-else-if="field.component === 'textarea'"
                        v-model:value="getFormState(block)[field.field]" :rows="3" :placeholder="field.placeholder" />

                      <a-input-number v-else-if="field.component === 'number'"
                        v-model:value="getFormState(block)[field.field]" :controls="false" style="width: 100%" />

                      <a-select v-else-if="field.component === 'select'"
                        v-model:value="getFormState(block)[field.field]" :placeholder="field.placeholder || '请选择'"
                        style="width: 100%">
                        <a-select-option v-for="option in field.options || []" :key="String(option.value)"
                          :value="option.value">
                          {{ option.label }}
                        </a-select-option>
                      </a-select>

                      <a-date-picker v-else-if="field.component === 'date'"
                        v-model:value="getFormState(block)[field.field]" value-format="YYYY-MM-DD" placeholder="请选择日期"
                        style="width: 100%" />

                      <a-upload v-else-if="field.component === 'upload'" :before-upload="() => false" :max-count="1">
                        <a-button>选择文件</a-button>
                      </a-upload>

                      <a-input v-else v-model:value="getFormState(block)[field.field]"
                        :placeholder="field.placeholder" />
                    </a-form-item>

                    <a-button type="primary" :loading="formSubmitting" @click="handleSubmitForm(block)">
                      提交
                    </a-button>
                  </a-form>
                </div>

                <div v-else-if="block.type === 'action'" class="action-block">
                  <a-button type="primary" :loading="actionExecuting" @click="handleExecuteAction(block)">
                    {{ block.label || '执行操作' }}
                  </a-button>
                </div>

                <a-button v-else-if="block.type === 'file'" type="primary" :disabled="!block.url">
                  <a :href="block.url" download>{{ block.name || '下载文件' }}</a>
                </a-button>

                <pre v-else>{{ block.content }}</pre>
              </div>
            </template>

            <pre v-else>{{ item.content }}</pre>
          </div>
        </div>
      </div>

      <div class="composer" :class="{ dragging: isDragging }" @dragenter.prevent="handleDragEnter"
        @dragover.prevent="handleDragEnter" @dragleave.prevent="handleDragLeave" @drop.prevent="handleDrop">
        <input ref="fileInputRef" type="file" accept=".xlsx,.xls,.csv" class="hidden-input"
          @change="handleFileChange" />

        <textarea v-model="input" class="composer-textarea" :placeholder="composerPlaceholder" rows="3"
          @keydown.enter="handleEnter" />

        <div class="composer-toolbar">
          <div class="composer-tools">
            <button type="button" class="icon-button attach-button" :class="{ loading: uploading }" @click="pickFile">
              <span class="paperclip-icon" />
            </button>

            <button type="button" class="send-button" :disabled="sending" @click="handleSend">
              <span class="send-arrow" />
            </button>
          </div>
        </div>

        <span v-if="isDragging" class="drop-hint">松开后上传并分析文件</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { nextTick, onMounted, ref } from 'vue'
import { message } from 'ant-design-vue'
import {
  createChatSession,
  deleteChatSession,
  listAgentModels,
  listChatMessages,
  listChatSessions,
  sendChatMessage,
  updateChatSessionTitle
} from '@/api/agent/chat'
import { executeAgentAction, submitAgentForm } from '@/api/agent/form'
import { analyzeShipmentInChat } from '@/api/agent/shipment'

type AnyObject = Record<string, any>

interface AgentModel {
  label: string
  value: string
  description?: string
  default?: boolean
}

interface ChatSession {
  id: string | number
  title: string
  modelName?: string
  updatedAt?: string
}

interface ChatMessage {
  id: string | number
  sessionId?: string | number
  role: 'user' | 'assistant' | string
  content: string
  blockResult?: AnyObject
}

const sessions = ref<ChatSession[]>([])
const messages = ref<ChatMessage[]>([])
const models = ref<AgentModel[]>([])
const selectedModel = ref('qwen2.5:7b')
const activeSessionId = ref<string | number>()
const input = ref('帮我计算 100*200*150cm，10箱，需要多少方')
const sending = ref(false)
const uploading = ref(false)
const formSubmitting = ref(false)
const actionExecuting = ref(false)
const isDragging = ref(false)
const messageListRef = ref<HTMLElement>()
const fileInputRef = ref<HTMLInputElement>()
const formStateMap = ref<Record<string, AnyObject>>({})
const editingSessionId = ref<string | number>()
const editingTitle = ref('')
const sessionTitleInputRef = ref()
const composerPlaceholder = ref('给 IFS 智能助手发送消息')

onMounted(async () => {
  await refreshModels()
  await refreshSessions()
  if (!sessions.value.length) {
    await handleCreateSession()
    return
  }
  await openSession(sessions.value[0].id)
})

async function refreshModels() {
  try {
    const response = await listAgentModels()
    models.value = unwrapData<AgentModel[]>(response, [])
    selectedModel.value = models.value.find(item => item.default)?.value || models.value[0]?.value || selectedModel.value
    composerPlaceholder.value = `给 ${models.value.find(item => item.value === selectedModel.value)?.label || '智能助手'} 发送消息`
  } catch (error) {
    models.value = [{ label: 'Qwen 2.5 7B', value: selectedModel.value, description: '默认模型', default: true }]
  }
}

async function refreshSessions() {
  const response = await listChatSessions()
  sessions.value = unwrapData<ChatSession[]>(response, [])
}

async function handleCreateSession() {
  const response = await createChatSession({ title: 'IFS 智能助手对话', modelName: selectedModel.value })
  const session = unwrapData<ChatSession>(response, {
    id: '',
    title: 'IFS 智能助手对话',
    modelName: selectedModel.value
  })
  await refreshSessions()
  await openSession(session.id)
}

async function openSession(sessionId: string | number) {
  activeSessionId.value = sessionId
  const session = sessions.value.find(item => item.id === sessionId)
  if (session && session.modelName) {
    selectedModel.value = session.modelName
  }
  composerPlaceholder.value = `给 ${models.value.find(item => item.value === selectedModel.value)?.label || '智能助手'} 发送消息`
  const response = await listChatMessages(sessionId)
  messages.value = unwrapData<ChatMessage[]>(response, [])
  await scrollToBottom()
}

async function handleDeleteSession(sessionId: string | number) {
  await deleteChatSession(sessionId)
  if (activeSessionId.value === sessionId) {
    activeSessionId.value = undefined
    messages.value = []
  }
  await refreshSessions()
  if (!activeSessionId.value && sessions.value.length) {
    await openSession(sessions.value[0].id)
  }
  message.success('Deleted')
}

function startRenameSession(session: ChatSession) {
  editingSessionId.value = session.id
  editingTitle.value = session.title
  nextTick(() => {
    const inputRef = Array.isArray(sessionTitleInputRef.value)
      ? sessionTitleInputRef.value[0]
      : sessionTitleInputRef.value

    inputRef?.focus?.()
  })
}

function cancelRenameSession() {
  editingSessionId.value = undefined
  editingTitle.value = ''
}

async function submitRenameSession(session: ChatSession) {
  if (editingSessionId.value !== session.id) {
    return
  }

  const title = editingTitle.value.trim()
  cancelRenameSession()
  if (!title || title === session.title) {
    return
  }

  try {
    await updateChatSessionTitle(session.id, title)
    await refreshSessions()
    message.success('Updated')
  } catch (error) {
    message.error((error as Error)?.message || 'Update failed')
  }
}

function handleEnter(event: KeyboardEvent) {
  if (event.shiftKey) return
  event.preventDefault()
  handleSend()
}

async function ensureSession() {
  if (!activeSessionId.value) {
    await handleCreateSession()
  }
  return activeSessionId.value
}

async function handleSend() {
  const text = input.value.trim()
  if (!text) {
    message.warning('请输入消息')
    return
  }

  const sessionId = await ensureSession()
  if (!sessionId) return

  messages.value.push({ id: Date.now(), sessionId, role: 'user', content: text })
  input.value = ''
  sending.value = true
  await scrollToBottom()

  try {
    const response = await sendChatMessage({ sessionId, message: text, modelName: selectedModel.value })
    const payload = unwrapData<AnyObject>(response, response as AnyObject)
    messages.value.push({
      id: payload.messageId,
      sessionId,
      role: 'assistant',
      content: payload.result.summary,
      blockResult: payload.result
    })
    await refreshSessions()
    await scrollToBottom()
  } catch (error) {
    message.error((error as Error).message || '发送失败')
  } finally {
    sending.value = false
  }
}

function pickFile() {
  fileInputRef.value?.click()
}

async function handleFileChange(event: Event) {
  const target = event.target as HTMLInputElement
  const file = target.files && target.files[0]
  if (file) await handleFile(file)
  target.value = ''
}

function handleDragEnter() {
  isDragging.value = true
}

function handleDragLeave(event: DragEvent) {
  const current = event.currentTarget as HTMLElement
  const related = event.relatedTarget as Node | null
  if (!related || !current.contains(related)) {
    isDragging.value = false
  }
}

async function handleDrop(event: DragEvent) {
  isDragging.value = false
  const file = event.dataTransfer && event.dataTransfer.files && event.dataTransfer.files[0]
  if (file) await handleFile(file)
}

async function handleFile(file: File) {
  if (!/\.(xlsx|xls|csv)$/i.test(file.name)) {
    message.warning('请选择 Excel 或 CSV 文件')
    return
  }

  const sessionId = await ensureSession()
  if (!sessionId) return

  uploading.value = true
  try {
    messages.value.push({
      id: Date.now(),
      sessionId,
      role: 'user',
      content: `上传文件：${file.name}，正在由服务端解析并分析。`
    })
    await scrollToBottom()

    const response = await analyzeShipmentInChat(sessionId, file, selectedModel.value)
    const payload = unwrapData<AnyObject>(response, response as AnyObject)
    messages.value.push({
      id: payload.messageId,
      sessionId,
      role: 'assistant',
      content: payload.result.summary,
      blockResult: payload.result
    })
    await refreshSessions()
    await scrollToBottom()
  } catch (error) {
    message.error((error as Error).message || '文件分析失败')
  } finally {
    uploading.value = false
  }
}

function getFormState(block: AnyObject) {
  const key = block.formCode || block.title || 'default'
  if (!formStateMap.value[key]) {
    formStateMap.value[key] = { ...(block.initialValues || {}) }
  }
  return formStateMap.value[key]
}

async function handleSubmitForm(block: AnyObject) {
  if (!activeSessionId.value) {
    message.warning('缺少会话 ID')
    return
  }

  const values = getFormState(block)
  const missing = (block.fields || []).find((field: AnyObject) => field.required && !values[field.field])
  if (missing) {
    message.warning('请填写' + missing.label)
    return
  }

  formSubmitting.value = true
  try {
    const response = await submitAgentForm({
      sessionId: activeSessionId.value,
      formCode: block.formCode,
      values
    })
    appendAgentResult(unwrapData<AnyObject>(response, response as AnyObject))
  } catch (error) {
    message.error((error as Error).message || '提交失败')
  } finally {
    formSubmitting.value = false
  }
}

async function handleExecuteAction(block: AnyObject) {
  if (!activeSessionId.value) {
    message.warning('缺少会话 ID')
    return
  }

  actionExecuting.value = true
  try {
    const response = await executeAgentAction({
      sessionId: activeSessionId.value,
      actionCode: block.actionCode,
      payload: block.payload || {}
    })
    appendAgentResult(unwrapData<AnyObject>(response, response as AnyObject))
  } catch (error) {
    message.error((error as Error).message || '执行失败')
  } finally {
    actionExecuting.value = false
  }
}

async function appendAgentResult(result: AnyObject) {
  if (!activeSessionId.value) return

  messages.value.push({
    id: Date.now(),
    sessionId: activeSessionId.value,
    role: 'assistant',
    content: result.summary,
    blockResult: result
  })
  await refreshSessions()
  await scrollToBottom()
}

function getBlockTableColumns(block: AnyObject) {
  return (block.columns || []).map((column: AnyObject) => ({
    title: column.label,
    dataIndex: column.field,
    key: column.field,
    ellipsis: true
  }))
}

function unwrapData<T>(response: any, fallback: T): T {
  if (response && Object.prototype.hasOwnProperty.call(response, 'data')) {
    return response.data as T
  }

  return (response || fallback) as T
}

async function scrollToBottom() {
  await nextTick()
  const wrap = messageListRef.value
  if (wrap) {
    wrap.scrollTop = wrap.scrollHeight
  }
}
</script>

<style scoped>
.agent-page {
  height: calc(100vh - 104px);
  display: grid;
  grid-template-columns: 300px minmax(0, 1fr);
  gap: 16px;
}

.session-panel,
.chat-panel {
  min-height: 0;
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
}

.panel-head {
  height: 56px;
  padding: 0 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #e5e7eb;
}

.model-select-wrap {
  padding: 10px 16px;
  border-bottom: 1px solid #e5e7eb;
  display: grid;
  gap: 6px;
}

.model-select-wrap span {
  color: #64748b;
  font-size: 12px;
}

.model-option-desc {
  display: block;
  color: #909399;
  font-size: 12px;
}

.session-list {
  height: calc(100% - 116px);
  overflow: auto;
}

.session-item {
  padding: 14px 16px;
  cursor: pointer;
  border-bottom: 1px solid #f1f5f9;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.session-item.active {
  background: #eef5ff;
}

.session-main {
  min-width: 0;
  flex: 1;
}

.session-title {
  display: block;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  cursor: text;
}

.session-title-input {
  width: 100%;
}

.session-main strong,
.session-main span {
  display: block;
}

.session-main span {
  margin-top: 6px;
  color: #64748b;
  font-size: 12px;
}

.delete-btn {
  flex-shrink: 0;
  color: #f56c6c;
}

.chat-panel {
  display: grid;
  grid-template-rows: minmax(0, 1fr) auto;
}

.message-list {
  padding: 18px;
  overflow: auto;
  min-height: 0;
}

.message {
  display: flex;
  margin-bottom: 14px;
}

.message.user {
  justify-content: flex-end;
}

.bubble {
  max-width: min(780px, 90%);
  padding: 14px;
  border-radius: 6px;
  background: #f8fafc;
}

.message.user .bubble {
  color: #fff;
  background: #409eff;
}

.composer {
  position: relative;
  display: grid;
  gap: 14px;
  margin: 14px;
  padding: 18px 18px 14px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 24px;
  background: linear-gradient(180deg, #363638 0%, #2e2e31 100%);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.04);
}

.composer.dragging {
  outline: 2px dashed rgba(108, 141, 255, 0.9);
  outline-offset: -10px;
}

.hidden-input {
  display: none;
}

.composer-textarea {
  width: 100%;
  min-height: 72px;
  resize: none;
  border: 0;
  outline: 0;
  padding: 0;
  background: transparent;
  color: #f4f7fb;
  font: inherit;
  line-height: 1.7;
}

.composer-textarea::placeholder {
  color: rgba(244, 247, 251, 0.46);
}

.composer-toolbar {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
}

.composer-tools {
  display: flex;
  align-items: center;
  gap: 10px;
}

.icon-button,
.send-button {
  border: 0;
  cursor: pointer;
}

.icon-button {
  width: 36px;
  height: 36px;
  display: grid;
  place-items: center;
  border-radius: 50%;
  background: transparent;
  color: #f4f7fb;
}

.icon-button.loading,
.icon-button:hover {
  background: rgba(255, 255, 255, 0.08);
}

.paperclip-icon {
  width: 13px;
  height: 13px;
  border: 2px solid currentColor;
  border-top-color: transparent;
  border-radius: 8px;
  transform: rotate(-35deg);
  position: relative;
}

.paperclip-icon::after {
  content: '';
  position: absolute;
  inset: 2px;
  border: 2px solid currentColor;
  border-top-color: transparent;
  border-radius: 8px;
}

.send-button {
  width: 40px;
  height: 40px;
  display: grid;
  place-items: center;
  border-radius: 50%;
  background: #4d69b8;
  color: #fff;
  box-shadow: 0 10px 20px rgba(77, 105, 184, 0.32);
}

.send-button:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.send-arrow {
  width: 12px;
  height: 12px;
  border-top: 2px solid currentColor;
  border-right: 2px solid currentColor;
  transform: rotate(-45deg) translate(-1px, 1px);
}

.send-arrow::after {
  content: '';
  display: block;
  width: 2px;
  height: 12px;
  margin-left: 4px;
  margin-top: -2px;
  background: currentColor;
  transform: rotate(45deg);
}

.drop-hint {
  position: absolute;
  inset: 8px;
  display: grid;
  place-items: center;
  pointer-events: none;
  color: #dce6ff;
  font-weight: 700;
}

.result-title,
.block {
  margin-bottom: 12px;
}

.block-title {
  display: block;
  margin-bottom: 12px;
}

.dynamic-form {
  padding: 12px;
  border-radius: 6px;
  background: #fff;
  border: 1px solid #e5e7eb;
}

.form-body {
  margin-top: 10px;
}

.action-block {
  padding: 12px;
  border-radius: 6px;
  background: #fff;
  border: 1px solid #e5e7eb;
}

.result-title span {
  display: block;
  margin-top: 6px;
  color: #64748b;
}

.metric-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 10px;
}

.metric {
  padding: 12px;
  border-radius: 6px;
  background: #fff;
  border: 1px solid #e5e7eb;
}

.metric small,
.metric strong {
  display: block;
}

pre {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: inherit;
  line-height: 1.7;
}
</style>
