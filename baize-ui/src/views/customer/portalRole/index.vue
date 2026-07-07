<template>
  <div class="app-container">
    <a-form
      v-show="showSearch"
      ref="queryRef"
      :model="queryParams"
      layout="inline"
      class="search-form"
    >
      <a-form-item label="角色名称" name="roleName">
        <a-input
          v-model:value="queryParams.roleName"
          allow-clear
          placeholder="请输入角色名称"
          style="width: 220px"
          @pressEnter="handleQuery"
        />
      </a-form-item>
      <a-form-item label="权限字符" name="roleKey">
        <a-input
          v-model:value="queryParams.roleKey"
          allow-clear
          placeholder="请输入权限字符"
          style="width: 220px"
          @pressEnter="handleQuery"
        />
      </a-form-item>
      <a-form-item label="状态" name="status">
        <a-select
          v-model:value="queryParams.status"
          allow-clear
          placeholder="请选择状态"
          style="width: 160px"
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
        <a-button type="primary" @click="handleAdd" v-hasPermi="['customer:portalRole:add']">新增</a-button>
        <a-button :disabled="single" @click="handleUpdate()" v-hasPermi="['customer:portalRole:edit']">修改</a-button>
        <a-button danger :disabled="multiple" @click="handleDelete()" v-hasPermi="['customer:portalRole:remove']">删除</a-button>
      </a-space>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList" />
    </div>

    <a-table
      :loading="loading"
      :data-source="roleList"
      :columns="roleColumns"
      :pagination="false"
      :row-selection="roleRowSelection"
      :scroll="{ x: 900 }"
      row-key="roleId"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <a-switch
            v-model:checked="record.status"
            checked-value="0"
            un-checked-value="1"
            @change="handleStatusChange(record)"
          />
        </template>
        <template v-else-if="column.key === 'createTime'">
          {{ parseTime(record.createTime) }}
        </template>
        <template v-else-if="column.key === 'action'">
          <a-space>
            <a-button type="link" @click="handleUpdate(record)" v-hasPermi="['customer:portalRole:edit']">修改</a-button>
            <a-button type="link" danger @click="handleDelete(record)" v-hasPermi="['customer:portalRole:remove']">删除</a-button>
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

    <a-modal v-model:open="open" :title="title" width="560px" :footer="null" destroy-on-close>
      <a-form ref="roleRef" :model="form" :rules="rules" :label-col="{ style: { width: '100px' } }">
        <a-form-item label="角色名称" name="roleName">
          <a-input v-model:value="form.roleName" placeholder="请输入角色名称" />
        </a-form-item>
        <a-form-item label="权限字符" name="roleKey">
          <a-input v-model:value="form.roleKey" placeholder="请输入权限字符" />
        </a-form-item>
        <a-form-item label="显示排序" name="roleSort">
          <a-input-number v-model:value="form.roleSort" :min="0" style="width: 100%" />
        </a-form-item>
        <a-form-item label="角色状态" name="status">
          <a-radio-group v-model:value="form.status">
            <a-radio v-for="dict in sys_normal_disable" :key="dict.value" :value="dict.value">
              {{ dict.label }}
            </a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="菜单权限">
          <a-space direction="vertical" style="width: 100%">
            <a-space>
              <a-checkbox v-model:checked="menuExpand" @change="handleCheckedTreeExpand">展开/折叠</a-checkbox>
              <a-checkbox v-model:checked="menuNodeAll" @change="handleCheckedTreeNodeAll">全选/全不选</a-checkbox>
            </a-space>
            <a-tree
              class="tree-border"
              checkable
              :tree-data="menuOptions"
              :field-names="{ title: 'label', key: 'id', children: 'children' }"
              :checked-keys="checkedMenuKeys"
              :expanded-keys="expandedMenuKeys"
              @check="handleTreeCheck"
              @update:expandedKeys="expandedMenuKeys = $event"
            />
          </a-space>
        </a-form-item>
        <a-form-item label="备注" name="remark">
          <a-textarea v-model:value="form.remark" :rows="3" placeholder="请输入备注" />
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

<script setup name="CustomerPortalRole">
import { addPortalRole, changePortalRoleStatus, delPortalRole, getPortalRole, listPortalRole, rolePortalMenuTreeselect, updatePortalRole } from "@/api/customer/portalRole";
import { listPortalMenu } from "@/api/customer/portalMenu";

const { proxy } = getCurrentInstance();
const { sys_normal_disable } = proxy.useDict("sys_normal_disable");

const roleColumns = [
  { title: "角色名称", dataIndex: "roleName", key: "roleName", minWidth: 160 },
  { title: "权限字符", dataIndex: "roleKey", key: "roleKey", minWidth: 180 },
  { title: "排序", dataIndex: "roleSort", key: "roleSort", width: 80, align: "center" },
  { title: "状态", dataIndex: "status", key: "status", width: 100, align: "center" },
  { title: "创建时间", dataIndex: "createTime", key: "createTime", width: 180, align: "center" },
  { title: "操作", key: "action", width: 160, align: "center" }
];

const roleList = ref([]);
const menuOptions = ref([]);
const loading = ref(false);
const showSearch = ref(true);
const open = ref(false);
const title = ref("");
const ids = ref([]);
const selectedRowKeys = ref([]);
const single = ref(true);
const multiple = ref(true);
const total = ref(0);
const menuExpand = ref(false);
const menuNodeAll = ref(false);
const checkedMenuKeys = ref([]);
const halfCheckedMenuKeys = ref([]);
const expandedMenuKeys = ref([]);

const roleRowSelection = computed(() => ({
  selectedRowKeys: selectedRowKeys.value,
  onChange: (keys, rows) => handleSelectionChange(rows, keys)
}));

const statusOptions = computed(() => {
  return (sys_normal_disable.value || []).map(item => ({
    label: item.label,
    value: item.value
  }));
});

const data = reactive({
  queryParams: {
    pageNum: 1,
    pageSize: 10,
    roleName: undefined,
    roleKey: undefined,
    status: undefined
  },
  form: {},
  rules: {
    roleName: [{ required: true, message: "角色名称不能为空", trigger: "blur" }],
    roleKey: [{ required: true, message: "权限字符不能为空", trigger: "blur" }],
    roleSort: [{ required: true, message: "显示排序不能为空", trigger: "blur" }]
  }
});

const { queryParams, form, rules } = toRefs(data);

function getAllTreeKeys(list = []) {
  const keys = [];
  list.forEach(item => {
    keys.push(item.id);
    if (item.children?.length) {
      keys.push(...getAllTreeKeys(item.children));
    }
  });
  return keys;
}

function getList() {
  loading.value = true;
  listPortalRole(queryParams.value).then(response => {
    roleList.value = response.data.rows || [];
    total.value = response.data.total || 0;
    loading.value = false;
  });
}

function getRoleMenuTreeselect(roleId) {
  if (!roleId) {
    return listPortalMenu().then(response => {
      menuOptions.value = proxy.handleProps(response.data || [], "menuId", "menuName");
      return { data: { checkedKeys: [] } };
    });
  }
  return rolePortalMenuTreeselect(roleId).then(response => {
    menuOptions.value = proxy.handleProps(response.data.menus || [], "menuId", "menuName");
    return response;
  });
}

function reset() {
  checkedMenuKeys.value = [];
  halfCheckedMenuKeys.value = [];
  expandedMenuKeys.value = [];
  menuExpand.value = false;
  menuNodeAll.value = false;
  form.value = {
    roleId: undefined,
    roleName: undefined,
    roleKey: undefined,
    roleSort: 0,
    status: "0",
    menuIds: [],
    remark: undefined
  };
  proxy.resetForm("roleRef");
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
  ids.value = selection.map(item => item.roleId);
  selectedRowKeys.value = keys;
  single.value = selection.length !== 1;
  multiple.value = !selection.length;
}

function handleCheckedTreeExpand() {
  expandedMenuKeys.value = menuExpand.value ? getAllTreeKeys(menuOptions.value) : [];
}

function handleCheckedTreeNodeAll() {
  checkedMenuKeys.value = menuNodeAll.value ? getAllTreeKeys(menuOptions.value) : [];
  halfCheckedMenuKeys.value = [];
}

function handleTreeCheck(checkedKeys, event) {
  checkedMenuKeys.value = checkedKeys;
  halfCheckedMenuKeys.value = event.halfCheckedKeys || [];
  const allKeys = getAllTreeKeys(menuOptions.value);
  menuNodeAll.value = allKeys.length > 0 && checkedMenuKeys.value.length === allKeys.length;
}

function getMenuAllCheckedKeys() {
  return Array.from(new Set([...halfCheckedMenuKeys.value, ...checkedMenuKeys.value]));
}

function handleAdd() {
  reset();
  getRoleMenuTreeselect(0).then(() => {
    open.value = true;
    title.value = "新增客户端角色";
  });
}

function handleUpdate(row) {
  reset();
  const roleId = row?.roleId || ids.value[0];
  const roleMenu = getRoleMenuTreeselect(roleId);
  getPortalRole(roleId).then(response => {
    form.value = response.data;
    form.value.roleSort = Number(form.value.roleSort);
    open.value = true;
    title.value = "修改客户端角色";
    roleMenu.then(res => {
      checkedMenuKeys.value = (res.data.checkedKeys || []).slice();
    });
  });
}

function submitForm() {
  proxy.$refs.roleRef.validate().then(() => {
    form.value.menuIds = getMenuAllCheckedKeys();
    const request = form.value.roleId ? updatePortalRole(form.value) : addPortalRole(form.value);
    request.then(() => {
      proxy.$modal.msgSuccess("保存成功");
      open.value = false;
      getList();
    });
  }).catch(() => {});
}

function handleStatusChange(row) {
  const text = row.status === "0" ? "启用" : "停用";
  proxy.$modal.confirm(`确认要${text}角色“${row.roleName}”吗？`).then(() => {
    return changePortalRoleStatus(row.roleId, row.status);
  }).then(() => {
    proxy.$modal.msgSuccess(`${text}成功`);
  }).catch(() => {
    row.status = row.status === "0" ? "1" : "0";
  });
}

function handleDelete(row) {
  const roleIds = row?.roleId || ids.value;
  proxy.$modal.confirm(`是否确认删除角色“${roleIds}”？`).then(() => {
    return delPortalRole(roleIds);
  }).then(() => {
    proxy.$modal.msgSuccess("删除成功");
    getList();
  }).catch(() => {});
}

function cancel() {
  open.value = false;
  reset();
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

.tree-border {
  width: 100%;
  max-height: 320px;
  overflow: auto;
  padding: 12px;
  border: 1px solid #f0f0f0;
  border-radius: 6px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  margin-top: 24px;
}
</style>
