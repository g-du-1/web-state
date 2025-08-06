package pagestate

import (
	"encoding/json"
	"net/http"
)

type CreatePagestateRequest struct {
	Url       string `json:"url"`
	ScrollPos int    `json:"scrollPos"`
}

type CreatePagestateResponse struct {
	Id        int    `json:"id"`
	Url       string `json:"url"`
	ScrollPos int    `json:"scrollPos"`
}

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) CreatePagestate(w http.ResponseWriter, r *http.Request) {
	var req CreatePagestateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	pagestate := Pagestate{
		Url:       req.Url,
		ScrollPos: req.ScrollPos,
	}

	createdPagestate, err := h.repo.CreatePagestate(r.Context(), pagestate)

	if err != nil {
		http.Error(w, "Failed to create pagestate", http.StatusInternalServerError)
		return
	}

	response := CreatePagestateResponse{
		Id:        createdPagestate.Id,
		Url:       createdPagestate.Url,
		ScrollPos: createdPagestate.ScrollPos,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
