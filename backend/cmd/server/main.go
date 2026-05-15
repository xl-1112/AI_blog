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
	uploadDir := getenv("UPLOAD_DIR", "uploads")
	adminToken := os.Getenv("ADMIN_TOKEN")
	if adminToken == "" {
		adminToken = "dev-admin-token"
		log.Println("ADMIN_TOKEN is not set; using development token: dev-admin-token")
	}

	store := cms.NewStore(dataPath)
	if err := store.Load(); err != nil {
		log.Fatalf("load content store: %v", err)
	}

	handler := cms.NewServer(store, cms.Config{
		AdminToken:    adminToken,
		UploadDir:     uploadDir,
		CORSOrigins:   splitCSV(getenv("CORS_ORIGIN", "http://127.0.0.1:5173,http://localhost:5173")),
		MaxUploadSize: 5 << 20,
	})

	log.Printf("Liang CMS API listening on http://127.0.0.1%s", addr)
	log.Printf("Content file: %s", dataPath)
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
