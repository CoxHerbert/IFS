<template>
  <main class="shipment-detail-page">
    <section class="page-head">
      <div>
        <a-button type="link" class="back-link" @click="router.push('/customer/shipment')">返回出货查询</a-button>
        <h1>{{ detail?.plan.shipmentNo || '出货计划详情' }}</h1>
        <p>查看单票出货计划、状态进度、货物明细和费用付款信息。</p>
      </div>
      <a-button :loading="loading" @click="loadDetail">刷新</a-button>
    </section>

    <a-alert v-if="errorMessage" :message="errorMessage" type="error" show-icon class="top-gap" />

    <section v-if="detail?.plan" class="detail-tabs-card top-gap">
      <a-tabs v-model:activeKey="activeTab">
        <a-tab-pane key="overview" tab="概览">
          <section class="summary-grid">
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
              <strong>{{ detail.plan.pol || '-' }} → {{ detail.plan.pod || '-' }}</strong>
            </article>
            <article class="summary-card">
              <span>出货单号</span>
              <strong>{{ detail.order?.orderNo || '未生成' }}</strong>
            </article>
          </section>

          <section class="panel-card inner-gap">
            <h3>基础信息</h3>
            <dl class="field-list">
              <div><dt>客户订单号</dt><dd>{{ detail.plan.orderNo || '-' }}</dd></div>
              <div><dt>计划开船</dt><dd>{{ detail.plan.plannedEtd || '-' }}</dd></div>
              <div><dt>实际开船</dt><dd>{{ detail.plan.actualEtd || '-' }}</dd></div>
              <div><dt>计划到港</dt><dd>{{ detail.plan.plannedEta || '-' }}</dd></div>
              <div><dt>实际到港</dt><dd>{{ detail.plan.actualEta || '-' }}</dd></div>
              <div><dt>总箱数</dt><dd>{{ detail.plan.totalCartons }}</dd></div>
              <div><dt>总体积</dt><dd>{{ detail.plan.totalVolume }} CBM</dd></div>
              <div><dt>总重量</dt><dd>{{ detail.plan.totalWeight }} KG</dd></div>
              <div><dt>创建时间</dt><dd>{{ formatTime(detail.plan.createTime) }}</dd></div>
              <div><dt>更新时间</dt><dd>{{ formatTime(detail.plan.updateTime) }}</dd></div>
            </dl>
          </section>
        </a-tab-pane>

        <a-tab-pane key="status" tab="状态进度">
          <section class="panel-card">
            <h3>状态时间线</h3>
            <div class="timeline">
              <div v-for="step in detail.statusFlow" :key="step.value" :class="['timeline-item', { active: step.active }]">
                <strong>{{ step.value }}</strong>
                <span>{{ step.label }}</span>
              </div>
            </div>
          </section>
        </a-tab-pane>

        <a-tab-pane key="containers" tab="柜型方案">
          <section class="panel-card">
            <h3>推荐方案</h3>
            <div v-if="detail.containers.length">
              <div v-for="container in detail.containers" :key="container.containerType" class="container-item">
                <strong>{{ container.containerType }} x {{ container.quantity }}</strong>
                <span>装载率 {{ container.loadRate }}%</span>
                <small>{{ container.remark }}</small>
              </div>
            </div>
            <a-empty v-else description="暂无柜型建议" />
          </section>
        </a-tab-pane>

        <a-tab-pane key="cargo" tab="货物明细">
          <section class="panel-card">
            <h3>货物明细</h3>
            <div class="cargo-grid">
              <article v-for="cargo in detail.cargoList" :key="cargo.cargoName + cargo.sku" class="cargo-item">
                <strong>{{ cargo.cargoName }}</strong>
                <span>{{ cargo.sku || '无 SKU' }}</span>
                <small>{{ cargo.cartons }} 箱 / {{ cargo.volumeCbm }} CBM / {{ cargo.weightKg }} KG</small>
              </article>
            </div>
          </section>
        </a-tab-pane>

        <a-tab-pane key="finance" tab="费用付款">
          <section class="panel-card">
            <h3>费用与付款</h3>
            <div class="finance-grid">
              <div class="finance-item">
                <span>应付总额</span>
                <strong>待维护</strong>
              </div>
              <div class="finance-item">
                <span>已付金额</span>
                <strong>待维护</strong>
              </div>
              <div class="finance-item">
                <span>付款状态</span>
                <strong>{{ paymentStatusLabel(detail.plan.paymentStatus) }}</strong>
              </div>
              <div class="finance-item">
                <span>付款金额</span>
                <strong>{{ money(detail.plan.paymentAmount) }}</strong>
              </div>
              <div class="finance-item">
                <span>计划结案</span>
                <strong>{{ detail.plan.status === '170' ? (detail.plan.paymentStatus === 'PAID' ? '已结案' : '待付款确认') : '运输未完成' }}</strong>
              </div>
            </div>
            <p class="muted-copy">后续可接入基础运费、附加费、付款记录和凭证审核；客户付款完成后，计划才进入最终结束状态。</p>
          </section>
        </a-tab-pane>

        <a-tab-pane key="exception" tab="异常费用">
          <section class="panel-card">
            <h3>查验与异常费用</h3>
            <a-empty description="暂无查验或异常费用记录" />
          </section>

          <section class="panel-card inner-gap">
            <h3>备注</h3>
            <p class="muted-copy">{{ detail.plan.remark || '暂无备注' }}</p>
          </section>
        </a-tab-pane>
      </a-tabs>
    </section>

    <a-empty v-else-if="!loading" class="top-gap" description="未找到出货计划" />
  </main>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getWorkspaceToken } from '@/api/workspace/auth'
import { getWorkspaceShipmentDetail, type ShipmentDetail } from '@/api/portal/shipment'

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const errorMessage = ref('')
const detail = ref<ShipmentDetail>()
const activeTab = ref('overview')

const shipmentId = computed(() => String(route.params.shipmentId || ''))
const currentStatus = computed(() => {
  const activeStatuses = detail.value?.statusFlow?.filter((item) => item.active) || []
  return activeStatuses[activeStatuses.length - 1]?.label || '待更新'
})

function paymentStatusLabel(value?: string) {
  if (value === 'PAID') return '已付款'
  if (value === 'PARTIAL') return '部分付款'
  return '未付款'
}

function money(value?: number) {
  return `¥${Number(value || 0).toFixed(2)}`
}

async function loadDetail() {
  const token = getWorkspaceToken()
  if (!token) {
    errorMessage.value = '请先登录客户中心'
    return
  }
  if (!shipmentId.value) {
    errorMessage.value = '缺少出货计划编号'
    return
  }
  loading.value = true
  errorMessage.value = ''
  try {
    const response = await getWorkspaceShipmentDetail(shipmentId.value, token)
    if (response.code !== 200 || !response.data) {
      errorMessage.value = response.msg || '详情加载失败'
      return
    }
    detail.value = response.data
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '详情加载失败'
  } finally {
    loading.value = false
  }
}

function formatTime(value?: string) {
  if (!value) {
    return '-'
  }
  return String(value).replace('T', ' ').slice(0, 19)
}

onMounted(loadDetail)
</script>

<style scoped>
.shipment-detail-page {
  min-height: 100%;
}

.page-head,
.summary-card,
.detail-tabs-card,
.panel-card {
  border-radius: 8px;
  border: 1px solid rgba(148, 163, 184, 0.18);
  background: #fff;
}

.page-head {
  padding: 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.back-link {
  padding: 0;
  height: auto;
  margin-bottom: 10px;
}

.page-head h1 {
  margin: 0;
  color: #0f172a;
}

.page-head p,
.muted-copy {
  margin: 8px 0 0;
  color: #64748b;
  line-height: 1.8;
}

.top-gap {
  margin-top: 18px;
}

.inner-gap {
  margin-top: 18px;
}

.detail-tabs-card {
  padding: 18px 22px 22px;
}

.detail-tabs-card :deep(.ant-tabs-nav) {
  margin-bottom: 18px;
}

.summary-grid,
.cargo-grid,
.finance-grid {
  display: grid;
  gap: 18px;
}

.summary-grid {
  grid-template-columns: repeat(4, minmax(0, 1fr));
}

.summary-card {
  min-height: 118px;
  padding: 20px;
}

.summary-card.strong {
  background: #111827;
}

.summary-card.strong span,
.summary-card.strong strong {
  color: #fff;
}

.summary-card span,
.finance-item span {
  font-size: 12px;
  color: #64748b;
}

.summary-card strong,
.finance-item strong {
  display: block;
  margin-top: 10px;
  font-size: 20px;
  color: #0f172a;
  overflow-wrap: anywhere;
}

.panel-card {
  padding: 22px;
}

.panel-card h3 {
  margin: 0 0 16px;
  color: #0f172a;
}

.field-list {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px 18px;
  margin: 0;
}

.field-list div {
  min-height: 68px;
  border-radius: 8px;
  background: #f8fafc;
  border: 1px solid rgba(148, 163, 184, 0.12);
  padding: 14px;
}

.field-list dt,
.field-list dd {
  margin: 0;
}

.field-list dt {
  color: #64748b;
  font-size: 12px;
}

.field-list dd {
  margin-top: 8px;
  color: #0f172a;
  font-weight: 700;
  overflow-wrap: anywhere;
}

.timeline,
.finance-grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.timeline-item,
.container-item,
.cargo-item,
.finance-item {
  border-radius: 8px;
  background: #f8fafc;
  border: 1px solid rgba(148, 163, 184, 0.12);
  padding: 14px;
}

.timeline-item.active {
  background: #eff6ff;
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
  .page-head {
    align-items: stretch;
    flex-direction: column;
  }

  .summary-grid,
  .cargo-grid,
  .field-list {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 640px) {
  .timeline,
  .finance-grid {
    grid-template-columns: 1fr;
  }
}
</style>
