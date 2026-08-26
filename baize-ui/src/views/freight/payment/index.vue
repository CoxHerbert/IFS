<template>
  <div class="app-container">
    <vxe-form :data="query" @submit="search" @reset="reset">
      <vxe-form-item title="出货计划" field="shipmentNo"><template #default><vxe-input v-model="query.shipmentNo" clearable /></template></vxe-form-item>
      <vxe-form-item title="客户" field="customerId"><template #default><vxe-select v-model="query.customerId" clearable filterable><vxe-option v-for="c in customers" :key="c.customerId" :value="c.customerId" :label="c.customerName" /></vxe-select></template></vxe-form-item>
      <vxe-form-item title="是否付款" field="paymentStatus"><template #default><vxe-select v-model="query.paymentStatus" clearable><vxe-option value="UNPAID" label="未付款"/><vxe-option value="PARTIAL" label="部分付款"/><vxe-option value="PAID" label="已付款"/></vxe-select></template></vxe-form-item>
      <vxe-form-item><template #default><vxe-button type="submit" status="primary">查询</vxe-button><vxe-button type="reset">重置</vxe-button></template></vxe-form-item>
    </vxe-form>
    <vxe-table border stripe :loading="loading" :data="rows">
      <vxe-column field="shipmentNo" title="出货计划" min-width="170"/>
      <vxe-column field="orderNo" title="客户单号" min-width="150"/>
      <vxe-column field="customerName" title="客户" min-width="150"/>
      <vxe-column title="应付金额" width="130" align="right"><template #default="{row}">{{ money(row.payableAmount) }}</template></vxe-column>
      <vxe-column title="已付金额" width="130" align="right"><template #default="{row}">{{ money(row.paidAmount) }}</template></vxe-column>
      <vxe-column title="未付金额" width="130" align="right"><template #default="{row}">{{ money(row.unpaidAmount) }}</template></vxe-column>
      <vxe-column title="付款状态" width="110"><template #default="{row}"><a-tag :color="statusColor[row.paymentStatus]">{{ statusLabel[row.paymentStatus] }}</a-tag></template></vxe-column>
      <vxe-column title="操作" width="130"><template #default="{row}"><vxe-button mode="text" v-hasPermi="['freight:payment:edit']" @click="edit(row)">费用明细</vxe-button></template></vxe-column>
    </vxe-table>
    <vxe-pager v-model:current-page="query.pageNum" v-model:page-size="query.pageSize" :total="total" @page-change="load"/>
    <vxe-modal v-model="open" title="应付费用明细" width="760" show-footer>
      <div class="charge-toolbar"><b>系统自动加总：¥{{ money(chargeTotal) }}</b><vxe-button size="mini" @click="addCharge">增加费用</vxe-button></div>
      <vxe-table border :data="form.items">
        <vxe-column title="费用名称"><template #default="{row}"><vxe-input v-model="row.feeName" placeholder="如海运费、报关费"/></template></vxe-column>
        <vxe-column title="金额" width="180"><template #default="{row}"><vxe-number-input v-model="row.amount" type="float" :min="0" :digits="2"/></template></vxe-column>
        <vxe-column title="备注"><template #default="{row}"><vxe-input v-model="row.remark"/></template></vxe-column>
        <vxe-column title="操作" width="80"><template #default="{rowIndex}"><vxe-button mode="text" status="error" @click="form.items.splice(rowIndex,1)">删除</vxe-button></template></vxe-column>
      </vxe-table>
      <template #footer><vxe-button @click="open=false">取消</vxe-button><vxe-button status="primary" @click="save">保存</vxe-button></template>
    </vxe-modal>
  </div>
</template>
<script setup>
import { computed,reactive,ref } from 'vue'
import { simpleCustomerOptions } from '@/api/simple/simple'
import { getShipmentCharges,listPaymentLedger,saveShipmentCharges } from '@/api/freight/payment'
const {proxy}=getCurrentInstance(),loading=ref(false),rows=ref([]),total=ref(0),customers=ref([]),open=ref(false)
const query=reactive({pageNum:1,pageSize:10,shipmentNo:'',customerId:'',paymentStatus:''}),form=reactive({shipmentId:'',items:[]})
const chargeTotal=computed(()=>form.items.reduce((sum,item)=>sum+Number(item.amount||0),0))
const statusLabel={UNPAID:'未付款',PARTIAL:'部分付款',PAID:'已付款'},statusColor={UNPAID:'red',PARTIAL:'orange',PAID:'green'},money=v=>Number(v||0).toFixed(2)
function load(){loading.value=true;listPaymentLedger(query).then(r=>{rows.value=r.data?.rows||[];total.value=r.data?.total||0}).finally(()=>loading.value=false)}
function search(){query.pageNum=1;load()} function reset(){Object.assign(query,{pageNum:1,pageSize:10,shipmentNo:'',customerId:'',paymentStatus:''});load()}
function addCharge(){form.items.push({feeName:'',amount:0,currency:'CNY',remark:''})}
function edit(row){form.shipmentId=row.shipmentId;getShipmentCharges(row.shipmentId).then(r=>{form.items=r.data||[];if(!form.items.length)addCharge();open.value=true})}
function save(){const items=form.items.filter(item=>item.feeName&&Number(item.amount)>0);saveShipmentCharges(form.shipmentId,items).then(()=>{proxy.$modal.msgSuccess('费用明细已更新');open.value=false;load()})}
simpleCustomerOptions({}).then(r=>customers.value=r.data||[]);load()
</script>
<style scoped>.app-container :deep(.vxe-pager){margin-top:12px}.charge-toolbar{display:flex;align-items:center;justify-content:space-between;margin-bottom:12px}</style>
