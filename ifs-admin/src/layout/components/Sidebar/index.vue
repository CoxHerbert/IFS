<template>
  <div class="sidebar-panel" :class="{ 'has-logo': showLogo }" :style="{
    backgroundColor:
      sideTheme === 'theme-dark'
        ? variables.menuBackground
        : variables.menuLightBackground
  }">
    <logo v-if="showLogo" :collapse="isCollapse" />

    <div :class="['sidebar-scrollbar', sideTheme]">
      <a-menu v-model:openKeys="openKeys" :selectedKeys="[activeMenu]" :inline-collapsed="isCollapse"
        :theme="sideTheme === 'theme-dark' ? 'dark' : 'light'" mode="inline" class="sidebar-menu"
        @click="handleMenuClick">
        <sidebar-item v-for="(route, index) in sidebarRouters" :key="route.path + index" :item="route"
          :base-path="route.path" />
      </a-menu>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from "vue";
import { useRoute } from "vue-router";
import { useStore } from "vuex";
import Logo from "./Logo";
import SidebarItem from "./SidebarItem";
import variables from "@/assets/styles/variables.module.scss";

const route = useRoute();
const store = useStore();

const openKeys = ref([]);

const sidebarRouters = computed(() => store.getters.sidebarRouters);
const showLogo = computed(() => store.state.settings.sidebarLogo);
const sideTheme = computed(() => store.state.settings.sideTheme);
const theme = computed(() => store.state.settings.theme);
const isCollapse = computed(() => !store.state.app.sidebar.opened);

const activeMenu = computed(() => {
  const { meta, path } = route;

  if (meta.activeMenu) {
    return meta.activeMenu;
  }

  return path;
});

function handleMenuClick() {
  // 路由跳转逻辑通常放在 SidebarItem 内部处理
}
</script>

<style scoped>
.sidebar-panel {
  overflow: hidden;
  transition: background-color 0.2s ease;
}

.sidebar-scrollbar {
  overflow-y: auto;
  overflow-x: hidden;
  padding: 10px 9px 18px;
}

.sidebar-menu {
  height: 100%;
  border-right: 0;
  background: transparent;
  font-size: 13px;
}

.sidebar-menu :deep(.ant-menu-item),
.sidebar-menu :deep(.ant-menu-submenu-title) {
  height: 43px;
  line-height: 43px;
  margin: 4px 0;
  padding-inline: 14px !important;
  border-radius: 12px;
  transition: color 0.18s ease, background-color 0.18s ease, transform 0.18s ease;
}

.sidebar-menu :deep(.ant-menu-item:hover),
.sidebar-menu :deep(.ant-menu-submenu-title:hover) {
  transform: translateX(2px);
}

.sidebar-menu :deep(.ant-menu-item .ant-menu-item-icon),
.sidebar-menu :deep(.ant-menu-submenu-title .ant-menu-item-icon) {
  font-size: 17px;
}

.sidebar-menu :deep(.ant-menu-inline .ant-menu-sub) {
  background: transparent;
}

.sidebar-menu :deep(.ant-menu-sub .ant-menu-item) {
  height: 39px;
  line-height: 39px;
  margin-block: 3px;
  padding-left: 42px !important;
  border-radius: 10px;
}

/* 替代 Element 的 scrollbar-wrapper */
.sidebar-scrollbar::-webkit-scrollbar {
  width: 6px;
}

.sidebar-scrollbar::-webkit-scrollbar-thumb {
  background: rgba(144, 147, 153, 0.3);
  border-radius: 3px;
}

.sidebar-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}

/* 让 AntD 菜单主题色跟随系统 theme */
.sidebar-menu :deep(.ant-menu-item-selected) {
  color: v-bind(theme);
  background: color-mix(in srgb, v-bind(theme) 12%, transparent);
  font-weight: 650;
}

.sidebar-menu :deep(.ant-menu-item-selected::after) {
  border-right-color: v-bind(theme);
}

.sidebar-menu :deep(.ant-menu-submenu-selected > .ant-menu-submenu-title) {
  color: v-bind(theme);
  font-weight: 650;
}

.sidebar-scrollbar::-webkit-scrollbar {
  width: 5px;
}
</style>
