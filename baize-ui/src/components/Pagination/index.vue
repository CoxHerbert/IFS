<template>
  <div :class="{ hidden }" class="pagination-container">
    <a-pagination
      v-model:current="currentPage"
      v-model:pageSize="pageSize"
      :total="total"
      :page-size-options="pageSizeOptions"
      :show-size-changer="showSizeChanger"
      :show-quick-jumper="showQuickJumper"
      :show-total="showTotalText"
      @change="handleCurrentChange"
      @showSizeChange="handleSizeChange"
    />
  </div>
</template>

<script setup>
import { scrollTo } from "@/utils/scroll-to";

const props = defineProps({
  total: {
    required: true,
    type: Number
  },
  page: {
    type: Number,
    default: 1
  },
  limit: {
    type: Number,
    default: 20
  },
  pageSizes: {
    type: Array,
    default() {
      return [10, 20, 30, 50];
    }
  },
  pagerCount: {
    type: Number,
    default: document.body.clientWidth < 992 ? 5 : 7
  },
  layout: {
    type: String,
    default: "total, sizes, prev, pager, next, jumper"
  },
  background: {
    type: Boolean,
    default: true
  },
  autoScroll: {
    type: Boolean,
    default: true
  },
  hidden: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits();

const currentPage = computed({
  get() {
    return props.page;
  },
  set(val) {
    emit("update:page", val);
  }
});

const pageSize = computed({
  get() {
    return props.limit;
  },
  set(val) {
    emit("update:limit", val);
  }
});

const pageSizeOptions = computed(() => props.pageSizes.map(String));
const showSizeChanger = computed(() => props.layout.includes("sizes"));
const showQuickJumper = computed(() => props.layout.includes("jumper"));
const showTotalText = computed(() => {
  if (!props.layout.includes("total")) {
    return undefined;
  }
  return (total) => `Total ${total}`;
});

function handleSizeChange(_current, size) {
  emit("pagination", { page: currentPage.value, limit: size });
  if (props.autoScroll) {
    scrollTo(0, 800);
  }
}

function handleCurrentChange(page, size) {
  emit("pagination", { page, limit: size });
  if (props.autoScroll) {
    scrollTo(0, 800);
  }
}
</script>

<style scoped>
.pagination-container {
  background: #fff;
  padding: 32px 16px;
}

.pagination-container.hidden {
  display: none;
}
</style>
