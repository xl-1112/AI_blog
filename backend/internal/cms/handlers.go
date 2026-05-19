package cms

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	JWTSecret     string
	UploadDir     string
	CORSOrigins   []string
	MaxUploadSize int64
}

type Server struct {
	store         *Store
	jwtSecret     []byte
	uploadDir     string
	corsOrigins   map[string]bool
	maxUploadSize int64
}

type contextKey string

const userContextKey contextKey = "cms-user"

func NewServer(store *Store, config Config) http.Handler {
	if config.MaxUploadSize <= 0 {
		config.MaxUploadSize = 5 << 20
	}
	server := &Server{
		store:         store,
		jwtSecret:     []byte(config.JWTSecret),
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

	mux.HandleFunc("/api/admin/login", server.handleLogin)
	mux.HandleFunc("/api/admin/logout", server.withAuth(server.handleLogout, RoleSuperAdmin, RoleAdmin, RoleEditor))
	mux.HandleFunc("/api/admin/me", server.withAuth(server.handleMe, RoleSuperAdmin, RoleAdmin, RoleEditor))
	mux.HandleFunc("/api/admin/dashboard", server.withAuth(server.handleDashboard, RoleSuperAdmin, RoleAdmin, RoleEditor))
	mux.HandleFunc("/api/admin/analytics", server.withAuth(server.handleAnalytics, RoleSuperAdmin, RoleAdmin, RoleEditor))
	mux.HandleFunc("/api/admin/content", server.withAuth(server.handleAdminContent, RoleSuperAdmin, RoleAdmin))
	mux.HandleFunc("/api/admin/site", server.withAuth(server.handleAdminSite, RoleSuperAdmin, RoleAdmin))
	mux.HandleFunc("/api/admin/home", server.withAuth(server.handleAdminHome, RoleSuperAdmin, RoleAdmin))
	mux.HandleFunc("/api/admin/about", server.withAuth(server.handleAdminAbout, RoleSuperAdmin, RoleAdmin))
	mux.HandleFunc("/api/admin/skills", server.withAuth(server.handleAdminSkills, RoleSuperAdmin, RoleAdmin))
	mux.HandleFunc("/api/admin/timeline", server.withAuth(server.handleAdminTimeline, RoleSuperAdmin, RoleAdmin))
	mux.HandleFunc("/api/admin/tags", server.withAuth(server.handleAdminTags, RoleSuperAdmin, RoleAdmin, RoleEditor))
	mux.HandleFunc("/api/admin/upload", server.withAuth(server.handleAdminUpload, RoleSuperAdmin, RoleAdmin, RoleEditor))
	mux.HandleFunc("/api/admin/logo", server.withAuth(server.handleAdminLogo, RoleSuperAdmin, RoleAdmin))
	mux.HandleFunc("/api/admin/articles", server.withAuth(server.handleAdminArticles, RoleSuperAdmin, RoleAdmin, RoleEditor))
	mux.HandleFunc("/api/admin/articles/", server.withAuth(server.handleAdminArticle, RoleSuperAdmin, RoleAdmin, RoleEditor))
	mux.HandleFunc("/api/admin/users", server.withAuth(server.handleAdminUsers, RoleSuperAdmin))
	mux.HandleFunc("/api/admin/users/", server.withAuth(server.handleAdminUser, RoleSuperAdmin))
	mux.HandleFunc("/api/admin/login-logs", server.withAuth(server.handleLoginLogs, RoleSuperAdmin))
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(config.UploadDir))))

	return server.withCORS(mux)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, false, http.MethodGet)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "time": time.Now().UTC()})
}

func (s *Server) handlePublicSite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, false, http.MethodGet)
		return
	}
	content, err := s.store.PublicContent()
	if err != nil {
		writeError(w, false, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, content)
}

func (s *Server) handlePublicTags(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, false, http.MethodGet)
		return
	}
	tags, err := s.store.Tags()
	if err != nil {
		writeError(w, false, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"tags": tags})
}

func (s *Server) handlePublicArticles(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/articles" {
		writeNotFound(w, false)
		return
	}
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, false, http.MethodGet)
		return
	}
	articles, _, err := s.store.ListArticles(ArticleQuery{Page: 1, PageSize: 10000, IncludeDraft: false})
	if err != nil {
		writeError(w, false, http.StatusInternalServerError, err.Error())
		return
	}
	summaries := make([]ArticleSummary, 0, len(articles))
	for _, article := range articles {
		summaries = append(summaries, article.SummaryView(false))
	}
	writeJSON(w, http.StatusOK, map[string]any{"articles": summaries})
}

func (s *Server) handlePublicArticle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, false, http.MethodGet)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/articles/")
	if id == "" {
		writeNotFound(w, false)
		return
	}
	article, err := s.store.Article(id, false)
	if err != nil {
		writeNotFound(w, false)
		return
	}
	_ = s.store.IncrementArticleView(article.ID)
	article.ViewCount++
	writeJSON(w, http.StatusOK, article)
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w, true, http.MethodPost)
		return
	}
	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := readJSON(r, &request); err != nil {
		writeError(w, true, http.StatusBadRequest, err.Error())
		return
	}
	user, err := s.store.Authenticate(request.Username, request.Password, clientIP(r), r.UserAgent())
	if err != nil {
		writeError(w, true, http.StatusUnauthorized, err.Error())
		return
	}
	token, err := s.issueToken(user)
	if err != nil {
		writeError(w, true, http.StatusInternalServerError, err.Error())
		return
	}
	writeAdminOK(w, map[string]any{"token": token, "userInfo": user})
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w, true, http.MethodPost)
		return
	}
	writeAdminOK(w, map[string]any{"ok": true})
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, true, http.MethodGet)
		return
	}
	writeAdminOK(w, currentUser(r))
}

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, true, http.MethodGet)
		return
	}
	data, err := s.store.Dashboard()
	if err != nil {
		writeError(w, true, http.StatusInternalServerError, err.Error())
		return
	}
	writeAdminOK(w, data)
}

func (s *Server) handleAnalytics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, true, http.MethodGet)
		return
	}
	data, err := s.store.Analytics()
	if err != nil {
		writeError(w, true, http.StatusInternalServerError, err.Error())
		return
	}
	writeAdminOK(w, data)
}

func (s *Server) handleAdminContent(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		content, err := s.store.AdminContent()
		if err != nil {
			writeError(w, true, http.StatusInternalServerError, err.Error())
			return
		}
		writeAdminOK(w, content)
	case http.MethodPut:
		var content Content
		if err := readJSON(r, &content); err != nil {
			writeError(w, true, http.StatusBadRequest, err.Error())
			return
		}
		if err := s.store.ReplaceContent(content); err != nil {
			writeError(w, true, http.StatusBadRequest, err.Error())
			return
		}
		next, _ := s.store.AdminContent()
		writeAdminOK(w, next)
	default:
		writeMethodNotAllowed(w, true, http.MethodGet, http.MethodPut)
	}
}

func (s *Server) handleAdminSite(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		site, err := s.store.Site()
		if err != nil {
			writeError(w, true, http.StatusInternalServerError, err.Error())
			return
		}
		writeAdminOK(w, site)
	case http.MethodPut:
		var site SiteSettings
		if err := readJSON(r, &site); err != nil {
			writeError(w, true, http.StatusBadRequest, err.Error())
			return
		}
		next, err := s.store.SaveSite(site)
		if err != nil {
			writeError(w, true, http.StatusBadRequest, err.Error())
			return
		}
		writeAdminOK(w, next)
	default:
		writeMethodNotAllowed(w, true, http.MethodGet, http.MethodPut)
	}
}

func (s *Server) handleAdminHome(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		home, err := s.store.Home()
		if err != nil {
			writeError(w, true, http.StatusInternalServerError, err.Error())
			return
		}
		writeAdminOK(w, home)
	case http.MethodPut:
		var home HomeContent
		if err := readJSON(r, &home); err != nil {
			writeError(w, true, http.StatusBadRequest, err.Error())
			return
		}
		next, err := s.store.SaveHome(home)
		if err != nil {
			writeError(w, true, http.StatusBadRequest, err.Error())
			return
		}
		writeAdminOK(w, next)
	default:
		writeMethodNotAllowed(w, true, http.MethodGet, http.MethodPut)
	}
}

func (s *Server) handleAdminAbout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		about, err := s.store.About()
		if err != nil {
			writeError(w, true, http.StatusInternalServerError, err.Error())
			return
		}
		writeAdminOK(w, about)
	case http.MethodPut:
		var about AboutContent
		if err := readJSON(r, &about); err != nil {
			writeError(w, true, http.StatusBadRequest, err.Error())
			return
		}
		next, err := s.store.SaveAbout(about)
		if err != nil {
			writeError(w, true, http.StatusBadRequest, err.Error())
			return
		}
		writeAdminOK(w, next)
	default:
		writeMethodNotAllowed(w, true, http.MethodGet, http.MethodPut)
	}
}

func (s *Server) handleAdminSkills(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		groups, err := s.store.Skills()
		if err != nil {
			writeError(w, true, http.StatusInternalServerError, err.Error())
			return
		}
		writeAdminOK(w, groups)
	case http.MethodPut:
		var groups []WorkStackGroup
		if err := readJSON(r, &groups); err != nil {
			writeError(w, true, http.StatusBadRequest, err.Error())
			return
		}
		if err := s.store.SaveSkills(groups); err != nil {
			writeError(w, true, http.StatusBadRequest, err.Error())
			return
		}
		next, _ := s.store.Skills()
		writeAdminOK(w, next)
	default:
		writeMethodNotAllowed(w, true, http.MethodGet, http.MethodPut)
	}
}

func (s *Server) handleAdminTimeline(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		items, err := s.store.Timeline()
		if err != nil {
			writeError(w, true, http.StatusInternalServerError, err.Error())
			return
		}
		writeAdminOK(w, items)
	case http.MethodPut:
		var items []ExperienceItem
		if err := readJSON(r, &items); err != nil {
			writeError(w, true, http.StatusBadRequest, err.Error())
			return
		}
		if err := s.store.SaveTimeline(items); err != nil {
			writeError(w, true, http.StatusBadRequest, err.Error())
			return
		}
		next, _ := s.store.Timeline()
		writeAdminOK(w, next)
	default:
		writeMethodNotAllowed(w, true, http.MethodGet, http.MethodPut)
	}
}

func (s *Server) handleAdminTags(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		list, err := s.store.TagsWithUsage()
		if err != nil {
			writeError(w, true, http.StatusInternalServerError, err.Error())
			return
		}
		writeAdminOK(w, PageResult{List: list, Total: len(list)})
	case http.MethodPut:
		var request struct {
			Tags []string `json:"tags"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, true, http.StatusBadRequest, err.Error())
			return
		}
		tags, err := s.store.SaveTags(request.Tags)
		if err != nil {
			writeError(w, true, http.StatusBadRequest, err.Error())
			return
		}
		writeAdminOK(w, map[string]any{"tags": tags})
	default:
		writeMethodNotAllowed(w, true, http.MethodGet, http.MethodPut)
	}
}

func (s *Server) handleAdminUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w, true, http.MethodPost)
		return
	}
	url, err := s.saveUpload(w, r, "file", r.FormValue("kind"))
	if err != nil {
		writeError(w, true, http.StatusBadRequest, err.Error())
		return
	}
	writeAdminOK(w, map[string]string{"url": url})
}

func (s *Server) handleAdminLogo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w, true, http.MethodPost)
		return
	}
	url, err := s.saveUpload(w, r, "logo", "logo")
	if err != nil {
		writeError(w, true, http.StatusBadRequest, err.Error())
		return
	}
	site, err := s.store.Site()
	if err == nil {
		site.LogoURL = url
		site, err = s.store.SaveSite(site)
	}
	if err != nil {
		writeError(w, true, http.StatusBadRequest, err.Error())
		return
	}
	writeAdminOK(w, map[string]any{"logoUrl": url, "site": site})
}

func (s *Server) handleAdminArticles(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/admin/articles" {
		writeNotFound(w, true)
		return
	}
	switch r.Method {
	case http.MethodGet:
		query := ArticleQuery{
			Keyword:      r.URL.Query().Get("keyword"),
			Tag:          r.URL.Query().Get("tag"),
			Status:       r.URL.Query().Get("status"),
			DateFrom:     r.URL.Query().Get("dateFrom"),
			DateTo:       r.URL.Query().Get("dateTo"),
			Page:         intQuery(r, "page", 1),
			PageSize:     intQuery(r, "pageSize", 10),
			IncludeDraft: true,
		}
		articles, total, err := s.store.ListArticles(query)
		if err != nil {
			writeError(w, true, http.StatusInternalServerError, err.Error())
			return
		}
		writeAdminOK(w, PageResult{List: articles, Total: total})
	case http.MethodPost:
		var article Article
		if err := readJSON(r, &article); err != nil {
			writeError(w, true, http.StatusBadRequest, err.Error())
			return
		}
		next, err := s.store.SaveArticle(article, "")
		if err != nil {
			writeError(w, true, http.StatusBadRequest, err.Error())
			return
		}
		writeAdminOK(w, next)
	default:
		writeMethodNotAllowed(w, true, http.MethodGet, http.MethodPost)
	}
}

func (s *Server) handleAdminArticle(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/admin/articles/")
	id = strings.TrimSpace(id)
	if id == "" {
		writeNotFound(w, true)
		return
	}
	if strings.HasSuffix(id, "/publish") {
		id = strings.TrimSuffix(id, "/publish")
		if r.Method != http.MethodPost {
			writeMethodNotAllowed(w, true, http.MethodPost)
			return
		}
		if err := s.store.SetArticleStatus(id, ArticlePublished); err != nil {
			writeError(w, true, http.StatusBadRequest, err.Error())
			return
		}
		article, _ := s.store.Article(id, true)
		writeAdminOK(w, article)
		return
	}
	switch r.Method {
	case http.MethodGet:
		article, err := s.store.Article(id, true)
		if err != nil {
			writeNotFound(w, true)
			return
		}
		writeAdminOK(w, article)
	case http.MethodPut:
		var article Article
		if err := readJSON(r, &article); err != nil {
			writeError(w, true, http.StatusBadRequest, err.Error())
			return
		}
		next, err := s.store.SaveArticle(article, id)
		if err != nil {
			writeError(w, true, http.StatusBadRequest, err.Error())
			return
		}
		writeAdminOK(w, next)
	case http.MethodDelete:
		if err := s.store.DeleteArticle(id); err != nil {
			writeError(w, true, http.StatusBadRequest, err.Error())
			return
		}
		writeAdminOK(w, map[string]any{"ok": true})
	default:
		writeMethodNotAllowed(w, true, http.MethodGet, http.MethodPut, http.MethodDelete)
	}
}

func (s *Server) handleAdminUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		users, total, err := s.store.ListUsers(intQuery(r, "page", 1), intQuery(r, "pageSize", 20))
		if err != nil {
			writeError(w, true, http.StatusInternalServerError, err.Error())
			return
		}
		writeAdminOK(w, PageResult{List: users, Total: total})
	case http.MethodPost:
		var request UserCreateRequest
		if err := readJSON(r, &request); err != nil {
			writeError(w, true, http.StatusBadRequest, err.Error())
			return
		}
		user, err := s.store.CreateUser(request)
		if err != nil {
			writeError(w, true, http.StatusBadRequest, err.Error())
			return
		}
		writeAdminOK(w, user)
	default:
		writeMethodNotAllowed(w, true, http.MethodGet, http.MethodPost)
	}
}

func (s *Server) handleAdminUser(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/admin/users/")
	parts := strings.Split(strings.Trim(path, "/"), "/")
	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		writeNotFound(w, true)
		return
	}
	if len(parts) == 2 && parts[1] == "reset-password" {
		if r.Method != http.MethodPost {
			writeMethodNotAllowed(w, true, http.MethodPost)
			return
		}
		var request struct {
			Password string `json:"password"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, true, http.StatusBadRequest, err.Error())
			return
		}
		if err := s.store.ResetPassword(id, request.Password); err != nil {
			writeError(w, true, http.StatusBadRequest, err.Error())
			return
		}
		writeAdminOK(w, map[string]any{"ok": true})
		return
	}
	if len(parts) == 2 && parts[1] == "status" {
		if r.Method != http.MethodPost {
			writeMethodNotAllowed(w, true, http.MethodPost)
			return
		}
		var request struct {
			Status string `json:"status"`
		}
		if err := readJSON(r, &request); err != nil {
			writeError(w, true, http.StatusBadRequest, err.Error())
			return
		}
		user, err := s.store.SetUserStatus(id, request.Status)
		if err != nil {
			writeError(w, true, http.StatusBadRequest, err.Error())
			return
		}
		writeAdminOK(w, user)
		return
	}
	switch r.Method {
	case http.MethodGet:
		user, err := s.store.UserByID(id)
		if err != nil {
			writeNotFound(w, true)
			return
		}
		writeAdminOK(w, user)
	case http.MethodPut:
		var request UserUpdateRequest
		if err := readJSON(r, &request); err != nil {
			writeError(w, true, http.StatusBadRequest, err.Error())
			return
		}
		user, err := s.store.UpdateUser(id, request)
		if err != nil {
			writeError(w, true, http.StatusBadRequest, err.Error())
			return
		}
		writeAdminOK(w, user)
	case http.MethodDelete:
		if err := s.store.DeleteUser(id); err != nil {
			writeError(w, true, http.StatusBadRequest, err.Error())
			return
		}
		writeAdminOK(w, map[string]any{"ok": true})
	default:
		writeMethodNotAllowed(w, true, http.MethodGet, http.MethodPut, http.MethodDelete)
	}
}

func (s *Server) handleLoginLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, true, http.MethodGet)
		return
	}
	logs, total, err := s.store.ListLoginLogs(intQuery(r, "page", 1), intQuery(r, "pageSize", 20))
	if err != nil {
		writeError(w, true, http.StatusInternalServerError, err.Error())
		return
	}
	writeAdminOK(w, PageResult{List: logs, Total: total})
}

func (s *Server) withAuth(next http.HandlerFunc, roles ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := s.authenticateRequest(r)
		if err != nil {
			writeError(w, true, http.StatusUnauthorized, err.Error())
			return
		}
		if len(roles) > 0 && !roleAllowed(user.Role, roles) {
			writeError(w, true, http.StatusForbidden, "权限不足")
			return
		}
		next(w, r.WithContext(context.WithValue(r.Context(), userContextKey, user)))
	}
}

func (s *Server) authenticateRequest(r *http.Request) (User, error) {
	tokenString := bearerToken(r.Header.Get("Authorization"))
	if tokenString == "" {
		return User{}, errors.New("登录已失效，请重新登录")
	}
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return User{}, errors.New("登录已失效，请重新登录")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return User{}, errors.New("登录已失效，请重新登录")
	}
	userIDFloat, ok := claims["userId"].(float64)
	if !ok {
		return User{}, errors.New("登录已失效，请重新登录")
	}
	user, err := s.store.UserByID(int64(userIDFloat))
	if err != nil || user.Status != UserActive {
		return User{}, errors.New("登录已失效，请重新登录")
	}
	return user, nil
}

func (s *Server) issueToken(user User) (string, error) {
	claims := jwt.MapClaims{
		"userId":   user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *Server) saveUpload(w http.ResponseWriter, r *http.Request, field string, kind string) (string, error) {
	r.Body = http.MaxBytesReader(w, r.Body, s.maxUploadSize)
	file, header, err := r.FormFile(field)
	if err != nil && field != "file" {
		file, header, err = r.FormFile("file")
	}
	if err != nil {
		return "", errors.New("file is required")
	}
	defer file.Close()
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" {
		if exts, err := mime.ExtensionsByType(header.Header.Get("Content-Type")); err == nil && len(exts) > 0 {
			ext = exts[0]
		}
	}
	if !allowedImageExt(ext) {
		return "", errors.New("file must be png, jpg, jpeg, webp, svg, or gif")
	}
	if err := os.MkdirAll(s.uploadDir, 0755); err != nil {
		return "", err
	}
	if kind == "" {
		kind = "image"
	}
	filename := kind + "-" + time.Now().UTC().Format("20060102150405") + ext
	destinationPath := filepath.Join(s.uploadDir, filename)
	destination, err := os.Create(destinationPath)
	if err != nil {
		return "", err
	}
	defer destination.Close()
	if _, err := io.Copy(destination, file); err != nil {
		return "", err
	}
	url := "/uploads/" + filename
	_ = s.store.SaveUpload(url, kind, filename)
	return url, nil
}

func currentUser(r *http.Request) User {
	user, _ := r.Context().Value(userContextKey).(User)
	return user
}

func roleAllowed(role string, allowed []string) bool {
	for _, value := range allowed {
		if subtle.ConstantTimeCompare([]byte(role), []byte(value)) == 1 {
			return true
		}
	}
	return false
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

func readJSON(r *http.Request, destination any) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(io.LimitReader(r.Body, 4<<20))
	decoder.DisallowUnknownFields()
	return decoder.Decode(destination)
}

func writeAdminOK(w http.ResponseWriter, payload any) {
	writeJSON(w, http.StatusOK, APIResponse{Code: 0, Message: "success", Data: payload})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, envelope bool, status int, message string) {
	if envelope {
		writeJSON(w, status, APIResponse{Code: status, Message: message})
		return
	}
	writeJSON(w, status, map[string]string{"error": message})
}

func writeMethodNotAllowed(w http.ResponseWriter, envelope bool, methods ...string) {
	w.Header().Set("Allow", strings.Join(methods, ", "))
	writeError(w, envelope, http.StatusMethodNotAllowed, "method not allowed")
}

func writeNotFound(w http.ResponseWriter, envelope bool) {
	writeError(w, envelope, http.StatusNotFound, "not found")
}

func allowedImageExt(ext string) bool {
	switch ext {
	case ".png", ".jpg", ".jpeg", ".webp", ".svg", ".gif":
		return true
	default:
		return false
	}
}

func bearerToken(header string) string {
	parts := strings.Fields(header)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}
	return ""
}

func intQuery(r *http.Request, key string, fallback int) int {
	value, err := strconv.Atoi(r.URL.Query().Get(key))
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}

func clientIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return strings.TrimSpace(strings.Split(forwarded, ",")[0])
	}
	host := r.RemoteAddr
	if index := strings.LastIndex(host, ":"); index > -1 {
		return host[:index]
	}
	return host
}
