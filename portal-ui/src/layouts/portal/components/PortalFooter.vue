<template>
  <a-layout-footer class="footer">
    <div>
      <strong>白泽国际物流</strong>
      <p>提供报价、航线、出货计划与客户协同入口。</p>
    </div>
    <div class="footer-meta">
      <span>quote@seawaypro.com</span>
      <span>400-888-2026</span>
    </div>
  </a-layout-footer>
</template>

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
