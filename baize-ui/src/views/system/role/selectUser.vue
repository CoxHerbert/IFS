<template>
   <!-- 授权用户 -->
   <a-modal v-model:open="visible" title="选择用户" width="800px" :style="{ top: '5vh' }" :mask-closable="false"
      @ok="handleSelectUser" @cancel="visible = false">
      <a-form ref="queryRef" :model="queryParams" layout="inline">
         <a-form-item label="用户名称" name="userName">
            <a-input v-model:value="queryParams.userName" placeholder="请输入用户名称" allow-clear
               @keyup.enter="handleQuery" />
         </a-form-item>

         <a-form-item label="手机号码" name="phonenumber">
            <a-input v-model:value="queryParams.phonenumber" placeholder="请输入手机号码" allow-clear
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

      <a-table style="margin-top: 12px" :columns="columns" :data-source="userList" row-key="userId" :pagination="false"
         :row-selection="rowSelection" :scroll="{ y: 260 }" @row="handleRow">
         <template #bodyCell="{ column, record }">
            <template v-if="column.dataIndex === 'status'">
               <dict-tag :options="sys_normal_disable" :value="record.status" />
            </template>

            <template v-else-if="column.dataIndex === 'createTime'">
               <span>{{ parseTime(record.createTime) }}</span>
            </template>
         </template>
      </a-table>

      <pagination v-show="total > 0" :total="total" v-model:page="queryParams.pageNum"
         v-model:limit="queryParams.pageSize" @pagination="getList" />

      <template #footer>
         <a-space>
            <a-button @click="visible = false">取 消</a-button>
            <a-button type="primary" @click="handleSelectUser">确 定</a-button>
         </a-space>
      </template>
   </a-modal>
</template>

<script setup name="SelectUser">
import {
   SearchOutlined,
   ReloadOutlined
} from "@ant-design/icons-vue";

import {
   authUserSelectAll,
   unallocatedUserList
} from "@/api/system/role";

const props = defineProps({
   roleId: {
      type: [Number, String]
   }
});

const emit = defineEmits(["ok"]);

const { proxy } = getCurrentInstance();
const { sys_normal_disable } = proxy.useDict("sys_normal_disable");

const queryRef = ref();

const userList = ref([]);
const visible = ref(false);
const total = ref(0);
const userIds = ref([]);

const queryParams = reactive({
   pageNum: 1,
   pageSize: 10,
   roleId: undefined,
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
   }
];

const rowSelection = computed(() => {
   return {
      selectedRowKeys: userIds.value,
      onChange: selectedRowKeys => {
         userIds.value = selectedRowKeys;
      }
   };
});

/** 显示弹框 */
function show() {
   queryParams.roleId = props.roleId;
   queryParams.pageNum = 1;
   userIds.value = [];
   getList();
   visible.value = true;
}

/** 点击行切换选择 */
function handleRow(record) {
   return {
      onClick: () => {
         const index = userIds.value.indexOf(record.userId);

         if (index > -1) {
            userIds.value.splice(index, 1);
         } else {
            userIds.value.push(record.userId);
         }
      }
   };
}

/** 查询表数据 */
function getList() {
   unallocatedUserList(queryParams).then(res => {
      userList.value = res.data.rows;
      total.value = res.data.total;
   });
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

/** 选择授权用户操作 */
function handleSelectUser() {
   const roleId = queryParams.roleId;
   const uIds = userIds.value.join(",");

   if (uIds == "") {
      proxy.$modal.msgError("请选择要分配的用户");
      return;
   }

   authUserSelectAll({
      roleId,
      userIds: uIds
   }).then(res => {
      proxy.$modal.msgSuccess(res.msg);

      if (res.code === 200) {
         visible.value = false;
         emit("ok");
      }
   });
}

defineExpose({
   show
});
</script>