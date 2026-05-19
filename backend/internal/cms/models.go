package cms

import (
	"regexp"
	"strings"
	"time"
	"unicode"
)

const (
	RoleSuperAdmin = "super_admin"
	RoleAdmin      = "admin"
	RoleEditor     = "editor"

	UserActive   = "active"
	UserDisabled = "disabled"

	ArticleDraft     = "draft"
	ArticlePublished = "published"
)

type APIResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type PageResult struct {
	List  any `json:"list"`
	Total int `json:"total"`
}

type Content struct {
	Site             SiteSettings     `json:"site"`
	Home             HomeContent      `json:"home"`
	Tags             []string         `json:"tags"`
	Articles         []Article        `json:"articles"`
	About            AboutContent     `json:"about"`
	UpdatedAt        time.Time        `json:"updatedAt"`
	FeaturedArticles []ArticleSummary `json:"featuredArticles,omitempty"`
}

type SiteSettings struct {
	Name          string  `json:"name"`
	SiteTitle     string  `json:"siteTitle"`
	Description   string  `json:"description"`
	Keywords      string  `json:"keywords"`
	LogoURL       string  `json:"logoUrl"`
	FaviconURL    string  `json:"faviconUrl"`
	ICP           string  `json:"icp"`
	AnalyticsCode string  `json:"analyticsCode"`
	Role          string  `json:"role"`
	Location      string  `json:"location"`
	Contact       Contact `json:"contact"`
}

type Contact struct {
	Email  string `json:"email"`
	GitHub string `json:"github"`
}

type HomeContent struct {
	IntroTitle         string   `json:"introTitle"`
	Subtitle           string   `json:"subtitle"`
	Intro              string   `json:"intro"`
	PrimaryCtaText     string   `json:"primaryCtaText"`
	SecondaryCtaText   string   `json:"secondaryCtaText"`
	FeaturedArticleIDs []string `json:"featuredArticleIds"`
}

type AboutContent struct {
	Name                    string           `json:"name"`
	Title                   string           `json:"title"`
	AvatarURL               string           `json:"avatarUrl"`
	ShowHero                bool             `json:"showHero"`
	ShowHeroBadge           bool             `json:"showHeroBadge"`
	ShowHeroTitle           bool             `json:"showHeroTitle"`
	ShowHeroSubtitle        bool             `json:"showHeroSubtitle"`
	ShowHeroDescription     bool             `json:"showHeroDescription"`
	HeroBadge               string           `json:"heroBadge"`
	HeroTitle               string           `json:"heroTitle"`
	HeroSubtitle            string           `json:"heroSubtitle"`
	Description             string           `json:"description"`
	RichDescription         string           `json:"richDescription"`
	ShowContact             bool             `json:"showContact"`
	ShowLocation            bool             `json:"showLocation"`
	ShowEmail               bool             `json:"showEmail"`
	ShowGithub              bool             `json:"showGithub"`
	ShowContactBadge        bool             `json:"showContactBadge"`
	ShowContactTitle        bool             `json:"showContactTitle"`
	ShowContactSubtitle     bool             `json:"showContactSubtitle"`
	ShowContactDescription  bool             `json:"showContactDescription"`
	ContactBadge            string           `json:"contactBadge"`
	ContactTitle            string           `json:"contactTitle"`
	ContactSubtitle         string           `json:"contactSubtitle"`
	ContactDescription      string           `json:"contactDescription"`
	ShowSkills              bool             `json:"showSkills"`
	ShowSkillsHeader        bool             `json:"showSkillsHeader"`
	ShowSkillsBadge         bool             `json:"showSkillsBadge"`
	ShowSkillsTitle         bool             `json:"showSkillsTitle"`
	ShowSkillsSubtitle      bool             `json:"showSkillsSubtitle"`
	ShowSkillsDescription   bool             `json:"showSkillsDescription"`
	SkillsBadge             string           `json:"skillsBadge"`
	SkillsTitle             string           `json:"skillsTitle"`
	SkillsSubtitle          string           `json:"skillsSubtitle"`
	SkillsDescription       string           `json:"skillsDescription"`
	ShowTimeline            bool             `json:"showTimeline"`
	ShowTimelineHeader      bool             `json:"showTimelineHeader"`
	ShowTimelineBadge       bool             `json:"showTimelineBadge"`
	ShowTimelineTitle       bool             `json:"showTimelineTitle"`
	ShowTimelineSubtitle    bool             `json:"showTimelineSubtitle"`
	ShowTimelineDescription bool             `json:"showTimelineDescription"`
	TimelineBadge           string           `json:"timelineBadge"`
	TimelineTitle           string           `json:"timelineTitle"`
	TimelineSubtitle        string           `json:"timelineSubtitle"`
	TimelineDescription     string           `json:"timelineDescription"`
	Contact                 Contact          `json:"contact"`
	City                    string           `json:"city"`
	WorkStack               []WorkStackGroup `json:"workStack"`
	Experience              []ExperienceItem `json:"experience"`
}

type WorkStackGroup struct {
	ID    int64    `json:"id,omitempty"`
	Title string   `json:"title"`
	Items []string `json:"items"`
	Sort  int      `json:"sort,omitempty"`
}

type ExperienceItem struct {
	ID     int64  `json:"id,omitempty"`
	Period string `json:"period"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	Sort   int    `json:"sort,omitempty"`
}

type Article struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Slug           string    `json:"slug"`
	Category       string    `json:"category"`
	Date           string    `json:"date"`
	Summary        string    `json:"summary"`
	CoverURL       string    `json:"coverUrl"`
	Tags           []string  `json:"tags"`
	Featured       bool      `json:"featured"`
	Draft          bool      `json:"draft"`
	Status         string    `json:"status"`
	Content        string    `json:"content"`
	SEOTitle       string    `json:"seoTitle"`
	SEOKeywords    string    `json:"seoKeywords"`
	SEODescription string    `json:"seoDescription"`
	ViewCount      int       `json:"viewCount"`
	WordCount      int       `json:"wordCount"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type ArticleSummary struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Slug      string   `json:"slug"`
	Category  string   `json:"category"`
	Date      string   `json:"date"`
	Summary   string   `json:"summary"`
	CoverURL  string   `json:"coverUrl"`
	Tags      []string `json:"tags"`
	Featured  bool     `json:"featured"`
	Draft     bool     `json:"draft,omitempty"`
	Status    string   `json:"status"`
	WordCount int      `json:"wordCount"`
	ViewCount int      `json:"viewCount"`
}

type User struct {
	ID          int64      `json:"id"`
	Username    string     `json:"username"`
	Nickname    string     `json:"nickname"`
	Email       string     `json:"email"`
	Role        string     `json:"role"`
	Status      string     `json:"status"`
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type LoginLog struct {
	ID        int64     `json:"id"`
	UserID    *int64    `json:"userId,omitempty"`
	Username  string    `json:"username"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"userAgent"`
	Success   bool      `json:"success"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"createdAt"`
}

type Dashboard struct {
	TotalArticles  int              `json:"totalArticles"`
	TotalTags      int              `json:"totalTags"`
	TotalViews     int              `json:"totalViews"`
	TodayViews     int              `json:"todayViews"`
	RecentArticles []ArticleSummary `json:"recentArticles"`
	ViewTrend      []TrendPoint     `json:"viewTrend"`
}

type Analytics struct {
	ViewTrend   []TrendPoint     `json:"viewTrend"`
	HotArticles []ArticleSummary `json:"hotArticles"`
}

type TrendPoint struct {
	Date  string `json:"date"`
	Views int    `json:"views"`
}

func (a Article) SummaryView(includeDraft bool) ArticleSummary {
	status := normalizeArticleStatus(a.Status, a.Draft)
	summary := ArticleSummary{
		ID:        a.ID,
		Title:     a.Title,
		Slug:      a.Slug,
		Category:  a.Category,
		Date:      a.Date,
		Summary:   a.Summary,
		CoverURL:  a.CoverURL,
		Tags:      append([]string(nil), a.Tags...),
		Featured:  a.Featured,
		Status:    status,
		WordCount: estimateWords(a.Title + " " + a.Summary + " " + strings.Join(a.Tags, " ") + " " + a.Content),
		ViewCount: a.ViewCount,
	}
	if includeDraft {
		summary.Draft = status == ArticleDraft
	}
	return summary
}

func normalizeArticleStatus(status string, draft bool) string {
	status = strings.TrimSpace(status)
	if status == ArticlePublished || (!draft && status == "") {
		return ArticlePublished
	}
	return ArticleDraft
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
