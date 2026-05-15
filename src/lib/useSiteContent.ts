import { useContext } from 'react'
import { ContentContext } from './content-context'

export function useSiteContent() {
  const context = useContext(ContentContext)
  if (!context) {
    throw new Error('useSiteContent must be used inside ContentProvider')
  }
  return context
}
