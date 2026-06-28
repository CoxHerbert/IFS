<template>
  <div class="app-container">
    <el-form :model="queryParams" ref="queryRef" :inline="true" v-show="showSearch" label-width="68px">
      <el-form-item label="客户名称" prop="customerName">
        <el-input v-model="queryParams.customerName" placeholder="请输入客户名称" clearable size="small" style="width: 220px" @keyup.enter="handleQuery" />
      </el-form-item>
      <el-form-item label="公司名称" prop="companyName">
        <el-input v-model="queryParams.companyName" placeholder="请输入公司名称" clearable size="small" style="width: 220px" @keyup.enter="handleQuery" />
      </el-form-item>
      <el-form-item label="联系电话" prop="phone">
        <el-input v-model="queryParams.phone" placeholder="请输入联系电话" clearable size="small" style="width: 180px" @keyup.enter="handleQuery" />
      </el-form-item>
      <el-form-item label="状态" prop="status">
        <el-select v-model="queryParams.status" placeholder="客户状态" clearable size="small" style="width: 140px">
          <el-option label="正常" value="0" />
          <el-option label="停用" value="1" />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" icon="Search" size="mini" @click="handleQuery">搜索</el-button>
        <el-button icon="Refresh" size="mini" @click="resetQuery">重置</el-button>
      </el-form-item>
    </el-form>

    <el-row :gutter="10" class="mb8">
      <el-col :span="1.5">
        <el-button type="primary" plain icon="Plus" size="mini" @click="handleAdd" v-hasPermi="['customer:customer:add']">新增</el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button type="success" plain icon="Edit" size="mini" :disabled="single" @click="handleUpdate" v-hasPermi="['customer:customer:edit']">修改</el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button type="danger" plain icon="Delete" size="mini" :disabled="multiple" @click="handleDelete" v-hasPermi="['customer:customer:remove']">删除</el-button>
      </el-col>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList"></right-toolbar>
    </el-row>

    <el-table v-loading="loading" :data="customerList" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="50" align="center" />
      <el-table-column label="客户编号" prop="customerNo" min-width="170" align="center" />
      <el-table-column label="客户名称" prop="customerName" min-width="150" align="center" :show-overflow-tooltip="true" />
      <el-table-column label="公司名称" prop="companyName" min-width="170" align="center" :show-overflow-tooltip="true" />
      <el-table-column label="默认联系人" prop="contactName" width="120" align="center" />
      <el-table-column label="联系电话" prop="phone" width="140" align="center" />
      <el-table-column label="邮箱" prop="email" min-width="160" align="center" :show-overflow-tooltip="true" />
      <el-table-column label="状态" prop="status" width="90" align="center">
        <template #default="scope">
          <el-tag :type="scope.row.status === '0' ? 'success' : 'info'">{{ scope.row.status === '0' ? '正常' : '停用' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" prop="createTime" width="160" align="center">
        <template #default="scope">
          <span>{{ parseTime(scope.row.createTime) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="280" align="center" fixed="right">
        <template #default="scope">
          <el-button size="mini" type="text" icon="Phone" @click="handleContacts(scope.row)">联系人</el-button>
          <el-button size="mini" type="text" icon="User" @click="handleAccounts(scope.row)">账号</el-button>
          <el-button size="mini" type="text" icon="Edit" @click="handleUpdate(scope.row)" v-hasPermi="['customer:customer:edit']">修改</el-button>
          <el-button size="mini" type="text" icon="Delete" @click="handleDelete(scope.row)" v-hasPermi="['customer:customer:remove']">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total > 0" :total="total" v-model:page="queryParams.pageNum" v-model:limit="queryParams.pageSize" @pagination="getList" />

    <el-dialog :title="title" v-model="open" width="640px" append-to-body>
      <el-form ref="customerRef" :model="form" :rules="rules" label-width="90px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="客户名称" prop="customerName">
              <el-input v-model="form.customerName" placeholder="请输入客户名称" maxlength="128" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="公司名称" prop="companyName">
              <el-input v-model="form.companyName" placeholder="请输入公司名称" maxlength="128" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="默认联系人" prop="contactName">
              <el-input v-model="form.contactName" placeholder="请输入默认联系人" maxlength="64" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="联系电话" prop="phone">
              <el-input v-model="form.phone" placeholder="请输入联系电话" maxlength="64" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="form.email" placeholder="请输入邮箱" maxlength="128" />
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

    <el-dialog :title="contactTitle" v-model="contactOpen" width="980px" append-to-body>
      <el-row :gutter="10" class="mb8">
        <el-col :span="1.5">
          <el-button type="primary" plain icon="Plus" size="mini" @click="handleContactAdd" v-hasPermi="['customer:customer:edit']">新增联系人</el-button>
        </el-col>
        <el-col :span="1.5">
          <el-button type="danger" plain icon="Delete" size="mini" :disabled="contactMultiple" @click="handleContactDelete" v-hasPermi="['customer:customer:remove']">删除</el-button>
        </el-col>
      </el-row>

      <el-table v-loading="contactLoading" :data="contactList" @selection-change="handleContactSelectionChange">
        <el-table-column type="selection" width="50" align="center" />
        <el-table-column label="姓名" prop="contactName" width="120" align="center" />
        <el-table-column label="职务" prop="position" width="120" align="center" />
        <el-table-column label="电话" prop="phone" width="140" align="center" />
        <el-table-column label="邮箱" prop="email" min-width="160" align="center" :show-overflow-tooltip="true" />
        <el-table-column label="微信" prop="wechat" width="130" align="center" />
        <el-table-column label="主要联系人" prop="isPrimary" width="110" align="center">
          <template #default="scope">
            <el-tag :type="scope.row.isPrimary === '1' ? 'warning' : 'info'">{{ scope.row.isPrimary === '1' ? '是' : '否' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" prop="status" width="90" align="center">
          <template #default="scope">
            <el-tag :type="scope.row.status === '0' ? 'success' : 'info'">{{ scope.row.status === '0' ? '正常' : '停用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="130" align="center" fixed="right">
          <template #default="scope">
            <el-button size="mini" type="text" icon="Edit" @click="handleContactUpdate(scope.row)" v-hasPermi="['customer:customer:edit']">修改</el-button>
            <el-button size="mini" type="text" icon="Delete" @click="handleContactDelete(scope.row)" v-hasPermi="['customer:customer:remove']">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <pagination v-show="contactTotal > 0" :total="contactTotal" v-model:page="contactQuery.pageNum" v-model:limit="contactQuery.pageSize" @pagination="getContactList" />
    </el-dialog>

    <el-dialog :title="contactFormTitle" v-model="contactFormOpen" width="620px" append-to-body>
      <el-form ref="contactRef" :model="contactForm" :rules="contactRules" label-width="90px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="姓名" prop="contactName">
              <el-input v-model="contactForm.contactName" placeholder="请输入联系人姓名" maxlength="64" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="职务" prop="position">
              <el-input v-model="contactForm.position" placeholder="请输入职务" maxlength="64" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="电话" prop="phone">
              <el-input v-model="contactForm.phone" placeholder="请输入电话" maxlength="64" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="contactForm.email" placeholder="请输入邮箱" maxlength="128" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="微信" prop="wechat">
              <el-input v-model="contactForm.wechat" placeholder="请输入微信" maxlength="64" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="主要联系人" prop="isPrimary">
              <el-radio-group v-model="contactForm.isPrimary">
                <el-radio label="0">否</el-radio>
                <el-radio label="1">是</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="contactForm.status">
            <el-radio label="0">正常</el-radio>
            <el-radio label="1">停用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="contactForm.remark" type="textarea" :rows="3" placeholder="请输入备注" maxlength="500" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button type="primary" @click="submitContactForm">确定</el-button>
          <el-button @click="contactFormOpen = false">取消</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup name="Customer">
import {
  listCustomer,
  getCustomer,
  addCustomer,
  updateCustomer,
  delCustomer,
  listCustomerContact,
  getCustomerContact,
  addCustomerContact,
  updateCustomerContact,
  delCustomerContact
} from "@/api/customer/customer";

const router = useRouter();
const { proxy } = getCurrentInstance();

const customerList = ref([]);
const open = ref(false);
const loading = ref(true);
const showSearch = ref(true);
const ids = ref([]);
const single = ref(true);
const multiple = ref(true);
const total = ref(0);
const title = ref("");

const currentCustomer = ref({});
const contactOpen = ref(false);
const contactLoading = ref(false);
const contactList = ref([]);
const contactIds = ref([]);
const contactMultiple = ref(true);
const contactTotal = ref(0);
const contactTitle = ref("");
const contactFormOpen = ref(false);
const contactFormTitle = ref("");

const data = reactive({
  form: {},
  queryParams: {
    pageNum: 1,
    pageSize: 10,
    customerName: undefined,
    companyName: undefined,
    phone: undefined,
    status: undefined
  },
  rules: {
    customerName: [{ required: true, message: "客户名称不能为空", trigger: "blur" }],
    email: [{ type: "email", message: "请输入正确的邮箱地址", trigger: ["blur", "change"] }]
  },
  contactQuery: {
    pageNum: 1,
    pageSize: 10
  },
  contactForm: {},
  contactRules: {
    contactName: [{ required: true, message: "联系人姓名不能为空", trigger: "blur" }],
    email: [{ type: "email", message: "请输入正确的邮箱地址", trigger: ["blur", "change"] }]
  }
});

const { queryParams, form, rules, contactQuery, contactForm, contactRules } = toRefs(data);

function getList() {
  loading.value = true;
  listCustomer(queryParams.value).then(response => {
    const data = response.data || {};
    customerList.value = data.rows || [];
    total.value = data.total || 0;
    loading.value = false;
  });
}

function reset() {
  form.value = {
    customerId: undefined,
    customerName: undefined,
    companyName: undefined,
    contactName: undefined,
    phone: undefined,
    email: undefined,
    status: "0",
    remark: undefined
  };
  proxy.resetForm("customerRef");
}

function resetContactForm() {
  contactForm.value = {
    contactId: undefined,
    customerId: currentCustomer.value.customerId,
    contactName: undefined,
    position: undefined,
    phone: undefined,
    email: undefined,
    wechat: undefined,
    isPrimary: "0",
    status: "0",
    remark: undefined
  };
  proxy.resetForm("contactRef");
}

function handleQuery() {
  queryParams.value.pageNum = 1;
  getList();
}

function resetQuery() {
  proxy.resetForm("queryRef");
  handleQuery();
}

function handleSelectionChange(selection) {
  ids.value = selection.map(item => item.customerId);
  single.value = selection.length != 1;
  multiple.value = !selection.length;
}

function handleAdd() {
  reset();
  open.value = true;
  title.value = "新增客户";
}

function handleUpdate(row) {
  reset();
  const customerId = row?.customerId || ids.value[0];
  getCustomer(customerId).then(response => {
    form.value = response.data;
    open.value = true;
    title.value = "修改客户";
  });
}

function handleAccounts(row) {
  router.push({ path: "/customer/customer-account", query: { customerId: row.customerId, customerName: row.customerName } });
}

function submitForm() {
  proxy.$refs["customerRef"].validate(valid => {
    if (valid) {
      const request = form.value.customerId ? updateCustomer(form.value) : addCustomer(form.value);
      request.then(() => {
        proxy.$modal.msgSuccess("保存成功");
        open.value = false;
        getList();
      });
    }
  });
}

function handleDelete(row) {
  const customerIds = row?.customerId || ids.value;
  proxy.$modal.confirm('是否确认删除客户编号为 "' + customerIds + '" 的数据项？').then(function() {
    return delCustomer(customerIds);
  }).then(() => {
    getList();
    proxy.$modal.msgSuccess("删除成功");
  }).catch(() => {});
}

function cancel() {
  open.value = false;
  reset();
}

function handleContacts(row) {
  currentCustomer.value = row;
  contactTitle.value = "联系人 - " + row.customerName;
  contactQuery.value.pageNum = 1;
  contactOpen.value = true;
  getContactList();
}

function getContactList() {
  contactLoading.value = true;
  listCustomerContact(currentCustomer.value.customerId, contactQuery.value).then(response => {
    const data = response.data || {};
    contactList.value = data.rows || [];
    contactTotal.value = data.total || 0;
    contactLoading.value = false;
  });
}

function handleContactSelectionChange(selection) {
  contactIds.value = selection.map(item => item.contactId);
  contactMultiple.value = !selection.length;
}

function handleContactAdd() {
  resetContactForm();
  contactFormOpen.value = true;
  contactFormTitle.value = "新增联系人";
}

function handleContactUpdate(row) {
  resetContactForm();
  getCustomerContact(row.contactId).then(response => {
    contactForm.value = response.data;
    contactFormOpen.value = true;
    contactFormTitle.value = "修改联系人";
  });
}

function submitContactForm() {
  proxy.$refs["contactRef"].validate(valid => {
    if (valid) {
      const request = contactForm.value.contactId ? updateCustomerContact(contactForm.value) : addCustomerContact(contactForm.value);
      request.then(() => {
        proxy.$modal.msgSuccess("保存成功");
        contactFormOpen.value = false;
        getContactList();
      });
    }
  });
}

function handleContactDelete(row) {
  const ids = row?.contactId || contactIds.value;
  proxy.$modal.confirm('是否确认删除联系人编号为 "' + ids + '" 的数据项？').then(function() {
    return delCustomerContact(ids);
  }).then(() => {
    getContactList();
    proxy.$modal.msgSuccess("删除成功");
  }).catch(() => {});
}

getList();
</script>
