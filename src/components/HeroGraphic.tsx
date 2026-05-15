import { ChartNoAxesCombined, ClipboardList, PenLine, Target } from 'lucide-react'

export function HeroGraphic() {
  return (
    <div className="relative hidden min-h-[520px] lg:block" aria-hidden="true">
      <div className="absolute right-8 top-8 h-96 w-96 rounded-full border border-dashed border-blue-500/25 animate-slow-spin" />
      <div className="absolute right-24 top-24 h-72 w-72 rounded-full bg-[radial-gradient(circle,var(--accent)_0%,transparent_62%)] opacity-[0.08] blur-2xl" />
      <div className="absolute right-8 top-16 grid h-80 w-80 place-items-center rounded-[2rem] border border-white/80 bg-white/70 shadow-[0_30px_70px_rgba(15,23,42,0.12)] backdrop-blur-xl">
        <div className="grid h-40 w-40 place-items-center rounded-[2rem] bg-[linear-gradient(135deg,var(--accent),var(--accent-secondary))] text-white shadow-[0_18px_48px_rgba(0,82,255,0.35)]">
          <ChartNoAxesCombined className="h-16 w-16" />
        </div>
        <div className="absolute -left-10 top-12 animate-float rounded-2xl border border-[var(--border)] bg-white px-5 py-4 shadow-xl">
          <div className="mb-3 flex items-center gap-2 text-xs font-semibold text-[var(--muted-foreground)]">
            <Target className="h-4 w-4 text-[var(--accent)]" />
            产品目标
          </div>
          <div className="h-2 w-28 rounded-full bg-slate-100">
            <div className="h-full w-20 rounded-full bg-[linear-gradient(90deg,var(--accent),var(--accent-secondary))]" />
          </div>
        </div>
        <div className="absolute -right-12 bottom-16 animate-float-delayed rounded-2xl border border-[var(--border)] bg-white px-5 py-4 shadow-xl">
          <div className="mb-3 flex items-center gap-2 text-xs font-semibold text-[var(--muted-foreground)]">
            <ClipboardList className="h-4 w-4 text-[var(--accent)]" />
            需求路线图
          </div>
          <div className="grid grid-cols-3 gap-2">
            {Array.from({ length: 9 }).map((_, index) => (
              <span key={index} className="h-2 w-2 rounded-full bg-blue-500/30" />
            ))}
          </div>
        </div>
        <div className="absolute bottom-6 left-8 grid h-16 w-16 place-items-center rounded-2xl bg-[var(--foreground)] text-white shadow-xl">
          <PenLine className="h-7 w-7" />
        </div>
      </div>
    </div>
  )
}
