<template>
  <div class="app-container">
    <a-form v-show="showSearch" ref="queryRef" :model="queryParams" layout="inline"
      :label-col="{ style: { width: '68px' } }" class="query-form">
      <a-form-item label="表名称" name="tableName">
        <a-input v-model:value="queryParams.tableName" placeholder="请输入表名称" allow-clear @pressEnter="handleQuery" />
      </a-form-item>

      <a-form-item label="表描述" name="tableComment">
        <a-input v-model:value="queryParams.tableComment" placeholder="请输入表描述" allow-clear @pressEnter="handleQuery" />
      </a-form-item>

      <a-form-item label="创建时间">
        <a-range-picker v-model:value="dateRange" value-format="YYYY-MM-DD" format="YYYY-MM-DD" style="width: 240px" />
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
        <a-button type="primary" ghost @click="handleGenTable" v-hasPermi="['tool:gen:code']">
          <template #icon>
            <DownloadOutlined />
          </template>
          生成
        </a-button>
      </a-col>

      <a-col>
        <a-button ghost @click="openImportTable" v-hasPermi="['tool:gen:import']">
          <template #icon>
            <UploadOutlined />
          </template>
          导入
        </a-button>
      </a-col>

      <a-col>
        <a-button type="primary" ghost :disabled="single" @click="handleEditTable" v-hasPermi="['tool:gen:edit']">
          <template #icon>
            <EditOutlined />
          </template>
          修改
        </a-button>
      </a-col>

      <a-col>
        <a-button danger ghost :disabled="multiple" @click="handleDelete" v-hasPermi="['tool:gen:remove']">
          <template #icon>
            <DeleteOutlined />
          </template>
          删除
        </a-button>
      </a-col>

      <a-col flex="auto">
        <right-toolbar v-model:showSearch="showSearch" @queryTable="getList" />
      </a-col>
    </a-row>

    <a-table :loading="loading" :data-source="tableList" :columns="columns" :row-key="record => record.tableId"
      :row-selection="rowSelection" :pagination="false" bordered size="middle">
      <template #bodyCell="{ column, record, index }">
        <template v-if="column.key === 'index'">
          {{ (queryParams.pageNum - 1) * queryParams.pageSize + index + 1 }}
        </template>

        <template v-else-if="column.key === 'createTime'">
          {{ parseTime(record.createTime, '{y}-{m}-{d}') }}
        </template>

        <template v-else-if="column.key === 'updateTime'">
          {{ parseTime(record.updateTime, '{y}-{m}-{d}') }}
        </template>

        <template v-else-if="column.key === 'operation'">
          <a-space>
            <a-button type="link" size="small" @click="handlePreview(record)" v-hasPermi="['tool:gen:preview']">
              <template #icon>
                <EyeOutlined />
              </template>
              预览
            </a-button>

            <a-button type="link" size="small" @click="handleEditTable(record)" v-hasPermi="['tool:gen:edit']">
              <template #icon>
                <EditOutlined />
              </template>
              编辑
            </a-button>

            <a-button type="link" size="small" danger @click="handleDelete(record)" v-hasPermi="['tool:gen:remove']">
              <template #icon>
                <DeleteOutlined />
              </template>
              删除
            </a-button>

            <a-button type="link" size="small" @click="handleSynchDb(record)" v-hasPermi="['tool:gen:edit']">
              <template #icon>
                <ReloadOutlined />
              </template>
              同步
            </a-button>

            <a-button type="link" size="small" @click="handleGenTable(record)" v-hasPermi="['tool:gen:code']">
              <template #icon>
                <DownloadOutlined />
              </template>
              生成代码
            </a-button>
          </a-space>
        </template>
      </template>
    </a-table>

    <div v-if="total > 0" class="pagination-wrapper">
      <a-pagination v-model:current="queryParams.pageNum" v-model:page-size="queryParams.pageSize" :total="total"
        show-size-changer :show-total="total => `共 ${total} 条`" @change="getList" @showSizeChange="getList" />
    </div>

    <!-- 预览界面 -->
    <a-modal v-model:open="preview.open" :title="preview.title" width="80%" :footer="null" class="preview-modal">
      <a-tabs v-model:activeKey="preview.activeName">
        <a-tab-pane v-for="(value, key) in preview.data" :key="getPreviewName(key)" :tab="getPreviewName(key)">
          <pre><code class="hljs" v-html="value"></code></pre>
        </a-tab-pane>
      </a-tabs>
    </a-modal>

    <import-table ref="importRef" @ok="handleQuery" />
  </div>
</template>

<script setup name="Gen">
import { ref, reactive, toRefs, computed, getCurrentInstance } from "vue";
import { useRoute } from "vue-router";
import { Modal, message } from "ant-design-vue";
import {
  SearchOutlined,
  ReloadOutlined,
  DownloadOutlined,
  UploadOutlined,
  EditOutlined,
  DeleteOutlined,
  EyeOutlined
} from "@ant-design/icons-vue";

import {
  listTable,
  previewTable,
  delTable,
  genCode,
  synchDb
} from "@/api/tool/gen";
import router from "@/router";
import importTable from "./importTable";

const route = useRoute();
const { proxy } = getCurrentInstance();

const tableList = ref([]);
const loading = ref(true);
const showSearch = ref(true);
const ids = ref([]);
const single = ref(true);
const multiple = ref(true);
const total = ref(0);
const tableNames = ref([]);
const dateRange = ref([]);
const queryRef = ref();
const importRef = ref();

const data = reactive({
  queryParams: {
    pageNum: 1,
    pageSize: 10,
    tableName: undefined,
    tableComment: undefined
  },
  preview: {
    open: false,
    title: "代码预览",
    data: {},
    activeName: "domain.java"
  }
});

const { queryParams, preview } = toRefs(data);

const columns = [
  {
    title: "序号",
    key: "index",
    width: 70,
    align: "center"
  },
  {
    title: "表名称",
    dataIndex: "tableName",
    key: "tableName",
    align: "center",
    ellipsis: true
  },
  {
    title: "表描述",
    dataIndex: "tableComment",
    key: "tableComment",
    align: "center",
    ellipsis: true
  },
  {
    title: "实体",
    dataIndex: "className",
    key: "className",
    align: "center",
    ellipsis: true
  },
  {
    title: "创建时间",
    dataIndex: "createTime",
    key: "createTime",
    align: "center",
    width: 150
  },
  {
    title: "更新时间",
    dataIndex: "updateTime",
    key: "updateTime",
    align: "center",
    width: 150
  },
  {
    title: "操作",
    key: "operation",
    align: "center",
    width: 330,
    fixed: "right"
  }
];

const rowSelection = computed(() => ({
  selectedRowKeys: ids.value,
  onChange: (selectedRowKeys, selectedRows) => {
    ids.value = selectedRowKeys;
    tableNames.value = selectedRows.map(item => item.tableName);
    single.value = selectedRows.length !== 1;
    multiple.value = !selectedRows.length;
  }
}));

/** 查询表集合 */
function getList() {
  loading.value = true;

  listTable(proxy.addDateRange(queryParams.value, dateRange.value))
    .then(response => {
      tableList.value = response.data.rows || [];
      total.value = response.data.total || 0;
    })
    .finally(() => {
      loading.value = false;
    });
}

/** 搜索按钮操作 */
function handleQuery() {
  queryParams.value.pageNum = 1;
  getList();
}

/** 生成代码操作 */
function handleGenTable(row = {}) {
  const tbNames = row.tableName || tableNames.value.join(",");

  if (!tbNames) {
    message.error("请选择要生成的数据");
    return;
  }

  if (row.genType === "1") {
    genCode(row.tableName).then(() => {
      message.success("成功生成到自定义路径：" + row.genPath);
    });
  } else {
    proxy.$download.zip("/tool/gen/batchGenCode?tables=" + tbNames, "ruoyi");
  }
}

/** 同步数据库操作 */
function handleSynchDb(row) {
  const tableName = row.tableName;

  Modal.confirm({
    title: "确认同步",
    content: `确认要强制同步"${tableName}"表结构吗？`,
    okText: "确定",
    cancelText: "取消",
    onOk() {
      return synchDb(tableName).then(() => {
        message.success("同步成功");
      });
    }
  });
}

/** 打开导入表弹窗 */
function openImportTable() {
  importRef.value?.show();
}

/** 重置按钮操作 */
function resetQuery() {
  dateRange.value = [];
  queryRef.value?.resetFields();
  handleQuery();
}

/** 预览按钮 */
function handlePreview(row) {
  previewTable(row.tableId).then(response => {
    preview.value.data = response.data || {};
    preview.value.open = true;
    preview.value.activeName = "domain.java";
  });
}

/** 修改按钮操作 */
function handleEditTable(row = {}) {
  const tableId = row.tableId || ids.value[0];

  router.push({
    path: "/tool/gen-edit/index",
    query: {
      tableId,
      pageNum: queryParams.value.pageNum
    }
  });
}

/** 删除按钮操作 */
function handleDelete(row = {}) {
  const tableIds = row.tableId || ids.value;

  Modal.confirm({
    title: "确认删除",
    content: `是否确认删除表编号为"${tableIds}"的数据项？`,
    okText: "确定",
    cancelText: "取消",
    onOk() {
      return delTable(tableIds).then(() => {
        getList();
        message.success("删除成功");
      });
    }
  });
}

function getPreviewName(key) {
  return key.substring(key.lastIndexOf("/") + 1, key.indexOf(".vm"));
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

.preview-modal :deep(.ant-modal-body) {
  max-height: 75vh;
  overflow: auto;
}

pre {
  margin: 0;
  padding: 12px;
  background: #f6f8fa;
  border-radius: 4px;
  overflow: auto;
}
</style>