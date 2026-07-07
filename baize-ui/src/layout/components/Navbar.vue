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
import { DownOutlined } from '@ant-design/icons-vue'
import { Modal } from 'ant-design-vue'
import { computed } from 'vue'
import { useStore } from 'vuex'
import Breadcrumb from '@/components/Breadcrumb/index.vue'
import TopNav from '@/components/TopNav/index.vue'
import Hamburger from '@/components/Hamburger/index.vue'
import Screenfull from '@/components/Screenfull/index.vue'
import SizeSelect from '@/components/SizeSelect/index.vue'
import HeaderSearch from '@/components/HeaderSearch/index.vue'

const store = useStore()
const getters = computed(() => store.getters)
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

function logout() {
  Modal.confirm({
    title: '提示',
    content: '确定注销并退出系统吗？',
    okText: '确定',
    cancelText: '取消',
    onOk: () =>
      store.dispatch('LogOut').then(() => {
        location.href = '/index'
      })
  })
}

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

  .errLog-container {
    display: inline-block;
    vertical-align: top;
  }

  .right-menu {
    float: right;
    height: 100%;
    line-height: 50px;
    display: flex;

    &:focus {
      outline: none;
    }

    .right-menu-item {
      display: inline-flex;
      align-items: center;
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
</style>
