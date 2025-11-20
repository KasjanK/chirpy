package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/KasjanK/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	platform 		string	
	db 				*database.Queries
	fileserverHits 	atomic.Int32
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("could not open db: %s", err)
	}
	dbQueries := database.New(db)

	platform := os.Getenv("PLATFORM")
	serveMux := http.NewServeMux()
	server := &http.Server{
		Handler: serveMux,
		Addr: ":8080",
	}

	apiCfg := apiConfig{
		platform: platform,
		db: dbQueries,
		fileserverHits: atomic.Int32{},
	}

	serveMux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.handlerReqCount)
	serveMux.HandleFunc("GET /api/healthz", handlerReadiness)
	serveMux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	serveMux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	serveMux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)

	log.Fatal(server.ListenAndServe())
}
