<template>
  <view class="page">
    <view class="intro"><text class="title">导入账单</text><text class="sub">解析支付平台原始账单，并为每次导入生成独立批次号</text></view>
    <view class="source-card"><view class="source-icon wechat">微</view><view class="source-main"><text class="source-title">微信支付账单</text><text class="source-desc">支持微信导出的 XLSX / CSV 文件</text></view><wd-button type="primary" size="small" @click="chooseFile('wechat')">选择文件</wd-button></view>
    <view class="source-card"><view class="source-icon ali">支</view><view class="source-main"><text class="source-title">支付宝账单</text><text class="source-desc">支持支付宝导出的 XLSX / CSV 文件</text></view><wd-button type="primary" size="small" @click="chooseFile('alipay')">选择文件</wd-button></view>
    <view v-if="preview" class="preview"><text class="preview-title">导入预览</text><text class="file-name">{{ preview.fileName }}</text><text class="batch-id">批次号：{{ preview.batchId }}</text><view class="stats"><text>可导入 {{ preview.records.length }} 笔</text><text>收入 ¥{{ total('income').toFixed(2) }}</text><text>支出 ¥{{ total('expense').toFixed(2) }}</text></view><text v-if="preview.skippedRows" class="skipped">已忽略 {{ preview.skippedRows }} 行无效或未完成记录</text><wd-button type="primary" block :disabled="importing" @click="confirmImport">{{ importing ? '导入中…' : '确认导入到本地' }}</wd-button></view>
    <view class="guide"><text class="guide-title">导入说明</text><view class="step"><text class="number">1</text><text>先从微信或支付宝导出账单文件</text></view><view class="step"><text class="number">2</text><text>确认预览中的流水数量和金额</text></view><view class="step"><text class="number">3</text><text>确认后每条记录会标记批次号，可用于追溯来源</text></view></view>
    <view class="privacy">🔒 文件在本机解析；仅当你开启云同步时，记录及批次标识才会上传</view>
    <view v-if="batches.length" class="history"><text class="history-title">最近导入批次</text><view v-for="batch in batches" :key="batch.batchId" class="batch-row"><view class="batch-heading"><text class="batch-source">{{ batch.source === 'wechat' ? '微信' : '支付宝' }} · {{ batch.fileName }}</text><text class="delete-batch" @tap="removeBatch(batch)">删除</text></view><text class="batch-code">{{ batch.batchId }}</text><text class="batch-meta">{{ formatTime(batch.importedAt) }} · 导入 {{ batch.importedCount }} 笔 · 重复 {{ batch.duplicateCount }} 笔</text></view></view>
    <LedgerTabBar current="import" />
  </view>
</template>

<script setup lang="ts">
import LedgerTabBar from '../../components/LedgerTabBar.vue'
import { computed, onMounted, ref } from 'vue'
import { createImportBatchId, parseLedgerStatement } from '../../utils/statementImport'
import { getCloudSyncEnabled, getAccount, getImportBatches, removeImportBatch, saveImportBatch, type LedgerImportBatch, type LedgerRecord } from '../../utils/ledger'

type ImportSource = 'wechat' | 'alipay'
const batches = ref<LedgerImportBatch[]>([])
const preview = ref<{ source: ImportSource; fileName: string; batchId: string; records: LedgerRecord[]; skippedRows: number } | null>(null)
const importing = ref(false)
const deletingBatch = ref('')
onMounted(() => { batches.value = getImportBatches() })
const total = (type: LedgerRecord['type']) => computed(() => preview.value?.records.filter(item => item.type === type).reduce((sum, item) => sum + item.amount, 0) || 0).value

function readH5File(file: File, source: ImportSource) {
  file.arrayBuffer().then(buffer => setPreview(buffer, file.name, source)).catch(() => showError('读取文件失败'))
}
async function setPreview(data: ArrayBuffer | Uint8Array, fileName: string, source: ImportSource) {
  try {
    const parsed = await parseLedgerStatement(data, source)
    preview.value = { source, fileName, batchId: createImportBatchId(source), ...parsed }
    uni.showToast({ title: '账单解析完成', icon: 'success' })
  } catch (error) { showError(error instanceof Error ? error.message : '账单解析失败') }
}
function showError(title: string) { uni.showToast({ title, icon: 'none', duration: 3000 }) }

function chooseFile(source: ImportSource) {
  // #ifdef MP-WEIXIN
  wx.chooseMessageFile({ count: 1, type: 'file', extension: ['xlsx', 'csv'], success: result => { const file = result.tempFiles[0]; if (!file) return; wx.getFileSystemManager().readFile({ filePath: file.path, success: read => setPreview(read.data as Uint8Array, file.name, source), fail: () => showError('读取文件失败') }) }, fail: () => {} })
  // #endif
  // #ifdef MP-ALIPAY
  my.chooseFile({ count: 1, success: result => { const path = result.filePaths?.[0]; if (!path) return; my.getFileSystemManager().readFile({ filePath: path, success: read => setPreview(read.data as Uint8Array, path.split('/').pop() || '支付宝账单', source), fail: () => showError('读取文件失败') }) }, fail: () => {} })
  // #endif
  // #ifdef H5
  uni.chooseFile({ count: 1, extension: ['.xlsx', '.csv'], success: result => { const file = result.tempFiles[0] as File | undefined; if (file) readH5File(file, source) }, fail: () => {} })
  // #endif
}

function confirmImport() {
  if (!preview.value || importing.value) return
  importing.value = true
  try {
    const batch = saveImportBatch(preview.value.records, preview.value.source, preview.value.fileName, preview.value.batchId)
    batches.value = getImportBatches()
    preview.value = null
    uni.showToast({ title: `已导入 ${batch.importedCount} 笔`, icon: 'success' })
  } catch { showError('保存账单失败') } finally { importing.value = false }
}
function formatTime(value: string) { return new Date(value).toLocaleString('zh-CN', { hour12: false }) }
async function removeBatch(batch: LedgerImportBatch) {
  if (deletingBatch.value) return
  uni.showModal({ title: '删除导入批次', content: `删除批次 ${batch.batchId} 后，该批次导入的 ${batch.importedCount} 笔账单也会一并删除，且无法恢复。`, confirmText: '删除批次', confirmColor: '#c25750', success: result => {
    if (!result.confirm) return
    deletingBatch.value = batch.batchId
    removeImportBatch(batch.batchId).then(() => {
      batches.value = getImportBatches()
      uni.showToast({ title: getCloudSyncEnabled() && getAccount()?.token ? '批次与云端账单已删除' : '批次及账单已删除', icon: 'success' })
    }).catch(error => showError(error instanceof Error ? `删除失败：${error.message}` : '删除失败，请检查网络后重试')).finally(() => { deletingBatch.value = '' })
  } })
}
</script>

<style scoped>
.page{padding:34rpx 32rpx 170rpx}.intro{display:flex;flex-direction:column;gap:10rpx;margin:4rpx 0 30rpx}.title{font-size:40rpx;font-weight:750}.sub{font-size:22rpx;color:#969da0}.source-card,.preview,.guide,.history{background:#fff;border-radius:22rpx;padding:24rpx;margin-bottom:16rpx}.source-card{display:flex;align-items:center;gap:16rpx}.source-icon{width:68rpx;height:68rpx;display:flex;align-items:center;justify-content:center;border-radius:20rpx;font-size:30rpx;font-weight:700}.wechat{background:#e7f6ef;color:#22a06b}.ali{background:#e8f3ff;color:#1677d2}.source-main{flex:1;display:flex;flex-direction:column;gap:8rpx}.source-title{font-size:24rpx;font-weight:650}.source-desc{font-size:18rpx;color:#a0a6a9}.preview-title,.guide-title,.history-title{display:block;font-size:24rpx;font-weight:700;margin-bottom:16rpx}.file-name,.batch-id,.batch-code{display:block;overflow-wrap:anywhere;color:#69736f;font-size:19rpx;margin:8rpx 0}.batch-id,.batch-code{font-family:monospace;color:#36745d}.stats{display:flex;justify-content:space-between;flex-wrap:wrap;gap:10rpx;padding:14rpx 0;color:#37413c;font-size:19rpx}.skipped{display:block;color:#9a7b3e;font-size:18rpx;margin-bottom:14rpx}.step{display:flex;align-items:flex-start;gap:14rpx;margin-top:16rpx;color:#737c7d;font-size:20rpx;line-height:1.6}.number{width:32rpx;height:32rpx;border-radius:50%;background:#eaf3ef;color:#347662;display:flex;align-items:center;justify-content:center;font-size:17rpx;flex:none}.privacy{text-align:center;color:#9ba3a4;font-size:18rpx;margin:25rpx 0}.batch-row{padding:16rpx 0;border-bottom:1rpx solid #f0f1f1}.batch-row:last-child{border-bottom:0}.batch-heading{display:flex;align-items:center;justify-content:space-between;gap:14rpx}.batch-source{display:block;font-size:20rpx;color:#343b38;flex:1;min-width:0;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}.delete-batch{padding:9rpx 16rpx;border-radius:12rpx;background:#fff0ef;color:#bd5b54;font-size:18rpx;flex:none}.batch-meta{display:block;font-size:17rpx;color:#9ba1a4;margin-top:7rpx}
</style>
