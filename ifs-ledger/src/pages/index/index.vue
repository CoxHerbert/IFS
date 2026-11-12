<template>
  <view class="screen">
    <view class="topbar">
      <view>
        <text class="eyebrow">PERSONAL FINANCE</text>
        <view class="title-row"><text class="title">轻记</text><picker mode="date" fields="month" :value="selectedMonth" @change="onMonthChange"><view class="month">{{ monthLabel }}⌄</view></picker></view>
      </view>
      <view class="avatar" @tap="openAccount">我</view>
    </view>

    <view class="summary-card">
      <view class="summary-head"><text>本月结余</text><text class="eye">◉</text></view>
      <view class="balance">¥ {{ formatMoney(balance) }}</view>
      <view class="summary-grid">
        <view><text class="sum-label"><text class="dot income-dot"/>收入</text><text class="sum-value">¥ {{ formatMoney(income) }}</text></view>
        <view><text class="sum-label"><text class="dot expense-dot"/>支出</text><text class="sum-value">¥ {{ formatMoney(expense) }}</text></view>
      </view>
      <view class="card-glow" />
    </view>

    <view class="section-head"><view><text class="section-title">收支分析</text><text class="section-sub">按本月本机账单统计</text></view><text class="more" @tap="openRecords">详情 →</text></view>
    <view class="analysis-card">
      <view class="analysis-top"><text>支出分类</text><text class="muted">本月</text></view>
      <view class="chart-row">
        <view class="ring" :style="{ background: categoryGradient }"><view class="ring-inner"><text class="ring-label">总支出</text><text class="ring-value">{{ compactMoney(expense) }}</text></view></view>
        <view class="categories">
          <view v-if="categories.length" class="categories"><view v-for="item in categories" :key="item.name" class="category-row"><view class="cat-left"><text class="cat-icon" :style="{ background: item.bg }">{{ item.icon }}</text><text>{{ item.name }}</text></view><text class="cat-amount">{{ item.percent }}%</text></view></view>
          <text v-else class="empty-category">本月暂无支出记录</text>
        </view>
      </view>
    </view>

    <view class="section-head recent-head"><view><text class="section-title">最近账单</text><text class="section-sub">本机共 {{ records.length }} 笔</text></view><text class="more" @tap="openRecords">查看全部 →</text></view>
    <view class="bill-card">
      <view v-for="item in recent" :key="item.name" class="bill-row"><view class="bill-icon" :style="{ background: item.bg }">{{ item.icon }}</view><view class="bill-main"><text class="bill-name">{{ item.name }}</text><text class="bill-meta">{{ item.category }} · {{ item.time }}</text></view><text class="bill-amount" :class="item.type">{{ item.type === 'income' ? '+' : '−' }}{{ item.amount }}</text></view>
      <view v-if="!recent.length" class="empty-recent">还没有账单记录，导入账单后会显示在这里</view>
    </view>

    <view class="bottom-space" />
    <LedgerTabBar current="ledger" />
  </view>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { appendRecords, fetchCloudRecords, getAccount, getCloudSyncEnabled, getRecords, type LedgerRecord } from '../../utils/ledger'
import LedgerTabBar from '../../components/LedgerTabBar.vue'

const records = ref<LedgerRecord[]>([])
onShow(async () => {
  records.value = getRecords()
  if (getCloudSyncEnabled() && getAccount()?.token) {
    try {
      const cloud = await fetchCloudRecords()
      appendRecords(cloud.records)
      records.value = getRecords()
    } catch (error) {
      console.warn('Could not refresh ledger records from cloud:', error)
    }
  }
})
const selectedMonth = ref(monthKey(new Date()))
const monthLabel = computed(() => {
  const [year, month] = selectedMonth.value.split('-').map(Number)
  return `${year}年${month}月`
})
function monthKey(date: Date) { return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}` }
function onMonthChange(event: { detail: { value: string } }) { selectedMonth.value = event.detail.value.slice(0, 7) }
const currentMonthRecords = computed(() => {
  return records.value.filter(item => item.date.slice(0, 7) === selectedMonth.value)
})
const expense = computed(() => currentMonthRecords.value.filter(item => item.type === 'expense').reduce((sum, item) => sum + item.amount, 0))
const income = computed(() => currentMonthRecords.value.filter(item => item.type === 'income').reduce((sum, item) => sum + item.amount, 0))
const balance = computed(() => income.value - expense.value)
const categoryColors = ['#fff0e6', '#e9f4ff', '#f1edff', '#edf5ee', '#fff4d9', '#fce8ed']
const categories = computed(() => {
  const totals = new Map<string, number>()
  currentMonthRecords.value.filter(item => item.type === 'expense').forEach(item => totals.set(item.category, (totals.get(item.category) || 0) + item.amount))
  return [...totals.entries()].sort((a, b) => b[1] - a[1]).map(([name, amount], index) => ({ name, icon: name.slice(0, 1) || '其', bg: categoryColors[index % categoryColors.length], percent: expense.value ? Math.round(amount / expense.value * 100) : 0, share: expense.value ? amount / expense.value : 0 }))
})
const categoryGradient = computed(() => {
  if (!categories.value.length) return '#e1ece7'
  let stop = 0
  return `conic-gradient(${categories.value.map((item, index) => {
    const start = stop
    stop = index === categories.value.length - 1 ? 100 : stop + item.share * 100
    return `${item.bg} ${start}% ${stop}%`
  }).join(', ')})`
})
const recent = computed(() => [...records.value].sort((a, b) => b.date.localeCompare(a.date)).slice(0, 3).map(item => ({ name: item.merchant, category: item.category, time: item.date.slice(11, 16) || item.date.slice(0, 10), amount: item.amount.toFixed(2), type: item.type, icon: item.category.slice(0, 1) || '其', bg: categoryColors[0] })))
function formatMoney(value: number) { return value.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) }
function compactMoney(value: number) { return value > 10000 ? `${(value / 10000).toFixed(1)}万` : value.toFixed(0) }
function openRecords() { uni.redirectTo({ url: '/pages/records/index' }) }
function openImport() { uni.redirectTo({ url: '/pages/import/index' }) }
function openAccount() { uni.redirectTo({ url: '/pages/account/index' }) }
</script>

<style scoped>
.screen{min-height:100vh;padding:56rpx 36rpx 0;background:#f7f8fa}.topbar{display:flex;align-items:center;justify-content:space-between;margin-bottom:34rpx}.eyebrow{font-size:18rpx;letter-spacing:3rpx;color:#9ba0aa}.title-row{display:flex;align-items:center;gap:20rpx;margin-top:6rpx}.title{font-size:48rpx;font-weight:750;letter-spacing:1rpx}.month{padding:10rpx 18rpx;border-radius:24rpx;background:#fff;color:#707783;font-size:22rpx}.avatar{width:72rpx;height:72rpx;border-radius:50%;display:flex;align-items:center;justify-content:center;background:#dceee8;color:#26765a;font-size:28rpx;font-weight:700}.summary-card{height:310rpx;padding:32rpx 34rpx;border-radius:28rpx;background:linear-gradient(135deg,#183f38,#235d4c 65%,#31735e);color:white;position:relative;overflow:hidden;box-shadow:0 16rpx 32rpx #24564524}.summary-head{display:flex;justify-content:space-between;color:#c5d8d1;font-size:24rpx}.eye{font-size:26rpx}.balance{font-size:54rpx;font-weight:700;letter-spacing:1rpx;margin-top:12rpx}.summary-grid{display:flex;gap:72rpx;margin-top:28rpx;position:relative;z-index:1}.summary-grid>view{display:flex;flex-direction:column;gap:8rpx}.sum-label{font-size:21rpx;color:#c4d9d1}.dot{display:inline-block;width:12rpx;height:12rpx;margin-right:9rpx;border-radius:50%}.income-dot{background:#9de3bf}.expense-dot{background:#ffd28c}.sum-value{font-size:26rpx;font-weight:600}.card-glow{position:absolute;width:360rpx;height:360rpx;border-radius:50%;right:-120rpx;top:120rpx;background:#7bc1a21c}.section-head{display:flex;align-items:end;justify-content:space-between;margin:42rpx 2rpx 22rpx}.section-head>view{display:flex;flex-direction:column;gap:7rpx}.section-title{font-size:30rpx;font-weight:700}.section-sub{font-size:21rpx;color:#a1a6ae}.more{font-size:22rpx;color:#719487}.analysis-card,.bill-card{border-radius:24rpx;background:#fff;padding:26rpx 28rpx;box-shadow:0 8rpx 24rpx #1e302a08}.analysis-top{display:flex;justify-content:space-between;font-size:23rpx;font-weight:600}.muted{font-weight:400;color:#a3a8af}.chart-row{display:flex;align-items:center;gap:38rpx;margin-top:24rpx}.ring{width:204rpx;height:204rpx;flex:none;border-radius:50%;display:flex;align-items:center;justify-content:center;background:conic-gradient(#79a99a 0 42%,#a5c7bb 42% 70%,#c5ddd4 70% 88%,#e1ece7 88% 100%)}.ring-inner{width:150rpx;height:150rpx;border-radius:50%;background:white;display:flex;flex-direction:column;align-items:center;justify-content:center;gap:7rpx}.ring-label{font-size:20rpx;color:#9ba29e}.ring-value{font-size:29rpx;font-weight:700}.categories{flex:1;display:flex;flex-direction:column;gap:16rpx}.category-row,.cat-left{display:flex;align-items:center}.category-row{justify-content:space-between;font-size:21rpx;color:#666e70}.cat-left{gap:12rpx}.cat-icon{width:38rpx;height:38rpx;border-radius:12rpx;display:flex;align-items:center;justify-content:center;font-size:18rpx;color:#4a655c}.cat-amount{font-weight:600;color:#59615e}.recent-head{margin-top:40rpx}.bill-card{padding:4rpx 26rpx}.bill-row{height:112rpx;display:flex;align-items:center;border-bottom:1rpx solid #f2f3f3}.bill-row:last-child{border:0}.bill-icon{width:66rpx;height:66rpx;border-radius:20rpx;display:flex;align-items:center;justify-content:center;color:#986a41;font-weight:700}.bill-main{margin-left:18rpx;display:flex;flex-direction:column;gap:8rpx;flex:1}.bill-name{font-size:24rpx;font-weight:600}.bill-meta{font-size:19rpx;color:#a3a8ad}.bill-amount{font-size:25rpx;font-weight:650}.expense{color:#252a29}.income{color:#3b9b72}.bottom-space{height:150rpx}.tabbar{position:fixed;z-index:5;left:0;right:0;bottom:0;height:112rpx;padding-bottom:env(safe-area-inset-bottom);background:#ffffffef;box-shadow:0 -8rpx 24rpx #233b3010;display:flex;align-items:center;justify-content:space-around}.tab{width:100rpx;display:flex;flex-direction:column;align-items:center;gap:7rpx;font-size:18rpx;color:#a2a8a8}.tab-icon{font-size:34rpx;line-height:34rpx}.tab.active{color:#347662}.add-button{width:76rpx;height:76rpx;border-radius:24rpx;background:#245a49;color:white;font-size:48rpx;display:flex;align-items:center;justify-content:center;transform:translateY(-18rpx);box-shadow:0 8rpx 18rpx #245a4933}
</style>
