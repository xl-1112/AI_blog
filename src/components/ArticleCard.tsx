import { Link } from 'react-router-dom'
import { Calendar, Clock, MoveUpRight } from 'lucide-react'
import { getReadingMinutes, resolveAssetUrl, type ArticleSummary } from '../lib/api'
import { cn } from '../lib/utils'

type ArticleCardProps = {
  article: ArticleSummary
  featured?: boolean
  className?: string
}

export function ArticleCard({ article, featured = false, className }: ArticleCardProps) {
  return (
    <Link
      to={`/site/articles/${article.slug || article.id}`}
      className={cn(
        'group block rounded-2xl border border-[var(--border)] bg-[var(--card)] p-6 shadow-[0_10px_30px_rgba(15,23,42,0.05)] transition-all duration-300 hover:-translate-y-1 hover:border-blue-300 hover:shadow-[0_20px_45px_rgba(15,23,42,0.1)]',
        featured && 'relative overflow-hidden',
        className,
      )}
    >
      {featured ? (
        <span className="absolute inset-x-0 top-0 h-1 bg-[linear-gradient(90deg,var(--accent),var(--accent-secondary))]" />
      ) : null}
      <article className="space-y-4">
        {article.coverUrl ? (
          <img
            src={resolveAssetUrl(article.coverUrl)}
            alt={`${article.title} 封面`}
            className="aspect-[16/9] w-full rounded-xl object-cover"
            loading="lazy"
          />
        ) : null}
        <div className="flex flex-wrap items-center gap-2">
          {article.category ? (
            <span className="rounded-full border border-slate-200 bg-slate-50 px-3 py-1 text-xs font-medium tracking-normal text-slate-600">
              {article.category}
            </span>
          ) : null}
        </div>
        <div className="space-y-2">
          <div className="flex items-start justify-between gap-4">
            <h2 className="text-xl font-semibold leading-snug tracking-normal text-[var(--foreground)] transition-colors group-hover:text-[var(--accent)]">
              {article.title}
            </h2>
            <MoveUpRight className="mt-1 h-5 w-5 flex-none text-[var(--muted-foreground)] transition-all group-hover:translate-x-0.5 group-hover:-translate-y-0.5 group-hover:text-[var(--accent)]" />
          </div>
          <p className="line-clamp-3 text-sm leading-7 text-[var(--muted-foreground)]">{article.summary}</p>
        </div>
        <div className="flex flex-wrap gap-4 text-xs text-[var(--muted-foreground)]">
          <span className="inline-flex items-center gap-1.5">
            <Calendar className="h-3.5 w-3.5" />
            {article.date}
          </span>
          <span className="inline-flex items-center gap-1.5">
            <Clock className="h-3.5 w-3.5" />
            {getReadingMinutes(article)} 分钟阅读
          </span>
          <span>{article.viewCount ?? 0} 次浏览</span>
        </div>
      </article>
    </Link>
  )
}
