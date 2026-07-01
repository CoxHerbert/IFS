<template>
  <main class="page">
    <a-card class="hero" :bordered="false">
      <div class="hero-inner">
        <div>
          <a-tag color="blue">航线资讯</a-tag>
          <h2>把运价、航线、政策变化做成持续更新的内容中心</h2>
          <p>
            货代官网的内容不是为了“写文章”，而是为了让客户知道你当前主做哪些线路、是否有现舱/舱位、哪些政策会影响出货。
          </p>
        </div>
        <a-button type="primary">发布动态</a-button>
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
          <h3>栏目标签</h3>
          <a-space wrap>
            <a-tag v-for="tag in filters" :key="tag" color="blue">{{ tag }}</a-tag>
          </a-space>
        </a-card>
      </a-col>
    </a-row>
  </main>
</template>

<script setup lang="ts">
interface NewsItem {
  title: string
  desc: string
  time: string
  tag: string
  tagColor: string
}

const items: NewsItem[] = [
  { title: '美线舱位更新', desc: '适合展示当前航线是否有价格波动或舱位变化。', time: '2026-06-27', tag: '美线', tagColor: 'green' },
  { title: '欧线操作提醒', desc: '适合发布清关、截单、截港等关键节点提醒。', time: '2026-06-26', tag: '欧线', tagColor: 'blue' },
  { title: '东南亚专线优惠', desc: '适合做短期促销、拼箱拼柜活动和限时报价。', time: '2026-06-24', tag: '促销', tagColor: 'gold' }
]

const filters: string[] = ['全部', '美线', '欧线', '东南亚', '中东', '空运', '海运']
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
