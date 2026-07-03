<template>
  <main class="shipment-page">
    <section class="lookup-panel">
      <div>
        <h1>出货查询</h1>
        <p>输入分享链接或分享令牌后即可查询当前出货计划。</p>
      </div>

      <a-form layout="vertical" :model="form" @finish="handleLookup">
        <a-form-item label="分享令牌" name="token" :rules="[{ required: true, message: '请输入分享令牌' }]">
          <a-input v-model:value="form.token" size="large" placeholder="请输入分享链接中的 token" />
        </a-form-item>
        <a-space>
          <a-button type="primary" html-type="submit" :loading="loading">查询出货</a-button>
          <a-button @click="fillFromRoute">读取地址参数</a-button>
        </a-space>
      </a-form>
    </section>

    <a-alert v-if="errorMessage" :message="errorMessage" type="error" show-icon class="top-gap" />

    <template v-if="detail?.plan">
      <section class="summary-grid top-gap">
        <article class="summary-card strong">
          <span>计划编号</span>
          <strong>{{ detail.plan.shipmentNo }}</strong>
        </article>
        <article class="summary-card">
          <span>当前状态</span>
          <strong>{{ currentStatus }}</strong>
        </article>
        <article class="summary-card">
          <span>航线</span>
          <strong>{{ detail.plan.pol || '-' }} 到 {{ detail.plan.pod || '-' }}</strong>
        </article>
        <article class="summary-card">
          <span>出货单号</span>
          <strong>{{ detail.order?.orderNo || '未生成' }}</strong>
        </article>
      </section>

      <section class="panel-grid top-gap">
        <article class="panel-card">
          <h3>状态时间线</h3>
          <div class="timeline">
            <div v-for="step in detail.statusFlow" :key="step.value" :class="['timeline-item', { active: step.active }]">
              <strong>{{ step.value }}</strong>
              <span>{{ step.label }}</span>
            </div>
          </div>
        </article>

        <article class="panel-card">
          <h3>推荐货柜</h3>
          <div v-for="container in detail.containers" :key="container.containerType" class="container-item">
            <strong>{{ container.containerType }} x {{ container.quantity }}</strong>
            <span>装载率 {{ container.loadRate }}%</span>
            <small>{{ container.remark }}</small>
          </div>
        </article>
      </section>

      <section class="panel-card top-gap">
        <h3>货物明细</h3>
        <div class="cargo-grid">
          <article v-for="cargo in detail.cargoList" :key="cargo.cargoName + cargo.sku" class="cargo-item">
            <strong>{{ cargo.cargoName }}</strong>
            <span>{{ cargo.sku || '无 SKU' }}</span>
            <small>{{ cargo.cartons }} 箱 / {{ cargo.volumeCbm }} CBM / {{ cargo.weightKg }} KG</small>
          </article>
        </div>
      </section>
    </template>
  </main>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute } from 'vue-router'
import { getPortalShipmentShare, type ShipmentDetail } from '@/api/portal/shipment'

const route = useRoute()
const loading = ref(false)
const errorMessage = ref('')
const detail = ref<ShipmentDetail>()
const form = reactive({
  token: '',
})

const currentStatus = computed(() => {
  const activeStatuses = detail.value?.statusFlow?.filter((item) => item.active) || []
  return activeStatuses[activeStatuses.length - 1]?.label || '待更新'
})

function normalizeToken(raw: string) {
  const value = raw.trim()
  if (!value) {
    return ''
  }
  const match = value.match(/\/shipment\/share\/([^/?#]+)/)
  return match?.[1] || value
}

function fillFromRoute() {
  const token = typeof route.query.token === 'string' ? route.query.token : ''
  form.token = normalizeToken(token)
}

async function handleLookup() {
  loading.value = true
  errorMessage.value = ''
  detail.value = undefined
  try {
    const response = await getPortalShipmentShare(normalizeToken(form.token))
    if (response.code !== 200 || !response.data) {
      errorMessage.value = response.msg || '未找到出货信息'
      return
    }
    detail.value = response.data
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '查询失败'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fillFromRoute()
  if (form.token) {
    handleLookup()
  }
})
</script>

<style scoped>
.shipment-page {
  min-height: 100%;
}

.lookup-panel,
.summary-card,
.panel-card {
  border-radius: 20px;
  border: 1px solid rgba(148, 163, 184, 0.18);
  box-shadow: 0 20px 42px rgba(15, 23, 42, 0.06);
  background: #fff;
}

.lookup-panel {
  padding: 24px;
  display: grid;
  grid-template-columns: minmax(0, 1.2fr) minmax(320px, 1fr);
  gap: 22px;
}

.summary-card span {
  letter-spacing: 0.08em;
  font-size: 11px;
  color: #64748b;
}

.lookup-panel h1 {
  margin: 14px 0 0;
  color: #0f172a;
  font-size: 30px;
}

.lookup-panel p {
  margin: 12px 0 0;
  color: #64748b;
  line-height: 1.8;
  max-width: 42ch;
}

.top-gap {
  margin-top: 18px;
}

.summary-grid,
.panel-grid,
.cargo-grid {
  display: grid;
  gap: 18px;
}

.summary-grid {
  grid-template-columns: repeat(4, minmax(0, 1fr));
}

.summary-card {
  min-height: 126px;
  padding: 20px;
}

.summary-card.strong {
  background: #111827;
}

.summary-card.strong span,
.summary-card.strong strong {
  color: #fff;
}

.summary-card strong {
  display: block;
  margin-top: 10px;
  font-size: 20px;
  color: #0f172a;
  overflow-wrap: anywhere;
}

.panel-grid {
  grid-template-columns: 1.4fr 1fr;
}

.panel-card {
  padding: 24px;
}

.panel-card h3 {
  margin: 0 0 16px;
  color: #0f172a;
  font-size: 20px;
}

.timeline {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.timeline-item,
.container-item,
.cargo-item {
  border-radius: 16px;
  background: #f8fafc;
}

.timeline-item {
  min-height: 96px;
  padding: 16px;
  color: #64748b;
  border: 1px solid rgba(148, 163, 184, 0.12);
}

.timeline-item.active {
  background: #eff6ff;
  border: 1px solid rgba(59, 130, 246, 0.2);
  color: #1d4ed8;
}

.timeline-item strong,
.timeline-item span,
.container-item strong,
.container-item span,
.container-item small,
.cargo-item strong,
.cargo-item span,
.cargo-item small {
  display: block;
}

.timeline-item strong {
  margin-bottom: 8px;
  font-size: 20px;
}

.container-item,
.cargo-item {
  padding: 16px;
  border: 1px solid rgba(148, 163, 184, 0.12);
}

.container-item + .container-item {
  margin-top: 12px;
}

.container-item span,
.cargo-item span,
.cargo-item small {
  margin-top: 8px;
  color: #64748b;
}

.cargo-grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

@media (max-width: 1100px) {
  .lookup-panel,
  .panel-grid,
  .summary-grid,
  .cargo-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 640px) {
  .timeline {
    grid-template-columns: 1fr;
  }
}
</style>
