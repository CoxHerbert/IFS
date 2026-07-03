<template>
  <div class="app-container">
    <el-form :model="queryParams" ref="queryRef" :inline="true" v-show="showSearch">
      <el-form-item label="菜单名称" prop="menuName">
        <el-input
          v-model="queryParams.menuName"
          placeholder="请输入菜单名称"
          clearable
          size="small"
          @keyup.enter="handleQuery"
        />
      </el-form-item>
      <el-form-item label="状态" prop="status">
        <el-select v-model="queryParams.status" placeholder="请选择状态" clearable size="small">
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
        <el-button type="primary" plain icon="Plus" size="mini" @click="handleAdd" v-hasPermi="['customer:portalMenu:add']">
          新增
        </el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button type="info" plain icon="Sort" size="mini" @click="toggleExpandAll">展开/折叠</el-button>
      </el-col>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList" />
    </el-row>

    <el-table
      v-if="refreshTable"
      v-loading="loading"
      :data="menuList"
      row-key="menuId"
      :default-expand-all="isExpandAll"
      :tree-props="{ children: 'children' }"
    >
      <el-table-column prop="menuName" label="菜单名称" min-width="180" />
      <el-table-column prop="menuType" label="类型" width="90" align="center">
        <template #default="scope">
          <el-tag v-if="scope.row.menuType === 'M'">目录</el-tag>
          <el-tag v-else-if="scope.row.menuType === 'C'" type="success">菜单</el-tag>
          <el-tag v-else type="info">按钮</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="path" label="路由地址" min-width="140" />
      <el-table-column prop="component" label="组件标识" min-width="180" />
      <el-table-column prop="perms" label="权限标识" min-width="180" />
      <el-table-column prop="icon" label="图标" width="120" align="center" />
      <el-table-column prop="status" label="状态" width="90" align="center">
        <template #default="scope">
          <dict-tag :options="sys_normal_disable" :value="scope.row.status" />
        </template>
      </el-table-column>
      <el-table-column label="操作" width="210" align="center">
        <template #default="scope">
          <el-button size="mini" type="text" icon="Edit" @click="handleUpdate(scope.row)" v-hasPermi="['customer:portalMenu:edit']">
            修改
          </el-button>
          <el-button size="mini" type="text" icon="Plus" @click="handleAdd(scope.row)" v-hasPermi="['customer:portalMenu:add']">
            新增
          </el-button>
          <el-button size="mini" type="text" icon="Delete" @click="handleDelete(scope.row)" v-hasPermi="['customer:portalMenu:remove']">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog :title="title" v-model="open" width="680px" append-to-body>
      <el-form ref="menuRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="上级菜单">
          <tree-select
            v-model:value="form.parentId"
            :options="menuOptions"
            :objMap="{ value: 'menuId', label: 'menuName', children: 'children' }"
            placeholder="请选择上级菜单"
          />
        </el-form-item>
        <el-form-item label="菜单类型" prop="menuType">
          <el-radio-group v-model="form.menuType">
            <el-radio label="M">目录</el-radio>
            <el-radio label="C">菜单</el-radio>
            <el-radio label="F">按钮</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-row>
          <el-col :span="12">
            <el-form-item label="菜单名称" prop="menuName">
              <el-input v-model="form.menuName" placeholder="请输入菜单名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="显示排序" prop="orderNum">
              <el-input-number v-model="form.orderNum" controls-position="right" :min="0" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12" v-if="form.menuType !== 'F'">
            <el-form-item label="路由地址" prop="path">
              <el-input v-model="form.path" placeholder="例如：workspace" />
            </el-form-item>
          </el-col>
          <el-col :span="12" v-if="form.menuType === 'C'">
            <el-form-item label="组件标识" prop="component">
              <el-select v-model="form.component" placeholder="请选择客户端页面" style="width: 100%">
                <el-option label="工作台" value="workspace/dashboard" />
                <el-option label="账号资料" value="workspace/account-profile" />
                <el-option label="出货查询" value="workspace/shipment-tracking" />
                <el-option label="智能出货助手" value="workspace/shipment-assistant" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row v-if="form.menuType === 'C'">
          <el-col :span="12">
            <el-form-item label="是否缓存">
              <el-radio-group v-model="form.isCache">
                <el-radio label="0">缓存</el-radio>
                <el-radio label="1">不缓存</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12" v-if="form.menuType !== 'M'">
            <el-form-item label="权限标识">
              <el-input v-model="form.perms" placeholder="例如：portal:shipment:view" />
            </el-form-item>
          </el-col>
          <el-col :span="12" v-if="form.menuType !== 'F'">
            <el-form-item label="图标">
              <el-select v-model="form.icon" placeholder="请选择图标" clearable style="width: 100%">
                <el-option label="工作台" value="AppstoreOutlined" />
                <el-option label="账号" value="ProfileOutlined" />
                <el-option label="出货" value="RadarChartOutlined" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row v-if="form.menuType !== 'F'">
          <el-col :span="12">
            <el-form-item label="显示状态">
              <el-radio-group v-model="form.visible">
                <el-radio label="0">显示</el-radio>
                <el-radio label="1">隐藏</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="菜单状态">
              <el-radio-group v-model="form.status">
                <el-radio v-for="dict in sys_normal_disable" :key="dict.value" :label="dict.value">{{ dict.label }}</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
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

<script setup name="CustomerPortalMenu">
import { addPortalMenu, delPortalMenu, getPortalMenu, listPortalMenu, updatePortalMenu } from "@/api/customer/portalMenu";

const { proxy } = getCurrentInstance();
const { sys_normal_disable } = proxy.useDict("sys_normal_disable");

const menuList = ref([]);
const menuOptions = ref([]);
const loading = ref(false);
const showSearch = ref(true);
const open = ref(false);
const title = ref("");
const isExpandAll = ref(false);
const refreshTable = ref(true);

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

function getList() {
  loading.value = true;
  listPortalMenu(queryParams.value).then(response => {
    menuList.value = response.data || [];
    loading.value = false;
  });
}

function getTreeselect() {
  listPortalMenu().then(response => {
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
  refreshTable.value = false;
  isExpandAll.value = !isExpandAll.value;
  nextTick(() => {
    refreshTable.value = true;
  });
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
  proxy.$refs["menuRef"].validate(valid => {
    if (!valid) {
      return;
    }
    const request = form.value.menuId ? updatePortalMenu(form.value) : addPortalMenu(form.value);
    request.then(() => {
      proxy.$modal.msgSuccess("保存成功");
      open.value = false;
      getList();
    });
  });
}

function handleDelete(row) {
  proxy.$modal.confirm('是否确认删除菜单“' + row.menuName + '”？').then(() => {
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
