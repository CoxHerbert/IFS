<template>
  <a-drawer v-model:open="showSettings" placement="right" width="300px" :closable="false">
    <div class="setting-drawer-title">
      <h3 class="drawer-title">主题风格设置</h3>
    </div>

    <div class="setting-drawer-block-checbox">
      <div class="setting-drawer-block-checbox-item" @click="handleTheme('theme-dark')">
        <img src="@/assets/images/dark.svg" alt="dark" />

        <div v-if="sideTheme === 'theme-dark'" class="setting-drawer-block-checbox-selectIcon" style="display: block">
          <CheckOutlined :style="{ color: theme }" />
        </div>
      </div>

      <div class="setting-drawer-block-checbox-item" @click="handleTheme('theme-light')">
        <img src="@/assets/images/light.svg" alt="light" />

        <div
          v-if="sideTheme === 'theme-light'"
          class="setting-drawer-block-checbox-selectIcon"
          style="display: block"
        >
          <CheckOutlined :style="{ color: theme }" />
        </div>
      </div>
    </div>

    <div class="drawer-item">
      <span>主题颜色</span>
      <span class="comp-style theme-color-control">
        <input v-model="theme" class="theme-color-input" type="color" @input="themeChange(theme)" />
      </span>
    </div>

    <div class="theme-color-presets">
      <button
        v-for="color in predefineColors"
        :key="color"
        class="theme-color-preset"
        type="button"
        :aria-label="`select theme color ${color}`"
        :style="{ backgroundColor: color }"
        @click="themeChange(color)"
      />
    </div>

    <a-divider />

    <h3 class="drawer-title">系统布局配置</h3>

    <div class="drawer-item">
      <span>开启 TopNav</span>
      <span class="comp-style">
        <a-switch v-model:checked="topNav" class="drawer-switch" />
      </span>
    </div>

    <div class="drawer-item">
      <span>开启 Tags-Views</span>
      <span class="comp-style">
        <a-switch v-model:checked="tagsView" class="drawer-switch" />
      </span>
    </div>

    <div class="drawer-item">
      <span>固定 Header</span>
      <span class="comp-style">
        <a-switch v-model:checked="fixedHeader" class="drawer-switch" />
      </span>
    </div>

    <div class="drawer-item">
      <span>显示 Logo</span>
      <span class="comp-style">
        <a-switch v-model:checked="sidebarLogo" class="drawer-switch" />
      </span>
    </div>

    <div class="drawer-item">
      <span>动态标题</span>
      <span class="comp-style">
        <a-switch v-model:checked="dynamicTitle" class="drawer-switch" />
      </span>
    </div>

    <a-divider />

    <a-space>
      <a-button type="primary" size="small" ghost @click="saveSetting">
        <template #icon>
          <SaveOutlined />
        </template>
        保存配置
      </a-button>

      <a-button size="small" ghost @click="resetSetting">
        <template #icon>
          <ReloadOutlined />
        </template>
        重置配置
      </a-button>
    </a-space>
  </a-drawer>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";
import { useStore } from "vuex";
import { message } from "ant-design-vue";
import {
  CheckOutlined,
  SaveOutlined,
  ReloadOutlined
} from "@ant-design/icons-vue";
import { useDynamicTitle } from "@/utils/dynamicTitle";

const store = useStore();

const showSettings = ref(false);
const theme = ref(store.state.settings.theme);
const sideTheme = ref(store.state.settings.sideTheme);
const storeSettings = computed(() => store.state.settings);

const predefineColors = ref([
  "#409EFF",
  "#ff4500",
  "#ff8c00",
  "#ffd700",
  "#90ee90",
  "#00ced1",
  "#1e90ff",
  "#c71585"
]);

/** 是否需要 topnav */
const topNav = computed({
  get: () => storeSettings.value.topNav,
  set: val => {
    store.dispatch("settings/changeSetting", {
      key: "topNav",
      value: val
    });

    if (!val) {
      store.commit("SET_SIDEBAR_ROUTERS", store.state.permission.defaultRoutes);
    }
  }
});

/** 是否需要 tagview */
const tagsView = computed({
  get: () => storeSettings.value.tagsView,
  set: val => {
    store.dispatch("settings/changeSetting", {
      key: "tagsView",
      value: val
    });
  }
});

/** 是否需要固定头部 */
const fixedHeader = computed({
  get: () => storeSettings.value.fixedHeader,
  set: val => {
    store.dispatch("settings/changeSetting", {
      key: "fixedHeader",
      value: val
    });
  }
});

/** 是否需要侧边栏 logo */
const sidebarLogo = computed({
  get: () => storeSettings.value.sidebarLogo,
  set: val => {
    store.dispatch("settings/changeSetting", {
      key: "sidebarLogo",
      value: val
    });
  }
});

/** 是否需要动态网页 title */
const dynamicTitle = computed({
  get: () => storeSettings.value.dynamicTitle,
  set: val => {
    store.dispatch("settings/changeSetting", {
      key: "dynamicTitle",
      value: val
    });

    useDynamicTitle();
  }
});

function normalizeColorValue(val: any) {
  if (!val) return theme.value;

  if (typeof val === "string") {
    return val;
  }

  if (typeof val.toHexString === "function") {
    return val.toHexString();
  }

  return String(val);
}

function themeChange(val: any) {
  const color = normalizeColorValue(val);

  store.dispatch("settings/changeSetting", {
    key: "theme",
    value: color
  });

  theme.value = color;
}

function handleTheme(val: string) {
  store.dispatch("settings/changeSetting", {
    key: "sideTheme",
    value: val
  });

  sideTheme.value = val;
}

function saveSetting() {
  const hide = message.loading("正在保存到本地，请稍候...", 0);

  const layoutSetting = {
    topNav: storeSettings.value.topNav,
    tagsView: storeSettings.value.tagsView,
    fixedHeader: storeSettings.value.fixedHeader,
    sidebarLogo: storeSettings.value.sidebarLogo,
    dynamicTitle: storeSettings.value.dynamicTitle,
    sideTheme: storeSettings.value.sideTheme,
    theme: storeSettings.value.theme
  };

  localStorage.setItem("layout-setting", JSON.stringify(layoutSetting));

  window.setTimeout(() => {
    hide();
    message.success("保存成功");
  }, 1000);
}

function resetSetting() {
  const hide = message.loading("正在清除设置缓存并刷新，请稍候...", 0);

  localStorage.removeItem("layout-setting");

  window.setTimeout(() => {
    hide();
    window.location.reload();
  }, 1000);
}

function openSetting() {
  showSettings.value = true;
}

defineExpose({
  openSetting
});
</script>

<style lang="scss" scoped>
.setting-drawer-title {
  margin-bottom: 12px;
  color: rgba(0, 0, 0, 0.85);
  line-height: 22px;
  font-weight: bold;

  .drawer-title {
    font-size: 14px;
  }
}

.setting-drawer-block-checbox {
  display: flex;
  justify-content: flex-start;
  align-items: center;
  margin-top: 10px;
  margin-bottom: 20px;

  .setting-drawer-block-checbox-item {
    position: relative;
    margin-right: 16px;
    border-radius: 2px;
    cursor: pointer;

    img {
      width: 48px;
      height: 48px;
    }

    .custom-img {
      width: 48px;
      height: 38px;
      border-radius: 5px;
      box-shadow: 1px 1px 2px #898484;
    }

    .setting-drawer-block-checbox-selectIcon {
      position: absolute;
      top: 0;
      right: 0;
      width: 100%;
      height: 100%;
      padding-top: 15px;
      padding-left: 24px;
      color: #1890ff;
      font-weight: 700;
      font-size: 14px;
    }
  }
}

.drawer-item {
  color: rgba(0, 0, 0, 0.65);
  padding: 12px 0;
  font-size: 14px;

  .comp-style {
    float: right;
    margin: -3px 8px 0 0;
  }
}

.theme-color-control {
  display: inline-flex;
  align-items: center;
}

.theme-color-input {
  width: 36px;
  height: 24px;
  padding: 0;
  border: none;
  background: transparent;
  cursor: pointer;
}

.theme-color-presets {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding-bottom: 12px;
}

.theme-color-preset {
  width: 20px;
  height: 20px;
  padding: 0;
  border: 1px solid rgba(0, 0, 0, 0.15);
  border-radius: 50%;
  cursor: pointer;
}
</style>
