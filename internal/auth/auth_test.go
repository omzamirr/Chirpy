package auth

import (
    "testing"
    "time"

    "github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
    userID := uuid.New()
    secret := "mysecret"

    token, err := MakeJWT(userID, secret, time.Hour)
    if err != nil {
        t.Fatalf("unexpected error making JWT: %v", err)
    }

    gotID, err := ValidateJWT(token, secret)
    if err != nil {
        t.Fatalf("unexpected error validating JWT: %v", err)
    }
    if gotID != userID {
        t.Errorf("expected %v, got %v", userID, gotID)
    }

}


//expired token
func TestExpiredJWT(t *testing.T) {
    userID := uuid.New()
    secret := "mysecret"

    token, err := MakeJWT(userID, secret, -time.Hour)
    if err != nil {
        t.Fatalf("unexpected error making JWT: %v", err)
    }

    _, err = ValidateJWT(token, secret)
    if err == nil {
        t.Fatal("expected error for expired token, got nil")
    }
}


//wrong secret
func TestWrongSecretJWT(t *testing.T) {
    userID := uuid.New()

    token, err := MakeJWT(userID, "correctsecret", time.Hour)
    if err != nil {
        t.Fatalf("unexpected error making JWT: %v", err)
    }

    _, err = ValidateJWT(token, "wrongsecret")
    if err == nil {
        t.Fatal("expected error for wrong secret, got nil")
    }
}