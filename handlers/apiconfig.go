package handlers

import (
	"fmt"
	"net/http"
	"os"
)



type ApiConfig struct {
	fileserverHits int
	JWTSecret string
}

var Cfg = ApiConfig{
	JWTSecret: os.Getenv("JWT_SECRET"),
}


func (cfg *ApiConfig) MiddlewareHitsInc(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++

		h.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) ResetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	cfg.fileserverHits = 0
	w.Write([]byte("Hits Count Reset"))
}


func (cfg *ApiConfig) HandlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(fmt.Sprintf(`
	
	<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>
	
	`, cfg.fileserverHits)))
}
