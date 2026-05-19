import { Link } from 'react-router-dom'
import { ArrowRight, Sparkles } from 'lucide-react'
import { ArticleCard } from '../components/ArticleCard'
import { ButtonLink } from '../components/Button'
import { Card, GradientCard } from '../components/Card'
import { HeroGraphic } from '../components/HeroGraphic'
import { PageError, PageLoading } from '../components/PageState'
import { SectionLabel } from '../components/SectionLabel'
import { useSiteContent } from '../lib/useSiteContent'

export function HomePage() {
  const { content, loading, error, reload } = useSiteContent()

  if (loading && !content) {
    return <PageLoading title="正在从 Go 服务加载首页内容..." />
  }

  if (error || !content) {
    return <PageError title="首页内容加载失败" description={error ?? '请确认 Go 服务已经启动。'} onAction={reload} />
  }

  const highlighted = content.featuredArticles[0] ?? content.articles[0]
  const recentArticles = content.articles.slice(0, 2)
  const productArticleCount = content.articles.filter((article) => article.tags.some((tag) => tag.includes('产品'))).length
  const techArticleCount = content.articles.filter((article) => article.tags.some((tag) => tag.includes('技术'))).length
  const metrics = [
    { label: '文章总数', value: String(content.articles.length).padStart(2, '0') },
    { label: '产品文章', value: String(productArticleCount).padStart(2, '0') },
    { label: '技术文章', value: techArticleCount > 0 ? String(techArticleCount).padStart(2, '0') : '偶尔' },
  ]

  return (
    <main>
      <section className="relative overflow-hidden px-5 py-16 sm:py-20 lg:py-24">
        <div className="pointer-events-none absolute left-[-12rem] top-8 h-96 w-96 rounded-full bg-blue-500/10 blur-[120px]" />
        <div className="pointer-events-none absolute right-[-10rem] top-32 h-96 w-96 rounded-full bg-sky-400/10 blur-[120px]" />
        <div className="mx-auto grid max-w-6xl items-center gap-12 lg:grid-cols-[1.1fr_0.9fr]">
          <div className="relative z-10 space-y-9">
            <SectionLabel pulse>产品笔记</SectionLabel>
            <div className="space-y-6">
              <h1 className="max-w-3xl font-display text-[3.2rem] leading-[1.04] tracking-normal text-[var(--foreground)] sm:text-6xl lg:text-[5.2rem]">
                {content.home.introTitle}
                <span className="relative ml-3 inline-block">
                  <span className="gradient-text">{content.site.role}</span>
                  <span className="absolute -bottom-1 left-0 h-3 w-full rounded bg-blue-500/10 sm:-bottom-2 sm:h-4" />
                </span>
              </h1>
              <p className="max-w-2xl text-lg leading-8 text-[var(--muted-foreground)]">{content.home.intro}</p>
            </div>
            <div className="flex flex-col gap-3 sm:flex-row">
              <ButtonLink href="#/site/articles" showArrow>
                {content.home.primaryCtaText || '阅读文章'}
              </ButtonLink>
              <ButtonLink href="#/site/about" variant="secondary">
                {content.home.secondaryCtaText || '关于我'}
              </ButtonLink>
            </div>
            <div className="grid max-w-2xl grid-cols-3 gap-3">
              {metrics.map((metric) => (
                <Card key={metric.label} className="rounded-xl p-4">
                  <div className="font-display text-2xl text-[var(--foreground)]">{metric.value}</div>
                  <div className="mt-1 text-xs text-[var(--muted-foreground)]">{metric.label}</div>
                </Card>
              ))}
            </div>
          </div>
          <HeroGraphic />
        </div>
      </section>

      <section className="px-5 pb-20">
        <div className="mx-auto grid max-w-6xl gap-6 lg:grid-cols-[0.8fr_1.2fr]">
          {highlighted ? (
            <GradientCard>
              <div className="space-y-5 p-7">
                <div className="inline-flex h-12 w-12 items-center justify-center rounded-2xl bg-[linear-gradient(135deg,var(--accent),var(--accent-secondary))] text-white shadow-[0_8px_24px_rgba(0,82,255,0.32)]">
                  <Sparkles className="h-5 w-5" />
                </div>
                <div>
                  <p className="font-mono text-xs tracking-normal text-[var(--accent)]">精选文章</p>
                  <h2 className="mt-3 text-2xl font-semibold tracking-normal text-[var(--foreground)]">{highlighted.title}</h2>
                  <p className="mt-3 text-sm leading-7 text-[var(--muted-foreground)]">{highlighted.summary}</p>
                </div>
                <Link
                  to={`/site/articles/${highlighted.slug || highlighted.id}`}
                  className="inline-flex min-h-11 items-center gap-2 text-sm font-semibold text-[var(--accent)] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--ring)]"
                >
                  查看精选文章
                  <ArrowRight className="h-4 w-4" />
                </Link>
              </div>
            </GradientCard>
          ) : null}

          <div className="grid gap-4 md:grid-cols-2">
            {recentArticles.map((article) => (
              <ArticleCard key={article.id} article={article} featured={article.featured} />
            ))}
          </div>
        </div>
      </section>
    </main>
  )
}
