<template>
   <div class="app-container">
      <a-form ref="queryRef" :model="queryParams" v-show="showSearch" layout="inline">
         <a-form-item label="用户名称" name="userName">
            <a-input v-model:value="queryParams.userName" placeholder="请输入用户名称" allow-clear style="width: 240px"
               @keyup.enter="handleQuery" />
         </a-form-item>

         <a-form-item label="手机号码" name="phonenumber">
            <a-input v-model:value="queryParams.phonenumber" placeholder="请输入手机号码" allow-clear style="width: 240px"
               @keyup.enter="handleQuery" />
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
            <a-button type="primary" ghost @click="openSelectUser" v-hasPermi="['system:role:add']">
               <template #icon>
                  <PlusOutlined />
               </template>
               添加用户
            </a-button>
         </a-col>

         <a-col>
            <a-button danger ghost :disabled="multiple" @click="cancelAuthUserAll" v-hasPermi="['system:role:remove']">
               <template #icon>
                  <CloseCircleOutlined />
               </template>
               批量取消授权
            </a-button>
         </a-col>

         <a-col>
            <a-button ghost @click="handleClose">
               <template #icon>
                  <CloseOutlined />
               </template>
               关闭
            </a-button>
         </a-col>

         <right-toolbar v-model:showSearch="showSearch" @queryTable="getList" />
      </a-row>

      <a-table :loading="loading" :columns="columns" :data-source="userList" row-key="userId" :pagination="false"
         :row-selection="rowSelection">
         <template #bodyCell="{ column, record }">
            <template v-if="column.dataIndex === 'status'">
               <dict-tag :options="sys_normal_disable" :value="record.status" />
            </template>

            <template v-else-if="column.dataIndex === 'createTime'">
               <span>{{ parseTime(record.createTime) }}</span>
            </template>

            <template v-else-if="column.key === 'action'">
               <a-button type="link" size="small" @click="cancelAuthUser(record)" v-hasPermi="['system:role:remove']">
                  <template #icon>
                     <CloseCircleOutlined />
                  </template>
                  取消授权
               </a-button>
            </template>
         </template>
      </a-table>

      <pagination v-show="total > 0" :total="total" v-model:page="queryParams.pageNum"
         v-model:limit="queryParams.pageSize" @pagination="getList" />

      <select-user ref="selectRef" :roleId="queryParams.roleId" @ok="handleQuery" />
   </div>
</template>

<script setup name="AuthUser">
import {
   SearchOutlined,
   ReloadOutlined,
   PlusOutlined,
   CloseOutlined,
   CloseCircleOutlined
} from "@ant-design/icons-vue";

import selectUser from "./selectUser";
import {
   allocatedUserList,
   authUserCancel,
   authUserCancelAll
} from "@/api/system/role";

const route = useRoute();
const { proxy } = getCurrentInstance();
const { sys_normal_disable } = proxy.useDict("sys_normal_disable");

const queryRef = ref();
const selectRef = ref();

const userList = ref([]);
const loading = ref(true);
const showSearch = ref(true);
const multiple = ref(true);
const total = ref(0);
const userIds = ref([]);

const queryParams = reactive({
   pageNum: 1,
   pageSize: 10,
   roleId: route.params.roleId,
   userName: undefined,
   phonenumber: undefined
});

const columns = [
   {
      title: "用户名称",
      dataIndex: "userName",
      key: "userName",
      ellipsis: true
   },
   {
      title: "用户昵称",
      dataIndex: "nickName",
      key: "nickName",
      ellipsis: true
   },
   {
      title: "邮箱",
      dataIndex: "email",
      key: "email",
      ellipsis: true
   },
   {
      title: "手机",
      dataIndex: "phonenumber",
      key: "phonenumber",
      ellipsis: true
   },
   {
      title: "状态",
      dataIndex: "status",
      key: "status",
      align: "center"
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
      selectedRowKeys: userIds.value,
      onChange: selectedRowKeys => {
         userIds.value = selectedRowKeys;
         multiple.value = !selectedRowKeys.length;
      }
   };
});

/** 查询授权用户列表 */
function getList() {
   loading.value = true;

   allocatedUserList(queryParams).then(response => {
      userList.value = response.data.rows;
      total.value = response.data.total;
      loading.value = false;
   });
}

/** 返回按钮 */
function handleClose() {
   const obj = { path: "/system/role" };
   proxy.$tab.closeOpenPage(obj);
}

/** 搜索按钮操作 */
function handleQuery() {
   queryParams.pageNum = 1;
   getList();
}

/** 重置按钮操作 */
function resetQuery() {
   queryRef.value?.resetFields?.();
   handleQuery();
}

/** 打开授权用户表弹窗 */
function openSelectUser() {
   selectRef.value?.show?.();
}

/** 取消授权按钮操作 */
function cancelAuthUser(row) {
   proxy.$modal
      .confirm('确认要取消该用户"' + row.userName + '"角色吗？')
      .then(function () {
         return authUserCancel({
            userId: row.userId,
            roleId: queryParams.roleId
         });
      })
      .then(() => {
         getList();
         proxy.$modal.msgSuccess("取消授权成功");
      })
      .catch(() => { });
}

/** 批量取消授权按钮操作 */
function cancelAuthUserAll() {
   const roleId = queryParams.roleId;
   const uIds = userIds.value.join(",");

   proxy.$modal
      .confirm("是否取消选中用户授权数据项?")
      .then(function () {
         return authUserCancelAll({
            roleId,
            userIds: uIds
         });
      })
      .then(() => {
         getList();
         proxy.$modal.msgSuccess("取消授权成功");
      })
      .catch(() => { });
}

getList();
</script>