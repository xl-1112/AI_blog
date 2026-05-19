import { useEffect, useState } from 'react'
import { Link, Navigate, useParams } from 'react-router-dom'
import { ArrowLeft, Calendar, Clock, Eye } from 'lucide-react'
import { MarkdownContent } from '../components/MarkdownContent'
import { PageError, PageLoading } from '../components/PageState'
import { fetchArticle, getArticleWordCount, getReadingMinutes, resolveAssetUrl, type Article } from '../lib/api'
import { useSiteContent } from '../lib/useSiteContent'

export function ArticleDetailPage() {
  const { id } = useParams()
  const { content } = useSiteContent()
  const [article, setArticle] = useState<Article | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const articleId = id ?? ''
    if (!articleId) return

    const controller = new AbortController()

    async function loadArticle() {
      setLoading(true)
      setError(null)
      try {
        setArticle(await fetchArticle(articleId, controller.signal))
      } catch (err) {
        if (controller.signal.aborted) return
        setError(err instanceof Error ? err.message : '文章加载失败')
      } finally {
        if (!controller.signal.aborted) {
          setLoading(false)
        }
      }
    }

    void loadArticle()

    return () => controller.abort()
  }, [id])

  useEffect(() => {
    if (!article) return
    document.title = article.seoTitle || `${article.title} | Liang Blog`
    const description = article.seoDescription || article.summary
    let meta = document.querySelector('meta[name="description"]')
    if (!meta) {
      meta = document.createElement('meta')
      meta.setAttribute('name', 'description')
      document.head.appendChild(meta)
    }
    meta.setAttribute('content', description)
  }, [article])

  if (!id) {
    return <Navigate to="/site/articles" replace />
  }

  if (loading) {
    return <PageLoading title="正在从 Go 服务加载文章详情..." />
  }

  if (error || !article) {
    return (
      <PageError
        title="文章加载失败"
        description={error ?? '没有找到这篇文章。'}
        actionLabel="返回文章列表"
        onAction={() => {
          window.location.hash = '#/site/articles'
        }}
      />
    )
  }

  return (
    <main className="px-5 py-12 sm:py-16">
      <article className="mx-auto max-w-3xl">
        <Link
          to="/site/articles"
          className="mb-10 inline-flex min-h-11 items-center gap-2 rounded-full text-sm font-semibold text-[var(--muted-foreground)] transition-colors hover:text-[var(--accent)] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--ring)]"
        >
          <ArrowLeft className="h-4 w-4" />
          返回文章列表
        </Link>

        <header className="mb-10 space-y-5">
          <div className="flex flex-wrap gap-2">
            {[article.category].filter(Boolean).map((tag) => (
              <span
                key={tag}
                className="rounded-full border border-blue-500/20 bg-blue-500/5 px-3 py-1 text-xs font-medium tracking-normal text-[var(--accent)]"
              >
                {tag}
              </span>
            ))}
          </div>
          <h1 className="font-display text-5xl leading-tight text-[var(--foreground)] sm:text-6xl">{article.title}</h1>
          <p className="text-lg leading-8 text-[var(--muted-foreground)]">{article.summary}</p>
          {article.coverUrl ? (
            <img src={resolveAssetUrl(article.coverUrl)} alt={`${article.title} 封面`} className="aspect-[16/9] w-full rounded-2xl object-cover" />
          ) : null}
          <div className="flex flex-wrap gap-4 border-y border-[var(--border)] py-4 text-sm text-[var(--muted-foreground)]">
            <span className="inline-flex items-center gap-2">
              <Calendar className="h-4 w-4" />
              {article.date}
            </span>
            <span className="inline-flex items-center gap-2">
              <Clock className="h-4 w-4" />
              {getReadingMinutes(article)} 分钟阅读
            </span>
            <span>{getArticleWordCount(article).toLocaleString()} 字</span>
            <span className="inline-flex items-center gap-2">
              <Eye className="h-4 w-4" />
              {article.viewCount ?? 0} 次浏览
            </span>
          </div>
        </header>

        <MarkdownContent content={article.content} />

        <footer className="mt-14 flex items-center justify-between border-t border-[var(--border)] pt-6 text-sm text-[var(--muted-foreground)]">
          <Link to="/site/articles" className="inline-flex items-center gap-2 transition-colors hover:text-[var(--accent)]">
            <ArrowLeft className="h-4 w-4" />
            返回列表
          </Link>
          <span>{content?.site.name ?? 'Liang'} · Blog</span>
        </footer>
      </article>
    </main>
  )
}
