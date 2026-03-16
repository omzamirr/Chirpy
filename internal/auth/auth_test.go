package auth

import (
    "testing"
    "time"
    "net/http"
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


func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name      string
		headers   http.Header
		wantToken string
		wantErr   bool
	}{
		{
			name: "valid bearer token",
			headers: http.Header{
				"Authorization": []string{"Bearer abc123"},
			},
			wantToken: "abc123",
			wantErr:   false,
		},
		{
			name:      "missing authorization header",
			headers:   http.Header{},
			wantToken: "",
			wantErr:   true,
		},
		{
			name: "bad prefix",
			headers: http.Header{
				"Authorization": []string{"Token abc123"},
			},
			wantToken: "",
			wantErr:   true,
		},
		{
			name: "empty bearer token",
			headers: http.Header{
				"Authorization": []string{"Bearer "},
			},
			wantToken: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := GetBearerToken(tt.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotToken != tt.wantToken {
				t.Errorf("GetBearerToken() got = %v, want %v", gotToken, tt.wantToken)
			}
		})
	}
}
