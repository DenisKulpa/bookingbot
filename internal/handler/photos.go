package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/DenisKulpa/bookingbot/internal/repository"
	"github.com/go-chi/chi/v5"
)

type PhotoHandler struct {
	repo       *repository.PhotoRepository
	uploadsDir string // абсолютный путь к папке uploads/apartments
}

func NewPhotoHandler(repo *repository.PhotoRepository) *PhotoHandler {
	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(filename), "..", "..")
	uploadsDir := filepath.Join(root, "uploads", "apartments")
	return &PhotoHandler{repo: repo, uploadsDir: uploadsDir}
}

// POST /api/apartments/{id}/photos
// Content-Type: multipart/form-data, поле "photo"
func (h *PhotoHandler) Upload(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	apartmentID, err := strconv.Atoi(idStr)
	if err != nil || apartmentID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid apartment id")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB
		writeError(w, http.StatusBadRequest, "file too large or bad form")
		return
	}

	file, header, err := r.FormFile("photo")
	if err != nil {
		writeError(w, http.StatusBadRequest, "field 'photo' is required")
		return
	}
	defer file.Close()

	// Проверка расширения
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowed[ext] {
		writeError(w, http.StatusBadRequest, "only jpg, png, webp allowed")
		return
	}

	// Создать папку uploads/apartments/{apartment_id}/
	dir := filepath.Join(h.uploadsDir, strconv.Itoa(apartmentID))
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("PhotoHandler.Upload: mkdir %s: %v", dir, err)
		writeError(w, http.StatusInternalServerError, "failed to create upload dir")
		return
	}

	// Уникальное имя файла: timestamp + оригинальное имя
	filename := fmt.Sprintf("%d_%s", time.Now().UnixMilli(), header.Filename)
	absPath := filepath.Join(dir, filename)
	relPath := filepath.ToSlash(filepath.Join("uploads", "apartments", strconv.Itoa(apartmentID), filename))
	publicURL := "/" + relPath

	dst, err := os.Create(absPath)
	if err != nil {
		log.Printf("PhotoHandler.Upload: create file: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to save file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		log.Printf("PhotoHandler.Upload: copy: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to write file")
		return
	}

	// Определить sort_order (следующий по порядку)
	photos, _ := h.repo.GetByApartment(r.Context(), apartmentID)
	sortOrder := len(photos)

	photo, err := h.repo.Add(r.Context(), apartmentID, relPath, publicURL, sortOrder)
	if err != nil {
		log.Printf("PhotoHandler.Upload: db: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to save photo record")
		return
	}

	writeJSON(w, http.StatusCreated, photo)
}

// GET /api/apartments/{id}/photos
func (h *PhotoHandler) List(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	apartmentID, err := strconv.Atoi(idStr)
	if err != nil || apartmentID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid apartment id")
		return
	}

	photos, err := h.repo.GetByApartment(r.Context(), apartmentID)
	if err != nil {
		log.Printf("PhotoHandler.List: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to fetch photos")
		return
	}

	writeJSON(w, http.StatusOK, photos)
}

// DELETE /api/photos/{id}
func (h *PhotoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	photoID, err := strconv.Atoi(idStr)
	if err != nil || photoID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid photo id")
		return
	}

	if err := h.repo.Delete(r.Context(), photoID); err != nil {
		log.Printf("PhotoHandler.Delete: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to delete photo")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
