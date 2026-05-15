import type { HTMLAttributes } from 'react'
import { cn } from '../lib/utils'

export function Card({ className, ...props }: HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn(
        'rounded-2xl border border-[var(--border)] bg-[var(--card)] shadow-[0_10px_30px_rgba(15,23,42,0.06)] transition-all duration-300 hover:-translate-y-1 hover:shadow-[0_20px_45px_rgba(15,23,42,0.1)]',
        className,
      )}
      {...props}
    />
  )
}

export function GradientCard({ className, children, ...props }: HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn('rounded-2xl bg-[linear-gradient(135deg,var(--accent),var(--accent-secondary))] p-[2px]', className)}
      {...props}
    >
      <div className="h-full rounded-[calc(1rem-2px)] bg-[var(--card)]">{children}</div>
    </div>
  )
}
