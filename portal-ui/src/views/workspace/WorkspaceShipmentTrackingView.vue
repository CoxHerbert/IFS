<template>
  <main class="shipment-page">
    <section class="page-head">
      <div>
        <h1>出货查询</h1>
        <p>查看当前客户账号对应的出货计划、货物明细和状态进度。</p>
      </div>
      <a-space>
        <a-input-search
          v-model:value="query.shipmentNo"
          placeholder="计划编号"
          enter-button="搜索"
          allow-clear
          @search="handleSearch"
        />
        <a-button @click="loadList">刷新</a-button>
      </a-space>
    </section>

    <a-alert v-if="errorMessage" :message="errorMessage" type="error" show-icon class="top-gap" />

    <a-table
      class="top-gap"
      row-key="shipmentId"
      :loading="loading"
      :columns="columns"
      :data-source="shipmentList"
      :pagination="pagination"
      @change="handleTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'route'">
          {{ record.pol || '-' }} → {{ record.pod || '-' }}
        </template>
        <template v-else-if="column.key === 'cargo'">
          {{ record.totalCartons }} 箱 / {{ record.totalVolume }} CBM / {{ record.totalWeight }} KG
        </template>
        <template v-else-if="column.key === 'status'">
          <a-tag color="blue">{{ statusLabel(record.status) }}</a-tag>
        </template>
        <template v-else-if="column.key === 'action'">
          <a-button type="link" @click="openDetail(record.shipmentId)">详情</a-button>
        </template>
      </template>
    </a-table>

    <section class="lookup-panel top-gap">
      <div>
        <h2>分享令牌查询</h2>
        <p>如果运营给了分享链接，也可以输入 token 查看单票出货。</p>
      </div>
      <a-input-group compact class="token-search">
        <a-input v-model:value="shareToken" placeholder="分享 token 或完整分享链接" />
        <a-button type="primary" :loading="shareLoading" @click="lookupShare">查询</a-button>
      </a-input-group>
    </section>

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
          <strong>{{ detail.plan.pol || '-' }} → {{ detail.plan.pod || '-' }}</strong>
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
          <h3>推荐方案</h3>
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
import { message as antMessage } from 'ant-design-vue'
import { getWorkspaceToken } from '@/api/workspace/auth'
import {
  getPortalShipmentShare,
  getWorkspaceShipmentDetail,
  listWorkspaceShipments,
  type ShipmentDetail,
  type ShipmentPlan,
} from '@/api/portal/shipment'

const loading = ref(false)
const shareLoading = ref(false)
const errorMessage = ref('')
const shipmentList = ref<ShipmentPlan[]>([])
const total = ref(0)
const detail = ref<ShipmentDetail>()
const shareToken = ref('')
const query = reactive({
  pageNum: 1,
  pageSize: 10,
  shipmentNo: '',
})

const columns = [
  { title: '计划编号', dataIndex: 'shipmentNo', key: 'shipmentNo' },
  { title: '客户订单号', dataIndex: 'orderNo', key: 'orderNo' },
  { title: '航线', key: 'route' },
  { title: '货量', key: 'cargo' },
  { title: '状态', key: 'status' },
  { title: '操作', key: 'action', width: 90 },
]

const pagination = computed(() => ({
  current: query.pageNum,
  pageSize: query.pageSize,
  total: total.value,
  showSizeChanger: true,
}))

const currentStatus = computed(() => {
  const activeStatuses = detail.value?.statusFlow?.filter((item) => item.active) || []
  return activeStatuses[activeStatuses.length - 1]?.label || '待更新'
})

function statusLabel(status: string) {
  const map: Record<string, string> = {
    '10': '计划已创建',
    '20': '出货计划已确认',
    '30': '等待客户发货',
    '40': '已提货/已送仓',
    '50': '仓库已收货',
    '60': '已入仓/码头进仓',
    '70': '订舱处理中',
    '80': '舱位已确认',
    '90': '报关资料已收齐',
    '100': '报关已放行',
    '110': '已装柜',
    '120': '已进港/码头放行',
    '130': '船舶已开船',
    '140': '目的港已到港',
    '150': '目的港清关中',
    '160': '目的港已清关',
    '170': '已派送/已签收',
    '900': '异常处理中',
  }
  return map[status] || status
}

async function loadList() {
  const token = getWorkspaceToken()
  if (!token) {
    errorMessage.value = '请先登录客户中心'
    return
  }
  loading.value = true
  errorMessage.value = ''
  try {
    const response = await listWorkspaceShipments(query, token)
    if (response.code !== 200 || !response.data) {
      errorMessage.value = response.msg || '出货计划加载失败'
      return
    }
    shipmentList.value = response.data.rows || []
    total.value = response.data.total || 0
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '出货计划加载失败'
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  query.pageNum = 1
  loadList()
}

function handleTableChange(pager: { current?: number; pageSize?: number }) {
  query.pageNum = pager.current || 1
  query.pageSize = pager.pageSize || 10
  loadList()
}

async function openDetail(shipmentId: string) {
  const token = getWorkspaceToken()
  if (!token) return
  try {
    const response = await getWorkspaceShipmentDetail(shipmentId, token)
    if (response.code !== 200 || !response.data) {
      antMessage.warning(response.msg || '未找到出货计划')
      return
    }
    detail.value = response.data
  } catch (error) {
    antMessage.error(error instanceof Error ? error.message : '详情加载失败')
  }
}

function normalizeToken(raw: string) {
  const value = raw.trim()
  const match = value.match(/\/shipment\/share\/([^/?#]+)/)
  return match?.[1] || value
}

async function lookupShare() {
  const token = normalizeToken(shareToken.value)
  if (!token) {
    antMessage.warning('请输入分享 token')
    return
  }
  shareLoading.value = true
  try {
    const response = await getPortalShipmentShare(token)
    if (response.code !== 200 || !response.data) {
      antMessage.warning(response.msg || '未找到出货信息')
      return
    }
    detail.value = response.data
  } catch (error) {
    antMessage.error(error instanceof Error ? error.message : '查询失败')
  } finally {
    shareLoading.value = false
  }
}

onMounted(loadList)
</script>

<style scoped>
.shipment-page {
  min-height: 100%;
}

.page-head,
.lookup-panel,
.summary-card,
.panel-card {
  border-radius: 8px;
  border: 1px solid rgba(148, 163, 184, 0.18);
  background: #fff;
}

.page-head,
.lookup-panel {
  padding: 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.page-head h1,
.lookup-panel h2 {
  margin: 0;
  color: #0f172a;
}

.page-head p,
.lookup-panel p {
  margin: 8px 0 0;
  color: #64748b;
}

.token-search {
  width: min(520px, 100%);
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

.summary-card span {
  font-size: 12px;
  color: #64748b;
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
  padding: 22px;
}

.panel-card h3 {
  margin: 0 0 16px;
  color: #0f172a;
}

.timeline {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.timeline-item,
.container-item,
.cargo-item {
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
  .page-head,
  .lookup-panel {
    align-items: stretch;
    flex-direction: column;
  }

  .panel-grid,
  .summary-grid,
  .cargo-grid {
    grid-template-columns: 1fr;
  }
}
</style>
