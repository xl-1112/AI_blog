import { useState, type ComponentProps, type ReactNode } from 'react'
import { Check, Copy } from 'lucide-react'

function extractText(node: ReactNode): string {
  if (typeof node === 'string' || typeof node === 'number') return String(node)
  if (Array.isArray(node)) return node.map(extractText).join('')
  if (node && typeof node === 'object' && 'props' in node) {
    const props = node.props as { children?: ReactNode }
    return extractText(props.children)
  }

  return ''
}

export function CodeBlock({ children, ...props }: ComponentProps<'pre'>) {
  const [copied, setCopied] = useState(false)
  const code = extractText(children)

  async function copyCode() {
    await navigator.clipboard.writeText(code)
    setCopied(true)
    window.setTimeout(() => setCopied(false), 1400)
  }

  return (
    <div className="code-block-wrapper group relative my-6">
      <pre {...props}>{children}</pre>
      <button
        type="button"
        onClick={copyCode}
        className="absolute right-3 top-3 inline-flex h-9 w-9 items-center justify-center rounded-lg border border-white/10 bg-white/10 text-white/80 opacity-0 backdrop-blur transition-all hover:bg-white/20 hover:text-white group-hover:opacity-100 focus-visible:opacity-100 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/50"
        aria-label={copied ? '代码已复制' : '复制代码'}
        title={copied ? '已复制' : '复制代码'}
      >
        {copied ? <Check className="h-4 w-4" /> : <Copy className="h-4 w-4" />}
      </button>
    </div>
  )
}
