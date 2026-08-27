<template>
  <div class="app-container">
    <a-form ref="queryRef" :model="queryParams" v-show="showSearch" layout="inline">
      <a-form-item label="角色名称" name="roleName">
        <a-input v-model:value="queryParams.roleName" placeholder="请输入角色名称" allow-clear style="width: 240px"
          @keyup.enter="handleQuery" />
      </a-form-item>

      <a-form-item label="权限字符" name="roleKey">
        <a-input v-model:value="queryParams.roleKey" placeholder="请输入权限字符" allow-clear style="width: 240px"
          @keyup.enter="handleQuery" />
      </a-form-item>

      <a-form-item label="状态" name="status">
        <a-select v-model:value="queryParams.status" placeholder="角色状态" allow-clear style="width: 240px">
          <a-select-option v-for="dict in sys_normal_disable" :key="dict.value" :value="dict.value">
            {{ dict.label }}
          </a-select-option>
        </a-select>
      </a-form-item>

      <a-form-item label="创建时间">
        <a-range-picker v-model:value="dateRange" value-format="YYYY-MM-DD" style="width: 240px" />
      </a-form-item>

      <a-form-item>
        <a-space>
          <a-button type="primary" @click="handleQuery">
            <template #icon>
              <SearchOutlined />
            </template>
            搜索
          </a-button>

          <a-button @click="resetQuery">
            <template #icon>
              <ReloadOutlined />
            </template>
            重置
          </a-button>
        </a-space>
      </a-form-item>
    </a-form>

    <a-row :gutter="10" class="mb8" style="margin: 12px 0">
      <a-col>
        <a-button type="primary" ghost @click="handleAdd" v-hasPermi="['system:role:add']">
          <template #icon>
            <PlusOutlined />
          </template>
          新增
        </a-button>
      </a-col>

      <a-col>
        <a-button ghost :disabled="single" @click="handleUpdate" v-hasPermi="['system:role:edit']">
          <template #icon>
            <EditOutlined />
          </template>
          修改
        </a-button>
      </a-col>

      <a-col>
        <a-button danger ghost :disabled="multiple" @click="handleDelete" v-hasPermi="['system:role:remove']">
          <template #icon>
            <DeleteOutlined />
          </template>
          删除
        </a-button>
      </a-col>

      <a-col>
        <a-button ghost @click="handleExport" v-hasPermi="['system:role:export']">
          <template #icon>
            <DownloadOutlined />
          </template>
          导出
        </a-button>
      </a-col>

      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList" />
    </a-row>

    <a-table :loading="loading" :columns="columns" :data-source="roleList" row-key="roleId" :pagination="false"
      :row-selection="rowSelection">
      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex === 'status'">
          <a-switch v-model:checked="record.status" checked-value="0" un-checked-value="1"
            @change="() => handleStatusChange(record)" />
        </template>

        <template v-else-if="column.dataIndex === 'createTime'">
          <span>{{ parseTime(record.createTime) }}</span>
        </template>

        <template v-else-if="column.key === 'action'">
          <a-space v-if="record.roleId !== 1">
            <a-button type="link" size="small" @click="handleUpdate(record)" v-hasPermi="['system:role:edit']">
              <template #icon>
                <EditOutlined />
              </template>
              修改
            </a-button>

            <a-button type="link" size="small" danger @click="handleDelete(record)" v-hasPermi="['system:role:remove']">
              <template #icon>
                <DeleteOutlined />
              </template>
              删除
            </a-button>

            <a-dropdown>
              <a-button type="link" size="small" v-hasPermi="['system:role:edit']">
                <DoubleRightOutlined />
                更多
              </a-button>

              <template #overlay>
                <a-menu @click="({ key }) => handleCommand(key, record)">
                  <a-menu-item key="handleDataScope" v-hasPermi="['system:role:edit']">
                    <CheckCircleOutlined />
                    数据权限
                  </a-menu-item>

                  <a-menu-item key="handleAuthUser" v-hasPermi="['system:role:edit']">
                    <UserOutlined />
                    分配用户
                  </a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
          </a-space>
        </template>
      </template>
    </a-table>

    <pagination v-show="total > 0" :total="total" v-model:page="queryParams.pageNum"
      v-model:limit="queryParams.pageSize" @pagination="getList" />

    <!-- 添加或修改角色配置对话框 -->
    <a-modal v-model:open="open" :title="title" width="500px" :mask-closable="false" @ok="submitForm" @cancel="cancel">
      <a-form ref="roleRef" :model="form" :rules="rules" :label-col="{ style: { width: '100px' } }">
        <a-form-item label="角色名称" name="roleName">
          <a-input v-model:value="form.roleName" placeholder="请输入角色名称" />
        </a-form-item>

        <a-form-item name="roleKey">
          <template #label>
            <span>
              <a-tooltip content="控制器中定义的权限字符，如：@PreAuthorize(`@ss.hasRole('admin')`)">
                <QuestionCircleOutlined />
              </a-tooltip>
              权限字符
            </span>
          </template>

          <a-input v-model:value="form.roleKey" placeholder="请输入权限字符" />
        </a-form-item>

        <a-form-item label="角色顺序" name="roleSort">
          <a-input-number v-model:value="form.roleSort" :min="0" style="width: 100%" />
        </a-form-item>

        <a-form-item label="状态" name="status">
          <a-radio-group v-model:value="form.status">
            <a-radio v-for="dict in sys_normal_disable" :key="dict.value" :value="dict.value">
              {{ dict.label }}
            </a-radio>
          </a-radio-group>
        </a-form-item>

        <a-form-item label="菜单权限">
          <a-space direction="vertical" style="width: 100%">
            <a-space>
              <a-checkbox v-model:checked="menuExpand"
                @change="event => handleCheckedTreeExpand(event.target.checked, 'menu')">
                展开/折叠
              </a-checkbox>

              <a-checkbox v-model:checked="menuNodeAll"
                @change="event => handleCheckedTreeNodeAll(event.target.checked, 'menu')">
                全选/全不选
              </a-checkbox>

              <a-checkbox v-model:checked="form.menuCheckStrictly"
                @change="event => handleCheckedTreeConnect(event.target.checked, 'menu')">
                父子联动
              </a-checkbox>
            </a-space>

            <a-tree class="tree-border" checkable :tree-data="menuOptions"
              :field-names="{ title: 'label', key: 'id', children: 'children' }" :checked-keys="menuCheckedKeys"
              :expanded-keys="menuExpandedKeys" :check-strictly="!form.menuCheckStrictly" @check="handleMenuCheck"
              @expand="keys => (menuExpandedKeys = keys)" />
          </a-space>
        </a-form-item>

        <a-form-item label="备注" name="remark">
          <a-textarea v-model:value="form.remark" placeholder="请输入内容" />
        </a-form-item>
      </a-form>

      <template #footer>
        <a-space>
          <a-button @click="cancel">取 消</a-button>
          <a-button type="primary" @click="submitForm">确 定</a-button>
        </a-space>
      </template>
    </a-modal>

    <!-- 分配角色数据权限对话框 -->
    <a-modal v-model:open="openDataScope" :title="title" width="500px" :mask-closable="false" @ok="submitDataScope"
      @cancel="cancelDataScope">
      <a-form :model="form" :label-col="{ style: { width: '80px' } }">
        <a-form-item label="角色名称">
          <a-input v-model:value="form.roleName" disabled />
        </a-form-item>

        <a-form-item label="权限字符">
          <a-input v-model:value="form.roleKey" disabled />
        </a-form-item>

        <a-form-item label="权限范围">
          <a-select v-model:value="form.dataScope" @change="dataScopeSelectChange" style="width: 100%">
            <a-select-option v-for="item in dataScopeOptions" :key="item.value" :value="item.value">
              {{ item.label }}
            </a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item label="数据权限" v-show="form.dataScope == 2">
          <a-space direction="vertical" style="width: 100%">
            <a-space>
              <a-checkbox v-model:checked="deptExpand"
                @change="event => handleCheckedTreeExpand(event.target.checked, 'dept')">
                展开/折叠
              </a-checkbox>

              <a-checkbox v-model:checked="deptNodeAll"
                @change="event => handleCheckedTreeNodeAll(event.target.checked, 'dept')">
                全选/全不选
              </a-checkbox>

              <a-checkbox v-model:checked="form.deptCheckStrictly"
                @change="event => handleCheckedTreeConnect(event.target.checked, 'dept')">
                父子联动
              </a-checkbox>
            </a-space>

            <a-tree class="tree-border" checkable :tree-data="deptOptions"
              :field-names="{ title: 'label', key: 'id', children: 'children' }" :checked-keys="deptCheckedKeys"
              :expanded-keys="deptExpandedKeys" :check-strictly="!form.deptCheckStrictly" @check="handleDeptCheck"
              @expand="keys => (deptExpandedKeys = keys)" />
          </a-space>
        </a-form-item>
      </a-form>

      <template #footer>
        <a-space>
          <a-button @click="cancelDataScope">取 消</a-button>
          <a-button type="primary" @click="submitDataScope">确 定</a-button>
        </a-space>
      </template>
    </a-modal>
  </div>
</template>

<script setup name="Role">
import {
  SearchOutlined,
  ReloadOutlined,
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  DownloadOutlined,
  DoubleRightOutlined,
  CheckCircleOutlined,
  UserOutlined,
  QuestionCircleOutlined
} from "@ant-design/icons-vue";

import {
  addRole,
  changeRoleStatus,
  dataScope,
  delRole,
  getRole,
  listRole,
  updateRole
} from "@/api/system/role";

import {
  roleMenuTreeselect,
  treeselect as menuTreeselect
} from "@/api/system/menu";

import {
  treeselect as deptTreeselect,
  roleDeptTreeselect
} from "@/api/system/dept";

const router = useRouter();
const { proxy } = getCurrentInstance();
const { sys_normal_disable } = proxy.useDict("sys_normal_disable");

const queryRef = ref();
const roleRef = ref();

const roleList = ref([]);
const open = ref(false);
const loading = ref(true);
const showSearch = ref(true);
const ids = ref([]);
const single = ref(true);
const multiple = ref(true);
const total = ref(0);
const title = ref("");
const dateRange = ref([]);

const menuOptions = ref([]);
const menuExpand = ref(false);
const menuNodeAll = ref(false);
const menuCheckedKeys = ref([]);
const menuHalfCheckedKeys = ref([]);
const menuExpandedKeys = ref([]);

const deptExpand = ref(true);
const deptNodeAll = ref(false);
const deptOptions = ref([]);
const deptCheckedKeys = ref([]);
const deptHalfCheckedKeys = ref([]);
const deptExpandedKeys = ref([]);

const openDataScope = ref(false);

/** 数据范围选项 */
const dataScopeOptions = ref([
  { value: "1", label: "全部数据权限" },
  { value: "2", label: "自定数据权限" },
  { value: "3", label: "本部门数据权限" },
  { value: "4", label: "本部门及以下数据权限" },
  { value: "5", label: "仅本人数据权限" }
]);

const data = reactive({
  form: {},
  queryParams: {
    pageNum: 1,
    pageSize: 10,
    roleName: undefined,
    roleKey: undefined,
    status: undefined
  },
  rules: {
    roleName: [
      {
        required: true,
        message: "角色名称不能为空",
        trigger: "blur"
      }
    ],
    roleKey: [
      {
        required: true,
        message: "权限字符不能为空",
        trigger: "blur"
      }
    ],
    roleSort: [
      {
        required: true,
        message: "角色顺序不能为空",
        trigger: "blur"
      }
    ]
  }
});

const { queryParams, form, rules } = toRefs(data);

const columns = [
  {
    title: "角色编号",
    dataIndex: "roleId",
    key: "roleId",
    width: 120
  },
  {
    title: "角色名称",
    dataIndex: "roleName",
    key: "roleName",
    width: 150,
    ellipsis: true
  },
  {
    title: "权限字符",
    dataIndex: "roleKey",
    key: "roleKey",
    width: 150,
    ellipsis: true
  },
  {
    title: "显示顺序",
    dataIndex: "roleSort",
    key: "roleSort",
    width: 100
  },
  {
    title: "状态",
    dataIndex: "status",
    key: "status",
    align: "center",
    width: 100
  },
  {
    title: "创建时间",
    dataIndex: "createTime",
    key: "createTime",
    align: "center",
    width: 180
  },
  {
    title: "操作",
    key: "action",
    align: "center"
  }
];

const rowSelection = computed(() => {
  return {
    selectedRowKeys: ids.value,
    onChange: selectedRowKeys => {
      ids.value = selectedRowKeys;
      single.value = selectedRowKeys.length !== 1;
      multiple.value = !selectedRowKeys.length;
    }
  };
});

/** 查询角色列表 */
function getList() {
  loading.value = true;

  listRole(proxy.addDateRange(queryParams.value, dateRange.value)).then(response => {
    roleList.value = response.data.rows;
    total.value = response.data.total;
    loading.value = false;
  });
}

/** 搜索按钮操作 */
function handleQuery() {
  queryParams.value.pageNum = 1;
  getList();
}

/** 重置按钮操作 */
function resetQuery() {
  dateRange.value = [];
  queryRef.value?.resetFields?.();
  handleQuery();
}

/** 删除按钮操作 */
function handleDelete(row = {}) {
  const roleIds = row.roleId || ids.value;

  proxy.$modal
    .confirm('是否确认删除角色编号为"' + roleIds + '"的数据项?')
    .then(function () {
      return delRole(roleIds);
    })
    .then(() => {
      getList();
      proxy.$modal.msgSuccess("删除成功");
    })
    .catch(() => { });
}

/** 导出按钮操作 */
function handleExport() {
  proxy.download(
    "system/role/export",
    {
      ...queryParams.value
    },
    `role_${new Date().getTime()}.xlsx`
  );
}

/** 角色状态修改 */
function handleStatusChange(row) {
  const text = row.status === "0" ? "启用" : "停用";

  proxy.$modal
    .confirm('确认要"' + text + '""' + row.roleName + '"角色吗?')
    .then(function () {
      return changeRoleStatus(row.roleId, row.status);
    })
    .then(() => {
      proxy.$modal.msgSuccess(text + "成功");
    })
    .catch(function () {
      row.status = row.status === "0" ? "1" : "0";
    });
}

/** 更多操作 */
function handleCommand(command, row) {
  switch (command) {
    case "handleDataScope":
      handleDataScope(row);
      break;
    case "handleAuthUser":
      handleAuthUser(row);
      break;
    default:
      break;
  }
}

/** 分配用户 */
function handleAuthUser(row) {
  router.push("/system/role-auth/user/" + row.roleId);
}

/** 查询菜单树结构 */
function getMenuTreeselect() {
  menuTreeselect().then(response => {
    menuOptions.value = proxy.handleProps(response.data, "menuId", "menuName");
    menuExpandedKeys.value = getAllTreeKeys(menuOptions.value);
  });
}

/** 根据角色ID查询菜单树结构 */
function getRoleMenuTreeselect(roleId) {
  return roleMenuTreeselect(roleId).then(response => {
    menuOptions.value = proxy.handleProps(response.data.menus, "menuId", "menuName");
    menuExpandedKeys.value = getAllTreeKeys(menuOptions.value);
    return response;
  });
}

/** 根据角色ID查询部门树结构 */
function getRoleDeptTreeselect(roleId) {
  return roleDeptTreeselect(roleId).then(response => {
    deptOptions.value = proxy.handleProps(response.data.depts, "deptId", "deptName");
    deptExpandedKeys.value = getAllTreeKeys(deptOptions.value);
    return response;
  });
}

/** 递归获取所有树节点 key */
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

/** 重置新增的表单以及其他数据 */
function reset() {
  menuCheckedKeys.value = [];
  menuHalfCheckedKeys.value = [];
  deptCheckedKeys.value = [];
  deptHalfCheckedKeys.value = [];

  menuExpand.value = false;
  menuNodeAll.value = false;
  deptExpand.value = true;
  deptNodeAll.value = false;

  form.value = {
    roleId: undefined,
    roleName: undefined,
    roleKey: undefined,
    roleSort: 0,
    status: "0",
    menuIds: [],
    deptIds: [],
    menuCheckStrictly: true,
    deptCheckStrictly: true,
    remark: undefined
  };

  nextTick(() => {
    roleRef.value?.resetFields?.();
  });
}

/** 添加角色 */
function handleAdd() {
  reset();
  getMenuTreeselect();
  open.value = true;
  title.value = "添加角色";
}

/** 修改角色 */
function handleUpdate(row = {}) {
  reset();

  const roleId = row.roleId || ids.value[0];
  const roleMenu = getRoleMenuTreeselect(roleId);

  getRole(roleId).then(response => {
    form.value = response.data;
    form.value.roleSort = Number(form.value.roleSort);
    open.value = true;

    nextTick(() => {
      roleMenu.then(res => {
        menuCheckedKeys.value = res.data.checkedKeys || [];
      });
    });

    title.value = "修改角色";
  });
}

/** 树权限（展开/折叠） */
function handleCheckedTreeExpand(value, type) {
  if (type === "menu") {
    menuExpandedKeys.value = value ? getAllTreeKeys(menuOptions.value) : [];
  } else if (type === "dept") {
    deptExpandedKeys.value = value ? getAllTreeKeys(deptOptions.value) : [];
  }
}

/** 树权限（全选/全不选） */
function handleCheckedTreeNodeAll(value, type) {
  if (type === "menu") {
    menuCheckedKeys.value = value ? getAllTreeKeys(menuOptions.value) : [];
    menuHalfCheckedKeys.value = [];
  } else if (type === "dept") {
    deptCheckedKeys.value = value ? getAllTreeKeys(deptOptions.value) : [];
    deptHalfCheckedKeys.value = [];
  }
}

/** 树权限（父子联动） */
function handleCheckedTreeConnect(value, type) {
  if (type === "menu") {
    form.value.menuCheckStrictly = !!value;
  } else if (type === "dept") {
    form.value.deptCheckStrictly = !!value;
  }
}

/** 菜单树选择 */
function handleMenuCheck(checkedKeys, info) {
  if (Array.isArray(checkedKeys)) {
    menuCheckedKeys.value = checkedKeys;
    menuHalfCheckedKeys.value = info.halfCheckedKeys || [];
  } else {
    menuCheckedKeys.value = checkedKeys.checked || [];
    menuHalfCheckedKeys.value = checkedKeys.halfChecked || [];
  }
}

/** 部门树选择 */
function handleDeptCheck(checkedKeys, info) {
  if (Array.isArray(checkedKeys)) {
    deptCheckedKeys.value = checkedKeys;
    deptHalfCheckedKeys.value = info.halfCheckedKeys || [];
  } else {
    deptCheckedKeys.value = checkedKeys.checked || [];
    deptHalfCheckedKeys.value = checkedKeys.halfChecked || [];
  }
}

/** 所有菜单节点数据 */
function getMenuAllCheckedKeys() {
  return [...menuCheckedKeys.value, ...menuHalfCheckedKeys.value];
}

/** 所有部门节点数据 */
function getDeptAllCheckedKeys() {
  return [...deptCheckedKeys.value, ...deptHalfCheckedKeys.value];
}

/** 提交按钮 */
async function submitForm() {
  try {
    await roleRef.value.validate();

    form.value.menuIds = getMenuAllCheckedKeys();

    if (form.value.roleId != undefined) {
      updateRole(form.value).then(() => {
        proxy.$modal.msgSuccess("修改成功");
        open.value = false;
        getList();
      });
    } else {
      addRole(form.value).then(() => {
        proxy.$modal.msgSuccess("新增成功");
        open.value = false;
        getList();
      });
    }
  } catch (error) {
    // 表单校验未通过
  }
}

/** 取消按钮 */
function cancel() {
  open.value = false;
  reset();
}

/** 选择角色权限范围触发 */
function dataScopeSelectChange(value) {
  if (value !== "2") {
    deptCheckedKeys.value = [];
    deptHalfCheckedKeys.value = [];
  }
}

/** 分配数据权限操作 */
function handleDataScope(row) {
  reset();

  const roleDeptTreeselect = getRoleDeptTreeselect(row.roleId);

  getRole(row.roleId).then(response => {
    form.value = response.data;
    openDataScope.value = true;

    nextTick(() => {
      roleDeptTreeselect.then(res => {
        deptCheckedKeys.value = res.data.checkedKeys || [];
      });
    });

    title.value = "分配数据权限";
  });
}

/** 提交按钮（数据权限） */
function submitDataScope() {
  if (form.value.roleId != undefined) {
    form.value.deptIds = getDeptAllCheckedKeys();

    dataScope(form.value).then(() => {
      proxy.$modal.msgSuccess("修改成功");
      openDataScope.value = false;
      getList();
    });
  }
}

/** 取消按钮（数据权限） */
function cancelDataScope() {
  openDataScope.value = false;
  reset();
}

getList();
</script>