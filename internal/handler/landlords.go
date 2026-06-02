package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/DenisKulpa/bookingbot/internal/repository"
	"github.com/go-chi/chi/v5"
)

type LandlordHandler struct {
	repo *repository.UserRepository
}

func NewLandlordHandler(repo *repository.UserRepository) *LandlordHandler {
	return &LandlordHandler{repo: repo}
}

// GET /api/landlords
func (h *LandlordHandler) List(w http.ResponseWriter, r *http.Request) {
	landlords, err := h.repo.ListLandlords(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, landlords)
}

// POST /api/landlords
// Body: { "telegram_id", "username", "first_name", "last_name", "phone", "company_name", "description" }
func (h *LandlordHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		TelegramID  int64  `json:"telegram_id"`
		Username    string `json:"username"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Phone       string `json:"phone"`
		CompanyName string `json:"company_name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if body.TelegramID == 0 {
		writeError(w, http.StatusBadRequest, "telegram_id is required")
		return
	}

	landlord, err := h.repo.CreateLandlord(r.Context(),
		body.TelegramID, body.Username, body.FirstName, body.LastName,
		body.Phone, body.CompanyName, body.Description,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, landlord)
}

// GET /api/landlords/{id}
func (h *LandlordHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	landlord, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if landlord == nil || landlord.Role != "landlord" {
		writeError(w, http.StatusNotFound, "landlord not found")
		return
	}
	writeJSON(w, http.StatusOK, landlord)
}

// PUT /api/landlords/{id}
// Body: { "phone", "company_name", "description" }
func (h *LandlordHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var body struct {
		Phone       string `json:"phone"`
		CompanyName string `json:"company_name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	landlord, err := h.repo.UpdateLandlord(r.Context(), id, body.Phone, body.CompanyName, body.Description)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if landlord == nil {
		writeError(w, http.StatusNotFound, "landlord not found")
		return
	}
	writeJSON(w, http.StatusOK, landlord)
}

// DELETE /api/landlords/{id}
func (h *LandlordHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.repo.DeleteLandlord(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GET /api/landlords/{id}/apartments
func (h *LandlordHandler) Apartments(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	apts, err := h.repo.ListApartmentsByLandlord(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, apts)
}
