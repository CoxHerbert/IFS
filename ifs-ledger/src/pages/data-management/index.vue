<template>
  <view class="page">
    <view class="navigation-bar"><view class="nav-back" @tap="goBack"><wd-icon name="arrow-left" size="20" /></view><text class="nav-title">数据管理</text><view class="nav-placeholder" /></view>
    <view class="page-heading"><text class="page-subtitle">账单保存与同步</text></view>

    <view class="section-title">保存方式</view>
    <view class="data-card">
      <view class="menu-row"><view class="row-icon local"><wd-icon name="safe" size="20" /></view><view class="row-copy"><text class="row-title">本机保存</text></view><text class="row-status">已开启</text></view>
      <view class="row-divider" />
      <view class="menu-row"><view class="row-icon cloud"><wd-icon name="cloud" size="20" /></view><view class="row-copy"><text class="row-title">云端同步</text></view><wd-switch v-if="account" v-model="cloudEnabled" size="20px" @change="onCloudToggle" /><text v-else class="row-status muted">需登录</text></view>
      <view class="info-row">{{ account ? '可随时关闭云同步；关闭后，云端账单会删除，本机数据保留。' : '未登录或未开启同步时，账单数据只保存在本机。' }}</view>
      <view v-if="account && cloudEnabled" class="menu-row action-row" @tap="syncNow"><view class="row-icon sync-icon"><wd-icon name="sync" size="20" /></view><view class="row-copy"><text class="row-title">立即同步</text></view><view class="row-action">同步 <wd-icon name="arrow-right" size="15" /></view></view>
    </view>

    <view class="section-title">本机数据</view>
    <view class="data-card">
      <view class="menu-row"><view class="row-icon records-icon"><wd-icon name="record" size="20" /></view><view class="row-copy"><text class="row-title">账单记录</text></view><text class="row-status">{{ recordCount }} 笔</text></view>
      <view class="row-divider" />
      <view class="menu-row" @tap="clearLocalData"><view class="row-icon clear-icon"><wd-icon name="delete" size="20" /></view><view class="row-copy"><text class="row-title">清除本机数据</text></view><view class="danger">清除 <wd-icon name="arrow-right" size="15" /></view></view>
    </view>

    <view class="footnote">关闭云同步不会删除本机账单；清除本机数据仅影响当前设备。</view>

    <view v-if="cloudDisableConfirmOpen" class="editor-mask confirm-mask" @tap.self="cancelDisableCloudSync">
      <view class="confirm-dialog" @tap.stop>
        <text class="confirm-title">关闭云同步？</text>
        <text class="confirm-message">关闭后，服务器会删除此账号已上传的账单，本机账单会保留。</text>
        <view class="confirm-actions"><button class="confirm-cancel" @tap="cancelDisableCloudSync">暂不关闭</button><button class="confirm-danger" :disabled="cloudTogglePending" @tap="confirmDisableCloudSync">{{ cloudTogglePending ? '正在关闭…' : '关闭同步并删除云端数据' }}</button></view>
      </view>
    </view>

    <LedgerTabBar current="account" />
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import LedgerTabBar from '../../components/LedgerTabBar.vue'
import { appendRecords, fetchCloudImportBatches, fetchCloudRecords, getAccount, getCloudSyncEnabled, getRecords, mergeImportBatches, updateCloudConsent, uploadCloudImportBatches, uploadCloudRecords, type LocalAccount } from '../../utils/ledger'
import { applyTheme } from '../../utils/theme'

const account = ref<LocalAccount | null>(null)
const cloudEnabled = ref(false)
const recordCount = ref(0)
const cloudDisableConfirmOpen = ref(false)
const cloudTogglePending = ref(false)

onShow(() => {
  applyTheme()
  account.value = getAccount()
  cloudEnabled.value = account.value ? getCloudSyncEnabled() : false
  recordCount.value = getRecords().length
})

function goBack() { uni.navigateBack() }

async function syncNow() {
  try {
    await uploadCloudRecords(account.value)
    const cloud = await fetchCloudRecords(account.value)
    appendRecords(cloud.records)
    const batchResult = await fetchCloudImportBatches(account.value)
    await uploadCloudImportBatches(account.value, mergeImportBatches(batchResult.batches))
    recordCount.value = getRecords().length
    uni.showToast({ title: '账单已同步', icon: 'success' })
  } catch (error) {
    uni.showToast({ title: error instanceof Error ? error.message : '同步失败', icon: 'none' })
  }
}

function onCloudToggle(event: { value?: boolean }) {
  const enabled = Boolean(event.value ?? cloudEnabled.value)
  if (!account.value) { cloudEnabled.value = false; return }
  if (!enabled) { cloudEnabled.value = true; cloudDisableConfirmOpen.value = true; return }
  uni.showModal({ title: '开启云端同步？', content: '开启后，账单数据将上传到云端并在已登录设备间同步。', confirmText: '开启同步', success: result => {
    if (!result.confirm) { cloudEnabled.value = false; return }
    void (async () => {
      try {
        await updateCloudConsent(true, account.value)
        cloudEnabled.value = true
        await uploadCloudRecords(account.value)
        const cloud = await fetchCloudRecords(account.value)
        appendRecords(cloud.records)
        const batchResult = await fetchCloudImportBatches(account.value)
        await uploadCloudImportBatches(account.value, mergeImportBatches(batchResult.batches))
        recordCount.value = getRecords().length
        uni.showToast({ title: '云同步已开启', icon: 'success' })
      } catch (error) {
        cloudEnabled.value = getCloudSyncEnabled()
        uni.showToast({ title: error instanceof Error ? error.message : '同步失败，请重试', icon: 'none' })
      }
    })()
  } })
}

function cancelDisableCloudSync() { cloudDisableConfirmOpen.value = false; cloudEnabled.value = true }
async function confirmDisableCloudSync() {
  if (cloudTogglePending.value || !account.value) return
  cloudTogglePending.value = true
  try {
    await updateCloudConsent(false, account.value)
    cloudEnabled.value = false
    cloudDisableConfirmOpen.value = false
    uni.showToast({ title: '已关闭同步并清除云端账单', icon: 'success' })
  } catch (error) {
    cloudEnabled.value = true
    uni.showToast({ title: error instanceof Error ? error.message : '操作失败', icon: 'none' })
  } finally { cloudTogglePending.value = false }
}

function clearLocalData() {
  uni.showModal({ title: '清除本机账单？', content: '此操作会删除当前设备保存的全部账单，且无法撤销。', confirmColor: '#c84b46', success: result => {
    if (!result.confirm) return
    uni.removeStorageSync('ifs-ledger-records')
    uni.removeStorageSync('ifs-ledger-import-batches')
    recordCount.value = 0
    uni.showToast({ title: '本机账单已清除', icon: 'success' })
  } })
}
</script>

<style scoped>
.page{min-height:100vh;box-sizing:border-box;padding:0 32rpx 170rpx;background:var(--theme-page,#f6f5fa)}.navigation-bar{position:relative;box-sizing:border-box;height:calc(88rpx + env(safe-area-inset-top));margin:0 -32rpx;padding:env(safe-area-inset-top) 28rpx 0;display:flex;align-items:center;justify-content:space-between;background:var(--theme-page,#f6f5fa);border-bottom:1rpx solid rgba(68,60,105,.06)}.nav-back,.nav-placeholder{width:64rpx;height:64rpx;display:flex;align-items:center;justify-content:center;flex:none}.nav-back{border-radius:20rpx;color:var(--theme-primary,#6556b2);background:#fff;box-shadow:0 4rpx 14rpx rgba(54,45,92,.06)}.nav-title{position:absolute;left:50%;transform:translateX(-50%);font-size:32.5rpx;font-weight:700;color:#302c45;white-space:nowrap}.page-heading{padding:22rpx 4rpx 16rpx}.page-subtitle{display:block;font-size:21.3rpx;color:#9290a1}.section-title{margin:24rpx 6rpx 13rpx;font-size:25.8rpx;font-weight:700;color:#5b566e}.data-card{padding:0 22rpx;overflow:hidden;border-radius:24rpx;background:#fff;box-shadow:0 8rpx 24rpx rgba(54,45,92,.035)}.menu-row{min-height:104rpx;display:flex;align-items:center;gap:16rpx}.row-icon{width:54rpx;height:54rpx;border-radius:17rpx;display:flex;align-items:center;justify-content:center;font-size:23.5rpx;font-weight:650;flex:none}.local,.sync-icon{background:var(--theme-soft,#efedf8);color:var(--theme-primary,#6556b2)}.cloud{background:#edf2ff;color:#5870b9}.records-icon{background:#fff4e7;color:#a67841}.clear-icon{background:#fff0ef;color:#c25750}.row-copy{flex:1;min-width:0;display:flex;flex-direction:column;gap:7rpx}.row-title{font-size:24.6rpx;font-weight:600;color:#303a35}.row-sub{font-size:20.2rpx;color:#9894a5}.row-status{font-size:21.3rpx;color:#418b70;flex:none}.row-status.muted{color:#a1a8a4}.row-divider{height:1rpx;background:#efedf3;margin-left:70rpx}.info-row{padding:18rpx 4rpx 20rpx;color:#8b879d;font-size:20.2rpx;line-height:1.65}.action-row{border-top:1rpx solid #efedf3}.row-action{font-size:21.3rpx;color:var(--theme-primary,#6556b2);flex:none}.danger{font-size:21.3rpx;color:#c25750;flex:none}.footnote{padding:28rpx 14rpx;color:#a0a7a3;font-size:20.2rpx;line-height:1.7;text-align:center}.editor-mask{position:fixed;z-index:1200;inset:0;display:flex;align-items:center;justify-content:center;padding:32rpx;background:rgba(36,30,59,.42)}.confirm-dialog{box-sizing:border-box;width:100%;max-width:640rpx;padding:38rpx 34rpx 30rpx;border-radius:26rpx;background:#fff}.confirm-title{display:block;text-align:center;font-size:33.6rpx;font-weight:700;color:#302c45}.confirm-message{display:block;margin:22rpx 0 30rpx;text-align:center;font-size:24.6rpx;line-height:1.65;color:#858194}.confirm-actions{display:flex;flex-direction:column;gap:12rpx}.confirm-cancel,.confirm-danger{box-sizing:border-box;width:100%;min-height:82rpx;margin:0;padding:18rpx 24rpx;border:0;border-radius:16rpx;font-size:24.6rpx;white-space:normal}.confirm-cancel{background:#f3f2f7;color:#5f5b70}.confirm-danger{background:#c25750;color:#fff}.confirm-cancel::after,.confirm-danger::after{border:0}

/* 提升页面文字可读性 */
.nav-title{font-size:37rpx}.page-subtitle{font-size:23.5rpx}.section-title{font-size:30.2rpx}.row-title{font-size:28rpx}.row-sub{font-size:22.4rpx}.row-status,.row-action,.danger{font-size:23.5rpx}</style>


