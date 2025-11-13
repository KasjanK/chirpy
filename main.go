package main

import (
	"log"
	"net/http"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func main() {
	serveMux := http.NewServeMux()
	server := &http.Server{
		Handler: serveMux,
		Addr: ":8080",
	}

	serveMux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	serveMux.HandleFunc("/healthz", handlerReadiness)
	log.Fatal(server.ListenAndServe())
}
