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

type GetPagestateResponse struct {
	Id        int    `json:"id"`
	Url       string `json:"url"`
	ScrollPos int    `json:"scrollPos"`
}

type GetAllPagestatesResponse struct {
	Pagestates []GetPagestateResponse `json:"pagestates"`
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

	_ = json.NewDecoder(r.Body).Decode(&req)

	pagestate := Pagestate{
		Url:       req.Url,
		ScrollPos: req.ScrollPos,
	}

	createdPagestate, _ := h.repo.CreatePagestate(r.Context(), pagestate)

	response := CreatePagestateResponse(createdPagestate)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetPagestate(w http.ResponseWriter, r *http.Request) {
	pagestates, err := h.repo.GetAllPagestates(r.Context())

	if err != nil {
		http.Error(w, "failed to retrieve page states", http.StatusInternalServerError)
		return
	}

	response := GetAllPagestatesResponse{
		Pagestates: make([]GetPagestateResponse, len(pagestates)),
	}

	for i, pagestate := range pagestates {
		response.Pagestates[i] = GetPagestateResponse(pagestate)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
