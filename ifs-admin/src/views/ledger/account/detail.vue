<template>
  <div class="app-container">
    <div class="page-heading">
      <div>
        <h2>个人台账</h2>
        <span>{{ account.loginId || route.query.loginId }} · UID {{ account.uid }} · {{ account.cloudSyncEnabled ? '已开启云同步' : '云同步已关闭' }}</span>
      </div>
      <a-button @click="router.back()">返回用户列表</a-button>
    </div>
    <a-alert message="以下是该用户上传至云端的账单。用户设备上的本地数据不会显示在这里。删除记录会立即从云端移除。" type="info" show-icon class="mb8" />
    <vxe-table border stripe :loading="loading" :data="rows">
      <vxe-column field="date" title="账单时间" min-width="180" />
      <vxe-column field="merchant" title="商户 / 说明" min-width="200" />
      <vxe-column field="category" title="分类" width="130" />
      <vxe-column title="来源" width="120"><template #default="{row}">{{ sourceLabel(row.source) }}</template></vxe-column>
      <vxe-column field="importBatchId" title="导入批次" min-width="210"><template #default="{row}"><span v-if="row.importBatchId" :title="row.importFileName">{{ row.importBatchId }}</span><span v-else>—</span></template></vxe-column>
      <vxe-column title="收支类型" width="120"><template #default="{row}"><a-tag :color="row.type === 'income' ? 'green' : 'orange'">{{ row.type === 'income' ? '收入' : '支出' }}</a-tag></template></vxe-column>
      <vxe-column title="金额" width="140" align="right"><template #default="{row}">{{ row.type === 'income' ? '+' : '−' }}¥{{ money(row.amount) }}</template></vxe-column>
      <vxe-column title="操作" width="110" fixed="right"><template #default="{row}"><vxe-button mode="text" status="danger" v-hasPermi="['ledger:account:removeRecord']" @click="removeRecord(row)">删除</vxe-button></template></vxe-column>
    </vxe-table>
    <vxe-pager v-model:current-page="query.pageNum" v-model:page-size="query.pageSize" :total="total" @page-change="load" />
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { deleteLedgerAccountRecord, getLedgerAccountRecords } from '@/api/ledger/account'

const { proxy } = getCurrentInstance()
const route = useRoute()
const router = useRouter()
const loading = ref(false)
const account = ref({})
const rows = ref([])
const total = ref(0)
const query = reactive({ pageNum: 1, pageSize: 10 })
const money = value => Number(value || 0).toFixed(2)
const sourceLabel = source => ({ wechat: '微信', alipay: '支付宝', manual: '手动记账' }[source] || source)

function load() {
  loading.value = true
  getLedgerAccountRecords(route.params.accountId, query).then(response => {
    account.value = response.data?.account || {}
    rows.value = response.data?.rows || []
    total.value = response.data?.total || 0
  }).finally(() => { loading.value = false })
}
function removeRecord(row) {
  proxy.$modal.confirm(`确认删除“${row.merchant}”的账单记录吗？此操作不可撤销。`)
    .then(() => deleteLedgerAccountRecord(route.params.accountId, row.id))
    .then(() => { proxy.$modal.msgSuccess('账单记录已删除'); load() })
    .catch(() => {})
}

load()
</script>

<style scoped>
.page-heading{display:flex;align-items:center;justify-content:space-between;margin-bottom:16px}.page-heading h2{margin:0 0 6px;font-size:20px}.page-heading span{color:#8c8c8c}.app-container :deep(.vxe-pager){margin-top:12px}
</style>
