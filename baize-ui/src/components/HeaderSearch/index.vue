<template>
  <div :class="{ show }" class="header-search">
    <svg-icon class-name="search-icon" icon-class="search" @click.stop="click" />
    <a-select
      v-show="show"
      ref="headerSearchSelectRef"
      v-model:value="search"
      show-search
      :filter-option="false"
      :options="selectOptions"
      placeholder="Search"
      class="header-search-select"
      @search="querySearch"
      @select="change"
    />
  </div>
</template>

<script setup>
import Fuse from "fuse.js";
import { getNormalPath } from "@/utils/ruoyi";
import { isHttp } from "@/utils/validate";

const search = ref("");
const options = ref([]);
const searchPool = ref([]);
const show = ref(false);
const fuse = ref(undefined);
const headerSearchSelectRef = ref(null);
const store = useStore();
const router = useRouter();
const routes = computed(() => store.getters.permission_routes);

const selectOptions = computed(() => {
  return options.value.map(option => ({
    value: option.item.path,
    label: option.item.title.join(" > "),
    item: option.item
  }));
});

function click() {
  show.value = !show.value;
  if (show.value) {
    nextTick(() => {
      headerSearchSelectRef.value?.focus?.();
    });
  }
}

function close() {
  headerSearchSelectRef.value?.blur?.();
  options.value = [];
  show.value = false;
}

function change(path, option) {
  const targetPath = option.item.path || path;
  if (isHttp(targetPath)) {
    const pindex = targetPath.indexOf("http");
    window.open(targetPath.substr(pindex, targetPath.length), "_blank");
  } else {
    router.push(targetPath);
  }
  search.value = "";
  options.value = [];
  nextTick(() => {
    show.value = false;
  });
}

function initFuse(list) {
  fuse.value = new Fuse(list, {
    shouldSort: true,
    threshold: 0.4,
    location: 0,
    distance: 100,
    maxPatternLength: 32,
    minMatchCharLength: 1,
    keys: [
      { name: "title", weight: 0.7 },
      { name: "path", weight: 0.3 }
    ]
  });
}

function generateRoutes(routesValue, basePath = "", prefixTitle = []) {
  let res = [];
  for (const r of routesValue) {
    if (r.hidden) {
      continue;
    }
    const p = r.path.length > 0 && r.path[0] === "/" ? r.path : "/" + r.path;
    const data = {
      path: !isHttp(r.path) ? getNormalPath(basePath + p) : r.path,
      title: [...prefixTitle]
    };
    if (r.meta && r.meta.title) {
      data.title = [...data.title, r.meta.title];
      if (r.redirect !== "noRedirect") {
        res.push(data);
      }
    }
    if (r.children) {
      const tempRoutes = generateRoutes(r.children, data.path, data.title);
      if (tempRoutes.length >= 1) {
        res = [...res, ...tempRoutes];
      }
    }
  }
  return res;
}

function querySearch(query) {
  if (query !== "") {
    options.value = fuse.value?.search(query) || [];
  } else {
    options.value = [];
  }
}

onMounted(() => {
  searchPool.value = generateRoutes(routes.value);
});

watchEffect(() => {
  searchPool.value = generateRoutes(routes.value);
});

watch(show, (value) => {
  if (value) {
    document.body.addEventListener("click", close);
  } else {
    document.body.removeEventListener("click", close);
  }
});

watch(searchPool, (list) => {
  initFuse(list);
});
</script>

<style lang="scss" scoped>
.header-search {
  font-size: 0 !important;

  .search-icon {
    cursor: pointer;
    font-size: 18px;
    vertical-align: middle;
  }

  .header-search-select {
    font-size: 18px;
    width: 210px;
    margin-left: 10px;
    vertical-align: middle;
  }
}
</style>
