<template>
  <div class="app-container shipment-page">
    <el-form :model="queryParams" ref="queryRef" :inline="true" v-show="showSearch" label-width="78px">
      <el-form-item label="计划编号" prop="shipmentNo">
        <el-input v-model="queryParams.shipmentNo" placeholder="请输入计划编号" clearable size="small" style="width: 180px" @keyup.enter="handleQuery" />
      </el-form-item>
      <el-form-item label="客户名称" prop="customerName">
        <el-input v-model="queryParams.customerName" placeholder="请输入客户名称" clearable size="small" style="width: 180px" @keyup.enter="handleQuery" />
      </el-form-item>
      <el-form-item label="目的港" prop="pod">
        <el-input v-model="queryParams.pod" placeholder="请输入目的港" clearable size="small" style="width: 160px" @keyup.enter="handleQuery" />
      </el-form-item>
      <el-form-item label="状态" prop="status">
        <el-select v-model="queryParams.status" placeholder="请选择" clearable size="small" style="width: 190px">
          <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" icon="Search" size="mini" @click="handleQuery">搜索</el-button>
        <el-button icon="Refresh" size="mini" @click="resetQuery">重置</el-button>
      </el-form-item>
    </el-form>

    <el-row :gutter="10" class="mb8">
      <el-col :span="1.5">
        <el-button type="primary" plain icon="Upload" size="mini" @click="handleImport" v-hasPermi="['freight:shipment:import']">导入清单</el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button type="danger" plain icon="Delete" size="mini" :disabled="multiple" @click="handleDelete" v-hasPermi="['freight:shipment:remove']">删除</el-button>
      </el-col>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList"></right-toolbar>
    </el-row>

    <el-table v-loading="loading" :data="shipmentList" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="50" align="center" />
      <el-table-column label="计划编号" prop="shipmentNo" min-width="170" align="center" />
      <el-table-column label="客户" prop="customerName" min-width="150" align="center" :show-overflow-tooltip="true" />
      <el-table-column label="客户订单号" prop="orderNo" min-width="140" align="center" />
      <el-table-column label="航线" min-width="170" align="center">
        <template #default="scope">{{ scope.row.pol || '-' }} → {{ scope.row.pod || '-' }}</template>
      </el-table-column>
      <el-table-column label="货量" width="170" align="center">
        <template #default="scope">{{ scope.row.totalCartons }}箱 / {{ scope.row.totalVolume }}CBM</template>
      </el-table-column>
      <el-table-column label="重量" width="110" align="center">
        <template #default="scope">{{ scope.row.totalWeight }}KG</template>
      </el-table-column>
      <el-table-column label="状态" prop="status" width="150" align="center">
        <template #default="scope">
          <el-tag :type="statusTag(scope.row.status)">{{ statusLabel(scope.row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" prop="createTime" width="160" align="center">
        <template #default="scope">{{ parseTime(scope.row.createTime) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="310" align="center" fixed="right">
        <template #default="scope">
          <el-button size="mini" type="text" icon="View" @click="handleDetail(scope.row)">详情</el-button>
          <el-button size="mini" type="text" icon="Edit" @click="handleStatus(scope.row)" v-hasPermi="['freight:shipment:edit']">状态</el-button>
          <el-button size="mini" type="text" icon="DocumentChecked" @click="handleConfirm(scope.row)" v-hasPermi="['freight:shipment:confirm']">出货单</el-button>
          <el-button size="mini" type="text" icon="Share" @click="handleShare(scope.row)" v-hasPermi="['freight:shipment:share']">分享</el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total > 0" :total="total" v-model:page="queryParams.pageNum" v-model:limit="queryParams.pageSize" @pagination="getList" />

    <el-dialog title="导入出货清单并生成智能计划" v-model="importOpen" width="980px" append-to-body>
      <el-form ref="importRef" :model="importForm" :rules="importRules" label-width="92px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="客户" prop="customerId">
              <el-select v-model="importForm.customerId" filterable remote clearable reserve-keyword placeholder="搜索客户" :remote-method="loadCustomerOptions" :loading="customerLoading" style="width: 100%" @change="handleCustomerChange">
                <el-option v-for="item in customerList" :key="item.customerId" :label="item.customerName + ' / ' + (item.companyName || '-')" :value="item.customerId" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="客户单号" prop="orderNo">
              <el-input v-model="importForm.orderNo" placeholder="客户订单号/参考号" maxlength="64" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="8">
            <el-form-item label="起运港" prop="pol">
              <el-input v-model="importForm.pol" placeholder="如 SHANGHAI" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="目的港" prop="pod">
              <el-input v-model="importForm.pod" placeholder="如 LOS ANGELES" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="偏好柜型" prop="preferredType">
              <el-select v-model="importForm.preferredType" clearable placeholder="系统自动">
                <el-option label="20GP" value="20GP" />
                <el-option label="40GP" value="40GP" />
                <el-option label="40HQ" value="40HQ" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <div class="cargo-toolbar">
          <span>货物明细</span>
          <el-button type="primary" plain icon="Plus" size="mini" @click="addCargoRow">新增一行</el-button>
        </div>
        <el-table :data="importForm.cargoList" border size="small">
          <el-table-column label="货名" min-width="160">
            <template #default="scope"><el-input v-model="scope.row.cargoName" placeholder="货物名称" /></template>
          </el-table-column>
          <el-table-column label="SKU/唛头" min-width="120">
            <template #default="scope"><el-input v-model="scope.row.sku" placeholder="可选" /></template>
          </el-table-column>
          <el-table-column label="箱数" width="100">
            <template #default="scope"><el-input-number v-model="scope.row.cartons" :min="0" :controls="false" style="width: 78px" /></template>
          </el-table-column>
          <el-table-column label="重量KG" width="120">
            <template #default="scope"><el-input-number v-model="scope.row.weightKg" :min="0" :precision="2" :controls="false" style="width: 98px" /></template>
          </el-table-column>
          <el-table-column label="体积CBM" width="120">
            <template #default="scope"><el-input-number v-model="scope.row.volumeCbm" :min="0" :precision="2" :controls="false" style="width: 98px" /></template>
          </el-table-column>
          <el-table-column label="操作" width="80" align="center">
            <template #default="scope"><el-button type="text" icon="Delete" @click="removeCargoRow(scope.$index)">删除</el-button></template>
          </el-table-column>
        </el-table>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button type="primary" @click="submitImport">生成计划</el-button>
          <el-button @click="importOpen = false">取消</el-button>
        </div>
      </template>
    </el-dialog>

    <el-dialog title="出货计划详情" v-model="detailOpen" width="980px" append-to-body>
      <template v-if="detail.plan">
        <el-descriptions :column="3" border>
          <el-descriptions-item label="计划编号">{{ detail.plan.shipmentNo }}</el-descriptions-item>
          <el-descriptions-item label="客户">{{ detail.plan.customerName }}</el-descriptions-item>
          <el-descriptions-item label="状态">{{ statusLabel(detail.plan.status) }}</el-descriptions-item>
          <el-descriptions-item label="航线">{{ detail.plan.pol }} → {{ detail.plan.pod }}</el-descriptions-item>
          <el-descriptions-item label="货量">{{ detail.plan.totalCartons }}箱 / {{ detail.plan.totalVolume }}CBM</el-descriptions-item>
          <el-descriptions-item label="重量">{{ detail.plan.totalWeight }}KG</el-descriptions-item>
          <el-descriptions-item label="出货单">{{ detail.order?.orderNo || '未生成' }}</el-descriptions-item>
          <el-descriptions-item label="计划ETD">{{ detail.plan.plannedEtd || '-' }}</el-descriptions-item>
          <el-descriptions-item label="计划ETA">{{ detail.plan.plannedEta || '-' }}</el-descriptions-item>
        </el-descriptions>

        <h4>智能货柜建议</h4>
        <el-table :data="detail.containers" border size="small">
          <el-table-column label="柜型" prop="containerType" align="center" />
          <el-table-column label="柜量" prop="quantity" align="center" />
          <el-table-column label="装载率" align="center"><template #default="scope">{{ scope.row.loadRate }}%</template></el-table-column>
          <el-table-column label="说明" prop="remark" />
        </el-table>

        <h4>货物明细</h4>
        <el-table :data="detail.cargoList" border size="small">
          <el-table-column label="货名" prop="cargoName" />
          <el-table-column label="SKU/唛头" prop="sku" />
          <el-table-column label="箱数" prop="cartons" align="center" />
          <el-table-column label="重量KG" prop="weightKg" align="center" />
          <el-table-column label="体积CBM" prop="volumeCbm" align="center" />
        </el-table>
      </template>
    </el-dialog>

    <el-dialog title="维护客户可见状态" v-model="statusOpen" width="520px" append-to-body>
      <el-form :model="statusForm" label-width="90px">
        <el-form-item label="当前状态">
          <el-select v-model="statusForm.status" style="width: 100%">
            <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="实际ETD">
          <el-date-picker v-model="statusForm.actualEtd" value-format="YYYY-MM-DD" type="date" placeholder="选择日期" style="width: 100%" />
        </el-form-item>
        <el-form-item label="实际ETA">
          <el-date-picker v-model="statusForm.actualEta" value-format="YYYY-MM-DD" type="date" placeholder="选择日期" style="width: 100%" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="statusForm.remark" type="textarea" :rows="3" placeholder="给客户看的简短说明" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button type="primary" @click="submitStatus">保存</el-button>
          <el-button @click="statusOpen = false">取消</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup name="FreightShipment">
import {
  listShipment,
  importShipment,
  getShipment,
  updateShipmentStatus,
  confirmShipment,
  getShipmentShare,
  delShipment
} from "@/api/freight/shipment";
import { customerOptions } from "@/api/customer/customer";

const { proxy } = getCurrentInstance();
const portalBaseUrl = (import.meta.env.VITE_PORTAL_BASE_URL || window.location.origin).replace(/\/$/, "");

const statusOptions = [
  { value: "10", label: "计划已创建", tag: "" },
  { value: "20", label: "出货计划已确认", tag: "" },
  { value: "30", label: "等待客户发货", tag: "info" },
  { value: "40", label: "已提货/已送仓", tag: "warning" },
  { value: "50", label: "仓库已收货", tag: "warning" },
  { value: "60", label: "已入仓/码头进仓", tag: "warning" },
  { value: "70", label: "订舱处理中", tag: "warning" },
  { value: "80", label: "舱位已确认", tag: "" },
  { value: "90", label: "报关资料已收齐", tag: "" },
  { value: "100", label: "报关已放行", tag: "success" },
  { value: "110", label: "已装柜", tag: "success" },
  { value: "120", label: "已进港/码头放行", tag: "success" },
  { value: "130", label: "船舶已开航", tag: "success" },
  { value: "140", label: "目的港已到港", tag: "success" },
  { value: "150", label: "目的港清关中", tag: "warning" },
  { value: "160", label: "目的港已清关", tag: "success" },
  { value: "170", label: "已派送/已签收", tag: "success" },
  { value: "900", label: "异常处理中", tag: "danger" }
];

const shipmentList = ref([]);
const loading = ref(true);
const showSearch = ref(true);
const ids = ref([]);
const multiple = ref(true);
const total = ref(0);
const importOpen = ref(false);
const detailOpen = ref(false);
const statusOpen = ref(false);
const customerLoading = ref(false);
const customerList = ref([]);
const currentShipment = ref({});

const data = reactive({
  queryParams: {
    pageNum: 1,
    pageSize: 10,
    shipmentNo: undefined,
    customerName: undefined,
    pod: undefined,
    status: undefined
  },
  importForm: {
    customerId: undefined,
    customerName: undefined,
    orderNo: undefined,
    pol: undefined,
    pod: undefined,
    preferredType: undefined,
    cargoList: []
  },
  importRules: {
    customerId: [{ required: true, message: "请选择客户", trigger: "change" }]
  },
  detail: {},
  statusForm: {
    status: "10",
    actualEtd: undefined,
    actualEta: undefined,
    remark: undefined
  }
});

const { queryParams, importForm, importRules, detail, statusForm } = toRefs(data);

function getList() {
  loading.value = true;
  listShipment(queryParams.value).then(response => {
    const data = response.data || {};
    shipmentList.value = data.rows || [];
    total.value = data.total || 0;
    loading.value = false;
  });
}

function handleQuery() {
  queryParams.value.pageNum = 1;
  getList();
}

function resetQuery() {
  proxy.resetForm("queryRef");
  handleQuery();
}

function handleSelectionChange(selection) {
  ids.value = selection.map(item => item.shipmentId);
  multiple.value = !selection.length;
}

function resetImport() {
  importForm.value = {
    customerId: undefined,
    customerName: undefined,
    orderNo: undefined,
    pol: undefined,
    pod: undefined,
    preferredType: undefined,
    cargoList: []
  };
  addCargoRow();
}

function handleImport() {
  resetImport();
  importOpen.value = true;
}

function loadCustomerOptions(keyword) {
  customerLoading.value = true;
  customerOptions({ keyword }).then(response => {
    customerList.value = response.data || [];
    customerLoading.value = false;
  });
}

function handleCustomerChange(customerId) {
  const customer = customerList.value.find(item => item.customerId === customerId);
  importForm.value.customerName = customer?.customerName || "";
}

function addCargoRow() {
  importForm.value.cargoList.push({
    cargoName: "",
    sku: "",
    cartons: 0,
    weightKg: 0,
    volumeCbm: 0
  });
}

function removeCargoRow(index) {
  importForm.value.cargoList.splice(index, 1);
  if (!importForm.value.cargoList.length) addCargoRow();
}

function submitImport() {
  proxy.$refs["importRef"].validate(valid => {
    if (!valid) return;
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
    actualEtd: row.actualEtd,
    actualEta: row.actualEta,
    remark: row.remark
  };
  statusOpen.value = true;
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
    proxy.$modal.msgSuccess("出货单已生成：" + response.data.orderNo);
    getList();
  });
}

function handleShare(row) {
  getShipmentShare(row.shipmentId).then(response => {
    const url = portalBaseUrl + response.data.shareUrl;
    navigator.clipboard?.writeText(url);
    proxy.$modal.alert("客户免登录分享链接已复制：\n" + url);
  });
}

function handleDelete(row) {
  const shipmentIds = row?.shipmentId || ids.value;
  proxy.$modal.confirm('是否确认删除出货计划编号为 "' + shipmentIds + '" 的数据项？').then(function() {
    return delShipment(shipmentIds);
  }).then(() => {
    getList();
    proxy.$modal.msgSuccess("删除成功");
  }).catch(() => {});
}

function statusLabel(status) {
  return statusOptions.find(item => item.value === status)?.label || status;
}

function statusTag(status) {
  return statusOptions.find(item => item.value === status)?.tag || "";
}

getList();
</script>

<style scoped>
.shipment-page h4 {
  margin: 18px 0 10px;
  font-size: 15px;
  color: #303133;
}

.cargo-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin: 8px 0 10px;
  font-weight: 600;
}
</style>
