package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)
const secretKey = "hee" 
func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	token ,err := MakeJWT(userID,secretKey,time.Hour)
	if err != nil {
		t.Errorf("FAIL %v",err)
		return 
	}
	if token != "" {
		return
	}
}

func TestValidateJWT(t *testing.T) {
	//test make
	userID := uuid.New()
	token ,_ := MakeJWT(userID,secretKey,time.Hour)
	validid , _:= ValidateJWT(token,secretKey)
	if validid != userID{
		t.Errorf("FAIL %v != %v",validid,userID)
		return
	}
	fmt.Println("pass make")

}

/*func TestValidateJWTValidString(t *testing.T){
	//invalid string
	userID := uuid.New()
	token ,_ := MakeJWT(userID,secretKey,time.Hour)
	_ , err := ValidateJWT("eiei",secretKey)
	if err == nil {
		t.Errorf("FAIL eiei != %v",token)
		return
	}
}

func TestValidateJWTExpToken(t *testing.T){
	//expire token
	userID := uuid.New()
	token ,_ := MakeJWT(userID,secretKey,5 *time.Second)
	time.Sleep(5 * time.Second)
	_ , err := ValidateJWT(token,secretKey)
	if err == nil {
		t.Errorf("expected error for expired token, got nil")
		return
	}
}

func TestValidateJWTWrongSecret(t *testing.T){
	//wrong secret
	userID := uuid.New()
	token ,_ := MakeJWT(userID,secretKey,time.Hour)
	_ , err := ValidateJWT(token,"wrong")
	if err == nil {
		t.Errorf("FAIL")
		return
	}
}*/

