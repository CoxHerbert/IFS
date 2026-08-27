<template>
  <div class="app-container">
    <a-form
      v-show="showSearch"
      ref="queryRef"
      :model="queryParams"
      layout="inline"
      class="search-form"
    >
      <a-form-item label="联系人" name="contactName">
        <a-input
          v-model:value="queryParams.contactName"
          allow-clear
          placeholder="请输入联系人"
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
          style="width: 220px"
          @pressEnter="handleQuery"
        />
      </a-form-item>
      <a-form-item label="状态" name="status">
        <a-select
          v-model:value="queryParams.status"
          allow-clear
          placeholder="线索状态"
          style="width: 180px"
          :options="statusDictOptions"
        />
      </a-form-item>
      <a-form-item label="提交时间">
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
        <a-button :disabled="single" @click="handleUpdate()" v-hasPermi="['portal:contact:edit']">跟进</a-button>
        <a-button danger :disabled="multiple" @click="handleDelete()" v-hasPermi="['portal:contact:remove']">删除</a-button>
        <a-button @click="handleExport" v-hasPermi="['portal:contact:export']">导出</a-button>
      </a-space>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList" />
    </div>

    <vxe-table
      ref="contactTableRef"
      border
      stripe
      auto-resize
      show-overflow="title"
      :row-config="{ keyField: 'contactId' }"
      :loading="loading"
      :data="contactList"
      :checkbox-config="{ reserve: true }"
      @checkbox-change="handleSelectionChange"
      @checkbox-all="handleSelectionChange"
    >
      <vxe-column type="checkbox" width="55" align="center" />
      <vxe-column field="leadNo" title="线索编号" width="180" align="center" />
      <vxe-column field="contactName" title="联系人" width="120" align="center" />
      <vxe-column field="companyName" title="公司名称" min-width="150" align="center" />
      <vxe-column field="phone" title="联系电话" width="140" align="center" />
      <vxe-column field="email" title="邮箱" min-width="160" align="center" />
      <vxe-column field="route" title="目标航线" min-width="160" align="center" />
      <vxe-column field="status" title="状态" width="100" align="center">
        <template #default="{ row }">
          <dict-tag :options="portal_contact_status" :value="row.status" />
        </template>
      </vxe-column>
      <vxe-column field="source" title="来源" width="140" align="center" />
      <vxe-column field="createTime" title="提交时间" width="180" align="center">
        <template #default="{ row }">
          <span>{{ parseTime(row.createTime) }}</span>
        </template>
      </vxe-column>
      <vxe-column title="操作" width="190" align="center" fixed="right">
        <template #default="{ row }">
          <a-space>
            <a-button type="link" @click="handleView(row)">详情</a-button>
            <a-button type="link" @click="handleUpdate(row)" v-hasPermi="['portal:contact:edit']">跟进</a-button>
            <a-button type="link" danger @click="handleDelete(row)" v-hasPermi="['portal:contact:remove']">删除</a-button>
          </a-space>
        </template>
      </vxe-column>
    </vxe-table>

    <pagination
      v-show="total > 0"
      v-model:page="queryParams.pageNum"
      v-model:limit="queryParams.pageSize"
      :total="total"
      @pagination="getList"
    />

    <a-modal v-model:open="open" :title="title" width="720px" :footer="null" destroy-on-close>
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="线索编号">{{ form.leadNo }}</a-descriptions-item>
        <a-descriptions-item label="来源">{{ form.source }}</a-descriptions-item>
        <a-descriptions-item label="联系人">{{ form.contactName }}</a-descriptions-item>
        <a-descriptions-item label="公司名称">{{ form.companyName }}</a-descriptions-item>
        <a-descriptions-item label="联系电话">{{ form.phone }}</a-descriptions-item>
        <a-descriptions-item label="邮箱">{{ form.email }}</a-descriptions-item>
        <a-descriptions-item label="目标航线">{{ form.route }}</a-descriptions-item>
        <a-descriptions-item label="货物信息">{{ form.cargoInfo }}</a-descriptions-item>
        <a-descriptions-item label="提交 IP">{{ form.ipAddr }}</a-descriptions-item>
        <a-descriptions-item label="提交时间">{{ parseTime(form.createTime) }}</a-descriptions-item>
        <a-descriptions-item label="需求说明" :span="2">{{ form.message }}</a-descriptions-item>
      </a-descriptions>

      <a-form ref="contactRef" :model="form" :label-col="{ style: { width: '80px' } }" class="follow-form">
        <a-form-item label="跟进状态" name="status">
          <a-radio-group v-model:value="form.status" :disabled="readonly">
            <a-radio v-for="item in portal_contact_status" :key="item.value" :value="item.value">
              {{ item.label }}
            </a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="备注" name="remark">
          <a-textarea
            v-model:value="form.remark"
            :rows="4"
            placeholder="请输入跟进备注"
            :disabled="readonly"
          />
        </a-form-item>
      </a-form>

      <div class="modal-footer">
        <a-space>
          <a-button v-if="!readonly" type="primary" @click="submitForm">确定</a-button>
          <a-button @click="cancel">关闭</a-button>
        </a-space>
      </div>
    </a-modal>
  </div>
</template>

<script setup name="CustomerContact">
import { listContact, getContact, updateContact, delContact } from "@/api/customer/contact";

const { proxy } = getCurrentInstance();
const { portal_contact_status } = proxy.useDict("portal_contact_status");

const statusDictOptions = computed(() => {
  return (portal_contact_status.value || portal_contact_status || []).map(item => ({
    label: item.label,
    value: item.value
  }));
});

const contactList = ref([]);
const contactTableRef = ref();
const open = ref(false);
const readonly = ref(false);
const loading = ref(true);
const showSearch = ref(true);
const ids = ref([]);
const single = ref(true);
const multiple = ref(true);
const total = ref(0);
const title = ref("");
const dateRange = ref([]);

const data = reactive({
  form: {},
  queryParams: {
    pageNum: 1,
    pageSize: 10,
    contactName: undefined,
    companyName: undefined,
    phone: undefined,
    status: undefined
  }
});

const { queryParams, form } = toRefs(data);

function getList() {
  loading.value = true;
  listContact(proxy.addDateRange(queryParams.value, dateRange.value)).then(response => {
    const result = response.data || {};
    contactList.value = result.rows || [];
    total.value = result.total || 0;
    loading.value = false;
  });
}

function cancel() {
  open.value = false;
  reset();
}

function reset() {
  form.value = {
    contactId: undefined,
    leadNo: undefined,
    contactName: undefined,
    companyName: undefined,
    phone: undefined,
    email: undefined,
    route: undefined,
    cargoInfo: undefined,
    message: undefined,
    source: undefined,
    status: "10",
    ipAddr: undefined,
    createTime: undefined,
    remark: undefined
  };
}

function handleQuery() {
  queryParams.value.pageNum = 1;
  getList();
}

function resetQuery() {
  dateRange.value = [];
  proxy.resetForm("queryRef");
  contactTableRef.value?.clearCheckboxRow();
  handleQuery();
}

function handleSelectionChange() {
  const records = contactTableRef.value?.getCheckboxRecords() || [];
  ids.value = records.map(item => item.contactId);
  single.value = records.length !== 1;
  multiple.value = !records.length;
}

function handleView(row) {
  reset();
  readonly.value = true;
  getContact(row.contactId).then(response => {
    form.value = response.data;
    open.value = true;
    title.value = "线索详情";
  });
}

function handleUpdate(row) {
  reset();
  readonly.value = false;
  const contactId = row?.contactId || ids.value[0];
  getContact(contactId).then(response => {
    form.value = response.data;
    open.value = true;
    title.value = "线索跟进";
  });
}

function submitForm() {
  updateContact({
    contactId: form.value.contactId,
    status: form.value.status,
    remark: form.value.remark
  }).then(() => {
    proxy.$modal.msgSuccess("修改成功");
    open.value = false;
    getList();
  });
}

function handleDelete(row) {
  const contactIds = row?.contactId || ids.value;
  proxy.$modal.confirm(`是否确认删除线索编号为“${contactIds}”的数据项？`).then(() => {
    return delContact(contactIds);
  }).then(() => {
    getList();
    proxy.$modal.msgSuccess("删除成功");
  }).catch(() => {});
}

function handleExport() {
  proxy.download("portal/contact/export", {
    ...proxy.addDateRange(queryParams.value, dateRange.value)
  }, `portal_contact_${new Date().getTime()}.xlsx`);
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

.follow-form {
  margin-top: 18px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  margin-top: 24px;
}
</style>
