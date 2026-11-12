<template>
  <section class="app-main">
    <router-view v-slot="{ Component, route }">
      <transition name="fade-transform" mode="out-in">
        <keep-alive :include="cachedViews">
          <component :is="Component" :key="route.path"/>
        </keep-alive>
      </transition>
    </router-view>
  </section>
</template>

<script setup>
let store = useStore()
const cachedViews = computed(() => {
    return store.state.tagsView.cachedViews
})
</script>

<style lang="scss" scoped>
.app-main {
  /* Match the fixed navigation header. */
  min-height: calc(100vh - 58px);
  width: 100%;
  position: relative;
  overflow: hidden;
}

.fixed-header + .app-main {
  padding-top: 58px;
}

.hasTagsView {
  .app-main {
    /* 100 = 58px navbar + 42px tags view. */
    min-height: calc(100vh - 100px);
  }

  .fixed-header + .app-main {
    padding-top: 100px;
  }
}
</style>

<style lang="scss">
// fix css style bug in open el-dialog
.el-popup-parent--hidden {
  .fixed-header {
    padding-right: 17px;
  }
}
</style>
