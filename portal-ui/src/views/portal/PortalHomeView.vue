<template>
  <main class="page">
    <section class="hero">
      <div class="hero-copy" :style="heroStyle">
        <a-tag color="blue">国际货代官网</a-tag>
        <h1>让客户快速知道你能运什么、发到哪里、怎么联系你。</h1>
        <p>
          首页只保留最重要的信息：核心航线、常用工具、服务能力和询价入口。
        </p>

        <div class="hero-actions">
          <router-link to="/contact">
            <a-button type="primary" size="large">立即询价</a-button>
          </router-link>
          <a-button size="large" class="ghost" @click="handleToolAction('shipment-agent')">智能分析</a-button>
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
          <h3>重点航线</h3>
          <div class="route-list">
            <div v-for="item in routes" :key="item.name" class="route-item">
              <strong>{{ item.name }}</strong>
              <span>{{ item.meta }}</span>
            </div>
          </div>
        </a-card>

        <a-card class="side-card compact" :bordered="false">
          <h3>常见场景</h3>
          <div class="tag-list">
            <a-tag v-for="item in scenarios" :key="item" color="blue">{{ item }}</a-tag>
          </div>
        </a-card>
      </div>
    </section>

    <section class="section">
      <div class="section-head">
        <h2>快捷工具</h2>
        <p>先用工具判断方案，再进入人工报价。</p>
      </div>

      <div class="tool-grid">
        <a-card v-for="item in tools" :key="item.title" class="tool-card" :bordered="false">
          <component :is="item.icon" class="icon" />
          <h3>{{ item.title }}</h3>
          <p>{{ item.desc }}</p>
          <button type="button" class="tool-action" @click="handleToolAction(item.action)">{{ item.cta }}</button>
        </a-card>
      </div>
    </section>

    <section class="section">
      <div class="section-head">
        <h2>核心服务</h2>
        <p>客户最常问的能力，直接放在首页说明。</p>
      </div>

      <div class="service-grid">
        <a-card v-for="item in services" :key="item.title" class="service-card" :bordered="false">
          <component :is="item.icon" class="icon" />
          <h3>{{ item.title }}</h3>
          <p>{{ item.desc }}</p>
        </a-card>
      </div>
    </section>

    <section class="section final-section">
      <a-row :gutter="[18, 18]">
        <a-col :xs="24" :lg="10">
          <a-card class="panel contact-panel" :bordered="false">
            <div class="section-head">
              <h2>为什么现在联系</h2>
              <p>信息越完整，报价和方案返回越快。</p>
            </div>

            <div class="bullet-list">
              <div v-for="item in contactReasons" :key="item.title" class="bullet-item">
                <strong>{{ item.title }}</strong>
                <span>{{ item.desc }}</span>
              </div>
            </div>
          </a-card>
        </a-col>

        <a-col :xs="24" :lg="14">
          <a-card class="panel quote-panel" :bordered="false">
            <div class="section-head">
              <h2>快速询价</h2>
              <p>留下基础信息，销售跟进即可继续推进。</p>
            </div>

            <a-form layout="vertical" :model="quoteForm" @finish="handleQuoteSubmit">
              <a-row :gutter="14">
                <a-col :xs="24" :md="12">
                  <a-form-item label="联系人">
                    <a-input v-model:value="quoteForm.contactName" placeholder="请输入联系人" />
                  </a-form-item>
                </a-col>
                <a-col :xs="24" :md="12">
                  <a-form-item label="联系方式">
                    <a-input v-model:value="quoteForm.contact" placeholder="电话 / 微信 / 邮箱" />
                  </a-form-item>
                </a-col>
              </a-row>

              <a-row :gutter="14">
                <a-col :xs="24" :md="12">
                  <a-form-item label="起运地">
                    <a-input v-model:value="quoteForm.origin" placeholder="如：上海 / 深圳" />
                  </a-form-item>
                </a-col>
                <a-col :xs="24" :md="12">
                  <a-form-item label="目的地">
                    <a-input v-model:value="quoteForm.destination" placeholder="如：洛杉矶 / 汉堡" />
                  </a-form-item>
                </a-col>
              </a-row>

              <a-form-item label="货物信息">
                <a-textarea
                  v-model:value="quoteForm.goods"
                  :rows="4"
                  placeholder="填写品名、体积、重量、箱数、时效要求等"
                />
              </a-form-item>

              <div class="quote-actions">
                <a-button type="primary" size="large" html-type="submit" :loading="quoteSubmitting">
                  提交询价
                </a-button>
                <a-button size="large" @click="handleToolAction('assistant')">先问智能助手</a-button>
              </div>
            </a-form>
          </a-card>
        </a-col>
      </a-row>
    </section>
  </main>
</template>

<script setup lang="ts">
import {
  CalculatorOutlined,
  CommentOutlined,
  ContainerOutlined,
  GlobalOutlined,
  InboxOutlined,
  ThunderboltOutlined
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { reactive, ref, type Component } from 'vue'
import { useRouter } from 'vue-router'
import heroImage from '@/assets/hero.jpg'
import { submitContact } from '@/api/portal/contact'

interface StatItem {
  value: string
  label: string
}

interface RouteItem {
  name: string
  meta: string
}

interface ToolItem {
  title: string
  desc: string
  cta: string
  icon: Component
  action: ToolAction
}

interface ServiceItem {
  title: string
  desc: string
  icon: Component
}

interface ContactReason {
  title: string
  desc: string
}

interface QuoteForm {
  contactName: string
  origin: string
  destination: string
  goods: string
  contact: string
}

type ToolAction = 'shipment-agent' | 'contact' | 'assistant'

const router = useRouter()

const heroStyle = {
  backgroundImage: `linear-gradient(135deg, rgba(7, 23, 47, 0.94), rgba(16, 103, 200, 0.82)), url(${heroImage})`
}

const stats: StatItem[] = [
  { value: '24h', label: '快速响应' },
  { value: '12+', label: '重点航线' },
  { value: '1000+', label: '月度询盘承接' }
]

const routes: RouteItem[] = [
  { name: '中国 -> 美国', meta: '美西 / 美东 / FBA' },
  { name: '中国 -> 欧洲', meta: '汉堡 / 鹿特丹 / 安特卫普' },
  { name: '中国 -> 东南亚', meta: '新加坡 / 曼谷 / 胡志明' },
  { name: '中国 -> 中东', meta: '迪拜 / 杰贝阿里' }
]

const scenarios = ['整柜 FCL', '拼箱 LCL', '空运急货', '跨境电商', '门到门', '带电产品']

const tools: ToolItem[] = [
  {
    title: '出货计划分析',
    desc: '上传 Excel / CSV，快速得到体积、重量和柜型建议。',
    cta: '立即分析',
    icon: CalculatorOutlined,
    action: 'shipment-agent'
  },
  {
    title: '装柜与拼箱判断',
    desc: '先判断整柜还是拼箱，减少来回沟通。',
    cta: '开始测算',
    icon: InboxOutlined,
    action: 'shipment-agent'
  },
  {
    title: '智能物流问答',
    desc: '先问清路线、时效、清关和报价逻辑。',
    cta: '打开助手',
    icon: CommentOutlined,
    action: 'assistant'
  }
]

const services: ServiceItem[] = [
  { title: '海运整柜', desc: '报价、订舱、拖车、报关和目的港协同。', icon: ContainerOutlined },
  { title: '海运拼箱', desc: '适合小批量货物，按体积或重量灵活组合。', icon: InboxOutlined },
  { title: '空运快件', desc: '适合高时效订单、样品和紧急补货。', icon: ThunderboltOutlined },
  { title: '跨境物流', desc: '支持海外仓、双清和门到门方案。', icon: GlobalOutlined }
]

const contactReasons: ContactReason[] = [
  { title: '报价更快', desc: '留下起运地、目的地和货物信息后可直接进入报价。' },
  { title: '方案更准', desc: '先明确箱数、体积、重量和时效，减少反复确认。' },
  { title: '沟通更短', desc: '工具先分析，销售只接复杂问题和最终方案。' }
]

const quoteForm = reactive<QuoteForm>({
  contactName: '',
  origin: '',
  destination: '',
  goods: '',
  contact: ''
})

const quoteSubmitting = ref(false)

function handleToolAction(action: ToolAction) {
  if (action === 'shipment-agent') {
    router.push('/shipment-agent')
    return
  }
  if (action === 'contact') {
    router.push('/contact')
    return
  }
  window.dispatchEvent(new CustomEvent('portal-agent:open'))
}

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
.tool-card,
.service-card,
.panel {
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

.ghost {
  color: #fff;
  background: rgba(255, 255, 255, 0.12);
  border-color: rgba(255, 255, 255, 0.18);
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
.tool-card :deep(.ant-card-body),
.service-card :deep(.ant-card-body),
.panel :deep(.ant-card-body) {
  padding: 24px;
}

.side-card h3,
.tool-card h3,
.service-card h3,
.section-head h2 {
  margin: 0;
}

.route-list,
.bullet-list {
  display: grid;
  gap: 12px;
  margin-top: 16px;
}

.route-item,
.bullet-item {
  border-radius: 14px;
  border: 1px solid rgba(148, 163, 184, 0.16);
  background: #f8fafc;
  padding: 14px 16px;
}

.route-item strong,
.route-item span,
.bullet-item strong,
.bullet-item span {
  display: block;
}

.route-item span,
.bullet-item span,
.section-head p,
.tool-card p,
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

.tool-grid,
.service-grid {
  display: grid;
  gap: 18px;
}

.tool-grid {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.service-grid {
  grid-template-columns: repeat(4, minmax(0, 1fr));
}

.tool-card,
.service-card,
.panel {
  background: #fff;
}

.icon {
  font-size: 28px;
  color: #1677ff;
}

.tool-card h3,
.service-card h3 {
  margin-top: 16px;
  font-size: 20px;
}

.tool-action {
  margin-top: 18px;
  padding: 0;
  border: 0;
  background: transparent;
  color: #1677ff;
  font-weight: 700;
  cursor: pointer;
}

.final-section {
  padding-bottom: 8px;
}

.quote-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

@media (max-width: 1100px) {
  .hero,
  .tool-grid,
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

  .quote-actions {
    display: grid;
  }
}
</style>
