package main

import (
	"log"
	"net/http"
)



var apiCfg apiConfig
	//

func main() {
	mux := http.NewServeMux()
	h := http.FileServer(http.Dir("."))

	mux.Handle("/app/",(http.StripPrefix("/app", apiCfg.middlewareMetricsInc(h))))
	//mux.HandleFunc("/healthz",app)
	//mux.HandleFunc("/metrics",apiCfg.hit)
	//mux.HandleFunc("/reset",apiCfg.reset)
	mux.HandleFunc("GET /api/healthz", app)
	mux.HandleFunc("GET /admin/metrics", apiCfg.hit)
	mux.HandleFunc("POST /admin/reset", apiCfg.reset)
	mux.HandleFunc("POST /api/validate_chirp", validate)
	sv := http.Server{Handler: mux,Addr: ":8080"}
	log.Println("Serving files from . on port: 8080")
	sv.ListenAndServe()
}

