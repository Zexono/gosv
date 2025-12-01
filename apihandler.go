package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

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
	w.Header().Set("Content-Type","text/html; charset=utf-8")
	w.WriteHeader(200)
	//head := "<h1>Welcome, Chirpy Admin</h1>"
	body := fmt.Sprintf(`
<html>

<body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
</body>

</html>
	`,cfg.fileserverHits.Load())
	//_,_ = w.Write([]byte(head))
	_,_ = w.Write([]byte(body))
	//fmt.Fprintf(w,"Hits: %d", cfg.fileserverHits.Load())
	
}

func (cfg *apiConfig) reset(w http.ResponseWriter,_ *http.Request){
	w.Header().Set("Content-Type","text/plain; charset=utf-8")
	w.WriteHeader(200)
	cfg.fileserverHits.Store(0)
	//_,_ = w.Write([]byte("Hits:"+string(num)))
}


func validate(w http.ResponseWriter, r *http.Request){
	type parameters struct {
    	Body string `json:"body"`
	}
	type returnVals struct {
		Error error `json:"error"`
        Valid bool `json:"valid"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
	var resBody returnVals
    if err != nil {
		resBody = returnVals{
			Error: err,
		}
		//w.WriteHeader(500)
		dat, err := json.Marshal(resBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
		}
    	w.Header().Set("Content-Type", "application/json")
    	w.WriteHeader(500)
    	w.Write(dat)

    }else {
		//validate time
		if len(params.Body) <= 140{
			resBody = returnVals{
			Valid: true,
			}
			w.WriteHeader(200)
		}else{
			resBody = returnVals{
				Error: fmt.Errorf("chirp is too long"),
			}
			w.WriteHeader(400)
		}

    	dat, err := json.Marshal(resBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
		}
    	w.Header().Set("Content-Type", "application/json")
    	w.Write(dat)

	}	
	
}