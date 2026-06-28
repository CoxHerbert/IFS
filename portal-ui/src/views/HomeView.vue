<template>
  <main class="page">
    <section class="hero">
      <div class="hero-copy" :style="heroStyle">
        <a-tag color="blue">货代销售官网 / Vue 3 / Ant Design Vue</a-tag>
        <h2>让客户 30 秒看懂你做哪条航线、能解决什么问题、怎么联系你</h2>
        <p>
          这是一套面向货代、物流、跨境运输销售场景的官网模板。
          首页直接服务获客，突出报价入口、航线能力、运输优势、案例信任和联系转化。
        </p>

        <div class="hero-actions">
          <router-link to="/contact">
            <a-button type="primary" size="large">立即询价</a-button>
          </router-link>
          <a-button size="large" class="ghost">查看航线</a-button>
        </div>

        <div class="stats">
          <a-card v-for="item in stats" :key="item.label" class="stat-card" :bordered="false">
            <div class="value">{{ item.value }}</div>
            <div class="label">{{ item.label }}</div>
          </a-card>
        </div>
      </div>

      <div class="hero-media">
        <img :src="heroImage" alt="freight forwarding" />
      </div>
    </section>

    <section class="section">
      <div class="section-head">
        <h3>核心服务</h3>
        <p>把客户最常问的能力放在最前面，减少沟通成本。</p>
      </div>

      <a-row :gutter="[18, 18]">
        <a-col v-for="item in serviceCards" :key="item.title" :xs="24" :md="12" :lg="6">
          <a-card class="feature" :bordered="false">
            <component :is="item.icon" class="icon" />
            <h4>{{ item.title }}</h4>
            <p>{{ item.desc }}</p>
          </a-card>
        </a-col>
      </a-row>
    </section>

    <section class="section alt">
      <a-row :gutter="[18, 18]">
        <a-col :xs="24" :lg="14">
          <a-card class="panel" :bordered="false">
            <div class="section-head">
              <h3>热门航线</h3>
              <p>把高频国家和港口列出来，客户一眼就知道你做哪里。</p>
            </div>

            <div class="route-grid">
              <div v-for="route in routes" :key="route.name" class="route-card">
                <div class="route-name">{{ route.name }}</div>
                <div class="route-meta">{{ route.meta }}</div>
                <div class="route-desc">{{ route.desc }}</div>
              </div>
            </div>
          </a-card>
        </a-col>

        <a-col :xs="24" :lg="10">
          <a-card class="panel quote-panel" :bordered="false">
            <div class="section-head">
              <h3>快速询价</h3>
              <p>客户越快留资，你的销售转化越高。</p>
            </div>

            <a-form layout="vertical" :model="quoteForm" @finish="handleQuoteSubmit">
              <a-form-item label="联系人">
                <a-input v-model:value="quoteForm.contactName" placeholder="请输入联系人姓名" />
              </a-form-item>
              <a-form-item label="起运地">
                <a-input v-model:value="quoteForm.origin" placeholder="例如：上海 / 深圳 / 宁波" />
              </a-form-item>
              <a-form-item label="目的地">
                <a-input v-model:value="quoteForm.destination" placeholder="例如：洛杉矶 / 汉堡 / 新加坡" />
              </a-form-item>
              <a-form-item label="货物信息">
                <a-textarea v-model:value="quoteForm.goods" :rows="3" placeholder="品名、体积、重量、是否带电 / 带磁 / 危化" />
              </a-form-item>
              <a-form-item label="联系方式">
                <a-input v-model:value="quoteForm.contact" placeholder="电话 / 微信 / 邮箱" />
              </a-form-item>
              <a-button type="primary" block size="large" html-type="submit" :loading="quoteSubmitting">
                提交询价
              </a-button>
            </a-form>
          </a-card>
        </a-col>
      </a-row>
    </section>

    <section class="section">
      <a-row :gutter="[18, 18]">
        <a-col :xs="24" :lg="12">
          <a-card class="panel" :bordered="false">
            <div class="section-head">
              <h3>客户最关心的事</h3>
              <p>用 FAQ 降低顾虑，销售沟通更轻。</p>
            </div>

            <a-collapse accordion>
              <a-collapse-panel v-for="item in faqs" :key="item.q" :header="item.q">
                <p>{{ item.a }}</p>
              </a-collapse-panel>
            </a-collapse>
          </a-card>
        </a-col>

        <a-col :xs="24" :lg="12">
          <a-card class="panel" :bordered="false">
            <div class="section-head">
              <h3>成交案例</h3>
              <p>用场景证明你能把货发走，而不只是介绍自己。</p>
            </div>

            <a-list :data-source="cases" item-layout="horizontal">
              <template #renderItem="{ item }">
                <a-list-item>
                  <a-list-item-meta :description="item.desc">
                    <template #title>
                      <span class="case-title">{{ item.title }}</span>
                    </template>
                  </a-list-item-meta>
                  <a-tag :color="item.tagColor">{{ item.tag }}</a-tag>
                </a-list-item>
              </template>
            </a-list>
          </a-card>
        </a-col>
      </a-row>
    </section>
  </main>
</template>

<script setup lang="ts">
import {
  ContainerOutlined,
  DashboardOutlined,
  GlobalOutlined,
  ThunderboltOutlined
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { reactive, ref, type Component } from 'vue'
import heroImage from '@/assets/hero.jpg'
import { submitContact } from '@/api/contact'

interface StatItem {
  value: string
  label: string
}

interface ServiceCard {
  title: string
  desc: string
  icon: Component
}

interface RouteCard {
  name: string
  meta: string
  desc: string
}

interface FaqItem {
  q: string
  a: string
}

interface CaseItem {
  title: string
  desc: string
  tag: string
  tagColor: string
}

interface QuoteForm {
  contactName: string
  origin: string
  destination: string
  goods: string
  contact: string
}

const heroStyle = {
  backgroundImage: `linear-gradient(135deg, rgba(7, 23, 47, 0.94), rgba(16, 103, 200, 0.82)), url(${heroImage})`
}

const stats: StatItem[] = [
  { value: '24h', label: '快速报价' },
  { value: '12+', label: '重点航线' },
  { value: '1000+', label: '月询盘承接' }
]

const serviceCards: ServiceCard[] = [
  { title: '海运整柜', desc: 'FCL 报价、订舱、拖车、报关和目的港协同。', icon: ContainerOutlined },
  { title: '海运拼箱', desc: '适合小批量货物，按立方或重量灵活组合。', icon: DashboardOutlined },
  { title: '空运快件', desc: '适合高时效订单、样品、急货和电商补货。', icon: ThunderboltOutlined },
  { title: '跨境物流', desc: '支持海外仓、双清、门到门和多段运输方案。', icon: GlobalOutlined }
]

const routes: RouteCard[] = [
  { name: '中国 - 美国', meta: '西海岸 / 东海岸 / FBA', desc: '适合普货、带电产品和电商补货。' },
  { name: '中国 - 欧洲', meta: '汉堡 / 鹿特丹 / 安特卫普', desc: '支持整柜、拼箱和铁路联运。' },
  { name: '中国 - 东南亚', meta: '新加坡 / 曼谷 / 胡志明', desc: '适合中小企业高频出货。' },
  { name: '中国 - 中东', meta: '迪拜 / 杰贝阿里', desc: '可承接贸易和项目型货物。' }
]

const faqs: FaqItem[] = [
  {
    q: '你们支持哪些货物类型？',
    a: '可展示普货、带电、危险品、超尺寸、样品和跨境电商等能力，实际业务可按你自己的服务范围替换。'
  },
  {
    q: '报价一般多久能回复？',
    a: '官网建议明确承诺，比如“工作时间 30 分钟内响应”，这类话术对销售转化很重要。'
  },
  {
    q: '客户怎样最快联系到你？',
    a: '把电话、微信、表单和 WhatsApp/邮件入口同时放在首屏和页尾，降低流失。'
  }
]

const cases: CaseItem[] = [
  { title: '深圳到洛杉矶整柜项目', desc: '从订舱、拖车到目的港交付，全流程对接。', tag: 'FCL', tagColor: 'green' },
  { title: '宁波到汉堡拼箱方案', desc: '按货量动态拼柜，帮助客户控制成本。', tag: 'LCL', tagColor: 'blue' },
  { title: '上海到新加坡空运急货', desc: '优先提货、优先排舱，控制时效。', tag: 'AIR', tagColor: 'gold' }
]

const quoteForm = reactive<QuoteForm>({
  contactName: '',
  origin: '',
  destination: '',
  goods: '',
  contact: ''
})

const quoteSubmitting = ref(false)

async function handleQuoteSubmit() {
  if (!quoteForm.contactName || !quoteForm.contact || !quoteForm.goods) {
    message.warning('请填写联系人、联系方式和货物信息')
    return
  }

  quoteSubmitting.value = true
  try {
    const result = await submitContact({
      contactName: quoteForm.contactName,
      phone: quoteForm.contact,
      route: [quoteForm.origin, quoteForm.destination].filter(Boolean).join(' - '),
      cargoInfo: quoteForm.goods,
      message: quoteForm.goods,
      source: 'portal-home-quote'
    })
    if (result.code !== 200) {
      throw new Error(result.msg || '提交失败')
    }
    message.success(`提交成功，线索编号：${result.data?.leadNo || '已生成'}`)
    Object.assign(quoteForm, {
      contactName: '',
      origin: '',
      destination: '',
      goods: '',
      contact: ''
    })
  } catch (error) {
    message.error(error instanceof Error ? error.message : '提交失败，请稍后再试')
  } finally {
    quoteSubmitting.value = false
  }
}
</script>

<style scoped>
.page {
  width: min(1240px, calc(100% - 32px));
  margin: 0 auto;
  padding: 28px 0 40px;
}

.hero {
  display: grid;
  grid-template-columns: 1.08fr 0.92fr;
  gap: 22px;
  align-items: stretch;
}

.hero-copy,
.hero-media,
.panel,
.feature {
  border-radius: 22px;
  box-shadow: 0 18px 40px rgba(16, 35, 63, 0.08);
}

.hero-copy {
  padding: 48px;
  color: #fff;
  background-position: center;
  background-size: cover;
  background-repeat: no-repeat;
}

.hero-copy h2 {
  margin: 18px 0 0;
  font-size: clamp(36px, 4.8vw, 62px);
  line-height: 1.06;
}

.hero-copy p {
  max-width: 58ch;
  margin: 18px 0 0;
  color: rgba(255, 255, 255, 0.86);
  line-height: 1.8;
}

.hero-actions {
  display: flex;
  gap: 14px;
  flex-wrap: wrap;
  margin-top: 26px;
}

.ghost {
  color: #fff;
  background: rgba(255, 255, 255, 0.12);
  border-color: rgba(255, 255, 255, 0.18);
}

.stats {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  margin-top: 32px;
}

.stat-card {
  background: rgba(255, 255, 255, 0.1);
}

.value {
  font-size: 28px;
  font-weight: 800;
}

.label {
  margin-top: 6px;
  color: rgba(255, 255, 255, 0.72);
}

.hero-media {
  overflow: hidden;
  background: #fff;
}

.hero-media img {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.section {
  padding-top: 28px;
}

.section-head {
  margin-bottom: 16px;
}

.section-head h3 {
  margin: 0;
  font-size: 24px;
}

.section-head p {
  margin: 8px 0 0;
  color: #66748b;
}

.feature {
  height: 100%;
}

.feature :deep(.ant-card-body) {
  min-height: 190px;
}

.icon {
  font-size: 26px;
  color: #1677ff;
}

.feature h4 {
  margin: 18px 0 0;
  font-size: 18px;
}

.feature p,
.route-desc,
.quote-panel :deep(.ant-form-item-label > label),
.panel :deep(.ant-collapse-content-box),
.case-title {
  color: #66748b;
}

.feature p {
  margin: 10px 0 0;
  line-height: 1.7;
}

.alt {
  margin-top: 6px;
}

.panel {
  min-height: 100%;
}

.panel :deep(.ant-card-body) {
  padding: 24px;
}

.route-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px;
}

.route-card {
  padding: 18px;
  border-radius: 18px;
  background: #f6f9ff;
}

.route-name {
  font-size: 18px;
  font-weight: 700;
}

.route-meta {
  margin-top: 8px;
  color: #1677ff;
  font-size: 13px;
}

.route-desc {
  margin-top: 8px;
  line-height: 1.7;
}

.quote-panel :deep(.ant-form-item) {
  margin-bottom: 14px;
}

@media (max-width: 1120px) {
  .hero {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 760px) {
  .page {
    width: min(100% - 20px, 1240px);
  }

  .hero-copy {
    padding: 28px;
  }

  .stats,
  .route-grid {
    grid-template-columns: 1fr;
  }
}
</style>
