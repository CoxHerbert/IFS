<template>
  <main class="share-page">
    <a-spin :spinning="loading">
      <template v-if="detail?.plan">
        <section class="mobile-hero">
          <div class="hero-copy">
            <span>出货追踪</span>
            <h1>{{ detail.plan.shipmentNo }}</h1>
            <p>{{ detail.plan.pol || '-' }} -> {{ detail.plan.pod || '-' }}</p>
          </div>
          <a-tag class="status-tag" :color="statusColor(detail.plan.status)">{{ currentStatus }}</a-tag>
        </section>

        <section class="section-block">
          <div class="section-title">
            <h2>出货概览</h2>
            <span>{{ detail.plan.customerName || '客户' }}</span>
          </div>
          <div class="summary-list">
            <div class="summary-item strong">
              <span>当前状态</span>
              <strong>{{ currentStatus }}</strong>
            </div>
            <div class="summary-item">
              <span>客户参考号</span>
              <strong>{{ detail.plan.orderNo || '-' }}</strong>
            </div>
            <div class="summary-item">
              <span>出货单号</span>
              <strong>{{ detail.order?.orderNo || '未生成' }}</strong>
            </div>
            <div class="summary-item">
              <span>货量</span>
              <strong>{{ detail.plan.totalCartons }} 箱 / {{ detail.plan.totalVolume }} CBM</strong>
            </div>
            <div class="summary-item">
              <span>总重量</span>
              <strong>{{ detail.plan.totalWeight }} KG</strong>
            </div>
            <div class="summary-item">
              <span>计划状态</span>
              <strong>{{ statusLabel(detail.plan.status) }}</strong>
            </div>
            <div class="summary-item">
              <span>付款状态</span>
              <strong>{{ paymentStatusLabel(detail.plan.paymentStatus) }}</strong>
            </div>
            <div class="summary-item">
              <span>付款金额</span>
              <strong>{{ money(detail.plan.paymentAmount) }}</strong>
            </div>
          </div>
        </section>

        <section class="content-grid">
          <article class="section-block">
            <div class="section-title">
              <h2>时间节点</h2>
              <span>ETD / ETA</span>
            </div>
            <dl class="field-list">
              <div>
                <dt>计划开船</dt>
                <dd>{{ detail.plan.plannedEtd || '-' }}</dd>
              </div>
              <div>
                <dt>实际开船</dt>
                <dd>{{ detail.plan.actualEtd || '-' }}</dd>
              </div>
              <div>
                <dt>计划到港</dt>
                <dd>{{ detail.plan.plannedEta || '-' }}</dd>
              </div>
              <div>
                <dt>实际到港</dt>
                <dd>{{ detail.plan.actualEta || '-' }}</dd>
              </div>
              <div>
                <dt>起运港</dt>
                <dd>{{ detail.plan.pol || '-' }}</dd>
              </div>
              <div>
                <dt>目的港</dt>
                <dd>{{ detail.plan.pod || '-' }}</dd>
              </div>
            </dl>
          </article>

          <article class="section-block">
            <div class="section-title">
              <h2>柜型方案</h2>
              <span>{{ detail.containers.length }} 项</span>
            </div>
            <div v-if="detail.containers.length" class="stack-list">
              <div v-for="item in detail.containers" :key="item.containerType" class="info-row">
                <strong>{{ item.containerType }} x {{ item.quantity }}</strong>
                <span>装载率 {{ item.loadRate }}%</span>
                <small>{{ item.remark || '暂无备注' }}</small>
              </div>
            </div>
            <a-empty v-else description="暂无柜型信息" />
          </article>
        </section>

        <section class="section-block">
          <div class="section-title">
            <h2>状态进度</h2>
            <span>{{ activeStepCount }}/{{ detail.statusFlow.length }}</span>
          </div>
          <div class="status-flow">
            <article v-for="step in detail.statusFlow" :key="step.value"
              :class="['status-chip', { active: step.active }]">
              <strong>{{ step.label }}</strong>
              <span>{{ step.value }}</span>
            </article>
          </div>
        </section>

        <section class="section-block">
          <div class="section-title">
            <h2>货物明细</h2>
            <span>{{ detail.cargoList.length }} 项</span>
          </div>
          <div class="cargo-list">
            <article v-for="item in detail.cargoList" :key="item.cargoName + item.sku" class="cargo-card">
              <div>
                <strong>{{ item.cargoName }}</strong>
                <span>{{ item.sku || '无 SKU' }}</span>
              </div>
              <dl>
                <div>
                  <dt>箱数</dt>
                  <dd>{{ item.cartons }}</dd>
                </div>
                <div>
                  <dt>体积</dt>
                  <dd>{{ item.volumeCbm }} CBM</dd>
                </div>
                <div>
                  <dt>重量</dt>
                  <dd>{{ item.weightKg }} KG</dd>
                </div>
              </dl>
            </article>
          </div>
        </section>

        <section class="section-block">
          <div class="section-title">
            <h2>备注</h2>
            <span>Remark</span>
          </div>
          <p class="remark-copy">{{ detail.plan.remark || '暂无备注' }}</p>
        </section>
      </template>

      <a-result v-else-if="!loading" status="warning" title="未找到这票货的分享信息">
        <template #subTitle>请确认链接是否完整，或联系您的货代专员。</template>
      </a-result>
    </a-spin>
  </main>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { getPortalShipmentShare, type ShipmentDetail } from '@/api/portal/shipment'

const route = useRoute()
const loading = ref(true)
const detail = ref<ShipmentDetail>()

const currentStatus = computed(() => {
  const activeStatuses = detail.value?.statusFlow?.filter((item) => item.active) || []
  const status = activeStatuses[activeStatuses.length - 1]
  return status?.label || '待更新'
})

const activeStepCount = computed(() => detail.value?.statusFlow?.filter((item) => item.active).length || 0)

function statusLabel(status: string) {
  const item = detail.value?.statusFlow?.find((step) => step.value === status)
  return item?.label || status || '-'
}

function statusColor(status: string) {
  if (status === '900') return 'red'
  if (Number(status) >= 170) return 'green'
  if (Number(status) >= 100) return 'cyan'
  if (Number(status) >= 60) return 'blue'
  return 'gold'
}

function paymentStatusLabel(value?: string) {
  if (value === 'PAID') return '已付款'
  if (value === 'PARTIAL') return '部分付款'
  return '未付款'
}

function money(value?: number) {
  return `¥${Number(value || 0).toFixed(2)}`
}

onMounted(async () => {
  try {
    const token = String(route.params.token || '')
    const response = await getPortalShipmentShare(token)
    if (response.code === 200 && response.data) {
      detail.value = response.data
    }
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.share-page {
  width: min(1040px, calc(100% - 24px));
  margin: 0 auto;
  padding: 18px 0 42px;
}

.mobile-hero,
.section-block {
  border-radius: 8px;
  border: 1px solid rgba(15, 23, 42, 0.08);
  background: #fff;
  box-shadow: 0 10px 28px rgba(15, 23, 42, 0.06);
}

.mobile-hero {
  min-height: 168px;
  padding: 20px;
  display: grid;
  align-content: end;
  gap: 16px;
  background:
    linear-gradient(rgba(15, 23, 42, 0.68), rgba(15, 23, 42, 0.42));
  color: #fff;
}

.hero-copy span {
  color: rgba(255, 255, 255, 0.78);
  font-size: 12px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.hero-copy h1 {
  margin: 8px 0;
  font-size: clamp(26px, 8vw, 40px);
  overflow-wrap: anywhere;
}

.hero-copy p {
  margin: 0;
  font-size: 16px;
}

.status-tag {
  width: fit-content;
  margin: 0;
  font-size: 13px;
  padding: 4px 10px;
}

.section-block {
  margin-top: 14px;
  padding: 16px;
}

.section-title {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 14px;
}

.section-title h2 {
  margin: 0;
  color: #0f172a;
  font-size: 18px;
}

.section-title span,
.summary-item span,
.field-list dt,
.info-row span,
.info-row small,
.cargo-card span,
.cargo-card dt,
.remark-copy,
.status-chip span {
  color: #64748b;
}

.summary-list,
.content-grid,
.cargo-list {
  display: grid;
  gap: 12px;
}

.summary-list {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.summary-item,
.info-row,
.cargo-card {
  border-radius: 8px;
  background: #f8fafc;
  border: 1px solid rgba(148, 163, 184, 0.14);
  padding: 14px;
}

.summary-item.strong {
  background: #111827;
}

.summary-item.strong span,
.summary-item.strong strong {
  color: #fff;
}

.summary-item span,
.summary-item strong,
.info-row strong,
.info-row span,
.info-row small,
.cargo-card strong,
.cargo-card span,
.status-chip strong,
.status-chip span {
  display: block;
}

.summary-item strong {
  margin-top: 8px;
  color: #0f172a;
  font-size: 16px;
  overflow-wrap: anywhere;
}

.field-list {
  display: grid;
  gap: 10px;
  margin: 0;
}

.field-list div {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid #edf1f7;
}

.field-list div:last-child {
  border-bottom: 0;
}

.field-list dt,
.field-list dd,
.cargo-card dl,
.cargo-card dt,
.cargo-card dd {
  margin: 0;
}

.field-list dd {
  color: #0f172a;
  font-weight: 700;
  text-align: right;
}

.info-row+.info-row {
  margin-top: 10px;
}

.info-row span,
.info-row small,
.cargo-card span {
  margin-top: 6px;
}

.status-flow {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.status-chip {
  min-width: 0;
  padding: 9px 12px;
  border-radius: 999px;
  border: 1px solid rgba(148, 163, 184, 0.18);
  background: #f8fafc;
}

.status-chip strong {
  color: #0f172a;
  font-size: 13px;
  line-height: 1.2;
}

.status-chip span {
  margin-top: 3px;
  font-size: 11px;
  line-height: 1;
}

.status-chip.active {
  border-color: rgba(37, 99, 235, 0.22);
  background: #eff6ff;
}

.status-chip.active strong,
.status-chip.active span {
  color: #1d4ed8;
}

.cargo-card {
  display: grid;
  gap: 12px;
}

.cargo-card dl {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8px;
}

.cargo-card dd {
  margin-top: 4px;
  color: #0f172a;
  font-weight: 700;
}

.remark-copy {
  margin: 0;
  line-height: 1.8;
  overflow-wrap: anywhere;
}

@media (min-width: 760px) {
  .share-page {
    padding-top: 30px;
  }

  .mobile-hero {
    min-height: 220px;
    grid-template-columns: minmax(0, 1fr) auto;
    align-items: end;
  }

  .content-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .summary-list {
    grid-template-columns: repeat(4, minmax(0, 1fr));
  }

  .cargo-list {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 420px) {

  .summary-list,
  .cargo-card dl {
    grid-template-columns: 1fr;
  }

  .field-list div {
    display: grid;
  }

  .field-list dd {
    text-align: left;
  }
}
</style>
