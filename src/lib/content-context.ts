import { createContext } from 'react'
import type { SiteContent } from './api'

export type ContentState = {
  content: SiteContent | null
  loading: boolean
  error: string | null
  reload: () => void
}

export const ContentContext = createContext<ContentState | null>(null)
