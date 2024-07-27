package main

import (
	"log"
	"net/http"

	handlers "github.com/theyosefegy/chriby/handlers"
)

const filepathRoot = "."
const port = "8080"

func main() {
	mux := http.NewServeMux()
	cfg := handlers.ApiConfig{}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// app/*
	mux.Handle("/app/*", http.StripPrefix("/app", cfg.MiddlewareHitsInc(http.FileServer(http.Dir(filepathRoot)))))
	
	// api/*
	mux.HandleFunc("GET /api/healthz", handlers.ReadinessHandler)
	mux.HandleFunc("GET /api/reset", cfg.ResetHandler)

	mux.HandleFunc("POST /api/chirps", handlers.PostChirpHandler)
	mux.HandleFunc("GET /api/chirps", handlers.GetChirpHandler)
	
	// admin/*
	mux.HandleFunc("GET /admin/metrics", cfg.HandlerMetrics)

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

