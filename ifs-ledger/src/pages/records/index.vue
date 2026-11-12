<template>
  <view class="page"><view class="heading"><picker mode="date" fields="month" :value="selectedMonth" @change="onMonthChange"><view class="month">{{ monthLabel }}　⌄</view></picker><button class="add-button" @tap="openEditor">＋ 手动记支出</button></view><view class="totals"><view><text>收入</text><text class="in">+ ¥{{ income.toFixed(2) }}</text></view><view><text>支出</text><text class="out">− ¥{{ expense.toFixed(2) }}</text></view></view><view v-if="filteredRecords.length" class="list"><view v-for="item in filteredRecords" :key="item.id" class="row"><view class="icon">{{ item.category.slice(0,1) }}</view><view class="body"><text class="merchant">{{ item.merchant }}</text><text class="meta">{{ item.date }} · {{ sourceName(item.source) }} · {{ item.category }}</text><text v-if="item.importBatchId" class="batch">批次 {{ item.importBatchId }}</text></view><text :class="item.type">{{ item.type === 'income' ? '+' : '−' }}{{ item.amount.toFixed(2) }}</text></view></view><view v-else class="empty">{{ selectedMonth }}没有账单，可切换月份或手动添加支出</view><view v-if="editorVisible" class="mask" @tap.self="closeEditor"><view class="editor"><view class="editor-heading"><text class="editor-title">添加支出</text><text class="close" @tap="closeEditor">×</text></view><text class="field-label">金额（元）</text><input v-model="form.amount" class="field-input amount-input" type="digit" placeholder="0.00" /><text class="field-label">分类</text><view class="categories"><text v-for="category in categories" :key="category" :class="['category-option', { selected: form.category === category }]" @tap="form.category = category">{{ category }}</text></view><text class="field-label">商户 / 备注</text><input v-model="form.merchant" class="field-input" maxlength="80" placeholder="例如：午餐、交通" /><text class="field-label">时间</text><picker mode="date" :value="form.date" @change="onDateChange"><view class="field-input picker-value">{{ form.date }}　›</view></picker><view class="editor-actions"><button class="cancel-button" @tap="closeEditor">取消</button><button class="save-button" @tap="saveExpense">保存支出</button></view></view></view><LedgerTabBar current="records" /></view>
</template>
<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { getRecords, saveRecords, type LedgerRecord } from '../../utils/ledger'
import LedgerTabBar from '../../components/LedgerTabBar.vue'
const records = ref<LedgerRecord[]>([])
onShow(() => { records.value = getRecords() })
const today = () => { const d = new Date(); return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}` }
const selectedMonth = ref(today().slice(0, 7))
const monthLabel = computed(() => { const [year, month] = selectedMonth.value.split('-'); return `${year}年${Number(month)}月` })
const filteredRecords = computed(() => records.value.filter(item => item.date.slice(0, 7) === selectedMonth.value).sort((a, b) => b.date.localeCompare(a.date)))
const editorVisible = ref(false)
const categories = ['餐饮', '交通', '购物', '日用', '住房', '医疗', '娱乐', '其他']
const form = reactive({ amount: '', category: '餐饮', merchant: '', date: today() })
const income = computed(() => filteredRecords.value.filter(item => item.type === 'income').reduce((sum, item) => sum + item.amount, 0))
const expense = computed(() => filteredRecords.value.filter(item => item.type === 'expense').reduce((sum, item) => sum + item.amount, 0))
function sourceName(source: LedgerRecord['source']) { return source === 'wechat' ? '微信' : source === 'alipay' ? '支付宝' : '手动记账' }
function openEditor() { form.amount = ''; form.category = '餐饮'; form.merchant = ''; form.date = today(); editorVisible.value = true }
function closeEditor() { editorVisible.value = false }
function onDateChange(event: { detail: { value: string } }) { form.date = event.detail.value }
function onMonthChange(event: { detail: { value: string } }) { selectedMonth.value = event.detail.value.slice(0, 7) }
function saveExpense() {
  const amount = Number(form.amount)
  const merchant = form.merchant.trim()
  if (!Number.isFinite(amount) || amount <= 0) { uni.showToast({ title: '请输入大于 0 的金额', icon: 'none' }); return }
  if (!merchant) { uni.showToast({ title: '请填写商户或备注', icon: 'none' }); return }
  const record: LedgerRecord = { id: `manual-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`, date: `${form.date} ${new Date().toTimeString().slice(0, 8)}`, merchant, category: form.category, amount: Math.round(amount * 100) / 100, type: 'expense', source: 'manual' }
  const updated = [...getRecords(), record].sort((a, b) => b.date.localeCompare(a.date))
  saveRecords(updated)
  records.value = updated
  closeEditor()
  uni.showToast({ title: '支出已添加', icon: 'success' })
}
</script>
<style scoped>
.page{padding:30rpx 32rpx 170rpx}.heading{display:flex;align-items:center;justify-content:space-between;padding:5rpx 0 20rpx}.month{font-size:26rpx;font-weight:650}.add-button{margin:0;padding:0 20rpx;height:64rpx;line-height:64rpx;border-radius:18rpx;background:#347662;color:#fff;font-size:20rpx}.add-button::after,.cancel-button::after,.save-button::after{border:0}.totals{display:flex;background:#fff;border-radius:22rpx;padding:28rpx;margin-bottom:24rpx}.totals>view{flex:1;display:flex;flex-direction:column;gap:12rpx;color:#93999d;font-size:21rpx}.in,.out{font-size:25rpx;font-weight:650;color:#343b38}.list{background:white;border-radius:22rpx;padding:0 24rpx}.row{min-height:110rpx;display:flex;align-items:center;border-bottom:1rpx solid #f0f1f1;gap:16rpx}.row:last-child{border:0}.icon{width:58rpx;height:58rpx;border-radius:18rpx;background:#eaf3ef;color:#437565;display:flex;align-items:center;justify-content:center}.body{flex:1;display:flex;flex-direction:column;gap:8rpx}.merchant{font-size:23rpx;font-weight:600}.meta{font-size:18rpx;color:#9ba1a4}.batch{font-size:16rpx;color:#4e8b71;font-family:monospace}.expense{color:#333;font-weight:650}.income{color:#39966e;font-weight:650}.empty{margin-top:180rpx;text-align:center;color:#9ba1a4;font-size:23rpx;line-height:1.8}.mask{position:fixed;z-index:1000;inset:0;background:rgba(20,34,27,.45);display:flex;align-items:flex-end}.editor{box-sizing:border-box;width:100%;padding:30rpx 34rpx calc(30rpx + env(safe-area-inset-bottom));border-radius:30rpx 30rpx 0 0;background:#f8faf9}.editor-heading{display:flex;justify-content:space-between;align-items:center;margin-bottom:22rpx}.editor-title{font-size:30rpx;font-weight:700;color:#25372f}.close{font-size:44rpx;color:#9aa39e;padding:0 10rpx}.field-label{display:block;font-size:19rpx;font-weight:650;color:#68736d;margin:18rpx 0 9rpx}.field-input{box-sizing:border-box;width:100%;height:78rpx;padding:0 20rpx;background:#fff;border:1rpx solid #e7ece9;border-radius:16rpx;font-size:22rpx;color:#25372f}.amount-input{font-size:32rpx;font-weight:700}.categories{display:flex;flex-wrap:wrap;gap:12rpx}.category-option{padding:12rpx 20rpx;border-radius:14rpx;background:#fff;color:#66716b;font-size:19rpx}.category-option.selected{background:#e5f1eb;color:#347662;font-weight:650}.picker-value{display:flex;align-items:center;color:#34423a}.editor-actions{display:flex;gap:16rpx;margin-top:28rpx}.cancel-button,.save-button{flex:1;height:82rpx;line-height:82rpx;margin:0;border-radius:18rpx;font-size:22rpx}.cancel-button{background:#fff;color:#59645d}.save-button{background:#347662;color:#fff}
</style>
