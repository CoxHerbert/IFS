<template>
  <main class="page">
    <section class="hero">
      <div class="hero-copy" :style="heroStyle">
        <a-tag color="blue">{{ copy.tag }}</a-tag>
        <h1>{{ copy.title }}</h1>
        <p>{{ copy.description }}</p>

        <div class="hero-actions">
          <router-link :to="localePath('/contact')">
            <a-button type="primary" size="large">{{ copy.quote }}</a-button>
          </router-link>
        </div>

        <div class="stats">
          <div v-for="item in stats" :key="item.label" class="stat-card">
            <strong>{{ item.value }}</strong>
            <span>{{ item.label }}</span>
          </div>
        </div>
      </div>

      <div class="hero-side">
        <a-card class="side-card" :bordered="false">
          <h3>{{ copy.routesTitle }}</h3>
          <div class="route-list">
            <div v-for="item in routes" :key="item.name" class="route-item">
              <strong>{{ item.name }}</strong>
              <span>{{ item.meta }}</span>
            </div>
          </div>
        </a-card>

        <a-card class="side-card compact" :bordered="false">
          <h3>{{ copy.scenariosTitle }}</h3>
          <div class="tag-list">
            <a-tag v-for="item in scenarios" :key="item" color="blue">{{ item }}</a-tag>
          </div>
        </a-card>
      </div>
    </section>

    <section class="section">
      <div class="section-head">
        <h2>{{ copy.servicesTitle }}</h2>
        <p>{{ copy.servicesDescription }}</p>
      </div>

      <div class="service-grid">
        <a-card v-for="item in services" :key="item.title" class="service-card" :bordered="false">
          <component :is="item.icon" class="icon" />
          <h3>{{ item.title }}</h3>
          <p>{{ item.desc }}</p>
        </a-card>
      </div>
    </section>

  </main>
</template>

<script setup lang="ts">
import {
  ContainerOutlined,
  GlobalOutlined,
  InboxOutlined,
  ThunderboltOutlined
} from '@ant-design/icons-vue'
import { computed, type Component } from 'vue'
import heroImage from '@/assets/hero.jpg'
import { usePortalI18n } from '@/i18n'

interface StatItem {
  value: string
  label: string
}

interface RouteItem {
  name: string
  meta: string
}

interface ServiceItem {
  title: string
  desc: string
  icon: Component
}

const heroStyle = {
  backgroundImage: `linear-gradient(135deg, rgba(7, 23, 47, 0.94), rgba(16, 103, 200, 0.82)), url(${heroImage})`
}

const { locale, localePath } = usePortalI18n()
const isEnglish = computed(() => locale.value === 'en')
const copy = computed(() => isEnglish.value ? {
  tag: 'INTERNATIONAL FREIGHT FORWARDING', title: 'Move cargo from China with a clear, workable shipping plan.',
  description: 'Compare core routes and freight services, then send your cargo details for a tailored quotation.', quote: 'Get a Freight Quote',
  routesTitle: 'Key Routes', scenariosTitle: 'Common Shipments', servicesTitle: 'Core Services', servicesDescription: 'The freight capabilities customers ask about most.',
} : {
  tag: '国际货代官网', title: '让客户快速知道你能运什么、发到哪里、怎么联系你。',
  description: '首页只保留最重要的信息：核心航线、服务能力和询价入口。', quote: '立即询价',
  routesTitle: '重点航线', scenariosTitle: '常见场景', servicesTitle: '核心服务', servicesDescription: '客户最常问的能力，直接放在首页说明。',
})

const stats = computed<StatItem[]>(() => isEnglish.value ? [
  { value: '24h', label: 'Prompt response' }, { value: '12+', label: 'Key routes' }, { value: 'FCL/LCL', label: 'Flexible options' },
] : [
  { value: '24h', label: '快速响应' }, { value: '12+', label: '重点航线' }, { value: 'FCL/LCL', label: '灵活运输方案' },
])

const routes = computed<RouteItem[]>(() => isEnglish.value ? [
  { name: 'China → United States', meta: 'West Coast / East Coast / FBA' }, { name: 'China → Europe', meta: 'Hamburg / Rotterdam / Antwerp' },
  { name: 'China → Southeast Asia', meta: 'Singapore / Bangkok / Ho Chi Minh City' }, { name: 'China → Middle East', meta: 'Dubai / Jebel Ali' },
] : [
  { name: '中国 -> 美国', meta: '美西 / 美东 / FBA' },
  { name: '中国 -> 欧洲', meta: '汉堡 / 鹿特丹 / 安特卫普' },
  { name: '中国 -> 东南亚', meta: '新加坡 / 曼谷 / 胡志明' },
  { name: '中国 -> 中东', meta: '迪拜 / 杰贝阿里' }
])

const scenarios = computed(() => isEnglish.value ? ['FCL', 'LCL', 'Urgent air cargo', 'E-commerce', 'Door to door', 'Battery cargo'] : ['整柜 FCL', '拼箱 LCL', '空运急货', '跨境电商', '门到门', '带电产品'])

const services = computed<ServiceItem[]>(() => isEnglish.value ? [
  { title: 'FCL Ocean Freight', desc: 'Quotation, booking, trucking, export declaration and destination coordination.', icon: ContainerOutlined },
  { title: 'LCL Ocean Freight', desc: 'Flexible consolidation for cargo that does not require a full container.', icon: InboxOutlined },
  { title: 'Air Freight', desc: 'Faster transport for urgent orders, samples and time-sensitive replenishment.', icon: ThunderboltOutlined },
  { title: 'Cross-border Logistics', desc: 'Warehousing, customs coordination and door-to-door delivery options.', icon: GlobalOutlined },
] : [
  { title: '海运整柜', desc: '报价、订舱、拖车、报关和目的港协同。', icon: ContainerOutlined },
  { title: '海运拼箱', desc: '适合小批量货物，按体积或重量灵活组合。', icon: InboxOutlined },
  { title: '空运快件', desc: '适合高时效订单、样品和紧急补货。', icon: ThunderboltOutlined },
  { title: '跨境物流', desc: '支持海外仓、双清和门到门方案。', icon: GlobalOutlined }
])

</script>

<style scoped>
.page {
  width: min(1180px, calc(100% - 32px));
  margin: 0 auto;
  padding: 28px 0 40px;
}

.hero {
  display: grid;
  grid-template-columns: minmax(0, 1.15fr) minmax(320px, 0.85fr);
  gap: 20px;
  align-items: stretch;
}

.hero-copy,
.side-card,
.service-card {
  border-radius: 20px;
  box-shadow: 0 18px 40px rgba(16, 35, 63, 0.08);
}

.hero-copy {
  padding: 42px;
  color: #fff;
  background-position: center;
  background-size: cover;
  background-repeat: no-repeat;
}

.hero-copy h1 {
  margin: 18px 0 0;
  font-size: clamp(34px, 4.6vw, 58px);
  line-height: 1.08;
}

.hero-copy p {
  max-width: 46ch;
  margin: 16px 0 0;
  color: rgba(255, 255, 255, 0.86);
  line-height: 1.8;
}

.hero-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 14px;
  margin-top: 24px;
}

.stats {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  margin-top: 28px;
}

.stat-card {
  border-radius: 16px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.1);
}

.stat-card strong,
.stat-card span {
  display: block;
}

.stat-card strong {
  font-size: 26px;
}

.stat-card span {
  margin-top: 6px;
  color: rgba(255, 255, 255, 0.72);
}

.hero-side {
  display: grid;
  gap: 20px;
}

.side-card {
  background: #fff;
}

.side-card :deep(.ant-card-body),
.service-card :deep(.ant-card-body) {
  padding: 24px;
}

.side-card h3,
.service-card h3,
.section-head h2 {
  margin: 0;
}

.route-list {
  display: grid;
  gap: 12px;
  margin-top: 16px;
}

.route-item {
  border-radius: 14px;
  border: 1px solid rgba(148, 163, 184, 0.16);
  background: #f8fafc;
  padding: 14px 16px;
}

.route-item strong,
.route-item span {
  display: block;
}

.route-item span,
.section-head p,
.service-card p {
  margin-top: 6px;
  color: #64748b;
}

.tag-list {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 14px;
}

.section {
  padding-top: 28px;
}

.section-head {
  margin-bottom: 16px;
}

.section-head h2 {
  font-size: 24px;
}

.service-grid {
  display: grid;
  gap: 18px;
}

.service-grid {
  grid-template-columns: repeat(4, minmax(0, 1fr));
}

.service-card {
  background: #fff;
}

.icon {
  font-size: 28px;
  color: #1677ff;
}

.service-card h3 {
  margin-top: 16px;
  font-size: 20px;
}

@media (max-width: 1100px) {
  .hero,
  .service-grid {
    grid-template-columns: 1fr;
  }

  .hero-copy {
    padding: 34px 28px;
  }
}

@media (max-width: 640px) {
  .page {
    width: min(100%, calc(100% - 20px));
    padding-top: 18px;
  }

  .stats {
    grid-template-columns: 1fr;
  }

  .hero-copy h1 {
    font-size: 34px;
  }

}
</style>
