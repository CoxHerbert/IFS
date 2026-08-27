<template>
  <div class="app-container">
    <a-form ref="queryRef" :model="queryParams" layout="inline" class="search-form">
      <a-form-item label="登录地址" name="ipaddr">
        <a-input
          v-model:value="queryParams.ipaddr"
          allow-clear
          placeholder="请输入登录地址"
          @pressEnter="handleQuery"
        />
      </a-form-item>
      <a-form-item label="用户名称" name="userName">
        <a-input
          v-model:value="queryParams.userName"
          allow-clear
          placeholder="请输入用户名称"
          @pressEnter="handleQuery"
        />
      </a-form-item>
      <a-form-item>
        <a-space>
          <a-button type="primary" @click="handleQuery">搜索</a-button>
          <a-button @click="resetQuery">重置</a-button>
        </a-space>
      </a-form-item>
    </a-form>

    <a-table
      :loading="loading"
      :data-source="pageList"
      :columns="columns"
      :pagination="false"
      row-key="tokenId"
    >
      <template #bodyCell="{ column, record, index }">
        <template v-if="column.key === 'index'">
          {{ (pageNum - 1) * pageSize + index + 1 }}
        </template>
        <template v-else-if="column.key === 'loginTime'">
          {{ parseTime(record.loginTime) }}
        </template>
        <template v-else-if="column.key === 'action'">
          <a-button type="link" danger @click="handleForceLogout(record)" v-hasPermi="['monitor:online:forceLogout']">强退</a-button>
        </template>
      </template>
    </a-table>

    <pagination v-show="total > 0" v-model:page="pageNum" v-model:limit="pageSize" :total="total" />
  </div>
</template>

<script setup name="Online">
import { forceLogout, list as initData } from "@/api/monitor/online";

const { proxy } = getCurrentInstance();

const columns = [
  { title: "序号", key: "index", width: 80, align: "center" },
  { title: "会话编号", dataIndex: "tokenId", key: "tokenId", ellipsis: true, align: "center" },
  { title: "登录名称", dataIndex: "userName", key: "userName", ellipsis: true, align: "center" },
  { title: "所属部门", dataIndex: "deptName", key: "deptName", ellipsis: true, align: "center" },
  { title: "主机", dataIndex: "ipaddr", key: "ipaddr", ellipsis: true, align: "center" },
  { title: "登录地点", dataIndex: "loginLocation", key: "loginLocation", ellipsis: true, align: "center" },
  { title: "操作系统", dataIndex: "os", key: "os", ellipsis: true, align: "center" },
  { title: "浏览器", dataIndex: "browser", key: "browser", ellipsis: true, align: "center" },
  { title: "登录时间", dataIndex: "loginTime", key: "loginTime", width: 180, align: "center" },
  { title: "操作", key: "action", align: "center" }
];

const onlineList = ref([]);
const loading = ref(true);
const total = ref(0);
const pageNum = ref(1);
const pageSize = ref(10);

const queryParams = ref({
  ipaddr: undefined,
  userName: undefined
});

const pageList = computed(() => {
  return onlineList.value.slice((pageNum.value - 1) * pageSize.value, pageNum.value * pageSize.value);
});

function getList() {
  loading.value = true;
  initData(queryParams.value).then(response => {
    onlineList.value = response.data.rows;
    total.value = response.data.total;
    loading.value = false;
  });
}

function handleQuery() {
  pageNum.value = 1;
  getList();
}

function resetQuery() {
  proxy.resetForm("queryRef");
  handleQuery();
}

function handleForceLogout(row) {
  proxy.$modal.confirm(`是否确认强退名称为“${row.userName}”的用户？`).then(() => {
    return forceLogout(row.tokenId);
  }).then(() => {
    getList();
    proxy.$modal.msgSuccess("删除成功");
  }).catch(() => {});
}

getList();
</script>

<style scoped>
.search-form {
  margin-bottom: 16px;
}
</style>
