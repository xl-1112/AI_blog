import axios, { type AxiosError } from 'axios'
import type {
  AboutContent,
  AnalyticsData,
  Article,
  ArticleSummary,
  DashboardData,
  ExperienceItem,
  HomeContent,
  LoginLog,
  PageResult,
  SiteContent,
  SiteSettings,
  TagUsage,
  User,
  UserStatus,
  WorkStackGroup,
} from '../types'
import { useAuthStore } from '../store/auth'

export const API_BASE_URL = (import.meta.env.VITE_API_BASE_URL ?? 'http://127.0.0.1:8080').replace(/\/$/, '')

type Envelope<T> = {
  code: number
  message: string
  data: T
}

export const http = axios.create({
  baseURL: API_BASE_URL,
  timeout: 15000,
})

http.interceptors.request.use((config) => {
  const token = useAuthStore.getState().token
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

http.interceptors.response.use(
  (response) => {
    const payload = response.data as Envelope<unknown>
    if (payload && typeof payload === 'object' && 'code' in payload && 'message' in payload) {
      if (payload.code !== 0) {
        throw new Error(payload.message || '请求失败')
      }
      response.data = payload.data
    }
    return response
  },
  (error: AxiosError<Envelope<unknown>>) => {
    if (error.response?.status === 401) {
      useAuthStore.getState().logout()
      if (window.location.hash !== '#/login') {
        window.location.hash = '#/login'
      }
    }
    const message = error.response?.data?.message || error.message || '请求失败'
    return Promise.reject(new Error(message))
  },
)

export function resolveAssetUrl(path?: string) {
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

export const publicApi = {
  async site() {
    const { data } = await http.get<SiteContent>('/api/site')
    return data
  },
  async articles() {
    const { data } = await http.get<{ articles: ArticleSummary[] }>('/api/articles')
    return data.articles
  },
  async article(id: string) {
    const { data } = await http.get<Article>(`/api/articles/${encodeURIComponent(id)}`)
    return data
  },
}

export const adminApi = {
  async login(payload: { username: string; password: string }) {
    const { data } = await http.post<{ token: string; userInfo: User }>('/api/admin/login', payload)
    return data
  },
  async me() {
    const { data } = await http.get<User>('/api/admin/me')
    return data
  },
  async dashboard() {
    const { data } = await http.get<DashboardData>('/api/admin/dashboard')
    return data
  },
  async analytics() {
    const { data } = await http.get<AnalyticsData>('/api/admin/analytics')
    return data
  },
  async content() {
    const { data } = await http.get<Omit<SiteContent, 'featuredArticles'>>('/api/admin/content')
    return data
  },
  async articles(params: Record<string, unknown>) {
    const { data } = await http.get<PageResult<Article>>('/api/admin/articles', { params })
    return data
  },
  async article(id: string) {
    const { data } = await http.get<Article>(`/api/admin/articles/${encodeURIComponent(id)}`)
    return data
  },
  async saveArticle(article: Article, originalId?: string) {
    const path = originalId ? `/api/admin/articles/${encodeURIComponent(originalId)}` : '/api/admin/articles'
    const { data } = originalId ? await http.put<Article>(path, article) : await http.post<Article>(path, article)
    return data
  },
  async deleteArticle(id: string) {
    await http.delete(`/api/admin/articles/${encodeURIComponent(id)}`)
  },
  async publishArticle(id: string) {
    const { data } = await http.post<Article>(`/api/admin/articles/${encodeURIComponent(id)}/publish`)
    return data
  },
  async site() {
    const { data } = await http.get<SiteSettings>('/api/admin/site')
    return data
  },
  async saveSite(site: SiteSettings) {
    const { data } = await http.put<SiteSettings>('/api/admin/site', site)
    return data
  },
  async home() {
    const { data } = await http.get<HomeContent>('/api/admin/home')
    return data
  },
  async saveHome(home: HomeContent) {
    const { data } = await http.put<HomeContent>('/api/admin/home', home)
    return data
  },
  async about() {
    const { data } = await http.get<AboutContent>('/api/admin/about')
    return data
  },
  async saveAbout(about: Partial<AboutContent>) {
    const { data } = await http.put<AboutContent>('/api/admin/about', about)
    return data
  },
  async skills() {
    const { data } = await http.get<WorkStackGroup[]>('/api/admin/skills')
    return data
  },
  async saveSkills(skills: WorkStackGroup[]) {
    const { data } = await http.put<WorkStackGroup[]>('/api/admin/skills', skills)
    return data
  },
  async timeline() {
    const { data } = await http.get<ExperienceItem[]>('/api/admin/timeline')
    return data
  },
  async saveTimeline(items: ExperienceItem[]) {
    const { data } = await http.put<ExperienceItem[]>('/api/admin/timeline', items)
    return data
  },
  async tags() {
    const { data } = await http.get<PageResult<TagUsage>>('/api/admin/tags')
    return data
  },
  async saveTags(tags: string[]) {
    const { data } = await http.put<{ tags: string[] }>('/api/admin/tags', { tags })
    return data
  },
  async upload(file: File, kind = 'image') {
    const form = new FormData()
    form.append('file', file)
    form.append('kind', kind)
    const { data } = await http.post<{ url: string }>('/api/admin/upload', form)
    return data.url
  },
  async users(params: Record<string, unknown>) {
    const { data } = await http.get<PageResult<User>>('/api/admin/users', { params })
    return data
  },
  async createUser(payload: { username: string; nickname: string; email: string; password: string; role: string; status: UserStatus }) {
    const { data } = await http.post<User>('/api/admin/users', payload)
    return data
  },
  async updateUser(id: number, payload: { nickname: string; email: string; role: string; status: UserStatus }) {
    const { data } = await http.put<User>(`/api/admin/users/${id}`, payload)
    return data
  },
  async deleteUser(id: number) {
    await http.delete(`/api/admin/users/${id}`)
  },
  async resetPassword(id: number, password: string) {
    await http.post(`/api/admin/users/${id}/reset-password`, { password })
  },
  async setUserStatus(id: number, status: UserStatus) {
    const { data } = await http.post<User>(`/api/admin/users/${id}/status`, { status })
    return data
  },
  async loginLogs(params: Record<string, unknown>) {
    const { data } = await http.get<PageResult<LoginLog>>('/api/admin/login-logs', { params })
    return data
  },
}
