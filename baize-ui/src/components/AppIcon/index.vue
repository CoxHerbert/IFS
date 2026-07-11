<template><Icon v-if="iconData" :icon="iconData" :width="size" :height="size" :color="color || undefined" class="app-icon"/><svg-icon v-else-if="icon" :icon-class="icon" class="app-icon"/></template>
<script setup lang="ts">
import { computed } from 'vue';import { Icon } from '@iconify/vue';import { resolveIconData } from './iconRegistry'
const props=withDefaults(defineProps<{icon?:string;size?:string|number;color?:string}>(),{icon:'',size:'1em',color:''});const iconData=computed(()=>resolveIconData(props.icon))
</script><style scoped>.app-icon{display:inline-block;vertical-align:-0.15em;flex-shrink:0}</style>
