<template>
  <main class="page">
    <section class="hero">
      <div class="hero-copy" :style="heroStyle">
        <a-tag color="blue">Vue 3 / Ant Design Vue / Vite</a-tag>
        <h2>一个能直接用于企业门户的前端工程模板</h2>
        <p>
          这套工程预置了门户首页、新闻、服务、关于四个页面，以及统一头部和底部。
          你可以直接替换内容、接接口、扩栏目，不需要重新搭架子。
        </p>

        <div class="hero-actions">
          <a-button type="primary" size="large">查看门户</a-button>
          <a-button size="large" class="ghost">下载资料</a-button>
        </div>

        <div class="stats">
          <a-card v-for="item in stats" :key="item.label" class="stat-card" :bordered="false">
            <div class="value">{{ item.value }}</div>
            <div class="label">{{ item.label }}</div>
          </a-card>
        </div>
      </div>

      <div class="hero-media">
        <img :src="heroImage" alt="portal preview" />
      </div>
    </section>

    <section class="section">
      <div class="section-head">
        <h3>快捷入口</h3>
        <p>把高频功能放到首页，用户少点一步。</p>
      </div>

      <a-row :gutter="[18, 18]">
        <a-col v-for="item in quickLinks" :key="item.title" :xs="24" :md="12" :lg="6">
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
              <h3>最新动态</h3>
              <p>公告、资讯、活动都能放这里。</p>
            </div>
            <a-list :data-source="news" item-layout="horizontal">
              <template #renderItem="{ item }">
                <a-list-item>
                  <a-list-item-meta :description="item.desc">
                    <template #title>
                      <a href="javascript:void(0)">{{ item.title }}</a>
                    </template>
                  </a-list-item-meta>
                  <a-tag :color="item.tagColor">{{ item.tag }}</a-tag>
                </a-list-item>
              </template>
            </a-list>
          </a-card>
        </a-col>

        <a-col :xs="24" :lg="10">
          <a-card class="panel" :bordered="false">
            <div class="section-head">
              <h3>系统说明</h3>
              <p>门户工程的基础能力已经搭好。</p>
            </div>
            <a-timeline>
              <a-timeline-item v-for="item in timeline" :key="item.title" :color="item.color">
                <strong>{{ item.title }}</strong>
                <p>{{ item.desc }}</p>
              </a-timeline-item>
            </a-timeline>
          </a-card>
        </a-col>
      </a-row>
    </section>
  </main>
</template>

<script setup>
import { AppstoreOutlined, DatabaseOutlined, TeamOutlined, ThunderboltOutlined } from '@ant-design/icons-vue'
import heroImage from '@/assets/hero.jpg'

const heroStyle = {
  backgroundImage: `linear-gradient(135deg, rgba(7, 23, 47, 0.94), rgba(16, 103, 200, 0.82)), url(${heroImage})`
}

const stats = [
  { value: '24/7', label: '在线服务' },
  { value: '12+', label: '核心模块' },
  { value: '99.9%', label: '可用性' }
]

const quickLinks = [
  { title: '统一导航', desc: '首页、新闻、服务、关于分栏清晰。', icon: AppstoreOutlined },
  { title: '内容发布', desc: '适合资讯、通知、专题和活动。', icon: DatabaseOutlined },
  { title: '运营展示', desc: '适合品牌介绍、案例和数据概览。', icon: TeamOutlined },
  { title: '安全访问', desc: '后续可无缝接入登录与权限系统。', icon: ThunderboltOutlined }
]

const news = [
  { title: '门户工程已创建完成', desc: '可直接进入开发或接接口。', tag: '上线', tagColor: 'green' },
  { title: '支持响应式布局', desc: '手机、平板、桌面都能跑。', tag: '布局', tagColor: 'blue' },
  { title: '首页内容可替换', desc: '模块和文案都方便扩展。', tag: '模板', tagColor: 'gold' }
]

const timeline = [
  { title: '首页', desc: '用于首屏展示品牌和入口。', color: 'blue' },
  { title: '新闻', desc: '用于公告、动态和内容发布。', color: 'green' },
  { title: '服务', desc: '用于业务系统和常用工具入口。', color: 'orange' }
]
</script>

<style scoped>
.page {
  width: min(1240px, calc(100% - 32px));
  margin: 0 auto;
  padding: 28px 0 40px;
}

.hero {
  display: grid;
  grid-template-columns: 1.05fr 0.95fr;
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

.feature p {
  margin: 10px 0 0;
  color: #66748b;
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

  .stats {
    grid-template-columns: 1fr;
  }
}
</style>
