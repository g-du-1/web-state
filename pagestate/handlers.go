package pagestate

import (
	"encoding/json"
	"net/http"
	"time"
)

type SavePagesstateRequest struct {
	Url         string `json:"url"`
	ScrollPos   int    `json:"scrollPos"`
	VisibleText string `json:"visibleText"`
}

type PagestateResponse struct {
	Id          int       `json:"id"`
	Url         string    `json:"url"`
	ScrollPos   int       `json:"scrollPos"`
	VisibleText string    `json:"visibleText"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) SavePageState(w http.ResponseWriter, r *http.Request) {
	var req SavePagesstateRequest

	_ = json.NewDecoder(r.Body).Decode(&req)

	pagestate := Pagestate{
		Url:         req.Url,
		ScrollPos:   req.ScrollPos,
		VisibleText: req.VisibleText,
	}

	createdPagestate, _ := h.repo.CreatePagestate(r.Context(), pagestate)

	response := PagestateResponse(createdPagestate)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
