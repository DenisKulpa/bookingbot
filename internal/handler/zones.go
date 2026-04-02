package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/DenisKulpa/bookingbot/internal/repository"
	"github.com/go-chi/chi/v5"
)

type ZoneHandler struct {
	repo *repository.ZoneRepository
}

func NewZoneHandler(repo *repository.ZoneRepository) *ZoneHandler {
	return &ZoneHandler{repo: repo}
}

// GET
func (h *ZoneHandler) GetDistricts(w http.ResponseWriter, r *http.Request) {
	zones, err := h.repo.GetTopLevel(r.Context())
	if err != nil {
		log.Printf("GetDistricts error: %v", err)
		return
	}
	writeJSON(w, http.StatusOK, zones)
}

// GET /api/districts/{id}
func (h *ZoneHandler) GetDistrictDetail(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid district id")
		return
	}

	detail, err := h.repo.GetDistrictDetail(r.Context(), id)
	if err != nil {
		log.Printf("GetDistrictDetail error (id=%d): %v", id, err)
		writeError(w, http.StatusInternalServerError, "failed to fetch district detail")
		return
	}

	writeJSON(w, http.StatusOK, detail)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
