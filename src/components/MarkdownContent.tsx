import type { ReactNode } from 'react'
import { CodeBlock } from './CodeBlock'

type Block =
  | { type: 'heading'; level: number; text: string }
  | { type: 'paragraph'; text: string }
  | { type: 'blockquote'; text: string }
  | { type: 'list'; ordered: boolean; items: string[] }
  | { type: 'code'; code: string; language?: string }

type MarkdownContentProps = {
  content: string
  className?: string
}

export function MarkdownContent({ content, className = 'prose-blog' }: MarkdownContentProps) {
  return (
    <div className={className}>
      {parseMarkdown(content).map((block, index) => (
        <MarkdownBlock key={`${block.type}-${index}`} block={block} />
      ))}
    </div>
  )
}

function MarkdownBlock({ block }: { block: Block }) {
  switch (block.type) {
    case 'heading': {
      const children = renderInline(block.text)
      if (block.level === 1) return <h1>{children}</h1>
      if (block.level === 2) return <h2>{children}</h2>
      if (block.level === 3) return <h3>{children}</h3>
      return <h4>{children}</h4>
    }
    case 'blockquote':
      return <blockquote>{renderInline(block.text)}</blockquote>
    case 'list': {
      const items = block.items.map((item) => <li key={item}>{renderInline(item)}</li>)
      return block.ordered ? <ol>{items}</ol> : <ul>{items}</ul>
    }
    case 'code':
      return (
        <CodeBlock>
          <code className={block.language ? `language-${block.language}` : undefined}>{block.code}</code>
        </CodeBlock>
      )
    case 'paragraph':
    default:
      return <p>{renderInline(block.text)}</p>
  }
}

function parseMarkdown(markdown: string): Block[] {
  const lines = markdown.replace(/\r\n/g, '\n').split('\n')
  const blocks: Block[] = []
  const paragraph: string[] = []

  function flushParagraph() {
    const text = paragraph.join(' ').trim()
    if (text) {
      blocks.push({ type: 'paragraph', text })
    }
    paragraph.length = 0
  }

  for (let index = 0; index < lines.length; index += 1) {
    const line = lines[index]
    const trimmed = line.trim()

    if (!trimmed) {
      flushParagraph()
      continue
    }

    const codeFence = trimmed.match(/^```(\w+)?/)
    if (codeFence) {
      flushParagraph()
      const codeLines: string[] = []
      index += 1
      while (index < lines.length && !lines[index].trim().startsWith('```')) {
        codeLines.push(lines[index])
        index += 1
      }
      blocks.push({ type: 'code', code: codeLines.join('\n'), language: codeFence[1] })
      continue
    }

    const heading = trimmed.match(/^(#{1,6})\s+(.+)$/)
    if (heading) {
      flushParagraph()
      blocks.push({ type: 'heading', level: Math.min(heading[1].length, 4), text: heading[2].trim() })
      continue
    }

    const quote = trimmed.match(/^>\s?(.*)$/)
    if (quote) {
      flushParagraph()
      const quoteLines = [quote[1]]
      while (index + 1 < lines.length) {
        const next = lines[index + 1].trim().match(/^>\s?(.*)$/)
        if (!next) break
        quoteLines.push(next[1])
        index += 1
      }
      blocks.push({ type: 'blockquote', text: quoteLines.join(' ').trim() })
      continue
    }

    const list = trimmed.match(/^((?:[-*+])|\d+\.)\s+(.+)$/)
    if (list) {
      flushParagraph()
      const ordered = /\d+\./.test(list[1])
      const items = [list[2].trim()]
      while (index + 1 < lines.length) {
        const next = lines[index + 1].trim().match(/^((?:[-*+])|\d+\.)\s+(.+)$/)
        if (!next || /\d+\./.test(next[1]) !== ordered) break
        items.push(next[2].trim())
        index += 1
      }
      blocks.push({ type: 'list', ordered, items })
      continue
    }

    paragraph.push(trimmed)
  }

  flushParagraph()
  return blocks
}

function renderInline(text: string): ReactNode[] {
  const nodes: ReactNode[] = []
  const pattern = /(\*\*[^*]+\*\*|`[^`]+`|\[[^\]]+\]\([^)]+\))/g
  let lastIndex = 0
  let match: RegExpExecArray | null

  while ((match = pattern.exec(text))) {
    if (match.index > lastIndex) {
      nodes.push(text.slice(lastIndex, match.index))
    }

    const token = match[0]
    if (token.startsWith('**')) {
      nodes.push(<strong key={nodes.length}>{token.slice(2, -2)}</strong>)
    } else if (token.startsWith('`')) {
      nodes.push(<code key={nodes.length}>{token.slice(1, -1)}</code>)
    } else {
      const link = token.match(/^\[([^\]]+)\]\(([^)]+)\)$/)
      if (link) {
        const href = link[2]
        nodes.push(
          <a key={nodes.length} href={href} target={href.startsWith('http') ? '_blank' : undefined} rel="noreferrer">
            {link[1]}
          </a>,
        )
      }
    }

    lastIndex = pattern.lastIndex
  }

  if (lastIndex < text.length) {
    nodes.push(text.slice(lastIndex))
  }

  return nodes
}
