export type Role = 'super_admin' | 'admin' | 'editor'
export type UserStatus = 'active' | 'disabled'
export type ArticleStatus = 'draft' | 'published'

export type Contact = {
  email: string
  github: string
}

export type SiteSettings = {
  name: string
  siteTitle: string
  description: string
  keywords: string
  logoUrl: string
  faviconUrl: string
  icp: string
  analyticsCode: string
  role: string
  location: string
  contact: Contact
}

export type HomeContent = {
  introTitle: string
  subtitle: string
  intro: string
  primaryCtaText: string
  secondaryCtaText: string
  featuredArticleIds: string[]
}

export type WorkStackGroup = {
  id?: number
  title: string
  items: string[]
  sort?: number
}

export type ExperienceItem = {
  id?: number
  period: string
  title: string
  body: string
  sort?: number
}

export type AboutContent = {
  name: string
  title: string
  avatarUrl: string
  showHero: boolean
  showHeroBadge: boolean
  showHeroTitle: boolean
  showHeroSubtitle: boolean
  showHeroDescription: boolean
  heroBadge: string
  heroTitle: string
  heroSubtitle: string
  description: string
  richDescription: string
  showContact: boolean
  showLocation: boolean
  showEmail: boolean
  showGithub: boolean
  showContactBadge: boolean
  showContactTitle: boolean
  showContactSubtitle: boolean
  showContactDescription: boolean
  contactBadge: string
  contactTitle: string
  contactSubtitle: string
  contactDescription: string
  showSkills: boolean
  showSkillsHeader: boolean
  showSkillsBadge: boolean
  showSkillsTitle: boolean
  showSkillsSubtitle: boolean
  showSkillsDescription: boolean
  skillsBadge: string
  skillsTitle: string
  skillsSubtitle: string
  skillsDescription: string
  showTimeline: boolean
  showTimelineHeader: boolean
  showTimelineBadge: boolean
  showTimelineTitle: boolean
  showTimelineSubtitle: boolean
  showTimelineDescription: boolean
  timelineBadge: string
  timelineTitle: string
  timelineSubtitle: string
  timelineDescription: string
  city: string
  contact: Contact
  workStack: WorkStackGroup[]
  experience: ExperienceItem[]
}

export type ArticleSummary = {
  id: string
  title: string
  slug: string
  category: string
  date: string
  summary: string
  coverUrl: string
  tags: string[]
  featured: boolean
  draft?: boolean
  status: ArticleStatus
  wordCount: number
  viewCount: number
}

export type Article = ArticleSummary & {
  content: string
  seoTitle: string
  seoKeywords: string
  seoDescription: string
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

export type User = {
  id: number
  username: string
  nickname: string
  email: string
  role: Role
  status: UserStatus
  lastLoginAt?: string
  createdAt: string
  updatedAt: string
}

export type LoginLog = {
  id: number
  userId?: number
  username: string
  ip: string
  userAgent: string
  success: boolean
  reason: string
  createdAt: string
}

export type TrendPoint = {
  date: string
  views: number
}

export type DashboardData = {
  totalArticles: number
  totalTags: number
  totalViews: number
  todayViews: number
  recentArticles: ArticleSummary[]
  viewTrend: TrendPoint[]
}

export type AnalyticsData = {
  viewTrend: TrendPoint[]
  hotArticles: ArticleSummary[]
}

export type PageResult<T> = {
  list: T[]
  total: number
}

export type TagUsage = {
  name: string
  useCount: number
  createdAt: string
}
