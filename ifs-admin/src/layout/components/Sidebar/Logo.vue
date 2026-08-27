<template>
  <div class="sidebar-logo-container" :class="{ collapse }"
    :style="{ backgroundColor: sideTheme === 'theme-dark' ? variables.menuBackground : variables.menuLightBackground }">
    <transition name="sidebarLogoFade">
      <router-link v-if="collapse" key="collapse" class="sidebar-logo-link" to="/index">
        <img class="brand-logo" src="/logo.svg" alt="IFS logo" />
      </router-link>
      <router-link v-else key="expand" class="sidebar-logo-link" to="/index">
        <img class="brand-logo" src="/logo.svg" alt="IFS logo" />
        <span class="sidebar-title"
          :style="{ color: sideTheme === 'theme-dark' ? variables.logoTitleColor : variables.logoLightTitleColor }">
          {{ title }}
        </span>
      </router-link>
    </transition>
  </div>
</template>

<script setup>
import variables from '@/assets/styles/variables.module.scss'

defineProps({
  collapse: {
    type: Boolean,
    required: true,
  },
})

const title = ref('IFS 管理系统')
const store = useStore()
const sideTheme = computed(() => store.state.settings.sideTheme)
</script>

<style lang="scss" scoped>
.sidebarLogoFade-enter-active {
  transition: opacity 0.9s;
}

.sidebarLogoFade-enter,
.sidebarLogoFade-leave-to {
  opacity: 0;
}

.sidebar-logo-container {
  position: relative;
  width: 100%;
  height: 58px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 12px;
  overflow: hidden;
  border-bottom: 1px solid rgba(148, 163, 184, 0.12);

  .sidebar-logo-link {
    width: 100%;
    display: flex;
    align-items: center;
    gap: 12px;
    overflow: hidden;
  }

  .brand-logo {
    width: 36px;
    height: 36px;
    flex: 0 0 auto;
    object-fit: contain;
    transition: transform 0.2s ease, box-shadow 0.2s ease;
  }

  .sidebar-title {
    margin-left: 10px;
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    font-size: 15px;
    font-weight: 700;
    letter-spacing: 0.02em;
    line-height: 1;
  }

  &.collapse {
    padding: 0;

    .sidebar-logo-link {
      justify-content: center;
    }
  }

  &:hover .brand-logo {
    transform: translateY(-1px);
  }
}
</style>
