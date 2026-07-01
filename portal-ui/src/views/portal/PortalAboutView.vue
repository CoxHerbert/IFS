<template>
  <main class="page">
    <a-card class="hero" :bordered="false">
      <a-tag color="gold">关于我们</a-tag>
      <h2>这是一套面向货代销售增长的官网底座</h2>
      <p>
        它的目标很明确：让客户快速了解你的航线、服务能力和联系方式。
        后续你可以继续接接口、接 CMS、接登录或者做多语言版本。
      </p>
    </a-card>

    <a-row :gutter="[18, 18]" class="section">
      <a-col :xs="24" :lg="14">
        <a-card class="panel" :bordered="false">
          <a-timeline>
            <a-timeline-item color="blue">首页突出询价入口和航线能力。</a-timeline-item>
            <a-timeline-item color="green">新闻页承载航线动态和政策提醒。</a-timeline-item>
            <a-timeline-item color="orange">服务页展示整柜、拼箱、空运和海外仓。</a-timeline-item>
            <a-timeline-item color="purple">关于页补充信任、资质和联系方式。</a-timeline-item>
          </a-timeline>
        </a-card>
      </a-col>

      <a-col :xs="24" :lg="10">
        <a-card class="panel contact" :bordered="false">
          <h3>联系信息</h3>
          <p>邮箱：quote@seawaypro.com</p>
          <p>电话：400-888-2026</p>
          <p>微信：SeaWayPro</p>
          <p>地址：上海市徐汇区</p>
        </a-card>
      </a-col>
    </a-row>
  </main>
</template>

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
  background: linear-gradient(135deg, #ffffff, #f3f8ff);
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

.contact p {
  color: #66748b;
}

@media (max-width: 760px) {
  .page {
    width: min(100% - 20px, 1240px);
  }
}
</style>
