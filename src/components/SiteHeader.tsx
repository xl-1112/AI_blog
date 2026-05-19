import { NavLink } from 'react-router-dom'
import { BookOpen, Home, UserRound } from 'lucide-react'
import { resolveAssetUrl } from '../lib/api'
import { cn } from '../lib/utils'
import { useSiteContent } from '../lib/useSiteContent'

const links = [
  { to: '/site', label: '首页', icon: Home },
  { to: '/site/articles', label: '文章', icon: BookOpen },
  { to: '/site/about', label: '关于我', icon: UserRound },
]

export function SiteHeader() {
  const { content } = useSiteContent()
  const logoUrl = resolveAssetUrl(content?.site.logoUrl)
  const siteName = content?.site.name ?? 'Liang'

  return (
    <header className="sticky top-0 z-40 border-b border-[var(--border)] bg-[color-mix(in_srgb,var(--background)_86%,transparent)] backdrop-blur-xl">
      <div className="mx-auto flex max-w-6xl items-center justify-between gap-4 px-5 py-4">
        <NavLink to="/site" className="inline-flex min-w-0 items-center focus-visible:outline-none">
          <span className="flex h-12 w-36 shrink-0 items-center justify-start overflow-visible">
            {logoUrl ? (
              <img src={logoUrl} alt={`${siteName} logo`} className="h-full w-full object-contain object-left" />
            ) : (
              <span className="font-display text-2xl text-[var(--foreground)]">{siteName}</span>
            )}
          </span>
        </NavLink>

        <nav className="flex items-center gap-1 rounded-full border border-[var(--border)] bg-white/75 p-1 shadow-sm">
          {links.map((link) => {
            const Icon = link.icon
            return (
              <NavLink
                key={link.to}
                to={link.to}
                end={link.to === '/site'}
                className={({ isActive }) =>
                  cn(
                    'inline-flex min-h-10 items-center gap-2 rounded-full border border-transparent px-3 text-sm font-medium text-[var(--muted-foreground)] transition-all hover:bg-[var(--muted)] hover:text-[var(--foreground)] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--ring)] sm:px-4',
                    isActive &&
                      'border-blue-500/20 bg-blue-500/10 text-[var(--accent)] shadow-[0_6px_18px_rgba(0,82,255,0.12)] hover:bg-blue-500/15 hover:text-[var(--accent)]',
                  )
                }
              >
                <Icon className="h-4 w-4" />
                <span>{link.label}</span>
              </NavLink>
            )
          })}
        </nav>
      </div>
    </header>
  )
}
