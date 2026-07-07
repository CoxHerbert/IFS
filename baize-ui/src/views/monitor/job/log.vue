<template>
  <div class="app-container">
    <a-form v-show="showSearch" ref="queryRef" :model="queryParams" layout="inline" class="search-form">
      <a-form-item label="任务名称" name="jobName">
        <a-input v-model:value="queryParams.jobName" allow-clear placeholder="请输入任务名称" style="width: 240px" @pressEnter="handleQuery" />
      </a-form-item>
      <a-form-item label="任务组名" name="jobGroup">
        <a-select v-model:value="queryParams.jobGroup" allow-clear placeholder="请选择任务组名" style="width: 240px" :options="jobGroupOptions" />
      </a-form-item>
      <a-form-item label="执行状态" name="status">
        <a-select v-model:value="queryParams.status" allow-clear placeholder="请选择执行状态" style="width: 240px" :options="statusOptions" />
      </a-form-item>
      <a-form-item label="执行时间">
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
        <a-button danger :disabled="multiple" @click="handleDelete()" v-hasPermi="['monitor:job:remove']">删除</a-button>
        <a-button danger @click="handleClean" v-hasPermi="['monitor:job:remove']">清空</a-button>
        <a-button @click="handleExport" v-hasPermi="['monitor:job:export']">导出</a-button>
        <a-button @click="handleClose">关闭</a-button>
      </a-space>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList" />
    </div>

    <a-table
      :loading="loading"
      :data-source="jobLogList"
      :columns="columns"
      :pagination="false"
      :row-selection="rowSelection"
      :scroll="{ x: 1400 }"
      row-key="jobLogId"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'jobGroup'">
          <dict-tag :options="sys_job_group" :value="record.jobGroup" />
        </template>
        <template v-else-if="column.key === 'status'">
          <dict-tag :options="sys_common_status" :value="record.status" />
        </template>
        <template v-else-if="column.key === 'createTime'">
          {{ parseTime(record.createTime) }}
        </template>
        <template v-else-if="column.key === 'action'">
          <a-button type="link" @click="handleView(record)" v-hasPermi="['monitor:job:query']">详情</a-button>
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

    <a-modal v-model:open="open" title="调度日志详情" width="700px" :footer="null" destroy-on-close>
      <a-descriptions :column="2" bordered size="small">
        <a-descriptions-item label="日志序号">{{ form.jobLogId }}</a-descriptions-item>
        <a-descriptions-item label="任务名称">{{ form.jobName }}</a-descriptions-item>
        <a-descriptions-item label="任务分组">{{ form.jobGroup }}</a-descriptions-item>
        <a-descriptions-item label="执行时间">{{ form.createTime }}</a-descriptions-item>
        <a-descriptions-item label="调用方法" :span="2">{{ form.invokeTarget }}</a-descriptions-item>
        <a-descriptions-item label="日志信息" :span="2">{{ form.jobMessage }}</a-descriptions-item>
        <a-descriptions-item label="执行状态">{{ form.status == 0 ? "正常" : "失败" }}</a-descriptions-item>
        <a-descriptions-item v-if="form.status == 1" label="异常信息" :span="2">{{ form.exceptionInfo }}</a-descriptions-item>
      </a-descriptions>
      <div class="modal-footer">
        <a-button @click="open = false">关闭</a-button>
      </div>
    </a-modal>
  </div>
</template>

<script setup name="JobLog">
import { getJob } from "@/api/monitor/job";
import { listJobLog, delJobLog, cleanJobLog } from "@/api/monitor/jobLog";

const { proxy } = getCurrentInstance();
const { sys_common_status, sys_job_group } = proxy.useDict("sys_common_status", "sys_job_group");

const columns = [
  { title: "日志编号", dataIndex: "jobLogId", key: "jobLogId", width: 80, align: "center" },
  { title: "任务名称", dataIndex: "jobName", key: "jobName", ellipsis: true, align: "center" },
  { title: "任务组名", dataIndex: "jobGroup", key: "jobGroup", ellipsis: true, align: "center" },
  { title: "调用目标字符串", dataIndex: "invokeTarget", key: "invokeTarget", ellipsis: true, align: "center" },
  { title: "日志信息", dataIndex: "jobMessage", key: "jobMessage", ellipsis: true, align: "center" },
  { title: "执行状态", dataIndex: "status", key: "status", align: "center" },
  { title: "执行时间", dataIndex: "createTime", key: "createTime", width: 180, align: "center" },
  { title: "操作", key: "action", align: "center" }
];

const jobLogList = ref([]);
const open = ref(false);
const loading = ref(true);
const showSearch = ref(true);
const ids = ref([]);
const selectedRowKeys = ref([]);
const multiple = ref(true);
const total = ref(0);
const dateRange = ref([]);
const route = useRoute();

const rowSelection = computed(() => ({
  selectedRowKeys: selectedRowKeys.value,
  onChange: (keys, rows) => handleSelectionChange(rows, keys)
}));

const jobGroupOptions = computed(() => (sys_job_group.value || []).map(item => ({ label: item.label, value: item.value })));
const statusOptions = computed(() => (sys_common_status.value || []).map(item => ({ label: item.label, value: item.value })));

const data = reactive({
  form: {},
  queryParams: {
    pageNum: 1,
    pageSize: 10,
    jobName: undefined,
    jobGroup: undefined,
    status: undefined
  }
});

const { queryParams, form } = toRefs(data);

function getList() {
  loading.value = true;
  listJobLog(proxy.addDateRange(queryParams.value, dateRange.value)).then(response => {
    jobLogList.value = response.rows;
    total.value = response.total;
    loading.value = false;
  });
}

function handleClose() {
  const obj = { path: "/monitor/job" };
  proxy.$tab.closeOpenPage(obj);
}

function handleQuery() {
  queryParams.value.pageNum = 1;
  getList();
}

function resetQuery() {
  dateRange.value = [];
  proxy.resetForm("queryRef");
  handleQuery();
}

function handleSelectionChange(selection, keys) {
  ids.value = selection.map(item => item.jobLogId);
  selectedRowKeys.value = keys;
  multiple.value = !selection.length;
}

function handleView(row) {
  open.value = true;
  form.value = row;
}

function handleDelete() {
  proxy.$modal.confirm(`是否确认删除调度日志编号为“${ids.value}”的数据项？`).then(() => {
    return delJobLog(ids.value);
  }).then(() => {
    getList();
    proxy.$modal.msgSuccess("删除成功");
  }).catch(() => {});
}

function handleClean() {
  proxy.$modal.confirm("是否确认清空所有调度日志数据项？").then(() => {
    return cleanJobLog();
  }).then(() => {
    getList();
    proxy.$modal.msgSuccess("清空成功");
  }).catch(() => {});
}

function handleExport() {
  proxy.download("monitor/jobLog/export", { ...queryParams.value }, `job_log_${new Date().getTime()}.xlsx`);
}

(() => {
  const jobId = route.query.jobId;
  if (jobId !== undefined && jobId != 0) {
    getJob(jobId).then(response => {
      queryParams.value.jobName = response.data.jobName;
      queryParams.value.jobGroup = response.data.jobGroup;
      getList();
    });
  } else {
    getList();
  }
})();
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
