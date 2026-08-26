<template>
  <div class="app-container shipment-detail" v-loading="loading">
    <div class="page-header">
      <div><vxe-button @click="router.back()">返回</vxe-button><span class="page-title">出货计划详情</span></div>
      <div>
        <vxe-button v-if="editing" @click="editing = false">取消</vxe-button>
        <vxe-button v-if="editing" status="primary" :loading="saving" @click="save">保存</vxe-button>
        <vxe-button v-else status="primary" v-hasPermi="['freight:shipment:edit']" @click="startEdit">编辑</vxe-button>
      </div>
    </div>
    <a-card title="计划信息" :bordered="false">
      <a-descriptions v-if="!editing" :column="3" bordered>
        <a-descriptions-item label="出货计划号">{{ plan.shipmentNo || '-' }}</a-descriptions-item>
        <a-descriptions-item label="客户单号">{{ plan.orderNo || '-' }}</a-descriptions-item>
        <a-descriptions-item label="客户">{{ plan.customerName || '-' }}</a-descriptions-item>
        <a-descriptions-item label="起运港">{{ plan.pol || '-' }}</a-descriptions-item>
        <a-descriptions-item label="目的港">{{ plan.pod || '-' }}</a-descriptions-item>
        <a-descriptions-item label="业务员">{{ plan.salesUserName || '-' }}</a-descriptions-item>
        <a-descriptions-item label="船名">{{ plan.vesselName || '-' }}</a-descriptions-item>
        <a-descriptions-item label="航次">{{ plan.voyageNo || '-' }}</a-descriptions-item>
        <a-descriptions-item label="计划离港">{{ plan.plannedEtd || '-' }}</a-descriptions-item>
        <a-descriptions-item label="计划到港">{{ plan.plannedEta || '-' }}</a-descriptions-item>
        <a-descriptions-item label="实际离港">{{ plan.actualEtd || '-' }}</a-descriptions-item>
        <a-descriptions-item label="实际到港">{{ plan.actualEta || '-' }}</a-descriptions-item>
        <a-descriptions-item label="状态">{{ statusLabel }}</a-descriptions-item>
        <a-descriptions-item label="收款状态">{{ paymentStatusLabel }}</a-descriptions-item>
        <a-descriptions-item label="已收金额">{{ plan.currency || 'CNY' }} {{ Number(plan.paymentAmount || 0).toFixed(2) }}</a-descriptions-item>
        <a-descriptions-item label="装载方式">{{ loadingModeLabel }}</a-descriptions-item>
        <a-descriptions-item label="结算币种">{{ plan.currency || 'CNY' }}</a-descriptions-item>
        <a-descriptions-item label="贸易条款">{{ plan.tradeTerm || '-' }}</a-descriptions-item>
        <a-descriptions-item label="运输范围">{{ deliveryTypeLabel(plan.deliveryType) }}</a-descriptions-item>
        <a-descriptions-item v-if="plan.pickupAddress" label="提货地址" :span="3">{{ plan.pickupAddress }}</a-descriptions-item>
        <a-descriptions-item v-if="plan.handoverLocation" label="交货地点" :span="3">{{ plan.handoverLocation }}</a-descriptions-item>
        <a-descriptions-item v-if="plan.deliveryAddress" label="送货地址" :span="3">{{ plan.deliveryAddress }}</a-descriptions-item>
        <a-descriptions-item v-if="plan.clearanceParty" label="目的国清关">{{ partyLabel(plan.clearanceParty) }}</a-descriptions-item>
        <a-descriptions-item v-if="plan.dutyPayer" label="税费承担方">{{ partyLabel(plan.dutyPayer) }}</a-descriptions-item>
        <a-descriptions-item label="总箱数">{{ plan.totalCartons || 0 }}</a-descriptions-item>
        <a-descriptions-item label="总重量">{{ plan.totalWeight || 0 }} kg</a-descriptions-item>
        <a-descriptions-item label="总体积">{{ plan.totalVolume || 0 }} m³</a-descriptions-item>
        <a-descriptions-item label="备注" :span="3">{{ plan.remark || '-' }}</a-descriptions-item>
      </a-descriptions>
      <vxe-form v-else :data="form" title-width="90">
        <vxe-form-item title="出货计划号" span="8"><vxe-input :model-value="plan.shipmentNo" disabled /></vxe-form-item>
        <vxe-form-item title="客户单号" span="8"><vxe-input :model-value="plan.orderNo" disabled /></vxe-form-item>
        <vxe-form-item title="客户" span="8"><vxe-select v-model="form.customerId" filterable><vxe-option v-for="item in customers" :key="item.customerId" :value="item.customerId" :label="`${item.customerName} / ${item.companyName || '-'}`" /></vxe-select></vxe-form-item>
        <vxe-form-item title="业务员" span="8"><vxe-input :model-value="plan.salesUserName || '-'" disabled /></vxe-form-item>
        <vxe-form-item title="当前状态" span="8"><vxe-input :model-value="statusLabel" disabled /></vxe-form-item>
        <vxe-form-item title="起运港" span="8"><vxe-input v-model="form.pol" /></vxe-form-item>
        <vxe-form-item title="目的港" span="8"><vxe-input v-model="form.pod" /></vxe-form-item>
        <vxe-form-item title="计划离港" span="8"><vxe-date-picker v-model="form.plannedEtd" type="date" value-format="YYYY-MM-DD" /></vxe-form-item>
        <vxe-form-item title="计划到港" span="8"><vxe-date-picker v-model="form.plannedEta" type="date" value-format="YYYY-MM-DD" /></vxe-form-item>
        <vxe-form-item title="实际离港" span="8"><vxe-input :model-value="plan.actualEtd || '-'" disabled /></vxe-form-item>
        <vxe-form-item title="实际到港" span="8"><vxe-input :model-value="plan.actualEta || '-'" disabled /></vxe-form-item>
        <vxe-form-item title="船名" span="8"><vxe-input v-model="form.vesselName" /></vxe-form-item>
        <vxe-form-item title="航次" span="8"><vxe-input v-model="form.voyageNo" /></vxe-form-item>
        <vxe-form-item title="装载方式" span="8"><vxe-select v-model="form.preferredType" clearable placeholder="系统自动"><vxe-option v-for="item in containerTypeOptions" :key="item.value" :value="item.value" :label="containerOptionLabel(item)" /></vxe-select></vxe-form-item>
        <vxe-form-item title="结算币种" span="8"><vxe-select v-model="form.currency"><vxe-option v-for="item in currencyOptions" :key="item.value" :value="item.value" :label="item.label" /></vxe-select></vxe-form-item>
        <vxe-form-item title="付款金额" span="8"><vxe-number-input v-model="form.paymentAmount" type="float" :min="0" :digits="2" controls /></vxe-form-item>
        <vxe-form-item title="收款状态" span="8"><vxe-input :model-value="paymentStatusLabel" disabled /></vxe-form-item>
        <vxe-form-item title="贸易条款" span="8"><vxe-select v-model="form.tradeTerm" clearable><vxe-option v-for="item in tradeTermOptions" :key="item" :value="item" :label="item" /></vxe-select></vxe-form-item>
        <vxe-form-item title="运输范围" span="8"><vxe-select v-model="form.deliveryType" clearable><vxe-option v-for="item in deliveryTypeOptions" :key="item.value" :value="item.value" :label="item.label" /></vxe-select></vxe-form-item>
        <vxe-form-item v-if="showPickupFields" title="提货地址" span="24"><vxe-textarea v-model="form.pickupAddress" placeholder="请输入工厂、仓库或约定提货地址" /></vxe-form-item>
        <vxe-form-item v-if="form.tradeTerm === 'FCA'" title="交货地点" span="24"><vxe-input v-model="form.handoverLocation" /></vxe-form-item>
        <vxe-form-item v-if="showDeliveryFields" title="送货地址" span="24"><vxe-textarea v-model="form.deliveryAddress" placeholder="请输入最终收货地址" /></vxe-form-item>
        <vxe-form-item v-if="form.tradeTerm === 'DDP'" title="目的国清关" span="12"><vxe-select v-model="form.clearanceParty"><vxe-option value="SELLER" label="卖方负责" /><vxe-option value="SERVICE_PROVIDER" label="物流服务商负责" /></vxe-select></vxe-form-item>
        <vxe-form-item v-if="form.tradeTerm === 'DDP'" title="税费承担方" span="12"><vxe-select v-model="form.dutyPayer"><vxe-option value="SELLER" label="卖方承担" /><vxe-option value="SERVICE_PROVIDER" label="物流服务商代缴" /></vxe-select></vxe-form-item>
        <vxe-form-item title="总箱数" span="8"><vxe-input :model-value="editSummary.cartons" disabled /></vxe-form-item>
        <vxe-form-item title="总重量(kg)" span="8"><vxe-input :model-value="editSummary.weight" disabled /></vxe-form-item>
        <vxe-form-item title="总体积(m³)" span="8"><vxe-input :model-value="editSummary.volume" disabled /></vxe-form-item>
        <vxe-form-item title="备注" span="24"><vxe-textarea v-model="form.remark" :rows="3" /></vxe-form-item>
      </vxe-form>
    </a-card>
    <a-card title="货物明细" :bordered="false" class="section-card">
      <template #extra><vxe-button v-if="editing" size="mini" @click="addCargo">增加货物</vxe-button></template>
      <vxe-table border stripe :data="editing ? form.cargoList : detail.cargoList">
        <vxe-column field="cargoName" title="货物名称" min-width="150"><template v-if="editing" #default="{ row }"><vxe-input v-model="row.cargoName" /></template></vxe-column>
        <vxe-column field="sku" title="SKU" min-width="110"><template v-if="editing" #default="{ row }"><vxe-input v-model="row.sku" /></template></vxe-column>
        <vxe-column field="packageType" title="包装" width="110"><template v-if="editing" #default="{ row }"><vxe-input v-model="row.packageType" /></template></vxe-column>
        <vxe-column field="quantity" title="数量" width="100"><template v-if="editing" #default="{ row }"><vxe-number-input v-model="row.quantity" :min="0" /></template></vxe-column>
        <vxe-column field="cartons" title="箱数" width="100"><template v-if="editing" #default="{ row }"><vxe-number-input v-model="row.cartons" :min="0" /></template></vxe-column>
        <vxe-column field="lengthCm" title="长(cm)" width="105"><template v-if="editing" #default="{ row }"><vxe-number-input v-model="row.lengthCm" type="float" :min="0" :digits="2" /></template></vxe-column>
        <vxe-column field="widthCm" title="宽(cm)" width="105"><template v-if="editing" #default="{ row }"><vxe-number-input v-model="row.widthCm" type="float" :min="0" :digits="2" /></template></vxe-column>
        <vxe-column field="heightCm" title="高(cm)" width="105"><template v-if="editing" #default="{ row }"><vxe-number-input v-model="row.heightCm" type="float" :min="0" :digits="2" /></template></vxe-column>
        <vxe-column field="unitWeightKg" title="单个重量(kg)" width="135"><template v-if="editing" #default="{ row }"><vxe-number-input v-model="row.unitWeightKg" type="float" :min="0" :digits="4" /></template></vxe-column>
        <vxe-column field="unitVolumeCbm" title="单个体积(m³)" width="140"><template v-if="editing" #default="{ row }"><vxe-number-input v-model="row.unitVolumeCbm" type="float" :min="0" :digits="6" /></template></vxe-column>
        <vxe-column field="weightKg" title="总重量(kg)" width="125"><template #default="{ row }">{{ editing ? cargoTotalWeight(row) : row.weightKg }}</template></vxe-column>
        <vxe-column field="volumeCbm" title="总体积(m³)" width="125"><template #default="{ row }">{{ editing ? cargoTotalVolume(row) : row.volumeCbm }}</template></vxe-column>
        <vxe-column v-if="editing" title="操作" width="80"><template #default="{ rowIndex }"><vxe-button mode="text" status="error" @click="removeCargo(rowIndex)">删除</vxe-button></template></vxe-column>
      </vxe-table>
    </a-card>
    <a-card title="推荐柜型" :bordered="false" class="section-card">
      <vxe-table border :data="detail.containers || []"><vxe-column field="containerType" title="柜型" /><vxe-column field="quantity" title="数量" /><vxe-column field="maxVolume" title="可装体积(m³)" /><vxe-column field="maxWeight" title="可装重量(kg)" /><vxe-column field="usedVolume" title="使用体积(m³)" /><vxe-column field="usedWeight" title="使用重量(kg)" /><vxe-column field="loadRate" title="装载率(%)" /></vxe-table>
    </a-card>
    <a-card title="收付款记录" :bordered="false" class="section-card">
      <vxe-table border :data="detail.payments || []"><vxe-column field="paymentTime" title="付款时间" /><vxe-column field="amount" title="金额" /><vxe-column field="currency" title="币种" /><vxe-column field="paymentMethod" title="付款方式" /><vxe-column field="remark" title="备注" /></vxe-table>
    </a-card>
  </div>
</template>
<script setup name="FreightShipmentDetail">
import { computed, getCurrentInstance, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getShipment, updateShipment } from '@/api/freight/shipment'
import { simpleCustomerOptions } from '@/api/simple/simple'
const route = useRoute(), router = useRouter(), { proxy } = getCurrentInstance()
const { freight_container_type } = proxy.useDict('freight_container_type')
const loading = ref(false), saving = ref(false), editing = ref(false), customers = ref([])
const detail = ref({ plan: {}, cargoList: [], containers: [], payments: [] })
const form = reactive({ customerId: '', customerName: '', pol: '', pod: '', plannedEtd: '', plannedEta: '', vesselName: '', voyageNo: '', preferredType: '', currency: 'CNY', paymentAmount: 0, tradeTerm: '', deliveryType: '', pickupAddress: '', deliveryAddress: '', handoverLocation: '', clearanceParty: '', dutyPayer: '', remark: '', cargoList: [] })
const currencyOptions = [{ value: 'CNY', label: 'CNY 人民币' }, { value: 'USD', label: 'USD 美元' }, { value: 'EUR', label: 'EUR 欧元' }, { value: 'GBP', label: 'GBP 英镑' }, { value: 'JPY', label: 'JPY 日元' }, { value: 'HKD', label: 'HKD 港币' }, { value: 'SGD', label: 'SGD 新加坡元' }, { value: 'AUD', label: 'AUD 澳大利亚元' }, { value: 'CAD', label: 'CAD 加拿大元' }]
const tradeTermOptions = ['EXW', 'FCA', 'FOB', 'CFR', 'CIF', 'CPT', 'CIP', 'DAP', 'DPU', 'DDP']
const deliveryTypeOptions = [{ value: 'PORT_TO_PORT', label: '港到港' }, { value: 'DOOR_TO_PORT', label: '门到港' }, { value: 'PORT_TO_DOOR', label: '港到门' }, { value: 'DOOR_TO_DOOR', label: '门到门' }]
const deliveryTypeLabel = value => deliveryTypeOptions.find(item => item.value === value)?.label || value || '-'
const defaultContainerCapacities = { LCL: '约15CBM / 3000KG', '20GP': '约28CBM / 21700KG', '40GP': '约58CBM / 26500KG', '40HQ': '约68CBM / 26500KG' }
const containerTypeOptions = computed(() => freight_container_type.value || [])
const containerOptionLabel = item => `${item.label}（${item.remark || defaultContainerCapacities[item.value] || '容量未配置'}）`
const partyLabel = value => ({ SELLER: '卖方', BUYER: '买方', SERVICE_PROVIDER: '物流服务商' }[value] || value)
const plan = computed(() => detail.value.plan || {})
const statusLabel = computed(() => (detail.value.statusFlow || []).find(item => item.value === plan.value.status)?.label || plan.value.status || '-')
const paymentStatusLabel = computed(() => ({ UNPAID: '未收款', PARTIAL: '部分收款', PAID: '已收款' }[plan.value.paymentStatus] || plan.value.paymentStatus || '未收款'))
const loadingMode = computed(() => detail.value.containers?.[0]?.containerType || '')
const loadingModeLabel = computed(() => loadingMode.value === 'LCL' ? 'LCL 拼箱' : (loadingMode.value ? `${loadingMode.value} 整柜` : '系统自动'))
const showPickupFields = computed(() => form.tradeTerm === 'EXW' || ['DOOR_TO_PORT', 'DOOR_TO_DOOR'].includes(form.deliveryType))
const showDeliveryFields = computed(() => ['DAP', 'DPU', 'DDP'].includes(form.tradeTerm) || ['PORT_TO_DOOR', 'DOOR_TO_DOOR'].includes(form.deliveryType))
function load() { loading.value = true; getShipment(route.params.shipmentId).then(r => { detail.value = r.data || detail.value; if (route.query.edit === '1') startEdit() }).finally(() => loading.value = false) }
function startEdit() { const p = plan.value; Object.assign(form, { customerId: p.customerId || '', customerName: p.customerName || '', pol: p.pol || '', pod: p.pod || '', plannedEtd: p.plannedEtd || '', plannedEta: p.plannedEta || '', vesselName: p.vesselName || '', voyageNo: p.voyageNo || '', preferredType: loadingMode.value, currency: p.currency || 'CNY', paymentAmount: Number(p.paymentAmount || 0), tradeTerm: p.tradeTerm || '', deliveryType: p.deliveryType || '', pickupAddress: p.pickupAddress || '', deliveryAddress: p.deliveryAddress || '', handoverLocation: p.handoverLocation || '', clearanceParty: p.clearanceParty || '', dutyPayer: p.dutyPayer || '', remark: p.remark || '', cargoList: (detail.value.cargoList || []).map(item => ({ ...item })) }); editing.value = true }
function addCargo() { form.cargoList.push({ cargoName: '', sku: '', packageType: '', quantity: 0, cartons: 0, weightKg: 0, volumeCbm: 0, lengthCm: 0, widthCm: 0, heightCm: 0, unitWeightKg: 0, unitVolumeCbm: 0 }) }
function removeCargo(index) { form.cargoList.splice(index, 1); if (!form.cargoList.length) addCargo() }
function cargoCount(row) { return Number(row.quantity || row.cartons || 0) }
function cargoUnitVolume(row) { return Number(row.unitVolumeCbm || 0) || Number(row.lengthCm || 0) * Number(row.widthCm || 0) * Number(row.heightCm || 0) / 1000000 }
function cargoTotalWeight(row) { return Number(Number(row.unitWeightKg || 0) * cargoCount(row)).toFixed(2) }
function cargoTotalVolume(row) { return Number(cargoUnitVolume(row) * cargoCount(row)).toFixed(4) }
const editSummary = computed(() => { const total = form.cargoList.reduce((sum, item) => ({ cartons: sum.cartons + Number(item.cartons || 0), weight: sum.weight + Number(cargoTotalWeight(item)), volume: sum.volume + Number(cargoTotalVolume(item)) }), { cartons: 0, weight: 0, volume: 0 }); return { cartons: total.cartons, weight: total.weight.toFixed(2), volume: total.volume.toFixed(4) } })
function save() { if (!form.customerId) return proxy.$modal.msgWarning('请选择客户'); if (form.plannedEtd && form.plannedEta && form.plannedEta < form.plannedEtd) return proxy.$modal.msgWarning('计划到港日期不能早于计划离港日期'); const cargoList = form.cargoList.filter(i => i.cargoName).map(item => ({ ...item, unitVolumeCbm: cargoUnitVolume(item), weightKg: Number(cargoTotalWeight(item)), volumeCbm: Number(cargoTotalVolume(item)) })); if (!cargoList.length) return proxy.$modal.msgWarning('请至少填写一条货物明细'); saving.value = true; updateShipment(route.params.shipmentId, { ...form, cargoList }).then(r => { detail.value = r.data || detail.value; editing.value = false; proxy.$modal.msgSuccess('出货计划已更新') }).finally(() => saving.value = false) }
simpleCustomerOptions({}).then(r => customers.value = r.data || [])
load()
</script>
<style scoped>
.shipment-detail{background:#f5f7fa;min-height:calc(100vh - 84px)}.page-header{display:flex;align-items:center;justify-content:space-between;margin-bottom:16px}.page-title{margin-left:12px;font-size:20px;font-weight:600}.section-card{margin-top:16px}
</style>
