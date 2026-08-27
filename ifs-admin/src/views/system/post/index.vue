<template>
  <div class="app-container">
    <a-form v-show="showSearch" ref="queryRef" :model="queryParams" layout="inline" class="search-form">
      <a-form-item label="岗位编码" name="postCode">
        <a-input v-model:value="queryParams.postCode" allow-clear placeholder="请输入岗位编码" @pressEnter="handleQuery" />
      </a-form-item>
      <a-form-item label="岗位名称" name="postName">
        <a-input v-model:value="queryParams.postName" allow-clear placeholder="请输入岗位名称" @pressEnter="handleQuery" />
      </a-form-item>
      <a-form-item label="状态" name="status">
        <a-select v-model:value="queryParams.status" allow-clear placeholder="岗位状态" style="width: 160px" :options="statusOptions" />
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
        <a-button type="primary" @click="handleAdd" v-hasPermi="['system:post:add']">新增</a-button>
        <a-button :disabled="single" @click="handleUpdate()" v-hasPermi="['system:post:edit']">修改</a-button>
        <a-button danger :disabled="multiple" @click="handleDelete()" v-hasPermi="['system:post:remove']">删除</a-button>
        <a-button @click="handleExport" v-hasPermi="['system:post:export']">导出</a-button>
      </a-space>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList" />
    </div>

    <a-table
      :loading="loading"
      :data-source="postList"
      :columns="columns"
      :pagination="false"
      :row-selection="rowSelection"
      row-key="postId"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <dict-tag :options="sys_normal_disable" :value="record.status" />
        </template>
        <template v-else-if="column.key === 'createTime'">
          {{ parseTime(record.createTime) }}
        </template>
        <template v-else-if="column.key === 'action'">
          <a-space>
            <a-button type="link" @click="handleUpdate(record)" v-hasPermi="['system:post:edit']">修改</a-button>
            <a-button type="link" danger @click="handleDelete(record)" v-hasPermi="['system:post:remove']">删除</a-button>
          </a-space>
        </template>
      </template>
    </a-table>

    <pagination v-show="total > 0" v-model:page="queryParams.pageNum" v-model:limit="queryParams.pageSize" :total="total" @pagination="getList" />

    <a-modal v-model:open="open" :title="title" width="500px" :footer="null" destroy-on-close>
      <a-form ref="postRef" :model="form" :rules="rules" :label-col="{ style: { width: '80px' } }">
        <a-form-item label="岗位名称" name="postName">
          <a-input v-model:value="form.postName" placeholder="请输入岗位名称" />
        </a-form-item>
        <a-form-item label="岗位编码" name="postCode">
          <a-input v-model:value="form.postCode" placeholder="请输入编码名称" />
        </a-form-item>
        <a-form-item label="岗位顺序" name="postSort">
          <a-input-number v-model:value="form.postSort" :min="0" style="width: 100%" />
        </a-form-item>
        <a-form-item label="岗位状态" name="status">
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

<script setup name="Post">
import { listPost, addPost, delPost, getPost, updatePost } from "@/api/system/post";

const { proxy } = getCurrentInstance();
const { sys_normal_disable } = proxy.useDict("sys_normal_disable");

const columns = [
  { title: "岗位编号", dataIndex: "postId", key: "postId", align: "center" },
  { title: "岗位编码", dataIndex: "postCode", key: "postCode", align: "center" },
  { title: "岗位名称", dataIndex: "postName", key: "postName", align: "center" },
  { title: "岗位排序", dataIndex: "postSort", key: "postSort", align: "center" },
  { title: "状态", dataIndex: "status", key: "status", align: "center" },
  { title: "创建时间", dataIndex: "createTime", key: "createTime", width: 180, align: "center" },
  { title: "操作", key: "action", align: "center" }
];

const postList = ref([]);
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

const statusOptions = computed(() => (sys_normal_disable.value || []).map(item => ({ label: item.label, value: item.value })));

const data = reactive({
  form: {},
  queryParams: {
    pageNum: 1,
    pageSize: 10,
    postCode: undefined,
    postName: undefined,
    status: undefined
  },
  rules: {
    postName: [{ required: true, message: "岗位名称不能为空", trigger: "blur" }],
    postCode: [{ required: true, message: "岗位编码不能为空", trigger: "blur" }],
    postSort: [{ required: true, message: "岗位顺序不能为空", trigger: "blur" }]
  }
});

const { queryParams, form, rules } = toRefs(data);

function getList() {
  loading.value = true;
  listPost(queryParams.value).then(response => {
    postList.value = response.data.rows;
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
    postId: undefined,
    postCode: undefined,
    postName: undefined,
    postSort: 0,
    status: "0",
    remark: undefined
  };
  proxy.resetForm("postRef");
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
  ids.value = selection.map(item => item.postId);
  selectedRowKeys.value = keys;
  single.value = selection.length != 1;
  multiple.value = !selection.length;
}

function handleAdd() {
  reset();
  open.value = true;
  title.value = "添加岗位";
}

function handleUpdate(row) {
  reset();
  const postId = row?.postId || ids.value;
  getPost(postId).then(response => {
    form.value = response.data;
    open.value = true;
    title.value = "修改岗位";
  });
}

function submitForm() {
  proxy.$refs.postRef.validate().then(() => {
    const request = form.value.postId != undefined ? updatePost(form.value) : addPost(form.value);
    request.then(() => {
      proxy.$modal.msgSuccess(form.value.postId != undefined ? "修改成功" : "新增成功");
      open.value = false;
      getList();
    });
  }).catch(() => {});
}

function handleDelete(row) {
  const postIds = row?.postId || ids.value;
  proxy.$modal.confirm(`是否确认删除岗位编号为“${postIds}”的数据项？`).then(() => delPost(postIds)).then(() => {
    getList();
    proxy.$modal.msgSuccess("删除成功");
  }).catch(() => {});
}

function handleExport() {
  proxy.download("system/post/export", { ...queryParams.value }, `post_${new Date().getTime()}.xlsx`);
}

getList();
</script>

<style scoped>
.search-form { margin-bottom: 16px; }
.toolbar-row { display: flex; align-items: center; justify-content: space-between; gap: 16px; }
.modal-footer { display: flex; justify-content: flex-end; margin-top: 24px; }
</style>
