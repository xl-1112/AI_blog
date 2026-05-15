import { useEffect, useState } from 'react'
import { Link, Navigate, useParams } from 'react-router-dom'
import { ArrowLeft, Calendar, Clock } from 'lucide-react'
import { MarkdownContent } from '../components/MarkdownContent'
import { PageError, PageLoading } from '../components/PageState'
import { fetchArticle, getArticleWordCount, getReadingMinutes, type Article } from '../lib/api'
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

  if (!id) {
    return <Navigate to="/articles" replace />
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
          window.location.hash = '#/articles'
        }}
      />
    )
  }

  return (
    <main className="px-5 py-12 sm:py-16">
      <article className="mx-auto max-w-3xl">
        <Link
          to="/articles"
          className="mb-10 inline-flex min-h-11 items-center gap-2 rounded-full text-sm font-semibold text-[var(--muted-foreground)] transition-colors hover:text-[var(--accent)] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--ring)]"
        >
          <ArrowLeft className="h-4 w-4" />
          返回文章列表
        </Link>

        <header className="mb-10 space-y-5">
          <div className="flex flex-wrap gap-2">
            {article.tags.map((tag) => (
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
          </div>
        </header>

        <MarkdownContent content={article.content} />

        <footer className="mt-14 flex items-center justify-between border-t border-[var(--border)] pt-6 text-sm text-[var(--muted-foreground)]">
          <Link to="/articles" className="inline-flex items-center gap-2 transition-colors hover:text-[var(--accent)]">
            <ArrowLeft className="h-4 w-4" />
            返回列表
          </Link>
          <span>{content?.site.name ?? 'Liang'} · Blog</span>
        </footer>
      </article>
    </main>
  )
}
