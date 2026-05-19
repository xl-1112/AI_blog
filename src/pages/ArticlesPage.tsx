import { useEffect, useMemo, useState } from 'react'
import { Search } from 'lucide-react'
import { ArticleCard } from '../components/ArticleCard'
import { PageError, PageLoading } from '../components/PageState'
import { SectionLabel } from '../components/SectionLabel'
import { fetchArticles, type ArticleSummary } from '../lib/api'
import { cn } from '../lib/utils'
import { useSiteContent } from '../lib/useSiteContent'

export function ArticlesPage() {
  const { content, loading: contentLoading, error: contentError, reload } = useSiteContent()
  const [query, setQuery] = useState('')
  const [activeTag, setActiveTag] = useState('全部')
  const [articles, setArticles] = useState<ArticleSummary[]>([])
  const [articlesLoading, setArticlesLoading] = useState(true)
  const [articlesError, setArticlesError] = useState<string | null>(null)

  useEffect(() => {
    const controller = new AbortController()

    async function loadArticles() {
      setArticlesLoading(true)
      setArticlesError(null)
      try {
        setArticles(await fetchArticles(controller.signal))
      } catch (err) {
        if (controller.signal.aborted) return
        setArticlesError(err instanceof Error ? err.message : '文章加载失败')
      } finally {
        if (!controller.signal.aborted) {
          setArticlesLoading(false)
        }
      }
    }

    void loadArticles()

    return () => controller.abort()
  }, [])

  const tags = useMemo(() => content?.tags ?? [], [content?.tags])

  useEffect(() => {
    if (activeTag !== '全部' && !tags.includes(activeTag)) {
      setActiveTag('全部')
    }
  }, [activeTag, tags])

  const filteredArticles = useMemo(() => {
    const keyword = query.trim().toLowerCase()

    return articles.filter((article) => {
      const matchesTag = activeTag === '全部' || article.category === activeTag
      const matchesQuery =
        !keyword ||
        article.title.toLowerCase().includes(keyword) ||
        article.summary.toLowerCase().includes(keyword) ||
        article.category.toLowerCase().includes(keyword)

      return matchesTag && matchesQuery
    })
  }, [activeTag, articles, query])

  if ((contentLoading && !content) || articlesLoading) {
    return <PageLoading title="正在从 Go 服务加载文章..." />
  }

  if (contentError || articlesError || !content) {
    return (
      <PageError
        title="文章内容加载失败"
        description={contentError ?? articlesError ?? '请确认 Go 服务已经启动。'}
        onAction={() => {
          reload()
          window.location.reload()
        }}
      />
    )
  }

  return (
    <main className="px-5 py-14 sm:py-20">
      <div className="mx-auto max-w-5xl">
        <div className="mb-10 space-y-5">
          <SectionLabel pulse>文章</SectionLabel>
          <div className="flex flex-col justify-between gap-5 md:flex-row md:items-end">
            <div>
              <h1 className="font-display text-5xl leading-tight text-[var(--foreground)]">所有文章</h1>
              <p className="mt-4 max-w-2xl leading-8 text-[var(--muted-foreground)]">
                按标题、摘要和标签搜索。这里主要记录产品思考、项目复盘和方法论，也会保留少量技术相关文章。
              </p>
            </div>
            <span className="font-mono text-sm tracking-normal text-[var(--muted-foreground)]">
              {filteredArticles.length} / {articles.length}
            </span>
          </div>
        </div>

        <div className="mb-8 rounded-2xl border border-[var(--border)] bg-white/80 p-3 shadow-sm backdrop-blur">
          <label className="relative block">
            <Search className="pointer-events-none absolute left-4 top-1/2 h-4 w-4 -translate-y-1/2 text-[var(--muted-foreground)]" />
            <input
              value={query}
              onChange={(event) => setQuery(event.target.value)}
              placeholder="搜索文章标题、摘要或标签..."
              className="h-12 w-full rounded-xl border border-transparent bg-[var(--muted)] pl-11 pr-4 text-sm text-[var(--foreground)] outline-none transition focus:border-blue-300 focus:bg-white focus:ring-2 focus:ring-blue-500/20"
            />
          </label>
          <div className="mt-3 flex flex-wrap gap-2">
            {['全部', ...tags].map((tag) => (
              <button
                key={tag}
                type="button"
                onClick={() => setActiveTag(tag)}
                className={cn(
                  'min-h-10 rounded-full border px-4 text-sm font-medium tracking-normal transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--ring)]',
                  activeTag === tag
                    ? 'border-blue-500/20 bg-blue-500/10 text-[var(--accent)] shadow-[0_6px_18px_rgba(0,82,255,0.12)]'
                    : 'border-[var(--border)] bg-white text-[var(--muted-foreground)] hover:border-blue-300 hover:text-[var(--accent)]',
                )}
              >
                {tag}
              </button>
            ))}
          </div>
        </div>

        {filteredArticles.length > 0 ? (
          <div className="grid gap-5 md:grid-cols-2">
            {filteredArticles.map((article) => (
              <ArticleCard key={article.id} article={article} featured={article.featured} />
            ))}
          </div>
        ) : (
          <div className="rounded-2xl border border-dashed border-[var(--border)] bg-white/60 px-6 py-16 text-center text-[var(--muted-foreground)]">
            没有找到匹配的文章。
          </div>
        )}
      </div>
    </main>
  )
}
