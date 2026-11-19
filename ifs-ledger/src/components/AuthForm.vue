<template>
  <view class="auth-page">
    <view class="brand"><image class="brand-logo" src="/static/ledger-logo.svg" mode="aspectFit" /><text class="brand-name">轻记</text></view>
    <view class="hero"><text class="eyebrow">PERSONAL FINANCE</text><text class="headline">{{ isRegister ? '从今天开始，轻松记账' : '欢迎回来' }}</text><text class="subtitle">{{ isRegister ? '创建账号是可选的，本机记账无需注册' : '登录后可选择是否开启云端同步' }}</text></view>
    <view class="form-card">
      <text class="form-title">{{ isRegister ? '创建账号' : '账号登录' }}</text>
      <text class="form-caption">使用手机号或邮箱{{ isRegister ? '注册' : '登录' }}</text>
      <view v-if="isRegister && inviterUid" class="inviter-note"><wd-icon name="user" size="15" color="var(--theme-primary, #6556b2)" /><text>邀请人 UID：{{ inviterUid }}</text></view>
      <wd-input v-model="loginId" label="账号" placeholder="手机号 / 邮箱" prefix-icon="user" clearable />
      <wd-input v-model="password" label="密码" placeholder="请输入密码" type="password" prefix-icon="lock" clearable />
      <wd-input v-if="isRegister" v-model="confirmPassword" label="确认密码" placeholder="请再次输入密码" type="password" prefix-icon="lock" clearable />
      <view v-if="isRegister" class="captcha-row"><wd-input v-model="captchaCode" label="验证码" placeholder="请输入验证码" prefix-icon="safe" clearable /><image v-if="captchaImage" class="captcha-image" :src="captchaImage" mode="aspectFit" @tap="loadCaptcha" /></view>
      <text v-if="isRegister" class="password-hint">密码长度至少 8 位</text>
      <view v-if="errorMessage" class="error-message">{{ errorMessage }}</view>
      <wd-button block type="primary" :loading="loading" @click="submit">{{ loading ? '请稍候…' : (isRegister ? '注册并登录' : '登录') }}</wd-button>
      <view class="switch-line"><text>{{ isRegister ? '已经有账号？' : '还没有账号？' }}</text><text class="switch-link" @tap="switchMode">{{ isRegister ? '去登录' : '立即注册' }}</text></view>
    </view>
    <view class="local-option" @tap="continueLocally"><text class="local-icon">◉</text><text>暂不登录，继续本地使用</text><text class="arrow">→</text></view>
    <text class="privacy-note">新账号默认开启云同步，可随时关闭。</text>
  </view>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { appendRecords, authenticate, fetchCloudImportBatches, fetchCloudRecords, fetchRegisterCaptcha, mergeImportBatches, uploadCloudImportBatches, uploadCloudRecords } from '../utils/ledger'

const props = defineProps<{ mode: 'login' | 'register'; inviterUid?: string }>()
const isRegister = props.mode === 'register'
const loginId = ref('')
const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const errorMessage = ref('')
const captchaID = ref('')
const captchaCode = ref('')
const captchaImage = ref('')

onMounted(() => { if (isRegister) void loadCaptcha() })

async function loadCaptcha() {
  try {
    const captcha = await fetchRegisterCaptcha()
    captchaID.value = captcha.uuid
    captchaImage.value = captcha.img.startsWith('data:') ? captcha.img : `data:image/png;base64,${captcha.img}`
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '验证码加载失败'
  }
}

function switchMode() {
  uni.redirectTo({ url: isRegister ? '/pages/login/index' : '/pages/register/index' })
}

function continueLocally() {
  uni.navigateBack({ delta: 1, fail: () => uni.reLaunch({ url: '/pages/index/index' }) })
}

async function submit() {
  errorMessage.value = ''
  const id = loginId.value.trim()
  if (!id || !password.value) {
    errorMessage.value = '请输入账号和密码'
    return
  }
  if (isRegister && !/^(?:\+?86)?1[3-9]\d{9}$/.test(id) && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(id)) {
    errorMessage.value = '请输入正确的手机号或邮箱'
    return
  }
  if (isRegister && password.value.length < 8) {
    errorMessage.value = '密码长度至少为 8 位'
    return
  }
  if (isRegister && password.value !== confirmPassword.value) {
    errorMessage.value = '两次输入的密码不一致'
    return
  }
  if (isRegister && (!captchaID.value || !captchaCode.value.trim())) {
    errorMessage.value = '请输入验证码'
    return
  }
  loading.value = true
  try {
    const result = await authenticate(id, password.value, isRegister, props.inviterUid || '', isRegister ? { captchaId: captchaID.value, captchaCode: captchaCode.value.trim() } : undefined)
    if (result.cloudSyncEnabled) {
      const cloud = await fetchCloudRecords(result.account)
      appendRecords(cloud.records)
      await uploadCloudRecords(result.account)
      const cloudBatches = await fetchCloudImportBatches(result.account)
      const mergedBatches = mergeImportBatches(cloudBatches.batches)
      await uploadCloudImportBatches(result.account, mergedBatches)
    }
    uni.showToast({ title: isRegister ? '账号已创建' : (result.cloudSyncEnabled ? '登录并同步完成' : '登录成功'), icon: 'success' })
    setTimeout(() => uni.reLaunch({ url: '/pages/index/index' }), 450)
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '操作失败，请稍后重试'
    if (isRegister) { captchaCode.value = ''; void loadCaptcha() }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-page{min-height:100vh;padding:46rpx 38rpx 56rpx;background:linear-gradient(180deg,#f0eef8 0,var(--theme-page, #f6f5fa) 43%,var(--theme-page, #f6f5fa) 100%);display:flex;flex-direction:column}.brand{display:flex;align-items:center;gap:14rpx}.brand-logo{width:62rpx;height:62rpx;border-radius:20rpx;box-shadow:0 8rpx 18rpx #40387428}.brand-name{font-size:28rpx;font-weight:700;color:var(--theme-primary, #51458e)}.hero{display:flex;flex-direction:column;margin:82rpx 0 38rpx}.eyebrow{font-size:20.2rpx;letter-spacing:4rpx;color:#9186c5}.headline{font-size:48.2rpx;font-weight:750;color:#29263a;margin-top:15rpx;line-height:1.35}.subtitle{font-size:24.6rpx;color:#8b9690;margin-top:13rpx}.form-card{background:#fff;border-radius:28rpx;padding:34rpx 27rpx 28rpx;box-shadow:0 14rpx 45rpx #42367a0d}.form-title{font-size:32.5rpx;font-weight:700}.form-caption{display:block;color:#a0a7a2;font-size:22.4rpx;margin:8rpx 0 20rpx}.password-hint{display:block;color:#9da5a0;font-size:20.2rpx;margin:10rpx 0 18rpx}.error-message{padding:15rpx 16rpx;margin:15rpx 0;border-radius:14rpx;background:#fff1ef;color:#bd554d;font-size:21.3rpx}.switch-line{display:flex;justify-content:center;gap:8rpx;margin-top:25rpx;color:#9ba39e;font-size:22.4rpx}.switch-link{color:var(--theme-primary, #6556b2);font-weight:650}.local-option{margin-top:26rpx;padding:23rpx 22rpx;border-radius:20rpx;background:var(--theme-soft, #efedf8);display:flex;align-items:center;gap:14rpx;color:var(--theme-primary, #5a4d9c);font-size:23.5rpx;font-weight:600}.local-icon{font-size:29.1rpx}.arrow{margin-left:auto;font-size:28rpx}.privacy-note{margin-top:auto;padding-top:36rpx;text-align:center;color:#9ca49f;font-size:20.2rpx;line-height:1.7}
</style>

<style scoped>
.inviter-note{display:flex;align-items:center;gap:8rpx;padding:12rpx 14rpx;margin:-6rpx 0 16rpx;border-radius:13rpx;background:var(--theme-soft,#efedf8);color:var(--theme-primary,#6556b2);font-size:20.2rpx}
.captcha-row{display:flex;align-items:center;gap:14rpx;margin:8rpx 0}.captcha-row :deep(.wd-input){flex:1;min-width:0}.captcha-image{width:190rpx;height:78rpx;border:1rpx solid #e9e7f0;border-radius:12rpx;background:#fff;flex:none}
</style>
