<template>
  <div class="app-container">
    <a-form
      v-show="showSearch"
      ref="queryRef"
      :model="queryParams"
      layout="inline"
      class="search-form"
    >
      <a-form-item label="客户" name="customerId">
        <a-select
          v-model:value="queryParams.customerId"
          show-search
          allow-clear
          :filter-option="false"
          placeholder="请选择客户"
          style="width: 260px"
          :options="customerSelectOptions"
          @search="loadCustomerOptions"
        />
      </a-form-item>
      <a-form-item label="账号" name="username">
        <a-input
          v-model:value="queryParams.username"
          allow-clear
          placeholder="请输入账号"
          style="width: 180px"
          @pressEnter="handleQuery"
        />
      </a-form-item>
      <a-form-item label="姓名" name="realName">
        <a-input
          v-model:value="queryParams.realName"
          allow-clear
          placeholder="请输入姓名"
          style="width: 160px"
          @pressEnter="handleQuery"
        />
      </a-form-item>
      <a-form-item label="状态" name="status">
        <a-select
          v-model:value="queryParams.status"
          allow-clear
          placeholder="账号状态"
          style="width: 140px"
          :options="statusOptions"
        />
      </a-form-item>
      <a-form-item>
        <a-space>
          <a-button type="primary" @click="handleQuery">搜索</a-button>
          <a-button @click="resetQuery">重置</a-button>
        </a-space>
      </a-form-item>
    </a-form>

    <a-alert
      v-if="route.query.customerName"
      :message="`当前客户：${route.query.customerName}`"
      type="info"
      show-icon
      class="mb8"
    />

    <div class="toolbar-row mb8">
      <a-space>
        <a-button type="primary" @click="handleAdd" v-hasPermi="['customer:account:add']">新增</a-button>
        <a-button :disabled="single" @click="handleUpdate()" v-hasPermi="['customer:account:edit']">修改</a-button>
        <a-button danger :disabled="multiple" @click="handleDelete()" v-hasPermi="['customer:account:remove']">删除</a-button>
      </a-space>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList" />
    </div>

    <a-table
      :loading="loading"
      :data-source="accountList"
      :columns="accountColumns"
      :pagination="false"
      :row-selection="accountRowSelection"
      :scroll="{ x: 1700 }"
      row-key="accountId"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'isMain'">
          <a-tag :color="record.isMain === '1' ? 'gold' : 'default'">
            {{ record.isMain === '1' ? '是' : '否' }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'status'">
          <a-tag :color="record.status === '0' ? 'success' : 'default'">
            {{ record.status === '0' ? '正常' : '停用' }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'lastLoginTime'">
          {{ parseTime(record.lastLoginTime) }}
        </template>
        <template v-else-if="column.key === 'action'">
          <a-space>
            <a-button type="link" @click="handleUpdate(record)" v-hasPermi="['customer:account:edit']">修改</a-button>
            <a-button type="link" @click="handleRoleAssign(record)" v-hasPermi="['customer:account:edit']">分配角色</a-button>
            <a-button type="link" @click="handleResetPwd(record)" v-hasPermi="['customer:account:resetPwd']">重置密码</a-button>
            <a-button type="link" danger @click="handleDelete(record)" v-hasPermi="['customer:account:remove']">删除</a-button>
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

    <a-modal v-model:open="open" :title="title" width="640px" :footer="null" destroy-on-close>
      <a-form ref="accountRef" :model="form" :rules="rules" :label-col="{ style: { width: '90px' } }">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="客户" name="customerId">
              <a-select
                v-model:value="form.customerId"
                show-search
                :filter-option="false"
                placeholder="请选择客户"
                :disabled="!!form.accountId"
                :options="customerSelectOptions"
                @search="loadCustomerOptions"
              />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="账号" name="username">
              <a-input v-model:value="form.username" placeholder="请输入登录账号" :maxlength="64" :disabled="!!form.accountId" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item v-if="!form.accountId" label="密码" name="password">
              <a-input-password v-model:value="form.password" placeholder="请输入登录密码" :maxlength="20" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="姓名" name="realName">
              <a-input v-model:value="form.realName" placeholder="请输入姓名" :maxlength="64" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="手机" name="phone">
              <a-input v-model:value="form.phone" placeholder="请输入手机号" :maxlength="64" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="邮箱" name="email">
              <a-input v-model:value="form.email" placeholder="请输入邮箱" :maxlength="128" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="主账号" name="isMain">
              <a-radio-group v-model:value="form.isMain">
                <a-radio value="0">否</a-radio>
                <a-radio value="1">是</a-radio>
              </a-radio-group>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="状态" name="status">
              <a-radio-group v-model:value="form.status">
                <a-radio value="0">正常</a-radio>
                <a-radio value="1">停用</a-radio>
              </a-radio-group>
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="备注" name="remark">
          <a-textarea v-model:value="form.remark" placeholder="请输入备注" :rows="3" :maxlength="500" />
        </a-form-item>
      </a-form>
      <div class="modal-footer">
        <a-space>
          <a-button type="primary" @click="submitForm">确定</a-button>
          <a-button @click="cancel">取消</a-button>
        </a-space>
      </div>
    </a-modal>

    <a-modal v-model:open="roleOpen" title="分配客户端角色" width="520px" :footer="null" destroy-on-close>
      <a-form :label-col="{ style: { width: '90px' } }">
        <a-form-item label="登录账号">
          <a-input :value="currentRoleAccount.username" disabled />
        </a-form-item>
        <a-form-item label="客户名称">
          <a-input :value="currentRoleAccount.customerName" disabled />
        </a-form-item>
        <a-form-item label="角色配置">
          <a-checkbox-group v-model:value="roleForm.roleIds">
            <a-space direction="vertical">
              <a-checkbox v-for="item in roleOptionsList" :key="item.roleId" :value="String(item.roleId)">
                {{ item.roleName }}
              </a-checkbox>
            </a-space>
          </a-checkbox-group>
        </a-form-item>
      </a-form>
      <div class="modal-footer">
        <a-space>
          <a-button type="primary" @click="submitRoleForm">确定</a-button>
          <a-button @click="roleOpen = false">取消</a-button>
        </a-space>
      </div>
    </a-modal>
  </div>
</template>

<script setup name="CustomerAccount">
import { listAccount, getAccount, addAccount, updateAccount, updateAccountRoles, resetAccountPwd, delAccount } from "@/api/customer/account";
import { customerOptions } from "@/api/customer/customer";
import { listPortalRoleOptions } from "@/api/customer/portalRole";

const route = useRoute();
const { proxy } = getCurrentInstance();

const statusOptions = [
  { label: "正常", value: "0" },
  { label: "停用", value: "1" }
];

const accountColumns = [
  { title: "账号", dataIndex: "username", key: "username", width: 140, align: "center" },
  { title: "客户名称", dataIndex: "customerName", key: "customerName", minWidth: 150, ellipsis: true, align: "center" },
  { title: "客户端角色", dataIndex: "roleNames", key: "roleNames", minWidth: 160, ellipsis: true, align: "center" },
  { title: "公司名称", dataIndex: "companyName", key: "companyName", minWidth: 170, ellipsis: true, align: "center" },
  { title: "姓名", dataIndex: "realName", key: "realName", width: 110, align: "center" },
  { title: "手机", dataIndex: "phone", key: "phone", width: 140, align: "center" },
  { title: "邮箱", dataIndex: "email", key: "email", minWidth: 160, ellipsis: true, align: "center" },
  { title: "主账号", dataIndex: "isMain", key: "isMain", width: 90, align: "center" },
  { title: "状态", dataIndex: "status", key: "status", width: 90, align: "center" },
  { title: "最后登录", dataIndex: "lastLoginTime", key: "lastLoginTime", width: 160, align: "center" },
  { title: "操作", key: "action", width: 260, align: "center", fixed: "right" }
];

const accountList = ref([]);
const customerOptionsList = ref([]);
const roleOptionsList = ref([]);
const open = ref(false);
const roleOpen = ref(false);
const loading = ref(true);
const showSearch = ref(true);
const ids = ref([]);
const selectedRowKeys = ref([]);
const single = ref(true);
const multiple = ref(true);
const total = ref(0);
const title = ref("");
const currentRoleAccount = ref({});

const accountRowSelection = computed(() => ({
  selectedRowKeys: selectedRowKeys.value,
  onChange: (keys, rows) => handleSelectionChange(rows, keys)
}));

const customerSelectOptions = computed(() => {
  return customerOptionsList.value.map(item => ({
    label: customerLabel(item),
    value: item.customerId
  }));
});

const data = reactive({
  form: {},
  roleForm: {
    accountId: undefined,
    roleIds: []
  },
  queryParams: {
    pageNum: 1,
    pageSize: 10,
    customerId: route.query.customerId,
    username: undefined,
    realName: undefined,
    status: undefined
  },
  rules: {
    customerId: [{ required: true, message: "客户不能为空", trigger: "change" }],
    username: [{ required: true, message: "账号不能为空", trigger: "blur" }],
    password: [
      { required: true, message: "密码不能为空", trigger: "blur" },
      { min: 5, max: 20, message: "密码长度必须介于 5 和 20 之间", trigger: "blur" }
    ],
    email: [{ type: "email", message: "请输入正确的邮箱地址", trigger: ["blur", "change"] }]
  }
});

const { queryParams, form, rules, roleForm } = toRefs(data);

function customerLabel(item) {
  return item.companyName ? `${item.customerName}（${item.companyName}）` : item.customerName;
}

function loadCustomerOptions(keyword = "") {
  customerOptions({ keyword }).then(response => {
    customerOptionsList.value = response.data || [];
  });
}

function ensureCustomerOption(row) {
  if (!row?.customerId) {
    return;
  }
  const exists = customerOptionsList.value.some(item => item.customerId == row.customerId);
  if (!exists) {
    customerOptionsList.value.unshift({
      customerId: row.customerId,
      customerNo: row.customerNo,
      customerName: row.customerName,
      companyName: row.companyName
    });
  }
}

function getList() {
  loading.value = true;
  listAccount(queryParams.value).then(response => {
    const result = response.data || {};
    accountList.value = result.rows || [];
    total.value = result.total || 0;
    loading.value = false;
  });
}

function reset() {
  form.value = {
    accountId: undefined,
    customerId: queryParams.value.customerId,
    username: undefined,
    password: undefined,
    realName: undefined,
    phone: undefined,
    email: undefined,
    isMain: "0",
    status: "0",
    remark: undefined
  };
  proxy.resetForm("accountRef");
}

function handleQuery() {
  queryParams.value.pageNum = 1;
  getList();
}

function resetQuery() {
  proxy.resetForm("queryRef");
  queryParams.value.customerId = route.query.customerId;
  handleQuery();
}

function handleSelectionChange(selection, keys) {
  ids.value = selection.map(item => item.accountId);
  selectedRowKeys.value = keys;
  single.value = selection.length !== 1;
  multiple.value = !selection.length;
}

function handleAdd() {
  reset();
  open.value = true;
  title.value = "新增客户账号";
}

function handleUpdate(row) {
  reset();
  const accountId = row?.accountId || ids.value[0];
  getAccount(accountId).then(response => {
    form.value = response.data.account;
    ensureCustomerOption(form.value);
    open.value = true;
    title.value = "修改客户账号";
  });
}

function loadRoleOptions() {
  listPortalRoleOptions().then(response => {
    roleOptionsList.value = response.data || [];
  });
}

function handleRoleAssign(row) {
  getAccount(row.accountId).then(response => {
    currentRoleAccount.value = response.data.account || {};
    roleForm.value.accountId = row.accountId;
    roleForm.value.roleIds = (response.data.roleIds || []).map(item => String(item));
    roleOpen.value = true;
  });
}

function submitRoleForm() {
  updateAccountRoles(roleForm.value.accountId, roleForm.value.roleIds).then(() => {
    proxy.$modal.msgSuccess("保存成功");
    roleOpen.value = false;
    getList();
  });
}

function submitForm() {
  proxy.$refs.accountRef.validate().then(() => {
    const request = form.value.accountId ? updateAccount(form.value) : addAccount(form.value);
    request.then(() => {
      proxy.$modal.msgSuccess("保存成功");
      open.value = false;
      getList();
    });
  }).catch(() => {});
}

function handleResetPwd(row) {
  proxy.$prompt(`请输入账号“${row.username}”的新密码`).then((value) => {
    if (!/^.{5,20}$/.test(value)) {
      proxy.$modal.msgError("密码长度必须介于 5 和 20 之间");
      return;
    }
    resetAccountPwd(row.accountId, value).then(() => {
      proxy.$modal.msgSuccess("密码已重置");
    });
  }).catch(() => {});
}

function handleDelete(row) {
  const accountIds = row?.accountId || ids.value;
  proxy.$modal.confirm(`是否确认删除客户账号编号为“${accountIds}”的数据项？`).then(() => {
    return delAccount(accountIds);
  }).then(() => {
    getList();
    proxy.$modal.msgSuccess("删除成功");
  }).catch(() => {});
}

function cancel() {
  open.value = false;
  reset();
}

loadCustomerOptions("");
loadRoleOptions();
if (route.query.customerId && route.query.customerName) {
  customerOptionsList.value = [{
    customerId: route.query.customerId,
    customerName: route.query.customerName,
    companyName: ""
  }];
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
