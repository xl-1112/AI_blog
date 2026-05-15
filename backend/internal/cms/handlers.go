package cms

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Config struct {
	AdminToken    string
	UploadDir     string
	CORSOrigins   []string
	MaxUploadSize int64
}

type Server struct {
	store         *Store
	adminToken    string
	uploadDir     string
	corsOrigins   map[string]bool
	maxUploadSize int64
}

func NewServer(store *Store, config Config) http.Handler {
	if config.MaxUploadSize <= 0 {
		config.MaxUploadSize = 5 << 20
	}

	server := &Server{
		store:         store,
		adminToken:    config.AdminToken,
		uploadDir:     config.UploadDir,
		corsOrigins:   make(map[string]bool),
		maxUploadSize: config.MaxUploadSize,
	}
	for _, origin := range config.CORSOrigins {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			server.corsOrigins[origin] = true
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", server.handleHealth)
	mux.HandleFunc("/api/site", server.handlePublicSite)
	mux.HandleFunc("/api/tags", server.handlePublicTags)
	mux.HandleFunc("/api/articles", server.handlePublicArticles)
	mux.HandleFunc("/api/articles/", server.handlePublicArticle)
	mux.HandleFunc("/api/admin/content", server.handleAdminContent)
	mux.HandleFunc("/api/admin/site", server.handleAdminSite)
	mux.HandleFunc("/api/admin/home", server.handleAdminHome)
	mux.HandleFunc("/api/admin/about", server.handleAdminAbout)
	mux.HandleFunc("/api/admin/tags", server.handleAdminTags)
	mux.HandleFunc("/api/admin/logo", server.handleAdminLogo)
	mux.HandleFunc("/api/admin/articles", server.handleAdminArticles)
	mux.HandleFunc("/api/admin/articles/", server.handleAdminArticle)
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(config.UploadDir))))

	return server.withCORS(mux)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, http.MethodGet)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "time": time.Now().UTC()})
}

func (s *Server) handlePublicSite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, http.MethodGet)
		return
	}
	writeJSON(w, http.StatusOK, s.store.Snapshot().Public())
}

func (s *Server) handlePublicTags(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, http.MethodGet)
		return
	}
	content := s.store.Snapshot().Public()
	writeJSON(w, http.StatusOK, map[string]any{"tags": content.Tags})
}

func (s *Server) handlePublicArticles(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/articles" {
		writeNotFound(w)
		return
	}
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, http.MethodGet)
		return
	}

	content := s.store.Snapshot().Public()
	articles := make([]ArticleSummary, 0, len(content.Articles))
	for _, article := range content.Articles {
		articles = append(articles, article.SummaryView(false))
	}
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].Date > articles[j].Date
	})
	writeJSON(w, http.StatusOK, map[string]any{"articles": articles})
}

func (s *Server) handlePublicArticle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, http.MethodGet)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/articles/")
	if id == "" {
		writeNotFound(w)
		return
	}
	for _, article := range s.store.Snapshot().Articles {
		if article.ID == id && !article.Draft {
			writeJSON(w, http.StatusOK, article)
			return
		}
	}
	writeNotFound(w)
}

func (s *Server) handleAdminContent(w http.ResponseWriter, r *http.Request) {
	if !s.requireAdmin(w, r) {
		return
	}

	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, s.store.Snapshot())
	case http.MethodPut:
		var content Content
		if err := readJSON(r, &content); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		next, err := s.store.Replace(content)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, next)
	default:
		writeMethodNotAllowed(w, http.MethodGet, http.MethodPut)
	}
}

func (s *Server) handleAdminSite(w http.ResponseWriter, r *http.Request) {
	if !s.requireAdmin(w, r) {
		return
	}
	if r.Method != http.MethodPut {
		writeMethodNotAllowed(w, http.MethodPut)
		return
	}

	var site SiteSettings
	if err := readJSON(r, &site); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	next, err := s.store.Update(func(content *Content) error {
		content.Site = site
		return nil
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, next.Site)
}

func (s *Server) handleAdminHome(w http.ResponseWriter, r *http.Request) {
	if !s.requireAdmin(w, r) {
		return
	}
	if r.Method != http.MethodPut {
		writeMethodNotAllowed(w, http.MethodPut)
		return
	}

	var home HomeContent
	if err := readJSON(r, &home); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	next, err := s.store.Update(func(content *Content) error {
		content.Home = home
		return nil
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, next.Home)
}

func (s *Server) handleAdminAbout(w http.ResponseWriter, r *http.Request) {
	if !s.requireAdmin(w, r) {
		return
	}
	if r.Method != http.MethodPut {
		writeMethodNotAllowed(w, http.MethodPut)
		return
	}

	var about AboutContent
	if err := readJSON(r, &about); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	next, err := s.store.Update(func(content *Content) error {
		content.About = about
		return nil
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, next.About)
}

func (s *Server) handleAdminTags(w http.ResponseWriter, r *http.Request) {
	if !s.requireAdmin(w, r) {
		return
	}
	if r.Method != http.MethodPut {
		writeMethodNotAllowed(w, http.MethodPut)
		return
	}

	var request struct {
		Tags []string `json:"tags"`
	}
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	next, err := s.store.Update(func(content *Content) error {
		content.Tags = request.Tags
		return nil
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"tags": next.Tags})
}

func (s *Server) handleAdminLogo(w http.ResponseWriter, r *http.Request) {
	if !s.requireAdmin(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w, http.MethodPost)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, s.maxUploadSize)
	file, header, err := r.FormFile("logo")
	if err != nil {
		writeError(w, http.StatusBadRequest, "logo file is required")
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" {
		if exts, err := mime.ExtensionsByType(header.Header.Get("Content-Type")); err == nil && len(exts) > 0 {
			ext = exts[0]
		}
	}
	if !allowedImageExt(ext) {
		writeError(w, http.StatusBadRequest, "logo must be png, jpg, jpeg, webp, svg, or gif")
		return
	}

	if err := os.MkdirAll(s.uploadDir, 0755); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	filename := "logo-" + time.Now().UTC().Format("20060102150405") + ext
	destinationPath := filepath.Join(s.uploadDir, filename)
	destination, err := os.Create(destinationPath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer destination.Close()

	if _, err := io.Copy(destination, file); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	logoURL := "/uploads/" + filename
	next, err := s.store.Update(func(content *Content) error {
		content.Site.LogoURL = logoURL
		return nil
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"logoUrl": logoURL, "site": next.Site})
}

func (s *Server) handleAdminArticles(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/admin/articles" {
		writeNotFound(w)
		return
	}
	if !s.requireAdmin(w, r) {
		return
	}

	switch r.Method {
	case http.MethodGet:
		content := s.store.Snapshot()
		articles := make([]ArticleSummary, 0, len(content.Articles))
		for _, article := range content.Articles {
			articles = append(articles, article.SummaryView(true))
		}
		sort.Slice(articles, func(i, j int) bool {
			return articles[i].Date > articles[j].Date
		})
		writeJSON(w, http.StatusOK, map[string]any{"articles": articles})
	case http.MethodPost:
		var article Article
		if err := readJSON(r, &article); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		var created Article
		_, err := s.store.Update(func(content *Content) error {
			if article.ID == "" {
				article.ID = slugify(article.Title)
			}
			if article.ID == "" {
				article.ID = "article-" + time.Now().UTC().Format("20060102150405")
			}
			for _, existing := range content.Articles {
				if existing.ID == article.ID {
					return errors.New("article already exists: " + article.ID)
				}
			}
			now := time.Now().UTC()
			article.CreatedAt = now
			article.UpdatedAt = now
			content.Articles = append(content.Articles, article)
			if article.Featured {
				content.Home.FeaturedArticleIDs = append(content.Home.FeaturedArticleIDs, article.ID)
			}
			created = article
			return nil
		})
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, created)
	default:
		writeMethodNotAllowed(w, http.MethodGet, http.MethodPost)
	}
}

func (s *Server) handleAdminArticle(w http.ResponseWriter, r *http.Request) {
	if !s.requireAdmin(w, r) {
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/admin/articles/")
	id = strings.TrimSpace(id)
	if id == "" {
		writeNotFound(w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		for _, article := range s.store.Snapshot().Articles {
			if article.ID == id {
				writeJSON(w, http.StatusOK, article)
				return
			}
		}
		writeNotFound(w)
	case http.MethodPut:
		var article Article
		if err := readJSON(r, &article); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		article.ID = id
		var updated Article
		_, err := s.store.Update(func(content *Content) error {
			for index, existing := range content.Articles {
				if existing.ID == id {
					article.CreatedAt = existing.CreatedAt
					article.UpdatedAt = time.Now().UTC()
					content.Articles[index] = article
					if article.Featured {
						content.Home.FeaturedArticleIDs = append(content.Home.FeaturedArticleIDs, article.ID)
					} else {
						content.Home.FeaturedArticleIDs = removeString(content.Home.FeaturedArticleIDs, article.ID)
					}
					updated = article
					return nil
				}
			}
			return errors.New("article not found: " + id)
		})
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, updated)
	case http.MethodDelete:
		_, err := s.store.Update(func(content *Content) error {
			next := content.Articles[:0]
			found := false
			for _, article := range content.Articles {
				if article.ID == id {
					found = true
					continue
				}
				next = append(next, article)
			}
			if !found {
				return errors.New("article not found: " + id)
			}
			content.Articles = next
			content.Home.FeaturedArticleIDs = removeString(content.Home.FeaturedArticleIDs, id)
			return nil
		})
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		writeMethodNotAllowed(w, http.MethodGet, http.MethodPut, http.MethodDelete)
	}
}

func (s *Server) withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && (s.corsOrigins["*"] || s.corsOrigins[origin]) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Admin-Token")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) requireAdmin(w http.ResponseWriter, r *http.Request) bool {
	token := bearerToken(r.Header.Get("Authorization"))
	if token == "" {
		token = r.Header.Get("X-Admin-Token")
	}
	if token == "" || subtle.ConstantTimeCompare([]byte(token), []byte(s.adminToken)) != 1 {
		writeError(w, http.StatusUnauthorized, "admin token is required")
		return false
	}
	return true
}

func bearerToken(header string) string {
	if header == "" {
		return ""
	}
	parts := strings.Fields(header)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}
	return ""
}

func readJSON(r *http.Request, destination any) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(io.LimitReader(r.Body, 4<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeMethodNotAllowed(w http.ResponseWriter, methods ...string) {
	w.Header().Set("Allow", strings.Join(methods, ", "))
	writeError(w, http.StatusMethodNotAllowed, "method not allowed")
}

func writeNotFound(w http.ResponseWriter) {
	writeError(w, http.StatusNotFound, "not found")
}

func allowedImageExt(ext string) bool {
	switch ext {
	case ".png", ".jpg", ".jpeg", ".webp", ".svg", ".gif":
		return true
	default:
		return false
	}
}

func removeString(values []string, target string) []string {
	next := values[:0]
	for _, value := range values {
		if value != target {
			next = append(next, value)
		}
	}
	return next
}
