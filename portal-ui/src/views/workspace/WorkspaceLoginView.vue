<template>
  <main class="login-page">
    <section class="login-panel">
      <div class="login-copy">
        <p>客户端入口</p>
        <h2>登录后进入专属工作台</h2>
        <span>支持同一客户配置多个登录账号，分别给业务、财务和操作人员使用。</span>
      </div>

      <a-form layout="vertical" :model="form" class="login-form" @finish="handleLogin">
        <a-form-item label="账号" name="username" :rules="[{ required: true, message: '请输入账号' }]">
          <a-input v-model:value="form.username" size="large" placeholder="请输入账号" />
        </a-form-item>
        <a-form-item label="密码" name="password" :rules="[{ required: true, message: '请输入密码' }]">
          <a-input-password v-model:value="form.password" size="large" placeholder="请输入密码" />
        </a-form-item>
        <a-alert v-if="errorMessage" :message="errorMessage" type="error" show-icon />
        <a-button type="primary" html-type="submit" size="large" block :loading="loading">登录</a-button>
      </a-form>
    </section>
  </main>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { setWorkspaceToken, workspaceLogin } from '@/api/workspace/auth'

const router = useRouter()
const loading = ref(false)
const errorMessage = ref('')
const form = reactive({
  username: '',
  password: '',
})

async function handleLogin() {
  loading.value = true
  errorMessage.value = ''
  try {
    const response = await workspaceLogin(form.username, form.password)
    if (response.code !== 200 || !response.data?.token) {
      errorMessage.value = response.msg || '登录失败'
      return
    }
    setWorkspaceToken(response.data.token)
    router.push('/customer')
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '登录失败'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: 64px 24px;
  background:
    radial-gradient(circle at top left, rgba(14, 165, 233, 0.08), transparent 24%),
    linear-gradient(180deg, #f5f8fc 0%, #eef3f8 100%);
}

.login-panel {
  width: min(920px, 100%);
  display: grid;
  grid-template-columns: 1fr 380px;
  gap: 36px;
  padding: 36px;
  background:
    radial-gradient(circle at top right, rgba(56, 189, 248, 0.08), transparent 24%),
    rgba(255, 255, 255, 0.94);
  border: 1px solid rgba(16, 35, 63, 0.08);
  border-radius: 24px;
  box-shadow: 0 24px 60px rgba(37, 76, 124, 0.12);
}

.login-copy {
  align-self: center;
}

.login-copy p {
  margin: 0 0 14px;
  color: #0284c7;
  font-weight: 700;
  letter-spacing: 0.08em;
}

.login-copy h2 {
  margin: 0;
  font-size: 34px;
  line-height: 1.2;
  color: #0f172a;
}

.login-copy span {
  display: block;
  margin-top: 18px;
  color: #5f6f85;
  line-height: 1.8;
}

.login-form {
  padding: 8px 0;
  border-radius: 18px;
}

.login-form :deep(.ant-input-affix-wrapper),
.login-form :deep(.ant-input) {
  border-radius: 12px;
}

.login-form :deep(.ant-alert) {
  margin-bottom: 16px;
}

.login-form :deep(.ant-btn-primary) {
  border-radius: 12px;
}

@media (max-width: 760px) {
  .login-panel {
    grid-template-columns: 1fr;
    padding: 24px;
  }

  .login-copy h2 {
    font-size: 26px;
  }
}
</style>
