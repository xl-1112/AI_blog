import Parser from "rss-parser";
import fs from "fs/promises";
import crypto from "crypto";

type Source = {
  name: string;
  type: "rss";
  url: string;
};

type NewsItem = {
  id: string;
  source: string;
  title: string;
  link: string;
  publishedAt?: string;
  summary?: string;
};

const parser = new Parser();

function hash(input: string) {
  return crypto.createHash("sha256").update(input).digest("hex");
}

async function main() {
  const raw = await fs.readFile("scripts/sources.json", "utf-8");
  const sourcesConfig = JSON.parse(raw);

  const sources: Source[] = [
    ...sourcesConfig.official,
    ...sourcesConfig.media
  ];

  const items: NewsItem[] = [];

  for (const source of sources) {
    try {
      const feed = await parser.parseURL(source.url);

      for (const item of feed.items.slice(0, 10)) {
        if (!item.title || !item.link) continue;

        items.push({
          id: hash(item.link),
          source: source.name,
          title: item.title,
          link: item.link,
          publishedAt: item.pubDate || item.isoDate,
          summary: item.contentSnippet || item.summary || ""
        });
      }
    } catch (error) {
      console.error(`Failed to fetch ${source.name}`, error);
    }
  }

  const uniqueMap = new Map<string, NewsItem>();

  for (const item of items) {
    uniqueMap.set(item.link, item);
  }

  const uniqueItems = Array.from(uniqueMap.values());

  await fs.mkdir("content", { recursive: true });
  await fs.writeFile(
    "content/latest-ai-news.json",
    JSON.stringify(uniqueItems, null, 2),
    "utf-8"
  );

  console.log(`Fetched ${uniqueItems.length} news items.`);
}

main();