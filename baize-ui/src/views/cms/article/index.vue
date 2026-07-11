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

    <a-table
      :loading="loading"
      :data-source="articleList"
      :columns="columns"
      :pagination="false"
      :row-selection="rowSelection"
      row-key="articleId"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'category'">
          <a-tag color="blue">{{ record.category }}</a-tag>
        </template>
        <template v-else-if="column.key === 'status'">
          <a-tag :color="record.status === '0' ? 'green' : 'orange'">{{ record.status === '0' ? '已发布' : '草稿' }}</a-tag>
        </template>
        <template v-else-if="column.key === 'publishTime'">
          {{ record.publishTime || '-' }}
        </template>
        <template v-else-if="column.key === 'action'">
          <a-space>
            <a-button type="link" @click="handlePreview(record)">预览</a-button>
            <a-button type="link" @click="handleUpdate(record)" v-hasPermi="['cms:article:edit']">修改</a-button>
            <a-button type="link" danger @click="handleDelete(record)" v-hasPermi="['cms:article:remove']">删除</a-button>
          </a-space>
        </template>
      </template>
    </a-table>

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
import { addArticle, delArticle, getArticle, listArticle, updateArticle } from '@/api/cms/article'

const { proxy } = getCurrentInstance()

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

const columns = [
  { title: '标题', dataIndex: 'title', key: 'title', ellipsis: true },
  { title: '栏目', dataIndex: 'category', key: 'category', width: 110, align: 'center' },
  { title: '状态', dataIndex: 'status', key: 'status', width: 100, align: 'center' },
  { title: '排序', dataIndex: 'sort', key: 'sort', width: 80, align: 'center' },
  { title: '发布时间', dataIndex: 'publishTime', key: 'publishTime', width: 170, align: 'center' },
  { title: '更新人', dataIndex: 'updateBy', key: 'updateBy', width: 100, align: 'center' },
  { title: '操作', key: 'action', width: 190, align: 'center' }
]

const articleList = ref([])
const open = ref(false)
const loading = ref(true)
const showSearch = ref(true)
const ids = ref([])
const selectedRowKeys = ref([])
const single = ref(true)
const multiple = ref(true)
const total = ref(0)
const title = ref('')

const rowSelection = computed(() => ({
  selectedRowKeys: selectedRowKeys.value,
  onChange: (keys, rows) => handleSelectionChange(rows, keys)
}))

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
  listArticle(queryParams.value).then(response => {
    articleList.value = response.rows || []
    total.value = response.total || 0
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
  handleQuery()
}

function handleSelectionChange(selection, keys) {
  ids.value = selection.map(item => item.articleId)
  selectedRowKeys.value = keys
  single.value = selection.length !== 1
  multiple.value = !selection.length
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
  window.open('/news/' + row.slug, '_blank')
}

getList()
</script>

<style scoped>
.search-form { margin-bottom: 16px; }
.toolbar-row { display: flex; align-items: center; justify-content: space-between; gap: 16px; }
.modal-footer { display: flex; justify-content: flex-end; margin-top: 24px; }
</style>
