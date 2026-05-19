请使用 ai-news-blog Skill，根据 content/latest-ai-news.json 中的资讯，生成一篇中文 AI 日报博客文章。

要求：

1. 只选择最近 24 小时内较重要的 AI 资讯。
2. 优先选择官方发布、产品更新、模型发布、开源项目、行业重大变化。
3. 合并重复报道。
4. 不要全文翻译或复制原文。
5. 每条资讯都要保留来源名称和链接。
6. 文章面向产品经理、开发者和 AI 工具使用者。
7. 生成 Markdown 文件，保存到 source/_posts/。
8. 文件名格式：YYYY-MM-DD-ai-news.md。
9. frontmatter 中 draft 默认为 true。
10. 文章末尾加入“信息来源”列表。