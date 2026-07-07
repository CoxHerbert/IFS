<template>
   <div class="app-container">
      <h4 class="form-header h4">基本信息</h4>

      <a-form :model="form" :label-col="{ style: { width: '80px' } }">
         <a-row>
            <a-col :span="8" :offset="2">
               <a-form-item label="用户昵称" name="nickName">
                  <a-input v-model:value="form.nickName" disabled />
               </a-form-item>
            </a-col>

            <a-col :span="8" :offset="2">
               <a-form-item label="登录账号" name="userName">
                  <a-input v-model:value="form.userName" disabled />
               </a-form-item>
            </a-col>
         </a-row>
      </a-form>

      <h4 class="form-header h4">角色信息</h4>

      <a-table :loading="loading" :columns="columns" :data-source="pageRoles" row-key="roleId" :pagination="false"
         :row-selection="rowSelection" @row="handleRow">
         <template #bodyCell="{ column, record, index }">
            <template v-if="column.key === 'index'">
               <span>{{ (pageNum - 1) * pageSize + index + 1 }}</span>
            </template>

            <template v-else-if="column.dataIndex === 'createTime'">
               <span>{{ parseTime(record.createTime) }}</span>
            </template>
         </template>
      </a-table>

      <pagination v-show="total > 0" :total="total" v-model:page="pageNum" v-model:limit="pageSize" />

      <a-form :label-col="{ style: { width: '100px' } }">
         <a-form-item style="text-align: center; margin-left: -120px; margin-top: 30px">
            <a-space>
               <a-button type="primary" @click="submitForm">
                  提交
               </a-button>

               <a-button @click="close">
                  返回
               </a-button>
            </a-space>
         </a-form-item>
      </a-form>
   </div>
</template>

<script setup name="AuthRole">
import { getAuthRole, updateAuthRole } from "@/api/system/user";

const route = useRoute();
const { proxy } = getCurrentInstance();

const loading = ref(true);
const total = ref(0);
const pageNum = ref(1);
const pageSize = ref(10);
const roleIds = ref([]);
const roles = ref([]);

const form = ref({
   nickName: undefined,
   userName: undefined,
   userId: undefined
});

const columns = [
   {
      title: "序号",
      key: "index",
      align: "center",
      width: 80
   },
   {
      title: "角色编号",
      dataIndex: "roleId",
      key: "roleId",
      align: "center"
   },
   {
      title: "角色名称",
      dataIndex: "roleName",
      key: "roleName",
      align: "center"
   },
   {
      title: "权限字符",
      dataIndex: "roleKey",
      key: "roleKey",
      align: "center"
   },
   {
      title: "创建时间",
      dataIndex: "createTime",
      key: "createTime",
      align: "center",
      width: 180
   }
];

const pageRoles = computed(() => {
   return roles.value.slice(
      (pageNum.value - 1) * pageSize.value,
      pageNum.value * pageSize.value
   );
});

const rowSelection = computed(() => {
   return {
      selectedRowKeys: roleIds.value,
      preserveSelectedRowKeys: true,
      onChange: selectedRowKeys => {
         roleIds.value = selectedRowKeys;
      }
   };
});

/** 单击选中行数据 */
function handleRow(record) {
   return {
      onClick: () => {
         const index = roleIds.value.indexOf(record.roleId);

         if (index > -1) {
            roleIds.value.splice(index, 1);
         } else {
            roleIds.value.push(record.roleId);
         }
      }
   };
}

/** 关闭按钮 */
function close() {
   const obj = { path: "/system/user" };
   proxy.$tab.closeOpenPage(obj);
}

/** 提交按钮 */
function submitForm() {
   const userId = form.value.userId;
   const rIds = roleIds.value.join(",");

   updateAuthRole({
      userId,
      roleIds: rIds
   }).then(() => {
      proxy.$modal.msgSuccess("授权成功");
      close();
   });
}

function init() {
   const userId = route.params && route.params.userId;

   if (userId) {
      loading.value = true;

      getAuthRole(userId).then(response => {
         form.value = response.data.user;
         roles.value = response.data.roles;
         total.value = roles.value.length;
         roleIds.value = response.data.roleIds || [];
         loading.value = false;
      });
   }
}

init();
</script>