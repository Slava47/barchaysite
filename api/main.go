package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"
)

// ── Data types ──────────────────────────────────────────────────────────────

// Category represents a menu category.
type Category struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	NameZh string `json:"nameZh,omitempty"`
}

// Item represents a menu item.
type Item struct {
	ID              string   `json:"id"`
	Category        string   `json:"category"`
	Name            string   `json:"name"`
	NameZh          string   `json:"nameZh,omitempty"`
	Price           string   `json:"price"`
	Description     string   `json:"description,omitempty"`
	FullDescription string   `json:"fullDescription,omitempty"`
	Image           string   `json:"image,omitempty"`
	Images          []string `json:"images,omitempty"`
	Tags            []string `json:"tags,omitempty"`
}

// MenuData represents the full menu structure.
type MenuData struct {
	Categories []Category `json:"categories"`
	Items      []Item     `json:"items"`
}

// QuizOption represents a quiz answer option.
type QuizOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// QuizQuestion represents a quiz question.
type QuizQuestion struct {
	Text    string       `json:"text"`
	Options []QuizOption `json:"options"`
}

// QuizData represents the quiz structure.
type QuizData struct {
	Questions []QuizQuestion `json:"questions"`
}

// ── Server ──────────────────────────────────────────────────────────────────

// Server holds the state and configuration for the API server.
type Server struct {
	mu         sync.RWMutex
	menu       *MenuData
	quiz       *QuizData
	dataDir    string
	uploadsDir string
	publicURL  string
	apiKey     string
	allowedOrigins []string
}

// NewServer creates a new Server with the given configuration.
func NewServer(dataDir, uploadsDir, publicURL, apiKey string, allowedOrigins []string) (*Server, error) {
	s := &Server{
		dataDir:    dataDir,
		uploadsDir: uploadsDir,
		publicURL:  publicURL,
		apiKey:     apiKey,
		allowedOrigins: allowedOrigins,
	}

	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}
	if err := os.MkdirAll(uploadsDir, 0o755); err != nil {
		return nil, fmt.Errorf("create uploads dir: %w", err)
	}

	if err := s.loadMenu(); err != nil {
		return nil, fmt.Errorf("load menu: %w", err)
	}

	s.quiz = defaultQuizData()

	return s, nil
}

// loadMenu reads menu.json from dataDir. If missing, creates a default file.
func (s *Server) loadMenu() error {
	path := filepath.Join(s.dataDir, "menu.json")

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		s.menu = defaultMenuData()
		return s.saveMenuLocked()
	}
	if err != nil {
		return err
	}

	var menu MenuData
	if err := json.Unmarshal(data, &menu); err != nil {
		return fmt.Errorf("parse menu.json: %w", err)
	}
	s.menu = &menu
	return nil
}

// saveMenuLocked atomically writes menu.json. Caller must hold s.mu.
func (s *Server) saveMenuLocked() error {
	data, err := json.MarshalIndent(s.menu, "", "  ")
	if err != nil {
		return err
	}

	path := filepath.Join(s.dataDir, "menu.json")
	tmpPath := path + ".tmp"

	if err := os.WriteFile(tmpPath, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}

// ── HTTP Handlers ───────────────────────────────────────────────────────────

// Handler returns the top-level HTTP handler with all routes and middleware.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /menu", s.handleGetMenu)
	mux.HandleFunc("GET /quiz", s.handleGetQuiz)
	mux.HandleFunc("GET /health", s.handleGetHealth)
	mux.HandleFunc("POST /admin/items-with-images", s.handleAdminAddItem)

	// Serve uploaded files
	uploadsFS := http.StripPrefix("/uploads/", http.FileServer(http.Dir(s.uploadsDir)))
	mux.Handle("GET /uploads/", uploadsFS)

	return s.corsMiddleware(mux)
}

func (s *Server) handleGetHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func (s *Server) handleGetMenu(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(s.menu); err != nil {
		log.Printf("encode menu: %v", err)
	}
}

func (s *Server) handleGetQuiz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(s.quiz); err != nil {
		log.Printf("encode quiz: %v", err)
	}
}

func (s *Server) handleAdminAddItem(w http.ResponseWriter, r *http.Request) {
	// 1. Check API key
	if s.apiKey == "" {
		http.Error(w, `{"error":"server misconfigured: ADMIN_API_KEY not set"}`, http.StatusInternalServerError)
		return
	}
	key := r.Header.Get("X-API-Key")
	if key == "" || key != s.apiKey {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// 2. Parse multipart form (max 50 MB)
	const maxBody = 50 << 20
	r.Body = http.MaxBytesReader(w, r.Body, maxBody)
	if err := r.ParseMultipartForm(maxBody); err != nil {
		http.Error(w, `{"error":"request too large or invalid multipart form"}`, http.StatusBadRequest)
		return
	}
	defer func() {
		if r.MultipartForm != nil {
			_ = r.MultipartForm.RemoveAll()
		}
	}()

	// 3. Read "item" part
	itemJSON := r.FormValue("item")
	if itemJSON == "" {
		http.Error(w, `{"error":"missing 'item' part"}`, http.StatusBadRequest)
		return
	}

	var item Item
	if err := json.Unmarshal([]byte(itemJSON), &item); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"invalid item JSON: %s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	// 4. Validate required fields
	if item.ID == "" || item.Category == "" || item.Name == "" || item.Price == "" {
		http.Error(w, `{"error":"missing required fields: id, category, name, price"}`, http.StatusBadRequest)
		return
	}

	// 5. Check for duplicate ID
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, existing := range s.menu.Items {
		if existing.ID == item.ID {
			http.Error(w, `{"error":"item with this id already exists"}`, http.StatusConflict)
			return
		}
	}

	// 6. Process uploaded images
	var imageURLs []string
	files := r.MultipartForm.File["images"]
	for _, fh := range files {
		if fh.Size > 10<<20 {
			http.Error(w, `{"error":"individual file exceeds 10 MB limit"}`, http.StatusBadRequest)
			return
		}

		url, err := s.saveUploadedFile(fh)
		if err != nil {
			log.Printf("save upload: %v", err)
			http.Error(w, fmt.Sprintf(`{"error":"failed to save image: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		imageURLs = append(imageURLs, url)
	}

	if len(imageURLs) > 0 {
		item.Image = imageURLs[0]
		item.Images = imageURLs
	}

	// 7. Add to menu and save
	s.menu.Items = append(s.menu.Items, item)
	if err := s.saveMenuLocked(); err != nil {
		// Rollback: remove the item we just added
		s.menu.Items = s.menu.Items[:len(s.menu.Items)-1]
		log.Printf("save menu: %v", err)
		http.Error(w, `{"error":"failed to save menu"}`, http.StatusInternalServerError)
		return
	}

	// 8. Return created item
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(item)
}

// saveUploadedFile stores an image file and returns its public URL.
func (s *Server) saveUploadedFile(fh *multipart.FileHeader) (string, error) {
	// Validate MIME type
	ct := fh.Header.Get("Content-Type")
	if ct == "" {
		ct = "application/octet-stream"
	}
	mediaType, _, _ := mime.ParseMediaType(ct)
	if !strings.HasPrefix(mediaType, "image/") {
		return "", fmt.Errorf("unsupported file type: %s", mediaType)
	}

	// Determine extension from MIME type
	ext := ".bin"
	switch mediaType {
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	case "image/webp":
		ext = ".webp"
	case "image/gif":
		ext = ".gif"
	default:
		exts, _ := mime.ExtensionsByType(mediaType)
		if len(exts) > 0 {
			ext = exts[0]
		}
	}

	filename := uuid.New().String() + ext
	destPath := filepath.Join(s.uploadsDir, filename)

	src, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		_ = os.Remove(destPath)
		return "", err
	}

	publicURL := s.publicURL + "/uploads/" + filename
	return publicURL, nil
}

// ── CORS Middleware ─────────────────────────────────────────────────────────

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if s.isAllowedOrigin(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) isAllowedOrigin(origin string) bool {
	if origin == "" {
		return false
	}
	for _, allowed := range s.allowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
	}
	return false
}

// ── Default Data ────────────────────────────────────────────────────────────

func defaultMenuData() *MenuData {
	return &MenuData{
		Categories: []Category{
			{ID: "cold", Name: "Холодные коктейли", NameZh: "冷饮"},
			{ID: "hot", Name: "Горячие коктейли", NameZh: "热饮"},
			{ID: "alco", Name: "Алкогольные коктейли", NameZh: "酒饮"},
		},
		Items: []Item{},
	}
}

func defaultQuizData() *QuizData {
	return &QuizData{
		Questions: []QuizQuestion{
			{
				Text: "Какой вкус Вам ближе?",
				Options: []QuizOption{
					{Label: "Сладкий, мягкий", Value: "сладкий"},
					{Label: "Кислый, освежающий", Value: "кислый"},
					{Label: "Крепкий, насыщенный", Value: "крепкий"},
					{Label: "Пряный, согревающий", Value: "пряный"},
				},
			},
			{
				Text: "Предпочитаете горячее или холодное?",
				Options: []QuizOption{
					{Label: "Горячее", Value: "тёплый"},
					{Label: "Холодное", Value: "холодный"},
					{Label: "Без разницы", Value: "любой"},
				},
			},
			{
				Text: "Хотите ли Вы алкогольный напиток?",
				Options: []QuizOption{
					{Label: "Да", Value: "алкогольный"},
					{Label: "Нет", Value: "безалкогольный"},
					{Label: "Не важно", Value: "любой_алко"},
				},
			},
			{
				Text: "Какое настроение Вы хотите создать?",
				Options: []QuizOption{
					{Label: "Расслабиться и отдохнуть", Value: "мягкий"},
					{Label: "Взбодриться", Value: "бодрящий"},
					{Label: "Удивить себя чем-то необычным", Value: "необычный"},
					{Label: "Насладиться классикой", Value: "классический"},
				},
			},
			{
				Text: "Какие ноты Вам нравятся больше?",
				Options: []QuizOption{
					{Label: "Цветочные, травяные", Value: "цветочный"},
					{Label: "Фруктовые, ягодные", Value: "фруктовый"},
					{Label: "Дымные, землистые", Value: "дымный"},
					{Label: "Десертные, молочные", Value: "десертный"},
				},
			},
		},
	}
}

// ── Env helpers ─────────────────────────────────────────────────────────────

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// ── Main ────────────────────────────────────────────────────────────────────

func main() {
	port := envOrDefault("PORT", "8080")
	dataDir := envOrDefault("DATA_DIR", "data")
	uploadsDir := envOrDefault("UPLOADS_DIR", "uploads")
	publicURL := envOrDefault("PUBLIC_URL", "http://localhost:"+port)
	apiKey := os.Getenv("ADMIN_API_KEY")
	allowedOriginsStr := envOrDefault("ALLOWED_ORIGINS", "https://barchay.rudenko.tech")

	allowedOrigins := strings.Split(allowedOriginsStr, ",")
	for i := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
	}

	if apiKey == "" {
		log.Println("WARNING: ADMIN_API_KEY is not set. Admin endpoints will reject all requests.")
	}

	srv, err := NewServer(dataDir, uploadsDir, publicURL, apiKey, allowedOrigins)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	addr := ":" + port
	log.Printf("Starting Barchay API on %s", addr)
	log.Printf("  Data dir:    %s", dataDir)
	log.Printf("  Uploads dir: %s", uploadsDir)
	log.Printf("  Public URL:  %s", publicURL)

	if err := http.ListenAndServe(addr, srv.Handler()); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
