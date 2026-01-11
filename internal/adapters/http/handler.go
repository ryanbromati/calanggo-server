package http

import (
	"encoding/json"
	"net/http"

	"calanggo-server/internal/core/ports"
)

type LinkHandler struct {
	service ports.LinkService
}

func NewLinkHandler(service ports.LinkService) *LinkHandler {
	return &LinkHandler{service: service}
}

type createLinkRequest struct {
	OriginalURL string `json:"original_url"`
}

type createLinkResponse struct {
	ShortCode   string `json:"short_code"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// CreateShortLink - POST /api/v1/shorten
func (h *LinkHandler) CreateShortLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req createLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	link, err := h.service.CreateShortLink(r.Context(), req.OriginalURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := createLinkResponse{
		ShortCode:   link.Shortened,
		ShortURL:    "http://" + r.Host + "/" + link.Shortened,
		OriginalURL: link.Original,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// Redirect - GET /{code}
func (h *LinkHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		http.Error(w, "Code is required", http.StatusBadRequest)
		return
	}

	originalURL, err := h.service.GetOriginalURL(r.Context(), code)
	if err != nil {
		http.Error(w, "Link not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}
