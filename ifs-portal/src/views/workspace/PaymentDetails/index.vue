<template>
  <div class="payment-page">
    <section class="page-head"><div><p>PAYMENTS</p><h1>付款明细</h1><span>查看当前客户各出货计划的应付、已付和未付金额。</span></div><a-select v-model:value="status" allow-clear placeholder="全部状态" :options="statusOptions" @change="search"/></section>
    <section class="summary"><article><span>应付合计</span><strong>¥{{ money(summary.payable) }}</strong></article><article><span>已付合计</span><strong>¥{{ money(summary.paid) }}</strong></article><article><span>未付合计</span><strong class="danger">¥{{ money(summary.unpaid) }}</strong></article></section>
    <section class="table-card"><a-table :loading="loading" :data-source="rows" :pagination="pagination" row-key="shipmentId" @change="pageChanged"><a-table-column title="出货计划" data-index="shipmentNo"/><a-table-column title="客户单号" data-index="orderNo"/><a-table-column title="应付金额"><template #default="{record}">¥{{money(record.payableAmount)}}</template></a-table-column><a-table-column title="已付金额"><template #default="{record}">¥{{money(record.paidAmount)}}</template></a-table-column><a-table-column title="未付金额"><template #default="{record}">¥{{money(record.unpaidAmount)}}</template></a-table-column><a-table-column title="状态"><template #default="{record}"><a-tag :color="paymentStatusColor(record.paymentStatus)">{{paymentStatusLabel(record.paymentStatus)}}</a-tag></template></a-table-column></a-table></section>
  </div>
</template>
<script setup lang="ts">
import { computed,onMounted,ref } from 'vue'
import { listWorkspacePayments,type WorkspacePaymentItem } from '@/api/workspace/payment'
const loading=ref(false),rows=ref<WorkspacePaymentItem[]>([]),total=ref(0),page=ref(1),status=ref<string>()
const labels={UNPAID:'未付款',PARTIAL:'部分付款',PAID:'已付款'},colors={UNPAID:'red',PARTIAL:'orange',PAID:'green'}
const statusOptions=Object.entries(labels).map(([value,label])=>({value,label})),pagination=computed(()=>({current:page.value,pageSize:10,total:total.value}))
const summary=computed(()=>rows.value.reduce((s,r)=>({payable:s.payable+Number(r.payableAmount||0),paid:s.paid+Number(r.paidAmount||0),unpaid:s.unpaid+Number(r.unpaidAmount||0)}),{payable:0,paid:0,unpaid:0}))
const money=(v:number)=>Number(v||0).toFixed(2)
const paymentStatusLabel=(value:WorkspacePaymentItem['paymentStatus'])=>labels[value]||value
const paymentStatusColor=(value:WorkspacePaymentItem['paymentStatus'])=>colors[value]||'default'
async function load(){loading.value=true;try{const r=await listWorkspacePayments({pageNum:page.value,pageSize:10,paymentStatus:status.value});rows.value=r.data?.rows||[];total.value=r.data?.total||0}finally{loading.value=false}}
function search(){page.value=1;load()} function pageChanged(p:{current?:number}){page.value=p.current||1;load()} onMounted(load)
</script>
<style scoped>.payment-page{display:grid;gap:18px}.page-head,.summary article,.table-card{border:1px solid var(--ws-border);border-radius:16px;background:var(--ws-surface);padding:20px}.page-head{display:flex;justify-content:space-between;align-items:end}.page-head p{margin:0;color:#0284c7;font-weight:700}.page-head h1{margin:6px 0}.page-head span,.summary span{color:var(--ws-text-muted)}.page-head :deep(.ant-select){width:180px}.summary{display:grid;grid-template-columns:repeat(3,1fr);gap:14px}.summary article{display:grid;gap:8px}.summary strong{font-size:24px}.danger{color:#dc2626}@media(max-width:700px){.page-head{align-items:start;gap:16px;flex-direction:column}.summary{grid-template-columns:1fr}.table-card{overflow:auto}}</style>
