package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	httpAdapter "calanggo-server/internal/adapters/http"
	"calanggo-server/internal/adapters/repository"
	"calanggo-server/internal/core/services"
)

func TestShortenFlow(t *testing.T) {
	// Setup (Igual ao main.go)
	repo := repository.NewMemoryRepository()
	service := services.NewLinkService(repo)
	handler := httpAdapter.NewLinkHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/shorten", handler.CreateShortLink)
	mux.HandleFunc("GET /{code}", handler.Redirect)

	// Servidor de Teste (roda na porta disponível)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	client := ts.Client()
	originalURL := "https://bytebytego.com"

	// 1. Encurtar
	reqBody, _ := json.Marshal(map[string]string{
		"original_url": originalURL,
	})

	resp, err := client.Post(ts.URL+"/api/v1/shorten", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to shorten: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status 201, got %d", resp.StatusCode)
	}

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	shortCode := result["short_code"]
	if shortCode == "" {
		t.Fatal("Short code is empty")
	}
	t.Logf("Short Code: %s", shortCode)

	// 2. Redirecionar
	// Por padrão o client Go segue redirects. Vamos desabilitar para verificar o 302.
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	redirectURL := ts.URL + "/" + shortCode
	resp2, err := client.Get(redirectURL)
	if err != nil {
		t.Fatalf("Failed to get redirect: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusFound { // 302
		t.Fatalf("Expected status 302, got %d", resp2.StatusCode)
	}

	location := resp2.Header.Get("Location")
	if location != originalURL {
		t.Errorf("Expected redirect to %s, got %s", originalURL, location)
	}
}
