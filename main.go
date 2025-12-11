package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Zexono/gosv/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)



var apiCfg apiConfig
	//
const root = "."
const port = "8080"

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	jwt_secret := os.Getenv("SECRET")
	api_key := os.Getenv("POLKA_KEY")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		print(err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	apiCfg.db = dbQueries
	apiCfg.platform = platform
	apiCfg.secret = jwt_secret
	apiCfg.polkakey = api_key
	
	
	mux := http.NewServeMux()
	h := http.FileServer(http.Dir(root))

	mux.Handle("/app/",(http.StripPrefix("/app", apiCfg.middlewareMetricsInc(h))))
	//mux.HandleFunc("/healthz",app)
	//mux.HandleFunc("/metrics",apiCfg.hit)
	//mux.HandleFunc("/reset",apiCfg.reset)
	mux.HandleFunc("GET /api/healthz", app)
	mux.HandleFunc("GET /admin/metrics", apiCfg.hit)
	//mux.HandleFunc("POST /admin/reset", apiCfg.reset)
	//mux.HandleFunc("POST /api/validate_chirp", chirpsValidate)
	mux.HandleFunc("POST /admin/reset", userResetHandler)

	mux.HandleFunc("POST /api/users", userHandler)
	mux.HandleFunc("POST /api/u", userGetHandler)
	mux.HandleFunc("POST /api/login", userLoginHandler)

	mux.HandleFunc("POST /api/chirps", chirpsHandler)
	mux.HandleFunc("GET /api/chirps", chirpGetHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", chirpGetByIDHandler)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", chirpDeleteByIDHandler)

	mux.HandleFunc("POST /api/refresh", refreshEndpoint)
	mux.HandleFunc("POST /api/revoke", revokeEndpoint)

	mux.HandleFunc("PUT /api/users", updateOwnUsernamePassword)

	mux.HandleFunc("POST /api/polka/webhooks", userPolkaWebhooks)

	sv := http.Server{Handler: mux,Addr: ":"+port}
	log.Printf("Serving files from %s on port: %s",root,port)
	sv.ListenAndServe()
}

