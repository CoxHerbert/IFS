<template>
  <a-layout class="customer-shell" :has-sider="true">
    <a-layout-sider
      class="customer-sidebar"
      :width="280"
      :collapsed-width="0"
      :breakpoint="'lg'"
      theme="light"
    >
      <div class="brand-block">
        <div class="brand-mark">SP</div>
        <div>
          <p>Customer Workspace</p>
          <h1>客户中心</h1>
        </div>
      </div>

      <div class="sidebar-copy">
        <span>国际货运代理</span>
        <strong>登录后查看账号资料与后续业务能力</strong>
      </div>

      <a-menu v-model:selectedKeys="selectedKeys" mode="inline" class="sidebar-menu">
        <a-menu-item key="customer-center">
          <router-link to="/customer-center">我的首页</router-link>
        </a-menu-item>
      </a-menu>

      <div class="sidebar-footer">
        <router-link to="/">
          <a-button block>返回官网</a-button>
        </router-link>
      </div>
    </a-layout-sider>

    <a-layout class="customer-main">
      <a-layout-header class="customer-header">
        <div>
          <p class="eyebrow">Portal Access</p>
          <h2>{{ pageTitle }}</h2>
        </div>

        <a-space :size="12" class="header-actions">
          <router-link to="/contact">
            <a-button>提交需求</a-button>
          </router-link>
          <a-button type="primary" ghost @click="logout">退出登录</a-button>
        </a-space>
      </a-layout-header>

      <a-layout-content class="customer-content">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { removeCustomerToken } from '@/api/customer'

const route = useRoute()
const router = useRouter()

const selectedKeys = computed<string[]>(() => [String(route.name || 'customer-center')])
const pageTitle = computed(() => (route.name === 'customer-center' ? '客户中心' : '客户门户'))

function logout() {
  removeCustomerToken()
  router.push('/customer-login')
}
</script>

<style scoped>
.customer-shell {
  min-height: 100vh;
  background:
    radial-gradient(circle at top left, rgba(11, 79, 167, 0.16), transparent 28%),
    linear-gradient(180deg, #eef4fb 0%, #e7eef8 100%);
}

.customer-sidebar {
  overflow: auto;
  padding: 28px 22px;
  color: #e7efff;
  background:
    linear-gradient(180deg, rgba(8, 24, 49, 0.96), rgba(12, 45, 89, 0.96)),
    linear-gradient(135deg, rgba(82, 166, 255, 0.18), transparent);
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.brand-block {
  display: flex;
  align-items: center;
  gap: 14px;
}

.brand-mark {
  width: 48px;
  height: 48px;
  border-radius: 16px;
  display: grid;
  place-items: center;
  font-size: 18px;
  font-weight: 800;
  color: #0d2a53;
  background: linear-gradient(135deg, #ffffff, #8cc8ff);
}

.brand-block p,
.sidebar-copy span,
.eyebrow {
  margin: 0;
  text-transform: uppercase;
  letter-spacing: 0.12em;
  font-size: 11px;
}

.brand-block p,
.sidebar-copy span {
  color: rgba(231, 239, 255, 0.65);
}

.brand-block h1,
.sidebar-copy strong,
.customer-header h2 {
  margin: 6px 0 0;
}

.brand-block h1 {
  font-size: 24px;
  color: #fff;
}

.sidebar-copy {
  padding: 18px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.04);
}

.sidebar-copy strong {
  display: block;
  color: #fff;
  line-height: 1.6;
  font-size: 18px;
}

.sidebar-menu {
  flex: 1;
  color: #fff;
  background: transparent;
  border-inline-end: 0;
}

.sidebar-menu :deep(.ant-menu-item) {
  border-radius: 12px;
  margin-inline: 0;
}

.sidebar-menu :deep(.ant-menu-item-selected) {
  background: rgba(255, 255, 255, 0.12);
}

.sidebar-menu :deep(.ant-menu-item a) {
  color: #eef5ff;
}

.sidebar-footer {
  padding-top: 12px;
}

.customer-main {
  min-width: 0;
  background: transparent;
}

.customer-header {
  min-height: 88px;
  height: auto;
  line-height: 1;
  padding: 28px 36px 12px;
  background: transparent;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
}

.eyebrow {
  color: #5d7496;
}

.customer-header h2 {
  font-size: 30px;
  color: #10233f;
}

.header-actions {
  flex-wrap: wrap;
}

.customer-content {
  min-width: 0;
  max-width: 1440px;
  padding: 0 36px 36px;
}

@media (max-width: 992px) {
  .customer-header {
    padding-top: 20px;
    flex-direction: column;
    align-items: flex-start;
  }

  .customer-content {
    padding: 0 20px 24px;
  }
}
</style>
