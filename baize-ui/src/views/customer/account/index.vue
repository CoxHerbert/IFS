<template>
  <div class="app-container">
    <el-form :model="queryParams" ref="queryRef" :inline="true" v-show="showSearch" label-width="68px">
      <el-form-item label="客户" prop="customerId">
        <el-select
          v-model="queryParams.customerId"
          filterable
          remote
          clearable
          reserve-keyword
          placeholder="请选择客户"
          :remote-method="loadCustomerOptions"
          style="width: 260px"
        >
          <el-option v-for="item in customerOptionsList" :key="item.customerId" :label="customerLabel(item)" :value="item.customerId" />
        </el-select>
      </el-form-item>
      <el-form-item label="账号" prop="username">
        <el-input v-model="queryParams.username" placeholder="请输入账号" clearable size="small" style="width: 180px" @keyup.enter="handleQuery" />
      </el-form-item>
      <el-form-item label="姓名" prop="realName">
        <el-input v-model="queryParams.realName" placeholder="请输入姓名" clearable size="small" style="width: 160px" @keyup.enter="handleQuery" />
      </el-form-item>
      <el-form-item label="状态" prop="status">
        <el-select v-model="queryParams.status" placeholder="账号状态" clearable size="small" style="width: 140px">
          <el-option label="正常" value="0" />
          <el-option label="停用" value="1" />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" icon="Search" size="mini" @click="handleQuery">搜索</el-button>
        <el-button icon="Refresh" size="mini" @click="resetQuery">重置</el-button>
      </el-form-item>
    </el-form>

    <el-alert v-if="route.query.customerName" :title="'当前客户：' + route.query.customerName" type="info" show-icon class="mb8" />

    <el-row :gutter="10" class="mb8">
      <el-col :span="1.5">
        <el-button type="primary" plain icon="Plus" size="mini" @click="handleAdd" v-hasPermi="['customer:account:add']">新增</el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button type="success" plain icon="Edit" size="mini" :disabled="single" @click="handleUpdate" v-hasPermi="['customer:account:edit']">修改</el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button type="danger" plain icon="Delete" size="mini" :disabled="multiple" @click="handleDelete" v-hasPermi="['customer:account:remove']">删除</el-button>
      </el-col>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList"></right-toolbar>
    </el-row>

    <el-table v-loading="loading" :data="accountList" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="50" align="center" />
      <el-table-column label="账号" prop="username" min-width="140" align="center" />
      <el-table-column label="客户名称" prop="customerName" min-width="150" align="center" :show-overflow-tooltip="true" />
      <el-table-column label="公司名称" prop="companyName" min-width="170" align="center" :show-overflow-tooltip="true" />
      <el-table-column label="姓名" prop="realName" width="110" align="center" />
      <el-table-column label="手机" prop="phone" width="140" align="center" />
      <el-table-column label="邮箱" prop="email" min-width="160" align="center" :show-overflow-tooltip="true" />
      <el-table-column label="主账号" prop="isMain" width="90" align="center">
        <template #default="scope">
          <el-tag :type="scope.row.isMain === '1' ? 'warning' : 'info'">{{ scope.row.isMain === '1' ? '是' : '否' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="状态" prop="status" width="90" align="center">
        <template #default="scope">
          <el-tag :type="scope.row.status === '0' ? 'success' : 'info'">{{ scope.row.status === '0' ? '正常' : '停用' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="最后登录" prop="lastLoginTime" width="160" align="center">
        <template #default="scope">
          <span>{{ parseTime(scope.row.lastLoginTime) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="230" align="center" fixed="right">
        <template #default="scope">
          <el-button size="mini" type="text" icon="Edit" @click="handleUpdate(scope.row)" v-hasPermi="['customer:account:edit']">修改</el-button>
          <el-button size="mini" type="text" icon="Key" @click="handleResetPwd(scope.row)" v-hasPermi="['customer:account:resetPwd']">重置密码</el-button>
          <el-button size="mini" type="text" icon="Delete" @click="handleDelete(scope.row)" v-hasPermi="['customer:account:remove']">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total > 0" :total="total" v-model:page="queryParams.pageNum" v-model:limit="queryParams.pageSize" @pagination="getList" />

    <el-dialog :title="title" v-model="open" width="640px" append-to-body>
      <el-form ref="accountRef" :model="form" :rules="rules" label-width="90px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="客户" prop="customerId">
              <el-select
                v-model="form.customerId"
                filterable
                remote
                reserve-keyword
                placeholder="请选择客户"
                :remote-method="loadCustomerOptions"
                :disabled="!!form.accountId"
                style="width: 100%"
              >
                <el-option v-for="item in customerOptionsList" :key="item.customerId" :label="customerLabel(item)" :value="item.customerId" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="账号" prop="username">
              <el-input v-model="form.username" placeholder="请输入登录账号" maxlength="64" :disabled="!!form.accountId" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item v-if="!form.accountId" label="密码" prop="password">
              <el-input v-model="form.password" placeholder="请输入登录密码" type="password" show-password maxlength="20" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="姓名" prop="realName">
              <el-input v-model="form.realName" placeholder="请输入姓名" maxlength="64" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="手机" prop="phone">
              <el-input v-model="form.phone" placeholder="请输入手机号" maxlength="64" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="form.email" placeholder="请输入邮箱" maxlength="128" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="主账号" prop="isMain">
              <el-radio-group v-model="form.isMain">
                <el-radio label="0">否</el-radio>
                <el-radio label="1">是</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-radio-group v-model="form.status">
                <el-radio label="0">正常</el-radio>
                <el-radio label="1">停用</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="form.remark" type="textarea" :rows="3" placeholder="请输入备注" maxlength="500" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button type="primary" @click="submitForm">确定</el-button>
          <el-button @click="cancel">取消</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup name="CustomerAccount">
import { listAccount, getAccount, addAccount, updateAccount, resetAccountPwd, delAccount } from "@/api/customer/account";
import { customerOptions } from "@/api/customer/customer";

const route = useRoute();
const { proxy } = getCurrentInstance();

const accountList = ref([]);
const customerOptionsList = ref([]);
const open = ref(false);
const loading = ref(true);
const showSearch = ref(true);
const ids = ref([]);
const single = ref(true);
const multiple = ref(true);
const total = ref(0);
const title = ref("");

const data = reactive({
  form: {},
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
    password: [{ required: true, message: "密码不能为空", trigger: "blur" }, { min: 5, max: 20, message: "密码长度必须介于 5 和 20 之间", trigger: "blur" }],
    email: [{ type: "email", message: "请输入正确的邮箱地址", trigger: ["blur", "change"] }]
  }
});

const { queryParams, form, rules } = toRefs(data);

function customerLabel(item) {
  return item.companyName ? `${item.customerName}（${item.companyName}）` : item.customerName;
}

function loadCustomerOptions(keyword) {
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
    const data = response.data || {};
    accountList.value = data.rows || [];
    total.value = data.total || 0;
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

function handleSelectionChange(selection) {
  ids.value = selection.map(item => item.accountId);
  single.value = selection.length != 1;
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
    form.value = response.data;
    ensureCustomerOption(form.value);
    open.value = true;
    title.value = "修改客户账号";
  });
}

function submitForm() {
  proxy.$refs["accountRef"].validate(valid => {
    if (valid) {
      const request = form.value.accountId ? updateAccount(form.value) : addAccount(form.value);
      request.then(() => {
        proxy.$modal.msgSuccess("保存成功");
        open.value = false;
        getList();
      });
    }
  });
}

function handleResetPwd(row) {
  proxy.$prompt('请输入账号 "' + row.username + '" 的新密码', "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    closeOnClickModal: false,
    inputPattern: /^.{5,20}$/,
    inputErrorMessage: "密码长度必须介于 5 和 20 之间"
  }).then(({ value }) => {
    resetAccountPwd(row.accountId, value).then(() => {
      proxy.$modal.msgSuccess("密码已重置");
    });
  }).catch(() => {});
}

function handleDelete(row) {
  const accountIds = row?.accountId || ids.value;
  proxy.$modal.confirm('是否确认删除客户账号编号为 "' + accountIds + '" 的数据项？').then(function() {
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
if (route.query.customerId && route.query.customerName) {
  customerOptionsList.value = [{
    customerId: route.query.customerId,
    customerName: route.query.customerName,
    companyName: ""
  }];
}
getList();
</script>
