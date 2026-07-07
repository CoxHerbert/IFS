<template>
  <div
    class="sidebar-logo-container"
    :class="{ collapse }"
    :style="{ backgroundColor: sideTheme === 'theme-dark' ? variables.menuBackground : variables.menuLightBackground }"
  >
    <transition name="sidebarLogoFade">
      <router-link v-if="collapse" key="collapse" class="sidebar-logo-link" to="/">
        <span class="brand-mark" :class="sideTheme">BZ</span>
      </router-link>
      <router-link v-else key="expand" class="sidebar-logo-link" to="/">
        <span class="brand-mark" :class="sideTheme">BZ</span>
        <span
          class="sidebar-title"
          :style="{ color: sideTheme === 'theme-dark' ? variables.logoTitleColor : variables.logoLightTitleColor }"
        >
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

  .brand-mark {
    width: 36px;
    height: 36px;
    border-radius: 12px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    flex: 0 0 auto;
    font-size: 14px;
    font-weight: 800;
    letter-spacing: 0.08em;
    box-shadow: 0 10px 24px rgba(15, 23, 42, 0.18);
    transition: transform 0.2s ease, box-shadow 0.2s ease;

    &.theme-dark {
      background: linear-gradient(135deg, #f8fafc 0%, #7dd3fc 48%, #38bdf8 100%);
      color: #0f172a;
    }

    &.theme-light {
      background: linear-gradient(135deg, #0f172a 0%, #1d4ed8 100%);
      color: #ffffff;
      box-shadow: 0 10px 24px rgba(37, 99, 235, 0.18);
    }
  }

  .sidebar-title {
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

  &:hover .brand-mark {
    transform: translateY(-1px);
  }
}
</style>
