package auth


import (
    "time"
	"strings"
	"errors"
	"net/http"

    "github.com/alexedwards/argon2id"
    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
)


func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err!= nil {
		return "", err
	}
	
	return hash, nil
}


func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	
	return match, nil
	
}


func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	secretKey := []byte(tokenSecret)
	
	token:= jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:   "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})

	return token.SignedString(secretKey)
}


func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims:= jwt.RegisteredClaims{}
	token, err:= jwt.ParseWithClaims(
		tokenString, &claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		},
	)
	if err != nil {
		return uuid.Nil, err
	}

	userIDString, err:= token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	id, err:= uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil

}


func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("authorization header format must be: Bearer <token>")
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
    if token == "" {
        return "", errors.New("bearer token is empty")
    }

    return token, nil
}
