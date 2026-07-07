<template>
  <div class="app-container">
    <a-form v-show="showSearch" ref="queryRef" :model="queryParams" layout="inline" class="search-form">
      <a-form-item label="字典名称" name="dictName">
        <a-input v-model:value="queryParams.dictName" allow-clear placeholder="请输入字典名称" style="width: 240px" @pressEnter="handleQuery" />
      </a-form-item>
      <a-form-item label="字典类型" name="dictType">
        <a-input v-model:value="queryParams.dictType" allow-clear placeholder="请输入字典类型" style="width: 240px" @pressEnter="handleQuery" />
      </a-form-item>
      <a-form-item label="状态" name="status">
        <a-select v-model:value="queryParams.status" allow-clear placeholder="字典状态" style="width: 240px" :options="statusOptions" />
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
        <a-button type="primary" @click="handleAdd" v-hasPermi="['system:dict:add']">新增</a-button>
        <a-button :disabled="single" @click="handleUpdate()" v-hasPermi="['system:dict:edit']">修改</a-button>
        <a-button danger :disabled="multiple" @click="handleDelete()" v-hasPermi="['system:dict:remove']">删除</a-button>
        <a-button @click="handleExport" v-hasPermi="['system:dict:export']">导出</a-button>
        <a-button danger @click="handleRefreshCache" v-hasPermi="['system:dict:remove']">刷新缓存</a-button>
      </a-space>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList" />
    </div>

    <a-table
      :loading="loading"
      :data-source="typeList"
      :columns="columns"
      :pagination="false"
      :row-selection="rowSelection"
      row-key="dictId"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'dictType'">
          <router-link :to="'/system/dict-data/index/' + record.dictId" class="link-type">{{ record.dictType }}</router-link>
        </template>
        <template v-else-if="column.key === 'status'">
          <dict-tag :options="sys_normal_disable" :value="record.status" />
        </template>
        <template v-else-if="column.key === 'createTime'">
          {{ parseTime(record.createTime) }}
        </template>
        <template v-else-if="column.key === 'action'">
          <a-space>
            <a-button type="link" @click="handleUpdate(record)" v-hasPermi="['system:dict:edit']">修改</a-button>
            <a-button type="link" danger @click="handleDelete(record)" v-hasPermi="['system:dict:remove']">删除</a-button>
          </a-space>
        </template>
      </template>
    </a-table>

    <pagination v-show="total > 0" v-model:page="queryParams.pageNum" v-model:limit="queryParams.pageSize" :total="total" @pagination="getList" />

    <a-modal v-model:open="open" :title="title" width="500px" :footer="null" destroy-on-close>
      <a-form ref="dictRef" :model="form" :rules="rules" :label-col="{ style: { width: '80px' } }">
        <a-form-item label="字典名称" name="dictName">
          <a-input v-model:value="form.dictName" placeholder="请输入字典名称" />
        </a-form-item>
        <a-form-item label="字典类型" name="dictType">
          <a-input v-model:value="form.dictType" placeholder="请输入字典类型" />
        </a-form-item>
        <a-form-item label="状态" name="status">
          <a-radio-group v-model:value="form.status">
            <a-radio v-for="dict in sys_normal_disable" :key="dict.value" :value="dict.value">{{ dict.label }}</a-radio>
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

<script setup name="Dict">
import { listType, getType, delType, addType, updateType, refreshCache } from "@/api/system/dict/type";

const { proxy } = getCurrentInstance();
const { sys_normal_disable } = proxy.useDict("sys_normal_disable");

const columns = [
  { title: "字典编号", dataIndex: "dictId", key: "dictId", align: "center" },
  { title: "字典名称", dataIndex: "dictName", key: "dictName", ellipsis: true, align: "center" },
  { title: "字典类型", key: "dictType", ellipsis: true, align: "center" },
  { title: "状态", dataIndex: "status", key: "status", align: "center" },
  { title: "备注", dataIndex: "remark", key: "remark", ellipsis: true, align: "center" },
  { title: "创建时间", dataIndex: "createTime", key: "createTime", width: 180, align: "center" },
  { title: "操作", key: "action", align: "center" }
];

const typeList = ref([]);
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

const statusOptions = computed(() => (sys_normal_disable.value || []).map(item => ({ label: item.label, value: item.value })));

const data = reactive({
  form: {},
  queryParams: {
    pageNum: 1,
    pageSize: 10,
    dictName: undefined,
    dictType: undefined,
    status: undefined
  },
  rules: {
    dictName: [{ required: true, message: "字典名称不能为空", trigger: "blur" }],
    dictType: [{ required: true, message: "字典类型不能为空", trigger: "blur" }]
  }
});

const { queryParams, form, rules } = toRefs(data);

function getList() {
  loading.value = true;
  listType(proxy.addDateRange(queryParams.value, dateRange.value)).then(response => {
    typeList.value = response.data.rows;
    total.value = response.data.total;
    loading.value = false;
  });
}

function cancel() {
  open.value = false;
  reset();
}

function reset() {
  form.value = {
    dictId: undefined,
    dictName: undefined,
    dictType: undefined,
    status: "0",
    remark: undefined
  };
  proxy.resetForm("dictRef");
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

function handleAdd() {
  reset();
  open.value = true;
  title.value = "添加字典类型";
}

function handleSelectionChange(selection, keys) {
  ids.value = selection.map(item => item.dictId);
  selectedRowKeys.value = keys;
  single.value = selection.length != 1;
  multiple.value = !selection.length;
}

function handleUpdate(row) {
  reset();
  const dictId = row?.dictId || ids.value;
  getType(dictId).then(response => {
    form.value = response.data;
    open.value = true;
    title.value = "修改字典类型";
  });
}

function submitForm() {
  proxy.$refs.dictRef.validate().then(() => {
    const request = form.value.dictId != undefined ? updateType(form.value) : addType(form.value);
    request.then(() => {
      proxy.$modal.msgSuccess(form.value.dictId != undefined ? "修改成功" : "新增成功");
      open.value = false;
      getList();
    });
  }).catch(() => {});
}

function handleDelete(row) {
  const dictIds = row?.dictId || ids.value;
  proxy.$modal.confirm(`是否确认删除字典编号为“${dictIds}”的数据项？`).then(() => delType(dictIds)).then(() => {
    getList();
    proxy.$modal.msgSuccess("删除成功");
  }).catch(() => {});
}

function handleExport() {
  proxy.download("system/dict/type/export", { ...queryParams.value }, `dict_${new Date().getTime()}.xlsx`);
}

function handleRefreshCache() {
  refreshCache().then(() => {
    proxy.$modal.msgSuccess("刷新成功");
  });
}

getList();
</script>

<style scoped>
.search-form { margin-bottom: 16px; }
.toolbar-row { display: flex; align-items: center; justify-content: space-between; gap: 16px; }
.modal-footer { display: flex; justify-content: flex-end; margin-top: 24px; }
</style>
