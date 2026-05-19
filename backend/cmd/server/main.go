package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"liang-blog-backend/internal/cms"
)

func main() {
	addr := getenv("ADDR", ":8080")
	dataPath := getenv("DATA_PATH", "data/site.json")
	databasePath := getenv("DATABASE_PATH", "data/blog.db")
	uploadDir := getenv("UPLOAD_DIR", "uploads")
	initialPassword := os.Getenv("ADMIN_INITIAL_PASSWORD")
	if initialPassword == "" {
		initialPassword = "dev-admin-token"
		log.Println("ADMIN_INITIAL_PASSWORD is not set; using development password for admin: dev-admin-token")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-jwt-secret-change-me"
		log.Println("JWT_SECRET is not set; using development JWT secret")
	}

	store, err := cms.OpenStore(databasePath, dataPath, initialPassword)
	if err != nil {
		log.Fatalf("open database store: %v", err)
	}
	defer store.Close()

	handler := cms.NewServer(store, cms.Config{
		JWTSecret:     jwtSecret,
		UploadDir:     uploadDir,
		CORSOrigins:   splitCSV(getenv("CORS_ORIGIN", "http://127.0.0.1:5173,http://localhost:5173")),
		MaxUploadSize: 5 << 20,
	})

	log.Printf("Liang CMS API listening on http://127.0.0.1%s", addr)
	log.Printf("SQLite database: %s", databasePath)
	log.Printf("Seed content file: %s", dataPath)
	log.Printf("Uploads dir: %s", uploadDir)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}

func getenv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}
