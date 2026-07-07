<template>
  <div class="app-container">
    <a-form v-show="showSearch" ref="queryRef" :model="queryParams" layout="inline" class="search-form">
      <a-form-item label="系统模块" name="title">
        <a-input v-model:value="queryParams.title" allow-clear placeholder="请输入系统模块" style="width: 240px" @pressEnter="handleQuery" />
      </a-form-item>
      <a-form-item label="操作人员" name="operName">
        <a-input v-model:value="queryParams.operName" allow-clear placeholder="请输入操作人员" style="width: 240px" @pressEnter="handleQuery" />
      </a-form-item>
      <a-form-item label="类型" name="businessType">
        <a-select v-model:value="queryParams.businessType" allow-clear placeholder="操作类型" style="width: 240px" :options="operTypeOptions" />
      </a-form-item>
      <a-form-item label="状态" name="status">
        <a-select v-model:value="queryParams.status" allow-clear placeholder="操作状态" style="width: 240px" :options="statusOptions" />
      </a-form-item>
      <a-form-item label="操作时间">
        <a-range-picker v-model:value="dateRange" value-format="YYYY-MM-DD" />
      </a-form-item>
      <a-form-item>
        <a-space>
          <a-button type="primary" @click="handleQuery">搜索</a-button>
          <a-button @click="resetQuery">重置</a-button>
        </a-space>
      </a-form-item>
    </a-form>

    <div class="toolbar-row mb8">
      <a-space>
        <a-button danger :disabled="multiple" @click="handleDelete()" v-hasPermi="['system:operlog:remove']">删除</a-button>
        <a-button danger @click="handleClean" v-hasPermi="['system:operlog:remove']">清空</a-button>
        <a-button @click="handleExport" v-hasPermi="['system:operlog:export']">导出</a-button>
      </a-space>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList" />
    </div>

    <a-table
      :loading="loading"
      :data-source="operlogList"
      :columns="columns"
      :pagination="false"
      :row-selection="rowSelection"
      :scroll="{ x: 1500 }"
      row-key="operId"
      @change="handleTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'businessType'">
          <dict-tag :options="sys_oper_type" :value="record.businessType" />
        </template>
        <template v-else-if="column.key === 'status'">
          <dict-tag :options="sys_common_status" :value="record.status" />
        </template>
        <template v-else-if="column.key === 'operTime'">
          {{ parseTime(record.operTime) }}
        </template>
        <template v-else-if="column.key === 'action'">
          <a-button type="link" @click="handleView(record)" v-hasPermi="['system:operlog:query']">详情</a-button>
        </template>
      </template>
    </a-table>

    <pagination
      v-show="total > 0"
      v-model:page="queryParams.pageNum"
      v-model:limit="queryParams.pageSize"
      :total="total"
      @pagination="getList"
    />

    <a-modal v-model:open="open" title="操作日志详情" width="700px" :footer="null" destroy-on-close>
      <a-descriptions :column="2" bordered size="small">
        <a-descriptions-item label="操作模块">{{ form.title }} / {{ typeFormat(form) }}</a-descriptions-item>
        <a-descriptions-item label="登录信息">{{ form.operName }} / {{ form.operIp }} / {{ form.operLocation }}</a-descriptions-item>
        <a-descriptions-item label="请求地址">{{ form.operUrl }}</a-descriptions-item>
        <a-descriptions-item label="请求方式">{{ form.requestMethod }}</a-descriptions-item>
        <a-descriptions-item label="操作方法" :span="2">{{ form.method }}</a-descriptions-item>
        <a-descriptions-item label="请求参数" :span="2">{{ form.operParam }}</a-descriptions-item>
        <a-descriptions-item label="返回参数" :span="2">{{ form.jsonResult }}</a-descriptions-item>
        <a-descriptions-item label="操作状态">{{ form.status === 0 ? "正常" : "失败" }}</a-descriptions-item>
        <a-descriptions-item label="操作时间">{{ parseTime(form.operTime) }}</a-descriptions-item>
        <a-descriptions-item v-if="form.status === 1" label="异常信息" :span="2">{{ form.errorMsg }}</a-descriptions-item>
      </a-descriptions>
      <div class="modal-footer">
        <a-button @click="open = false">关闭</a-button>
      </div>
    </a-modal>
  </div>
</template>

<script setup name="Operlog">
import { list, delOperlog, cleanOperlog } from "@/api/monitor/operlog";

const { proxy } = getCurrentInstance();
const { sys_oper_type, sys_common_status } = proxy.useDict("sys_oper_type", "sys_common_status");

const columns = [
  { title: "日志编号", dataIndex: "operId", key: "operId", align: "center" },
  { title: "系统模块", dataIndex: "title", key: "title", align: "center" },
  { title: "操作类型", dataIndex: "businessType", key: "businessType", align: "center" },
  { title: "请求方式", dataIndex: "requestMethod", key: "requestMethod", align: "center" },
  { title: "操作人员", dataIndex: "operName", key: "operName", ellipsis: true, sorter: true, width: 100, align: "center" },
  { title: "主机", dataIndex: "operIp", key: "operIp", ellipsis: true, width: 130, align: "center" },
  { title: "操作状态", dataIndex: "status", key: "status", align: "center" },
  { title: "操作日期", dataIndex: "operTime", key: "operTime", sorter: true, width: 180, align: "center" },
  { title: "操作", key: "action", align: "center" }
];

const operlogList = ref([]);
const open = ref(false);
const loading = ref(true);
const showSearch = ref(true);
const ids = ref([]);
const selectedRowKeys = ref([]);
const multiple = ref(true);
const total = ref(0);
const dateRange = ref([]);

const rowSelection = computed(() => ({
  selectedRowKeys: selectedRowKeys.value,
  onChange: (keys, rows) => handleSelectionChange(rows, keys)
}));

const operTypeOptions = computed(() => (sys_oper_type.value || []).map(item => ({ label: item.label, value: item.value })));
const statusOptions = computed(() => (sys_common_status.value || []).map(item => ({ label: item.label, value: item.value })));

const data = reactive({
  form: {},
  queryParams: {
    pageNum: 1,
    pageSize: 10,
    title: undefined,
    operName: undefined,
    businessType: undefined,
    status: undefined,
    orderByColumn: undefined,
    isAsc: undefined
  }
});

const { queryParams, form } = toRefs(data);

function getList() {
  loading.value = true;
  list(proxy.addDateRange(queryParams.value, dateRange.value)).then(response => {
    operlogList.value = response.data.rows;
    total.value = response.data.total;
    loading.value = false;
  });
}

function typeFormat(row) {
  return proxy.selectDictLabel(sys_oper_type, row.businessType);
}

function handleQuery() {
  queryParams.value.pageNum = 1;
  getList();
}

function resetQuery() {
  dateRange.value = [];
  proxy.resetForm("queryRef");
  queryParams.value.orderByColumn = "operTime";
  queryParams.value.isAsc = "descending";
  handleQuery();
}

function handleSelectionChange(selection, keys) {
  ids.value = selection.map(item => item.operId);
  selectedRowKeys.value = keys;
  multiple.value = !selection.length;
}

function handleTableChange(_pagination, _filters, sorter) {
  queryParams.value.orderByColumn = sorter.field;
  queryParams.value.isAsc = sorter.order;
  getList();
}

function handleView(row) {
  open.value = true;
  form.value = row;
}

function handleDelete(row) {
  const operIds = row?.operId || ids.value;
  proxy.$modal.confirm(`是否确认删除日志编号为“${operIds}”的数据项？`).then(() => {
    return delOperlog(operIds);
  }).then(() => {
    getList();
    proxy.$modal.msgSuccess("删除成功");
  }).catch(() => {});
}

function handleClean() {
  proxy.$modal.confirm("是否确认清空所有操作日志数据项？").then(() => {
    return cleanOperlog();
  }).then(() => {
    getList();
    proxy.$modal.msgSuccess("清空成功");
  }).catch(() => {});
}

function handleExport() {
  proxy.download("monitor/operlog/export", { ...queryParams.value }, `config_${new Date().getTime()}.xlsx`);
}

getList();
</script>

<style scoped>
.search-form {
  margin-bottom: 16px;
}

.toolbar-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  margin-top: 24px;
}
</style>
