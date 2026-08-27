<template>
  <div>
    <a-dropdown>
      <div class="size-icon--style">
        <svg-icon class-name="size-icon" icon-class="size" />
      </div>
      <template #overlay>
        <a-menu @click="handleMenuClick">
          <a-menu-item v-for="item of sizeOptions" :key="item.value" :disabled="size === item.value">
            {{ item.label }}
          </a-menu-item>
        </a-menu>
      </template>
    </a-dropdown>
  </div>
</template>

<script setup lang="ts">
import { message } from "ant-design-vue";
import { computed, nextTick, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useStore } from "vuex";

const store = useStore();
const size = computed(() => store.getters.size);
const route = useRoute();
const router = useRouter();
const sizeOptions = ref([
  { label: "Large", value: "large" },
  { label: "Medium", value: "medium" },
  { label: "Small", value: "small" },
  { label: "Mini", value: "mini" }
]);

function refreshView() {
  store.dispatch("tagsView/delAllCachedViews", route);
  const { fullPath } = route;
  nextTick(() => {
    router.replace({
      path: "/redirect" + fullPath
    });
  });
}

function handleSetSize(sizeValue: string) {
  store.dispatch("app/setSize", sizeValue);
  refreshView();
  message.success("Switch size success");
}

function handleMenuClick({ key }) {
  handleSetSize(key);
}
</script>

<style lang="scss" scoped>
.size-icon--style {
  font-size: 18px;
  line-height: 50px;
  padding-right: 7px;
}
</style>
