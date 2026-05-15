import { Mail, MapPin, MoveUpRight } from 'lucide-react'
import { Card } from '../components/Card'
import { PageError, PageLoading } from '../components/PageState'
import { SectionLabel } from '../components/SectionLabel'
import { useSiteContent } from '../lib/useSiteContent'

function GitHubIcon({ className }: { className?: string }) {
  return (
    <svg viewBox="0 0 24 24" aria-hidden="true" className={className} fill="currentColor">
      <path d="M12 2C6.48 2 2 6.58 2 12.26c0 4.52 2.87 8.36 6.84 9.72.5.1.68-.22.68-.49 0-.24-.01-1.04-.01-1.89-2.78.62-3.37-1.22-3.37-1.22-.45-1.19-1.11-1.5-1.11-1.5-.91-.64.07-.63.07-.63 1.01.07 1.54 1.06 1.54 1.06.89 1.57 2.34 1.12 2.91.85.09-.66.35-1.12.63-1.37-2.22-.26-4.55-1.14-4.55-5.07 0-1.12.39-2.03 1.03-2.75-.1-.26-.45-1.3.1-2.71 0 0 .84-.28 2.75 1.05A9.3 9.3 0 0 1 12 6.99c.85 0 1.7.12 2.5.34 1.9-1.33 2.74-1.05 2.74-1.05.55 1.41.2 2.45.1 2.71.64.72 1.03 1.63 1.03 2.75 0 3.94-2.34 4.8-4.57 5.06.36.32.68.94.68 1.9 0 1.37-.01 2.47-.01 2.8 0 .27.18.6.69.49A10.14 10.14 0 0 0 22 12.26C22 6.58 17.52 2 12 2Z" />
    </svg>
  )
}

export function AboutPage() {
  const { content, loading, error, reload } = useSiteContent()

  if (loading && !content) {
    return <PageLoading title="正在从 Go 服务加载关于我..." />
  }

  if (error || !content) {
    return <PageError title="关于我加载失败" description={error ?? '请确认 Go 服务已经启动。'} onAction={reload} />
  }

  const { site, about } = content

  return (
    <main className="px-5 py-14 sm:py-20">
      <div className="mx-auto max-w-6xl space-y-14">
        <section className="grid gap-10 lg:grid-cols-[0.9fr_1.1fr]">
          <div className="space-y-7">
            <SectionLabel pulse>关于我</SectionLabel>
            <div className="space-y-5">
              <h1 className="font-display text-5xl leading-tight text-[var(--foreground)] sm:text-6xl">
                {site.role} {site.name}
              </h1>
              <p className="text-lg leading-8 text-[var(--muted-foreground)]">{site.description}</p>
              <p className="leading-8 text-[var(--muted-foreground)]">{about.description}</p>
            </div>
          </div>

          <Card className="p-7">
            <div className="mb-6">
              <p className="font-mono text-xs tracking-normal text-[var(--accent)]">联系我</p>
              <h2 className="mt-2 text-2xl font-semibold tracking-normal text-[var(--foreground)]">欢迎交流产品、增长和项目协作</h2>
            </div>
            <div className="grid gap-5 sm:grid-cols-2">
              <div className="rounded-2xl bg-[var(--muted)] p-5">
                <MapPin className="mb-4 h-5 w-5 text-[var(--accent)]" />
                <p className="text-sm text-[var(--muted-foreground)]">位置</p>
                <p className="mt-1 font-semibold text-[var(--foreground)]">{site.location}</p>
              </div>
              <a
                href={`mailto:${site.contact.email}`}
                className="group rounded-2xl bg-[var(--muted)] p-5 transition-all hover:-translate-y-0.5 hover:shadow-md focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--ring)]"
              >
                <Mail className="mb-4 h-5 w-5 text-[var(--accent)]" />
                <p className="text-sm text-[var(--muted-foreground)]">邮箱</p>
                <p className="mt-1 break-all font-semibold text-[var(--foreground)] group-hover:text-[var(--accent)]">
                  {site.contact.email}
                </p>
              </a>
              <a
                href={site.contact.github}
                target="_blank"
                rel="noreferrer"
                className="group rounded-2xl bg-[var(--muted)] p-5 transition-all hover:-translate-y-0.5 hover:shadow-md focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--ring)] sm:col-span-2"
              >
                <GitHubIcon className="mb-4 h-5 w-5 text-[var(--accent)]" />
                <p className="text-sm text-[var(--muted-foreground)]">GitHub</p>
                <p className="mt-1 inline-flex items-center gap-2 break-all font-semibold text-[var(--foreground)] group-hover:text-[var(--accent)]">
                  {site.contact.github}
                  <MoveUpRight className="h-4 w-4 flex-none" />
                </p>
              </a>
            </div>
          </Card>
        </section>

        <section className="space-y-6">
          <div className="space-y-4">
            <SectionLabel>技术栈</SectionLabel>
            <h2 className="font-display text-4xl text-[var(--foreground)] sm:text-5xl">产品经理的工作栈</h2>
            <p className="max-w-2xl leading-8 text-[var(--muted-foreground)]">
              以产品方法为核心，结合设计协作、数据分析和基础技术理解，保证需求判断和落地过程都更稳。
            </p>
          </div>
          <div className="grid gap-5 lg:grid-cols-3">
            {about.workStack.map((group) => (
              <Card key={group.title} className="p-6">
                <h3 className="text-xl font-semibold tracking-normal text-[var(--foreground)]">{group.title}</h3>
                <div className="mt-5 flex flex-wrap gap-2">
                  {group.items.map((item) => (
                    <span
                      key={item}
                      className="rounded-full border border-blue-500/20 bg-blue-500/5 px-3 py-1 text-xs font-medium text-[var(--accent)]"
                    >
                      {item}
                    </span>
                  ))}
                </div>
              </Card>
            ))}
          </div>
        </section>

        <section className="space-y-6">
          <div className="space-y-4">
            <SectionLabel>路线图</SectionLabel>
            <h2 className="font-display text-4xl text-[var(--foreground)] sm:text-5xl">工作经历</h2>
          </div>
          <div className="relative rounded-2xl border border-[var(--border)] bg-white/75 p-6 shadow-[0_10px_30px_rgba(15,23,42,0.05)]">
            <div className="absolute bottom-8 left-8 top-8 hidden w-px bg-[var(--border)] sm:block" />
            <div className="space-y-5">
              {about.experience.map((item, index) => (
                <div key={`${item.period}-${item.title}`} className="relative grid gap-4 sm:grid-cols-[6rem_1fr] sm:pl-8">
                  <div className="flex items-center gap-3 sm:block">
                    <span className="relative z-10 grid h-9 w-9 place-items-center rounded-full bg-[linear-gradient(135deg,var(--accent),var(--accent-secondary))] font-mono text-xs text-white shadow-[0_8px_24px_rgba(0,82,255,0.28)] sm:absolute sm:left-[-1.1rem] sm:top-1">
                      {index + 1}
                    </span>
                    <span className="font-mono text-xs tracking-normal text-[var(--muted-foreground)]">{item.period}</span>
                  </div>
                  <Card className="rounded-xl p-5">
                    <h3 className="text-xl font-semibold tracking-normal text-[var(--foreground)]">{item.title}</h3>
                    <p className="mt-2 leading-7 text-[var(--muted-foreground)]">{item.body}</p>
                  </Card>
                </div>
              ))}
            </div>
          </div>
        </section>
      </div>
    </main>
  )
}
