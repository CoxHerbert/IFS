<template>
  <div class="app-container">
    <a-form v-show="showSearch" ref="queryRef" :model="queryParams" layout="inline" class="search-form">
      <a-form-item label="标题" name="title">
        <a-input v-model:value="queryParams.title" allow-clear placeholder="请输入标题" @pressEnter="handleQuery" />
      </a-form-item>
      <a-form-item label="栏目" name="category">
        <a-select v-model:value="queryParams.category" allow-clear placeholder="请选择栏目" style="width: 150px" :options="categoryOptions" />
      </a-form-item>
      <a-form-item label="状态" name="status">
        <a-select v-model:value="queryParams.status" allow-clear placeholder="请选择状态" style="width: 130px" :options="statusOptions" />
      </a-form-item>
      <a-form-item>
        <a-space>
          <a-button type="primary" @click="handleQuery">搜索</a-button>
          <a-button @click="resetQuery">重置</a-button>
        </a-space>
      </a-form-item>
    </a-form>

    <div class="toolbar-row mb8">
      <a-space>
        <a-button type="primary" @click="handleAdd" v-hasPermi="['cms:article:add']">新增</a-button>
        <a-button :disabled="single" @click="handleUpdate()" v-hasPermi="['cms:article:edit']">修改</a-button>
        <a-button danger :disabled="multiple" @click="handleDelete()" v-hasPermi="['cms:article:remove']">删除</a-button>
      </a-space>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList" />
    </div>

    <vxe-table
      ref="articleTableRef"
      border
      stripe
      auto-resize
      show-overflow="title"
      :loading="loading"
      :data="articleList"
      :row-config="{ keyField: 'articleId' }"
      :checkbox-config="{ reserve: true }"
      @checkbox-change="handleSelectionChange"
      @checkbox-all="handleSelectionChange"
    >
      <vxe-column type="checkbox" width="55" align="center" />
      <vxe-column field="coverUrl" title="封面" width="110" align="center">
        <template #default="{ row }">
          <a-image
            v-if="row.coverUrl"
            :src="row.coverUrl"
            :width="72"
            :height="45"
            class="cover-thumbnail"
          />
          <span v-else>-</span>
        </template>
      </vxe-column>
      <vxe-column field="title" title="标题" min-width="220" />
      <vxe-column field="category" title="栏目" width="110" align="center">
        <template #default="{ row }">
          <a-tag color="blue">{{ row.category }}</a-tag>
        </template>
      </vxe-column>
      <vxe-column field="status" title="状态" width="100" align="center">
        <template #default="{ row }">
          <a-tag :color="row.status === '0' ? 'green' : 'orange'">{{ row.status === '0' ? '已发布' : '草稿' }}</a-tag>
        </template>
      </vxe-column>
      <vxe-column field="sort" title="排序" width="80" align="center" />
      <vxe-column field="publishTime" title="发布时间" width="250" align="center">
        <template #default="{ row }">
          <span v-if="row.publishTime">{{ formatDateTime(row.publishTime) }}</span>
          <span v-else-if="row.createTime" class="draft-time">
            未发布（创建于 {{ formatDateTime(row.createTime) }}）
          </span>
          <span v-else>-</span>
        </template>
      </vxe-column>
      <vxe-column field="updateBy" title="更新人" width="100" align="center" />
      <vxe-column title="操作" width="230" align="center" fixed="right">
        <template #default="{ row }">
          <a-space>
            <a-button type="link" @click="handlePreview(row)">预览</a-button>
            <a-button type="link" @click="handleUpdate(row)" v-hasPermi="['cms:article:edit']">修改</a-button>
            <a-button type="link" danger @click="handleDelete(row)" v-hasPermi="['cms:article:remove']">删除</a-button>
          </a-space>
        </template>
      </vxe-column>
    </vxe-table>

    <pagination v-show="total > 0" v-model:page="queryParams.pageNum" v-model:limit="queryParams.pageSize" :total="total" @pagination="getList" />

    <a-modal v-model:open="open" :title="title" width="980px" :footer="null" destroy-on-close>
      <a-form ref="articleRef" :model="form" :rules="rules" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="16">
            <a-form-item label="标题" name="title">
              <a-input v-model:value="form.title" placeholder="请输入文章标题" />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item label="栏目" name="category">
              <a-select v-model:value="form.category" placeholder="请选择栏目" :options="categoryOptions" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="访问标识 Slug" name="slug">
              <a-input v-model:value="form.slug" placeholder="留空则按标题生成" />
            </a-form-item>
          </a-col>
          <a-col :span="6">
            <a-form-item label="排序" name="sort">
              <a-input-number v-model:value="form.sort" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="6">
            <a-form-item label="状态" name="status">
              <a-select v-model:value="form.status" :options="statusOptions" />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item label="摘要" name="summary">
              <a-textarea v-model:value="form.summary" :rows="3" placeholder="用于列表页展示" />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item label="封面图" name="coverUrl">
              <div class="cover-field">
                <a-upload
                  accept="image/png,image/jpeg,image/gif,image/webp"
                  :show-upload-list="false"
                  :custom-request="uploadCover"
                  :before-upload="beforeCoverUpload"
                >
                  <div v-if="form.coverUrl" class="cover-preview">
                    <img :src="form.coverUrl" alt="文章封面预览" />
                    <div class="cover-preview__mask">更换封面</div>
                  </div>
                  <a-button v-else :loading="coverUploading">上传封面</a-button>
                </a-upload>
                <div class="cover-actions">
                  <span>建议使用 16:9 图片，支持 PNG、JPG、GIF、WEBP，最大 5MB</span>
                  <a-button v-if="form.coverUrl" type="link" danger @click="form.coverUrl = ''">移除封面</a-button>
                </div>
              </div>
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item label="正文" name="content">
              <rich-text-editor v-model="form.content" placeholder="请输入文章正文" />
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
      <div class="modal-footer">
        <a-space>
          <a-button type="primary" @click="submitForm">保存</a-button>
          <a-button @click="cancel">取消</a-button>
        </a-space>
      </div>
    </a-modal>
  </div>
</template>

<script setup name="CmsArticle">
import RichTextEditor from '@/components/RichTextEditor/index.vue'
import { addArticle, delArticle, getArticle, listArticle, updateArticle, uploadArticleImage } from '@/api/cms/article'

const { proxy } = getCurrentInstance()
const portalBaseUrl = (import.meta.env.VITE_PORTAL_BASE_URL || window.location.origin).replace(/\/$/, '')

const categoryOptions = [
  { label: '美线', value: '美线' },
  { label: '欧线', value: '欧线' },
  { label: '东南亚', value: '东南亚' },
  { label: '中东', value: '中东' },
  { label: '拼箱', value: '拼箱' },
  { label: '政策', value: '政策' }
]
const statusOptions = [
  { label: '已发布', value: '0' },
  { label: '草稿', value: '1' }
]

const articleList = ref([])
const articleTableRef = ref()
const open = ref(false)
const loading = ref(true)
const showSearch = ref(true)
const ids = ref([])
const single = ref(true)
const multiple = ref(true)
const total = ref(0)
const title = ref('')
const coverUploading = ref(false)

const data = reactive({
  form: {},
  queryParams: {
    pageNum: 1,
    pageSize: 10,
    title: undefined,
    category: undefined,
    status: undefined
  },
  rules: {
    title: [{ required: true, message: '标题不能为空', trigger: 'blur' }],
    category: [{ required: true, message: '栏目不能为空', trigger: 'change' }],
    content: [{ required: true, message: '正文不能为空', trigger: 'blur' }]
  }
})

const { queryParams, form, rules } = toRefs(data)

function getList() {
  loading.value = true
  listArticle(queryParams.value)
    .then(response => {
      const result = response.data || {}
      articleList.value = result.rows || []
      total.value = result.total || 0
    })
    .finally(() => {
      loading.value = false
    })
}

function reset() {
  form.value = {
    articleId: undefined,
    title: undefined,
    slug: undefined,
    summary: undefined,
    category: undefined,
    coverUrl: '',
    content: '',
    status: '0',
    sort: 0
  }
  proxy.resetForm('articleRef')
}

function cancel() {
  open.value = false
  reset()
}

function handleQuery() {
  queryParams.value.pageNum = 1
  getList()
}

function resetQuery() {
  proxy.resetForm('queryRef')
  articleTableRef.value?.clearCheckboxRow()
  handleQuery()
}

function handleSelectionChange() {
  const records = articleTableRef.value?.getCheckboxRecords() || []
  ids.value = records.map(item => item.articleId)
  single.value = records.length !== 1
  multiple.value = !records.length
}

function handleAdd() {
  reset()
  open.value = true
  title.value = '新增文章'
}

function handleUpdate(row) {
  reset()
  const articleId = row?.articleId || ids.value[0]
  getArticle(articleId).then(response => {
    form.value = response.data || {}
    open.value = true
    title.value = '修改文章'
  })
}

function submitForm() {
  proxy.$refs.articleRef.validate().then(() => {
    const request = form.value.articleId ? updateArticle(form.value) : addArticle(form.value)
    request.then(() => {
      proxy.$modal.msgSuccess(form.value.articleId ? '修改成功' : '新增成功')
      open.value = false
      getList()
    })
  }).catch(() => {})
}

function handleDelete(row) {
  const articleIds = row?.articleId || ids.value
  proxy.$modal.confirm(`是否确认删除文章编号为"${articleIds}"的数据项？`).then(() => delArticle(articleIds)).then(() => {
    getList()
    proxy.$modal.msgSuccess('删除成功')
  }).catch(() => {})
}

function handlePreview(row) {
  if (!row.slug) {
    proxy.$modal.msgWarning('请先保存生成访问标识')
    return
  }
  if (row.status !== '0') {
    proxy.$modal.msgWarning('草稿尚未对门户发布，请发布后再预览')
    return
  }
  window.open(`${portalBaseUrl}/news/${encodeURIComponent(row.slug)}`, '_blank', 'noopener,noreferrer')
}

function beforeCoverUpload(file) {
  const allowedTypes = ['image/png', 'image/jpeg', 'image/gif', 'image/webp']
  if (!allowedTypes.includes(file.type)) {
    proxy.$modal.msgError('封面仅支持 PNG、JPG、GIF、WEBP 格式')
    return false
  }
  if (file.size > 5 * 1024 * 1024) {
    proxy.$modal.msgError('封面图片不能超过 5MB')
    return false
  }
  return true
}

async function uploadCover({ file, onSuccess, onError }) {
  coverUploading.value = true
  try {
    const response = await uploadArticleImage(file)
    form.value.coverUrl = response.data?.url || ''
    if (!form.value.coverUrl) throw new Error('上传结果中缺少图片地址')
    proxy.$modal.msgSuccess('封面上传成功')
    onSuccess?.(response)
  } catch (error) {
    onError?.(error)
  } finally {
    coverUploading.value = false
  }
}

function formatDateTime(value) {
  if (value === undefined || value === null || value === '') return '-'

  let date
  if (typeof value === 'number' || /^\d+$/.test(String(value))) {
    const timestamp = Number(value)
    date = new Date(timestamp < 1e12 ? timestamp * 1000 : timestamp)
  } else {
    date = new Date(String(value).replace(/-/g, '/'))
  }

  if (Number.isNaN(date.getTime())) return String(value)
  const pad = number => String(number).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`
}

getList()
</script>

<style scoped>
.search-form { margin-bottom: 16px; }
.toolbar-row { display: flex; align-items: center; justify-content: space-between; gap: 16px; }
.modal-footer { display: flex; justify-content: flex-end; margin-top: 24px; }
.draft-time { color: #8c8c8c; }
.cover-thumbnail { border-radius: 4px; object-fit: cover; }
.cover-field { display: grid; gap: 8px; justify-items: start; }
.cover-preview { position: relative; width: 240px; aspect-ratio: 16 / 9; overflow: hidden; border-radius: 6px; background: #f5f5f5; cursor: pointer; }
.cover-preview img { width: 100%; height: 100%; object-fit: cover; }
.cover-preview__mask { position: absolute; inset: 0; display: grid; place-items: center; color: #fff; background: rgba(0, 0, 0, 0.45); opacity: 0; transition: opacity 0.2s; }
.cover-preview:hover .cover-preview__mask { opacity: 1; }
.cover-actions { display: flex; align-items: center; gap: 8px; color: #8c8c8c; font-size: 12px; }
</style>
