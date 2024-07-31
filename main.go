package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	handlers "github.com/theyosefegy/chriby/handlers"
)



const filepathRoot = "./assets"
const port = "8080"



func main() {
	godotenv.Load()

	// by default, godotenv will look for a file named .env in the current directory

	mux := http.NewServeMux()


	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}


	// app/*
	mux.Handle("/app/*", http.StripPrefix("/app", handlers.Cfg.MiddlewareHitsInc(http.FileServer(http.Dir(filepathRoot)))))
	
	// api/*
	mux.HandleFunc("GET /api/healthz", handlers.ReadinessHandler)
	mux.HandleFunc("GET /api/reset", handlers.Cfg.ResetHandler)

	// Chirps's Endpoints
	mux.HandleFunc("POST /api/chirps", handlers.PostChirpHandler)
	mux.HandleFunc("GET /api/chirps", handlers.GetChirpHandler)
	mux.HandleFunc("GET /api/chirp/", handlers.GetChripByIDHandler)

	// User's Endpoints
	mux.HandleFunc("POST /api/users", handlers.PostUserHandler)
	mux.HandleFunc("GET /api/users/", handlers.GetUsersHandler)
	mux.HandleFunc("GET /api/user/", handlers.GetUserByIdHandler)
	mux.HandleFunc("PUT /api/users/", handlers.UpdateUserHandler)

	
	mux.HandleFunc("POST /api/login", handlers.PostLoginHandler)
	// admin/*
	mux.HandleFunc("GET /admin/metrics", handlers.Cfg.HandlerMetrics)

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}


