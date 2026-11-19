<template>
  <view class="page">
    <view class="profile-card">
      <view class="profile-cover"><view class="cover-orb orb-one" /><view class="cover-orb orb-two" /></view>
      <view class="profile-main">
        <view class="avatar-wrap" @tap="openAvatarPicker"><view :class="['avatar', { 'frame-sparkle': profile.avatarFrameKey === 'sparkle' }]"><image class="avatar-image" :src="profileAvatarSource" mode="aspectFit" /></view><view class="avatar-edit-badge"><wd-icon name="edit" size="12" color="#fff" /></view></view>
        <view class="profile-copy">
          <view class="nickname-line" @tap="openNicknameEditor"><text class="profile-title">{{ profileName }}</text><wd-icon name="edit" size="14" color="var(--theme-primary, #6556b2)" /></view>
          <view class="profile-sub uid-copy" @tap="copyUID"><text>UID：{{ account?.uid || profile.uid }}</text><wd-icon name="copy" size="13" /></view>
          <text v-if="account?.inviterUid" class="profile-sub inviter-line">邀请人 UID：{{ account.inviterUid }}</text>
        </view>
      </view>
      <view class="profile-quick-info"><view class="quick-item"><text class="quick-label">账号状态</text><text class="quick-value">{{ account ? '已登录' : '本机使用' }}</text></view><view class="quick-divider" /><view class="quick-item"><text class="quick-label">账单记录</text><text class="quick-value">{{ recordCount }} 笔</text></view></view>
    </view>

    <view v-if="avatarPickerOpen" class="editor-mask" @tap="avatarPickerOpen = false">
      <view class="avatar-picker-sheet" @tap.stop>
        <view class="editor-heading"><view><text class="editor-title">选择头像</text><text class="sheet-subtitle">选一个喜欢的形象</text></view><wd-icon class="editor-close" name="close" size="20" @tap="avatarPickerOpen = false" /></view>
        <view class="avatar-grid">
          <view :class="['avatar-choice', { selected: !profileDraft.avatarKey }]" @tap="profileDraft.avatarKey = ''"><image class="avatar-sample avatar-default-logo" :src="defaultProfileAvatar" mode="aspectFit" /><text>轻记</text><wd-icon v-if="!profileDraft.avatarKey" name="check" size="14" /></view>
          <view v-for="item in profileAvatars" :key="item.key" :class="['avatar-choice', { selected: profileDraft.avatarKey === item.key }]" @tap="profileDraft.avatarKey = item.key"><image class="avatar-sample" :src="item.src" mode="aspectFit" /><text>{{ item.label }}</text><wd-icon v-if="profileDraft.avatarKey === item.key" name="check" size="14" /></view>
        </view>
        <text class="editor-label frame-label">头像框</text>
        <view class="frame-options">
          <view :class="['frame-option', { selected: !profileDraft.avatarFrameKey }]" @tap="profileDraft.avatarFrameKey = ''"><view class="frame-preview"><image :src="profileAvatarSource" mode="aspectFit" /></view><text>默认</text><wd-icon v-if="!profileDraft.avatarFrameKey" name="check" size="13" /></view>
          <view v-if="referralSummary.avatarFrameUnlocked" :class="['frame-option', { selected: profileDraft.avatarFrameKey === 'sparkle' }]" @tap="profileDraft.avatarFrameKey = 'sparkle'"><view class="frame-preview frame-sparkle"><image :src="profileAvatarSource" mode="aspectFit" /></view><text>星光头像框</text><wd-icon v-if="profileDraft.avatarFrameKey === 'sparkle'" name="check" size="13" /></view>
          <view v-else class="frame-locked"><wd-icon name="lock" size="15" /><text>邀请 3 位好友注册，解锁星光头像框（{{ referralSummary.invitedCount }}/3）</text></view>
        </view>
        <wd-button block type="primary" :loading="profileSaving" @click="saveProfile">使用这个头像</wd-button>
      </view>
    </view>

    <view v-if="nicknameEditorOpen" class="editor-mask" @tap="nicknameEditorOpen = false">
      <view class="nickname-sheet" @tap.stop>
        <view class="editor-heading"><view><text class="editor-title">修改昵称</text><text class="sheet-subtitle">最多 30 个字符</text></view><wd-icon class="editor-close" name="close" size="20" @tap="nicknameEditorOpen = false" /></view>
        <input v-model="profileDraft.nickname" class="editor-input" maxlength="30" placeholder="请输入昵称" focus />
        <text class="editor-note">{{ account ? '昵称会保存到账号并在其他设备显示' : '昵称保存在当前设备，登录后可同步到账号' }}</text>
        <wd-button block type="primary" :loading="profileSaving" @click="saveProfile">保存昵称</wd-button>
      </view>
    </view>

    <view v-if="profileEditorOpen" class="editor-mask" @tap="profileEditorOpen = false">
      <view class="profile-editor" @tap.stop>
        <view class="editor-heading"><text class="editor-title">编辑个人资料</text><wd-icon class="editor-close" name="close" size="20" @tap="profileEditorOpen = false" /></view>
        <text class="editor-label">昵称</text><input v-model="profileDraft.nickname" class="editor-input" maxlength="30" placeholder="请输入昵称" />
        <text class="editor-label">手机号</text><input v-model="profileDraft.phone" class="editor-input" type="number" maxlength="20" placeholder="选填" />
        <text class="editor-label">邮箱</text><input v-model="profileDraft.email" class="editor-input" type="text" maxlength="120" placeholder="选填" />
        <text class="editor-note">UID {{ account?.uid || profile.uid }} · {{ account ? '昵称和头像保存到账号，联系方式仅保存在本机' : '登录后可将昵称和头像保存到账号' }}</text>
        <wd-button block type="primary" :loading="profileSaving" @click="saveProfile">{{ profileSaving ? '保存中…' : '保存资料' }}</wd-button>
      </view>
    </view>

    <view v-if="cloudDisableConfirmOpen" class="editor-mask confirm-mask" @tap.self="cancelDisableCloudSync">
      <view class="confirm-dialog" @tap.stop>
        <text class="confirm-title">关闭云同步？</text>
        <text class="confirm-message">关闭后，服务器会删除此账号已上传的账单，本机账单会保留。</text>
        <view class="confirm-actions">
          <button class="confirm-cancel" @tap="cancelDisableCloudSync">暂不关闭</button>
          <button class="confirm-danger" :disabled="cloudTogglePending" @tap="confirmDisableCloudSync">{{ cloudTogglePending ? '正在关闭…' : '关闭同步并删除云端数据' }}</button>
        </view>
      </view>
    </view>

    <view v-if="accountDeleteOpen" class="editor-mask confirm-mask" @tap.self="closeDeleteAccount">
      <view class="account-delete-dialog" @tap.stop>
        <view class="editor-heading"><view><text class="editor-title">注销记账账号</text><text class="sheet-subtitle">此操作无法撤销</text></view><wd-icon class="editor-close" name="close" size="20" @tap="closeDeleteAccount" /></view>
        <text class="delete-warning">确认后将永久删除账号、个人资料和云端账单。保存在本机的账单会保留。</text>
        <text class="editor-label">输入当前账号密码</text>
        <input v-model="deletePassword" class="editor-input" password maxlength="72" placeholder="请输入密码以确认" />
        <view class="delete-actions"><wd-button plain @click="closeDeleteAccount">取消</wd-button><wd-button type="error" :loading="accountDeleting" @click="confirmDeleteAccount">永久注销</wd-button></view>
      </view>
    </view>

    <view v-if="passwordEditorOpen" class="editor-mask confirm-mask" @tap.self="closePasswordEditor">
      <view class="password-dialog" @tap.stop>
        <view class="editor-heading"><text class="editor-title">修改密码</text><wd-icon class="editor-close" name="close" size="20" @tap="closePasswordEditor" /></view>
        <text class="editor-label">当前密码</text><input v-model="passwordDraft.oldPassword" class="editor-input" password maxlength="72" placeholder="请输入当前密码" />
        <text class="editor-label">新密码</text><input v-model="passwordDraft.newPassword" class="editor-input" password maxlength="72" placeholder="至少 8 个字符" />
        <view class="delete-actions"><wd-button plain @click="closePasswordEditor">取消</wd-button><wd-button type="primary" :loading="passwordSaving" @click="submitPasswordChange">保存</wd-button></view>
      </view>
    </view>

    <view class="section-title"><wd-icon name="share" size="19" /><text>邀请好友</text></view>
    <view class="data-card invite-card">
      <view class="invite-summary"><view class="invite-count"><text class="invite-number">{{ referralSummary.invitedCount }}</text><text class="invite-caption">位好友已注册</text></view><view class="invite-reward"><wd-icon name="star-fill" size="18" color="#d7a63c" /><text>邀请星 {{ referralSummary.inviteStars }}</text></view></view>
      <view class="invite-progress"><view class="progress-track"><view class="progress-fill" :style="{ width: `${Math.min(referralSummary.invitedCount / 3 * 100, 100)}%` }" /></view><text>{{ referralSummary.avatarFrameUnlocked ? '已解锁星光头像框' : `再邀请 ${Math.max(3 - referralSummary.invitedCount, 0)} 位好友解锁星光头像框` }}</text></view>
      <button v-if="account" class="share-invite-button" open-type="share" @tap="shareInvite"><wd-icon name="share" size="17" color="#fff" />邀请好友注册</button>
      <view v-else class="share-invite-button login-invite" @tap="openLogin">登录 <wd-icon name="arrow-right" size="15" color="#fff" /></view>
    </view>

    <view class="section-title"><wd-icon name="user" size="19" /><text>账号与资料</text></view>
    <view class="data-card">
      <view class="menu-row" @tap="openProfileEditor"><view class="row-icon profile-icon"><wd-icon name="user" size="20" /></view><view class="row-copy"><text class="row-title">个人信息</text><text v-if="profile.phone || profile.email" class="row-sub">{{ [profile.phone, profile.email].filter(Boolean).join(' · ') }}</text></view><view class="row-action">编辑 <wd-icon name="arrow-right" size="15" /></view></view>
      <view class="row-divider" />
      <view class="menu-row" @tap="account ? logout() : openLogin()"><view class="row-icon account-icon"><wd-icon :name="account ? 'poweroff' : 'idcard'" size="20" /></view><view class="row-copy"><text class="row-title">{{ account ? '退出登录' : '记账账号' }}</text><text class="row-sub">{{ account ? account.loginId : '登录或注册' }}</text></view><view class="row-action">{{ account ? '退出' : '登录 / 注册' }} <wd-icon name="arrow-right" size="15" /></view></view>
      <view v-if="account" class="row-divider" />
      <view v-if="account" class="menu-row" @tap="openPasswordEditor"><view class="row-icon account-icon"><wd-icon name="lock" size="20" /></view><view class="row-copy"><text class="row-title">修改密码</text></view><view class="row-action">修改 <wd-icon name="arrow-right" size="15" /></view></view>
      <view v-if="account" class="row-divider" />
      <view v-if="account" class="menu-row" @tap="openDeleteAccount"><view class="row-icon delete-account-icon"><wd-icon name="delete" size="20" /></view><view class="row-copy"><text class="row-title danger-title">注销账号</text></view><view class="danger">注销 <wd-icon name="arrow-right" size="15" /></view></view>
    </view>

    <view class="section-title"><wd-icon name="safe" size="19" /><text>数据管理</text></view>
    <view class="data-card">
      <view class="menu-row" @tap="openDataManagement"><view class="row-icon local"><wd-icon name="list" size="20" /></view><view class="row-copy"><text class="row-title">数据管理</text></view><view class="row-action">管理 <wd-icon name="arrow-right" size="15" /></view></view>
    </view>

    <view class="section-title"><wd-icon name="settings" size="19" /><text>偏好设置</text></view>
    <view class="data-card">
      <view class="menu-row" @tap="themePickerOpen = true"><view class="row-icon theme-icon"><wd-icon name="font-colors" size="20" /></view><view class="row-copy"><text class="row-title">主题设置</text></view><view class="theme-current"><view class="theme-dot" :style="{ background: currentTheme.primary }" />{{ currentTheme.label }} <wd-icon name="arrow-right" size="15" /></view></view>
    </view>

    <view v-if="themePickerOpen" class="editor-mask theme-mask" @tap.self="themePickerOpen = false">
      <view class="theme-picker">
        <view class="editor-heading"><text class="editor-title">选择主题</text><wd-icon class="editor-close" name="close" size="20" @tap="themePickerOpen = false" /></view>
        <view v-for="item in themes" :key="item.name" :class="['theme-option', { active: themeName === item.name }]" @tap="chooseTheme(item.name)">
          <view class="theme-swatches"><view :style="{ background: item.primary }" /><view :style="{ background: item.soft }" /><view :style="{ background: item.page }" /></view>
          <text class="theme-label">{{ item.label }}</text>
          <wd-icon v-if="themeName === item.name" name="check" size="18" color="var(--theme-primary, #6556b2)" />
        </view>
      </view>
    </view>

    <view class="footnote">账号、数据与主题设置集中管理。</view>
  <LedgerTabBar current="account" /></view>
</template>


<script setup lang="ts">
import LedgerTabBar from '../../components/LedgerTabBar.vue'
import { computed, reactive, ref } from 'vue'
import { onShareAppMessage, onShow } from '@dcloudio/uni-app'
import { appendRecords, clearLocalProfile, deleteLedgerAccount, fetchAccountProfile, fetchCloudImportBatches, fetchCloudRecords, fetchReferralSummary, getAccount, getCloudSyncEnabled, getLocalProfile, getRecords, mergeImportBatches, saveAccountProfile, setAccount, setLocalProfile, updateCloudConsent, updateLedgerPassword, uploadCloudImportBatches, uploadCloudRecords, type LocalAccount, type LocalProfile, type ReferralSummary } from '../../utils/ledger'
import { applyTheme, getTheme, setTheme, themes, type ThemeName } from '../../utils/theme'
import { defaultProfileAvatar, getProfileAvatarSource, profileAvatars } from '../../utils/profileAvatars'

const account = ref<LocalAccount | null>(null)
const profile = ref<LocalProfile>(getLocalProfile())
const profileDraft = reactive({ nickname: '', phone: '', email: '', avatarKey: '', avatarFrameKey: '' })
const profileEditorOpen = ref(false)
const avatarPickerOpen = ref(false)
const nicknameEditorOpen = ref(false)
const profileSaving = ref(false)
const cloudEnabled = ref(false)
const cloudDisableConfirmOpen = ref(false)
const cloudTogglePending = ref(false)
const accountDeleteOpen = ref(false)
const accountDeleting = ref(false)
const deletePassword = ref('')
const passwordEditorOpen = ref(false)
const passwordSaving = ref(false)
const passwordDraft = reactive({ oldPassword: '', newPassword: '' })
const recordCount = ref(0)
const referralSummary = ref<ReferralSummary>({ invitedCount: 0, inviteStars: 0, avatarFrameUnlocked: false, nextFrameAt: 3 })
const themeName = ref<ThemeName>(getTheme())
const themePickerOpen = ref(false)
const currentTheme = computed(() => themes.find(theme => theme.name === themeName.value) || themes[0])
const profileName = computed(() => profile.value.nickname || account.value?.displayName || '本地使用中')
const profileAvatarSource = computed(() => getProfileAvatarSource(profile.value.avatarKey) || profile.value.avatar || '')
onShow(async () => {
  applyTheme()
  account.value = getAccount()
  profile.value = getLocalProfile()
  cloudEnabled.value = account.value ? getCloudSyncEnabled() : false
  recordCount.value = getRecords().length
  themeName.value = getTheme()
  if (account.value?.token) {
    try {
      const result = await fetchAccountProfile(account.value)
      profile.value = { ...profile.value, nickname: result.profile.nickname, avatarKey: result.profile.avatarKey, avatarFrameKey: result.profile.avatarFrameKey }
      setLocalProfile(profile.value)
    } catch (error) { console.warn('Could not load account profile:', error) }
    try { referralSummary.value = await fetchReferralSummary(account.value) } catch (error) { console.warn('Could not load referral summary:', error) }
  } else {
    referralSummary.value = { invitedCount: 0, inviteStars: 0, avatarFrameUnlocked: false, nextFrameAt: 3 }
  }
})

onShareAppMessage(() => ({
  title: '和我一起用轻记，轻松管理每一笔',
  path: account.value ? `/pages/register/index?inviter=${encodeURIComponent(account.value.uid)}` : '/pages/register/index',
}))

function openProfileEditor() {
  Object.assign(profileDraft, { nickname: profile.value.nickname, phone: profile.value.phone, email: profile.value.email, avatarKey: profile.value.avatarKey || '', avatarFrameKey: profile.value.avatarFrameKey || '' })
  profileEditorOpen.value = true
}
function openAvatarPicker() {
  profileDraft.avatarKey = profile.value.avatarKey || ''
  profileDraft.avatarFrameKey = profile.value.avatarFrameKey || ''
  avatarPickerOpen.value = true
}
function openNicknameEditor() {
  profileDraft.nickname = profile.value.nickname || ''
  nicknameEditorOpen.value = true
}
async function saveProfile() {
  if (profileSaving.value) return
  const email = profileDraft.email.trim()
  if (email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
    uni.showToast({ title: '邮箱格式不正确', icon: 'none' })
    return
  }
  profileSaving.value = true
  try {
    const nickname = profileDraft.nickname.trim()
    if (account.value?.token) {
      const result = await saveAccountProfile({ nickname, avatarKey: profileDraft.avatarKey, avatarFrameKey: profileDraft.avatarFrameKey }, account.value)
      profileDraft.nickname = result.profile.nickname
      profileDraft.avatarKey = result.profile.avatarKey
      profileDraft.avatarFrameKey = result.profile.avatarFrameKey
    }
    profile.value = { ...profile.value, nickname, phone: profileDraft.phone.trim(), email, avatarKey: profileDraft.avatarKey, avatarFrameKey: profileDraft.avatarFrameKey, avatar: '' }
    setLocalProfile(profile.value)
    profileEditorOpen.value = false
    avatarPickerOpen.value = false
    nicknameEditorOpen.value = false
    uni.showToast({ title: account.value ? '个人资料已保存到账号' : '个人资料已保存到本机', icon: 'success' })
  } catch (error) {
    uni.showToast({ title: error instanceof Error ? error.message : '个人资料保存失败', icon: 'none' })
  } finally { profileSaving.value = false }
}
function copyUID() {
  uni.setClipboardData({ data: account.value?.uid || profile.value.uid, success: () => uni.showToast({ title: 'UID 已复制', icon: 'success' }) })
}
function chooseTheme(name: ThemeName) {
  themeName.value = name
  setTheme(name)
  themePickerOpen.value = false
  uni.showToast({ title: '主题已应用', icon: 'success' })
}
function openLogin() { uni.navigateTo({ url: '/pages/login/index' }) }
function openDataManagement() { uni.navigateTo({ url: '/pages/data-management/index' }) }
function openImport() { uni.navigateTo({ url: '/pages/import/index' }) }
function shareInvite() {
  // #ifdef H5
  if (!account.value) return
  const inviteUrl = new URL('/pages/register/index', window.location.origin)
  inviteUrl.searchParams.set('inviter', account.value.uid)
  const shareData = { title: '和我一起用轻记，轻松管理每一笔', text: '注册轻记，记账更轻松', url: inviteUrl.toString() }
  if (navigator.share) {
    void navigator.share(shareData).catch(() => {})
  } else {
    uni.setClipboardData({ data: inviteUrl.toString(), success: () => uni.showToast({ title: '邀请链接已复制', icon: 'success' }) })
  }
  // #endif
}
function logout() {
  setAccount(null)
  account.value = null
  clearLocalProfile()
  profile.value = getLocalProfile()
  cloudEnabled.value = false
  uni.showToast({ title: '已退出登录', icon: 'none' })
}
function openDeleteAccount() {
  deletePassword.value = ''
  accountDeleteOpen.value = true
}
function openPasswordEditor() {
  Object.assign(passwordDraft, { oldPassword: '', newPassword: '' })
  passwordEditorOpen.value = true
}
function closePasswordEditor() {
  if (passwordSaving.value) return
  passwordEditorOpen.value = false
  Object.assign(passwordDraft, { oldPassword: '', newPassword: '' })
}
async function submitPasswordChange() {
  if (!account.value || passwordSaving.value) return
  if (!passwordDraft.oldPassword || !passwordDraft.newPassword) {
    uni.showToast({ title: '请填写完整', icon: 'none' })
    return
  }
  if (passwordDraft.newPassword.length < 8 || passwordDraft.newPassword.length > 72) {
    uni.showToast({ title: '新密码需为 8 到 72 个字符', icon: 'none' })
    return
  }
  passwordSaving.value = true
  try {
    await updateLedgerPassword(passwordDraft.oldPassword, passwordDraft.newPassword, account.value)
    passwordEditorOpen.value = false
    Object.assign(passwordDraft, { oldPassword: '', newPassword: '' })
    uni.showToast({ title: '密码已修改', icon: 'success' })
  } catch (error) {
    uni.showToast({ title: error instanceof Error ? error.message : '修改失败，请重试', icon: 'none' })
  } finally { passwordSaving.value = false }
}
function closeDeleteAccount() {
  if (accountDeleting.value) return
  accountDeleteOpen.value = false
  deletePassword.value = ''
}
async function confirmDeleteAccount() {
  if (!account.value || accountDeleting.value) return
  if (!deletePassword.value) {
    uni.showToast({ title: '请输入账号密码', icon: 'none' })
    return
  }
  accountDeleting.value = true
  try {
    await deleteLedgerAccount(deletePassword.value, account.value)
    setAccount(null)
    account.value = null
    clearLocalProfile()
    profile.value = getLocalProfile()
    cloudEnabled.value = false
    accountDeleteOpen.value = false
    deletePassword.value = ''
    uni.showToast({ title: '账号已注销，本机账单已保留', icon: 'success' })
  } catch (error) {
    uni.showToast({ title: error instanceof Error ? error.message : '注销失败，请重试', icon: 'none' })
  } finally { accountDeleting.value = false }
}
async function syncNow() {
  try {
    await uploadCloudRecords(account.value)
    const cloud = await fetchCloudRecords(account.value)
    appendRecords(cloud.records)
    const cloudBatches = await fetchCloudImportBatches(account.value)
    const mergedBatches = mergeImportBatches(cloudBatches.batches)
    await uploadCloudImportBatches(account.value, mergedBatches)
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
          const cloudBatches = await fetchCloudImportBatches(account.value)
          const mergedBatches = mergeImportBatches(cloudBatches.batches)
          await uploadCloudImportBatches(account.value, mergedBatches)
          recordCount.value = getRecords().length
          uni.showToast({ title: '云同步已开启', icon: 'success' })
        } catch (error) {
          cloudEnabled.value = getCloudSyncEnabled()
          uni.showToast({ title: error instanceof Error ? error.message : '同步失败，请重试', icon: 'none' })
        }
      })()
    } })
  } else {
    cloudEnabled.value = true
    cloudDisableConfirmOpen.value = true
  }
}
function cancelDisableCloudSync() {
  cloudDisableConfirmOpen.value = false
  cloudEnabled.value = true
}
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
    recordCount.value = 0
    uni.showToast({ title: '本机账单已清除', icon: 'success' })
  } })
}
</script>

<style scoped>
.page{padding:30rpx 32rpx 170rpx}.profile-card,.data-card{background:#fff;border-radius:24rpx;padding:25rpx 23rpx}.profile-card{position:relative;display:block;overflow:hidden;margin-bottom:32rpx;padding:0;box-shadow:0 12rpx 34rpx rgba(54,45,92,.09)}.profile-cover{position:relative;height:108rpx;overflow:hidden;padding:20rpx 24rpx;background:linear-gradient(120deg,var(--theme-primary,#6556b2),#9184d2);box-sizing:border-box}.cover-mark{position:relative;z-index:1;color:rgba(255,255,255,.86);font-size:21.3rpx;letter-spacing:2rpx}.cover-orb{position:absolute;border-radius:50%;background:rgba(255,255,255,.1)}.orb-one{width:160rpx;height:160rpx;right:24rpx;top:-90rpx}.orb-two{width:86rpx;height:86rpx;right:118rpx;top:56rpx;background:rgba(255,255,255,.08)}.profile-main{position:relative;display:flex;align-items:center;gap:20rpx;padding:0 24rpx 18rpx}.avatar-wrap{position:relative;margin-top:-34rpx;flex:none}.avatar{width:104rpx;height:104rpx;border-radius:50%;background:#fff;display:flex;align-items:center;justify-content:center;flex:none;border:6rpx solid #fff;box-sizing:content-box}.profile-copy{flex:1;min-width:0;display:flex;flex-direction:column;gap:8rpx}.profile-title{font-size:32.5rpx;font-weight:750;color:#302c45}.nickname-line{display:flex;align-items:center;gap:8rpx}.edit-hint{font-size:19rpx;color:#9994aa}.avatar-edit-badge{position:absolute;right:-2rpx;bottom:0;width:32rpx;height:32rpx;border-radius:50%;display:flex;align-items:center;justify-content:center;background:var(--theme-primary,#6556b2);border:3rpx solid #fff}.profile-quick-info{display:flex;align-items:center;gap:24rpx;padding:16rpx 24rpx 20rpx;border-top:1rpx solid #f1eff5}.quick-item{display:flex;flex-direction:column;gap:5rpx}.quick-label{font-size:17.9rpx;color:#9a96a9}.quick-value{font-size:21.3rpx;color:#49445d;font-weight:650}.quick-divider{height:36rpx;width:1rpx;background:#e9e6ef}.profile-card-hint{margin-left:auto;font-size:19rpx;color:var(--theme-primary,#6556b2);background:var(--theme-soft,#efedf8);padding:8rpx 13rpx;border-radius:999rpx}.profile-sub,.row-sub{font-size:21.3rpx;color:#9894a5;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}.editor-mask{position:fixed;z-index:1000;inset:0;background:rgba(36,30,59,.42);display:flex;align-items:flex-end}.profile-editor,.theme-picker{width:100%;padding:34rpx 34rpx calc(34rpx + env(safe-area-inset-bottom));border-radius:30rpx 30rpx 0 0;background:#faf9fd;box-shadow:0 -14rpx 48rpx rgba(36,30,59,.16)}.editor-heading{display:flex;align-items:center;justify-content:space-between;margin-bottom:28rpx}.editor-title{font-size:33.6rpx;font-weight:700;color:#302c45}.editor-close{font-size:47rpx;color:#9aa39e;padding:0 8rpx}.editor-label{display:block;margin:18rpx 0 10rpx;font-size:22.4rpx;font-weight:600;color:#68647b}.editor-input{height:88rpx;padding:0 22rpx;border:1rpx solid #e9e7f0;border-radius:18rpx;background:#fff;font-size:25.8rpx;color:#302c45}.editor-note{display:block;margin:20rpx 0 24rpx;color:#919b95;font-size:20.2rpx}.section-title{display:flex;align-items:center;gap:10rpx;font-size:25.8rpx;font-weight:700;color:#5b566e;letter-spacing:0;margin:30rpx 6rpx 13rpx}.section-index{display:inline-flex;align-items:center;justify-content:center;width:38rpx;height:30rpx;border-radius:9rpx;background:var(--theme-soft,#efedf8);color:var(--theme-primary,#6556b2);font-size:16.8rpx;font-weight:700;letter-spacing:0}.section-caption{margin-left:auto;color:#aaa6b5;font-size:17.9rpx;font-weight:400;letter-spacing:0}.data-card{padding:0 22rpx;overflow:hidden;box-shadow:0 8rpx 24rpx rgba(54,45,92,.035)}.menu-row{min-height:104rpx;display:flex;align-items:center;gap:16rpx}.row-icon{width:54rpx;height:54rpx;border-radius:17rpx;display:flex;align-items:center;justify-content:center;font-size:23.5rpx;font-weight:650;flex:none}.profile-icon{background:var(--theme-soft, #efedf8);color:var(--theme-primary, #6556b2)}.account-icon{background:#eef1ff;color:#6674b8}.local{background:var(--theme-soft, #efedf8);color:var(--theme-primary, #6556b2)}.cloud{background:#edf2ff;color:#5870b9}.sync-icon{background:var(--theme-soft, #efedf8);color:var(--theme-primary, #6556b2);font-size:33.6rpx}.records-icon{background:#fff4e7;color:#a67841}.clear-icon{background:#fff0ef;color:#c25750}.theme-icon{background:var(--theme-soft, #efedf8);color:var(--theme-primary, #6556b2)}.row-copy{flex:1;min-width:0;display:flex;flex-direction:column;gap:7rpx}.row-title{font-size:24.6rpx;font-weight:600;color:#303a35}.row-sub{font-size:20.2rpx}.row-action,.theme-current{font-size:21.3rpx;color:var(--theme-primary, #6556b2);flex:none}.theme-current{display:flex;align-items:center;gap:8rpx}.theme-dot{width:18rpx;height:18rpx;border-radius:50%}.row-status{font-size:21.3rpx;color:#418b70;flex:none}.row-status.muted{color:#a1a8a4}.row-divider{height:1rpx;background:#efedf3;margin-left:70rpx}.info-row{padding:18rpx 4rpx 20rpx;color:#8b879d;font-size:20.2rpx;line-height:1.65}.action-row{border-top:1rpx solid #efedf3}.danger{font-size:21.3rpx;color:#c25750;flex:none}.footnote{font-size:20.2rpx;color:#a0a7a3;line-height:1.7;text-align:center;padding:28rpx 14rpx}.theme-mask{z-index:1100}.theme-picker{box-sizing:border-box;max-height:82vh;overflow-y:auto}.theme-option{min-height:94rpx;display:flex;align-items:center;gap:18rpx;padding:14rpx 16rpx;border:1rpx solid #eeecf4;border-radius:18rpx;margin-top:12rpx;background:#fff}.theme-option.active{border-color:var(--theme-primary, #6556b2);background:var(--theme-soft, #efedf8)}.theme-swatches{display:flex;width:76rpx}.theme-swatches view{width:30rpx;height:30rpx;border:3rpx solid #fff;border-radius:50%;margin-left:-7rpx}.theme-swatches view:first-child{margin-left:0}.theme-label{flex:1;font-size:24.6rpx;color:#403b55}
.editor-mask{box-sizing:border-box;padding:0;z-index:1200}.avatar-picker-sheet,.nickname-sheet{box-sizing:border-box;width:100vw;max-width:none;flex:0 0 100vw;margin:0;padding:34rpx 30rpx calc(28rpx + env(safe-area-inset-bottom));border-radius:30rpx 30rpx 0 0;background:#faf9fd;box-shadow:0 -14rpx 48rpx rgba(36,30,59,.16)}.avatar-picker-sheet{max-height:78vh;overflow-y:auto}.nickname-sheet{padding-bottom:calc(30rpx + env(safe-area-inset-bottom))}.avatar-grid{grid-template-columns:repeat(4,minmax(0,1fr));gap:10rpx}.avatar-choice{min-height:112rpx}.avatar-choice .avatar-sample{width:68rpx;height:68rpx}.sheet-subtitle{display:block;margin-top:7rpx;font-size:20.2rpx;color:#9691a5}.editor-heading{margin-bottom:22rpx}.nickname-sheet .editor-input{box-sizing:border-box;width:100%;margin:6rpx 0 0}.nickname-sheet .editor-note{line-height:1.5}.confirm-mask{align-items:center;justify-content:center;padding:32rpx}.confirm-dialog{box-sizing:border-box;width:100%;max-width:640rpx;padding:38rpx 34rpx 30rpx;border-radius:26rpx;background:#fff;box-shadow:0 20rpx 70rpx rgba(36,30,59,.2)}.confirm-title{display:block;text-align:center;font-size:33.6rpx;font-weight:700;color:#302c45}.confirm-message{display:block;margin:22rpx 0 30rpx;text-align:center;font-size:24.6rpx;line-height:1.65;color:#858194}.confirm-actions{display:flex;flex-direction:column;gap:12rpx}.confirm-cancel,.confirm-danger{box-sizing:border-box;width:100%;min-height:82rpx;margin:0;padding:18rpx 24rpx;border:0;border-radius:16rpx;font-size:24.6rpx;line-height:1.4;white-space:normal}.confirm-cancel{background:#f3f2f7;color:#5f5b70}.confirm-danger{background:#c25750;color:#fff}.confirm-cancel::after,.confirm-danger::after{border:0}.avatar{overflow:hidden;padding:0}.avatar-image{display:block;width:100%;height:100%;border-radius:50%}.profile-editor{box-sizing:border-box;max-height:90vh;overflow-y:auto}.avatar-edit-row{display:flex;align-items:center;gap:16rpx;padding:12rpx 0 8rpx}.avatar-preview{width:84rpx;height:84rpx;font-size:34.7rpx}.avatar-edit-copy{flex:1;min-width:0;display:flex;flex-direction:column;gap:7rpx}.avatar-reset{display:block;text-align:right;margin:8rpx 2rpx 12rpx;color:var(--theme-primary, #6556b2);font-size:20.2rpx}.avatar-grid{display:grid;grid-template-columns:repeat(4,minmax(0,1fr));gap:12rpx}.avatar-choice{position:relative;box-sizing:border-box;min-width:0;min-height:126rpx;display:flex;flex-direction:column;align-items:center;justify-content:center;gap:7rpx;border:1rpx solid #eeecf4;border-radius:16rpx;background:#fff;color:#777286;font-size:19rpx}.avatar-choice.selected{border-color:var(--theme-primary, #6556b2);background:var(--theme-soft, #efedf8);color:var(--theme-primary, #6556b2)}.avatar-choice .avatar-sample{width:74rpx;height:74rpx;border-radius:23rpx;overflow:hidden}.avatar-default{display:flex;align-items:center;justify-content:center;background:var(--theme-soft, #efedf8);color:var(--theme-primary, #6556b2);font-size:30.2rpx;font-weight:700}.avatar-choice .wd-icon{position:absolute;right:7rpx;top:6rpx;color:var(--theme-primary, #6556b2)}.account-delete-dialog{box-sizing:border-box;width:min(100%,640rpx);padding:34rpx;border-radius:28rpx;background:#fff;box-shadow:0 20rpx 70rpx rgba(36,30,59,.2)}.delete-warning{display:block;padding:20rpx;border-radius:16rpx;background:#fff4f1;color:#a45b52;font-size:22.4rpx;line-height:1.65}.delete-actions{display:flex;gap:14rpx;margin-top:24rpx}.delete-actions .wd-button{flex:1}.delete-account-icon{background:#fff0ef;color:#c25750}.danger-title{color:#c25750}.inviter-line{display:block;margin-top:5rpx}.frame-label{margin-top:22rpx}.frame-options{display:flex;flex-wrap:wrap;gap:12rpx;margin-bottom:22rpx}.frame-option{position:relative;min-width:190rpx;box-sizing:border-box;display:flex;align-items:center;gap:14rpx;padding:12rpx 18rpx;border:1rpx solid #eeecf4;border-radius:16rpx;background:#fff;color:#777286;font-size:20.2rpx}.frame-option.selected{border-color:var(--theme-primary,#6556b2);background:var(--theme-soft,#efedf8);color:var(--theme-primary,#6556b2)}.frame-option>wd-icon{margin-left:auto}.frame-preview{width:54rpx;height:54rpx;padding:3rpx;border:3rpx solid #dedbe9;border-radius:50%;box-sizing:border-box}.frame-preview image{display:block;width:100%;height:100%;border-radius:50%}.frame-sparkle{border-color:#e5b943!important;box-shadow:0 0 0 3rpx #fff,0 0 0 6rpx #d9a931!important}.frame-locked{width:100%;display:flex;align-items:center;gap:10rpx;padding:16rpx;border-radius:14rpx;background:#f3f1f8;color:#8d879e;font-size:19rpx;line-height:1.45}.invite-card{padding:24rpx}.invite-summary{display:flex;align-items:center;justify-content:space-between}.invite-count{display:flex;align-items:baseline;gap:10rpx}.invite-number{font-size:49.3rpx;font-weight:750;color:var(--theme-primary,#6556b2)}.invite-caption{font-size:21.3rpx;color:#888499}.invite-reward{display:flex;align-items:center;gap:8rpx;padding:12rpx 16rpx;border-radius:999rpx;background:#fff8e7;color:#9a7220;font-size:20.2rpx;font-weight:650}.invite-progress{display:flex;flex-direction:column;gap:10rpx;margin:18rpx 0 20rpx;color:#8c879a;font-size:19rpx}.progress-track{height:10rpx;overflow:hidden;border-radius:999rpx;background:#efedf4}.progress-fill{height:100%;border-radius:inherit;background:linear-gradient(90deg,#d9a931,#f0ce70);transition:width .2s ease}.share-invite-button{box-sizing:border-box;width:100%;height:78rpx;margin:0;padding:0 20rpx;border:0;border-radius:17rpx;background:linear-gradient(110deg,var(--theme-gradient-start,#332e5d),var(--theme-primary,#6556b2));color:#fff;display:flex;align-items:center;justify-content:center;gap:10rpx;font-size:23.5rpx;font-weight:650}.share-invite-button::after{border:0}.login-invite{justify-content:space-between}
/* 账号页文字层级 */
.profile-title{font-size:37rpx!important}.profile-sub{font-size:23.5rpx!important}.edit-hint{font-size:21.3rpx!important}.quick-label{font-size:21.3rpx!important}.quick-value{font-size:25.8rpx!important}.profile-card-hint{font-size:21.3rpx!important}.section-title{font-size:30.2rpx!important}.row-title{font-size:28rpx!important}.row-sub{font-size:22.4rpx!important}.row-action,.danger,.theme-current{font-size:23.5rpx!important}.invite-caption,.invite-progress,.invite-reward{font-size:22.4rpx!important}
.password-dialog{box-sizing:border-box;width:min(100%,640rpx);padding:34rpx;border-radius:28rpx;background:#fff;box-shadow:0 20rpx 70rpx rgba(36,30,59,.2)}
</style>


.avatar-choice .avatar-sample{width:78rpx;height:78rpx;border-radius:23rpx;overflow:hidden}.avatar-default-logo{background:#fff}.avatar-picker-sheet,.nickname-sheet{box-sizing:border-box;width:100%;padding:34rpx 34rpx calc(34rpx + env(safe-area-inset-bottom));border-radius:30rpx 30rpx 0 0;background:#faf9fd;box-shadow:0 -14rpx 48rpx rgba(36,30,59,.16)}.avatar-picker-sheet{max-height:84vh;overflow-y:auto}.sheet-subtitle{display:block;margin-top:7rpx;font-size:20.2rpx;color:#9691a5}.profile-main .avatar-image{display:block;width:100%;height:100%;border-radius:50%}








