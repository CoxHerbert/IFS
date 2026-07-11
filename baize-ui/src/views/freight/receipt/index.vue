<template>
  <div class="app-container receipt-page">
    <vxe-form :data="query" @submit="search" @reset="reset">
      <vxe-form-item field="receiptNo" title="收款单号"><template #default><vxe-input v-model="query.receiptNo" clearable /></template></vxe-form-item>
      <vxe-form-item field="customerName" title="客户"><template #default><vxe-input v-model="query.customerName" clearable /></template></vxe-form-item>
      <vxe-form-item><template #default><vxe-button type="submit" status="primary">查询</vxe-button><vxe-button type="reset">重置</vxe-button></template></vxe-form-item>
    </vxe-form>
    <vxe-button status="primary" class="add-button" @click="openCreate">新增收款</vxe-button>
    <vxe-table border stripe :loading="loading" :data="rows">
      <vxe-column field="receiptNo" title="收款单号" width="180" />
      <vxe-column field="customerName" title="客户" min-width="150" />
      <vxe-column title="收款金额" width="140" align="right"><template #default="{row}">{{ row.currency }} {{ money(row.amount) }}</template></vxe-column>
      <vxe-column title="已核销" width="130" align="right"><template #default="{row}">{{ money(row.allocatedAmount) }}</template></vxe-column>
      <vxe-column title="核销状态" width="110"><template #default="{row}">{{ statusLabel[row.status] || row.status }}</template></vxe-column>
      <vxe-column field="receiptTime" title="收款时间" width="165" />
      <vxe-column field="paymentMethod" title="收款方式" width="110" />
      <vxe-column title="凭证" width="100"><template #default="{row}"><a v-if="row.voucherUrl" :href="row.voucherUrl" target="_blank">查看</a><span v-else>-</span></template></vxe-column>
      <vxe-column title="操作" width="130"><template #default="{row}"><vxe-button mode="text" @click="showDetail(row)">详情</vxe-button><vxe-button mode="text" status="error" @click="remove(row)">删除</vxe-button></template></vxe-column>
    </vxe-table>
    <vxe-pager v-model:current-page="query.pageNum" v-model:page-size="query.pageSize" :total="total" @page-change="load" />

    <vxe-modal v-model="createOpen" title="新增收款" width="760" show-footer>
      <vxe-form :data="form" title-width="90">
        <vxe-form-item title="客户" span="24" :item-render="{}"><template #default><vxe-select v-model="form.customerId" filterable @change="customerChanged"><vxe-option v-for="c in customers" :key="c.customerId" :value="c.customerId" :label="c.customerName+' / '+(c.companyName||'-')" /></vxe-select></template></vxe-form-item>
        <vxe-form-item title="收款金额" span="12" :item-render="{}"><template #default><vxe-number-input v-model="form.amount" type="float" min="0" :digits="2" /></template></vxe-form-item>
        <vxe-form-item title="币种" span="12" :item-render="{}"><template #default><vxe-input v-model="form.currency" /></template></vxe-form-item>
        <vxe-form-item title="收款时间" span="12" :item-render="{}"><template #default><vxe-date-picker v-model="form.receiptTime" type="datetime" value-format="YYYY-MM-DD HH:mm:ss" /></template></vxe-form-item>
        <vxe-form-item title="收款方式" span="12" :item-render="{}"><template #default><vxe-select v-model="form.paymentMethod"><vxe-option value="BANK_TRANSFER" label="银行转账"/><vxe-option value="CASH" label="现金"/><vxe-option value="OTHER" label="其他"/></vxe-select></template></vxe-form-item>
        <vxe-form-item title="收款凭证" span="24" :item-render="{}"><template #default><a-upload v-model:file-list="files" :before-upload="()=>false" :max-count="1" accept=".pdf,.png,.jpg,.jpeg"><a-button>选择文件</a-button></a-upload><span class="hint"> PDF/PNG/JPG，最大10MB</span></template></vxe-form-item>
        <vxe-form-item title="备注" span="24" :item-render="{}"><template #default><vxe-textarea v-model="form.remark" /></template></vxe-form-item>
      </vxe-form>
      <div class="allocation-title"><b>核销出货计划</b><vxe-button size="mini" @click="addAllocation">增加计划</vxe-button></div>
      <vxe-table border :data="form.allocations">
        <vxe-column title="出货计划"><template #default="{row}"><vxe-select v-model="row.shipmentId" filterable><vxe-option v-for="s in shipments" :key="s.shipmentId" :value="s.shipmentId" :label="s.shipmentNo+' / '+(s.orderNo||'-')" /></vxe-select></template></vxe-column>
        <vxe-column title="核销金额" width="200"><template #default="{row}"><vxe-number-input v-model="row.allocatedAmount" type="float" min="0" :digits="2" /></template></vxe-column>
        <vxe-column title="操作" width="80"><template #default="{rowIndex}"><vxe-button mode="text" status="error" @click="form.allocations.splice(rowIndex,1)">删除</vxe-button></template></vxe-column>
      </vxe-table>
      <template #footer><vxe-button @click="createOpen=false">取消</vxe-button><vxe-button status="primary" :loading="saving" @click="submit">保存</vxe-button></template>
    </vxe-modal>

    <vxe-modal v-model="detailOpen" title="收款详情" width="700" :show-footer="false">
      <div v-if="detail" class="detail-grid"><div>收款单号：{{detail.receiptNo}}</div><div>客户：{{detail.customerName}}</div><div>金额：{{detail.currency}} {{money(detail.amount)}}</div><div>已核销：{{money(detail.allocatedAmount)}}</div></div>
      <vxe-table border :data="detail?.allocations||[]"><vxe-column field="shipmentNo" title="出货计划"/><vxe-column field="allocatedAmount" title="核销金额"/></vxe-table>
    </vxe-modal>
  </div>
</template>

<script setup name="FreightReceipt">
import { getCurrentInstance, reactive, ref } from 'vue'
import { listReceipt,getReceipt,addReceipt,delReceipt } from '@/api/freight/receipt'
import { customerOptions } from '@/api/customer/customer'
import { listShipment } from '@/api/freight/shipment'
const {proxy}=getCurrentInstance(); const loading=ref(false),saving=ref(false),createOpen=ref(false),detailOpen=ref(false)
const rows=ref([]),total=ref(0),customers=ref([]),shipments=ref([]),files=ref([]),detail=ref(null)
const query=reactive({pageNum:1,pageSize:10,receiptNo:'',customerName:''})
const form=reactive({customerId:'',customerName:'',amount:0,currency:'CNY',receiptTime:'',paymentMethod:'BANK_TRANSFER',remark:'',allocations:[]})
const statusLabel={UNALLOCATED:'未核销',PARTIAL:'部分核销',ALLOCATED:'已核销'}; const money=v=>Number(v||0).toFixed(2)
function load(){loading.value=true;listReceipt(query).then(r=>{rows.value=r.data?.rows||[];total.value=r.data?.total||0}).finally(()=>loading.value=false)}
function search(){query.pageNum=1;load()} function reset(){Object.assign(query,{pageNum:1,pageSize:10,receiptNo:'',customerName:''});load()}
function openCreate(){Object.assign(form,{customerId:'',customerName:'',amount:0,currency:'CNY',receiptTime:'',paymentMethod:'BANK_TRANSFER',remark:'',allocations:[{shipmentId:'',allocatedAmount:0}]});files.value=[];shipments.value=[];createOpen.value=true;customerOptions({}).then(r=>customers.value=r.data||[])}
function customerChanged({value}){const c=customers.value.find(x=>x.customerId===value);form.customerName=c?.customerName||'';form.allocations=[{shipmentId:'',allocatedAmount:form.amount||0}];listShipment({pageNum:1,pageSize:100,customerId:value}).then(r=>shipments.value=r.data?.rows||[])}
function addAllocation(){form.allocations.push({shipmentId:'',allocatedAmount:0})}
function submit(){if(!form.customerId||Number(form.amount)<=0){proxy.$modal.msgWarning('请选择客户并填写收款金额');return}const allocations=form.allocations.filter(x=>x.shipmentId&&Number(x.allocatedAmount)>0);if(allocations.reduce((n,x)=>n+Number(x.allocatedAmount),0)>Number(form.amount)){proxy.$modal.msgWarning('核销合计不能超过收款金额');return}const fd=new FormData();['customerId','customerName','amount','currency','receiptTime','paymentMethod','remark'].forEach(k=>fd.append(k,String(form[k]??'')));fd.append('allocations',JSON.stringify(allocations));const file=files.value[0]?.originFileObj;if(file)fd.append('voucher',file);saving.value=true;addReceipt(fd).then(()=>{proxy.$modal.msgSuccess('收款单已创建');createOpen.value=false;load()}).finally(()=>saving.value=false)}
function showDetail(row){getReceipt(row.receiptId).then(r=>{detail.value=r.data;detailOpen.value=true})}
function remove(row){proxy.$modal.confirm('确认删除该收款单及其核销记录？').then(()=>delReceipt(row.receiptId)).then(()=>{proxy.$modal.msgSuccess('删除成功');load()}).catch(()=>{})}
load()
</script>
<style scoped>.add-button{margin:10px 0}.receipt-page :deep(.vxe-pager){margin-top:12px}.allocation-title{display:flex;justify-content:space-between;align-items:center;margin:12px 0}.hint{color:#909399;font-size:12px}.detail-grid{display:grid;grid-template-columns:1fr 1fr;gap:12px;margin-bottom:16px}</style>
