<template>
  <a-card title="货物明细预览" :bordered="false" class="panel-card">
    <vxe-table :data="cargoList" border stripe height="420" size="small">
      <vxe-column field="rowIndex" title="行号" width="80" />
      <vxe-column field="sku" title="SKU" min-width="120" />
      <vxe-column field="productName" title="品名" min-width="160" />
      <vxe-column field="qty" title="箱数" width="100" />
      <vxe-column field="length" title="长(cm)" width="100" />
      <vxe-column field="width" title="宽(cm)" width="100" />
      <vxe-column field="height" title="高(cm)" width="100" />
      <vxe-column field="weight" title="重量" width="100" />
      <vxe-column field="cbm" title="CBM" width="110" />
    </vxe-table>
  </a-card>
</template>

<script setup lang="ts">
import type { StandardCargoItem } from '@/types/shipment'

defineProps<{
  cargoList: StandardCargoItem[]
}>()
</script>

<style scoped>
.panel-card {
  border-radius: 8px;
}
</style>
