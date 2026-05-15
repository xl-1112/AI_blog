import { Button } from './Button'
import { Card } from './Card'

type PageStateProps = {
  title: string
  description?: string
  actionLabel?: string
  onAction?: () => void
}

export function PageLoading({ title = '正在加载内容...' }: Partial<PageStateProps>) {
  return (
    <main className="px-5 py-20">
      <Card className="mx-auto max-w-xl p-8 text-center">
        <div className="mx-auto mb-5 h-10 w-10 animate-pulse rounded-full bg-[linear-gradient(135deg,var(--accent),var(--accent-secondary))]" />
        <p className="font-medium text-[var(--foreground)]">{title}</p>
      </Card>
    </main>
  )
}

export function PageError({ title, description, actionLabel = '重新加载', onAction }: PageStateProps) {
  return (
    <main className="px-5 py-20">
      <Card className="mx-auto max-w-xl p-8 text-center">
        <h1 className="text-2xl font-semibold text-[var(--foreground)]">{title}</h1>
        {description ? <p className="mt-3 leading-7 text-[var(--muted-foreground)]">{description}</p> : null}
        {onAction ? (
          <Button type="button" className="mt-6" onClick={onAction}>
            {actionLabel}
          </Button>
        ) : null}
      </Card>
    </main>
  )
}
