<template>
  <a-layout class="portal-shell">
    <PortalHeader />

    <a-layout-content>
      <router-view />
    </a-layout-content>

    <PortalFooter />
  </a-layout>
</template>

<script setup lang="ts">
import PortalFooter from './components/PortalFooter.vue'
import PortalHeader from './components/PortalHeader.vue'
</script>
