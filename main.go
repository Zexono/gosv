package main

import (
	"net/http"
)




func main() {
	mux := http.NewServeMux()
	sv := http.Server{Handler: mux,Addr: ":8080"}
	sv.ListenAndServe()
}