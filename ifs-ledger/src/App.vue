<script setup lang="ts">
import { onLaunch, onShow } from '@dcloudio/uni-app'
import { applyTheme } from './utils/theme'
import './styles/index.scss'

onLaunch(() => applyTheme())
onShow(() => applyTheme())
</script>

<style>
page { min-height: 100%; background: var(--theme-page, #f6f5fa); color: #282638; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; --wot-blue-1: #f2f0fb; --wot-blue-2: #e8e4f8; --wot-blue-3: #d1c9f0; --wot-blue-4: #b0a4e4; --wot-blue-5: #8b7bd2; --wot-blue-6: #6d5dc0; --wot-blue-7: #5a4da7; --wot-blue-8: #483e86; --wot-blue-9: #393267; --wot-blue-10: #2b264e; }
view, text, button { box-sizing: border-box; }
button::after { border: 0; }

/* Refine the shared H5 navigation header while keeping its safe-area layout. */
.uni-page-head {
  height: calc(52px + env(safe-area-inset-top)) !important;
  padding: calc(8px + env(safe-area-inset-top)) 14px 8px !important;
  border-bottom: 1px solid rgba(82, 70, 142, 0.09);
  box-shadow: 0 5px 18px rgba(45, 39, 78, 0.045);
  backdrop-filter: saturate(150%) blur(18px);
  -webkit-backdrop-filter: saturate(150%) blur(18px);
}
.uni-page-head ~ .uni-placeholder {
  height: calc(52px + env(safe-area-inset-top)) !important;
}
.uni-page-head .uni-page-head__title {
  color: #302c45;
  font-size:19px;
  font-weight: 700;
  letter-spacing: 0.2px;
}
.uni-page-head .uni-page-head-hd,
.uni-page-head .uni-page-head-ft,
.uni-page-head .uni-btn-icon {
  color: var(--theme-primary, #6556b2);
}
.uni-page-head .uni-page-head-hd .uni-btn-icon {
  font-size:21.3px;
  font-weight: 600;
}
uni-page-head[uni-page-head-type='default'] ~ uni-page-wrapper {
  height: calc(100% - 52px - env(safe-area-inset-top)) !important;
}
</style>
