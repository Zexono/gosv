package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

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
}
