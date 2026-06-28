<template>
  <main class="center-page">
    <section class="center-header">
      <div>
        <p>客户中心</p>
        <h2>{{ profile?.customerName || '我的客户资料' }}</h2>
        <span>{{ profile?.companyName || '登录后可查看客户账号信息' }}</span>
      </div>
      <a-button @click="logout">退出登录</a-button>
    </section>

    <a-spin :spinning="loading">
      <section v-if="profile" class="info-grid">
        <div class="info-item">
          <span>客户编号</span>
          <strong>{{ profile.customerNo }}</strong>
        </div>
        <div class="info-item">
          <span>登录账号</span>
          <strong>{{ profile.username }}</strong>
        </div>
        <div class="info-item">
          <span>账号姓名</span>
          <strong>{{ profile.realName || '-' }}</strong>
        </div>
        <div class="info-item">
          <span>主账号</span>
          <strong>{{ profile.isMain === '1' ? '是' : '否' }}</strong>
        </div>
        <div class="info-item">
          <span>联系电话</span>
          <strong>{{ profile.phone || '-' }}</strong>
        </div>
        <div class="info-item">
          <span>邮箱</span>
          <strong>{{ profile.email || '-' }}</strong>
        </div>
      </section>

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

function logout() {
  removeCustomerToken()
  router.push('/customer-login')
}
</script>

<style scoped>
.center-page {
  min-height: calc(100vh - 190px);
  padding: 48px 40px 72px;
}

.center-header {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 24px;
  max-width: 1120px;
  margin: 0 auto 24px;
}

.center-header p {
  margin: 0 0 10px;
  color: #1677ff;
  font-weight: 700;
}

.center-header h2 {
  margin: 0;
  font-size: 32px;
}

.center-header span {
  display: block;
  margin-top: 10px;
  color: #66748b;
}

.info-grid {
  max-width: 1120px;
  margin: 0 auto;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
}

.info-item {
  min-height: 108px;
  padding: 20px;
  background: #fff;
  border: 1px solid rgba(16, 35, 63, 0.08);
  border-radius: 8px;
}

.info-item span {
  display: block;
  color: #66748b;
  margin-bottom: 12px;
}

.info-item strong {
  display: block;
  overflow-wrap: anywhere;
  font-size: 18px;
}

@media (max-width: 900px) {
  .center-header {
    align-items: flex-start;
    flex-direction: column;
  }

  .info-grid {
    grid-template-columns: 1fr;
  }
}
</style>
