<template>
  <main class="share-page">
    <a-spin :spinning="loading">
      <template v-if="detail?.plan">
        <section class="shipment-hero">
          <div>
            <p>Shipment Tracking</p>
            <h1>{{ detail.plan.shipmentNo }}</h1>
            <span>{{ detail.plan.pol || '-' }} → {{ detail.plan.pod || '-' }}</span>
          </div>
          <div class="status-pill">{{ currentStatus }}</div>
        </section>

        <section class="timeline">
          <article v-for="step in detail.statusFlow" :key="step.value" :class="{ active: step.active }">
            <strong>{{ step.value }}</strong>
            <span>{{ step.label }}</span>
          </article>
        </section>

        <section class="metric-grid">
          <article>
            <span>客户参考号</span>
            <strong>{{ detail.plan.orderNo || '-' }}</strong>
          </article>
          <article>
            <span>出货单</span>
            <strong>{{ detail.order?.orderNo || '生成中' }}</strong>
          </article>
          <article>
            <span>货量</span>
            <strong>{{ detail.plan.totalCartons }}箱 / {{ detail.plan.totalVolume }}CBM</strong>
          </article>
          <article>
            <span>重量</span>
            <strong>{{ detail.plan.totalWeight }}KG</strong>
          </article>
        </section>

        <section class="panel-grid">
          <article class="panel">
            <h2>计划节点</h2>
            <dl>
              <div><dt>计划开船</dt><dd>{{ detail.plan.plannedEtd || '-' }}</dd></div>
              <div><dt>实际开船</dt><dd>{{ detail.plan.actualEtd || '-' }}</dd></div>
              <div><dt>计划到港</dt><dd>{{ detail.plan.plannedEta || '-' }}</dd></div>
              <div><dt>实际到港</dt><dd>{{ detail.plan.actualEta || '-' }}</dd></div>
            </dl>
          </article>

          <article class="panel">
            <h2>推荐货柜</h2>
            <div v-for="item in detail.containers" :key="item.containerType" class="container-row">
              <strong>{{ item.containerType }} × {{ item.quantity }}</strong>
              <span>装载率 {{ item.loadRate }}%</span>
              <small>{{ item.remark }}</small>
            </div>
          </article>
        </section>

        <section class="panel">
          <h2>货物明细</h2>
          <div class="cargo-list">
            <article v-for="item in detail.cargoList" :key="item.cargoName + item.sku">
              <strong>{{ item.cargoName }}</strong>
              <span>{{ item.sku || '无SKU' }}</span>
              <small>{{ item.cartons }}箱 / {{ item.volumeCbm }}CBM / {{ item.weightKg }}KG</small>
            </article>
          </div>
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
import { getShipmentShare, type ShipmentDetail } from '@/api/customer'

const route = useRoute()
const loading = ref(true)
const detail = ref<ShipmentDetail>()

const currentStatus = computed(() => {
  const activeStatuses = detail.value?.statusFlow?.filter(item => item.active) || []
  const status = activeStatuses[activeStatuses.length - 1]
  return status?.label || '待更新'
})

onMounted(async () => {
  try {
    const token = String(route.params.token || '')
    const response = await getShipmentShare(token)
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
  width: min(1120px, calc(100% - 32px));
  margin: 0 auto;
  padding: 34px 0 56px;
}

.shipment-hero {
  min-height: 220px;
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 24px;
  padding: 30px;
  border-radius: 22px;
  background:
    linear-gradient(rgba(8, 24, 45, 0.55), rgba(8, 24, 45, 0.35)),
    url('@/assets/hero.jpg') center/cover;
  color: #fff;
}

.shipment-hero p {
  margin: 0;
  font-size: 12px;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: rgba(255, 255, 255, 0.78);
}

.shipment-hero h1 {
  margin: 8px 0;
  font-size: 38px;
  letter-spacing: 0;
}

.shipment-hero span {
  font-size: 18px;
}

.status-pill {
  padding: 10px 16px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.16);
  backdrop-filter: blur(8px);
  font-weight: 700;
}

.timeline {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
  margin: 18px 0;
}

.timeline article,
.metric-grid article,
.panel {
  border: 1px solid rgba(16, 35, 63, 0.08);
  background: rgba(255, 255, 255, 0.92);
  box-shadow: 0 16px 34px rgba(21, 52, 93, 0.08);
}

.timeline article {
  min-height: 92px;
  padding: 16px;
  border-radius: 14px;
  color: #7a8799;
}

.timeline article.active {
  border-color: rgba(20, 126, 92, 0.35);
  background: #f3fbf7;
  color: #147e5c;
}

.timeline strong,
.timeline span {
  display: block;
}

.timeline strong {
  margin-bottom: 8px;
  font-size: 20px;
}

.metric-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 14px;
}

.metric-grid article {
  min-height: 112px;
  padding: 18px;
  border-radius: 14px;
}

.metric-grid span,
.panel dt,
.container-row span,
.cargo-list span,
.cargo-list small {
  color: #66748b;
}

.metric-grid strong {
  display: block;
  margin-top: 10px;
  font-size: 18px;
  overflow-wrap: anywhere;
}

.panel-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  margin-top: 16px;
}

.panel {
  padding: 22px;
  border-radius: 16px;
  margin-top: 16px;
}

.panel-grid .panel {
  margin-top: 0;
}

.panel h2 {
  margin: 0 0 16px;
  font-size: 20px;
}

dl {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
  margin: 0;
}

dt,
dd {
  margin: 0;
}

dd {
  margin-top: 6px;
  font-weight: 700;
}

.container-row,
.cargo-list article {
  padding: 14px 0;
  border-top: 1px solid #edf1f7;
}

.container-row:first-of-type,
.cargo-list article:first-child {
  border-top: 0;
}

.container-row strong,
.container-row span,
.container-row small,
.cargo-list strong,
.cargo-list span,
.cargo-list small {
  display: block;
}

.container-row small {
  margin-top: 6px;
  color: #8b96a8;
}

.cargo-list {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0 18px;
}

@media (max-width: 820px) {
  .shipment-hero,
  .panel-grid {
    display: block;
  }

  .status-pill {
    display: inline-block;
    margin-top: 18px;
  }

  .timeline,
  .metric-grid,
  .cargo-list,
  dl {
    grid-template-columns: 1fr;
  }
}
</style>
