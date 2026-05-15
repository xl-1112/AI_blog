import type { ComponentProps } from 'react'
import { CodeBlock } from './CodeBlock'

export const mdxComponents = {
  pre: CodeBlock,
  a: (props: ComponentProps<'a'>) => (
    <a {...props} target={props.href?.startsWith('http') ? '_blank' : props.target} rel="noreferrer" />
  ),
}
