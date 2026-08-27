<template>
  <!-- 导入表 -->
  <a-modal v-model:open="visible" title="导入表" width="800px" :body-style="{ paddingTop: '16px' }" @ok="handleImportTable"
    @cancel="visible = false">
    <a-form ref="queryRef" :model="queryParams" layout="inline">
      <a-form-item label="表名称" name="tableName">
        <a-input v-model:value="queryParams.tableName" placeholder="请输入表名称" allow-clear @pressEnter="handleQuery" />
      </a-form-item>

      <a-form-item label="表描述" name="tableComment">
        <a-input v-model:value="queryParams.tableComment" placeholder="请输入表描述" allow-clear @pressEnter="handleQuery" />
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

    <a-table ref="tableRef" class="import-table" :data-source="dbTableList" :columns="columns"
      :row-key="record => record.tableName" :row-selection="rowSelection" :pagination="false" :scroll="{ y: 260 }"
      size="middle" @row="handleRow">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'createTime'">
          {{ parseTime(record.createTime, '{y}-{m}-{d}') }}
        </template>

        <template v-if="column.key === 'updateTime'">
          {{ parseTime(record.updateTime, '{y}-{m}-{d}') }}
        </template>
      </template>
    </a-table>

    <div v-if="total > 0" class="pagination-wrapper">
      <a-pagination v-model:current="queryParams.pageNum" v-model:page-size="queryParams.pageSize" :total="total"
        :show-size-changer="true" :show-total="total => `共 ${total} 条`" @change="getList" @showSizeChange="getList" />
    </div>

    <template #footer>
      <a-button @click="visible = false">取 消</a-button>
      <a-button type="primary" @click="handleImportTable">确 定</a-button>
    </template>
  </a-modal>
</template>

<script setup>
import { ref, reactive, getCurrentInstance, computed } from "vue";
import { message } from "ant-design-vue";
import { SearchOutlined, ReloadOutlined } from "@ant-design/icons-vue";
import { listDbTable, importTable } from "@/api/tool/gen";

const total = ref(0);
const visible = ref(false);
const tables = ref([]);
const dbTableList = ref([]);
const queryRef = ref();

const { proxy } = getCurrentInstance();

const queryParams = reactive({
  pageNum: 1,
  pageSize: 10,
  tableName: undefined,
  tableComment: undefined
});

const emit = defineEmits(["ok"]);

const columns = [
  {
    title: "表名称",
    dataIndex: "tableName",
    key: "tableName",
    ellipsis: true
  },
  {
    title: "表描述",
    dataIndex: "tableComment",
    key: "tableComment",
    ellipsis: true
  },
  {
    title: "创建时间",
    dataIndex: "createTime",
    key: "createTime"
  },
  {
    title: "更新时间",
    dataIndex: "updateTime",
    key: "updateTime"
  }
];

const rowSelection = computed(() => ({
  selectedRowKeys: tables.value,
  onChange: selectedRowKeys => {
    tables.value = selectedRowKeys;
  }
}));

/** 查询参数列表 */
function show() {
  getList();
  visible.value = true;
}

/** 单击选择行 */
function handleRow(record) {
  return {
    onClick: () => {
      const tableName = record.tableName;
      const index = tables.value.indexOf(tableName);

      if (index > -1) {
        tables.value.splice(index, 1);
      } else {
        tables.value.push(tableName);
      }
    }
  };
}

/** 查询表数据 */
function getList() {
  listDbTable(queryParams).then(res => {
    dbTableList.value = res.data.rows || [];
    total.value = res.data.total || 0;
  });
}

/** 搜索按钮操作 */
function handleQuery() {
  queryParams.pageNum = 1;
  getList();
}

/** 重置按钮操作 */
function resetQuery() {
  queryRef.value?.resetFields();

  queryParams.pageNum = 1;
  queryParams.pageSize = 10;

  handleQuery();
}

/** 导入按钮操作 */
function handleImportTable() {
  const tableNames = tables.value.join(",");

  if (!tableNames) {
    message.error("请选择要导入的表");
    return;
  }

  importTable({ tables: tableNames }).then(res => {
    message.success(res.msg || "导入成功");

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

<style scoped>
.import-table {
  margin-top: 16px;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>