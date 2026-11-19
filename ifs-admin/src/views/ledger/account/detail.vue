<template>
  <div class="app-container">
    <div class="page-heading">
      <div>
        <h2>个人台账</h2>
        <span>{{ account.loginId || route.query.loginId }} · UID {{ account.uid }} · {{ account.cloudSyncEnabled ? '已开启云同步' : '云同步已关闭' }}</span>
      </div>
      <div class="heading-actions"><a-button @click="viewMode = viewMode === 'records' ? 'batches' : 'records'">{{ viewMode === 'records' ? '导入批次' : '个人台账' }}</a-button><a-button @click="router.back()">返回用户列表</a-button></div>
    </div>
    <a-alert v-if="viewMode === 'records'" message="以下是该用户上传至云端的账单。用户设备上的本地数据不会显示在这里。删除记录会立即从云端移除。" type="info" show-icon class="mb8" />
    <a-alert v-else message="此处仅显示该用户上传至云端的导入批次。删除批次会同时删除批次内的云端账单。" type="info" show-icon class="mb8" />
    <vxe-table v-if="viewMode === 'records'" border stripe :loading="loading" :data="rows">
      <vxe-column field="date" title="账单时间" min-width="180" />
      <vxe-column field="merchant" title="商户 / 说明" min-width="200" />
      <vxe-column field="category" title="分类" width="130" />
      <vxe-column title="来源" width="120"><template #default="{row}">{{ sourceLabel(row.source) }}</template></vxe-column>
      <vxe-column field="importBatchId" title="导入批次" min-width="210"><template #default="{row}"><span v-if="row.importBatchId" :title="row.importFileName">{{ row.importBatchId }}</span><span v-else>—</span></template></vxe-column>
      <vxe-column title="收支类型" width="120"><template #default="{row}"><a-tag :color="row.type === 'income' ? 'green' : 'orange'">{{ row.type === 'income' ? '收入' : '支出' }}</a-tag></template></vxe-column>
      <vxe-column title="金额" width="140" align="right"><template #default="{row}">{{ row.type === 'income' ? '+' : '−' }}¥{{ money(row.amount) }}</template></vxe-column>
      <vxe-column title="操作" width="110" fixed="right"><template #default="{row}"><vxe-button mode="text" status="danger" v-hasPermi="['ledger:account:removeRecord']" @click="removeRecord(row)">删除</vxe-button></template></vxe-column>
    </vxe-table>
    <vxe-table v-else border stripe :loading="loading" :data="batchRows">
      <vxe-column field="batchId" title="批次编号" min-width="220" />
      <vxe-column title="来源" width="110"><template #default="{row}">{{ sourceLabel(row.source) }}</template></vxe-column>
      <vxe-column field="fileName" title="文件名" min-width="220" show-overflow />
      <vxe-column field="importedAt" title="导入时间" min-width="180" />
      <vxe-column field="importedCount" title="导入数量" width="110" align="right" />
      <vxe-column field="duplicateCount" title="重复数量" width="110" align="right" />
      <vxe-column field="recordCount" title="云端账单" width="110" align="right" />
      <vxe-column title="操作" width="110" fixed="right"><template #default="{row}"><vxe-button mode="text" status="danger" v-hasPermi="['ledger:account:removeRecord']" @click="removeBatch(row)">删除</vxe-button></template></vxe-column>
    </vxe-table>
    <vxe-pager v-model:current-page="query.pageNum" v-model:page-size="query.pageSize" :total="total" @page-change="load" />
  </div>
</template>

<script setup>
import { reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { deleteLedgerAccountImportBatch, deleteLedgerAccountRecord, getLedgerAccountImportBatches, getLedgerAccountRecords } from '@/api/ledger/account'

const { proxy } = getCurrentInstance()
const route = useRoute()
const router = useRouter()
const loading = ref(false)
const account = ref({})
const rows = ref([])
const batchRows = ref([])
const total = ref(0)
const viewMode = ref('records')
const query = reactive({ pageNum: 1, pageSize: 10 })
const money = value => Number(value || 0).toFixed(2)
const sourceLabel = source => ({ wechat: '微信', alipay: '支付宝', manual: '手动记账' }[source] || source)

function load() {
  loading.value = true
  const loadRequest = viewMode.value === 'records'
    ? getLedgerAccountRecords(route.params.accountId, query)
    : getLedgerAccountImportBatches(route.params.accountId, query)
  loadRequest.then(response => {
    account.value = response.data?.account || {}
    if (viewMode.value === 'records') rows.value = response.data?.rows || []
    else batchRows.value = response.data?.rows || []
    total.value = response.data?.total || 0
  }).finally(() => { loading.value = false })
}
watch(viewMode, () => { query.pageNum = 1; load() })
function removeRecord(row) {
  proxy.$modal.confirm(`确认删除“${row.merchant}”的账单记录吗？此操作不可撤销。`)
    .then(() => deleteLedgerAccountRecord(route.params.accountId, row.id))
    .then(() => { proxy.$modal.msgSuccess('账单记录已删除'); load() })
    .catch(() => {})
}
function removeBatch(row) {
  proxy.$modal.confirm(`确认删除导入批次“${row.batchId}”吗？批次内的云端账单也会一并删除。`)
    .then(() => deleteLedgerAccountImportBatch(route.params.accountId, row.batchId))
    .then(() => { proxy.$modal.msgSuccess('导入批次已删除'); load() })
    .catch(() => {})
}

load()
</script>

<style scoped>
.page-heading{display:flex;align-items:center;justify-content:space-between;margin-bottom:16px}.page-heading h2{margin:0 0 6px;font-size:20px}.page-heading span{color:#8c8c8c}.heading-actions{display:flex;gap:10px}.app-container :deep(.vxe-pager){margin-top:12px}
</style>
