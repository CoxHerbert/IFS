<template>
  <view class="page">
    <view class="profile-card">
      <view class="avatar">{{ profileName.slice(0, 1) || '轻' }}</view>
      <view class="profile-copy"><text class="profile-title">{{ profileName }}</text><text class="profile-sub" @tap="copyUID">UID：{{ account?.uid || profile.uid }}　复制</text></view>
    </view>

    <view v-if="profileEditorOpen" class="editor-mask" @tap="profileEditorOpen = false">
      <view class="profile-editor" @tap.stop>
        <view class="editor-heading"><text class="editor-title">编辑个人资料</text><text class="editor-close" @tap="profileEditorOpen = false">×</text></view>
        <text class="editor-label">昵称</text><input v-model="profileDraft.nickname" class="editor-input" maxlength="30" placeholder="请输入昵称" />
        <text class="editor-label">手机号</text><input v-model="profileDraft.phone" class="editor-input" type="number" maxlength="20" placeholder="选填" />
        <text class="editor-label">邮箱</text><input v-model="profileDraft.email" class="editor-input" type="text" maxlength="120" placeholder="选填" />
        <text class="editor-note">UID {{ account?.uid || profile.uid }} · 个人资料仅保存在本机</text>
        <wd-button block type="primary" @click="saveProfile">保存资料</wd-button>
      </view>
    </view>

    <view class="section-title">账号与资料</view>
    <view class="data-card">
      <view class="menu-row" @tap="openProfileEditor"><view class="row-icon profile-icon">人</view><view class="row-copy"><text class="row-title">个人信息</text><text class="row-sub">{{ [profile.phone, profile.email].filter(Boolean).join(' · ') || '编辑昵称、手机号和邮箱' }}</text></view><text class="row-action">编辑　›</text></view>
      <view class="row-divider" />
      <view class="menu-row" @tap="account ? logout() : openLogin()"><view class="row-icon account-icon">账</view><view class="row-copy"><text class="row-title">记账账号</text><text class="row-sub">{{ account ? account.loginId : '登录或注册均为自愿' }}</text></view><text class="row-action">{{ account ? '退出登录' : '登录 / 注册' }}　›</text></view>
    </view>

    <view class="section-title">数据保存与同步</view>
    <view class="data-card">
      <view class="menu-row"><view class="row-icon local">本</view><view class="row-copy"><text class="row-title">本机保存</text><text class="row-sub">账单默认保存在当前设备</text></view><text class="row-status">已开启</text></view>
      <view class="row-divider" />
      <view class="menu-row"><view class="row-icon cloud">云</view><view class="row-copy"><text class="row-title">云端同步</text><text class="row-sub">{{ account ? (cloudEnabled ? '账单同步已开启' : '关闭时账单不会上传') : '登录后可选开启，默认关闭' }}</text></view><wd-switch v-if="account" v-model="cloudEnabled" size="20px" @change="onCloudToggle" /><text v-else class="row-status muted">需登录</text></view>
      <view class="info-row">{{ account ? '可随时关闭云同步；关闭后，云端账单会删除，本机数据保留。' : '未登录或未开启同步时，账单数据只保存在本机。' }}</view>
      <view v-if="account && cloudEnabled" class="menu-row action-row" @tap="syncNow"><view class="row-icon sync-icon">↻</view><view class="row-copy"><text class="row-title">立即同步</text><text class="row-sub">同步本机账单并回显云端数据</text></view><text class="row-action">同步　›</text></view>
    </view>

    <view class="section-title">本机数据</view>
    <view class="data-card storage-card"><view class="menu-row"><view class="row-icon records-icon">账</view><view class="row-copy"><text class="row-title">账单记录</text><text class="row-sub">当前设备保存 {{ recordCount }} 笔</text></view><text class="row-status">{{ recordCount }} 笔</text></view><view class="row-divider" /><view class="menu-row" @tap="clearLocalData"><view class="row-icon clear-icon">删</view><view class="row-copy"><text class="row-title">清除本机数据</text><text class="row-sub">删除本设备的全部账单记录</text></view><text class="danger">清除　›</text></view></view>

    <view class="footnote">账单导入文件在本机处理；是否开启云同步由你单独决定。</view>
  <LedgerTabBar current="account" /></view>
</template>


<script setup lang="ts">
import LedgerTabBar from '../../components/LedgerTabBar.vue'
import { computed, reactive, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { appendRecords, fetchCloudRecords, getAccount, getCloudSyncEnabled, getLocalProfile, getRecords, setAccount, setLocalProfile, updateCloudConsent, uploadCloudRecords, type LocalAccount, type LocalProfile } from '../../utils/ledger'

const account = ref<LocalAccount | null>(null)
const profile = ref<LocalProfile>(getLocalProfile())
const profileDraft = reactive({ nickname: '', phone: '', email: '' })
const profileEditorOpen = ref(false)
const cloudEnabled = ref(false)
const recordCount = ref(0)
const profileName = computed(() => profile.value.nickname || account.value?.displayName || '本地使用中')
onShow(() => {
  account.value = getAccount()
  profile.value = getLocalProfile()
  cloudEnabled.value = account.value ? getCloudSyncEnabled() : false
  recordCount.value = getRecords().length
})

function openProfileEditor() {
  Object.assign(profileDraft, { nickname: profile.value.nickname, phone: profile.value.phone, email: profile.value.email })
  profileEditorOpen.value = true
}
function saveProfile() {
  const email = profileDraft.email.trim()
  if (email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
    uni.showToast({ title: '邮箱格式不正确', icon: 'none' })
    return
  }
  profile.value = { ...profile.value, nickname: profileDraft.nickname.trim(), phone: profileDraft.phone.trim(), email }
  setLocalProfile(profile.value)
  profileEditorOpen.value = false
  uni.showToast({ title: '个人资料已保存到本机', icon: 'success' })
}
function copyUID() {
  uni.setClipboardData({ data: account.value?.uid || profile.value.uid, success: () => uni.showToast({ title: 'UID 已复制', icon: 'success' }) })
}
function openLogin() { uni.navigateTo({ url: '/pages/login/index' }) }
function logout() {
  setAccount(null)
  account.value = null
  cloudEnabled.value = false
  uni.showToast({ title: '已退出登录', icon: 'none' })
}
async function syncNow() {
  try {
    await uploadCloudRecords(account.value)
    const cloud = await fetchCloudRecords(account.value)
    appendRecords(cloud.records)
    recordCount.value = getRecords().length
    uni.showToast({ title: '账单已同步', icon: 'success' })
  } catch (error) {
    uni.showToast({ title: error instanceof Error ? error.message : '同步失败', icon: 'none' })
  }
}
function onCloudToggle(event: { value?: boolean }) {
  const enabled = Boolean(event.value ?? cloudEnabled.value)
  if (!account.value) {
    cloudEnabled.value = false
    return
  }
  if (enabled) {
    uni.showModal({ title: '开启云端同步？', content: '开启后，账单数据将上传到云端并在已登录设备间同步。', confirmText: '开启同步', success: result => {
      if (!result.confirm) { cloudEnabled.value = false; return }
      void (async () => {
        try {
          await updateCloudConsent(true, account.value)
          cloudEnabled.value = true
          await uploadCloudRecords(account.value)
          const cloud = await fetchCloudRecords(account.value)
          appendRecords(cloud.records)
          recordCount.value = getRecords().length
          uni.showToast({ title: '云同步已开启', icon: 'success' })
        } catch (error) {
          cloudEnabled.value = getCloudSyncEnabled()
          uni.showToast({ title: error instanceof Error ? error.message : '同步失败，请重试', icon: 'none' })
        }
      })()
    } })
  } else {
    uni.showModal({ title: '关闭云同步？', content: '关闭后，服务器会删除此账号已上传的账单；本机账单会保留。', confirmText: '关闭并删除云端数据', success: result => {
      if (!result.confirm) { cloudEnabled.value = true; return }
      void updateCloudConsent(false, account.value).then(() => {
        cloudEnabled.value = false
        uni.showToast({ title: '已关闭同步并清除云端账单', icon: 'success' })
      }).catch(error => {
        cloudEnabled.value = true
        uni.showToast({ title: error instanceof Error ? error.message : '操作失败', icon: 'none' })
      })
    } })
  }
}
function clearLocalData() {
  uni.showModal({ title: '清除本机账单？', content: '此操作会删除当前设备保存的全部账单，且无法撤销。', confirmColor: '#c84b46', success: result => {
    if (!result.confirm) return
    uni.removeStorageSync('ifs-ledger-records')
    recordCount.value = 0
    uni.showToast({ title: '本机账单已清除', icon: 'success' })
  } })
}
</script>

<style scoped>
.page{padding:30rpx 32rpx 170rpx}.profile-card,.data-card{background:#fff;border-radius:24rpx;padding:25rpx 23rpx}.profile-card{display:flex;align-items:center;gap:16rpx;margin-bottom:32rpx;box-shadow:0 8rpx 24rpx rgba(30,52,43,.04)}.avatar{width:76rpx;height:76rpx;border-radius:50%;background:#e5f1ec;color:#347662;display:flex;align-items:center;justify-content:center;font-size:30rpx;font-weight:700;flex:none}.profile-copy{flex:1;min-width:0;display:flex;flex-direction:column;gap:8rpx}.profile-title{font-size:25rpx;font-weight:700}.profile-sub,.row-sub{font-size:19rpx;color:#9aa1a4;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}.editor-mask{position:fixed;z-index:1000;inset:0;background:rgba(19,32,26,.42);display:flex;align-items:flex-end}.profile-editor{width:100%;padding:34rpx 34rpx calc(34rpx + env(safe-area-inset-bottom));border-radius:30rpx 30rpx 0 0;background:#f8faf9;box-shadow:0 -14rpx 48rpx rgba(18,42,31,.16)}.editor-heading{display:flex;align-items:center;justify-content:space-between;margin-bottom:28rpx}.editor-title{font-size:30rpx;font-weight:700;color:#25372f}.editor-close{font-size:42rpx;color:#9aa39e;padding:0 8rpx}.editor-label{display:block;margin:18rpx 0 10rpx;font-size:20rpx;font-weight:600;color:#66716b}.editor-input{height:88rpx;padding:0 22rpx;border:1rpx solid #e8eeea;border-radius:18rpx;background:#fff;font-size:23rpx;color:#25372f}.editor-note{display:block;margin:20rpx 0 24rpx;color:#919b95;font-size:18rpx}.section-title{font-size:23rpx;font-weight:700;color:#89938d;letter-spacing:1rpx;margin:28rpx 6rpx 13rpx}.data-card{padding:0 22rpx;overflow:hidden;box-shadow:0 8rpx 24rpx rgba(30,52,43,.035)}.menu-row{min-height:104rpx;display:flex;align-items:center;gap:16rpx}.row-icon{width:54rpx;height:54rpx;border-radius:17rpx;display:flex;align-items:center;justify-content:center;font-size:21rpx;font-weight:650;flex:none}.profile-icon{background:#eaf4ef;color:#397a61}.account-icon{background:#eef1ff;color:#6674b8}.local{background:#eaf4ef;color:#397a61}.cloud{background:#edf2ff;color:#5870b9}.sync-icon{background:#eaf5ef;color:#347662;font-size:30rpx}.records-icon{background:#fff4e7;color:#a67841}.clear-icon{background:#fff0ef;color:#c25750}.row-copy{flex:1;min-width:0;display:flex;flex-direction:column;gap:7rpx}.row-title{font-size:22rpx;font-weight:600;color:#303a35}.row-sub{font-size:18rpx}.row-action{font-size:19rpx;color:#347662;flex:none}.row-status{font-size:19rpx;color:#43856c;flex:none}.row-status.muted{color:#a1a8a4}.row-divider{height:1rpx;background:#f0f2f1;margin-left:70rpx}.info-row{padding:18rpx 4rpx 20rpx;color:#8b9590;font-size:18rpx;line-height:1.65}.action-row{border-top:1rpx solid #f0f2f1}.danger{font-size:19rpx;color:#c25750;flex:none}.footnote{font-size:18rpx;color:#a0a7a3;line-height:1.7;text-align:center;padding:28rpx 14rpx}
</style>
