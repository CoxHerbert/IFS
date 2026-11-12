<script setup lang="ts">
import './styles/index.scss'
</script>

<style>
page { min-height: 100%; background: #f7f8fa; color: #202329; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; }
view, text, button { box-sizing: border-box; }
button::after { border: 0; }

/* Refine the shared H5 navigation header while keeping its safe-area layout. */
.uni-page-head {
  height: calc(52px + env(safe-area-inset-top)) !important;
  padding: calc(8px + env(safe-area-inset-top)) 14px 8px !important;
  border-bottom: 1px solid rgba(38, 82, 64, 0.08);
  box-shadow: 0 5px 18px rgba(30, 57, 45, 0.045);
  backdrop-filter: saturate(150%) blur(18px);
  -webkit-backdrop-filter: saturate(150%) blur(18px);
}
.uni-page-head ~ .uni-placeholder {
  height: calc(52px + env(safe-area-inset-top)) !important;
}
.uni-page-head .uni-page-head__title {
  color: #24372e;
  font-size: 17px;
  font-weight: 700;
  letter-spacing: 0.2px;
}
.uni-page-head .uni-page-head-hd,
.uni-page-head .uni-page-head-ft,
.uni-page-head .uni-btn-icon {
  color: #347662;
}
.uni-page-head .uni-page-head-hd .uni-btn-icon {
  font-size: 19px;
  font-weight: 600;
}
uni-page-head[uni-page-head-type='default'] ~ uni-page-wrapper {
  height: calc(100% - 52px - env(safe-area-inset-top)) !important;
}
</style>
