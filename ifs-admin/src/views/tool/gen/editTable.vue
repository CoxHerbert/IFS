<template>
  <a-card>
    <a-tabs v-model:activeKey="activeName">
      <a-tab-pane key="basic" tab="基本信息">
        <basic-info-form ref="basicInfo" :info="info" />
      </a-tab-pane>

      <a-tab-pane key="columnInfo" tab="字段信息">
        <a-table ref="dragTable" :data-source="columns" :columns="tableColumns" :row-key="record => record.columnId"
          :pagination="false" :scroll="{ y: tableHeight }" size="middle" bordered>
          <template #bodyCell="{ column, record, index }">
            <template v-if="column.key === 'index'">
              {{ index + 1 }}
            </template>

            <template v-if="column.key === 'columnComment'">
              <a-input v-model:value="record.columnComment" />
            </template>

            <template v-if="column.key === 'javaType'">
              <a-select v-model:value="record.javaType" style="width: 100%">
                <a-select-option value="Long">Long</a-select-option>
                <a-select-option value="String">String</a-select-option>
                <a-select-option value="Integer">Integer</a-select-option>
                <a-select-option value="Double">Double</a-select-option>
                <a-select-option value="BigDecimal">BigDecimal</a-select-option>
                <a-select-option value="Date">Date</a-select-option>
              </a-select>
            </template>

            <template v-if="column.key === 'javaField'">
              <a-input v-model:value="record.javaField" />
            </template>

            <template v-if="column.key === 'isInsert'">
              <a-checkbox :checked="record.isInsert === '1'"
                @change="e => (record.isInsert = e.target.checked ? '1' : '0')" />
            </template>

            <template v-if="column.key === 'isEdit'">
              <a-checkbox :checked="record.isEdit === '1'"
                @change="e => (record.isEdit = e.target.checked ? '1' : '0')" />
            </template>

            <template v-if="column.key === 'isList'">
              <a-checkbox :checked="record.isList === '1'"
                @change="e => (record.isList = e.target.checked ? '1' : '0')" />
            </template>

            <template v-if="column.key === 'isQuery'">
              <a-checkbox :checked="record.isQuery === '1'"
                @change="e => (record.isQuery = e.target.checked ? '1' : '0')" />
            </template>

            <template v-if="column.key === 'queryType'">
              <a-select v-model:value="record.queryType" style="width: 100%">
                <a-select-option value="EQ">=</a-select-option>
                <a-select-option value="NE">!=</a-select-option>
                <a-select-option value="GT">&gt;</a-select-option>
                <a-select-option value="GTE">&gt;=</a-select-option>
                <a-select-option value="LT">&lt;</a-select-option>
                <a-select-option value="LTE">&lt;=</a-select-option>
                <a-select-option value="LIKE">LIKE</a-select-option>
                <a-select-option value="BETWEEN">BETWEEN</a-select-option>
              </a-select>
            </template>

            <template v-if="column.key === 'isRequired'">
              <a-checkbox :checked="record.isRequired === '1'"
                @change="e => (record.isRequired = e.target.checked ? '1' : '0')" />
            </template>

            <template v-if="column.key === 'htmlType'">
              <a-select v-model:value="record.htmlType" style="width: 100%">
                <a-select-option value="input">文本框</a-select-option>
                <a-select-option value="textarea">文本域</a-select-option>
                <a-select-option value="select">下拉框</a-select-option>
                <a-select-option value="radio">单选框</a-select-option>
                <a-select-option value="checkbox">复选框</a-select-option>
                <a-select-option value="datetime">日期控件</a-select-option>
                <a-select-option value="imageUpload">图片上传</a-select-option>
                <a-select-option value="fileUpload">文件上传</a-select-option>
                <a-select-option value="editor">富文本控件</a-select-option>
              </a-select>
            </template>

            <template v-if="column.key === 'dictType'">
              <a-select v-model:value="record.dictType" allow-clear show-search placeholder="请选择" style="width: 100%"
                :filter-option="filterDictOption">
                <a-select-option v-for="dict in dictOptions" :key="dict.dictType" :value="dict.dictType"
                  :label="dict.dictName">
                  <div class="dict-option">
                    <span>{{ dict.dictName }}</span>
                    <span class="dict-type">{{ dict.dictType }}</span>
                  </div>
                </a-select-option>
              </a-select>
            </template>
          </template>
        </a-table>
      </a-tab-pane>

      <a-tab-pane key="genInfo" tab="生成信息">
        <gen-info-form ref="genInfo" :info="info" :tables="tables" />
      </a-tab-pane>
    </a-tabs>

    <a-form>
      <a-form-item class="footer-actions">
        <a-button type="primary" @click="submitForm">提交</a-button>
        <a-button @click="close">返回</a-button>
      </a-form-item>
    </a-form>
  </a-card>
</template>

<script setup name="GenEdit">
import { ref, getCurrentInstance } from "vue";
import { useRoute } from "vue-router";
import { message } from "ant-design-vue";
import { getGenTable, updateGenTable } from "@/api/tool/gen";
import { optionselect as getDictOptionselect } from "@/api/system/dict/type";
import basicInfoForm from "./basicInfoForm";
import genInfoForm from "./genInfoForm";

const route = useRoute();
const { proxy } = getCurrentInstance();

const activeName = ref("columnInfo");
const tableHeight = ref(document.documentElement.scrollHeight - 245);
const tables = ref([]);
const columns = ref([]);
const dictOptions = ref([]);
const info = ref({});

const tableColumns = [
  {
    title: "序号",
    key: "index",
    width: 70,
    align: "center"
  },
  {
    title: "字段列名",
    dataIndex: "columnName",
    key: "columnName",
    width: 140,
    ellipsis: true
  },
  {
    title: "字段描述",
    key: "columnComment",
    width: 160
  },
  {
    title: "物理类型",
    dataIndex: "columnType",
    key: "columnType",
    width: 140,
    ellipsis: true
  },
  {
    title: "Java类型",
    key: "javaType",
    width: 150
  },
  {
    title: "java属性",
    key: "javaField",
    width: 150
  },
  {
    title: "插入",
    key: "isInsert",
    width: 80,
    align: "center"
  },
  {
    title: "编辑",
    key: "isEdit",
    width: 80,
    align: "center"
  },
  {
    title: "列表",
    key: "isList",
    width: 80,
    align: "center"
  },
  {
    title: "查询",
    key: "isQuery",
    width: 80,
    align: "center"
  },
  {
    title: "查询方式",
    key: "queryType",
    width: 140
  },
  {
    title: "必填",
    key: "isRequired",
    width: 80,
    align: "center"
  },
  {
    title: "显示类型",
    key: "htmlType",
    width: 160
  },
  {
    title: "字典类型",
    key: "dictType",
    width: 180
  }
];

function filterDictOption(input, option) {
  const label = option.label || "";
  const value = option.value || "";
  return (
    label.toLowerCase().includes(input.toLowerCase()) ||
    value.toLowerCase().includes(input.toLowerCase())
  );
}

/** 提交按钮 */
function submitForm() {
  const basicForm = proxy.$refs.basicInfo?.$refs.basicInfoForm;
  const genForm = proxy.$refs.genInfo?.$refs.genInfoForm;

  Promise.all([basicForm, genForm].map(getFormPromise)).then(res => {
    const validateResult = res.every(item => !!item);

    if (validateResult) {
      const genTable = Object.assign({}, info.value);

      genTable.columns = columns.value;
      genTable.params = {
        treeCode: genTable.treeCode,
        treeName: genTable.treeName,
        treeParentCode: genTable.treeParentCode,
        parentMenuId: genTable.parentMenuId
      };

      updateGenTable(genTable).then(res => {
        message.success(res.msg || "提交成功");

        if (res.code === 200) {
          close();
        }
      });
    } else {
      message.error("表单校验未通过，请重新检查提交内容");
    }
  });
}

function getFormPromise(form) {
  if (!form) {
    return Promise.resolve(false);
  }

  return form
    .validate()
    .then(() => true)
    .catch(() => false);
}

function close() {
  const obj = {
    path: "/tool/gen",
    query: {
      t: Date.now(),
      pageNum: route.query.pageNum
    }
  };

  proxy.$tab.closeOpenPage(obj);
}

function init() {
  const tableId = route.query && route.query.tableId;

  if (tableId) {
    // 获取表详细信息
    getGenTable(tableId).then(res => {
      columns.value = res.data.rows || [];
      info.value = res.data.info || {};
      tables.value = res.data.tables || [];
    });

    // 查询字典下拉列表
    getDictOptionselect().then(response => {
      dictOptions.value = response.data || [];
    });
  }
}

init();
</script>

<style scoped>
.footer-actions {
  text-align: center;
  margin-top: 10px;
}

.footer-actions :deep(.ant-form-item-control-input-content) {
  display: flex;
  justify-content: center;
  gap: 8px;
}

.dict-option {
  display: flex;
  justify-content: space-between;
}

.dict-type {
  color: #8492a6;
  font-size: 13px;
}
</style>