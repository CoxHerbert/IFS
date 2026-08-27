<template>
  <main class="page">
    <template v-if="activeArticle">
      <router-link to="/news" class="back-link">返回新闻资讯</router-link>
      <article class="article-detail">
        <a-image
          v-if="activeArticle.coverUrl"
          :src="activeArticle.coverUrl"
          :alt="activeArticle.title"
          :preview="false"
          class="detail-cover"
        />
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

      <section v-if="articles.length" class="category-tabs">
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

      <a-alert
        v-if="loadError"
        class="load-state"
        type="error"
        show-icon
        :message="loadError"
      />
      <a-empty v-else-if="!loading && articles.length === 0" class="load-state" description="暂无新闻资讯" />

      <section v-else-if="articles.length" class="content-grid">
        <div class="main-column">
          <div class="section-title">
            <span>头条新闻</span>
          </div>
          <router-link v-if="featured" :to="`/news/${featured.slug}`" class="featured-card">
            <div v-if="featured.coverUrl" class="featured-cover">
              <a-image
                :src="featured.coverUrl"
                :alt="featured.title"
                :preview="false"
                class="cover-image"
              />
            </div>
            <div class="featured-content">
              <a-tag :color="featured.color">{{ featured.category }}</a-tag>
              <h2>{{ featured.title }}</h2>
              <p>{{ featured.summary }}</p>
              <div class="article-meta">
                <span>{{ featured.publishedAt }}</span>
                <span>{{ featured.author }}</span>
                <span>{{ featured.readingTime }}</span>
              </div>
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
              :class="{ 'article-card--with-cover': item.coverUrl }"
            >
              <div v-if="item.coverUrl" class="article-cover">
                <a-image
                  :src="item.coverUrl"
                  :alt="item.title"
                  :preview="false"
                  class="cover-image"
                />
              </div>
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
import { getArticleBySlug, listArticles, type ArticleItem } from '@/api/portal/article'

interface Article {
  slug: string
  title: string
  summary: string
  category: string
  coverUrl: string
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
const loading = ref(false)
const loadError = ref('')

const categories = computed(() => [
  '全部',
  ...new Set(articles.value.map((item) => item.category).filter(Boolean)),
])

const featured = computed(() => articles.value[0])

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
  loading.value = true
  loadError.value = ''
  try {
    const rows = await listArticles()
    articles.value = rows.map(normalizeArticle)
    if (!categories.value.includes(activeCategory.value)) {
      activeCategory.value = '全部'
    }
  } catch (error) {
    articles.value = []
    loadError.value = error instanceof Error ? error.message : '新闻资讯加载失败'
  } finally {
    loading.value = false
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

function normalizeArticle(item: ArticleItem): Article {
  const content = String(item.content || '')
  const paragraphs = content.split(/\n+/).map((text) => text.trim()).filter(Boolean)
  return {
    slug: item.slug,
    title: item.title,
    summary: item.summary || '',
    category: item.category || '资讯',
    coverUrl: item.coverUrl || '',
    color: categoryColor(item.category),
    publishedAt: formatArticleTime(item.publishTime),
    readingTime: estimateReadingTime(content),
    author: item.updateBy || item.createBy || 'IFS 航线团队',
    content: paragraphs,
    contentHtml: normalizeContentHtml(content),
    searchText: stripHtml(content),
  }
}

function formatArticleTime(value?: string | number) {
  if (value === undefined || value === null || value === '') return '未发布'

  const text = String(value)
  if (/^\d+$/.test(text)) {
    const timestamp = Number(text)
    const date = new Date(timestamp < 1e12 ? timestamp * 1000 : timestamp)
    if (!Number.isNaN(date.getTime())) {
      const pad = (number: number) => String(number).padStart(2, '0')
      return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`
    }
  }

  return text.replace('T', ' ').slice(0, 19)
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

.featured-card { overflow: hidden; }
.featured-content { padding: 0; }
.featured-cover { display: block; width: 160px; height: 100px; overflow: hidden; border-radius: 6px; align-self: start; }
.detail-cover { display: block; width: 100%; aspect-ratio: 16 / 9; overflow: hidden; }
.detail-cover :deep(.ant-image-img) { width: 100%; height: 100%; object-fit: cover; }
.cover-image { display: block; width: 100%; height: 100%; }
.cover-image :deep(.ant-image-img) { display: block; width: 100%; height: 100%; object-fit: cover; object-position: center; }
.detail-cover { margin-bottom: 24px; border-radius: 6px; }

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

.load-state {
  margin-top: 18px;
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

.article-card--with-cover { grid-template-columns: 160px minmax(0, 1fr); }
.article-card--with-cover .article-meta { grid-column: 2; }
.article-cover { grid-row: 1 / span 2; width: 160px; height: 100px; overflow: hidden; border-radius: 6px; align-self: start; }

.featured-card { display: grid; grid-template-columns: 160px minmax(0, 1fr); gap: 20px; padding: 24px; }

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

  .article-card--with-cover { grid-template-columns: 1fr; }
  .article-card--with-cover .article-meta { grid-column: auto; }
  .article-cover { grid-row: auto; width: 160px; height: 100px; aspect-ratio: auto; }
  .featured-card { grid-template-columns: 1fr; }
  .featured-cover { width: 160px; height: 100px; aspect-ratio: auto; }
}
</style>
