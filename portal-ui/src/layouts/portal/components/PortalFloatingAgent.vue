<template>
  <div class="floating-agent" :class="{ expanded: isOpen }" :style="floatingStyle">
    <transition name="agent-panel">
      <section v-if="isOpen" class="agent-panel">
        <header class="agent-panel-head" @mousedown="startDrag" @touchstart.passive="startTouchDrag">
          <div>
            <strong>智能助手</strong>
            <span>在线解答航线、报价和出货问题</span>
          </div>
          <button type="button" class="panel-close" @click="isOpen = false">×</button>
        </header>

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
          <input
            ref="fileInputRef"
            type="file"
            accept=".xlsx,.xls,.csv"
            class="hidden-input"
            @change="handleFileChange"
          />
          <textarea
            v-model="input"
            class="composer-textarea"
            :placeholder="composerPlaceholder"
            rows="3"
            @keydown.enter="handleEnter"
          />
          <div class="composer-toolbar">
            <div class="composer-tools">
              <button
                type="button"
                class="icon-button"
                :class="{ loading: uploading }"
                @click="pickFile"
              >
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
    </transition>

    <button
      type="button"
      class="floating-trigger"
      @click="toggleOpen"
      @mousedown="startDrag"
      @touchstart.passive="startTouchDrag"
    >
      <span class="trigger-dot" />
      <span>{{ isOpen ? '收起助手' : '智能助手' }}</span>
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref } from 'vue'
import { message as antMessage } from 'ant-design-vue'
import AgentResultRenderer from '@/components/agent-renderer/AgentResultRenderer.vue'
import {
  analyzeShipmentInChat,
  createChatSession,
  listChatMessages,
  listChatSessions,
  sendChatMessage,
} from '@/api/chat'
import type { AgentResult, ChatMessage, ChatSession } from '@/types/agent'

const isOpen = ref(false)
const sessions = ref<ChatSession[]>([])
const messages = ref<ChatMessage[]>([])
const activeSessionId = ref<number>()
const input = ref('')
const sending = ref(false)
const uploading = ref(false)
const isDragging = ref(false)
const messageListRef = ref<HTMLElement>()
const fileInputRef = ref<HTMLInputElement>()
const composerPlaceholder = '给 IFS 智能助手发送消息'
const offsetX = ref(24)
const offsetY = ref(24)
const dragState = ref<{ startX: number; startY: number; baseX: number; baseY: number }>()

const floatingStyle = computed(() => ({
  right: `${offsetX.value}px`,
  bottom: `${offsetY.value}px`,
}))

onMounted(() => {
  window.addEventListener('portal-agent:open', handleOpenEvent as EventListener)
  window.addEventListener('mousemove', handleDragMove)
  window.addEventListener('mouseup', stopDrag)
  window.addEventListener('touchmove', handleTouchDragMove, { passive: false })
  window.addEventListener('touchend', stopDrag)
})

onBeforeUnmount(() => {
  window.removeEventListener('portal-agent:open', handleOpenEvent as EventListener)
  window.removeEventListener('mousemove', handleDragMove)
  window.removeEventListener('mouseup', stopDrag)
  window.removeEventListener('touchmove', handleTouchDragMove)
  window.removeEventListener('touchend', stopDrag)
})

function handleOpenEvent() {
  openPanel()
}

async function openPanel() {
  isOpen.value = true
  if (!activeSessionId.value) {
    await bootstrapSession()
  } else {
    await scrollToBottom()
  }
}

async function toggleOpen() {
  if (isOpen.value) {
    isOpen.value = false
    return
  }
  await openPanel()
}

async function bootstrapSession() {
  await refreshSessions()
  if (!sessions.value.length) {
    await handleCreateSession()
    return
  }
  await openSession(sessions.value[0].id)
}

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
      content: `上传文件：${file.name}，正在分析。`,
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

function startDrag(event: MouseEvent) {
  if (shouldSkipDrag(event.target as HTMLElement | null)) {
    return
  }
  dragState.value = {
    startX: event.clientX,
    startY: event.clientY,
    baseX: offsetX.value,
    baseY: offsetY.value,
  }
}

function startTouchDrag(event: TouchEvent) {
  if (shouldSkipDrag(event.target as HTMLElement | null)) {
    return
  }
  const touch = event.touches[0]
  if (!touch) {
    return
  }
  dragState.value = {
    startX: touch.clientX,
    startY: touch.clientY,
    baseX: offsetX.value,
    baseY: offsetY.value,
  }
}

function handleDragMove(event: MouseEvent) {
  if (!dragState.value) {
    return
  }
  applyDragPosition(event.clientX, event.clientY)
}

function handleTouchDragMove(event: TouchEvent) {
  if (!dragState.value) {
    return
  }
  const touch = event.touches[0]
  if (!touch) {
    return
  }
  event.preventDefault()
  applyDragPosition(touch.clientX, touch.clientY)
}

function applyDragPosition(clientX: number, clientY: number) {
  const state = dragState.value
  if (!state) {
    return
  }
  const nextX = state.baseX - (clientX - state.startX)
  const nextY = state.baseY - (clientY - state.startY)
  offsetX.value = clamp(nextX, 12, Math.max(12, window.innerWidth - 96))
  offsetY.value = clamp(nextY, 12, Math.max(12, window.innerHeight - 72))
}

function stopDrag() {
  dragState.value = undefined
}

function clamp(value: number, min: number, max: number) {
  return Math.min(Math.max(value, min), max)
}

function shouldSkipDrag(target: HTMLElement | null) {
  return !!target?.closest('.panel-close')
}
</script>

<style scoped>
.floating-agent {
  position: fixed;
  right: 24px;
  bottom: 24px;
  z-index: 40;
  display: grid;
  justify-items: end;
  gap: 14px;
}

.agent-panel {
  width: min(560px, calc(100vw - 32px));
  height: min(760px, calc(100vh - 108px));
  max-width: calc(100vw - 32px);
  max-height: calc(100vh - 88px);
  display: grid;
  grid-template-rows: auto minmax(0, 1fr) auto;
  border-radius: 28px;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.96);
  border: 1px solid rgba(15, 23, 42, 0.08);
  box-shadow: 0 28px 60px rgba(15, 23, 42, 0.22);
  backdrop-filter: blur(18px);
}

.agent-panel-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  padding: 18px 18px 12px;
  cursor: move;
  user-select: none;
}

.agent-panel-head > div {
  min-width: 0;
}

.agent-panel-head strong,
.agent-panel-head span {
  display: block;
}

.agent-panel-head strong {
  color: #0f172a;
  font-size: 16px;
}

.agent-panel-head span {
  margin-top: 6px;
  color: #64748b;
  font-size: 12px;
  line-height: 1.4;
}

.panel-close {
  flex: 0 0 auto;
  width: 32px;
  height: 32px;
  border: 0;
  border-radius: 50%;
  background: #f1f5f9;
  color: #0f172a;
  cursor: pointer;
}

.message-list {
  min-height: 0;
  overflow: auto;
  padding: 0 18px 18px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.message {
  display: flex;
  min-width: 0;
  width: 100%;
}

.message.assistant {
  justify-content: flex-start;
}

.message.user {
  justify-content: flex-end;
}

.bubble {
  min-width: 0;
  max-width: 92%;
  border-radius: 16px;
  padding: 12px 14px;
  background: #f8fafc;
  border: 1px solid rgba(148, 163, 184, 0.16);
  overflow: hidden;
}

.message.assistant .bubble {
  width: 100%;
  max-width: 100%;
}

.message.user .bubble {
  color: #fff;
  background: #1677ff;
  border-color: #1677ff;
}

pre {
  margin: 0;
  white-space: pre-wrap;
  overflow-wrap: anywhere;
  word-break: break-word;
  font-family: inherit;
  line-height: 1.7;
}

.composer {
  position: relative;
  display: grid;
  gap: 14px;
  min-width: 0;
  margin: 0 14px 14px;
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
  min-width: 0;
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
  justify-content: flex-end;
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

.floating-trigger {
  min-width: 132px;
  height: 54px;
  padding: 0 18px;
  display: inline-flex;
  align-items: center;
  gap: 10px;
  border: 0;
  border-radius: 999px;
  background: linear-gradient(135deg, #153a84, #315fcd);
  color: #fff;
  font-weight: 700;
  cursor: pointer;
  box-shadow: 0 18px 34px rgba(49, 95, 205, 0.28);
  user-select: none;
  touch-action: none;
}

.trigger-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: #9cc2ff;
  box-shadow: 0 0 0 6px rgba(156, 194, 255, 0.18);
}

.agent-panel-enter-active,
.agent-panel-leave-active {
  transition: all 0.22s ease;
}

.agent-panel-enter-from,
.agent-panel-leave-to {
  opacity: 0;
  transform: translateY(12px) scale(0.98);
}

@media (max-width: 760px) {
  .floating-agent {
    right: 12px;
    bottom: 12px;
    max-width: calc(100vw - 24px);
  }

  .agent-panel {
    width: min(520px, calc(100vw - 24px));
    height: min(78vh, 720px);
    border-radius: 20px;
  }

  .floating-trigger {
    justify-self: end;
  }
}
</style>
