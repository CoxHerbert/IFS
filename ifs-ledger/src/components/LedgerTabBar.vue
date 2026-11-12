<template>
  <view class="tabbar">
    <view class="tabbar-inner">
      <view v-for="item in leftItems" :key="item.key" class="tab" :class="{ active: current === item.key }" @tap="go(item)">
        <view class="icon-wrap"><text class="tab-icon">{{ item.icon }}</text></view><text class="tab-label">{{ item.label }}</text>
      </view>
      <view class="add-button" @tap="addRecord"><text>＋</text></view>
      <view v-for="item in rightItems" :key="item.key" class="tab" :class="{ active: current === item.key }" @tap="go(item)">
        <view class="icon-wrap"><text class="tab-icon">{{ item.icon }}</text></view><text class="tab-label">{{ item.label }}</text>
      </view>
    </view>
  </view>
</template>
<script setup lang="ts">
type TabKey = 'ledger' | 'records' | 'import' | 'account'
type TabItem = { key: TabKey; label: string; icon: string; url: string }
defineProps<{ current: TabKey }>()
const leftItems: TabItem[] = [
  { key: 'ledger', label: '账本', icon: '◈', url: '/pages/index/index' },
  { key: 'records', label: '明细', icon: '▤', url: '/pages/records/index' },
]
const rightItems: TabItem[] = [
  { key: 'import', label: '导入', icon: '⇧', url: '/pages/import/index' },
  { key: 'account', label: '我的', icon: '○', url: '/pages/account/index' },
]
function go(item: TabItem) { uni.redirectTo({ url: item.url }) }
function addRecord() { uni.showToast({ title: '手动记账即将开放', icon: 'none' }) }
</script>
<style scoped>
.tabbar{position:fixed;z-index:99;left:0;right:0;bottom:0;padding:0}.tabbar-inner{width:100%;height:calc(112rpx + env(safe-area-inset-bottom));box-sizing:border-box;padding:0 14rpx env(safe-area-inset-bottom);border:0;border-top:1rpx solid #edf1ee;border-radius:30rpx 30rpx 0 0;background:rgba(255,255,255,.97);box-shadow:0 -10rpx 36rpx rgba(27,55,43,.11);display:grid;grid-template-columns:repeat(5,minmax(0,1fr));align-items:center;backdrop-filter:blur(18px)}.tab{min-width:0;display:flex;flex-direction:column;align-items:center;justify-content:center;gap:3rpx;color:#9ba49f;transition:color .18s ease,transform .18s ease}.tab:active{transform:scale(.94)}.icon-wrap{width:56rpx;height:44rpx;display:flex;align-items:center;justify-content:center;border-radius:20rpx;transition:background .18s ease}.tab-icon{font-size:34rpx;line-height:38rpx;font-weight:600}.tab-label{font-size:19rpx;line-height:26rpx;letter-spacing:1rpx}.tab.active{color:#24654e;font-weight:650}.tab.active .icon-wrap{background:#eaf4ef}.add-button{justify-self:center;width:82rpx;height:82rpx;margin-top:-28rpx;border:6rpx solid #fff;border-radius:29rpx;background:linear-gradient(145deg,#347b5d,#1d5643);color:#fff;font-size:44rpx;display:flex;align-items:center;justify-content:center;box-shadow:0 10rpx 22rpx rgba(36,90,73,.3)}.add-button:active{transform:scale(.94)}
</style>
