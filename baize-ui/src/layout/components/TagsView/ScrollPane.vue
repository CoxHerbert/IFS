<template>
  <div ref="scrollContainer" class="scroll-container" @wheel.prevent="handleScroll">
    <div class="scroll-content">
      <slot />
    </div>
  </div>
</template>

<script setup>
const tagAndTagSpacing = ref(4);
const scrollContainer = ref();

const emits = defineEmits(["scroll"]);

const store = useStore();
const visitedViews = computed(() => store.state.tagsView.visitedViews);

onMounted(() => {
  scrollContainer.value?.addEventListener("scroll", emitScroll, true);
});

onBeforeUnmount(() => {
  scrollContainer.value?.removeEventListener("scroll", emitScroll, true);
});

function handleScroll(e) {
  const wrapper = scrollContainer.value;

  if (!wrapper) return;

  const delta = e.deltaX || e.deltaY || 0;
  wrapper.scrollLeft += delta;
}

function emitScroll() {
  emits("scroll");
}

function moveToTarget(currentTag) {
  const container = scrollContainer.value;

  if (!container) return;

  const containerWidth = container.offsetWidth;
  const scrollWrapper = container;

  let firstTag = null;
  let lastTag = null;

  if (visitedViews.value.length > 0) {
    firstTag = visitedViews.value[0];
    lastTag = visitedViews.value[visitedViews.value.length - 1];
  }

  if (firstTag === currentTag) {
    scrollWrapper.scrollLeft = 0;
    return;
  }

  if (lastTag === currentTag) {
    scrollWrapper.scrollLeft = scrollWrapper.scrollWidth - containerWidth;
    return;
  }

  const tagListDom = document.getElementsByClassName("tags-view-item");
  const currentIndex = visitedViews.value.findIndex(item => item === currentTag);

  let prevTag = null;
  let nextTag = null;

  for (const k in tagListDom) {
    if (k !== "length" && Object.hasOwnProperty.call(tagListDom, k)) {
      if (tagListDom[k].dataset.path === visitedViews.value[currentIndex - 1]?.path) {
        prevTag = tagListDom[k];
      }

      if (tagListDom[k].dataset.path === visitedViews.value[currentIndex + 1]?.path) {
        nextTag = tagListDom[k];
      }
    }
  }

  if (!prevTag || !nextTag) return;

  const afterNextTagOffsetLeft =
    nextTag.offsetLeft + nextTag.offsetWidth + tagAndTagSpacing.value;

  const beforePrevTagOffsetLeft =
    prevTag.offsetLeft - tagAndTagSpacing.value;

  if (afterNextTagOffsetLeft > scrollWrapper.scrollLeft + containerWidth) {
    scrollWrapper.scrollLeft = afterNextTagOffsetLeft - containerWidth;
  } else if (beforePrevTagOffsetLeft < scrollWrapper.scrollLeft) {
    scrollWrapper.scrollLeft = beforePrevTagOffsetLeft;
  }
}

defineExpose({
  moveToTarget
});
</script>

<style lang="scss" scoped>
.scroll-container {
  white-space: nowrap;
  position: relative;
  overflow-x: auto;
  overflow-y: hidden;
  width: 100%;
  height: 34px;
  scroll-behavior: smooth;
}

.scroll-content {
  display: inline-flex;
  min-width: 100%;
  height: 100%;
  align-items: center;
}

/* 隐藏横向滚动条，保持和 scroll-container 类似的视觉 */
.scroll-container::-webkit-scrollbar {
  height: 0;
}

.scroll-container {
  scrollbar-width: none;
}
</style>
