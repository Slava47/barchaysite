package main

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// newTestServer creates a Server backed by temporary directories.
func newTestServer(t *testing.T) *Server {
	t.Helper()
	dataDir := t.TempDir()
	uploadsDir := t.TempDir()

	srv, err := NewServer(dataDir, uploadsDir, "http://localhost:8080", "test-api-key", []string{"https://barchay.rudenko.tech", "*"})
	if err != nil {
		t.Fatalf("NewServer: %v", err)
	}
	return srv
}

// ── GET /health ─────────────────────────────────────────────────────────────

func TestHealthEndpoint(t *testing.T) {
	srv := newTestServer(t)
	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/health")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body["status"] != "ok" {
		t.Fatalf("expected status ok, got %q", body["status"])
	}
}

// ── GET /menu ───────────────────────────────────────────────────────────────

func TestMenuEndpoint(t *testing.T) {
	srv := newTestServer(t)
	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/menu")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var menu MenuData
	if err := json.NewDecoder(resp.Body).Decode(&menu); err != nil {
		t.Fatal(err)
	}

	if len(menu.Categories) != 3 {
		t.Fatalf("expected 3 default categories, got %d", len(menu.Categories))
	}
}

// ── GET /quiz ───────────────────────────────────────────────────────────────

func TestQuizEndpoint(t *testing.T) {
	srv := newTestServer(t)
	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/quiz")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var quiz QuizData
	if err := json.NewDecoder(resp.Body).Decode(&quiz); err != nil {
		t.Fatal(err)
	}

	if len(quiz.Questions) != 5 {
		t.Fatalf("expected 5 quiz questions, got %d", len(quiz.Questions))
	}
}

// ── CORS ────────────────────────────────────────────────────────────────────

func TestCORSHeaders(t *testing.T) {
	srv := newTestServer(t)
	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", ts.URL+"/menu", nil)
	req.Header.Set("Origin", "https://barchay.rudenko.tech")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	acao := resp.Header.Get("Access-Control-Allow-Origin")
	if acao != "https://barchay.rudenko.tech" {
		t.Fatalf("expected CORS origin header, got %q", acao)
	}
}

func TestCORSPreflight(t *testing.T) {
	srv := newTestServer(t)
	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("OPTIONS", ts.URL+"/menu", nil)
	req.Header.Set("Origin", "https://barchay.rudenko.tech")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.StatusCode)
	}
}

func TestCORSDisallowedOrigin(t *testing.T) {
	srv, err := NewServer(t.TempDir(), t.TempDir(), "http://localhost:8080", "test-key", []string{"https://barchay.rudenko.tech"})
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", ts.URL+"/menu", nil)
	req.Header.Set("Origin", "https://evil.example.com")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	acao := resp.Header.Get("Access-Control-Allow-Origin")
	if acao != "" {
		t.Fatalf("expected no CORS header for disallowed origin, got %q", acao)
	}
}

// ── POST /admin/items-with-images ───────────────────────────────────────────

func TestAdminAddItemNoAuth(t *testing.T) {
	srv := newTestServer(t)
	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("item", `{"id":"t1","category":"cold","name":"Test","price":"100"}`)
	_ = writer.Close()

	resp, err := http.Post(ts.URL+"/admin/items-with-images", writer.FormDataContentType(), body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestAdminAddItemWrongKey(t *testing.T) {
	srv := newTestServer(t)
	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("item", `{"id":"t1","category":"cold","name":"Test","price":"100"}`)
	_ = writer.Close()

	req, _ := http.NewRequest("POST", ts.URL+"/admin/items-with-images", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-API-Key", "wrong-key")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestAdminAddItemSuccess(t *testing.T) {
	srv := newTestServer(t)
	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("item", `{"id":"t1","category":"cold","name":"Test Tea","price":"250","tags":["сладкий","холодный"]}`)
	_ = writer.Close()

	req, _ := http.NewRequest("POST", ts.URL+"/admin/items-with-images", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-API-Key", "test-api-key")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 201, got %d: %s", resp.StatusCode, string(b))
	}

	var item Item
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		t.Fatal(err)
	}
	if item.ID != "t1" || item.Name != "Test Tea" {
		t.Fatalf("unexpected item: %+v", item)
	}

	// Verify the item appears in GET /menu
	menuResp, _ := http.Get(ts.URL + "/menu")
	defer menuResp.Body.Close()
	var menu MenuData
	_ = json.NewDecoder(menuResp.Body).Decode(&menu)
	if len(menu.Items) != 1 {
		t.Fatalf("expected 1 item in menu, got %d", len(menu.Items))
	}
}

func TestAdminAddItemDuplicate(t *testing.T) {
	srv := newTestServer(t)
	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	addItem := func() *http.Response {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("item", `{"id":"dup1","category":"cold","name":"Dup","price":"100"}`)
		_ = writer.Close()

		req, _ := http.NewRequest("POST", ts.URL+"/admin/items-with-images", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("X-API-Key", "test-api-key")
		resp, _ := http.DefaultClient.Do(req)
		return resp
	}

	resp1 := addItem()
	resp1.Body.Close()
	if resp1.StatusCode != http.StatusCreated {
		t.Fatalf("first add: expected 201, got %d", resp1.StatusCode)
	}

	resp2 := addItem()
	resp2.Body.Close()
	if resp2.StatusCode != http.StatusConflict {
		t.Fatalf("duplicate add: expected 409, got %d", resp2.StatusCode)
	}
}

func TestAdminAddItemWithImage(t *testing.T) {
	srv := newTestServer(t)
	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("item", `{"id":"img1","category":"hot","name":"Image Test","price":"300"}`)

	// Create a fake JPEG file part with explicit Content-Type header
	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition", `form-data; name="images"; filename="test.jpg"`)
	partHeader.Set("Content-Type", "image/jpeg")
	part, _ := writer.CreatePart(partHeader)
	// JPEG magic bytes + padding
	_, _ = part.Write([]byte{0xFF, 0xD8, 0xFF, 0xE0})
	_, _ = part.Write(make([]byte, 100))

	_ = writer.Close()

	req, _ := http.NewRequest("POST", ts.URL+"/admin/items-with-images", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-API-Key", "test-api-key")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 201, got %d: %s", resp.StatusCode, string(b))
	}

	var item Item
	_ = json.NewDecoder(resp.Body).Decode(&item)

	if item.Image == "" {
		t.Fatal("expected image URL to be set")
	}
	if len(item.Images) != 1 {
		t.Fatalf("expected 1 image, got %d", len(item.Images))
	}
	if !strings.HasPrefix(item.Image, "http://localhost:8080/uploads/") {
		t.Fatalf("unexpected image URL: %s", item.Image)
	}

	// Verify file was saved to disk
	files, _ := os.ReadDir(srv.uploadsDir)
	if len(files) != 1 {
		t.Fatalf("expected 1 uploaded file, got %d", len(files))
	}
}

func TestAdminAddItemMissingFields(t *testing.T) {
	srv := newTestServer(t)
	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("item", `{"id":"","category":"","name":"","price":""}`)
	_ = writer.Close()

	req, _ := http.NewRequest("POST", ts.URL+"/admin/items-with-images", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-API-Key", "test-api-key")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

// ── menu.json persistence ───────────────────────────────────────────────────

func TestMenuPersistence(t *testing.T) {
	dataDir := t.TempDir()
	uploadsDir := t.TempDir()

	srv1, err := NewServer(dataDir, uploadsDir, "http://localhost:8080", "key", []string{"*"})
	if err != nil {
		t.Fatal(err)
	}

	// Add an item directly
	srv1.mu.Lock()
	srv1.menu.Items = append(srv1.menu.Items, Item{
		ID:       "persist1",
		Category: "cold",
		Name:     "Persistent Tea",
		Price:    "100",
	})
	_ = srv1.saveMenuLocked()
	srv1.mu.Unlock()

	// Create a new server with the same data dir — it should load the saved data
	srv2, err := NewServer(dataDir, uploadsDir, "http://localhost:8080", "key", []string{"*"})
	if err != nil {
		t.Fatal(err)
	}

	srv2.mu.RLock()
	defer srv2.mu.RUnlock()
	if len(srv2.menu.Items) != 1 || srv2.menu.Items[0].ID != "persist1" {
		t.Fatalf("expected persisted item, got %+v", srv2.menu.Items)
	}

	// Verify the file exists on disk
	data, err := os.ReadFile(filepath.Join(dataDir, "menu.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "persist1") {
		t.Fatal("expected persist1 in menu.json")
	}
}
