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

// CreateLinkRequest represents the request body for creating a short link
type CreateLinkRequest struct {
	OriginalURL string `json:"original_url" example:"https://www.google.com"`
}

// CreateLinkResponse represents the response for a created short link
type CreateLinkResponse struct {
	ShortCode   string `json:"short_code" example:"7rB8u"`
	ShortURL    string `json:"short_url" example:"http://localhost:8080/7rB8u"`
	OriginalURL string `json:"original_url" example:"https://www.google.com"`
}

// CreateShortLink godoc
// @Summary      Cria um link curto
// @Description  Recebe uma URL longa e retorna o link encurtado com o código gerado.
// @Tags         links
// @Accept       json
// @Produce      json
// @Param        request body CreateLinkRequest true "URL para encurtar"
// @Success      201  {object}  CreateLinkResponse
// @Failure      400  {string}  string "Invalid request body"
// @Failure      500  {string}  string "Internal Server Error"
// @Router       /api/v1/shorten [post]
func (h *LinkHandler) CreateShortLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	link, err := h.service.CreateShortLink(r.Context(), req.OriginalURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := CreateLinkResponse{
		ShortCode:   link.Shortened,
		ShortURL:    "http://" + r.Host + "/" + link.Shortened,
		OriginalURL: link.Original,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// Redirect godoc
// @Summary      Redireciona para a URL original
// @Description  Busca o código curto e redireciona o usuário (HTTP 302).
// @Tags         links
// @Param        code path string true "Código curto do link"
// @Success      302  {string}  string "Location header com a URL original"
// @Failure      404  {string}  string "Link not found"
// @Router       /{code} [get]
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
