<template>
  <main class="page">
    <a-card class="hero" :bordered="false">
      <a-tag color="cyan">服务中心</a-tag>
      <h2>把常用系统和服务入口集中起来</h2>
      <p>适合做办事入口、下载中心、在线客服、工单系统和常用工具聚合页。</p>
    </a-card>

    <a-row :gutter="[18, 18]" class="section">
      <a-col v-for="item in services" :key="item.title" :xs="24" :md="12" :lg="8">
        <a-card class="panel" :bordered="false">
          <component :is="item.icon" class="icon" />
          <h3>{{ item.title }}</h3>
          <p>{{ item.desc }}</p>
        </a-card>
      </a-col>
    </a-row>
  </main>
</template>

<script setup>
import { AppstoreOutlined, CustomerServiceOutlined, DownloadOutlined } from '@ant-design/icons-vue'

const services = [
  { title: '在线办理', desc: '把用户最常用的申请和审批入口前置。', icon: CustomerServiceOutlined },
  { title: '资料下载', desc: '常见表单、说明书和附件集中管理。', icon: DownloadOutlined },
  { title: '系统导航', desc: '对接内外部系统，减少重复跳转。', icon: AppstoreOutlined }
]
</script>

<style scoped>
.page {
  width: min(1240px, calc(100% - 32px));
  margin: 0 auto;
  padding: 28px 0 40px;
}

.hero,
.panel {
  border-radius: 22px;
  box-shadow: 0 18px 40px rgba(16, 35, 63, 0.08);
}

.hero {
  padding: 28px;
  background: linear-gradient(135deg, #071a33, #0f65c3);
  color: #fff;
}

.hero h2 {
  margin: 14px 0 0;
  font-size: clamp(28px, 3vw, 44px);
}

.hero p {
  max-width: 56ch;
  margin: 14px 0 0;
  color: rgba(255, 255, 255, 0.82);
  line-height: 1.8;
}

.section {
  margin-top: 6px;
}

.panel {
  height: 100%;
}

.panel :deep(.ant-card-body) {
  min-height: 210px;
}

.icon {
  font-size: 26px;
  color: #1677ff;
}

.panel h3 {
  margin: 18px 0 0;
  font-size: 20px;
}

.panel p {
  margin: 10px 0 0;
  color: #66748b;
  line-height: 1.7;
}

@media (max-width: 760px) {
  .page {
    width: min(100% - 20px, 1240px);
  }
}
</style>
