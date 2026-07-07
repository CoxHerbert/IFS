<template>
  <div class="app-container">
    <a-form v-show="showSearch" ref="queryRef" :model="queryParams" layout="inline" class="search-form">
      <a-form-item label="公告标题" name="noticeTitle">
        <a-input v-model:value="queryParams.noticeTitle" allow-clear placeholder="请输入公告标题" @pressEnter="handleQuery" />
      </a-form-item>
      <a-form-item label="操作人员" name="createBy">
        <a-input v-model:value="queryParams.createBy" allow-clear placeholder="请输入操作人员" @pressEnter="handleQuery" />
      </a-form-item>
      <a-form-item label="类型" name="noticeType">
        <a-select v-model:value="queryParams.noticeType" allow-clear placeholder="公告类型" style="width: 160px" :options="noticeTypeOptions" />
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
        <a-button type="primary" @click="handleAdd" v-hasPermi="['system:notice:add']">新增</a-button>
        <a-button :disabled="single" @click="handleUpdate()" v-hasPermi="['system:notice:edit']">修改</a-button>
        <a-button danger :disabled="multiple" @click="handleDelete()" v-hasPermi="['system:notice:remove']">删除</a-button>
      </a-space>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList" />
    </div>

    <a-table
      :loading="loading"
      :data-source="noticeList"
      :columns="columns"
      :pagination="false"
      :row-selection="rowSelection"
      row-key="noticeId"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'noticeType'">
          <dict-tag :options="sys_notice_type" :value="record.noticeType" />
        </template>
        <template v-else-if="column.key === 'status'">
          <dict-tag :options="sys_notice_status" :value="record.status" />
        </template>
        <template v-else-if="column.key === 'createTime'">
          {{ parseTime(record.createTime, '{y}-{m}-{d}') }}
        </template>
        <template v-else-if="column.key === 'action'">
          <a-space>
            <a-button type="link" @click="handleUpdate(record)" v-hasPermi="['system:notice:edit']">修改</a-button>
            <a-button type="link" danger @click="handleDelete(record)" v-hasPermi="['system:notice:remove']">删除</a-button>
          </a-space>
        </template>
      </template>
    </a-table>

    <pagination v-show="total > 0" v-model:page="queryParams.pageNum" v-model:limit="queryParams.pageSize" :total="total" @pagination="getList" />

    <a-modal v-model:open="open" :title="title" width="780px" :footer="null" destroy-on-close>
      <a-form ref="noticeRef" :model="form" :rules="rules" :label-col="{ style: { width: '80px' } }">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="公告标题" name="noticeTitle">
              <a-input v-model:value="form.noticeTitle" placeholder="请输入公告标题" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="公告类型" name="noticeType">
              <a-select v-model:value="form.noticeType" placeholder="请选择" :options="noticeTypeOptions" />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item label="状态" name="status">
              <a-radio-group v-model:value="form.status">
                <a-radio v-for="dict in sys_notice_status" :key="dict.value" :value="dict.value">{{ dict.label }}</a-radio>
              </a-radio-group>
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item label="内容" name="noticeContent">
              <a-textarea v-model:value="form.noticeContent" :rows="6" placeholder="请输入内容" />
            </a-form-item>
          </a-col>
        </a-row>
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

<script setup name="Notice">
import { listNotice, getNotice, delNotice, addNotice, updateNotice } from "@/api/system/notice";

const { proxy } = getCurrentInstance();
const { sys_notice_status, sys_notice_type } = proxy.useDict("sys_notice_status", "sys_notice_type");

const columns = [
  { title: "序号", dataIndex: "noticeId", key: "noticeId", width: 100, align: "center" },
  { title: "公告标题", dataIndex: "noticeTitle", key: "noticeTitle", ellipsis: true, align: "center" },
  { title: "公告类型", dataIndex: "noticeType", key: "noticeType", width: 100, align: "center" },
  { title: "状态", dataIndex: "status", key: "status", width: 100, align: "center" },
  { title: "创建者", dataIndex: "createBy", key: "createBy", width: 100, align: "center" },
  { title: "创建时间", dataIndex: "createTime", key: "createTime", width: 120, align: "center" },
  { title: "操作", key: "action", align: "center" }
];

const noticeList = ref([]);
const open = ref(false);
const loading = ref(true);
const showSearch = ref(true);
const ids = ref([]);
const selectedRowKeys = ref([]);
const single = ref(true);
const multiple = ref(true);
const total = ref(0);
const title = ref("");

const rowSelection = computed(() => ({
  selectedRowKeys: selectedRowKeys.value,
  onChange: (keys, rows) => handleSelectionChange(rows, keys)
}));

const noticeTypeOptions = computed(() => (sys_notice_type.value || []).map(item => ({ label: item.label, value: item.value })));

const data = reactive({
  form: {},
  queryParams: {
    pageNum: 1,
    pageSize: 10,
    noticeTitle: undefined,
    createBy: undefined,
    status: undefined
  },
  rules: {
    noticeTitle: [{ required: true, message: "公告标题不能为空", trigger: "blur" }],
    noticeType: [{ required: true, message: "公告类型不能为空", trigger: "change" }]
  }
});

const { queryParams, form, rules } = toRefs(data);

function getList() {
  loading.value = true;
  listNotice(queryParams.value).then(response => {
    noticeList.value = response.rows;
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
    noticeId: undefined,
    noticeTitle: undefined,
    noticeType: undefined,
    noticeContent: undefined,
    status: "0"
  };
  proxy.resetForm("noticeRef");
}

function handleQuery() {
  queryParams.value.pageNum = 1;
  getList();
}

function resetQuery() {
  proxy.resetForm("queryRef");
  handleQuery();
}

function handleSelectionChange(selection, keys) {
  ids.value = selection.map(item => item.noticeId);
  selectedRowKeys.value = keys;
  single.value = selection.length != 1;
  multiple.value = !selection.length;
}

function handleAdd() {
  reset();
  open.value = true;
  title.value = "添加公告";
}

function handleUpdate(row) {
  reset();
  const noticeId = row?.noticeId || ids.value;
  getNotice(noticeId).then(response => {
    form.value = response.data;
    open.value = true;
    title.value = "修改公告";
  });
}

function submitForm() {
  proxy.$refs.noticeRef.validate().then(() => {
    const request = form.value.noticeId != undefined ? updateNotice(form.value) : addNotice(form.value);
    request.then(() => {
      proxy.$modal.msgSuccess(form.value.noticeId != undefined ? "修改成功" : "新增成功");
      open.value = false;
      getList();
    });
  }).catch(() => {});
}

function handleDelete(row) {
  const noticeIds = row?.noticeId || ids.value;
  proxy.$modal.confirm(`是否确认删除公告编号为“${noticeIds}”的数据项？`).then(() => delNotice(noticeIds)).then(() => {
    getList();
    proxy.$modal.msgSuccess("删除成功");
  }).catch(() => {});
}

getList();
</script>

<style scoped>
.search-form { margin-bottom: 16px; }
.toolbar-row { display: flex; align-items: center; justify-content: space-between; gap: 16px; }
.modal-footer { display: flex; justify-content: flex-end; margin-top: 24px; }
</style>
