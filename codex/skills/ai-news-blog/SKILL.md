# AI News Blog Skill

## Purpose

This skill helps generate a daily AI news article for a personal blog.

It should:
1. Collect recent AI news from trusted sources.
2. Filter low-value or repeated news.
3. Summarize the key points.
4. Rewrite into an original Chinese blog article.
5. Add source attribution.
6. Generate SEO title, summary, tags, and Markdown frontmatter.
7. Save the article as a Markdown blog post.

## Trusted Source Types

Prefer:
- Official company blogs
- Product release notes
- Research lab announcements
- GitHub release notes
- Reputable technology media
- Academic paper pages

Avoid:
- Low-quality repost sites
- Unverified rumors
- Pure marketing content
- Content without original source links

## Article Style

Write in Chinese.

Tone:
- Clear
- Practical
- Product-manager friendly
- Avoid exaggerated claims
- Avoid copying original paragraphs

Article structure:

1. Title
2. Frontmatter
3. 今日 AI 资讯摘要
4. 重点新闻解读
5. 对产品经理 / 开发者的影响
6. 值得关注的后续趋势
7. 信息来源

## Copyright Rules

Do not copy full articles.
Do not translate full original articles.
Only summarize, explain, and cite sources.
Use short quotes only when necessary.

## Output Format

Generate Markdown with this frontmatter:

---
title:
date:
categories:
  - AI资讯
tags:
  - AI
  - OpenAI
  - 产品经理
  - 人工智能
summary:
cover:
draft: true
---

## Quality Check

Before saving the article, check:

- Does the article contain at least 3 verified sources?
- Are repeated news items merged?
- Are source links included?
- Is the article original rather than copied?
- Is the title clear and not clickbait?
- Is draft set to true unless explicitly told to publish?