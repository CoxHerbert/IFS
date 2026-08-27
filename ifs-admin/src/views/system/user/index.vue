<template>
  <div class="app-container">
    <div class="user-page-layout">
      <aside class="user-page-aside">
        <div class="panel-shell">
          <div class="panel-title">部门筛选</div>
          <div class="head-container">
            <vxe-input v-model="deptName" clearable placeholder="请输入部门名称" prefix-icon="vxe-icon-search" />
          </div>
          <div class="dept-tree-wrap">
            <a-tree
              :tree-data="filteredDeptOptions"
              :field-names="{ title: 'label', key: 'id', children: 'children' }"
              :expanded-keys="deptExpandedKeys"
              :auto-expand-parent="true"
              @expand="keys => (deptExpandedKeys = keys)"
              @select="handleNodeSelect"
            />
          </div>
        </div>
      </aside>

      <section class="user-page-main">
        <div class="panel-shell">
          <vxe-form :data="queryParams" @submit="handleQuery" @reset="resetQuery">
            <vxe-form-item title="用户名称" field="userName" span="6" :item-render="{}">
              <template #default>
                <vxe-input v-model="queryParams.userName" clearable placeholder="请输入用户名称" />
              </template>
            </vxe-form-item>
            <vxe-form-item title="手机号码" field="phonenumber" span="6" :item-render="{}">
              <template #default>
                <vxe-input v-model="queryParams.phonenumber" clearable placeholder="请输入手机号码" />
              </template>
            </vxe-form-item>
            <vxe-form-item title="状态" field="status" span="6" :item-render="{}">
              <template #default>
                <vxe-select v-model="queryParams.status" clearable placeholder="请选择状态">
                  <vxe-option v-for="item in sysNormalDisableOptions" :key="item.value" :value="item.value" :label="item.label" />
                </vxe-select>
              </template>
            </vxe-form-item>
            <vxe-form-item title="创建时间" field="dateRange" span="6" :item-render="{}">
              <template #default>
                <vxe-date-range-picker v-model="dateRange" value-format="YYYY-MM-DD" clearable />
              </template>
            </vxe-form-item>
            <vxe-form-item span="24" align="left" :item-render="{}" class-name="search-actions">
              <template #default>
                <vxe-button type="submit" status="primary">搜索</vxe-button>
                <vxe-button type="reset">重置</vxe-button>
                <vxe-button @click="showSearch = !showSearch">{{ showSearch ? "隐藏搜索" : "显示搜索" }}</vxe-button>
              </template>
            </vxe-form-item>
          </vxe-form>

          <div v-show="showSearch" class="toolbar-spacer"></div>

          <vxe-toolbar class="user-toolbar">
            <template #buttons>
              <vxe-button status="primary" @click="handleAdd" v-hasPermi="['system:user:add']">新增</vxe-button>
              <vxe-button :disabled="single" @click="handleUpdate()" v-hasPermi="['system:user:edit']">修改</vxe-button>
              <vxe-button status="error" :disabled="multiple" @click="handleDelete()" v-hasPermi="['system:user:remove']">
                删除
              </vxe-button>
              <vxe-button @click="handleImport" v-hasPermi="['system:user:import']">导入</vxe-button>
              <vxe-button @click="handleExport" v-hasPermi="['system:user:export']">导出</vxe-button>
            </template>
            <template #tools>
              <right-toolbar v-model:showSearch="showSearch" @queryTable="getList" :columns="columns" />
            </template>
          </vxe-toolbar>

          <vxe-table
            ref="userTableRef"
            border
            stripe
            auto-resize
            show-overflow="title"
            :loading="loading"
            :data="userList"
            :row-config="{ keyField: 'userId' }"
            :checkbox-config="{ reserve: true }"
            @checkbox-change="handleSelectionChange"
            @checkbox-all="handleSelectionChange"
          >
            <vxe-column type="checkbox" width="54" align="center" />
            <vxe-column v-if="columns[0].visible" field="userId" title="用户编号" width="90" align="center" />
            <vxe-column v-if="columns[1].visible" field="userName" title="用户名称" min-width="140" align="center" />
            <vxe-column v-if="columns[2].visible" field="nickName" title="用户昵称" min-width="140" align="center" />
            <vxe-column v-if="columns[3].visible" field="deptName" title="部门" min-width="160" align="center" />
            <vxe-column v-if="columns[4].visible" field="phonenumber" title="手机号码" width="140" align="center" />
            <vxe-column v-if="columns[5].visible" field="status" title="状态" width="110" align="center">
              <template #default="{ row }">
                <a-switch
                  v-model:checked="row.status"
                  checked-value="0"
                  un-checked-value="1"
                  @change="() => handleStatusChange(row)"
                />
              </template>
            </vxe-column>
            <vxe-column v-if="columns[6].visible" field="createTime" title="创建时间" width="170" align="center">
              <template #default="{ row }">
                <span>{{ parseTime(row.createTime) }}</span>
              </template>
            </vxe-column>
            <vxe-column title="操作" width="220" align="center" fixed="right">
              <template #default="{ row }">
                <div v-if="row.userId !== 1" class="action-links">
                  <vxe-button mode="text" status="primary" @click="handleUpdate(row)" v-hasPermi="['system:user:edit']">
                    修改
                  </vxe-button>
                  <vxe-button mode="text" status="error" @click="handleDelete(row)" v-hasPermi="['system:user:remove']">
                    删除
                  </vxe-button>
                  <span v-hasPermi="['system:user:resetPwd', 'system:user:edit']">
                    <a-dropdown>
                    <a class="more-link" @click.prevent>更多</a>
                    <template #overlay>
                      <a-menu @click="({ key }) => handleCommand(key, row)">
                        <a-menu-item key="handleResetPwd" v-hasPermi="['system:user:resetPwd']">重置密码</a-menu-item>
                        <a-menu-item key="handleAuthRole" v-hasPermi="['system:user:edit']">分配角色</a-menu-item>
                      </a-menu>
                    </template>
                    </a-dropdown>
                  </span>
                </div>
              </template>
            </vxe-column>
          </vxe-table>

          <div class="pager-wrap">
            <vxe-pager
              :current-page="queryParams.pageNum"
              :page-size="queryParams.pageSize"
              :total="total"
              :page-sizes="[10, 20, 30, 50]"
              :layouts="['PrevPage', 'JumpNumber', 'NextPage', 'Sizes', 'FullJump', 'Total']"
              @page-change="handlePageChange"
            />
          </div>
        </div>
      </section>
    </div>

    <vxe-modal v-model="open" :title="title" width="720" show-footer esc-closable mask-closable="false">
      <vxe-form ref="userFormRef" :data="form" :rules="rules" title-width="90" @submit="submitForm">
        <vxe-form-item title="用户昵称" field="nickName" span="12" :item-render="{}">
          <template #default>
            <vxe-input v-model="form.nickName" maxlength="30" />
          </template>
        </vxe-form-item>
        <vxe-form-item title="归属部门" field="deptId" span="12" :item-render="{}">
          <template #default>
            <vxe-tree-select
              v-model="form.deptId"
              :options="deptOptions"
              :option-props="{ value: 'id', label: 'label', children: 'children' }"
              clearable
              transfer
            />
          </template>
        </vxe-form-item>
        <vxe-form-item title="手机号码" field="phonenumber" span="12" :item-render="{}">
          <template #default>
            <vxe-input v-model="form.phonenumber" maxlength="11" />
          </template>
        </vxe-form-item>
        <vxe-form-item title="邮箱" field="email" span="12" :item-render="{}">
          <template #default>
            <vxe-input v-model="form.email" maxlength="50" />
          </template>
        </vxe-form-item>
        <vxe-form-item v-if="form.userId === undefined" title="用户名称" field="userName" span="12" :item-render="{}">
          <template #default>
            <vxe-input v-model="form.userName" maxlength="30" />
          </template>
        </vxe-form-item>
        <vxe-form-item v-if="form.userId === undefined" title="用户密码" field="password" span="12" :item-render="{}">
          <template #default>
            <vxe-input v-model="form.password" type="password" maxlength="20" />
          </template>
        </vxe-form-item>
        <vxe-form-item title="用户性别" field="sex" span="12" :item-render="{}">
          <template #default>
            <vxe-select v-model="form.sex">
              <vxe-option v-for="item in sysUserSexOptions" :key="item.value" :value="item.value" :label="item.label" />
            </vxe-select>
          </template>
        </vxe-form-item>
        <vxe-form-item title="状态" field="status" span="12" :item-render="{}">
          <template #default>
            <vxe-radio-group v-model="form.status">
              <vxe-radio v-for="item in sysNormalDisableOptions" :key="item.value" :label="item.value" :content="item.label" />
            </vxe-radio-group>
          </template>
        </vxe-form-item>
        <vxe-form-item title="岗位" field="postIds" span="12" :item-render="{}">
          <template #default>
            <vxe-select v-model="form.postIds" multiple clearable>
              <vxe-option
                v-for="item in postOptions"
                :key="item.postId"
                :value="item.postId"
                :label="item.postName"
                :disabled="item.status === 1"
              />
            </vxe-select>
          </template>
        </vxe-form-item>
        <vxe-form-item title="角色" field="roleIds" span="12" :item-render="{}">
          <template #default>
            <vxe-select v-model="form.roleIds" multiple clearable>
              <vxe-option
                v-for="item in roleOptions"
                :key="item.roleId"
                :value="item.roleId"
                :label="item.roleName"
                :disabled="item.status === 1"
              />
            </vxe-select>
          </template>
        </vxe-form-item>
        <vxe-form-item title="备注" field="remark" span="24" :item-render="{}">
          <template #default>
            <vxe-textarea v-model="form.remark" rows="3" />
          </template>
        </vxe-form-item>
      </vxe-form>
      <template #footer>
        <div class="modal-footer">
          <vxe-button @click="cancel">取消</vxe-button>
          <vxe-button status="primary" @click="submitForm">确定</vxe-button>
        </div>
      </template>
    </vxe-modal>

    <a-modal
      v-model:open="upload.open"
      :title="upload.title"
      width="400px"
      :mask-closable="false"
      @ok="submitFileForm"
      @cancel="upload.open = false"
    >
      <a-upload-dragger
        ref="uploadRef"
        v-model:file-list="uploadFileList"
        name="file"
        :max-count="1"
        accept=".xlsx,.xls"
        :headers="upload.headers"
        :action="upload.url + '?updateSupport=' + upload.updateSupport"
        :disabled="upload.isUploading"
        :before-upload="beforeUpload"
        @change="handleUploadChange"
      >
        <p class="ant-upload-drag-icon">
          <InboxOutlined />
        </p>
        <p class="ant-upload-text">将文件拖到此处，或点击上传</p>
        <p class="ant-upload-hint">仅允许导入 xls、xlsx 格式文件。</p>
      </a-upload-dragger>

      <div class="upload-tips">
        <span>仅允许导入 xls、xlsx 格式文件。</span>
        <a-button type="link" size="small" @click="importTemplate">下载模板</a-button>
      </div>

      <template #footer>
        <a-space>
          <a-button @click="upload.open = false">取消</a-button>
          <a-button type="primary" :loading="upload.isUploading" @click="submitFileForm">确定</a-button>
        </a-space>
      </template>
    </a-modal>
  </div>
</template>

<script setup name="User">
import { InboxOutlined } from "@ant-design/icons-vue";
import { getToken } from "@/utils/auth";
import { listDept } from "@/api/system/dept";
import {
  addUser,
  changeUserStatus,
  delUser,
  getUser,
  listUser,
  resetUserPwd,
  updateUser
} from "@/api/system/user";

const router = useRouter();
const { proxy } = getCurrentInstance();
const { sys_normal_disable, sys_user_sex } = proxy.useDict("sys_normal_disable", "sys_user_sex");

const uploadRef = ref();
const userTableRef = ref();
const userFormRef = ref();

const loading = ref(false);
const showSearch = ref(true);
const open = ref(false);
const title = ref("");
const total = ref(0);
const deptName = ref("");
const dateRange = ref([]);
const userList = ref([]);
const ids = ref([]);
const single = ref(true);
const multiple = ref(true);
const deptOptions = ref([]);
const deptExpandedKeys = ref([]);
const postOptions = ref([]);
const roleOptions = ref([]);
const uploadFileList = ref([]);
const initPassword = ref(undefined);

const columns = ref([
  { key: 0, label: "用户编号", visible: true },
  { key: 1, label: "用户名称", visible: true },
  { key: 2, label: "用户昵称", visible: true },
  { key: 3, label: "部门", visible: true },
  { key: 4, label: "手机号码", visible: true },
  { key: 5, label: "状态", visible: true },
  { key: 6, label: "创建时间", visible: true }
]);

const queryParams = reactive({
  pageNum: 1,
  pageSize: 10,
  userName: undefined,
  phonenumber: undefined,
  status: undefined,
  deptId: undefined
});

const form = reactive({
  userId: undefined,
  deptId: undefined,
  userName: undefined,
  nickName: undefined,
  password: undefined,
  phonenumber: undefined,
  email: undefined,
  sex: "2",
  status: "0",
  remark: undefined,
  postIds: [],
  roleIds: []
});

const rules = {
  userName: [
    { required: true, message: "用户名称不能为空" },
    { min: 2, max: 20, message: "用户名称长度必须介于 2 和 20 之间" }
  ],
  nickName: [{ required: true, message: "用户昵称不能为空" }],
  password: [
    { required: true, message: "用户密码不能为空" },
    { min: 5, max: 20, message: "用户密码长度必须介于 5 和 20 之间" }
  ],
  email: [{ pattern: /^[^\s@]+@[^\s@]+\.[^\s@]+$/, message: "请输入正确的邮箱地址" }],
  phonenumber: [{ pattern: /^1[3-9][0-9]\d{8}$/, message: "请输入正确的手机号码" }]
};

const upload = reactive({
  open: false,
  title: "",
  isUploading: false,
  updateSupport: 0,
  headers: {
    Authorization: "Bearer " + getToken()
  },
  url: import.meta.env.VITE_APP_BASE_API + "system/user/importData"
});

const sysNormalDisableOptions = computed(() => sys_normal_disable.value || []);
const sysUserSexOptions = computed(() => sys_user_sex.value || []);

const filteredDeptOptions = computed(() => {
  if (!deptName.value) {
    return deptOptions.value;
  }
  return filterTreeByLabel(deptOptions.value, deptName.value);
});

function filterTreeByLabel(tree = [], keyword) {
  return tree
    .map(item => {
      const children = item.children?.length ? filterTreeByLabel(item.children, keyword) : [];
      if (item.label?.includes(keyword) || children.length) {
        return {
          ...item,
          children
        };
      }
      return null;
    })
    .filter(Boolean);
}

function getAllTreeKeys(tree = []) {
  const keys = [];
  const loop = list => {
    list.forEach(item => {
      keys.push(item.id);
      if (item.children?.length) {
        loop(item.children);
      }
    });
  };
  loop(tree);
  return keys;
}

function getTreeselect() {
  listDept().then(response => {
    deptOptions.value = proxy.handleProps(response.data, "deptId", "deptName");
    deptExpandedKeys.value = getAllTreeKeys(deptOptions.value);
  });
}

function getList() {
  loading.value = true;
  listUser(proxy.addDateRange(queryParams, dateRange.value)).then(res => {
    loading.value = false;
    userList.value = res.data.rows || [];
    total.value = res.data.total || 0;
  }).catch(() => {
    loading.value = false;
  });
}

function handleNodeSelect(selectedKeys) {
  queryParams.deptId = selectedKeys?.[0] || undefined;
  queryParams.pageNum = 1;
  getList();
}

function handleQuery() {
  queryParams.pageNum = 1;
  getList();
}

function resetQuery() {
  dateRange.value = [];
  queryParams.userName = undefined;
  queryParams.phonenumber = undefined;
  queryParams.status = undefined;
  queryParams.deptId = undefined;
  getList();
}

function handleSelectionChange() {
  const records = userTableRef.value?.getCheckboxRecords() || [];
  ids.value = records.map(item => item.userId);
  single.value = records.length !== 1;
  multiple.value = !records.length;
}

function handlePageChange({ currentPage, pageSize }) {
  queryParams.pageNum = currentPage;
  queryParams.pageSize = pageSize;
  getList();
}

function handleDelete(row = {}) {
  const userIds = row.userId || ids.value;
  proxy.$modal.confirm(`是否确认删除用户编号为“${userIds}”的数据项？`).then(() => {
    return delUser(userIds);
  }).then(() => {
    getList();
    proxy.$modal.msgSuccess("删除成功");
  }).catch(() => {});
}

function handleExport() {
  proxy.download(
    "system/user/export",
    { ...proxy.addDateRange(queryParams, dateRange.value) },
    `user_${Date.now()}.xlsx`
  );
}

function handleStatusChange(row) {
  const text = row.status === "0" ? "启用" : "停用";
  proxy.$modal.confirm(`确认要${text}“${row.userName}”用户吗？`).then(() => {
    return changeUserStatus(row.userId, row.status);
  }).then(() => {
    proxy.$modal.msgSuccess(`${text}成功`);
  }).catch(() => {
    row.status = row.status === "0" ? "1" : "0";
  });
}

function handleCommand(command, row) {
  if (command === "handleResetPwd") {
    handleResetPwd(row);
  }
  if (command === "handleAuthRole") {
    handleAuthRole(row);
  }
}

function handleAuthRole(row) {
  router.push("/system/user-auth/role/" + row.userId);
}

function handleResetPwd(row) {
  proxy.$prompt(`请输入“${row.userName}”的新密码`, "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    closeOnClickModal: false,
    inputPattern: /^.{5,20}$/,
    inputErrorMessage: "用户密码长度必须介于 5 和 20 之间"
  }).then(({ value }) => {
    resetUserPwd(row.userId, value).then(() => {
      proxy.$modal.msgSuccess(`修改成功，新密码是：${value}`);
    });
  }).catch(() => {});
}

function handleImport() {
  upload.title = "用户导入";
  upload.open = true;
  uploadFileList.value = [];
}

function importTemplate() {
  proxy.download("system/user/importTemplate", {}, `user_template_${Date.now()}.xlsx`);
}

function beforeUpload() {
  return false;
}

function handleUploadChange(info) {
  uploadFileList.value = info.fileList.slice(-1);
  if (info.file.status === "uploading") {
    upload.isUploading = true;
  }
  if (info.file.status === "done") {
    upload.open = false;
    upload.isUploading = false;
    uploadFileList.value = [];
    proxy.$alert(info.file.response.msg, "导入结果", { dangerouslyUseHTMLString: true });
    getList();
  }
  if (info.file.status === "error") {
    upload.isUploading = false;
    proxy.$modal.msgError("上传失败");
  }
}

function submitFileForm() {
  if (!uploadFileList.value.length) {
    proxy.$modal.msgError("请选择要上传的文件");
    return;
  }
  upload.isUploading = true;
  uploadRef.value?.upload?.(uploadFileList.value[0]);
}

function resetFormModel() {
  Object.assign(form, {
    userId: undefined,
    deptId: undefined,
    userName: undefined,
    nickName: undefined,
    password: undefined,
    phonenumber: undefined,
    email: undefined,
    sex: "2",
    status: "0",
    remark: undefined,
    postIds: [],
    roleIds: []
  });
}

function cancel() {
  open.value = false;
  resetFormModel();
}

function handleAdd() {
  resetFormModel();
  getTreeselect();
  getUser().then(response => {
    postOptions.value = response.data.posts || [];
    roleOptions.value = response.data.roles || [];
    open.value = true;
    title.value = "添加用户";
    if (initPassword.value) {
      form.password = initPassword.value;
    }
  });
}

function handleUpdate(row = {}) {
  resetFormModel();
  getTreeselect();
  const userId = row.userId || ids.value[0];
  getUser(userId).then(response => {
    Object.assign(form, response.data.sysUser || {});
    postOptions.value = response.data.posts || [];
    roleOptions.value = response.data.roles || [];
    form.postIds = response.data.postIds || [];
    form.roleIds = response.data.roleIds || [];
    form.password = "";
    open.value = true;
    title.value = "修改用户";
  });
}

async function submitForm() {
  const errMap = await userFormRef.value?.validate?.().catch(err => err);
  if (errMap) {
    return;
  }
  const request = form.userId !== undefined ? updateUser(form) : addUser(form);
  request.then(() => {
    proxy.$modal.msgSuccess(form.userId !== undefined ? "修改成功" : "新增成功");
    open.value = false;
    getList();
  });
}

getTreeselect();
getList();
</script>

<style scoped>
.user-page-layout {
  display: grid;
  grid-template-columns: 280px minmax(0, 1fr);
  gap: 20px;
  align-items: start;
}

.user-page-aside,
.user-page-main {
  min-width: 0;
}

.panel-shell {
  padding: 18px;
  border-radius: 14px;
  background: #fff;
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.06);
}

.panel-title {
  margin-bottom: 16px;
  color: #1f2937;
  font-size: 15px;
  font-weight: 700;
}

.head-container {
  margin-bottom: 16px;
}

.dept-tree-wrap {
  max-height: calc(100vh - 260px);
  overflow: auto;
  padding-right: 4px;
}

.toolbar-spacer {
  height: 4px;
}

.user-toolbar {
  padding-top: 0;
}

.pager-wrap {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.action-links {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.more-link {
  color: #1677ff;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

.upload-tips {
  margin-top: 12px;
  text-align: center;
}

@media (max-width: 992px) {
  .user-page-layout {
    grid-template-columns: 1fr;
  }

  .dept-tree-wrap {
    max-height: 320px;
  }
}
</style>
