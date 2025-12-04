package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Zexono/gosv/internal/database"
	"github.com/google/uuid"
)


type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user_db,err := apiCfg.db.CreateUser(context.Background(),params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	user := User{
		ID: user_db.ID,
		CreatedAt: user_db.CreatedAt,
		UpdatedAt: user_db.UpdatedAt,
		Email: user_db.Email,
	}


	respondWithJSON(w, http.StatusCreated,user)
}

func userResetHandler(w http.ResponseWriter, _ *http.Request) {
	if apiCfg.platform != "dev" {
		respondWithJSON(w, http.StatusForbidden,"")
		return
	}
	apiCfg.db.DeleteAllUser(context.Background())
	log.Println("Delete all user")
	//apiCfg.db.DeleteAllChirp(context.Background())
}

func userGetHandler(w http.ResponseWriter, _ *http.Request) {
	if apiCfg.platform != "dev" {
		respondWithJSON(w, http.StatusForbidden,"")
		return
	}
	db ,err := apiCfg.db.GetAllUser(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}
	log.Println(db)
	//apiCfg.db.DeleteAllChirp(context.Background())
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body     string    `json:"body"`
	User_id   uuid.UUID `json:"user_id"`
}

func chirpsHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		User_id   uuid.UUID `json:"user_id"`
	}
	//type returnVals struct {
	//	Valid bool `json:"valid"`
		//Cleaned_body string `json:"cleaned_body"`
	//}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	chirp_db ,err:= apiCfg.db.CreateChirp(context.Background(),database.CreateChirpParams{Body: params.Body,UserID: params.User_id})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}
	chirp := Chirp{
		ID: chirp_db.ID,
		CreatedAt: chirp_db.CreatedAt,
		UpdatedAt: chirp_db.UpdatedAt,
		Body: chirp_db.Body,
		User_id: chirp_db.UserID,
	}
	/*respondWithJSON(w, http.StatusOK, returnVals{
		Valid: true,
		Cleaned_body: checkBadword(params.Body),
	})*/
	respondWithJSON(w, http.StatusCreated,chirp)
}
