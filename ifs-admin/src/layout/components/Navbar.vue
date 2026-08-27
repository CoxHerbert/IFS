<template>
  <div class="navbar">
    <hamburger id="hamburger-container" :is-active="getters.sidebar.opened" class="hamburger-container"
      @toggleClick="toggleSideBar" />
    <breadcrumb id="breadcrumb-container" v-if="!store.state.settings.topNav" class="breadcrumb-container" />
    <top-nav id="topmenu-container" v-if="store.state.settings.topNav" class="topmenu-container" />

    <div class="right-menu">
      <template v-if="getters.device !== 'mobile'">
        <header-search id="header-search" class="right-menu-item" />
        <screenfull id="screenfull" class="right-menu-item hover-effect" />
        <a-tooltip placement="bottom" title="布局大小">
          <size-select id="size-select" class="right-menu-item hover-effect" />
        </a-tooltip>
      </template>

      <a-popover trigger="click" placement="bottomRight" overlay-class-name="notification-popover" @openChange="handleNotificationOpen">
        <template #content>
          <div class="notification-panel">
            <div class="notification-header">
              <span>消息通知</span>
              <a-button type="link" size="small" :disabled="!unreadCount" @click="handleReadAll">全部已读</a-button>
            </div>
            <div v-if="notificationLoading" class="notification-state">加载中...</div>
            <div v-else-if="!notificationList.length" class="notification-state">暂无通知</div>
            <div v-else class="notification-list">
              <div
                v-for="item in notificationList"
                :key="item.notificationId"
                class="notification-item"
                :class="{ unread: item.readFlag === '0' }"
                @click="handleNotificationClick(item)"
              >
                <div class="notification-item-title">
                  <span>{{ item.title }}</span>
                  <span class="notification-item-time">{{ formatNotificationTime(item.createTime) }}</span>
                </div>
                <div class="notification-item-content">{{ item.content }}</div>
              </div>
            </div>
          </div>
        </template>
        <div class="right-menu-item hover-effect notification-trigger">
          <a-badge :count="unreadCount" :overflow-count="99" size="small">
            <bell-outlined />
          </a-badge>
        </div>
      </a-popover>

      <div class="avatar-container">
        <a-dropdown :trigger="['click']" class="right-menu-item hover-effect">
          <div class="avatar-wrapper">
            <img :src="getters.avatar" class="user-avatar" />
            <div class="avatar-meta" v-if="getters.device !== 'mobile'">
              <span class="user-name">{{ getters.name || '管理员' }}</span>
              <span class="user-role">{{ userRoleLabel }}</span>
            </div>
            <DownOutlined class="avatar-arrow" />
          </div>
          <template #overlay>
            <a-menu>
              <a-menu-item key="profile">
                <router-link to="/user/profile">个人中心</router-link>
              </a-menu-item>
              <a-menu-divider />
              <a-menu-item key="logout" @click="logout">
                <span>退出登录</span>
              </a-menu-item>
            </a-menu>
          </template>
        </a-dropdown>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { BellOutlined, DownOutlined } from '@ant-design/icons-vue'
import { Modal, message } from 'ant-design-vue'
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'
import Breadcrumb from '@/components/Breadcrumb/index.vue'
import TopNav from '@/components/TopNav/index.vue'
import Hamburger from '@/components/Hamburger/index.vue'
import Screenfull from '@/components/Screenfull/index.vue'
import SizeSelect from '@/components/SizeSelect/index.vue'
import HeaderSearch from '@/components/HeaderSearch/index.vue'
import { parseTime } from '@/utils/ruoyi'
import {
  getUnreadNotificationCount,
  listNotification,
  readAllNotification,
  readNotification
} from '@/api/system/notification'

const store = useStore()
const router = useRouter()
const getters = computed(() => store.getters)
const unreadCount = ref(0)
const notificationLoading = ref(false)
const notificationList = ref([])
let pollingTimer: number | undefined

const userRoleLabel = computed(() => {
  const roles = getters.value.roles || []
  if (!roles.length) {
    return '当前账号'
  }
  return roles.includes('admin') ? '系统管理员' : '已登录'
})

function toggleSideBar() {
  store.dispatch('app/toggleSideBar')
}

function fetchUnreadCount() {
  getUnreadNotificationCount().then(response => {
    unreadCount.value = response.data?.unreadCount || 0
  })
}

function formatNotificationTime(value) {
  return parseTime(value, '{m}-{d} {h}:{i}') || ''
}

function fetchNotificationList() {
  notificationLoading.value = true
  listNotification({ pageNum: 1, pageSize: 8 })
    .then(response => {
      const data = response.data || {}
      notificationList.value = data.rows || []
    })
    .finally(() => {
      notificationLoading.value = false
    })
}

function handleNotificationOpen(open: boolean) {
  if (open) {
    fetchNotificationList()
    fetchUnreadCount()
  }
}

function handleNotificationClick(item) {
  const wasUnread = item.readFlag === '0'
  const done = wasUnread ? readNotification(item.notificationId) : Promise.resolve()
  done.then(() => {
    item.readFlag = '1'
    if (wasUnread) {
      unreadCount.value = Math.max(0, unreadCount.value - 1)
    }
    if (item.bizType === 'shipment') {
      router.push('/freight/shipment')
    }
  })
}

function handleReadAll(event) {
  event?.stopPropagation?.()
  readAllNotification().then(() => {
    notificationList.value = notificationList.value.map(item => ({ ...item, readFlag: '1' }))
    unreadCount.value = 0
    message.success('已全部标记为已读')
  })
}

function logout() {
  Modal.confirm({
    title: '提示',
    content: '确认注销并退出系统吗？',
    okText: '确定',
    cancelText: '取消',
    onOk: () =>
      store.dispatch('LogOut').then(() => {
        location.href = '/index'
      })
  })
}

onMounted(() => {
  fetchUnreadCount()
  pollingTimer = window.setInterval(fetchUnreadCount, 60000)
})

onBeforeUnmount(() => {
  if (pollingTimer) {
    window.clearInterval(pollingTimer)
  }
})
</script>

<style lang="scss" scoped>
.navbar {
  height: 50px;
  overflow: hidden;
  position: relative;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);

  .hamburger-container {
    line-height: 46px;
    height: 100%;
    float: left;
    cursor: pointer;
    transition: background 0.3s;
    -webkit-tap-highlight-color: transparent;

    &:hover {
      background: rgba(0, 0, 0, 0.025);
    }
  }

  .breadcrumb-container {
    float: left;
  }

  .topmenu-container {
    position: absolute;
    left: 50px;
  }

  .right-menu {
    float: right;
    height: 100%;
    line-height: 50px;
    display: flex;
    align-items: center;

    &:focus {
      outline: none;
    }

    .right-menu-item {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      padding: 0 8px;
      height: 100%;
      font-size: 18px;
      color: #5a5e66;
      vertical-align: text-bottom;

      &.hover-effect {
        cursor: pointer;
        transition: background 0.3s;

        &:hover {
          background: rgba(0, 0, 0, 0.025);
        }
      }
    }

    .notification-trigger {
      width: 42px;
    }

    .avatar-container {
      margin-right: 20px;
      display: flex;
      align-items: center;

      .avatar-wrapper {
        display: flex;
        align-items: center;
        gap: 10px;
        min-width: 0;
        height: 50px;
        padding: 0 12px;

        .user-avatar {
          cursor: pointer;
          width: 34px;
          height: 34px;
          border-radius: 12px;
          object-fit: cover;
          box-shadow: 0 6px 16px rgba(15, 23, 42, 0.12);
        }

        .avatar-meta {
          display: flex;
          min-width: 0;
          flex-direction: column;
          justify-content: center;
          line-height: 1.2;
        }

        .user-name {
          max-width: 120px;
          overflow: hidden;
          color: #1f2937;
          font-size: 14px;
          font-weight: 600;
          text-overflow: ellipsis;
          white-space: nowrap;
        }

        .user-role {
          color: #8c98a8;
          font-size: 12px;
        }

        .avatar-arrow {
          font-size: 12px;
          color: #8c98a8;
        }
      }
    }
  }
}

.notification-panel {
  width: 360px;
}

.notification-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
  font-weight: 600;
}

.notification-state {
  padding: 20px 0;
  color: #8c98a8;
  text-align: center;
}

.notification-list {
  max-height: 360px;
  overflow-y: auto;
}

.notification-item {
  padding: 10px 0;
  border-top: 1px solid #f0f0f0;
  cursor: pointer;
}

.notification-item.unread .notification-item-title span:first-child {
  color: #1677ff;
}

.notification-item-title {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 4px;
  font-size: 13px;
  font-weight: 600;
}

.notification-item-time {
  flex-shrink: 0;
  color: #8c98a8;
  font-size: 12px;
  font-weight: 400;
}

.notification-item-content {
  color: #5a5e66;
  line-height: 1.5;
}
</style>
