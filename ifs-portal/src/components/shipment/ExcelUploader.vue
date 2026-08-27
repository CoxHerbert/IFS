<template>
  <a-upload-dragger
    accept=".xlsx,.xls,.csv"
    :before-upload="beforeUpload"
    :show-upload-list="false"
  >
    <p class="upload-title">上传出货 Excel</p>
    <p class="upload-desc">支持 .xlsx、.xls、.csv，文件将上传到服务端解析并计算 CBM。</p>
  </a-upload-dragger>
</template>

<script setup lang="ts">
const emit = defineEmits<{
  (event: 'select', file: File): void
}>()

function beforeUpload(file: File) {
  emit('select', file)
  return false
}
</script>

<style scoped>
.upload-title {
  margin: 0;
  font-size: 18px;
  font-weight: 700;
}

.upload-desc {
  margin: 8px 0 0;
  color: #64748b;
}
</style>
