<template>
  <template v-if="!item.hidden">
    <template v-if="
      hasOneShowingChild(item.children, item) &&
      (!onlyOneChild.children || onlyOneChild.noShowingChildren) &&
      !item.alwaysShow
    ">
      <a-menu-item
        v-if="onlyOneChild.meta"
        v-bind="attrs"
        :key="resolveMenuKey(onlyOneChild.path)"
        :class="{ 'submenu-title-noDropdown': !isNest }"
      >
        <template #icon>
          <svg-icon :icon-class="onlyOneChild.meta.icon || (item.meta && item.meta.icon)" />
        </template>

        <app-link :to="resolvePath(onlyOneChild.path, onlyOneChild.query)">
          <span>{{ onlyOneChild.meta.title }}</span>
        </app-link>
      </a-menu-item>
    </template>

    <a-sub-menu v-else v-bind="attrs" :key="resolveMenuKey(item.path)">
      <template #icon>
        <svg-icon :icon-class="item.meta && item.meta.icon" />
      </template>

      <template #title>
        <span>{{ item.meta && item.meta.title }}</span>
      </template>

      <sidebar-item v-for="child in item.children" :key="child.path" :is-nest="true" :item="child"
        :base-path="resolvePath(child.path)" class="nest-menu" />
    </a-sub-menu>
  </template>
</template>

<script setup>
import { ref, useAttrs } from "vue";
import { isExternal } from "@/utils/validate";
import AppLink from "./Link";
import { getNormalPath } from "@/utils/ruoyi";

defineOptions({
  inheritAttrs: false
});

const props = defineProps({
  // route object
  item: {
    type: Object,
    required: true
  },
  isNest: {
    type: Boolean,
    default: false
  },
  basePath: {
    type: String,
    default: ""
  }
});

const attrs = useAttrs();
const onlyOneChild = ref({});

function hasOneShowingChild(children = [], parent) {
  if (!children) {
    children = [];
  }

  const showingChildren = children.filter(item => {
    if (item.hidden) {
      return false;
    }

    // Temp set(will be used if only has one showing child)
    onlyOneChild.value = item;
    return true;
  });

  // When there is only one child router, the child router is displayed by default
  if (showingChildren.length === 1) {
    return true;
  }

  // Show parent if there are no child router to display
  if (showingChildren.length === 0) {
    onlyOneChild.value = {
      ...parent,
      path: "",
      noShowingChildren: true
    };

    return true;
  }

  return false;
}

function resolvePath(routePath, routeQuery) {
  if (isExternal(routePath)) {
    return routePath;
  }

  if (isExternal(props.basePath)) {
    return props.basePath;
  }

  if (routeQuery) {
    const query = JSON.parse(routeQuery);

    return {
      path: getNormalPath(props.basePath + "/" + routePath),
      query
    };
  }

  return getNormalPath(props.basePath + "/" + routePath);
}

function resolveMenuKey(routePath) {
  const resolved = resolvePath(routePath);

  if (typeof resolved === "string") {
    return resolved;
  }

  return resolved.path;
}
</script>
