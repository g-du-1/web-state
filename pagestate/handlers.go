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
	UpdatedAt   time.Time `json:"updatedAt"`
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

	createdPagestate, _ := h.repo.SavePagestate(r.Context(), pagestate)

	response := PagestateResponse(createdPagestate)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetPageState(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")

	pagestate, err := h.repo.GetPagestate(r.Context(), url)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response := PagestateResponse(pagestate)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetAllPageStates(w http.ResponseWriter, r *http.Request) {
	pagestates, _ := h.repo.GetAllPagestates(r.Context())

	response := make([]PagestateResponse, len(pagestates))

	for i, pagestate := range pagestates {
		response[i] = PagestateResponse(pagestate)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) DeleteAllPageStates(w http.ResponseWriter, r *http.Request) {
	h.repo.DeleteAllPageStates(r.Context())

	w.WriteHeader(http.StatusOK)
}
