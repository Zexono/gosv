package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/Zexono/gosv/internal/auth"
	"github.com/Zexono/gosv/internal/database"
	"github.com/google/uuid"
)


type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	ChirpyRed bool		`json:"is_chirpy_red"`
	//Token	  string	`json:"token"`
	//RefreshToken string `json:"refresh_token"`
	//Password  string	`json:"password"`
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email 	  string `json:"email"`
		Password  string	`json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	pass ,err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "password hash error", err)
		return
	}

	user_db,err := apiCfg.db.CreateUser(context.Background(),database.CreateUserParams{
		Email: params.Email,
		HashedPassword: pass,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	

	user := User{
		ID: user_db.ID,
		CreatedAt: user_db.CreatedAt,
		UpdatedAt: user_db.UpdatedAt,
		Email: user_db.Email,
		ChirpyRed: user_db.IsChirpyRed,
		//Password: user_db.HashedPassword,
	}


	respondWithJSON(w, http.StatusCreated,user)
}

func userLoginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email 	  string `json:"email"`
		Password  string `json:"password"`
		//ExpiresIn int `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user_db , err := apiCfg.db.GetUserByEmail(context.Background(),params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	pass_match ,err := auth.CheckPasswordHash(params.Password,user_db.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}
	if pass_match {
		//var duration time.Duration
		/*if params.ExpiresIn <= 0 || params.ExpiresIn > 3600  {
			duration = time.Duration(3600 * time.Second)
		}else {
			duration = time.Duration(params.ExpiresIn) * time.Second
		}*/

		duration := time.Hour
		ac_token , err := auth.MakeJWT(user_db.ID,apiCfg.secret,duration)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "make JWT auth err", err)
			return
		}

		rf_token ,err:= auth.MakeRefreshToken()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "make refresh token auth err", err)
			return
		}
		rf_db,err := apiCfg.db.CreateRefresh_tokens(context.Background(),database.CreateRefresh_tokensParams{
			Token: rf_token,
			UserID: user_db.ID,
		})
		if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong with refresh token", err)
		return
	}
		

		user := User{
		ID: user_db.ID,
		CreatedAt: user_db.CreatedAt,
		UpdatedAt: user_db.UpdatedAt,
		Email: user_db.Email,
		ChirpyRed: user_db.IsChirpyRed,
		//Password: user_db.HashedPassword,
		}

		token_response := response{
			User: user,
			Token: ac_token,
			RefreshToken: rf_db.Token,
		}
		respondWithJSON(w, http.StatusOK,token_response)
	}else {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

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

func userPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	type data struct{
		ID uuid.UUID `json:"user_id"`
	}
	type parameters struct {
		Event 	  string `json:"event"`
		data `json:"data"`
	}

	apikey,err := auth.GetAPIKey(r.Header)
	if apikey !=  apiCfg.polkakey{
		respondWithError(w, http.StatusUnauthorized, "Unauthorize api key", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		respondWithError(w,http.StatusNoContent,"event is not upgraded",nil)
		return
	}

	_,err = apiCfg.db.UpdateUserChirpyred(context.Background(),params.data.ID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find user", err)
		return
	}

	respondWithJSON(w,http.StatusNoContent,nil)

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
		//User_id   uuid.UUID `json:"user_id"`
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

	bearer_token ,err:= auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauth bearer token", err)
		return
	}

	valid_uid,err := auth.ValidateJWT(bearer_token,apiCfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "JWT invalid", err)
		return
	}	

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	chirp_db ,err:= apiCfg.db.CreateChirp(context.Background(),database.CreateChirpParams{
		Body: params.Body,
		UserID: valid_uid})
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

	author := r.URL.Query().Get("author_id")
	sort_by := r.URL.Query().Get("sort")

	var uid uuid.UUID
	if author != "" {
		uid,err = uuid.Parse(author)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't Parse uuid", err)
			return
		}
	}
	

	var c []Chirp
	for _, v := range chirp_db {
		if uid != uuid.Nil && v.UserID != uid {
			continue
		}
		c = append(c, Chirp{			
			ID: v.ID,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			Body: v.Body,
			User_id: v.UserID,
		})
	}

	if sort_by == "desc" {
		sort.Slice(c,func(i, j int) bool { return c[i].CreatedAt.After(c[j].CreatedAt) })
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

func chirpDeleteByIDHandler(w http.ResponseWriter, r *http.Request) {
	ac_token,err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find token", err)
		return
	}

	user_id,err := auth.ValidateJWT(ac_token,apiCfg.secret)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Invalid access token", err)
		return
	}


	chirp_id,err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get Chirp ID", err)
		return
	}

	chirp_db,err := apiCfg.db.GetChirpByID(context.Background(),chirp_id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't found Chirp", err)
		return
	}
	if chirp_db.UserID != user_id {
		respondWithError(w,http.StatusForbidden, "can't delete because it not your chirp",nil)
		return
	}


	apiCfg.db.DeleteChirpByID(context.Background(),database.DeleteChirpByIDParams{
		UserID: user_id,
		ID: chirp_id,
	})


	respondWithJSON(w, http.StatusNoContent,nil)
	
}

func refreshEndpoint(w http.ResponseWriter, r *http.Request) {
	type returnVals struct {
		Token string `json:"token"`
	}
	
	
	rf_token,err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find token", err)
		return
	}

	user_db, err := apiCfg.db.GetUserFromRefreshToken(context.Background(),rf_token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get user Refresh token", err)
		return
	}

	//old filter
	//if rf_token_db.ExpiresAt.Before(time.Now()) || !rf_token_db.RevokedAt.Valid {
	//	respondWithError(w, http.StatusUnauthorized, "token expire", err)
	//	return
	//}

	rf_token,err = auth.MakeJWT(user_db.ID,apiCfg.secret,time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "something wrong with make JWT", err)
		return
	}
	value := returnVals{
		Token: rf_token,
	}

	respondWithJSON(w, http.StatusOK,value)
	
}

func revokeEndpoint(w http.ResponseWriter, r *http.Request) {
	
	rf_token,err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find token", err)
		return
	}

	err = apiCfg.db.UpdateRefreshTokenRevoke(context.Background(),rf_token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't Revoke refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent,nil)
	
}

func updateOwnUsernamePassword(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email 	  string `json:"email"`
		Password  string `json:"password"`
	}

	ac_token,err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find token", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hash_pass,err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Hash password error", err)
		return
	}

	valid_userid,err := auth.ValidateJWT(ac_token,apiCfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid access token", err)
		return
	}

	update_user_db, err := apiCfg.db.UpdateUsernamePassword(context.Background(),database.UpdateUsernamePasswordParams{
		ID: valid_userid,
		Email: params.Email,
		HashedPassword: hash_pass,
	})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid user ID", err)
		return
	}

	update_user := User{
		ID: update_user_db.ID,
		Email: update_user_db.Email,
		CreatedAt: update_user_db.CreatedAt,
		UpdatedAt: update_user_db.UpdatedAt,
	}

	respondWithJSON(w,http.StatusOK,update_user)

}
