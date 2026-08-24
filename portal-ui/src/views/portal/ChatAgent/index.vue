<template>
  <main class="chat-page">
    <aside class="session-panel">
      <div class="session-head">
        <strong>IFS 智能助手</strong>
        <a-button type="primary" size="small" @click="handleCreateSession">新建</a-button>
      </div>
      <a-list :data-source="sessions" class="session-list">
        <template #renderItem="{ item }">
          <a-list-item class="session-item" :class="{ active: item.id === activeSessionId }"
            @click="openSession(item.id)">
            <div class="session-main">
              <a-input v-if="editingSessionId === item.id" ref="sessionTitleInputRef" v-model:value="editingTitle"
                size="small" class="session-title-input" :maxlength="80" @click.stop
                @press-enter="submitRenameSession(item)" @blur="submitRenameSession(item)"
                @keydown.esc.stop.prevent="cancelRenameSession" />
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
        <a-empty v-if="!messages.length" description="开始一段对话" />
        <article v-for="message in messages" :key="message.id" class="message" :class="message.role">
          <div class="bubble">
            <template v-if="message.role === 'assistant' && message.blockResult">
              <AgentResultRenderer :result="message.blockResult" :session-id="activeSessionId"
                @submitted="appendAgentResult" />
            </template>
            <pre v-else>{{ message.content }}</pre>
          </div>
        </article>
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
const input = ref('')
const sending = ref(false)
const uploading = ref(false)
const isDragging = ref(false)
const editingSessionId = ref<number>()
const editingTitle = ref('')
const messageListRef = ref<HTMLElement>()
const fileInputRef = ref<HTMLInputElement>()
const sessionTitleInputRef = ref<{ focus: () => void }>()
const composerPlaceholder = '给 IFS 智能助手发送消息'

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
  const session = await createChatSession({ title: 'IFS 智能助手对话' })
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
  antMessage.success('对话已删除')
}

function startRenameSession(session: ChatSession) {
  editingSessionId.value = session.id
  editingTitle.value = session.title
  nextTick(() => sessionTitleInputRef.value?.focus())
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
    antMessage.error(error instanceof Error ? error.message : '鏇存柊澶辫触')
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
    const response = await sendChatMessage({ sessionId, message: text })
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

    const response = await analyzeShipmentInChat({ sessionId, file })
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
    antMessage.error(error instanceof Error ? error.message : '鏂囦欢鍒嗘瀽澶辫触')
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
  gap: 14px;
  margin: 16px;
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
  inset: 10px;
  display: grid;
  place-items: center;
  pointer-events: none;
  color: #dce6ff;
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
    margin: 12px;
    padding: 16px;
  }

  .composer-toolbar {
    flex-direction: column;
    align-items: stretch;
  }

  .composer-tools {
    justify-content: flex-end;
  }
}
</style>
