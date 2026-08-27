<template>
  <main class="not-found-page">
    <section>
      <p class="status-code">404</p>
      <h1>页面不存在</h1>
      <p>你访问的页面不存在或已被移除。</p>
      <router-link to="/">返回首页</router-link>
    </section>
  </main>
</template>

<script setup lang="ts">
defineOptions({ name: 'NotFoundPage' })
</script>

<style scoped>
.not-found-page {
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: 24px;
  color: #0f172a;
  background: #f8fafc;
  text-align: center;
}

.status-code {
  margin: 0;
  color: #2563eb;
  font-size: clamp(72px, 16vw, 160px);
  font-weight: 800;
  line-height: 1;
}

h1 {
  margin: 20px 0 8px;
  font-size: 28px;
}

section > p:not(.status-code) {
  margin: 0 0 24px;
  color: #64748b;
}

a {
  display: inline-flex;
  padding: 10px 18px;
  border-radius: 8px;
  color: #fff;
  background: #2563eb;
  text-decoration: none;
}
</style>
