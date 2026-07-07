<template>
  <div class="app-container shipment-page">
    <a-form v-show="showSearch" ref="queryRef" :model="queryParams" layout="inline" class="shipment-search-form">
      <a-form-item label="计划编号" name="shipmentNo">
        <a-input v-model:value="queryParams.shipmentNo" placeholder="请输入计划编号" allow-clear style="width: 180px"
          @pressEnter="handleQuery" />
      </a-form-item>

      <a-form-item label="客户名称" name="customerName">
        <a-input v-model:value="queryParams.customerName" placeholder="请输入客户名称" allow-clear style="width: 180px"
          @pressEnter="handleQuery" />
      </a-form-item>

      <a-form-item label="目的港" name="pod">
        <a-input v-model:value="queryParams.pod" placeholder="请输入目的港" allow-clear style="width: 160px"
          @pressEnter="handleQuery" />
      </a-form-item>

      <a-form-item label="状态" name="status">
        <a-select v-model:value="queryParams.status" placeholder="请选择" allow-clear style="width: 190px">
          <a-select-option v-for="item in statusOptions" :key="item.value" :value="item.value">
            {{ item.label }}
          </a-select-option>
        </a-select>
      </a-form-item>

      <a-form-item>
        <a-space>
          <a-button type="primary" @click="handleQuery">
            <template #icon>
              <SearchOutlined />
            </template>
            搜索
          </a-button>
          <a-button @click="resetQuery">
            <template #icon>
              <ReloadOutlined />
            </template>
            重置
          </a-button>
        </a-space>
      </a-form-item>
    </a-form>

    <a-row :gutter="10" class="mb8 shipment-toolbar">
      <a-col>
        <a-button type="primary" ghost @click="handleImport" v-hasPermi="['freight:shipment:import']">
          <template #icon>
            <UploadOutlined />
          </template>
          导入清单
        </a-button>
      </a-col>

      <a-col>
        <a-button danger ghost :disabled="multiple" @click="handleDelete" v-hasPermi="['freight:shipment:remove']">
          <template #icon>
            <DeleteOutlined />
          </template>
          删除
        </a-button>
      </a-col>

      <a-col>
        <a-button ghost :disabled="ids.length !== 1" @click="handleBindCustomer()"
          v-hasPermi="['freight:shipment:edit']">
          <template #icon>
            <UserOutlined />
          </template>
          绑定客户
        </a-button>
      </a-col>

      <a-col flex="auto" class="toolbar-right">
        <right-toolbar v-model:showSearch="showSearch" @queryTable="getList" />
      </a-col>
    </a-row>

    <a-table row-key="shipmentId" :loading="loading" :data-source="shipmentList" :pagination="false"
      :row-selection="rowSelection" bordered size="middle" :scroll="{ x: 1500 }">
      <a-table-column title="计划编号" data-index="shipmentNo" :width="170" align="center" />
      <a-table-column title="客户" data-index="customerName" :width="150" align="center" ellipsis />
      <a-table-column title="客户订单号" data-index="orderNo" :width="140" align="center" />

      <a-table-column title="航线" :width="170" align="center">
        <template #default="{ record }">
          {{ record.pol || '-' }} → {{ record.pod || '-' }}
        </template>
      </a-table-column>

      <a-table-column title="货量" :width="170" align="center">
        <template #default="{ record }">
          {{ record.totalCartons }}箱 / {{ record.totalVolume }}CBM
        </template>
      </a-table-column>

      <a-table-column title="重量" :width="110" align="center">
        <template #default="{ record }">
          {{ record.totalWeight }}KG
        </template>
      </a-table-column>

      <a-table-column title="状态" data-index="status" :width="150" align="center">
        <template #default="{ record }">
          <a-tag :color="statusTag(record.status)">
            {{ statusLabel(record.status) }}
          </a-tag>
        </template>
      </a-table-column>

      <a-table-column title="创建时间" data-index="createTime" :width="160" align="center">
        <template #default="{ record }">
          {{ parseTime(record.createTime) }}
        </template>
      </a-table-column>

      <a-table-column title="操作" :width="310" align="center" fixed="right">
        <template #default="{ record }">
          <a-space>
            <a-button type="link" size="small" @click="handleDetail(record)">
              <template #icon>
                <EyeOutlined />
              </template>
              详情
            </a-button>
            <a-button type="link" size="small" @click="handleStatus(record)" v-hasPermi="['freight:shipment:edit']">
              <template #icon>
                <EditOutlined />
              </template>
              状态
            </a-button>
            <a-button type="link" size="small" @click="handleConfirm(record)" v-hasPermi="['freight:shipment:confirm']">
              <template #icon>
                <FileDoneOutlined />
              </template>
              出货单
            </a-button>
            <a-button type="link" size="small" @click="handleShare(record)" v-hasPermi="['freight:shipment:share']">
              <template #icon>
                <ShareAltOutlined />
              </template>
              分享
            </a-button>
          </a-space>
        </template>
      </a-table-column>
    </a-table>

    <div v-show="total > 0" class="shipment-pagination">
      <a-pagination v-model:current="queryParams.pageNum" v-model:page-size="queryParams.pageSize" :total="total"
        show-size-changer show-quick-jumper :show-total="total => `共 ${total} 条`" @change="handlePageChange"
        @showSizeChange="handlePageChange" />
    </div>

    <a-modal v-model:open="importOpen" title="导入出货清单并生成智能计划" width="980px" :mask-closable="false" destroy-on-close
      @ok="submitImport">
      <a-form ref="importRef" :model="importForm" :rules="importRules" :label-col="{ style: { width: '92px' } }">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="客户" name="customerId">
              <a-select v-model:value="importForm.customerId" show-search allow-clear placeholder="搜索客户"
                :filter-option="false" :loading="customerLoading" style="width: 100%" @search="loadCustomerOptions"
                @change="handleCustomerChange">
                <a-select-option v-for="item in customerList" :key="item.customerId" :value="item.customerId">
                  {{ item.customerName + ' / ' + (item.companyName || '-') }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>

          <a-col :span="12">
            <a-form-item label="客户单号" name="orderNo">
              <a-input v-model:value="importForm.orderNo" placeholder="客户订单号/参考号" :maxlength="64" />
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="16">
          <a-col :span="8">
            <a-form-item label="起运港" name="pol">
              <a-input v-model:value="importForm.pol" placeholder="如 SHANGHAI" />
            </a-form-item>
          </a-col>

          <a-col :span="8">
            <a-form-item label="目的港" name="pod">
              <a-input v-model:value="importForm.pod" placeholder="如 LOS ANGELES" />
            </a-form-item>
          </a-col>

          <a-col :span="8">
            <a-form-item label="偏好柜型" name="preferredType">
              <a-select v-model:value="importForm.preferredType" allow-clear placeholder="系统自动">
                <a-select-option value="20GP">20GP</a-select-option>
                <a-select-option value="40GP">40GP</a-select-option>
                <a-select-option value="40HQ">40HQ</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>

        <div class="cargo-toolbar">
          <span>货物明细</span>
          <a-button type="primary" ghost size="small" @click="addCargoRow">
            <template #icon>
              <PlusOutlined />
            </template>
            新增一行
          </a-button>
        </div>

        <a-table :data-source="importForm.cargoList" :pagination="false" row-key="rowKey" bordered size="small">
          <a-table-column title="货名" :width="160">
            <template #default="{ record }">
              <a-input v-model:value="record.cargoName" placeholder="货物名称" />
            </template>
          </a-table-column>

          <a-table-column title="SKU/唛头" :width="120">
            <template #default="{ record }">
              <a-input v-model:value="record.sku" placeholder="可选" />
            </template>
          </a-table-column>

          <a-table-column title="箱数" :width="100">
            <template #default="{ record }">
              <a-input-number v-model:value="record.cartons" :min="0" :controls="false" style="width: 78px" />
            </template>
          </a-table-column>

          <a-table-column title="重量KG" :width="120">
            <template #default="{ record }">
              <a-input-number v-model:value="record.weightKg" :min="0" :precision="2" :controls="false"
                style="width: 98px" />
            </template>
          </a-table-column>

          <a-table-column title="体积CBM" :width="120">
            <template #default="{ record }">
              <a-input-number v-model:value="record.volumeCbm" :min="0" :precision="2" :controls="false"
                style="width: 98px" />
            </template>
          </a-table-column>

          <a-table-column title="操作" :width="80" align="center">
            <template #default="{ index }">
              <a-button type="link" danger size="small" @click="removeCargoRow(index)">
                <template #icon>
                  <DeleteOutlined />
                </template>
                删除
              </a-button>
            </template>
          </a-table-column>
        </a-table>
      </a-form>

      <template #footer>
        <a-space>
          <a-button @click="importOpen = false">取消</a-button>
          <a-button type="primary" @click="submitImport">生成计划</a-button>
        </a-space>
      </template>
    </a-modal>

    <a-modal v-model:open="detailOpen" title="出货计划详情" width="980px" :footer="null" destroy-on-close>
      <template v-if="detail.plan">
        <a-descriptions :column="3" bordered size="small">
          <a-descriptions-item label="计划编号">{{ detail.plan.shipmentNo }}</a-descriptions-item>
          <a-descriptions-item label="客户">{{ detail.plan.customerName }}</a-descriptions-item>
          <a-descriptions-item label="状态">{{ statusLabel(detail.plan.status) }}</a-descriptions-item>
          <a-descriptions-item label="航线">{{ detail.plan.pol }} → {{ detail.plan.pod }}</a-descriptions-item>
          <a-descriptions-item label="货量">
            {{ detail.plan.totalCartons }}箱 / {{ detail.plan.totalVolume }}CBM
          </a-descriptions-item>
          <a-descriptions-item label="重量">{{ detail.plan.totalWeight }}KG</a-descriptions-item>
          <a-descriptions-item label="出货单">{{ detail.order?.orderNo || '未生成' }}</a-descriptions-item>
          <a-descriptions-item label="计划ETD">{{ detail.plan.plannedEtd || '-' }}</a-descriptions-item>
          <a-descriptions-item label="计划ETA">{{ detail.plan.plannedEta || '-' }}</a-descriptions-item>
        </a-descriptions>

        <h4>智能货柜建议</h4>
        <a-table :data-source="detail.containers" :pagination="false" row-key="containerType" bordered size="small">
          <a-table-column title="柜型" data-index="containerType" align="center" />
          <a-table-column title="柜量" data-index="quantity" align="center" />
          <a-table-column title="装载率" align="center">
            <template #default="{ record }">
              {{ record.loadRate }}%
            </template>
          </a-table-column>
          <a-table-column title="说明" data-index="remark" />
        </a-table>

        <h4>货物明细</h4>
        <a-table :data-source="detail.cargoList" :pagination="false" row-key="cargoId" bordered size="small">
          <a-table-column title="货名" data-index="cargoName" />
          <a-table-column title="SKU/唛头" data-index="sku" />
          <a-table-column title="箱数" data-index="cartons" align="center" />
          <a-table-column title="重量KG" data-index="weightKg" align="center" />
          <a-table-column title="体积CBM" data-index="volumeCbm" align="center" />
        </a-table>
      </template>
    </a-modal>

    <a-modal v-model:open="statusOpen" title="维护客户可见状态" width="520px" :mask-closable="false" destroy-on-close>
      <a-form :model="statusForm" :label-col="{ style: { width: '90px' } }">
        <a-form-item label="当前状态">
          <a-select v-model:value="statusForm.status" style="width: 100%">
            <a-select-option v-for="item in statusOptions" :key="item.value" :value="item.value">
              {{ item.label }}
            </a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item label="实际ETD">
          <a-date-picker v-model:value="statusForm.actualEtd" value-format="YYYY-MM-DD" placeholder="选择日期"
            style="width: 100%" />
        </a-form-item>

        <a-form-item label="实际ETA">
          <a-date-picker v-model:value="statusForm.actualEta" value-format="YYYY-MM-DD" placeholder="选择日期"
            style="width: 100%" />
        </a-form-item>

        <a-form-item label="备注">
          <a-textarea v-model:value="statusForm.remark" :rows="3" placeholder="给客户看的简短说明" />
        </a-form-item>
      </a-form>

      <template #footer>
        <a-space>
          <a-button @click="statusOpen = false">取消</a-button>
          <a-button type="primary" @click="submitStatus">保存</a-button>
        </a-space>
      </template>
    </a-modal>

    <a-modal v-model:open="bindCustomerOpen" title="绑定客户" width="520px" :mask-closable="false" destroy-on-close>
      <a-form :model="bindCustomerForm" :label-col="{ style: { width: '90px' } }">
        <a-form-item label="客户">
          <a-select v-model:value="bindCustomerForm.customerId" show-search allow-clear placeholder="搜索客户"
            :filter-option="false" :loading="customerLoading" style="width: 100%" @search="loadCustomerOptions"
            @change="handleBindCustomerChange">
            <a-select-option v-for="item in customerList" :key="item.customerId" :value="item.customerId">
              {{ item.customerName + ' / ' + (item.companyName || '-') }}
            </a-select-option>
          </a-select>
        </a-form-item>
      </a-form>

      <template #footer>
        <a-space>
          <a-button @click="bindCustomerOpen = false">取消</a-button>
          <a-button type="primary" @click="submitBindCustomer">保存</a-button>
        </a-space>
      </template>
    </a-modal>
  </div>
</template>

<script setup name="FreightShipment">
import { computed, reactive, ref, toRefs } from "vue";
import { Modal, message } from "ant-design-vue";
import {
  DeleteOutlined,
  EditOutlined,
  EyeOutlined,
  FileDoneOutlined,
  PlusOutlined,
  ReloadOutlined,
  SearchOutlined,
  ShareAltOutlined,
  UploadOutlined,
  UserOutlined
} from "@ant-design/icons-vue";
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

const portalBaseUrl = (import.meta.env.VITE_PORTAL_BASE_URL || window.location.origin).replace(/\/$/, "");

const queryRef = ref();
const importRef = ref();

const statusOptions = [
  { value: "10", label: "计划已创建", tag: "processing" },
  { value: "20", label: "出货计划已确认", tag: "processing" },
  { value: "30", label: "等待客户发货", tag: "default" },
  { value: "40", label: "已提货/已送仓", tag: "warning" },
  { value: "50", label: "仓库已收货", tag: "warning" },
  { value: "60", label: "已入仓/码头进仓", tag: "warning" },
  { value: "70", label: "订舱处理中", tag: "warning" },
  { value: "80", label: "舱位已确认", tag: "processing" },
  { value: "90", label: "报关资料已收齐", tag: "processing" },
  { value: "100", label: "报关已放行", tag: "success" },
  { value: "110", label: "已装柜", tag: "success" },
  { value: "120", label: "已进港/码头放行", tag: "success" },
  { value: "130", label: "船舶已开航", tag: "success" },
  { value: "140", label: "目的港已到港", tag: "success" },
  { value: "150", label: "目的港清关中", tag: "warning" },
  { value: "160", label: "目的港已清关", tag: "success" },
  { value: "170", label: "已派送/已签收", tag: "success" },
  { value: "900", label: "异常处理中", tag: "error" }
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
const bindCustomerOpen = ref(false);
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
  },
  bindCustomerForm: {
    shipmentId: undefined,
    customerId: undefined,
    customerName: undefined
  }
});

const { queryParams, importForm, importRules, detail, statusForm, bindCustomerForm } = toRefs(data);

const rowSelection = computed(() => ({
  selectedRowKeys: ids.value,
  onChange: handleSelectionChange
}));

function getList() {
  loading.value = true;
  listShipment(queryParams.value)
    .then(response => {
      const data = response.data || {};
      shipmentList.value = data.rows || [];
      total.value = data.total || 0;
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
  queryRef.value?.resetFields();
  handleQuery();
}

function handlePageChange(page, pageSize) {
  queryParams.value.pageNum = page;
  queryParams.value.pageSize = pageSize;
  getList();
}

function handleSelectionChange(selectedRowKeys) {
  ids.value = selectedRowKeys;
  multiple.value = !selectedRowKeys.length;
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
  if (!importForm.value.cargoList.length) addCargoRow();
}

function submitImport() {
  importRef.value
    ?.validate()
    .then(() => {
      const cargoList = importForm.value.cargoList.filter(item => item.cargoName);
      if (!cargoList.length) {
        message.warning("请至少填写一条货物明细");
        return;
      }

      importShipment({ ...importForm.value, cargoList }).then(response => {
        detail.value = response.data || {};
        importOpen.value = false;
        detailOpen.value = true;
        message.success("出货计划已生成");
        getList();
      });
    })
    .catch(() => { });
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

function handleBindCustomer(row) {
  currentShipment.value = row || shipmentList.value.find(item => item.shipmentId === ids.value[0]) || {};
  bindCustomerForm.value = {
    shipmentId: currentShipment.value.shipmentId,
    customerId: currentShipment.value.customerId || undefined,
    customerName: currentShipment.value.customerName || undefined
  };
  loadCustomerOptions("");
  bindCustomerOpen.value = true;
}

function handleBindCustomerChange(customerId) {
  const customer = customerList.value.find(item => item.customerId === customerId);
  bindCustomerForm.value.customerName = customer?.customerName || "";
}

function submitBindCustomer() {
  if (!bindCustomerForm.value.shipmentId || !bindCustomerForm.value.customerId) {
    message.warning("请选择客户");
    return;
  }

  bindShipmentCustomer(bindCustomerForm.value.shipmentId, {
    customerId: bindCustomerForm.value.customerId,
    customerName: bindCustomerForm.value.customerName
  }).then(() => {
    message.success("客户已绑定");
    bindCustomerOpen.value = false;
    getList();
  });
}

function submitStatus() {
  updateShipmentStatus(currentShipment.value.shipmentId, statusForm.value).then(() => {
    message.success("状态已更新");
    statusOpen.value = false;
    getList();
  });
}

function handleConfirm(row) {
  confirmShipment(row.shipmentId).then(response => {
    message.success("出货单已生成：" + response.data.orderNo);
    getList();
  });
}

function handleShare(row) {
  getShipmentShare(row.shipmentId).then(response => {
    const url = portalBaseUrl + response.data.shareUrl;
    navigator.clipboard?.writeText(url);
    Modal.info({
      title: "客户免登录分享链接已复制",
      content: url
    });
  });
}

function handleDelete(row) {
  const shipmentIds = row?.shipmentId || ids.value;
  Modal.confirm({
    title: "确认删除",
    content: `是否确认删除出货计划编号为 "${shipmentIds}" 的数据项？`,
    okText: "确定",
    cancelText: "取消",
    okType: "danger",
    onOk() {
      return delShipment(shipmentIds).then(() => {
        getList();
        message.success("删除成功");
      });
    }
  });
}

function statusLabel(status) {
  return statusOptions.find(item => item.value === status)?.label || status;
}

function statusTag(status) {
  return statusOptions.find(item => item.value === status)?.tag || "default";
}

getList();
</script>

<style scoped>
.shipment-search-form {
  row-gap: 12px;
  margin-bottom: 16px;
}

.shipment-toolbar {
  margin-bottom: 8px;
}

.toolbar-right {
  display: flex;
  justify-content: flex-end;
}

.shipment-pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

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
