package main

import (
	"encoding/json"
	"io"
	"minio-go-s3/internal/service"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func (app *application) CreatePasteHandler(w http.ResponseWriter, r *http.Request) {
	const maxBody = 10 << 20
	const defaultTTL = 24 * time.Hour
	const maxTTL = 30 * 24 * time.Hour

	r.Body = http.MaxBytesReader(w, r.Body, maxBody)
	defer r.Body.Close()

	content, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
		return
	}

	filename := r.URL.Query().Get("filename")
	idempotencyKey := r.Header.Get("Idempotency-Key")

	// TTL parsing
	ttl := defaultTTL
	if ttlStr := r.URL.Query().Get("ttl"); ttlStr != "" {
		parsed, err := time.ParseDuration(ttlStr)
		if err != nil || parsed <= 0 {
			http.Error(w, "invalid ttl", http.StatusBadRequest)
			return
		}
		if parsed > maxTTL {
			http.Error(w, "ttl too large", http.StatusBadRequest)
			return
		}
		ttl = parsed
	}

	meta, err := app.pasteSvc.CreatePaste(r.Context(), content, filename, idempotencyKey, ttl)
	if err != nil {
		switch err {
		case service.ErrEmptyPaste:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	resp := map[string]string{
		"id":  meta.ID,
		"url": app.config.baseURL + "/v1/pastes/" + meta.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

}

func (app *application) GetPasteHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, "missing paste ID", http.StatusBadRequest)
		return
	}

	meta, data, err := app.pasteSvc.GetPaste(r.Context(), id)
	if err != nil {
		if err == service.ErrExpiredPaste {
			http.Error(w, "paste expired", http.StatusGone)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if meta == nil {
		http.Error(w, "paste not found", http.StatusNotFound)
		return
	}

	// Optional: set content type if you stored it in metadata
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", "inline; filename="+meta.Filename)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
