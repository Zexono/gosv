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

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		print(err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	apiCfg.db = dbQueries

	mux := http.NewServeMux()
	h := http.FileServer(http.Dir("."))

	mux.Handle("/app/",(http.StripPrefix("/app", apiCfg.middlewareMetricsInc(h))))
	//mux.HandleFunc("/healthz",app)
	//mux.HandleFunc("/metrics",apiCfg.hit)
	//mux.HandleFunc("/reset",apiCfg.reset)
	mux.HandleFunc("GET /api/healthz", app)
	mux.HandleFunc("GET /admin/metrics", apiCfg.hit)
	mux.HandleFunc("POST /admin/reset", apiCfg.reset)
	mux.HandleFunc("POST /api/validate_chirp", chirpsValidate)
	sv := http.Server{Handler: mux,Addr: ":8080"}
	log.Println("Serving files from . on port: 8080")
	sv.ListenAndServe()
}

