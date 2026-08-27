<template>
   <div class="app-container">
      <a-form ref="queryRef" :model="queryParams" layout="inline" v-show="showSearch">
         <a-form-item label="部门名称" name="deptName">
            <a-input v-model:value="queryParams.deptName" placeholder="请输入部门名称" allow-clear
               @keyup.enter="handleQuery" />
         </a-form-item>

         <a-form-item label="状态" name="status">
            <a-select v-model:value="queryParams.status" placeholder="部门状态" allow-clear style="width: 180px">
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
            <a-button type="primary" ghost @click="handleAdd" v-hasPermi="['system:dept:add']">
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

      <a-table v-if="refreshTable" :loading="loading" :columns="columns" :data-source="deptList" row-key="deptId"
         :pagination="false" :default-expand-all-rows="isExpandAll" :children-column-name="'children'">
         <template #bodyCell="{ column, record }">
            <template v-if="column.dataIndex === 'status'">
               <dict-tag :options="sys_normal_disable" :value="record.status" />
            </template>

            <template v-else-if="column.dataIndex === 'createTime'">
               <span>{{ parseTime(record.createTime) }}</span>
            </template>

            <template v-else-if="column.key === 'action'">
               <a-space>
                  <a-button type="link" size="small" @click="handleUpdate(record)" v-hasPermi="['system:dept:edit']">
                     <template #icon>
                        <EditOutlined />
                     </template>
                     修改
                  </a-button>

                  <a-button type="link" size="small" @click="handleAdd(record)" v-hasPermi="['system:dept:add']">
                     <template #icon>
                        <PlusOutlined />
                     </template>
                     新增
                  </a-button>

                  <a-button v-if="record.parentId != 0" type="link" size="small" danger @click="handleDelete(record)"
                     v-hasPermi="['system:dept:remove']">
                     <template #icon>
                        <DeleteOutlined />
                     </template>
                     删除
                  </a-button>
               </a-space>
            </template>
         </template>
      </a-table>

      <!-- 添加或修改部门对话框 -->
      <a-modal v-model:open="open" :title="title" width="600px" :mask-closable="false" @ok="submitForm"
         @cancel="cancel">
         <a-form ref="deptRef" :model="form" :rules="rules" :label-col="{ style: { width: '80px' } }">
            <a-row :gutter="16">
               <a-col :span="24" v-if="form.parentId != 0">
                  <a-form-item label="上级部门" name="parentId">
                     <a-tree-select v-model:value="form.parentId" :tree-data="deptOptions" :field-names="{
                        value: 'deptId',
                        label: 'deptName',
                        children: 'children'
                     }" placeholder="选择上级部门" allow-clear tree-default-expand-all style="width: 100%" />
                  </a-form-item>
               </a-col>

               <a-col :span="12">
                  <a-form-item label="部门名称" name="deptName">
                     <a-input v-model:value="form.deptName" placeholder="请输入部门名称" />
                  </a-form-item>
               </a-col>

               <a-col :span="12">
                  <a-form-item label="显示排序" name="orderNum">
                     <a-input-number v-model:value="form.orderNum" :min="0" style="width: 100%" />
                  </a-form-item>
               </a-col>

               <a-col :span="12">
                  <a-form-item label="负责人" name="leader">
                     <a-input v-model:value="form.leader" placeholder="请输入负责人" :maxlength="20" />
                  </a-form-item>
               </a-col>

               <a-col :span="12">
                  <a-form-item label="联系电话" name="phone">
                     <a-input v-model:value="form.phone" placeholder="请输入联系电话" :maxlength="11" />
                  </a-form-item>
               </a-col>

               <a-col :span="12">
                  <a-form-item label="邮箱" name="email">
                     <a-input v-model:value="form.email" placeholder="请输入邮箱" :maxlength="50" />
                  </a-form-item>
               </a-col>

               <a-col :span="12">
                  <a-form-item label="部门状态" name="status">
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

<script setup name="Dept">
import {
   SearchOutlined,
   ReloadOutlined,
   PlusOutlined,
   SwapOutlined,
   EditOutlined,
   DeleteOutlined
} from "@ant-design/icons-vue";

import {
   listDept,
   getDept,
   delDept,
   addDept,
   updateDept
} from "@/api/system/dept";

const { proxy } = getCurrentInstance();
const { sys_normal_disable } = proxy.useDict("sys_normal_disable");

const queryRef = ref();
const deptRef = ref();

const deptList = ref([]);
const open = ref(false);
const loading = ref(true);
const showSearch = ref(true);
const title = ref("");
const deptOptions = ref([]);
const isExpandAll = ref(true);
const refreshTable = ref(true);

const columns = [
   {
      title: "部门名称",
      dataIndex: "deptName",
      key: "deptName",
      width: 260
   },
   {
      title: "排序",
      dataIndex: "orderNum",
      key: "orderNum",
      width: 200
   },
   {
      title: "状态",
      dataIndex: "status",
      key: "status",
      width: 100
   },
   {
      title: "创建时间",
      dataIndex: "createTime",
      key: "createTime",
      align: "center",
      width: 200
   },
   {
      title: "操作",
      key: "action",
      align: "center"
   }
];

const data = reactive({
   form: {},
   queryParams: {
      deptName: undefined,
      status: undefined
   },
   rules: {
      parentId: [
         {
            required: true,
            message: "上级部门不能为空",
            trigger: "change"
         }
      ],
      deptName: [
         {
            required: true,
            message: "部门名称不能为空",
            trigger: "blur"
         }
      ],
      orderNum: [
         {
            required: true,
            message: "显示排序不能为空",
            trigger: "blur"
         }
      ],
      email: [
         {
            type: "email",
            message: "请输入正确的邮箱地址",
            trigger: ["blur", "change"]
         }
      ],
      phone: [
         {
            pattern: /^1[3|4|5|6|7|8|9][0-9]\d{8}$/,
            message: "请输入正确的手机号码",
            trigger: "blur"
         }
      ]
   }
});

const { queryParams, form, rules } = toRefs(data);

/** 查询部门列表 */
function getList() {
   loading.value = true;
   listDept(queryParams.value).then(response => {
      deptList.value = proxy.handleTree(response.data, "deptId");
      loading.value = false;
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
      deptId: undefined,
      parentId: undefined,
      deptName: undefined,
      orderNum: 0,
      leader: undefined,
      phone: undefined,
      email: undefined,
      status: "0"
   };

   nextTick(() => {
      deptRef.value?.resetFields?.();
   });
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

   await listDept().then(response => {
      deptOptions.value = proxy.handleTree(response.data, "deptId");
   });

   if (row != undefined) {
      form.value.parentId = row.deptId;
   }

   open.value = true;
   title.value = "添加部门";
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

   await listDept(row.deptId).then(response => {
      const list = response.data.filter(dept => {
         return !(
            dept.deptId == row.deptId ||
            dept.ancestors.split(",").indexOf(row.deptId) != -1
         );
      });

      deptOptions.value = proxy.handleTree(list, "deptId");
   });

   getDept(row.deptId).then(response => {
      form.value = response.data;
      open.value = true;
      title.value = "修改部门";
   });
}

/** 提交按钮 */
async function submitForm() {
   try {
      await deptRef.value.validate();

      if (form.value.deptId != undefined) {
         updateDept(form.value).then(() => {
            proxy.$modal.msgSuccess("修改成功");
            open.value = false;
            getList();
         });
      } else {
         addDept(form.value).then(() => {
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
      .confirm('是否确认删除名称为"' + row.deptName + '"的数据项?')
      .then(function () {
         return delDept(row.deptId);
      })
      .then(() => {
         getList();
         proxy.$modal.msgSuccess("删除成功");
      })
      .catch(() => { });
}

getList();
</script>