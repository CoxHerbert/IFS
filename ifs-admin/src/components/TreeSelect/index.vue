<template>
  <div class="tree-select">
    <a-tree-select
      v-model:value="valueId"
      style="width: 100%"
      :tree-data="options"
      :field-names="{ value: objMap.value, label: objMap.label, children: objMap.children }"
      :tree-default-expanded-keys="defaultExpandedKey"
      :tree-node-filter-prop="objMap.label"
      :placeholder="placeholder"
      :allow-clear="true"
      :show-search="true"
      @clear="clearHandle"
    />
  </div>
</template>

<script setup>
const props = defineProps({
  objMap: {
    type: Object,
    default: () => ({
      value: "id",
      label: "label",
      children: "children"
    })
  },
  accordion: {
    type: Boolean,
    default: false
  },
  value: {
    type: [String, Number],
    default: ""
  },
  options: {
    type: Array,
    default: () => []
  },
  placeholder: {
    type: String,
    default: ""
  }
});

const emit = defineEmits(["update:value"]);

const valueId = computed({
  get: () => props.value,
  set: (val) => {
    emit("update:value", val);
  }
});

const defaultExpandedKey = ref([]);

function initHandle() {
  const selectedValue = valueId.value;
  if (selectedValue !== null && typeof selectedValue !== "undefined" && selectedValue !== "") {
    defaultExpandedKey.value = [selectedValue];
  }
}

function clearHandle() {
  valueId.value = "";
  defaultExpandedKey.value = [];
}

onMounted(() => {
  initHandle();
});

watch(valueId, () => {
  initHandle();
});
</script>
