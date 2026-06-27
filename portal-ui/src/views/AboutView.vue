<template>
  <main class="page">
    <a-card class="hero" :bordered="false">
      <a-tag color="gold">关于我们</a-tag>
      <h2>这是一个可直接继续开发的门户工程底座</h2>
      <p>如果你后续想接 CMS、登录、权限、消息中心或者多语言，这个结构都能继续往上扩。</p>
    </a-card>

    <a-row :gutter="[18, 18]" class="section">
      <a-col :xs="24" :lg="14">
        <a-card class="panel" :bordered="false">
          <a-timeline>
            <a-timeline-item color="blue">首屏已准备好展示品牌和业务价值。</a-timeline-item>
            <a-timeline-item color="green">栏目结构支持新闻、服务和专题扩展。</a-timeline-item>
            <a-timeline-item color="orange">可继续接入后端接口和内容管理系统。</a-timeline-item>
          </a-timeline>
        </a-card>
      </a-col>

      <a-col :xs="24" :lg="10">
        <a-card class="panel contact" :bordered="false">
          <h3>联系信息</h3>
          <p>邮箱：portal@example.com</p>
          <p>电话：400-888-2026</p>
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
