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
              <a-tag color="blue">{{ profile.isMain === '1' ? '主账号' : '子账号' }}</a-tag>
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
  </main>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import {
  type WorkspaceAccount,
  getWorkspaceProfile,
  getWorkspaceProfileCache,
  normalizeWorkspaceProfile,
  setWorkspaceProfileCache,
} from '@/api/workspace/auth'

const loading = ref(true)
const profile = ref<WorkspaceAccount>()

const initials = computed(() => {
  const raw = profile.value?.realName || profile.value?.customerName || 'CU'
  return raw.slice(0, 2).toUpperCase()
})

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
.profile-shell {
  display: grid;
  grid-template-columns: 320px minmax(0, 1fr);
  gap: 20px;
}

.profile-aside,
.profile-main {
  padding: 26px;
  border-radius: 22px;
  border: 1px solid rgba(148, 163, 184, 0.18);
  background: #fff;
  box-shadow: 0 20px 42px rgba(15, 23, 42, 0.06);
}

.profile-aside {
  background:
    radial-gradient(circle at top right, rgba(255, 255, 255, 0.12), transparent 28%),
    linear-gradient(180deg, #0f172a, #1e293b);
  color: #fff;
}

.avatar-badge {
  width: 72px;
  height: 72px;
  border-radius: 22px;
  display: grid;
  place-items: center;
  background: linear-gradient(135deg, #f8fafc, #93c5fd);
  color: #0f172a;
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
  color: #64748b;
}

.profile-aside span {
  display: block;
  margin-top: 10px;
  color: rgba(255, 255, 255, 0.8);
}

.profile-aside p {
  margin: 10px 0 0;
  color: rgba(226, 232, 240, 0.68);
  line-height: 1.8;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding-bottom: 20px;
  border-bottom: 1px solid rgba(148, 163, 184, 0.16);
}

.section-header small,
.detail-row span {
  letter-spacing: 0.08em;
  font-size: 11px;
}

.section-header h2 {
  font-size: 28px;
  color: #0f172a;
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
  background: linear-gradient(180deg, #f8fbff, #f1f5f9);
}

.detail-row span,
.detail-row strong {
  display: block;
}

.detail-row strong {
  margin-top: 10px;
  color: #0f172a;
  font-size: 20px;
  overflow-wrap: anywhere;
}

@media (max-width: 980px) {
  .profile-shell,
  .detail-grid {
    grid-template-columns: 1fr;
  }
}
</style>
