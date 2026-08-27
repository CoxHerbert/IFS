<template>
  <a-card title="出货汇总" :bordered="false" class="panel-card">
    <a-row :gutter="[16, 16]">
      <a-col :xs="24" :sm="8">
        <a-statistic title="总箱数" :value="summary.totalQty" />
      </a-col>
      <a-col :xs="24" :sm="8">
        <a-statistic title="总体积" :value="`${summary.totalCBM} CBM`" />
      </a-col>
      <a-col :xs="24" :sm="8">
        <a-statistic title="建议柜型" :value="summary.containerSuggestion || '-'" />
      </a-col>
    </a-row>
  </a-card>
</template>

<script setup lang="ts">
import type { ShipmentSummary } from '@/types/shipment'

defineProps<{
  summary: ShipmentSummary
}>()
</script>

<style scoped>
.panel-card {
  border-radius: 8px;
}
</style>
