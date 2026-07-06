<template>
  <main class="agent-chat-page">
    <aside class="session-panel">
      <div class="session-head">
        <div>
          <strong>Agent 对话</strong>
          <span>客户专属会话</span>
        </div>
        <a-button type="primary" size="small" @click="handleCreateSession">新建</a-button>
      </div>
      <div class="model-select-wrap">
        <span>当前模型</span>
        <a-select v-model:value="selectedModel" size="small" class="model-select" :options="modelSelectOptions" />
      </div>

      <a-list :data-source="sessions" class="session-list">
        <template #renderItem="{ item }">
          <a-list-item
            class="session-item"
            :class="{ active: item.id === activeSessionId }"
            @click="openSession(item.id)"
          >
            <div class="session-main">
              <a-input
                v-if="editingSessionId === item.id"
                v-model:value="editingTitle"
                size="small"
                class="session-title-input"
                :maxlength="80"
                @click.stop
                @press-enter="submitRenameSession(item)"
                @blur="submitRenameSession(item)"
                @keydown.esc.stop.prevent="cancelRenameSession"
              />
              <strong v-else class="session-title" @click.stop="startRenameSession(item)">{{ item.title }}</strong>
              <span>{{ item.updatedAt || item.modelName }}</span>
            </div>
            <template #actions>
              <a-popconfirm title="确定删除这个对话吗？" @confirm="handleDeleteSession(item.id)">
                <a-button type="link" danger size="small" @click.stop>删除</a-button>
              </a-popconfirm>
            </template>
          </a-list-item>
        </template>
      </a-list>
    </aside>

    <section class="chat-panel">
      <div ref="messageListRef" class="message-list">
        <a-empty v-if="!messages.length" description="开始一段客户 Agent 对话" />

        <article v-for="message in messages" :key="message.id" class="message" :class="message.role">
          <div class="bubble">
            <template v-if="message.role === 'assistant' && message.blockResult">
              <AgentResultRenderer
                :result="message.blockResult"
                :session-id="activeSessionId"
                @submitted="appendAgentResult"
              />
            </template>
            <pre v-else>{{ message.content }}</pre>
          </div>
        </article>
      </div>

      <div
        class="composer"
        :class="{ dragging: isDragging }"
        @dragenter.prevent="handleDragEnter"
        @dragover.prevent="handleDragEnter"
        @dragleave.prevent="handleDragLeave"
        @drop.prevent="handleDrop"
      >
        <input
          ref="fileInputRef"
          type="file"
          accept=".xlsx,.xls,.csv"
          class="hidden-input"
          @change="handleFileChange"
        />
        <a-button class="attach-button" :loading="uploading" @click="pickFile">选择文件</a-button>
        <a-textarea
          v-model:value="input"
          :auto-size="{ minRows: 2, maxRows: 5 }"
          placeholder="输入问题，或拖入 Excel/CSV 出货计划"
          @keydown.enter="handleEnter"
        />
        <a-button type="primary" :loading="sending" @click="handleSend">发送</a-button>
        <span v-if="isDragging" class="drop-hint">松开后上传并分析文件</span>
      </div>
    </section>
  </main>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, ref } from 'vue'
import { message as antMessage } from 'ant-design-vue'
import AgentResultRenderer from '@/components/agent-renderer/AgentResultRenderer.vue'
import {
  analyzeShipmentInChat,
  createChatSession,
  deleteChatSession,
  listAgentModels,
  listChatMessages,
  listChatSessions,
  sendChatMessage,
  updateChatSessionTitle,
  type AgentModelOption,
} from '@/api/chat'
import type { AgentResult, ChatMessage, ChatSession } from '@/types/agent'

const sessions = ref<ChatSession[]>([])
const messages = ref<ChatMessage[]>([])
const models = ref<AgentModelOption[]>([])
const selectedModel = ref('qwen2.5:7b')
const activeSessionId = ref<number>()
const input = ref('')
const sending = ref(false)
const uploading = ref(false)
const isDragging = ref(false)
const editingSessionId = ref<number>()
const editingTitle = ref('')
const messageListRef = ref<HTMLElement>()
const fileInputRef = ref<HTMLInputElement>()

onMounted(async () => {
  await refreshModels()
  await refreshSessions()
  if (sessions.value.length) {
    await openSession(sessions.value[0].id)
  }
})

const modelSelectOptions = computed(() =>
  models.value.map((item) => ({
    label: item.label,
    value: item.value,
    title: item.description,
  })),
)

async function refreshModels() {
  try {
    models.value = await listAgentModels()
    selectedModel.value = models.value.find((item) => item.default)?.value || models.value[0]?.value || selectedModel.value
  } catch (_error) {
    models.value = [{ label: 'Qwen 2.5 7B', value: selectedModel.value, description: '默认模型', default: true }]
  }
}

async function refreshSessions() {
  sessions.value = await listChatSessions()
}

async function handleCreateSession() {
  const session = await createChatSession({ title: '客户 Agent 对话', modelName: selectedModel.value })
  await refreshSessions()
  await openSession(session.id)
}

async function openSession(sessionId: number) {
  activeSessionId.value = sessionId
  const session = sessions.value.find((item) => item.id === sessionId)
  if (session?.modelName) {
    selectedModel.value = session.modelName
  }
  messages.value = await listChatMessages(sessionId)
  await scrollToBottom()
}

async function handleDeleteSession(sessionId: number) {
  await deleteChatSession(sessionId)
  if (activeSessionId.value === sessionId) {
    activeSessionId.value = undefined
    messages.value = []
  }
  await refreshSessions()
  if (!activeSessionId.value && sessions.value.length) {
    await openSession(sessions.value[0].id)
  }
  antMessage.success('对话已删除')
}

function startRenameSession(session: ChatSession) {
  editingSessionId.value = session.id
  editingTitle.value = session.title
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
    antMessage.success('对话标题已更新')
  } catch (error) {
    antMessage.error(error instanceof Error ? error.message : '更新失败')
  }
}

function handleEnter(event: KeyboardEvent) {
  if (event.shiftKey) {
    return
  }
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
    antMessage.warning('请输入消息')
    return
  }
  const sessionId = await ensureSession()
  if (!sessionId) {
    return
  }

  messages.value.push({ id: Date.now(), sessionId, role: 'user', content: text, createdAt: '' })
  input.value = ''
  sending.value = true
  await scrollToBottom()

  try {
    const response = await sendChatMessage({ sessionId, message: text, modelName: selectedModel.value })
    messages.value.push({
      id: response.messageId,
      sessionId,
      role: 'assistant',
      content: response.result.summary,
      blockResult: response.result,
      createdAt: '',
    })
    await refreshSessions()
    await scrollToBottom()
  } catch (error) {
    antMessage.error(error instanceof Error ? error.message : '发送失败')
  } finally {
    sending.value = false
  }
}

function pickFile() {
  fileInputRef.value?.click()
}

async function handleFileChange(event: Event) {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    await handleFile(file)
  }
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
  const file = event.dataTransfer?.files?.[0]
  if (file) {
    await handleFile(file)
  }
}

async function handleFile(file: File) {
  if (!/\.(xlsx|xls|csv)$/i.test(file.name)) {
    antMessage.warning('请选择 Excel 或 CSV 文件')
    return
  }
  const sessionId = await ensureSession()
  if (!sessionId) {
    return
  }

  uploading.value = true
  try {
    messages.value.push({
      id: Date.now(),
      sessionId,
      role: 'user',
      content: `上传文件：${file.name}，正在由服务端解析并分析。`,
      createdAt: '',
    })
    await scrollToBottom()

    const response = await analyzeShipmentInChat({ sessionId, file, modelName: selectedModel.value })
    messages.value.push({
      id: response.messageId,
      sessionId,
      role: 'assistant',
      content: response.result.summary,
      blockResult: response.result,
      createdAt: '',
    })
    await refreshSessions()
    await scrollToBottom()
  } catch (error) {
    antMessage.error(error instanceof Error ? error.message : '文件分析失败')
  } finally {
    uploading.value = false
  }
}

async function appendAgentResult(result: AgentResult) {
  if (!activeSessionId.value) {
    return
  }
  messages.value.push({
    id: Date.now(),
    sessionId: activeSessionId.value,
    role: 'assistant',
    content: result.summary,
    blockResult: result,
    createdAt: '',
  })
  await refreshSessions()
  await scrollToBottom()
}

async function scrollToBottom() {
  await nextTick()
  if (messageListRef.value) {
    messageListRef.value.scrollTop = messageListRef.value.scrollHeight
  }
}
</script>

<style scoped>
.agent-chat-page {
  height: 100%;
  min-height: 0;
  display: grid;
  grid-template-columns: 280px minmax(0, 1fr);
  gap: 16px;
}

.session-panel,
.chat-panel {
  min-height: 0;
  border: 1px solid rgba(15, 23, 42, 0.08);
  border-radius: 8px;
  background: #fff;
  box-shadow: 0 12px 30px rgba(15, 23, 42, 0.05);
}

.session-panel {
  overflow: hidden;
}

.session-head {
  height: 64px;
  padding: 12px 14px;
  border-bottom: 1px solid rgba(15, 23, 42, 0.08);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.session-head strong,
.session-head span {
  display: block;
}

.session-head strong {
  color: #0f172a;
  font-size: 15px;
}

.session-head span {
  margin-top: 4px;
  color: #64748b;
  font-size: 12px;
}

.model-select-wrap {
  padding: 10px 14px;
  border-bottom: 1px solid rgba(15, 23, 42, 0.08);
  display: grid;
  gap: 6px;
}

.model-select-wrap span {
  color: #64748b;
  font-size: 12px;
}

.model-select {
  width: 100%;
}

.session-list {
  height: calc(100% - 124px);
  overflow: auto;
}

.session-item {
  cursor: pointer;
  padding-inline: 14px;
}

.session-item.active {
  background: #eaf6ff;
}

.session-main {
  min-width: 0;
  flex: 1;
}

.session-title {
  display: block;
  max-width: 100%;
  color: #0f172a;
  font-size: 14px;
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  cursor: text;
}

.session-main span {
  display: block;
  margin-top: 6px;
  color: #64748b;
  font-size: 12px;
}

.session-title-input {
  width: 100%;
}

.chat-panel {
  display: grid;
  grid-template-rows: minmax(0, 1fr) auto;
  overflow: hidden;
}

.message-list {
  min-height: 0;
  overflow: auto;
  padding: 22px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.message {
  display: flex;
}

.message.user {
  justify-content: flex-end;
}

.bubble {
  max-width: min(780px, 88%);
  border: 1px solid rgba(148, 163, 184, 0.18);
  border-radius: 8px;
  padding: 14px 16px;
  background: #f8fafc;
}

.message.user .bubble {
  color: #fff;
  background: #1677ff;
  border-color: #1677ff;
}

pre {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: inherit;
  line-height: 1.7;
}

.composer {
  position: relative;
  display: grid;
  grid-template-columns: auto minmax(0, 1fr) auto;
  gap: 12px;
  padding: 14px;
  border-top: 1px solid rgba(15, 23, 42, 0.08);
  background: #fff;
}

.composer.dragging {
  background: #eef6ff;
  outline: 2px dashed #1677ff;
  outline-offset: -8px;
}

.hidden-input {
  display: none;
}

.attach-button {
  align-self: end;
}

.drop-hint {
  position: absolute;
  inset: 10px;
  display: grid;
  place-items: center;
  color: #1677ff;
  font-weight: 700;
  pointer-events: none;
}

@media (max-width: 920px) {
  .agent-chat-page {
    height: auto;
    min-height: 720px;
    grid-template-columns: 1fr;
  }

  .session-panel {
    height: 240px;
  }

  .composer {
    grid-template-columns: 1fr;
  }
}
</style>
