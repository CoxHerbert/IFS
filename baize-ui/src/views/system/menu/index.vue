<template>
   <div class="app-container">
      <a-form ref="queryRef" :model="queryParams" layout="inline" v-show="showSearch">
         <a-form-item label="菜单名称" name="menuName">
            <a-input v-model:value="queryParams.menuName" placeholder="请输入菜单名称" allow-clear
               @keyup.enter="handleQuery" />
         </a-form-item>

         <a-form-item label="状态" name="status">
            <a-select v-model:value="queryParams.status" placeholder="菜单状态" allow-clear style="width: 180px">
               <a-select-option v-for="dict in sys_normal_disable" :key="dict.value" :value="dict.value">
                  {{ dict.label }}
               </a-select-option>
            </a-select>
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
            <a-button type="primary" ghost @click="handleAdd" v-hasPermi="['system:menu:add']">
               <template #icon>
                  <PlusOutlined />
               </template>
               新增
            </a-button>
         </a-col>

         <a-col>
            <a-button type="default" ghost @click="toggleExpandAll">
               <template #icon>
                  <SwapOutlined />
               </template>
               展开/折叠
            </a-button>
         </a-col>

         <right-toolbar v-model:showSearch="showSearch" @queryTable="getList" />
      </a-row>

      <a-table v-if="refreshTable" :loading="loading" :columns="columns" :data-source="menuList" row-key="menuId"
         :pagination="false" :default-expand-all-rows="isExpandAll" :children-column-name="'children'">
         <template #bodyCell="{ column, record }">
            <template v-if="column.dataIndex === 'icon'">
               <svg-icon :icon-class="record.icon" />
            </template>

            <template v-else-if="column.dataIndex === 'status'">
               <dict-tag :options="sys_normal_disable" :value="record.status" />
            </template>

            <template v-else-if="column.dataIndex === 'createTime'">
               <span>{{ parseTime(record.createTime) }}</span>
            </template>

            <template v-else-if="column.key === 'action'">
               <a-space>
                  <a-button type="link" size="small" @click="handleUpdate(record)" v-hasPermi="['system:menu:edit']">
                     <template #icon>
                        <EditOutlined />
                     </template>
                     修改
                  </a-button>

                  <a-button type="link" size="small" @click="handleAdd(record)" v-hasPermi="['system:menu:add']">
                     <template #icon>
                        <PlusOutlined />
                     </template>
                     新增
                  </a-button>

                  <a-button type="link" size="small" danger @click="handleDelete(record)"
                     v-hasPermi="['system:menu:remove']">
                     <template #icon>
                        <DeleteOutlined />
                     </template>
                     删除
                  </a-button>
               </a-space>
            </template>
         </template>
      </a-table>

      <!-- 添加或修改菜单对话框 -->
      <a-modal v-model:open="open" :title="title" width="680px" :mask-closable="false" @ok="submitForm"
         @cancel="cancel">
         <a-form ref="menuRef" :model="form" :rules="rules" :label-col="{ style: { width: '100px' } }">
            <a-row :gutter="16">
               <a-col :span="24">
                  <a-form-item label="上级菜单" name="parentId">
                     <a-tree-select v-model:value="form.parentId" :tree-data="menuOptions" :field-names="{
                        value: 'menuId',
                        label: 'menuName',
                        children: 'children'
                     }" placeholder="选择上级菜单" allow-clear tree-default-expand-all style="width: 100%" />
                  </a-form-item>
               </a-col>

               <a-col :span="24">
                  <a-form-item label="菜单类型" name="menuType">
                     <a-radio-group v-model:value="form.menuType">
                        <a-radio value="M">目录</a-radio>
                        <a-radio value="C">菜单</a-radio>
                        <a-radio value="F">按钮</a-radio>
                     </a-radio-group>
                  </a-form-item>
               </a-col>

               <a-col :span="24" v-if="form.menuType != 'F'">
                  <a-form-item label="菜单图标" name="icon">
                     <a-popover placement="bottomLeft" trigger="click" v-model:open="showChooseIcon"
                        @openChange="handleIconPopoverOpenChange">
                        <template #content>
                           <div style="width: 540px">
                              <icon-select ref="iconSelectRef" @selected="selected" />
                           </div>
                        </template>

                        <a-input v-model:value="form.icon" placeholder="点击选择图标" readonly>
                           <template #prefix>
                              <svg-icon v-if="form.icon" :icon-class="form.icon" style="height: 32px; width: 16px" />
                              <SearchOutlined v-else />
                           </template>
                        </a-input>
                     </a-popover>
                  </a-form-item>
               </a-col>

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

               <a-col :span="12" v-if="form.menuType != 'F'">
                  <a-form-item>
                     <template #label>
                        <span>
                           <a-tooltip content="选择是外链则路由地址需要以`http(s)://`开头">
                              <QuestionCircleOutlined />
                           </a-tooltip>
                           是否外链
                        </span>
                     </template>

                     <a-radio-group v-model:value="form.isFrame">
                        <a-radio value="0">是</a-radio>
                        <a-radio value="1">否</a-radio>
                     </a-radio-group>
                  </a-form-item>
               </a-col>

               <a-col :span="12" v-if="form.menuType != 'F'">
                  <a-form-item name="path">
                     <template #label>
                        <span>
                           <a-tooltip content="访问的路由地址，如：`user`，如外网地址需内链访问则以`http(s)://`开头">
                              <QuestionCircleOutlined />
                           </a-tooltip>
                           路由地址
                        </span>
                     </template>

                     <a-input v-model:value="form.path" placeholder="请输入路由地址" />
                  </a-form-item>
               </a-col>

               <a-col :span="12" v-if="form.menuType == 'C'">
                  <a-form-item name="component">
                     <template #label>
                        <span>
                           <a-tooltip content="访问的组件路径，如：`system/user/index`，默认在`views`目录下">
                              <QuestionCircleOutlined />
                           </a-tooltip>
                           组件路径
                        </span>
                     </template>

                     <a-input v-model:value="form.component" placeholder="请输入组件路径" />
                  </a-form-item>
               </a-col>

               <a-col :span="12" v-if="form.menuType != 'M'">
                  <a-form-item name="perms">
                     <template #label>
                        <span>
                           <a-tooltip content="控制器中定义的权限字符，如：@PreAuthorize(`@ss.hasPermi('system:user:list')`)">
                              <QuestionCircleOutlined />
                           </a-tooltip>
                           权限字符
                        </span>
                     </template>

                     <a-input v-model:value="form.perms" placeholder="请输入权限标识" :maxlength="100" />
                  </a-form-item>
               </a-col>

               <a-col :span="12" v-if="form.menuType == 'C'">
                  <a-form-item name="query">
                     <template #label>
                        <span>
                           <a-tooltip content='访问路由的默认传递参数，如：{"id": 1, "name": "ry"}'>
                              <QuestionCircleOutlined />
                           </a-tooltip>
                           路由参数
                        </span>
                     </template>

                     <a-input v-model:value="form.query" placeholder="请输入路由参数" :maxlength="255" />
                  </a-form-item>
               </a-col>

               <a-col :span="12" v-if="form.menuType == 'C'">
                  <a-form-item name="isCache">
                     <template #label>
                        <span>
                           <a-tooltip content="选择是则会被`keep-alive`缓存，需要匹配组件的`name`和地址保持一致">
                              <QuestionCircleOutlined />
                           </a-tooltip>
                           是否缓存
                        </span>
                     </template>

                     <a-radio-group v-model:value="form.isCache">
                        <a-radio value="0">缓存</a-radio>
                        <a-radio value="1">不缓存</a-radio>
                     </a-radio-group>
                  </a-form-item>
               </a-col>

               <a-col :span="12" v-if="form.menuType != 'F'">
                  <a-form-item name="visible">
                     <template #label>
                        <span>
                           <a-tooltip content="选择隐藏则路由将不会出现在侧边栏，但仍然可以访问">
                              <QuestionCircleOutlined />
                           </a-tooltip>
                           显示状态
                        </span>
                     </template>

                     <a-radio-group v-model:value="form.visible">
                        <a-radio v-for="dict in sys_show_hide" :key="dict.value" :value="dict.value">
                           {{ dict.label }}
                        </a-radio>
                     </a-radio-group>
                  </a-form-item>
               </a-col>

               <a-col :span="12" v-if="form.menuType != 'F'">
                  <a-form-item name="status">
                     <template #label>
                        <span>
                           <a-tooltip content="选择停用则路由将不会出现在侧边栏，也不能被访问">
                              <QuestionCircleOutlined />
                           </a-tooltip>
                           菜单状态
                        </span>
                     </template>

                     <a-radio-group v-model:value="form.status">
                        <a-radio v-for="dict in sys_normal_disable" :key="dict.value" :value="dict.value">
                           {{ dict.label }}
                        </a-radio>
                     </a-radio-group>
                  </a-form-item>
               </a-col>
            </a-row>
         </a-form>

         <template #footer>
            <a-space>
               <a-button @click="cancel">取 消</a-button>
               <a-button type="primary" @click="submitForm">确 定</a-button>
            </a-space>
         </template>
      </a-modal>
   </div>
</template>

<script setup name="Menu">
import {
   SearchOutlined,
   ReloadOutlined,
   PlusOutlined,
   SwapOutlined,
   EditOutlined,
   DeleteOutlined,
   QuestionCircleOutlined
} from "@ant-design/icons-vue";

import {
   addMenu,
   delMenu,
   getMenu,
   listMenu,
   updateMenu
} from "@/api/system/menu";

import SvgIcon from "@/components/SvgIcon";
import IconSelect from "@/components/IconSelect";

const { proxy } = getCurrentInstance();
const { sys_show_hide, sys_normal_disable } = proxy.useDict(
   "sys_show_hide",
   "sys_normal_disable"
);

const queryRef = ref();
const menuRef = ref();

const menuList = ref([]);
const open = ref(false);
const loading = ref(true);
const showSearch = ref(true);
const title = ref("");
const menuOptions = ref([]);
const isExpandAll = ref(false);
const refreshTable = ref(true);
const showChooseIcon = ref(false);
const iconSelectRef = ref(null);

const columns = [
   {
      title: "菜单名称",
      dataIndex: "menuName",
      key: "menuName",
      width: 160,
      ellipsis: true
   },
   {
      title: "图标",
      dataIndex: "icon",
      key: "icon",
      align: "center",
      width: 100
   },
   {
      title: "排序",
      dataIndex: "orderNum",
      key: "orderNum",
      width: 60
   },
   {
      title: "权限标识",
      dataIndex: "perms",
      key: "perms",
      ellipsis: true
   },
   {
      title: "组件路径",
      dataIndex: "component",
      key: "component",
      ellipsis: true
   },
   {
      title: "状态",
      dataIndex: "status",
      key: "status",
      width: 80
   },
   {
      title: "创建时间",
      dataIndex: "createTime",
      key: "createTime",
      align: "center"
   },
   {
      title: "操作",
      key: "action",
      align: "center",
      width: 200
   }
];

const data = reactive({
   form: {},
   queryParams: {
      menuName: undefined,
      status: undefined
   },
   rules: {
      menuName: [
         {
            required: true,
            message: "菜单名称不能为空",
            trigger: "blur"
         }
      ],
      orderNum: [
         {
            required: true,
            message: "菜单顺序不能为空",
            trigger: "blur"
         }
      ],
      path: [
         {
            required: true,
            message: "路由地址不能为空",
            trigger: "blur"
         }
      ]
   }
});

const { queryParams, form, rules } = toRefs(data);

/** 查询菜单列表 */
function getList() {
   loading.value = true;

   listMenu(queryParams.value).then(response => {
      menuList.value = proxy.handleTree(response.data, "menuId");
      loading.value = false;
   });
}

/** 查询菜单下拉树结构 */
async function getTreeselect() {
   menuOptions.value = [];

   await listMenu().then(response => {
      const menu = {
         menuId: 0,
         menuName: "主类目",
         children: []
      };

      menu.children = proxy.handleTree(response.data, "menuId");
      menuOptions.value.push(menu);
   });
}

/** 取消按钮 */
function cancel() {
   open.value = false;
   reset();
}

/** 表单重置 */
function reset() {
   form.value = {
      menuId: undefined,
      parentId: 0,
      menuName: undefined,
      icon: undefined,
      menuType: "M",
      orderNum: undefined,
      isFrame: "1",
      isCache: "0",
      visible: "0",
      status: "0"
   };

   nextTick(() => {
      menuRef.value?.resetFields?.();
   });
}

/** 图标弹窗打开状态变化 */
function handleIconPopoverOpenChange(visible) {
   if (visible) {
      showSelectIcon();
   } else {
      showChooseIcon.value = false;
   }
}

/** 展示下拉图标 */
function showSelectIcon() {
   iconSelectRef.value?.reset?.();
   showChooseIcon.value = true;
}

/** 选择图标 */
function selected(name) {
   form.value.icon = name;
   showChooseIcon.value = false;
}

/** 搜索按钮操作 */
function handleQuery() {
   getList();
}

/** 重置按钮操作 */
function resetQuery() {
   queryRef.value?.resetFields?.();
   handleQuery();
}

/** 新增按钮操作 */
async function handleAdd(row) {
   reset();

   await getTreeselect();

   if (row != null && row.menuId) {
      form.value.parentId = row.menuId;
   } else {
      form.value.parentId = 0;
   }

   open.value = true;
   title.value = "添加菜单";
}

/** 展开/折叠操作 */
function toggleExpandAll() {
   refreshTable.value = false;
   isExpandAll.value = !isExpandAll.value;

   nextTick(() => {
      refreshTable.value = true;
   });
}

/** 修改按钮操作 */
async function handleUpdate(row) {
   reset();

   await getTreeselect();

   getMenu(row.menuId).then(response => {
      form.value = response.data;
      open.value = true;
      title.value = "修改菜单";
   });
}

/** 提交按钮 */
async function submitForm() {
   try {
      await menuRef.value.validate();

      if (form.value.menuId != undefined) {
         updateMenu(form.value).then(() => {
            proxy.$modal.msgSuccess("修改成功");
            open.value = false;
            getList();
         });
      } else {
         addMenu(form.value).then(() => {
            proxy.$modal.msgSuccess("新增成功");
            open.value = false;
            getList();
         });
      }
   } catch (error) {
      // 表单校验未通过
   }
}

/** 删除按钮操作 */
function handleDelete(row) {
   proxy.$modal
      .confirm('是否确认删除名称为"' + row.menuName + '"的数据项?')
      .then(function () {
         return delMenu(row.menuId);
      })
      .then(() => {
         getList();
         proxy.$modal.msgSuccess("删除成功");
      })
      .catch(() => { });
}

getList();
</script>