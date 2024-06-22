package api

import (
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/tondro1/actual-test/internal/auth"
)

type authedHandler func(http.ResponseWriter, *http.Request, *auth.UserClaims)
type validatedHandler func(http.ResponseWriter, *http.Request, bool)

func (api *ApiCfg) Authenticate(next authedHandler) validatedHandler {
	return func(w http.ResponseWriter, r *http.Request, valid bool) {
		if !valid {
			log.Println("token not valid")
			next(w, r, &auth.UserClaims{UserId: ""})
			return
		}
		var accessTokenString string
		for _, cookie := range r.Cookies() {
			if cookie.Name == "accessToken" {
				accessTokenString = cookie.Value
				break
			}
		}

		if accessTokenString == "" {
			next(w, r, &auth.UserClaims{UserId: ""})
			return
		}

		accessToken, err := jwt.ParseWithClaims(accessTokenString, &auth.UserClaims{}, func(t *jwt.Token) (interface{}, error) {
			pub, err := auth.GetPublicKey()
			if err != nil {
				return nil, err
			}
			return pub, nil
		})

		if !accessToken.Valid {
			log.Println("token not valid")
			next(w, r, &auth.UserClaims{UserId: ""})
			return
		}
		if err != nil {
			log.Println(err)
			next(w, r, &auth.UserClaims{UserId: ""})
			return
		}
		
		next(w, r, accessToken.Claims.(*auth.UserClaims))
	}
}

func (api *ApiCfg) Validate(next validatedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var accessTokenString string
		for _, cookie := range r.Cookies() {
			if cookie.Name == "accessToken" {
				accessTokenString = cookie.Value
				break
			}
		}

		if accessTokenString == "" {
			next(w, r, false)
			return
		}

		_, err := api.DB.GetLoggedOutToken(r.Context(), accessTokenString)

		if err != nil && strings.Contains(err.Error(), "no rows in result set") {
			next(w, r, true)
			return
		} else if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			log.Println("token found")
			next(w, r, false)
			return
		}
	}
}