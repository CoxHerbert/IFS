<template>
  <div class="app-container">
    <el-form :model="queryParams" ref="queryRef" :inline="true" v-show="showSearch" label-width="68px">
      <el-form-item label="联系人" prop="contactName">
        <el-input
          v-model="queryParams.contactName"
          placeholder="请输入联系人"
          clearable
          size="small"
          style="width: 220px"
          @keyup.enter="handleQuery"
        />
      </el-form-item>
      <el-form-item label="公司名称" prop="companyName">
        <el-input
          v-model="queryParams.companyName"
          placeholder="请输入公司名称"
          clearable
          size="small"
          style="width: 220px"
          @keyup.enter="handleQuery"
        />
      </el-form-item>
      <el-form-item label="联系电话" prop="phone">
        <el-input
          v-model="queryParams.phone"
          placeholder="请输入联系电话"
          clearable
          size="small"
          style="width: 220px"
          @keyup.enter="handleQuery"
        />
      </el-form-item>
      <el-form-item label="状态" prop="status">
        <el-select v-model="queryParams.status" placeholder="线索状态" clearable size="small" style="width: 180px">
          <el-option v-for="item in portal_contact_status" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
      </el-form-item>
      <el-form-item label="提交时间">
        <el-date-picker
          v-model="dateRange"
          size="small"
          style="width: 240px"
          value-format="YYYY-MM-DD"
          type="daterange"
          range-separator="-"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
        />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" icon="Search" size="mini" @click="handleQuery">搜索</el-button>
        <el-button icon="Refresh" size="mini" @click="resetQuery">重置</el-button>
      </el-form-item>
    </el-form>

    <el-row :gutter="10" class="mb8">
      <el-col :span="1.5">
        <el-button
          type="success"
          plain
          icon="Edit"
          size="mini"
          :disabled="single"
          @click="handleUpdate"
          v-hasPermi="['portal:contact:edit']"
        >跟进</el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button
          type="danger"
          plain
          icon="Delete"
          size="mini"
          :disabled="multiple"
          @click="handleDelete"
          v-hasPermi="['portal:contact:remove']"
        >删除</el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button
          type="warning"
          plain
          icon="Download"
          size="mini"
          @click="handleExport"
          v-hasPermi="['portal:contact:export']"
        >导出</el-button>
      </el-col>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList"></right-toolbar>
    </el-row>

    <vxe-table
      ref="contactTableRef"
      border
      stripe
      auto-resize
      show-overflow
      row-id="contactId"
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
      <vxe-column title="操作" width="170" align="center" fixed="right">
        <template #default="{ row }">
          <el-button size="mini" type="text" icon="View" @click="handleView(row)">详情</el-button>
          <el-button
            size="mini"
            type="text"
            icon="Edit"
            @click="handleUpdate(row)"
            v-hasPermi="['portal:contact:edit']"
          >跟进</el-button>
          <el-button
            size="mini"
            type="text"
            icon="Delete"
            @click="handleDelete(row)"
            v-hasPermi="['portal:contact:remove']"
          >删除</el-button>
        </template>
      </vxe-column>
    </vxe-table>

    <pagination
      v-show="total > 0"
      :total="total"
      v-model:page="queryParams.pageNum"
      v-model:limit="queryParams.pageSize"
      @pagination="getList"
    />

    <el-dialog :title="title" v-model="open" width="720px" append-to-body>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="线索编号">{{ form.leadNo }}</el-descriptions-item>
        <el-descriptions-item label="来源">{{ form.source }}</el-descriptions-item>
        <el-descriptions-item label="联系人">{{ form.contactName }}</el-descriptions-item>
        <el-descriptions-item label="公司名称">{{ form.companyName }}</el-descriptions-item>
        <el-descriptions-item label="联系电话">{{ form.phone }}</el-descriptions-item>
        <el-descriptions-item label="邮箱">{{ form.email }}</el-descriptions-item>
        <el-descriptions-item label="目标航线">{{ form.route }}</el-descriptions-item>
        <el-descriptions-item label="货物信息">{{ form.cargoInfo }}</el-descriptions-item>
        <el-descriptions-item label="提交IP">{{ form.ipAddr }}</el-descriptions-item>
        <el-descriptions-item label="提交时间">{{ parseTime(form.createTime) }}</el-descriptions-item>
        <el-descriptions-item label="需求说明" :span="2">{{ form.message }}</el-descriptions-item>
      </el-descriptions>

      <el-form ref="contactRef" :model="form" label-width="80px" class="follow-form">
        <el-form-item label="跟进状态" prop="status">
          <el-radio-group v-model="form.status" :disabled="readonly">
            <el-radio v-for="item in portal_contact_status" :key="item.value" :label="item.value">{{ item.label }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input
            v-model="form.remark"
            type="textarea"
            :rows="4"
            placeholder="请输入跟进备注"
            :disabled="readonly"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <div class="dialog-footer">
          <el-button v-if="!readonly" type="primary" @click="submitForm">确 定</el-button>
          <el-button @click="cancel">关 闭</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup name="CustomerContact">
import { listContact, getContact, updateContact, delContact } from "@/api/customer/contact";

const { proxy } = getCurrentInstance();
const { portal_contact_status } = proxy.useDict("portal_contact_status");

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
    const data = response.data || {};
    contactList.value = data.rows || [];
    total.value = data.total || 0;
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
  single.value = records.length != 1;
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
  proxy.$modal.confirm('是否确认删除线索编号为"' + contactIds + '"的数据项？').then(function() {
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
.follow-form {
  margin-top: 18px;
}
</style>
