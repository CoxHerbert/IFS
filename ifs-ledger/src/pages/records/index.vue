<template>
  <view class="page">
    <view class="heading">
      <text class="page-title">全部账单</text>
      <button class="add-button" @tap="openEditor">＋ 手动记支出</button>
    </view>

    <view class="filters">
      <view class="filter-row">
        <text class="filter-label">时间</text>
        <view class="options">
          <text v-for="option in dateOptions" :key="option.value" :class="['option', { selected: dateMode === option.value }]" @tap="dateMode = option.value">{{ option.label }}</text>
        </view>
        <view v-if="dateMode !== 'all'" class="period-picker" @tap="periodPickerVisible = true">
          <text>{{ periodLabel }}</text><wd-icon name="caret-down" size="13" color="var(--theme-primary, #756ab2)" />
        </view>
      </view>
      <view class="filter-row">
        <text class="filter-label">收支</text>
        <view class="options">
          <text v-for="option in typeOptions" :key="option.value" :class="['option', { selected: typeFilter === option.value }]" @tap="typeFilter = option.value">{{ option.label }}</text>
        </view>
      </view>
      <view class="filter-row">
        <text class="filter-label">来源</text>
        <view class="options">
          <text v-for="option in sourceOptions" :key="option.value" :class="['option', { selected: sourceFilter === option.value }]" @tap="sourceFilter = option.value">{{ option.label }}</text>
        </view>
      </view>
      <view class="filter-row category-filter">
        <text class="filter-label">分类</text>
        <scroll-view scroll-x class="category-scroll" :show-scrollbar="false">
          <view class="options category-options">
            <text v-for="category in categoryOptions" :key="category" :class="['option', { selected: categoryFilter === category }]" @tap="categoryFilter = category">{{ category }}</text>
          </view>
        </scroll-view>
      </view>
    </view>

    <view class="totals">
      <view><text>收入</text><text class="in">+ ¥{{ income.toFixed(2) }}</text></view>
      <view><text>支出</text><text class="out">− ¥{{ expense.toFixed(2) }}</text></view>
      <view><text>记录</text><text class="count">{{ filteredRecords.length }} 笔</text></view>
    </view>

    <view v-if="filteredRecords.length" class="list">
      <view v-for="item in filteredRecords" :key="item.id" class="row">
        <view class="icon">{{ item.category.slice(0,1) }}</view>
        <view class="body">
          <text class="merchant">{{ item.merchant }}</text>
          <text class="meta">{{ item.date }} · {{ sourceName(item.source) }} · {{ item.category }}</text>
          <text v-if="item.importBatchId" class="batch">批次 {{ item.importBatchId }}</text>
        </view>
        <view class="row-actions"><text :class="item.type">{{ item.type === 'income' ? '+' : '−' }}{{ item.amount.toFixed(2) }}</text><view class="record-actions"><text class="edit-record" @tap="openEdit(item)">编辑</text><text class="delete-record" @tap="deleteRecord(item)">删除</text></view></view>
      </view>
    </view>
    <view v-else class="empty">当前筛选条件下没有账单</view>

    <view v-if="editorVisible" class="mask" @tap.self="closeEditor">
      <view class="editor">
        <view class="editor-heading"><text class="editor-title">{{ editingRecordId ? '编辑账单' : '添加支出' }}</text><text class="close" @tap="closeEditor">×</text></view>
        <text class="field-label">金额（元）</text><input v-model="form.amount" class="field-input amount-input" type="digit" placeholder="0.00" />
        <template v-if="editingRecordId"><text class="field-label">收支类型</text><view class="categories"><text v-for="option in editTypeOptions" :key="option.value" :class="['category-option', { selected: form.type === option.value }]" @tap="form.type = option.value">{{ option.label }}</text></view></template>
        <text class="field-label">分类</text>
        <view class="categories"><text v-for="category in manualCategories" :key="category" :class="['category-option', { selected: form.category === category }]" @tap="form.category = category">{{ category }}</text></view>
        <text class="field-label">商户 / 备注</text><input v-model="form.merchant" class="field-input" maxlength="80" placeholder="例如：午餐、交通" />
        <text class="field-label">时间</text><view class="field-input picker-value" @tap="openEditorDatePicker">{{ form.date }}<wd-icon name="calendar-line" size="18" color="var(--theme-primary, #756ab2)" /></view>
        <view class="editor-actions"><button class="cancel-button" @tap="closeEditor">取消</button><button class="save-button" :disabled="savingRecord" @tap="saveExpense">{{ savingRecord ? '保存中…' : editingRecordId ? '保存修改' : '保存支出' }}</button></view>
      </view>
    </view>
    <wd-datetime-picker :visible="periodPickerVisible" :model-value="periodPickerValue" :type="dateMode === 'year' ? 'year' : 'year-month'" title="选择查看时间" :root-portal="true" :z-index="2000" @update:visible="periodPickerVisible = $event" @confirm="onPeriodConfirm" />
    <wd-datetime-picker :visible="editorDatePickerVisible" :model-value="editorDatePickerValue" type="date" title="选择支出日期" :root-portal="true" :z-index="2100" @update:visible="editorDatePickerVisible = $event" @confirm="onEditorDateConfirm" />
    <LedgerTabBar current="ledger" />
  </view>
</template>

<script setup lang="ts">
import { computed, nextTick, reactive, ref } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import { addLedgerRecord, getRecords, removeLedgerRecord, updateLedgerRecord, type LedgerRecord } from '../../utils/ledger'
import LedgerTabBar from '../../components/LedgerTabBar.vue'

type DateMode = 'month' | 'year' | 'all'
type TypeFilter = 'all' | LedgerRecord['type']
type SourceFilter = 'all' | LedgerRecord['source']

const records = ref<LedgerRecord[]>([])
onShow(() => { records.value = getRecords() })
onLoad(options => { if (options?.manual === '1') openEditor() })
const today = () => { const d = new Date(); return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}` }
const dateMode = ref<DateMode>('month')
const selectedMonth = ref(today().slice(0, 7))
const selectedYear = ref(today().slice(0, 4))
const typeFilter = ref<TypeFilter>('all')
const sourceFilter = ref<SourceFilter>('all')
const categoryFilter = ref('全部')
const dateOptions: { value: DateMode; label: string }[] = [{ value: 'month', label: '月' }, { value: 'year', label: '年' }, { value: 'all', label: '全部' }]
const typeOptions: { value: TypeFilter; label: string }[] = [{ value: 'all', label: '全部' }, { value: 'expense', label: '支出' }, { value: 'income', label: '收入' }]
const sourceOptions: { value: SourceFilter; label: string }[] = [{ value: 'all', label: '全部' }, { value: 'wechat', label: '微信' }, { value: 'alipay', label: '支付宝' }, { value: 'manual', label: '手动' }]
const periodLabel = computed(() => dateMode.value === 'year' ? `${selectedYear.value}年` : `${selectedMonth.value.slice(0, 4)}年${Number(selectedMonth.value.slice(5, 7))}月`)
const categoryOptions = computed(() => ['全部', ...new Set(records.value.map(item => item.category).filter(Boolean))])
const filteredRecords = computed(() => records.value.filter(item => {
  if (dateMode.value === 'month' && item.date.slice(0, 7) !== selectedMonth.value) return false
  if (dateMode.value === 'year' && item.date.slice(0, 4) !== selectedYear.value) return false
  if (typeFilter.value !== 'all' && item.type !== typeFilter.value) return false
  if (sourceFilter.value !== 'all' && item.source !== sourceFilter.value) return false
  if (categoryFilter.value !== '全部' && item.category !== categoryFilter.value) return false
  return true
}).sort((a, b) => b.date.localeCompare(a.date)))
const income = computed(() => filteredRecords.value.filter(item => item.type === 'income').reduce((sum, item) => sum + item.amount, 0))
const expense = computed(() => filteredRecords.value.filter(item => item.type === 'expense').reduce((sum, item) => sum + item.amount, 0))
const editorVisible = ref(false)
const editingRecordId = ref<string | null>(null)
const savingRecord = ref(false)
const periodPickerVisible = ref(false)
const editorDatePickerVisible = ref(false)
const manualCategories = ['餐饮', '交通', '购物', '日用', '住房', '医疗', '娱乐', '其他']
const form = reactive({ amount: '', category: '餐饮', merchant: '', date: today(), type: 'expense' as LedgerRecord['type'] })
const editTypeOptions: { value: LedgerRecord['type']; label: string }[] = [{ value: 'expense', label: '支出' }, { value: 'income', label: '收入' }]
const periodPickerValue = computed(() => new Date(`${dateMode.value === 'year' ? selectedYear.value : `${selectedMonth.value}-01`}T00:00:00`).getTime())
const editorDatePickerValue = computed(() => new Date(`${form.date}T00:00:00`).getTime())

function onPeriodConfirm(event: { value: number | string }) {
  const date = new Date(Number(event.value))
  if (dateMode.value === 'year') selectedYear.value = String(date.getFullYear())
  else selectedMonth.value = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`
}
function sourceName(source: LedgerRecord['source']) { return source === 'wechat' ? '微信' : source === 'alipay' ? '支付宝' : '手动记账' }
function openEditor() { editingRecordId.value = null; form.amount = ''; form.category = '餐饮'; form.merchant = ''; form.date = today(); form.type = 'expense'; editorVisible.value = true }
function openEdit(item: LedgerRecord) {
  editingRecordId.value = item.id
  form.amount = String(item.amount)
  form.category = item.category
  form.merchant = item.merchant
  form.date = item.date.slice(0, 10)
  form.type = item.type
  editorVisible.value = true
}
function closeEditor() { editorVisible.value = false }
async function openEditorDatePicker() {
  editorDatePickerVisible.value = false
  await nextTick()
  editorDatePickerVisible.value = true
}
function onEditorDateConfirm(event: { value: number | string }) {
  const date = new Date(Number(event.value))
  form.date = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}
async function saveExpense() {
  if (savingRecord.value) return
  const amount = Number(form.amount)
  const merchant = form.merchant.trim()
  if (!Number.isFinite(amount) || amount <= 0) { uni.showToast({ title: '请输入大于 0 的金额', icon: 'none' }); return }
  if (!merchant) { uni.showToast({ title: '请填写商户或备注', icon: 'none' }); return }
  savingRecord.value = true
  try {
    const roundedAmount = Math.round(amount * 100) / 100
    let updated: LedgerRecord[]
    if (editingRecordId.value) {
      const current = getRecords().find(item => item.id === editingRecordId.value)
      if (!current) throw new Error('账单已不存在，请刷新后重试')
      const time = current.date.slice(11, 19) || new Date().toTimeString().slice(0, 8)
      const record: LedgerRecord = { ...current, date: `${form.date} ${time}`, merchant, category: form.category, amount: roundedAmount, type: form.type }
      updated = await updateLedgerRecord(record)
    } else {
      const record: LedgerRecord = { id: `manual-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`, date: `${form.date} ${new Date().toTimeString().slice(0, 8)}`, merchant, category: form.category, amount: roundedAmount, type: 'expense', source: 'manual' }
      updated = await addLedgerRecord(record)
    }
    records.value = updated
    closeEditor()
    uni.showToast({ title: editingRecordId.value ? '账单已修改' : '支出已添加', icon: 'success' })
  } catch (error) {
    uni.showToast({ title: error instanceof Error ? error.message : '保存账单失败', icon: 'none' })
  } finally { savingRecord.value = false }
}
function deleteRecord(item: LedgerRecord) {
  uni.showModal({
    title: '删除这笔账单？',
    content: `将删除「${item.merchant}」${item.type === 'expense' ? '支出' : '收入'} ¥${item.amount.toFixed(2)}，此操作无法撤销。`,
    confirmText: '删除',
    confirmColor: '#c25750',
    success: result => {
      if (!result.confirm) return
      removeLedgerRecord(item.id).then(updated => {
        records.value = updated
        uni.showToast({ title: '账单已删除', icon: 'success' })
      }).catch(error => uni.showToast({ title: error instanceof Error ? error.message : '删除账单失败', icon: 'none' }))
    },
  })
}
</script>

<style scoped>
.page{padding:30rpx 32rpx 170rpx}.heading{display:flex;align-items:center;justify-content:space-between;padding:5rpx 0 20rpx}.page-title{font-size:43rpx;font-weight:750;color:#302c45}.add-button{margin:0;padding:0 20rpx;height:64rpx;line-height:64rpx;border-radius:18rpx;background:var(--theme-primary, #6556b2);color:#fff;font-size:26.9rpx}.add-button::after,.cancel-button::after,.save-button::after{border:0}
.filters{padding:12rpx 18rpx;margin-bottom:20rpx;background:#fff;border-radius:20rpx;box-shadow:0 6rpx 20rpx rgba(54,45,92,.035)}.filter-row{display:flex;align-items:center;gap:12rpx;min-height:64rpx}.filter-row+.filter-row{border-top:1rpx solid #f1eff6}.filter-label{width:60rpx;flex:none;color:#8b879d;font-size:25.6rpx}.options{display:flex;align-items:center;gap:8rpx;flex-wrap:nowrap}.option{padding:9rpx 16rpx;border-radius:13rpx;color:#706c80;font-size:24.2rpx;white-space:nowrap}.option.selected{background:var(--theme-soft, #efedf8);color:var(--theme-primary, #6556b2);font-weight:650}.period-picker{display:flex;align-items:center;gap:6rpx;padding:8rpx 12rpx;border-radius:12rpx;background:#f8f7fb;color:var(--theme-primary, #51458e);font-size:24.2rpx;white-space:nowrap}.category-scroll{flex:1;min-width:0;white-space:nowrap}.category-options{display:inline-flex;padding-right:10rpx}
.totals{display:flex;background:#fff;border-radius:22rpx;padding:24rpx;margin-bottom:24rpx;box-shadow:0 6rpx 20rpx rgba(54,45,92,.035)}.totals>view{flex:1;display:flex;flex-direction:column;gap:10rpx;color:#93999d;font-size:25.6rpx}.in,.out,.count{font-size:31rpx;font-weight:650;color:#343b38}.list{background:white;border-radius:22rpx;padding:0 24rpx}.row{min-height:110rpx;display:flex;align-items:center;border-bottom:1rpx solid #efedf3;gap:16rpx}.row:last-child{border:0}.icon{width:58rpx;height:58rpx;border-radius:18rpx;background:var(--theme-soft, #efedf8);color:var(--theme-primary, #6556b2);display:flex;align-items:center;justify-content:center}.body{flex:1;display:flex;flex-direction:column;gap:8rpx}.merchant{font-size:31rpx;font-weight:600}.meta{font-size:24.2rpx;color:#9ba1a4}.batch{font-size:21.5rpx;color:var(--theme-primary, #6556b2);font-family:monospace}.row-actions{display:flex;flex-direction:column;align-items:flex-end;gap:8rpx;flex:none}.record-actions{display:flex;gap:8rpx}.edit-record{padding:5rpx 12rpx;border-radius:10rpx;background:var(--theme-soft, #efedf8);color:var(--theme-primary, #6556b2);font-size:21.5rpx}.expense{color:#333;font-weight:650}.income{color:#39966e;font-weight:650}.delete-record{padding:5rpx 12rpx;border-radius:10rpx;background:#fff0ef;color:#bd5b54;font-size:21.5rpx}.empty{margin-top:150rpx;text-align:center;color:#9ba1a4;font-size:31rpx;line-height:1.8}.mask{position:fixed;z-index:1000;inset:0;background:rgba(36,30,59,.45);display:flex;align-items:flex-end}.editor{box-sizing:border-box;width:100%;padding:30rpx 34rpx calc(30rpx + env(safe-area-inset-bottom));border-radius:30rpx 30rpx 0 0;background:#faf9fd}.editor-heading{display:flex;justify-content:space-between;align-items:center;margin-bottom:22rpx}.editor-title{font-size:40.3rpx;font-weight:700;color:#302c45}.close{font-size:59.2rpx;color:#9aa39e;padding:0 10rpx}.field-label{display:block;font-size:25.6rpx;font-weight:650;color:#68647b;margin:18rpx 0 9rpx}.field-input{box-sizing:border-box;width:100%;height:78rpx;padding:0 20rpx;background:#fff;border:1rpx solid #e9e7f0;border-radius:16rpx;font-size:29.5rpx;color:#302c45}.amount-input{font-size:43rpx;font-weight:700}.categories{display:flex;flex-wrap:wrap;gap:12rpx}.category-option{padding:12rpx 20rpx;border-radius:14rpx;background:#fff;color:#68647b;font-size:25.6rpx}.category-option.selected{background:var(--theme-soft, #e9e6f5);color:var(--theme-primary, #6556b2);font-weight:650}.picker-value{display:flex;align-items:center;color:#302c45}.editor-actions{display:flex;gap:16rpx;margin-top:28rpx}.cancel-button,.save-button{flex:1;height:82rpx;line-height:82rpx;margin:0;border-radius:18rpx;font-size:29.5rpx}.cancel-button{background:#fff;color:#5c5870}.save-button{background:var(--theme-primary, #6556b2);color:#fff}
</style>
