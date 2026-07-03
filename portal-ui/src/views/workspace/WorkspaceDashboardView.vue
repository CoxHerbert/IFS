<template>
  <main class="workspace-page">
    <a-spin :spinning="loading">
      <template v-if="profile">
        <section class="hero-grid">
          <article class="hero-card">
            <p class="hero-kicker">客户协同中心</p>
            <h1>{{ profile.customerName || '客户端工作台' }}</h1>
            <span>{{ profile.companyName || '当前账号尚未维护公司名称' }}</span>
            <div class="hero-actions">
              <router-link to="/customer/account">
                <a-button size="large">查看账号资料</a-button>
              </router-link>
              <router-link to="/customer/shipment">
                <a-button size="large">出货查询</a-button>
              </router-link>
              <router-link to="/customer/shipment-assistant">
                <a-button type="primary" size="large">智能出货助手</a-button>
              </router-link>
            </div>
          </article>

          <article class="hero-side">
            <div>
              <small>客户编号</small>
              <strong>{{ profile.customerNo || '-' }}</strong>
            </div>
            <div>
              <small>登录账号</small>
              <strong>{{ profile.username || '-' }}</strong>
            </div>
          </article>
        </section>

        <section class="metric-grid">
          <article class="metric-card">
            <span>联系人</span>
            <strong>{{ profile.realName || '-' }}</strong>
          </article>
          <article class="metric-card">
            <span>联系电话</span>
            <strong>{{ profile.phone || '-' }}</strong>
          </article>
          <article class="metric-card">
            <span>邮箱地址</span>
            <strong>{{ profile.email || '-' }}</strong>
          </article>
          <article class="metric-card">
            <span>账号类型</span>
            <strong>{{ profile.isMain === '1' ? '主账号' : '子账号' }}</strong>
          </article>
        </section>

        <section class="panel-grid">
          <article class="panel-card">
            <h3>常用入口</h3>
            <ul>
              <li>账号资料维护</li>
              <li>出货进度查询</li>
              <li>Excel 出货测算</li>
            </ul>
          </article>

          <article class="panel-card">
            <h3>当前可做</h3>
            <p>整理货物明细、估算整柜数量、评估散货体积，先把重复计算的工作前移到客户端。</p>
          </article>
        </section>
      </template>
    </a-spin>
  </main>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import {
  type WorkspaceAccount,
  getWorkspaceProfile,
  getWorkspaceProfileCache,
  normalizeWorkspaceProfile,
  setWorkspaceProfileCache,
} from '@/api/workspace/auth'

const loading = ref(true)
const profile = ref<WorkspaceAccount>()

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
.workspace-page {
  min-height: 100%;
}

.hero-grid,
.metric-grid,
.panel-grid {
  display: grid;
  gap: 18px;
}

.hero-grid {
  grid-template-columns: minmax(0, 1.8fr) minmax(280px, 1fr);
}

.hero-card,
.hero-side,
.metric-card,
.panel-card {
  border-radius: 22px;
  border: 1px solid rgba(15, 23, 42, 0.08);
  box-shadow: 0 12px 28px rgba(15, 23, 42, 0.05);
  background: #fff;
}

.hero-card {
  padding: 30px;
  color: #0f172a;
}

.hero-kicker,
.hero-side small,
.metric-card span {
  letter-spacing: 0.08em;
  font-size: 11px;
}

.hero-kicker {
  margin: 0;
  color: #64748b;
  font-weight: 700;
}

.hero-card h1 {
  margin: 14px 0 0;
  font-size: clamp(28px, 4vw, 38px);
}

.hero-card span {
  display: block;
  margin-top: 12px;
  max-width: 40ch;
  color: #475569;
  line-height: 1.8;
}

.hero-actions {
  display: flex;
  gap: 12px;
  margin-top: 24px;
  flex-wrap: wrap;
}

.hero-actions :deep(.ant-btn) {
  min-width: 136px;
  border-radius: 12px;
}

.hero-actions :deep(.ant-btn-primary) {
  background: #111827;
  border-color: #111827;
  box-shadow: none;
}

.hero-side {
  padding: 20px;
  display: grid;
  gap: 16px;
}

.hero-side div {
  padding: 18px;
  border-radius: 16px;
  background: #f8fafc;
}

.hero-side small {
  display: block;
  color: #64748b;
}

.hero-side strong,
.metric-card strong {
  display: block;
  margin-top: 10px;
  color: #1c1917;
  font-size: 22px;
  overflow-wrap: anywhere;
}

.metric-grid {
  margin-top: 18px;
  grid-template-columns: repeat(4, minmax(0, 1fr));
}

.metric-card {
  min-height: 140px;
  padding: 22px;
}

.metric-card span {
  display: block;
  color: #64748b;
}

.panel-grid {
  margin-top: 18px;
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.panel-card {
  padding: 24px;
}

.panel-card h3 {
  margin: 0 0 14px;
  font-size: 20px;
}

.panel-card p,
.panel-card ul {
  margin: 0;
  color: #0f172a;
  line-height: 1.8;
}

.panel-card ul {
  padding-left: 18px;
}

@media (max-width: 1100px) {
  .hero-grid,
  .panel-grid {
    grid-template-columns: 1fr;
  }

  .metric-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 640px) {
  .metric-grid {
    grid-template-columns: 1fr;
  }
}
</style>
