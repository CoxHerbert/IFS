<template>
  <main class="page">
    <a-card class="hero" :bordered="false">
      <a-tag color="cyan">服务能力</a-tag>
      <h2>把客户最关心的货代服务，整理成清晰可见的模块</h2>
      <p>客户通常不会研究你公司介绍，而是先看你能不能办事。这个页面适合放服务能力和操作流程。</p>
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

    <a-card class="panel section" :bordered="false">
      <div class="section-head">
        <h3>标准操作流程</h3>
        <p>把流程写清楚，客户会更敢下单。</p>
      </div>

      <a-steps :current="1" direction="horizontal" responsive>
        <a-step title="提交询价" description="客户提交起运地、目的地和货物信息。" />
        <a-step title="快速报价" description="销售给出方案、价格和时效。" />
        <a-step title="确认订舱" description="确认船期、订舱、拖车和单证。" />
        <a-step title="出运交付" description="跟踪货物，直到到港或签收。" />
      </a-steps>
    </a-card>
  </main>
</template>

<script setup lang="ts">
import {
  ContainerOutlined,
  CustomerServiceOutlined,
  DownloadOutlined,
  GlobalOutlined,
  SecurityScanOutlined,
  SwapOutlined
} from '@ant-design/icons-vue'
import type { Component } from 'vue'

interface ServiceItem {
  title: string
  desc: string
  icon: Component
}

const services: ServiceItem[] = [
  { title: '海运整柜', desc: '整柜订舱、拖车、报关、单证和目的港协调。', icon: ContainerOutlined },
  { title: '海运拼箱', desc: '小批量货物灵活拼箱，适合成本敏感型客户。', icon: SwapOutlined },
  { title: '空运快件', desc: '高时效货物、样品和紧急补货快速处理。', icon: CustomerServiceOutlined },
  { title: '海外仓与转运', desc: '门到门、海外仓、转运和最后一公里方案。', icon: GlobalOutlined },
  { title: '货物安全', desc: '包装建议、保险、合规提醒和风险提示。', icon: SecurityScanOutlined },
  { title: '资料下载', desc: '报价表、装箱清单、报关资料和操作说明。', icon: DownloadOutlined }
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
  margin-top: 18px;
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

.section-head {
  margin-bottom: 20px;
}

.section-head h3 {
  margin: 0;
  font-size: 22px;
}

.section-head p {
  margin: 8px 0 0;
  color: #66748b;
}

@media (max-width: 760px) {
  .page {
    width: min(100% - 20px, 1240px);
  }
}
</style>
