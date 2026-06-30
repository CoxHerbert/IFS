<template>
  <main class="center-page">
    <a-spin :spinning="loading">
      <template v-if="profile">
        <section class="hero-card">
          <div>
            <p>Account Overview</p>
            <h1>{{ profile.customerName || '我的客户资料' }}</h1>
            <span>{{ profile.companyName || '登录后可查看客户账号信息' }}</span>
          </div>

          <div class="hero-metrics">
            <div>
              <small>客户编号</small>
              <strong>{{ profile.customerNo }}</strong>
            </div>
            <div>
              <small>登录账号</small>
              <strong>{{ profile.username }}</strong>
            </div>
          </div>
        </section>

        <section class="info-grid">
          <article class="info-card accent">
            <span>账号姓名</span>
            <strong>{{ profile.realName || '-' }}</strong>
          </article>
          <article class="info-card">
            <span>主账号</span>
            <strong>{{ profile.isMain === '1' ? '是' : '否' }}</strong>
          </article>
          <article class="info-card">
            <span>联系电话</span>
            <strong>{{ profile.phone || '-' }}</strong>
          </article>
          <article class="info-card">
            <span>邮箱</span>
            <strong>{{ profile.email || '-' }}</strong>
          </article>
        </section>

        <section class="summary-panel">
          <div>
            <p>下一步</p>
            <h3>客户门户基础资料已就绪</h3>
            <span>后续可以在这里接入出货导入、装柜测算、出货计划和历史记录查询。</span>
          </div>
        </section>
      </template>

      <a-result v-else-if="!loading" status="warning" title="未能读取客户信息">
        <template #extra>
          <a-button type="primary" @click="router.push('/customer-login')">重新登录</a-button>
        </template>
      </a-result>
    </a-spin>
  </main>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { type CustomerAccount, getCustomerProfile, removeCustomerToken } from '@/api/customer'

const router = useRouter()
const loading = ref(true)
const profile = ref<CustomerAccount>()

onMounted(async () => {
  try {
    const response = await getCustomerProfile()
    if (response.code !== 200 || !response.data) {
      removeCustomerToken()
      router.push('/customer-login')
      return
    }
    profile.value = response.data
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.center-page {
  min-height: 100%;
}

.hero-card,
.summary-panel,
.info-card {
  border: 1px solid rgba(16, 35, 63, 0.08);
  box-shadow: 0 18px 40px rgba(21, 52, 93, 0.08);
}

.hero-card {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 24px;
  padding: 28px 30px;
  border-radius: 28px;
  background:
    radial-gradient(circle at top right, rgba(255, 255, 255, 0.18), transparent 24%),
    linear-gradient(135deg, #0d315f, #1e6fbe);
}

.hero-card p,
.summary-panel p {
  margin: 0;
  text-transform: uppercase;
  letter-spacing: 0.12em;
  font-size: 11px;
}

.hero-card p {
  color: rgba(255, 255, 255, 0.72);
}

.hero-card h1,
.summary-panel h3 {
  margin: 10px 0 0;
}

.hero-card h1 {
  font-size: clamp(30px, 4vw, 40px);
  color: #fff;
}

.hero-card span {
  display: block;
  margin-top: 10px;
  color: rgba(255, 255, 255, 0.82);
  max-width: 48ch;
}

.hero-metrics {
  display: grid;
  grid-template-columns: repeat(2, minmax(140px, 1fr));
  gap: 14px;
  width: min(100%, 360px);
}

.hero-metrics div {
  padding: 18px;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.12);
  backdrop-filter: blur(10px);
}

.hero-metrics small {
  display: block;
  color: rgba(255, 255, 255, 0.68);
  margin-bottom: 8px;
}

.hero-metrics strong {
  color: #fff;
  font-size: 18px;
  overflow-wrap: anywhere;
}

.info-grid {
  margin-top: 18px;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.info-card {
  min-height: 130px;
  padding: 24px;
  background: #fff;
  border-radius: 22px;
}

.info-card.accent {
  background: linear-gradient(180deg, #ffffff, #f4f8fd);
}

.info-card span {
  display: block;
  color: #66748b;
  margin-bottom: 12px;
}

.info-card strong {
  display: block;
  overflow-wrap: anywhere;
  font-size: 24px;
  line-height: 1.3;
}

.summary-panel {
  margin-top: 18px;
  padding: 24px 28px;
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.88);
}

.summary-panel p {
  color: #5d7496;
}

.summary-panel h3 {
  color: #10233f;
  font-size: 24px;
}

.summary-panel span {
  display: block;
  margin-top: 10px;
  color: #66748b;
  line-height: 1.8;
}

@media (max-width: 900px) {
  .hero-card {
    flex-direction: column;
    align-items: flex-start;
  }

  .info-grid {
    grid-template-columns: 1fr;
  }

  .hero-metrics {
    grid-template-columns: 1fr;
    width: 100%;
  }
}
</style>
