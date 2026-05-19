import { API_BASE_URL, adminApi, estimateWords, getArticleWordCount, getReadingMinutes, publicApi, resolveAssetUrl } from '../api/client'
import type {
  AboutContent,
  Article,
  ArticleSummary,
  ExperienceItem,
  HomeContent,
  SiteContent,
  SiteSettings,
  WorkStackGroup,
} from '../types'

export { API_BASE_URL, adminApi, estimateWords, getArticleWordCount, getReadingMinutes, resolveAssetUrl }
export type { AboutContent, Article, ArticleSummary, ExperienceItem, HomeContent, SiteContent, SiteSettings, WorkStackGroup }

export type AdminContent = Omit<SiteContent, 'featuredArticles'>

export async function fetchSiteContent(signal?: AbortSignal) {
  void signal
  return publicApi.site()
}

export async function fetchArticles(signal?: AbortSignal) {
  void signal
  return publicApi.articles()
}

export async function fetchArticle(id: string, signal?: AbortSignal) {
  void signal
  return publicApi.article(id)
}

export async function fetchAdminContent(token: string, signal?: AbortSignal) {
  void token
  void signal
  return adminApi.content()
}

export async function saveAdminSite(token: string, site: SiteSettings) {
  void token
  return adminApi.saveSite(site)
}

export async function saveAdminHome(token: string, home: HomeContent) {
  void token
  return adminApi.saveHome(home)
}

export async function saveAdminAbout(token: string, about: AboutContent) {
  void token
  return adminApi.saveAbout(about)
}

export async function saveAdminTags(token: string, tags: string[]) {
  void token
  return adminApi.saveTags(tags)
}

export async function saveAdminArticle(token: string, article: Article, originalId?: string) {
  void token
  return adminApi.saveArticle(article, originalId)
}

export async function deleteAdminArticle(token: string, id: string) {
  void token
  return adminApi.deleteArticle(id)
}

export async function uploadAdminLogo(token: string, file: File) {
  void token
  const logoUrl = await adminApi.upload(file, 'logo')
  const site = await adminApi.site()
  return { logoUrl, site: { ...site, logoUrl } }
}

