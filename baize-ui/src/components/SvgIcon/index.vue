<template>
  <Icon v-if="iconData" :icon="iconData" :class="svgClass" aria-hidden="true" />
  <svg v-else :class="svgClass" aria-hidden="true">
    <use :xlink:href="iconName" :fill="color" />
  </svg>
</template>

<script>
import { Icon } from '@iconify/vue'
import { resolveIconData } from '@/components/AppIcon/iconRegistry'
export default defineComponent({
	components: { Icon },
  props: {
    iconClass: {
      type: String,
      required: true
    },
    className: {
      type: String,
      default: ''
    },
    color: {
      type: String,
      default: ''
    },
  },
  setup(props) {
    return {
      iconName: computed(() => `#icon-${props.iconClass}`),
	  iconData: computed(() => resolveIconData(props.iconClass)),
      svgClass: computed(() => {
        if (props.className) {
          return `svg-icon ${props.className}`
        }
        return 'svg-icon'
      })
    }
  }
})
</script>

<style scope lang="scss">
.nav-icon {
  display: inline-block;
  font-size: 15px;
  margin-right: 12px;
  position: relative;
}

.svg-icon {
  width: 1em;
  height: 1em;
  position: relative;
  fill: currentColor;
  vertical-align: -2px;
}
</style>
