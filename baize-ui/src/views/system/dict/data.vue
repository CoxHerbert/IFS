<template>
  <div class="app-container">
    <a-form v-show="showSearch" ref="queryRef" :model="queryParams" layout="inline" class="search-form">
      <a-form-item label="字典名称" name="dictType">
        <a-select v-model:value="queryParams.dictType" style="width: 180px" :options="typeSelectOptions" />
      </a-form-item>
      <a-form-item label="字典标签" name="dictLabel">
        <a-input v-model:value="queryParams.dictLabel" allow-clear placeholder="请输入字典标签" @pressEnter="handleQuery" />
      </a-form-item>
      <a-form-item label="状态" name="status">
        <a-select v-model:value="queryParams.status" allow-clear placeholder="数据状态" style="width: 160px" :options="statusOptions" />
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
        <a-button @click="handleClose">关闭</a-button>
      </a-space>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList" />
    </div>

    <a-table
      :loading="loading"
      :data-source="dataList"
      :columns="columns"
      :pagination="false"
      :row-selection="rowSelection"
      row-key="dictCode"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'dictLabel'">
          <span v-if="record.listClass === '' || record.listClass === 'default'">{{ record.dictLabel }}</span>
          <a-tag v-else :color="tagColor(record.listClass)">{{ record.dictLabel }}</a-tag>
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
      <a-form ref="dataRef" :model="form" :rules="rules" :label-col="{ style: { width: '80px' } }">
        <a-form-item label="字典类型">
          <a-input v-model:value="form.dictType" disabled />
        </a-form-item>
        <a-form-item label="数据标签" name="dictLabel">
          <a-input v-model:value="form.dictLabel" placeholder="请输入数据标签" />
        </a-form-item>
        <a-form-item label="数据键值" name="dictValue">
          <a-input v-model:value="form.dictValue" placeholder="请输入数据键值" />
        </a-form-item>
        <a-form-item label="样式属性" name="cssClass">
          <a-input v-model:value="form.cssClass" placeholder="请输入样式属性" />
        </a-form-item>
        <a-form-item label="显示排序" name="dictSort">
          <a-input-number v-model:value="form.dictSort" :min="0" style="width: 100%" />
        </a-form-item>
        <a-form-item label="回显样式" name="listClass">
          <a-select v-model:value="form.listClass" :options="listClassOptions" />
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

<script setup name="Data">
import { listType, getType } from "@/api/system/dict/type";
import { listData, getData, delData, addData, updateData } from "@/api/system/dict/data";

const { proxy } = getCurrentInstance();
const { sys_normal_disable } = proxy.useDict("sys_normal_disable");

const columns = [
  { title: "字典编码", dataIndex: "dictCode", key: "dictCode", align: "center" },
  { title: "字典标签", dataIndex: "dictLabel", key: "dictLabel", align: "center" },
  { title: "字典键值", dataIndex: "dictValue", key: "dictValue", align: "center" },
  { title: "字典排序", dataIndex: "dictSort", key: "dictSort", align: "center" },
  { title: "状态", dataIndex: "status", key: "status", align: "center" },
  { title: "备注", dataIndex: "remark", key: "remark", ellipsis: true, align: "center" },
  { title: "创建时间", dataIndex: "createTime", key: "createTime", width: 180, align: "center" },
  { title: "操作", key: "action", width: 150, align: "center" }
];

const dataList = ref([]);
const open = ref(false);
const loading = ref(true);
const showSearch = ref(true);
const ids = ref([]);
const selectedRowKeys = ref([]);
const single = ref(true);
const multiple = ref(true);
const total = ref(0);
const title = ref("");
const defaultDictType = ref("");
const typeOptions = ref([]);
const route = useRoute();

const listClassOptions = [
  { value: "default", label: "默认" },
  { value: "primary", label: "主要" },
  { value: "success", label: "成功" },
  { value: "info", label: "信息" },
  { value: "warning", label: "警告" },
  { value: "danger", label: "危险" }
];

const rowSelection = computed(() => ({
  selectedRowKeys: selectedRowKeys.value,
  onChange: (keys, rows) => handleSelectionChange(rows, keys)
}));

const typeSelectOptions = computed(() => typeOptions.value.map(item => ({ label: item.dictName, value: item.dictType })));
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
    dictLabel: [{ required: true, message: "数据标签不能为空", trigger: "blur" }],
    dictValue: [{ required: true, message: "数据键值不能为空", trigger: "blur" }],
    dictSort: [{ required: true, message: "数据顺序不能为空", trigger: "blur" }]
  }
});

const { queryParams, form, rules } = toRefs(data);

function getTypes(dictId) {
  getType(dictId).then(response => {
    queryParams.value.dictType = response.data.dictType;
    defaultDictType.value = response.data.dictType;
    getList();
  });
}

function getTypeList() {
  listType().then(response => {
    typeOptions.value = response.data.rows;
  });
}

function getList() {
  loading.value = true;
  listData(queryParams.value).then(response => {
    dataList.value = response.data.rows;
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
    dictCode: undefined,
    dictLabel: undefined,
    dictValue: undefined,
    cssClass: undefined,
    listClass: "default",
    dictSort: 0,
    status: "0",
    remark: undefined
  };
  proxy.resetForm("dataRef");
}

function handleQuery() {
  queryParams.value.pageNum = 1;
  getList();
}

function handleClose() {
  const obj = { path: "/system/dict" };
  proxy.$tab.closeOpenPage(obj);
}

function resetQuery() {
  proxy.resetForm("queryRef");
  queryParams.value.dictType = defaultDictType.value;
  handleQuery();
}

function handleAdd() {
  reset();
  open.value = true;
  title.value = "添加字典数据";
  form.value.dictType = queryParams.value.dictType;
}

function handleSelectionChange(selection, keys) {
  ids.value = selection.map(item => item.dictCode);
  selectedRowKeys.value = keys;
  single.value = selection.length != 1;
  multiple.value = !selection.length;
}

function handleUpdate(row) {
  reset();
  const dictCode = row?.dictCode || ids.value;
  getData(dictCode).then(response => {
    form.value = response.data;
    open.value = true;
    title.value = "修改字典数据";
  });
}

function submitForm() {
  proxy.$refs.dataRef.validate().then(() => {
    const request = form.value.dictCode != undefined ? updateData(form.value) : addData(form.value);
    request.then(() => {
      proxy.$modal.msgSuccess(form.value.dictCode != undefined ? "修改成功" : "新增成功");
      open.value = false;
      getList();
    });
  }).catch(() => {});
}

function handleDelete(row) {
  const dictCodes = row?.dictCode || ids.value;
  proxy.$modal.confirm(`是否确认删除字典编码为“${dictCodes}”的数据项？`).then(() => delData(dictCodes)).then(() => {
    getList();
    proxy.$modal.msgSuccess("删除成功");
  }).catch(() => {});
}

function handleExport() {
  proxy.download("system/dict/data/export", { ...queryParams.value }, `dict_data_${new Date().getTime()}.xlsx`);
}

function tagColor(type) {
  const map = { success: "success", info: "default", warning: "warning", danger: "error", primary: "processing" };
  return map[type] || "default";
}

getTypes(route.params && route.params.dictId);
getTypeList();
</script>

<style scoped>
.search-form { margin-bottom: 16px; }
.toolbar-row { display: flex; align-items: center; justify-content: space-between; gap: 16px; }
.modal-footer { display: flex; justify-content: flex-end; margin-top: 24px; }
</style>
