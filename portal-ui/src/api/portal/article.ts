export interface ArticleItem {
  articleId?: string
  title: string
  slug: string
  summary: string
  category: string
  coverUrl?: string
  content: string
  status?: string
  sort?: number
  publishTime?: string
  createBy?: string
  updateBy?: string
}

export async function listArticles(query: Record<string, unknown> = {}) {
  const params = new URLSearchParams()
  Object.entries({ pageNum: 1, pageSize: 20, ...query }).forEach(([key, value]) => {
    const text = value == null ? '' : String(value)
    if (text !== '') {
      params.set(key, text)
    }
  })
  const response = await fetch(`/portal/articles?${params.toString()}`)
  const data = await response.json()
  if (!response.ok || String(data.code) !== '200') {
    throw new Error(data.msg || '文章加载失败')
  }
  return (data.rows || []) as ArticleItem[]
}

export async function getArticleBySlug(slug: string) {
  const response = await fetch(`/portal/articles/${encodeURIComponent(slug)}`)
  const data = await response.json()
  if (!response.ok || String(data.code) !== '200') {
    throw new Error(data.msg || '文章加载失败')
  }
  return data.data as ArticleItem
}
