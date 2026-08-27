<template>
  <a-menu class="topmenu-container" mode="horizontal" :selectedKeys="[activeMenu]" @select="handleMenuSelect">
    <a-menu-item
      v-for="(item, index) in topMenus"
      v-if="index < visibleNumber"
      :key="item.path"
      :style="{ '--theme': theme }"
    >
      <svg-icon :icon-class="item.meta.icon" />
      {{ item.meta.title }}
    </a-menu-item>

    <a-sub-menu v-if="topMenus.length > visibleNumber" key="more" :style="{ '--theme': theme }">
      <template #title>更多菜单</template>
      <a-menu-item v-for="(item, index) in topMenus" v-if="index >= visibleNumber" :key="item.path">
        <svg-icon :icon-class="item.meta.icon" />
        {{ item.meta.title }}
      </a-menu-item>
    </a-sub-menu>
  </a-menu>
</template>

<script setup>
import { constantRoutes } from "@/router";
import { isHttp } from "@/utils/validate";

const visibleNumber = ref(null);
const isFrist = ref(null);
const currentIndex = ref(null);

const store = useStore();
const route = useRoute();
const router = useRouter();

const theme = computed(() => store.state.settings.theme);
const routers = computed(() => store.state.permission.topbarRouters);

const topMenus = computed(() => {
  const menus = [];
  routers.value.map((menu) => {
    if (menu.hidden !== true) {
      if (menu.path === "/") {
        menus.push(menu.children[0]);
      } else {
        menus.push(menu);
      }
    }
  });
  return menus;
});

const childrenMenus = computed(() => {
  const menus = [];
  routers.value.map((routerItem) => {
    for (const item in routerItem.children) {
      if (routerItem.children[item].parentPath === undefined) {
        if (routerItem.path === "/") {
          routerItem.children[item].path = "/redirect/" + routerItem.children[item].path;
        } else if (!isHttp(routerItem.children[item].path)) {
          routerItem.children[item].path = routerItem.path + "/" + routerItem.children[item].path;
        }
        routerItem.children[item].parentPath = routerItem.path;
      }
      menus.push(routerItem.children[item]);
    }
  });
  return constantRoutes.concat(menus);
});

const activeMenu = computed(() => {
  const path = route.path;
  let activePath = defaultRouter.value;
  if (path !== undefined && path.lastIndexOf("/") > 0) {
    const tmpPath = path.substring(1, path.length);
    activePath = "/" + tmpPath.substring(0, tmpPath.indexOf("/"));
  } else if (path === "/index" || path === "") {
    if (!isFrist.value) {
      isFrist.value = true;
    } else {
      activePath = "index";
    }
  }
  const routesValue = activeRoutes(activePath);
  if (routesValue.length === 0) {
    activePath = currentIndex.value || defaultRouter.value;
    activeRoutes(activePath);
  }
  return activePath;
});

const defaultRouter = computed(() => {
  let routerValue;
  Object.keys(routers.value).some((key) => {
    if (!routers.value[key].hidden) {
      routerValue = routers.value[key].path;
      return true;
    }
    return false;
  });
  return routerValue;
});

function setVisibleNumber() {
  const width = document.body.getBoundingClientRect().width / 3;
  visibleNumber.value = parseInt(width / 85);
}

function handleMenuSelect({ key }) {
  handleSelect(key);
}

function handleSelect(key) {
  currentIndex.value = key;
  if (isHttp(key)) {
    window.open(key, "_blank");
  } else if (key.indexOf("/redirect") !== -1) {
    router.push({ path: key.replace("/redirect", "") });
  } else {
    activeRoutes(key);
  }
}

function activeRoutes(key) {
  const routesValue = [];
  if (childrenMenus.value && childrenMenus.value.length > 0) {
    childrenMenus.value.map((item) => {
      if (key === item.parentPath || (key === "index" && item.path === "")) {
        routesValue.push(item);
      }
    });
  }
  if (routesValue.length > 0) {
    store.commit("SET_SIDEBAR_ROUTERS", routesValue);
  }
  return routesValue;
}

onMounted(() => {
  window.addEventListener("resize", setVisibleNumber);
  setVisibleNumber();
});

onBeforeUnmount(() => {
  window.removeEventListener("resize", setVisibleNumber);
});
</script>

<style lang="scss">
.topmenu-container > .ant-menu-item {
  height: 50px !important;
  line-height: 50px !important;
  color: #999093 !important;
  padding: 0 5px !important;
  margin: 0 10px !important;
}

.topmenu-container > .ant-menu-item-selected,
.topmenu-container > .ant-menu-submenu-selected > .ant-menu-submenu-title {
  border-bottom: 2px solid #{'var(--theme)'} !important;
  color: #303133 !important;
}

.topmenu-container > .ant-menu-submenu > .ant-menu-submenu-title {
  height: 50px !important;
  line-height: 50px !important;
  color: #999093 !important;
  padding: 0 5px !important;
  margin: 0 10px !important;
}
</style>
