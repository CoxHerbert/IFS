<template>
  <div class="app-container">
    <a-form
      v-show="showSearch"
      ref="queryRef"
      :model="queryParams"
      layout="inline"
      class="search-form"
    >
      <a-form-item label="客户名称" name="customerName">
        <a-input
          v-model:value="queryParams.customerName"
          allow-clear
          placeholder="请输入客户名称"
          style="width: 220px"
          @pressEnter="handleQuery"
        />
      </a-form-item>
      <a-form-item label="公司名称" name="companyName">
        <a-input
          v-model:value="queryParams.companyName"
          allow-clear
          placeholder="请输入公司名称"
          style="width: 220px"
          @pressEnter="handleQuery"
        />
      </a-form-item>
      <a-form-item label="联系电话" name="phone">
        <a-input
          v-model:value="queryParams.phone"
          allow-clear
          placeholder="请输入联系电话"
          style="width: 180px"
          @pressEnter="handleQuery"
        />
      </a-form-item>
      <a-form-item label="状态" name="status">
        <a-select
          v-model:value="queryParams.status"
          allow-clear
          placeholder="客户状态"
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

    <div class="toolbar-row mb8">
      <a-space>
        <a-button type="primary" @click="handleAdd" v-hasPermi="['customer:customer:add']">新增</a-button>
        <a-button :disabled="single" @click="handleUpdate()" v-hasPermi="['customer:customer:edit']">修改</a-button>
        <a-button danger :disabled="multiple" @click="handleDelete()" v-hasPermi="['customer:customer:remove']">删除</a-button>
      </a-space>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList" />
    </div>

    <a-table
      :loading="loading"
      :data-source="customerList"
      :columns="customerColumns"
      :pagination="false"
      :row-selection="customerRowSelection"
      :scroll="{ x: 1400 }"
      row-key="customerId"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <a-tag :color="record.status === '0' ? 'success' : 'default'">
            {{ record.status === '0' ? '正常' : '停用' }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'createTime'">
          {{ parseTime(record.createTime) }}
        </template>
        <template v-else-if="column.key === 'action'">
          <a-space>
            <a-button type="link" @click="handleContacts(record)">联系人</a-button>
            <a-button type="link" @click="handleAccounts(record)">账号</a-button>
            <a-button type="link" @click="handleUpdate(record)" v-hasPermi="['customer:customer:edit']">修改</a-button>
            <a-button type="link" danger @click="handleDelete(record)" v-hasPermi="['customer:customer:remove']">删除</a-button>
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
      <a-form ref="customerRef" :model="form" :rules="rules" :label-col="{ style: { width: '90px' } }">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="客户名称" name="customerName">
              <a-input v-model:value="form.customerName" placeholder="请输入客户名称" :maxlength="128" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="公司名称" name="companyName">
              <a-input v-model:value="form.companyName" placeholder="请输入公司名称" :maxlength="128" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="默认联系人" name="contactName">
              <a-input v-model:value="form.contactName" placeholder="请输入默认联系人" :maxlength="64" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="联系电话" name="phone">
              <a-input v-model:value="form.phone" placeholder="请输入联系电话" :maxlength="64" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="邮箱" name="email">
              <a-input v-model:value="form.email" placeholder="请输入邮箱" :maxlength="128" />
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

    <a-modal v-model:open="contactOpen" :title="contactTitle" width="980px" :footer="null" destroy-on-close>
      <div class="toolbar-row mb8">
        <a-space>
          <a-button type="primary" @click="handleContactAdd" v-hasPermi="['customer:customer:edit']">新增联系人</a-button>
          <a-button danger :disabled="contactMultiple" @click="handleContactDelete()" v-hasPermi="['customer:customer:remove']">删除</a-button>
        </a-space>
      </div>

      <a-table
        :loading="contactLoading"
        :data-source="contactList"
        :columns="contactColumns"
        :pagination="false"
        :row-selection="contactRowSelection"
        :scroll="{ x: 1100 }"
        row-key="contactId"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'isPrimary'">
            <a-tag :color="record.isPrimary === '1' ? 'gold' : 'default'">
              {{ record.isPrimary === '1' ? '是' : '否' }}
            </a-tag>
          </template>
          <template v-else-if="column.key === 'status'">
            <a-tag :color="record.status === '0' ? 'success' : 'default'">
              {{ record.status === '0' ? '正常' : '停用' }}
            </a-tag>
          </template>
          <template v-else-if="column.key === 'action'">
            <a-space>
              <a-button type="link" @click="handleContactUpdate(record)" v-hasPermi="['customer:customer:edit']">修改</a-button>
              <a-button type="link" danger @click="handleContactDelete(record)" v-hasPermi="['customer:customer:remove']">删除</a-button>
            </a-space>
          </template>
        </template>
      </a-table>

      <pagination
        v-show="contactTotal > 0"
        v-model:page="contactQuery.pageNum"
        v-model:limit="contactQuery.pageSize"
        :total="contactTotal"
        @pagination="getContactList"
      />
    </a-modal>

    <a-modal v-model:open="contactFormOpen" :title="contactFormTitle" width="620px" :footer="null" destroy-on-close>
      <a-form ref="contactRef" :model="contactForm" :rules="contactRules" :label-col="{ style: { width: '90px' } }">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="姓名" name="contactName">
              <a-input v-model:value="contactForm.contactName" placeholder="请输入联系人姓名" :maxlength="64" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="职务" name="position">
              <a-input v-model:value="contactForm.position" placeholder="请输入职务" :maxlength="64" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="电话" name="phone">
              <a-input v-model:value="contactForm.phone" placeholder="请输入电话" :maxlength="64" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="邮箱" name="email">
              <a-input v-model:value="contactForm.email" placeholder="请输入邮箱" :maxlength="128" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="微信" name="wechat">
              <a-input v-model:value="contactForm.wechat" placeholder="请输入微信" :maxlength="64" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="主要联系人" name="isPrimary">
              <a-radio-group v-model:value="contactForm.isPrimary">
                <a-radio value="0">否</a-radio>
                <a-radio value="1">是</a-radio>
              </a-radio-group>
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="状态" name="status">
          <a-radio-group v-model:value="contactForm.status">
            <a-radio value="0">正常</a-radio>
            <a-radio value="1">停用</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="备注" name="remark">
          <a-textarea v-model:value="contactForm.remark" placeholder="请输入备注" :rows="3" :maxlength="500" />
        </a-form-item>
      </a-form>
      <div class="modal-footer">
        <a-space>
          <a-button type="primary" @click="submitContactForm">确定</a-button>
          <a-button @click="contactFormOpen = false">取消</a-button>
        </a-space>
      </div>
    </a-modal>
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

const statusOptions = [
  { label: "正常", value: "0" },
  { label: "停用", value: "1" }
];

const customerColumns = [
  { title: "客户编号", dataIndex: "customerNo", key: "customerNo", width: 170, align: "center" },
  { title: "客户名称", dataIndex: "customerName", key: "customerName", minWidth: 150, ellipsis: true, align: "center" },
  { title: "公司名称", dataIndex: "companyName", key: "companyName", minWidth: 170, ellipsis: true, align: "center" },
  { title: "默认联系人", dataIndex: "contactName", key: "contactName", width: 120, align: "center" },
  { title: "联系电话", dataIndex: "phone", key: "phone", width: 140, align: "center" },
  { title: "邮箱", dataIndex: "email", key: "email", minWidth: 160, ellipsis: true, align: "center" },
  { title: "状态", dataIndex: "status", key: "status", width: 90, align: "center" },
  { title: "创建时间", dataIndex: "createTime", key: "createTime", width: 160, align: "center" },
  { title: "操作", key: "action", width: 280, align: "center", fixed: "right" }
];

const contactColumns = [
  { title: "姓名", dataIndex: "contactName", key: "contactName", width: 120, align: "center" },
  { title: "职务", dataIndex: "position", key: "position", width: 120, align: "center" },
  { title: "电话", dataIndex: "phone", key: "phone", width: 140, align: "center" },
  { title: "邮箱", dataIndex: "email", key: "email", minWidth: 160, ellipsis: true, align: "center" },
  { title: "微信", dataIndex: "wechat", key: "wechat", width: 130, align: "center" },
  { title: "主要联系人", dataIndex: "isPrimary", key: "isPrimary", width: 110, align: "center" },
  { title: "状态", dataIndex: "status", key: "status", width: 90, align: "center" },
  { title: "操作", key: "action", width: 130, align: "center", fixed: "right" }
];

const customerList = ref([]);
const open = ref(false);
const loading = ref(true);
const showSearch = ref(true);
const ids = ref([]);
const selectedRowKeys = ref([]);
const single = ref(true);
const multiple = ref(true);
const total = ref(0);
const title = ref("");

const currentCustomer = ref({});
const contactOpen = ref(false);
const contactLoading = ref(false);
const contactList = ref([]);
const contactIds = ref([]);
const selectedContactRowKeys = ref([]);
const contactMultiple = ref(true);
const contactTotal = ref(0);
const contactTitle = ref("");
const contactFormOpen = ref(false);
const contactFormTitle = ref("");

const customerRowSelection = computed(() => ({
  selectedRowKeys: selectedRowKeys.value,
  onChange: (keys, rows) => handleSelectionChange(rows, keys)
}));

const contactRowSelection = computed(() => ({
  selectedRowKeys: selectedContactRowKeys.value,
  onChange: (keys, rows) => handleContactSelectionChange(rows, keys)
}));

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
    const result = response.data || {};
    customerList.value = result.rows || [];
    total.value = result.total || 0;
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

function handleSelectionChange(selection, keys) {
  ids.value = selection.map(item => item.customerId);
  selectedRowKeys.value = keys;
  single.value = selection.length !== 1;
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
  proxy.$refs.customerRef.validate().then(() => {
    const request = form.value.customerId ? updateCustomer(form.value) : addCustomer(form.value);
    request.then(() => {
      proxy.$modal.msgSuccess("保存成功");
      open.value = false;
      getList();
    });
  }).catch(() => {});
}

function handleDelete(row) {
  const customerIds = row?.customerId || ids.value;
  proxy.$modal.confirm(`是否确认删除客户编号为“${customerIds}”的数据项？`).then(() => {
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
  contactTitle.value = `联系人 - ${row.customerName}`;
  contactQuery.value.pageNum = 1;
  selectedContactRowKeys.value = [];
  contactIds.value = [];
  contactOpen.value = true;
  getContactList();
}

function getContactList() {
  contactLoading.value = true;
  listCustomerContact(currentCustomer.value.customerId, contactQuery.value).then(response => {
    const result = response.data || {};
    contactList.value = result.rows || [];
    contactTotal.value = result.total || 0;
    contactLoading.value = false;
  });
}

function handleContactSelectionChange(selection, keys) {
  contactIds.value = selection.map(item => item.contactId);
  selectedContactRowKeys.value = keys;
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
  proxy.$refs.contactRef.validate().then(() => {
    const request = contactForm.value.contactId ? updateCustomerContact(contactForm.value) : addCustomerContact(contactForm.value);
    request.then(() => {
      proxy.$modal.msgSuccess("保存成功");
      contactFormOpen.value = false;
      getContactList();
    });
  }).catch(() => {});
}

function handleContactDelete(row) {
  const targetIds = row?.contactId || contactIds.value;
  proxy.$modal.confirm(`是否确认删除联系人编号为“${targetIds}”的数据项？`).then(() => {
    return delCustomerContact(targetIds);
  }).then(() => {
    getContactList();
    proxy.$modal.msgSuccess("删除成功");
  }).catch(() => {});
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
