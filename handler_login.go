package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/KasjanK/chirpy/internal/auth"
	"github.com/KasjanK/chirpy/internal/database"
)

func (cfg apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password 		 string `json:"password"`
		Email 	 		 string `json:"email"`
	}

	type response struct {
		User
		Token 	     string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	
	val, err := auth.CheckPassword(params.Password, user.HashedPassword)
	if err != nil || !val {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	expiration := time.Hour
	
	token, err := auth.MakeJWT(user.ID, cfg.secret, expiration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create JWT", err)
		return
	}

	refreshToken := auth.MakeRefreshToken()

	_, err = cfg.db.CreateToken(r.Context(), database.CreateTokenParams{
		Token: refreshToken,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not store refresh token in db", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID: user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email: user.Email,
		},
		Token: token,
		RefreshToken: refreshToken,
	})
}
