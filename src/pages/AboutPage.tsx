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

function escapeHtml(value: string) {
  return value.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
}

function experienceBodyHtml(body: string) {
  if (/<[a-z][\s\S]*>/i.test(body)) {
    return { __html: body }
  }

  return { __html: escapeHtml(body).replace(/\n/g, '<br />') }
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
  const heroTitle = about.heroTitle || about.title || `${site.role} ${site.name}`
  const heroSubtitle = about.heroSubtitle || site.description
  const skillsBadge = about.skillsBadge || '技术栈'
  const skillsTitle = about.skillsTitle || '产品经理的工作栈'
  const skillsSubtitle = about.skillsSubtitle || ''
  const skillsDescription =
    about.skillsDescription || '以产品方法为核心，结合设计协作、数据分析和基础技术理解，保证需求判断和落地过程都更稳。'
  const timelineBadge = about.timelineBadge || '路线图'
  const timelineTitle = about.timelineTitle || '工作经历'
  const timelineSubtitle = about.timelineSubtitle || ''
  const timelineDescription = about.timelineDescription || ''
  const showHero = about.showHero ?? true
  const showHeroBadge = about.showHeroBadge ?? true
  const showHeroTitle = about.showHeroTitle ?? true
  const showHeroSubtitle = about.showHeroSubtitle ?? true
  const showHeroDescription = about.showHeroDescription ?? true
  const showContact = about.showContact ?? true
  const showLocation = about.showLocation ?? true
  const showEmail = about.showEmail ?? true
  const showGithub = about.showGithub ?? true
  const showContactBadge = about.showContactBadge ?? true
  const showContactTitle = about.showContactTitle ?? true
  const showContactSubtitle = about.showContactSubtitle ?? true
  const showContactDescription = about.showContactDescription ?? true
  const showSkills = about.showSkills ?? true
  const showSkillsBadge = about.showSkillsBadge ?? true
  const showSkillsTitle = about.showSkillsTitle ?? true
  const showSkillsSubtitle = about.showSkillsSubtitle ?? true
  const showSkillsDescription = about.showSkillsDescription ?? true
  const showTimeline = about.showTimeline ?? true
  const showTimelineBadge = about.showTimelineBadge ?? true
  const showTimelineTitle = about.showTimelineTitle ?? true
  const showTimelineSubtitle = about.showTimelineSubtitle ?? true
  const showTimelineDescription = about.showTimelineDescription ?? true
  const location = about.city || site.location
  const email = about.contact?.email || site.contact.email
  const github = about.contact?.github || site.contact.github
  const showIntroSection = showHero || showContact
  const introGridClass = showHero && showContact ? 'grid gap-10 lg:grid-cols-[0.9fr_1.1fr]' : 'space-y-10'

  return (
    <main className="px-5 py-14 sm:py-20">
      <div className="mx-auto max-w-6xl space-y-14">
        {showIntroSection ? (
          <section className={introGridClass}>
            {showHero ? (
              <div className="space-y-7">
                {showHeroBadge ? <SectionLabel pulse>{about.heroBadge || '关于我'}</SectionLabel> : null}
                <div className="space-y-5">
                  {showHeroTitle ? (
                    <h1 className="font-display text-5xl leading-tight text-[var(--foreground)] sm:text-6xl">
                      {heroTitle}
                    </h1>
                  ) : null}
                  {showHeroSubtitle ? <p className="text-lg leading-8 text-[var(--muted-foreground)]">{heroSubtitle}</p> : null}
                  {showHeroDescription ? <p className="leading-8 text-[var(--muted-foreground)]">{about.description}</p> : null}
                </div>
              </div>
            ) : null}

            {showContact ? (
              <Card className="p-7">
                {showContactBadge || showContactTitle || showContactSubtitle || showContactDescription ? (
                  <div className="mb-6 space-y-2">
                    {showContactBadge ? <p className="font-mono text-xs tracking-normal text-[var(--accent)]">{about.contactBadge || '联系我'}</p> : null}
                    {showContactTitle ? (
                      <h2 className="text-2xl font-semibold tracking-normal text-[var(--foreground)]">
                        {about.contactTitle || '欢迎交流产品、增长和项目协作'}
                      </h2>
                    ) : null}
                    {showContactSubtitle && about.contactSubtitle ? (
                      <p className="text-sm font-medium text-[var(--foreground)]">{about.contactSubtitle}</p>
                    ) : null}
                    {showContactDescription && about.contactDescription ? (
                      <p className="leading-7 text-[var(--muted-foreground)]">{about.contactDescription}</p>
                    ) : null}
                  </div>
                ) : null}
                <div className="grid gap-5 sm:grid-cols-2">
                  {showLocation ? (
                    <div className="rounded-2xl bg-[var(--muted)] p-5">
                      <MapPin className="mb-4 h-5 w-5 text-[var(--accent)]" />
                      <p className="text-sm text-[var(--muted-foreground)]">位置</p>
                      <p className="mt-1 font-semibold text-[var(--foreground)]">{location}</p>
                    </div>
                  ) : null}
                  {showEmail ? (
                    <a
                      href={`mailto:${email}`}
                      className="group rounded-2xl bg-[var(--muted)] p-5 transition-all hover:-translate-y-0.5 hover:shadow-md focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--ring)]"
                    >
                      <Mail className="mb-4 h-5 w-5 text-[var(--accent)]" />
                      <p className="text-sm text-[var(--muted-foreground)]">邮箱</p>
                      <p className="mt-1 break-all font-semibold text-[var(--foreground)] group-hover:text-[var(--accent)]">
                        {email}
                      </p>
                    </a>
                  ) : null}
                  {showGithub ? (
                    <a
                      href={github}
                      target="_blank"
                      rel="noreferrer"
                      className="group rounded-2xl bg-[var(--muted)] p-5 transition-all hover:-translate-y-0.5 hover:shadow-md focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--ring)] sm:col-span-2"
                    >
                      <GitHubIcon className="mb-4 h-5 w-5 text-[var(--accent)]" />
                      <p className="text-sm text-[var(--muted-foreground)]">GitHub</p>
                      <p className="mt-1 inline-flex items-center gap-2 break-all font-semibold text-[var(--foreground)] group-hover:text-[var(--accent)]">
                        {github}
                        <MoveUpRight className="h-4 w-4 flex-none" />
                      </p>
                    </a>
                  ) : null}
                </div>
              </Card>
            ) : null}
          </section>
        ) : null}

        {showSkills ? (
          <section className="space-y-6">
            {showSkillsBadge || showSkillsTitle || showSkillsSubtitle || showSkillsDescription ? (
              <div className="space-y-4">
                {showSkillsBadge ? <SectionLabel>{skillsBadge}</SectionLabel> : null}
                {showSkillsTitle ? <h2 className="font-display text-4xl text-[var(--foreground)] sm:text-5xl">{skillsTitle}</h2> : null}
                {showSkillsSubtitle && skillsSubtitle ? <p className="max-w-2xl text-lg leading-8 text-[var(--foreground)]">{skillsSubtitle}</p> : null}
                {showSkillsDescription ? <p className="max-w-2xl leading-8 text-[var(--muted-foreground)]">
                  {skillsDescription}
                </p> : null}
              </div>
            ) : null}
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
        ) : null}

        {showTimeline ? (
          <section className="space-y-6">
            {showTimelineBadge || showTimelineTitle || showTimelineSubtitle || showTimelineDescription ? (
              <div className="space-y-4">
                {showTimelineBadge ? <SectionLabel>{timelineBadge}</SectionLabel> : null}
                {showTimelineTitle ? <h2 className="font-display text-4xl text-[var(--foreground)] sm:text-5xl">{timelineTitle}</h2> : null}
                {showTimelineSubtitle && timelineSubtitle ? <p className="max-w-2xl text-lg leading-8 text-[var(--foreground)]">{timelineSubtitle}</p> : null}
                {showTimelineDescription && timelineDescription ? (
                  <p className="max-w-2xl leading-8 text-[var(--muted-foreground)]">{timelineDescription}</p>
                ) : null}
              </div>
            ) : null}
            <div className="relative rounded-2xl border border-[var(--border)] bg-white/75 p-6 shadow-[0_10px_30px_rgba(15,23,42,0.05)]">
              <div className="absolute bottom-8 left-8 top-8 hidden w-px bg-[var(--border)] sm:block" />
              <div className="space-y-5">
                {about.experience.map((item) => (
                  <div key={`${item.period}-${item.title}`} className="relative grid gap-4 sm:grid-cols-[7rem_1fr] sm:pl-10">
                    <div className="flex items-center gap-3 sm:block">
                      <span className="relative z-10 block h-3 w-3 flex-none rounded-full bg-[var(--accent)] ring-4 ring-blue-500/10 sm:absolute sm:left-[-0.4rem] sm:top-2.5" />
                      <span className="font-mono text-xs tracking-normal text-[var(--muted-foreground)]">{item.period}</span>
                    </div>
                    <Card className="rounded-xl p-5">
                      <h3 className="text-xl font-semibold tracking-normal text-[var(--foreground)]">{item.title}</h3>
                      <div className="timeline-rich-text mt-2" dangerouslySetInnerHTML={experienceBodyHtml(item.body)} />
                    </Card>
                  </div>
                ))}
              </div>
            </div>
          </section>
        ) : null}
      </div>
    </main>
  )
}
