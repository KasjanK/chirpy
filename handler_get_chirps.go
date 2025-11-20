package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirpFeed, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not get chirps", err)
		return
	}

	var chirpList []Chirp

	for _, chirp := range chirpFeed {
		chirpList = append(chirpList, Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirpList)

}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")

	id, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not get chirp id", err)
		return
	}

	chirp, err := cfg.db.GetChirpByID(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}
	respondWithJSON(w, http.StatusOK, Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.ID,
	})
}
