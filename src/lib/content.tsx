import { useEffect, useMemo, useState, type ReactNode } from 'react'
import { fetchSiteContent, type SiteContent } from './api'
import { ContentContext, type ContentState } from './content-context'

export function ContentProvider({ children }: { children: ReactNode }) {
  const [content, setContent] = useState<SiteContent | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [reloadKey, setReloadKey] = useState(0)

  useEffect(() => {
    const controller = new AbortController()

    async function loadContent() {
      setLoading(true)
      setError(null)
      try {
        const nextContent = await fetchSiteContent(controller.signal)
        setContent(nextContent)
      } catch (err) {
        if (controller.signal.aborted) return
        setError(err instanceof Error ? err.message : '内容加载失败')
      } finally {
        if (!controller.signal.aborted) {
          setLoading(false)
        }
      }
    }

    void loadContent()

    return () => controller.abort()
  }, [reloadKey])

  const value = useMemo<ContentState>(
    () => ({
      content,
      loading,
      error,
      reload: () => setReloadKey((key) => key + 1),
    }),
    [content, error, loading],
  )

  return <ContentContext.Provider value={value}>{children}</ContentContext.Provider>
}
