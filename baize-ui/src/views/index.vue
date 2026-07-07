<template>
  <div class="app-container home-page">
    <section class="hero-panel">
      <div class="hero-copy">
        <span class="hero-kicker">后台管理首页</span>
        <h1>欢迎回来，{{ userName }}</h1>
        <p>
          这里聚合了货代业务、客户管理、Agent 对话和系统运维的常用入口，适合作为日常操作起点。
        </p>
        <div class="hero-meta">
          <span>当前版本 {{ version }}</span>
          <span>系统角色 {{ roleText }}</span>
          <span>{{ todayText }}</span>
        </div>
        <div class="hero-actions">
          <el-button type="primary" size="mini" @click="goRoute('/freight/shipment')">进入出货计划</el-button>
          <el-button size="mini" @click="goRoute('/customer/customer')">查看客户资料</el-button>
          <el-button size="mini" @click="goRoute('/agent/chat')">打开 Agent 对话</el-button>
        </div>
      </div>

      <div class="hero-side">
        <div class="pulse-card">
          <div class="pulse-head">
            <span>工作台状态</span>
            <strong>就绪</strong>
          </div>
          <ul class="pulse-list">
            <li>货代业务、客户模块、系统模块已集成到同一后台。</li>
            <li>客户端工作台配置已支持菜单、角色与账号联动。</li>
            <li>建议优先从出货计划、客户资料、Agent 对话三个入口开始日常操作。</li>
          </ul>
        </div>
      </div>
    </section>

    <el-row :gutter="16" class="stat-row">
      <el-col v-for="item in statCards" :key="item.title" :xs="24" :sm="12" :lg="6">
        <div class="stat-card">
          <span class="stat-label">{{ item.title }}</span>
          <strong class="stat-value">{{ item.value }}</strong>
          <p>{{ item.desc }}</p>
        </div>
      </el-col>
    </el-row>

    <section class="section-block">
      <div class="section-head">
        <div>
          <h2>常用入口</h2>
          <p>按后台日常操作频率整理，减少层层展开菜单。</p>
        </div>
      </div>

      <el-row :gutter="16">
        <el-col v-for="item in quickEntries" :key="item.title" :xs="24" :sm="12" :lg="8">
          <div class="entry-card" @click="goRoute(item.path)">
            <div class="entry-top">
              <span class="entry-badge">{{ item.group }}</span>
              <strong>{{ item.title }}</strong>
            </div>
            <p>{{ item.desc }}</p>
            <span class="entry-link">打开页面</span>
          </div>
        </el-col>
      </el-row>
    </section>

    <el-row :gutter="16" class="lower-row">
      <el-col :xs="24" :lg="14">
        <div class="panel-card">
          <div class="panel-head">
            <h2>业务分区</h2>
            <p>把当前系统里的核心模块按用途归并，便于快速定位。</p>
          </div>
          <div class="domain-list">
            <div v-for="item in domainCards" :key="item.title" class="domain-card">
              <div class="domain-title">
                <strong>{{ item.title }}</strong>
                <span>{{ item.tag }}</span>
              </div>
              <p>{{ item.desc }}</p>
              <div class="domain-links">
                <el-button
                  v-for="link in item.links"
                  :key="link.title"
                  text
                  size="small"
                  @click.stop="goRoute(link.path)"
                >
                  {{ link.title }}
                </el-button>
              </div>
            </div>
          </div>
        </div>
      </el-col>

      <el-col :xs="24" :lg="10">
        <div class="panel-card info-card">
          <div class="panel-head">
            <h2>系统信息</h2>
            <p>当前环境与使用建议。</p>
          </div>
          <div class="info-grid">
            <div class="info-item">
              <span>系统名称</span>
              <strong>IFS 管理系统</strong>
            </div>
            <div class="info-item">
              <span>当前版本</span>
              <strong>{{ version }}</strong>
            </div>
            <div class="info-item">
              <span>前端框架</span>
              <strong>Vue 3 + Element Plus</strong>
            </div>
            <div class="info-item">
              <span>后端框架</span>
              <strong>Go + Gin</strong>
            </div>
          </div>
          <div class="notice-box">
            <strong>建议操作顺序</strong>
            <ol>
              <li>先在客户资料里确认客户与负责业务员归属。</li>
              <li>再进入出货计划维护状态、分享进度或绑定客户。</li>
              <li>需要辅助分析时，直接进入 Agent 对话页面。</li>
            </ol>
          </div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup name="Index">
import configs from '../../package.json'

const router = useRouter()
const store = useStore()

const version = configs.version
const userName = computed(() => store.state.user.name || '管理员')
const roleText = computed(() => {
  const roles = store.state.user.roles || []
  return roles.length ? roles.join(' / ') : '未加载角色'
})

const todayText = computed(() => {
  const now = new Date()
  const year = now.getFullYear()
  const month = `${now.getMonth() + 1}`.padStart(2, '0')
  const day = `${now.getDate()}`.padStart(2, '0')
  return `${year}-${month}-${day}`
})

const statCards = [
  { title: '核心业务域', value: '4', desc: '货代业务、客户管理、系统模块、开发工具。' },
  { title: '常用快捷入口', value: '6', desc: '首页直接进入最常用的页面，减少菜单层级。' },
  { title: '客户端能力', value: '3', desc: '账号、角色、菜单都已纳入后台统一配置。' },
  { title: '协作方式', value: 'Agent', desc: '后台支持 AI 辅助对话，用于出货分析与问答。' },
]

const quickEntries = [
  { title: '出货计划', group: '货代业务', desc: '维护出货计划、客户绑定、状态流转与分享。', path: '/freight/shipment' },
  { title: 'Agent 对话', group: '智能助手', desc: '用于出货分析、单据辅助和业务问答。', path: '/agent/chat' },
  { title: '客户资料', group: '客户管理', desc: '管理客户信息、业务员归属与联系资料。', path: '/customer/customer' },
  { title: '客户账号', group: '客户管理', desc: '维护客户端登录账号及角色分配。', path: '/customer/account' },
  { title: '客户端菜单', group: '工作台配置', desc: '配置客户端菜单树和页面入口。', path: '/customer/portalMenu' },
  { title: '客户端角色', group: '工作台配置', desc: '配置客户端角色和菜单权限。', path: '/customer/portalRole' },
]

const domainCards = [
  {
    title: '货代业务',
    tag: '业务执行',
    desc: '围绕出货计划、状态维护、客户分享与执行跟进展开，适合运营和业务员日常使用。',
    links: [
      { title: '出货计划', path: '/freight/shipment' },
      { title: 'Agent 对话', path: '/agent/chat' },
    ],
  },
  {
    title: '客户管理',
    tag: '客户资产',
    desc: '统一管理官网线索、客户资料和客户账号，保证客户归属与客户端权限一致。',
    links: [
      { title: '官网线索', path: '/customer/contact' },
      { title: '客户资料', path: '/customer/customer' },
      { title: '客户账号', path: '/customer/account' },
    ],
  },
  {
    title: '客户端工作台配置',
    tag: '权限编排',
    desc: '通过菜单和角色管理，控制客户端能看到什么、能进入哪些页面。',
    links: [
      { title: '客户端菜单', path: '/customer/portalMenu' },
      { title: '客户端角色', path: '/customer/portalRole' },
    ],
  },
  {
    title: '系统与工具',
    tag: '运维开发',
    desc: '系统管理、监控、代码生成和接口文档用于支持后台持续迭代。',
    links: [
      { title: '服务监控', path: '/monitor/server' },
      { title: 'Swagger', path: '/tool/swagger' },
      { title: '代码生成', path: '/tool/gen' },
    ],
  },
]

function goRoute(path) {
  router.push(path)
}
</script>

<style scoped lang="scss">
.home-page {
  min-height: calc(100vh - 84px);
  background:
    radial-gradient(circle at top left, rgba(37, 99, 235, 0.14), transparent 22%),
    radial-gradient(circle at right bottom, rgba(14, 165, 233, 0.12), transparent 24%),
    linear-gradient(180deg, #f4f7fb 0%, #eef3f9 100%);
}

.hero-panel {
  display: grid;
  grid-template-columns: minmax(0, 1.7fr) minmax(280px, 0.9fr);
  gap: 18px;
  margin-bottom: 16px;
}

.hero-copy,
.pulse-card,
.stat-card,
.entry-card,
.panel-card {
  border: 1px solid rgba(148, 163, 184, 0.18);
  background: rgba(255, 255, 255, 0.9);
  box-shadow: 0 16px 40px rgba(15, 23, 42, 0.06);
}

.hero-copy {
  padding: 28px 30px;
  border-radius: 24px;
}

.hero-kicker {
  display: inline-flex;
  align-items: center;
  height: 28px;
  padding: 0 12px;
  border-radius: 999px;
  background: rgba(15, 23, 42, 0.06);
  color: #0f172a;
  font-size: 12px;
  font-weight: 600;
  letter-spacing: 0.08em;
}

.hero-copy h1 {
  margin: 16px 0 10px;
  color: #0f172a;
  font-size: 34px;
  line-height: 1.15;
}

.hero-copy p {
  margin: 0;
  max-width: 760px;
  color: #475569;
  font-size: 15px;
  line-height: 1.8;
}

.hero-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 18px;
}

.hero-meta span {
  display: inline-flex;
  align-items: center;
  min-height: 32px;
  padding: 0 12px;
  border-radius: 999px;
  background: #f8fafc;
  color: #334155;
  font-size: 13px;
}

.hero-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 22px;
}

.hero-side {
  min-width: 0;
}

.pulse-card {
  height: 100%;
  padding: 22px 22px 18px;
  border-radius: 24px;
  background:
    linear-gradient(180deg, rgba(15, 23, 42, 0.96) 0%, rgba(30, 41, 59, 0.94) 100%);
  color: #e2e8f0;
}

.pulse-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 18px;
}

.pulse-head span {
  font-size: 14px;
  color: rgba(226, 232, 240, 0.8);
}

.pulse-head strong {
  display: inline-flex;
  align-items: center;
  height: 30px;
  padding: 0 12px;
  border-radius: 999px;
  background: rgba(34, 197, 94, 0.16);
  color: #86efac;
  font-size: 13px;
}

.pulse-list {
  margin: 0;
  padding-left: 18px;
  color: rgba(226, 232, 240, 0.9);
  line-height: 1.8;
}

.stat-row {
  margin-bottom: 16px;
}

.stat-card {
  height: 100%;
  min-height: 138px;
  padding: 20px 20px 18px;
  border-radius: 20px;
}

.stat-label {
  display: block;
  color: #64748b;
  font-size: 13px;
}

.stat-value {
  display: block;
  margin: 8px 0 10px;
  color: #0f172a;
  font-size: 32px;
  line-height: 1;
}

.stat-card p {
  margin: 0;
  color: #475569;
  line-height: 1.75;
}

.section-block {
  margin-bottom: 16px;
}

.section-head,
.panel-head {
  margin-bottom: 14px;
}

.section-head h2,
.panel-head h2 {
  margin: 0 0 6px;
  color: #0f172a;
  font-size: 22px;
}

.section-head p,
.panel-head p {
  margin: 0;
  color: #64748b;
  line-height: 1.7;
}

.entry-card {
  height: 100%;
  padding: 20px;
  border-radius: 20px;
  cursor: pointer;
  transition:
    transform 0.2s ease,
    box-shadow 0.2s ease,
    border-color 0.2s ease;
}

.entry-card:hover {
  transform: translateY(-3px);
  border-color: rgba(37, 99, 235, 0.26);
  box-shadow: 0 22px 40px rgba(37, 99, 235, 0.12);
}

.entry-top {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 10px;
}

.entry-badge {
  display: inline-flex;
  align-items: center;
  width: fit-content;
  height: 26px;
  padding: 0 10px;
  border-radius: 999px;
  background: #eff6ff;
  color: #1d4ed8;
  font-size: 12px;
  font-weight: 600;
}

.entry-top strong {
  color: #0f172a;
  font-size: 18px;
}

.entry-card p {
  margin: 0;
  color: #475569;
  line-height: 1.8;
}

.entry-link {
  display: inline-block;
  margin-top: 16px;
  color: #1d4ed8;
  font-weight: 600;
}

.lower-row {
  margin-bottom: 6px;
}

.panel-card {
  height: 100%;
  padding: 22px;
  border-radius: 24px;
}

.domain-list {
  display: grid;
  gap: 14px;
}

.domain-card {
  padding: 18px;
  border-radius: 18px;
  background: #f8fafc;
  border: 1px solid rgba(148, 163, 184, 0.14);
}

.domain-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 10px;
}

.domain-title strong {
  color: #0f172a;
  font-size: 17px;
}

.domain-title span {
  display: inline-flex;
  align-items: center;
  height: 26px;
  padding: 0 10px;
  border-radius: 999px;
  background: #e2e8f0;
  color: #334155;
  font-size: 12px;
}

.domain-card p {
  margin: 0;
  color: #475569;
  line-height: 1.8;
}

.domain-links {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 14px;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.info-item {
  padding: 16px;
  border-radius: 18px;
  background: #f8fafc;
  border: 1px solid rgba(148, 163, 184, 0.14);
}

.info-item span {
  display: block;
  margin-bottom: 8px;
  color: #64748b;
  font-size: 13px;
}

.info-item strong {
  color: #0f172a;
  font-size: 15px;
  line-height: 1.6;
}

.notice-box {
  margin-top: 16px;
  padding: 18px;
  border-radius: 18px;
  background: linear-gradient(180deg, #eff6ff 0%, #f8fafc 100%);
  border: 1px solid rgba(96, 165, 250, 0.18);
}

.notice-box strong {
  display: block;
  margin-bottom: 10px;
  color: #0f172a;
  font-size: 16px;
}

.notice-box ol {
  margin: 0;
  padding-left: 18px;
  color: #475569;
  line-height: 1.9;
}

@media (max-width: 1200px) {
  .hero-panel {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .home-page {
    padding-bottom: 8px;
  }

  .hero-copy {
    padding: 22px 18px;
  }

  .hero-copy h1 {
    font-size: 28px;
  }

  .panel-card,
  .pulse-card,
  .stat-card,
  .entry-card {
    border-radius: 18px;
  }

  .info-grid {
    grid-template-columns: 1fr;
  }
}
</style>
