<template><Icon v-if="iconName" :icon="iconName" :width="size" :height="size" :color="color || undefined" class="app-icon"/><svg-icon v-else-if="icon" :icon-class="icon" class="app-icon"/></template>
<script setup lang="ts">
import { computed } from 'vue';import { Icon } from '@iconify/vue';import { resolveIconName } from './iconRegistry'
const props=withDefaults(defineProps<{icon?:string;size?:string|number;color?:string}>(),{icon:'',size:'1em',color:''});const iconName=computed(()=>resolveIconName(props.icon))
</script><style scoped>.app-icon{display:inline-block;vertical-align:-0.15em;flex-shrink:0}</style>
