package middleware

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

func MakeToken(userID string) string {
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": userID})
	return tokenString
}

func GetTokenClaims(r *http.Request) string {
	_, claims, _ := jwtauth.FromContext(r.Context())
	return claims["user_id"].(string)
}
