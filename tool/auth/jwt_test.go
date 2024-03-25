package auth

import (
	"testing"
)

func TestClaims(t *testing.T) {
	u := &UserClaims{ClaimsId: 1, Name: "lin", Other: map[string]string{"auth_token": "xxxxxxxxxxx"}}
	token, err := Claims(u)
	if err != nil {
		t.Error(err)
	}
	t.Log(len(token), token)
}

func TestParse(t *testing.T) {
	signedString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTEwNzU4MDMsImlhdCI6MTcxMTA3NTUwMywiQ2xhaW1zSWQiOjEsIk5hbWUiOiJsaW4iLCJPdGhlciI6eyJhdXRoX3Rva2VuIjoieHh4eHh4eHh4eHgifX0.7Rf09LCV4gAMHslweNFcYGxHK5JAW_6ojXqI6iEL3tg"
	c, err := Parse(signedString)
	if err != nil {
		t.Error(err)
	}
	t.Log(c)
}
