<template>
  <div class="app-container">
    <el-form :model="queryParams" ref="queryRef" v-show="showSearch" :inline="true">
      <el-form-item label="角色名称" prop="roleName">
        <el-input v-model="queryParams.roleName" placeholder="请输入角色名称" clearable size="small" style="width: 220px" @keyup.enter="handleQuery" />
      </el-form-item>
      <el-form-item label="权限字符" prop="roleKey">
        <el-input v-model="queryParams.roleKey" placeholder="请输入权限字符" clearable size="small" style="width: 220px" @keyup.enter="handleQuery" />
      </el-form-item>
      <el-form-item label="状态" prop="status">
        <el-select v-model="queryParams.status" placeholder="请选择状态" clearable size="small" style="width: 160px">
          <el-option v-for="dict in sys_normal_disable" :key="dict.value" :label="dict.label" :value="dict.value" />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" icon="Search" size="mini" @click="handleQuery">搜索</el-button>
        <el-button icon="Refresh" size="mini" @click="resetQuery">重置</el-button>
      </el-form-item>
    </el-form>

    <el-row :gutter="10" class="mb8">
      <el-col :span="1.5">
        <el-button type="primary" plain icon="Plus" size="mini" @click="handleAdd" v-hasPermi="['customer:portalRole:add']">新增</el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button type="success" plain icon="Edit" size="mini" :disabled="single" @click="handleUpdate" v-hasPermi="['customer:portalRole:edit']">修改</el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button type="danger" plain icon="Delete" size="mini" :disabled="multiple" @click="handleDelete" v-hasPermi="['customer:portalRole:remove']">删除</el-button>
      </el-col>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList"></right-toolbar>
    </el-row>

    <el-table v-loading="loading" :data="roleList" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" align="center" />
      <el-table-column label="角色名称" prop="roleName" min-width="160" />
      <el-table-column label="权限字符" prop="roleKey" min-width="180" />
      <el-table-column label="排序" prop="roleSort" width="80" align="center" />
      <el-table-column label="状态" width="100" align="center">
        <template #default="scope">
          <el-switch v-model="scope.row.status" active-value="0" inactive-value="1" @change="handleStatusChange(scope.row)" />
        </template>
      </el-table-column>
      <el-table-column label="创建时间" prop="createTime" width="180" align="center">
        <template #default="scope">
          <span>{{ parseTime(scope.row.createTime) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="160" align="center">
        <template #default="scope">
          <el-button size="mini" type="text" icon="Edit" @click="handleUpdate(scope.row)" v-hasPermi="['customer:portalRole:edit']">修改</el-button>
          <el-button size="mini" type="text" icon="Delete" @click="handleDelete(scope.row)" v-hasPermi="['customer:portalRole:remove']">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total > 0" :total="total" v-model:page="queryParams.pageNum" v-model:limit="queryParams.pageSize" @pagination="getList" />

    <el-dialog :title="title" v-model="open" width="560px" append-to-body>
      <el-form ref="roleRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="角色名称" prop="roleName">
          <el-input v-model="form.roleName" placeholder="请输入角色名称" />
        </el-form-item>
        <el-form-item label="权限字符" prop="roleKey">
          <el-input v-model="form.roleKey" placeholder="请输入权限字符" />
        </el-form-item>
        <el-form-item label="显示排序" prop="roleSort">
          <el-input-number v-model="form.roleSort" controls-position="right" :min="0" />
        </el-form-item>
        <el-form-item label="角色状态">
          <el-radio-group v-model="form.status">
            <el-radio v-for="dict in sys_normal_disable" :key="dict.value" :label="dict.value">{{ dict.label }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="菜单权限">
          <el-checkbox v-model="menuExpand" @change="handleCheckedTreeExpand">展开/折叠</el-checkbox>
          <el-checkbox v-model="menuNodeAll" @change="handleCheckedTreeNodeAll">全选/全不选</el-checkbox>
          <el-tree
            class="tree-border"
            :data="menuOptions"
            show-checkbox
            ref="menuRef"
            node-key="id"
            :props="{ label: 'label', children: 'children' }"
          />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" :rows="3" placeholder="请输入备注" />
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

<script setup name="CustomerPortalRole">
import { addPortalRole, changePortalRoleStatus, delPortalRole, getPortalRole, listPortalRole, rolePortalMenuTreeselect, updatePortalRole } from "@/api/customer/portalRole";
import { listPortalMenu } from "@/api/customer/portalMenu";

const { proxy } = getCurrentInstance();
const { sys_normal_disable } = proxy.useDict("sys_normal_disable");

const roleList = ref([]);
const menuOptions = ref([]);
const loading = ref(false);
const showSearch = ref(true);
const open = ref(false);
const title = ref("");
const ids = ref([]);
const single = ref(true);
const multiple = ref(true);
const total = ref(0);
const menuExpand = ref(false);
const menuNodeAll = ref(false);
const menuRef = ref(null);

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
  if (menuRef.value) {
    menuRef.value.setCheckedKeys([]);
  }
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

function handleSelectionChange(selection) {
  ids.value = selection.map(item => item.roleId);
  single.value = selection.length !== 1;
  multiple.value = !selection.length;
}

function handleCheckedTreeExpand(value) {
  const treeList = menuOptions.value;
  for (let i = 0; i < treeList.length; i++) {
    menuRef.value.store.nodesMap[treeList[i].id].expanded = value;
  }
}

function handleCheckedTreeNodeAll(value) {
  menuRef.value.setCheckedNodes(value ? menuOptions.value : []);
}

function getMenuAllCheckedKeys() {
  let checkedKeys = menuRef.value.getCheckedKeys();
  let halfCheckedKeys = menuRef.value.getHalfCheckedKeys();
  checkedKeys.unshift.apply(checkedKeys, halfCheckedKeys);
  return checkedKeys;
}

function handleAdd() {
  reset();
  getRoleMenuTreeselect(0);
  open.value = true;
  title.value = "新增客户端角色";
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
    nextTick(() => {
      roleMenu.then(res => {
        (res.data.checkedKeys || []).forEach((id) => {
          menuRef.value.setChecked(id, true, false);
        });
      });
    });
  });
}

function submitForm() {
  proxy.$refs["roleRef"].validate(valid => {
    if (!valid) {
      return;
    }
    form.value.menuIds = getMenuAllCheckedKeys();
    const request = form.value.roleId ? updatePortalRole(form.value) : addPortalRole(form.value);
    request.then(() => {
      proxy.$modal.msgSuccess("保存成功");
      open.value = false;
      getList();
    });
  });
}

function handleStatusChange(row) {
  const text = row.status === "0" ? "启用" : "停用";
  proxy.$modal.confirm('确认要' + text + '角色 "' + row.roleName + '" 吗？').then(() => {
    return changePortalRoleStatus(row.roleId, row.status);
  }).then(() => {
    proxy.$modal.msgSuccess(text + "成功");
  }).catch(() => {
    row.status = row.status === "0" ? "1" : "0";
  });
}

function handleDelete(row) {
  const roleIds = row?.roleId || ids.value;
  proxy.$modal.confirm('是否确认删除角色 "' + roleIds + '"？').then(() => {
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
