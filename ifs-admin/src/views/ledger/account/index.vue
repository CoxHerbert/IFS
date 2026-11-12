<template>
  <div class="app-container">
    <a-alert message="这里只能查看用户主动开启云同步后上传的台账。仅保存在用户设备上的本地账单不会上传或显示在后台。" type="info" show-icon class="mb8" />
    <vxe-form :data="query" @submit="search" @reset="reset">
      <vxe-form-item title="登录账号" field="loginId"><template #default><vxe-input v-model="query.loginId" clearable placeholder="手机号或邮箱" /></template></vxe-form-item>
      <vxe-form-item><template #default><vxe-button type="submit" status="primary">查询</vxe-button><vxe-button type="reset">重置</vxe-button></template></vxe-form-item>
    </vxe-form>
    <vxe-table border stripe :loading="loading" :data="rows">
      <vxe-column field="uid" title="UID" min-width="180" />
      <vxe-column field="loginId" title="登录账号" min-width="220" />
      <vxe-column field="createTime" title="注册时间" min-width="180" />
      <vxe-column title="云同步" width="120"><template #default="{row}"><a-tag :color="row.cloudSyncEnabled ? 'green' : 'default'">{{ row.cloudSyncEnabled ? '已开启' : '已关闭' }}</a-tag></template></vxe-column>
      <vxe-column title="云端账单" width="120" align="right"><template #default="{row}">{{ row.recordCount }}</template></vxe-column>
      <vxe-column title="收入合计" width="150" align="right"><template #default="{row}">¥{{ money(row.incomeTotal) }}</template></vxe-column>
      <vxe-column title="支出合计" width="150" align="right"><template #default="{row}">¥{{ money(row.expenseTotal) }}</template></vxe-column>
      <vxe-column title="操作" width="150" fixed="right"><template #default="{row}"><vxe-button mode="text" status="primary" @click="openLedger(row)">个人台账</vxe-button></template></vxe-column>
    </vxe-table>
    <vxe-pager v-model:current-page="query.pageNum" v-model:page-size="query.pageSize" :total="total" @page-change="load" />
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { listLedgerAccounts } from '@/api/ledger/account'

const router = useRouter()
const loading = ref(false)
const rows = ref([])
const total = ref(0)
const query = reactive({ pageNum: 1, pageSize: 10, loginId: '' })
const money = value => Number(value || 0).toFixed(2)

function load() {
  loading.value = true
  listLedgerAccounts(query).then(response => {
    rows.value = response.data?.rows || []
    total.value = response.data?.total || 0
  }).finally(() => { loading.value = false })
}
function search() { query.pageNum = 1; load() }
function reset() { Object.assign(query, { pageNum: 1, pageSize: 10, loginId: '' }); load() }
function openLedger(row) {
  router.push({ path: `/ledger/account-detail/${row.accountId}`, query: { loginId: row.loginId } })
}

load()
</script>

<style scoped>
.app-container :deep(.vxe-pager){margin-top:12px}
</style>
