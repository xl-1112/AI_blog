package cms

import (
	"errors"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"
)

type Content struct {
	Site      SiteSettings `json:"site"`
	Home      HomeContent  `json:"home"`
	Tags      []string     `json:"tags"`
	Articles  []Article    `json:"articles"`
	About     AboutContent `json:"about"`
	UpdatedAt time.Time    `json:"updatedAt"`
}

type SiteSettings struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	LogoURL     string  `json:"logoUrl"`
	Role        string  `json:"role"`
	Location    string  `json:"location"`
	Contact     Contact `json:"contact"`
}

type Contact struct {
	Email  string `json:"email"`
	GitHub string `json:"github"`
}

type HomeContent struct {
	IntroTitle         string   `json:"introTitle"`
	Intro              string   `json:"intro"`
	FeaturedArticleIDs []string `json:"featuredArticleIds"`
}

type AboutContent struct {
	Description string           `json:"description"`
	WorkStack   []WorkStackGroup `json:"workStack"`
	Experience  []ExperienceItem `json:"experience"`
}

type WorkStackGroup struct {
	Title string   `json:"title"`
	Items []string `json:"items"`
}

type ExperienceItem struct {
	Period string `json:"period"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type Article struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Date      string    `json:"date"`
	Summary   string    `json:"summary"`
	Tags      []string  `json:"tags"`
	Featured  bool      `json:"featured"`
	Draft     bool      `json:"draft"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ArticleSummary struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Date      string   `json:"date"`
	Summary   string   `json:"summary"`
	Tags      []string `json:"tags"`
	Featured  bool     `json:"featured"`
	Draft     bool     `json:"draft,omitempty"`
	WordCount int      `json:"wordCount"`
}

type PublicContent struct {
	Content
	FeaturedArticles []ArticleSummary `json:"featuredArticles"`
}

func DefaultContent(now time.Time) Content {
	return Content{
		Site: SiteSettings{
			Name:        "Liang",
			Description: "一个产品经理的个人网站，记录产品思考、项目复盘和少量技术笔记。",
			LogoURL:     "/uploads/blog_logo_header.png",
			Role:        "产品经理",
			Location:    "广州",
			Contact: Contact{
				Email:  "xl1258763@gmail.com",
				GitHub: "https://github.com/xl-1112",
			},
		},
		Home: HomeContent{
			IntroTitle: "Liang 的产品思考",
			Intro:      "这里沉淀产品经理视角下的用户洞察、需求判断、产品设计、项目复盘，也会偶尔记录技术相关的学习笔记。",
			FeaturedArticleIDs: []string{
				"product-feedback-loop",
			},
		},
		Tags: []string{"产品复盘", "用户洞察", "信息架构", "产品设计", "技术笔记"},
		Articles: []Article{
			{
				ID:       "product-feedback-loop",
				Title:    "从用户反馈到产品迭代的闭环记录",
				Date:     "2026-05-15",
				Summary:  "一篇产品经理方向的占位文章，展示如何记录用户问题、需求判断、方案取舍和上线复盘。",
				Tags:     []string{"产品复盘", "用户洞察"},
				Featured: true,
				Content:  "## 背景\n\n产品迭代不只是收集需求，更重要的是把反馈放回真实场景里理解。\n\n## 方法\n\n1. 记录用户原话和触发场景。\n2. 对齐业务目标和优先级。\n3. 输出方案、指标和验收口径。\n4. 上线后复盘数据和体验反馈。\n",
			},
			{
				ID:      "personal-site-ia-review",
				Title:   "个人网站的信息架构设计复盘",
				Date:    "2026-05-12",
				Summary: "从产品经理视角复盘个人网站的首版结构，记录目标用户、内容层级和后续迭代方向。",
				Tags:    []string{"信息架构", "产品设计"},
				Content: "## 目标\n\n这个网站的首版目标是让访问者快速理解我的身份、关注方向和文章内容。\n\n## 结构\n\n首页负责建立第一印象，文章列表负责承载内容沉淀，关于我负责补充工作经历和联系方式。\n",
			},
		},
		About: AboutContent{
			Description: "我关注从用户问题到产品方案的完整链路：先把场景和目标讲清楚，再通过数据、流程和协作把方案推进落地。",
			WorkStack: []WorkStackGroup{
				{Title: "产品方法", Items: []string{"用户研究", "需求分析", "竞品拆解", "PRD 编写"}},
				{Title: "设计协作", Items: []string{"信息架构", "交互原型", "体验走查", "设计评审"}},
				{Title: "数据与技术", Items: []string{"指标体系", "SQL 基础", "埋点分析", "接口理解"}},
			},
			Experience: []ExperienceItem{
				{Period: "现在", Title: "产品经理", Body: "围绕业务目标拆解产品问题，推进需求从调研、方案、评审到上线复盘的完整闭环。"},
				{Period: "持续积累", Title: "产品设计与项目协同", Body: "沉淀用户场景、流程设计、跨团队协作和版本节奏管理，让想法稳定落到真实体验里。"},
				{Period: "长期方向", Title: "产品 + 技术理解", Body: "保持对技术实现、系统边界和数据链路的理解，帮助产品判断更接近真实约束。"},
			},
		},
		UpdatedAt: now,
	}
}

func (c *Content) Normalize(now time.Time) {
	c.Site.Name = strings.TrimSpace(c.Site.Name)
	c.Site.Description = strings.TrimSpace(c.Site.Description)
	c.Site.LogoURL = strings.TrimSpace(c.Site.LogoURL)
	c.Site.Role = strings.TrimSpace(c.Site.Role)
	c.Site.Location = strings.TrimSpace(c.Site.Location)
	c.Site.Contact.Email = strings.TrimSpace(c.Site.Contact.Email)
	c.Site.Contact.GitHub = strings.TrimSpace(c.Site.Contact.GitHub)
	c.Home.IntroTitle = strings.TrimSpace(c.Home.IntroTitle)
	c.Home.Intro = strings.TrimSpace(c.Home.Intro)
	c.About.Description = strings.TrimSpace(c.About.Description)

	c.Tags = uniqueClean(c.Tags)

	featured := make(map[string]bool)
	for _, id := range c.Home.FeaturedArticleIDs {
		id = strings.TrimSpace(id)
		if id != "" {
			featured[id] = true
		}
	}

	for i := range c.Articles {
		a := &c.Articles[i]
		a.ID = strings.TrimSpace(a.ID)
		a.Title = strings.TrimSpace(a.Title)
		a.Date = strings.TrimSpace(a.Date)
		a.Summary = strings.TrimSpace(a.Summary)
		a.Content = strings.TrimSpace(a.Content)
		a.Tags = uniqueClean(a.Tags)
		if a.ID == "" {
			a.ID = slugify(a.Title)
		}
		if a.ID == "" {
			a.ID = "article-" + now.Format("20060102150405")
		}
		if a.Date == "" {
			a.Date = now.Format("2006-01-02")
		}
		if a.CreatedAt.IsZero() {
			a.CreatedAt = now
		}
		if a.UpdatedAt.IsZero() {
			a.UpdatedAt = now
		}
		if featured[a.ID] {
			a.Featured = true
		} else {
			a.Featured = false
		}
		c.Tags = uniqueClean(append(c.Tags, a.Tags...))
	}

	c.Home.FeaturedArticleIDs = filterExistingArticleIDs(c.Articles, featured)

	for i := range c.About.WorkStack {
		c.About.WorkStack[i].Title = strings.TrimSpace(c.About.WorkStack[i].Title)
		c.About.WorkStack[i].Items = uniqueClean(c.About.WorkStack[i].Items)
	}
	for i := range c.About.Experience {
		c.About.Experience[i].Period = strings.TrimSpace(c.About.Experience[i].Period)
		c.About.Experience[i].Title = strings.TrimSpace(c.About.Experience[i].Title)
		c.About.Experience[i].Body = strings.TrimSpace(c.About.Experience[i].Body)
	}
	c.UpdatedAt = now
}

func (c Content) Validate() error {
	if c.Site.Name == "" {
		return errors.New("site.name is required")
	}
	if c.Site.Contact.Email != "" && !strings.Contains(c.Site.Contact.Email, "@") {
		return errors.New("site.contact.email must be a valid email address")
	}
	if c.Site.Contact.GitHub != "" {
		u, err := url.ParseRequestURI(c.Site.Contact.GitHub)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
			return errors.New("site.contact.github must be an http(s) URL")
		}
	}

	seen := make(map[string]bool)
	for _, article := range c.Articles {
		if article.ID == "" {
			return errors.New("article.id is required")
		}
		if seen[article.ID] {
			return errors.New("article.id must be unique: " + article.ID)
		}
		seen[article.ID] = true
		if article.Title == "" {
			return errors.New("article.title is required: " + article.ID)
		}
		if _, err := time.Parse("2006-01-02", article.Date); err != nil {
			return errors.New("article.date must use YYYY-MM-DD: " + article.ID)
		}
	}

	for _, id := range c.Home.FeaturedArticleIDs {
		if !seen[id] {
			return errors.New("home.featuredArticleIds contains unknown article id: " + id)
		}
	}

	return nil
}

func (c Content) Public() PublicContent {
	public := c
	public.Articles = make([]Article, 0, len(c.Articles))
	featuredSet := make(map[string]bool)
	for _, id := range c.Home.FeaturedArticleIDs {
		featuredSet[id] = true
	}

	featured := make([]ArticleSummary, 0)
	for _, article := range c.Articles {
		if article.Draft {
			continue
		}
		public.Articles = append(public.Articles, article)
		if article.Featured || featuredSet[article.ID] {
			featured = append(featured, article.SummaryView(false))
		}
	}

	sort.Slice(public.Articles, func(i, j int) bool {
		return public.Articles[i].Date > public.Articles[j].Date
	})

	return PublicContent{
		Content:          public,
		FeaturedArticles: featured,
	}
}

func (a Article) SummaryView(includeDraft bool) ArticleSummary {
	summary := ArticleSummary{
		ID:        a.ID,
		Title:     a.Title,
		Date:      a.Date,
		Summary:   a.Summary,
		Tags:      append([]string(nil), a.Tags...),
		Featured:  a.Featured,
		WordCount: estimateWords(a.Title + " " + a.Summary + " " + strings.Join(a.Tags, " ") + " " + a.Content),
	}
	if includeDraft {
		summary.Draft = a.Draft
	}
	return summary
}

func estimateWords(source string) int {
	latin := regexp.MustCompile(`[A-Za-z0-9]+`).FindAllString(source, -1)
	count := len(latin)
	for _, r := range source {
		if unicode.Is(unicode.Han, r) {
			count++
		}
	}
	return count
}

func filterExistingArticleIDs(articles []Article, featured map[string]bool) []string {
	ids := make([]string, 0, len(featured))
	seen := make(map[string]bool)
	for _, article := range articles {
		if featured[article.ID] && !seen[article.ID] {
			ids = append(ids, article.ID)
			seen[article.ID] = true
		}
	}
	return ids
}

func uniqueClean(values []string) []string {
	out := make([]string, 0, len(values))
	seen := make(map[string]bool)
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		out = append(out, value)
	}
	return out
}

func slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var b strings.Builder
	lastDash := false
	for _, r := range value {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			lastDash = false
		case unicode.Is(unicode.Han, r):
			b.WriteRune(r)
			lastDash = false
		default:
			if !lastDash && b.Len() > 0 {
				b.WriteByte('-')
				lastDash = true
			}
		}
	}
	return strings.Trim(b.String(), "-")
}
