package cms

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func OpenStore(databasePath string, seedPath string, initialPassword string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(databasePath), 0755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", databasePath)
	if err != nil {
		return nil, err
	}
	store := &Store{db: db}
	if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := store.Migrate(); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := store.SeedFromJSON(seedPath); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := store.EnsureEditableDefaults(); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := store.EnsureSuperAdmin(initialPassword); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (s *Store) Close() error {
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) Migrate() error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			nickname TEXT NOT NULL,
			email TEXT NOT NULL DEFAULT '',
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL,
			status TEXT NOT NULL,
			last_login_at TEXT,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS login_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			username TEXT NOT NULL,
			ip TEXT NOT NULL,
			user_agent TEXT NOT NULL,
			success INTEGER NOT NULL,
			reason TEXT NOT NULL,
			created_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS articles (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			slug TEXT NOT NULL UNIQUE,
			category TEXT NOT NULL,
			date TEXT NOT NULL,
			summary TEXT NOT NULL,
			cover_url TEXT NOT NULL DEFAULT '',
			tags_json TEXT NOT NULL,
			featured INTEGER NOT NULL DEFAULT 0,
			status TEXT NOT NULL,
			content TEXT NOT NULL,
			seo_title TEXT NOT NULL DEFAULT '',
			seo_keywords TEXT NOT NULL DEFAULT '',
			seo_description TEXT NOT NULL DEFAULT '',
			view_count INTEGER NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS tags (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			use_count INTEGER NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS site_settings (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			name TEXT NOT NULL,
			site_title TEXT NOT NULL,
			description TEXT NOT NULL,
			keywords TEXT NOT NULL DEFAULT '',
			logo_url TEXT NOT NULL DEFAULT '',
			favicon_url TEXT NOT NULL DEFAULT '',
			icp TEXT NOT NULL DEFAULT '',
			analytics_code TEXT NOT NULL DEFAULT '',
			role TEXT NOT NULL DEFAULT '',
			location TEXT NOT NULL DEFAULT '',
			email TEXT NOT NULL DEFAULT '',
			github TEXT NOT NULL DEFAULT '',
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS home_settings (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			intro_title TEXT NOT NULL,
			subtitle TEXT NOT NULL DEFAULT '',
			intro TEXT NOT NULL,
			primary_cta_text TEXT NOT NULL DEFAULT '',
			secondary_cta_text TEXT NOT NULL DEFAULT '',
			featured_article_ids_json TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS about_profile (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			name TEXT NOT NULL,
			title TEXT NOT NULL,
			avatar_url TEXT NOT NULL DEFAULT '',
			show_hero INTEGER NOT NULL DEFAULT 1,
			show_hero_badge INTEGER NOT NULL DEFAULT 1,
			show_hero_title INTEGER NOT NULL DEFAULT 1,
			show_hero_subtitle INTEGER NOT NULL DEFAULT 1,
			show_hero_description INTEGER NOT NULL DEFAULT 1,
			hero_badge TEXT NOT NULL DEFAULT '关于我',
			hero_title TEXT NOT NULL DEFAULT '',
			hero_subtitle TEXT NOT NULL DEFAULT '',
			description TEXT NOT NULL,
			rich_description TEXT NOT NULL DEFAULT '',
			show_contact INTEGER NOT NULL DEFAULT 1,
			show_location INTEGER NOT NULL DEFAULT 1,
			show_email INTEGER NOT NULL DEFAULT 1,
			show_github INTEGER NOT NULL DEFAULT 1,
			show_contact_badge INTEGER NOT NULL DEFAULT 1,
			show_contact_title INTEGER NOT NULL DEFAULT 1,
			show_contact_subtitle INTEGER NOT NULL DEFAULT 1,
			show_contact_description INTEGER NOT NULL DEFAULT 1,
			contact_badge TEXT NOT NULL DEFAULT '联系我',
			contact_title TEXT NOT NULL DEFAULT '欢迎交流产品、增长和项目协作',
			contact_subtitle TEXT NOT NULL DEFAULT '',
			contact_description TEXT NOT NULL DEFAULT '',
			show_skills INTEGER NOT NULL DEFAULT 1,
			show_skills_header INTEGER NOT NULL DEFAULT 1,
			show_skills_badge INTEGER NOT NULL DEFAULT 1,
			show_skills_title INTEGER NOT NULL DEFAULT 1,
			show_skills_subtitle INTEGER NOT NULL DEFAULT 1,
			show_skills_description INTEGER NOT NULL DEFAULT 1,
			skills_badge TEXT NOT NULL DEFAULT '技术栈',
			skills_title TEXT NOT NULL DEFAULT '产品经理的工作栈',
			skills_subtitle TEXT NOT NULL DEFAULT '',
			skills_description TEXT NOT NULL DEFAULT '以产品方法为核心，结合设计协作、数据分析和基础技术理解，保证需求判断和落地过程都更稳。',
			show_timeline INTEGER NOT NULL DEFAULT 1,
			show_timeline_header INTEGER NOT NULL DEFAULT 1,
			show_timeline_badge INTEGER NOT NULL DEFAULT 1,
			show_timeline_title INTEGER NOT NULL DEFAULT 1,
			show_timeline_subtitle INTEGER NOT NULL DEFAULT 1,
			show_timeline_description INTEGER NOT NULL DEFAULT 1,
			timeline_badge TEXT NOT NULL DEFAULT '路线图',
			timeline_title TEXT NOT NULL DEFAULT '工作经历',
			timeline_subtitle TEXT NOT NULL DEFAULT '',
			timeline_description TEXT NOT NULL DEFAULT '',
			city TEXT NOT NULL DEFAULT '',
			email TEXT NOT NULL DEFAULT '',
			github TEXT NOT NULL DEFAULT '',
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS skills (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			items_json TEXT NOT NULL,
			sort INTEGER NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS timeline_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			period TEXT NOT NULL,
			title TEXT NOT NULL,
			body TEXT NOT NULL,
			sort INTEGER NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS article_view_events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			article_id TEXT NOT NULL,
			view_date TEXT NOT NULL,
			views INTEGER NOT NULL,
			created_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS uploads (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			url TEXT NOT NULL,
			kind TEXT NOT NULL,
			filename TEXT NOT NULL,
			created_at TEXT NOT NULL
		)`,
	}
	for _, statement := range statements {
		if _, err := s.db.Exec(statement); err != nil {
			return err
		}
	}
	return s.ensureAboutProfileColumns()
}

func (s *Store) ensureAboutProfileColumns() error {
	columns, err := s.tableColumns("about_profile")
	if err != nil {
		return err
	}
	missing := map[string]string{
		"show_hero":                 `ALTER TABLE about_profile ADD COLUMN show_hero INTEGER NOT NULL DEFAULT 1`,
		"show_hero_badge":           `ALTER TABLE about_profile ADD COLUMN show_hero_badge INTEGER NOT NULL DEFAULT 1`,
		"show_hero_title":           `ALTER TABLE about_profile ADD COLUMN show_hero_title INTEGER NOT NULL DEFAULT 1`,
		"show_hero_subtitle":        `ALTER TABLE about_profile ADD COLUMN show_hero_subtitle INTEGER NOT NULL DEFAULT 1`,
		"show_hero_description":     `ALTER TABLE about_profile ADD COLUMN show_hero_description INTEGER NOT NULL DEFAULT 1`,
		"hero_badge":                `ALTER TABLE about_profile ADD COLUMN hero_badge TEXT NOT NULL DEFAULT '关于我'`,
		"hero_title":                `ALTER TABLE about_profile ADD COLUMN hero_title TEXT NOT NULL DEFAULT ''`,
		"hero_subtitle":             `ALTER TABLE about_profile ADD COLUMN hero_subtitle TEXT NOT NULL DEFAULT ''`,
		"show_contact":              `ALTER TABLE about_profile ADD COLUMN show_contact INTEGER NOT NULL DEFAULT 1`,
		"show_location":             `ALTER TABLE about_profile ADD COLUMN show_location INTEGER NOT NULL DEFAULT 1`,
		"show_email":                `ALTER TABLE about_profile ADD COLUMN show_email INTEGER NOT NULL DEFAULT 1`,
		"show_github":               `ALTER TABLE about_profile ADD COLUMN show_github INTEGER NOT NULL DEFAULT 1`,
		"show_contact_badge":        `ALTER TABLE about_profile ADD COLUMN show_contact_badge INTEGER NOT NULL DEFAULT 1`,
		"show_contact_title":        `ALTER TABLE about_profile ADD COLUMN show_contact_title INTEGER NOT NULL DEFAULT 1`,
		"show_contact_subtitle":     `ALTER TABLE about_profile ADD COLUMN show_contact_subtitle INTEGER NOT NULL DEFAULT 1`,
		"show_contact_description":  `ALTER TABLE about_profile ADD COLUMN show_contact_description INTEGER NOT NULL DEFAULT 1`,
		"contact_badge":             `ALTER TABLE about_profile ADD COLUMN contact_badge TEXT NOT NULL DEFAULT '联系我'`,
		"contact_title":             `ALTER TABLE about_profile ADD COLUMN contact_title TEXT NOT NULL DEFAULT '欢迎交流产品、增长和项目协作'`,
		"contact_subtitle":          `ALTER TABLE about_profile ADD COLUMN contact_subtitle TEXT NOT NULL DEFAULT ''`,
		"contact_description":       `ALTER TABLE about_profile ADD COLUMN contact_description TEXT NOT NULL DEFAULT ''`,
		"show_skills":               `ALTER TABLE about_profile ADD COLUMN show_skills INTEGER NOT NULL DEFAULT 1`,
		"show_skills_header":        `ALTER TABLE about_profile ADD COLUMN show_skills_header INTEGER NOT NULL DEFAULT 1`,
		"show_skills_badge":         `ALTER TABLE about_profile ADD COLUMN show_skills_badge INTEGER NOT NULL DEFAULT 1`,
		"show_skills_title":         `ALTER TABLE about_profile ADD COLUMN show_skills_title INTEGER NOT NULL DEFAULT 1`,
		"show_skills_subtitle":      `ALTER TABLE about_profile ADD COLUMN show_skills_subtitle INTEGER NOT NULL DEFAULT 1`,
		"show_skills_description":   `ALTER TABLE about_profile ADD COLUMN show_skills_description INTEGER NOT NULL DEFAULT 1`,
		"skills_badge":              `ALTER TABLE about_profile ADD COLUMN skills_badge TEXT NOT NULL DEFAULT '技术栈'`,
		"skills_title":              `ALTER TABLE about_profile ADD COLUMN skills_title TEXT NOT NULL DEFAULT '产品经理的工作栈'`,
		"skills_subtitle":           `ALTER TABLE about_profile ADD COLUMN skills_subtitle TEXT NOT NULL DEFAULT ''`,
		"skills_description":        `ALTER TABLE about_profile ADD COLUMN skills_description TEXT NOT NULL DEFAULT '以产品方法为核心，结合设计协作、数据分析和基础技术理解，保证需求判断和落地过程都更稳。'`,
		"show_timeline":             `ALTER TABLE about_profile ADD COLUMN show_timeline INTEGER NOT NULL DEFAULT 1`,
		"show_timeline_header":      `ALTER TABLE about_profile ADD COLUMN show_timeline_header INTEGER NOT NULL DEFAULT 1`,
		"show_timeline_badge":       `ALTER TABLE about_profile ADD COLUMN show_timeline_badge INTEGER NOT NULL DEFAULT 1`,
		"show_timeline_title":       `ALTER TABLE about_profile ADD COLUMN show_timeline_title INTEGER NOT NULL DEFAULT 1`,
		"show_timeline_subtitle":    `ALTER TABLE about_profile ADD COLUMN show_timeline_subtitle INTEGER NOT NULL DEFAULT 1`,
		"show_timeline_description": `ALTER TABLE about_profile ADD COLUMN show_timeline_description INTEGER NOT NULL DEFAULT 1`,
		"timeline_badge":            `ALTER TABLE about_profile ADD COLUMN timeline_badge TEXT NOT NULL DEFAULT '路线图'`,
		"timeline_title":            `ALTER TABLE about_profile ADD COLUMN timeline_title TEXT NOT NULL DEFAULT '工作经历'`,
		"timeline_subtitle":         `ALTER TABLE about_profile ADD COLUMN timeline_subtitle TEXT NOT NULL DEFAULT ''`,
		"timeline_description":      `ALTER TABLE about_profile ADD COLUMN timeline_description TEXT NOT NULL DEFAULT ''`,
	}
	for name, statement := range missing {
		if columns[name] {
			continue
		}
		if _, err := s.db.Exec(statement); err != nil {
			return err
		}
	}
	if _, err := s.db.Exec(`UPDATE about_profile SET hero_title = title WHERE TRIM(hero_title) = ''`); err != nil {
		return err
	}
	if _, err := s.db.Exec(`UPDATE about_profile SET hero_subtitle = COALESCE((SELECT description FROM site_settings WHERE id = 1), hero_subtitle) WHERE TRIM(hero_subtitle) = ''`); err != nil {
		return err
	}
	return nil
}

func (s *Store) tableColumns(table string) (map[string]bool, error) {
	rows, err := s.db.Query(`PRAGMA table_info(` + table + `)`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns := make(map[string]bool)
	for rows.Next() {
		var cid int
		var name, columnType string
		var notNull int
		var defaultValue any
		var pk int
		if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultValue, &pk); err != nil {
			return nil, err
		}
		columns[name] = true
	}
	return columns, rows.Err()
}

func (s *Store) SeedFromJSON(seedPath string) error {
	var count int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM site_settings`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	content := fallbackContent(time.Now().UTC())
	if data, err := os.ReadFile(seedPath); err == nil {
		_ = json.Unmarshal(data, &content)
	}
	return s.ReplaceContent(content)
}

func (s *Store) EnsureEditableDefaults() error {
	now := time.Now().UTC()
	var count int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM tags`).Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		if _, err := s.SaveTags([]string{"产品复盘", "用户洞察", "信息架构", "产品设计", "技术笔记"}); err != nil {
			return err
		}
	}
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM skills`).Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		defaults := []WorkStackGroup{
			{Title: "产品方法", Items: []string{"用户研究", "需求分析", "竞品拆解", "PRD 编写"}},
			{Title: "设计协作", Items: []string{"信息架构", "交互原型", "体验走查", "设计评审"}},
			{Title: "数据与技术", Items: []string{"指标体系", "SQL 基础", "埋点分析", "接口理解"}},
		}
		tx, err := s.db.Begin()
		if err != nil {
			return err
		}
		for index, item := range defaults {
			if err := insertSkill(tx, item, index, now); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM timeline_items`).Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		defaults := []ExperienceItem{
			{Period: "现在", Title: "广东东开创新智慧科技有限公司", Body: "围绕业务目标拆解产品问题，推进需求从调研、方案、评审到上线复盘的完整闭环。"},
			{Period: "持续积累", Title: "产品设计与项目协同", Body: "沉淀用户场景、流程设计、跨团队协作和版本节奏管理，让想法稳定落到真实体验里。"},
			{Period: "长期方向", Title: "产品 + 技术理解", Body: "保持对技术实现、系统边界和数据链路的理解，帮助产品判断更接近真实约束。"},
		}
		tx, err := s.db.Begin()
		if err != nil {
			return err
		}
		for index, item := range defaults {
			if err := insertTimeline(tx, item, index, now); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
		return tx.Commit()
	}
	return nil
}

func fallbackContent(now time.Time) Content {
	return Content{
		Site: SiteSettings{
			Name:        "Liang",
			SiteTitle:   "Liang | Blog",
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
			IntroTitle:         "Liang",
			Subtitle:           "产品经理",
			Intro:              "这里沉淀产品经理视角下的用户洞察、需求判断、产品设计、项目复盘，也会偶尔记录技术相关的学习笔记。",
			PrimaryCtaText:     "阅读文章",
			SecondaryCtaText:   "关于我",
			FeaturedArticleIDs: []string{"product-feedback-loop"},
		},
		Tags: []string{"产品复盘", "用户洞察", "信息架构", "产品设计", "技术笔记"},
		Articles: []Article{{
			ID:        "product-feedback-loop",
			Title:     "产品心得",
			Slug:      "product-feedback-loop",
			Category:  "产品复盘",
			Date:      now.Format("2006-01-02"),
			Summary:   "一篇产品经理方向的占位文章。",
			Tags:      []string{"产品复盘"},
			Featured:  true,
			Status:    ArticlePublished,
			Content:   "## 背景\n\n产品迭代不只是收集需求，更重要的是把反馈放回真实场景里理解。",
			CreatedAt: now,
			UpdatedAt: now,
		}},
		About: AboutContent{
			Name:                    "Liang",
			Title:                   "产品经理 Liang",
			ShowHero:                true,
			ShowHeroBadge:           true,
			ShowHeroTitle:           true,
			ShowHeroSubtitle:        true,
			ShowHeroDescription:     true,
			HeroBadge:               "关于我",
			HeroTitle:               "产品经理 Liang",
			HeroSubtitle:            "一个产品经理的个人网站，记录产品思考、项目复盘和少量技术笔记。",
			Description:             "我关注从用户问题到产品方案的完整链路。",
			ShowContact:             true,
			ShowLocation:            true,
			ShowEmail:               true,
			ShowGithub:              true,
			ShowContactBadge:        true,
			ShowContactTitle:        true,
			ShowContactSubtitle:     true,
			ShowContactDescription:  true,
			ContactBadge:            "联系我",
			ContactTitle:            "欢迎交流产品、增长和项目协作",
			ContactSubtitle:         "",
			ContactDescription:      "",
			ShowSkills:              true,
			ShowSkillsHeader:        true,
			ShowSkillsBadge:         true,
			ShowSkillsTitle:         true,
			ShowSkillsSubtitle:      true,
			ShowSkillsDescription:   true,
			SkillsBadge:             "技术栈",
			SkillsTitle:             "产品经理的工作栈",
			SkillsSubtitle:          "",
			SkillsDescription:       "以产品方法为核心，结合设计协作、数据分析和基础技术理解，保证需求判断和落地过程都更稳。",
			ShowTimeline:            true,
			ShowTimelineHeader:      true,
			ShowTimelineBadge:       true,
			ShowTimelineTitle:       true,
			ShowTimelineSubtitle:    true,
			ShowTimelineDescription: true,
			TimelineBadge:           "路线图",
			TimelineTitle:           "工作经历",
			TimelineSubtitle:        "",
			TimelineDescription:     "",
			City:                    "广州",
			Contact: Contact{
				Email:  "xl1258763@gmail.com",
				GitHub: "https://github.com/xl-1112",
			},
			WorkStack:  []WorkStackGroup{{Title: "产品方法", Items: []string{"用户研究", "需求分析", "PRD 编写"}}},
			Experience: []ExperienceItem{{Period: "现在", Title: "产品经理", Body: "推进需求从调研、方案、评审到上线复盘。"}},
		},
		UpdatedAt: now,
	}
}

func (s *Store) ReplaceContent(content Content) error {
	normalizeContent(&content, time.Now().UTC())
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, table := range []string{"site_settings", "home_settings", "about_profile", "articles", "tags", "skills", "timeline_items"} {
		if _, err := tx.Exec(`DELETE FROM ` + table); err != nil {
			return err
		}
	}

	now := time.Now().UTC()
	if err := upsertSite(tx, content.Site, now); err != nil {
		return err
	}
	if err := upsertHome(tx, content.Home, now); err != nil {
		return err
	}
	if err := upsertAbout(tx, content.About, now); err != nil {
		return err
	}
	for _, article := range content.Articles {
		if err := upsertArticle(tx, article, now); err != nil {
			return err
		}
	}
	if err := replaceTags(tx, content.Tags, content.Articles, now); err != nil {
		return err
	}
	for index, group := range content.About.WorkStack {
		if err := insertSkill(tx, group, index, now); err != nil {
			return err
		}
	}
	for index, item := range content.About.Experience {
		if err := insertTimeline(tx, item, index, now); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func normalizeContent(content *Content, now time.Time) {
	if content.Site.SiteTitle == "" {
		content.Site.SiteTitle = content.Site.Name + " | Blog"
	}
	if content.Site.Name == "" {
		content.Site.Name = "Liang"
	}
	if content.Home.Subtitle == "" {
		content.Home.Subtitle = content.Site.Role
	}
	if content.Home.PrimaryCtaText == "" {
		content.Home.PrimaryCtaText = "阅读文章"
	}
	if content.Home.SecondaryCtaText == "" {
		content.Home.SecondaryCtaText = "关于我"
	}
	if content.About.Name == "" {
		content.About.Name = content.Site.Name
	}
	if content.About.Title == "" {
		content.About.Title = content.Site.Role + " " + content.Site.Name
	}
	applyAboutVisibilityDefaults(&content.About)
	if content.About.HeroBadge == "" {
		content.About.HeroBadge = "关于我"
	}
	if content.About.HeroTitle == "" {
		content.About.HeroTitle = content.About.Title
	}
	if content.About.HeroSubtitle == "" {
		content.About.HeroSubtitle = content.Site.Description
	}
	if content.About.SkillsBadge == "" {
		content.About.SkillsBadge = "技术栈"
	}
	if content.About.SkillsTitle == "" {
		content.About.SkillsTitle = "产品经理的工作栈"
	}
	if content.About.SkillsDescription == "" {
		content.About.SkillsDescription = "以产品方法为核心，结合设计协作、数据分析和基础技术理解，保证需求判断和落地过程都更稳。"
	}
	if content.About.TimelineBadge == "" {
		content.About.TimelineBadge = "路线图"
	}
	if content.About.TimelineTitle == "" {
		content.About.TimelineTitle = "工作经历"
	}
	if content.About.City == "" {
		content.About.City = content.Site.Location
	}
	if content.About.Contact.Email == "" {
		content.About.Contact.Email = content.Site.Contact.Email
	}
	if content.About.Contact.GitHub == "" {
		content.About.Contact.GitHub = content.Site.Contact.GitHub
	}
	for index := range content.Articles {
		article := &content.Articles[index]
		if article.ID == "" {
			article.ID = slugify(article.Title)
		}
		if article.Slug == "" {
			article.Slug = article.ID
		}
		if article.Category == "" && len(article.Tags) > 0 {
			article.Category = article.Tags[0]
		}
		if article.Category != "" {
			article.Tags = []string{article.Category}
		}
		if article.Date == "" {
			article.Date = now.Format("2006-01-02")
		}
		article.Status = normalizeArticleStatus(article.Status, article.Draft)
		article.Draft = article.Status == ArticleDraft
		if article.CreatedAt.IsZero() {
			article.CreatedAt = now
		}
		if article.UpdatedAt.IsZero() {
			article.UpdatedAt = now
		}
		content.Tags = uniqueClean(append(content.Tags, article.Tags...))
	}
	content.Tags = uniqueClean(content.Tags)
}

func (s *Store) EnsureSuperAdmin(initialPassword string) error {
	var count int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(initialPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	_, err = s.db.Exec(
		`INSERT INTO users (username,nickname,email,password_hash,role,status,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?)`,
		"admin", "超级管理员", "", string(hash), RoleSuperAdmin, UserActive, now, now,
	)
	return err
}

func (s *Store) PublicContent() (Content, error) {
	content, err := s.AdminContent()
	if err != nil {
		return Content{}, err
	}
	published := make([]Article, 0)
	for _, article := range content.Articles {
		if article.Status == ArticlePublished {
			published = append(published, article)
		}
	}
	content.Articles = published
	featuredSet := make(map[string]bool)
	for _, id := range content.Home.FeaturedArticleIDs {
		featuredSet[id] = true
	}
	content.FeaturedArticles = nil
	for _, article := range published {
		if article.Featured || featuredSet[article.ID] {
			content.FeaturedArticles = append(content.FeaturedArticles, article.SummaryView(false))
		}
	}
	return content, nil
}

func (s *Store) AdminContent() (Content, error) {
	site, err := s.Site()
	if err != nil {
		return Content{}, err
	}
	home, err := s.Home()
	if err != nil {
		return Content{}, err
	}
	about, err := s.About()
	if err != nil {
		return Content{}, err
	}
	tags, err := s.Tags()
	if err != nil {
		return Content{}, err
	}
	articles, _, err := s.ListArticles(ArticleQuery{Page: 1, PageSize: 10000, IncludeDraft: true})
	if err != nil {
		return Content{}, err
	}
	updatedAt := time.Now().UTC()
	return Content{Site: site, Home: home, Tags: tags, Articles: articles, About: about, UpdatedAt: updatedAt}, nil
}

type ArticleQuery struct {
	Keyword      string
	Tag          string
	Status       string
	DateFrom     string
	DateTo       string
	Page         int
	PageSize     int
	IncludeDraft bool
}

func (s *Store) ListArticles(query ArticleQuery) ([]Article, int, error) {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}
	rows, err := s.db.Query(`SELECT id,title,slug,category,date,summary,cover_url,tags_json,featured,status,content,seo_title,seo_keywords,seo_description,view_count,created_at,updated_at FROM articles ORDER BY date DESC, updated_at DESC`)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	all := make([]Article, 0)
	for rows.Next() {
		article, err := scanArticle(rows)
		if err != nil {
			return nil, 0, err
		}
		if !query.IncludeDraft && article.Status != ArticlePublished {
			continue
		}
		if query.Keyword != "" {
			needle := strings.ToLower(query.Keyword)
			haystack := strings.ToLower(article.Title + " " + article.Summary + " " + article.Category + " " + strings.Join(article.Tags, " "))
			if !strings.Contains(haystack, needle) {
				continue
			}
		}
		if query.Tag != "" && article.Category != query.Tag {
			continue
		}
		if query.Status != "" && article.Status != query.Status {
			continue
		}
		if query.DateFrom != "" && article.Date < query.DateFrom {
			continue
		}
		if query.DateTo != "" && article.Date > query.DateTo {
			continue
		}
		all = append(all, article)
	}
	total := len(all)
	start := (query.Page - 1) * query.PageSize
	if start >= total {
		return []Article{}, total, nil
	}
	end := start + query.PageSize
	if end > total {
		end = total
	}
	return all[start:end], total, nil
}

func (s *Store) Article(idOrSlug string, includeDraft bool) (Article, error) {
	row := s.db.QueryRow(`SELECT id,title,slug,category,date,summary,cover_url,tags_json,featured,status,content,seo_title,seo_keywords,seo_description,view_count,created_at,updated_at FROM articles WHERE id = ? OR slug = ?`, idOrSlug, idOrSlug)
	article, err := scanArticle(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Article{}, err
		}
		return Article{}, err
	}
	if !includeDraft && article.Status != ArticlePublished {
		return Article{}, sql.ErrNoRows
	}
	return article, nil
}

func (s *Store) SaveArticle(article Article, originalID string) (Article, error) {
	now := time.Now().UTC()
	if article.ID == "" {
		article.ID = slugify(article.Title)
	}
	if article.Slug == "" {
		article.Slug = article.ID
	}
	if article.Category == "" && len(article.Tags) > 0 {
		article.Category = article.Tags[0]
	}
	if article.Category != "" {
		article.Tags = []string{article.Category}
	}
	if article.Date == "" {
		article.Date = now.Format("2006-01-02")
	}
	article.Status = normalizeArticleStatus(article.Status, article.Draft)
	article.Draft = article.Status == ArticleDraft
	article.WordCount = estimateWords(article.Title + " " + article.Summary + " " + article.Content)
	tagsJSON, _ := json.Marshal(article.Tags)
	if originalID == "" {
		originalID = article.ID
	}
	var existingCreated string
	var existingViews int
	err := s.db.QueryRow(`SELECT created_at, view_count FROM articles WHERE id = ?`, originalID).Scan(&existingCreated, &existingViews)
	if errors.Is(err, sql.ErrNoRows) {
		article.CreatedAt = now
		article.UpdatedAt = now
		_, err = s.db.Exec(`INSERT INTO articles (id,title,slug,category,date,summary,cover_url,tags_json,featured,status,content,seo_title,seo_keywords,seo_description,view_count,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
			article.ID, article.Title, article.Slug, article.Category, article.Date, article.Summary, article.CoverURL, string(tagsJSON), boolInt(article.Featured), article.Status, article.Content, article.SEOTitle, article.SEOKeywords, article.SEODescription, article.ViewCount, formatTime(article.CreatedAt), formatTime(article.UpdatedAt))
	} else if err == nil {
		article.CreatedAt = parseTime(existingCreated)
		article.UpdatedAt = now
		article.ViewCount = existingViews
		_, err = s.db.Exec(`UPDATE articles SET id=?, title=?, slug=?, category=?, date=?, summary=?, cover_url=?, tags_json=?, featured=?, status=?, content=?, seo_title=?, seo_keywords=?, seo_description=?, view_count=?, updated_at=? WHERE id=?`,
			article.ID, article.Title, article.Slug, article.Category, article.Date, article.Summary, article.CoverURL, string(tagsJSON), boolInt(article.Featured), article.Status, article.Content, article.SEOTitle, article.SEOKeywords, article.SEODescription, article.ViewCount, formatTime(article.UpdatedAt), originalID)
	}
	if err != nil {
		return Article{}, err
	}
	if err := s.rebuildTagUseCounts(); err != nil {
		return Article{}, err
	}
	return s.Article(article.ID, true)
}

func (s *Store) DeleteArticle(id string) error {
	if _, err := s.db.Exec(`DELETE FROM articles WHERE id = ?`, id); err != nil {
		return err
	}
	return s.rebuildTagUseCounts()
}

func (s *Store) SetArticleStatus(id string, status string) error {
	status = normalizeArticleStatus(status, status == ArticleDraft)
	_, err := s.db.Exec(`UPDATE articles SET status=?, updated_at=? WHERE id=?`, status, formatTime(time.Now().UTC()), id)
	return err
}

func (s *Store) IncrementArticleView(id string) error {
	now := time.Now().UTC()
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`UPDATE articles SET view_count = view_count + 1 WHERE id = ?`, id); err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO article_view_events (article_id, view_date, views, created_at) VALUES (?,?,?,?)`, id, now.Format("2006-01-02"), 1, formatTime(now))
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Store) Site() (SiteSettings, error) {
	var site SiteSettings
	var updated string
	err := s.db.QueryRow(`SELECT name,site_title,description,keywords,logo_url,favicon_url,icp,analytics_code,role,location,email,github,updated_at FROM site_settings WHERE id = 1`).
		Scan(&site.Name, &site.SiteTitle, &site.Description, &site.Keywords, &site.LogoURL, &site.FaviconURL, &site.ICP, &site.AnalyticsCode, &site.Role, &site.Location, &site.Contact.Email, &site.Contact.GitHub, &updated)
	return site, err
}

func (s *Store) SaveSite(site SiteSettings) (SiteSettings, error) {
	if site.Name == "" {
		site.Name = "Liang"
	}
	if site.SiteTitle == "" {
		site.SiteTitle = site.Name + " | Blog"
	}
	if _, err := s.db.Exec(`INSERT INTO site_settings (id,name,site_title,description,keywords,logo_url,favicon_url,icp,analytics_code,role,location,email,github,updated_at) VALUES (1,?,?,?,?,?,?,?,?,?,?,?,?,?)
		ON CONFLICT(id) DO UPDATE SET name=excluded.name,site_title=excluded.site_title,description=excluded.description,keywords=excluded.keywords,logo_url=excluded.logo_url,favicon_url=excluded.favicon_url,icp=excluded.icp,analytics_code=excluded.analytics_code,role=excluded.role,location=excluded.location,email=excluded.email,github=excluded.github,updated_at=excluded.updated_at`,
		site.Name, site.SiteTitle, site.Description, site.Keywords, site.LogoURL, site.FaviconURL, site.ICP, site.AnalyticsCode, site.Role, site.Location, site.Contact.Email, site.Contact.GitHub, formatTime(time.Now().UTC())); err != nil {
		return SiteSettings{}, err
	}
	return s.Site()
}

func (s *Store) Home() (HomeContent, error) {
	var home HomeContent
	var idsJSON string
	err := s.db.QueryRow(`SELECT intro_title,subtitle,intro,primary_cta_text,secondary_cta_text,featured_article_ids_json FROM home_settings WHERE id = 1`).
		Scan(&home.IntroTitle, &home.Subtitle, &home.Intro, &home.PrimaryCtaText, &home.SecondaryCtaText, &idsJSON)
	if err != nil {
		return HomeContent{}, err
	}
	_ = json.Unmarshal([]byte(idsJSON), &home.FeaturedArticleIDs)
	return home, nil
}

func (s *Store) SaveHome(home HomeContent) (HomeContent, error) {
	idsJSON, _ := json.Marshal(uniqueClean(home.FeaturedArticleIDs))
	tx, err := s.db.Begin()
	if err != nil {
		return HomeContent{}, err
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`INSERT INTO home_settings (id,intro_title,subtitle,intro,primary_cta_text,secondary_cta_text,featured_article_ids_json,updated_at) VALUES (1,?,?,?,?,?,?,?)
		ON CONFLICT(id) DO UPDATE SET intro_title=excluded.intro_title,subtitle=excluded.subtitle,intro=excluded.intro,primary_cta_text=excluded.primary_cta_text,secondary_cta_text=excluded.secondary_cta_text,featured_article_ids_json=excluded.featured_article_ids_json,updated_at=excluded.updated_at`,
		home.IntroTitle, home.Subtitle, home.Intro, home.PrimaryCtaText, home.SecondaryCtaText, string(idsJSON), formatTime(time.Now().UTC())); err != nil {
		return HomeContent{}, err
	}
	if _, err := tx.Exec(`UPDATE articles SET featured = 0`); err != nil {
		return HomeContent{}, err
	}
	for _, id := range uniqueClean(home.FeaturedArticleIDs) {
		if _, err := tx.Exec(`UPDATE articles SET featured = 1 WHERE id = ?`, id); err != nil {
			return HomeContent{}, err
		}
	}
	if err := tx.Commit(); err != nil {
		return HomeContent{}, err
	}
	return s.Home()
}

func (s *Store) About() (AboutContent, error) {
	var about AboutContent
	var showHero, showHeroBadge, showHeroTitle, showHeroSubtitle, showHeroDescription int
	var showContact, showLocation, showEmail, showGithub int
	var showContactBadge, showContactTitle, showContactSubtitle, showContactDescription int
	var showSkills, showSkillsHeader, showSkillsBadge, showSkillsTitle, showSkillsSubtitle, showSkillsDescription int
	var showTimeline, showTimelineHeader, showTimelineBadge, showTimelineTitle, showTimelineSubtitle, showTimelineDescription int
	var updated string
	err := s.db.QueryRow(`SELECT `+strings.Join(aboutProfileColumns(), ",")+` FROM about_profile WHERE id = 1`).
		Scan(
			&about.Name, &about.Title, &about.AvatarURL,
			&showHero, &showHeroBadge, &showHeroTitle, &showHeroSubtitle, &showHeroDescription,
			&about.HeroBadge, &about.HeroTitle, &about.HeroSubtitle, &about.Description, &about.RichDescription,
			&showContact, &showLocation, &showEmail, &showGithub,
			&showContactBadge, &showContactTitle, &showContactSubtitle, &showContactDescription,
			&about.ContactBadge, &about.ContactTitle, &about.ContactSubtitle, &about.ContactDescription,
			&showSkills, &showSkillsHeader, &showSkillsBadge, &showSkillsTitle, &showSkillsSubtitle, &showSkillsDescription,
			&about.SkillsBadge, &about.SkillsTitle, &about.SkillsSubtitle, &about.SkillsDescription,
			&showTimeline, &showTimelineHeader, &showTimelineBadge, &showTimelineTitle, &showTimelineSubtitle, &showTimelineDescription,
			&about.TimelineBadge, &about.TimelineTitle, &about.TimelineSubtitle, &about.TimelineDescription,
			&about.City, &about.Contact.Email, &about.Contact.GitHub, &updated,
		)
	if err != nil {
		return AboutContent{}, err
	}
	about.ShowHero = showHero == 1
	about.ShowHeroBadge = showHeroBadge == 1
	about.ShowHeroTitle = showHeroTitle == 1
	about.ShowHeroSubtitle = showHeroSubtitle == 1
	about.ShowHeroDescription = showHeroDescription == 1
	about.ShowContact = showContact == 1
	about.ShowLocation = showLocation == 1
	about.ShowEmail = showEmail == 1
	about.ShowGithub = showGithub == 1
	about.ShowContactBadge = showContactBadge == 1
	about.ShowContactTitle = showContactTitle == 1
	about.ShowContactSubtitle = showContactSubtitle == 1
	about.ShowContactDescription = showContactDescription == 1
	about.ShowSkills = showSkills == 1
	about.ShowSkillsHeader = showSkillsHeader == 1
	about.ShowSkillsBadge = showSkillsBadge == 1
	about.ShowSkillsTitle = showSkillsTitle == 1
	about.ShowSkillsSubtitle = showSkillsSubtitle == 1
	about.ShowSkillsDescription = showSkillsDescription == 1
	about.ShowTimeline = showTimeline == 1
	about.ShowTimelineHeader = showTimelineHeader == 1
	about.ShowTimelineBadge = showTimelineBadge == 1
	about.ShowTimelineTitle = showTimelineTitle == 1
	about.ShowTimelineSubtitle = showTimelineSubtitle == 1
	about.ShowTimelineDescription = showTimelineDescription == 1
	applyAboutDefaults(&about)
	about.WorkStack, err = s.Skills()
	if err != nil {
		return AboutContent{}, err
	}
	about.Experience, err = s.Timeline()
	return about, err
}

func (s *Store) SaveAbout(about AboutContent) (AboutContent, error) {
	applyAboutDefaults(&about)
	if _, err := s.db.Exec(aboutUpsertSQL(), aboutProfileValues(about, time.Now().UTC())...); err != nil {
		return AboutContent{}, err
	}
	if about.WorkStack != nil {
		if err := s.SaveSkills(about.WorkStack); err != nil {
			return AboutContent{}, err
		}
	}
	if about.Experience != nil {
		if err := s.SaveTimeline(about.Experience); err != nil {
			return AboutContent{}, err
		}
	}
	return s.About()
}

func applyAboutDefaults(about *AboutContent) {
	if about.Name == "" {
		about.Name = "Liang"
	}
	if about.Title == "" {
		about.Title = "产品经理 Liang"
	}
	if about.HeroBadge == "" {
		about.HeroBadge = "关于我"
	}
	if about.HeroTitle == "" {
		about.HeroTitle = about.Title
	}
	if about.HeroSubtitle == "" {
		about.HeroSubtitle = "一个产品经理的个人网站，记录产品思考、项目复盘和少量技术笔记。"
	}
	if about.ContactBadge == "" {
		about.ContactBadge = "联系我"
	}
	if about.ContactTitle == "" {
		about.ContactTitle = "欢迎交流产品、增长和项目协作"
	}
	if about.SkillsBadge == "" {
		about.SkillsBadge = "技术栈"
	}
	if about.SkillsTitle == "" {
		about.SkillsTitle = "产品经理的工作栈"
	}
	if about.SkillsDescription == "" {
		about.SkillsDescription = "以产品方法为核心，结合设计协作、数据分析和基础技术理解，保证需求判断和落地过程都更稳。"
	}
	if about.TimelineBadge == "" {
		about.TimelineBadge = "路线图"
	}
	if about.TimelineTitle == "" {
		about.TimelineTitle = "工作经历"
	}
}

func applyAboutVisibilityDefaults(about *AboutContent) {
	about.ShowHero = true
	about.ShowHeroBadge = true
	about.ShowHeroTitle = true
	about.ShowHeroSubtitle = true
	about.ShowHeroDescription = true
	about.ShowContact = true
	about.ShowLocation = true
	about.ShowEmail = true
	about.ShowGithub = true
	about.ShowContactBadge = true
	about.ShowContactTitle = true
	about.ShowContactSubtitle = true
	about.ShowContactDescription = true
	about.ShowSkills = true
	about.ShowSkillsHeader = true
	about.ShowSkillsBadge = true
	about.ShowSkillsTitle = true
	about.ShowSkillsSubtitle = true
	about.ShowSkillsDescription = true
	about.ShowTimeline = true
	about.ShowTimelineHeader = true
	about.ShowTimelineBadge = true
	about.ShowTimelineTitle = true
	about.ShowTimelineSubtitle = true
	about.ShowTimelineDescription = true
}

func (s *Store) Skills() ([]WorkStackGroup, error) {
	rows, err := s.db.Query(`SELECT id,title,items_json,sort FROM skills ORDER BY sort ASC, id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]WorkStackGroup, 0)
	for rows.Next() {
		var item WorkStackGroup
		var itemsJSON string
		if err := rows.Scan(&item.ID, &item.Title, &itemsJSON, &item.Sort); err != nil {
			return nil, err
		}
		_ = json.Unmarshal([]byte(itemsJSON), &item.Items)
		items = append(items, item)
	}
	return items, nil
}

func (s *Store) SaveSkills(groups []WorkStackGroup) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`DELETE FROM skills`); err != nil {
		return err
	}
	now := time.Now().UTC()
	for index, group := range groups {
		if err := insertSkill(tx, group, index, now); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) Timeline() ([]ExperienceItem, error) {
	rows, err := s.db.Query(`SELECT id,period,title,body,sort FROM timeline_items ORDER BY sort ASC, id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]ExperienceItem, 0)
	for rows.Next() {
		var item ExperienceItem
		if err := rows.Scan(&item.ID, &item.Period, &item.Title, &item.Body, &item.Sort); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *Store) SaveTimeline(items []ExperienceItem) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`DELETE FROM timeline_items`); err != nil {
		return err
	}
	now := time.Now().UTC()
	for index, item := range items {
		if err := insertTimeline(tx, item, index, now); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) Tags() ([]string, error) {
	rows, err := s.db.Query(`SELECT name FROM tags ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tags := make([]string, 0)
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (s *Store) TagsWithUsage() ([]map[string]any, error) {
	rows, err := s.db.Query(`SELECT name,use_count,created_at FROM tags ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := make([]map[string]any, 0)
	for rows.Next() {
		var name, createdAt string
		var useCount int
		if err := rows.Scan(&name, &useCount, &createdAt); err != nil {
			return nil, err
		}
		list = append(list, map[string]any{"name": name, "useCount": useCount, "createdAt": createdAt})
	}
	return list, nil
}

func (s *Store) SaveTags(tags []string) ([]string, error) {
	articles, _, err := s.ListArticles(ArticleQuery{Page: 1, PageSize: 10000, IncludeDraft: true})
	if err != nil {
		return nil, err
	}
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	if err := replaceTags(tx, tags, articles, time.Now().UTC()); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return s.Tags()
}

func (s *Store) Dashboard() (Dashboard, error) {
	articles, _, err := s.ListArticles(ArticleQuery{Page: 1, PageSize: 5, IncludeDraft: true})
	if err != nil {
		return Dashboard{}, err
	}
	var totalArticles, totalTags, totalViews, todayViews int
	_ = s.db.QueryRow(`SELECT COUNT(*) FROM articles`).Scan(&totalArticles)
	_ = s.db.QueryRow(`SELECT COUNT(*) FROM tags`).Scan(&totalTags)
	_ = s.db.QueryRow(`SELECT COALESCE(SUM(view_count),0) FROM articles`).Scan(&totalViews)
	_ = s.db.QueryRow(`SELECT COALESCE(SUM(views),0) FROM article_view_events WHERE view_date = ?`, time.Now().UTC().Format("2006-01-02")).Scan(&todayViews)
	recent := make([]ArticleSummary, 0, len(articles))
	for _, article := range articles {
		recent = append(recent, article.SummaryView(true))
	}
	trend, err := s.ViewTrend(30)
	if err != nil {
		return Dashboard{}, err
	}
	return Dashboard{TotalArticles: totalArticles, TotalTags: totalTags, TotalViews: totalViews, TodayViews: todayViews, RecentArticles: recent, ViewTrend: trend}, nil
}

func (s *Store) Analytics() (Analytics, error) {
	trend, err := s.ViewTrend(30)
	if err != nil {
		return Analytics{}, err
	}
	rows, err := s.db.Query(`SELECT id,title,slug,category,date,summary,cover_url,tags_json,featured,status,content,seo_title,seo_keywords,seo_description,view_count,created_at,updated_at FROM articles ORDER BY view_count DESC LIMIT 10`)
	if err != nil {
		return Analytics{}, err
	}
	defer rows.Close()
	hot := make([]ArticleSummary, 0)
	for rows.Next() {
		article, err := scanArticle(rows)
		if err != nil {
			return Analytics{}, err
		}
		hot = append(hot, article.SummaryView(true))
	}
	return Analytics{ViewTrend: trend, HotArticles: hot}, nil
}

func (s *Store) ViewTrend(days int) ([]TrendPoint, error) {
	points := make([]TrendPoint, 0, days)
	start := time.Now().UTC().AddDate(0, 0, -days+1)
	viewsByDate := make(map[string]int)
	rows, err := s.db.Query(`SELECT view_date, COALESCE(SUM(views),0) FROM article_view_events WHERE view_date >= ? GROUP BY view_date`, start.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var date string
		var views int
		if err := rows.Scan(&date, &views); err != nil {
			return nil, err
		}
		viewsByDate[date] = views
	}
	for i := 0; i < days; i++ {
		date := start.AddDate(0, 0, i).Format("2006-01-02")
		points = append(points, TrendPoint{Date: date, Views: viewsByDate[date]})
	}
	return points, nil
}

func (s *Store) SaveUpload(url, kind, filename string) error {
	_, err := s.db.Exec(`INSERT INTO uploads (url,kind,filename,created_at) VALUES (?,?,?,?)`, url, kind, filename, formatTime(time.Now().UTC()))
	return err
}

func upsertSite(tx *sql.Tx, site SiteSettings, now time.Time) error {
	_, err := tx.Exec(`INSERT INTO site_settings (id,name,site_title,description,keywords,logo_url,favicon_url,icp,analytics_code,role,location,email,github,updated_at) VALUES (1,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		site.Name, site.SiteTitle, site.Description, site.Keywords, site.LogoURL, site.FaviconURL, site.ICP, site.AnalyticsCode, site.Role, site.Location, site.Contact.Email, site.Contact.GitHub, formatTime(now))
	return err
}

func upsertHome(tx *sql.Tx, home HomeContent, now time.Time) error {
	idsJSON, _ := json.Marshal(home.FeaturedArticleIDs)
	_, err := tx.Exec(`INSERT INTO home_settings (id,intro_title,subtitle,intro,primary_cta_text,secondary_cta_text,featured_article_ids_json,updated_at) VALUES (1,?,?,?,?,?,?,?)`,
		home.IntroTitle, home.Subtitle, home.Intro, home.PrimaryCtaText, home.SecondaryCtaText, string(idsJSON), formatTime(now))
	return err
}

func aboutProfileColumns() []string {
	return []string{
		"name", "title", "avatar_url",
		"show_hero", "show_hero_badge", "show_hero_title", "show_hero_subtitle", "show_hero_description",
		"hero_badge", "hero_title", "hero_subtitle", "description", "rich_description",
		"show_contact", "show_location", "show_email", "show_github",
		"show_contact_badge", "show_contact_title", "show_contact_subtitle", "show_contact_description",
		"contact_badge", "contact_title", "contact_subtitle", "contact_description",
		"show_skills", "show_skills_header", "show_skills_badge", "show_skills_title", "show_skills_subtitle", "show_skills_description",
		"skills_badge", "skills_title", "skills_subtitle", "skills_description",
		"show_timeline", "show_timeline_header", "show_timeline_badge", "show_timeline_title", "show_timeline_subtitle", "show_timeline_description",
		"timeline_badge", "timeline_title", "timeline_subtitle", "timeline_description",
		"city", "email", "github", "updated_at",
	}
}

func aboutUpsertSQL() string {
	columns := aboutProfileColumns()
	placeholders := strings.TrimRight(strings.Repeat("?,", len(columns)), ",")
	assignments := make([]string, 0, len(columns))
	for _, column := range columns {
		assignments = append(assignments, column+"=excluded."+column)
	}
	return `INSERT INTO about_profile (id,` + strings.Join(columns, ",") + `) VALUES (1,` + placeholders + `)
		ON CONFLICT(id) DO UPDATE SET ` + strings.Join(assignments, ",")
}

func aboutProfileValues(about AboutContent, now time.Time) []any {
	return []any{
		about.Name, about.Title, about.AvatarURL,
		boolInt(about.ShowHero), boolInt(about.ShowHeroBadge), boolInt(about.ShowHeroTitle), boolInt(about.ShowHeroSubtitle), boolInt(about.ShowHeroDescription),
		about.HeroBadge, about.HeroTitle, about.HeroSubtitle, about.Description, about.RichDescription,
		boolInt(about.ShowContact), boolInt(about.ShowLocation), boolInt(about.ShowEmail), boolInt(about.ShowGithub),
		boolInt(about.ShowContactBadge), boolInt(about.ShowContactTitle), boolInt(about.ShowContactSubtitle), boolInt(about.ShowContactDescription),
		about.ContactBadge, about.ContactTitle, about.ContactSubtitle, about.ContactDescription,
		boolInt(about.ShowSkills), boolInt(about.ShowSkillsHeader), boolInt(about.ShowSkillsBadge), boolInt(about.ShowSkillsTitle), boolInt(about.ShowSkillsSubtitle), boolInt(about.ShowSkillsDescription),
		about.SkillsBadge, about.SkillsTitle, about.SkillsSubtitle, about.SkillsDescription,
		boolInt(about.ShowTimeline), boolInt(about.ShowTimelineHeader), boolInt(about.ShowTimelineBadge), boolInt(about.ShowTimelineTitle), boolInt(about.ShowTimelineSubtitle), boolInt(about.ShowTimelineDescription),
		about.TimelineBadge, about.TimelineTitle, about.TimelineSubtitle, about.TimelineDescription,
		about.City, about.Contact.Email, about.Contact.GitHub, formatTime(now),
	}
}

func upsertAbout(tx *sql.Tx, about AboutContent, now time.Time) error {
	applyAboutDefaults(&about)
	_, err := tx.Exec(aboutUpsertSQL(), aboutProfileValues(about, now)...)
	return err
}

func upsertArticle(tx *sql.Tx, article Article, now time.Time) error {
	article.Status = normalizeArticleStatus(article.Status, article.Draft)
	if article.Category == "" && len(article.Tags) > 0 {
		article.Category = article.Tags[0]
	}
	if article.Category != "" {
		article.Tags = []string{article.Category}
	}
	tagsJSON, _ := json.Marshal(uniqueClean(article.Tags))
	_, err := tx.Exec(`INSERT INTO articles (id,title,slug,category,date,summary,cover_url,tags_json,featured,status,content,seo_title,seo_keywords,seo_description,view_count,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		article.ID, article.Title, article.Slug, article.Category, article.Date, article.Summary, article.CoverURL, string(tagsJSON), boolInt(article.Featured), article.Status, article.Content, article.SEOTitle, article.SEOKeywords, article.SEODescription, article.ViewCount, formatTime(article.CreatedAt), formatTime(article.UpdatedAt))
	return err
}

func insertSkill(tx *sql.Tx, group WorkStackGroup, index int, now time.Time) error {
	itemsJSON, _ := json.Marshal(uniqueClean(group.Items))
	_, err := tx.Exec(`INSERT INTO skills (title,items_json,sort,created_at,updated_at) VALUES (?,?,?,?,?)`, group.Title, string(itemsJSON), index, formatTime(now), formatTime(now))
	return err
}

func insertTimeline(tx *sql.Tx, item ExperienceItem, index int, now time.Time) error {
	_, err := tx.Exec(`INSERT INTO timeline_items (period,title,body,sort,created_at,updated_at) VALUES (?,?,?,?,?,?)`, item.Period, item.Title, item.Body, index, formatTime(now), formatTime(now))
	return err
}

func replaceTags(tx *sql.Tx, tags []string, articles []Article, now time.Time) error {
	counts := make(map[string]int)
	for _, article := range articles {
		if article.Category != "" {
			counts[article.Category]++
		}
	}
	tags = uniqueClean(tags)
	for tag := range counts {
		tags = uniqueClean(append(tags, tag))
	}
	sort.Strings(tags)
	if _, err := tx.Exec(`DELETE FROM tags`); err != nil {
		return err
	}
	for _, tag := range tags {
		if _, err := tx.Exec(`INSERT INTO tags (name,use_count,created_at,updated_at) VALUES (?,?,?,?)`, tag, counts[tag], formatTime(now), formatTime(now)); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) rebuildTagUseCounts() error {
	tags, _ := s.Tags()
	articles, _, err := s.ListArticles(ArticleQuery{Page: 1, PageSize: 10000, IncludeDraft: true})
	if err != nil {
		return err
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := replaceTags(tx, tags, articles, time.Now().UTC()); err != nil {
		return err
	}
	return tx.Commit()
}

type scanner interface {
	Scan(dest ...any) error
}

func scanArticle(row scanner) (Article, error) {
	var article Article
	var tagsJSON string
	var featured int
	var createdAt, updatedAt string
	err := row.Scan(&article.ID, &article.Title, &article.Slug, &article.Category, &article.Date, &article.Summary, &article.CoverURL, &tagsJSON, &featured, &article.Status, &article.Content, &article.SEOTitle, &article.SEOKeywords, &article.SEODescription, &article.ViewCount, &createdAt, &updatedAt)
	if err != nil {
		return Article{}, err
	}
	_ = json.Unmarshal([]byte(tagsJSON), &article.Tags)
	article.Featured = featured == 1
	article.Draft = article.Status == ArticleDraft
	article.CreatedAt = parseTime(createdAt)
	article.UpdatedAt = parseTime(updatedAt)
	article.WordCount = estimateWords(article.Title + " " + article.Summary + " " + article.Content)
	return article, nil
}

func formatTime(value time.Time) string {
	if value.IsZero() {
		value = time.Now().UTC()
	}
	return value.UTC().Format(time.RFC3339Nano)
}

func parseTime(value string) time.Time {
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return time.Time{}
	}
	return parsed
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
