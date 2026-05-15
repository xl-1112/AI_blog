import type { AnchorHTMLAttributes, ButtonHTMLAttributes, ReactNode } from 'react'
import { cva, type VariantProps } from 'class-variance-authority'
import { ArrowRight } from 'lucide-react'
import { cn } from '../lib/utils'

const buttonVariants = cva(
  'group inline-flex min-h-12 items-center justify-center gap-2 rounded-xl px-5 py-3 text-sm font-semibold transition-all duration-200 ease-out focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--ring)] focus-visible:ring-offset-2 focus-visible:ring-offset-[var(--background)] active:scale-[0.98]',
  {
    variants: {
      variant: {
        primary:
          'bg-[linear-gradient(135deg,var(--accent),var(--accent-secondary))] text-white shadow-[0_4px_14px_rgba(0,82,255,0.25)] hover:-translate-y-0.5 hover:brightness-110 hover:shadow-[0_8px_24px_rgba(0,82,255,0.35)]',
        secondary:
          'border border-[var(--border)] bg-white/70 text-[var(--foreground)] shadow-sm hover:-translate-y-0.5 hover:border-blue-300 hover:bg-[var(--muted)] hover:shadow-md',
        ghost:
          'text-[var(--muted-foreground)] hover:bg-[var(--muted)] hover:text-[var(--foreground)]',
      },
      size: {
        default: 'min-h-12',
        sm: 'min-h-10 px-4 py-2',
      },
    },
    defaultVariants: {
      variant: 'primary',
      size: 'default',
    },
  },
)

type ButtonBaseProps = VariantProps<typeof buttonVariants> & {
  children: ReactNode
  className?: string
  showArrow?: boolean
}

type ButtonProps = ButtonBaseProps & ButtonHTMLAttributes<HTMLButtonElement>
type ButtonLinkProps = ButtonBaseProps & AnchorHTMLAttributes<HTMLAnchorElement>

export function Button({ children, className, variant, size, showArrow = false, ...props }: ButtonProps) {
  return (
    <button className={cn(buttonVariants({ variant, size }), className)} {...props}>
      {children}
      {showArrow ? <ArrowRight className="h-4 w-4 transition-transform group-hover:translate-x-1" /> : null}
    </button>
  )
}

export function ButtonLink({ children, className, variant, size, showArrow = false, ...props }: ButtonLinkProps) {
  return (
    <a className={cn(buttonVariants({ variant, size }), className)} {...props}>
      {children}
      {showArrow ? <ArrowRight className="h-4 w-4 transition-transform group-hover:translate-x-1" /> : null}
    </a>
  )
}
