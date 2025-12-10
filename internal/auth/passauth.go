package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// ErrNoAuthHeaderIncluded -
var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")

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
		return uuid.Nil, err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != "chirpy" {
		return uuid.Nil, errors.New("invalid issuer")
	}

	id, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}
	return id, nil
}

func GetBearerToken(headers http.Header) (string, error){
	TOKEN_STRING := headers.Get("Authorization")
	if TOKEN_STRING == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	if TOKEN_STRING != "" {
		TOKEN_AUTH :=  strings.TrimPrefix(TOKEN_STRING, "Bearer ")
		return TOKEN_AUTH,nil
	}

	return "",fmt.Errorf("no token string")
}

func MakeRefreshToken() (string, error){
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b),nil
}