<template>
  <div class="top-right-btn">
    <a-space>
      <a-tooltip :title="showSearch ? '隐藏搜索' : '显示搜索'">
        <a-button shape="circle" @click="toggleSearch"><SearchOutlined /></a-button>
      </a-tooltip>
      <a-tooltip title="刷新">
        <a-button shape="circle" @click="refresh"><ReloadOutlined /></a-button>
      </a-tooltip>
      <a-tooltip v-if="columns" title="显隐列">
        <a-button shape="circle" @click="showColumn"><MenuOutlined /></a-button>
      </a-tooltip>
    </a-space>
    <a-modal v-model:open="open" :title="title" :footer="null" destroy-on-close>
      <a-checkbox-group v-model:value="value" class="column-group" @change="dataChange">
        <a-checkbox v-for="item in columns" :key="item.key" :value="item.key">
          {{ item.title || item.label || item.key }}
        </a-checkbox>
      </a-checkbox-group>
    </a-modal>
  </div>
</template>

<script setup>
import { MenuOutlined, ReloadOutlined, SearchOutlined } from "@ant-design/icons-vue";

const props = defineProps({
  showSearch: {
    type: Boolean,
    default: true
  },
  columns: {
    type: Array
  }
});

const emits = defineEmits(["update:showSearch", "queryTable"]);

const value = ref([]);
const title = ref("显示/隐藏");
const open = ref(false);

function toggleSearch() {
  emits("update:showSearch", !props.showSearch);
}

function refresh() {
  emits("queryTable");
}

function dataChange(data) {
  for (const item in props.columns) {
    const key = props.columns[item].key;
    props.columns[item].visible = data.includes(key);
  }
}

function showColumn() {
  open.value = true;
}

if (props.columns) {
  value.value = props.columns.filter(item => item.visible !== false).map(item => item.key);
}
</script>

<style lang="scss" scoped>
.column-group {
  display: grid;
  gap: 12px;
}
</style>
