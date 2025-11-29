package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)


	var apiCfg apiConfig

func main() {
	mux := http.NewServeMux()
	h := http.FileServer(http.Dir("."))

	mux.Handle("/app/",(http.StripPrefix("/app", apiCfg.middlewareMetricsInc(h))))
	//mux.HandleFunc("/healthz",app)
	//mux.HandleFunc("/metrics",apiCfg.hit)
	//mux.HandleFunc("/reset",apiCfg.reset)
	mux.HandleFunc("GET /healthz", app)
	mux.HandleFunc("GET /metrics", apiCfg.hit)
	mux.HandleFunc("POST /reset", apiCfg.reset)
	sv := http.Server{Handler: mux,Addr: ":8080"}
	sv.ListenAndServe()
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w,r)
	})
}

func app(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","text/plain; charset=utf-8")
	w.WriteHeader(200)
	_,_ = w.Write([]byte("OK"))
}

func (cfg *apiConfig) hit(w http.ResponseWriter,_ *http.Request){
	w.Header().Set("Content-Type","text/plain; charset=utf-8")
	w.WriteHeader(200)
	fmt.Fprintf(w, "Hits: %d", cfg.fileserverHits.Load())
	//_,_ = w.Write([]byte("Hits:"+string(num)))
}

func (cfg *apiConfig) reset(w http.ResponseWriter,_ *http.Request){
	w.Header().Set("Content-Type","text/plain; charset=utf-8")
	w.WriteHeader(200)
	cfg.fileserverHits.Store(0)
	//_,_ = w.Write([]byte("Hits:"+string(num)))
}