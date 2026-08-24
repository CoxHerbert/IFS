<template>
  <main class="account-page">
    <a-spin :spinning="loading">
      <template v-if="profile">
        <section class="profile-shell">
          <article class="profile-aside">
            <div class="avatar-badge">{{ initials }}</div>
            <h1>{{ profile.realName || profile.customerName || '客户账号' }}</h1>
            <span>{{ profile.username }}</span>
            <p>{{ profile.companyName || '当前账号尚未维护公司名称' }}</p>
          </article>

          <article class="profile-main">
            <div class="section-header">
              <div>
                <small>账号资料</small>
                <h2>客户端账号资料</h2>
              </div>
              <div class="section-actions">
                <a-tag color="default">{{ profile.isMain === '1' ? '主账号' : '子账号' }}</a-tag>
                <a-button @click="openProfileEditor">编辑资料</a-button>
                <a-button type="primary" ghost @click="openPasswordEditor">更新密码</a-button>
              </div>
            </div>

            <div class="detail-grid">
              <div class="detail-row">
                <span>客户编号</span>
                <strong>{{ profile.customerNo || '-' }}</strong>
              </div>
              <div class="detail-row">
                <span>客户名称</span>
                <strong>{{ profile.customerName || '-' }}</strong>
              </div>
              <div class="detail-row">
                <span>联系人姓名</span>
                <strong>{{ profile.realName || '-' }}</strong>
              </div>
              <div class="detail-row">
                <span>联系电话</span>
                <strong>{{ profile.phone || '-' }}</strong>
              </div>
              <div class="detail-row">
                <span>邮箱地址</span>
                <strong>{{ profile.email || '-' }}</strong>
              </div>
              <div class="detail-row">
                <span>账号状态</span>
                <strong>{{ profile.status === '0' ? '正常' : '停用' }}</strong>
              </div>
            </div>
          </article>
        </section>
      </template>
    </a-spin>

    <a-modal
      v-model:open="profileModalOpen"
      title="编辑资料"
      :confirm-loading="submittingProfile"
      ok-text="保存"
      cancel-text="取消"
      @ok="submitProfileUpdate"
      @cancel="resetProfileForm"
    >
      <a-form layout="vertical">
        <a-form-item label="联系人姓名" required>
          <a-input v-model:value="profileForm.realName" placeholder="请输入联系人姓名" maxlength="30" />
        </a-form-item>
        <a-form-item label="联系电话">
          <a-input v-model:value="profileForm.phone" placeholder="请输入联系电话" maxlength="20" />
        </a-form-item>
        <a-form-item label="邮箱地址">
          <a-input v-model:value="profileForm.email" placeholder="请输入邮箱地址" maxlength="64" />
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal
      v-model:open="passwordModalOpen"
      title="更新密码"
      :confirm-loading="submittingPassword"
      ok-text="确认更新"
      cancel-text="取消"
      @ok="submitPasswordUpdate"
      @cancel="resetPasswordForm"
    >
      <a-form layout="vertical">
        <a-form-item label="旧密码" required>
          <a-input-password v-model:value="passwordForm.oldPassword" placeholder="请输入当前密码" />
        </a-form-item>
        <a-form-item label="新密码" required>
          <a-input-password v-model:value="passwordForm.newPassword" placeholder="请输入新密码" />
        </a-form-item>
        <a-form-item label="确认新密码" required>
          <a-input-password v-model:value="passwordForm.confirmPassword" placeholder="请再次输入新密码" />
        </a-form-item>
      </a-form>
    </a-modal>
  </main>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { message as antMessage } from 'ant-design-vue'
import {
  type WorkspaceAccount,
  getWorkspaceProfile,
  getWorkspaceProfileCache,
  normalizeWorkspaceProfile,
  setWorkspaceProfileCache,
  updateWorkspacePassword,
  updateWorkspaceProfile,
} from '@/api/workspace/auth'

const loading = ref(true)
const profile = ref<WorkspaceAccount>()

const profileModalOpen = ref(false)
const passwordModalOpen = ref(false)
const submittingProfile = ref(false)
const submittingPassword = ref(false)

const profileForm = reactive({
  realName: '',
  phone: '',
  email: '',
})

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})

const initials = computed(() => {
  const raw = profile.value?.realName || profile.value?.customerName || 'CU'
  return raw.slice(0, 2).toUpperCase()
})

function syncProfileCache(user: WorkspaceAccount) {
  profile.value = user
  const cachedProfile = getWorkspaceProfileCache()
  setWorkspaceProfileCache({
    user,
    roles: cachedProfile?.roles || [],
    permissions: cachedProfile?.permissions || [],
  })
}

function openProfileEditor() {
  profileForm.realName = profile.value?.realName || ''
  profileForm.phone = profile.value?.phone || ''
  profileForm.email = profile.value?.email || ''
  profileModalOpen.value = true
}

function resetProfileForm() {
  profileModalOpen.value = false
  profileForm.realName = ''
  profileForm.phone = ''
  profileForm.email = ''
}

function openPasswordEditor() {
  passwordModalOpen.value = true
}

function resetPasswordForm() {
  passwordModalOpen.value = false
  passwordForm.oldPassword = ''
  passwordForm.newPassword = ''
  passwordForm.confirmPassword = ''
}

function isEmailValid(email: string) {
  if (!email) return true
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)
}

async function submitProfileUpdate() {
  const realName = profileForm.realName.trim()
  const phone = profileForm.phone.trim()
  const email = profileForm.email.trim()

  if (!realName) {
    antMessage.warning('请输入联系人姓名')
    return
  }
  if (email && !isEmailValid(email)) {
    antMessage.warning('邮箱格式不正确')
    return
  }

  submittingProfile.value = true
  try {
    const response = await updateWorkspaceProfile({ realName, phone, email })
    if (response.code !== 200 || !response.data) {
      throw new Error(response.msg || '保存失败')
    }
    syncProfileCache(response.data)
    antMessage.success('资料已更新')
    resetProfileForm()
  } catch (error) {
    antMessage.error(error instanceof Error ? error.message : '保存失败')
  } finally {
    submittingProfile.value = false
  }
}

async function submitPasswordUpdate() {
  if (!passwordForm.oldPassword || !passwordForm.newPassword || !passwordForm.confirmPassword) {
    antMessage.warning('请填写完整的密码信息')
    return
  }
  if (passwordForm.newPassword.length < 6) {
    antMessage.warning('新密码至少需要 6 位')
    return
  }
  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    antMessage.warning('两次输入的新密码不一致')
    return
  }
  if (passwordForm.oldPassword === passwordForm.newPassword) {
    antMessage.warning('新密码不能与旧密码相同')
    return
  }

  submittingPassword.value = true
  try {
    const response = await updateWorkspacePassword({ ...passwordForm })
    if (response.code !== 200) {
      throw new Error(response.msg || '更新失败')
    }
    antMessage.success('密码已更新')
    resetPasswordForm()
  } catch (error) {
    antMessage.error(error instanceof Error ? error.message : '更新失败')
  } finally {
    submittingPassword.value = false
  }
}

onMounted(async () => {
  const cachedProfile = getWorkspaceProfileCache()
  if (cachedProfile?.user) {
    profile.value = cachedProfile.user
    loading.value = false
    return
  }
  try {
    const response = await getWorkspaceProfile()
    const normalizedProfile = normalizeWorkspaceProfile(response.data)
    if (response.code === 200 && normalizedProfile?.user) {
      setWorkspaceProfileCache(normalizedProfile)
      profile.value = normalizedProfile.user
    }
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.account-page {
  min-height: 100%;
}

.profile-shell {
  display: grid;
  grid-template-columns: 320px minmax(0, 1fr);
  gap: 20px;
}

.profile-aside,
.profile-main {
  padding: 26px;
  border-radius: 22px;
  border: 1px solid rgba(0, 0, 0, 0.1);
  background: #fff;
  box-shadow: 0 12px 28px rgba(0, 0, 0, 0.04);
}

.profile-aside {
  color: #111111;
}

.avatar-badge {
  width: 72px;
  height: 72px;
  border-radius: 22px;
  display: grid;
  place-items: center;
  background: #f3f4f6;
  color: #111111;
  border: 1px solid rgba(0, 0, 0, 0.08);
  font-size: 26px;
  font-weight: 800;
}

.profile-aside h1,
.section-header h2 {
  margin: 18px 0 0;
}

.profile-aside span,
.profile-aside p,
.section-header small,
.detail-row span {
  color: #6b7280;
}

.profile-aside span {
  display: block;
  margin-top: 10px;
  color: #4b5563;
}

.profile-aside p {
  margin: 10px 0 0;
  color: #6b7280;
  line-height: 1.8;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding-bottom: 20px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.08);
}

.section-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.section-actions :deep(.ant-btn) {
  border-color: #111111;
  color: #111111;
}

.section-actions :deep(.ant-btn-primary.ant-btn-background-ghost) {
  border-color: #111111;
  color: #111111;
}

.section-actions :deep(.ant-btn:hover) {
  border-color: #000000;
  color: #000000;
}

.section-header small,
.detail-row span {
  letter-spacing: 0.08em;
  font-size: 11px;
}

.section-header h2 {
  font-size: 28px;
  color: #111111;
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 18px;
  margin-top: 22px;
}

.detail-row {
  padding: 18px;
  border-radius: 16px;
  border: 1px solid rgba(0, 0, 0, 0.06);
  background: #ffffff;
}

.detail-row span,
.detail-row strong {
  display: block;
}

.detail-row strong {
  margin-top: 10px;
  color: #111111;
  font-size: 20px;
  overflow-wrap: anywhere;
}

@media (max-width: 980px) {
  .profile-shell,
  .detail-grid {
    grid-template-columns: 1fr;
  }

  .section-header {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
