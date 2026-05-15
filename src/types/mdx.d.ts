declare module '*.mdx' {
  import type { ComponentType } from 'react'

  export const frontmatter: {
    id?: string
    title: string
    date: string
    summary: string
    tags: string[]
    featured?: boolean
  }

  const MDXComponent: ComponentType<{ components?: Record<string, ComponentType<unknown> | keyof JSX.IntrinsicElements> }>
  export default MDXComponent
}
