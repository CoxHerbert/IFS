<template>
  <a-popover v-model:open="open" placement="bottomLeft" trigger="click" overlay-class-name="app-icon-picker-overlay">
    <template #content>
      <div class="icon-picker-panel">
        <a-input v-model:value="keyword" allow-clear placeholder="搜索图标名称">
          <template #prefix><SearchOutlined /></template>
        </a-input>
        <div class="icon-picker-grid">
          <button v-for="name in filteredIcons" :key="name" type="button" class="icon-option"
            :class="{ selected: name === modelValue }" :title="name" @click="choose(name)">
            <app-icon :icon="name" :size="22" />
            <span>{{ shortName(name) }}</span>
          </button>
        </div>
        <div v-if="!filteredIcons.length" class="empty-state">没有匹配的图标</div>
      </div>
    </template>
    <a-input :value="modelValue" readonly allow-clear placeholder="点击选择 Iconify 图标" @clear.stop="clear">
      <template #prefix><app-icon v-if="modelValue" :icon="modelValue" :size="20" /><SearchOutlined v-else /></template>
    </a-input>
  </a-popover>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { SearchOutlined } from '@ant-design/icons-vue'
import { iconNames } from '@/components/AppIcon/iconRegistry'
const props=defineProps<{modelValue?:string}>();const emit=defineEmits<{(event:'update:modelValue',value:string|undefined):void}>();const open=ref(false),keyword=ref('')
const filteredIcons=computed(()=>{const value=keyword.value.trim().toLowerCase();return value?iconNames.filter(name=>name.toLowerCase().includes(value)):iconNames})
function shortName(name:string){return name.replace(/^mdi:/,'')}
function choose(name:string){emit('update:modelValue',name);open.value=false}
function clear(){emit('update:modelValue',undefined)}
</script>

<style scoped>
.icon-picker-panel{width:560px;max-width:calc(100vw - 48px)}
.icon-picker-grid{display:grid;grid-template-columns:repeat(4,minmax(0,1fr));gap:8px;max-height:320px;margin-top:12px;padding-right:4px;overflow-y:auto}
.icon-option{display:flex;min-width:0;align-items:center;gap:8px;padding:10px;border:1px solid #e5e7eb;border-radius:8px;background:#fff;color:#374151;cursor:pointer;transition:.15s ease}
.icon-option:hover{border-color:#69b1ff;background:#f0f7ff;color:#1677ff}.icon-option.selected{border-color:#1677ff;background:#e6f4ff;color:#1677ff;box-shadow:0 0 0 1px #1677ff inset}
.icon-option span{overflow:hidden;font-size:12px;text-overflow:ellipsis;white-space:nowrap}.empty-state{padding:36px;text-align:center;color:#8c8c8c}
@media(max-width:720px){.icon-picker-grid{grid-template-columns:repeat(2,minmax(0,1fr))}}
</style>
