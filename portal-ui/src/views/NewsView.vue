<template>
  <main class="page">
    <a-card class="hero" :bordered="false">
      <div class="hero-inner">
        <div>
          <a-tag color="blue">新闻中心</a-tag>
          <h2>把通知、动态和专题统一收口</h2>
          <p>门户最常见的事情，就是让内容可读、可查、可持续更新。这个页面用来承载列表、筛选和详情入口。</p>
        </div>
        <a-button type="primary">发布资讯</a-button>
      </div>
    </a-card>

    <a-row :gutter="[18, 18]" class="section">
      <a-col :xs="24" :lg="16">
        <a-card class="panel" :bordered="false">
          <a-list :data-source="items" item-layout="horizontal">
            <template #renderItem="{ item }">
              <a-list-item>
                <a-list-item-meta :description="item.desc">
                  <template #title>
                    <a href="javascript:void(0)">{{ item.title }}</a>
                  </template>
                </a-list-item-meta>
                <a-space direction="vertical" size="small">
                  <a-tag :color="item.tagColor">{{ item.tag }}</a-tag>
                  <span class="meta">{{ item.time }}</span>
                </a-space>
              </a-list-item>
            </template>
          </a-list>
        </a-card>
      </a-col>

      <a-col :xs="24" :lg="8">
        <a-card class="panel" :bordered="false">
          <h3>筛选栏目</h3>
          <a-space wrap>
            <a-tag v-for="tag in filters" :key="tag" color="blue">{{ tag }}</a-tag>
          </a-space>
        </a-card>
      </a-col>
    </a-row>
  </main>
</template>

<script setup>
const items = [
  { title: '门户首页模板完成', desc: '首屏、栏目、动态和联系区都已搭好。', time: '2026-06-27', tag: '公告', tagColor: 'green' },
  { title: '内容中心支持扩展', desc: '后续可以接 CMS、接口或静态数据。', time: '2026-06-26', tag: '资讯', tagColor: 'blue' },
  { title: '专题页面可继续追加', desc: '节庆、活动、政策专题都能做。', time: '2026-06-24', tag: '专题', tagColor: 'gold' }
]

const filters = ['全部', '公告', '资讯', '专题', '活动', '政策']
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
  background: linear-gradient(135deg, #ffffff, #eef5ff);
}

.hero-inner {
  display: flex;
  align-items: end;
  justify-content: space-between;
  gap: 20px;
}

.hero h2 {
  margin: 14px 0 0;
  font-size: clamp(28px, 3vw, 44px);
}

.hero p {
  max-width: 56ch;
  margin: 14px 0 0;
  color: #66748b;
  line-height: 1.8;
}

.section {
  margin-top: 6px;
}

.panel h3 {
  margin: 0 0 16px;
  font-size: 20px;
}

.meta {
  color: #7a879c;
  font-size: 12px;
}

@media (max-width: 760px) {
  .page {
    width: min(100% - 20px, 1240px);
  }

  .hero-inner {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
