package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/DenisKulpa/bookingbot/internal/repository"
	"github.com/go-chi/chi/v5"
)

type ApartmentHandler struct {
	repo *repository.ApartmentRepository
}

func NewApartmentHandler(repo *repository.ApartmentRepository) *ApartmentHandler {
	return &ApartmentHandler{repo: repo}
}

// GET /api/districts/{id}/apartments?available=true
func (h *ApartmentHandler) GetApartments(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	zoneID, err := strconv.Atoi(idStr)
	if err != nil || zoneID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid district id")
		return
	}

	onlyAvailable := r.URL.Query().Get("available") == "true"

	data, err := h.repo.GetByZone(r.Context(), zoneID, onlyAvailable)
	if err != nil {
		log.Printf("GetApartments error (zone=%d): %v", zoneID, err)
		writeError(w, http.StatusInternalServerError, "failed to fetch apartments")
		return
	}

	writeJSON(w, http.StatusOK, data)
}
