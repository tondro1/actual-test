package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

type authedHandler func(http.ResponseWriter, *http.Request, *UserClaims)

func Authenticate(next authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var accessTokenString string
		for _, cookie := range r.Cookies() {
			if cookie.Name == "accessToken" {
				accessTokenString = cookie.Value
				break
			}
		}

		if accessTokenString == "" {
			next(w, r, &UserClaims{UserId: ""})
			return
		}

		accessToken, err := jwt.ParseWithClaims(accessTokenString, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
			pub, err := getPublicKey()
			if err != nil {
				return nil, err
			}
			return pub, nil
		})
		if err != nil {
			log.Println(err)
			next(w, r, &UserClaims{UserId: ""})
			return
		}
		if !accessToken.Valid {
			log.Println("token not valid")
			next(w, r, &UserClaims{UserId: ""})
			return
		}
		
		next(w, r, accessToken.Claims.(*UserClaims))
	}
}

func getSecretKey() (*rsa.PrivateKey, error) {
	SECRET_KEY_PATH := os.Getenv("PRIVATE_KEY_PATH")
	data, err := os.ReadFile(SECRET_KEY_PATH)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch priv := priv.(type) {
	case *rsa.PrivateKey:
		return priv, nil
	default:
		return nil, errors.New("private key type is not RSA")
	}

}

func getPublicKey() (*rsa.PublicKey, error) {
	PUBLIC_KEY_PATH := os.Getenv("PUBLIC_KEY_PATH")
	data, err := os.ReadFile(PUBLIC_KEY_PATH)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		return nil, errors.New("pub key type is not RSA")
	}
}