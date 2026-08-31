<template>
  <a-layout-footer class="footer">
    <div>
      <strong>{{ t('brandName') }}</strong>
      <p>{{ t('footerText') }}</p>
    </div>
    <div class="footer-meta">
      <router-link :to="localePath('/routes/china-to-usa')">{{ t('routeLink') }}</router-link>
      <router-link :to="localePath('/contact')">{{ t('getQuote') }}</router-link>
      <span>quote@seawaypro.com</span>
      <span>400-888-2026</span>
    </div>
  </a-layout-footer>
</template>

<script setup lang="ts">
import { usePortalI18n } from '@/i18n'
const { t, localePath } = usePortalI18n()
</script>

<style scoped>
.footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  padding: 28px 40px;
  background: rgba(255, 255, 255, 0.88);
  border-top: 1px solid rgba(16, 35, 63, 0.08);
}

.footer p {
  margin: 6px 0 0;
  color: #66748b;
}

.footer-meta {
  display: flex;
  gap: 18px;
  color: #55657d;
}

@media (max-width: 960px) {
  .footer {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
