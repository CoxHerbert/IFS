<template>
  <div class="app-container shipment-page">
    <vxe-form v-show="showSearch" :data="queryParams" class="shipment-search-form" @submit="handleQuery"
      @reset="resetQuery">
      <vxe-form-item title="计划编号" field="shipmentNo" span="6" :item-render="{}">
        <template #default>
          <vxe-input v-model="queryParams.shipmentNo" clearable placeholder="请输入计划编号" />
        </template>
      </vxe-form-item>
      <vxe-form-item title="客户名称" field="customerName" span="6" :item-render="{}">
        <template #default>
          <vxe-input v-model="queryParams.customerName" clearable placeholder="请输入客户名称" />
        </template>
      </vxe-form-item>
      <vxe-form-item title="目的港" field="pod" span="6" :item-render="{}">
        <template #default>
          <vxe-input v-model="queryParams.pod" clearable placeholder="请输入目的港" />
        </template>
      </vxe-form-item>
      <vxe-form-item title="状态" field="status" span="6" :item-render="{}">
        <template #default>
          <vxe-select v-model="queryParams.status" clearable placeholder="请选择状态">
            <vxe-option v-for="item in statusOptions" :key="item.value" :value="item.value" :label="item.label" />
          </vxe-select>
        </template>
      </vxe-form-item>
      <vxe-form-item span="24" align="left" :item-render="{}" class-name="search-actions">
        <template #default>
          <vxe-button type="submit" status="primary">搜索</vxe-button>
          <vxe-button type="reset">重置</vxe-button>
        </template>
      </vxe-form-item>
    </vxe-form>

    <div v-show="showSearch" class="toolbar-spacer"></div>

    <vxe-toolbar class="shipment-toolbar">
      <template #buttons>
        <vxe-button status="primary" @click="handleImport" v-hasPermi="['freight:shipment:import']">导入清单</vxe-button>
        <vxe-button status="error" :disabled="multiple" @click="handleDelete"
          v-hasPermi="['freight:shipment:remove']">删除</vxe-button>
        <vxe-button :disabled="single" @click="handleBindCustomer()"
          v-hasPermi="['freight:shipment:edit']">绑定客户</vxe-button>
      </template>
      <template #tools>
        <right-toolbar v-model:showSearch="showSearch" @queryTable="getList" />
      </template>
    </vxe-toolbar>

    <vxe-table ref="shipmentTableRef" border stripe auto-resize show-overflow="title" :loading="loading"
      :data="shipmentList" :row-config="{ keyField: 'shipmentId' }" :checkbox-config="{ reserve: true }"
      @checkbox-change="handleSelectionChange" @checkbox-all="handleSelectionChange">
      <vxe-column type="checkbox" width="54" align="center" />
      <vxe-column field="shipmentNo" title="计划编号" width="170" align="center" />
      <vxe-column field="customerName" title="客户" width="150" align="center" show-overflow="tooltip" />
      <vxe-column field="orderNo" title="客户订单号" width="150" align="center" />
      <vxe-column field="paymentStatus" title="付款状态" width="120" align="center">
        <template #default="{ row }">
          {{ paymentStatusOptions.find(item => item.value === row.paymentStatus)?.label || "未付款" }}
        </template>
      </vxe-column>
      <vxe-column field="paymentAmount" title="付款金额" width="130" align="center">
        <template #default="{ row }">
          {{ Number(row.paymentAmount || 0).toFixed(2) }}
        </template>
      </vxe-column>
      <vxe-column title="航线" width="170" align="center">
        <template #default="{ row }">
          {{ row.pol || "-" }} -> {{ row.pod || "-" }}
        </template>
      </vxe-column>
      <vxe-column title="货量" width="180" align="center">
        <template #default="{ row }">
          {{ row.totalCartons }}箱 / {{ row.totalVolume }}CBM
        </template>
      </vxe-column>
      <vxe-column title="重量" width="120" align="center">
        <template #default="{ row }">
          {{ row.totalWeight }}KG
        </template>
      </vxe-column>
      <vxe-column field="status" title="状态" width="150" align="center">
        <template #default="{ row }">
          <dict-tag :options="freight_shipment_status" :value="row.status" />
        </template>
      </vxe-column>
      <vxe-column field="createTime" title="创建时间" width="170" align="center">
        <template #default="{ row }">
          {{ parseTime(row.createTime) }}
        </template>
      </vxe-column>
      <vxe-column title="操作" width="240" align="center" fixed="right">
        <template #default="{ row }">
          <div class="action-links">
            <vxe-button mode="text" status="primary" @click="handleDetail(row)">详情</vxe-button>
            <vxe-button mode="text" status="primary" @click="handleStatus(row)"
              v-hasPermi="['freight:shipment:edit']">状态</vxe-button>
            <vxe-button mode="text" status="primary" @click="handleConfirm(row)"
              v-hasPermi="['freight:shipment:confirm']">出货单</vxe-button>
            <vxe-button mode="text" status="primary" @click="handleShare(row)"
              v-hasPermi="['freight:shipment:share']">分享</vxe-button>
          </div>
        </template>
      </vxe-column>
    </vxe-table>

    <div v-show="total > 0" class="shipment-pagination">
      <vxe-pager :current-page="queryParams.pageNum" :page-size="queryParams.pageSize" :total="total"
        :page-sizes="[10, 20, 30, 50]" :layouts="['PrevPage', 'JumpNumber', 'NextPage', 'Sizes', 'FullJump', 'Total']"
        @page-change="handlePageChange" />
    </div>

    <vxe-modal v-model="importOpen" title="导入出货清单并生成智能计划" width="980" show-footer esc-closable mask-closable="false">
      <vxe-form :data="importForm" title-width="88">
        <vxe-form-item title="客户" field="customerId" span="12" :item-render="{}">
          <template #default>
            <vxe-select v-model="importForm.customerId" clearable filterable placeholder="搜索客户"
              :loading="customerLoading" @change="({ value }) => handleCustomerChange(value)">
              <vxe-option v-for="item in customerList" :key="item.customerId" :value="item.customerId"
                :label="item.customerName + ' / ' + (item.companyName || '-')" />
            </vxe-select>
          </template>
        </vxe-form-item>
        <vxe-form-item title="客户单号" field="orderNo" span="12" :item-render="{}">
          <template #default>
            <vxe-input v-model="importForm.orderNo" maxlength="64" placeholder="客户订单号/参考号" />
          </template>
        </vxe-form-item>
        <vxe-form-item title="起运港" field="pol" span="8" :item-render="{}">
          <template #default>
            <vxe-input v-model="importForm.pol" placeholder="如 SHANGHAI" />
          </template>
        </vxe-form-item>
        <vxe-form-item title="目的港" field="pod" span="8" :item-render="{}">
          <template #default>
            <vxe-input v-model="importForm.pod" placeholder="如 LOS ANGELES" />
          </template>
        </vxe-form-item>
        <vxe-form-item title="偏好柜型" field="preferredType" span="8" :item-render="{}">
          <template #default>
            <vxe-select v-model="importForm.preferredType" clearable placeholder="系统自动">
              <vxe-option value="20GP" label="20GP" />
              <vxe-option value="40GP" label="40GP" />
              <vxe-option value="40HQ" label="40HQ" />
            </vxe-select>
          </template>
        </vxe-form-item>
        <vxe-form-item title="付款状态" field="paymentStatus" span="8" :item-render="{}">
          <template #default>
            <vxe-select v-model="importForm.paymentStatus" placeholder="请选择付款状态">
              <vxe-option v-for="item in paymentStatusOptions" :key="item.value" :value="item.value" :label="item.label" />
            </vxe-select>
          </template>
        </vxe-form-item>
        <vxe-form-item title="付款金额" field="paymentAmount" span="8" :item-render="{}">
          <template #default>
            <vxe-number-input v-model="importForm.paymentAmount" type="float" min="0" :digits="2" controls />
          </template>
        </vxe-form-item>
      </vxe-form>

      <div class="cargo-toolbar">
        <span>货物明细</span>
        <vxe-button status="primary" @click="addCargoRow">新增一行</vxe-button>
      </div>

      <vxe-table border stripe auto-resize show-overflow="title" :data="importForm.cargoList"
        :row-config="{ keyField: 'rowKey' }">
        <vxe-column field="cargoName" title="货名" min-width="170">
          <template #default="{ row }">
            <vxe-input v-model="row.cargoName" placeholder="货物名称" />
          </template>
        </vxe-column>
        <vxe-column field="sku" title="SKU/唛头" width="140">
          <template #default="{ row }">
            <vxe-input v-model="row.sku" placeholder="可选" />
          </template>
        </vxe-column>
        <vxe-column field="cartons" title="箱数" width="110" align="center">
          <template #default="{ row }">
            <vxe-number-input v-model="row.cartons" type="integer" min="0" controls />
          </template>
        </vxe-column>
        <vxe-column field="weightKg" title="重量KG" width="140" align="center">
          <template #default="{ row }">
            <vxe-number-input v-model="row.weightKg" type="float" min="0" :digits="2" controls />
          </template>
        </vxe-column>
        <vxe-column field="volumeCbm" title="体积CBM" width="140" align="center">
          <template #default="{ row }">
            <vxe-number-input v-model="row.volumeCbm" type="float" min="0" :digits="2" controls />
          </template>
        </vxe-column>
        <vxe-column title="操作" width="90" align="center">
          <template #default="{ rowIndex }">
            <vxe-button mode="text" status="error" @click="removeCargoRow(rowIndex)">删除</vxe-button>
          </template>
        </vxe-column>
      </vxe-table>

      <template #footer>
        <div class="modal-footer">
          <vxe-button @click="importOpen = false">取消</vxe-button>
          <vxe-button status="primary" @click="submitImport">生成计划</vxe-button>
        </div>
      </template>
    </vxe-modal>

    <vxe-modal v-model="detailOpen" title="出货计划详情" width="980" esc-closable :show-footer="false">
      <template v-if="detail.plan">
        <div class="detail-summary">
          <div class="detail-item">
            <span class="detail-label">付款状态</span>
            <span class="detail-value">{{ paymentStatusOptions.find(item => item.value === detail.plan.paymentStatus)?.label || "未付款" }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">付款金额</span>
            <span class="detail-value">{{ Number(detail.plan.paymentAmount || 0).toFixed(2) }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">计划编号</span>
            <span class="detail-value">{{ detail.plan.shipmentNo || "-" }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">客户</span>
            <span class="detail-value">{{ detail.plan.customerName || "-" }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">状态</span>
            <span class="detail-value"><dict-tag :options="freight_shipment_status"
                :value="detail.plan.status" /></span>
          </div>
          <div class="detail-item">
            <span class="detail-label">航线</span>
            <span class="detail-value">{{ detail.plan.pol || "-" }} -> {{ detail.plan.pod || "-" }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">货量</span>
            <span class="detail-value">{{ detail.plan.totalCartons }}箱 / {{ detail.plan.totalVolume }}CBM</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">重量</span>
            <span class="detail-value">{{ detail.plan.totalWeight }}KG</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">出货单</span>
            <span class="detail-value">{{ detail.order?.orderNo || "未生成" }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">计划ETD</span>
            <span class="detail-value">{{ detail.plan.plannedEtd || "-" }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">计划ETA</span>
            <span class="detail-value">{{ detail.plan.plannedEta || "-" }}</span>
          </div>
        </div>

        <h4>智能装柜建议</h4>
        <vxe-table border stripe auto-resize show-overflow="title" :data="detail.containers || []"
          :row-config="{ keyField: 'containerType' }">
          <vxe-column field="containerType" title="柜型" align="center" />
          <vxe-column field="quantity" title="柜量" align="center" />
          <vxe-column title="装载率" align="center">
            <template #default="{ row }">
              {{ row.loadRate }}%
            </template>
          </vxe-column>
          <vxe-column field="remark" title="说明" />
        </vxe-table>

        <h4>货物明细</h4>
        <vxe-table border stripe auto-resize show-overflow="title" :data="detail.cargoList || []"
          :row-config="{ keyField: 'cargoId' }">
          <vxe-column field="cargoName" title="货名" />
          <vxe-column field="sku" title="SKU/唛头" />
          <vxe-column field="cartons" title="箱数" align="center" />
          <vxe-column field="weightKg" title="重量KG" align="center" />
          <vxe-column field="volumeCbm" title="体积CBM" align="center" />
        </vxe-table>
      </template>
    </vxe-modal>

    <vxe-modal v-model="statusOpen" title="维护客户可见状态" width="560" show-footer esc-closable mask-closable="false">
      <vxe-form :data="statusForm" title-width="90">
        <vxe-form-item title="当前状态" field="status" span="24" :item-render="{}">
          <template #default>
            <vxe-select v-model="statusForm.status" placeholder="请选择状态">
              <vxe-option v-for="item in statusOptions" :key="item.value" :value="item.value" :label="item.label" />
            </vxe-select>
          </template>
        </vxe-form-item>
        <vxe-form-item title="实际ETD" field="actualEtd" span="12" :item-render="{}">
          <template #default>
            <vxe-date-picker v-model="statusForm.actualEtd" value-format="YYYY-MM-DD" type="date" placeholder="选择日期"
              clearable />
          </template>
        </vxe-form-item>
        <vxe-form-item title="实际ETA" field="actualEta" span="12" :item-render="{}">
          <template #default>
            <vxe-date-picker v-model="statusForm.actualEta" value-format="YYYY-MM-DD" type="date" placeholder="选择日期"
              clearable />
          </template>
        </vxe-form-item>
        <vxe-form-item title="备注" field="remark" span="24" :item-render="{}">
          <template #default>
            <vxe-textarea v-model="statusForm.remark" rows="3" placeholder="给客户看的简短说明" />
          </template>
        </vxe-form-item>
      </vxe-form>

      <template #footer>
        <div class="modal-footer">
          <vxe-button @click="statusOpen = false">取消</vxe-button>
          <vxe-button status="primary" @click="submitStatus">保存</vxe-button>
        </div>
      </template>
    </vxe-modal>

    <vxe-modal v-model="bindCustomerOpen" title="绑定客户" width="560" show-footer esc-closable mask-closable="false">
      <vxe-form :data="bindCustomerForm" title-width="90">
        <vxe-form-item title="客户" field="customerId" span="24" :item-render="{}">
          <template #default>
            <vxe-select v-model="bindCustomerForm.customerId" clearable filterable placeholder="搜索客户"
              :loading="customerLoading" @change="({ value }) => handleBindCustomerChange(value)">
              <vxe-option v-for="item in customerList" :key="item.customerId" :value="item.customerId"
                :label="item.customerName + ' / ' + (item.companyName || '-')" />
            </vxe-select>
          </template>
        </vxe-form-item>
      </vxe-form>

      <template #footer>
        <div class="modal-footer">
          <vxe-button @click="bindCustomerOpen = false">取消</vxe-button>
          <vxe-button status="primary" @click="submitBindCustomer">保存</vxe-button>
        </div>
      </template>
    </vxe-modal>

    <vxe-modal v-model="shareOpen" title="客户免登录分享链接" width="680" show-footer esc-closable>
      <div class="share-modal">
        <div class="share-tip">链接已自动复制，点击下方按钮可再次复制。</div>
        <vxe-input v-model="shareUrl" readonly />
      </div>
      <template #footer>
        <div class="modal-footer">
          <vxe-button @click="shareOpen = false">关闭</vxe-button>
          <vxe-button status="primary" @click="copyShareUrl(shareUrl)">复制链接</vxe-button>
        </div>
      </template>
    </vxe-modal>
  </div>
</template>

<script setup name="FreightShipment">
import { computed, getCurrentInstance, reactive, ref, toRefs } from "vue";
import {
  listShipment,
  importShipment,
  getShipment,
  updateShipmentStatus,
  bindShipmentCustomer,
  confirmShipment,
  getShipmentShare,
  delShipment
} from "@/api/freight/shipment";
import { customerOptions } from "@/api/customer/customer";

const { proxy } = getCurrentInstance();
const { freight_shipment_status } = proxy.useDict("freight_shipment_status");
const portalBaseUrl = (import.meta.env.VITE_PORTAL_BASE_URL || window.location.origin).replace(/\/$/, "");
const paymentStatusOptions = [
  { value: "UNPAID", label: "未付款" },
  { value: "PARTIAL", label: "部分付款" },
  { value: "PAID", label: "已付款" }
];

const shipmentTableRef = ref();

const statusOptions = computed(() =>
  (freight_shipment_status.value || []).map(item => ({
    value: item.value,
    label: item.label
  }))
);

const shipmentList = ref([]);
const loading = ref(true);
const showSearch = ref(true);
const ids = ref([]);
const single = ref(true);
const multiple = ref(true);
const total = ref(0);
const importOpen = ref(false);
const detailOpen = ref(false);
const statusOpen = ref(false);
const bindCustomerOpen = ref(false);
const shareOpen = ref(false);
const customerLoading = ref(false);
const customerList = ref([]);
const currentShipment = ref({});
const shareUrl = ref("");

const data = reactive({
  queryParams: {
    pageNum: 1,
    pageSize: 10,
    shipmentNo: "",
    customerName: "",
    pod: "",
    status: ""
  },
  importForm: {
    customerId: "",
    customerName: "",
    orderNo: "",
    pol: "",
    pod: "",
    paymentStatus: "UNPAID",
    paymentAmount: 0,
    preferredType: "",
    cargoList: []
  },
  detail: {},
  statusForm: {
    status: "10",
    actualEtd: "",
    actualEta: "",
    paymentStatus: "UNPAID",
    paymentAmount: 0,
    remark: ""
  },
  bindCustomerForm: {
    shipmentId: "",
    customerId: "",
    customerName: ""
  }
});

const { queryParams, importForm, detail, statusForm, bindCustomerForm } = toRefs(data);

function getList() {
  loading.value = true;
  const params = {
    ...queryParams.value,
    shipmentNo: queryParams.value.shipmentNo || undefined,
    customerName: queryParams.value.customerName || undefined,
    pod: queryParams.value.pod || undefined,
    status: queryParams.value.status || undefined
  };
  listShipment(params)
    .then(response => {
      const tableData = response.data || {};
      shipmentList.value = tableData.rows || [];
      total.value = tableData.total || 0;
    })
    .finally(() => {
      loading.value = false;
    });
}

function handleQuery() {
  queryParams.value.pageNum = 1;
  getList();
}

function resetQuery() {
  queryParams.value = {
    pageNum: 1,
    pageSize: 10,
    shipmentNo: "",
    customerName: "",
    pod: "",
    status: ""
  };
  shipmentTableRef.value?.clearCheckboxRow();
  ids.value = [];
  single.value = true;
  multiple.value = true;
  getList();
}

function handlePageChange({ currentPage, pageSize }) {
  queryParams.value.pageNum = currentPage;
  queryParams.value.pageSize = pageSize;
  getList();
}

function handleSelectionChange() {
  const records = shipmentTableRef.value?.getCheckboxRecords?.() || [];
  ids.value = records.map(item => item.shipmentId);
  single.value = records.length !== 1;
  multiple.value = !records.length;
}

function createCargoRow() {
  return {
    rowKey: `${Date.now()}_${Math.random().toString(16).slice(2)}`,
    cargoName: "",
    sku: "",
    cartons: 0,
    weightKg: 0,
    volumeCbm: 0
  };
}

function resetImport() {
  importForm.value = {
    customerId: "",
    customerName: "",
    orderNo: "",
    pol: "",
    pod: "",
    paymentStatus: "UNPAID",
    paymentAmount: 0,
    preferredType: "",
    cargoList: [createCargoRow()]
  };
}

function handleImport() {
  resetImport();
  importOpen.value = true;
  loadCustomerOptions("");
}

function loadCustomerOptions(keyword = "") {
  customerLoading.value = true;
  customerOptions({ keyword })
    .then(response => {
      customerList.value = response.data || [];
    })
    .finally(() => {
      customerLoading.value = false;
    });
}

function handleCustomerChange(customerId) {
  const customer = customerList.value.find(item => item.customerId === customerId);
  importForm.value.customerName = customer?.customerName || "";
}

function addCargoRow() {
  importForm.value.cargoList.push(createCargoRow());
}

function removeCargoRow(index) {
  importForm.value.cargoList.splice(index, 1);
  if (!importForm.value.cargoList.length) {
    addCargoRow();
  }
}

function submitImport() {
  if (!importForm.value.customerId) {
    proxy.$modal.msgWarning("请选择客户");
    return;
  }
  const cargoList = importForm.value.cargoList.filter(item => item.cargoName);
  if (!cargoList.length) {
    proxy.$modal.msgWarning("请至少填写一条货物明细");
    return;
  }

  importShipment({ ...importForm.value, cargoList }).then(response => {
    detail.value = response.data || {};
    importOpen.value = false;
    detailOpen.value = true;
    proxy.$modal.msgSuccess("出货计划已生成");
    getList();
  });
}

function handleDetail(row) {
  getShipment(row.shipmentId).then(response => {
    detail.value = response.data || {};
    detailOpen.value = true;
  });
}

function handleStatus(row) {
  currentShipment.value = row;
  statusForm.value = {
    status: row.status || "10",
    actualEtd: row.actualEtd || "",
    actualEta: row.actualEta || "",
    paymentStatus: row.paymentStatus || "UNPAID",
    paymentAmount: Number(row.paymentAmount || 0),
    remark: row.remark || ""
  };
  statusOpen.value = true;
}

function handleBindCustomer(row) {
  currentShipment.value = row || shipmentList.value.find(item => item.shipmentId === ids.value[0]) || {};
  bindCustomerForm.value = {
    shipmentId: currentShipment.value.shipmentId,
    customerId: currentShipment.value.customerId || "",
    customerName: currentShipment.value.customerName || ""
  };
  bindCustomerOpen.value = true;
  loadCustomerOptions("");
}

function handleBindCustomerChange(customerId) {
  const customer = customerList.value.find(item => item.customerId === customerId);
  bindCustomerForm.value.customerName = customer?.customerName || "";
}

function submitBindCustomer() {
  if (!bindCustomerForm.value.shipmentId || !bindCustomerForm.value.customerId) {
    proxy.$modal.msgWarning("请选择客户");
    return;
  }

  bindShipmentCustomer(bindCustomerForm.value.shipmentId, {
    customerId: bindCustomerForm.value.customerId,
    customerName: bindCustomerForm.value.customerName
  }).then(() => {
    proxy.$modal.msgSuccess("客户已绑定");
    bindCustomerOpen.value = false;
    getList();
  });
}

function submitStatus() {
  updateShipmentStatus(currentShipment.value.shipmentId, statusForm.value).then(() => {
    proxy.$modal.msgSuccess("状态已更新");
    statusOpen.value = false;
    getList();
  });
}

function handleConfirm(row) {
  confirmShipment(row.shipmentId).then(response => {
    proxy.$modal.msgSuccess(`出货单已生成：${response.data.orderNo}`);
    getList();
  });
}

function handleShare(row) {
  getShipmentShare(row.shipmentId).then(response => {
    shareUrl.value = portalBaseUrl + response.data.shareUrl;
    shareOpen.value = true;
    copyShareUrl(shareUrl.value);
  });
}

async function copyShareUrl(url) {
  if (!url) {
    return;
  }
  try {
    if (navigator.clipboard?.writeText) {
      await navigator.clipboard.writeText(url);
    } else {
      const input = document.createElement("textarea");
      input.value = url;
      input.style.position = "fixed";
      input.style.opacity = "0";
      document.body.appendChild(input);
      input.focus();
      input.select();
      document.execCommand("copy");
      document.body.removeChild(input);
    }
    proxy.$modal.msgSuccess("分享链接已复制");
  } catch (error) {
    proxy.$modal.msgWarning("复制失败，请手动复制链接");
  }
}

function handleDelete(row) {
  const shipmentIds = row?.shipmentId || ids.value;
  proxy.$modal.confirm(`是否确认删除出货计划编号为“${shipmentIds}”的数据项？`)
    .then(() => delShipment(shipmentIds))
    .then(() => {
      getList();
      proxy.$modal.msgSuccess("删除成功");
    })
    .catch(() => { });
}

getList();
</script>

<style scoped>
.shipment-search-form {
  margin-bottom: 8px;
}

.toolbar-spacer {
  height: 8px;
}

.shipment-toolbar {
  margin-bottom: 12px;
}

.shipment-pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.action-links {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
}

.cargo-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin: 12px 0;
  font-weight: 600;
}

.detail-summary {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  margin-bottom: 18px;
}

.detail-item {
  padding: 12px 14px;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  background: #fafcff;
}

.detail-label {
  display: block;
  margin-bottom: 6px;
  color: #909399;
  font-size: 13px;
}

.detail-value {
  color: #303133;
  font-weight: 500;
  word-break: break-all;
}

.shipment-page h4 {
  margin: 18px 0 10px;
  font-size: 15px;
  color: #303133;
}

.share-modal {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.share-tip {
  color: #606266;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

@media (max-width: 900px) {
  .detail-summary {
    grid-template-columns: 1fr;
  }
}
</style>
