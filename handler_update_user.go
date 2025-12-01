package main

import (
	"encoding/json"
	"net/http"

	"github.com/KasjanK/chirpy/internal/auth"
	"github.com/KasjanK/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email 	 string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not find token", err)
		return
	}

	id, err := auth.ValidateJWT(accessToken, cfg.secret)
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

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not hash password", err)
		return 
	}

	user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		Email: params.Email,
		HashedPassword: hashedPassword,
		ID: id,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not update user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID: user.ID,
			Email: user.Email,
			UpdatedAt: user.UpdatedAt,
			CreatedAt: user.CreatedAt,
		},
	})
}
