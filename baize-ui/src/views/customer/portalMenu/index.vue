<template>
  <div class="app-container">
    <a-form
      v-show="showSearch"
      ref="queryRef"
      :model="queryParams"
      layout="inline"
      class="search-form"
    >
      <a-form-item label="菜单名称" name="menuName">
        <a-input
          v-model:value="queryParams.menuName"
          allow-clear
          placeholder="请输入菜单名称"
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
        <a-button type="primary" @click="handleAdd()" v-hasPermi="['customer:portalMenu:add']">新增</a-button>
        <a-button @click="toggleExpandAll">展开/折叠</a-button>
      </a-space>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList" />
    </div>

    <a-table
      :loading="loading"
      :data-source="menuList"
      :columns="menuColumns"
      :pagination="false"
      :default-expand-all-rows="false"
      :expanded-row-keys="expandedRowKeys"
      row-key="menuId"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'menuType'">
          <a-tag v-if="record.menuType === 'M'">目录</a-tag>
          <a-tag v-else-if="record.menuType === 'C'" color="success">菜单</a-tag>
          <a-tag v-else color="default">按钮</a-tag>
        </template>
        <template v-else-if="column.key === 'status'">
          <dict-tag :options="sys_normal_disable" :value="record.status" />
        </template>
        <template v-else-if="column.key === 'action'">
          <a-space>
            <a-button type="link" @click="handleUpdate(record)" v-hasPermi="['customer:portalMenu:edit']">修改</a-button>
            <a-button type="link" @click="handleAdd(record)" v-hasPermi="['customer:portalMenu:add']">新增</a-button>
            <a-button type="link" danger @click="handleDelete(record)" v-hasPermi="['customer:portalMenu:remove']">删除</a-button>
          </a-space>
        </template>
      </template>
    </a-table>

    <a-modal v-model:open="open" :title="title" width="680px" :footer="null" destroy-on-close>
      <a-form ref="menuRef" :model="form" :rules="rules" :label-col="{ style: { width: '100px' } }">
        <a-form-item label="上级菜单" name="parentId">
          <tree-select
            v-model:value="form.parentId"
            :options="menuOptions"
            :objMap="{ value: 'menuId', label: 'menuName', children: 'children' }"
            placeholder="请选择上级菜单"
          />
        </a-form-item>
        <a-form-item label="菜单类型" name="menuType">
          <a-radio-group v-model:value="form.menuType">
            <a-radio value="M">目录</a-radio>
            <a-radio value="C">菜单</a-radio>
            <a-radio value="F">按钮</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="菜单名称" name="menuName">
              <a-input v-model:value="form.menuName" placeholder="请输入菜单名称" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="显示排序" name="orderNum">
              <a-input-number v-model:value="form.orderNum" :min="0" style="width: 100%" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col v-if="form.menuType !== 'F'" :span="12">
            <a-form-item label="路由地址" name="path">
              <a-input v-model:value="form.path" placeholder="例如：workspace" />
            </a-form-item>
          </a-col>
          <a-col v-if="form.menuType === 'C'" :span="12">
            <a-form-item label="组件标识" name="component">
              <a-select
                v-model:value="form.component"
                placeholder="请选择客户端页面"
                :options="componentOptions"
              />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row v-if="form.menuType === 'C'" :gutter="16">
          <a-col :span="12">
            <a-form-item label="是否缓存" name="isCache">
              <a-radio-group v-model:value="form.isCache">
                <a-radio value="0">缓存</a-radio>
                <a-radio value="1">不缓存</a-radio>
              </a-radio-group>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col v-if="form.menuType !== 'M'" :span="12">
            <a-form-item label="权限标识" name="perms">
              <a-input v-model:value="form.perms" placeholder="例如：portal:shipment:view" />
            </a-form-item>
          </a-col>
          <a-col v-if="form.menuType !== 'F'" :span="12">
            <a-form-item label="图标" name="icon">
              <a-select
                v-model:value="form.icon"
                allow-clear
                placeholder="请选择图标"
                :options="iconOptions"
              />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row v-if="form.menuType !== 'F'" :gutter="16">
          <a-col :span="12">
            <a-form-item label="显示状态" name="visible">
              <a-radio-group v-model:value="form.visible">
                <a-radio value="0">显示</a-radio>
                <a-radio value="1">隐藏</a-radio>
              </a-radio-group>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="菜单状态" name="status">
              <a-radio-group v-model:value="form.status">
                <a-radio v-for="dict in sys_normal_disable" :key="dict.value" :value="dict.value">
                  {{ dict.label }}
                </a-radio>
              </a-radio-group>
            </a-form-item>
          </a-col>
        </a-row>
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

<script setup name="CustomerPortalMenu">
import { addPortalMenu, delPortalMenu, getPortalMenu, listPortalMenu, updatePortalMenu } from "@/api/customer/portalMenu";

const { proxy } = getCurrentInstance();
const { sys_normal_disable } = proxy.useDict("sys_normal_disable");

const menuColumns = [
  { title: "菜单名称", dataIndex: "menuName", key: "menuName", minWidth: 180 },
  { title: "类型", dataIndex: "menuType", key: "menuType", width: 90, align: "center" },
  { title: "路由地址", dataIndex: "path", key: "path", minWidth: 140 },
  { title: "组件标识", dataIndex: "component", key: "component", minWidth: 180 },
  { title: "权限标识", dataIndex: "perms", key: "perms", minWidth: 180 },
  { title: "图标", dataIndex: "icon", key: "icon", width: 120, align: "center" },
  { title: "状态", dataIndex: "status", key: "status", width: 90, align: "center" },
  { title: "操作", key: "action", width: 210, align: "center" }
];

const componentOptions = [
  { label: "工作台", value: "workspace/dashboard" },
  { label: "账号资料", value: "workspace/account-profile" },
  { label: "出货查询", value: "workspace/shipment-tracking" },
  { label: "智能出货助手", value: "workspace/shipment-assistant" }
];

const iconOptions = [
  { label: "工作台", value: "AppstoreOutlined" },
  { label: "账号", value: "ProfileOutlined" },
  { label: "出货", value: "RadarChartOutlined" }
];

const menuList = ref([]);
const menuOptions = ref([]);
const loading = ref(false);
const showSearch = ref(true);
const open = ref(false);
const title = ref("");
const isExpandAll = ref(false);
const expandedRowKeys = ref([]);

const statusOptions = computed(() => {
  return (sys_normal_disable.value || []).map(item => ({
    label: item.label,
    value: item.value
  }));
});

const data = reactive({
  queryParams: {
    menuName: undefined,
    status: undefined
  },
  form: {},
  rules: {
    menuName: [{ required: true, message: "菜单名称不能为空", trigger: "blur" }],
    orderNum: [{ required: true, message: "显示排序不能为空", trigger: "blur" }],
    path: [{ required: true, message: "路由地址不能为空", trigger: "blur" }]
  }
});

const { queryParams, form, rules } = toRefs(data);

function getAllMenuKeys(list = []) {
  const keys = [];
  list.forEach(item => {
    keys.push(item.menuId);
    if (item.children?.length) {
      keys.push(...getAllMenuKeys(item.children));
    }
  });
  return keys;
}

function getList() {
  loading.value = true;
  listPortalMenu(queryParams.value).then(response => {
    menuList.value = response.data || [];
    expandedRowKeys.value = isExpandAll.value ? getAllMenuKeys(menuList.value) : [];
    loading.value = false;
  });
}

function getTreeselect() {
  return listPortalMenu().then(response => {
    menuOptions.value = [{
      menuId: 0,
      menuName: "主类目",
      children: response.data || []
    }];
  });
}

function reset() {
  form.value = {
    menuId: undefined,
    parentId: 0,
    menuName: undefined,
    orderNum: 0,
    path: undefined,
    component: undefined,
    menuType: "C",
    isCache: "0",
    visible: "0",
    status: "0",
    perms: undefined,
    icon: undefined,
    remark: undefined
  };
  proxy.resetForm("menuRef");
}

function handleQuery() {
  getList();
}

function resetQuery() {
  proxy.resetForm("queryRef");
  handleQuery();
}

function toggleExpandAll() {
  isExpandAll.value = !isExpandAll.value;
  expandedRowKeys.value = isExpandAll.value ? getAllMenuKeys(menuList.value) : [];
}

async function handleAdd(row) {
  reset();
  await getTreeselect();
  form.value.parentId = row?.menuId || 0;
  open.value = true;
  title.value = "新增客户端菜单";
}

async function handleUpdate(row) {
  reset();
  await getTreeselect();
  getPortalMenu(row.menuId).then(response => {
    form.value = response.data;
    open.value = true;
    title.value = "修改客户端菜单";
  });
}

function submitForm() {
  proxy.$refs.menuRef.validate().then(() => {
    const request = form.value.menuId ? updatePortalMenu(form.value) : addPortalMenu(form.value);
    request.then(() => {
      proxy.$modal.msgSuccess("保存成功");
      open.value = false;
      getList();
    });
  }).catch(() => {});
}

function handleDelete(row) {
  proxy.$modal.confirm(`是否确认删除菜单“${row.menuName}”？`).then(() => {
    return delPortalMenu(row.menuId);
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

.modal-footer {
  display: flex;
  justify-content: flex-end;
  margin-top: 24px;
}
</style>
