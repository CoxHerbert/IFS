<template>
   <div class="app-container">
      <a-form v-show="showSearch" ref="queryRef" :model="queryParams" layout="inline"
         :label-col="{ style: { width: '68px' } }" class="query-form">
         <a-form-item label="任务名称" name="jobName">
            <a-input v-model:value="queryParams.jobName" placeholder="请输入任务名称" allow-clear @pressEnter="handleQuery" />
         </a-form-item>

         <a-form-item label="任务组名" name="jobGroup">
            <a-select v-model:value="queryParams.jobGroup" placeholder="请选择任务组名" allow-clear style="width: 180px">
               <a-select-option v-for="dict in sys_job_group" :key="dict.value" :value="dict.value">
                  {{ dict.label }}
               </a-select-option>
            </a-select>
         </a-form-item>

         <a-form-item label="任务状态" name="status">
            <a-select v-model:value="queryParams.status" placeholder="请选择任务状态" allow-clear style="width: 180px">
               <a-select-option v-for="dict in sys_job_status" :key="dict.value" :value="dict.value">
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

      <a-row :gutter="10" class="mb8 toolbar-row">
         <a-col>
            <a-button type="primary" ghost @click="handleAdd" v-hasPermi="['monitor:job:add']">
               <template #icon>
                  <PlusOutlined />
               </template>
               新增
            </a-button>
         </a-col>

         <a-col>
            <a-button type="primary" ghost :disabled="single" @click="handleUpdate" v-hasPermi="['monitor:job:edit']">
               <template #icon>
                  <EditOutlined />
               </template>
               修改
            </a-button>
         </a-col>

         <a-col>
            <a-button danger ghost :disabled="multiple" @click="handleDelete" v-hasPermi="['monitor:job:remove']">
               <template #icon>
                  <DeleteOutlined />
               </template>
               删除
            </a-button>
         </a-col>

         <a-col>
            <a-button ghost @click="handleExport" v-hasPermi="['monitor:job:export']">
               <template #icon>
                  <DownloadOutlined />
               </template>
               导出
            </a-button>
         </a-col>

         <a-col>
            <a-button ghost @click="handleJobLog" v-hasPermi="['monitor:job:query']">
               <template #icon>
                  <ProfileOutlined />
               </template>
               日志
            </a-button>
         </a-col>

         <a-col flex="auto">
            <right-toolbar v-model:showSearch="showSearch" @queryTable="getList" />
         </a-col>
      </a-row>

      <a-table :loading="loading" :data-source="jobList" :columns="columns" :row-key="record => record.jobId"
         :row-selection="rowSelection" :pagination="false" bordered size="middle">
         <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'jobGroup'">
               <dict-tag :options="sys_job_group" :value="record.jobGroup" />
            </template>

            <template v-else-if="column.key === 'status'">
               <a-switch :checked="record.status === '0'" checked-value="0" un-checked-value="1"
                  @change="checked => handleStatusChange(record, checked)" />
            </template>

            <template v-else-if="column.key === 'operation'">
               <a-space>
                  <a-button type="link" size="small" @click="handleUpdate(record)" v-hasPermi="['monitor:job:edit']">
                     <template #icon>
                        <EditOutlined />
                     </template>
                     修改
                  </a-button>

                  <a-button type="link" size="small" danger @click="handleDelete(record)"
                     v-hasPermi="['monitor:job:remove']">
                     <template #icon>
                        <DeleteOutlined />
                     </template>
                     删除
                  </a-button>

                  <a-dropdown>
                     <a-button type="link" size="small" v-hasPermi="['monitor:job:changeStatus', 'monitor:job:query']">
                        更多
                        <DownOutlined />
                     </a-button>

                     <template #overlay>
                        <a-menu @click="({ key }) => handleCommand(key, record)">
                           <a-menu-item key="handleRun" v-hasPermi="['monitor:job:changeStatus']">
                              <CaretRightOutlined />
                              执行一次
                           </a-menu-item>
                           <a-menu-item key="handleView" v-hasPermi="['monitor:job:query']">
                              <EyeOutlined />
                              任务详细
                           </a-menu-item>
                           <a-menu-item key="handleJobLog" v-hasPermi="['monitor:job:query']">
                              <ProfileOutlined />
                              调度日志
                           </a-menu-item>
                        </a-menu>
                     </template>
                  </a-dropdown>
               </a-space>
            </template>
         </template>
      </a-table>

      <div v-if="total > 0" class="pagination-wrapper">
         <a-pagination v-model:current="queryParams.pageNum" v-model:page-size="queryParams.pageSize" :total="total"
            show-size-changer :show-total="total => `共 ${total} 条`" @change="getList" @showSizeChange="getList" />
      </div>

      <!-- 添加或修改定时任务对话框 -->
      <a-modal v-model:open="open" :title="title" width="800px" @ok="submitForm" @cancel="cancel">
         <a-form ref="jobRef" :model="form" :rules="rules" :label-col="{ style: { width: '120px' } }">
            <a-row :gutter="16">
               <a-col :span="12">
                  <a-form-item label="任务名称" name="jobName">
                     <a-input v-model:value="form.jobName" placeholder="请输入任务名称" />
                  </a-form-item>
               </a-col>

               <a-col :span="12">
                  <a-form-item label="任务分组" name="jobGroup">
                     <a-select v-model:value="form.jobGroup" placeholder="请选择" style="width: 100%">
                        <a-select-option v-for="dict in sys_job_group" :key="dict.value" :value="dict.value">
                           {{ dict.label }}
                        </a-select-option>
                     </a-select>
                  </a-form-item>
               </a-col>

               <a-col :span="24">
                  <a-form-item name="invokeTarget">
                     <template #label>
                        <span>调用方法</span>
                     </template>
                     <a-input v-model:value="form.invokeTarget" placeholder="请输入调用目标字符串" />
                  </a-form-item>
               </a-col>

               <a-col :span="24">
                  <a-form-item label="cron表达式" name="cronExpression">
                     <a-input v-model:value="form.cronExpression" placeholder="请输入cron执行表达式">
                        <template #addonAfter>
                           <a-button type="primary" @click="handleShowCron">
                              生成表达式
                              <ClockCircleOutlined />
                           </a-button>
                        </template>
                     </a-input>
                  </a-form-item>
               </a-col>

               <a-col :span="24">
                  <a-form-item>
                     <template #label>
                        <span>参数</span>
                     </template>
                     <a-input v-model:value="form.jobParams" placeholder="请输入方法参数" />
                  </a-form-item>
               </a-col>

               <a-col :span="12">
                  <a-form-item label="状态">
                     <a-radio-group v-model:value="form.status">
                        <a-radio v-for="dict in sys_job_status" :key="dict.value" :value="dict.value">
                           {{ dict.label }}
                        </a-radio>
                     </a-radio-group>
                  </a-form-item>
               </a-col>
            </a-row>
         </a-form>

         <template #footer>
            <a-button @click="cancel">取 消</a-button>
            <a-button type="primary" @click="submitForm">确 定</a-button>
         </template>
      </a-modal>

      <!-- 任务日志详细 -->
      <a-modal v-model:open="openView" title="任务详细" width="700px" :footer="null">
         <a-form :model="form" :label-col="{ style: { width: '120px' } }" size="small">
            <a-row :gutter="16">
               <a-col :span="12">
                  <a-form-item label="任务编号：">{{ form.jobId }}</a-form-item>
                  <a-form-item label="任务名称：">{{ form.jobName }}</a-form-item>
               </a-col>

               <a-col :span="12">
                  <a-form-item label="任务分组：">{{ jobGroupFormat(form) }}</a-form-item>
                  <a-form-item label="创建时间：">{{ form.createTime }}</a-form-item>
               </a-col>

               <a-col :span="12">
                  <a-form-item label="cron表达式：">{{ form.cronExpression }}</a-form-item>
               </a-col>

               <a-col :span="12">
                  <a-form-item label="下次执行时间：">{{ parseTime(form.nextValidTime) }}</a-form-item>
               </a-col>

               <a-col :span="24">
                  <a-form-item label="调用目标方法：">{{ form.invokeTarget }}</a-form-item>
               </a-col>

               <a-col :span="12">
                  <a-form-item label="任务状态：">
                     <div v-if="form.status == 0">正常</div>
                     <div v-else-if="form.status == 1">失败</div>
                  </a-form-item>
               </a-col>

               <a-col :span="12">
                  <a-form-item label="是否并发：">
                     <div v-if="form.concurrent == 0">允许</div>
                     <div v-else-if="form.concurrent == 1">禁止</div>
                  </a-form-item>
               </a-col>

               <a-col :span="12">
                  <a-form-item label="执行策略：">
                     <div v-if="form.misfirePolicy == 0">默认策略</div>
                     <div v-else-if="form.misfirePolicy == 1">立即执行</div>
                     <div v-else-if="form.misfirePolicy == 2">执行一次</div>
                     <div v-else-if="form.misfirePolicy == 3">放弃执行</div>
                  </a-form-item>
               </a-col>
            </a-row>
         </a-form>

         <template #footer>
            <a-button @click="openView = false">关 闭</a-button>
         </template>
      </a-modal>
   </div>
</template>

<script setup name="Job">
import { ref, reactive, toRefs, computed, getCurrentInstance } from "vue";
import { useRouter } from "vue-router";
import { Modal, message } from "ant-design-vue";
import {
   SearchOutlined,
   ReloadOutlined,
   PlusOutlined,
   EditOutlined,
   DeleteOutlined,
   DownloadOutlined,
   ProfileOutlined,
   DownOutlined,
   CaretRightOutlined,
   EyeOutlined,
   ClockCircleOutlined
} from "@ant-design/icons-vue";
import {
   listJob,
   getJob,
   delJob,
   addJob,
   updateJob,
   runJob,
   changeJobStatus
} from "@/api/monitor/job";

const router = useRouter();
const { proxy } = getCurrentInstance();
const { sys_job_group, sys_job_status } = proxy.useDict("sys_job_group", "sys_job_status");

const jobList = ref([]);
const open = ref(false);
const loading = ref(true);
const showSearch = ref(true);
const ids = ref([]);
const single = ref(true);
const multiple = ref(true);
const total = ref(0);
const title = ref("");
const openView = ref(false);
const openCron = ref(false);
const expression = ref("");
const queryRef = ref();
const jobRef = ref();

const data = reactive({
   form: {},
   queryParams: {
      pageNum: 1,
      pageSize: 10,
      jobName: undefined,
      jobGroup: undefined,
      status: undefined
   },
   rules: {
      jobName: [{ required: true, message: "任务名称不能为空", trigger: "blur" }],
      invokeTarget: [{ required: true, message: "调用目标字符串不能为空", trigger: "blur" }],
      cronExpression: [{ required: true, message: "cron执行表达式不能为空", trigger: "blur" }]
   }
});

const { queryParams, form, rules } = toRefs(data);

const columns = [
   {
      title: "任务编号",
      dataIndex: "jobId",
      key: "jobId",
      width: 100,
      align: "center"
   },
   {
      title: "任务名称",
      dataIndex: "jobName",
      key: "jobName",
      align: "center",
      ellipsis: true
   },
   {
      title: "任务组名",
      dataIndex: "jobGroup",
      key: "jobGroup",
      align: "center"
   },
   {
      title: "调用目标字符串",
      dataIndex: "invokeTarget",
      key: "invokeTarget",
      align: "center",
      ellipsis: true
   },
   {
      title: "cron执行表达式",
      dataIndex: "cronExpression",
      key: "cronExpression",
      align: "center",
      ellipsis: true
   },
   {
      title: "参数",
      dataIndex: "jobParams",
      key: "jobParams",
      align: "center",
      ellipsis: true
   },
   {
      title: "状态",
      dataIndex: "status",
      key: "status",
      align: "center"
   },
   {
      title: "操作",
      key: "operation",
      align: "center",
      width: 220,
      fixed: "right"
   }
];

const rowSelection = computed(() => ({
   selectedRowKeys: ids.value,
   onChange: (selectedRowKeys, selectedRows) => {
      ids.value = selectedRowKeys;
      single.value = selectedRows.length !== 1;
      multiple.value = !selectedRows.length;
   }
}));

/** 查询定时任务列表 */
function getList() {
   loading.value = true;

   listJob(queryParams.value)
      .then(response => {
         jobList.value = response.data.rows || [];
         total.value = response.data.total || 0;
      })
      .finally(() => {
         loading.value = false;
      });
}

/** 任务组名字典翻译 */
function jobGroupFormat(row) {
   return proxy.selectDictLabel(sys_job_group, row.jobGroup);
}

/** 取消按钮 */
function cancel() {
   open.value = false;
   reset();
}

/** 表单重置 */
function reset() {
   form.value = {
      jobId: undefined,
      jobName: undefined,
      jobGroup: undefined,
      invokeTarget: undefined,
      cronExpression: undefined,
      jobParams: undefined,
      status: "0"
   };

   jobRef.value?.resetFields();
}

/** 搜索按钮操作 */
function handleQuery() {
   queryParams.value.pageNum = 1;
   getList();
}

/** 重置按钮操作 */
function resetQuery() {
   queryRef.value?.resetFields();
   handleQuery();
}

// 更多操作触发
function handleCommand(command, row) {
   switch (command) {
      case "handleRun":
         handleRun(row);
         break;
      case "handleView":
         handleView(row);
         break;
      case "handleJobLog":
         handleJobLog(row);
         break;
      default:
         break;
   }
}

// 任务状态修改
function handleStatusChange(row, checkedValue) {
   const oldStatus = row.status;
   const newStatus = checkedValue;
   row.status = newStatus;

   const text = newStatus === "0" ? "启用" : "停用";

   Modal.confirm({
      title: "确认操作",
      content: `确认要"${text}""${row.jobName}"任务吗?`,
      okText: "确定",
      cancelText: "取消",
      onOk() {
         return changeJobStatus(row.jobId, newStatus).then(() => {
            message.success(text + "成功");
         });
      },
      onCancel() {
         row.status = oldStatus;
      }
   });
}

/* 立即执行一次 */
function handleRun(row) {
   Modal.confirm({
      title: "确认执行",
      content: `确认要立即执行一次"${row.jobName}"任务吗?`,
      okText: "确定",
      cancelText: "取消",
      onOk() {
         return runJob(row.jobId, row.jobGroup).then(() => {
            message.success("执行成功");
         });
      }
   });
}

/** 任务详细信息 */
function handleView(row) {
   getJob(row.jobId).then(response => {
      form.value = response.data || {};
      openView.value = true;
   });
}

/** cron表达式按钮操作 */
function handleShowCron() {
   expression.value = form.value.cronExpression;
   openCron.value = true;
}

/** 确定后回传值 */
function crontabFill(value) {
   form.value.cronExpression = value;
}

/** 任务日志列表查询 */
function handleJobLog(row = {}) {
   const jobId = row.jobId || 0;
   router.push({ path: "/monitor/job-log/index", query: { jobId } });
}

/** 新增按钮操作 */
function handleAdd() {
   reset();
   open.value = true;
   title.value = "添加任务";
}

/** 修改按钮操作 */
function handleUpdate(row = {}) {
   reset();
   const jobId = row.jobId || ids.value[0];

   getJob(jobId).then(response => {
      form.value = response.data || {};
      open.value = true;
      title.value = "修改任务";
   });
}

/** 提交按钮 */
function submitForm() {
   jobRef.value
      ?.validate()
      .then(() => {
         if (form.value.jobId !== undefined) {
            updateJob(form.value).then(() => {
               message.success("修改成功");
               open.value = false;
               getList();
            });
         } else {
            addJob(form.value).then(() => {
               message.success("新增成功");
               open.value = false;
               getList();
            });
         }
      })
      .catch(() => { });
}

/** 删除按钮操作 */
function handleDelete(row = {}) {
   const jobIds = row.jobId || ids.value;

   Modal.confirm({
      title: "确认删除",
      content: `是否确认删除定时任务编号为"${jobIds}"的数据项?`,
      okText: "确定",
      cancelText: "取消",
      onOk() {
         return delJob(jobIds).then(() => {
            getList();
            message.success("删除成功");
         });
      }
   });
}

/** 导出按钮操作 */
function handleExport() {
   proxy.download(
      "monitor/job/export",
      {
         ...queryParams.value
      },
      `job_${new Date().getTime()}.xlsx`
   );
}

getList();
</script>

<style scoped>
.query-form {
   margin-bottom: 16px;
}

.toolbar-row {
   margin-bottom: 8px;
}

.pagination-wrapper {
   display: flex;
   justify-content: flex-end;
   margin-top: 16px;
}
</style>
