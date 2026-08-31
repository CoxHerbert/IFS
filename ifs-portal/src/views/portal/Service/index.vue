<template>
  <main class="page">
    <a-card class="hero" :bordered="false">
      <a-tag color="cyan">{{ copy.tag }}</a-tag>
      <h1>{{ copy.title }}</h1>
      <p>{{ copy.description }}</p>
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
        <h2>{{ copy.processTitle }}</h2>
        <p>{{ copy.processDescription }}</p>
      </div>

      <a-steps :current="1" direction="horizontal" responsive>
        <a-step v-for="step in steps" :key="step.title" :title="step.title" :description="step.description" />
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
import { computed, type Component } from 'vue'
import { usePortalI18n } from '@/i18n'

interface ServiceItem {
  title: string
  desc: string
  icon: Component
}

const { locale } = usePortalI18n()
const isEnglish = computed(() => locale.value === 'en')
const copy = computed(() => isEnglish.value ? {
  tag: 'FREIGHT SERVICES', title: 'International freight services built around your cargo', description: 'Compare transport modes, handling scope and the steps from enquiry to delivery.',
  processTitle: 'How shipping works', processDescription: 'A clear process keeps responsibilities, documents and milestones aligned.',
} : { tag: '服务能力', title: '把客户最关心的货代服务，整理成清晰可见的模块', description: '客户通常先确认服务能力与操作范围，再决定运输方案。', processTitle: '标准操作流程', processDescription: '清晰说明每个阶段，让客户更放心地安排出货。' })

const services = computed<ServiceItem[]>(() => isEnglish.value ? [
  { title: 'FCL Ocean Freight', desc: 'Container booking, trucking, export declaration, documentation and destination coordination.', icon: ContainerOutlined },
  { title: 'LCL Ocean Freight', desc: 'Flexible consolidation for smaller commercial shipments and cost-sensitive cargo.', icon: SwapOutlined },
  { title: 'Air Freight', desc: 'Fast handling for urgent cargo, samples and time-sensitive replenishment.', icon: CustomerServiceOutlined },
  { title: 'Warehousing & Distribution', desc: 'Door-to-door, overseas warehousing, transfer and final-mile options.', icon: GlobalOutlined },
  { title: 'Cargo Risk Support', desc: 'Packaging, insurance, compliance and operational risk reminders.', icon: SecurityScanOutlined },
  { title: 'Shipping Documents', desc: 'Support for quotations, packing lists, customs files and operating instructions.', icon: DownloadOutlined },
] : [
  { title: '海运整柜', desc: '整柜订舱、拖车、报关、单证和目的港协调。', icon: ContainerOutlined },
  { title: '海运拼箱', desc: '小批量货物灵活拼箱，适合成本敏感型客户。', icon: SwapOutlined },
  { title: '空运快件', desc: '高时效货物、样品和紧急补货快速处理。', icon: CustomerServiceOutlined },
  { title: '海外仓与转运', desc: '门到门、海外仓、转运和最后一公里方案。', icon: GlobalOutlined },
  { title: '货物安全', desc: '包装建议、保险、合规提醒和风险提示。', icon: SecurityScanOutlined },
  { title: '资料下载', desc: '报价表、装箱清单、报关资料和操作说明。', icon: DownloadOutlined }
])
const steps = computed(() => isEnglish.value ? [
  { title: 'Send Enquiry', description: 'Provide origin, destination and cargo details.' }, { title: 'Review Options', description: 'Compare route, price factors and estimated timing.' },
  { title: 'Confirm Booking', description: 'Confirm schedule, booking, trucking and documents.' }, { title: 'Ship & Deliver', description: 'Track milestones through arrival or final delivery.' },
] : [
  { title: '提交询价', description: '客户提交起运地、目的地和货物信息。' }, { title: '方案报价', description: '销售给出方案、价格因素和预计时效。' },
  { title: '确认订舱', description: '确认船期、订舱、拖车和单证。' }, { title: '出运交付', description: '跟踪货物，直到到港或签收。' },
])
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

.hero h1 {
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
