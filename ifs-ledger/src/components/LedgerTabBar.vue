<template>
  <view class="tabbar">
    <view class="tabbar-inner">
      <view v-for="item in leftItems" :key="item.key" class="tab" :class="{ active: current === item.key }" @tap="go(item)">
        <view class="icon-wrap"><image v-if="item.key === 'ledger'" class="tab-logo" src="/static/ledger-logo.svg" mode="aspectFit" /><wd-icon v-else :name="item.icon" class="tab-icon" size="21" /></view><text class="tab-label">{{ item.label }}</text>
      </view>
      <view class="add-action" @tap="addRecord"><view class="add-button"><wd-icon name="plus" size="25" color="#fff" /></view><text class="add-label">手动记账</text></view>
      <view v-for="item in rightItems" :key="item.key" class="tab" :class="{ active: current === item.key }" @tap="go(item)">
        <view class="icon-wrap"><image v-if="item.key === 'ledger'" class="tab-logo" src="/static/ledger-logo.svg" mode="aspectFit" /><wd-icon v-else :name="item.icon" class="tab-icon" size="21" /></view><text class="tab-label">{{ item.label }}</text>
      </view>
    </view>
  </view>
</template>
<script setup lang="ts">
type TabKey = 'ledger' | 'records' | 'import' | 'account'
type TabItem = { key: TabKey; label: string; icon: string; url: string }
defineProps<{ current: TabKey }>()
const leftItems: TabItem[] = [
  { key: 'ledger', label: '账本', icon: '', url: '/pages/index/index' },
]
const rightItems: TabItem[] = [
  { key: 'account', label: '我的', icon: 'user', url: '/pages/account/index' },
]
function go(item: TabItem) { uni.redirectTo({ url: item.url }) }
function addRecord() { uni.navigateTo({ url: '/pages/records/index?manual=1' }) }
</script>
<style scoped>
.tabbar{position:fixed;z-index:99;left:0;right:0;bottom:0;padding:0}.tabbar-inner{width:100%;height:calc(112rpx + env(safe-area-inset-bottom));box-sizing:border-box;padding:0 14rpx env(safe-area-inset-bottom);border:0;border-top:1rpx solid #ebe9f2;border-radius:30rpx 30rpx 0 0;background:rgba(255,255,255,.97);box-shadow:0 -10rpx 36rpx rgba(54,45,92,.12);display:grid;grid-template-columns:repeat(3,minmax(0,1fr));align-items:center;backdrop-filter:blur(18px)}.tab{min-width:0;display:flex;flex-direction:column;align-items:center;justify-content:center;gap:3rpx;color:#9591a6;transition:color .18s ease,transform .18s ease}.tab:active{transform:scale(.94)}.icon-wrap{width:56rpx;height:44rpx;display:flex;align-items:center;justify-content:center;border-radius:20rpx;transition:background .18s ease}.tab-icon{font-size:34rpx;line-height:38rpx;font-weight:600}.tab-logo{width:38rpx;height:38rpx;border-radius:12rpx}.tab-label{font-size:19rpx;line-height:26rpx;letter-spacing:1rpx}.tab.active{color:var(--theme-primary, #51458e);font-weight:650}.tab.active .icon-wrap{background:var(--theme-soft, #efedf8)}.add-action{justify-self:center;align-self:center;display:flex;flex-direction:column;align-items:center;gap:0;color:var(--theme-primary,#51458e)}.add-button{width:74rpx;height:74rpx;margin-top:-24rpx;border:6rpx solid #fff;border-radius:26rpx;background:linear-gradient(145deg,var(--theme-gradient-end, #7769c7),var(--theme-primary, #51458e));display:flex;align-items:center;justify-content:center;box-shadow:0 10rpx 22rpx rgba(71,59,130,.28)}.add-label{font-size:17rpx;line-height:24rpx;white-space:nowrap}.add-action:active .add-button{transform:scale(.94)}
</style>

