import { cn } from '../lib/utils'

type SectionLabelProps = {
  children: string
  className?: string
  pulse?: boolean
}

export function SectionLabel({ children, className, pulse = false }: SectionLabelProps) {
  return (
    <div
      className={cn(
        'inline-flex items-center gap-3 rounded-full border border-blue-500/30 bg-blue-500/5 px-5 py-2',
        className,
      )}
    >
      <span className={cn('h-2 w-2 rounded-full bg-[var(--accent)]', pulse && 'animate-pulse-dot')} />
      <span className="font-mono text-xs uppercase tracking-[0.15em] text-[var(--accent)]">{children}</span>
    </div>
  )
}
