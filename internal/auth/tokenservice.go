package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/tondro1/actual-test/internal/database"
)

type UserClaims struct {
	UserId string `json:"id"`
	jwt.StandardClaims
}

func NewAccessToken(user database.User) (string, error) {
	now := time.Now().Local()

	uclaims := UserClaims{
		UserId: uuidToString(user.ID),
		StandardClaims: jwt.StandardClaims{
			Issuer: "localhost",
			Subject: user.Username,
			IssuedAt: now.Unix(),
			ExpiresAt: now.Add(time.Hour * 24).Unix(),
		},
	}

	signedAccessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, uclaims)
	secretKey, err := getSecretKey()
	if err != nil {
		return "", err
	}

	return signedAccessToken.SignedString(secretKey)
}

// func NewRefreshToken() (string, error) {
// 	now := time.Now().Local()
// 	rclaims := jwt.StandardClaims{
// 		Issuer: "localhost",
// 		IssuedAt: now.Unix(),
// 		ExpiresAt: now.Add(time.Hour * 48).Unix(),
// 	}

// 	signedRegreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, rclaims)
// 	secretKey, err := getSecretKey()
// 	if err != nil {
// 		return "", err
// 	}

// 	return signedRegreshToken.SignedString(secretKey)
// }

func uuidToString(uuid pgtype.UUID) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid.Bytes[0:4], uuid.Bytes[4:6], uuid.Bytes[6:8], uuid.Bytes[8:10], uuid.Bytes[10:16])
}