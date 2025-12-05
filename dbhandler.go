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
	err := apiCfg.db.DeleteAllUser(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}
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

func chirpGetHandler(w http.ResponseWriter, r *http.Request) {
	chirp_db ,err := apiCfg.db.GetAllChirp(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get Chirp", err)
		return
	}
	//log.Println(db)

	
	var c []Chirp
	for _, v := range chirp_db {
		/*c[i] = Chirp{
			ID: v.ID,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			Body: v.Body,
			User_id: v.UserID,
		}*/
		c = append(c, Chirp{			
			ID: v.ID,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			Body: v.Body,
			User_id: v.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK,c)
	
}

func chirpGetByIDHandler(w http.ResponseWriter, r *http.Request) {

	uid,err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get Chirp", err)
		return
	}

	chirp_db ,err := apiCfg.db.GetChirpByID(context.Background(),uid)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get Chirp", err)
		return
	}
	//log.Println(db)
		c := Chirp{			
			ID: chirp_db.ID,
			CreatedAt: chirp_db.CreatedAt,
			UpdatedAt: chirp_db.UpdatedAt,
			Body: chirp_db.Body,
			User_id: chirp_db.UserID,
		}

	respondWithJSON(w, http.StatusOK,c)
	
}
