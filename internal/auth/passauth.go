package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error){
	hashpass , err := argon2id.CreateHash(password,argon2id.DefaultParams)
	if err != nil {
		return "",err
	}
	return hashpass,nil
}

func CheckPasswordHash(password, hash string) (bool, error){
	match,err := argon2id.ComparePasswordAndHash(password,hash)
	if err != nil {
		return false,err
	}
	return match,nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error){
	return jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject: uuid.UUID.String(userID)}).SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error){
	claimsStruct := jwt.RegisteredClaims{}
	token , err := jwt.ParseWithClaims(tokenString,&claimsStruct,func(token *jwt.Token) (any, error) {
	return []byte(tokenSecret), nil
	})
	if err != nil {
		log.Fatal(err)
	}

	userID, err := token.Claims.GetSubject()

	if err != nil{
		log.Fatal("unknown claims type, cannot proceed")
	}

	id ,err := uuid.Parse(userID)

	if err != nil{
		log.Fatal("wtf just happend")
	}

	return  id,nil
}

func GetBearerToken(headers http.Header) (string, error){
	TOKEN_STRING := headers.Get("Authorization")
	if TOKEN_STRING != "" {
		TOKEN_AUTH :=  strings.TrimPrefix(TOKEN_STRING, "Bearer ")
		return TOKEN_AUTH,nil
	}

	return "",fmt.Errorf("no token string")
}