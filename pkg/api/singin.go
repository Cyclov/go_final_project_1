package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type SigninReq struct {
	Password string `json:"password"`
}

type SigninResp struct {
	Token string `json:"token"`
}

func singinHandler(w http.ResponseWriter, r *http.Request) {

	var req SigninReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		JsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	pass := os.Getenv("TODO_PASSWORD")

	if len(pass) > 0 && req.Password != pass {
		JsonError(w, http.StatusUnauthorized, "invalid password")
		return
	}

	secret := []byte(pass)

	jwtToken := jwt.New(jwt.SigningMethodHS256)
	signedToken, err := jwtToken.SignedString(secret)
	if err != nil {
		fmt.Printf("failed to sign jwt: %s\n", err)
	}

	WriteJson(w, SigninResp{Token: signedToken})

}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var jwtstring string
			cookie, err := r.Cookie("token")
			if err == nil {
				jwtstring = cookie.Value
			}
			var valid bool
			token, err := jwt.Parse(jwtstring, func(t *jwt.Token) (interface{}, error) {
				secret := []byte(pass)
				return secret, nil
			})
			valid = err == nil && token.Valid

			if !valid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
