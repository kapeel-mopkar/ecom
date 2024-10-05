package auth

import "testing"

func TestCreateJWT(t *testing.T) {
	token, err := CreateJWT([]byte("secret"), 123)
	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}
	if token == "" {
		t.Error("expected token to be not empty")
	}
}
