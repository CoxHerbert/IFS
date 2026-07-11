<template>
  <div class="rich-text-editor">
    <Toolbar
      class="rich-text-editor__toolbar"
      :editor="editorRef"
      :default-config="toolbarConfig"
      mode="default"
    />
    <Editor
      v-model="editorValue"
      class="rich-text-editor__body"
      :default-config="editorConfig"
      mode="default"
      @onCreated="handleCreated"
      @onChange="handleChange"
    />
  </div>
</template>

<script setup>
import '@wangeditor/editor/dist/css/style.css'
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'
import { onBeforeUnmount, ref, shallowRef, watch } from 'vue'
import request from '@/utils/request'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: '请输入正文内容'
  }
})

const emit = defineEmits(['update:modelValue'])

const editorRef = shallowRef()
const editorValue = ref(props.modelValue || '')

watch(
  () => props.modelValue,
  value => {
    const nextValue = value || ''
    if (nextValue !== editorValue.value) {
      editorValue.value = nextValue
    }
  }
)

watch(editorValue, value => {
  if (value !== props.modelValue) {
    emit('update:modelValue', value)
  }
})

const toolbarConfig = {
  excludeKeys: ['fullScreen']
}

const editorConfig = {
  placeholder: props.placeholder,
  readOnly: false,
  MENU_CONF: {
    uploadImage: {
      fieldName: 'file',
      maxFileSize: 5 * 1024 * 1024,
      allowedFileTypes: ['image/*'],
      async customUpload(file, insertFn) {
        const formData = new FormData()
        formData.append('file', file)
        const response = await request({
          url: '/cms/article/upload-image',
          method: 'post',
          data: formData,
          headers: { 'Content-Type': 'multipart/form-data' },
          timeout: 30000
        })
        const url = response.data?.url
        if (url) {
          insertFn(url, file.name, url)
        }
      }
    }
  }
}

function handleCreated(editor) {
  editorRef.value = editor
  editor.enable?.()
}

function handleChange(editor) {
  editorValue.value = editor.getHtml()
}

onBeforeUnmount(() => {
  editorRef.value?.destroy()
})
</script>

<style scoped>
.rich-text-editor {
  border: 1px solid #d9d9d9;
  border-radius: 6px;
  overflow: hidden;
  background: #fff;
}

.rich-text-editor__toolbar {
  border-bottom: 1px solid #edf0f5;
}

.rich-text-editor__body {
  height: 320px;
  overflow-y: hidden;
}

.rich-text-editor :deep(.w-e-text-container) {
  min-height: 320px;
}
</style>
