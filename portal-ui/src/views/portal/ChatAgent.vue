<template>
  <main class="chat-page">
    <aside class="session-panel">
      <div class="session-head">
        <strong>IFS 智能助手</strong>
        <a-button type="primary" size="small" @click="handleCreateSession">新建</a-button>
      </div>
      <a-list :data-source="sessions" class="session-list">
        <template #renderItem="{ item }">
          <a-list-item class="session-item" :class="{ active: item.id === activeSessionId }" @click="openSession(item.id)">
            <a-list-item-meta :title="item.title" :description="item.updatedAt || item.modelName" />
            <template #actions>
              <a-button type="link" size="small" @click.stop="handleRenameSession(item)">重命名</a-button>
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
        <a-empty v-if="!messages.length" description="开始一段对话" />
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
        <input ref="fileInputRef" type="file" accept=".xlsx,.xls,.csv" class="hidden-input" @change="handleFileChange" />
        <a-button class="attach-button" :loading="uploading" @click="pickFile">选择文件</a-button>
        <a-textarea
          v-model:value="input"
          :auto-size="{ minRows: 2, maxRows: 5 }"
          placeholder="输入消息，或拖入 Excel/CSV 文件"
          @keydown.enter="handleEnter"
        />
        <a-button type="primary" :loading="sending" @click="handleSend">发送</a-button>
        <span v-if="isDragging" class="drop-hint">松开后上传并分析文件</span>
      </div>
    </section>
  </main>
</template>

<script setup lang="ts">
import { nextTick, onMounted, ref } from 'vue'
import { message as antMessage } from 'ant-design-vue'
import AgentResultRenderer from '@/components/agent-renderer/AgentResultRenderer.vue'
import {
  analyzeShipmentInChat,
  createChatSession,
  deleteChatSession,
  listChatMessages,
  listChatSessions,
  sendChatMessage,
  updateChatSessionTitle,
} from '@/api/chat'
import type { ChatMessage, ChatSession } from '@/types/agent'

const sessions = ref<ChatSession[]>([])
const messages = ref<ChatMessage[]>([])
const activeSessionId = ref<number>()
const input = ref('帮我计算 100*200*150cm，10箱，需要多少方')
const sending = ref(false)
const uploading = ref(false)
const isDragging = ref(false)
const messageListRef = ref<HTMLElement>()
const fileInputRef = ref<HTMLInputElement>()

onMounted(async () => {
  await refreshSessions()
  if (!sessions.value.length) {
    await handleCreateSession()
    return
  }
  await openSession(sessions.value[0].id)
})

async function refreshSessions() {
  sessions.value = await listChatSessions()
}

async function handleCreateSession() {
  const session = await createChatSession({ title: 'IFS 智能助手对话', modelName: 'qwen2.5:7b' })
  await refreshSessions()
  await openSession(session.id)
}

async function openSession(sessionId: number) {
  activeSessionId.value = sessionId
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
  antMessage.success('已删除对话')
}

async function handleRenameSession(session: ChatSession) {
  const title = window.prompt('请输入新的对话名称', session.title)?.trim()
  if (!title || title === session.title) {
    return
  }
  try {
    await updateChatSessionTitle(session.id, title)
    await refreshSessions()
    antMessage.success('对话名称已更新')
  } catch (error) {
    antMessage.error(error instanceof Error ? error.message : '更新失败')
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
    antMessage.warning('请输入消息')
    return
  }
  const sessionId = await ensureSession()
  if (!sessionId) return

  messages.value.push({ id: Date.now(), sessionId, role: 'user', content: text, createdAt: '' })
  input.value = ''
  sending.value = true
  await scrollToBottom()

  try {
    const response = await sendChatMessage({ sessionId, message: text, modelName: 'qwen2.5:7b' })
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
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (file) {
    await handleFile(file)
  }
  input.value = ''
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
  if (!sessionId) return

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

    const response = await analyzeShipmentInChat({ sessionId, file, modelName: 'qwen2.5:7b' })
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

async function appendAgentResult(result: import('@/types/agent').AgentResult) {
  if (!activeSessionId.value) return
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
.chat-page {
  width: min(1240px, calc(100% - 40px));
  height: calc(100vh - 132px);
  margin: 28px auto;
  display: grid;
  grid-template-columns: 300px minmax(0, 1fr);
  gap: 18px;
}

.session-panel,
.chat-panel {
  min-height: 0;
  border: 1px solid rgba(15, 23, 42, 0.08);
  background: #fff;
  border-radius: 8px;
}

.session-panel {
  overflow: hidden;
}

.session-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid rgba(15, 23, 42, 0.08);
}

.session-list {
  height: calc(100% - 58px);
  overflow: auto;
}

.session-item {
  cursor: pointer;
  padding-inline: 16px;
}

.session-item.active {
  background: #eaf2ff;
}

.chat-panel {
  display: grid;
  grid-template-rows: minmax(0, 1fr) auto;
}

.message-list {
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
  max-width: min(760px, 88%);
  border-radius: 8px;
  padding: 14px 16px;
  background: #f8fafc;
  border: 1px solid rgba(148, 163, 184, 0.18);
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
  padding: 16px;
  border-top: 1px solid rgba(15, 23, 42, 0.08);
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
  pointer-events: none;
  color: #1677ff;
  font-weight: 700;
}

@media (max-width: 820px) {
  .chat-page {
    height: auto;
    grid-template-columns: 1fr;
  }

  .session-panel {
    height: 260px;
  }

  .chat-panel {
    min-height: 620px;
  }

  .composer {
    grid-template-columns: 1fr;
  }
}
</style>
