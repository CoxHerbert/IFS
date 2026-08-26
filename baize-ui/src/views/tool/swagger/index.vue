<template>
  <div class="swagger-page">
    <i-frame :src="url" />
  </div>
</template>

<script setup>
import iFrame from '@/components/iFrame'

const url = ref(import.meta.env.VITE_APP_BASE_API + "/swagger/index.html")
</script>

<style scoped>
.swagger-page {
  width: 100%;
  min-width: 0;
  overflow: hidden;
  background: #fff;
}
</style>
