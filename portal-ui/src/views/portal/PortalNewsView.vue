<template>
  <main class="page">
    <template v-if="activeArticle">
      <router-link to="/news" class="back-link">返回新闻资讯</router-link>
      <article class="article-detail">
        <a-tag :color="activeArticle.color">{{ activeArticle.category }}</a-tag>
        <h1>{{ activeArticle.title }}</h1>
        <p class="lead">{{ activeArticle.summary }}</p>
        <div class="article-meta">
          <span>{{ activeArticle.publishedAt }}</span>
          <span>{{ activeArticle.readingTime }}</span>
          <span>{{ activeArticle.author }}</span>
        </div>
        <div class="article-body" v-html="activeArticle.contentHtml"></div>
      </article>
    </template>

    <template v-else>
      <section class="hero">
        <div>
          <a-tag color="blue">新闻资讯</a-tag>
          <h1>关注航线、运价、舱位和政策变化</h1>
          <p>持续更新国际物流新闻、航线动态、操作提醒和市场观察，帮助客户及时判断出货节奏。</p>
        </div>
        <a-input-search
          v-model:value="keyword"
          class="search"
          placeholder="搜索新闻、航线、港口、政策"
          allow-clear
        />
      </section>

      <section class="category-tabs">
        <button
          v-for="category in categories"
          :key="category"
          type="button"
          :class="{ active: category === activeCategory }"
          @click="activeCategory = category"
        >
          {{ category }}
        </button>
      </section>

      <section class="content-grid">
        <div class="main-column">
          <div class="section-title">
            <span>头条新闻</span>
          </div>
          <router-link :to="`/news/${featured.slug}`" class="featured-card">
            <a-tag :color="featured.color">{{ featured.category }}</a-tag>
            <h2>{{ featured.title }}</h2>
            <p>{{ featured.summary }}</p>
            <div class="article-meta">
              <span>{{ featured.publishedAt }}</span>
              <span>{{ featured.author }}</span>
              <span>{{ featured.readingTime }}</span>
            </div>
          </router-link>

          <div class="section-title">
            <span>最新资讯</span>
          </div>
          <div class="article-list">
            <router-link
              v-for="item in filteredArticles"
              :key="item.slug"
              :to="`/news/${item.slug}`"
              class="article-card"
            >
              <div>
                <a-tag :color="item.color">{{ item.category }}</a-tag>
                <h3>{{ item.title }}</h3>
                <p>{{ item.summary }}</p>
              </div>
              <div class="article-meta">
                <span>{{ item.publishedAt }}</span>
                <span>{{ item.author }}</span>
                <span>{{ item.readingTime }}</span>
              </div>
            </router-link>
            <a-empty v-if="filteredArticles.length === 0" description="暂无匹配内容" />
          </div>
        </div>

        <aside class="side-column">
          <section class="side-panel">
            <h3>热门分类</h3>
            <div class="tag-cloud">
              <button
                v-for="category in categories.slice(1)"
                :key="category"
                type="button"
                @click="activeCategory = category"
              >
                {{ category }}
              </button>
            </div>
          </section>

          <section class="side-panel">
            <h3>新闻速递</h3>
            <router-link
              v-for="item in articles.slice(0, 4)"
              :key="item.slug"
              :to="`/news/${item.slug}`"
              class="mini-link"
            >
              <strong>{{ item.title }}</strong>
              <span>{{ item.publishedAt }}</span>
            </router-link>
          </section>
        </aside>
      </section>
    </template>
  </main>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import DOMPurify from 'dompurify'
import { getArticleBySlug, listArticles } from '@/api/portal/article'

interface Article {
  slug: string
  title: string
  summary: string
  category: string
  color: string
  publishedAt: string
  readingTime: string
  author: string
  content: string[]
  contentHtml?: string
  searchText?: string
}

const route = useRoute()
const keyword = ref('')
const activeCategory = ref('全部')
const articles = ref<Article[]>([])
const activeDetail = ref<Article>()

const categories = ['全部', '美线', '欧线', '东南亚', '中东', '拼箱', '政策']

const fallbackArticles: Article[] = [
  {
    slug: 'us-west-space-update',
    title: '美西舱位与截港节奏更新',
    summary: '近期美西主要港口舱位相对稳定，但部分船司截单时间提前，建议提前确认装柜和报关资料。',
    category: '美线',
    color: 'green',
    publishedAt: '2026-06-27',
    readingTime: '3 分钟阅读',
    author: 'IFS 航线团队',
    content: [
      '美西航线近期整体舱位保持稳定，洛杉矶、长滩方向仍是询价和出货最集中的目的港。',
      '建议客户在确认出货计划时同步准备装箱单、发票和申报要素，避免截单前资料不齐影响订舱。',
      '如货物涉及带电、液体、粉末或品牌属性，请在询价时提前说明，便于判断是否需要额外审核。'
    ]
  },
  {
    slug: 'eu-customs-reminder',
    title: '欧线清关资料准备提醒',
    summary: '欧洲方向对品名、材质、用途和 HS Code 的一致性要求较高，资料准备越完整，目的港衔接越顺畅。',
    category: '欧线',
    color: 'blue',
    publishedAt: '2026-06-26',
    readingTime: '4 分钟阅读',
    author: 'IFS 操作中心',
    content: [
      '欧线出货建议在订舱前确认品名、材质、用途、品牌和申报价值，减少目的港补资料的概率。',
      '如果涉及多个 SKU，需要确保箱单、发票和唛头信息一致。',
      '对时效敏感的订单，建议提前沟通目的港清关方式和派送要求。'
    ]
  },
  {
    slug: 'sea-lcl-promo',
    title: '东南亚拼箱短期优惠线路',
    summary: '新加坡、曼谷、胡志明方向适合小批量补货，近期可优先评估拼箱或小柜方案。',
    category: '东南亚',
    color: 'gold',
    publishedAt: '2026-06-24',
    readingTime: '2 分钟阅读',
    author: 'IFS 销售支持',
    content: [
      '东南亚方向小票货近期可以优先评估拼箱方案，适合样品、补货和多 SKU 小批量订单。',
      '如果总体积接近整柜安全装载范围，可同步比较小柜和拼箱成本。',
      '询价时提供件数、箱规、重量和目的城市，可以更快得到可执行方案。'
    ]
  },
  {
    slug: 'middle-east-booking-note',
    title: '中东航线订舱与目的港费用提示',
    summary: '杰贝阿里方向订舱需求稳定，目的港费用和收货人资料建议提前核对。',
    category: '中东',
    color: 'purple',
    publishedAt: '2026-06-21',
    readingTime: '3 分钟阅读',
    author: 'IFS 航线团队',
    content: [
      '中东方向订舱前建议确认收货人资料、目的港清关方式和是否需要门到门服务。',
      '部分目的港费用会随船司和服务条款变化，报价时需确认费用边界。',
      '对于工程类、设备类货物，建议提前提供尺寸和重量，便于评估装载风险。'
    ]
  }
]

const featured = computed(() => articles.value[0] || fallbackArticles[0])

const activeArticle = computed(() => {
  const slug = String(route.params.slug || '')
  return activeDetail.value || articles.value.find((item) => item.slug === slug)
})

const filteredArticles = computed(() => {
  const text = keyword.value.trim().toLowerCase()
  return articles.value.filter((item) => {
    const categoryMatched = activeCategory.value === '全部' || item.category === activeCategory.value
    const keywordMatched =
      !text ||
      [item.title, item.summary, item.category, item.searchText || item.content.join(' ')].some((value) => value.toLowerCase().includes(text))
    return categoryMatched && keywordMatched
  })
})

onMounted(loadArticles)

watch(
  () => route.params.slug,
  () => loadDetail(),
)

async function loadArticles() {
  try {
    const rows = await listArticles()
    articles.value = rows.length ? rows.map(normalizeArticle) : fallbackArticles
  } catch (_error) {
    articles.value = fallbackArticles
  }
  await loadDetail()
}

async function loadDetail() {
  const slug = String(route.params.slug || '')
  activeDetail.value = undefined
  if (!slug) {
    return
  }
  try {
    activeDetail.value = normalizeArticle(await getArticleBySlug(slug))
  } catch (_error) {
    activeDetail.value = articles.value.find((item) => item.slug === slug)
  }
}

function normalizeArticle(item: any): Article {
  const content = String(item.content || '')
  const paragraphs = content.split(/\n+/).map((text) => text.trim()).filter(Boolean)
  return {
    slug: item.slug,
    title: item.title,
    summary: item.summary || '',
    category: item.category || '资讯',
    color: categoryColor(item.category),
    publishedAt: (item.publishTime || '').slice(0, 10) || '未发布',
    readingTime: estimateReadingTime(content),
    author: item.updateBy || item.createBy || 'IFS 航线团队',
    content: paragraphs,
    contentHtml: normalizeContentHtml(content),
    searchText: stripHtml(content),
  }
}

function normalizeContentHtml(content: string) {
  if (/<[a-z][\s\S]*>/i.test(content)) {
    return sanitizeHtml(content)
  }
  const html = content
    .split(/\n+/)
    .map((text) => text.trim())
    .filter(Boolean)
    .map((text) => `<p>${escapeHtml(text)}</p>`)
    .join('')
  return sanitizeHtml(html)
}

function sanitizeHtml(content: string) {
  return DOMPurify.sanitize(content, {
    USE_PROFILES: { html: true },
    ADD_ATTR: ['target'],
  })
}

function stripHtml(content: string) {
  return content.replace(/<[^>]+>/g, ' ')
}

function escapeHtml(content: string) {
  return content
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;')
}

function categoryColor(category: string) {
  const colorMap: Record<string, string> = {
    美线: 'green',
    欧线: 'blue',
    东南亚: 'gold',
    中东: 'purple',
    拼箱: 'cyan',
    政策: 'red',
  }
  return colorMap[category] || 'blue'
}

function estimateReadingTime(content: string) {
  const length = String(content || '').length
  return `${Math.max(1, Math.ceil(length / 500))} 分钟阅读`
}
</script>

<style scoped>
.page {
  width: min(1240px, calc(100% - 32px));
  margin: 0 auto;
  padding: 28px 0 40px;
}

.hero,
.featured-card,
.article-card,
.side-panel,
.article-detail {
  border-radius: 8px;
  background: #fff;
  box-shadow: 0 18px 40px rgba(16, 35, 63, 0.08);
}

.hero {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 360px;
  gap: 24px;
  align-items: end;
  padding: 32px;
  background: linear-gradient(135deg, #ffffff, #eef5ff);
}

.hero h1 {
  max-width: 820px;
  margin: 14px 0 0;
  font-size: clamp(30px, 3.5vw, 48px);
  line-height: 1.16;
}

.hero p,
.featured-card p,
.article-card p,
.lead {
  color: #64748b;
  line-height: 1.8;
}

.search {
  width: 100%;
}

.category-tabs {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin: 18px 0;
}

.category-tabs button,
.tag-cloud button {
  border: 1px solid rgba(22, 119, 255, 0.18);
  border-radius: 999px;
  padding: 8px 14px;
  background: #fff;
  color: #334155;
  cursor: pointer;
}

.category-tabs button.active {
  background: #1677ff;
  color: #fff;
}

.content-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 320px;
  gap: 18px;
}

.main-column,
.article-list,
.side-column,
.tag-cloud {
  display: grid;
  gap: 14px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 10px;
  color: #0f172a;
  font-size: 18px;
  font-weight: 800;
}

.section-title::before {
  content: '';
  width: 4px;
  height: 18px;
  border-radius: 999px;
  background: #1677ff;
}

.featured-card,
.article-card,
.side-panel,
.article-detail {
  display: block;
  padding: 24px;
  color: inherit;
}

.featured-card h2,
.article-card h3,
.article-detail h1,
.side-panel h3 {
  margin: 12px 0 0;
  color: #0f172a;
}

.article-card {
  display: grid;
  gap: 12px;
  border-left: 3px solid transparent;
  transition: border-color 0.2s ease, transform 0.2s ease;
}

.article-card:hover,
.featured-card:hover {
  transform: translateY(-2px);
}

.article-card:hover {
  border-left-color: #1677ff;
}

.article-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  color: #7a879c;
  font-size: 12px;
}

.mini-link {
  display: grid;
  gap: 5px;
  color: inherit;
}

.mini-link span {
  color: #7a879c;
  font-size: 12px;
}

.back-link {
  display: inline-block;
  margin-bottom: 14px;
  color: #1677ff;
  font-weight: 700;
}

.article-detail {
  max-width: 860px;
  margin: 0 auto;
}

.article-detail h1 {
  font-size: clamp(30px, 4vw, 52px);
}

.article-body {
  margin-top: 24px;
  color: #334155;
  font-size: 16px;
  line-height: 2;
}

@media (max-width: 960px) {
  .hero,
  .content-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 640px) {
  .page {
    width: min(100% - 20px, 1240px);
  }

  .hero,
  .featured-card,
  .article-card,
  .side-panel,
  .article-detail {
    padding: 20px;
  }
}
</style>
