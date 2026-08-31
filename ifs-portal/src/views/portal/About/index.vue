<template>
  <main class="page">
    <a-card class="hero" :bordered="false">
      <a-tag color="gold">{{ copy.tag }}</a-tag>
      <h1>{{ copy.title }}</h1>
      <p>{{ copy.description }}</p>
    </a-card>

    <a-row :gutter="[18, 18]" class="section">
      <a-col :xs="24" :lg="14">
        <a-card class="panel" :bordered="false">
          <a-timeline>
            <a-timeline-item v-for="(item, index) in copy.points" :key="item" :color="colors[index]">{{ item }}</a-timeline-item>
          </a-timeline>
        </a-card>
      </a-col>

      <a-col :xs="24" :lg="10">
        <a-card class="panel contact" :bordered="false">
          <h2>{{ copy.contact }}</h2>
          <p>{{ copy.email }}: quote@seawaypro.com</p>
          <p>{{ copy.phone }}: 400-888-2026</p>
          <p>WeChat: SeaWayPro</p>
          <p>{{ copy.address }}: {{ copy.addressValue }}</p>
        </a-card>
      </a-col>
    </a-row>
  </main>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { usePortalI18n } from '@/i18n'
const { locale } = usePortalI18n()
const colors = ['blue', 'green', 'orange', 'purple']
const copy = computed(() => locale.value === 'en' ? {
  tag: 'ABOUT IFS', title: 'A practical international freight forwarding partner',
  description: 'We help customers understand available routes, transport options, documentation and the next step required to move cargo from China.',
  points: ['Start with cargo details and the required delivery outcome.', 'Compare ocean, air, FCL, LCL and door-to-door options.', 'Coordinate booking, documents, customs milestones and delivery.', 'Keep communication clear from enquiry through shipment completion.'],
  contact: 'Contact information', email: 'Email', phone: 'Phone', address: 'Address', addressValue: 'Xuhui District, Shanghai, China',
} : {
  tag: '关于我们', title: '专注于提供可执行的国际货运方案', description: '我们帮助客户了解航线、运输方式、单证要求和下一步安排，让中国出口运输更清晰。',
  points: ['从货物信息和交付目标开始确认需求。', '比较海运、空运、整柜、拼箱与门到门方案。', '协调订舱、单证、清关节点和最终交付。', '从询价到运输完成保持清晰沟通。'],
  contact: '联系信息', email: '邮箱', phone: '电话', address: '地址', addressValue: '上海市徐汇区',
})
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
  background: linear-gradient(135deg, #ffffff, #f3f8ff);
}

.hero h1 {
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
