<template>
  <div class="app-container">
    <a-form v-show="showSearch" ref="queryRef" :model="queryParams" layout="inline" class="search-form">
      <a-form-item label="通知标题" name="title">
        <a-input v-model:value="queryParams.title" allow-clear placeholder="请输入通知标题" @pressEnter="handleQuery" />
      </a-form-item>
      <a-form-item label="业务类型" name="bizType">
        <a-select v-model:value="queryParams.bizType" allow-clear placeholder="请选择业务类型" style="width: 160px">
          <a-select-option v-for="item in bizTypeOptions" :key="item.value" :value="item.value">
            {{ item.label }}
          </a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="读取状态" name="readFlag">
        <a-select v-model:value="queryParams.readFlag" allow-clear placeholder="请选择读取状态" style="width: 160px">
          <a-select-option value="0">未读</a-select-option>
          <a-select-option value="1">已读</a-select-option>
        </a-select>
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
        <a-button type="primary" :disabled="single" @click="handleRead" v-hasPermi="['system:notification:edit']">
          标记已读
        </a-button>
        <a-button @click="handleReadAll" v-hasPermi="['system:notification:edit']">全部已读</a-button>
        <a-button danger :disabled="multiple" @click="handleDelete()" v-hasPermi="['system:notification:remove']">
          删除
        </a-button>
      </a-space>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList" />
    </div>

    <a-table
      :loading="loading"
      :data-source="notificationList"
      :columns="columns"
      :pagination="false"
      :row-selection="rowSelection"
      row-key="notificationId"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'readFlag'">
          <a-tag :color="record.readFlag === '0' ? 'processing' : 'default'">
            {{ record.readFlag === '0' ? '未读' : '已读' }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'bizType'">
          {{ bizTypeLabel(record.bizType) }}
        </template>
        <template v-else-if="column.key === 'createTime'">
          {{ parseTime(record.createTime, '{y}-{m}-{d} {h}:{i}') }}
        </template>
        <template v-else-if="column.key === 'action'">
          <a-space>
            <a-button
              v-if="record.readFlag === '0'"
              type="link"
              @click="handleRead(record)"
              v-hasPermi="['system:notification:edit']"
            >
              标记已读
            </a-button>
            <a-button type="link" danger @click="handleDelete(record)" v-hasPermi="['system:notification:remove']">
              删除
            </a-button>
          </a-space>
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

<script setup name="Notification">
import { delNotification, listNotification, readAllNotification, readNotification } from "@/api/system/notification";

const { proxy } = getCurrentInstance();

const columns = [
  { title: "编号", dataIndex: "notificationId", key: "notificationId", width: 120, align: "center" },
  { title: "通知标题", dataIndex: "title", key: "title", width: 220, ellipsis: true },
  { title: "通知内容", dataIndex: "content", key: "content", ellipsis: true },
  { title: "业务类型", dataIndex: "bizType", key: "bizType", width: 120, align: "center" },
  { title: "状态", dataIndex: "readFlag", key: "readFlag", width: 100, align: "center" },
  { title: "发送人", dataIndex: "createBy", key: "createBy", width: 120, align: "center" },
  { title: "通知时间", dataIndex: "createTime", key: "createTime", width: 180, align: "center" },
  { title: "操作", key: "action", width: 180, align: "center" }
];

const bizTypeOptions = [
  { label: "出货计划", value: "shipment" }
];

const notificationList = ref([]);
const loading = ref(true);
const showSearch = ref(true);
const ids = ref([]);
const selectedRowKeys = ref([]);
const single = ref(true);
const multiple = ref(true);
const total = ref(0);

const rowSelection = computed(() => ({
  selectedRowKeys: selectedRowKeys.value,
  onChange: (keys, rows) => handleSelectionChange(rows, keys)
}));

const data = reactive({
  queryParams: {
    pageNum: 1,
    pageSize: 10,
    title: undefined,
    bizType: undefined,
    readFlag: undefined
  }
});

const { queryParams } = toRefs(data);

function getList() {
  loading.value = true;
  listNotification(queryParams.value).then(response => {
    const data = response.data || {};
    notificationList.value = data.rows || [];
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

function handleSelectionChange(selection, keys) {
  ids.value = selection.map(item => item.notificationId);
  selectedRowKeys.value = keys;
  single.value = selection.length !== 1;
  multiple.value = !selection.length;
}

function handleRead(row) {
  const notificationId = row?.notificationId || ids.value[0];
  if (!notificationId) {
    return;
  }
  readNotification(notificationId).then(() => {
    proxy.$modal.msgSuccess("已标记为已读");
    getList();
  });
}

function handleReadAll() {
  readAllNotification().then(() => {
    proxy.$modal.msgSuccess("全部通知已标记为已读");
    getList();
  });
}

function handleDelete(row) {
  const notificationIds = row?.notificationId || ids.value;
  proxy.$modal.confirm(`是否确认删除通知编号为“${notificationIds}”的数据项？`).then(() => delNotification(notificationIds)).then(() => {
    proxy.$modal.msgSuccess("删除成功");
    getList();
  }).catch(() => {});
}

function bizTypeLabel(value) {
  return bizTypeOptions.find(item => item.value === value)?.label || value || "-";
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
