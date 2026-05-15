const API_BASE_URL = (import.meta.env.VITE_API_BASE_URL ?? 'http://127.0.0.1:8080').replace(/\/$/, '')

export type SiteSettings = {
  name: string
  description: string
  logoUrl: string
  role: string
  location: string
  contact: {
    email: string
    github: string
  }
}

export type HomeContent = {
  introTitle: string
  intro: string
  featuredArticleIds: string[]
}

export type WorkStackGroup = {
  title: string
  items: string[]
}

export type ExperienceItem = {
  period: string
  title: string
  body: string
}

export type AboutContent = {
  description: string
  workStack: WorkStackGroup[]
  experience: ExperienceItem[]
}

export type ArticleSummary = {
  id: string
  title: string
  date: string
  summary: string
  tags: string[]
  featured: boolean
  draft?: boolean
  wordCount?: number
}

export type Article = ArticleSummary & {
  content: string
  createdAt?: string
  updatedAt?: string
}

export type SiteContent = {
  site: SiteSettings
  home: HomeContent
  tags: string[]
  articles: Article[]
  about: AboutContent
  updatedAt: string
  featuredArticles: ArticleSummary[]
}

export function resolveAssetUrl(path: string | undefined) {
  if (!path) return ''
  if (/^https?:\/\//i.test(path) || path.startsWith('data:')) return path
  return `${API_BASE_URL}${path.startsWith('/') ? path : `/${path}`}`
}

export function estimateWords(source: string) {
  const chinese = source.match(/[\u4e00-\u9fa5]/g)?.length ?? 0
  const latin = source.match(/[A-Za-z0-9]+/g)?.length ?? 0
  return chinese + latin
}

export function getArticleWordCount(article: Pick<Article, 'title' | 'summary' | 'tags'> & { content?: string; wordCount?: number }) {
  return article.wordCount ?? estimateWords(`${article.title} ${article.summary} ${article.tags.join(' ')} ${article.content ?? ''}`)
}

export function getReadingMinutes(article: Pick<Article, 'title' | 'summary' | 'tags'> & { content?: string; wordCount?: number }) {
  return Math.max(1, Math.ceil(getArticleWordCount(article) / 400))
}

export async function fetchSiteContent(signal?: AbortSignal) {
  return request<SiteContent>('/api/site', signal)
}

export async function fetchArticles(signal?: AbortSignal) {
  const data = await request<{ articles: ArticleSummary[] }>('/api/articles', signal)
  return data.articles
}

export async function fetchArticle(id: string, signal?: AbortSignal) {
  return request<Article>(`/api/articles/${encodeURIComponent(id)}`, signal)
}

async function request<T>(path: string, signal?: AbortSignal): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    headers: {
      Accept: 'application/json',
    },
    signal,
  })

  if (!response.ok) {
    const message = await response.text().catch(() => '')
    throw new Error(message || `Request failed: ${response.status}`)
  }

  return response.json() as Promise<T>
}
