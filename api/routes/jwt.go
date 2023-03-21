package routes

import (
	// "encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"

	// "log"
	// "net/http"
	"time"
)

type Claims struct {
	UserId  string `json:"user_id"`
	Rol     int    `json:"rol"`
	Empresa int    `json:"empresa_id"`
	jwt.RegisteredClaims
}

type ClaimsInvitation struct {
	Id        string `json:"superior_id"`
	Rol       int    `json:"rol"`
	EmpresaId int    `json:"empresa_id"`
	Email     string `json:"email"`
	jwt.RegisteredClaims
}

var sampleSecretKey = []byte(viper.GetString("JWT_SECRET"))

func GenerateInvitationJWT(id string, rol int, empresaId int, email string) (string, error) {
	claims := &ClaimsInvitation{
		Id:        id,
		Rol:       rol,
		EmpresaId: empresaId,
		Email:     email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ExtractClaimsInvitation(tokenString string) (*ClaimsInvitation, error) {
	claims := &ClaimsInvitation{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(tokenKey *jwt.Token) (interface{}, error) {
		return sampleSecretKey, nil
	})
	if err != nil {
		return claims, err
	}
	return claims, err
}

func GenerateJWT(userId string, rol int, empresaId int) (string, error) {
	claims := &Claims{
		UserId:  userId,
		Rol:     rol,
		Empresa: empresaId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(100 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ExtractClaims(tokenString string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(tokenKey *jwt.Token) (interface{}, error) {
		return sampleSecretKey, nil
	})
	if err != nil {
		return claims, err
	}
	return claims, err
}
