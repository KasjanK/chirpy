package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/KasjanK/chirpy/internal/auth"
	"github.com/KasjanK/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Body      string    `json:"body"`
    UserID    uuid.UUID `json:"user_id"`
}

func (cfg apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string 	 `json:"body"`
	}


	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not find JWT", err)
		return
	}

	id, err := auth.ValidateJWT(bearerToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not validate JWT", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameters", err)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleanBody := replaceProfanities(params.Body)

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
			Body: cleanBody,
			UserID: id,
		})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create chirp", err)
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: cleanBody,
		UserID: chirp.UserID,
	},
)
}

func replaceProfanities(text string) string {
	profanities := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(text, " ")
	for i, word := range words {
		lowerWord := strings.ToLower(word)
		if slices.Contains(profanities, lowerWord) {
			words[i] = "****"	
		}	
	}
	return strings.Join(words, " ")
}
