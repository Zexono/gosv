package main

import (
	"net/http"
)




func main() {
	mux := http.NewServeMux()
	h := http.FileServer(http.Dir("."))
	mux.Handle("/app/",http.StripPrefix("/app", h))
	mux.HandleFunc("/healthz",app)
	sv := http.Server{Handler: mux,Addr: ":8080"}
	sv.ListenAndServe()
}

func app(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","text/plain; charset=utf-8")
	w.WriteHeader(200)
	_,_ = w.Write([]byte("OK"))
}