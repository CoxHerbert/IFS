<template>
  <div>
    <template v-for="(item, index) in options" :key="item.value">
      <template v-if="values.includes(item.value)">
        <span
          v-if="item.elTagType === 'default' || item.elTagType === ''"
          :index="index"
          :class="item.elTagType"
        >{{ item.label }}</span>
        <a-tag
          v-else
          :index="index"
          :color="tagColor(item.elTagType)"
          :class="item.elTagType"
        >{{ item.label }}</a-tag>
      </template>
    </template>
  </div>
</template>

<script setup>
const props = defineProps({
  options: {
    type: Array,
    default: null
  },
  value: [Number, String, Array]
});

const values = computed(() => {
  if (props.value !== null && typeof props.value !== "undefined") {
    return Array.isArray(props.value) ? props.value : [String(props.value)];
  }
  return [];
});

function tagColor(type) {
  const map = {
    primary: "processing",
    success: "success",
    info: "default",
    warning: "warning",
    danger: "error"
  };
  return map[type] || "default";
}
</script>

<style scoped>
.ant-tag + .ant-tag {
  margin-left: 10px;
}
</style>
