package main

import (
	"net/http"

	"github.com/KasjanK/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get JWT", err)
		return
	}
	 
	_, err = cfg.db.RevokeRefreshTokenByToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w,http.StatusInternalServerError, "Could not revoke token", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
