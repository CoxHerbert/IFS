<template>
  <div class="app-container">
    <a-form v-show="showSearch" ref="queryRef" :model="queryParams" layout="inline" class="search-form">
      <a-form-item label="参数名称" name="configName">
        <a-input v-model:value="queryParams.configName" allow-clear placeholder="请输入参数名称" style="width: 240px" @pressEnter="handleQuery" />
      </a-form-item>
      <a-form-item label="参数键名" name="configKey">
        <a-input v-model:value="queryParams.configKey" allow-clear placeholder="请输入参数键名" style="width: 240px" @pressEnter="handleQuery" />
      </a-form-item>
      <a-form-item label="系统内置" name="configType">
        <a-select v-model:value="queryParams.configType" allow-clear placeholder="系统内置" style="width: 160px" :options="yesNoOptions" />
      </a-form-item>
      <a-form-item label="创建时间">
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
        <a-button type="primary" @click="handleAdd" v-hasPermi="['system:config:add']">新增</a-button>
        <a-button :disabled="single" @click="handleUpdate()" v-hasPermi="['system:config:edit']">修改</a-button>
        <a-button danger :disabled="multiple" @click="handleDelete()" v-hasPermi="['system:config:remove']">删除</a-button>
        <a-button @click="handleExport" v-hasPermi="['system:config:export']">导出</a-button>
        <a-button danger @click="handleRefreshCache" v-hasPermi="['system:config:remove']">刷新缓存</a-button>
      </a-space>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList" />
    </div>

    <a-table
      :loading="loading"
      :data-source="configList"
      :columns="columns"
      :pagination="false"
      :row-selection="rowSelection"
      :scroll="{ x: 1400 }"
      row-key="configId"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'configType'">
          <dict-tag :options="sys_yes_no" :value="record.configType" />
        </template>
        <template v-else-if="column.key === 'createTime'">
          {{ parseTime(record.createTime) }}
        </template>
        <template v-else-if="column.key === 'action'">
          <a-space>
            <a-button type="link" @click="handleUpdate(record)" v-hasPermi="['system:config:edit']">修改</a-button>
            <a-button type="link" danger @click="handleDelete(record)" v-hasPermi="['system:config:remove']">删除</a-button>
          </a-space>
        </template>
      </template>
    </a-table>

    <pagination v-show="total > 0" v-model:page="queryParams.pageNum" v-model:limit="queryParams.pageSize" :total="total" @pagination="getList" />

    <a-modal v-model:open="open" :title="title" width="500px" :footer="null" destroy-on-close>
      <a-form ref="configRef" :model="form" :rules="rules" :label-col="{ style: { width: '80px' } }">
        <a-form-item label="参数名称" name="configName">
          <a-input v-model:value="form.configName" placeholder="请输入参数名称" />
        </a-form-item>
        <a-form-item label="参数键名" name="configKey">
          <a-input v-model:value="form.configKey" placeholder="请输入参数键名" />
        </a-form-item>
        <a-form-item label="参数键值" name="configValue">
          <a-input v-model:value="form.configValue" placeholder="请输入参数键值" />
        </a-form-item>
        <a-form-item label="系统内置" name="configType">
          <a-radio-group v-model:value="form.configType">
            <a-radio v-for="dict in sys_yes_no" :key="dict.value" :value="dict.value">{{ dict.label }}</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="备注" name="remark">
          <a-textarea v-model:value="form.remark" placeholder="请输入内容" />
        </a-form-item>
      </a-form>
      <div class="modal-footer">
        <a-space>
          <a-button type="primary" @click="submitForm">确定</a-button>
          <a-button @click="cancel">取消</a-button>
        </a-space>
      </div>
    </a-modal>
  </div>
</template>

<script setup name="Config">
import { listConfig, getConfig, delConfig, addConfig, updateConfig, refreshCache } from "@/api/system/config";

const { proxy } = getCurrentInstance();
const { sys_yes_no } = proxy.useDict("sys_yes_no");

const columns = [
  { title: "参数主键", dataIndex: "configId", key: "configId", align: "center" },
  { title: "参数名称", dataIndex: "configName", key: "configName", ellipsis: true, align: "center" },
  { title: "参数键名", dataIndex: "configKey", key: "configKey", ellipsis: true, align: "center" },
  { title: "参数键值", dataIndex: "configValue", key: "configValue", align: "center" },
  { title: "系统内置", dataIndex: "configType", key: "configType", align: "center" },
  { title: "备注", dataIndex: "remark", key: "remark", ellipsis: true, align: "center" },
  { title: "创建时间", dataIndex: "createTime", key: "createTime", width: 180, align: "center" },
  { title: "操作", key: "action", width: 150, align: "center" }
];

const configList = ref([]);
const open = ref(false);
const loading = ref(true);
const showSearch = ref(true);
const ids = ref([]);
const selectedRowKeys = ref([]);
const single = ref(true);
const multiple = ref(true);
const total = ref(0);
const title = ref("");
const dateRange = ref([]);

const rowSelection = computed(() => ({
  selectedRowKeys: selectedRowKeys.value,
  onChange: (keys, rows) => handleSelectionChange(rows, keys)
}));

const yesNoOptions = computed(() => (sys_yes_no.value || []).map(item => ({ label: item.label, value: item.value })));

const data = reactive({
  form: {},
  queryParams: {
    pageNum: 1,
    pageSize: 10,
    configName: undefined,
    configKey: undefined,
    configType: undefined
  },
  rules: {
    configName: [{ required: true, message: "参数名称不能为空", trigger: "blur" }],
    configKey: [{ required: true, message: "参数键名不能为空", trigger: "blur" }],
    configValue: [{ required: true, message: "参数键值不能为空", trigger: "blur" }]
  }
});

const { queryParams, form, rules } = toRefs(data);

function getList() {
  loading.value = true;
  listConfig(proxy.addDateRange(queryParams.value, dateRange.value)).then(response => {
    configList.value = response.rows;
    total.value = response.total;
    loading.value = false;
  });
}

function cancel() {
  open.value = false;
  reset();
}

function reset() {
  form.value = {
    configId: undefined,
    configName: undefined,
    configKey: undefined,
    configValue: undefined,
    configType: "Y",
    remark: undefined
  };
  proxy.resetForm("configRef");
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
  ids.value = selection.map(item => item.configId);
  selectedRowKeys.value = keys;
  single.value = selection.length != 1;
  multiple.value = !selection.length;
}

function handleAdd() {
  reset();
  open.value = true;
  title.value = "添加参数";
}

function handleUpdate(row) {
  reset();
  const configId = row?.configId || ids.value;
  getConfig(configId).then(response => {
    form.value = response.data;
    open.value = true;
    title.value = "修改参数";
  });
}

function submitForm() {
  proxy.$refs.configRef.validate().then(() => {
    const request = form.value.configId != undefined ? updateConfig(form.value) : addConfig(form.value);
    request.then(() => {
      proxy.$modal.msgSuccess(form.value.configId != undefined ? "修改成功" : "新增成功");
      open.value = false;
      getList();
    });
  }).catch(() => {});
}

function handleDelete(row) {
  const configIds = row?.configId || ids.value;
  proxy.$modal.confirm(`是否确认删除参数编号为“${configIds}”的数据项？`).then(() => delConfig(configIds)).then(() => {
    getList();
    proxy.$modal.msgSuccess("删除成功");
  }).catch(() => {});
}

function handleExport() {
  proxy.download("system/config/export", { ...queryParams.value }, `config_${new Date().getTime()}.xlsx`);
}

function handleRefreshCache() {
  refreshCache().then(() => {
    proxy.$modal.msgSuccess("刷新缓存成功");
  });
}

getList();
</script>

<style scoped>
.search-form { margin-bottom: 16px; }
.toolbar-row { display: flex; align-items: center; justify-content: space-between; gap: 16px; }
.modal-footer { display: flex; justify-content: flex-end; margin-top: 24px; }
</style>
