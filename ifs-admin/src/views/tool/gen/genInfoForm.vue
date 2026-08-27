<template>
  <a-form ref="genInfoForm" :model="info" :rules="rules" :label-col="{ style: { width: '150px' } }">
    <a-row :gutter="16">
      <a-col :span="12">
        <a-form-item name="tplCategory">
          <template #label>生成模板</template>

          <a-select v-model:value="info.tplCategory" @change="tplSelectChange" style="width: 100%">
            <a-select-option value="crud">
              单表（增删改查）
            </a-select-option>
            <a-select-option value="tree">
              树表（增删改查）
            </a-select-option>
            <a-select-option value="sub">
              主子表（增删改查）
            </a-select-option>
          </a-select>
        </a-form-item>
      </a-col>

      <a-col :span="12">
        <a-form-item name="packageName">
          <template #label>
            生成包路径
            <a-tooltip content="生成在哪个java包下，例如 com.ruoyi.system">
              <QuestionCircleOutlined />
            </a-tooltip>
          </template>

          <a-input v-model:value="info.packageName" />
        </a-form-item>
      </a-col>

      <a-col :span="12">
        <a-form-item name="moduleName">
          <template #label>
            生成模块名
            <a-tooltip content="可理解为子系统名，例如 system">
              <QuestionCircleOutlined />
            </a-tooltip>
          </template>

          <a-input v-model:value="info.moduleName" />
        </a-form-item>
      </a-col>

      <a-col :span="12">
        <a-form-item name="businessName">
          <template #label>
            生成业务名
            <a-tooltip content="可理解为功能英文名，例如 user">
              <QuestionCircleOutlined />
            </a-tooltip>
          </template>

          <a-input v-model:value="info.businessName" />
        </a-form-item>
      </a-col>

      <a-col :span="12">
        <a-form-item name="functionName">
          <template #label>
            生成功能名
            <a-tooltip content="用作类描述，例如 用户">
              <QuestionCircleOutlined />
            </a-tooltip>
          </template>

          <a-input v-model:value="info.functionName" />
        </a-form-item>
      </a-col>

      <a-col :span="12">
        <a-form-item name="parentMenuId">
          <template #label>
            上级菜单
            <a-tooltip content="分配到指定菜单下，例如 系统管理">
              <QuestionCircleOutlined />
            </a-tooltip>
          </template>

          <a-tree-select v-model:value="info.parentMenuId" :tree-data="menuOptions" :field-names="{
            value: 'menuId',
            label: 'menuName',
            children: 'children'
          }" placeholder="请选择系统菜单" allow-clear tree-default-expand-all style="width: 100%" />
        </a-form-item>
      </a-col>

      <a-col :span="12">
        <a-form-item name="genType">
          <template #label>
            生成代码方式
            <a-tooltip content="默认为zip压缩包下载，也可以自定义生成路径">
              <QuestionCircleOutlined />
            </a-tooltip>
          </template>

          <a-radio-group v-model:value="info.genType">
            <a-radio value="0">zip压缩包</a-radio>
            <a-radio value="1">自定义路径</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-col>

      <a-col :span="24" v-if="info.genType == '1'">
        <a-form-item name="genPath">
          <template #label>
            自定义路径
            <a-tooltip content="填写磁盘绝对路径，若不填写，则生成到当前Web项目下">
              <QuestionCircleOutlined />
            </a-tooltip>
          </template>

          <a-input v-model:value="info.genPath">
            <template #addonAfter>
              <a-dropdown>
                <a-button type="primary">
                  最近路径快速选择
                  <DownOutlined />
                </a-button>

                <template #overlay>
                  <a-menu>
                    <a-menu-item @click="info.genPath = '/'">
                      恢复默认的生成基础路径
                    </a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>
            </template>
          </a-input>
        </a-form-item>
      </a-col>
    </a-row>

    <template v-if="info.tplCategory == 'tree'">
      <h4 class="form-header">其他信息</h4>

      <a-row :gutter="16">
        <a-col :span="12">
          <a-form-item name="treeCode">
            <template #label>
              树编码字段
              <a-tooltip content="树显示的编码字段名， 如：dept_id">
                <QuestionCircleOutlined />
              </a-tooltip>
            </template>

            <a-select v-model:value="info.treeCode" placeholder="请选择" style="width: 100%">
              <a-select-option v-for="(column, index) in info.columns" :key="index" :value="column.columnName">
                {{ column.columnName + "：" + column.columnComment }}
              </a-select-option>
            </a-select>
          </a-form-item>
        </a-col>

        <a-col :span="12">
          <a-form-item name="treeParentCode">
            <template #label>
              树父编码字段
              <a-tooltip content="树显示的父编码字段名， 如：parent_Id">
                <QuestionCircleOutlined />
              </a-tooltip>
            </template>

            <a-select v-model:value="info.treeParentCode" placeholder="请选择" style="width: 100%">
              <a-select-option v-for="(column, index) in info.columns" :key="index" :value="column.columnName">
                {{ column.columnName + "：" + column.columnComment }}
              </a-select-option>
            </a-select>
          </a-form-item>
        </a-col>

        <a-col :span="12">
          <a-form-item name="treeName">
            <template #label>
              树名称字段
              <a-tooltip content="树节点的显示名称字段名， 如：dept_name">
                <QuestionCircleOutlined />
              </a-tooltip>
            </template>

            <a-select v-model:value="info.treeName" placeholder="请选择" style="width: 100%">
              <a-select-option v-for="(column, index) in info.columns" :key="index" :value="column.columnName">
                {{ column.columnName + "：" + column.columnComment }}
              </a-select-option>
            </a-select>
          </a-form-item>
        </a-col>
      </a-row>
    </template>

    <template v-if="info.tplCategory == 'sub'">
      <h4 class="form-header">关联信息</h4>

      <a-row :gutter="16">
        <a-col :span="12">
          <a-form-item name="subTableName">
            <template #label>
              关联子表的表名
              <a-tooltip content="关联子表的表名， 如：sys_user">
                <QuestionCircleOutlined />
              </a-tooltip>
            </template>

            <a-select v-model:value="info.subTableName" placeholder="请选择" @change="subSelectChange" style="width: 100%">
              <a-select-option v-for="(table, index) in tables" :key="index" :value="table.tableName">
                {{ table.tableName + "：" + table.tableComment }}
              </a-select-option>
            </a-select>
          </a-form-item>
        </a-col>

        <a-col :span="12">
          <a-form-item name="subTableFkName">
            <template #label>
              子表关联的外键名
              <a-tooltip content="子表关联的外键名， 如：user_id">
                <QuestionCircleOutlined />
              </a-tooltip>
            </template>

            <a-select v-model:value="info.subTableFkName" placeholder="请选择" style="width: 100%">
              <a-select-option v-for="(column, index) in subColumns" :key="index" :value="column.columnName">
                {{ column.columnName + "：" + column.columnComment }}
              </a-select-option>
            </a-select>
          </a-form-item>
        </a-col>
      </a-row>
    </template>
  </a-form>
</template>

<script setup>
import {
  QuestionCircleOutlined,
  DownOutlined
} from "@ant-design/icons-vue";

import { listMenu } from "@/api/system/menu";

const subColumns = ref([]);
const menuOptions = ref([]);
const { proxy } = getCurrentInstance();

const props = defineProps({
  info: {
    type: Object,
    default: () => ({})
  },
  tables: {
    type: Array,
    default: () => []
  }
});

// 表单校验
const rules = ref({
  tplCategory: [
    {
      required: true,
      message: "请选择生成模板",
      trigger: "change"
    }
  ],
  packageName: [
    {
      required: true,
      message: "请输入生成包路径",
      trigger: "blur"
    }
  ],
  moduleName: [
    {
      required: true,
      message: "请输入生成模块名",
      trigger: "blur"
    }
  ],
  businessName: [
    {
      required: true,
      message: "请输入生成业务名",
      trigger: "blur"
    }
  ],
  functionName: [
    {
      required: true,
      message: "请输入生成功能名",
      trigger: "blur"
    }
  ]
});

function subSelectChange() {
  props.info.subTableFkName = "";
}

function tplSelectChange(value) {
  if (value !== "sub") {
    props.info.subTableName = "";
    props.info.subTableFkName = "";
  }
}

function setSubTableColumns(value) {
  for (const item of props.tables || []) {
    if (value === item.tableName) {
      subColumns.value = item.columns || [];
      break;
    }
  }
}

/** 查询菜单下拉树结构 */
function getMenuTreeselect() {
  listMenu().then(response => {
    menuOptions.value = proxy.handleTree(response.data, "menuId");
  });
}

watch(
  () => props.info.subTableName,
  val => {
    setSubTableColumns(val);
  }
);

getMenuTreeselect();
</script>