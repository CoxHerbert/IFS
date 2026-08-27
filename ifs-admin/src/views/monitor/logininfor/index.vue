<template>
  <div class="app-container">
    <a-form v-show="showSearch" ref="queryRef" :model="queryParams" layout="inline" class="search-form">
      <a-form-item label="登录地址" name="ipAddr">
        <a-input
          v-model:value="queryParams.ipAddr"
          allow-clear
          placeholder="请输入登录地址"
          style="width: 240px"
          @pressEnter="handleQuery"
        />
      </a-form-item>
      <a-form-item label="用户名称" name="userName">
        <a-input
          v-model:value="queryParams.userName"
          allow-clear
          placeholder="请输入用户名称"
          style="width: 240px"
          @pressEnter="handleQuery"
        />
      </a-form-item>
      <a-form-item label="状态" name="status">
        <a-select
          v-model:value="queryParams.status"
          allow-clear
          placeholder="登录状态"
          style="width: 240px"
          :options="statusOptions"
        />
      </a-form-item>
      <a-form-item label="登录时间">
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
        <a-button danger :disabled="multiple" @click="handleDelete()" v-hasPermi="['system:logininfor:remove']">删除</a-button>
        <a-button danger @click="handleClean" v-hasPermi="['system:logininfor:remove']">清空</a-button>
        <a-button @click="handleExport" v-hasPermi="['system:logininfor:export']">导出</a-button>
      </a-space>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList" />
    </div>

    <a-table
      :loading="loading"
      :data-source="logininforList"
      :columns="columns"
      :pagination="false"
      :row-selection="rowSelection"
      :scroll="{ x: 1500 }"
      row-key="infoId"
      @change="handleTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <dict-tag :options="sys_common_status" :value="record.status" />
        </template>
        <template v-else-if="column.key === 'loginTime'">
          {{ parseTime(record.loginTime) }}
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
  </div>
</template>

<script setup name="Logininfor">
import { list, delLogininfor, cleanLogininfor } from "@/api/monitor/logininfor";

const { proxy } = getCurrentInstance();
const { sys_common_status } = proxy.useDict("sys_common_status");

const columns = [
  { title: "访问编号", dataIndex: "infoId", key: "infoId", align: "center" },
  { title: "用户名称", dataIndex: "userName", key: "userName", ellipsis: true, sorter: true, align: "center" },
  { title: "地址", dataIndex: "ipAddr", key: "ipAddr", ellipsis: true, align: "center" },
  { title: "登录地点", dataIndex: "loginLocation", key: "loginLocation", ellipsis: true, align: "center" },
  { title: "操作系统", dataIndex: "os", key: "os", ellipsis: true, align: "center" },
  { title: "浏览器", dataIndex: "browser", key: "browser", ellipsis: true, align: "center" },
  { title: "登录状态", dataIndex: "status", key: "status", align: "center" },
  { title: "描述", dataIndex: "msg", key: "msg", align: "center" },
  { title: "访问时间", dataIndex: "loginTime", key: "loginTime", sorter: true, width: 180, align: "center" }
];

const logininforList = ref([]);
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

const statusOptions = computed(() => {
  return (sys_common_status.value || []).map(item => ({
    label: item.label,
    value: item.value
  }));
});

const queryParams = ref({
  pageNum: 1,
  pageSize: 10,
  ipAddr: undefined,
  userName: undefined,
  status: undefined,
  orderByColumn: undefined,
  isAsc: undefined
});

function getList() {
  loading.value = true;
  list(proxy.addDateRange(queryParams.value, dateRange.value)).then(response => {
    logininforList.value = response.data.rows;
    total.value = response.data.total;
    loading.value = false;
  });
}

function handleQuery() {
  queryParams.value.pageNum = 1;
  getList();
}

function resetQuery() {
  dateRange.value = [];
  proxy.resetForm("queryRef");
  queryParams.value.orderByColumn = "loginTime";
  queryParams.value.isAsc = "descending";
  handleQuery();
}

function handleSelectionChange(selection, keys) {
  ids.value = selection.map(item => item.infoId);
  selectedRowKeys.value = keys;
  multiple.value = !selection.length;
}

function handleTableChange(_pagination, _filters, sorter) {
  queryParams.value.orderByColumn = sorter.field;
  queryParams.value.isAsc = sorter.order;
  getList();
}

function handleDelete(row) {
  const infoIds = row?.infoId || ids.value;
  proxy.$modal.confirm(`是否确认删除访问编号为“${infoIds}”的数据项？`).then(() => {
    return delLogininfor(infoIds);
  }).then(() => {
    getList();
    proxy.$modal.msgSuccess("删除成功");
  }).catch(() => {});
}

function handleClean() {
  proxy.$modal.confirm("是否确认清空所有登录日志数据项？").then(() => {
    return cleanLogininfor();
  }).then(() => {
    getList();
    proxy.$modal.msgSuccess("清空成功");
  }).catch(() => {});
}

function handleExport() {
  proxy.download("monitor/logininfor/export", { ...queryParams.value }, `config_${new Date().getTime()}.xlsx`);
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
</style>
